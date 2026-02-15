package codegen

import (
	"fmt"
)

type ComposeTree struct {
	ParamNum   int
	ParamName  string
	TypeParams []TypeExpression
	IsLeaf     bool
	Left       *ComposeTree
	Right      *ComposeTree
}

func BuildComposeTree(nodes []*ComposeTree) *ComposeTree {

	if len(nodes) == 1 {
		return nodes[0]
	}

	partitionAt := (len(nodes) + 1) / 2

	left := BuildComposeTree(nodes[:partitionAt])
	right := BuildComposeTree(nodes[partitionAt:])

	return &ComposeTree{
		IsLeaf: false,
		Left:   left,
		Right:  right,
	}
}

func BuildCompositionTree(prefix string, node *ComposeTree) TypeDef {
	if node.IsLeaf {
		return TypeDef{Name: fmt.Sprintf("%v%v", prefix, node.ParamNum)}
	} else {

		left := BuildCompositionTree(prefix, node.Left)
		right := BuildCompositionTree(prefix, node.Right)
		return TypeDef{
			Name: "CompositionTree",
			TypeParams: []TypeExpression{
				left,
				right,
			},
		}
	}
}
