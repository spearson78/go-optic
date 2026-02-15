package optic_test

import (
	"fmt"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
)

func ExampleEqDeepT2() {

	data := []lo.Tuple2[map[string]int, map[string]int]{
		lo.T2(
			map[string]int{
				"alpha": 1,
			},
			map[string]int{
				"alpha": 1,
			},
		),
		lo.T2(
			map[string]int{
				"alpha": 1,
			},
			map[string]int{
				"beta": 2,
			},
		),
	}

	optic := Compose(
		TraverseSlice[lo.Tuple2[map[string]int, map[string]int]](),
		EqDeepT2[map[string]int](),
	)

	var res []bool = MustGet(
		SliceOf(
			optic,
			2,
		),
		data,
	)

	fmt.Println(res)
	//Output:
	//[true false]

}

func ExampleEqDeep() {

	data := []map[string]int{

		map[string]int{
			"alpha": 1,
		},

		map[string]int{
			"beta": 2,
		},
	}

	optic := Compose(
		TraverseSlice[map[string]int](),
		EqDeep(map[string]int{
			"alpha": 1,
		}),
	)

	var res []bool = MustGet(
		SliceOf(
			optic,
			2,
		),
		data,
	)

	fmt.Println(res)
	//Output:
	//[true false]

}

func ExampleEqDeepOp() {

	data := [][]map[string]int{
		{
			map[string]int{
				"alpha": 1,
			},
			map[string]int{
				"alpha": 1,
			},
		},
		{
			map[string]int{
				"alpha": 1,
			},
			map[string]int{
				"beta": 2,
			},
		},
	}

	optic := Compose(
		TraverseSlice[[]map[string]int](),
		EqDeepOp(
			FirstOrDefault(Index(TraverseSlice[map[string]int](), 0), nil),
			FirstOrDefault(Index(TraverseSlice[map[string]int](), 1), nil),
		),
	)

	var res []bool = MustGet(
		SliceOf(
			optic,
			2,
		),
		data,
	)

	fmt.Println(res)
	//Output:
	//[true false]

}
