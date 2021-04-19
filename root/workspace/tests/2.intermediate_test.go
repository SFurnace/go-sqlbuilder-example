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
	b := sqlbuilder.NewSelectBuilder()
	b.BuilderAs()


}

func TestSubQueryWithJoin(t *testing.T) {

}
