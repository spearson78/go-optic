package optic_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unicode"

	"github.com/samber/lo"

	. "github.com/spearson78/go-optic"
)

func TestCombinatorIndexedConsistency(t *testing.T) {

	o := TraverseSlice[int]()

	ValidateOpticTestPred(t, TakingWhileI(o, OpOnIx[int](Lte(3))), []int{1, 2, 30, 40, 50}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, TakingWhileI(o, OpOnIx[int](False[int]())), []int{1, 2, 30, 40, 50}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, TakingWhileI(o, OpOnIx[int](True[int]())), []int{1, 2, 30, 40, 50}, 10, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, DroppingWhileI(o, OpOnIx[int](Lte(3))), []int{1, 2, 30, 40, 50}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, DroppingWhileI(o, OpOnIx[int](False[int]())), []int{1, 2, 30, 40, 50}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, DroppingWhileI(o, OpOnIx[int](True[int]())), []int{1, 2, 30, 40, 50}, 10, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, TrimmingWhileI(o, OpOnIx[int](OrOp(Lt(1), Gt(3)))), []int{1, 2, 30, 2, 1}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, TrimmingWhileI(o, OpOnIx[int](False[int]())), []int{1, 2, 30, 2, 1}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, TrimmingWhileI(o, OpOnIx[int](True[int]())), []int{1, 2, 30, 2, 1}, 10, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, FilteredI(o, OpOnIx[int](OrOp(Lt(1), Gt(3)))), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, FilteredI(o, OpOnIx[int](False[int]())), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, FilteredI(o, OpOnIx[int](True[int]())), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, Index(o, 2), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, Index(o, 100), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, WithIndex(o), []int{1}, ValI(0, 10), EqDeepT2[[]int]())

	ValidateOpticTestPred(t, Indices(o, OrOp(Lt(1), Gt(3))), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, Indices(o, False[int]()), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, Indices[int](o, True[int]()), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, Indexing(o), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, ReIndexed(o, Mul(2), EqT2[int]()), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, SelfIndex(o, EqT2[int]()), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())

	ValidateOpticTestPred(
		t,
		OrderedI(
			o,
			Desc(
				OrderByI(
					OpOnIx[int](
						Identity[int](),
					),
				),
			),
		),
		[]int{5, 9, 3, 4, 2, 1, 8, 10, 6, 7},
		10,
		EqDeepT2[[]int](),
	)

	ValidateOpticTestPred(t, FirstOrDefaultI(o, -1, 10), []int{1, 20, 3}, 0, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, FirstOrDefaultI(o, -1, 10), []int{}, 0, EqDeepT2[[]int]())
}

func ExampleTakingWhileI() {

	data := map[int]string{
		1: "alpha",
		2: "beta",
		3: "gamma",
		4: "delta",
	}

	//Focus on elements until the index 3 is found.
	optic := TakingWhileI(

		TraverseMap[int, string](), OpI(func(index int, focus string) bool {
			return index < 3
		}))

	listResult := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(listResult)

	overResult := MustModify(optic, Op(strings.ToUpper), data)
	fmt.Println(overResult)

	//Output:
	//[alpha beta]
	//map[1:ALPHA 2:BETA 3:gamma 4:delta]
}

func ExampleDroppingWhileI() {

	data := map[int]string{
		1: "alpha",
		2: "beta",
		3: "gamma",
		4: "delta",
	}

	//drop elements until the index 3 is found.
	optic := DroppingWhileI(

		TraverseMap[int, string](), OpI(func(index int, focus string) bool {
			return index < 3
		}))

	listResult := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(listResult)

	overResult := MustModify(optic, Op(strings.ToUpper), data)
	fmt.Println(overResult)

	//Output:
	//[gamma delta]
	//map[1:alpha 2:beta 3:GAMMA 4:DELTA]
}

func ExampleTrimmingWhileI() {

	data := "Lorem ipsum delor"

	optic := TrimmingWhileI(

		TraverseString(), OpI(func(index int, focus rune) bool {
			return index < 6 || index > 10
		}))

	listResult := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(string(listResult))

	overResult := MustModify(optic, Op(unicode.ToUpper), data)
	fmt.Println(overResult)

	//Output:
	//ipsum
	//Lorem IPSUM delor
}

