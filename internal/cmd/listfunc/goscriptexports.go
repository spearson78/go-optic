package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"unicode"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/exp/oio"
)

func goFilesInSingleDir(filter Predicate[string, Pure]) Optic[string, string, string, oio.FileInfo, oio.FileInfo, ReturnMany, ReadOnly, UniDir, Err] {
	return RetM(Ro(Ud(EErr(Compose(
		oio.Stat(),
		Filtered(
			oio.TraverseFileInfo(),
			Compose(
				oio.FileInfoFullPath(),
				AndOp(
					Op(func(path string) bool {
						return filepath.Ext(path) == ".go"
					}),
					PredToOptic(filter),
				),
			),
		),
	)))))
}

func goscriptExports(rootDir string) {

	goFiles := goFilesInSingleDir(Op(func(path string) bool {
		return !strings.Contains(path, "_test")
	}))

	publicFunctionsInFile := Filtered(
		Compose(
			oio.FileInfoFullPath(),
			functionsInFile(),
		),
		Compose3(
			FuncDeclName(),
			Index(TraverseStringP[FncName, rune](), 0),
			Op(unicode.IsUpper),
		),
	)

	functionsInDir := ComposeLeft(
		goFiles,
		publicFunctionsInFile,
	)

	seq, err := Get(
		SeqIEOf(functionsInDir), rootDir)
	if err != nil {
		log.Fatal(err)
	}

	var lastIndex string
	for val := range seq {
		index, focus, err := val.Get()
		if err != nil {
			log.Fatal(index, err)
		}

		if index != lastIndex {
			fmt.Println("//" + index)
		}

		fmt.Printf("\"%v\":		reflect.ValueOf(%v", focus.Name.Name, focus.Name.Name)
		if focus.Type.TypeParams != nil {
			fmt.Print("[")
			first := true
			for _, v := range focus.Type.TypeParams.List {
				for range v.Names {
					if !first {
						fmt.Print(", ")
					}
					first = false
					fmt.Print("any")
				}
			}
			fmt.Print("]")
		}
		fmt.Println("),")
		lastIndex = index
	}
}
