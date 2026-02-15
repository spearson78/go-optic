package optic_test

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"reflect"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/samber/mo"
	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

func ExampleColIE() {

	data := []string{"alpha", "beta", "gamma"}

	col := ColIE[Err, int, string](
		func(ctx context.Context) SeqIE[int, string] {
			return func(yield func(ValueIE[int, string]) bool) {
				for i, v := range data {
					if !yield(ValIE(i, v, nil)) {
						break
					}
				}
			}
		},
		func(ctx context.Context, index int) SeqIE[int, string] {
			return func(yield func(ValueIE[int, string]) bool) {
				if index < 0 || index >= len(data) {
					yield(ValIE(index, "", errors.New("index out of range")))
					return
				}

				yield(ValIE(index, data[index], nil))
			}
		},
		IxMatchComparable[int](),
		func(ctx context.Context) (int, error) {
			return len(data), nil
		},
	)

	optic := TraverseColE[int, string, Err]()

	res, ok, err := GetFirst(Index(optic, 1), col)
	fmt.Println("Result:", res, ok, err)

	res, ok, err = GetFirst(Index(optic, 10), col)
	fmt.Println("Result:", res, ok, err)

	//Output:
	//Result: beta true <nil>
	//Result:  false index out of range
	//optic error path:
	//	Traverse
	//	Index(10)
}

func ExampleColE() {

	data := []string{"alpha", "beta", "gamma"}

	col := ColE[string](
		func(ctx context.Context) SeqE[string] {
			return func(yield func(ValueE[string]) bool) {
				for _, v := range data {
					if !yield(ValE(v, nil)) {
						break
					}
				}
			}
		},
		func(ctx context.Context) (int, error) {
			return len(data), nil
		},
	)

	optic := TraverseColE[int, string, Err]()

	res, ok, err := GetFirst(Index(optic, 1), col)
	fmt.Println("Result:", res, ok, err)

	res, ok, err = GetFirst(Index(optic, 10), col)
	fmt.Println("Result:", res, ok, err)

	//Output:
	//Result: beta true <nil>
	//Result:  false <nil>
}

func ExampleColI() {

	data := []string{"alpha", "beta", "gamma"}

	col := ColI[int, string](
		func(yield func(int, string) bool) {
			for i, v := range data {
				if !yield(i, v) {
					break
				}
			}
		},
		func(index int) iter.Seq2[int, string] {
			return func(yield func(int, string) bool) {
				if index < 0 || index >= len(data) {
					return
				}

				yield(index, data[index])
			}

		},
		IxMatchComparable[int](),
		func() int {
			return len(data)
		},
	)

	optic := TraverseCol[int, string]()

	res, ok := MustGetFirst(Index(optic, 1), col)
	fmt.Println("Result:", res, ok)

	res, ok = MustGetFirst(Index(optic, 10), col)
	fmt.Println("Result:", res, ok)

	//Output:
	//Result: beta true
	//Result:  false
}

func ExampleCol() {

	data := []string{"alpha", "beta", "gamma"}

	col := Col[string](
		func(yield func(string) bool) {
			for _, v := range data {
				if !yield(v) {
					break
				}
			}
		},
		func() int {
			return len(data)
		},
	)

	optic := TraverseCol[int, string]()

	res, ok := MustGetFirst(Index(optic, 1), col)
	fmt.Println("Result:", res, ok)

	res, ok = MustGetFirst(Index(optic, 10), col)
	fmt.Println("Result:", res, ok)

	//Output:
	//Result: beta true
	//Result:  false
}

func ExampleColSourceFocusErr() {

	data := []string{"3", "beta", "2"}

	//The Sort may generate an error in the collection during sorting.
	var sort Optic[Void, Collection[int, string, Err], Collection[int, string, Err], Collection[int, string, Err], Collection[int, string, Err], ReturnOne, ReadWrite, UniDir, Err] = OrderedCol[int](
		OrderBy(
			ParseInt[int](10, 0),
		),
	)

	//FilteredCol however generates a pure collection.
	var filteredPure Optic[Void, Collection[int, string, Pure], Collection[int, string, Pure], Collection[int, string, Pure], Collection[int, string, Pure], ReturnOne, ReadWrite, UniDir, Pure] = FilteredCol[int](Ne("beta"))

	//ColSourceFocusPure converts the Pure Collections to the given ERR type.
	var filteredErr Optic[Void, Collection[int, string, Err], Collection[int, string, Err], Collection[int, string, Err], Collection[int, string, Err], ReturnOne, ReadWrite, UniDir, Err] = ColSourceFocusErr(filteredPure)

	//res, err := Modify(
	//	ColFocusErr(SliceToCol[string]()),
	//	Compose(
	//		filteredPure,
	//		sort, //<----- Compile error here, sort expects a Collection[int,string,Err] but filteredPure provides Collection[int,string,Pure]
	//	),
	//	data,
	//)

	res, err := Modify(
		ColFocusErr(SliceToCol[string]()),
		Compose(
			filteredErr,
			sort,
		),
		data,
	)

	fmt.Println(res, err)

	//Output:
	//[2 3] <nil>

}

