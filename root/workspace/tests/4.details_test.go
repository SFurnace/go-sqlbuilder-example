package tests

import (
	"fmt"
	"testing"

	"github.com/huandu/go-sqlbuilder"
)

func TestInterpolate(t *testing.T) {
	b := sqlbuilder.Select("uin", "appId").From(CustomerTable)
	b.Where(b.In("uin", "1", "2", "3", "4"))
	b.Where(b.In("userName", sqlbuilder.List([]string{"name0", "name1", "name2"})))

	fmt.Println(sqlbuilder.MySQL.Interpolate(b.Build()))
	// SELECT uin, appId FROM t_customer WHERE uin IN ('1', '2', '3', '4') AND userName IN ('name0', 'name1', 'name2')
}
