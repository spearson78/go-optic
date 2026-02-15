package optic_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/samber/mo"

	. "github.com/spearson78/go-optic"
)

func FuzzTraverseSlice(f *testing.F) {
	f.Add(uint8(0), int64(0), int(10))
	f.Add(uint8(1), int64(1), int(20))

	f.Fuzz(func(t *testing.T, len uint8, seed int64, newVal int) {

		r := rand.New(rand.NewSource(seed))
		data := make([]int, 0, len)
		for i := byte(0); i < len; i++ {
			data = append(data, r.Int())
		}

		ValidateOpticTestPred(t, TraverseSlice[int](), data, newVal, EqDeepT2[[]int]())

	})
}

func ExampleTraverseSlice() {

	slice := []int{5, 10, 15}

	//Technically multiple elements of an Traversal may have the same index so Index focuses multiple values. MustGetFirst selects the first one.
	elemOf, found := MustGetFirst(Index(TraverseSlice[int](), 1), slice)
	fmt.Println(elemOf, found)

	result := MustModify(TraverseSlice[int](), Mul(10), slice)
	fmt.Println(result)

	//Output:10 true
	//[50 100 150]
}

func ExampleTraverseSliceP() {

	slice := []string{"1", "2", "3"}

	var overResult []int
	overResult, err := Modify(
		TraverseSliceP[string, int](),
		ParseInt[int](10, 32),
		slice,
	)
	fmt.Println(overResult, err)

	//Output:[1 2 3] <nil>
}

func ExampleAtSlice() {

	data := []int{1, 2, 3}

	optic := AtSlice[int](1, EqT2[int]())

	var getResult mo.Option[int] = MustGet(
		optic,
		data,
	)
	fmt.Println(getResult)

	var setResult []int = MustSet(
		optic,
		mo.Some(10),
		data,
	)
	fmt.Println(setResult)

	optic = AtSlice[int](10, EqT2[int]())

	getResult = MustGet(
		optic,
		data,
	)
	fmt.Println(getResult)

	setResult = MustSet(
		optic,
		mo.Some(10),
		data,
	)
	fmt.Println(setResult)

	//Output:
	//{true 2}
	//[1 10 3]
	//{false 0}
	//[1 2 3 0 0 0 0 0 0 0 10]
}

func ExampleSliceToCol() {

	data := lo.T2(1, []string{
		"alpha",
		"beta",
		"gamma",
		"delta",
	})

	//See: [FilteredSlice] for a more convenient slice filter function.
	var getRes []string = MustGet(
		Compose4(
			T2B[int, []string](),
			SliceToCol[string](),
			FilteredCol[int](Ne("beta")),
			AsReverseGet(SliceToCol[string]()),
		),
		data,
	)
	fmt.Println(getRes)

	var modifyRes lo.Tuple2[int, []string] = MustModify(
		Compose(
			T2B[int, []string](),
			SliceToCol[string](),
		),
		FilteredCol[int](Ne("beta")),
		data,
	)
	fmt.Println(modifyRes)

	//Output:
	//[alpha gamma delta]
	//{1 [alpha gamma delta]}
}

func ExampleSliceColTypeP() {

	data := lo.T2(1, []string{
		"1",
		"2",
		"3",
		"4",
	})

	var result lo.Tuple2[int, []int]
	var err error
	result, err = Modify(
		Compose(
			T2BP[int, []string, []int](),
			TraverseColType(SliceColTypeP[string, int]()),
		),
		Compose(
			ParseInt[int](10, 0),
			Mul(2),
		),
		data,
	)

	fmt.Println(result, err)

	// Output:
	//{1 [2 4 6 8]} <nil>
}

func ExampleSliceTail() {

	data := []string{
		"alpha",
		"beta",
		"gamma",
		"delta",
	}

	res, ok := MustGetFirst(
		SliceTail[string](),
		data,
	)
	fmt.Println(res, ok)

	var modifyRes []string = MustModify(
		Compose(
			SliceTail[string](),
			TraverseSlice[string](),
		),
		Op(strings.ToUpper),
		data,
	)
	fmt.Println(modifyRes)

	// Output:
	// [beta gamma delta] true
	// [alpha BETA GAMMA DELTA]
}

func TestAtSliceNegativeIndex(t *testing.T) {

	if res, err := Get(
		AtSlice[int](-2, EqT2[int]()),
		[]int{1, 2, 3},
	); err != nil || res.MustGet() != 2 {
		t.Fatalf("%v , %v", res, err)
	}

	if res, err := Get(
		AtSlice[int](-4, EqT2[int]()),
		[]int{1, 2, 3},
	); err != nil || res.IsPresent() {
		t.Fatalf("%v , %v", res, err)
	}

	if res, err := Set(
		AtSlice[int](-2, EqT2[int]()),
		mo.Some(99),
		[]int{1, 2, 3},
	); err != nil || !reflect.DeepEqual(res, []int{1, 99, 3}) {
		t.Fatalf("%v , %v", res, err)
	}

	if res, err := Set(
		AtSlice[int](-5, EqT2[int]()),
		mo.Some(99),
		[]int{1, 2, 3},
	); err != nil || !reflect.DeepEqual(res, []int{99, 0, 1, 2, 3}) {
		t.Fatalf("%v , %v", res, err)
	}

}

func ExampleSliceOf() {

	data := map[int]string{
		1: "alpha",
		2: "beta",
		3: "gamma",
		4: "delta",
	}

	optic := SliceOf(TraverseMap[int, string](), 10)

	viewResult := MustGet(optic, data)
	fmt.Println(viewResult)

	modifyResult := MustModify(optic, Op(func(focus []string) []string {
		slices.Reverse(focus)
		return focus
	}), data)
	fmt.Println(modifyResult)

	//Output:
	//[alpha beta gamma delta]
	//map[1:delta 2:gamma 3:beta 4:alpha]

}

func ExampleSliceOfP() {

	data := map[int]string{
		1: "alpha",
		2: "beta",
		3: "gamma",
		4: "delta",
	}

	//Polymorphic map traversal with a source type of map[int]string and result type of map[int]int
	optic := SliceOfP(TraverseMapP[int, string, int](), len(data))

	var viewResult []string
	viewResult, err := Get(optic, data)
	fmt.Println(viewResult, err)

	//Note the return type is map[int]int not map[int]string
	var modifyResult map[int]int
	modifyResult, err = Modify(optic, Op(func(focus []string) []int {
		var modified []int
		for _, str := range focus {
			modified = append(modified, len(str))
		}

		return modified
	}), data)
	fmt.Println(modifyResult, err)

	//Output:
	//[alpha beta gamma delta] <nil>
	//map[1:5 2:4 3:5 4:5] <nil>
}
