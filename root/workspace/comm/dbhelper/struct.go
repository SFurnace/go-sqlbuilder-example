package dbhelper

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/huandu/go-sqlbuilder"

	ecmlog "pers.drcz/tests/sqlbuilder/comm/log"
)

// Struct 对 sqlbuilder.Struct 进行了封装，使其更易使用
type Struct struct {
	*sqlbuilder.Struct
	typ reflect.Type
}

// NewStruct ...
func NewStruct(val interface{}) *Struct {
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Struct {
		return nil
	}

	return &Struct{Struct: sqlbuilder.NewStruct(val), typ: typ}
}

// ScanRow ...
func (s *Struct) ScanRow(row *sql.Row, destPtr interface{}) error {
	return s.ScanRowForTag(row, "", destPtr)
}

// ScanRowForTag ...
func (s *Struct) ScanRowForTag(row *sql.Row, tag string, destPtr interface{}) error {
	dTyp := reflect.TypeOf(destPtr)
	if dTyp.Kind() != reflect.Ptr || dTyp.Elem() != s.typ {
		return fmt.Errorf("invalid dest type: %v", dTyp)
	}

	if err := row.Scan(s.AddrForTag(tag, destPtr)...); err != nil {
		return err
	}
	return nil
}

// ScanRows ...
func (s *Struct) ScanRows(rows *sql.Rows, destPtr interface{}) error {
	return s.ScanRowsForTag(rows, "", destPtr)
}

// ScanRowsForTag ...
func (s *Struct) ScanRowsForTag(rows *sql.Rows, tag string, destPtr interface{}) error {
	dTyp := reflect.TypeOf(destPtr)
	if dTyp.Kind() != reflect.Ptr || dTyp.Elem().Kind() != reflect.Slice || dTyp.Elem().Elem() != s.typ {
		return fmt.Errorf("invalid dest type: %v", dTyp)
	}

	var (
		dVal = reflect.ValueOf(destPtr).Elem()
		err  error
	)
	for rows.Next() {
		tmp := reflect.New(s.typ)
		if err = rows.Scan(s.AddrForTag(tag, tmp.Interface())...); err != nil {
			return err
		}
		dVal.Set(reflect.Append(dVal, tmp.Elem()))
	}
	if rows.Err() != nil {
		return err
	}
	return nil
}

// Query ...
func (s *Struct) Query(
	ctx context.Context, db Executor, result interface{}, expr string, args ...interface{},
) error {
	return s.TagQuery(ctx, db, result, "", expr, args...)
}

// QueryRow ...
func (s *Struct) QueryRow(
	ctx context.Context, db Executor, result interface{}, expr string, args ...interface{},
) error {
	return s.TagQueryRow(ctx, db, result, "", expr, args...)
}

// TagQuery ...
func (s *Struct) TagQuery(
	ctx context.Context, db Executor, result interface{}, tag, expr string, args ...interface{},
) error {
	rows, err := Query(ctx, db, expr, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err = s.ScanRowsForTag(rows, tag, result); err != nil {
		ecmlog.ErrorEx(ctx, "ScanRowsForTag failed", "err", err, "expr", expr, "args", args)
		return err
	}
	return nil
}

// TagQueryRow ...
func (s *Struct) TagQueryRow(
	ctx context.Context, db Executor, result interface{}, tag, expr string, args ...interface{},
) error {
	if err := s.ScanRowForTag(QueryRow(ctx, db, expr, args...), tag, result); err != nil {
		ecmlog.ErrorEx(ctx, "ScanRowForTag failed", "err", err, "expr", expr, "args", args)
		return err
	}
	return nil
}

// Exec ...
func (s *Struct) Exec(ctx context.Context, db Executor, expr string, args ...interface{}) (sql.Result, error) {
	return Exec(ctx, db, expr, args...)
}
