package optic_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/mo"

	. "github.com/spearson78/go-optic"
)

func ExampleOp() {

	strLen := Op[string, int](func(s string) int {
		return len(s)
	})

	viewResult := MustGet(strLen, "hello world")
	fmt.Println(viewResult)

	//Ops are optics and can be composed.
	lenLt4 := Compose(strLen, Lt(4))

	data := []string{
		"lorem",
		"ipsum",
		"dolor",
		"sit",
		"amet",
	}

	filterResult := MustGet(SliceOf(Filtered(TraverseSlice[string](), lenLt4), len(data)), data)
	fmt.Println(filterResult)

	upper := Op(strings.ToUpper)
	upperResult := MustModify(TraverseSlice[string](), upper, data)
	fmt.Println(upperResult)

	//Output:
	//11
	//[sit]
	//[LOREM IPSUM DOLOR SIT AMET]
}

func ExampleOpE() {

	parseInt := OpE[string, int](func(ctx context.Context, s string) (int, error) {
		i, err := strconv.ParseInt(s, 10, 32)
		return int(i), err
	})

	viewResult, err := Get(parseInt, "11")
	fmt.Println(viewResult, err)

	lt4 := Compose(parseInt, Lt(4))

	data := []string{
		"1",
		"2",
		"3",
		"4",
		"5",
	}

	filterResult, err := Get(SliceOf(Filtered(TraverseSlice[string](), lt4), len(data)), data)
	fmt.Println(filterResult, err)

	var overResult []int
	overResult, err = Modify(TraverseSliceP[string, int](), parseInt, data)
	fmt.Println(overResult, err)

	badData := []string{
		"1",
		"2",
		"three",
		"4",
		"5",
	}

	_, err = Modify(TraverseSliceP[string, int](), parseInt, badData)
	fmt.Println(err.Error())

	//Output:
	//11 <nil>
	//[1 2 3] <nil>
	//[1 2 3 4 5] <nil>
	//strconv.ParseInt: parsing "three": invalid syntax
	//optic error path:
	//	func1
	//	ValueI[int,string].value
	//	Traverse
}

func ExampleOpI() {

	indexFizzBuzz := OpI[int, string, string](func(index int, s string) string {

		index = index + 1

		if index%15 == 0 {
			return "fizz buzz"
		}

		if index%3 == 0 {
			return "fizz"
		}

		if index%5 == 0 {
			return "buzz"
		}

		return s
	})

	data := []string{
		"alpha",
		"beta",
		"gamma",
		"delta",
		"epsilon",
		"zeta",
	}

	var overResult []string = MustModifyI(TraverseSlice[string](), indexFizzBuzz, data)
	fmt.Println(overResult)

	//Index Ops can be composed.
	eqFizz := ComposeLeft(indexFizzBuzz, Ne("fizz"))

	filterResult := MustGet(SliceOf(FilteredI(TraverseSlice[string](), eqFizz), len(data)), data)
	fmt.Println(filterResult)

	//Output:
	//[alpha beta fizz delta buzz fizz]
	//[alpha beta delta epsilon]
}

func ExampleOpIE() {

	parseIntWithIndex := OpIE[int, string, int](func(ctx context.Context, index int, s string) (int, error) {
		i, err := strconv.ParseInt(s, 10, 32)
		return int(i) + (index * 100), err
	})

	data := []string{
		"1",
		"2",
		"3",
		"4",
		"5",
	}

	var overResult []int
	overResult, err := ModifyI(TraverseSliceP[string, int](), parseIntWithIndex, data)
	fmt.Println(overResult, err)

	badData := []string{
		"1",
		"2",
		"three",
		"4",
		"5",
	}

	_, err = ModifyI(TraverseSliceP[string, int](), parseIntWithIndex, badData)
	fmt.Println(err.Error())

	//Output:
	//[1 102 203 304 405] <nil>
	//strconv.ParseInt: parsing "three": invalid syntax
	//optic error path:
	//	func1
	//	Traverse
}