func ExampleFilteredI() {

	data := "Lorem ipsum delor"

	optic := FilteredI(

		TraverseString(), OpI(func(index int, focus rune) bool {
			return index%2 == 0
		}))

	listResult := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(string(listResult))

	overResult := MustModify(optic, Op(unicode.ToUpper), data)
	fmt.Println(overResult)

	//Output:
	//Lrmismdlr
	//LoReM IpSuM DeLoR
}

func ExampleIndex() {

	data := map[int]string{
		1: "alpha",
		2: "beta",
		3: "gamma",
		4: "delta",
	}

	optic := Index(TraverseMap[int, string](), 3)

	viewResult, found := MustGetFirst(optic, data)
	fmt.Println(viewResult, found)

	var modifyResult map[int]string = MustModify(optic, Op(strings.ToUpper), data)
	fmt.Println(modifyResult)

	//Note: index cannot add new elements

	setResult := MustSet(Index(TraverseMap[int, string](), 5), "epsilon", data)
	fmt.Println(setResult)

	//Output:
	//gamma true
	//map[1:alpha 2:beta 3:GAMMA 4:delta]
	//map[1:alpha 2:beta 3:gamma 4:delta]

}

func ExampleWithIndex() {
	data := map[string]int{
		"alpha": 1,
		"beta":  2,
		"gamma": 3,
		"delta": 4,
	}

	optic := SliceOf(WithIndex(TraverseMap[string, int]()), len(data))

	var result []ValueI[string, int] = MustGet(optic, data)

	fmt.Println(result)
	//Output:
	//[alpha:1 beta:2 delta:4 gamma:3]
}

func ExampleIndices() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	optic := Indices(

		TraverseSlice[string](), Op(func(index int) bool {
			return index%2 == 0
		}))

	viewResult := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(viewResult)

	var modifyResult []string = MustModify(optic, Op(strings.ToUpper), data)
	fmt.Println(modifyResult)

	//Output:
	//[alpha gamma]
	//[ALPHA beta GAMMA delta]

}

func ExampleIndexing() {

	data := map[string]int{
		"alpha": 1,
		"beta":  2,
		"gamma": 3,
		"delta": 4,
	}

	optic := Indexing(TraverseMap[string, int]())

	viewResult := MustGet(SliceOf(WithIndex(optic), 4), data)
	fmt.Println(viewResult)

	var modifyResult map[string]int = MustModifyI(optic, OpI(func(index int, focus int) int {
		if index%2 == 0 {
			return focus * 10
		}
		return focus
	}), data)
	fmt.Println(modifyResult)

	//Output:
	//[0:1 1:2 2:4 3:3]
	//map[alpha:10 beta:2 delta:40 gamma:3]
}

func ExampleReIndexed() {

	data := map[string]int{
		"alpha": 1,
		"beta":  2,
		"gamma": 3,
		"delta": 4,
	}

	optic := ReIndexed(
		TraverseMap[string, int](),
		Op(func(index string) int {
			return len(index)
		}),
		EqT2[int](),
	)

	//Focus on elements with an index > 4 where the index is now the string length of the key
	viewResult := MustGet(SliceOf(WithIndex(Indices(optic, Gt(4))), 4), data)
	fmt.Println(viewResult)

	//Focus on elements with an index > 4 where the index is now the string length of the key
	var modifyResult map[string]int = MustModify(Indices(optic, Gt(4)), Mul(10), data)
	fmt.Println(modifyResult)

	//Output:
	//[5:1 5:4 5:3]
	//map[alpha:10 beta:2 delta:40 gamma:30]
}

func ExampleSelfIndex() {

	data := map[string]lo.Tuple2[int, int]{
		"alpha": lo.T2(1, 2),
		"beta":  lo.T2(3, 4),
		"gamma": lo.T2(5, 6),
		"delta": lo.T2(7, 8),
	}

	//This optic focuses on the second element of the tuple but makes the whole tuple available in the index.
	optic := ComposeLeft(
		SelfIndex(
			TraverseMap[string, lo.Tuple2[int, int]](),
			EqDeepT2[lo.Tuple2[int, int]](),
		), T2B[int, int]())

	viewResult := MustGet(SliceOf(WithIndex(optic), 4), data)
	fmt.Println(viewResult)

	var modifyResult map[string]lo.Tuple2[int, int] = MustModifyI(optic, OpI(func(index lo.Tuple2[int, int], focus int) int {
		return index.A + index.B
	}), data)
	fmt.Println(modifyResult)

	//Output:
	//[{1 2}:2 {3 4}:4 {7 8}:8 {5 6}:6]
	//map[alpha:{1 3} beta:{3 7} delta:{7 15} gamma:{5 11}]
}

