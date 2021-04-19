package tests

import (
	"fmt"
	"testing"

	"github.com/huandu/go-sqlbuilder"
)

func TestJoinSimple(t *testing.T) {
	b := sqlbuilder.NewSelectBuilder()
	b.Select("tc.appId", "tn.zone", "td.instanceType", b.As("COUNT(*)", "num")).
		From(b.As(CustomerTable, "tc")).
		Join(b.As(DeviceTable, "td"), "tc.appId = td.appId").
		Join(b.As(NodeTable, "tn"), "tn.zone = td.zone").
		GroupBy("tc.appId", "tn.zone", "td.instanceType").
		Having(b.G("COUNT(*)", 1)) // 查询实例数量大于 1 的机型信息

	expr, args := b.Build()
	fmt.Println(expr)
	// SELECT tc.appId, tn.zone, td.instanceType, COUNT(*) AS num FROM t_customer AS tc
	//   JOIN t_device AS td ON tc.appId = td.appId
	//   JOIN t_node AS tn ON tn.zone = td.zone
	// GROUP BY tc.appId, tn.zone, td.instanceType
	// HAVING COUNT(*) > ?
	fmt.Println(args)
	// [1]
}

func TestSubQuery(t *testing.T) {
	s := sqlbuilder.Select("appId").From(CustomerTable)
	s.Where(s.Like("userName", "%tencent%"))

	b := sqlbuilder.Select("instanceId", "appId", "zone").From(DeviceTable)
	b.Where(b.In("appId", s))

	expr, args := b.Build()
	fmt.Println(expr)
	// SELECT instanceId, appId, zone FROM t_device WHERE appId IN (SELECT appId FROM t_customer WHERE userName LIKE ?)
	fmt.Println(args)
	// [%tencent%]
}

func TestSubQueryWithJoin(t *testing.T) {
	s := sqlbuilder.Select("appId").From(CustomerTable)
	s.Where(s.Like("userName", "%tencent%"))

	b := sqlbuilder.NewSelectBuilder()
	b.Select("instanceId").From(b.As(DeviceTable, "td")).Join(b.BuilderAs(s, "tc"), "tc.appId = td.appId")

	expr, args := b.Build()
	fmt.Println(expr)
	// SELECT instanceId FROM t_device AS td JOIN (SELECT appId FROM t_customer WHERE userName LIKE ?) AS tc ON tc.appId = td.appId
	fmt.Println(args)
	// [%tencent%]
}
