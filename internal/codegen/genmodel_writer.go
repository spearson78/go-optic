package codegen

import (
	"fmt"
	"io"
	"log"
	"strings"
)

func WriteGoFile(w io.Writer, fd *FileDef) {

	fmt.Fprintf(w, "package %v\n\n", fd.Package)

	if len(fd.Imports) > 0 || len(fd.DotImports) > 0 {
		fmt.Fprint(w, "import (\n")

		for _, imp := range fd.DotImports {
			fmt.Fprintf(w, "\t. \"%v\"\n", imp)
		}

		for _, imp := range fd.Imports {
			fmt.Fprintf(w, "\t\"%v\"\n", imp)
		}
		fmt.Fprint(w, ")\n")
	}

	for _, str := range fd.Structs {
		fmt.Fprintf(w, "type %v", str.Name)

		if len(str.TypeParams) > 0 {
			fmt.Fprint(w, "[")
			for i, tp := range str.TypeParams {
				if i != 0 {
					fmt.Fprint(w, ",")
				}
				WriteTypeDefSignature(w, tp)
				fmt.Fprint(w, " ")
				if td, ok := tp.(TypeDef); ok && td.Constraint != "" {
					fmt.Fprint(w, td.Constraint)
				} else {
					fmt.Fprint(w, "any")
				}

			}

			fmt.Fprint(w, "]")
		}

		fmt.Fprintf(w, " struct {\n")

		for _, fld := range str.Fields {
			fmt.Fprint(w, "\t")
			if fld.Name != "" {
				fmt.Fprintf(w, "%v ", fld.Name)
			}

			WriteTypeDef(w, fld.Type)
			fmt.Fprint(w, "\n")
		}

		fmt.Fprint(w, "}\n")

		for _, mthd := range str.Methods {
			for _, doc := range mthd.Docs {
				fmt.Fprint(w, "//")
				fmt.Fprint(w, doc)
				fmt.Fprint(w, "\n")
			}
			fmt.Fprint(w, "func (s *")
			fmt.Fprint(w, str.Name)

			if len(str.TypeParams) > 0 {
				fmt.Fprint(w, "[")
				for i, tp := range str.TypeParams {
					if i != 0 {
						fmt.Fprint(w, ",")
					}
					WriteTypeDefSignature(w, tp)
				}
				fmt.Fprint(w, "]")
			}

			fmt.Fprint(w, ")")
			fmt.Fprint(w, mthd.Name)
			fmt.Fprint(w, "(")
			for i, p := range mthd.Params {
				if i != 0 {
					fmt.Fprint(w, ",")
				}
				fmt.Fprintf(w, "%v ", p.Name)
				WriteTypeDef(w, p.Type)
			}
			fmt.Fprint(w, ")")

			if len(mthd.ReturnTypes) > 1 {
				fmt.Fprintf(w, "(")
				for i, r := range mthd.ReturnTypes {
					if i != 0 {
						fmt.Fprint(w, ",")
					}
					WriteTypeDef(w, r)
				}

				fmt.Fprintf(w, ")")
			} else if len(mthd.ReturnTypes) == 1 {
				WriteTypeDef(w, mthd.ReturnTypes[0])
			}

			fmt.Fprintf(w, "{\n")

			for _, stmnt := range mthd.Body {
				WriteStatement(w, stmnt)
				fmt.Fprint(w, "\n")
			}

			fmt.Fprintf(w, "}\n")
		}
	}

	for _, fnc := range fd.Funcs {
		for _, doc := range fnc.Docs {
			fmt.Fprint(w, "//")
			fmt.Fprint(w, doc)
			fmt.Fprint(w, "\n")
		}
		fmt.Fprintf(w, "func %v", fnc.Name)
		if len(fnc.TypeParams) > 0 {
			fmt.Fprint(w, "[")
			for i, tp := range fnc.TypeParams {
				if i != 0 {
					fmt.Fprint(w, ",")
				}
				WriteTypeDefSignature(w, tp)
				fmt.Fprint(w, " ")
				if td, ok := tp.(TypeDef); ok && td.Constraint != "" {
					fmt.Fprint(w, td.Constraint)
				} else {
					fmt.Fprint(w, "any")
				}
			}

			fmt.Fprint(w, "]")
		}

		fmt.Fprintf(w, "(")
		for i, p := range fnc.Params {
			if i != 0 {
				fmt.Fprint(w, ",")
			}
			fmt.Fprintf(w, "%v ", p.Name)
			WriteTypeDef(w, p.Type)
		}

		fmt.Fprintf(w, ")")

		if len(fnc.ReturnTypes) > 1 {
			fmt.Fprintf(w, "(")
			for i, r := range fnc.ReturnTypes {
				if i != 0 {
					fmt.Fprint(w, ",")
				}
				WriteTypeDef(w, r)
			}

			fmt.Fprintf(w, ")")
		} else if len(fnc.ReturnTypes) == 1 {
			WriteTypeDef(w, fnc.ReturnTypes[0])
		}

		fmt.Fprintf(w, "{\n")

		for _, stmnt := range fnc.Body {
			WriteStatement(w, stmnt)
			fmt.Fprint(w, "\n")
		}

		fmt.Fprintf(w, "}\n")
	}

	for _, vr := range fd.Vars {

		fmt.Fprintf(w, "var %v ", vr.Name)
		if vr.Type != nil {
			WriteTypeDefSignature(w, vr.Type)
		}
		fmt.Fprint(w, " = ")

		WriteExpression(w, vr.Value)
		fmt.Fprintf(w, "\n")
	}

}