func ExampleOrderedI() {

	data := map[string]int{
		"z": 1,
		"y": 2,
		"x": 3,
		"w": 4,
	}

	optic := Taking(
		OrderedI(
			TraverseMap[string, int](),
			OrderBy(ValueIIndex[string, int]()),
		),
		2,
	)

	var getResult map[string]int
	getResult, err := Get(MapOf(optic, 5), data)
	fmt.Println(getResult, err)

	var modifyResult map[string]int = MustModify(optic, Mul(100), data)
	fmt.Println(modifyResult)

	//Output:
	//map[w:4 x:3] <nil>
	//map[w:400 x:300 y:2 z:1]
}

func ExampleFirstOrDefaultI() {

	optic := TraverseSlice[string]()

	index, result := MustGetI(FirstOrDefaultI(optic, -1, "default"), []string{})
	fmt.Println(index, result)

	//Output:
	//-1 default

}

func TestCompositeIOrderBy(t *testing.T) {

	data := []lo.Tuple2[int, int]{
		lo.T2(2, 2),
		lo.T2(3, 4),
		lo.T2(1, 2),
		lo.T2(3, 3),
	}

	res, err := Get(
		SliceOf(
			OrderedI(
				TraverseSlice[lo.Tuple2[int, int]](),
				OrderBy2(
					OrderByI(OpToOpI[int](T2A[int, int]())),
					OrderByI(OpToOpI[int](T2B[int, int]())),
				),
			),
			len(data),
		),
		data,
	)

	if err != nil || !reflect.DeepEqual(res, []lo.Tuple2[int, int]{
		lo.T2(1, 2),
		lo.T2(2, 2),
		lo.T2(3, 3),
		lo.T2(3, 4),
	}) {
		t.Fatal(res, err)
	}

}

func TestCompositeOrderBy(t *testing.T) {

	data := []lo.Tuple2[int, int]{
		lo.T2(2, 2),
		lo.T2(3, 4),
		lo.T2(1, 2),
		lo.T2(3, 3),
	}

	orderBy := Ordered(
		TraverseSlice[lo.Tuple2[int, int]](),
		OrderBy2(
			OrderBy(T2A[int, int]()),
			OrderBy(T2B[int, int]()),
		),
	)

	res, err := Get(
		SliceOf(
			orderBy,
			len(data),
		),
		data,
	)

	if err != nil || !reflect.DeepEqual(res, []lo.Tuple2[int, int]{
		lo.T2(1, 2),
		lo.T2(2, 2),
		lo.T2(3, 3),
		lo.T2(3, 4),
	}) {
		t.Fatal(res, err)
	}

	if explain := orderBy.AsExpr().Short(); explain != `Ordered(Traverse,OrderByN(OrderBy(TupleElement(0)) , OrderBy(TupleElement(1))))` {
		t.Fatal(explain)
	}

}

func ExampleForEachI() {

	lookupData := map[string]string{
		"alpha": "1",
		"beta":  "2",
		"gamma": "3",
	}

	lookupNames := ComposeLeft(Lookup(TraverseMap[string, string](), lookupData), Some[string]())

	optic := ForEachI(
		IxMapLeft[string, Void](lookupNames.AsIxMatch()),
		EErr(lookupNames),
		ParseInt[int](10, 32),
	)

	var res Collection[string, string, Err]
	res, err := Modify(optic, Mul(10), "alpha")

	fmt.Println(res, err)
	//Output:
	//Col[alpha:10] <nil>
}

func ExampleEditDistanceI() {

	data := []lo.Tuple2[string, string]{
		lo.T2("alpha", "alpha"),
		lo.T2("alpha", "Alpha"),
		lo.T2("alpha", "lapha"),
	}

	optic := Compose(
		TraverseSlice[lo.Tuple2[string, string]](),
		EditDistanceI(TraverseString(), EditOSA, EqT2[int](), EqT2[rune](), 10),
	)

	res := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(res)

	//Output:
	//[0 1 2]

}
