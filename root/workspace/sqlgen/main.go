package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	// parameters
	version    = flag.Int("v", 1, "version")
	name       = flag.String("t", "", "result struct type name")
	ormPrefix  = flag.String("p", "", "prefix of orm object name")
	table      = flag.String("tn", "", "table name")
	tableVar   = flag.String("tv", "", "variable contains the table name")
	dbVar      = flag.String("db", "db", "variable of db object")
	outputFile = flag.String("out", "stdin", "output file path")
	outputPkg  = flag.String("pkg", "", "output package")
	mapStr     = flag.String("map", "", "generate result mappers, format like: member:type;member:type...")
	grouperStr = flag.String("group", "", "generate result groupers, format like: member:type;member:type...")

	// calculated
	structName, structFullName, ormName, ormStruct, tableStr string
	outPkg, extFilePath                                      string
	outFile                                                  *os.File
	converterMap, grouperMap                                 = make(map[string]string), make(map[string]string)
)

func main() {
	checkParam()

	switch *version {
	case 1:
		generateContent("tmpl/v1.tmpl")
	case 2:
		generateContent("tmpl/v2.tmpl")
	}

	var (
		extFile *os.File
		err     error
	)
	if _, err = os.Stat(extFilePath); os.IsNotExist(err) {
		if extFile, err = os.Create(extFilePath); err == nil {
			_, _ = extFile.WriteString(fmt.Sprintf("package %s\n", outPkg))
		}
	}
}

/* generator */

func generateContent(file string) {
	defer outFile.Close()

	buf := new(bytes.Buffer)
	tmpl := template.Must(template.ParseFS(embeddedTemplates, file))
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
		"ormStruct":   ormStruct,
		"converters":  converterMap,
		"groupers":    grouperMap,
	})
	if err != nil {
		panic(err)
	}

	_, _ = outFile.Write(buf.Bytes())
}

/* check helper */

func checkParam() {
	flag.Parse()
	checkVersion()
	checkName()
	checkOrmName()
	checkTableStr()
	checkOutput()
	checkMappers()
	checkGroupers()
}

func checkVersion() {
	switch *version {
	case 1, 2:
	default:
		failedExit("unknown version: %d", *version)
	}
}

func checkName() {
	ss := strings.Split(*name, ".")
	switch {
	case len(ss) == 1 && token.IsIdentifier(*name):
		structName, structFullName = *name, *name
	case len(ss) == 2 && token.IsIdentifier(ss[0]), token.IsIdentifier(ss[1]):
		structName, structFullName = ss[1], *name
	default:
		failedExit("invalid struct type name: %s", *name)
	}
}

func checkOrmName() {
	ss := strings.Split(*name, ".")
	origin := ss[len(ss)-1]
	ormName = *ormPrefix + origin
	ormStruct = strings.ToLower(origin[0:1]) + origin[1:] + "Model" // 首字母小写
}

func checkTableStr() {
	if *table == "" && *tableVar == "" {
		failedExit("empty table name")
	}
	switch {
	case token.IsIdentifier(*tableVar):
		tableStr = *tableVar
	case *table != "":
		tableStr = fmt.Sprintf(`"%s"`, *table)
	default:
		failedExit("invalid table name")
	}
}

func checkOutput() {
	switch *outputFile {
	case "stdin":
		outFile, outPkg = os.Stdout, "unknown"
	default:
		rel, err := filepath.Rel(".", *outputFile)
		if err != nil {
			failedExit("invalid output file path")
		}

		outFile, err = os.Create(rel)
		if err != nil {
			failedExit("can't create output file")
		}

		d, f := filepath.Split(rel)
		e := filepath.Ext(f)
		extFilePath = filepath.Join(d, fmt.Sprintf("%s_ext%s", strings.TrimSuffix(f, e), e))
		p, _ := filepath.Abs(rel)
		outPkg = filepath.Base(filepath.Dir(p))
	}

	if *outputPkg != "" {
		outPkg = *outputPkg
	}
}

func checkMappers() {
	if *mapStr != "" {
		for _, pair := range strings.Split(*mapStr, ";") {
			mem, typ := checkTypeToMemberStr(pair)
			if _, seen := converterMap[mem]; seen {
				failedExit("duplicated mapper: %s", mem)
			}

			converterMap[mem] = typ
		}
	}
}

func checkGroupers() {
	if *grouperStr != "" {
		for _, pair := range strings.Split(*grouperStr, ";") {
			mem, typ := checkTypeToMemberStr(pair)
			if _, seen := grouperMap[mem]; seen {
				failedExit("duplicated grouper: %s", mem)
			}

			grouperMap[mem] = typ
		}
	}
}

func checkTypeToMemberStr(str string) (string, string) {
	ss := strings.Split(str, ":")
	if len(ss) != 2 {
		failedExit("invalid pair: %s", str)
	}
	if !isValidMemberType(ss[1]) {
		failedExit("invalid struct member type: %s", ss[1])
	}
	if !token.IsIdentifier(ss[0]) {
		failedExit("invalid member name: %s", ss[0])
	}
	return ss[0], ss[1]
}
