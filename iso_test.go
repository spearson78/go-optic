package optic_test

import (
	"context"
	"fmt"
	"strconv"

	. "github.com/spearson78/go-optic"
)

func ExampleInvoluted() {

	//Negating a value is its own inverse. We can use an involuted here.
	negate := Involuted[int](
		func(i int) int {
			return -i
		},
		ExprCustom("ExampleInvoluted"),
	)

	viewResult := MustGet(negate, 1)
	fmt.Println(viewResult)

	//For involuted isos ReView is equivalent to View
	reViewResult := MustReverseGet(negate, 1)
	fmt.Println(reViewResult)

	data := []int{1, 2, 3}

	//The Add(1) is operating in a negated context the negative value is converted back to a positive number in the result.
	//This has the effect of making Add(1) subtract 1 from the value
	overResult := MustModify(Compose(TraverseSlice[int](), negate), Add(1), data)
	fmt.Println(overResult)

	//Output:
	//-1
	//-1
	//[0 1 2]
}

func ExampleIso() {

	//Converting between celsius an iso is lossless.
	toFahrenheit := Iso[float64, float64](
		func(celsius float64) float64 {
			return (celsius * (9.0 / 5.0)) + 32.0
		},
		func(fahrenheit float64) float64 {
			return (fahrenheit - 32.0) * (5.0 / 9.0)
		},
		ExprCustom("ExampleIso"),
	)

	viewResult := MustGet(toFahrenheit, 32.0)
	fmt.Println(viewResult)

	//ReView peforms the reverse fahrenheit to celsius comversion
	reViewResult := MustReverseGet(toFahrenheit, 89.6)
	fmt.Println(reViewResult)

	data := map[string]float64{
		"freeze": 0.0,
		"body":   37.0,
		"boil":   100.0,
	}

	//Isos can be composed to perform operations within the converted context.
	//In this case Add(1.0) will add 1 fahrenheit to the values
	//But the result is returned back to celsius
	overComposedResult := MustModify(Compose(TraverseMap[string, float64](), toFahrenheit), Add(1.0), data)
	fmt.Println(overComposedResult)

	//Isos can be used directly as the operation in the Over action
	//This converts the values to fahrenheit in the result
	overResult := MustModify(TraverseMap[string, float64](), toFahrenheit, data)
	fmt.Println(overResult)

	//Output:
	//89.6
	//32
	//map[body:37.555555555555564 boil:100.55555555555556 freeze:0.5555555555555556]
	//map[body:98.60000000000001 boil:212 freeze:32]
}

func ExampleIsoE() {

	//Converting between string and int is lossless in the cases it can be parsed
	//An IsoR can be used to handle the error conditions
	parseInt := IsoE[string, int](
		func(ctx context.Context, source string) (int, error) {
			intVal, err := strconv.ParseInt(source, 10, 32)
			return int(intVal), err
		},
		func(ctx context.Context, i int) (string, error) {
			return strconv.Itoa(i), nil
		},
		ExprCustom("ExampleIsoR"),
	)

	viewResult, err := Get(parseInt, "1")
	fmt.Println(viewResult, err)

	_, err = Get(parseInt, "one")
	fmt.Println(err.Error())

	//ReView performs the reverse int to string
	reViewResult, err := ReverseGet(parseInt, 1)
	fmt.Println(reViewResult, err)

	data := []string{"1", "2", "3"}

	//Isos can be composed to perform operations within the converted context.
	//In this case Add(1) will add 1 the int values
	//But the result is returned back to string
	var overComposedResult []string
	overComposedResult, err = Modify(Compose(TraverseSlice[string](), parseInt), Add(1), data)
	fmt.Println(overComposedResult, err)

	//Isos can be used directly as the operation in the Over action
	//This converts the values to int in the result
	var overResult []int
	overResult, err = Modify(TraverseSliceP[string, int](), parseInt, data)
	fmt.Println(overResult, err)

	//Output:
	//1 <nil>
	//strconv.ParseInt: parsing "one": invalid syntax
	//optic error path:
	//	Custom(ExampleIsoR)
	//
	//1 <nil>
	//[2 3 4] <nil>
	//[1 2 3] <nil>
}

func ExampleIsoP() {

	//Polymorphic Isos can be useful when combined with other polymorphic optics
	toFloat := IsoP[int, float64, float64, float64](
		//getter
		func(source int) float64 {
			return float64(source)
		},
		//reverse
		func(focus float64) float64 {
			return focus
		},
		ExprCustom("ExampleIsoP"),
	)

	var viewResult float64 = MustGet(toFloat, 1)
	fmt.Println(viewResult)

	//ReView passes calls the reverse function.
	var reViewResult float64 = MustReverseGet(toFloat, 1.5)
	fmt.Println(reViewResult)

	//The polymorphic traversal expects a conversion from int to float64. A non polymorphic iso would convert back to the source type
	var overResult []float64 = MustModify(Compose(TraverseSliceP[int, float64](), toFloat), Add(0.5), []int{1, 2, 3})
	fmt.Println(overResult)

	//Output:
	//1
	//1.5
	//[1.5 2.5 3.5]
}

func ExampleIsoEP() {

	//Polymorphic Isos can be useful when combined with other polymorphic optics
	parseInt := IsoEP[string, int, int, int](
		func(ctx context.Context, source string) (int, error) {
			intVal, err := strconv.ParseInt(source, 10, 32)
			return int(intVal), err
		},
		func(ctx context.Context, focus int) (int, error) {
			return focus, ctx.Err()
		},
		ExprCustom("ExampleIsoF"),
	)

	var viewResult int
	viewResult, err := Get(parseInt, "1")
	fmt.Println(viewResult, err)

	//ReView passes calls the reverse function.
	var reViewResult int
	reViewResult, err = ReverseGet(parseInt, 1)
	fmt.Println(reViewResult, err)

	//The polymorphic traversal expects a conversion from int to string. A non polymorphic iso would convert back to []string
	var overResult []int
	overResult, err = Modify(Compose(TraverseSliceP[string, int](), parseInt), Add(1), []string{"1", "2", "3"})
	fmt.Println(overResult, err)

	//Output:
	//1 <nil>
	//1 <nil>
	//[2 3 4] <nil>
}
