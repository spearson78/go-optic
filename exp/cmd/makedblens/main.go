package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"

	"github.com/spearson78/go-optic/internal/codegen"
)

func main() {
	parseArgs()
	astMain()
}

func astMain() {

	packageName := "main"
	goFileName := "../../examples/recipes/model.go"
	genFileName := "../../examples/recipes/model_generated.go"

	args := flag.Args()

	if len(args) > 0 {
		packageName = args[0]
		goFileName = args[1]
		genFileName = args[2]

	}

	fset := token.NewFileSet()

	var files []*ast.File

	if goFileName == "." {

		goFiles, err := os.ReadDir(goFileName)
		if err != nil {
			log.Fatal(err)
		}
		for _, goFile := range goFiles {
			if goFile.Name() != genFileName && !goFile.IsDir() && strings.HasSuffix(goFile.Name(), ".go") && !strings.HasSuffix(goFile.Name(), "_test.go") {
				node, err := parser.ParseFile(fset, goFile.Name(), nil, parser.ParseComments)
				if err != nil {
					log.Fatal(err)
				}
				files = append(files, node)
			}
		}
	} else {
		node, err := parser.ParseFile(fset, goFileName, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, node)
	}

	w, err := os.Create(genFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	dm := Parse(files)

	fd := DmToFd(&dm, packageName, *rootObjName)

	if packageName != "osql" {
		fd.Imports = append(fd.Imports, "github.com/spearson78/go-optic/exp/osql")
	}

	for _, v := range imports {
		fd.Imports = append(fd.Imports, v)
	}

	codegen.WriteGoFile(w, &fd)
}
