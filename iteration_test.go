package optic_test

import (
	"context"
	"fmt"
	"iter"
	"log"
	"strconv"

	. "github.com/spearson78/go-optic"
)

func ExampleIteration() {
	//This iteration iterates of the elements of an []int
	iterateSlice := Iteration[[]int, int](
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
		ExprCustom("ExampleIteration"),
	)

	for val := range MustGet(SeqOf(Compose(iterateSlice, Mul(10))), []int{1, 2, 3}) {
		fmt.Println(val)
	}

	//Output:
	//10
	//20
	//30
}

func ExampleIterationI() {
	//This Iteration iterates of the elements of an []int including the index of the element
	iterateSlice := IterationI[int, []int, int](
		//iter
		func(source []int) iter.Seq2[int, int] {
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
		//ixget
		func(source []int, index int) iter.Seq2[int, int] {
			return func(yield func(index int, focus int) bool) {
				yield(index, source[index])
			}
		},
		//ix match
		func(indexA, indexB int) bool {
			return indexA == indexB
		},
		ExprCustom("ExampleIterationI"),
	)

	for i, val := range MustGet(SeqIOf(ComposeLeft(iterateSlice, Mul(10))), []int{1, 2, 3}) {
		fmt.Println(i, val)
	}

	//Output:
	//0 10
	//1 20
	//2 30
}

func ExampleIterationE() {
	//This Iteration converts the elements of a string slice to ints, reporting any conversion errors encountered.
	//It is non polymorphic so any modification actions still return a []string
	parseInts := IterationE[[]string, int](
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
		ExprCustom("ExampleIterationE"),
	)

	seq, err := Get(SeqEOf(Compose(parseInts, Mul(10))), []string{"1", "2", "3"})
	if err != nil {
		log.Fatal(err)
	}

	for val := range seq {
		res, err := val.Get()
		fmt.Println(res, err)
	}

	seq, err = Get(SeqEOf(Compose(parseInts, Mul(10))), []string{"1", "two", "3"})
	if err != nil {
		log.Fatal(err)
	}
	for val := range seq {
		res, err := val.Get()
		fmt.Println(res, err)
	}

	//10 <nil>
	//20 <nil>
	//30 <nil>
	//10 <nil>
	//0 strconv.ParseInt: parsing "two": invalid syntax
	//optic error path:
	//	Custom(ExampleIterationR)
	//30 <nil>
}

func ExampleIterationIE() {

	parseInts := IterationIE[int, []string, int](
		func(ctx context.Context, source []string) SeqIE[int, int] {
			return func(yield func(ValueIE[int, int]) bool) {
				for index, strFocus := range source {
					intFocus, err := strconv.ParseInt(strFocus, 10, 0)
					//errors are reported to the caller by yielding. The caller can decide to stop iteration by returning false from the yield function.
					if !yield(ValIE(index, int(intFocus), err)) {
						//If yield returns false iteration must be stopped
						break
					}
				}
			}
		},
		//Length getter nil will use the iterate function to calculate the length,
		nil,
		//ixget
		func(ctx context.Context, index int, source []string) SeqIE[int, int] {
			return func(yield func(ValueIE[int, int]) bool) {
				intFocus, err := strconv.ParseInt(source[index], 10, 0)
				//errors are reported to the caller by yielding
				yield(ValIE(index, int(intFocus), err))
			}
		},
		IxMatchComparable[int](),
		ExprCustom("ExampleIterationIE"),
	)

	seq, err := Get(SeqIEOf(Compose(parseInts, Mul(10))), []string{"1", "2", "3"})
	if err != nil {
		log.Fatal(err)
	}
	for val := range seq {
		index, res, err := val.Get()
		fmt.Println(index, res, err)
	}

	seq, err = Get(SeqIEOf(Compose(parseInts, Mul(10))), []string{"1", "two", "3"})
	if err != nil {
		log.Fatal(err)
	}
	for val := range seq {
		index, res, err := val.Get()
		fmt.Println(index, res, err)
	}

	//0 10 <nil>
	//1 20 <nil>
	//2 30 <nil>
	//3 10 <nil>
	//4 0 strconv.ParseInt: parsing "two": invalid syntax
	//optic error path:
	//	Custom(ExampleIterationF)
	//5 30 <nil>
}