func ExampleColSourcePure() {

	data := lo.T2("alpha", []string{"3", "2", "1"})

	//The Sort may generate an error in the collection during sorting.
	//The Collection has a CompositionTree asit's ERR.
	sort := OrderedCol[int](
		OrderBy(
			ParseInt[int](10, 0),
		),
	)

	sortSlice := Compose3(
		ColFocusErr(SliceToCol[string]()), //ColFocusErr pure conversion from slice to collection to focus on collection Err
		sort,
		ColSourceErr(AsReverseGet(SliceToCol[string]())), //ColSourcePure pure conversion from collection to slice to source from a collection Err
	)

	res, err := Modify(
		T2B[string, []string](),
		sortSlice,
		data,
	)

	fmt.Println(res, err)

	//Output:
	//{alpha [1 2 3]} <nil>
}

func ExampleColSourceErr() {

	data := lo.T2("alpha", []string{"3", "2", "1"})

	//The Sort may generate an error in the collection during sorting.
	//The Collectionon has a CompositionTree asit's ERR.
	sort := OrderedCol[int](
		OrderBy(
			ParseInt[int](10, 0),
		),
	)

	sortSlice := Compose3(
		ColFocusErr(SliceToCol[string]()), //ColFocusErr pure conversion from slice to collection to focus on collection Err
		sort,
		ColSourceErr(AsReverseGet(SliceToCol[string]())), //ColSourceErr pure conversion from collection to slice to source from a collection Err
	)

	res, err := Modify(
		T2B[string, []string](),
		sortSlice,
		data,
	)

	fmt.Println(res, err)

	//Output:
	//{alpha [1 2 3]} <nil>
}

func ExampleColFocusPure() {

	data := lo.T2("alpha", []string{"3", "2", "1"})

	//The Sort may generate an error in the collection during sorting.
	//The Collection has a CompositionTree asit's ERR.
	sort := OrderedCol[int](
		OrderBy(
			ParseInt[int](10, 0),
		),
	)

	sortSlice := Compose3(
		ColFocusErr(SliceToCol[string]()), //ColFocusPure pure conversion from slice to collection to focus on collection Err
		sort,
		ColSourceErr(AsReverseGet(SliceToCol[string]())), //ColSourceErr pure conversion from collection to slice to source from a collection Err
	)

	res, err := Modify(
		T2B[string, []string](),
		sortSlice,
		data,
	)

	fmt.Println(res, err)

	//Output:
	//{alpha [1 2 3]} <nil>
}

func ExampleColFocusErr() {

	data := lo.T2("alpha", []string{"3", "2", "1"})

	//The Sort may generate an error in the collection during sorting.
	//The Collection has a CompositionTree asit's ERR.
	sort := OrderedCol[int](
		OrderBy(
			ParseInt[int](10, 0),
		),
	)

	sortSlice := Compose3(
		ColFocusErr(SliceToCol[string]()), //ColFocusErr pure conversion from slice to collection to focus on collection Err
		sort,
		ColSourceErr(AsReverseGet(SliceToCol[string]())), //ColSourceErr pure conversion from collection to slice to source from a collection Err
	)

	res, err := Modify(
		T2B[string, []string](),
		sortSlice,
		data,
	)

	fmt.Println(res, err)

	//Output:
	//{alpha [1 2 3]} <nil>
}

func ExampleColErr() {
	data := ValCol("2", "3", "3")

	//The Sort may generate an error in the collection during sorting.
	//The Collection has a CompositionTree asit's ERR.
	sort := OrderedCol[int](
		OrderBy(
			Compose(
				ParseInt[int](10, 0),
				Negate[int](),
			),
		),
	)

	//Simplify the error composition tree in the collection.
	sortErr := ColSourceFocusErr(sort)

	//res, err := Modify(
	//	Identity[Collection[int, string, Err]](),
	//	sort, //<----- Compile error here, The error composition tree in sort is not compatible with filteredErr.
	//	ColErr(data),
	//)

	res, err := Modify(
		Identity[Collection[int, string, Err]](),
		sortErr,
		ColErr(data),
	)

	fmt.Println(res, err)

	//Output:
	//Col[1:3 2:3 0:2] <nil>

}

func TestTraverseColConsistency(t *testing.T) {
	data := ValCol(10, 20, 30, 40, 50)
	ValidateOpticTestPred(t, TraverseCol[int, int](), data, 60, EqColT2[int, int, Pure](EqT2[int]()))

	ValidateOpticTestPred(t, TraverseColE[int, int, Err](), ColErr(data), 60, EqColT2[int, int, Err](EqT2[int]()))
}

func TestReverse(t *testing.T) {

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	if r := MustGet(ReversedSlice[int](), data); !reflect.DeepEqual(r, []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}) {
		t.Fatalf("MustView: %v", r)
	}

	if r := MustReverseGet(ReversedSlice[int](), data); !reflect.DeepEqual(r, []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}) {
		t.Fatalf("MustReView: %v", r)
	}

	if r := MustModify(T2B[int, []int](), ReversedSlice[int](), lo.T2(1, data)); !reflect.DeepEqual(r, lo.T2(1, []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1})) {
		t.Fatalf("MustOver: %v", r)
	}
}

