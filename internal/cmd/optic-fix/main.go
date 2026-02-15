package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"strings"
)

func main() {

	goFileName := "/home/spearson/projects/optic3/plated/actions_test.go"

	fset := token.NewFileSet()

	files := make(map[string]*ast.File)

	if goFileName == "." {

		goFiles, err := os.ReadDir(goFileName)
		if err != nil {
			log.Fatal(err)
		}
		for _, goFile := range goFiles {
			if !goFile.IsDir() && strings.HasSuffix(goFile.Name(), ".go") {
				node, err := parser.ParseFile(fset, goFile.Name(), nil, parser.ParseComments)
				if err != nil {
					log.Fatal(err)
				}
				files[goFile.Name()] = node
			}
		}
	} else {
		node, err := parser.ParseFile(fset, goFileName, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}
		files[goFileName] = node
	}

	for name, ast := range files {
		err := updateFile(fset, ast)
		if err != nil {
			log.Fatal(err)
		}

		w, err := os.Create(name)
		if err != nil {
			log.Fatal(err)
		}

		cfg := printer.Config{
			Mode:     printer.TabIndent | printer.UseSpaces,
			Tabwidth: 4,
		}

		err = cfg.Fprint(w, fset, ast)
		if err != nil {
			log.Fatal(err)
		}

		w.Close()
	}

}

