package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/huandu/go-sqlbuilder"
)

func TestSelectStruct(t *testing.T) {
	b := SCustomerEx.SelectFrom(CustomerTable)
	b.Where(b.Like("uin", "%tencent%"))

	expr, args := b.Build()
	fmt.Println(expr)
	// SELECT t_customer.userIndustry, t_customer.userArchitect, t_customer.userSeller, t_customer.picUrl,
	//   t_customer.industryGrade, t_customer.uin, t_customer.appId, t_customer.userName, t_customer.remarkName
	// FROM t_customer WHERE uin LIKE ?
	fmt.Println(args)
	// [%tencent%]

	var result []CustomerEx
	if err := SCustomerEx.Query(context.Background(), DB, &result, expr, args); err != nil { // 获取数据
		t.Fatal(err)
	}
	fmt.Println(result)
}

func TestSelectStructTag(t *testing.T) {
	b := SCustomer.SelectFromForTag("t_customer AS tc", "only_id").
		JoinWithOption(sqlbuilder.InnerJoin, "t_device AS td", "tc.appId = td.appId")

	expr, args := b.Build()
	fmt.Println(expr)
	// SELECT t_customer AS tc.uin, t_customer AS tc.appId FROM t_customer AS tc INNER JOIN t_device AS td ON tc.appId = td.appId
	fmt.Println(args)
	// []

	var result []Customer
	if err := SCustomer.Query(context.Background(), DB, &result, expr, args); err != nil {
		t.Fatal(err)
	}
}
