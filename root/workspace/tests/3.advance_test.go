package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/huandu/go-sqlbuilder"

	"pers.drcz/tests/sqlbuilder/comm/dbhelper"
)

func TestSelectStruct(t *testing.T) {
	b := dbhelper.S(CustomerEx{}).SelectFrom(CustomerTable)
	b.Where(b.Like("uin", "%tencent%"))

	expr, args := b.Build()
	fmt.Println(expr)
	// SELECT t_customer.userIndustry, t_customer.userArchitect, t_customer.userSeller, t_customer.picUrl, t_customer.industryGrade, t_customer.uin, t_customer.appId, t_customer.userName, t_customer.remarkName FROM t_customer WHERE uin LIKE ?
	fmt.Println(args)
	// [%tencent%]

	var result []CustomerEx
	if err := dbhelper.PullStructs(context.Background(), DB, &result, b); err != nil { // 获取数据
		t.Fatal(err)
	}
	fmt.Println(result)
}

func TestSelectStructTag(t *testing.T) {
	b := dbhelper.S(Customer{}).SelectFromForTag("t_customer", "only_id").
		JoinWithOption(sqlbuilder.InnerJoin, "t_device", "t_customer.appId = t_device.appId")

	expr, args := b.Build()
	fmt.Println(expr)
	// SELECT t_customer.uin, t_customer.appId FROM t_customer INNER JOIN t_device ON t_customer.appId = t_device.appId
	fmt.Println(args)
	// []

	var result []Customer
	if err := dbhelper.PullTagStructs(context.Background(), DB, "only_id", &result, b); err != nil {
		t.Fatal(err)
	}
}
