package optic_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	. "github.com/spearson78/go-optic"
)

func TestCollectionOpComposition(t *testing.T) {

	//I was concerned about the Optic Err getting out of sync with the

	data := 1215435

	optic := Compose(
		AsReverseGet(ParseInt[int](10, 0)),
		ColFocusErr(StringToCol()),
	)

	op := Compose(
		FilteredCol[int](
			EErr(Compose4(
				SliceOf(Identity[rune](), 1),
				IsoCast[[]rune, string](),
				ParseInt[int](10, 0),
				Odd[int](),
			)),
		),
		OrderedCol[int](
			EErr(OrderBy(Identity[rune]())),
		),
	)

	res, err := Modify(optic, op, data)

	if res != 11355 || err != nil {
		t.Fatal(res, err)
	}

}

func ExampleReverse() {

	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var result []int = MustModify(
		SliceToCol[int](),
		ReversedCol[int, int](),
		slice,
	)
	fmt.Println(result)

	//Output:
	//[10 9 8 7 6 5 4 3 2 1]
}

func ExampleReversedColP() {

	slice := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	var result []int
	result, err := Get(
		SliceOfP(
			Compose4(
				SliceToColP[string, int](),
				ReversedColP[int, string, int](),
				TraverseColP[int, string, int](),
				ParseIntP[int](10, 32),
			),
			10,
		),
		slice,
	)
	fmt.Println(result, err)

	//Output:
	//[10 9 8 7 6 5 4 3 2 1] <nil>
}

func ExampleFilteredCol() {

	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var result []int = MustModify(
		SliceToCol[int](),
		FilteredCol[int, int](
			Odd[int](),
		),
		slice,
	)

	fmt.Println(result)

	//Output:
	//[1 3 5 7 9]
}

func ExampleFilteredColI() {

	data := map[int]string{
		1: "alpha",
		2: "beta",
		3: "gamma",
		4: "delta",
	}

	var result map[int]string = MustModify(
		MapToCol[int, string](),
		FilteredColI[int, string](
			OpOnIx[string](
				Odd[int](),
			),
			IxMatchComparable[int](),
		),
		data,
	)

	fmt.Println(result)

	//Output:
	//map[1:alpha 3:gamma]
}

func ExampleAppend() {

	data := []string{
		"alpha",
		"beta",
	}

	var result []string = MustModify(
		SliceToCol[string](),
		Append("gamma", "delta"),
		data,
	)

	fmt.Println(result)

	//Output:
	//[alpha beta gamma delta]
}

func ExampleAppendCol() {

	data := []string{
		"alpha",
		"beta",
	}

	var result []string = MustModify(
		SliceToCol[string](),
		AppendCol(
			ValCol("gamma", "delta"),
		),
		data,
	)

	fmt.Println(result)

	//Output:
	//[alpha beta gamma delta]
}

func ExamplePrependCol() {

	data := []string{
		"gamma",
		"delta",
	}

	var result []string = MustModify(
		SliceToCol[string](),
		PrependCol(
			ValCol("alpha", "beta"),
		),
		data,
	)

	fmt.Println(result)

	//Output:
	//[alpha beta gamma delta]
}

func ExamplePrepend() {

	data := []string{
		"gamma",
		"delta",
	}

	var result []string = MustModify(
		SliceToCol[string](),
		Prepend("alpha", "beta"),
		data,
	)

	fmt.Println(result)

	//Output:
	//[alpha beta gamma delta]
}

func ExampleSort() {

	data := []int{3, 4, 7, 9, 10, 8, 2, 1, 5, 6}

	var result []int = MustModify(
		SliceToCol[int](),
		OrderedCol[int](
			OrderBy(
				Identity[int](),
			),
		),
		data,
	)

	fmt.Println(result)

	//Output:
	//[1 2 3 4 5 6 7 8 9 10]
}

func ExampleOrderedColI() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	var result []string = MustModify(
		SliceToCol[string](),
		OrderedColI[int, string](
			OrderByI(
				OpOnIx[string](
					Negate[int](),
				),
			),
			IxMatchComparable[int](),
		),
		data,
	)

	fmt.Println(result)

	//Output:
	//[delta gamma beta alpha]
}

func ExampleSubCol() {

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var result []int = MustModify(
		SliceToCol[int](),
		SubCol[int, int](2, 2),
		data,
	)

	fmt.Println(result)

	var negativeLengthResult []int = MustModify(
		SliceToCol[int](),
		SubCol[int, int](2, -2),
		data,
	)

	fmt.Println(negativeLengthResult)

	//Output:
	//[3 4]
	//[3 4 5 6 7 8]
}

func ExampleReIndexedCol() {

	data := ValCol("alpha", "beta", "gamma", "delta")

	optic := ReIndexedCol[string](Add(1))

	getRes := MustGet(optic, data)
	fmt.Println(getRes)

	modifyRes := MustModify(
		optic,
		FilteredColI(
			OpOnIx[string](
				Odd[int](),
			),
			IxMatchComparable[int](),
		),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//Col[1:alpha 2:beta 3:gamma 4:delta]
	//Col[0:alpha 2:gamma]
}

func ExampleReIndexedColI() {

	data := ValCol("alpha", "beta", "gamma", "delta")

	optic := ReIndexedColI[string](
		Compose(
			FirstOrDefault(Ignore(AsReverseGet(ParseInt[int](10, 0)), True[error]()), ""),
			PrependString(StringCol("index_")),
		),
		EqT2[int](),
		EqT2[string](),
	)

	getRes := MustGet(optic, data)
	fmt.Println(getRes)

	//Output:
	//Col[index_0:alpha index_1:beta index_2:gamma index_3:delta]
}

func TestAsColOpComposition(t *testing.T) {

	data := []string{"alpha", "beta", "gamma", "delta"}

	composition := EPure(ComposeLeft(
		Ordered(
			Filtered(
				TraverseCol[int, string](),
				Ne("beta"),
			),
			Desc(OrderBy(Identity[string]())),
		),
		Op(strings.ToUpper),
	))

	var result []string = MustModify(
		SliceToCol[string](),
		ColOf(composition),
		data,
	)

	if !reflect.DeepEqual(result, []string{"GAMMA", "DELTA", "ALPHA"}) {
		t.Fatal(result)
	}

	result = MustModify(
		Compose(
			SliceToCol[string](),
			Taking(
				Ordered(
					Filtered(
						TraverseCol[int, string](),
						Ne("beta"),
					),
					Desc(OrderBy(Identity[string]())),
				),
				2,
			),
		),
		Op(strings.ToUpper),
		data,
	)

	if !reflect.DeepEqual(result, []string{"alpha", "beta", "GAMMA", "DELTA"}) {
		t.Fatal(result)
	}

}
