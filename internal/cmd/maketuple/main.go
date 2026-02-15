package main

import (
	"fmt"
	"log"
	"os"

	. "github.com/spearson78/go-optic/internal/codegen"
)

var tuplePartNames = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}

func genTupleElements(fd *FileDef) {

	for tuple := 1; tuple < len(tuplePartNames); tuple++ {

		//Tuple Part Lenses
		for i := 0; i <= tuple; i++ {

			tnFnc := FuncDef{
				Docs: []string{
					fmt.Sprintf("T%v%v returns a [Lens] focusing on element %v of a [lo.Tuple%v]", tuple+1, tuplePartNames[i], i, tuple+1),
					"",
					fmt.Sprintf("See: [T%v%vP] for a polymorphic version", tuple+1, tuplePartNames[i]),
				},
				Name: fmt.Sprintf("T%v%v", tuple+1, tuplePartNames[i]),
			}

			tnpFnc := FuncDef{
				Docs: []string{
					fmt.Sprintf("T%v%vP returns a polymorphic [Lens] focusing on element %v of a [lo.Tuple%v]", tuple+1, tuplePartNames[i], i, tuple+1),
					"",
					fmt.Sprintf("See: [T%v%v] for a non polymorphic version", tuple+1, tuplePartNames[i]),
				},
				Name: fmt.Sprintf("T%v%vP", tuple+1, tuplePartNames[i]),
			}

			tupleType := TypeDef{
				Name: fmt.Sprintf("lo.Tuple%v", tuple+1),
			}

			tuplePType := TypeDef{
				Name: fmt.Sprintf("lo.Tuple%v", tuple+1),
			}

			partTypeDef := TypeDef{Name: tuplePartNames[i]}
			partTypeDefP := TypeDef{Name: fmt.Sprintf("%v2", tuplePartNames[i])}

			var tnCallTypeParams []TypeExpression
			var loTNParams []Expression

			for typeParamNum := 0; typeParamNum <= tuple; typeParamNum++ {

				typeDef := TypeDef{Name: tuplePartNames[typeParamNum]}

				tnFnc.TypeParams = append(tnFnc.TypeParams, typeDef)
				tupleType.TypeParams = append(tupleType.TypeParams, typeDef)

				tnpFnc.TypeParams = append(tnpFnc.TypeParams, typeDef)
				if typeParamNum == i {
					loTNParams = append(loTNParams, "focus")
					typeDefP := TypeDef{Name: fmt.Sprintf("%v2", tuplePartNames[typeParamNum])}
					tuplePType.TypeParams = append(tuplePType.TypeParams, typeDefP)
				} else {
					loTNParams = append(loTNParams, fmt.Sprintf("source.%v", tuplePartNames[typeParamNum]))
					tuplePType.TypeParams = append(tuplePType.TypeParams, typeDef)
				}

				tnCallTypeParams = append(tnCallTypeParams, typeDef)

			}

			tnCallTypeParams = append(tnCallTypeParams, partTypeDef)

			tnpFnc.TypeParams = append(tnpFnc.TypeParams, TypeDef{Name: fmt.Sprintf("%v2", tuplePartNames[i])})

			tnFnc.ReturnTypes = append(tnFnc.ReturnTypes, TypeDef{
				Name: "Optic",
				TypeParams: []TypeExpression{
					TypeDef{Name: "int"},
					tupleType,
					tupleType,
					partTypeDef,
					partTypeDef,
					TypeDef{Name: "ReturnOne"},
					TypeDef{Name: "ReadWrite"},
					TypeDef{Name: "UniDir"},
					TypeDef{Name: "Pure"},
				},
			})

			tnpFnc.ReturnTypes = []TypeExpression{
				TypeDef{
					Name: "Optic",
					TypeParams: []TypeExpression{
						TypeDef{Name: "int"},
						tupleType,
						tuplePType,
						partTypeDef,
						partTypeDefP,
						TypeDef{Name: "ReturnOne"},
						TypeDef{Name: "ReadWrite"},
						TypeDef{Name: "UniDir"},
						TypeDef{Name: "Pure"},
					},
				},
			}

			tnFnc.Body = []Statement{
				ReturnStmnt{
					Values: []Expression{
						CallExpr{
							Func:       tnpFnc.Name,
							TypeParams: tnCallTypeParams,
						},
					},
				},
			}

			ctxType := TypeDef{Name: "context.Context"}
			intType := TypeDef{Name: "int"}
			errType := TypeDef{Name: "error"}

			tnpFnc.Body = []Statement{
				ReturnStmnt{
					Values: []Expression{
						CallExpr{
							Func: "CombiLens",
							TypeParams: []TypeExpression{
								TypeDef{Name: "ReadWrite"},
								TypeDef{Name: "Pure"},
								TypeDef{Name: "int"},
								tupleType,
								tuplePType,
								partTypeDef,
								partTypeDefP,
							},
							Params: []Expression{
								FuncDef{
									Params: []Param{
										{
											Name: "ctx",
											Type: ctxType,
										},
										{
											Name: "source",
											Type: tupleType,
										},
									},
									ReturnTypes: []TypeExpression{
										intType,
										partTypeDef,
										errType,
									},
									Body: []Statement{
										ReturnStmnt{
											Values: []Expression{
												i,
												fmt.Sprintf("source.%v", tuplePartNames[i]),
												"nil",
											},
										},
									},
								},
								FuncDef{
									Params: []Param{
										{
											Name: "ctx",
											Type: ctxType,
										},
										{
											Name: "focus",
											Type: partTypeDefP,
										},
										{
											Name: "source",
											Type: tupleType,
										},
									},
									ReturnTypes: []TypeExpression{
										tuplePType,
										errType,
									},
									Body: []Statement{
										ReturnStmnt{
											Values: []Expression{
												CallExpr{
													Func:   fmt.Sprintf("lo.T%v", tuple+1),
													Params: loTNParams,
												},
												"nil",
											},
										},
									},
								},
								CallExpr{
									Func: "IxMatchComparable",
									TypeParams: []TypeExpression{
										intType,
									},
								},
								CallExpr{
									Func: "ExprDef",
									Params: []Expression{
										FuncDef{
											Params: []Param{
												{
													Name: "ot",
													Type: TypeDef{Name: "expr.OpticTypeExpr"},
												},
											},
											ReturnTypes: []TypeExpression{
												TypeDef{Name: "expr.OpticExpression"},
											},
											Body: []Statement{
												ReturnStmnt{
													Values: []Expression{
														StructExpr{
															Type: TypeDef{Name: "expr.TupleElement"},
															Fields: []AssignField{
																{
																	Name:  "OpticTypeExpr",
																	Value: "ot",
																},
																{
																	Name:  "Index",
																	Value: i,
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
						},
					},
				},
			}

			fd.Funcs = append(fd.Funcs, tnFnc)
			fd.Funcs = append(fd.Funcs, tnpFnc)

		}

	}
}

func genTupleDup(fd *FileDef) {

	for tuple := 1; tuple < len(tuplePartNames); tuple++ {

		dupFnc := FuncDef{
			Docs: []string{
				fmt.Sprintf("DupT%v returns a [Lens] focusing a [lo.Tuple%v] with all elements set to the source value.", tuple+1, tuple+1),
				"",
				" Note: under modification the first element of the tuple is used. All other elements are ignored.",
				"",
				fmt.Sprintf("See: [DupT%vP] for a polymorphic version", tuple+1),
			},
			Name: fmt.Sprintf("DupT%v", tuple+1),
			TypeParams: []TypeExpression{
				TypeDef{Name: "A"},
			},
		}

		dupFncP := FuncDef{
			Docs: []string{
				fmt.Sprintf("DupT%vP returns a polymorphic [Lens] focusing a [lo.Tuple%v] with all elements set to the source value.", tuple+1, tuple+1),
				"",
				" Note: under modification the first element of the tuple is used. All other elements are ignored.",
				"",
				fmt.Sprintf("See: [DupT%v] for a non polymorphic version", tuple+1),
			},
			Name: fmt.Sprintf("DupT%vP", tuple+1),
			TypeParams: []TypeExpression{
				TypeDef{Name: "A"},
				TypeDef{Name: "B"},
			},
		}

		tupleType := TypeDef{Name: fmt.Sprintf("lo.Tuple%v", tuple+1)}
		tupleTypeP := TypeDef{Name: fmt.Sprintf("lo.Tuple%v", tuple+1)}
		var getterParams []Expression

		for i := 0; i <= tuple; i++ {
			tupleType.TypeParams = append(tupleType.TypeParams, TypeDef{Name: "A"})
			tupleTypeP.TypeParams = append(tupleTypeP.TypeParams, TypeDef{Name: "B"})
			getterParams = append(getterParams, "source")
		}

		dupFnc.ReturnTypes = []TypeExpression{
			TypeDef{
				Name: "Optic",
				TypeParams: []TypeExpression{
					TypeDef{Name: "Void"},
					TypeDef{Name: "A"},
					TypeDef{Name: "A"},
					tupleType,
					tupleType,
					TypeDef{Name: "ReturnOne"},
					TypeDef{Name: "ReadWrite"},
					TypeDef{Name: "UniDir"},
					TypeDef{Name: "Pure"},
				},
			},
		}

		dupFncP.ReturnTypes = []TypeExpression{
			TypeDef{
				Name: "Optic",
				TypeParams: []TypeExpression{
					TypeDef{Name: "Void"},
					TypeDef{Name: "A"},
					TypeDef{Name: "B"},
					tupleType,
					tupleTypeP,
					TypeDef{Name: "ReturnOne"},
					TypeDef{Name: "ReadWrite"},
					TypeDef{Name: "UniDir"},
					TypeDef{Name: "Pure"},
				},
			},
		}

		getterFnc := FuncDef{
			Params: []Param{
				{
					Name: "source",
					Type: TypeDef{Name: "A"},
				},
			},
			ReturnTypes: []TypeExpression{
				tupleType,
			},
			Body: []Statement{
				ReturnStmnt{
					Values: []Expression{
						CallExpr{
							Func:   fmt.Sprintf("lo.T%v", tuple+1),
							Params: getterParams,
						},
					},
				},
			},
		}

		setterFnc := FuncDef{
			Params: []Param{
				{
					Name: "focus",
					Type: tupleTypeP,
				},
				{
					Name: "source",
					Type: TypeDef{Name: "A"},
				},
			},
			ReturnTypes: []TypeExpression{
				TypeDef{Name: "B"},
			},
			Body: []Statement{
				ReturnStmnt{
					Values: []Expression{
						"focus.A",
					},
				},
			},
		}

		exprDef := CallExpr{
			Func: "ExprDef",
			Params: []Expression{
				FuncDef{
					Params: []Param{
						{
							Name: "ot",
							Type: TypeDef{Name: "expr.OpticTypeExpr"},
						},
					},
					ReturnTypes: []TypeExpression{
						TypeDef{Name: "expr.OpticExpression"},
					},
					Body: []Statement{
						ReturnStmnt{
							Values: []Expression{
								StructExpr{
									Type: TypeDef{Name: "expr.TupleDup"},
									Fields: []AssignField{
										{
											Name:  "OpticTypeExpr",
											Value: "ot",
										},
										{
											Name:  "N",
											Value: tuple + 1,
										},
									},
								},
							},
						},
					},
				},
			},
		}

		dupFnc.Body = []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: dupFncP.Name,
						TypeParams: []TypeExpression{
							TypeDef{Name: "A"},
							TypeDef{Name: "A"},
						},
					},
				},
			},
		}

		dupFncP.Body = []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: "LensP",
						TypeParams: []TypeExpression{
							TypeDef{Name: "A"},
							TypeDef{Name: "B"},
							tupleType,
							tupleTypeP,
						},
						Params: []Expression{
							getterFnc,
							setterFnc,
							exprDef,
						},
					},
				},
			},
		}

		fd.Funcs = append(fd.Funcs, dupFnc, dupFncP)
	}
}

func genTupleTraverse(fd *FileDef) {

	for tuple := 1; tuple < len(tuplePartNames); tuple++ {

		travFnc := FuncDef{
			Docs: []string{
				fmt.Sprintf("TraverseT%v returns a [Traversal] focusing on the elements of a [lo.Tuple%v]", tuple+1, tuple+1),
				"",
				fmt.Sprintf("See: [TraverseT%vP] for a polymorphic version", tuple+1),
			},
			Name: fmt.Sprintf("TraverseT%v", tuple+1),
			TypeParams: []TypeExpression{
				TypeDef{Name: "A"},
			},
		}

		travPFnc := FuncDef{
			Docs: []string{
				fmt.Sprintf("TraverseT%vP returns a polymorphic [Traversal] focusing on the elements of a [lo.Tuple%v]", tuple+1, tuple+1),
				"",
				fmt.Sprintf("See: [TraverseT%v] for a non polymorphic version", tuple+1),
			},
			Name: fmt.Sprintf("TraverseT%vP", tuple+1),
			TypeParams: []TypeExpression{
				TypeDef{Name: "A"},
				TypeDef{Name: "B"},
			},
		}

		tupleType := TypeDef{Name: fmt.Sprintf("lo.Tuple%v", tuple+1)}
		tupleTypeP := TypeDef{Name: fmt.Sprintf("lo.Tuple%v", tuple+1)}
		var elementPtrs []Expression
		var setterRetParams []Expression

		for i := 0; i <= tuple; i++ {
			tupleType.TypeParams = append(tupleType.TypeParams, TypeDef{Name: "A"})
			tupleTypeP.TypeParams = append(tupleTypeP.TypeParams, TypeDef{Name: "B"})
			elementPtrs = append(elementPtrs, "&source."+tuplePartNames[i])
			setterRetParams = append(setterRetParams, fmt.Sprintf("ret[%v]", i))
		}

		travFnc.ReturnTypes = []TypeExpression{
			TypeDef{
				Name: "Optic",
				TypeParams: []TypeExpression{
					TypeDef{Name: "int"},
					tupleType,
					tupleType,
					TypeDef{Name: "A"},
					TypeDef{Name: "A"},
					TypeDef{Name: "ReturnMany"},
					TypeDef{Name: "ReadWrite"},
					TypeDef{Name: "UniDir"},
					TypeDef{Name: "Pure"},
				},
			},
		}

		travPFnc.ReturnTypes = []TypeExpression{
			TypeDef{
				Name: "Optic",
				TypeParams: []TypeExpression{
					TypeDef{Name: "int"},
					tupleType,
					tupleTypeP,
					TypeDef{Name: "A"},
					TypeDef{Name: "B"},
					TypeDef{Name: "ReturnMany"},
					TypeDef{Name: "ReadWrite"},
					TypeDef{Name: "UniDir"},
					TypeDef{Name: "Pure"},
				},
			},
		}

		intType := TypeDef{Name: "int"}
		ctxType := TypeDef{Name: "context.Context"}
		seqFType := TypeDef{
			Name: "SeqIE",
			TypeParams: []TypeExpression{
				intType,
				TypeDef{Name: "A"},
			},
		}
		errorType := TypeDef{Name: "error"}

		getterParams := []Expression{"ctx"}
		getterParams = append(getterParams, elementPtrs...)

		iterFnc := FuncDef{
			Params: []Param{
				{
					Name: "ctx",
					Type: ctxType,
				},
				{
					Name: "source",
					Type: tupleType,
				},
			},
			ReturnTypes: []TypeExpression{
				seqFType,
			},
			Body: []Statement{
				ReturnStmnt{
					Values: []Expression{
						CallExpr{
							Func:   "traverseTNIter",
							Params: getterParams,
						},
					},
				},
			},
		}

		lenFnc := FuncDef{
			Params: []Param{
				{
					Name: "ctx",
					Type: ctxType,
				},
				{
					Name: "source",
					Type: tupleType,
				},
			},
			ReturnTypes: []TypeExpression{
				intType,
				errorType,
			},
			Body: []Statement{
				ReturnStmnt{
					Values: []Expression{
						tuple + 1,
						"nil",
					},
				},
			},
		}

		modifyParams := []Expression{"ctx", "fmap"}
		modifyParams = append(modifyParams, elementPtrs...)

		modifyFnc := FuncDef{
			Params: []Param{
				{
					Name: "ctx",
					Type: ctxType,
				},
				{
					Name: "fmap",
					Type: TypeDef{
						Name: "func(index int,focus A) (B,error)",
					},
				},
				{
					Name: "source",
					Type: tupleType,
				},
			},
			ReturnTypes: []TypeExpression{
				tupleTypeP,
				errorType,
			},
			Body: []Statement{
				AssignVar{
					Declare: true,
					Vars:    []string{"ret", "err"},
					Value: CallExpr{
						Func:   "traverseTNModify",
						Params: modifyParams,
					},
				},
				ReturnStmnt{
					Values: []Expression{
						CallExpr{
							Func:   fmt.Sprintf("lo.T%v", tuple+1),
							Params: setterRetParams,
						},
						"err",
					},
				},
			},
		}

		ixGetParams := []Expression{"index"}
		ixGetParams = append(ixGetParams, elementPtrs...)

		ixGetFnc := FuncDef{
			Params: []Param{
				{
					Name: "ctx",
					Type: ctxType,
				},
				{
					Name: "index",
					Type: intType,
				},
				{
					Name: "source",
					Type: tupleType,
				},
			},
			ReturnTypes: []TypeExpression{
				seqFType,
			},
			Body: []Statement{
				ReturnStmnt{
					Values: []Expression{
						CallExpr{
							Func:   "traverseTNIxGet",
							Params: ixGetParams,
						},
					},
				},
			},
		}

		exprDef := CallExpr{
			Func: "ExprDef",
			Params: []Expression{
				FuncDef{
					Params: []Param{
						{
							Name: "ot",
							Type: TypeDef{Name: "expr.OpticTypeExpr"},
						},
					},
					ReturnTypes: []TypeExpression{
						TypeDef{Name: "expr.OpticExpression"},
					},
					Body: []Statement{
						ReturnStmnt{
							Values: []Expression{
								StructExpr{
									Type: TypeDef{Name: "expr.Traverse"},
									Fields: []AssignField{
										{
											Name:  "OpticTypeExpr",
											Value: "ot",
										},
									},
								},
							},
						},
					},
				},
			},
		}

		travFnc.Body = []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: travPFnc.Name,
						TypeParams: []TypeExpression{
							TypeDef{Name: "A"},
							TypeDef{Name: "A"},
						},
					},
				},
			},
		}

		travPFnc.Body = []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: "CombiTraversal",
						TypeParams: []TypeExpression{
							TypeDef{Name: "ReturnMany"},
							TypeDef{Name: "ReadWrite"},
							TypeDef{Name: "Pure"},
							intType,
							tupleType,
							tupleTypeP,
							TypeDef{Name: "A"},
							TypeDef{Name: "B"},
						},
						Params: []Expression{
							iterFnc,
							lenFnc,
							modifyFnc,
							ixGetFnc,
							CallExpr{
								Func: "IxMatchComparable",
								TypeParams: []TypeExpression{
									intType,
								},
							},
							exprDef,
						},
					},
				},
			},
		}

		fd.Funcs = append(fd.Funcs, travFnc)
		fd.Funcs = append(fd.Funcs, travPFnc)
	}
}

