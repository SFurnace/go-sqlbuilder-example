package tests

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// Tables
const (
	CustomerTable = "t_customer"
	NodeTable     = "t_node"
	DeviceTable   = "t_device"
)

var (
	DB *sql.DB
)

func init() {
	var err error
	const DSN = "tester:tester123@tcp(localhost:3306)/test"

	if DB, err = sql.Open("mysql", DSN); err != nil {
		panic(err)
	}
}
