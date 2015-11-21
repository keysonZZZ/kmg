package kmgThirdCloud

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/bronze1man/kmg/encoding/kmgJson"
	"github.com/bronze1man/kmg/kmgRand"
	"github.com/bronze1man/kmg/kmgXss"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var logPrefixAliyun = "[kmgThirdCloud SDKAliyun]"

type AliyunSDK struct {
	ConfigParam      *url.Values
	AccessKeyId      string
	AccessKeySecret  string
	SecurityGroupId  string
	InstancePaidType AliyunPaidType
	InstancePassword string
	Region           string //Alias RegionId
	InstanceName     string
	ImageName        string //Alias ImageId
	FlavorName       string //Alias InstanceType
}

type AliyunRespond struct {
	RequestId       string
	RegionId        string
	PageNumber      int
	PageSize        int
	TotalCount      int
	Images          map[string][]AliyunImage
	Regions         map[string][]AliyunRegion
	Instances       map[string][]AliyunInstance
	SecurityGroupId string
	InstanceId      string
	IpAddress       string
	//错误消息
	Code    string
	Message string
}

type AliyunPaidType string

const (
	//包年包月
	AliyunPaidTypePre AliyunPaidType = "PrePaid"
	//按量
	AliyunPaidTypePost AliyunPaidType = "PostPaid"
)

type AliyunInstanceStatus string

const (
	AliyunInstanceStatusRunning AliyunInstanceStatus = "Running"
	AliyunInstanceStatusStopped AliyunInstanceStatus = "Stopped"
	AliyunInstanceStatusDeleted AliyunInstanceStatus = "Deleted"
)

type AliyunImage struct {
	ImageId   string
	ImageName string
	OSName    string
	OSType    string
}

type AliyunRegion struct {
	RegionId  string
	LocalName string
}

type AliyunInstance struct {
	InstanceId         string
	Status             AliyunInstanceStatus
	InstanceName       string
	ImageId            string
	RegionId           string
	InstanceType       string
	PublicIpAddress    map[string][]string
	InstanceChargeType AliyunPaidType
	ExpiredTime        string
}

func (aliyunInstance AliyunInstance) getIp() string {
	l := aliyunInstance.PublicIpAddress
	if l == nil || len(l) == 0 {
		return ""
	}
	if len(l["IpAddress"]) == 0 {
		return ""
	}
	return l["IpAddress"][0]
}

func NewAliyunSDK(accessKeyId, accessKeySecret, securityGroupId string) *AliyunSDK {
	sdk := &AliyunSDK{
		ConfigParam:      &url.Values{},
		AccessKeyId:      accessKeyId,
		AccessKeySecret:  accessKeySecret,
		SecurityGroupId:  securityGroupId,
		Region:           "cn-shenzhen",
		ImageName:        "ubuntu1404_64_20G_aliaegis_20150325.vhd",
		FlavorName:       "ecs.s3.medium", //测试用 ecs.t1.small
		InstanceName:     "kmg-AliyunSDK-auto-build",
		InstancePaidType: AliyunPaidTypePre,
	}
	sdk.ConfigParam.Set("Format", "JSON")
	sdk.ConfigParam.Set("Version", "2014-05-26")
	sdk.ConfigParam.Set("AccessKeyId", accessKeyId)
	sdk.ConfigParam.Set("SignatureMethod", "HMAC-SHA1")
	sdk.ConfigParam.Set("SignatureVersion", "1.0")
	return sdk
}

func (sdk *AliyunSDK) signature(param *url.Values) string {
	param = combineParam(sdk.ConfigParam, param)
	param.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	param.Set("SignatureNonce", kmgRand.MustCryptoRandToNum(14))
	strSlice := []string{"GET", kmgXss.Urlv("/"), kmgXss.Urlv(param.Encode())}
	stringToSign := strings.Join(strSlice, "&")
	hash := hmac.New(sha1.New, []byte(sdk.AccessKeySecret+"&"))
	hash.Write([]byte(stringToSign))
	b := hash.Sum(nil)
	signature := base64.StdEncoding.EncodeToString(b)
	param.Set("Signature", signature)
	return "https://ecs.aliyuncs.com/?" + param.Encode()
}

func (sdk *AliyunSDK) MustCall(param *url.Values) *AliyunRespond {
	resp, isErr := sdk.Call(param)
	if isErr {
		panic(resp.Code + " " + resp.Message)
	}
	return resp
}

