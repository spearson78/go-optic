package optic_test

import (
	"fmt"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
)

func ExampleWithVar() {

	data := lo.T2(10, []int{1, 2, 3})

	result, err := Get(
		SliceOf(
			WithVar(
				Compose3(
					T2B[int, []int](),
					TraverseSlice[int](),
					AddOp(
						Var[int, int]("a"),
						Identity[int](),
					),
				),
				"a",
				T2A[int, []int](),
			),
			3,
		),
		data,
	)
	fmt.Println(result, err)

	//Output:
	//[11 12 13] <nil>
}

func ExampleVar() {

	data := lo.T2(10, []int{1, 2, 3})

	result, err := Get(
		SliceOf(
			WithVar(
				Compose3(
					T2B[int, []int](),
					TraverseSlice[int](),
					AddOp(
						Var[int, int]("a"),
						Identity[int](),
					),
				),
				"a",
				T2A[int, []int](),
			),
			3,
		),
		data,
	)
	fmt.Println(result, err)

	//Output:
	//[11 12 13] <nil>
}
