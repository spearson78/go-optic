package main

import (
	"fmt"
	"go/ast"
	"log"
	"strings"
	"unicode"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/exp/oio"
)

func checkDocs(rootDir string) {

	goFiles := goFilesInDir(Op(func(path string) bool {
		return !strings.Contains(path, "_test") &&
			!strings.HasPrefix(path, "/home/spearson/projects/optic3/exp/") &&
			!strings.HasPrefix(path, "/home/spearson/projects/optic3/expr/") &&
			!strings.HasPrefix(path, "/home/spearson/projects/optic3/internal/") &&
			!strings.HasPrefix(path, "/home/spearson/projects/optic3/examples/") &&
			!strings.HasPrefix(path, "/home/spearson/projects/optic3/optics-by-example/") &&
			!strings.HasPrefix(path, "/home/spearson/projects/optic3/cmd/")

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

	undocced := ComposeLeft(
		Filtered(
			Compose(
				goFiles,
				publicFunctionsInFile,
			),
			Compose(
				FuncDeclDoc(),
				Eq[*ast.CommentGroup](nil),
			),
		),
		FuncDeclName(),
	)

	seq, err := Get(
		SeqIEOf(
			undocced,
		),
		rootDir,
	)

	if err != nil {
		log.Fatal(err)
	}

	for val := range seq {
		index, focus, err := val.Get()
		if err != nil {
			log.Fatal(index, err)
		}

		fmt.Printf("%v %v Missing docs\n", index, focus)
	}

}
