package optic_test

import (
	"context"
	"fmt"
	"strconv"

	"github.com/samber/lo"

	. "github.com/spearson78/go-optic"
)

func ExampleGetter() {
	//This Getter focuses on the B element of a tuple
	tupleB := Getter[lo.Tuple2[int, int], int](
		func(source lo.Tuple2[int, int]) int {
			return source.B
		},
		ExprCustom("ExampleGetter"),
	)

	var viewResult int = MustGet(tupleB, lo.T2(10, 20))
	fmt.Println(viewResult)

	//Output:20
}

func ExampleGetterI() {
	//This Getter focuses on the B element of a tuple with an index
	tupleB := GetterI[int, lo.Tuple2[int, int], int](
		func(source lo.Tuple2[int, int]) (int, int) {
			return 1, source.B
		},
		func(indexA, indexB int) bool {
			return indexA == indexB
		},
		ExprCustom("ExampleGetterI"),
	)

	var viewIndex, viewResult int = MustGetI(tupleB, lo.T2(10, 20))
	fmt.Println(viewIndex, viewResult)

	//Output:
	//1 20
}

func ExampleGetterE() {
	//This Getter focuses on the B element of a string tuple and converts it to an int
	//reporting the error in case of conversion failure
	tupleB := GetterE[lo.Tuple2[string, string], int](
		func(ctx context.Context, source lo.Tuple2[string, string]) (int, error) {
			intFocus, err := strconv.ParseInt(source.B, 10, 32)
			return int(intFocus), err
		},
		ExprCustom("ExampleGetterR"),
	)

	//Note the result is an int even though the tuple element is a string.
	var viewResult int
	viewResult, err := Get(tupleB, lo.T2("10", "20"))
	fmt.Println(viewResult, err)

	_, err = Get(tupleB, lo.T2("10", "twenty"))
	fmt.Println(err)

	//Output:20 <nil>
	//strconv.ParseInt: parsing "twenty": invalid syntax
	//optic error path:
	//	Custom(ExampleGetterR)
}

func ExampleGetterIE() {
	//This Getter converts the b elements of a string tuple to an int, reporting any conversion errors encountered.
	//It is polymorphic so is able to return an int tuple even through the source type is a string string
	parseInt := CombiGetter[Err, int, lo.Tuple2[string, string], lo.Tuple2[string, string], int, int](
		func(ctx context.Context, source lo.Tuple2[string, string]) (int, int, error) {
			intFocus, err := strconv.ParseInt(source.B, 10, 32)
			return 1, int(intFocus), err
		},
		IxMatchComparable[int](),
		ExprCustom("ExampleGetterF"),
	)

	var result int
	result, err := Get(parseInt, lo.T2("alpha", "10"))

	fmt.Println(result, err)

	_, err = Get(parseInt, lo.T2("alpha", "ten"))

	fmt.Println(err.Error())

	//Output:
	//10 <nil>
	//strconv.ParseInt: parsing "ten": invalid syntax
	//optic error path:
	//	Custom(ExampleGetterF)
}
