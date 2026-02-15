package optic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/samber/lo"

	. "github.com/spearson78/go-optic"
)

func TestPredicateConsistency(t *testing.T) {

	ValidateOpticTest(t, True[int](), 1, true)
	ValidateOpticTest(t, False[int](), 1, false)

	ValidateOpticTest(t, Even[int](), 2, false)
	ValidateOpticTest(t, Gt[int](1), 2, false)

	ValidateOpticTestPred(t, Any(TraverseSlice[int](), Eq(2)), []int{1, 2, 3}, false, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, All(TraverseSlice[int](), Eq(2)), []int{1, 2, 3}, false, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, NotEmpty(TraverseSlice[int]()), []int{1, 2, 3}, false, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, Empty(TraverseSlice[int]()), []int{}, false, EqDeepT2[[]int]())
}

func ExamplePredGet() {

	res, err := PredGet(context.Background(), Gt(3), 7)
	fmt.Println(res, err)

	//PredGet converts an empty result to false
	res, err = PredGet(
		context.Background(),
		Compose(
			First(
				TraverseSlice[int](),
			),
			Gt(3),
		),
		[]int{},
	)
	fmt.Println(res, err)

	//Output:
	//true <nil>
	//false <nil>
}

func ExampleTrue() {

	data := []int{1, 2, 3, 4, 5}

	filtered := MustGet(
		SliceOf(
			Filtered(

				TraverseSlice[int](), True[int](),
			),
			len(data),
		),
		data,
	)
	fmt.Println(filtered)

	//Output:[1 2 3 4 5]
}

func ExampleFalse() {

	data := []int{1, 2, 3, 4, 5}

	filtered := MustGet(
		SliceOf(
			Filtered(

				TraverseSlice[int](), False[int](),
			),
			len(data),
		),
		data,
	)
	fmt.Println(filtered)

	//Output:[]
}

func ExampleEven() {

	data := []int{1, 2, 3, 4, 5}

	filtered := MustGet(SliceOf(Filtered(TraverseSlice[int](), Even[int]()), len(data)), data)
	fmt.Println(filtered)

	var overResult []bool = MustModify(TraverseSliceP[int, bool](), Even[int](), data)
	fmt.Println(overResult)

	//Output:[2 4]
	//[false true false true false]
}

func ExampleGt() {

	data := []int{1, 2, 3, 4, 5}

	var filtered []int = MustGet(SliceOf(Filtered(TraverseSlice[int](), Gt(3)), len(data)), data)
	fmt.Println(filtered)

	var overResult []bool = MustModify(TraverseSliceP[int, bool](), Gt(3), data)
	fmt.Println(overResult)

	//Output:[4 5]
	//[false false false true true]
}

func ExampleGtOp() {

	data := []lo.Tuple2[int, int]{
		lo.T2(1, 1),
		lo.T2(1, 2),
		lo.T2(2, 1),
		lo.T2(2, 2),
	}

	//returns a predicate that matches if A>B in the tuple
	//the type of this predicate is equivalent to func(lo.Tuple2[int,int])bool
	predicate := GtOp(T2A[int, int](), T2B[int, int]())

	//Returns an optic that filters a slice of int tuples using the predicate
	filterOptic := Filtered(TraverseSlice[lo.Tuple2[int, int]](), predicate)
	filtered := MustGet(SliceOf(filterOptic, len(data)), data)
	fmt.Println(filtered)

	var overResult []bool = Must(Modify(TraverseSliceP[lo.Tuple2[int, int], bool](), predicate, data))
	fmt.Println(overResult)

	//Output:[{2 1}]
	//[false false true false]
}

func ExampleGte() {

	data := []int{1, 2, 3, 4, 5}

	var filtered []int = MustGet(SliceOf(Filtered(TraverseSlice[int](), Gte(3)), len(data)), data)
	fmt.Println(filtered)

	var overResult []bool = MustModify(TraverseSliceP[int, bool](), Gte(3), data)
	fmt.Println(overResult)

	//Output:[3 4 5]
	//[false false true true true]
}

