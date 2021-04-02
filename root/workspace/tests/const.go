package tests

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	CustomerTable = "t_customer"
	NodeTable     = "t_node"
	DeviceTable   = "t_device"
)

const (
	customerNum = 100
	nodeNum     = 50
	deviceNum   = 1000
)

var (
	IndustryGrades = []string{"一级", "二级", "三级"}
	AreaList       = []string{
		"china-central", "china-east", "china-north", "china-northwest", "china-south", "china-southwest",
		"china-northeast",
	}
	ProvinceList = []string{
		"china-central-henan", "china-east-shandong", "china-south-guangdong",
		"china-southwest-sichuan", "china-central-hubei", "china-central-hunan",
		"china-southwest-chongqing", "china-south-guangxi", "china-southwest-yunnan",
	}
	CityList = []string{
		"china-central-henan-zhengzhou", "china-east-shandong-qingdao", "china-east-zhejiang-hangzhou",
		"china-east-jiangsu-nanjing", "china-east-shandong-jinan", "china-north-hebei-shijiazhuang",
		"china-north-beijing-beijing", "china-northwest-shanxi-xian", "china-east-shanghai-shanghai",
		"china-south-guangdong-guangzhou", "china-south-guangdong-shenzhen", "china-southwest-sichuan-chengdu",
		"china-central-hubei-wuhan", "china-south-guangdong-dongguan",
	}
	ISPList                 = []string{"CTCC;CUCC;CMCC", "CTCC", "CUCC", "CMCC"}
	InstanceFamilyTypesList = []string{
		`["S4", "SN3ne"]`, `["S4", "SN3ne", "S4e10G"]`, `["SN3ne"]`, `["S3", "SN3ne", "S4"]`,
		`["IT11", "S4", "SN3ne", "S4e10G"]`, `["SN3ne", "S4", "IT11"]`, `["S4", "S4e10G", "SN3ne"]`,
	}
	NodeStateList     = []string{"NORMAL", "OFFLINE", "SELLOUT"}
	InstanceStateList = []string{"LAUNCH_FAILED", "RUNNING", "TERMINATED", "STOPPED", "PENDING", "SHUTDOWN"}
	InstanceTypeList  = []string{
		"S4.MEDIUM4", "SN3ne.LARGE8", "S4.2XLARGE16", "S4.4XLARGE32", "SN3ne.SMALL2", "S4.8XLARGE64", "S3.SMALL2",
		"SN3ne.8XLARGE64", "SN3ne.2XLARGE32", "S4.LARGE8", "S4.4XLARGE64", "S4.6XLARGE48", "SN3ne.LARGE16",
		"S4.LARGE16", "S3.2XLARGE16", "SN3ne.6XLARGE64", "IT5.4XLARGE32", "S4.6XLARGE64", "SN3ne.4XLARGE64",
		"S4.2XLARGE32", "IT11.2XLARGE24", "S3.MEDIUM4", "S4.MEDIUM8", "S3.LARGE8", "SN3ne.SMALL8", "SN3ne.LARGE32",
	}
)

var (
	DB *sql.DB
)

func init() {
	var err error
	const DSN = "tester:tester123@tcp(localhost:3306)/test?charset=utf8"

	if DB, err = sql.Open("mysql", DSN); err != nil {
		panic(err)
	}
}
