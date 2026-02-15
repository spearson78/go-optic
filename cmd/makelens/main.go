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

	packageName := "main_test"
	goFileName := "./types_test.go"
	genFileName := "./types_generated_test.go"

	args := flag.Args()

	if len(args) > 0 {
		packageName = args[0]
		goFileName = args[1]
		genFileName = args[2]

	}

	combinators := []Combinator{
		//Check the history for a full set of combinators.
		//I decided against them in the end
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

	dm.Comparable = make(map[string]struct{}, len(comparables))

	for _, cmp := range comparables {
		dm.Comparable[cmp] = struct{}{}
	}

	fd := DmToFd(&dm, packageName, *rootObjName, combinators)

	for _, v := range imports {
		fd.Imports = append(fd.Imports, v)
	}

	for _, v := range dotImports {
		fd.DotImports = append(fd.DotImports, v)
	}

	codegen.WriteGoFile(w, &fd)
}