func ExampleOperator() {

	strLen := Operator[string, int](
		func(s string) int {
			return len(s)
		},
		ExprCustom("ExampleOperator"),
	)

	viewResult := MustGet(strLen, "hello world")
	fmt.Println(viewResult)

	//Ops are optics and can be composed.
	lenLt4 := Compose(strLen, Lt(4))

	data := []string{
		"lorem",
		"ipsum",
		"dolor",
		"sit",
		"amet",
	}

	filterResult := MustGet(SliceOf(Filtered(TraverseSlice[string](), lenLt4), len(data)), data)
	fmt.Println(filterResult)

	upper := Op(strings.ToUpper)
	upperResult := MustModify(TraverseSlice[string](), upper, data)
	fmt.Println(upperResult)

	//Output:
	//11
	//[sit]
	//[LOREM IPSUM DOLOR SIT AMET]
}

func ExampleOperatorE() {

	parseInt := OperatorE[string, int](
		func(ctx context.Context, s string) (int, error) {
			i, err := strconv.ParseInt(s, 10, 32)
			return int(i), err
		},
		ExprCustom("ExampleOperatorR"),
	)

	viewResult, err := Get(parseInt, "11")
	fmt.Println(viewResult, err)

	lt4 := Compose(parseInt, Lt(4))

	data := []string{
		"1",
		"2",
		"3",
		"4",
		"5",
	}

	filterResult, err := Get(SliceOf(Filtered(TraverseSlice[string](), lt4), len(data)), data)
	fmt.Println(filterResult, err)

	var overResult []int
	overResult, err = Modify(TraverseSliceP[string, int](), parseInt, data)
	fmt.Println(overResult, err)

	badData := []string{
		"1",
		"2",
		"three",
		"4",
		"5",
	}

	_, err = Modify(TraverseSliceP[string, int](), parseInt, badData)
	fmt.Println(err.Error())

	//Output:
	//11 <nil>
	//[1 2 3] <nil>
	//[1 2 3 4 5] <nil>
	//strconv.ParseInt: parsing "three": invalid syntax
	//optic error path:
	//	Custom(ExampleOperatorR)
	//	ValueI[int,string].value
	//	Traverse
}

func ExampleOperatorI() {

	indexFizzBuzz := OperatorI[int, string, string](
		func(index int, s string) string {

			index = index + 1

			if index%15 == 0 {
				return "fizz buzz"
			}

			if index%3 == 0 {
				return "fizz"
			}

			if index%5 == 0 {
				return "buzz"
			}

			return s
		},
		ExprCustom("ExampleOperatorI"),
	)

	data := []string{
		"alpha",
		"beta",
		"gamma",
		"delta",
		"epsilon",
		"zeta",
	}

	var overResult []string = MustModifyI(TraverseSlice[string](), indexFizzBuzz, data)
	fmt.Println(overResult)

	//Index Ops can be composed.
	eqFizz := ComposeLeft(indexFizzBuzz, Ne("fizz"))

	filterResult := MustGet(SliceOf(FilteredI(TraverseSlice[string](), eqFizz), len(data)), data)
	fmt.Println(filterResult)

	//Output:
	//[alpha beta fizz delta buzz fizz]
	//[alpha beta delta epsilon]
}

func ExampleOperatorIE() {

	parseIntWithIndex := OperatorIE[int, string, int, Err](
		func(ctx context.Context, index int, s string) (int, error) {
			i, err := strconv.ParseInt(s, 10, 32)
			return int(i) + (index * 100), err
		},
		ExprCustom("ExampleOperatorF"))

	data := []string{
		"1",
		"2",
		"3",
		"4",
		"5",
	}

	var overResult []int
	overResult, err := ModifyI(TraverseSliceP[string, int](), parseIntWithIndex, data)
	fmt.Println(overResult, err)

	badData := []string{
		"1",
		"2",
		"three",
		"4",
		"5",
	}

	_, err = ModifyI(TraverseSliceP[string, int](), parseIntWithIndex, badData)
	fmt.Println(err.Error())

	//Output:
	//[1 102 203 304 405] <nil>
	//strconv.ParseInt: parsing "three": invalid syntax
	//optic error path:
	//	Custom(ExampleOperatorF)
	//	Traverse
}