func WriteTypeDefSignature(w io.Writer, e TypeExpression) {

	switch t := e.(type) {
	case TypeDef:
		fmt.Fprint(w, t.Name)
		if len(t.TypeParams) > 0 {
			fmt.Fprint(w, "[")
			for i := range t.TypeParams {
				tp := t.TypeParams[i]
				if i != 0 {
					fmt.Fprint(w, ",")
				}
				WriteTypeDefSignature(w, tp)
			}
			fmt.Fprint(w, "]")
		}
	default:
		panic("WriteTypeDefSignature unknown type expression")
	}
}

func WriteTypeDef(w io.Writer, e TypeExpression) {

	switch t := e.(type) {
	case SliceDef:
		fmt.Fprint(w, "[]")
		WriteTypeDef(w, t.Type)
	case MapDef:
		fmt.Fprint(w, "map[")
		WriteTypeDef(w, t.Key)
		fmt.Fprint(w, "]")
		WriteTypeDef(w, t.Type)
	case Star:
		fmt.Fprint(w, "*")
		WriteTypeDef(w, t.Type)
	case TypeDef:
		fmt.Fprint(w, t.Name)
		if len(t.TypeParams) > 0 {
			fmt.Fprint(w, "[")
			for i, tp := range t.TypeParams {
				if i != 0 {
					fmt.Fprint(w, ",")
				}
				WriteTypeDef(w, tp)
			}
			fmt.Fprint(w, "]")
		}
	case FuncType:
		fmt.Fprint(w, "func(")
		for i, p := range t.Params {
			if i != 0 {
				fmt.Fprint(w, ",")
			}
			fmt.Fprint(w, p.Name)
			fmt.Fprint(w, " ")
			WriteTypeDef(w, p.Type)
		}
		fmt.Fprint(w, ") ")

		if len(t.ReturnTypes) > 1 {
			fmt.Fprintf(w, "(")
			for i, r := range t.ReturnTypes {
				if i != 0 {
					fmt.Fprint(w, ",")
				}
				WriteTypeDef(w, r)
			}

			fmt.Fprintf(w, ")")
		} else if len(t.ReturnTypes) == 1 {
			WriteTypeDef(w, t.ReturnTypes[0])
		}
	case FuncDef:
		fnc := t
		fmt.Fprintf(w, "func")
		fmt.Fprintf(w, "(")
		for i, p := range fnc.Params {
			if i != 0 {
				fmt.Fprint(w, ",")
			}
			if p.Name != "" {
				fmt.Fprintf(w, "%v ", p.Name)
			}
			WriteTypeDef(w, p.Type)
		}

		fmt.Fprintf(w, ")")

		if len(fnc.ReturnTypes) > 1 {
			fmt.Fprintf(w, "(")
			for i, r := range fnc.ReturnTypes {
				if i != 0 {
					fmt.Fprint(w, ",")
				}
				WriteTypeDef(w, r)
			}

			fmt.Fprintf(w, ")")
		} else if len(fnc.ReturnTypes) == 1 {
			WriteTypeDef(w, fnc.ReturnTypes[0])
		}
	default:
		log.Fatalf("WriteTypeDef unknown typeexpression type %T", t)
	}
}

