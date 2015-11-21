package kmgThirdCloud

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Instance struct {
	Ip          string
	Id          string //第三方只认他们自己产生的ID
	Name        string
	BelongToSDK SDK
}

//检查是否满足 interface
var _ SDK = (*RackspaceSDK)(nil)
var _ SDK = (*AliyunSDK)(nil)

//本接口中，为降低调用者复杂度，不提供机器「状态」，因为状态较为复杂，而调用者其实只关心正在可用的，即正在运行的机器
type SDK interface {
	//Field
	//Region       string //实例所在地区
	//InstanceName string
	//ImageName    string //操作系统
	//FlavorName   string //实例配置，如 4 CPU / 4GB RAM

	//以实例外网 Ip 为主键的方法，通常用这些已经够了
	//创建一个新实例，保证实例可用,出现错误 panic
	CreateInstance() (ip string)
	DeleteInstance(ip string)
	RenameInstanceByIp(name, ip string)
	ListAllRunningInstance() (ipInstanceMap map[string]Instance)

	//向云服务提供商申请分配新实例，http 请求后，立即返回实例ID，不要求实例当时可用
	AllocateNewInstance() (id string)
	//让一台实例变为可用状态
	//不断重试,彻底失败,会返回空字符串
	MakeInstanceAvailable(id string) (ip string)

	//以实例 Id 为主键的方法，一般不使用，当没有实例没有 ip 时可以使用
	ListAllInstance() (idInstanceMap map[string]Instance)
	DeleteInstanceById(id string)
}

type SDKCache struct {
	cache SDK
	Init  func() SDK
}

//Lazy Getter
func (sdkCache SDKCache) Get() SDK {
	if sdkCache.cache == nil {
		sdkCache.cache = sdkCache.Init()
	}
	return sdkCache.cache
}
