package optic_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/samber/lo"

	. "github.com/spearson78/go-optic"
)

func ExampleSliceOp() {

	data := lo.T2(
		1,
		[]string{
			"alpha",
			"beta",
			"gamma",
			"delta",
		},
	)

	optic := SliceOp(
		SubCol[int, string](1, -1),
	)

	var getRes []string = MustGet(
		Compose(
			T2B[int, []string](),
			optic,
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, []string] = MustModify(
		T2B[int, []string](),
		optic,
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//[beta gamma]
	//{1 [beta gamma]}
}

func ExampleFilteredSlice() {

	data := lo.T2(
		1,
		[]string{
			"alpha",
			"beta",
			"gamma",
			"delta",
		},
	)

	var getRes []string = MustGet(
		Compose(
			T2B[int, []string](),
			FilteredSlice(Ne("beta")),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, []string] = MustModify(
		T2B[int, []string](),
		FilteredSlice(Ne("beta")), //FilteredSlice can be used as an operation
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//[alpha gamma delta]
	//{1 [alpha gamma delta]}
}

func ExampleFilteredSliceI() {

	data := lo.T2(
		1,
		[]string{
			"alpha",
			"beta",
			"gamma",
			"delta",
		},
	)

	var getRes []string = MustGet(
		Compose(
			T2B[int, []string](),
			FilteredSliceI(NeI[string](1)),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, []string] = MustModify(
		T2B[int, []string](),
		FilteredSliceI(NeI[string](1)), //FilteredSliceI can be used as an operation
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//[alpha gamma delta]
	//{1 [alpha gamma delta]}
}

func ExampleAppendSlice() {

	data := lo.T2(
		1,
		[]string{
			"alpha",
			"beta",
		},
	)

	var getRes []string = MustGet(
		Compose(
			T2B[int, []string](),
			AppendSlice(ValCol("gamma", "delta")),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, []string] = MustModify(
		T2B[int, []string](),
		AppendSlice(ValCol("gamma", "delta")),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//[alpha beta gamma delta]
	//{1 [alpha beta gamma delta]}
}

func ExampleReversedSlice() {

	data := lo.T2(
		1,
		[]string{
			"alpha",
			"beta",
			"gamma",
			"delta",
		},
	)

	var getRes []string = MustGet(
		Compose(
			T2B[int, []string](),
			ReversedSlice[string](),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, []string] = MustModify(
		T2B[int, []string](),
		ReversedSlice[string](),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//[delta gamma beta alpha]
	//{1 [delta gamma beta alpha]}
}

func ExamplePrependSlice() {

	data := lo.T2(
		1,
		[]string{
			"gamma",
			"delta",
		},
	)

	var getRes []string = MustGet(
		Compose(
			T2B[int, []string](),
			PrependSlice(ValCol("alpha", "beta")),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, []string] = MustModify(
		T2B[int, []string](),
		PrependSlice(ValCol("alpha", "beta")),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//[alpha beta gamma delta]
	//{1 [alpha beta gamma delta]}
}

func ExampleOrderedSlice() {

	data := lo.T2(
		1,
		[]string{
			"alpha",
			"beta",
			"gamma",
			"delta",
		},
	)

	var getRes []string = MustGet(
		Compose(
			T2B[int, []string](),
			OrderedSlice(OrderBy(Identity[string]())),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, []string] = MustModify(
		T2B[int, []string](),
		OrderedSlice(OrderBy(Identity[string]())),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//[alpha beta delta gamma]
	//{1 [alpha beta delta gamma]}
}

func ExampleOrderedSliceI() {

	data := lo.T2(
		1,
		[]string{
			"alpha",
			"beta",
			"gamma",
			"delta",
		},
	)

	var getRes []string = MustGet(
		Compose(
			T2B[int, []string](),
			OrderedSliceI(
				Desc(
					OrderByI(
						ValueIIndex[int, string](),
					),
				),
			),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, []string] = MustModify(
		T2B[int, []string](),
		OrderedSliceI(Desc(OrderByI(ValueIIndex[int, string]()))),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//[delta gamma beta alpha]
	//{1 [delta gamma beta alpha]}
}

func ExampleSubSlice() {

	data := lo.T2(
		1,
		[]string{
			"alpha",
			"beta",
			"gamma",
			"delta",
		},
	)

	var getRes []string = MustGet(
		Compose(
			T2B[int, []string](),
			SubSlice[string](1, -1),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, []string] = MustModify(
		T2B[int, []string](),
		SubSlice[string](1, -1),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//[beta gamma]
	//{1 [beta gamma]}
}

func TestFilteredSlice(t *testing.T) {

	slice := [][]int{
		{1, 2, 3},
		{4},
		{5, 6},
	}

	//Delete any elements > 3
	if result := MustModify(TraverseSlice[[]int](), FilteredSlice(Lte(3)), slice); !reflect.DeepEqual(result, [][]int{
		{1, 2, 3},
		nil,
		nil,
	}) {
		t.Fatalf("Modify FilteredSlice %v", result)
	}

	if result := MustGet(SliceOf(Compose(TraverseSlice[[]int](), FilteredSlice(Lte(3))), 3), slice); !reflect.DeepEqual(result, [][]int{
		{1, 2, 3},
		nil,
		nil,
	}) {
		t.Fatalf("MustToSliceOf %v", result)
	}
}
