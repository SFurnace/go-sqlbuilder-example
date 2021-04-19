# go-sqlbuilder

> 文档：[link](https://pkg.go.dev/github.com/huandu/go-sqlbuilder)

go-sqlbuilder 提供了一组灵活且强大的 SQL 构造方法，帮助用户构造可被标准库中提供的`db#Query`、`db#Exec`、`Rows#Scan`等方法使用的参数。

[TOC]

## Pros

1. 避免手工构造 SQL，在需要依据复杂的条件逐步构造 SQL 时可以减少错误。
2. 与 Go 语言标准库`database/sql`搭配地很好，没有引入过多抽象。
3. 包含一个零配置的 ORM，可根据 Struct 类型信息构造出合适的 Builder。
4. 支持 Builder/FormatStyle/FreeStyle 等 SQL 构造方式。
5. 默认支持 MySQL/PostgreSQL/SQLite 三种 SQL 风格。

## Cons

1. 对非标准的 SQL 语句支持较弱，需要用户掌握如何通过`Var`接口自行组装这类语句或者使用 FormatStyle/FreeStyle 的写法。
2. 为了简化接口、方便使用，使用`Struct`组装 SQL 时会忽略各种错误，需要使用者保证参数的正确性。
3. 没有提供从`Rows`等结构体中提取查询结果的方便函数，需要自己实现。

## 用法介绍

### 说明

为了方便演示，我从 ecm_websvr 库中拿了几个数据结构并简化了一下作为要处理的数据。

演示使用的数据结构放在 [这里](./root/workspace/tests/datatype.go)。 各个例子在`root/workspace/test`文件夹下，且都可以用 Dockerfile 构建出的镜像运行。

本文中主要介绍以`Builder`形式使用 go-sqlbuilder。FormatStyle/FreeStyle 请自行通过文档了解。

### 基础用法

sqlbuilder 库提供了六个最基本的`Builder`，可以使用它们来构造 SQL 语句。

- CreateTableBuilder: Builder for CREATE TABLE.
- SelectBuilder: Builder for SELECT.
- InsertBuilder: Builder for INSERT.
- UpdateBuilder: Builder for UPDATE.
- DeleteBuilder: Builder for DELETE.
- UnionBuilder: Builder for UNION and UNION ALL.

每个 Builder 上都有`Build()`方法来返回最终的 SQL 字符串和参数列表。除此之外还有用来**收集参数**或者**构建 SQL** 的各类帮助函数。

Builder 方法会返回 Builder 本身方便链式调用，帮助函数一般会返回 string 用于组合语句。

---

#### Select Builder 用法举例

`NewSelectBuilder()`返回一个 SelectBuilder，`In`、`Like`等方法可以帮助构建 SQL 语句并收集参数。

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

当在 SQL 中使用`IN`表达式而参数是`slice`时，可以用 sqlbuilder.List 把参数包装起来。

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

`SelectBuilder`上的`Join`和`JoinWithOption`方法可以用来构造连接查询。可以用`As`指定表的别名。

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

当需要显示声明连接方式时，要使用`JoinWithOption`。

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

子查询的写法很简单，将其他`Builder`作为参数直接使用即可。

```go
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
```

在 From 子句或表连接中，用`BuilderAs`方法为其他`Builder`起别名就能将其查询结果做为一张表使用了。

```go
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
```

### 高级用法

#### Struct

---

#### Struct Tags

##### omitempty 默认只在 update 语句中生效

用`Struct`作为轻量级 ORM 的时候，可以使用 omitempty

##### omitempty 可以指定在多个 tag 下生效

---

#### Struct.Insert 不能用 sqlbuilder.List 作为参数

`Struct`上的`InsertInto`、`ReplaceInto`等方法，可以接受变长参数但是不接收`sqlbuilder.List`的返回值作为参数。 由于`Struct`
上的接口会忽略错误，调试时可能很难发现这种错误。_建议所有新增的 SQL 都在调试阶段用`Interpolate`函数检查一遍_

```go
// InsertInto creates a new`InsertBuilder`with table name using verb INSERT INTO.
// By default, all exported fields of s are set as columns by calling`InsertBuilder#Cols`,
// and value is added as a list of values by calling`InsertBuilder#Values`.
//
// InsertInto never returns any error.
// If the type of any item in value is not expected, it will be ignored.
// If value is an empty slice,`InsertBuilder#Values`will not be called.
func (s *Struct) InsertInto(table string, value ...interface{}) *InsertBuilder
```

---

#### dbhelper.Struct

仅使用`Struct.Addr`和`Rows.Scan`等方法获取查询结果的话，需要写很多结构重复的代码。为了简化代码，我在`sqlbuilder.Struct`之上封装了
`dbhelper.Struct`，方便获取 *sql.DB 或 *sql.Tx 返回的查询结果。

```go
// 执行 SQL 查询并将结果存放至 result 中，result 必须是用于生成 Struct 的结构体类型的切片指针
func (s *Struct) Query(ctx context.Context, db *sql.DB, result interface{}, expr string, args ...interface{}) error

// 执行 SQL 查询并将结果存放至 result 中，result 必须是用于生成 Struct 的结构体类型的指针
func (s *Struct) QueryRow(ctx context.Context, db *sql.DB, result interface{}, expr string, args ...interface{}) error

// 执行 SQL
func (s *Struct) Exec(ctx context.Context, db *sql.DB, expr string, args ...interface{}) (sql.Result, error)
```

_每个代表数据库返回结果的结构体都应有对应的 Struct方便他人复用，其名称惯例使用"S结构体名称"_

### 一些细节

#### Var() 与 SQL() 的使用

`Var`方法用于收集 SQL 语句中的参数，其接收一个参数值，返回一个 SQL 语句中使用的占位符。

```go
// Var returns a placeholder for value.
func (c *Cond) Var(value interface{}) string
```

`SQL`方法用于在 SQL 语句中插入任意内容。 为了明确`SQL`方法插入的位置，每种`Builder`中都会记录一个 marker，代表当前插入 SQL 的位置。

```go
// Select 语句使用的 marker 列表
const (
	selectMarkerInit injectionMarker = iota
	selectMarkerAfterSelect
	selectMarkerAfterFrom
	selectMarkerAfterJoin
	selectMarkerAfterWhere
	selectMarkerAfterGroupBy
	selectMarkerAfterOrderBy
	selectMarkerAfterLimit
	selectMarkerAfterFor
)
```

当你使用`Builder`上的构建方法时，就会改变其记录的 marker。

```go
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
```

`Var` 与 `SQL` 结合，就可以自定义 sqlbuilder 中没有支持的表达式。_为了方便复用，建议自定义的表达式统一放在 dbhelper/utils.go 里_

```go
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
```

---

#### 其他帮助函数

##### sqlbuilder.Interpolate

`Interpolate`函数可以将参数内插进 SQL 语句中，获得完整的 SQL 语句。在调试 SQL 语句或者 Driver 没有实现参数传递的情况下可以使用。

```go
func TestInterpolate(t *testing.T) {
	b := sqlbuilder.Select("uin", "appId").From(CustomerTable)
	b.Where(b.In("uin", "1", "2", "3", "4"))
	b.Where(b.In("userName", sqlbuilder.List([]string{"name0", "name1", "name2"})))

	fmt.Println(sqlbuilder.MySQL.Interpolate(b.Build()))
	// SELECT uin, appId FROM t_customer WHERE uin IN ('1', '2', '3', '4') AND userName IN ('name0', 'name1', 'name2')
}
```

##### sqlbuilder.Flatten

将其他类型的切片转换成`[]interface{}`的帮助函数。

---

#### NULL 字段的处理方法

虽然我们应尽量避免表中存在`Nullable`的字段，若不能避免的话（如使用了`JSON`/`TEXT`类型的字段而且 MySQL 版本比较旧），有如下两种处理方法。

1. 使用 database/sql 包中的`NullString`等类型
2. 使用`COALESCE`函数来排除查询结果中的`NULL`值

---

#### 数据库错误处理

根据数据库返回的错误信息对错误做处理是一个常见的诉求，使用社区维护的错误码库来判断错误类型是一种正确性较高的做法。 MySQL
的错误码库可以使用 [https://github.com/VividCortex/mysqlerr](https://github.com/VividCortex/mysqlerr).
