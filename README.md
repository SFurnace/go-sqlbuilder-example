# go-sqlbuilder

> 文档：[link](https://pkg.go.dev/github.com/huandu/go-sqlbuilder)

go-sqlbuilder 提供了一组灵活且强大的 SQL 构造方法，帮助用户构造可被标准库中提供的 `db#Query`、`db#Exec`、`Rows#Scan` 等方法使用的参数。

## Pros

1. 避免手工构造 SQL，在需要依据复杂的条件逐步构造 SQL 时可以减少错误。
2. 与 Go 语言标准库 `database/sql` 搭配地很好，没有引入过多抽象。
3. 包含一个零配置的 ORM，可根据 Struct 类型信息构造出合适的 Builder。
4. 支持 Builder/FormatStyle/FreeStyle 等 SQL 构造方式。
5. 默认支持 MySQL/PostgreSQL/SQLite 三种 SQL 风格。

## Cons

1. 对非标准的 SQL 语句支持较弱，需要用户掌握如何通过 `Var` 接口自行组装这类语句或者使用 FormatStyle/FreeStyle 的写法。
2. 为了简化接口、方便使用，使用 `Struct` 组装 SQL 时会忽略各种错误，需要使用者保证参数的正确性。

## 用法介绍

### 说明

为了方便演示，我从 ecm_websvr 库中拿了几个数据结构并简化了一下作为要处理的数据。

演示使用的数据结构放在 [这里](./root/workspace/tests/datatype.go)。 各个例子在 `root/workspace/test` 文件夹下，且都可以用 Dockerfile 构建出的镜像运行。

### 基础用法

sqlbuilder 库提供了六个最基本的 `Builder`，可以使用它们来构造 SQL 语句。

- CreateTableBuilder: Builder for CREATE TABLE.
- SelectBuilder: Builder for SELECT.
- InsertBuilder: Builder for INSERT.
- UpdateBuilder: Builder for UPDATE.
- DeleteBuilder: Builder for DELETE.
- UnionBuilder: Builder for UNION and UNION ALL.

每个 Builder 上都有 `Build()` 方法来返回最终的 SQL 字符串和参数列表。除此之外还有用来**收集参数**或者**构建 SQL** 的各类帮助函数。

Builder 方法会返回 Builder 本身方便链式调用，帮助函数一般会返回 string 用于组合语句。

---

#### Select Builder 用法举例

`NewSelectBuilder()` 返回一个 SelectBuilder，`In`、`Like`等方法可以帮助构建 SQL 语句并收集参数。

```go
func TestSelect0(t *testing.T) {
	b := sqlbuilder.NewSelectBuilder()
	b.Select("uin", "appId").From(CustomerTable).
		Where(b.In("uin", "2792294370"), b.Like("userName", "%tencent%")) // 多个条件之间是 AND 的关系

	expr, args := b.Build()
	fmt.Println(expr)
	// SELECT uin, appId FROM t_customer WHERE uin IN (?) AND userName LIKE ?
	fmt.Println(args)
	// [2792294370 %tencent%]
}
```

当在 SQL 中使用 `IN` 表达式而参数是 `slice` 时，可以用 sqlbuilder.List 把参数包装起来。

```go
func TestSelect1(t *testing.T) {
	b := sqlbuilder.Select("uin", "appId").From(CustomerTable)                      // 比 NewSelectBuilder 更简单的写法
	b.Where(b.In("uin", "1", "2", "3", "4"))                                        // In 接受变长参数
	b.Where(b.In("userName", sqlbuilder.List([]string{"name0", "name1", "name2"}))) // 多个 Where 之间也是 AND 的关系

	expr, args := b.Build()
	fmt.Println(expr)
	// SELECT uin, appId FROM t_customer WHERE uin IN (?, ?, ?, ?) AND userName IN (?, ?, ?)
	fmt.Println(args)
	// [1 2 3 4 name0 name1 name2]
}
```

SQL 中的“或”表达式用法如下，多个条件可以分开收集。

```go
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
```

---

#### 其他 Builder 用法举例

Update Builder 用法举例。

```go
func TestUpdateBasic(t *testing.T) {
	b := sqlbuilder.Update(NodeTable)
	b.Set(b.Assign("state", "OFFLINE")) // 可接受变长参数
	b.SetMore(b.Add("ispNum", 2))       // 添加赋值语句需要用 SetMore，用 Set 会覆盖掉之前的赋值
	b.Where(b.Like("zone", "%beijing%"))

	expr, args := b.Build()
	fmt.Println(expr)
	// UPDATE t_node SET state = ?, ispNum = ispNum + ? WHERE zone LIKE ?
	fmt.Println(args)
	// [OFFLINE 2 %beijing%]
}
```

Insert Builder 用法举例。

```go
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
```

### 中级用法

#### 表连接

`SelectBuilder` 上的 `Join` 和 `JoinWithOption` 方法可以用来构造连接查询。可以用 `As` 指定表的别名。

**注意**：表的别名不能通过参数的形式传给数据库。

```go
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
```

当需要显示声明连接方式时，要使用 `JoinWithOption`。

```go
// JoinWithOption sets expressions of JOIN with an option.
//
// It builds a JOIN expression like
//     option JOIN table ON onExpr[0] AND onExpr[1] ...
//
// Here is a list of supported options.
//     - FullJoin: FULL JOIN
//     - FullOuterJoin: FULL OUTER JOIN
//     - InnerJoin: INNER JOIN
//     - LeftJoin: LEFT JOIN
//     - LeftOuterJoin: LEFT OUTER JOIN
//     - RightJoin: RIGHT JOIN
//     - RightOuterJoin: RIGHT OUTER JOIN
func (sb *SelectBuilder) JoinWithOption(option JoinOption, table string, onExpr ...string) *SelectBuilder {
	// ...
}
```

---

#### 子查询

子查询的写法很简单，只需用 `SelectBuilder` 上的 `BuilderAs` 方法为其他 `Builder` 起别名就能将其查询结果做为一张表使用了。

```go

```

子查询也能在表连接等地方直接使用。

```go

```

### 高级用法

#### Struct

---

#### dbhelper

### 一些细节

#### Struct omitempty 只在 update 语句中生效

---

#### Struct.Insert 不能用 sqlbuilder.List 作为参数

---

#### Var() 与 SQL() 的使用

---

#### 其他帮助函数

##### Interpolate

##### Flatten

---

#### NULL 字段的处理方法

---

#### 数据库错误字段库


