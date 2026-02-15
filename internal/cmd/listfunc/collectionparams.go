package main

import (
	"fmt"
	"log"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/exp/oio"
)

func collectionAsParams(rootDir string) {

	goFiles := goFilesInDir(Const[string](true))

	allFunctionsInFile := Compose(
		oio.FileInfoFullPath(),
		functionsInFile(),
	)

	functionsInDir := Compose(
		goFiles,
		allFunctionsInFile,
	)

	opticParams := Filtered(
		functionsInDir,

		Any(
			FncDeclParams(),
			AndOp(
				Compose(
					FieldType(),
					Eq("Optic"),
				),
				Any(
					Compose(
						FieldTypeParams(),
						ExprName(),
					),
					Eq("Collection"),
				),
			),
		),
	)

	res := opticParams

	//FncDeclParams(),
	//FieldTypeParams(),
	//DownCast[ast.Expr, *ast.IndexListExpr](),
	//TraverseIndexList(),

	seq, err := Get(
		SeqIEOf(res), rootDir)
	if err != nil {
		log.Fatal(err)
	}

	for val := range seq {
		index, focus, err := val.Get()
		if err != nil {
			log.Fatal(index, err)
		}
		fmt.Printf("%v -> %v %T\n", index, focus, focus)
	}
}
