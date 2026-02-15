package main

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/exp/oio"
	"github.com/spearson78/go-optic/otree"
)

type FncName string
type Category string

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory>")
		os.Exit(1)
	}

	rootDir := os.Args[1]

	//rootDir := "../../../"

	goscriptExports(rootDir)
	//classifyFuncs(rootDir)
	//checkExamples(rootDir)
	//checkDocs(rootDir)
}

func PositionFileName() Optic[Void, token.Position, token.Position, string, string, ReturnOne, ReadWrite, UniDir, Pure] {
	return FieldLens(func(source *token.Position) *string {
		return &source.Filename
	})
}

func CommentGroupPos() Optic[Void, *ast.CommentGroup, *ast.CommentGroup, token.Pos, token.Pos, ReturnOne, ReadOnly, UniDir, Pure] {
	return MethodGetter((*ast.CommentGroup).Pos)
}

func CommentGroupEnd() Optic[Void, *ast.CommentGroup, *ast.CommentGroup, token.Pos, token.Pos, ReturnOne, ReadOnly, UniDir, Pure] {
	return MethodGetter((*ast.CommentGroup).End)
}

func PathBase() Optic[Void, string, string, string, string, ReturnOne, ReadWrite, UniDir, Pure] {
	return Lens[string, string](
		func(source string) string {
			return path.Base(source)
		},
		func(focus, source string) string {
			sourceBase := path.Base(source)
			prefix := source[:len(source)-len(sourceBase)]
			return prefix + focus
		},
		ExprCustom("PathBase"),
	)
}

func PathExt() Optic[Void, string, string, string, string, ReturnOne, ReadWrite, UniDir, Pure] {
	return Lens[string, string](
		func(source string) string {
			return path.Ext(source)
		},
		func(focus, source string) string {
			sourceBase := path.Ext(source)
			prefix := source[:len(source)-len(sourceBase)]
			return prefix + focus
		},
		ExprCustom("PathExt"),
	)
}

func FileNameWithoutExt() Optic[Void, string, string, string, string, ReturnOne, ReadWrite, UniDir, Pure] {
	return Lens[string, string](
		func(source string) string {
			ext := path.Ext(source)
			return source[:len(source)-len(ext)]
		},
		func(focus, source string) string {
			ext := path.Ext(source)
			return focus + ext
		},
		ExprCustom("PathExt"),
	)
}

