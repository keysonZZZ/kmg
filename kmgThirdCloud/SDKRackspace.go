package kmgThirdCloud

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgSsh"
	"github.com/rackspace/gophercloud"
	OpenStackImages "github.com/rackspace/gophercloud/openstack/compute/v2/images"
	OpenStackServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/flavors"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
	"time"
)

var logPrefixRackspace = "[kmgThirdCloud SDKRackspace]"

type RackspaceInstanceStatus string

const (
	RackspaceInstanceStatusACTIVE  RackspaceInstanceStatus = "ACTIVE"
	RackspaceInstanceStatusBUILD   RackspaceInstanceStatus = "BUILD"
	RackspaceInstanceStatusDELETED RackspaceInstanceStatus = "DELETED"
	RackspaceInstanceStatusERROR   RackspaceInstanceStatus = "ERROR"
	RackspaceInstanceStatusUNKNOWN RackspaceInstanceStatus = "UNKNOWN"
)

type RackspaceSDK struct {
	Username     string
	APIKey       string
	p            *gophercloud.ServiceClient
	SSHKeyName   string
	Region       string
	InstanceName string
	ImageName    string
	FlavorName   string
}

func NewRackspaceSDK(username, apiKey, SSHKeyName string) *RackspaceSDK {
	if SSHKeyName == "" {
		panic("Empty SSHKeyName of instance don't allow!")
	}
	sdk := &RackspaceSDK{
		Username:     username,
		APIKey:       apiKey,
		SSHKeyName:   SSHKeyName,
		Region:       "HKG",
		InstanceName: "kmg-RackspaceSDK-auto-build",
		ImageName:    "Ubuntu 14.04 LTS (Trusty Tahr) (PVHVM)",
		FlavorName:   "4 GB General Purpose v1",
	}
	ao := gophercloud.AuthOptions{
		Username: sdk.Username,
		APIKey:   sdk.APIKey,
	}
	provider, err := rackspace.AuthenticatedClient(ao)
	handleErr(err)
	serviceClient, err := rackspace.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: sdk.Region,
	})
	handleErr(err)
	sdk.p = serviceClient
	return sdk
}

func (sdk *RackspaceSDK) CreateInstance() (ip string) {
	for i := 0; i < 12; i++ {
		id := sdk.AllocateNewInstance()
		ip = sdk.MakeInstanceAvailable(id)
		if ip == "" {
			sdk.DeleteInstanceById(id)
			continue
		} else {
			return ip
		}
	}
	if ip == "" {
		panic("[RackspaceSDK CreateInstance] Failed")
	}
	return ip
}

func (sdk *RackspaceSDK) AllocateNewInstance() (id string) {
	prefix := "[RackspaceSDK AllocateNewInstance]"
	f := func() string {
		opt := servers.CreateOpts{
			Name:       sdk.InstanceName,
			ImageName:  sdk.ImageName,
			FlavorName: sdk.FlavorName,
			KeyPair:    sdk.SSHKeyName,
		}
		result := servers.Create(sdk.p, opt)
		one, err := result.Extract()
		if err != nil {
			fmt.Println(prefix, err)
			return ""
		}
		return one.ID
	}
	for i := 0; i < 12; i++ {
		id = f()
		if id != "" {
			break
		} else {
			fmt.Println(prefix, "retry", i)
			time.Sleep(time.Second * 10)
		}
	}
	return id
}

//不断重试,彻底失败,会返回空字符串
func (sdk *RackspaceSDK) MakeInstanceAvailable(id string) (ip string) {
	prefix := "[RackspaceSDK AllocateNewInstance]"
	interval := time.Second * 3
	for i := 0; i < 24; i++ {
		ins := sdk.getInstance(id)
		fmt.Println(prefix, ins.AccessIPv4, ins.Status)
		if ins.Status != string(RackspaceInstanceStatusACTIVE) {
			time.Sleep(interval)
			continue
		}
		i = 0
		isReachable, _ := kmgSsh.AvailableCheck(&kmgSsh.RemoteServer{
			Address: ins.AccessIPv4,
		})
		if isReachable {
			return ins.AccessIPv4
		} else {
			time.Sleep(interval)
		}
	}
	return ""
}

