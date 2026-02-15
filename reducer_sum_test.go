package optic_test

import (
	"fmt"

	. "github.com/spearson78/go-optic"
)

func ExampleReducer() {

	data := []int{1, 2, 3}

	var reducer Reduction[int, Pure] = Reducer[int](
		func() int {
			return 0
		},
		func(state, appendVal int) int {
			return state + appendVal
		},
		ReducerExprCustom("ExampleReducer"),
	)

	//Reduce applies the Reducer to the given optic and returns the aggregated value.
	var res int
	var ok bool
	res, ok = MustGetFirst(Reduce(TraverseSlice[int](), reducer), data)

	fmt.Println(res, ok)
	//Output:6 true
}
