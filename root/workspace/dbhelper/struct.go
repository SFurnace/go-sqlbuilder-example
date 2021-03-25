package dbhelper

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/huandu/go-sqlbuilder"

	"pers.drcz/test/regex/log"
)

type Struct struct {
	*sqlbuilder.Struct
	typ reflect.Type
}

func NewStruct(val interface{}) *Struct {
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Struct {
		return nil
	}

	return &Struct{Struct: sqlbuilder.NewStruct(val), typ: typ}
}

func (s *Struct) ScanRow(row *sql.Row, destPtr interface{}) error {
	return s.ScanRowForTag(row, "", destPtr)
}

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

func (s *Struct) ScanRows(rows *sql.Rows, destPtr interface{}) error {
	return s.ScanRowsForTag(rows, "", destPtr)
}

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

func (s *Struct) Query(ctx context.Context, db *sql.DB, result interface{}, expr string, args ...interface{}) error {
	return s.TagQuery(ctx, db, result, "", expr, args...)
}

func (s *Struct) QueryRow(ctx context.Context, db *sql.DB, result interface{}, expr string, args ...interface{}) error {
	return s.TagQueryRow(ctx, db, result, "", expr, args...)
}

func (s *Struct) TagQuery(ctx context.Context, db *sql.DB, result interface{}, tag, expr string, args ...interface{}) error {
	rows, err := db.QueryContext(ctx, expr, args...)
	if err != nil {
		log.ErrorEx(ctx, "QueryContext failed", "err", err, "expr", expr, "args", args)
		return err
	}
	defer rows.Close()

	if err = s.ScanRowsForTag(rows, tag, result); err != nil {
		log.ErrorEx(ctx, "ScanRowsForTag failed", "err", err, "expr", expr, "args", args)
		return err
	}
	return nil
}

func (s *Struct) TagQueryRow(ctx context.Context, db *sql.DB, result interface{}, tag, expr string, args ...interface{}) error {
	row := db.QueryRowContext(ctx, expr, args...)
	if err := s.ScanRowForTag(row, tag, result); err != nil {
		log.ErrorEx(ctx, "ScanRowForTag failed", "err", err, "expr", expr, "args", args)
		return err
	}
	return nil
}

func (s *Struct) QueryTx(ctx context.Context, tx *sql.Tx, result interface{}, expr string, args ...interface{}) error {
	return s.TagQueryTx(ctx, tx, result, "", expr, args...)
}

func (s *Struct) QueryRowTx(ctx context.Context, tx *sql.Tx, result interface{}, expr string, args ...interface{}) error {
	return s.TagQueryRowTx(ctx, tx, result, "", expr, args...)
}

func (s *Struct) TagQueryTx(ctx context.Context, tx *sql.Tx, result interface{}, tag, expr string, args ...interface{}) error {
	rows, err := tx.QueryContext(ctx, expr, args...)
	if err != nil {
		log.ErrorEx(ctx, "QueryContext failed", "err", err, "expr", expr, "args", args)
		return err
	}
	defer rows.Close()

	if err = s.ScanRowsForTag(rows, tag, result); err != nil {
		log.ErrorEx(ctx, "ScanRowsForTag failed", "err", err, "expr", expr, "args", args)
		return err
	}
	return nil
}

func (s *Struct) TagQueryRowTx(ctx context.Context, tx *sql.Tx, result interface{}, tag, expr string, args ...interface{}) error {
	row := tx.QueryRowContext(ctx, expr, args...)
	if err := s.ScanRowForTag(row, tag, result); err != nil {
		log.ErrorEx(ctx, "ScanRowForTag failed", "err", err)
		return err
	}
	return nil
}

func (s *Struct) Exec(ctx context.Context, db *sql.DB, expr string, args ...interface{}) (sql.Result, error) {
	result, err := db.ExecContext(ctx, expr, args...)
	if err != nil {
		log.ErrorEx(ctx, "ExecContext failed", "err", err, "expr", expr, "args", args)
	}
	return result, err
}

func (s *Struct) ExecTx(ctx context.Context, tx *sql.Tx, expr string, args ...interface{}) (sql.Result, error) {
	result, err := tx.ExecContext(ctx, expr, args...)
	if err != nil {
		log.ErrorEx(ctx, "ExecContext failed", "err", err, "expr", expr, "args", args)
	}
	return result, err
}