func (sdk *AliyunSDK) Call(param *url.Values) (r *AliyunRespond, isErr bool) {
	u := sdk.signature(param)
	resp, err := http.Get(u)
	handleErr(err)
	body, err := ioutil.ReadAll(resp.Body)
	handleErr(err)
	aliyunResp := &AliyunRespond{}
	kmgJson.MustUnmarshal(body, aliyunResp)
	if aliyunResp.Code != "" || aliyunResp.Message != "" {
		isErr = true
	}
	r = aliyunResp
	return
}

func (sdk *AliyunSDK) runSingleAction(instanceId, actionName string, expectStatus AliyunInstanceStatus) {
	param := &url.Values{}
	param.Set("Action", actionName)
	param.Set("InstanceId", instanceId)
	sdk.MustCall(param)
	sdk.WaitUntil(instanceId, expectStatus)
}

//创建机器
//分配公网IP
//重启机器
func (sdk *AliyunSDK) CreateInstance() (ip string) {
	ip = sdk.MakeInstanceAvailable(sdk.AllocateNewInstance())
	if ip == "" {
		panic("[AliyunSDK CreateInstance] Failed")
	}
	return ip
}

func (sdk *AliyunSDK) AllocateNewInstance() (id string) {
	prefix := "[AliyunSDK AllocateNewInstance]"
	f := func() string {
		defer func() {
			r := recover()
			if r != nil {
				fmt.Println(prefix, "failed")
			}
		}()
		param := &url.Values{}
		param.Set("Action", "CreateInstance")
		param.Set("SecurityGroupId", sdk.SecurityGroupId)
		param.Set("RegionId", sdk.Region)
		param.Set("ImageId", sdk.ImageName)
		param.Set("InstanceType", sdk.FlavorName)
		param.Set("InstanceName", sdk.InstanceName)
		param.Set("InstanceChargeType", string(sdk.InstancePaidType))
		if sdk.InstancePaidType == AliyunPaidTypePre {
			param.Set("Period", "1")
		}
		param.Set("InternetChargeType", "PayByTraffic")
		param.Set("InternetMaxBandwidthIn", "200")
		param.Set("InternetMaxBandwidthOut", "100")
		if sdk.InstancePassword == "" {
			panic("Empty password of instance don't allow!")
		}
		param.Set("Password", sdk.InstancePassword)
		resp := sdk.MustCall(param)
		return resp.InstanceId
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
func (sdk *AliyunSDK) MakeInstanceAvailable(id string) (ip string) {
	prefix := "[AliyunSDK MakeInstanceAvailable]"
	f := func() string {
		defer func() {
			r := recover()
			if r != nil {
				fmt.Println(prefix, "failed")
			}
		}()
		ip := sdk.AllocatePublicIpAddress(id)
		sdk.rebootInstance(id)
		return ip
	}
	for i := 0; i < 12; i++ {
		ip = f()
		if ip != "" {
			break
		} else {
			fmt.Println(prefix, "retry", i)
			time.Sleep(time.Second * 10)
		}
	}
	return ip
}

func (sdk *AliyunSDK) IsInstanceRunning(id string) bool {
	all := sdk.getAllInstance()
	for _, i := range all {
		if i.InstanceId == id {
			return i.Status == AliyunInstanceStatusRunning
		}
	}
	return false
}

func (sdk *AliyunSDK) WaitUntil(instanceId string, status AliyunInstanceStatus) {
	for {
		var instance AliyunInstance
		for _, ins := range sdk.getAllInstance() {
			if ins.InstanceId == instanceId {
				instance = ins
				break
			}
		}
		fmt.Println(logPrefixAliyun, instance.InstanceId, instance.getIp(), instance.Status)
		if instance.Status == status {
			break
		}
		time.Sleep(time.Second)
	}
}

func (sdk *AliyunSDK) StartInstance(instanceId string) {
	if !sdk.IsInstanceRunning(instanceId) {
		sdk.runSingleAction(instanceId, "StartInstance", AliyunInstanceStatusRunning)
	}
}

func (sdk *AliyunSDK) StopInstance(instanceId string) {
	if sdk.IsInstanceRunning(instanceId) {
		sdk.runSingleAction(instanceId, "StopInstance", AliyunInstanceStatusStopped)
	}
}

func (sdk *AliyunSDK) rebootInstance(instanceId string) {
	//不用阿里云API自带的重启接口，那个是异步的
	sdk.StopInstance(instanceId)
	sdk.StartInstance(instanceId)
}

func (sdk *AliyunSDK) delete(instance AliyunInstance) {
	//包年包月不能直接释放，只能先关机了
	if instance.InstanceChargeType == AliyunPaidTypePre && instance.Status == AliyunInstanceStatusRunning {
		sdk.StopInstance(instance.InstanceId)
		return
	}
	sdk.StopInstance(instance.InstanceId)
	param := &url.Values{}
	param.Set("Action", "DeleteInstance")
	param.Set("InstanceId", instance.InstanceId)
	sdk.MustCall(param)
}

func (sdk *AliyunSDK) DeleteInstanceById(id string) {
	for _, instance := range sdk.getAllInstance() {
		if instance.InstanceId == id {
			sdk.delete(instance)
		}
	}
}

func (sdk *AliyunSDK) DeleteInstance(ip string) {
	for _, instance := range sdk.getAllInstance() {
		if instance.getIp() == ip {
			sdk.delete(instance)
		}
	}
}

func (sdk *AliyunSDK) RenameInstanceByIp(name, ip string) {
	instance, exist := sdk.ListAllRunningInstance()[ip]
	if !exist {
		return
	}
	param := &url.Values{}
	param.Set("Action", "ModifyInstanceAttribute")
	param.Set("InstanceId", instance.Id)
	param.Set("InstanceName", name)
	sdk.MustCall(param)
}

func (sdk *AliyunSDK) ListAllRunningInstance() (ipInstanceMap map[string]Instance) {
	all := sdk.getAllInstance()
	ipInstanceMap = map[string]Instance{}
	for _, i := range all {
		if i.Status != AliyunInstanceStatusRunning {
			fmt.Println(logPrefixAliyun, i.Status, i.InstanceName, i.getIp())
			continue
		}
		ip := i.getIp()
		if ip == "" {
			continue
		}
		ipInstanceMap[ip] = Instance{
			Id:          i.InstanceId,
			Ip:          ip,
			Name:        i.InstanceName,
			BelongToSDK: sdk,
		}
	}
	return
}

func (sdk *AliyunSDK) ListAllInstance() (idInstanceMap map[string]Instance) {
	all := sdk.getAllInstance()
	idInstanceMap = map[string]Instance{}
	for _, i := range all {
		idInstanceMap[i.InstanceId] = Instance{
			Id:          i.InstanceId,
			Ip:          i.getIp(),
			Name:        i.InstanceName,
			BelongToSDK: sdk,
		}
	}
	return
}

func (sdk *AliyunSDK) getAllInstance() []AliyunInstance {
	param := &url.Values{}
	pre := 50
	p := 1
	all := []AliyunInstance{}
	param.Set("Action", "DescribeInstances")
	param.Set("RegionId", sdk.Region)
	param.Set("PageSize", strconv.Itoa(pre))
	param.Set("PageNumber", strconv.Itoa(p))
	resp := sdk.MustCall(param)
	all = append(all, resp.Instances["Instance"]...)
	total := resp.TotalCount
	totalPage := total / pre
	if total%pre != 0 {
		totalPage++
	}
	p = 2
	for {
		if p > totalPage {
			break
		}
		param.Set("PageNumber", strconv.Itoa(p))
		resp := sdk.MustCall(param)
		all = append(all, resp.Instances["Instance"]...)
		p++
	}
	return all
}

func (sdk *AliyunSDK) AllocatePublicIpAddress(instanceId string) (ip string) {
	param := &url.Values{}
	param.Set("Action", "AllocatePublicIpAddress")
	param.Set("InstanceId", instanceId)
	resp := &AliyunRespond{}
	isErr := false
	for {
		resp, isErr = sdk.Call(param)
		time.Sleep(time.Second)
		//没有错误，表示分配成功
		if !isErr {
			break
		}
	}
	return resp.IpAddress
}

func (sdk *AliyunSDK) CreateSecurityGroup() (securityGroupId string) {
	param := &url.Values{}
	param.Set("Action", "CreateSecurityGroup")
	param.Set("RegionId", sdk.Region)
	resp := sdk.MustCall(param)
	securityGroupId = resp.SecurityGroupId
	return
}

func (sdk *AliyunSDK) GetAllRegion() []AliyunRegion {
	param := &url.Values{}
	param.Set("Action", "DescribeRegions")
	resp := sdk.MustCall(param)
	return resp.Regions["Region"]
}

func (sdk *AliyunSDK) GetAllImage() []AliyunImage {
	param := &url.Values{}
	param.Set("Action", "DescribeImages")
	param.Set("RegionId", "cn-shenzhen")
	param.Set("PageSize", strconv.Itoa(100))
	resp := sdk.MustCall(param)
	return resp.Images["Image"]
}

func combineParam(p ...*url.Values) *url.Values {
	list := []string{}
	for _, v := range p {
		list = append(list, v.Encode())
	}
	out, err := url.ParseQuery(strings.Join(list, "&"))
	handleErr(err)
	return &out
}
