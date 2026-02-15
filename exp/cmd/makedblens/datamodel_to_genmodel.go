package main

import (
	"strings"

	. "github.com/spearson78/go-optic/internal/codegen"
)

func DmToFd(dm *DataModel, packageName string, rootObjName string) FileDef {

	fd := FileDef{
		Package: packageName,
		Imports: []string{
			"database/sql",
			"github.com/spearson78/go-optic",
		},
	}

	oPrefix := strings.ToLower(rootObjName)
	upperOPrefix := rootObjName

	oStruct := StructDef{
		Name: oPrefix,
	}

	for i := range dm.Structs {
		str := &dm.Structs[i]

		if str.TypeParam == "" {

			buildWrapperStruct(upperOPrefix, oPrefix+"L", "l", &fd, str, &oStruct, dm.TypeAliases)
			//fd.Funcs = append(fd.Funcs, genFrom(oPrefix, str, "dbL", "l"))
			//buildWrapperStruct(upperOPrefix, oPrefix+"S", "s", &fd, str, &oStruct, dm.TypeAliases)
			//buildMapWrapperStruct(upperOPrefix, oPrefix+"M", "m", &fd, str, &oStruct, dm.TypeAliases)

			if str.Name == *sourceConstraint {
				buildOMethod(upperOPrefix, oPrefix, &fd, str, &oStruct, dm.TypeAliases, packageName == "osql")
			} else {
				buildTStruct(upperOPrefix, oPrefix, &fd, str, &oStruct, dm.TypeAliases)
			}

			//buildOpt(upperOPrefix, &fd, str, dm.TypeAliases)

		}
	}

	fd.Structs = append(fd.Structs, oStruct)

	/*
		fd.Funcs = append(fd.Funcs, FuncDef{
			Name:   upperOPrefix,
			Params: []Param{},
			ReturnTypes: []TypeExpression{
				Star{
					Type: TypeDef{
						Name: oStruct.Name,
					},
				},
			},
			Body: []Statement{
				ReturnStmnt{
					Value: AddressExpr{
						Target: StructExpr{
							Name: oStruct.Name,
						},
					},
				},
			},
		})
	*/

	fd.Vars = append(fd.Vars, VarDef{

		Name: upperOPrefix,
		Value: StructExpr{
			Type: TypeDef{Name: oStruct.Name},
		},
	})

	return fd
}

