package main

import (
	"bytes"
	"path/filepath"
	"text/template"
)

var tmpl = template.Must(template.ParseFS(embeddedTemplates, "tmpl/v1.tmpl"))

func genV1() {
	defer outFile.Close()

	buf := new(bytes.Buffer)
	err := tmpl.Execute(buf, map[string]interface{}{
		"dbHelperPkg": dbHelperPkg,
		"ecmLogPkg":   ecmLogPkg,
		"outPkg":      outPkg,
		"extFile":     filepath.Base(extFilePath),
		"table":       tableStr,
		"db":          *dbVar,
		"structName":  structName,
		"fullName":    structFullName,
		"ormName":     ormName,
		"converters":  converterMap,
		"groupers":    grouperMap,
	})
	if err != nil {
		panic(err)
	}

	_, _ = outFile.Write(buf.Bytes())
}
