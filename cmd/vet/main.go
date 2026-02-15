package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"strings"
)

func collectFiles(dir string, files map[string]*ast.File, fset *token.FileSet) error {

	goFiles, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, goFile := range goFiles {
		if goFile.IsDir() {
			err := collectFiles(path.Join(dir, goFile.Name()), files, fset)
			if err != nil {
				return err
			}
		} else if strings.HasSuffix(goFile.Name(), ".go") {
			node, err := parser.ParseFile(fset, path.Join(dir, goFile.Name()), nil, parser.ParseComments)
			if err != nil {
				return err
			}
			files[goFile.Name()] = node
		}
	}
	return nil
}

func main() {

	goFileName := os.Args[1]

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
	} else if goFileName == "..." {
		err := collectFiles(".", files, fset)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		node, err := parser.ParseFile(fset, goFileName, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}
		files[goFileName] = node
	}

	for _, ast := range files {
		err := vetFile(fset, ast)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func vetFile(fset *token.FileSet, f *ast.File) error {

	found := false
	importName := "optic"
	for _, impt := range f.Imports {
		if impt.Path.Value == `"github.com/spearson78/go-optic"` {
			found = true
			if impt.Name != nil {
				if impt.Name.Name == "." {
					importName = ""
				} else {
					importName = impt.Name.Name
				}
			}
			break
		}
	}

	if !found {
		return nil
	}

	for _, decl := range f.Decls {
		err := vetDecl(fset, decl, nil, nil, importName)
		if err != nil {
			return err
		}
	}

	return nil
}

func vetFuncDecl(fset *token.FileSet, f *ast.FuncDecl, importName string) error {

	var opticParams []string
	var reducerParams []string

	for _, p := range f.Type.Params.List {
		if il, ok := p.Type.(*ast.IndexListExpr); ok {
			if importName == "" {
				if id, ok := il.X.(*ast.Ident); ok {
					if id.Name == "Optic" || id.Name == "OpticRO" || id.Name == "Predicate" || id.Name == "PredicateE" || id.Name == "Operation" || id.Name == "PredicateI" || id.Name == "PredicateIE" || id.Name == "OperationI" {
						opticParams = append(opticParams, p.Names[0].Name)
					}

					if id.Name == "Reduction" || id.Name == "ReductionP" {
						reducerParams = append(reducerParams, p.Names[0].Name)
					}
				}
			} else {
				if sel, ok := il.X.(*ast.SelectorExpr); ok {
					if packageName, ok := sel.X.(*ast.Ident); ok {
						if packageName.Name == importName {
							id := sel.Sel

							if id.Name == "Optic" || id.Name == "OpticRO" || id.Name == "Predicate" || id.Name == "PredicateE" || id.Name == "Operation" || id.Name == "PredicateI" || id.Name == "PredicateIE" || id.Name == "OperationI" {
								opticParams = append(opticParams, p.Names[0].Name)
							}

							if id.Name == "Reduction" || id.Name == "ReductionP" {
								reducerParams = append(reducerParams, p.Names[0].Name)
							}
						}

					}
				}
			}

		}

	}

	return vetBlockStmt(fset, f.Body, opticParams, reducerParams, importName)
}

func vetBlockStmt(fset *token.FileSet, block *ast.BlockStmt, opticParams []string, reducerParams []string, importName string) error {
	for _, stmnt := range block.List {
		err := vetStmt(fset, stmnt, opticParams, reducerParams, importName)
		if err != nil {
			return err
		}
	}

	return nil
}

func vetStmt(fset *token.FileSet, stmnt ast.Stmt, opticParams []string, reducerParams []string, importName string) error {

	if stmnt == nil {
		return nil
	}

	switch t := stmnt.(type) {
	case *ast.ReturnStmt:
		for _, x := range t.Results {
			err := vetExpr(fset, x, opticParams, reducerParams, importName)
			if err != nil {
				return err
			}
		}
		return nil
	case *ast.ExprStmt:
		return vetExpr(fset, t.X, opticParams, reducerParams, importName)
	case *ast.AssignStmt:
		for _, x := range t.Rhs {
			err := vetExpr(fset, x, opticParams, reducerParams, importName)
			if err != nil {
				return err
			}
		}
		return nil
	case *ast.DeclStmt:
		return vetDecl(fset, t.Decl, opticParams, reducerParams, importName)
	case *ast.RangeStmt:
		return vetExpr(fset, t.X, opticParams, reducerParams, importName)
	case *ast.IfStmt:
		err := vetStmt(fset, t.Init, opticParams, reducerParams, importName)
		if err != nil {
			return err
		}

		err = vetExpr(fset, t.Cond, opticParams, reducerParams, importName)
		if err != nil {
			return err
		}

		err = vetStmt(fset, t.Else, opticParams, reducerParams, importName)
		if err != nil {
			return err
		}

		return vetBlockStmt(fset, t.Body, opticParams, reducerParams, importName)
	case *ast.DeferStmt:
		return vetExpr(fset, t.Call, opticParams, reducerParams, importName)
	case *ast.TypeSwitchStmt:
		err := vetStmt(fset, t.Init, opticParams, reducerParams, importName)
		if err != nil {
			return err
		}
		err = vetStmt(fset, t.Assign, opticParams, reducerParams, importName)
		if err != nil {
			return err
		}

		return vetBlockStmt(fset, t.Body, opticParams, reducerParams, importName)
	case *ast.CaseClause:
		for _, stmnt := range t.Body {
			err := vetStmt(fset, stmnt, opticParams, reducerParams, importName)
			if err != nil {
				return err
			}
		}
		return nil
	case *ast.ForStmt:
		err := vetStmt(fset, t.Init, opticParams, reducerParams, importName)
		if err != nil {
			return err
		}
		err = vetExpr(fset, t.Cond, opticParams, reducerParams, importName)
		if err != nil {
			return err
		}
		err = vetStmt(fset, t.Post, opticParams, reducerParams, importName)
		if err != nil {
			return err
		}

		return vetBlockStmt(fset, t.Body, opticParams, reducerParams, importName)
	case *ast.IncDecStmt:
		return nil
	case *ast.BranchStmt:
		return nil
	case *ast.SwitchStmt:
		err := vetStmt(fset, t.Init, opticParams, reducerParams, importName)
		if err != nil {
			return err
		}

		return vetBlockStmt(fset, t.Body, opticParams, reducerParams, importName)
	case *ast.BlockStmt:
		return vetBlockStmt(fset, t, opticParams, reducerParams, importName)
	default:
		return fmt.Errorf("unknown stmnt %T : %v", t, fset.Position(t.Pos()))
	}

}

func vetSpec(fset *token.FileSet, s ast.Spec, opticParams []string, reducerParams []string, importName string) error {
	switch t := s.(type) {
	case *ast.ImportSpec:
		return nil
	case *ast.TypeSpec:
		return nil
	case *ast.ValueSpec:
		for _, val := range t.Values {
			err := vetExpr(fset, val, opticParams, reducerParams, importName)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown spec %T : %v", t, fset.Position(t.Pos()))
	}
}

func vetDecl(fset *token.FileSet, d ast.Decl, opticParams []string, reducerParams []string, importName string) error {
	switch t := d.(type) {
	case *ast.GenDecl:
		for _, spec := range t.Specs {
			err := vetSpec(fset, spec, opticParams, reducerParams, importName)
			if err != nil {
				return err
			}
		}
		return nil
	case *ast.FuncDecl:
		return vetFuncDecl(fset, t, importName)
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
	case *ast.FuncLit:
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

func vetExpr(fset *token.FileSet, x ast.Expr, opticParams []string, reducerParams []string, importName string) error {
	if x == nil {
		return nil
	}

	switch t := x.(type) {
	case *ast.Ident:
		return nil
	case *ast.FuncLit:
		return vetBlockStmt(fset, t.Body, opticParams, reducerParams, importName)
	case *ast.IndexExpr:
		err := vetExpr(fset, t.X, opticParams, reducerParams, importName)
		if err != nil {
			return err
		}
		return vetExpr(fset, t.Index, opticParams, reducerParams, importName)
	case *ast.IndexListExpr:
		return vetExpr(fset, t.X, opticParams, reducerParams, importName)
	case *ast.KeyValueExpr:
		return vetExpr(fset, t.Value, opticParams, reducerParams, importName)
	case *ast.CompositeLit:
		for _, elt := range t.Elts {
			err := vetExpr(fset, elt, opticParams, reducerParams, importName)
			if err != nil {
				return err
			}
		}
		return nil
	case *ast.UnaryExpr:
		return vetExpr(fset, t.X, opticParams, reducerParams, importName)
	case *ast.BinaryExpr:
		err := vetExpr(fset, t.X, opticParams, reducerParams, importName)
		if err != nil {
			return err
		}
		return vetExpr(fset, t.Y, opticParams, reducerParams, importName)
	case *ast.BasicLit:
		return nil
	case *ast.SelectorExpr:
		err := vetExpr(fset, t.X, opticParams, reducerParams, importName)
		if err != nil {
			return err
		}
		return vetExpr(fset, t.Sel, opticParams, reducerParams, importName)
	case *ast.CallExpr:

		err := vetExpr(fset, t.Fun, opticParams, reducerParams, importName)
		if err != nil {
			return err
		}

		for _, arg := range t.Args {
			err := vetExpr(fset, arg, opticParams, reducerParams, importName)
			if err != nil {
				return err
			}
		}

		fncName, err := getFncName(fset, t.Fun)
		if err != nil {
			return err
		}

		if fncName == "ExprDef" {
			passedOptics := make(map[string]struct{})

			for i := 1; i < len(t.Args); i++ {
				if id, ok := t.Args[i].(*ast.Ident); ok {
					passedOptics[id.Name] = struct{}{}
				}
			}

			for _, opticParam := range opticParams {
				if _, found := passedOptics[opticParam]; !found {
					log.Println("Missing optic param", opticParam, fset.Position(t.Pos()))
				}
			}

		}

		if fncName == "ReducerExprDef" {
			passedReducers := make(map[string]struct{})

			for i := 1; i < len(t.Args); i++ {
				if id, ok := t.Args[i].(*ast.Ident); ok {
					passedReducers[id.Name] = struct{}{}
				}
			}

			for _, reducerParam := range reducerParams {
				if _, found := passedReducers[reducerParam]; !found {
					log.Println("Missing reducer param", reducerParam, fset.Position(t.Pos()))
				}
			}

		}

		return nil
	case *ast.TypeAssertExpr:
		return nil
	case *ast.StarExpr:
		return vetExpr(fset, t.X, opticParams, reducerParams, importName)
	case *ast.MapType:
		return nil
	case *ast.ArrayType:
		return nil
	case *ast.SliceExpr:

		err := vetExpr(fset, t.X, opticParams, reducerParams, importName)
		if err != nil {
			return err
		}

		err = vetExpr(fset, t.Low, opticParams, reducerParams, importName)
		if err != nil {
			return err
		}

		err = vetExpr(fset, t.High, opticParams, reducerParams, importName)
		if err != nil {
			return err
		}

		return vetExpr(fset, t.Max, opticParams, reducerParams, importName)

	case *ast.ParenExpr:
		return vetExpr(fset, t.X, opticParams, reducerParams, importName)
	default:
		return fmt.Errorf("unknown expr %T : %v", t, fset.Position(t.Pos()))
	}
}
