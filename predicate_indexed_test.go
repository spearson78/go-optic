package optic_test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/spearson78/go-optic"
)

func TestIxPredicateConsistency(t *testing.T) {

	ValidateOpticTest(t, EqI[int, int](5), ValI(1, 2), false)

}

func ExampleEqI() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	//This optic focuses each string in the slice
	optic := TraverseSlice[string]()

	var result []string
	result = MustGet(
		SliceOf(
			FilteredI(

				optic, EqI[string](2),
			),
			len(data),
		),
		data,
	)
	fmt.Println(result)

	//Output:
	//[gamma]
}

func ExampleNeI() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	//This optic focuses each string in the slice
	optic := TraverseSlice[string]()

	var result []string
	result = MustGet(
		SliceOf(
			FilteredI(

				optic, NeI[string](2),
			),
			len(data),
		),
		data,
	)
	fmt.Println(result)

	//Output:
	//[alpha beta delta]
}

func ExampleGtI() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	//This optic focuses each string in the slice
	optic := TraverseSlice[string]()

	var result []string
	result = MustGet(
		SliceOf(
			FilteredI(

				optic, GtI[string](2),
			),
			len(data),
		),
		data,
	)
	fmt.Println(result)

	//Output:
	//[delta]
}

func ExampleGteI() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	//This optic focuses each string in the slice
	optic := TraverseSlice[string]()

	var result []string
	result = MustGet(
		SliceOf(
			FilteredI(

				optic, GteI[string](2),
			),
			len(data),
		),
		data,
	)
	fmt.Println(result)

	//Output:
	//[gamma delta]
}

func ExampleLtI() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	//This optic focuses each string in the slice
	optic := TraverseSlice[string]()

	var result []string
	result = MustGet(
		SliceOf(
			FilteredI(

				optic, LtI[string](2),
			),
			len(data),
		),
		data,
	)
	fmt.Println(result)

	//Output:
	//[alpha beta]
}

func ExampleLteI() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	//This optic focuses each string in the slice
	optic := TraverseSlice[string]()

	var result []string
	result = MustGet(
		SliceOf(
			FilteredI(

				optic, LteI[string](2),
			),
			len(data),
		),
		data,
	)
	fmt.Println(result)

	//Output:
	//[alpha beta gamma]
}

func ExampleAnyI() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	//This optic focuses each string in the slice
	optic := TraverseSlice[string]()

	var result bool
	var err error
	result, err = Get(AnyI(optic, EqI[string](2)), data)
	fmt.Println(result, err)

	//Output:
	//true <nil>
}

func ExampleAllI() {

	data := []int{1, 20, 3, 40}

	//This optic focuses each int in the slice
	optic := TraverseSlice[int]()

	var result bool
	var err error
	result, err = Get(AllI(optic, OpI(func(index int, focus int) bool {
		//even index or value is > 5
		return index%2 == 0 || focus > 5
	})), data)
	fmt.Println(result, err)

	//Output:
	//true <nil>
}

func TestEvenIndex(t *testing.T) {

	data := []int{1, 20, 3, 40}

	//This optic focuses each int in the slice
	optic := TraverseSlice[int]()

	even := Even[int]()
	evenIndex := OpOnIx[int](even)

	if result := MustGet(SliceOf(FilteredI(optic, evenIndex), len(data)), data); !reflect.DeepEqual(result, []int{1, 3}) {
		t.Fatal("EvenIndex", result)
	}
}
