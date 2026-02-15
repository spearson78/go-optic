package optic_test

import (
	"fmt"

	"github.com/samber/lo"

	. "github.com/spearson78/go-optic"
)

func ExampleMapOp() {

	data := lo.T2(
		1,
		map[int]string{
			1: "alpha",
			2: "beta",
			3: "gamma",
			4: "delta",
		},
	)

	optic := MapOp(
		FilteredCol[int](Ne("beta")),
	)

	var getRes map[int]string = MustGet(
		Compose(
			T2B[int, map[int]string](),
			optic,
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, map[int]string] = MustModify(
		T2B[int, map[int]string](),
		optic,
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//map[1:alpha 3:gamma 4:delta]
	//{1 map[1:alpha 3:gamma 4:delta]}
}

func ExampleFilteredMap() {

	data := lo.T2(
		1,
		map[int]string{
			1: "alpha",
			2: "beta",
			3: "gamma",
			4: "delta",
		},
	)

	var getRes map[int]string = MustGet(
		Compose(
			T2B[int, map[int]string](),
			FilteredMap[int](Ne("beta")),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, map[int]string] = MustModify(
		T2B[int, map[int]string](),
		FilteredMap[int](Ne("beta")), //FilteredMap can be used as an operation
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//map[1:alpha 3:gamma 4:delta]
	//{1 map[1:alpha 3:gamma 4:delta]}
}

func ExampleFilteredMapI() {

	data := lo.T2(
		1,
		map[int]string{
			1: "alpha",
			2: "beta",
			3: "gamma",
			4: "delta",
		},
	)

	var getRes map[int]string = MustGet(
		Compose(
			T2B[int, map[int]string](),
			FilteredMapI(NeI[string](2)),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, map[int]string] = MustModify(
		T2B[int, map[int]string](),
		FilteredMapI(NeI[string](2)), //FilteredMapI can be used as an operation
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//map[1:alpha 3:gamma 4:delta]
	//{1 map[1:alpha 3:gamma 4:delta]}
}

func ExampleAppendMap() {

	data := lo.T2(
		1,
		map[int]string{
			1: "alpha",
			2: "beta",
		},
	)

	var getRes map[int]string = MustGet(
		Compose(
			T2B[int, map[int]string](),
			AppendMap[int](MapCol(map[int]string{3: "gamma", 4: "delta"})),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, map[int]string] = MustModify(
		T2B[int, map[int]string](),
		AppendMap[int](MapCol(map[int]string{3: "gamma", 4: "delta"})),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//map[1:alpha 2:beta 3:gamma 4:delta]
	//{1 map[1:alpha 2:beta 3:gamma 4:delta]}
}

func ExampleMapCol() {

	data := map[string]int{
		"alpha": 1,
		"beta":  2,
	}

	var result map[string]int = MustModify(
		MapToCol[string, int](),
		AppendCol(
			MapCol(map[string]int{
				"gamma": 3,
				"delta": 4,
			}),
		),
		data,
	)

	fmt.Println(result)

	//Output:
	//map[alpha:1 beta:2 delta:4 gamma:3]
}
