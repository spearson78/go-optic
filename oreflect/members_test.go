package oreflect_test

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/oreflect"
	"github.com/spearson78/go-optic/otree"
)

type Expr interface {
	expr()
}

type BinaryExpr struct {
	Left  Expr
	Op    string
	Right Expr
}

func (BinaryExpr) expr() {}
func (e BinaryExpr) String() string {
	return fmt.Sprintf("BinaryExpr(%v %v %v)}", e.Left, e.Op, e.Right)
}

type NumberLiteral float64

func (NumberLiteral) expr() {}
func (e NumberLiteral) String() string {
	return fmt.Sprintf("Num(%v)", float64(e))
}

func TestTraverseMembers(t *testing.T) {

	var expr Expr = BinaryExpr{
		Left:  NumberLiteral(1),
		Op:    "+",
		Right: NumberLiteral(2),
	}

	var setOp Expr = BinaryExpr{
		Left:  NumberLiteral(1),
		Op:    "-",
		Right: NumberLiteral(2),
	}

	var incNum Expr = BinaryExpr{
		Left:  NumberLiteral(2),
		Op:    "+",
		Right: NumberLiteral(3),
	}

	if res := MustGet(SliceOf(TraverseMembers[Expr, string](), 1), expr); !reflect.DeepEqual(res, []string{"+"}) {
		t.Fatal("BiPLate string", res)
	}

	if res := MustGet(SliceOf(TraverseMembers[Expr, string](), 1), expr); !reflect.DeepEqual(res, []string{"+"}) {
		t.Fatal("BiPLate string", res)
	}

	if res := MustGet(SliceOf(WithIndex(TraverseMembers[Expr, NumberLiteral]()), 2), expr); fmt.Sprint(res) != "[[.Left]:Num(1) [.Right]:Num(2)]" {
		t.Fatal("BiPLate NumberLiteral IxVals", res)
	}

	if res := MustSet(TraverseMembers[Expr, string](), "-", expr); !reflect.DeepEqual(res, setOp) {
		t.Fatal("BiPLate set  op", res)
	}

	if res := MustModify(TraverseMembers[Expr, NumberLiteral](), Op(func(i NumberLiteral) NumberLiteral {
		return i + 1
	}), expr); !reflect.DeepEqual(res, incNum) {
		t.Fatal("BiPLate inc num", res)
	}

}

type unexportedTest struct {
	a int
	b int
	c []int
	d map[string]int
	e any
}

func TestSetUnexportedField(t *testing.T) {

	//r := MustSet(TraverseMembers[*unexportedTest, int](), 1, &unexportedTest{})

	rs := reflect.ValueOf(unexportedTest{5, 5, []int{5, 5}, map[string]int{"a": 5, "b": 5}, 5})
	rs2 := reflect.New(rs.Type()).Elem()
	rs2.Set(rs)
	rf := rs2.Field(0)
	rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
	rf.SetInt(1)

	if !reflect.DeepEqual(rs2.Interface(), any(unexportedTest{1, 5, []int{5, 5}, map[string]int{"a": 5, "b": 5}, 5})) {
		t.Fatal(rs2)

	}

}

func TestSetUnexportedFieldTraverseMembers(t *testing.T) {

	r := MustSet(TraverseMembers[unexportedTest, int](), 1, unexportedTest{
		5,
		5,
		[]int{5, 5},
		map[string]int{"a": 5, "b": 5},
		unexportedTest{
			5,
			5,
			[]int{5, 5},
			map[string]int{"a": 5, "b": 5},
			5,
		},
	})

	if !reflect.DeepEqual(r, unexportedTest{
		1,
		1,
		[]int{1, 1},
		map[string]int{"a": 1, "b": 1},
		unexportedTest{
			1,
			1,
			[]int{1, 1},
			map[string]int{"a": 1, "b": 1},
			1,
		},
	}) {
		t.Fatal(r)

	}

}

type ptrTest struct {
	A *int
	B *int
}

