package optic_test

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/samber/lo"

	. "github.com/spearson78/go-optic"
)

func TestTraverseOf(t *testing.T) {
	runeBoth := TraverseT2[string]()

	x := WithComprehension(runeBoth)
	if res, err := Modify(x, Op(func(focus string) []string {
		return []string{
			strings.ToLower(focus),
			strings.ToUpper(focus),
		}
	}), lo.T2("a", "b")); err != nil || !reflect.DeepEqual(res, []lo.Tuple2[string, string]{
		lo.T2("a", "b"),
		lo.T2("a", "B"),
		lo.T2("A", "b"),
		lo.T2("A", "B"),
	}) {
		t.Fatal(res)
	}
}

func TestDeadline(t *testing.T) {

	infiniteTraversal := Traversal[Void, int](
		func(source Void) iter.Seq[int] {
			return func(yield func(focus int) bool) {
				i := 0
				for yield(i) {
					i++
				}
			}
		},
		nil,
		func(fmap func(focus int) int, source Void) Void {
			i := 0
			for {
				fmap(i)
				i++
			}
		},
		ExprCustom("TestDeadline"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := GetContext(ctx, SliceOf(infiniteTraversal, 0), Void{})

	if err == nil {
		t.Fatal("err == nil")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal("!errors.Is(err, context.DeadlineExceeded)")
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = ModifyContext(ctx, infiniteTraversal, Add(1), Void{})

	if err == nil {
		t.Fatal("err == nil")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal("!errors.Is(err, context.DeadlineExceeded)")
	}
}

func ExampleGet() {

	data := lo.T2("alpha", "2")

	//This optic focuses the A element of a tuple
	optic := T2A[string, string]()

	viewResult, err := Get(optic, data)
	fmt.Println(viewResult, err)

	//This optic focuses the A element of a tuple and converts it to an int
	composedOptic := Compose(T2A[string, string](), ParseInt[int](10, 32))

	_, err = Get(composedOptic, data)
	fmt.Println(err)

	//Output:
	//alpha <nil>
	//strconv.ParseInt: parsing "alpha": invalid syntax
	//optic error path:
	//	ParseInt(10,32)
	//	TupleElement(0)
}

func ExampleGetContext() {

	data := lo.T2("alpha", "2")

	//This optic focuses the A element of a tuple
	optic := T2A[string, string]()

	ctx, cancel := context.WithCancel(context.Background())

	viewResult, err := GetContext(ctx, optic, data)
	fmt.Println(viewResult, err)

	cancel()

	_, err = GetContext(ctx, optic, data)
	fmt.Println(err)

	//Output:
	//alpha <nil>
	//context canceled
	//optic error path:
	//	TupleElement(0)
}

func ExampleMustGet() {

	data := lo.T2("alpha", "2")

	//This optic focuses the A element of a tuple
	optic := T2A[string, string]()

	viewResult := MustGet(optic, data)
	fmt.Println(viewResult)

	//Output:
	//alpha
}

func ExampleReverseGet() {

	optic := ParseInt[int](10, 32)

	var err error
	var result string

	//ReView runs the optic "backwards" in this case converting an int back to a string.
	result, err = ReverseGet(optic, 1)
	fmt.Println(result, err)

	//Output:
	//1 <nil>
}

func ExampleReverseGetContext() {

	optic := ParseInt[int](10, 32)

	ctx, cancel := context.WithCancel(context.Background())

	var err error
	var result string

	//ReView runs the optic "backwards" in this case converting an int back to a string.
	result, err = ReverseGetContext(ctx, optic, 1)
	fmt.Println(result, err)

	cancel()

	_, err = ReverseGetContext(ctx, optic, 1)
	fmt.Println(err)

	//Output:
	//1 <nil>
	//context canceled
	//optic error path:
	//	ParseInt(10,32)
}

func ExampleMustReverseGet() {

	optic := Mul(10)

	//ReView runs the optic "backwards" in this case dividing by 10.
	var result int = MustReverseGet(optic, 10)
	fmt.Println(result)

	//Output:
	//1
}

func ExampleModify() {

	data := []string{"1", "2", "3", "4"}

	//This optic focuses each string in the slice parsed as an int
	optic := Compose(TraverseSlice[string](), ParseInt[int](10, 32))

	var result []string
	var err error
	//Note the operation works on an int but the result is a new []string with the applied modifications
	result, err = Modify(optic, Add[int](1), data)
	fmt.Println(result, err)

	//When an error is encountered the Over action short circuits and returns the error.
	errData := []string{"1", "two", "3", "4"}
	_, err = Modify(optic, Add[int](1), errData)
	fmt.Println(err)

	//Output:
	//[2 3 4 5] <nil>
	//strconv.ParseInt: parsing "two": invalid syntax
	//optic error path:
	//	ParseInt(10,32)
	//	Traverse
}

func ExampleModifyContext() {

	data := []string{"1", "2", "3", "4"}

	//This optic focuses each string in the slice parsed as an int
	optic := Compose(TraverseSlice[string](), ParseInt[int](10, 32))

	ctx, cancel := context.WithCancel(context.Background())

	var result []string
	var err error
	//Note the operation works on an int but the result is a new []string with the applied modifications
	result, err = ModifyContext(ctx, optic, Add[int](1), data)
	fmt.Println(result, err)

	cancel()

	_, err = ModifyContext(ctx, optic, Add[int](1), data)
	fmt.Println(err)

	//Output:
	//[2 3 4 5] <nil>
	//context canceled
	//optic error path:
	//	ParseInt(10,32)
	//	Traverse
}

func ExampleMustModify() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	//This optic focuses each string in the slice.
	optic := TraverseSlice[string]()

	//Note the operation works on an int but the result is a new []string with the applied modifications
	var result []string = MustModify(optic, Op(strings.ToUpper), data)
	fmt.Println(result)

	//Output:
	//[ALPHA BETA GAMMA DELTA]
}

func ExampleModifyCheck() {

	data := []int{1, 2, 3, 4}

	//This optic focuses each int in the slice
	optic := TraverseSlice[int]()

	var result []int
	var modified bool
	var err error

	result, modified, err = ModifyCheck(optic, Add[int](1), data)
	fmt.Println(result, modified, err)

	//When no elements are focused false is returned for modified
	emptyData := []int{}
	_, modified, _ = ModifyCheck(optic, Add[int](1), emptyData)
	fmt.Println(modified)

	//Output:
	//[2 3 4 5] true <nil>
	//false
}

func ExampleModifyCheckContext() {

	data := []int{1, 2, 3, 4}

	//This optic focuses each int in the slice
	optic := TraverseSlice[int]()

	ctx, cancel := context.WithCancel(context.Background())

	var result []int
	var modified bool
	var err error

	result, modified, err = ModifyCheckContext(ctx, optic, Add[int](1), data)
	fmt.Println(result, modified, err)

	cancel()

	_, _, err = ModifyCheckContext(ctx, optic, Add[int](1), data)
	fmt.Println(err)

	//Output:
	//[2 3 4 5] true <nil>
	//context canceled
	//optic error path:
	//	ValueI[int,int].value
	//	Traverse
}

func ExampleMustModifyCheck() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	//This optic focuses each string in the slice.
	optic := TraverseSlice[string]()

	var result []string
	var modified bool
	result, modified = MustModifyCheck(optic, Op(strings.ToUpper), data)
	fmt.Println(result, modified)

	_, modified = MustModifyCheck(optic, Op(strings.ToUpper), []string{})
	fmt.Println(modified)

	//Output:
	//[ALPHA BETA GAMMA DELTA] true
	//false
}

func ExampleGetFirst() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	optic := TraverseSlice[string]()

	viewResult, found, err := GetFirst(optic, data)
	fmt.Println(viewResult, found, err)

	_, found, err = GetFirst(optic, []string{})
	fmt.Println(found, err)

	//Output:
	//alpha true <nil>
	//false <nil>
}