func ExampleGteOp() {

	data := []lo.Tuple2[int, int]{
		lo.T2(1, 1),
		lo.T2(1, 2),
		lo.T2(2, 1),
		lo.T2(2, 2),
	}

	//returns a predicate that matches if A>=B in the tuple
	//the type of this predicate is equivalent to func(lo.Tuple2[int,int])bool
	predicate := GteOp(T2A[int, int](), T2B[int, int]())

	//Returns an optic that filters a slice of int tuples using the predicate
	filterOptic := Filtered(TraverseSlice[lo.Tuple2[int, int]](), predicate)
	filtered := MustGet(SliceOf(filterOptic, len(data)), data)
	fmt.Println(filtered)

	var overResult []bool = Must(Modify(TraverseSliceP[lo.Tuple2[int, int], bool](), predicate, data))
	fmt.Println(overResult)

	//Output:
	//[{1 1} {2 1} {2 2}]
	//[true false true true]
}

func ExampleLt() {

	data := []int{1, 2, 3, 4, 5}

	filtered := MustGet(SliceOf(Filtered(TraverseSlice[int](), Lt(3)), len(data)), data)
	fmt.Println(filtered)

	var overResult []bool = MustModify(TraverseSliceP[int, bool](), Lt(3), data)
	fmt.Println(overResult)

	//Output:[1 2]
	//[true true false false false]
}

func ExampleLtOp() {

	data := []lo.Tuple2[int, int]{
		lo.T2(1, 1),
		lo.T2(1, 2),
		lo.T2(2, 1),
		lo.T2(2, 2),
	}

	//returns a predicate that matches if A<B in the tuple
	//the type of this predicate is equivalent to func(lo.Tuple2[int,int])bool
	predicate := LtOp(T2A[int, int](), T2B[int, int]())

	//Returns an optic that filters a slice of int tuples using the predicate
	filterOptic := Filtered(TraverseSlice[lo.Tuple2[int, int]](), predicate)
	filtered := MustGet(SliceOf(filterOptic, len(data)), data)
	fmt.Println(filtered)

	var overResult []bool = Must(Modify(TraverseSliceP[lo.Tuple2[int, int], bool](), predicate, data))
	fmt.Println(overResult)

	//Output:[{1 2}]
	//[false true false false]
}

func ExampleLte() {

	data := []int{1, 2, 3, 4, 5}

	filtered := MustGet(SliceOf(Filtered(TraverseSlice[int](), Lte(3)), len(data)), data)
	fmt.Println(filtered)

	var overResult []bool = MustModify(TraverseSliceP[int, bool](), Lte(3), data)
	fmt.Println(overResult)

	//Output:[1 2 3]
	//[true true true false false]
}

func ExampleLteOp() {

	data := []lo.Tuple2[int, int]{
		lo.T2(1, 1),
		lo.T2(1, 2),
		lo.T2(2, 1),
		lo.T2(2, 2),
	}

	//returns a predicate that matches if A<=B in the tuple
	//the type of this predicate is equivalent to func(lo.Tuple2[int,int])bool
	predicate := LteOp(T2A[int, int](), T2B[int, int]())

	//Returns an optic that filters a slice of int tuples using the predicate
	filterOptic := Filtered(TraverseSlice[lo.Tuple2[int, int]](), predicate)
	filtered := MustGet(SliceOf(filterOptic, len(data)), data)
	fmt.Println(filtered)

	var overResult []bool = Must(Modify(TraverseSliceP[lo.Tuple2[int, int], bool](), predicate, data))
	fmt.Println(overResult)

	//Output:
	//[{1 1} {1 2} {2 2}]
	//[true true false true]
}

func ExampleEq() {

	data := []int{1, 2, 3, 4, 5}

	filtered := MustGet(SliceOf(Filtered(TraverseSlice[int](), Eq(3)), len(data)), data)
	fmt.Println(filtered)

	var overResult []bool = MustModify(TraverseSliceP[int, bool](), Eq(3), data)
	fmt.Println(overResult)

	//Output:[3]
	//[false false true false false]
}

func ExampleEqOp() {

	data := []lo.Tuple2[int, int]{
		lo.T2(1, 1),
		lo.T2(1, 2),
		lo.T2(2, 1),
		lo.T2(2, 2),
	}

	//returns a predicate that matches if A==B in the tuple
	//the type of this predicate is equivalent to func(lo.Tuple2[int,int])bool
	predicate := EqOp(T2A[int, int](), T2B[int, int]())

	//Returns an optic that filters a slice of int tuples using the predicate
	filterOptic := Filtered(TraverseSlice[lo.Tuple2[int, int]](), predicate)
	filtered := MustGet(SliceOf(filterOptic, len(data)), data)
	fmt.Println(filtered)

	var overResult []bool = Must(Modify(TraverseSliceP[lo.Tuple2[int, int], bool](), predicate, data))
	fmt.Println(overResult)

	//Output:[{1 1} {2 2}]
	//[true false false true]
}

