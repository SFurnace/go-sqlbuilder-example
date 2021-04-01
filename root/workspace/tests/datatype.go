package tests

import "pers.drcz/tests/sqlbuilder/comm/dbhelper"

/* 客户 */

type Customer struct {
	Uin          string `db:"uin"`
	AppID        int64  `db:"appId"`
	CustomerName string `db:"userName"`
	RemarkName   string `db:"remarkName"`
}

var SCustomer = dbhelper.NewStruct(Customer{})

type CustomerEx struct {
	Customer
	CustomerIndustry  string `db:"userIndustry"`
	CustomerArchitect string `db:"userArchitect"`
	CustomerSeller    string `db:"userSeller"`
	PicUrl            string `db:"picUrl"`
	IndustryGrade     string `db:"industryGrade"`
	TimeInfo
}

var SCustomerEx = dbhelper.NewStruct(CustomerEx{})

/* 节点 */

type Node struct {
	IdcID    int    `db:"idcId"`
	ZoneID   int    `db:"zoneId"`
	Zone     string `db:"zone"`
	RegionID int    `db:"regionId"`
	Region   string `db:"region"`
	State    string `db:"state"`
}

var SNode = dbhelper.NewStruct(Node{})

type NodeEx struct {
	Node
	Country             string `db:"country"`
	Area                string `db:"area"`
	Province            string `db:"province"`
	City                string `db:"city"`
	ISP                 string `db:"isp"`
	ISPNum              int    `db:"ispNum"`
	InstanceFamilyTypes string `db:"instanceFamilyTypes"`
	TimeInfo
}

var SNodeEx = dbhelper.NewStruct(NodeEx{})

/* 机器 */

type Device struct {
	InstanceID   string `db:"instanceId"`
	InstanceName string `db:"instanceName"`
	AppID        int64  `db:"appId"`
	Zone         string `db:"zone"`
	InstanceType string `db:"instanceType"`
	State        string `db:"state"`
}

var SDevice = dbhelper.NewStruct(Device{})

type DeviceEx struct {
	Device
	TimeInfo
	TerminateTime string `db:"terminateTime"`
}

var SDeviceEx = dbhelper.NewStruct(DeviceEx{})

/* Component */

type TimeInfo struct {
	CreateTime string `db:"createTime"`
	UpdateTime string `db:"updateTime"`
}
