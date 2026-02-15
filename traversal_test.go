package optic_test

import (
	"context"
	"fmt"
	"iter"
	"reflect"
	"strconv"
	"testing"

	. "github.com/spearson78/go-optic"
)

func TestSliceFieldLens(t *testing.T) {

	p := Person{
		Name:    "steve",
		Age:     45,
		Hobbies: []string{"A", "B", "C"},
	}
	hobbiesLens := FieldLens[Person, []string](func(s *Person) *[]string { return &s.Hobbies })
	traverseHobbies := TraverseSlice[string]()

	composed := Compose(hobbiesLens, traverseHobbies)

	if res, err := Modify(composed, Op(func(s string) string {
		return s + "!"
	}), p); err != nil || !reflect.DeepEqual(res.Hobbies, []string{"A!", "B!", "C!"}) {
		t.Fatalf("TraverseHobbies : %v", res)
	}

	if res, err := Modify(composed, OpE(func(ctx context.Context, s string) (string, error) {
		return "", fmt.Errorf("BAD %v", s)
	}), p); err == nil || len(res.Hobbies) > 0 {
		t.Fatalf("TraverseHobbies Compose(AsResult): %v", res)
	}
}

func TestTraversed(t *testing.T) {

	//[1, 2, 3] & traversed *~ 10
	// MulTilde(traversed(),10)([1,2,3])

	//A any, B any, FB ApplicativeType[B], TA TraversableType[A], TB TraversableType[B], FTB ApplicativeType[TB], FA ApplicativeType[A], TFB TraversableType[FB], RET, FUNC ApplicativeType[func(TB) TB]
	traversed := TraverseSlice[int]()

	if r, err := Modify(traversed, Op(func(i int) int {
		return i * 10
	}), []int{1, 2, 3}); err != nil || !reflect.DeepEqual(r, []int{10, 20, 30}) {
		t.Fatalf("Traversed1 : %v", r)
	}

	if r, err := Modify(traversed, Mul(10), []int{1, 2, 3}); err != nil || !reflect.DeepEqual(r, []int{10, 20, 30}) {
		t.Fatalf("Traversed2a : %v", r)
	}

	if r, err := Get(SliceOf(Compose(traversed, Mul(10)), 3), []int{1, 2, 3}); err != nil || !reflect.DeepEqual(r, []int{10, 20, 30}) {
		t.Fatalf("Traversed2b : %v", r)
	}

}

func ExampleTraversal() {
	//This traversal iterates of the elements of an []int
	traverseSlice := Traversal[[]int, int](
		//iter
		func(source []int) iter.Seq[int] {
			return func(yield func(focus int) bool) {
				for _, focus := range source {
					//If yield returns false iteration must be stopped
					if !yield(focus) {
						break
					}
				}
			}
		},
		//lengthGetter
		func(source []int) int {
			return len(source)
		},
		//modify
		func(fmap func(focus int) int, source []int) []int {
			//modify functions must not alter the source they must return a modified copy
			var ret []int
			for _, focus := range source {
				ret = append(ret, fmap(focus))
			}
			return ret
		},
		ExprCustom("ExampleTraversal"),
	)

	result := MustModify(traverseSlice, Mul(10), []int{1, 2, 3})

	fmt.Println(result)
	//Output:[10 20 30]
}

func ExampleTraversalI() {
	//This traversal iterates of the elements of an []int including the index of the element
	traverseSlice := TraversalI[int, []int, int](
		//iter
		func(source []int) SeqI[int, int] {
			return func(yield func(index int, focus int) bool) {
				for index, focus := range source {
					//If yield returns false iteration must be stopped
					if !yield(index, focus) {
						break
					}
				}
			}
		},
		//length getter
		func(source []int) int {
			return len(source)
		},
		//modify
		func(fmap func(index int, focus int) int, source []int) []int {
			//modify functions must not alter the source they must return a modified copy
			var ret []int
			for index, focus := range source {
				ret = append(ret, fmap(index, focus))
			}
			return ret
		},
		func(source []int, index int) iter.Seq2[int, int] {
			return func(yield func(index int, focus int) bool) {
				yield(index, source[index])
			}
		},
		func(indexA, indexB int) bool {
			return indexA == indexB
		},
		ExprCustom("ExampleTraversalI"),
	)

	result := MustModifyI(traverseSlice, OpI(func(index int, focus int) int {
		return index * focus
	}), []int{1, 2, 3})

	fmt.Println(result)
	//Output:[0 2 6]
}