func benchmarkBackwards(max int, b *testing.B) {

	data := RangeCol(0, max)

	for n := 0; n < b.N; n++ {
		counter := 0
		Modify(Reversed(TraverseCol[int, int]()), Op(func(v int) int {
			counter++
			return counter
		}), data)
	}
}

func BenchmarkBackwards1(b *testing.B)    { benchmarkBackwards(1, b) }
func BenchmarkBackwards2(b *testing.B)    { benchmarkBackwards(2, b) }
func BenchmarkBackwards3(b *testing.B)    { benchmarkBackwards(3, b) }
func BenchmarkBackwards10(b *testing.B)   { benchmarkBackwards(10, b) }
func BenchmarkBackwards20(b *testing.B)   { benchmarkBackwards(20, b) }
func BenchmarkBackwards40(b *testing.B)   { benchmarkBackwards(40, b) }
func BenchmarkBackwards100(b *testing.B)  { benchmarkBackwards(100, b) }
func BenchmarkBackwards1000(b *testing.B) { benchmarkBackwards(1000, b) }

func TestReversed(t *testing.T) {

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var metrics Metrics

	if r := MustGet(SliceOf(WithMetrics(Reversed(TraverseSlice[int]()), &metrics), 10), data); !reflect.DeepEqual(r, []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}) || !reflect.DeepEqual(metrics, Metrics{Focused: 10, Access: 1}) {
		t.Fatalf("MustViewReversed: %v : %v", r, metrics)
	}

	metrics = Metrics{}
	counter := 0

	if r := MustModifyI(Reversed(WithMetrics(TraverseSlice[int](), &metrics)), OpI(func(i int, v int) int {
		counter++
		return counter + (i * 100) + (v * 1000)
	}), data); !reflect.DeepEqual(r, []int{1010, 2109, 3208, 4307, 5406, 6505, 7604, 8703, 9802, 10901}) || !reflect.DeepEqual(metrics, Metrics{Focused: 20, Access: 2}) {
		t.Fatalf("MustIOver: %v : %v", r, metrics)
	}

	metrics = Metrics{}
	counter = 0

	if r := MustModifyI(Reversed(Compose(T2B[int, []int](), WithMetrics(TraverseSlice[int](), &metrics))), OpI(func(i int, v int) int {
		counter++
		return counter + (i * 100) + (v * 1000)
	}), lo.T2(1, data)); !reflect.DeepEqual(r, lo.T2(1, []int{1010, 2109, 3208, 4307, 5406, 6505, 7604, 8703, 9802, 10901})) || !reflect.DeepEqual(metrics, Metrics{Focused: 20, Access: 2}) {
		t.Fatalf("MustOver: %v : %v", r, metrics)
	}

	metrics = Metrics{}
	counter = 0

}