func updateFile(fset *token.FileSet, f *ast.File) error {

	for _, decl := range f.Decls {
		err := updateDecl(fset, decl)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateFuncDecl(fset *token.FileSet, f *ast.FuncDecl) error {
	return updateBlockStmt(fset, f.Body)
}

func updateBlockStmt(fset *token.FileSet, block *ast.BlockStmt) error {
	for _, stmnt := range block.List {
		err := updateStmt(fset, stmnt)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateStmt(fset *token.FileSet, stmnt ast.Stmt) error {

	if stmnt == nil {
		return nil
	}

	switch t := stmnt.(type) {
	case *ast.ReturnStmt:
		for _, x := range t.Results {
			err := updateExpr(fset, x)
			if err != nil {
				return err
			}
		}
		return nil
	case *ast.ExprStmt:
		return updateExpr(fset, t.X)
	case *ast.AssignStmt:
		for _, x := range t.Rhs {
			err := updateExpr(fset, x)
			if err != nil {
				return err
			}
		}
		return nil
	case *ast.DeclStmt:
		return updateDecl(fset, t.Decl)
	case *ast.RangeStmt:
		return updateExpr(fset, t.X)
	case *ast.IfStmt:
		err := updateStmt(fset, t.Init)
		if err != nil {
			return err
		}

		err = updateExpr(fset, t.Cond)
		if err != nil {
			return err
		}

		err = updateStmt(fset, t.Else)
		if err != nil {
			return err
		}

		return updateBlockStmt(fset, t.Body)
	case *ast.DeferStmt:
		return updateExpr(fset, t.Call)
	case *ast.TypeSwitchStmt:
		err := updateStmt(fset, t.Init)
		if err != nil {
			return err
		}
		err = updateStmt(fset, t.Assign)
		if err != nil {
			return err
		}

		return updateBlockStmt(fset, t.Body)
	case *ast.CaseClause:
		for _, stmnt := range t.Body {
			err := updateStmt(fset, stmnt)
			if err != nil {
				return err
			}
		}
		return nil
	case *ast.ForStmt:
		err := updateStmt(fset, t.Init)
		if err != nil {
			return err
		}
		err = updateExpr(fset, t.Cond)
		if err != nil {
			return err
		}
		err = updateStmt(fset, t.Post)
		if err != nil {
			return err
		}

		return updateBlockStmt(fset, t.Body)
	case *ast.IncDecStmt:
		return nil
	case *ast.BranchStmt:
		return nil
	case *ast.SwitchStmt:
		err := updateStmt(fset, t.Init)
		if err != nil {
			return err
		}

		return updateBlockStmt(fset, t.Body)
	default:
		return fmt.Errorf("unknown stmnt %T : %v", t, fset.Position(t.Pos()))
	}

}

func updateSpec(fset *token.FileSet, s ast.Spec) error {
	switch t := s.(type) {
	case *ast.ImportSpec:
		return nil
	case *ast.TypeSpec:
		return nil
	case *ast.ValueSpec:
		for _, val := range t.Values {
			err := updateExpr(fset, val)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown spec %T : %v", t, fset.Position(t.Pos()))
	}
}

func updateDecl(fset *token.FileSet, d ast.Decl) error {
	switch t := d.(type) {
	case *ast.GenDecl:
		for _, spec := range t.Specs {
			err := updateSpec(fset, spec)
			if err != nil {
				return err
			}
		}
		return nil
	case *ast.FuncDecl:
		return updateFuncDecl(fset, t)
	default:
		return fmt.Errorf("unknown decl %T : %v", t, fset.Position(t.Pos()))
	}
}

func getFncName(fset *token.FileSet, f ast.Expr) (string, error) {

	switch t := f.(type) {
	case *ast.Ident:
		return t.Name, nil
	case *ast.IndexExpr:
		return getFncName(fset, t.X)
	case *ast.IndexListExpr:
		return getFncName(fset, t.X)
	case *ast.SelectorExpr:
		return getFncName(fset, t.Sel)
	case *ast.CallExpr:
		return "<anon>", nil
	case *ast.ArrayType:
		name, err := getFncName(fset, t.Elt)
		return "[]" + name, err
	case *ast.StarExpr:
		name, err := getFncName(fset, t.X)
		return "*" + name, err
	case *ast.ParenExpr:
		return getFncName(fset, t.X)
	default:
		return "", fmt.Errorf("unknown fnc expr %T : Pos: %v", t, fset.Position(t.Pos()))

	}

}

func updateExpr(fset *token.FileSet, x ast.Expr) error {
	if x == nil {
		return nil
	}

	switch t := x.(type) {
	case *ast.Ident:
		return nil
	case *ast.FuncLit:
		return updateBlockStmt(fset, t.Body)
	case *ast.IndexExpr:
		err := updateExpr(fset, t.X)
		if err != nil {
			return err
		}
		return updateExpr(fset, t.Index)
	case *ast.IndexListExpr:
		return updateExpr(fset, t.X)
	case *ast.KeyValueExpr:
		return updateExpr(fset, t.Value)
	case *ast.CompositeLit:
		for _, elt := range t.Elts {
			err := updateExpr(fset, elt)
			if err != nil {
				return err
			}
		}
		return nil
	case *ast.UnaryExpr:
		return updateExpr(fset, t.X)
	case *ast.BinaryExpr:
		err := updateExpr(fset, t.X)
		if err != nil {
			return err
		}
		return updateExpr(fset, t.Y)
	case *ast.BasicLit:
		return nil
	case *ast.SelectorExpr:
		err := updateExpr(fset, t.X)
		if err != nil {
			return err
		}
		return updateExpr(fset, t.Sel)
	case *ast.CallExpr:

		err := updateExpr(fset, t.Fun)
		if err != nil {
			return err
		}

		for _, arg := range t.Args {
			err := updateExpr(fset, arg)
			if err != nil {
				return err
			}
		}

		fncName, err := getFncName(fset, t.Fun)
		if err != nil {
			return err
		}

		log.Println("Func", fncName, fset.Position(t.Pos()))

		if fncName == "MustIToSliceOf" {
			err := removeAction(t)
			if err != nil {
				return err
			}
		}

		return nil
	case *ast.TypeAssertExpr:
		return nil
	case *ast.StarExpr:
		return updateExpr(fset, t.X)
	case *ast.MapType:
		return nil
	case *ast.ArrayType:
		return nil
	case *ast.SliceExpr:

		err := updateExpr(fset, t.X)
		if err != nil {
			return err
		}

		err = updateExpr(fset, t.Low)
		if err != nil {
			return err
		}

		err = updateExpr(fset, t.High)
		if err != nil {
			return err
		}

		return updateExpr(fset, t.Max)

	case *ast.ParenExpr:
		return updateExpr(fset, t.X)
	default:
		return fmt.Errorf("unknown expr %T : %v", t, fset.Position(t.Pos()))
	}
}

func paramSwitch(call *ast.CallExpr) error {
	call.Args[0], call.Args[1] = call.Args[1], call.Args[0]
	return nil
}

func removeAction(call *ast.CallExpr) error {

	call.Fun = &ast.Ident{
		Name: "MustGet",
	}

	call.Args[0] = &ast.CallExpr{
		Fun: &ast.Ident{
			Name: "ISliceOf",
		},
		Args: []ast.Expr{
			call.Args[0],
			call.Args[2],
		},
	}
	call.Args = call.Args[0:2]

	return nil
}
