package expr

import (
	"fmt"
	"reflect"
)

type ReIndexedTree struct {
	OpticTypeExpr

	I     reflect.Type
	A     reflect.Type
	B     reflect.Type
	IxMap OpticExpression
}

func (e ReIndexedTree) Short() string {
	return fmt.Sprintf("ReIndexedTree(%v)", e.IxMap.Short())
}

func (e ReIndexedTree) String() string {
	return fmt.Sprintf("ReIndexedTree(%v)", e.IxMap.String())
}

type DiffTree struct {
	OpticTypeExpr

	I reflect.Type
	A reflect.Type

	Children  OpticExpression
	Threshold float64
	Distance  OpticExpression
	IxMatch   OpticExpression
}

func (e DiffTree) Short() string {
	return fmt.Sprintf("DiffTree(%v,%v)", e.Threshold, e.Distance.Short())
}

func (e DiffTree) String() string {
	return fmt.Sprintf("DiffTree(%v,%v)", e.Threshold, e.Distance.String())
}

type TreeOp struct {
	OpticTypeExpr
	Children OpticExpression
	Op       OpticExpression
}

func (e TreeOp) Short() string {
	return fmt.Sprintf("TreeOp(%v)", e.Op.Short())
}

func (e TreeOp) String() string {
	return fmt.Sprintf("TreeOp(%v)", e.Op.String())
}

type TreeMerge struct {
	OpticTypeExpr
	Children    OpticExpression
	ChildrenSeq OpticExpression
}

func (e TreeMerge) Short() string {
	return fmt.Sprintf("TreeMerge(%v)", e.Children.Short())
}

func (e TreeMerge) String() string {
	return fmt.Sprintf("TreeMerge(%v)", e.Children.String())
}

type ResolvePath struct {
	OpticTypeExpr
	Children OpticExpression
	Path     []any
}

func (e ResolvePath) Short() string {
	return fmt.Sprintf("ResolvePath(%v,%v)", e.Children.Short(), e.Path)
}

func (e ResolvePath) String() string {
	return fmt.Sprintf("ResolvePath(%v,%v)", e.Children.String(), e.Path)
}

type TopDownFiltered struct {
	OpticTypeExpr

	Children OpticExpression
	Pred     OpticExpression
}

func (e TopDownFiltered) Short() string {
	return fmt.Sprintf("TopDownFiltered(%v,%v)", e.Children.Short(), e.Pred.Short())
}

func (e TopDownFiltered) String() string {
	return fmt.Sprintf("TopDownFiltered(%v,%v)", e.Children.String(), e.Pred.String())
}

type BottomUpFiltered struct {
	OpticTypeExpr

	Children OpticExpression
	Pred     OpticExpression
}

func (e BottomUpFiltered) Short() string {
	return fmt.Sprintf("BottomUpFiltered(%v,%v)", e.Children.Short(), e.Pred.Short())
}

func (e BottomUpFiltered) String() string {
	return fmt.Sprintf("BottomUpFiltered(%v,%v)", e.Children.String(), e.Pred.String())
}

type BreadthFirstFiltered struct {
	OpticTypeExpr

	Children OpticExpression
	Pred     OpticExpression
}

func (e BreadthFirstFiltered) Short() string {
	return fmt.Sprintf("BreadthFirstFiltered(%v,%v)", e.Children.Short(), e.Pred.Short())
}

func (e BreadthFirstFiltered) String() string {
	return fmt.Sprintf("BreadthFirstFiltered(%v,%v)", e.Children.String(), e.Pred.String())
}

type TraverseTreeChildrenExpr struct {
	OpticTypeExpr

	IxMatch OpticExpression

	I reflect.Type
	A reflect.Type
}

func (e TraverseTreeChildrenExpr) Short() string {
	return "TraverseTreeChildren"
}

func (e TraverseTreeChildrenExpr) String() string {
	return "TraverseTreeChildren"
}

type WithChildPath struct {
	OpticTypeExpr

	Optic OpticExpression
	I     reflect.Type
	A     reflect.Type
}

func (e WithChildPath) Short() string {
	return "WithChildPath"
}

func (e WithChildPath) String() string {
	return "WithChildPath"
}