func TestHeapSort(t *testing.T) {

	data := []int{278, 9, 3, 98, 3, 4, 690, 72, 8, 9761, 20, 93, 5, 9, 67, 28, 73, 479, 2, 5, 89}

	var metrics Metrics

	heapSort := WithMetrics(Ordered(TraverseSlice[int](), OrderBy(Identity[int]())), &metrics)

	if lowest, ok := MustGetFirst(heapSort, data); !ok || metrics.Custom["comparisons"] != 38 || lowest != 2 {
		t.Fatal("MustPreView", ok, metrics.Custom["comparisons"], lowest)
	}

	metrics = Metrics{}
	if sorted := MustGet(SliceOf(heapSort, len(data)), data); metrics.Custom["comparisons"] != 118 || !reflect.DeepEqual(sorted, []int{2, 3, 3, 4, 5, 5, 8, 9, 9, 20, 28, 67, 72, 73, 89, 93, 98, 278, 479, 690, 9761}) {
		t.Fatal("MustToSliceOf", metrics.Custom["comparisons"], sorted)
	}

	metrics = Metrics{}
	if sorted := MustGet(SliceOf(Taking(heapSort, 1), len(data)), data); metrics.Custom["comparisons"] != 38 || !reflect.DeepEqual(sorted, []int{2}) {
		t.Fatal("Taking(1)", metrics.Custom["comparisons"], sorted)
	}

	metrics = Metrics{}
	if sorted := MustGet(SliceOf(Taking(heapSort, 2), len(data)), data); metrics.Custom["comparisons"] != 44 || !reflect.DeepEqual(sorted, []int{2, 3}) {
		t.Fatal("Taking(2)", metrics.Custom["comparisons"], sorted)
	}

	metrics = Metrics{}
	if sorted := MustGet(SliceOf(Taking(heapSort, 3), len(data)), data); metrics.Custom["comparisons"] != 50 || !reflect.DeepEqual(sorted, []int{2, 3, 3}) {
		t.Fatal("Taking(3)", metrics.Custom["comparisons"], sorted)
	}

	metrics = Metrics{}
	if sorted := MustGet(SliceOf(Taking(heapSort, 10), len(data)), data); metrics.Custom["comparisons"] != 90 || !reflect.DeepEqual(sorted, []int{2, 3, 3, 4, 5, 5, 8, 9, 9, 20}) {
		t.Fatal("Taking(10)", metrics.Custom["comparisons"], sorted)
	}

	metrics = Metrics{}
	if sorted := MustGet(SliceOf(Taking(heapSort, 18), len(data)), data); metrics.Custom["comparisons"] != 117 || !reflect.DeepEqual(sorted, []int{2, 3, 3, 4, 5, 5, 8, 9, 9, 20, 28, 67, 72, 73, 89, 93, 98, 278}) {
		t.Fatal("Taking(18)", metrics.Custom["comparisons"], sorted)
	}

	metrics = Metrics{}
	if sorted := MustGet(SliceOf(Taking(heapSort, 19), len(data)), data); metrics.Custom["comparisons"] != 118 || !reflect.DeepEqual(sorted, []int{2, 3, 3, 4, 5, 5, 8, 9, 9, 20, 28, 67, 72, 73, 89, 93, 98, 278, 479}) {
		t.Fatal("Taking(19)", metrics.Custom["comparisons"], sorted)
	}

	metrics = Metrics{}
	if sorted := MustGet(SliceOf(Taking(heapSort, 20), len(data)), data); metrics.Custom["comparisons"] != 118 || !reflect.DeepEqual(sorted, []int{2, 3, 3, 4, 5, 5, 8, 9, 9, 20, 28, 67, 72, 73, 89, 93, 98, 278, 479, 690}) {
		t.Fatal("Taking(20)", metrics.Custom["comparisons"], sorted)
	}

	metrics = Metrics{}
	if sorted := MustGet(SliceOf(Taking(heapSort, 21), len(data)), data); metrics.Custom["comparisons"] != 118 || !reflect.DeepEqual(sorted, []int{2, 3, 3, 4, 5, 5, 8, 9, 9, 20, 28, 67, 72, 73, 89, 93, 98, 278, 479, 690, 9761}) {
		t.Fatal("Taking(21)", metrics.Custom["comparisons"], sorted)
	}

	metrics = Metrics{}
	if sorted := MustGet(SliceOf(Taking(heapSort, 22), len(data)), data); metrics.Custom["comparisons"] != 118 || !reflect.DeepEqual(sorted, []int{2, 3, 3, 4, 5, 5, 8, 9, 9, 20, 28, 67, 72, 73, 89, 93, 98, 278, 479, 690, 9761}) {
		t.Fatal("Taking(22)", metrics.Custom["comparisons"], sorted)
	}

}

func TestHeapSortOver(t *testing.T) {

	tupleData := lo.T2("hello", []lo.Tuple2[string, int]{
		lo.T2("a", 10),
		lo.T2("b", 1),
		lo.T2("c", 2),
		lo.T2("d", 9),
		lo.T2("e", 3),
		lo.T2("f", 8),
		lo.T2("g", 4),
		lo.T2("h", 7),
		lo.T2("i", 5),
		lo.T2("j", 6),
	})

	var metrics Metrics

	optic := EPure(WithMetrics(Ordered(
		TraverseColP[int, lo.Tuple2[string, int], lo.Tuple2[string, int]](),
		OrderBy(T2B[string, int]()),
	), &metrics))
	heapSort := ColOf(optic)

	//heapSort applied to a slice
	if sorted := MustModify(T2B[string, []lo.Tuple2[string, int]](), SliceOp(heapSort), tupleData); metrics.Custom["comparisons"] != 39 || !reflect.DeepEqual(sorted, lo.T2("hello", []lo.Tuple2[string, int]{lo.T2("b", 1), lo.T2("c", 2), lo.T2("e", 3), lo.T2("g", 4), lo.T2("i", 5), lo.T2("j", 6), lo.T2("h", 7), lo.T2("f", 8), lo.T2("d", 9), lo.T2("a", 10)})) {
		t.Fatal("SliceOp Taking(2)", metrics, sorted)
	}

	i := 0
	encodeSeq := Op(func(f int) int {
		i++
		return i*1000 + f
	})

	//heapSort applied to an simple Optic (taking(2) )
	metrics = Metrics{}
	i = 0
	heapSortTraverse := WithMetrics(Ordered(TraverseSlice[lo.Tuple2[string, int]](), OrderBy(T2B[string, int]())), &metrics)
	if sorted := MustModify(Compose3(T2B[string, []lo.Tuple2[string, int]](), heapSortTraverse, T2B[string, int]()), encodeSeq, tupleData); metrics.Custom["comparisons"] != 39 || !reflect.DeepEqual(sorted, lo.T2("hello", []lo.Tuple2[string, int]{lo.T2("a", 10010), lo.T2("b", 1001), lo.T2("c", 2002), lo.T2("d", 9009), lo.T2("e", 3003), lo.T2("f", 8008), lo.T2("g", 4004), lo.T2("h", 7007), lo.T2("i", 5005), lo.T2("j", 6006)})) {
		t.Fatal("HeapSortOver TraverseSlice()", metrics, sorted)
	}

	//heapSort applied to an arbitrary Optic (taking(2) )
	metrics = Metrics{}
	i = 0
	heapSortTaking := WithMetrics(Ordered(Taking(TraverseSlice[lo.Tuple2[string, int]](), 2), OrderBy(T2B[string, int]())), &metrics)
	if sorted := MustModify(Compose3(T2B[string, []lo.Tuple2[string, int]](), heapSortTaking, T2B[string, int]()), encodeSeq, tupleData); metrics.Custom["comparisons"] != 1 || !reflect.DeepEqual(sorted, lo.T2("hello", []lo.Tuple2[string, int]{lo.T2("a", 2010), lo.T2("b", 1001), lo.T2("c", 2), lo.T2("d", 9), lo.T2("e", 3), lo.T2("f", 8), lo.T2("g", 4), lo.T2("h", 7), lo.T2("i", 5), lo.T2("j", 6)})) {
		t.Fatal("HeapSortOver Taking(2)", metrics, sorted)
	}

	//Composed SliceOp focuses in sorted order but does not apply the sort. An intermediate slice is created.
	metrics = Metrics{}
	i = 0
	if sorted := MustModify(
		T2B[string, []lo.Tuple2[string, int]](),
		SliceOp(heapSort),
		tupleData); metrics.Custom["comparisons"] != 39 || !reflect.DeepEqual(sorted, lo.T2("hello", []lo.Tuple2[string, int]{lo.T2("b", 1), lo.T2("c", 2), lo.T2("e", 3), lo.T2("g", 4), lo.T2("i", 5), lo.T2("j", 6), lo.T2("h", 7), lo.T2("f", 8), lo.T2("d", 9), lo.T2("a", 10)})) {
		t.Fatal("SliceOp()|Traverse", metrics, sorted)
	}

}

