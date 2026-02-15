package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	. "github.com/spearson78/go-optic/internal/codegen"
)

var typeParamNames = []string{
	"S",
	"T",
	"A",
	"B",
	"C",
	"D",
	"E",
	"F",
	"G",
	"H",
	"J",
	"K",
	"L",
	"M",
	"N",
	"O",
	"P",
	"Q",
	"R",
	"U",
	"V",
	"W",
	"X",
	"Y",
	"Z",
}

func buildReturnExpression(node *ComposeTree) Expression {
	if node.IsLeaf {
		return node.ParamName
	} else {

		left := buildReturnExpression(node.Left)
		right := buildReturnExpression(node.Right)
		return CallExpr{
			Func: "Compose",
			Params: []Expression{
				left,
				right,
			},
		}
	}
}

func main() {

	var genFileName = "../../../compose_generated.go"

	var fd FileDef
	fd.Package = "optic"

	for composeNum := 3; composeNum <= 10; composeNum++ {

		composeFnc := FuncDef{
			Docs: []string{
				" Compose" + strconv.Itoa(composeNum) + " returns an [Optic] composed of the " + strconv.Itoa(composeNum) + " input optics.",
				"",
				" Composition combines the optics such that the output of each optic is connected to the inputs of the next using the rightmost index.",
				"",
				" The composed optic is compatible with both view and modify actions.",
				"",
				" See:",
				"   - [Compose] for a version that takes 2 parameters.",
				"   - [Compose3] for a version that takes 3 parameters.",
				"   - [Compose4] for a version that takes 4 parameters.",
				"   - [Compose5] for a version that takes 5 parameters.",
				"   - [Compose6] for a version that takes 6 parameters.",
				"   - [Compose7] for a version that takes 7 parameters.",
				"   - [Compose8] for a version that takes 8 parameters.",
				"   - [Compose9] for a version that takes 9 parameters.",
				"   - [Compose10] for a version that takes 10 parameters.",
			},
			Name: fmt.Sprintf("Compose%v", composeNum),
		}

		var typeParams []TypeExpression

		for i := 1; i <= composeNum; i++ {
			typeParams = append(typeParams, TypeDef{Name: fmt.Sprintf("I%v", i)})
		}

		typeParamCount := 2 + (composeNum * 2)
		for i := 0; i < typeParamCount; i++ {
			typeParams = append(typeParams, TypeDef{Name: typeParamNames[i]})
		}

		for i := 1; i <= composeNum; i++ {
			typeParams = append(typeParams,
				TypeDef{Name: fmt.Sprintf("RET%v", i)},
				TypeDef{Name: fmt.Sprintf("RW%v", i)},
				TypeDef{Name: fmt.Sprintf("DIR%v", i)},
				TypeDef{Name: fmt.Sprintf("ERR%v", i)},
			)
		}

		composeFnc.TypeParams = typeParams

		var composeList []*ComposeTree

		for i := 1; i <= composeNum; i++ {

			name := fmt.Sprintf("o%v", i)
			composeTypeParams := []TypeExpression{
				TypeDef{Name: fmt.Sprintf("I%v", i)},   //I
				typeParams[composeNum+((i-1)*2)+0],     //S -> A
				typeParams[composeNum+((i-1)*2)+1],     //T -> B
				typeParams[composeNum+((i-1)*2)+2],     //A -> C
				typeParams[composeNum+((i-1)*2)+3],     //B -> D
				TypeDef{Name: fmt.Sprintf("RET%v", i)}, //RET
				TypeDef{Name: fmt.Sprintf("RW%v", i)},  //RW
				TypeDef{Name: fmt.Sprintf("DIR%v", i)}, //DIR
				TypeDef{Name: fmt.Sprintf("ERR%v", i)}, //ERR
			}

			composeFnc.Params = append(composeFnc.Params,
				Param{
					Name: name,
					Type: TypeDef{
						Name:       "Optic",
						TypeParams: composeTypeParams,
					},
				},
			)

			composeList = append(composeList, &ComposeTree{
				ParamNum:   i,
				ParamName:  fmt.Sprintf("o%v", i),
				TypeParams: composeTypeParams,
				IsLeaf:     true,
			})
		}

		returnTypeParams := []TypeExpression{
			TypeDef{Name: fmt.Sprintf("I%v", composeNum)}, //I
			TypeDef{Name: "S"},
			TypeDef{Name: "T"},
			typeParams[composeNum+(composeNum*2)], //Last 2 focus types
			typeParams[composeNum+(composeNum*2)+1],
		}

		composeTree := BuildComposeTree(composeList)

		ret := BuildCompositionTree("RET", composeTree)
		returnTypeParams = append(returnTypeParams, ret)

		rw := BuildCompositionTree("RW", composeTree)
		returnTypeParams = append(returnTypeParams, rw)

		dir := BuildCompositionTree("DIR", composeTree)
		returnTypeParams = append(returnTypeParams, dir)

		err := BuildCompositionTree("ERR", composeTree)
		returnTypeParams = append(returnTypeParams, err)

		composeFnc.ReturnTypes = []TypeExpression{
			TypeDef{
				Name:       "Optic",
				TypeParams: returnTypeParams,
			},
		}

		composeFnc.Body = append(composeFnc.Body, ReturnStmnt{
			Values: []Expression{buildReturnExpression(composeTree)},
		})

		fd.Funcs = append(fd.Funcs, composeFnc)
	}

	w, err := os.Create(genFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	WriteGoFile(w, &fd)

}
