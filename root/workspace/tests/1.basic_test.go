package tests

import (
	"fmt"
	"testing"

	"github.com/huandu/go-sqlbuilder"
)

func TestSelectBasic(t *testing.T) {
	b := sqlbuilder.NewSelectBuilder()
	b.Select("uin", "appId").From(CustomerTable).
		Where(b.In("uin", "2792294370"), b.Like("userName", "%tencent%")) // 多个条件之间是 AND 的关系

	expr, args := b.Build()
	fmt.Println(expr)
	// SELECT uin, appId FROM t_customer WHERE uin IN (?) AND userName LIKE ?
	fmt.Println(args)
	// [2792294370 %tencent%]
}

func TestSelectWithList(t *testing.T) {
	b := sqlbuilder.Select("uin", "appId").From(CustomerTable)                      // 比 NewSelectBuilder 更简单的写法
	b.Where(b.In("uin", "1", "2", "3", "4"))                                        // In 接受变长参数
	b.Where(b.In("userName", sqlbuilder.List([]string{"name0", "name1", "name2"}))) // slice 可以用 List 方法包装

	expr, args := b.Build()
	fmt.Println(expr)
	// SELECT uin, appId FROM t_customer WHERE uin IN (?, ?, ?, ?) AND userName IN (?, ?, ?)
	fmt.Println(args)
	// [1 2 3 4 name0 name1 name2]
}

func TestSelectWithOr(t *testing.T) {
	var (
		appIds     = []int64{3, 4, 5}
		zones      = []string{"zone0", "zone1"}
		conditions []string
	)

	b := sqlbuilder.NewSelectBuilder().From(DeviceTable)
	b.Select("appId", b.As("COUNT(*)", "num"))
	if len(appIds) > 0 { // 检查长度可以避免 SQL 中的语法错误，使用 List/Values/In 等方法的时候应该注意
		conditions = append(conditions, b.In("appId", sqlbuilder.List(appIds))) // 收集Or的条件
	}
	if len(zones) > 0 {
		conditions = append(conditions, b.In("zone", sqlbuilder.List(zones)))
	}
	b.Where(b.Or(conditions...)).Limit(10) // 用 OR 来连接

	expr, args := b.Build()
	fmt.Println(expr)
	// SELECT appId, COUNT(*) AS num FROM t_device WHERE (appId IN (?, ?, ?) OR zone IN (?, ?)) LIMIT 10
	fmt.Println(args)
	// [3 4 5 zone0 zone1]
}

func TestUpdateBasic(t *testing.T) {
	b := sqlbuilder.Update(NodeTable)
	b.Set(b.Assign("state", "OFFLINE"))
	b.SetMore(b.Add("ispNum", 2)) // 添加赋值语句需要用 SetMore，用 Set 会覆盖掉之前的赋值
	b.Where(b.Like("zone", "%beijing%"))

	expr, args := b.Build()
	fmt.Println(expr)
	// UPDATE t_node SET state = ?, ispNum = ispNum + ? WHERE zone LIKE ?
	fmt.Println(args)
	// [OFFLINE 2 %beijing%]
}

func TestInsertBasic(t *testing.T) {
	values := [][3]string{{"1", "2", "3"}, {"4", "5", "6"}}

	b := sqlbuilder.InsertIgnoreInto("a_table")
	b.Cols("col0", "col1", "col2")
	for i := range values { // 记得检查参数
		b.Values(values[i][0], values[i][1], values[i][2])
	}

	expr, args := b.Build()
	fmt.Println(expr)
	// INSERT IGNORE INTO a_table (col0, col1, col2) VALUES (?, ?, ?), (?, ?, ?)
	fmt.Println(args)
	// [1 2 3 4 5 6]
}
