package tests

import (
	"fmt"
	"testing"

	"github.com/huandu/go-sqlbuilder"
)

// 生成测试用的数据表
func TestCreateDB(t *testing.T) {
	builders := []sqlbuilder.Builder{
		sqlbuilder.CreateTable(CustomerTable).IfNotExists().Define(
			"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',",
			"`uin` bigint(20) NOT NULL DEFAULT '0' COMMENT '资源实际拥有者 uid',",
			"`appId` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户appId',",
			"`userName` varchar(255) NOT NULL COMMENT '用户名称',",
			"`userArchitect` varchar(255) NOT NULL COMMENT '架构师',",
			"`userSeller` varchar(255) NOT NULL COMMENT '销售员',",
			"`userIndustry` varchar(255) NOT NULL COMMENT '用户行业',",
			"`remarkName` varchar(255) NOT NULL DEFAULT '-' COMMENT '备注名称',",
			"`createTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',",
			"`updateTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',",
			"`picUrl` varchar(255) NOT NULL,",
			"`industryGrade` varchar(255) NOT NULL DEFAULT '',",
			"PRIMARY KEY (`id`),",
			"UNIQUE KEY `uin` (`uin`),",
			"UNIQUE KEY `index_appid` (`appId`)",
		),

		sqlbuilder.CreateTable(NodeTable).IfNotExists().Define(
			"`idcId` int(11) NOT NULL DEFAULT '0' COMMENT '机房id',",
			"`zoneId` int(11) NOT NULL DEFAULT '0' COMMENT 'zoneId',",
			"`zone` varchar(255) NOT NULL DEFAULT '' COMMENT 'zone',",
			"`regionId` int(11) NOT NULL DEFAULT '0' COMMENT 'regionId',",
			"`region` varchar(255) NOT NULL COMMENT 'region',",
			"`state` varchar(255) NOT NULL DEFAULT 'NORMAL',",
			"`country` varchar(255) NOT NULL DEFAULT '' COMMENT '国家代码',",
			"`area` varchar(255) NOT NULL DEFAULT '' COMMENT '区域代码',",
			"`province` varchar(255) NOT NULL DEFAULT '' COMMENT '省份',",
			"`city` varchar(255) NOT NULL DEFAULT '' COMMENT '城市',",
			"`isp` varchar(255) NOT NULL DEFAULT '' COMMENT '运营商',",
			"`ispNum` int(11) NOT NULL DEFAULT '0' COMMENT '节点支持运营商的数量',",
			"`instanceFamilyTypes` varchar(256) NOT NULL DEFAULT '' COMMENT '区域支持的机型列表',",
			"`createTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',",
			"`updateTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',",
			"PRIMARY KEY (`zoneId`),",
			"KEY `index_zone_to_node` (`zone`)",
		),

		sqlbuilder.CreateTable(DeviceTable).IfNotExists().Define(
			"`instanceId` varchar(255) NOT NULL COMMENT '实例ID',",
			"`instanceName` varchar(255) NOT NULL DEFAULT '' COMMENT '实例显示名称',",
			"`appId` bigint(20) NOT NULL DEFAULT '0' COMMENT 'appId',",
			"`zone` varchar(255) NOT NULL DEFAULT '' COMMENT 'zone',",
			"`instanceType` varchar(255) NOT NULL DEFAULT '' COMMENT '机型配置ID',",
			"`state` varchar(255) NOT NULL DEFAULT '' COMMENT '实例状态 UNKNOWN-未知状态, UPDATING-更新中, PENDING-创建中, LAUNCH_FAILED-创建失败, RUNNING-运行中, STOPPED-关机, STARTING-开机中, STOPPING-关机中, REBOOTING-重启中, SHUTDOWN-停止待销毁, TERMINATING-销毁中',",
			"`createTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',",
			"`updateTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',",
			"`terminateTime` datetime NOT NULL DEFAULT '2099-12-31 23:59:59',",
			"PRIMARY KEY (`instanceId`)",
		),
	}

	for _, b := range builders {
		expr, args := b.Build()
		if _, err := DB.Exec(expr, args...); err != nil {
			fmt.Println(err)
			fmt.Println(expr)
		} else {
			fmt.Println("create table ok!")
		}
	}
}

// 生成测试数据
func TestGenerateData(t *testing.T) {

}
