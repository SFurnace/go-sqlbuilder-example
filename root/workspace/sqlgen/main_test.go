package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGen(t *testing.T) {
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, map[string]interface{}{
		"dbHelperPkg": dbHelperPkg,
		"ecmLogPkg":   ecmLogPkg,
		"structName":  "Device",
		"ormName":     "SDevice",
		"converters":  map[string]string{"InstanceID": "string"},
		"groupers":    map[string]string{"AppID": "int64"},
	})

	fmt.Println(buf)
}
