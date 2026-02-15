package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/exp/oio"
)

func classifyFuncs(rootDir string) {

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
			ComposeLeft(
				functionsInFile(),
				FuncDeclName(),
			),
		),
		Compose(
			Index(TraverseStringP[FncName, rune](), 0),
			Op(unicode.IsUpper),
		),
	)

	functionsInDir := Compose(
		goFiles,
		publicFunctionsInFile,
	)

	classify := T2Of(
		Identity[FncName](),
		Op(func(fncName FncName) []Category {
			var cats []Category
			if strings.Contains(string(fncName), "Is") && !strings.Contains(string(fncName), "Iso") {
				cats = append(cats, "Is")
			}
			if strings.Contains(string(fncName), "As") && !strings.Contains(string(fncName), "Ast") {
				cats = append(cats, "As")
			}
			if strings.Contains(string(fncName), "To") && !strings.Contains(string(fncName), "Top") && !strings.Contains(string(fncName), "Tok") {
				cats = append(cats, "To")
			}
			if strings.Contains(string(fncName), "Of") && !strings.Contains(string(fncName), "Off") {
				cats = append(cats, "Of")
			}
			if strings.Contains(string(fncName), "By") && !strings.Contains(string(fncName), "Byte") {
				cats = append(cats, "By")
			}
			if strings.Contains(string(fncName), "From") {
				cats = append(cats, "From")
			}
			if strings.Contains(string(fncName), "Op") && !strings.Contains(string(fncName), "Ope") && !strings.Contains(string(fncName), "Opt") {
				cats = append(cats, "Op")
			}
			if strings.Contains(string(fncName), "Has") {
				cats = append(cats, "Has")
			}
			if strings.Contains(string(fncName), "With") {
				cats = append(cats, "With")
			}
			if strings.Contains(string(fncName), "Cast") {
				cats = append(cats, "Cast")
			}
			if strings.Contains(string(fncName), "Col") {
				cats = append(cats, "Col")
			}
			if strings.HasSuffix(string(fncName), "T2") {
				cats = append(cats, "T2")
			}
			if strings.Contains(string(fncName), "Lens") {
				cats = append(cats, "Lens")
			}
			if strings.Contains(string(fncName), "Travers") {
				cats = append(cats, "Travers")
			}

			if strings.Contains(string(fncName), "Iter") {
				cats = append(cats, "Iter")
			}

			if strings.Contains(string(fncName), "Getter") {
				cats = append(cats, "Getter")
			}

			if strings.Contains(string(fncName), "Oper") {
				cats = append(cats, "Oper")
			}

			if strings.Contains(string(fncName), "Prism") {
				cats = append(cats, "Prism")
			}

			if strings.Contains(string(fncName), "Apply") {
				cats = append(cats, "Apply")
			}

			if strings.Contains(string(fncName), "Pred") {
				cats = append(cats, "Pred")
			}

			if strings.Contains(string(fncName), "Bind") {
				cats = append(cats, "Bind")
			}
			if strings.Contains(string(fncName), "Ix") {
				cats = append(cats, "Ix")
			}

			if strings.Contains(string(fncName), "Combi") {
				cats = append(cats, "Combi")
			}

			if strings.Contains(string(fncName), "Eq") {
				cats = append(cats, "Eq")
			}

			if strings.Contains(string(fncName), "Min") {
				cats = append(cats, "Min")
			}

			if strings.Contains(string(fncName), "Max") {
				cats = append(cats, "Max")
			}

			if strings.Contains(string(fncName), "On") {
				cats = append(cats, "On")
			}

			if strings.HasSuffix(string(fncName), "IEP") {
				cats = append(cats, "I", "E", "P")
			} else if strings.HasSuffix(string(fncName), "EP") {
				cats = append(cats, "E", "P")
			} else if strings.HasSuffix(string(fncName), "IP") {
				cats = append(cats, "I", "P")
			} else if strings.HasSuffix(string(fncName), "P") {
				cats = append(cats, "P")
			} else if strings.HasSuffix(string(fncName), "IE") {
				cats = append(cats, "I", "E")
			} else if strings.HasSuffix(string(fncName), "E") {
				cats = append(cats, "E")
			} else if strings.HasSuffix(string(fncName), "I") {
				cats = append(cats, "I")
			}

			return cats

		}),
	)

	var filter Optic[int, lo.Tuple2[FncName, []Category], lo.Tuple2[FncName, []Category], bool, bool, ReturnOne, ReadOnly, UniDir, Pure]

	if len(os.Args) == 2 {
		filter = Indexing(True[lo.Tuple2[FncName, []Category]]())
	} else {

		categories := MustGet(
			SliceOf(
				Compose(
					TraverseSlice[string](),
					IsoCast[string, Category](),
				),
				10,
			),
			os.Args[2:],
		)

		filter = Ret1(Ro(Ud(EPure(Compose(
			T2B[FncName, []Category](),
			Any(
				TraverseSlice[Category](),
				In(categories...),
			),
		)))))
	}

	filterTag := Filtered(
		ComposeLeft(
			functionsInDir,
			classify,
		),
		filter,
	)

	seq, err := Get(
		SeqIEOf(filterTag), rootDir)
	if err != nil {
		log.Fatal(err)
	}

	for val := range seq {
		index, focus, err := val.Get()
		if err != nil {
			log.Fatal(index, err)
		}
		fmt.Printf("%v %v -> %v\n", index, focus.A, focus.B)
	}
}