func genTupleCol(fd *FileDef) {

	for tuple := 1; tuple < len(tuplePartNames); tuple++ {

		travFnc := FuncDef{
			Docs: []string{
				fmt.Sprintf("T%vToCol returns an [Iso] that converts a [lo.Tuple%v] to a [Collection]", tuple+1, tuple+1),
				"",
				"Note: under modification if the collection contains more elements than the tuple then the additional elements are discarded. If the collection contains less elements than the tuple then the tuple elements will default uninitialized values.",
				"",
				fmt.Sprintf("See: T%vPCol for a polymorphic version", tuple+1),
			},
			Name: fmt.Sprintf("T%vToCol", tuple+1),
			TypeParams: []TypeExpression{
				TypeDef{Name: "A"},
			},
		}

		travPFnc := FuncDef{
			Docs: []string{
				fmt.Sprintf("T%vCToolP returns a polymorphic [Iso] that converts a [lo.Tuple%v] to a [Collection]", tuple+1, tuple+1),
				"",
				"Note: under modification if the collection contains more elements than the tuple then the additional elements are discarded. If the collection contains less elements than the tuple then the tuple elements will default uninitialized values.",
				"",
				fmt.Sprintf("See: T%vToCol for a non polymorphic version", tuple+1),
			},
			Name: fmt.Sprintf("T%vToColP", tuple+1),
			TypeParams: []TypeExpression{
				TypeDef{Name: "A"},
				TypeDef{Name: "B"},
			},
		}

		tupleType := TypeDef{Name: fmt.Sprintf("lo.Tuple%v", tuple+1)}
		tupleTypeP := TypeDef{Name: fmt.Sprintf("lo.Tuple%v", tuple+1)}
		var elementPtrs []Expression
		var elements []Expression

		for i := 0; i <= tuple; i++ {
			tupleType.TypeParams = append(tupleType.TypeParams, TypeDef{Name: "A"})
			tupleTypeP.TypeParams = append(tupleTypeP.TypeParams, TypeDef{Name: "B"})
			elementPtrs = append(elementPtrs, "&source."+tuplePartNames[i])
			elements = append(elements, "&ret."+tuplePartNames[i])
		}

		travFnc.ReturnTypes = []TypeExpression{
			TypeDef{
				Name: "Optic",
				TypeParams: []TypeExpression{
					TypeDef{Name: "Void"},
					tupleType,
					tupleType,
					TypeDef{Name: "Collection[int,A,Pure]"},
					TypeDef{Name: "Collection[int,A,Pure]"},
					TypeDef{Name: "ReturnOne"},
					TypeDef{Name: "ReadWrite"},
					TypeDef{Name: "BiDir"},
					TypeDef{Name: "Pure"},
				},
			},
		}

		travPFnc.ReturnTypes = []TypeExpression{
			TypeDef{
				Name: "Optic",
				TypeParams: []TypeExpression{
					TypeDef{Name: "Void"},
					tupleType,
					tupleTypeP,
					TypeDef{Name: "Collection[int,A,Pure]"},
					TypeDef{Name: "Collection[int,B,Pure]"},
					TypeDef{Name: "ReturnOne"},
					TypeDef{Name: "ReadWrite"},
					TypeDef{Name: "BiDir"},
					TypeDef{Name: "Pure"},
				},
			},
		}

		intType := TypeDef{Name: "int"}
		ctxType := TypeDef{Name: "context.Context"}
		seqFType := TypeDef{
			Name: "Collection",
			TypeParams: []TypeExpression{
				intType,
				TypeDef{Name: "A"},
				TypeDef{Name: "Pure"},
			},
		}
		errorType := TypeDef{Name: "error"}

		getterParams := []Expression{"ctx"}
		getterParams = append(getterParams, elementPtrs...)

		getFnc := FuncDef{
			Params: []Param{
				{
					Name: "ctx",
					Type: ctxType,
				},
				{
					Name: "source",
					Type: tupleType,
				},
			},
			ReturnTypes: []TypeExpression{
				seqFType,
				errorType,
			},
			Body: []Statement{
				ReturnStmnt{
					Values: []Expression{
						CallExpr{
							Func:   "tnColGetter",
							Params: getterParams,
						},
						"nil",
					},
				},
			},
		}

		revGetParams := []Expression{"ctx", "focus"}
		revGetParams = append(revGetParams, elements...)

		revGetFnc := FuncDef{
			Params: []Param{
				{
					Name: "ctx",
					Type: ctxType,
				},
				{
					Name: "focus",
					Type: TypeDef{
						Name: "Collection[int, B, Pure]",
					},
				},
			},
			ReturnTypes: []TypeExpression{
				tupleTypeP,
				errorType,
			},
			Body: []Statement{
				VarDef{
					Name: "ret",
					Type: tupleTypeP,
				},
				AssignVar{
					Declare: true,
					Vars:    []string{"err"},
					Value: CallExpr{
						Func:   "tnColReverseGet",
						Params: revGetParams,
					},
				},
				ReturnStmnt{
					Values: []Expression{
						"ret",
						"err",
					},
				},
			},
		}

		exprDef := CallExpr{
			Func: "ExprDef",
			Params: []Expression{
				FuncDef{
					Params: []Param{
						{
							Name: "ot",
							Type: TypeDef{Name: "expr.OpticTypeExpr"},
						},
					},
					ReturnTypes: []TypeExpression{
						TypeDef{Name: "expr.OpticExpression"},
					},
					Body: []Statement{
						ReturnStmnt{
							Values: []Expression{
								StructExpr{
									Type: TypeDef{Name: "expr.ToCol"},
									Fields: []AssignField{
										{
											Name:  "OpticTypeExpr",
											Value: "ot",
										},
										{
											Name: "I",
											Value: CallExpr{
												Func: "reflect.TypeFor[int]",
											},
										},
										{
											Name: "A",
											Value: CallExpr{
												Func: "reflect.TypeFor[A]",
											},
										},
										{
											Name: "B",
											Value: CallExpr{
												Func: "reflect.TypeFor[B]",
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

		travFnc.Body = []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: travPFnc.Name,
						TypeParams: []TypeExpression{
							TypeDef{Name: "A"},
							TypeDef{Name: "A"},
						},
					},
				},
			},
		}

		travPFnc.Body = []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: "CombiIso",
						TypeParams: []TypeExpression{
							TypeDef{Name: "ReadWrite"},
							TypeDef{Name: "BiDir"},
							TypeDef{Name: "Pure"},
							tupleType,
							tupleTypeP,
							TypeDef{Name: "Collection[int, A, Pure]"},
							TypeDef{Name: "Collection[int, B, Pure]"},
						},
						Params: []Expression{
							getFnc,
							revGetFnc,
							exprDef,
						},
					},
				},
			},
		}

		fd.Funcs = append(fd.Funcs, travFnc)
		fd.Funcs = append(fd.Funcs, travPFnc)
	}
}

func genCol(fd *FileDef) {

	for tuple := 1; tuple < len(tuplePartNames); tuple++ {

		colFnc := FuncDef{
			Docs: []string{
				fmt.Sprintf("T%vColType returns a [CollectionType] wrapper for [lo.Tuple%v].", tuple+1, tuple+1),
				"",
				fmt.Sprintf("See: T%vColTypeP for a polymorphic version", tuple+1),
			},
			Name: fmt.Sprintf("T%vColType", tuple+1),
			TypeParams: []TypeExpression{
				TypeDef{Name: "A"},
			},
		}

		colPFnc := FuncDef{
			Docs: []string{
				fmt.Sprintf("T%vColTypeP returns a polymorphic [CollectionType] wrapper for [lo.Tuple%v].", tuple+1, tuple+1),
				"",
				fmt.Sprintf("See: T%vColType for a non polymorphic version", tuple+1),
			},
			Name: fmt.Sprintf("T%vColTypeP", tuple+1),
			TypeParams: []TypeExpression{
				TypeDef{Name: "A"},
				TypeDef{Name: "B"},
			},
		}

		tupleType := TypeDef{Name: fmt.Sprintf("lo.Tuple%v", tuple+1)}
		tupleTypeP := TypeDef{Name: fmt.Sprintf("lo.Tuple%v", tuple+1)}

		for i := 0; i <= tuple; i++ {
			tupleType.TypeParams = append(tupleType.TypeParams, TypeDef{Name: "A"})
			tupleTypeP.TypeParams = append(tupleTypeP.TypeParams, TypeDef{Name: "B"})
		}

		colFnc.ReturnTypes = []TypeExpression{
			TypeDef{
				Name: "CollectionType",
				TypeParams: []TypeExpression{
					TypeDef{Name: "int"},
					tupleType,
					tupleType,
					TypeDef{Name: "A"},
					TypeDef{Name: "A"},
					TypeDef{Name: "Pure"},
				},
			},
		}

		colPFnc.ReturnTypes = []TypeExpression{
			TypeDef{
				Name: "CollectionType",
				TypeParams: []TypeExpression{
					TypeDef{Name: "int"},
					tupleType,
					tupleTypeP,
					TypeDef{Name: "A"},
					TypeDef{Name: "B"},
					TypeDef{Name: "Pure"},
				},
			},
		}

		colFnc.Body = []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: colPFnc.Name,
						TypeParams: []TypeExpression{
							TypeDef{Name: "A"},
							TypeDef{Name: "A"},
						},
					},
				},
			},
		}

		colPFnc.Body = []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: "ColTypeP",
						Params: []Expression{
							CallExpr{
								Func: fmt.Sprintf("T%vToColP", tuple+1),
								TypeParams: []TypeExpression{
									TypeDef{Name: "A"},
									TypeDef{Name: "B"},
								},
							},
							CallExpr{
								Func: "AsReverseGet",
								Params: []Expression{
									CallExpr{
										Func: fmt.Sprintf("T%vToColP", tuple+1),
										TypeParams: []TypeExpression{
											TypeDef{Name: "B"},
											TypeDef{Name: "A"},
										},
									},
								},
							},
							CallExpr{
								Func: fmt.Sprintf("TraverseT%vP", tuple+1),
								TypeParams: []TypeExpression{
									TypeDef{Name: "A"},
									TypeDef{Name: "B"},
								},
							},
						},
					},
				},
			},
		}

		fd.Funcs = append(fd.Funcs, colFnc)
		fd.Funcs = append(fd.Funcs, colPFnc)
	}
}

