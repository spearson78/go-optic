package main

import (
	"slices"
	"strings"

	. "github.com/spearson78/go-optic/internal/codegen"
)

func typeNameString(typeName string, typePackage []string, typeParam []string) string {
	if len(typePackage) == 0 && len(typeParam) == 0 {
		return typeName
	}

	var sb strings.Builder

	if len(typePackage) > 0 {
		sb.WriteString(strings.Join(typePackage, "."))
		sb.WriteString(".")
	}

	sb.WriteString(typeName)

	if len(typeParam) > 0 {
		sb.WriteString("[")
		for i := len(typeParam) - 1; i >= 0; i-- {
			sb.WriteString(typeParam[i])
			if i != 0 {
				sb.WriteString(".")
			}
		}

		sb.WriteString("]")
	}

	str := sb.String()

	return str

}

func makePublic(name string) string {
	return strings.Title(name)
}

func getStructType(name string) string {
	if *packagePrefix == "" {
		return name
	}

	return *packagePrefix + "." + name
}

func isPrimitive(typeName string, typePackage []string, dm *DataModel) bool {

	if len(typePackage) == 0 {
		if alias, ok := dm.TypeAliases[typeName]; ok {
			return isPrimitive(alias, nil, dm)
		}

		_, ok := dm.Structs[typeName]
		if ok {
			return false
		}

		return true
	} else {
		//Treat anything in another package as a primitive
		return true
	}
}

func isComparable(typeName string, typePackage []string, dm *DataModel) bool {

	if len(typePackage) == 0 {

		if alias, ok := dm.TypeAliases[typeName]; ok {
			return isComparable(alias, nil, dm)
		}

		switch typeName {
		case "string":
			return true
		case "float64":
			return true
		case "float32":
			return true
		case "int64":
			return true
		case "uint64":
			return true
		case "int":
			return true
		case "uint":
			return true
		case "bool":
			return true
		case "int32":
			return true
		case "uint32":
			return true
		case "uint8":
			return true
		case "int8":
			return true
		case "uint16":
			return true
		case "int16":
			return true
		case "complex":
			return true
		default:
			_, ok := dm.Comparable[typeName]
			return ok
		}
	} else {
		fullName := typePackage[0] + "." + typeName
		_, ok := dm.Comparable[fullName]
		return ok
	}
}

func isReal(typeName string, typePackage []string, typeAliases map[string]string) bool {

	if len(typePackage) == 0 {
		switch typeName {
		case "float64":
			return true
		case "float32":
			return true
		case "int64":
			return true
		case "uint64":
			return true
		case "int":
			return true
		case "uint":
			return true
		case "int32":
			return true
		case "uint32":
			return true
		case "uint8":
			return true
		case "int8":
			return true
		case "uint16":
			return true
		case "int16":
			return true
		}
	}

	if alias, ok := typeAliases[typeName]; ok {
		return isReal(alias, nil, typeAliases)
	}

	return false
}

func isOrdered(typeName string, typePackage []string, typeAliases map[string]string) bool {

	if len(typePackage) == 0 {
		switch typeName {
		case "string":
			return true
		case "float64":
			return true
		case "float32":
			return true
		case "int64":
			return true
		case "uint64":
			return true
		case "int":
			return true
		case "uint":
			return true
		case "int32":
			return true
		case "uint32":
			return true
		case "uint8":
			return true
		case "int8":
			return true
		case "uint16":
			return true
		case "int16":
			return true
		}
	}

	if alias, ok := typeAliases[typeName]; ok {
		return isOrdered(alias, nil, typeAliases)
	}

	return false
}

func DmToFd(dm *DataModel, packageName string, rootObjName string, combinators []Combinator) FileDef {

	fd := FileDef{
		Package: packageName,
		Imports: []string{
			"github.com/samber/mo",
			"github.com/spearson78/go-optic",
		},
	}

	oPrefix := strings.ToLower(rootObjName)
	upperOPrefix := rootObjName

	oStruct := StructDef{
		Name: oPrefix,
	}

	structNames := make([]string, 0, len(dm.Structs))
	userTypes := make(map[string]struct{}, len(dm.Structs))
	for _, str := range dm.Structs {
		if str.Name != "" {
			structNames = append(structNames, str.Name)
		}
		userTypes[str.Name] = struct{}{}
	}

	slices.Sort(structNames)

	for _, strName := range structNames {
		str := dm.Structs[strName]

		if str.TypeParam != "" {
			continue
		}

		buildLens(upperOPrefix, &fd, str, &oStruct, dm, combinators)
		buildSlice(upperOPrefix, &fd, str)
		buildMap(upperOPrefix, &fd, str)
		buildOpt(upperOPrefix, &fd, str, dm)
	}

	fd.Structs = append(fd.Structs, oStruct)

	fd.Vars = append(fd.Vars, VarDef{
		Name: upperOPrefix,
		Value: StructExpr{
			Type: TypeDef{Name: oPrefix},
		},
	})

	return fd
}

func genFieldLens(str *Struct, fld *Field) CallExpr {

	var lensType = "optic.FieldLens"

	return CallExpr{
		Func: lensType,
		Params: []Expression{
			GenFuncExpr{
				Params: []Param{
					{
						Name: "x",
						Type: Star{Type: structToTypeDef(str)},
					},
				},
				ReturnTypes: []TypeExpression{
					Star{Type: fieldToTypeDef(fld)},
				},
				Body: []Statement{
					ReturnStmnt{
						Values: []Expression{
							AddressExpr{
								Target: BinaryExpr{
									Op:    ".",
									Left:  "x",
									Right: fld.Name,
								},
							},
						},
					},
				},
			},
		},
	}
}

