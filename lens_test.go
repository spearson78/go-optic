package optic_test

import (
	"context"
	"fmt"
	"strconv"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
)

func ExampleLens() {
	//This lens focuses on the B element of a tuple
	tupleB := Lens[lo.Tuple2[int, int], int](
		func(source lo.Tuple2[int, int]) int {
			return source.B
		},
		func(focus int, source lo.Tuple2[int, int]) lo.Tuple2[int, int] {
			source.B = focus
			return source
		},
		ExprCustom("ExampleLens"),
	)

	var viewResult int = MustGet(tupleB, lo.T2(10, 20))
	fmt.Println(viewResult)

	var overResult lo.Tuple2[int, int] = MustModify(tupleB, Mul(2), lo.T2(10, 20))
	fmt.Println(overResult)

	//Output:20
	//{10 40}
}

func ExampleLensI() {
	//This lens focuses on the B element of a tuple with an index
	tupleB := LensI[int, lo.Tuple2[int, int], int](
		func(source lo.Tuple2[int, int]) (int, int) {
			return 1, source.B
		},
		func(focus int, source lo.Tuple2[int, int]) lo.Tuple2[int, int] {
			source.B = focus
			return source
		},
		func(indexA, indexB int) bool {
			return indexA == indexB
		},
		ExprCustom("ExampleLensI"),
	)

	var viewIndex, viewResult int = MustGetI(tupleB, lo.T2(10, 20))
	fmt.Println(viewIndex, viewResult)

	var overResult lo.Tuple2[int, int] = MustModifyI(tupleB, OpI(func(index int, focus int) int {
		return index + focus
	}), lo.T2(10, 20))
	fmt.Println(overResult)

	//Output:1 20
	//{10 21}
}

func ExampleLensE() {
	//This lens focuses on the B element of a string tuple and converts it to an int
	//reporting the error in case of conversion failure
	tupleB := LensE[lo.Tuple2[string, string], int](
		func(ctx context.Context, source lo.Tuple2[string, string]) (int, error) {
			intFocus, err := strconv.ParseInt(source.B, 10, 32)
			return int(intFocus), err
		},
		func(ctx context.Context, focus int, source lo.Tuple2[string, string]) (lo.Tuple2[string, string], error) {
			source.B = strconv.Itoa(focus)
			return source, nil
		},
		ExprCustom("ExampleLensR"),
	)

	//Note the result is an int even though the tuple element is a string.
	var viewResult int
	viewResult, err := Get(tupleB, lo.T2("10", "20"))
	fmt.Println(viewResult, err)

	//Note the result is lo.Tuple[string,string] but the Mul operation acts on the focused int.
	var overResult lo.Tuple2[string, string]
	overResult, err = Modify(tupleB, Mul[int](2), lo.T2("10", "20"))
	fmt.Println(overResult, err)

	_, err = Modify(tupleB, Mul(2), lo.T2("10", "twenty"))
	fmt.Println(err.Error())

	//Output:20 <nil>
	//{10 40} <nil>
	//strconv.ParseInt: parsing "twenty": invalid syntax
	//optic error path:
	//	Custom(ExampleLensR)
}

func ExampleLensP() {
	//This lens is polymorphic. It converts the B element of a tuple from int to float64s.
	intToFloat := LensP[lo.Tuple2[string, int], lo.Tuple2[string, float64], int, float64](
		func(source lo.Tuple2[string, int]) int {
			return source.B
		},
		func(focus float64, source lo.Tuple2[string, int]) lo.Tuple2[string, float64] {
			return lo.T2(source.A, focus)
		},
		ExprCustom("ExampleLensP"),
	)

	var result lo.Tuple2[string, float64] = MustModify(intToFloat, Op(func(focus int) float64 {
		return float64(focus) + 0.5
	}), lo.T2("alpha", 10))

	fmt.Println(result)
	//Output:{alpha 10.5}
}

func ExampleLensIEP() {
	//This lens converts the b elements of a string tuple to an int, reporting any conversion errors encountered.
	//It is polymorphic so is able to return an int tuple even through the source type is a string string
	parseInt := CombiLens[ReadWrite, Err, int, lo.Tuple2[string, string], lo.Tuple2[string, int], int, int](
		func(ctx context.Context, source lo.Tuple2[string, string]) (int, int, error) {
			intFocus, err := strconv.ParseInt(source.B, 10, 32)
			return 1, int(intFocus), err
		},
		func(ctx context.Context, focus int, source lo.Tuple2[string, string]) (lo.Tuple2[string, int], error) {
			return lo.T2(source.A, focus), nil
		},
		IxMatchComparable[int](),
		ExprCustom("ExampleLensF"),
	)

	var result lo.Tuple2[string, int]
	result, err := Modify(parseInt, Mul(10), lo.T2("alpha", "10"))

	fmt.Println(result, err)

	_, err = Modify(parseInt, Mul(10), lo.T2("alpha", "ten"))

	fmt.Println(err.Error())

	//Output:{alpha 100} <nil>
	//strconv.ParseInt: parsing "ten": invalid syntax
	//optic error path:
	//	Custom(ExampleLensF)
}