func TestFilteredColModify(t *testing.T) {

	tupleData := lo.T2("hello", []lo.Tuple2[string, int]{
		lo.T2("a", 10),
		lo.T2("b", 1),
		lo.T2("c", 2),
		lo.T2("d", 9),
		lo.T2("e", 3),
		lo.T2("f", 8),
		lo.T2("g", 4),
		lo.T2("h", 7),
		lo.T2("i", 5),
		lo.T2("j", 6),
	})

	pred := Compose(
		T2B[string, int](),
		Gte(6),
	)

	filtered := FilteredCol[int](EPure(pred))

	//filtered applied to a slice
	if sorted, err := Modify(T2B[string, []lo.Tuple2[string, int]](), SliceOp(filtered), tupleData); err != nil || !reflect.DeepEqual(sorted, lo.T2("hello", []lo.Tuple2[string, int]{lo.T2("a", 10), lo.T2("d", 9), lo.T2("f", 8), lo.T2("h", 7), lo.T2("j", 6)})) {
		t.Fatal("SliceOp", sorted, err)
	}

	//filtered applied to an arbitrary Optic (taking(5) )
	if sorted, err := Modify(
		Compose3(
			T2B[string, []lo.Tuple2[string, int]](),
			Filtered(
				Taking(
					TraverseSlice[lo.Tuple2[string, int]](),
					5,
				),
				pred, //FilteredCol
			),
			T2B[string, int](),
		),
		Mul(2),
		tupleData,
	); err != nil || !reflect.DeepEqual(sorted, lo.T2("hello", []lo.Tuple2[string, int]{lo.T2("a", 20), lo.T2("b", 1), lo.T2("c", 2), lo.T2("d", 18), lo.T2("e", 3), lo.T2("f", 8), lo.T2("g", 4), lo.T2("h", 7), lo.T2("i", 5), lo.T2("j", 6)})) {
		t.Fatal("FilteredColModify Taking(5)", sorted, err)
	}

	//Composed SliceOp focuses a filtered but does not apply the filter. An intermediate slice is created.
	if sorted, err := Modify(
		Compose3(
			T2B[string, []lo.Tuple2[string, int]](),
			Filtered(
				TraverseSlice[lo.Tuple2[string, int]](),
				pred,
			),
			T2B[string, int](),
		),
		Mul(2),
		tupleData,
	); err != nil || !reflect.DeepEqual(sorted, lo.T2("hello", []lo.Tuple2[string, int]{lo.T2("a", 20), lo.T2("b", 1), lo.T2("c", 2), lo.T2("d", 18), lo.T2("e", 3), lo.T2("f", 16), lo.T2("g", 4), lo.T2("h", 14), lo.T2("i", 5), lo.T2("j", 12)})) {
		t.Fatal("SliceOp()|Traverse", sorted, err)
	}
}