func WriteStatement(w io.Writer, stmnt Statement) {
	switch t := stmnt.(type) {
	case ReturnStmnt:
		fmt.Fprint(w, "return ")
		for i, ret := range t.Values {
			if i != 0 {
				fmt.Fprint(w, ", ")
			}
			WriteExpression(w, ret)
		}
	case AssignVar:
		for i, v := range t.Vars {
			if i != 0 {
				fmt.Fprint(w, ", ")

			}
			fmt.Fprint(w, v)
		}

		if t.Declare {
			fmt.Fprint(w, " := ")
		} else {
			fmt.Fprint(w, " = ")
		}

		WriteExpression(w, t.Value)
	case IfStmnt:
		fmt.Fprint(w, "if ")
		WriteExpression(w, t.Condition)
		fmt.Fprint(w, "{\n")
		for _, s := range t.OnTrue {
			WriteStatement(w, s)
			fmt.Fprint(w, "\n")
		}
		if t.OnFalse != nil {
			fmt.Fprint(w, "} else {\n")
			for _, s := range t.OnFalse {
				WriteStatement(w, s)
				fmt.Fprint(w, "\n")
			}
		}

		fmt.Fprint(w, "}\n")
	case VarDef:
		fmt.Fprintf(w, "var %v", t.Name)
		if t.Type != nil {
			fmt.Fprint(w, " ")
			WriteTypeDefSignature(w, t.Type)
		}
		if t.Value != nil {
			fmt.Fprint(w, " = ")
			WriteExpression(w, t.Value)
		}
		fmt.Fprint(w, "\n")
	case CallExpr:
		writeTypeSig(w, t.Func, t.TypeParams)
		fmt.Fprint(w, "(")
		for i, p := range t.Params {
			if i != 0 {
				fmt.Fprint(w, ",")
			}
			WriteExpression(w, p)
		}
		fmt.Fprint(w, ")")
	case MethodCallExpr:
		WriteExpression(w, t.Receiver)
		fmt.Fprintf(w, ".%v(", t.Name)
		for i, p := range t.Params {
			if i != 0 {
				fmt.Fprint(w, ",")
			}
			WriteExpression(w, p)
		}
		fmt.Fprint(w, ")")
	default:
		log.Fatalf("unknown statement type %T", t)
	}
}

func writeTypeSig(w io.Writer, name Expression, typeParams []TypeExpression) {
	WriteExpression(w, name)
	if len(typeParams) > 0 {
		fmt.Fprint(w, "[")
		for i, tp := range typeParams {
			if i != 0 {
				fmt.Fprint(w, ",")
			}
			WriteTypeDefSignature(w, tp)
		}
		fmt.Fprint(w, "]")
	}
}

