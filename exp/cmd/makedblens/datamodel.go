package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"strings"
)

type DataModel struct {
	Structs     []Struct
	Interfaces  []Interface
	TypeAliases map[string]string
}

type Struct struct {
	Name      string
	TypeParam string
	Fields    []Field
}

type Field struct {
	Name        string
	Slice       bool
	Map         bool
	Pointer     bool
	MapKey      string
	TypeName    string
	TypeParam   string
	TypePackage []string
	StructTag   string
}

type Interface struct {
	Name string
}

func exprToPackage(path []string, x ast.Expr) []string {

	switch xt := x.(type) {
	case *ast.Ident:
		return append(path, xt.Name)
	case *ast.SelectorExpr:

		path = append(path, xt.Sel.Name)
		path = exprToPackage(path, xt.X)

		return path
	default:
		log.Fatalf("exprToPackage unknown type %v (%T) \n", x, x)
	}

	return path
}

func getMapKey(ft *ast.MapType) (string, error) {
	if mapKey, ok := ft.Key.(*ast.Ident); ok {
		return mapKey.Name, nil
	} else {
		return "", fmt.Errorf("unknown map key type : %T", ft.Key)
	}
}

func Parse(files []*ast.File) DataModel {

	var dm DataModel
	dm.TypeAliases = make(map[string]string)

	for _, f := range files {
		for _, d := range f.Decls {
			switch v := d.(type) {
			case *ast.GenDecl:
				if v.Tok == token.TYPE {
					for _, s := range v.Specs {
						ts := s.(*ast.TypeSpec)
						switch st := ts.Type.(type) {
						case *ast.StructType:

							if strings.HasPrefix(ts.Name.Name, "Const") {
								continue
							}

							var typeParams []string
							if ts.TypeParams != nil {
								for _, tp := range ts.TypeParams.List {
									log.Println(ts.Name.Name, tp.Names[0].Name)
									typeParams = append(typeParams, tp.Names[0].Name)
								}
							}

							//only a single type param is supported
							if len(typeParams) > 1 {
								continue
							}

							typeParam := ""
							if len(typeParams) == 1 {
								typeParam = typeParams[0]
							}

							dmStruct := Struct{
								Name:      ts.Name.Name,
								TypeParam: typeParam,
							}

							for _, f := range st.Fields.List {
								if f.Names == nil {
									continue //Anonymous struct
								}
								name := f.Names[0].Name

								var structTag string
								if f.Tag != nil {
									structTag = f.Tag.Value
								}

								switch ft := f.Type.(type) {
								case *ast.Ident:
									dmStruct.Fields = append(dmStruct.Fields, Field{
										Name:        name,
										Slice:       false,
										Pointer:     false,
										TypeName:    ft.Name,
										TypePackage: []string{},
										StructTag:   structTag,
									})
								case *ast.SelectorExpr:
									dmStruct.Fields = append(dmStruct.Fields, Field{
										Name:        name,
										Slice:       false,
										Pointer:     false,
										TypeName:    ft.Sel.Name,
										TypePackage: exprToPackage(nil, ft.X),
										StructTag:   structTag,
									})
								case *ast.ArrayType:

									switch elt := ft.Elt.(type) {
									case *ast.Ident:
										dmStruct.Fields = append(dmStruct.Fields, Field{
											Name:        name,
											Slice:       true,
											Pointer:     false,
											TypeName:    elt.Name,
											TypePackage: []string{},
											StructTag:   structTag,
										})
									case *ast.SelectorExpr:
										dmStruct.Fields = append(dmStruct.Fields, Field{
											Name:        name,
											Slice:       true,
											Pointer:     false,
											TypeName:    elt.Sel.Name,
											TypePackage: exprToPackage(nil, elt.X),
											StructTag:   structTag,
										})
									case *ast.StarExpr:
										switch stelt := elt.X.(type) {
										case *ast.Ident:
											dmStruct.Fields = append(dmStruct.Fields, Field{
												Name:        name,
												Slice:       true,
												Pointer:     true,
												TypeName:    stelt.Name,
												TypePackage: []string{},
												StructTag:   structTag,
											})
										case *ast.SelectorExpr:
											dmStruct.Fields = append(dmStruct.Fields, Field{
												Name:        name,
												Slice:       true,
												Pointer:     true,
												TypeName:    stelt.Sel.Name,
												TypePackage: exprToPackage(nil, elt.X),
												StructTag:   structTag,
											})
										}
									case *ast.IndexExpr:
										dmStruct.Fields = append(dmStruct.Fields, Field{
											Name:        name,
											Slice:       true,
											Pointer:     false,
											TypeName:    elt.X.(*ast.Ident).Name,
											TypeParam:   elt.Index.(*ast.Ident).Name,
											TypePackage: []string{},
											StructTag:   structTag,
										})
									default:
										log.Fatalf("\tField %v unknown array element type %v (%T) \n", f.Names, elt, elt)
									}
								case *ast.MapType:

									mapKey, err := getMapKey(ft)
									if err != nil {
										log.Fatal(err)
									}

									switch elt := ft.Value.(type) {
									case *ast.Ident:
										dmStruct.Fields = append(dmStruct.Fields, Field{
											Name:        name,
											Slice:       false,
											Pointer:     false,
											Map:         true,
											MapKey:      mapKey,
											TypeName:    elt.Name,
											TypePackage: []string{},
											StructTag:   structTag,
										})
									case *ast.SelectorExpr:
										dmStruct.Fields = append(dmStruct.Fields, Field{
											Name:        name,
											Slice:       false,
											Pointer:     false,
											Map:         true,
											MapKey:      mapKey,
											TypeName:    elt.Sel.Name,
											TypePackage: exprToPackage(nil, elt.X),
											StructTag:   structTag,
										})
									case *ast.StarExpr:
										switch stelt := elt.X.(type) {
										case *ast.Ident:
											dmStruct.Fields = append(dmStruct.Fields, Field{
												Name:        name,
												Slice:       false,
												Pointer:     true,
												Map:         true,
												MapKey:      mapKey,
												TypeName:    stelt.Name,
												TypePackage: []string{},
												StructTag:   structTag,
											})
										case *ast.SelectorExpr:
											dmStruct.Fields = append(dmStruct.Fields, Field{
												Name:        name,
												Slice:       false,
												Pointer:     true,
												Map:         true,
												MapKey:      mapKey,
												TypeName:    stelt.Sel.Name,
												TypePackage: exprToPackage(nil, elt.X),
												StructTag:   structTag,
											})
										}
									case *ast.IndexExpr:
										dmStruct.Fields = append(dmStruct.Fields, Field{
											Name:        name,
											Slice:       false,
											Pointer:     false,
											Map:         true,
											MapKey:      mapKey,
											TypeName:    elt.X.(*ast.Ident).Name,
											TypeParam:   elt.Index.(*ast.Ident).Name,
											TypePackage: []string{},
											StructTag:   structTag,
										})
									default:
										log.Fatalf("\tField %v unknown map element type %v (%T) \n", f.Names, elt, elt)
									}
								case *ast.StarExpr:
									switch stelt := ft.X.(type) {
									case *ast.Ident:
										dmStruct.Fields = append(dmStruct.Fields, Field{
											Name:        name,
											Slice:       false,
											Pointer:     true,
											TypeName:    stelt.Name,
											TypePackage: []string{},
											StructTag:   structTag,
										})
									case *ast.SelectorExpr:
										dmStruct.Fields = append(dmStruct.Fields, Field{
											Name:        name,
											Slice:       false,
											Pointer:     true,
											TypeName:    stelt.Sel.Name,
											TypePackage: exprToPackage(nil, stelt.X),
											StructTag:   structTag,
										})
									case *ast.IndexExpr:
										dmStruct.Fields = append(dmStruct.Fields, Field{
											Name:        name,
											Slice:       false,
											Pointer:     true,
											TypeName:    stelt.X.(*ast.Ident).Name,
											TypeParam:   stelt.Index.(*ast.Ident).Name,
											TypePackage: []string{},
											StructTag:   structTag,
										})
									default:
										log.Fatalf("\tField %v unknown star type %v (%T) \n", f.Names, stelt, stelt)
									}
								case *ast.IndexExpr:
									dmStruct.Fields = append(dmStruct.Fields, Field{
										Name:        name,
										Slice:       false,
										Pointer:     false,
										TypeName:    ft.X.(*ast.Ident).Name,
										TypeParam:   ft.Index.(*ast.Ident).Name,
										TypePackage: []string{},
										StructTag:   structTag,
									})

								default:
									log.Fatalf("\tStruct %v Field %v unknown type %v (%T) \n", name, f.Names, ft, ft)
								}

							}

							dm.Structs = append(dm.Structs, dmStruct)
						case *ast.InterfaceType:
							if strings.HasPrefix(ts.Name.Name, "Const") {
								continue
							}

							dmInterface := Interface{
								Name: ts.Name.Name,
							}

							dm.Interfaces = append(dm.Interfaces, dmInterface)
						case *ast.Ident:
							dm.TypeAliases[ts.Name.Name] = st.Name
							log.Printf("type alias %v : %v", ts.Name.Name, st.Name)
						default:
							log.Fatalf("TypeSpec  unknown type %v (%T) \n", ts, ts)
						}
					}
				}
			default:
				//log.Fatalf("Unknown type %v (%T) \n", v, v)
			}
		}
	}

	return dm
}
