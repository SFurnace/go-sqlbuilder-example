package dbhelper

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/huandu/go-sqlbuilder"
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

	if err := row.Scan(s.AddrForTag(tag, destPtr)); err != nil {
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