func WriteExpression(w io.Writer, expr Statement) {
	switch t := expr.(type) {
	case CallExpr:
		writeTypeSig(w, t.Func, t.TypeParams)
		fmt.Fprint(w, "(")
		for i, p := range t.Params {
			if i != 0 {
				fmt.Fprint(w, ",")
			}
			WriteExpression(w, p)
		}
		fmt.Fprint(w, ")")
	case MethodCallExpr:
		WriteExpression(w, t.Receiver)
		fmt.Fprintf(w, ".%v(", t.Name)
		for i, p := range t.Params {
			if i != 0 {
				fmt.Fprint(w, ",")
			}
			WriteExpression(w, p)
		}
		fmt.Fprint(w, ")")
	case GenFuncExpr:
		fmt.Fprint(w, "func(")
		for i, p := range t.Params {
			if i != 0 {
				fmt.Fprint(w, ",")
			}
			fmt.Fprint(w, p.Name)
			fmt.Fprint(w, " ")
			WriteTypeDef(w, p.Type)
		}
		fmt.Fprint(w, ")")

		if len(t.ReturnTypes) > 1 {
			fmt.Fprintf(w, "(")
			for i, r := range t.ReturnTypes {
				if i != 0 {
					fmt.Fprint(w, ",")
				}
				WriteTypeDef(w, r)
			}

			fmt.Fprintf(w, ")")
		} else if len(t.ReturnTypes) == 1 {
			WriteTypeDef(w, t.ReturnTypes[0])
		}

		fmt.Fprintf(w, "{\n")

		for _, stmnt := range t.Body {
			WriteStatement(w, stmnt)
			fmt.Fprint(w, "\n")
		}

		fmt.Fprintf(w, "}")
	case AddressExpr:
		fmt.Fprint(w, "&")
		WriteExpression(w, t.Target)
	case BinaryExpr:
		WriteExpression(w, t.Left)
		fmt.Fprint(w, t.Op)
		WriteExpression(w, t.Right)
	case string:
		fmt.Fprint(w, t)
	case int:
		fmt.Fprint(w, t)
	case StructExpr:
		writeTypeSig(w, t.Type.Name, t.Type.TypeParams)
		fmt.Fprint(w, "{\n")
		for _, fld := range t.Fields {
			fmt.Fprint(w, "\t")
			fmt.Fprint(w, fld.Name)
			fmt.Fprint(w, ":\t")
			WriteExpression(w, fld.Value)
			fmt.Fprint(w, ",\n")
		}
		fmt.Fprint(w, "}")
	case StringLiteral:
		fmt.Fprint(w, "\"")
		fmt.Fprint(w, strings.ReplaceAll(strings.ReplaceAll(string(t), `\`, `\\`), `"`, `\"`))
		fmt.Fprint(w, "\"")
	case DeRef:
		fmt.Fprint(w, "*(")
		WriteExpression(w, t.Value)
		fmt.Fprint(w, ")")
	case FuncDef:
		fnc := t
		fmt.Fprintf(w, "func")
		fmt.Fprintf(w, "(")
		for i, p := range fnc.Params {
			if i != 0 {
				fmt.Fprint(w, ",")
			}
			fmt.Fprintf(w, "%v ", p.Name)
			WriteTypeDef(w, p.Type)
		}

		fmt.Fprintf(w, ")")

		if len(fnc.ReturnTypes) > 1 {
			fmt.Fprintf(w, "(")
			for i, r := range fnc.ReturnTypes {
				if i != 0 {
					fmt.Fprint(w, ",")
				}
				WriteTypeDef(w, r)
			}

			fmt.Fprintf(w, ")")
		} else if len(fnc.ReturnTypes) == 1 {
			WriteTypeDef(w, fnc.ReturnTypes[0])
		}

		fmt.Fprintf(w, "{\n")

		for _, stmnt := range fnc.Body {
			WriteStatement(w, stmnt)
			fmt.Fprint(w, "\n")
		}

		fmt.Fprintf(w, "}")
	case SliceExpr:
		fmt.Fprint(w, "[]")
		writeTypeSig(w, t.Type.Name, t.Type.TypeParams)
		fmt.Fprint(w, "{")
		for i, fld := range t.Values {
			if i != 0 {
				fmt.Fprint(w, " , ")
			}
			WriteExpression(w, fld)

		}
		fmt.Fprint(w, "}")
	default:
		log.Fatalf("unknown expression type %T", t)
	}
}
