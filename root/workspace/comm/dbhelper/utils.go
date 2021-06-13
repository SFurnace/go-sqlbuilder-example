package dbhelper

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/huandu/go-sqlbuilder"

	ecmlog "pers.drcz/tests/sqlbuilder/comm/log"
)

const unknownURI = "()"

/* Type Definition */

type (
	Executor interface {
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	}

	Cond     = *sqlbuilder.SelectBuilder
	CondFunc func(b Cond)
)

var (
	boolType    = reflect.TypeOf(false)
	intType     = reflect.TypeOf(0)
	int64Type   = reflect.TypeOf(int64(0))
	float64Type = reflect.TypeOf(float64(0))
	stringType  = reflect.TypeOf("")
)

/* Simple SQL helper - 单行单列结果查询 */

func getValue(ctx context.Context, db Executor, b sqlbuilder.Builder, vt reflect.Type) (interface{}, error) {
	expr, args := b.Build()
	row := QueryRow(ctx, db, expr, args...)

	tmp := reflect.New(vt)
	if err := row.Scan(tmp.Interface()); err != nil {
		return nil, err
	}
	return tmp.Elem().Interface(), nil
}

// GetBool 查询单个 bool
func GetBool(ctx context.Context, db Executor, b sqlbuilder.Builder) (bool, error) {
	if v, err := getValue(ctx, db, b, boolType); err != nil {
		return false, err
	} else {
		return v.(bool), nil
	}
}

// GetInt 查询单个 int
func GetInt(ctx context.Context, db Executor, b sqlbuilder.Builder) (int, error) {
	if v, err := getValue(ctx, db, b, intType); err != nil {
		return 0, err
	} else {
		return v.(int), nil
	}
}

// GetInt64 查询单个 int64
func GetInt64(ctx context.Context, db Executor, b sqlbuilder.Builder) (int64, error) {
	if v, err := getValue(ctx, db, b, int64Type); err != nil {
		return 0, err
	} else {
		return v.(int64), nil
	}
}

// GetFloat64 查询单个 float64
func GetFloat64(ctx context.Context, db Executor, b sqlbuilder.Builder) (float64, error) {
	if v, err := getValue(ctx, db, b, float64Type); err != nil {
		return 0, err
	} else {
		return v.(float64), nil
	}
}

// GetString 查询单个 string
func GetString(ctx context.Context, db Executor, b sqlbuilder.Builder) (string, error) {
	if v, err := getValue(ctx, db, b, stringType); err != nil {
		return "", err
	} else {
		return v.(string), nil
	}
}

/* Simple SQL helper - 单列结果查询 */

