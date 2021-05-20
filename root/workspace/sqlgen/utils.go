package main

import (
	"fmt"
	"os"
)

const (
	dbHelperPkg = `"pers.drcz/tests/sqlbuilder/comm/dbhelper"`
	ecmLogPkg   = `ecmlog "pers.drcz/tests/sqlbuilder/comm/log"`
)

var (
	validMemberType = []string{
		"string", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "bool",
		"float32", "float64",
	}
)

func failedExit(reason string, v ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, reason+"\n", v...)
	os.Exit(1)
}

func isValidMemberType(v string) bool {
	for _, s := range validMemberType {
		if v == s {
			return true
		}
	}
	return false
}
