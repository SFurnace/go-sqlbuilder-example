package tests

import (
	"encoding/json"
	"time"
)

/* 客户 */

type (
	Customer struct {
		Uin          string `db:"uin" fieldtag:"only_id"`
		AppID        int64  `db:"appId" fieldtag:"only_id"`
		CustomerName string `db:"userName"`
		RemarkName   string `db:"remarkName"`
	}
	CustomerEx struct {
		Customer
		CustomerIndustry  string `db:"userIndustry"`
		CustomerArchitect string `db:"userArchitect"`
		CustomerSeller    string `db:"userSeller"`
		PicUrl            string `db:"picUrl"`
		IndustryGrade     string `db:"industryGrade"`
	}
	CustomerFull struct {
		CustomerEx
		TimeInfo
	}
)

/* 节点 */

type (
	Node struct {
		IdcID    int    `db:"idcId"`
		ZoneID   int    `db:"zoneId"`
		Zone     string `db:"zone"`
		RegionID int    `db:"regionId"`
		Region   string `db:"region"`
		State    string `db:"state"`
	}
	NodeEx struct {
		Node
		Country             string          `db:"country"`
		Area                string          `db:"area"`
		Province            string          `db:"province"`
		City                string          `db:"city"`
		ISP                 string          `db:"isp"`
		ISPNum              int             `db:"ispNum"`
		InstanceFamilyTypes json.RawMessage `db:"instanceFamilyTypes"`
	}
	NodeFull struct {
		NodeEx
		TimeInfo
	}
)

/* 机器 */

type (
	Device struct {
		InstanceID   string `db:"instanceId"`
		InstanceName string `db:"instanceName"`
		AppID        int64  `db:"appId"`
		Zone         string `db:"zone"`
		InstanceType string `db:"instanceType"`
		State        string `db:"state"`
	}
	DeviceEx struct {
		Device
	}
	DeviceFull struct {
		DeviceEx
		TimeInfo
		TerminateTime time.Time `db:"terminateTime"`
	}
)

/* Component */

type TimeInfo struct {
	CreateTime time.Time `db:"createTime" fieldopt:"omitempty"`
	UpdateTime time.Time `db:"updateTime" fieldopt:"omitempty"`
}