func TestSliceOpBackwards(t *testing.T) {

	tupleData := lo.T2("hello", []lo.Tuple2[string, int]{
		lo.T2("a", 10),
		lo.T2("b", 1),
		lo.T2("c", 2),
		lo.T2("d", 9),
		lo.T2("e", 3),
		lo.T2("f", 8),
		lo.T2("g", 4),
		lo.T2("h", 7),
		lo.T2("i", 5),
		lo.T2("j", 6),
	})

	optic := Reversed(TraverseCol[int, lo.Tuple2[string, int]]())
	reversed := ColOf(optic)

	//heapSort applied to a slice
	if sorted := MustModify(T2B[string, []lo.Tuple2[string, int]](), SliceOp(reversed), tupleData); !reflect.DeepEqual(sorted, lo.T2("hello", []lo.Tuple2[string, int]{lo.T2("j", 6), lo.T2("i", 5), lo.T2("h", 7), lo.T2("g", 4), lo.T2("f", 8), lo.T2("e", 3), lo.T2("d", 9), lo.T2("c", 2), lo.T2("b", 1), lo.T2("a", 10)})) {
		t.Fatal("SliceOp", sorted)
	}

	i := 0
	encodeSeq := Op(func(f int) int {
		i++
		return i*1000 + f
	})

	//backwards applied to an simple Optic (traverseslice)
	i = 0
	if sorted := MustModify(Compose3(T2B[string, []lo.Tuple2[string, int]](), Reversed(TraverseSlice[lo.Tuple2[string, int]]()), T2B[string, int]()), encodeSeq, tupleData); !reflect.DeepEqual(sorted, lo.T2("hello", []lo.Tuple2[string, int]{lo.T2("a", 10010), lo.T2("b", 9001), lo.T2("c", 8002), lo.T2("d", 7009), lo.T2("e", 6003), lo.T2("f", 5008), lo.T2("g", 4004), lo.T2("h", 3007), lo.T2("i", 2005), lo.T2("j", 1006)})) {
		t.Fatal("SliceOp TraverseSlice()", sorted)
	}

	//backwards applied to an arbitrary Optic (taking(2) )
	i = 0
	if sorted := MustModify(Compose3(T2B[string, []lo.Tuple2[string, int]](), Reversed(Taking(TraverseSlice[lo.Tuple2[string, int]](), 2)), T2B[string, int]()), encodeSeq, tupleData); !reflect.DeepEqual(sorted, lo.T2("hello", []lo.Tuple2[string, int]{lo.T2("a", 2010), lo.T2("b", 1001), lo.T2("c", 2), lo.T2("d", 9), lo.T2("e", 3), lo.T2("f", 8), lo.T2("g", 4), lo.T2("h", 7), lo.T2("i", 5), lo.T2("j", 6)})) {
		t.Fatal("SliceOp Taking(2)", sorted)
	}
}

func ExampleValCol() {

	data := []string{
		"alpha",
		"beta",
	}

	var result []string = MustModify(
		SliceToCol[string](),
		AppendCol(ValCol("gamma", "delta")),
		data,
	)

	fmt.Println(result)

	//Output:
	//[alpha beta gamma delta]
}

func ExampleValColI() {

	data := map[string]int{
		"alpha": 1,
		"beta":  2,
	}

	var result map[string]int = MustModify(
		MapToCol[string, int](),
		AppendCol(
			ValColI(
				IxMatchComparable[string](),
				ValI("gamma", 3),
				ValI("delta", 4),
			),
		),
		data,
	)

	fmt.Println(result)

	//Output:
	//map[alpha:1 beta:2 delta:4 gamma:3]
}

func ExampleValColE() {

	data := []string{
		"alpha",
		"beta",
	}

	s := ValColE[string](
		ValE("gamma", nil),
		ValE("", errors.New("sabotage")),
	)

	_, err := Modify(
		ColFocusErr(SliceToCol[string]()),
		AppendCol(
			s,
		),
		data,
	)

	fmt.Println(err)

	//Output:
	//sabotage
	//optic error path:
	//	ToCol
	//	ColFocusErr(ToCol(int,[]string,[]string,string,string))
}

func ExampleValColIE() {

	data := map[string]int{
		"alpha": 1,
		"beta":  2,
	}

	_, err := Modify(
		ColFocusErr(MapToCol[string, int]()),
		AppendCol(
			ValColIE[string, int, Err](
				IxMatchComparable[string](),
				ValIE("gamma", 3, nil),
				ValIE("delta", 4, errors.New("sabotage")),
			),
		),
		data,
	)

	fmt.Println(err)

	//Output:
	//sabotage
	//optic error path:
	//	ToCol
	//	ColFocusErr(ToCol(string,map[string]int,map[string]int,int,int))
}

func ExampleRangeCol() {

	data := []int{1, 2, 3}

	var result []int = MustModify(
		SliceToCol[int](),
		AppendCol(
			RangeCol(4, 10),
		),
		data,
	)

	fmt.Println(result)

	//[1 2 3 4 5 6 7 8 9 10]
}