func genOf(fd *FileDef) {

	for tuple := 1; tuple < len(tuplePartNames); tuple++ {

		tnOfFnc := FuncDef{
			Docs: []string{
				fmt.Sprintf("The T%vOf Combinator constructs a [lo.Tuple%v] whose elements are the focuses of the given [Optic]s", tuple+1, tuple+1),
				"",
				"Note: The number of focused tuples is limited by the optic that focuses the least elements.",
			},
			Name: fmt.Sprintf("T%vOf", tuple+1),
		}

		for i := 0; i <= tuple; i++ {
			tnOfFnc.TypeParams = append(tnOfFnc.TypeParams, TypeDef{
				Name: fmt.Sprintf("I%v", i),
			})
		}

		tnOfFnc.TypeParams = append(tnOfFnc.TypeParams, TypeDef{
			Name: "S",
		})

		for i := 0; i <= tuple; i++ {
			tnOfFnc.TypeParams = append(tnOfFnc.TypeParams, TypeDef{
				Name: fmt.Sprintf("A%v", i),
			})
		}

		for i := 0; i <= tuple; i++ {
			tnOfFnc.TypeParams = append(tnOfFnc.TypeParams, TypeDef{
				Name: fmt.Sprintf("RET%v", i),
			})
			tnOfFnc.TypeParams = append(tnOfFnc.TypeParams, TypeDef{
				Name: fmt.Sprintf("RW%v", i),
			})
			tnOfFnc.TypeParams = append(tnOfFnc.TypeParams, TypeDef{
				Name: fmt.Sprintf("DIR%v", i),
			})
			tnOfFnc.TypeParams = append(tnOfFnc.TypeParams, TypeDef{
				Name: fmt.Sprintf("ERR%v", i),
			})
		}

		retI := TypeDef{
			Name: fmt.Sprintf("lo.Tuple%v", tuple+1),
		}

		retA := TypeDef{
			Name: fmt.Sprintf("lo.Tuple%v", tuple+1),
		}

		var composeList []*ComposeTree

		for i := 0; i <= tuple; i++ {

			retI.TypeParams = append(retI.TypeParams, TypeDef{
				Name: fmt.Sprintf("I%v", i),
			})

			retA.TypeParams = append(retA.TypeParams, TypeDef{
				Name: fmt.Sprintf("A%v", i),
			})

			tnOfFnc.Params = append(tnOfFnc.Params, Param{
				Name: fmt.Sprintf("o%v", i),
				Type: TypeDef{
					Name: "Optic",
					TypeParams: []TypeExpression{
						TypeDef{Name: fmt.Sprintf("I%v", i)},
						TypeDef{Name: "S"},
						TypeDef{Name: "S"},
						TypeDef{Name: fmt.Sprintf("A%v", i)},
						TypeDef{Name: fmt.Sprintf("A%v", i)},
						TypeDef{Name: fmt.Sprintf("RET%v", i)},
						TypeDef{Name: fmt.Sprintf("RW%v", i)},
						TypeDef{Name: fmt.Sprintf("DIR%v", i)},
						TypeDef{Name: fmt.Sprintf("ERR%v", i)},
					},
				},
			})

			composeList = append(composeList, &ComposeTree{
				ParamNum: i,
				IsLeaf:   true,
			})

		}

		composeTree := BuildComposeTree(composeList)

		tnOfFnc.ReturnTypes = []TypeExpression{
			TypeDef{
				Name: "Optic",
				TypeParams: []TypeExpression{
					retI,
					TypeDef{Name: "S"},
					TypeDef{Name: "S"},
					retA,
					retA,
					BuildCompositionTree("RET", composeTree),
					BuildCompositionTree("RW", composeTree),
					TypeDef{Name: "UniDir"},
					BuildCompositionTree("ERR", composeTree),
				},
			},
		}

		sourceType := TypeDef{Name: "S"}
		ctxType := TypeDef{Name: "context.Context"}
		errType := TypeDef{Name: "error"}

		iterFunc := FuncDef{
			Params: []Param{
				{
					Name: "ctx",
					Type: ctxType,
				},
				{
					Name: "source",
					Type: sourceType,
				},
			},
			ReturnTypes: []TypeExpression{
				TypeDef{
					Name: "SeqIE",
					TypeParams: []TypeExpression{
						retI,
						retA,
					},
				},
			},
		}

		iterSeqFunc := FuncDef{
			Params: []Param{
				{
					Name: "yield",
					Type: FuncDef{
						Params: []Param{
							{
								Name: "focusHello",
								Type: TypeDef{
									Name: "ValueIE",
									TypeParams: []TypeExpression{
										retI,
										retA,
									},
								},
							},
						},
						ReturnTypes: []TypeExpression{
							TypeDef{Name: "bool"},
						},
					},
				},
			},
		}

		onRetI := CallExpr{
			Func: fmt.Sprintf("lo.T%v", tuple+1),
		}
		onRetA := CallExpr{
			Func: fmt.Sprintf("lo.T%v", tuple+1),
		}

		for i := 1; i <= tuple; i++ {
			iterSeqFunc.Body = append(iterSeqFunc.Body, AssignVar{
				Declare: true,
				Vars:    []string{fmt.Sprintf("next%v", i), fmt.Sprintf("stop%v", i)},
				Value: CallExpr{
					Func: "iter.Pull",
					Params: []Expression{
						CallExpr{
							Func: "iter.Seq",
							TypeParams: []TypeExpression{
								TypeDef{
									Name: "ValueIE",
									TypeParams: []TypeExpression{
										TypeDef{Name: fmt.Sprintf("I%v", i)},
										TypeDef{Name: fmt.Sprintf("A%v", i)},
									},
								},
							},
							Params: []Expression{
								CallExpr{
									Func: MethodCallExpr{
										Receiver: fmt.Sprintf("o%v", i),
										Name:     "AsIter",
									},
									Params: []Expression{
										"ctx",
										"source",
									},
								},
							},
						},
					},
				},
			})
			iterSeqFunc.Body = append(iterSeqFunc.Body, CallExpr{
				Func: fmt.Sprintf("defer stop%v", i),
			})

			onRetI.Params = append(onRetI.Params, fmt.Sprintf("i%v", i))
			onRetA.Params = append(onRetA.Params, fmt.Sprintf("a%v", i))
		}

		iterYieldFunc := FuncDef{
			Params: []Param{
				{
					Name: "val",
					Type: TypeDef{
						Name: "ValueIE",
						TypeParams: []TypeExpression{
							TypeDef{Name: "I0"},
							TypeDef{Name: "A0"},
						},
					},
				},
			},
			ReturnTypes: []TypeExpression{
				TypeDef{Name: "bool"},
			},
		}

		iterYieldFuncErrorCheck := IfStmnt{
			Condition: BinaryExpr{
				Left:  "err",
				Op:    "!=",
				Right: "nil",
			},
			OnTrue: []Statement{
				VarDef{
					Name: "i",
					Type: retI,
				},
				VarDef{
					Name: "a",
					Type: retA,
				},
				ReturnStmnt{
					Values: []Expression{
						CallExpr{
							Func: "yield",
							Params: []Expression{
								CallExpr{
									Func: "ValIE",
									Params: []Expression{
										"i",
										"a",
										"err",
									},
								},
							},
						},
					},
				},
			},
		}

		iterYieldFunc.Body = append(iterYieldFunc.Body, AssignVar{
			Declare: true,
			Vars:    []string{"i0", "a0", "err"},
			Value:   MethodCallExpr{Receiver: "val", Name: "Get"},
		})
		iterYieldFunc.Body = append(iterYieldFunc.Body, iterYieldFuncErrorCheck)

		iParams := []Expression{"i0"}
		aParams := []Expression{"a0"}

		for i := 1; i <= tuple; i++ {
			iterYieldFunc.Body = append(iterYieldFunc.Body,
				AssignVar{
					Declare: true,
					Vars: []string{
						fmt.Sprintf("v%v", i),
						"ok",
					},
					Value: CallExpr{
						Func: fmt.Sprintf("next%v", i),
					},
				},
			)
			iterYieldFunc.Body = append(iterYieldFunc.Body,
				IfStmnt{
					Condition: "!ok",
					OnTrue: []Statement{
						ReturnStmnt{
							Values: []Expression{
								"false",
							},
						},
					},
				},
			)
			iterYieldFunc.Body = append(iterYieldFunc.Body, AssignVar{
				Declare: true,
				Vars: []string{
					fmt.Sprintf("i%v", i),
					fmt.Sprintf("a%v", i),
					"err",
				},
				Value: MethodCallExpr{Receiver: fmt.Sprintf("v%v", i), Name: "Get"},
			})
			iterYieldFunc.Body = append(iterYieldFunc.Body, iterYieldFuncErrorCheck)

			iParams = append(iParams, fmt.Sprintf("i%v", i))
			aParams = append(aParams, fmt.Sprintf("a%v", i))
		}

		iterYieldFunc.Body = append(iterYieldFunc.Body,
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: "yield",
						Params: []Expression{
							CallExpr{
								Func: "ValIE",
								Params: []Expression{
									CallExpr{
										Func:   fmt.Sprintf("lo.T%v", tuple+1),
										Params: iParams,
									},
									CallExpr{
										Func:   fmt.Sprintf("lo.T%v", tuple+1),
										Params: aParams,
									},
									"nil",
								},
							},
						},
					},
				},
			},
		)

		iterSeqFunc.Body = append(iterSeqFunc.Body,
			CallExpr{
				Func: CallExpr{
					Func: MethodCallExpr{
						Receiver: "o0",
						Name:     "AsIter",
					},
					Params: []Expression{
						"ctx",
						"source",
					},
				},
				Params: []Expression{
					iterYieldFunc,
				},
			},
		)

		iterFunc.Body = append(iterFunc.Body, ReturnStmnt{
			Values: []Expression{
				iterSeqFunc,
			},
		})

		lengthGetterFunc := FuncDef{
			Params: []Param{
				{
					Name: "ctx",
					Type: ctxType,
				},
				{
					Name: "source",
					Type: sourceType,
				},
			},
			ReturnTypes: []TypeExpression{
				TypeDef{Name: "int"},
				errType,
			},
		}

		var lenRet []Expression

		for i := 0; i <= tuple; i++ {
			lengthGetterFunc.Body = append(lengthGetterFunc.Body,
				AssignVar{
					Declare: true,
					Vars:    []string{fmt.Sprintf("l%v", i), "err"},
					Value: CallExpr{
						Func: MethodCallExpr{
							Receiver: fmt.Sprintf("o%v", i),
							Name:     "AsLengthGetter",
						},
						Params: []Expression{
							"ctx",
							"source",
						},
					},
				},
				IfStmnt{
					Condition: BinaryExpr{
						Left:  "err",
						Op:    "!=",
						Right: "nil",
					},
					OnTrue: []Statement{
						ReturnStmnt{
							Values: []Expression{
								"0",
								"err",
							},
						},
					},
				},
			)

			lenRet = append(lenRet, fmt.Sprintf("l%v", i))
		}

		lengthGetterFunc.Body = append(lengthGetterFunc.Body,
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func:   "min",
						Params: lenRet,
					},
					"nil",
				},
			},
		)

		modifyFunc := FuncDef{
			Params: []Param{
				{
					Name: "ctx",
					Type: TypeDef{Name: "context.Context"},
				},
				{
					Name: "fmap",
					Type: FuncDef{
						Params: []Param{
							{
								Name: "index",
								Type: retI,
							},
							{
								Name: "focus",
								Type: retA,
							},
						},
						ReturnTypes: []TypeExpression{
							retA,
							errType,
						},
					},
				},
				{
					Name: "source",
					Type: TypeDef{Name: "S"},
				},
			},
			ReturnTypes: []TypeExpression{
				TypeDef{Name: "S"},
				errType,
			},
		}

		for i := 1; i <= tuple; i++ {
			modifyFunc.Body = append(modifyFunc.Body, AssignVar{
				Declare: true,
				Vars:    []string{fmt.Sprintf("next%v", i), fmt.Sprintf("stop%v", i)},
				Value: CallExpr{
					Func: "iter.Pull",
					Params: []Expression{
						CallExpr{
							Func: "iter.Seq",
							TypeParams: []TypeExpression{
								TypeDef{
									Name: "ValueIE",
									TypeParams: []TypeExpression{
										TypeDef{Name: fmt.Sprintf("I%v", i)},
										TypeDef{Name: fmt.Sprintf("A%v", i)},
									},
								},
							},
							Params: []Expression{
								CallExpr{
									Func: MethodCallExpr{
										Receiver: fmt.Sprintf("o%v", i),
										Name:     "AsIter",
									},
									Params: []Expression{
										"ctx",
										"source",
									},
								},
							},
						},
					},
				},
			})
			modifyFunc.Body = append(modifyFunc.Body, CallExpr{
				Func: fmt.Sprintf("defer stop%v", i),
			})
			modifyFunc.Body = append(modifyFunc.Body,
				VarDef{
					Name: fmt.Sprintf("mapping%v", i),
					Type: TypeDef{Name: fmt.Sprintf("[]A%v", i)},
				},
			)
		}

		modifyFmapFunc := FuncDef{
			Params: []Param{
				{
					Name: "index",
					Type: TypeDef{Name: "I0"},
				},
				{
					Name: "focus",
					Type: TypeDef{Name: "A0"},
				},
			},
			ReturnTypes: []TypeExpression{
				TypeDef{Name: "A0"},
				TypeDef{Name: "error"},
			},
		}

		ifNotOk := IfStmnt{
			Condition: "!ok",
			OnTrue: []Statement{
				ReturnStmnt{
					Values: []Expression{
						"focus",
						"nil",
					},
				},
			},
		}

		ifErrModify := IfStmnt{
			Condition: BinaryExpr{
				Left:  "err",
				Op:    "!=",
				Right: "nil",
			},
			OnTrue: []Statement{
				VarDef{
					Name: "a",
					Type: TypeDef{Name: "A0"},
				},
				ReturnStmnt{
					Values: []Expression{
						"a",
						"err",
					},
				},
			},
		}

		fmapIParams := []Expression{"index"}
		fmapAParams := []Expression{"focus"}

		for i := 1; i <= tuple; i++ {

			fmapIParams = append(fmapIParams, fmt.Sprintf("i%v", i))
			fmapAParams = append(fmapAParams, fmt.Sprintf("a%v", i))

			modifyFmapFunc.Body = append(modifyFmapFunc.Body,
				AssignVar{
					Declare: true,
					Vars:    []string{fmt.Sprintf("v%v", i), "ok"},
					Value: CallExpr{
						Func: fmt.Sprintf("next%v", i),
					},
				},
				ifNotOk,
				AssignVar{
					Declare: true,
					Vars:    []string{fmt.Sprintf("i%v", i), fmt.Sprintf("a%v", i), "err"},
					Value: MethodCallExpr{
						Receiver: fmt.Sprintf("v%v", i),
						Name:     "Get",
					},
				},
				ifErrModify,
			)
		}

		modifyFmapFunc.Body = append(modifyFmapFunc.Body,
			AssignVar{
				Declare: true,
				Vars:    []string{"mapped", "err"},
				Value: CallExpr{
					Func: "fmap",
					Params: []Expression{
						CallExpr{
							Func:   fmt.Sprintf("lo.T%v", tuple+1),
							Params: fmapIParams,
						},
						CallExpr{
							Func:   fmt.Sprintf("lo.T%v", tuple+1),
							Params: fmapAParams,
						},
					},
				},
			},
			ifErrModify,
		)

		for i := 1; i <= tuple; i++ {
			modifyFmapFunc.Body = append(modifyFmapFunc.Body,
				AssignVar{
					Declare: false,
					Vars:    []string{fmt.Sprintf("mapping%v", i)},
					Value: CallExpr{
						Func: "append",
						Params: []Expression{
							fmt.Sprintf("mapping%v", i),
							fmt.Sprintf("mapped.%v", tuplePartNames[i]),
						},
					},
				},
			)
		}

		modifyFmapFunc.Body = append(modifyFmapFunc.Body,
			ReturnStmnt{
				Values: []Expression{
					"mapped.A",
					"err",
				},
			},
		)

		modifyFunc.Body = append(modifyFunc.Body,
			AssignVar{
				Declare: true,
				Vars:    []string{"ret", "err"},
				Value: CallExpr{
					Func: MethodCallExpr{
						Receiver: "o0",
						Name:     "AsModify",
					},
					Params: []Expression{
						"ctx",
						modifyFmapFunc,
						"source",
					},
				},
			},
			IfStmnt{
				Condition: BinaryExpr{
					Left:  "err",
					Op:    "!=",
					Right: "nil",
				},
				OnTrue: []Statement{
					VarDef{
						Name: "s",
						Type: TypeDef{Name: "S"},
					},
					ReturnStmnt{
						Values: []Expression{
							"s",
							"err",
						},
					},
				},
			},
		)

		modifyFunc.Body = append(modifyFunc.Body,
			AssignVar{
				Declare: true,
				Vars:    []string{"i"},
				Value:   "0",
			},
		)

		for i := 1; i <= tuple; i++ {
			fmapMappedFunc := FuncDef{
				Params: []Param{
					{
						Name: "index",
						Type: TypeDef{Name: fmt.Sprintf("I%v", i)},
					},
					{
						Name: "focus",
						Type: TypeDef{Name: fmt.Sprintf("A%v", i)},
					},
				},
				ReturnTypes: []TypeExpression{
					TypeDef{Name: fmt.Sprintf("A%v", i)},
					TypeDef{Name: "error"},
				},
			}

			mapping := fmt.Sprintf("mapping%v", i)
			fmapMappedFunc.Body = append(fmapMappedFunc.Body,
				IfStmnt{
					Condition: BinaryExpr{
						Left: "i",
						Op:   ">=",
						Right: CallExpr{
							Func:   "len",
							Params: []Expression{mapping},
						},
					},
					OnTrue: []Statement{
						ReturnStmnt{
							Values: []Expression{
								"focus",
								"nil",
							},
						},
					},
				},
				ReturnStmnt{
					Values: []Expression{
						mapping + "[i]",
						"nil",
					},
				},
			)

			modifyFunc.Body = append(modifyFunc.Body,
				AssignVar{
					Declare: false,
					Vars:    []string{"ret", "err"},
					Value: CallExpr{
						Func: MethodCallExpr{
							Receiver: fmt.Sprintf("o%v", i),
							Name:     "AsModify",
						},
						Params: []Expression{
							"ctx",
							fmapMappedFunc,
							"ret",
						},
					},
				},
				IfStmnt{
					Condition: BinaryExpr{
						Left:  "err",
						Op:    "!=",
						Right: "nil",
					},
					OnTrue: []Statement{
						VarDef{
							Name: "s",
							Type: TypeDef{Name: "S"},
						},
						ReturnStmnt{
							Values: []Expression{
								"s",
								"err",
							},
						},
					},
				},
			)
		}

		modifyFunc.Body = append(modifyFunc.Body,
			ReturnStmnt{
				Values: []Expression{
					"ret",
					"nil",
				},
			},
		)

		ixMatchFunc := FuncDef{
			Params: []Param{
				{
					Name: "a",
					Type: retI,
				},
				{
					Name: "b",
					Type: retI,
				},
			},
			ReturnTypes: []TypeExpression{
				TypeDef{Name: "bool"},
			},
		}

		for i := 0; i <= tuple; i++ {
			ixMatchFunc.Body = append(ixMatchFunc.Body,
				IfStmnt{
					Condition: BinaryExpr{
						Left: CallExpr{
							Func: MethodCallExpr{
								Receiver: fmt.Sprintf("o%v", i),
								Name:     "AsIxMatch",
							},
							Params: []Expression{
								"a." + tuplePartNames[i],
								"b." + tuplePartNames[i],
							},
						},
						Op:    "!=",
						Right: "true",
					},
					OnTrue: []Statement{
						ReturnStmnt{
							Values: []Expression{
								"false",
							},
						},
					},
				},
			)
		}

		ixMatchFunc.Body = append(ixMatchFunc.Body,
			ReturnStmnt{
				Values: []Expression{
					"true",
				},
			},
		)

		exprDefElements := SliceExpr{
			Type: TypeDef{Name: "expr.OpticExpression"},
		}

		for i := 0; i <= tuple; i++ {
			exprDefElements.Values = append(exprDefElements.Values, MethodCallExpr{
				Receiver: fmt.Sprintf("o%v", i),
				Name:     "AsExpr",
			})
		}

		exprDefParams := []Expression{
			FuncDef{
				Params: []Param{
					{
						Name: "ot",
						Type: TypeDef{Name: "expr.OpticTypeExpr"},
					},
				},
				ReturnTypes: []TypeExpression{
					TypeDef{Name: "expr.OpticExpression"},
				},
				Body: []Statement{
					ReturnStmnt{
						Values: []Expression{
							StructExpr{
								Type: TypeDef{Name: "expr.TupleOf"},
								Fields: []AssignField{
									{
										Name:  "OpticTypeExpr",
										Value: "ot",
									},
									{
										Name:  "Elements",
										Value: exprDefElements,
									},
								},
							},
						},
					},
				},
			},
		}

		for i := 0; i <= tuple; i++ {
			exprDefParams = append(exprDefParams, fmt.Sprintf("o%v", i))
		}

		exprDef := CallExpr{
			Func:   "ExprDef",
			Params: exprDefParams,
		}

		tnOfFnc.Body = []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: "CombiTraversal",
						TypeParams: []TypeExpression{
							BuildCompositionTree("RET", composeTree),
							BuildCompositionTree("RW", composeTree),
							BuildCompositionTree("ERR", composeTree),
							retI,
							TypeDef{Name: "S"},
							TypeDef{Name: "S"},
							retA,
							retA,
						},
						Params: []Expression{
							iterFunc,
							lengthGetterFunc,
							modifyFunc,
							"nil",
							ixMatchFunc,
							exprDef,
						},
					},
				},
			},
		}

		fd.Funcs = append(fd.Funcs, tnOfFnc)
	}

}

