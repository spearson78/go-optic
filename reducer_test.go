package optic_test

import (
	"fmt"

	"github.com/samber/lo"

	. "github.com/spearson78/go-optic"
)

func ExampleAsReducer() {

	data := []float64{2, 3, 4}

	//A naive implementation of a mean function.
	naiveMean := BinaryOp(
		func(a, b float64) float64 {
			return (a + b) / 2.0
		},
		"ExampleAsReducer",
	)

	meanReducer := AsReducer(0.0, naiveMean)

	result, ok := MustGetFirst(Reduce(TraverseSlice[float64](), meanReducer), data)

	fmt.Println(result, ok)
	//Output:3 true

}

func ExampleReducerT2() {

	data := []int{10, 20, 30}

	compositeReducer := ReducerT2(Sum[int](), Product[int]())

	//This optic takes slice of single ints and converts it to slice of lo.Tuple2[int,int] where A and B have duplicated values.
	optic := Compose(TraverseSlice[int](), DupT2[int]())

	//Reduce applies both Reducers at the same time and returns a lo.Tuple2 with the aggregated results.
	var res lo.Tuple2[int, int]
	var ok bool
	res, ok = MustGetFirst(Reduce(optic, compositeReducer), data)

	fmt.Println(res, ok)
	//Output: {60 6000} true
}

func ExampleReducerT3() {

	data := []int{10, 20, 30}

	compositeReducer := ReducerT3(Sum[int](), Product[int](), MaxReducer[int]())

	//This optic takes slice of single ints and converts it to slice of lo.Tuple2[int,int] where A and B have duplicated values.
	optic := Compose(TraverseSlice[int](), DupT3[int]())

	//Reduce applies all Reducers at the same time and returns a lo.Tuple3 with the aggregated results.
	var res lo.Tuple3[int, int, int]
	var ok bool
	res, ok = MustGetFirst(Reduce(optic, compositeReducer), data)

	fmt.Println(res, ok)
	//Output: {60 6000 30} true
}

func ExampleFirstReducer() {

	data := []int{1, 2, 3, 4, 5, 6, 4, 3, 2, 1, 4}

	res, err := Get(
		SliceOf(
			Grouped(
				SelfIndex(TraverseSlice[int](), EqT2[int]()),
				FirstReducer[int](),
			),
			len(data),
		),
		data,
	)

	fmt.Println(res, err)

	//Output:
	//[1 2 3 4 5 6] <nil>

}
