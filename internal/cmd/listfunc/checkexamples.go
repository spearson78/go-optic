package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"strings"
	"unicode"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/exp/oio"
)

func checkExamples(rootDir string) {

	goFiles := goFilesInDir(Op(func(path string) bool {
		return !strings.Contains(path, "_test") &&
			!strings.HasPrefix(path, "/home/spearson/projects/optic3/exp/") &&
			!strings.HasPrefix(path, "/home/spearson/projects/optic3/expr/") &&
			!strings.HasPrefix(path, "/home/spearson/projects/optic3/internal/") &&
			!strings.HasPrefix(path, "/home/spearson/projects/optic3/examples/") &&
			!strings.HasPrefix(path, "/home/spearson/projects/optic3/optics-by-example/") &&
			!strings.HasPrefix(path, "/home/spearson/projects/optic3/cmd/")

	}))

	testFiles := goFilesInDir(Op(func(path string) bool {
		return strings.Contains(path, "_test") &&
			!strings.HasPrefix(path, "/home/spearson/projects/optic3/exp/") &&
			!strings.HasPrefix(path, "/home/spearson/projects/optic3/expr/") &&
			!strings.HasPrefix(path, "/home/spearson/projects/optic3/internal/") &&
			!strings.HasPrefix(path, "/home/spearson/projects/optic3/examples/") &&
			!strings.HasPrefix(path, "/home/spearson/projects/optic3/optics-by-example/") &&
			!strings.HasPrefix(path, "/home/spearson/projects/optic3/cmd/")
	}))

	publicFunctionsInFile := ComposeLeft(
		Filtered(
			Compose(
				oio.FileInfoFullPath(),
				functionsInFile(),
			),
			AndOp(
				Compose(
					FuncDeclRecv(),
					Eq[*ast.FieldList](nil),
				),
				Compose3(
					FuncDeclName(),
					Index(TraverseStringP[FncName, rune](), 0),
					Op(unicode.IsUpper),
				),
			),
		),
		FuncDeclName(),
	)

	expectedExamples := RetM(Ro(Ud(EErr(ReIndexed(
		ComposeLeft(
			Compose(
				goFiles,
				publicFunctionsInFile,
			),
			Compose(
				IsoCast[FncName, string](),
				PrependString(StringCol("Example")),
			),
		),
		AsModify(
			Compose(
				PositionFileName(),
				FileNameWithoutExt(),
			),
			AppendString(StringCol("_test")),
		),
		EqT2[token.Position](),
	)))))

	actualExamples := RetM(Ro(Ud(EErr(
		ComposeLeft(
			Compose(
				testFiles,
				Filtered(
					publicFunctionsInFile,
					Compose(
						IsoCast[FncName, string](),
						StringHasPrefix("Example"),
					),
				),
			),
			IsoCast[FncName, string](),
		),
	))))

	diff := DiffColT2I(
		3,
		EErr(DistanceI[token.Position, string](
			func(ia token.Position, a string, ib token.Position, b string) float64 {
				dist := 0
				if ia.Filename != ib.Filename {
					dist++
				}
				if a != b {
					dist += 1000
				}

				return float64(dist)
			},
		)),
		OpT2(func(a token.Position, b token.Position) bool {
			return a.Filename == b.Filename
		}),
		DiffNone,
		false,
	)

	collections := RetM(Ro(Ud(EErr(T2Of(
		ColOf(actualExamples),
		ColOf(expectedExamples),
	)))))

	diffSeq, err := Get(

		SeqIEOf(
			Compose(
				collections,
				diff,
			),
		),
		rootDir,
	)

	if err != nil {
		log.Fatal(err)
	}

	for val := range diffSeq {
		index, focus, err := val.Get()
		if err != nil {
			log.Fatal(index, err)
		}

		fncName, _ := focus.Get()
		if fncName == "" {
			fncName = index.BeforeValue
		}

		if index.Type == DiffModify {
			/*
				err = migrateFuncToFile(fncName, index.AfterIndex.Filename, index.BeforeIndex.Filename)
				if err != nil {
					log.Fatal(err)
					return false
				}
			*/
			continue
		}

		fmt.Printf("Expected: %v Current: %v Func: %v Type: %v\n", index.BeforeIndex.Filename, index.AfterIndex.Filename, fncName, index.Type)
	}

}
