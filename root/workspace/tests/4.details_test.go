package tests

import (
	"fmt"
	"testing"

	"github.com/huandu/go-sqlbuilder"
)

func TestSQLBase(t *testing.T) {
	b := sqlbuilder.NewSelectBuilder()
	b.SQL("/*Before ALL*/").
		Select("appId").SQL("/*After SELECT*/").
		From(CustomerTable).SQL("/*After FROM*/").
		Where(b.E("uin", "testUin")).SQL("/*After ALL*/")

	expr, args := b.Build()
	fmt.Println(expr)
	// /*Before ALL*/ SELECT appId /*After SELECT*/ FROM t_customer /*After FROM*/ WHERE uin = ? /*After ALL*/
	fmt.Println(args)
	// [testUin]
}

func REGEXP(c sqlbuilder.Cond, field, pat string) string {
	return fmt.Sprintf("%s REGEXP %s", field, c.Var(pat))
}

func TestCustomerFunc(t *testing.T) {
	b := sqlbuilder.Select("appId", "userName").From(CustomerTable)
	b.Where(REGEXP(b.Cond, "userName", ".*tencent.*"))

	expr, args := b.Build()
	fmt.Println(expr)
	// SELECT appId, userName FROM t_customer WHERE userName REGEXP ?
	fmt.Println(args)
	// [.*tencent.*]
}

func TestInterpolate(t *testing.T) {
	b := sqlbuilder.Select("uin", "appId").From(CustomerTable)
	b.Where(b.In("uin", "1", "2", "3", "4"))
	b.Where(b.In("userName", sqlbuilder.List([]string{"name0", "name1", "name2"})))

	fmt.Println(sqlbuilder.MySQL.Interpolate(b.Build()))
	// SELECT uin, appId FROM t_customer WHERE uin IN ('1', '2', '3', '4') AND userName IN ('name0', 'name1', 'name2')
}