func genPtrFieldLens(str *Struct, fld *Field) CallExpr {

	var lensType = "optic.PtrFieldLens"

	return CallExpr{
		Func: lensType,
		Params: []Expression{
			GenFuncExpr{
				Params: []Param{
					{
						Name: "x",
						Type: Star{Type: structToTypeDef(str)},
					},
				},
				ReturnTypes: []TypeExpression{
					Star{Type: fieldToTypeDef(fld)},
				},
				Body: []Statement{
					ReturnStmnt{
						Values: []Expression{
							AddressExpr{
								Target: BinaryExpr{
									Op:    ".",
									Left:  "x",
									Right: fld.Name,
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildSlice(oPrefix string, fd *FileDef, str *Struct) {
	eStruct := genSliceStruct(str, "s", "optic.Optic")

	eStruct.Methods = append(eStruct.Methods, genTraverseSliceMethod(oPrefix, str))
	eStruct.Methods = append(eStruct.Methods, genIndexSliceMethod(oPrefix, str, "l"))

	fd.Structs = append(fd.Structs, eStruct)
}

func buildMap(oPrefix string, fd *FileDef, str *Struct) {
	mStruct := genMapStruct(str)

	mStruct.Methods = append(mStruct.Methods, genTraverseMapMethod(oPrefix, str))
	mStruct.Methods = append(mStruct.Methods, genIndexMapMethod(oPrefix, str, "l"))

	fd.Structs = append(fd.Structs, mStruct)
}

func genRetOneMethod(oPrefix string, str *Struct, prefix string, structTypePrefix string, asType string) MethodDef {

	params := []Expression{
		CallExpr{
			Func: "optic.Identity",
			TypeParams: []TypeExpression{
				TypeDef{
					Name: structTypePrefix + getStructType(str.Name),
				},
			},
		},
	}

	if asType != "" {
		params = []Expression{
			CallExpr{
				Func:   asType,
				Params: params,
			},
		}
	}

	return MethodDef{
		Name: str.Name,
		ReturnTypes: []TypeExpression{
			Star{
				Type: TypeDef{
					Name: prefix + str.Name,
					TypeParams: []TypeExpression{
						TypeDef{
							Name: "optic.Void",
						},
						TypeDef{
							Name: str.Name,
						},
						TypeDef{
							Name: str.Name,
						},
						TypeDef{
							Name: "optic.ReturnOne",
						},
						TypeDef{
							Name: "optic.ReadWrite",
						},
						TypeDef{
							Name: "optic.BiDir",
						},
						TypeDef{
							Name: "optic.Pure",
						},
					},
				},
			},
		},
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: oPrefix + str.Name + "Of",
						TypeParams: []TypeExpression{
							TypeDef{
								Name: "optic.Void",
							},
							TypeDef{
								Name: structTypePrefix + getStructType(str.Name),
							},
							TypeDef{
								Name: structTypePrefix + getStructType(str.Name),
							},
							TypeDef{
								Name: "optic.ReturnOne",
							},
						},
						Params: params,
					},
				},
			},
		},
	}
}

func genRetOneExpr(name string, prefix string, index string, dir string, optic Expression) Expression {

	return AddressExpr{
		Target: StructExpr{
			Type: TypeDef{
				Name: prefix + name,
				TypeParams: []TypeExpression{
					TypeDef{
						Name: index,
					},
					TypeDef{
						Name: "S",
					},
					TypeDef{
						Name: "T",
					},
					TypeDef{
						Name: "RET",
					},
					TypeDef{
						Name: "RW",
					},
					TypeDef{
						Name: dir,
					},
					TypeDef{
						Name: "ERR",
					},
				},
			},
			Fields: []AssignField{
				{
					Name:  "Optic",
					Value: optic,
				},
			},
		},
	}
}

func genRetManyExpr(name string, prefix string, index string, optic Expression) Expression {

	return AddressExpr{
		Target: StructExpr{
			Type: TypeDef{
				Name: prefix + name,
				TypeParams: []TypeExpression{
					TypeDef{
						Name: index,
					},
					TypeDef{
						Name: "S",
					},
					TypeDef{
						Name: "T",
					},
					TypeDef{
						Name: "optic.ReturnMany",
					},
					TypeDef{
						Name: "RW",
					},
					TypeDef{
						Name: "optic.UniDir",
					},
					TypeDef{
						Name: "ERR",
					},
				},
			},
			Fields: []AssignField{
				{
					Name:  "Optic",
					Value: optic,
				},
			},
		},
	}
}

func genMapRetManyExpr(name string, prefix string, index string, optic Expression) Expression {

	return AddressExpr{
		Target: StructExpr{
			Type: TypeDef{
				Name: prefix + name,
				TypeParams: []TypeExpression{
					TypeDef{
						Name: index,
					},
					TypeDef{
						Name: "I",
					},
					TypeDef{
						Name: "S",
					},
					TypeDef{
						Name: "T",
					},
					TypeDef{
						Name: "optic.ReturnMany",
					},
					TypeDef{
						Name: "RW",
					},
					TypeDef{
						Name: "optic.UniDir",
					},
					TypeDef{
						Name: "ERR",
					},
				},
			},
			Fields: []AssignField{
				{
					Name: "Optic",
					Value: EErrL(RetL(RwL(Ud(CallExpr{
						Func: "optic.ComposeLeft",
						Params: []Expression{
							optic,
							CallExpr{
								Func: "optic.MapCol",
								TypeParams: []TypeExpression{
									TypeDef{
										Name: index,
									},
									TypeDef{
										Name: name,
									},
								},
							},
						},
					})))),
				},
				{
					Name:  "o",
					Value: optic,
				},
			},
		},
	}
}

func genRetOneFunc(oPrefix string, str *Struct, prefix string, targetWerapper string) FuncDef {

	targetName := str.Name
	if targetWerapper != "" {
		targetName = targetWerapper + "[" + targetName + "]"
	}

	return FuncDef{
		Name: oPrefix + str.Name + "Of",
		TypeParams: []TypeExpression{
			TypeDef{
				Name: "I",
			},
			TypeDef{
				Name:       "S",
				Constraint: *sourceConstraint,
			},
			TypeDef{
				Name:       "T",
				Constraint: *sourceConstraint,
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
				Name: prefix,
				Type: TypeDef{
					Name: "optic.Optic",
					TypeParams: []TypeExpression{
						TypeDef{
							Name: "I",
						},
						TypeDef{
							Name: "S",
						},
						TypeDef{
							Name: "T",
						},
						TypeDef{
							Name: targetName,
						},
						TypeDef{
							Name: targetName,
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
					Name: prefix + str.Name,
					TypeParams: []TypeExpression{
						TypeDef{
							Name: "I",
						},
						TypeDef{
							Name: "S",
						},
						TypeDef{
							Name: "T",
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
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{genRetOneExpr(str.Name, prefix, "I", "DIR", prefix)},
			},
		},
	}
}

func genSRetMany(fld *Field, optic Expression, ret string) Expression {

	prefix := "s"
	opticType := "optic.Optic"

	opticFieldName := strings.SplitN(opticType, ".", 2)[1]

	return AddressExpr{
		Target: StructExpr{
			Type: TypeDef{
				Name: prefix + fld.TypeName,
				TypeParams: []TypeExpression{
					TypeDef{
						Name: "I",
					},
					TypeDef{
						Name: "S",
					},
					TypeDef{
						Name: "T",
					},
					TypeDef{
						Name: ret,
					},
					TypeDef{
						Name: "RW",
					},
					TypeDef{
						Name: "optic.UniDir",
					},
					TypeDef{
						Name: "ERR",
					},
				},
			},
			Fields: []AssignField{
				{
					Name: opticFieldName,
					Value: RetL(RwL(Ud(EErrL(ComposeLeft(
						optic,
						CallExpr{
							Func: "optic.SliceToCol",
							TypeParams: []TypeExpression{
								TypeDef{
									Name: fld.TypeName,
								},
							},
						},
					))))),
				},
				{
					Name:  "o",
					Value: optic,
				},
			},
		},
	}
}

func genMRetMany(fld *Field, optic Expression) Expression {

	prefix := "m"

	return AddressExpr{
		Target: StructExpr{
			Type: TypeDef{
				Name: prefix + fld.TypeName,
				TypeParams: []TypeExpression{
					TypeDef{
						Name: fld.MapKey,
					},
					TypeDef{
						Name: "I",
					},
					TypeDef{
						Name: "S",
					},
					TypeDef{
						Name: "T",
					},
					TypeDef{
						Name: "RET",
					},
					TypeDef{
						Name: "RW",
					},
					TypeDef{
						Name: "optic.UniDir",
					},
					TypeDef{
						Name: "ERR",
					},
				},
			},
			Fields: []AssignField{
				{
					Name: "Optic",
					Value: RetL(RwL(Ud(EErrL(ComposeLeft(
						optic,
						CallExpr{
							Func: "optic.MapToCol",
							TypeParams: []TypeExpression{
								TypeDef{
									Name: fld.MapKey,
								},
								TypeDef{
									Name: fld.TypeName,
								},
							},
						},
					))))),
				},
				{
					Name:  "o",
					Value: optic,
				},
			},
		},
	}
}

func buildOpt(oPrefix string, fd *FileDef, str *Struct, dm *DataModel) {
	oStruct := genStructDefRetOne(str, "o", "", true)

	for _, fld := range str.Fields {

		if isPrimitive(fld.TypeName, fld.TypePackage, dm) {
			if fld.Map {
				oStruct.Methods = append(oStruct.Methods, genOptMPrimitiveMethod(str, &fld, "Map"))
			} else if fld.Slice {
				oStruct.Methods = append(oStruct.Methods, genOptSPrimitiveMethod(str, &fld, "Slice"))
			} else {
				if fld.Pointer {
					oStruct.Methods = append(oStruct.Methods, genOptPrimitiveMethod(str, &fld, "optic.Optic", ""))
				} else {
					oStruct.Methods = append(oStruct.Methods, genOptPrimitiveMethod(str, &fld, "optic.Optic", ""))
				}
			}
		} else if fld.Slice {
			oStruct.Methods = append(oStruct.Methods, genOptSliceMethod(str, &fld))
		} else if fld.Map {
			oStruct.Methods = append(oStruct.Methods, genOptMapMethod(str, &fld))
		} else {
			if fld.Pointer {
				oStruct.Methods = append(oStruct.Methods, genOptMethod(str, &fld, "o"))
			} else {

				oStruct.Methods = append(oStruct.Methods, genOptMethod(str, &fld, "l"))
			}
		}
	}

	if *sourceConstraint == "" {

		oStruct.Methods = append(oStruct.Methods, MethodDef{
			Name: "Some",
			ReturnTypes: []TypeExpression{
				TypeDef{
					Name: "*l" + str.Name,
					TypeParams: []TypeExpression{
						TypeDef{
							Name: "optic.Void",
						},
						TypeDef{
							Name: "mo.Option[" + str.Name + "]",
						},
						TypeDef{
							Name: "mo.Option[" + str.Name + "]",
						},
						TypeDef{
							Name: "optic.ReturnMany",
						},
						TypeDef{
							Name: "optic.ReadWrite",
						},
						TypeDef{
							Name: "optic.BiDir",
						},
						TypeDef{
							Name: "optic.Pure",
						},
					},
				},
			},
			Body: []Statement{
				ReturnStmnt{
					Values: []Expression{
						CallExpr{
							Func: oPrefix + str.Name + "Of",
							Params: []Expression{
								CallExpr{
									Func: "optic.Some",
									TypeParams: []TypeExpression{
										TypeDef{
											Name: str.Name,
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

	oStruct.Methods = append(oStruct.Methods, MethodDef{
		Name: "Option",
		ReturnTypes: []TypeExpression{
			TypeDef{
				Name: "optic.Optic",
				TypeParams: []TypeExpression{
					TypeDef{
						Name: "I",
					},
					TypeDef{
						Name: "S",
					},
					TypeDef{
						Name: "T",
					},
					TypeDef{
						Name: "mo.Option[" + str.Name + "]",
					},
					TypeDef{
						Name: "mo.Option[" + str.Name + "]",
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
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{
					RetL(RwL(DirL(EErrL(ComposeLeft(
						"s",
						CallExpr{
							Func: "optic.PtrOption",
							TypeParams: []TypeExpression{
								TypeDef{
									Name: str.Name,
								},
							},
						},
					)))))},
			},
		},
	})

	fd.Structs = append(fd.Structs, oStruct)
}

func genSliceStruct(str *Struct, prefix string, opticType string) StructDef {

	return StructDef{
		Name: prefix + str.Name,
		TypeParams: []TypeExpression{
			TypeDef{
				Name: "I",
			},
			TypeDef{
				Name:       "S",
				Constraint: *sourceConstraint,
			},
			TypeDef{
				Name:       "T",
				Constraint: *sourceConstraint,
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
		Fields: []FieldDef{
			{
				Name: "",
				Type: TypeDef{
					Name: opticType,
					TypeParams: []TypeExpression{
						TypeDef{
							Name: "I",
						},
						TypeDef{
							Name: "S",
						},
						TypeDef{
							Name: "T",
						},
						TypeDef{
							Name: "optic.Collection[int," + str.Name + ", optic.Pure]",
						},
						TypeDef{
							Name: "optic.Collection[int," + str.Name + ", optic.Pure]",
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
			{
				Name: "o",
				Type: TypeDef{
					Name: opticType,
					TypeParams: []TypeExpression{
						TypeDef{
							Name: "I",
						},
						TypeDef{
							Name: "S",
						},
						TypeDef{
							Name: "T",
						},
						TypeDef{
							Name: "[]" + str.Name,
						},
						TypeDef{
							Name: "[]" + str.Name,
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
	}

}

func genMapStruct(str *Struct) StructDef {

	return StructDef{
		Name: "m" + str.Name,
		TypeParams: []TypeExpression{
			TypeDef{
				Name:       "I",
				Constraint: "comparable",
			},
			TypeDef{
				Name: "J",
			},
			TypeDef{
				Name:       "S",
				Constraint: *sourceConstraint,
			},
			TypeDef{
				Name:       "T",
				Constraint: *sourceConstraint,
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
		Fields: []FieldDef{
			{
				Name: "",
				Type: TypeDef{
					Name: "optic.Optic",
					TypeParams: []TypeExpression{
						TypeDef{
							Name: "J",
						},
						TypeDef{
							Name: "S",
						},
						TypeDef{
							Name: "T",
						},
						TypeDef{
							Name: "optic.Collection[I," + str.Name + ", optic.Pure]",
						},
						TypeDef{
							Name: "optic.Collection[I," + str.Name + ", optic.Pure]",
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
			{
				Name: "o",
				Type: TypeDef{
					Name: "optic.Optic",
					TypeParams: []TypeExpression{
						TypeDef{
							Name: "J",
						},
						TypeDef{
							Name: "S",
						},
						TypeDef{
							Name: "T",
						},
						TypeDef{
							Name: "map[I]" + str.Name,
						},
						TypeDef{
							Name: "map[I]" + str.Name,
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
	}

}

func genStructDefRetOne(str *Struct, prefix string, targetWrapper string, ptr bool) StructDef {

	targetName := str.Name
	if targetWrapper != "" {
		targetName = targetWrapper + "[" + targetName + "]"
	}
	if ptr {
		targetName = "*" + targetName
	}

	return StructDef{
		Name: prefix + str.Name,
		TypeParams: []TypeExpression{
			TypeDef{
				Name: "I",
			},
			TypeDef{
				Name:       "S",
				Constraint: *sourceConstraint,
			},
			TypeDef{
				Name:       "T",
				Constraint: *sourceConstraint,
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
		Fields: []FieldDef{
			{
				Name: "",
				Type: TypeDef{
					Name: "optic.Optic",
					TypeParams: []TypeExpression{
						TypeDef{
							Name: "I",
						},
						TypeDef{
							Name: "S",
						},
						TypeDef{
							Name: "T",
						},
						TypeDef{
							Name: targetName,
						},
						TypeDef{
							Name: targetName,
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
	}

}

func genMPrimitiveMethod(str *Struct, fld *Field, name string) MethodDef {

	lensCall := genFieldLens(str, fld)

	return MethodDef{
		Name: makePublic(fld.Name),
		ReturnTypes: []TypeExpression{
			TypeDef{
				Name: "optic.MakeLens" + name,
				TypeParams: []TypeExpression{
					TypeDef{
						Name: fld.MapKey,
					},
					TypeDef{
						Name: "S",
					},
					TypeDef{
						Name: "T",
					},
					TypeDef{
						Name: fld.TypeName,
					},
					TypeDef{
						Name: "RET",
					},
					TypeDef{
						Name: "RW",
					},
					TypeDef{
						Name: "optic.UniDir",
					},
					TypeDef{
						Name: "ERR",
					},
				},
			},
		},
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: "optic.NewMakeLens" + name,
						Params: []Expression{
							RetL(RwL(Ud(EErrL(Compose("s", lensCall))))),
						},
					},
				},
			},
		},
	}
}

func genOptMPrimitiveMethod(str *Struct, fld *Field, name string) MethodDef {

	lensCall := genPtrFieldLens(str, fld)

	return MethodDef{
		Name: makePublic(fld.Name),
		ReturnTypes: []TypeExpression{
			TypeDef{
				Name: "optic.MakeLens" + name,
				TypeParams: []TypeExpression{
					TypeDef{
						Name: fld.MapKey,
					},
					TypeDef{
						Name: "S",
					},
					TypeDef{
						Name: "T",
					},
					TypeDef{
						Name: fld.TypeName,
					},
					TypeDef{
						Name: "optic.ReturnMany",
					},
					TypeDef{
						Name: "RW",
					},
					TypeDef{
						Name: "optic.UniDir",
					},
					TypeDef{
						Name: "ERR",
					},
				},
			},
		},
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: "optic.NewMakeLens" + name,
						Params: []Expression{
							RetM(RwL(Ud(EErrL(Compose("s", lensCall))))),
						},
					},
				},
			},
		},
	}
}

func genSPrimitiveMethod(str *Struct, fld *Field, name string) MethodDef {

	lensCall := genFieldLens(str, fld)

	return MethodDef{
		Name: makePublic(fld.Name),
		ReturnTypes: []TypeExpression{
			TypeDef{
				Name: "optic.MakeLens" + name,
				TypeParams: []TypeExpression{
					TypeDef{
						Name: "I",
					},
					TypeDef{
						Name: "S",
					},
					TypeDef{
						Name: "T",
					},
					TypeDef{
						Name: fld.TypeName,
					},
					TypeDef{
						Name: "RET",
					},
					TypeDef{
						Name: "RW",
					},
					TypeDef{
						Name: "optic.UniDir",
					},
					TypeDef{
						Name: "ERR",
					},
				},
			},
		},
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{CallExpr{
					Func: "optic.NewMakeLens" + name,
					Params: []Expression{
						RetL(RwL(Ud(EErrL(ComposeLeft("s", lensCall))))),
					},
				},
				},
			},
		},
	}
}

func genOptSPrimitiveMethod(str *Struct, fld *Field, name string) MethodDef {

	lensCall := genPtrFieldLens(str, fld)

	return MethodDef{
		Name: makePublic(fld.Name),
		ReturnTypes: []TypeExpression{
			TypeDef{
				Name: "optic.MakeLens" + name,
				TypeParams: []TypeExpression{
					TypeDef{
						Name: "I",
					},
					TypeDef{
						Name: "S",
					},
					TypeDef{
						Name: "T",
					},
					TypeDef{
						Name: fld.TypeName,
					},
					TypeDef{
						Name: "optic.ReturnMany",
					},
					TypeDef{
						Name: "RW",
					},
					TypeDef{
						Name: "optic.UniDir",
					},
					TypeDef{
						Name: "ERR",
					},
				},
			},
		},
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: "optic.NewMakeLens" + name,
						Params: []Expression{
							RetM(RwL(Ud(EErrL(ComposeLeft("s", lensCall))))),
						},
					},
				},
			},
		},
	}
}

func genLensPrimitiveMethod(str *Struct, fld *Field, opticType string, opticAs1 string, dm *DataModel) MethodDef {

	lensCall := genFieldLens(str, fld)

	opticExpr1 := "s"
	if opticAs1 != "" {
		opticExpr1 = opticAs1 + "(s)"
	}

	returnType := TypeDef{
		Name: opticType,
		TypeParams: []TypeExpression{
			TypeDef{
				Name: "I",
			},
			TypeDef{
				Name: "S",
			},
			TypeDef{
				Name: "T",
			},
			TypeDef{
				Name: typeNameString(fld.TypeName, fld.TypePackage, fld.TypeParam),
			},
			TypeDef{
				Name: typeNameString(fld.TypeName, fld.TypePackage, fld.TypeParam),
			},
			TypeDef{
				Name: "RET",
			},
			TypeDef{
				Name: "RW",
			},
			TypeDef{
				Name: "optic.UniDir",
			},
			TypeDef{
				Name: "ERR",
			},
		},
	}

	opticExpr := RetL(RwL(Ud(EErrL(ComposeLeft(opticExpr1, lensCall)))))

	if isReal(fld.TypeName, fld.TypePackage, dm.TypeAliases) {

		opticExpr = CallExpr{
			Func: "optic.NewMakeLensRealOps",
			Params: []Expression{
				opticExpr,
			},
		}

		returnType = TypeDef{
			Name: "optic.MakeLensRealOps",
			TypeParams: []TypeExpression{
				TypeDef{
					Name: "I",
				},
				TypeDef{
					Name: "S",
				},
				TypeDef{
					Name: "T",
				},
				TypeDef{
					Name: typeNameString(fld.TypeName, fld.TypePackage, fld.TypeParam),
				},
				TypeDef{
					Name: "RET",
				},
				TypeDef{
					Name: "RW",
				},
				TypeDef{
					Name: "optic.UniDir",
				},
				TypeDef{
					Name: "ERR",
				},
			},
		}

	} else if isOrdered(fld.TypeName, fld.TypePackage, dm.TypeAliases) {

		opticExpr = CallExpr{
			Func: "optic.NewMakeLensOrdOps",
			Params: []Expression{
				opticExpr,
			},
		}

		returnType = TypeDef{
			Name: "optic.MakeLensOrdOps",
			TypeParams: []TypeExpression{
				TypeDef{
					Name: "I",
				},
				TypeDef{
					Name: "S",
				},
				TypeDef{
					Name: "T",
				},
				TypeDef{
					Name: typeNameString(fld.TypeName, fld.TypePackage, fld.TypeParam),
				},
				TypeDef{
					Name: "RET",
				},
				TypeDef{
					Name: "RW",
				},
				TypeDef{
					Name: "optic.UniDir",
				},
				TypeDef{
					Name: "ERR",
				},
			},
		}

	} else if isComparable(fld.TypeName, fld.TypePackage, dm) {
		opticExpr = CallExpr{
			Func: "optic.NewMakeLensCmpOps",
			Params: []Expression{
				opticExpr,
			},
		}

		returnType = TypeDef{
			Name: "optic.MakeLensCmpOps",
			TypeParams: []TypeExpression{
				TypeDef{
					Name: "I",
				},
				TypeDef{
					Name: "S",
				},
				TypeDef{
					Name: "T",
				},
				TypeDef{
					Name: typeNameString(fld.TypeName, fld.TypePackage, fld.TypeParam),
				},
				TypeDef{
					Name: "RET",
				},
				TypeDef{
					Name: "RW",
				},
				TypeDef{
					Name: "optic.UniDir",
				},
				TypeDef{
					Name: "ERR",
				},
			},
		}
	}

	return MethodDef{
		Name: makePublic(fld.Name),
		ReturnTypes: []TypeExpression{
			returnType,
		},
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{opticExpr},
			},
		},
	}
}

func genOptPrimitiveMethod(str *Struct, fld *Field, opticType string, opticAs1 string) MethodDef {

	lensCall := genPtrFieldLens(str, fld)

	opticExpr1 := "s"
	if opticAs1 != "" {
		opticExpr1 = opticAs1 + "(s)"
	}

	return MethodDef{
		Name: makePublic(fld.Name),
		ReturnTypes: []TypeExpression{
			TypeDef{
				Name: opticType,
				TypeParams: []TypeExpression{
					TypeDef{
						Name: "I",
					},
					TypeDef{
						Name: "S",
					},
					TypeDef{
						Name: "T",
					},
					TypeDef{
						Name: typeNameString(fld.TypeName, fld.TypePackage, fld.TypeParam),
					},
					TypeDef{
						Name: typeNameString(fld.TypeName, fld.TypePackage, fld.TypeParam),
					},
					TypeDef{
						Name: "optic.ReturnMany",
					},
					TypeDef{
						Name: "RW",
					},
					TypeDef{
						Name: "optic.UniDir",
					},
					TypeDef{
						Name: "ERR",
					},
				},
			},
		},
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{RetM(RwL(Ud(EErrL(ComposeLeft(opticExpr1, lensCall)))))},
			},
		},
	}
}

func genLensSliceMethod(str *Struct, fld *Field) MethodDef {

	lensCall := genFieldLens(str, fld)

	opticExpr1 := "s"

	return MethodDef{
		Name: makePublic(fld.Name),
		ReturnTypes: []TypeExpression{
			Star{
				Type: TypeDef{
					Name: "s" + fld.TypeName,
					TypeParams: []TypeExpression{
						TypeDef{
							Name: "I",
						},
						TypeDef{
							Name: "S",
						},
						TypeDef{
							Name: "T",
						},
						TypeDef{
							Name: "RET",
						},
						TypeDef{
							Name: "RW",
						},
						TypeDef{
							Name: "optic.UniDir",
						},
						TypeDef{
							Name: "ERR",
						},
					},
				},
			},
		},
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{
					genSRetMany(
						fld,
						RetL(RwL(Ud(EErrL(ComposeLeft(opticExpr1, lensCall))))),
						"RET",
					),
				},
			},
		},
	}
}

func genOptSliceMethod(str *Struct, fld *Field) MethodDef {

	lensCall := genPtrFieldLens(str, fld)

	opticExpr1 := "s"

	return MethodDef{
		Name: makePublic(fld.Name),
		ReturnTypes: []TypeExpression{
			Star{
				Type: TypeDef{
					Name: "s" + fld.TypeName,
					TypeParams: []TypeExpression{
						TypeDef{
							Name: "I",
						},
						TypeDef{
							Name: "S",
						},
						TypeDef{
							Name: "T",
						},
						TypeDef{
							Name: "optic.ReturnMany",
						},
						TypeDef{
							Name: "RW",
						},
						TypeDef{
							Name: "optic.UniDir",
						},
						TypeDef{
							Name: "ERR",
						},
					},
				},
			},
		},
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{
					genSRetMany(
						fld,
						RetM(RwL(Ud(EErrL(ComposeLeft(opticExpr1, lensCall))))),
						"optic.ReturnMany",
					),
				},
			},
		},
	}
}

func genLensMapMethod(str *Struct, fld *Field) MethodDef {

	lensCall := genFieldLens(str, fld)

	opticExpr1 := "s"

	return MethodDef{
		Name: makePublic(fld.Name),
		ReturnTypes: []TypeExpression{
			Star{
				Type: TypeDef{
					Name: "m" + fld.TypeName,
					TypeParams: []TypeExpression{
						TypeDef{
							Name: fld.MapKey,
						},
						TypeDef{
							Name: "I",
						},
						TypeDef{
							Name: "S",
						},
						TypeDef{
							Name: "T",
						},
						TypeDef{
							Name: "RET",
						},
						TypeDef{
							Name: "RW",
						},
						TypeDef{
							Name: "optic.UniDir",
						},
						TypeDef{
							Name: "ERR",
						},
					},
				},
			},
		},
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{
					genMRetMany(
						fld,
						RetL(RwL(Ud(EErrL(ComposeLeft(opticExpr1, lensCall))))),
					),
				},
			},
		},
	}
}

func genOptMapMethod(str *Struct, fld *Field) MethodDef {

	lensCall := genPtrFieldLens(str, fld)

	opticExpr1 := "s"

	return MethodDef{
		Name: makePublic(fld.Name),
		ReturnTypes: []TypeExpression{
			Star{
				Type: TypeDef{
					Name: "m" + fld.TypeName,
					TypeParams: []TypeExpression{
						TypeDef{
							Name: fld.MapKey,
						},
						TypeDef{
							Name: "I",
						},
						TypeDef{
							Name: "S",
						},
						TypeDef{
							Name: "T",
						},
						TypeDef{
							Name: "optic.ReturnMany",
						},
						TypeDef{
							Name: "RW",
						},
						TypeDef{
							Name: "optic.UniDir",
						},
						TypeDef{
							Name: "ERR",
						},
					},
				},
			},
		},
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{genMapRetManyExpr(fld.TypeName, "m", fld.MapKey, RetM(RwL(Ud(EErrL(ComposeLeft(opticExpr1, lensCall))))))},
			},
		},
	}
}

func genLensMethod(str *Struct, fld *Field, prefix string) MethodDef {

	lensCall := genFieldLens(str, fld)

	opticExpr1 := "s"

	return MethodDef{
		Name: makePublic(fld.Name),
		ReturnTypes: []TypeExpression{
			Star{
				Type: TypeDef{
					Name: prefix + fld.TypeName,
					TypeParams: []TypeExpression{
						TypeDef{
							Name: "I",
						},
						TypeDef{
							Name: "S",
						},
						TypeDef{
							Name: "T",
						},
						TypeDef{
							Name: "RET",
						},
						TypeDef{
							Name: "RW",
						},
						TypeDef{
							Name: "optic.UniDir",
						},
						TypeDef{
							Name: "ERR",
						},
					},
				},
			},
		},
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{genRetOneExpr(fld.TypeName, prefix, "I", "optic.UniDir", RetL(RwL(Ud(EErrL(ComposeLeft(opticExpr1, lensCall))))))},
			},
		},
	}
}

func genOptMethod(str *Struct, fld *Field, prefix string) MethodDef {

	lensCall := genPtrFieldLens(str, fld)

	opticExpr1 := "s"

	return MethodDef{
		Name: makePublic(fld.Name),
		ReturnTypes: []TypeExpression{
			Star{
				Type: TypeDef{
					Name: prefix + fld.TypeName,
					TypeParams: []TypeExpression{
						TypeDef{
							Name: "I",
						},
						TypeDef{
							Name: "S",
						},
						TypeDef{
							Name: "T",
						},
						TypeDef{
							Name: "optic.ReturnMany",
						},
						TypeDef{
							Name: "RW",
						},
						TypeDef{
							Name: "optic.UniDir",
						},
						TypeDef{
							Name: "ERR",
						},
					},
				},
			},
		},
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{genRetManyExpr(fld.TypeName, prefix, "I", RetM(RwL(Ud(EErrL(ComposeLeft(opticExpr1, lensCall))))))},
			},
		},
	}
}

func genTraverseSliceMethod(oPrefix string, str *Struct) MethodDef {

	prefix := "l"

	opticExpr1 := "s.o"
	opticExpr2 := "optic.TraverseSlice[" + str.Name + "]()"

	return MethodDef{
		Name: "Traverse",
		ReturnTypes: []TypeExpression{
			TypeDef{
				Name: "*" + prefix + str.Name,
				TypeParams: []TypeExpression{
					TypeDef{
						Name: "int",
					},
					TypeDef{
						Name: "S",
					},
					TypeDef{
						Name: "T",
					},
					TypeDef{
						Name: "optic.ReturnMany",
					},
					TypeDef{
						Name: "RW",
					},
					TypeDef{
						Name: "optic.UniDir",
					},
					TypeDef{
						Name: "ERR",
					},
				},
			},
		},
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: oPrefix + str.Name + "Of",
						Params: []Expression{
							RetM(RwL(Ud(EErrL(Compose(opticExpr1, opticExpr2))))),
						},
					},
				},
			},
		},
	}
}

func genTraverseMapMethod(oPrefix string, str *Struct) MethodDef {

	prefix := "l"

	opticExpr1 := "s.o"
	opticExpr2 := "optic.TraverseMap[I," + str.Name + "]()"

	return MethodDef{
		Name: "Traverse",
		ReturnTypes: []TypeExpression{
			TypeDef{
				Name: "*" + prefix + str.Name,
				TypeParams: []TypeExpression{
					TypeDef{
						Name: "I",
					},
					TypeDef{
						Name: "S",
					},
					TypeDef{
						Name: "T",
					},
					TypeDef{
						Name: "optic.ReturnMany",
					},
					TypeDef{
						Name: "RW",
					},
					TypeDef{
						Name: "optic.UniDir",
					},
					TypeDef{
						Name: "ERR",
					},
				},
			},
		},
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: oPrefix + str.Name + "Of",
						Params: []Expression{
							RetM(RwL(Ud(EErrL(Compose(opticExpr1, opticExpr2))))),
						},
					},
				},
			},
		},
	}
}

func genIndexSliceMethod(oPrefix string, str *Struct, prefix string) MethodDef {

	return MethodDef{
		Name: "Nth",
		ReturnTypes: []TypeExpression{
			TypeDef{
				Name: "*" + prefix + str.Name,
				TypeParams: []TypeExpression{
					TypeDef{
						Name: "int",
					},
					TypeDef{
						Name: "S",
					},
					TypeDef{
						Name: "T",
					},
					TypeDef{
						Name: "optic.ReturnMany",
					},
					TypeDef{
						Name: "RW",
					},
					TypeDef{
						Name: "optic.UniDir",
					},
					TypeDef{
						Name: "ERR",
					},
				},
			},
		},
		Params: []Param{
			{
				Name: "index",
				Type: TypeDef{
					Name: "int",
				},
			},
		},
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: oPrefix + str.Name + "Of",
						Params: []Expression{
							CallExpr{
								Func: "optic.Index",
								Params: []Expression{
									"s.Traverse()",
									"index",
								},
							},
						},
					},
				},
			},
		},
	}
}

func genIndexMapMethod(oPrefix string, str *Struct, prefix string) MethodDef {

	return MethodDef{
		Name: "Key",
		ReturnTypes: []TypeExpression{
			TypeDef{
				Name: "*" + prefix + str.Name,
				TypeParams: []TypeExpression{
					TypeDef{
						Name: "I",
					},
					TypeDef{
						Name: "S",
					},
					TypeDef{
						Name: "T",
					},
					TypeDef{
						Name: "optic.ReturnMany",
					},
					TypeDef{
						Name: "RW",
					},
					TypeDef{
						Name: "optic.UniDir",
					},
					TypeDef{
						Name: "ERR",
					},
				},
			},
		},
		Params: []Param{
			{
				Name: "index",
				Type: TypeDef{
					Name: "I",
				},
			},
		},
		Body: []Statement{
			ReturnStmnt{
				Values: []Expression{
					CallExpr{
						Func: oPrefix + str.Name + "Of",
						Params: []Expression{
							CallExpr{
								Func: "optic.Index",
								Params: []Expression{
									"s.Traverse()",
									"index",
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildLens(oPrefix string, fd *FileDef, str *Struct, oStruct *StructDef, dm *DataModel, combinators []Combinator) {
	lStruct := genStructDefRetOne(str, "l", "", false)

	for _, fld := range str.Fields {

		if isPrimitive(fld.TypeName, fld.TypePackage, dm) {
			if fld.Map {
				lStruct.Methods = append(lStruct.Methods, genMPrimitiveMethod(str, &fld, "Map"))
			} else if fld.Slice {
				lStruct.Methods = append(lStruct.Methods, genSPrimitiveMethod(str, &fld, "Slice"))
			} else {
				if fld.Pointer {
					lStruct.Methods = append(lStruct.Methods, genLensPrimitiveMethod(str, &fld, "optic.Optic", "", dm))
				} else {
					lStruct.Methods = append(lStruct.Methods, genLensPrimitiveMethod(str, &fld, "optic.Optic", "", dm))
				}
			}
		} else if fld.Slice {
			lStruct.Methods = append(lStruct.Methods, genLensSliceMethod(str, &fld))
		} else if fld.Map {
			lStruct.Methods = append(lStruct.Methods, genLensMapMethod(str, &fld))
		} else {
			if fld.Pointer {
				lStruct.Methods = append(lStruct.Methods, genLensMethod(str, &fld, "o"))
			} else {

				lStruct.Methods = append(lStruct.Methods, genLensMethod(str, &fld, "l"))
			}
		}
	}

	for _, c := range combinators {

		var params []Param
		args := []Expression{
			"s",
		}

		for _, p := range c.Parameters {
			params = append(params, Param{
				Name: p.Name,
				Type: TypeDef{Name: strings.ReplaceAll(p.Type, "{A}", str.Name)},
			})
			args = append(args, p.Name)
		}

		call := CallExpr{
			Func:   c.Package + "." + c.Name,
			Params: args,
		}

		slices.Reverse(c.Reconstrain)

		for _, re := range c.Reconstrain {
			call = CallExpr{
				Func:   re,
				Params: []Expression{call},
			}
		}

		call = CallExpr{
			Func:   "O" + str.Name + "Of",
			Params: []Expression{call},
		}

		m := MethodDef{
			Name:   c.Name,
			Params: params,
			ReturnTypes: []TypeExpression{
				TypeDef{
					Name: "*" + lStruct.Name + c.ReturnTypeParams,
				},
			},
			Body: []Statement{
				ReturnStmnt{
					Values: []Expression{
						call,
					},
				},
			},
		}

		lStruct.Methods = append(lStruct.Methods, m)
	}

	fd.Structs = append(fd.Structs, lStruct)

	fd.Funcs = append(fd.Funcs, genRetOneFunc(oPrefix, str, "l", ""))

	if *sourceConstraint == "" || *sourceConstraint == str.Name {
		oStruct.Methods = append(oStruct.Methods, genRetOneMethod(oPrefix, str, "l", "", ""))
	}

}

func structToTypeDef(s *Struct) TypeDef {

	var td TypeDef
	td.Name = s.Name

	if s.TypeParam != "" {
		td.TypeParams = []TypeExpression{
			TypeDef{
				Name: s.TypeParam,
			},
		}
	}

	return td
}
func fieldToTypeDefRaw(f *Field) TypeDef {

	var td TypeDef
	name := strings.Join(f.TypePackage, ".")
	if len(f.TypePackage) > 0 {
		name = name + "."
	}
	td.Name = name + f.TypeName

	if f.TypeParam != nil {

		var sb strings.Builder
		for i := len(f.TypeParam) - 1; i >= 0; i-- {
			sb.WriteString(f.TypeParam[i])
			if i != 0 {
				sb.WriteString(".")
			}
		}

		td.TypeParams = []TypeExpression{
			TypeDef{
				Name: sb.String(),
			},
		}
	}

	return td
}

func fieldToTypeDef(f *Field) TypeExpression {

	var expr TypeExpression
	expr = fieldToTypeDefRaw(f)

	if f.Pointer {
		expr = Star{
			Type: expr,
		}
	}

	if f.Slice {
		expr = SliceDef{
			Type: expr,
		}
	}

	if f.Map {
		expr = MapDef{
			Key: TypeDef{
				Name: f.MapKey,
			},
			Type: expr,
		}
	}

	return expr
}
