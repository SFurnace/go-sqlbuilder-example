package main

import (
	"flag"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

var (
	// parameters
	version      = flag.Int("v", 1, "version")
	name         = flag.String("t", "", "result struct type name")
	ormPrefix    = flag.String("p", "", "prefix of orm object name")
	table        = flag.String("tn", "", "table name")
	tableVar     = flag.String("tv", "", "variable contains the table name")
	outputFile   = flag.String("out", "stdin", "output file path")
	outputPkg    = flag.String("pkg", "", "output package")
	converterStr = flag.String("conv", "", "generate result map converters, format like: type:member;type:member...")

	// calculated
	structName, ormName, tableStr, outPkg string
	outFile                               *os.File
	converterMap                          map[string]string
)

func main() {
	checkParam()

	switch *version {
	case 1:
		genV1()
	case 2:
		// TODO
	}
}

func checkParam() {
	flag.Parse()
	checkVersion()
	checkName()
	checkOrmName()
	checkTableStr()
	checkOutput()
	checkConverters()
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
	case len(ss) == 2 && token.IsIdentifier(ss[0]), token.IsIdentifier(ss[1]):
	default:
		failedExit("invalid struct type name: %s", *name)
	}

	structName = *name
}

func checkOrmName() {
	ss := strings.Split(*name, ".")
	ormName = *ormPrefix + ss[len(ss)-1]
}

func checkTableStr() {
	if (*table != "" && *tableVar != "") || (*table == "" && *tableVar == "") {
		failedExit("invalid table name")
	}
	tableStr = *table + *tableVar
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

		p, _ := filepath.Abs(rel)
		outPkg = filepath.Base(filepath.Dir(p))
	}

	if *outputPkg != "" {
		outPkg = *outputPkg
	}
}

func checkConverters() {
	for _, pair := range strings.Split(*converterStr, ";") {
		ss := strings.Split(pair, ":")
		if len(ss) != 2 {
			failedExit("invalid converter: %s", pair)
		}
		if !isValidMemberType(ss[0]) {
			failedExit("invalid struct member type: %s", ss[0])
		}
		if !token.IsIdentifier(ss[1]) {
			failedExit("invalid member name: %s", ss[1])
		}
		if _, seen := converterMap[ss[1]]; seen {
			failedExit("duplicated converter: %s", ss[1])
		}

		converterMap[ss[1]] = ss[0]
	}
}