func (sdk *RackspaceSDK) getInstance(id string) (s *OpenStackServers.Server) {
	for i := 0; i < 12; i++ {
		var err error
		s, err = servers.Get(sdk.p, id).Extract()
		if err != nil {
			fmt.Println("RackspaceSDK getInstance", err)
			time.Sleep(time.Second)
			continue
		}
	}
	return s
}

func (sdk *RackspaceSDK) RenameInstanceByIp(name, ip string) {
	instance, exist := sdk.ListAllRunningInstance()[ip]
	if exist {
		servers.Update(sdk.p, instance.Id, OpenStackServers.UpdateOpts{
			Name: name,
		})
	}
}

func (sdk *RackspaceSDK) DeleteInstance(ip string) {
	instance, exist := sdk.ListAllRunningInstance()[ip]
	if exist {
		servers.Delete(sdk.p, instance.Id)
	}
	for {
		_, exist := sdk.ListAllRunningInstance()[ip]
		if !exist {
			break
		}
		time.Sleep(time.Second)
	}
}

func (sdk *RackspaceSDK) DeleteInstanceById(id string) {
	instance, exist := sdk.ListAllInstance()[id]
	if exist {
		servers.Delete(sdk.p, instance.Id)
	}
	for {
		_, exist := sdk.ListAllInstance()[id]
		if !exist {
			break
		}
		time.Sleep(time.Second)
	}
}

func (sdk *RackspaceSDK) ListAllRunningInstance() (ipInstanceMap map[string]Instance) {
	ipInstanceMap = map[string]Instance{}
	err := servers.List(sdk.p, nil).EachPage(func(p pagination.Page) (quit bool, err error) {
		serverSlice, err := servers.ExtractServers(p)
		for _, server := range serverSlice {
			if server.Status != string(RackspaceInstanceStatusACTIVE) {
				fmt.Println(logPrefixRackspace, server.ID, server.AccessIPv4, server.Status)
				continue
			}
			ipInstanceMap[server.AccessIPv4] = Instance{
				Ip:          server.AccessIPv4,
				Id:          server.ID,
				Name:        server.Name,
				BelongToSDK: sdk,
			}
		}
		return true, err
	})
	handleErr(err)
	return
}

func (sdk *RackspaceSDK) ListAllInstance() (idInstanceMap map[string]Instance) {
	idInstanceMap = map[string]Instance{}
	err := servers.List(sdk.p, nil).EachPage(func(p pagination.Page) (quit bool, err error) {
		serverSlice, err := servers.ExtractServers(p)
		for _, server := range serverSlice {
			idInstanceMap[server.ID] = Instance{
				Ip:          server.AccessIPv4,
				Id:          server.ID,
				Name:        server.Name,
				BelongToSDK: sdk,
			}
		}
		return true, err
	})
	handleErr(err)
	return
}

func (sdk *RackspaceSDK) GetAllFlavor() []flavors.Flavor {
	out := []flavors.Flavor{}
	err := flavors.ListDetail(sdk.p, nil).EachPage(func(p pagination.Page) (quit bool, err error) {
		v, err := flavors.ExtractFlavors(p)
		out = append(out, v...)
		return true, err
	})
	handleErr(err)
	return out
}

func (sdk *RackspaceSDK) GetAllImage() []OpenStackImages.Image {
	out := []OpenStackImages.Image{}
	err := OpenStackImages.ListDetail(sdk.p, nil).EachPage(func(p pagination.Page) (quit bool, err error) {
		v, err := OpenStackImages.ExtractImages(p)
		out = append(out, v...)
		return true, err
	})
	handleErr(err)
	return out
}

func (sdk *RackspaceSDK) PrintAllImage() {
	all := sdk.GetAllImage()
	for _, i := range all {
		fmt.Println("ID", i.ID, "Name", i.Name)
	}
}

func (sdk *RackspaceSDK) PrintAllFlavor() {
	all := sdk.GetAllFlavor()
	for _, i := range all {
		fmt.Println("ID", i.ID, "Name", i.Name)
	}
}