func ExampleOpToOpI() {

	data := []int{2, 4, 6, 8}

	//This optic focuses each int in the slice
	optic := TraverseSlice[int]()

	//This operator will encode the slice index in the int in the 100s section
	op := AddOp(
		OpToOpI[int]( //IxOp enables non indexed operations to be combined with indexed operations.
			Identity[int](),
		),
		ComposeLeft(
			ValueIIndex[int, int](),
			Mul(100),
		),
	)

	var result []int = MustModifyI(
		optic,
		op,
		data,
	)
	fmt.Println(result)

	//Output:
	//[2 104 206 308]
}

func ExamplePredToOpI() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	//This optic focuses each int in the slice
	optic := TraverseSlice[string]()

	igt := GtI[string](0)
	lt := Ne("gamma")

	//Indexed and non indexed predicates can be mixed by wrapping the non indexed operator
	andOp := AndOp(igt, PredToOpI[int](lt))

	filter := FilteredI(optic, andOp)

	var result []string = MustGet(SliceOf(filter, len(data)), data)
	fmt.Println(result)

	//Output:
	//[beta delta]
}

func ExampleOpT2ToOpT2I() {

	data := []int{2, 4, 6, 8}

	//This optic focuses each int in the slice
	optic := TraverseSlice[int]()

	//This operator will encode the slice index of the int in the 100s section
	op := OpT2IToOptic(
		OpToOpI[int]( //OpToOpI enables Identity to be used with indexed operators
			Identity[int](),
		),
		OpT2ToOpT2I[Void, int](AddT2[int]()), //OpT2ToOpT2I enables Add to be used with indexed operators
		ComposeLeft(
			ValueIIndex[int, int](),
			Mul(100),
		),
	)

	var result []int = MustModifyI(
		optic,
		op,
		data,
	)
	fmt.Println(result)

	//Output:
	//[2 104 206 308]
}

func ExampleAsIxGet() {

	data := map[int]string{
		1: "alpha",
		2: "beta",
		3: "gamma",
		4: "delta",
	}

	optic := AsIxGet(TraverseMap[int, string]())

	var result mo.Option[string] = MustGet(optic, ValI(1, data))
	fmt.Println(result)

	result = MustGet(optic, ValI(10, data))
	fmt.Println(result)

	//Output:
	//{true alpha}
	//{false }

}

func ExampleOpOnIx() {

	data := []int{2, 4, 6, 8}

	//This optic focuses each int in the slice
	optic := TraverseSlice[int]()

	//This operator will encode the slice index in the int in the 100s section
	op := AddOp(
		OpToOpI[int]( //IxOp enables identity to be combined with other indexed operations
			Identity[int](),
		),
		ComposeLeft(
			OpOnIx[int]( //OpOnIx converts Identity to operate on the index instead of the focus
				Identity[int](),
			),
			Mul(100),
		),
	)

	var result []int = MustModifyI(
		optic,
		op,
		data,
	)
	fmt.Println(result)

	//Output:
	//[2 104 206 308]
}

func ExamplePredOnIx() {

	data := []int{
		10,
		20,
		30,
		40,
	}

	//PredIx converts Even to operate on the index rather than the focused value
	evenIndex := PredOnIx[int](Even[int]())

	result, err := Modify(FilteredI(TraverseSlice[int](), evenIndex), Mul(10), data)
	fmt.Println(result, err)

	//Output:
	//[100 20 300 40] <nil>
}
