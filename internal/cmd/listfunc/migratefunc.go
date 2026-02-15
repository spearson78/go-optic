package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"

	. "github.com/spearson78/go-optic"
)

func migrateFuncToFile(funcName string, fromFile string, toFile string) error {

	fset := token.NewFileSet()

	fromNode, err := parser.ParseFile(fset, fromFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	fncDecl, ok, err := GetFirst(
		Filtered(
			ComposeLeft(
				TraverseSlice[ast.Decl](),
				DownCast[ast.Decl, *ast.FuncDecl](),
			),
			Compose(
				FuncDeclName(),
				Eq(FncName(funcName)),
			),
		),
		fromNode.Decls,
	)

	if !ok {
		return fmt.Errorf("func not found %v in file %v", funcName, fromFile)
	}

	fromNode.Decls, err = Modify(
		Identity[[]ast.Decl](),
		FilteredSlice[ast.Decl](
			Compose3(
				DownCast[ast.Decl, *ast.FuncDecl](),
				FuncDeclName(),
				Eq(FncName(funcName)),
			),
		),
		fromNode.Decls,
	)

	if err != nil {
		return err
	}

	toNode, err := parser.ParseFile(fset, toFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	toNode.Decls = append(toNode.Decls, fncDecl)

	comments := MustGet(
		SliceOf(
			Filtered(
				TraverseSlice[*ast.CommentGroup](),
				AndOp(
					Compose(
						CommentGroupPos(),
						Gte(fncDecl.Pos()),
					),
					Compose(
						CommentGroupEnd(),
						Lte(fncDecl.End()),
					),
				),
			),
			2,
		),
		fromNode.Comments,
	)

	fromNode.Comments = MustModify(
		Identity[[]*ast.CommentGroup](),
		FilteredSlice[*ast.CommentGroup](
			AndOp(
				Compose(
					CommentGroupPos(),
					Gte(fncDecl.Pos()),
				),
				Compose(
					CommentGroupEnd(),
					Lte(fncDecl.End()),
				),
			),
		),
		fromNode.Comments,
	)

	toNode.Comments = append(toNode.Comments, comments...)

	cfg := printer.Config{
		Mode:     printer.TabIndent | printer.UseSpaces,
		Tabwidth: 4,
	}

	// write changed AST to file
	toF, err := os.Create(toFile)
	if err != nil {
		return err
	}
	defer toF.Close()
	if err := cfg.Fprint(toF, fset, toNode); err != nil {
		return err
	}

	fromF, err := os.Create(fromFile)
	if err != nil {
		return err
	}
	defer fromF.Close()
	if err := cfg.Fprint(fromF, fset, fromNode); err != nil {
		return err
	}

	return nil

}