func pullValues(ctx context.Context, db Executor, b sqlbuilder.Builder, vt reflect.Type) (interface{}, error) {
	expr, args := b.Build()
	rows, err := Query(ctx, db, expr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tmp, result := reflect.New(vt), reflect.MakeSlice(reflect.SliceOf(vt), 0, 0)
	for rows.Next() {
		if err = rows.Scan(tmp.Interface()); err != nil {
			return nil, err
		}
		result = reflect.Append(result, tmp.Elem())
	}

	return result.Interface(), nil
}

// PullBools 查询单列 bool
func PullBools(ctx context.Context, db Executor, b sqlbuilder.Builder) ([]bool, error) {
	if result, err := pullValues(ctx, db, b, boolType); err != nil {
		return nil, err
	} else {
		return result.([]bool), nil
	}
}

// PullInts 查询单列 int
func PullInts(ctx context.Context, db Executor, b sqlbuilder.Builder) ([]int, error) {
	if result, err := pullValues(ctx, db, b, intType); err != nil {
		return nil, err
	} else {
		return result.([]int), nil
	}
}

// PullInt64s 查询单列 int64
func PullInt64s(ctx context.Context, db Executor, b sqlbuilder.Builder) ([]int64, error) {
	if result, err := pullValues(ctx, db, b, int64Type); err != nil {
		return nil, err
	} else {
		return result.([]int64), nil
	}
}

// PullFloat64s 查询单列 float64
func PullFloat64s(ctx context.Context, db Executor, b sqlbuilder.Builder) ([]float64, error) {
	if result, err := pullValues(ctx, db, b, float64Type); err != nil {
		return nil, err
	} else {
		return result.([]float64), nil
	}
}

// PullStrings 查询单列字符串
func PullStrings(ctx context.Context, db Executor, b sqlbuilder.Builder) ([]string, error) {
	if result, err := pullValues(ctx, db, b, stringType); err != nil {
		return nil, err
	} else {
		return result.([]string), nil
	}
}

/* Simple SQL helper - 结构体结果查询 */

// GetTagStruct 查询单个结构体
func GetTagStruct(ctx context.Context, db Executor, tag string, out interface{}, b sqlbuilder.Builder) error {
	expr, args := b.Build()
	return S(out).TagQueryRow(ctx, db, out, tag, expr, args...)
}

// PullTagStructs 查询结构体slice
func PullTagStructs(ctx context.Context, db Executor, tag string, out interface{}, b sqlbuilder.Builder) error {
	expr, args := b.Build()
	return S(out).TagQuery(ctx, db, out, tag, expr, args...)
}

// GetStruct 查询单个结构体
func GetStruct(ctx context.Context, db Executor, out interface{}, b sqlbuilder.Builder) error {
	return GetTagStruct(ctx, db, "", out, b)
}

// PullStructs 查询结构体slice
func PullStructs(ctx context.Context, db Executor, out interface{}, b sqlbuilder.Builder) error {
	return PullTagStructs(ctx, db, "", out, b)
}

/* SQL Execute Helper */

// Query 执行查询
func Query(ctx context.Context, db Executor, expr string, args ...interface{}) (*sql.Rows, error) {
	// start := time.Now()
	rows, err := db.QueryContext(ctx, expr, args...)
	// go reportDBCall(start, time.Now(), err)
	// go alarmDBError(start, expr, args, unknownURI, err)
	if err != nil {
		ecmlog.ErrorEx(ctx, "QueryContext failed", "err", err, "expr", expr, "args", args)
		return nil, err
	}
	return rows, nil
}

// QueryRow 执行查询
func QueryRow(ctx context.Context, db Executor, expr string, args ...interface{}) *sql.Row {
	// start := time.Now()
	row := db.QueryRowContext(ctx, expr, args...)
	err := row.Err()
	// go reportDBCall(start, time.Now(), err)
	// go alarmDBError(start, expr, args, unknownURI, err)
	if err != nil {
		ecmlog.ErrorEx(ctx, "QueryRowContext failed", "err", err, "expr", expr, "args", args)
	}
	return row
}

// Exec 执行 SQL
func Exec(ctx context.Context, db Executor, expr string, args ...interface{}) (sql.Result, error) {
	// start := time.Now()
	result, err := db.ExecContext(ctx, expr, args...)
	// go reportDBCall(start, time.Now(), err)
	// go alarmDBError(start, expr, args, unknownURI, err)
	if err != nil {
		ecmlog.ErrorEx(ctx, "ExecContext failed", "err", err, "expr", expr, "args", args)
	} else {
		ecmlog.InfoEx(ctx, "ExecContext ok", "expr", expr, "args", args)
	}
	return result, err
}

/* Tx Helper */

// TxCallback 事务回调
type TxCallback func(ctx context.Context, tx *sql.Tx) error

// TxWrapper 事务代码的帮助函数
func TxWrapper(ctx context.Context, db *sql.DB, opts *sql.TxOptions, callback TxCallback) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		ecmlog.ErrorEx(ctx, "BeginTx failed", "err", err)
		return err
	}

	if err = callback(ctx, tx); err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			ecmlog.ErrorEx(ctx, "Rollback failed", "err", err2)
		}
		return err
	}

	return tx.Commit()
}

/* Other Helper Functions */

func dereferencedType(t reflect.Type) reflect.Type {
	for k := t.Kind(); k == reflect.Ptr || k == reflect.Slice; k = t.Kind() {
		t = t.Elem()
	}
	return t
}
