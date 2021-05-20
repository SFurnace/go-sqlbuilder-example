package main

import (
	"bytes"
	"text/template"
)

var tmpl = template.Must(template.ParseFS(embeddedTemplates, "tmpl/v1.tmpl"))

func genV1() {
	defer outFile.Close()

	buf := new(bytes.Buffer)
	err := tmpl.Execute(buf, map[string]interface{}{
		"dbHelperPkg": dbHelperPkg,
		"ecmLogPkg":   ecmLogPkg,
		"structName":  structName,
		"ormName":     ormName,
		"converters":  converterMap,
		"groupers":    grouperMap,
	})
	if err != nil {
		panic(err)
	}

	_, _ = outFile.Write(buf.Bytes())
}
