package pkg

//go:generate sqlgen -t pkg.Customer -tn t_customer -out dao/customer.go -db EcmAdminDB.DB -map AppID:int64;Uin:int64 -group IndustryGrade:string
type Customer struct {
	Uin               int64  `db:"uin"`
	AppID             int64  `db:"appId"`
	CustomerName      string `db:"userName"`
	CustomerIndustry  string `db:"userIndustry"`
	CustomerArchitect string `db:"userArchitect"`
	CustomerSeller    string `db:"userSeller"`
	RemarkName        string `db:"remarkName"`
	PicUrl            string `db:"picUrl"`        // 客户对应的url
	IndustryGrade     string `db:"industryGrade"` // 客户行业定级
}

//go:generate sqlgen -t pkg.Node -tn t_node_info -out dao/node.go -db EcmDB.DB -map ZoneID:int;Zone:string -group ISPNum:int;State:string
type Node struct {
	ZoneID              int    `db:"zoneId"`   // zone_id
	Country             string `db:"country"`  // 国家
	Area                string `db:"area"`     // 区域
	Province            string `db:"province"` // 省份
	City                string `db:"city"`     // 城市
	Zone                string `db:"zone"`     // zone
	RegionID            int    `db:"regionId"`
	ISP                 string `db:"isp"` // 运营商
	Region              string `db:"region"`
	ISPNum              int    `db:"ispNum"`          // 运营商数量
	NodeDescription     string `db:"nodeDescription"` // 节点描述
	State               string `db:"state"`           // NORMAL/LIMITED/SELLOUT/OFFLINE
	CreateTime          string `db:"createTime"`
	UpdateTime          string `db:"updateTime"`
	IdcId               int    `db:"idcId"` // 机房id
	InstanceFamilyTypes string `db:"instanceFamilyTypes"`
	OBNDSuported        int    `db:"OBNDSuported"` // 机房是否支持obnd
	CTCCCrid            string `db:"ctccCrid"`     // 电信crid
	CMCCCrid            string `db:"cmccCrid"`     // 移动crid
	CUCCCrid            string `db:"cuccCrid"`     // 联通crid
	LbFlag              int    `db:"lbFlag"`       // 机房是否支持lb
	SDWANSupported      int    `db:"sdwanSupported"`
	UnifiedCTCCCrid     string `db:"unifiedCtccCrid"` // 电信全国统一分级CRID
	UnifiedCMCCCrid     string `db:"unifiedCmccCrid"` // 移动全国统一分级CRID
	UnifiedCUCCCrid     string `db:"unifiedCuccCrid"` // 联通全国统一分级CRID
	IPV6Supported       int    `db:"ipv6Supported"`   // 是否支持ipv6
}

//go:generate sqlgen -t pkg.Device -tn t_device_info -out dao/device.go -db EcmDB.DB -map InstanceID:string -group AppID:int64;Zone:string;Region:string
type Device struct {
	InstanceID           string `db:"instanceId"`
	ProjectID            int64  `db:"projectId"`
	InstanceName         string `db:"instanceName"`
	HostName             string `db:"hostName"`
	AppID                int64  `db:"appId"`
	ModuleID             string `db:"moduleId"`
	Zone                 string `db:"zone"`
	Region               string `db:"region"`
	DeviceLanIp          string `db:"deviceLanIp"`
	DeviceWanIp          string `db:"deviceWanIp"`
	ImageID              string `db:"imageId"`
	InstanceType         string `db:"instanceType"`
	Bandwidth            int64  `db:"bandwidth"`
	BandwidthIn          int64  `db:"bandwidthIn"`
	SystemDiskSize       int    `db:"systemDiskSize"`
	DataDiskSize         int    `db:"dataDiskSize"`
	State                string `db:"state"`
	ISP                  string `db:"isp"`
	CreateTime           string `db:"createTime"`
	UpdateTime           string `db:"updateTime"`
	TerminateTime        string `db:"terminateTime"`
	LatestOperation      string `db:"latestOperation"`
	LatestOperationState string `db:"latestOperationState"`
	RestrictState        string `db:"restrictState"`
	UUID                 string `db:"uuid"`
	ExpireTime           string `db:"expireTime"`   // 过期时间（预付费有效）
	IsolatedTime         string `db:"isolatedTime"` // 隔离时间（预付费有效）
	RenewFlag            int    `db:"renewFlag"`    // 自动续费标志
	OrderID              int64  `db:"orderId"`      // 新购或者变配的order
	PayMode              int    `db:"payMode"`      // 支付方式
	SystemDiskID         string `db:"systemDiskId"`
	DataDiskID           string `db:"dataDiskId"`
	NewFlag              int    `db:"newFlag"` // 是否展示新实例标志
	VpcID                string `db:"vpcId"`
	SubnetID             string `db:"subnetId"`
	BillingType          int    `db:"billingType"` // 计费方式，0，按cpu计费，1，按小时计费，2，按月计费
	PosId                string `db:"posId"`
	RackId               string `db:"rackId"`
	SwitchId             string `db:"switchId"`
}