func ExampleIn() {

	vowels := In('a', 'e', 'i', 'o', 'u')

	index, vowel, found := MustGetFirstI(Filtered(TraverseString(), vowels), "hello world")
	fmt.Println(index, string(vowel), found)
	//Output:
	//1 e true
}

func ExampleNe() {

	data := []int{1, 2, 3, 4, 5}

	filtered := MustGet(SliceOf(Filtered(TraverseSlice[int](), Ne(3)), len(data)), data)
	fmt.Println(filtered)

	var overResult []bool = MustModify(TraverseSliceP[int, bool](), Ne(3), data)
	fmt.Println(overResult)

	//Output:[1 2 4 5]
	//[true true false true true]
}

func ExampleNeOp() {

	data := []lo.Tuple2[int, int]{
		lo.T2(1, 1),
		lo.T2(1, 2),
		lo.T2(2, 1),
		lo.T2(2, 2),
	}

	//returns a predicate that matches if A!=B in the tuple
	//the type of this predicate is equivalent to func(lo.Tuple2[int,int])bool
	predicate := NeOp(T2A[int, int](), T2B[int, int]())

	//Returns an optic that filters a slice of int tuples using the predicate
	filterOptic := Filtered(TraverseSlice[lo.Tuple2[int, int]](), predicate)
	filtered := MustGet(SliceOf(filterOptic, len(data)), data)
	fmt.Println(filtered)

	var overResult []bool = Must(Modify(TraverseSliceP[lo.Tuple2[int, int], bool](), predicate, data))
	fmt.Println(overResult)

	//Output:[{1 2} {2 1}]
	//[false true true false]
}

func ExampleAny() {

	data := []string{"1", "2", "3", "4"}

	//This optic focuses each string in the slice
	optic := TraverseSlice[string]()

	var result bool
	var err error
	result, err = Get(Any(optic, Eq("3")), data)
	fmt.Println(result, err)

	//Output:
	//true <nil>
}

func ExampleContains() {

	data := []string{"1", "2", "3", "4"}

	//This optic focuses each string in the slice
	optic := TraverseSlice[string]()

	var result bool
	var err error
	result, err = Get(Contains(optic, "3"), data)
	fmt.Println(result, err)

	//Output:
	//true <nil>
}

func ExampleAll() {

	data := []int{10, 20, 30, 40}

	//This optic focuses each int in the slice
	optic := TraverseSlice[int]()

	var result bool
	var err error
	result, err = Get(All(optic, Gt(5)), data)
	fmt.Println(result, err)

	//Output:
	//true <nil>
}

func ExampleNotEmpty() {

	var result bool = MustGet(
		NotEmpty(TraverseSlice[string]()),
		[]string{"alpha", "beta", "gamma", "delta"},
	)
	fmt.Println(result)

	result = MustGet(
		NotEmpty(TraverseSlice[string]()),
		[]string{},
	)
	fmt.Println(result)

	//Output:
	//true
	//false
}

func ExampleEmpty() {

	var result bool = MustGet(
		Empty(TraverseSlice[string]()),
		[]string{"alpha", "beta", "gamma", "delta"},
	)
	fmt.Println(result)

	result = MustGet(
		Empty(TraverseSlice[string]()),
		[]string{},
	)
	fmt.Println(result)

	//Output:
	//false
	//true
}
func ExamplePredT2ToOptic() {

	data := []lo.Tuple2[bool, bool]{
		lo.T2(false, false),
		lo.T2(false, true),
		lo.T2(true, false),
		lo.T2(true, true),
	}

	//returns a predicate that matches if A && B in the tuple
	predicate := PredT2ToOptic(T2A[bool, bool](), AndT2(), T2B[bool, bool]())

	filtered := MustGet(
		SliceOf(
			Filtered(
				TraverseSlice[lo.Tuple2[bool, bool]](), predicate,
			),
			len(data),
		),
		data,
	)
	fmt.Println(filtered)

	var overResult []bool = MustModify(
		TraverseSliceP[lo.Tuple2[bool, bool], bool](),
		predicate,
		data,
	)
	fmt.Println(overResult)

	//Output:
	//[{true true}]
	//[false false false true]
}