func buildOMethod(oPrefix string, prefix string, fd *FileDef, str *Struct, oStruct *StructDef, typeAliases map[string]string, osql bool) {

	tableName := "osql.Table"
	if osql {
		tableName = "Table"
	}

	for _, field := range str.Fields {

		m := MethodDef{
			Name: field.Name,
		}

		indexType := "optic.Void"
		if field.Map {
			indexType = field.MapKey
		}

		m.ReturnTypes = append(m.ReturnTypes, Star{
			Type: TypeDef{
				Name: prefix + "T" + field.TypeName,
				TypeParams: []TypeExpression{
					TypeDef{
						Name: indexType,
					},
					TypeDef{
						Name: "optic.ReturnOne",
					},
					TypeDef{
						Name: "optic.ReadWrite",
					},
					TypeDef{
						Name: "optic.UniDir",
					},
				},
			},
		})

		m.Body = []Statement{
			ReturnStmnt{
				Values: []Expression{
					AddressExpr{

						Target: StructExpr{
							Type: TypeDef{
								Name: prefix + "T" + field.TypeName,
								TypeParams: []TypeExpression{
									TypeDef{
										Name: indexType,
									},
									TypeDef{
										Name: "optic.ReturnOne",
									},
									TypeDef{
										Name: "optic.ReadWrite",
									},
									TypeDef{
										Name: "optic.UniDir",
									},
								},
							},
							Fields: []AssignField{
								{
									Name: "Optic",
									Value: CallExpr{
										Func: tableName,
										TypeParams: []TypeExpression{
											TypeDef{
												Name: indexType,
											},
											TypeDef{
												Name: field.TypeName,
											},
											TypeDef{
												Name: *sourceConstraint,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		oStruct.Methods = append(oStruct.Methods, m)

		fd.Funcs = append(fd.Funcs, FuncDef{
			Name: "DB" + field.TypeName + "From",
			TypeParams: []TypeExpression{
				TypeDef{
					Name:       "I",
					Constraint: "comparable",
				},
				TypeDef{
					Name: "RET",
				},
				TypeDef{
					Name: "RW",
				},
				TypeDef{
					Name: "DIR",
				},
				TypeDef{
					Name: "ERR",
				},
			},
			Params: []Param{
				{
					Name: "o",
					Type: TypeDef{
						Name: "optic.Optic",
						TypeParams: []TypeExpression{
							TypeDef{
								Name: "I",
							},
							TypeDef{
								Name: "*sql.DB",
							},
							TypeDef{
								Name: "*sql.DB",
							},
							TypeDef{
								Name: field.TypeName,
							},
							TypeDef{
								Name: field.TypeName,
							},
							TypeDef{
								Name: "RET",
							},
							TypeDef{
								Name: "RW",
							},
							TypeDef{
								Name: "DIR",
							},
							TypeDef{
								Name: "ERR",
							},
						},
					},
				},
			},
			ReturnTypes: []TypeExpression{
				Star{
					Type: TypeDef{
						Name: prefix + "L" + field.TypeName,
						TypeParams: []TypeExpression{
							TypeDef{
								Name: "I",
							},
							TypeDef{
								Name: "RET",
							},
							TypeDef{
								Name: "RW",
							},
							TypeDef{
								Name: "DIR",
							},
						},
					},
				},
			},
			Body: []Statement{
				ReturnStmnt{
					Values: []Expression{
						AddressExpr{
							Target: StructExpr{
								Type: TypeDef{
									Name: prefix + "L" + field.TypeName,
									TypeParams: []TypeExpression{
										TypeDef{
											Name: "I",
										},
										TypeDef{
											Name: "RET",
										},
										TypeDef{
											Name: "RW",
										},
										TypeDef{
											Name: "DIR",
										},
									},
								},
								Fields: []AssignField{
									{
										Name: "l" + field.TypeName,
										Value: DeRef{
											Value: CallExpr{
												Func: "O" + field.Name + "From",
												Params: []Expression{
													eErr("o"),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		})
	}
}

func buildTStruct(upperOPrefix string, lowerPrefix string, fd *FileDef, str *Struct, oStruct *StructDef, typeAliases map[string]string) {

	tStruct := StructDef{
		Name: lowerPrefix + "T" + str.Name,
		TypeParams: []TypeExpression{
			TypeDef{
				Name:       "I",
				Constraint: "comparable",
			},
			TypeDef{
				Name:       "RET",
				Constraint: "optic.TReturnOne",
			},
			TypeDef{
				Name: "RW",
			},
			TypeDef{
				Name: "DIR",
			},
		},
		Fields: []FieldDef{
			{
				Name: "",
				Type: TypeDef{
					Name: "optic.Optic",
					TypeParams: []TypeExpression{
						TypeDef{
							Name: "optic.Void",
						},
						TypeDef{
							Name: "*sql.DB",
						},
						TypeDef{
							Name: "*sql.DB",
						},
						TypeDef{
							Name: "optic.Collection",
							TypeParams: []TypeExpression{
								TypeDef{
									Name: "I",
								},
								TypeDef{
									Name: str.Name,
								},
								TypeDef{
									Name: "optic.Err",
								},
							},
						},
						TypeDef{
							Name: "optic.Collection",
							TypeParams: []TypeExpression{
								TypeDef{
									Name: "I",
								},
								TypeDef{
									Name: str.Name,
								},
								TypeDef{
									Name: "optic.Err",
								},
							},
						},
						TypeDef{
							Name: "RET",
						},
						TypeDef{
							Name: "RW",
						},
						TypeDef{
							Name: "DIR",
						},
						TypeDef{
							Name: "optic.Err",
						},
					},
				},
			},
		},
	}

	tStruct.Methods = append(tStruct.Methods, MethodDef{
		Name: "Traverse",
		ReturnTypes: []TypeExpression{
			Star{
				Type: TypeDef{
					Name: lowerPrefix + "L" + str.Name,
					TypeParams: []TypeExpression{
						TypeDef{
							Name: "I",
						},
						TypeDef{
							Name: "optic.ReturnMany",
						},
						TypeDef{
							Name: "optic.CompositionTree",
							TypeParams: []TypeExpression{
								TypeDef{
									Name: "RW",
								},
								TypeDef{
									Name: "optic.ReadWrite",
								},
							},
						},
						TypeDef{
							Name: "optic.UniDir",
						},
					},
				},
			},
		},
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: upperOPrefix + str.Name + "From",
						Params: []Expression{
							composeRetMUd(
								"s.Optic",
								CallExpr{
									Func: "optic.TraverseColE",
									TypeParams: []TypeExpression{
										TypeDef{
											Name: "I",
										},
										TypeDef{
											Name: str.Name,
										},
										TypeDef{
											Name: "optic.Err",
										},
									},
								},
							),
						},
					},
				},
			},
		},
	})

	fd.Structs = append(fd.Structs, tStruct)
}

func eErr(o Expression) Expression {
	return CallExpr{
		Func: "optic.EErr",
		Params: []Expression{
			o,
		},
	}
}

func buildWrapperStruct(oPrefix string, prefix string, parentPrefix string, fd *FileDef, str *Struct, oStruct *StructDef, typeAliases map[string]string) {

	lStruct := StructDef{
		Name: prefix + str.Name,
		TypeParams: []TypeExpression{
			TypeDef{
				Name:       "I",
				Constraint: "comparable",
			},
			TypeDef{
				Name: "RET",
			},
			TypeDef{
				Name: "RW",
			},
			TypeDef{
				Name: "DIR",
			},
		},
		Fields: []FieldDef{
			{
				Name: "",
				Type: TypeDef{
					Name: parentPrefix + str.Name,
					TypeParams: []TypeExpression{
						TypeDef{
							Name: "I",
						},
						TypeDef{
							Name: "*sql.DB",
						},
						TypeDef{
							Name: "*sql.DB",
						},
						TypeDef{
							Name: "RET",
						},
						TypeDef{
							Name: "RW",
						},
						TypeDef{
							Name: "DIR",
						},
						TypeDef{
							Name: "optic.Err",
						},
					},
				},
			},
		},
	}

	fd.Structs = append(fd.Structs, lStruct)
}

func composeRetMUd(params ...Expression) Expression {
	return CallExpr{
		Func: "optic.RetM",
		Params: []Expression{
			CallExpr{
				Func: "optic.Ud",
				Params: []Expression{
					CallExpr{
						Func:   "optic.Compose",
						Params: params,
					},
				},
			},
		},
	}
}