func ExampleGetFirstContext() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	optic := TraverseSlice[string]()

	ctx, cancel := context.WithCancel(context.Background())

	viewResult, found, err := GetFirstContext(ctx, optic, data)
	fmt.Println(viewResult, found, err)

	cancel()

	_, _, err = GetFirstContext(ctx, optic, data)
	fmt.Println(err)

	//Output:
	//alpha true <nil>
	//context canceled
	//optic error path:
	//	Traverse
}

func ExampleMustGetFirst() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	optic := TraverseSlice[string]()

	viewResult, found := MustGetFirst(optic, data)
	fmt.Println(viewResult, found)

	_, found = MustGetFirst(optic, []string{})
	fmt.Println(found)

	//Output:
	//alpha true
	//false
}

func ExampleSet() {

	data := []string{"1", "2", "3", "4"}

	//This optic focuses each string in the slice parsed as an int
	optic := Compose(TraverseSlice[string](), ParseInt[int](10, 32))

	var result []string
	var err error
	//Note the operation works on an int but the result is a new []string with the applied modifications
	result, err = Set(optic, 100, data)
	fmt.Println(result, err)

	//Set does not need to parse the existing value so is able to update values that
	//would cause an error when using the over action.
	errData := []string{"1", "two", "3", "4"}
	result, err = Set(optic, 100, errData)
	fmt.Println(result, err)

	//Output:
	//[100 100 100 100] <nil>
	//[100 100 100 100] <nil>
}

func ExampleSetContext() {

	data := []string{"1", "2", "3", "4"}

	//This optic focuses each string in the slice parsed as an int
	optic := Compose(TraverseSlice[string](), ParseInt[int](10, 32))

	ctx, cancel := context.WithCancel(context.Background())

	var result []string
	var err error
	//Note the operation works on an int but the result is a new []string with the applied modifications
	result, err = SetContext(ctx, optic, 100, data)
	fmt.Println(result, err)

	cancel()

	_, err = SetContext(ctx, optic, 100, data)
	fmt.Println(err)

	//Output:
	//[100 100 100 100] <nil>
	//context canceled
	//optic error path:
	//	ParseInt(10,32)
	//	Traverse
}

func ExampleMustSet() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	//This optic focuses each string in the slice.
	optic := TraverseSlice[string]()

	var result []string = MustSet(optic, "epsilon", data)
	fmt.Println(result)

	//Output:
	//[epsilon epsilon epsilon epsilon]
}