func ExampleTraverseCol() {

	col := ValCol(5, 10, 15)

	result, err := Modify(
		TraverseCol[int, int](),
		Mul(10),
		col,
	)
	fmt.Println(result, err)

	//Output:
	//Col[0:50 1:100 2:150] <nil>

}

func ExampleTraverseColI() {

	col := ValColI(
		func(a, b []int) bool { return a[0] == b[0] },
		ValI([]int{0}, 5),
		ValI([]int{1}, 10),
		ValI([]int{2}, 15),
	)

	result, err := Modify(
		TraverseColI[[]int, int](col.AsIxMatch()),
		Mul(10),
		col,
	)
	fmt.Println(result, err)

	//Output:
	//Col[[0]:50 [1]:100 [2]:150] <nil>
}

func ExampleTraverseColE() {

	col := ValColE(
		ValE(1, nil),
		ValE(0, errors.New("sabotage")),
		ValE(3, nil),
	)

	result, err := Modify(
		TraverseColE[int, int, Err](),
		Mul(10),
		col,
	)
	fmt.Println(result, err)

	//Output:
	//<nil> sabotage
	//optic error path:
	//	Traverse
}

func ExampleTraverseColP() {

	col := ValCol("1", "2", "3")

	var overResult Collection[int, int, Pure]
	overResult, err := Modify(
		TraverseColP[int, string, int](),
		ParseInt[int](10, 32),
		col,
	)
	fmt.Println(overResult, err)

	//Output:
	//Col[0:1 1:2 2:3] <nil>
}

func ExampleTraverseColEP() {

	col := ValColE(
		ValE("1", nil),
		ValE("", errors.New("sabotage")),
		ValE("3", nil),
	)

	result, err := Modify(
		Compose(
			TraverseColEP[int, string, int, Err](),
			ParseIntP[int](10, 0),
		),
		Mul(10),
		col,
	)
	fmt.Println(result, err)

	//Output:
	//<nil> sabotage
	//optic error path:
	//	Traverse
}

func ExampleTraverseColIEP() {

	col := ValColE(
		ValE("1", nil),
		ValE("", errors.New("sabotage")),
		ValE("3", nil),
	)

	result, err := Modify(
		Compose(
			TraverseColIEP[int, string, int, Err](IxMatchComparable[int]()),
			ParseIntP[int](10, 0),
		),
		Mul(10),
		col,
	)
	fmt.Println(result, err)

	//Output:
	//<nil> sabotage
	//optic error path:
	//	Traverse
}

func TestDiffCol(t *testing.T) {

	before := ValCol(1, 2, 3, 4)
	after := ValCol(3, 2, 1, 8)

	diff := DiffColT2[int, int](0.5, Distance(func(a, b int) float64 {
		return float64(a - b)
	}), EqT2[int](), DiffNone, true)

	if res := MustGet(SliceOf(WithIndex(diff), 3), lo.T2(after, before)); !MustGet(EqDeepT2[[]ValueI[Diff[int, int], mo.Option[int]]](), lo.T2(res, []ValueI[Diff[int, int], mo.Option[int]]{
		ValI(
			Diff[int, int]{
				Type:        DiffModify,
				BeforeIndex: 2,
				AfterIndex:  0,
				BeforePos:   2,
				AfterPos:    0,
				BeforeValue: 3,
				Distance:    0,
			},
			mo.Some(3),
		),
		ValI(
			Diff[int, int]{
				Type:        DiffModify,
				BeforeIndex: 0,
				AfterIndex:  2,
				BeforePos:   0,
				AfterPos:    2,
				BeforeValue: 1,
				Distance:    0,
			},
			mo.Some(1),
		),
		ValI(
			Diff[int, int]{
				Type:        DiffRemove,
				BeforeIndex: 3,
				BeforePos:   3,
				BeforeValue: 4,
			},
			mo.None[int](),
		),
		ValI(
			Diff[int, int]{
				Type:       DiffAdd,
				AfterIndex: 3,
				AfterPos:   3,
			},
			mo.Some(8),
		),
	})) {
		t.Fatal(res)
	}

}

func TestDiffColDup(t *testing.T) {

	before := ValCol("alpha", "beta")
	after := ValCol("alpha", "alpha", "alpha")

	diff := DiffColT2[int](0.5, DistancePercent(EditDistance(TraverseString(), EditOSA, EqT2[rune](), 10), Length(TraverseString())), EqT2[int](), DiffNone, true)

	if res := MustGet(SliceOf(WithIndex(diff), 3), lo.T2(after, before)); !MustGet(EqDeepT2[[]ValueI[Diff[int, string], mo.Option[string]]](), lo.T2(res, []ValueI[Diff[int, string], mo.Option[string]]{
		ValI(
			Diff[int, string]{
				Type:        DiffRemove,
				BeforeIndex: 1,
				BeforePos:   1,
				BeforeValue: "beta",
			},
			mo.None[string](),
		),
		ValI(
			Diff[int, string]{
				Type:       DiffAdd,
				AfterIndex: 1,
				AfterPos:   1,
			},
			mo.Some("alpha"),
		),
		ValI(
			Diff[int, string]{
				Type:       DiffAdd,
				AfterIndex: 2,
				AfterPos:   2,
			},
			mo.Some("alpha"),
		),
	})) {
		t.Fatal(res)
	}
}