func TestPtrTest(t *testing.T) {

	i := 1

	data := ptrTest{
		A: &i,
		B: &i,
	}

	addr1 := reflect.ValueOf(&data).Elem().Field(0).UnsafePointer()

	addr2 := reflect.ValueOf(&data).Elem().Field(1).UnsafePointer()

	if addr1 != addr2 {
		t.Fatal(addr1, addr2)
	}
}

type anyPtrTest struct {
	A any
	B any
}

func TestAnyPtrTest(t *testing.T) {

	i := 1

	data := anyPtrTest{
		A: &i,
		B: &i,
	}

	addr1 := reflect.ValueOf(&data).Elem().Field(0).Elem().UnsafePointer()

	addr2 := reflect.ValueOf(&data).Elem().Field(1).Elem().UnsafePointer()

	if addr1 != addr2 {
		t.Fatal(addr1, addr2)
	}
}

func TestTraverseMembersNil(t *testing.T) {

	i := 1

	data := ptrTest{
		A: &i,
		B: nil,
	}

	optic := TraverseMembers[ptrTest, *int]()

	result := MustGet(SliceOf(optic, 2), data)

	if !reflect.DeepEqual(result, []*int{&i, nil}) {
		t.Fatal(result)
	}

	modifyResult := MustModify(
		Compose(
			optic,
			TraversePtrP[int, int](),
		),
		Add(1),
		data,
	)

	v := 2
	if !reflect.DeepEqual(modifyResult, ptrTest{
		A: &v,
		B: nil,
	}) {
		t.Fatal(modifyResult)
	}

}

func TestMemberTree(t *testing.T) {

	var expr Expr = BinaryExpr{
		Left:  NumberLiteral(1),
		Op:    "+",
		Right: NumberLiteral(2),
	}

	if res := MustGet(SliceOf(otree.TopDown(TraverseMembers[Expr, Expr]()), 10), expr); !reflect.DeepEqual(res, []Expr{expr, NumberLiteral(1), NumberLiteral(2)}) {
		t.Fatal(res)
	}

	if res := MustModify(Compose(otree.TopDown(TraverseMembers[Expr, Expr]()), DownCast[Expr, NumberLiteral]()), Add(NumberLiteral(10)), expr); !reflect.DeepEqual(res, BinaryExpr{
		Left:  NumberLiteral(11),
		Op:    "+",
		Right: NumberLiteral(12),
	}) {
		t.Fatal(res)
	}
}

func TestRewriteMemberTree(t *testing.T) {

	var expr Expr = BinaryExpr{
		Left: BinaryExpr{
			Left:  NumberLiteral(1),
			Op:    "+",
			Right: NumberLiteral(2),
		},
		Op:    "+",
		Right: NumberLiteral(3),
	}

	exprOp := FieldLens(func(source *BinaryExpr) *string { return &source.Op })

	opName := Compose(
		DownCast[Expr, BinaryExpr](),
		exprOp,
	)

	value := Compose(
		DownCast[Expr, NumberLiteral](),
		IsoCast[NumberLiteral, int](),
	)

	//Rewrite handler that evaluates the "+" function
	rewriteHandler := otree.RewriteOp(func(tree Expr) (Expr, bool) {
		//If all the children are literal values we may be able to evaluate the expression
		childrenAllValues, _ := MustGetFirst(All(TraverseMembers[Expr, Expr](), NotEmpty(value)), tree)
		if childrenAllValues {

			funcName, _ := MustGetFirst(opName, tree)

			switch funcName {
			case "+":
				sumParams, _ := MustGetFirst(Reduce(Compose(TraverseMembers[Expr, Expr](), value), Sum[int]()), tree)
				return NumberLiteral(sumParams), true
			}
		}

		return tree, false
	})

	newTree := MustGet(otree.Rewrite(TraverseMembers[Expr, Expr](), rewriteHandler), expr)
	if !reflect.DeepEqual(newTree, NumberLiteral(6)) {
		t.Fatal(newTree)
	}

	data := lo.T2("alpha", expr)

	newTuple := MustModify(T2B[string, Expr](), otree.Rewrite(TraverseMembers[Expr, Expr](), rewriteHandler), data)

	if !reflect.DeepEqual(newTuple, lo.T2("alpha", Expr(NumberLiteral(6)))) {
		t.Fatal(newTuple)
	}
}
