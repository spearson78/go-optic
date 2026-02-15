package optic_test

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
)

func TestLogicConsistency(t *testing.T) {

	ValidateOpticTest(t, And(true), true, false)
	ValidateOpticTest(t, And(true), false, false)
	ValidateOpticTest(t, And(false), true, false)
	ValidateOpticTest(t, And(false), false, false)

	ValidateOpticTest(t, OpT2ToOptic(T2A[bool, bool](), AndT2(), T2A[bool, bool]()), lo.T2(false, false), false)
	ValidateOpticTest(t, PredT2ToOptic(T2A[bool, bool](), AndT2(), T2A[bool, bool]()), lo.T2(false, true), false)
	ValidateOpticTest(t, PredT2ToOptic(T2A[bool, bool](), AndT2(), T2A[bool, bool]()), lo.T2(true, false), false)
	ValidateOpticTest(t, PredT2ToOptic(T2A[bool, bool](), AndT2(), T2A[bool, bool]()), lo.T2(true, true), false)

	ValidateOpticTest(t, Or(true), true, false)
	ValidateOpticTest(t, Or(true), false, false)
	ValidateOpticTest(t, Or(false), true, false)
	ValidateOpticTest(t, Or(false), false, false)

	ValidateOpticTest(t, Not(), true, false)
	ValidateOpticTest(t, Not(), false, false)
	ValidateOpticTest(t, Not(), true, true)
	ValidateOpticTest(t, Not(), false, true)

	ValidateOpticTest(t, NotOp(Identity[bool]()), false, false)
	ValidateOpticTest(t, NotOp(Identity[bool]()), true, false)

	ValidateOpticTestPred(
		t,
		If(
			True[[]int](),
			Filtered(TraverseSlice[int](), Even[int]()),
			Filtered(TraverseSlice[int](), Odd[int]()),
		),
		[]int{1, 2, 3, 4},
		2,
		EqDeepT2[[]int](),
	)

	ValidateOpticTestPred(
		t,
		If(
			False[[]int](),
			Filtered(TraverseSlice[int](), Even[int]()),
			Filtered(TraverseSlice[int](), Odd[int]()),
		),
		[]int{1, 2, 3, 4},
		3,
		EqDeepT2[[]int](),
	)

}

func ExampleAndT2() {

	data := []bool{true, false, true, true, true}

	//The And function can be converted to a Reducer to detect if all focused values are true in a traversal where the empty traversal is considered to be all true values.
	result, ok := MustGetFirst(Reduce(TraverseSlice[bool](), AsReducer(true, AndT2())), data)
	fmt.Println(result, ok)

	//Output:false true
}

func ExampleAndOp() {

	data := []lo.Tuple2[int, int]{
		lo.T2(1, 1),
		lo.T2(1, 2),
		lo.T2(2, 1),
		lo.T2(2, 2),
		lo.T2(3, 1),
	}

	//returns a predicate that matches if A>B in the tuple
	predicateGt := GtOp(T2A[int, int](), T2B[int, int]())
	//returns a predicate that matches if A==2 in the tuple
	predicateEq2 := Compose(T2A[int, int](), Eq(2))

	//These 2 predicates are compatible as they have an int source type. They can be combined into an And predicate
	predicateAnd := AndOp(predicateGt, predicateEq2)

	//Returns an optic that filters a slice of int tuples using the predicate
	filterOptic := Filtered(TraverseSlice[lo.Tuple2[int, int]](), predicateAnd)
	filtered := MustGet(SliceOf(filterOptic, len(data)), data)
	fmt.Println(filtered)

	//We can focus on this composite filter predicate and modify the focused element by adding 10
	var overResult []lo.Tuple2[int, int] = MustModify(Compose(filterOptic, T2B[int, int]()), Add(10), data)
	fmt.Println(overResult)

	//Output:[{2 1}]
	//[{1 1} {1 2} {2 11} {2 2} {3 1}]
}

func ExampleOrT2() {

	data := []bool{true, false, true, true, true}

	//The Or function can be converted to a Reducer to detect if any focused values are true in a traversal where the empty traversal is considered to be all true values.
	result, ok := MustGetFirst(Reduce(TraverseSlice[bool](), AsReducer(true, OrT2())), data)
	fmt.Println(result, ok)

	//Output:true true
}

func ExampleOrOp() {

	data := []lo.Tuple2[int, int]{
		lo.T2(1, 1),
		lo.T2(1, 2),
		lo.T2(2, 1),
		lo.T2(2, 2),
		lo.T2(3, 1),
	}

	//returns a predicate that matches if A>B in the tuple
	predicateGt := GtOp(T2A[int, int](), T2B[int, int]())
	//returns a predicate that matches if A==2 in the tuple
	predicateEq2 := Compose(T2A[int, int](), Eq(2))

	//These 2 predicates are compatible as they have an int source type. They can be combined into an or predicate
	predicateAnd := OrOp(predicateGt, predicateEq2)

	//Returns an optic that filters a slice of int tuples using the predicate
	filterOptic := Filtered(TraverseSlice[lo.Tuple2[int, int]](), predicateAnd)
	filtered := MustGet(SliceOf(filterOptic, len(data)), data)
	fmt.Println(filtered)

	//We can focus on this composite filter predicate and modify the focused element by adding 10
	var overResult []lo.Tuple2[int, int] = MustModify(Compose(filterOptic, T2B[int, int]()), Add(10), data)
	fmt.Println(overResult)

	//Output:[{2 1} {2 2} {3 1}]
	//[{1 1} {1 2} {2 11} {2 12} {3 11}]
}

func ExampleNot() {

	data := []int{1, 2, 3, 4, 5}

	notPredicate := Compose(Gt(3), Not())

	var filtered []int = MustGet(SliceOf(Filtered(TraverseSlice[int](), notPredicate), len(data)), data)
	fmt.Println(filtered)

	var overResult []bool = MustModify(TraverseSliceP[int, bool](), notPredicate, data)
	fmt.Println(overResult)

	//Output:[1 2 3]
	//[true true true false false]
}

func ExampleNotOp() {

	data := []int{1, 2, 3, 4, 5}

	notPredicate := NotOp(Gt(3))

	var filtered []int = MustGet(SliceOf(Filtered(TraverseSlice[int](), notPredicate), len(data)), data)
	fmt.Println(filtered)

	var overResult []bool = MustModify(TraverseSliceP[int, bool](), notPredicate, data)
	fmt.Println(overResult)

	//Output:[1 2 3]
	//[true true true false false]
}

func ExampleIf() {

	data := []int{1, 2, 3, 0, 4, 5}

	ixOp := If(
		OpOnIx[int](Even[int]()), //Check for even index
		OpToOpI[int](Mul(10)),    //Multiple even indexes by 10
		OpToOpI[int](Mul(100)),   //Multiple odd indexes by 100
	)

	var result []int = MustModifyI(TraverseSlice[int](), ixOp, data)
	fmt.Println(result)

	//Output:
	//[10 200 30 0 40 500]
}