func ExampleColOf() {

	data := map[int][]string{
		1: []string{"gamma", "beta"},
		2: []string{"alpha", "delta"},
	}

	optic := ColOf(
		EPure(Compose(
			TraverseMap[int, []string](),
			TraverseSlice[string](),
		)),
	)

	var viewResult Collection[int, string, Pure] = MustGet(optic, data)
	fmt.Println(viewResult)

	op := OrderedCol[int](OrderBy(Identity[string]()))

	modifyResult := MustModify(optic, op, data)
	fmt.Println(modifyResult)

	//Output:
	//Col[0:gamma 1:beta 0:alpha 1:delta]
	//map[1:[alpha beta] 2:[delta gamma]]
}

func ExampleColOfP() {

	data := map[int]string{
		1: "alpha",
		2: "beta",
		3: "gamma",
		4: "delta",
	}

	//Polymorphic map traversal with a source type of map[int]string and result type of map[int]int
	optic := ColOfP(TraverseMapP[int, string, int]())

	var viewResult Collection[int, string, Err]
	viewResult, err := Get(optic, data)
	fmt.Println(viewResult, err)

	//Note the return type is map[int]int not map[int]string
	var modifyResult map[int]int
	modifyResult, err = Modify(optic, OpE(func(ctx context.Context, focus Collection[int, string, Err]) (Collection[int, int, Err], error) {
		var modified []ValueE[int]
		seq, err := Get(SeqEOf(TraverseColE[int, string, Err]()), focus)
		if err != nil {
			return nil, err
		}
		for v := range seq {
			str, err := v.Get()
			modified = append(modified, ValE(len(str), err))
		}

		return ValColE(modified...), nil
	}), data)
	fmt.Println(modifyResult, err)

	//Output:
	//Col[1:alpha 2:beta 3:gamma 4:delta] <nil>
	//map[1:5 2:4 3:5 4:5] <nil>
}

func TestReIndexedColNonIso(t *testing.T) {

	//ReIndexCol calls reverseget in the iso set method this means that a non iso ixmap can only be used as a getter.

	optic := ReIndexedCol[string](Lens[int, int](
		func(source int) int {
			return 1
		},
		func(focus, source int) int {
			return 10
		},
		ExprCustom("TestReIndexedCol"),
	))

	opticType := optic.OpticType()

	if opticType != expr.OpticTypeGetter {
		t.Fatal(opticType.String())
	}
}

func TestColOfFiltered(t *testing.T) {

	data := []string{
		"alpha",
		"beta",
		"gamma",
		"delta",
	}

	predicate := Ne("beta")

	filtered := ColOf(
		EPure(Filtered(
			TraverseCol[int, string](),
			predicate,
		)),
	)

	var res []string = MustModify(
		SliceToCol[string](),
		filtered,
		data,
	)
	fmt.Println(res)

	//Output:
	//[alpha gamma delta]

	if !reflect.DeepEqual(res, []string{"alpha", "gamma", "delta"}) {
		t.Fatal(res)
	}
}

func TestColOfModifyAbortHandler(t *testing.T) {

	//Tests whether the abortHandler works with lazy collections returned by ColOf

	result, err := Modify(
		ColOf(
			EErr(Compose(
				TraverseCol[int, string](),
				ComposeLeft(
					//This traversal uses the abortHandler
					Traversal[string, rune](
						func(source string) iter.Seq[rune] {
							return func(yield func(rune) bool) {
								for _, v := range source {
									if !yield(v) {
										break
									}
								}
							}
						},
						func(source string) int {
							return len([]rune(source))
						},
						func(fmap func(focus rune) rune, source string) string {
							var ret strings.Builder
							for _, v := range source {
								newV := fmap(v)
								ret.WriteRune(newV)
							}
							return ret.String()
						},
						ExprCustom("handleAbortTraversal"),
					),
					Error[rune, rune](errors.New("sabotage")),
				),
			)),
		),
		Op(func(c Collection[int, rune, Err]) Collection[int, rune, Err] {
			return ColErr(ValCol[rune]('a', 'l', 'p', 'h', 'a'))
		}),
		ValCol("beta"),
	)

	if result.String() != `Col[<sabotage
optic error path:
	Error(sabotage)
	Custom(handleAbortTraversal)
>]` {
		t.Fatal(result, err)
	}

}

func TestTraverseColModify(t *testing.T) {

	data := ValCol[int](1, 2, 3)

	result, ok, err := ModifyCheckContext(
		context.Background(),
		TraverseCol[int, int](),
		Add(1),
		data,
	)

	if err != nil || !ok || !MustGet(EqColT2[int, int, Pure](EqT2[int]()), lo.T2(ValCol[int](2, 3, 4), result)) {
		t.Fatal(result, ok, err)
	}
}