func genDelve(fd *FileDef) {

	for tuple := 1; tuple < len(tuplePartNames); tuple++ {

		delveFnc := FuncDef{
			Docs: []string{
				fmt.Sprintf("The DelveT%v combinator focuses each element of a [lo.Tuple%v] using the given [Optic]", tuple+1, tuple+1),
			},
			Name: fmt.Sprintf("DelveT%v", tuple+1),
			Params: []Param{
				Param{
					Name: "o",
					Type: TypeDef{
						Name: "Optic",
						TypeParams: []TypeExpression{
							TypeDef{Name: "I"},
							TypeDef{Name: "S"},
							TypeDef{Name: "S"},
							TypeDef{Name: "A"},
							TypeDef{Name: "A"},
							TypeDef{Name: "RET"},
							TypeDef{Name: "RW"},
							TypeDef{Name: "DIR"},
							TypeDef{Name: "ERR"},
						},
					},
				},
			},
			TypeParams: []TypeExpression{
				TypeDef{Name: "I"},
				TypeDef{Name: "S"},
				TypeDef{Name: "A"},
				TypeDef{Name: "RET"},
				TypeDef{Name: "RW"},
				TypeDef{Name: "DIR"},
				TypeDef{Name: "ERR"},
			},
		}

		var tnOfTypeParams []TypeExpression
		for i := 0; i <= tuple; i++ {
			tnOfTypeParams = append(tnOfTypeParams, TypeDef{Name: "S"})
		}

		tupleTypeI := TypeDef{Name: fmt.Sprintf("lo.Tuple%v", tuple+1)}
		tupleTypeS := TypeDef{Name: fmt.Sprintf("lo.Tuple%v", tuple+1)}
		tupleTypeA := TypeDef{Name: fmt.Sprintf("lo.Tuple%v", tuple+1)}
		var tnOf []Expression

		for i := 0; i <= tuple; i++ {

			tupleTypeI.TypeParams = append(tupleTypeI.TypeParams, TypeDef{Name: "I"})
			tupleTypeS.TypeParams = append(tupleTypeS.TypeParams, TypeDef{Name: "S"})
			tupleTypeA.TypeParams = append(tupleTypeA.TypeParams, TypeDef{Name: "A"})

			tnOf = append(tnOf, CallExpr{
				Func: "Compose",
				Params: []Expression{
					CallExpr{
						Func:       fmt.Sprintf("T%v%v", tuple+1, tuplePartNames[i]),
						TypeParams: tnOfTypeParams,
					},
					"o",
				},
			})
		}

		delveFnc.ReturnTypes = []TypeExpression{
			TypeDef{
				Name: "Optic",
				TypeParams: []TypeExpression{
					tupleTypeI,
					tupleTypeS,
					tupleTypeS,
					tupleTypeA,
					tupleTypeA,
					TypeDef{Name: "RET"},
					TypeDef{Name: "RW"},
					TypeDef{Name: "UniDir"},
					TypeDef{Name: "ERR"},
				},
			},
		}

		delveFnc.Body = []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: "UnsafeReconstrain",
						TypeParams: []TypeExpression{
							TypeDef{Name: "RET"},
							TypeDef{Name: "RW"},
							TypeDef{Name: "UniDir"},
							TypeDef{Name: "ERR"},
						},
						Params: []Expression{
							CallExpr{
								Func:   fmt.Sprintf("T%vOf", tuple+1),
								Params: tnOf,
							},
						},
					},
				},
			},
		}

		fd.Funcs = append(fd.Funcs, delveFnc)
	}
}

func main() {

	var genFileName = "../../../tuple_generated.go"

	var fd FileDef
	fd.Package = "optic"

	fd.Imports = []string{
		"context",
		"iter",
		"reflect",
		"github.com/samber/lo",
		"github.com/spearson78/go-optic/expr",
	}

	genTupleTraverse(&fd)
	genTupleDup(&fd)
	genTupleElements(&fd)
	genTupleCol(&fd)
	genCol(&fd)
	genOf(&fd)
	genDelve(&fd)

	w, err := os.Create(genFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	WriteGoFile(w, &fd)

}
