package main

import (
	"embed"
	"fmt"
	"os"
)

const (
	dbHelperPkg = `"git.code.oa.com/IOT_EC/ecm_websvr_proj/src/comm/dbhelper"`
	ecmLogPkg   = `"git.code.oa.com/IOT_EC/ecm_websvr_proj/src/comm/ecmlog"`
)

var (
	validMemberType = []string{
		"string", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "bool",
		"float32", "float64",
	}

	//go:embed tmpl/*
	embeddedTemplates embed.FS
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