func goFilesInDir(filter Predicate[string, Pure]) Optic[*otree.PathNode[string], string, string, oio.FileInfo, oio.FileInfo, ReturnMany, ReadOnly, UniDir, Err] {
	return RetM(Ro(Ud(EErr(Compose(
		oio.Stat(),
		Filtered(
			otree.TopDown(oio.TraverseFileInfo()),
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

func FncDeclParams() Optic[int, *ast.FuncDecl, *ast.FuncDecl, *ast.Field, *ast.Field, ReturnMany, ReadOnly, UniDir, Err] {
	return RetM(Ro(Ud(EErr(Compose(
		PtrFieldLensE(
			func(source *ast.FuncDecl) *[]*ast.Field {
				return &source.Type.Params.List
			},
		),
		TraverseSlice[*ast.Field](),
	)))))
}

func ExprName() Optic[Void, ast.Expr, ast.Expr, string, string, ReturnMany, ReadOnly, UniDir, Err] {
	return RetM(Ro(Ud(EErr(Coalesce(
		Compose4(
			DownCast[ast.Expr, *ast.IndexListExpr](),
			PtrFieldLensE(
				func(source *ast.IndexListExpr) *ast.Expr {
					return &source.X
				},
			),
			DownCast[ast.Expr, *ast.Ident](),
			IdentName(),
		),
		Compose(
			DownCast[ast.Expr, *ast.Ident](),
			IdentName(),
		),
	)))))
}

func FieldType() Optic[Void, *ast.Field, *ast.Field, string, string, ReturnMany, ReadOnly, UniDir, Err] {

	return RetM(Ro(Ud(EErr(Compose(
		PtrFieldLensE(
			func(source *ast.Field) *ast.Expr {
				return &source.Type
			},
		),
		ExprName(),
	)))))

}

func TraverseIndexList() Optic[int, *ast.IndexListExpr, *ast.IndexListExpr, ast.Expr, ast.Expr, ReturnMany, ReadOnly, UniDir, Err] {
	return RetM(Ro(Ud(EErr(Compose(
		PtrFieldLensE(
			func(source *ast.IndexListExpr) *[]ast.Expr {
				return &source.Indices
			},
		),
		TraverseSlice[ast.Expr](),
	)))))
}

func FieldTypeParams() Optic[int, *ast.Field, *ast.Field, ast.Expr, ast.Expr, ReturnMany, ReadOnly, UniDir, Err] {
	return RetM(Ro(Ud(EErr(
		Compose3(
			PtrFieldLensE(
				func(source *ast.Field) *ast.Expr {
					return &source.Type
				},
			),
			DownCast[ast.Expr, *ast.IndexListExpr](),
			TraverseIndexList(),
		),
	))))

}

func IdentName() Optic[Void, *ast.Ident, *ast.Ident, string, string, ReturnOne, ReadWrite, UniDir, Err] {
	return PtrFieldLensE(
		func(source *ast.Ident) *string {
			return &source.Name
		},
	)
}

func FieldNames() Optic[Void, *ast.Field, *ast.Field, string, string, ReturnMany, ReadOnly, UniDir, Err] {
	return RetM(Ro(Ud(EErr(
		Compose3(
			PtrFieldLensE(
				func(source *ast.Field) *[]*ast.Ident {
					return &source.Names
				},
			),
			TraverseSlice[*ast.Ident](),
			IdentName(),
		),
	))))
}

func FieldTypes() Optic[Void, *ast.Field, *ast.Field, ast.Expr, ast.Expr, ReturnOne, ReadOnly, UniDir, Err] {
	return Ro(PtrFieldLensE(
		func(source *ast.Field) *ast.Expr {
			return &source.Type
		},
	))
}

func FuncDeclName() Optic[Void, *ast.FuncDecl, *ast.FuncDecl, FncName, FncName, ReturnOne, ReadOnly, UniDir, Err] {
	return EErr(Ret1(Ro(
		Compose(
			PtrFieldLensE(
				func(source *ast.FuncDecl) **ast.Ident {
					return &source.Name
				},
			),
			PtrFieldLensE(
				func(source *ast.Ident) *FncName {
					return (*FncName)(&source.Name)
				},
			),
		))))
}

func FuncDeclRecv() Optic[Void, *ast.FuncDecl, *ast.FuncDecl, *ast.FieldList, *ast.FieldList, ReturnMany, ReadOnly, UniDir, Pure] {
	return Ro(PtrFieldLens(
		func(source *ast.FuncDecl) **ast.FieldList {
			return &source.Recv
		},
	))
}

func FuncDeclDoc() Optic[Void, *ast.FuncDecl, *ast.FuncDecl, *ast.CommentGroup, *ast.CommentGroup, ReturnOne, ReadOnly, UniDir, Err] {
	return Ro(PtrFieldLensE(
		func(source *ast.FuncDecl) **ast.CommentGroup {
			return (&source.Doc)
		},
	))
}

func functionsInFile() Optic[token.Position, string, string, *ast.FuncDecl, *ast.FuncDecl, ReturnMany, ReadOnly, UniDir, Err] {
	return IterationIE[token.Position, string, *ast.FuncDecl](
		func(ctx context.Context, source string) SeqIE[token.Position, *ast.FuncDecl] {
			return func(yield func(val ValueIE[token.Position, *ast.FuncDecl]) bool) {

				fset := token.NewFileSet() // FileSet is needed to parse Go files
				node, err := parser.ParseFile(fset, source, nil, parser.ParseComments)
				if err != nil {
					yield(ValIE(token.Position{
						Filename: source,
					}, (*ast.FuncDecl)(nil), err))
					return
				}

				for _, n := range node.Decls {
					if funcDecl, ok := n.(*ast.FuncDecl); ok {
						p := fset.Position(n.Pos())
						if !yield(ValIE(p, funcDecl, nil)) {
							return
						}
					}
				}
			}
		},
		nil,
		nil,
		IxMatchComparable[token.Position](),
		ExprCustom("listFunctionsInFile"),
	)
}
