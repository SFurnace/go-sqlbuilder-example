package pkg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/huandu/go-sqlbuilder"

	"pers.drcz/tests/sqlbuilder/comm/dbhelper"
	ecmlog "pers.drcz/tests/sqlbuilder/comm/log"
)

var db *sql.DB

// SCustomer ORM object to Customer
var SCustomer = dbhelper.NewStruct(Customer{})

// AddCustomer ...
func AddCustomer(ctx context.Context, way, tag string, objects ...interface{}) error {
	var b *sqlbuilder.InsertBuilder
	switch way {
	case dbhelper.Insert:
		b = SCustomer.InsertIntoForTag("t_customer", tag, objects...)
	case dbhelper.InsertIgnore:
		b = SCustomer.InsertIgnoreIntoForTag("t_customer", tag, objects...)
	case dbhelper.Replace:
		b = SCustomer.ReplaceIntoForTag("t_customer", tag, objects...)
	default:
		return fmt.Errorf("invalid insert way: %s", way)
	}

	expr, args := b.Build()
	if _, err := SCustomer.Exec(ctx, db, expr, args...); err != nil {
		ecmlog.ErrorEx(ctx, "Exec failed", "err", err)
		return err
	}
	return nil
}

// DeleteCustomer ...
func DeleteCustomer(ctx context.Context, cond dbhelper.DelCondFunc) error {
	b := SCustomer.DeleteFrom("t_customer")
	cond(b)

	expr, args := b.Build()
	if _, err := SCustomer.Exec(ctx, db, expr, args...); err != nil {
		ecmlog.ErrorEx(ctx, "Exec failed", "err", err)
		return err
	}
	return nil
}

// UpdateCustomer ...
func UpdateCustomer(ctx context.Context, cond dbhelper.UpdateCondFunc) error {
	b := sqlbuilder.NewUpdateBuilder().Update("t_customer")
	cond(b)

	expr, args := b.Build()
	if _, err := SCustomer.Exec(ctx, db, expr, args...); err != nil {
		ecmlog.ErrorEx(ctx, "Exec failed", "err", err)
		return err
	}
	return nil
}

// GetTagCustomer ...
func GetTagCustomer(ctx context.Context, tag string, cond dbhelper.CondFunc) (*Customer, error) {
	b := SCustomer.SelectFromForTag("t_customer", tag)
	cond(b)

	var result Customer
	expr, args := b.Build()
	if err := SCustomer.TagQueryRow(ctx, db, &result, tag, expr, args...); err != nil {
		ecmlog.ErrorEx(ctx, "TagQueryRow failed", "err", err)
		return nil, err
	}
	return &result, nil
}

// GetCustomer ...
func GetCustomer(ctx context.Context, cond dbhelper.CondFunc) (*Customer, error) {
	return GetTagCustomer(ctx, "", cond)
}

// PullTagCustomer ...
func PullTagCustomer(ctx context.Context, tag string, cond dbhelper.CondFunc) ([]Customer, error) {
	b := SCustomer.SelectFromForTag("t_customer", tag)
	cond(b)

	var result []Customer
	expr, args := b.Build()
	if err := SCustomer.TagQuery(ctx, db, &result, tag, expr, args...); err != nil {
		ecmlog.ErrorEx(ctx, "TagQuery failed", "err", err)
		return nil, err
	}
	return result, nil
}

// PullCustomer ...
func PullCustomer(ctx context.Context, cond dbhelper.CondFunc) ([]Customer, error) {
	return PullTagCustomer(ctx, "", cond)
}

// MapNameToCustomer ...
func MapNameToCustomer(customers []Customer, err error) (map[string]*Customer, error) {
	if err != nil {
		return nil, err
	}

	result := make(map[string]*Customer, len(customers))
	for i := range customers {
		result[customers[i].CustomerName] = &customers[i]
	}
	return result, nil
}

// GroupCustomerByName ...
func GroupCustomerByName(objs []Customer, err error) (map[string][]Customer, error) {
	if err != nil {
		return nil, err
	}

	result := make(map[string][]Customer)
	for i := range objs {
		result[objs[i].CustomerName] = append(result[objs[i].CustomerName], objs[i])
	}
	return result, nil
}