func ExampleTraversalE() {
	//This traversal converts the elements of a string slice to ints, reporting any conversion errors encountered.
	//It is non polymorphic so any modification actions still return a []string
	parseInts := TraversalE[[]string, int](
		//iter
		func(ctx context.Context, source []string) iter.Seq2[int, error] {
			return func(yield func(focus int, err error) bool) {
				for _, strFocus := range source {
					intFocus, err := strconv.ParseInt(strFocus, 10, 32)
					//errors are reported to the caller by yielding. The caller can decide to stop iteration by returning false from the yield function.
					if !yield(int(intFocus), err) {
						//If yield returns false iteration must be stopped
						break
					}
				}
			}
		},
		//Length getter, nil will use the iterator to calculate the length
		nil,
		//modify
		func(ctx context.Context, fmap func(focus int) (int, error), source []string) ([]string, error) {
			//Modify functions must not alter the source they must return a modified copy
			//This traversal is non polymorphic so even though we convert the focus to ints we have to convert back to strings for the return
			var ret []string
			for _, strFocus := range source {
				intFocus, err := strconv.ParseInt(strFocus, 10, 32)
				if err != nil {
					//In contrast to iter modify function must fail the modify operation immediately otherwise the returned []string would be incomplete
					return nil, err
				}
				newVal, err := fmap(int(intFocus))
				if err != nil {
					//In contrast to iter modify function must fail the modify operation immediately otherwise the returned []string would be incomplete
					return nil, err
				}

				//This traversal is non polymorphic so we must convert back string
				ret = append(ret, strconv.Itoa(newVal))
			}

			return ret, nil
		},
		ExprCustom("ExampleTraversalR"),
	)

	//Note that result is an []string but the Mul operation acts on integer values
	var result []string
	result, err := Modify(parseInts, Mul[int](10), []string{"1", "2", "3"})

	fmt.Println(result, err)

	_, err = Modify(parseInts, Mul(10), []string{"1", "two", "3"})

	fmt.Println(err.Error())

	//Output:[10 20 30] <nil>
	//strconv.ParseInt: parsing "two": invalid syntax
	//optic error path:
	//	Custom(ExampleTraversalR)
}

func ExampleTraversalP() {
	//This traversal is polymorphic. It converts a slice of ints to float64s.
	intToFloat := TraversalP[[]int, []float64, int, float64](
		//iter
		func(source []int) iter.Seq[int] {
			return func(yield func(focus int) bool) {
				//If yield returns false iteration must be stopped
				for _, focus := range source {
					if !yield(focus) {
						break
					}
				}
			}
		},
		//length getter
		func(source []int) int {
			return len(source)
		},
		//modify
		func(fmap func(focus int) float64, source []int) []float64 {
			//Modify functions must not alter the source they must return a modified copy
			var ret []float64
			for _, focus := range source {
				ret = append(ret, fmap(focus))
			}
			return ret
		},
		ExprCustom("ExampleTraversalP"),
	)

	result := MustModify(intToFloat, Op(func(focus int) float64 {
		return float64(focus) + 0.5
	}), []int{1, 2, 3})

	fmt.Println(result)
	//Output:[1.5 2.5 3.5]
}

func ExampleTraversalIEP() {
	//This traversal converts the elements of a string slice to ints, reporting any conversion errors encountered.
	//It is polymorphic so is able to return an []int even through the source type is []string
	parseInts := TraversalIEP[int, []string, []int, int, int](
		func(ctx context.Context, source []string) SeqIE[int, int] {
			return func(yield func(ValueIE[int, int]) bool) {
				for index, strFocus := range source {
					intFocus, err := strconv.ParseInt(strFocus, 10, 32)
					//errors are reported to the caller by yielding. The caller can decide to stop iteration by returning false from the yield function.
					if !yield(ValIE(index, int(intFocus), err)) {
						//If yield returns false iteration must be stopped
						break
					}
				}
			}
		},
		//Length getter nil will use the iteratoe function to calculate the length,
		nil,
		func(ctx context.Context, fmap func(index int, focus int) (int, error), source []string) ([]int, error) {
			var ret []int
			for index, strFocus := range source {
				intFocus, err := strconv.ParseInt(strFocus, 10, 32)
				if err != nil {
					//We must fail the modify operation immediately otherwise the returned []int would be incomplete
					return nil, err
				}
				newVal, err := fmap(index, int(intFocus))
				if err != nil {
					//We must fail the modify operation immediately otherwise the returned []int would be incomplete
					return nil, err
				}

				ret = append(ret, newVal)
			}

			return ret, nil
		},
		func(ctx context.Context, index int, source []string) SeqIE[int, int] {
			return func(yield func(ValueIE[int, int]) bool) {
				strFocus := source[index]
				intFocus, err := strconv.ParseInt(strFocus, 10, 32)
				//errors are reported to the caller by yielding. The caller can decide to stop iteration by returning false from the yield function.
				yield(ValIE(index, int(intFocus), err))
			}
		},
		IxMatchComparable[int](),
		ExprCustom("ExampleTraversalF"),
	)

	var result []int
	result, err := Modify(parseInts, Mul(10), []string{"1", "2", "3"})

	fmt.Println(result, err)

	_, err = Modify(parseInts, Mul(10), []string{"1", "two", "3"})

	fmt.Println(err.Error())

	//Output:[10 20 30] <nil>
	//strconv.ParseInt: parsing "two": invalid syntax
	//optic error path:
	//	Custom(ExampleTraversalF)
}
