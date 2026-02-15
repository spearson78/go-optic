package optic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/samber/mo"

	. "github.com/spearson78/go-optic"
)

func ExampleGetI() {

	data := map[int]string{
		10: "alpha",
		20: "beta",
		30: "gamma",
		40: "delta",
	}

	optic := AtMap[int, string](30)

	var index int
	var viewResult mo.Option[string]
	var err error
	index, viewResult, err = GetI(optic, data)
	fmt.Println(index, viewResult, err)

	//Output:
	//30 {true gamma} <nil>
}

func ExampleGetContextI() {

	data := map[int]string{
		10: "alpha",
		20: "beta",
		30: "gamma",
		40: "delta",
	}

	optic := AtMap[int, string](30)

	ctx, cancel := context.WithCancel(context.Background())

	var index int
	var viewResult mo.Option[string]
	var err error
	index, viewResult, err = GetContextI(ctx, optic, data)
	fmt.Println(index, viewResult, err)

	cancel()

	_, _, err = GetContextI(ctx, optic, data)
	fmt.Println(err)

	//Output:
	//30 {true gamma} <nil>
	//context canceled
	//optic error path:
	//	Const(30)
	//	IgnoreWrite(Const[map[int]string](30))
	//	TupleOf(Identity , IgnoreWrite(Const[map[int]string](30)))
}

func ExampleMustGetI() {

	data := map[int]string{
		10: "alpha",
		20: "beta",
		30: "gamma",
		40: "delta",
	}

	optic := AtMap[int, string](30)

	var index int
	var viewResult mo.Option[string]
	index, viewResult = MustGetI(optic, data)
	fmt.Println(index, viewResult)

	//Output:
	//30 {true gamma}
}

func ExampleModifyI() {

	data := []int{1, 2, 3, 4}

	optic := TraverseSlice[int]()

	var result []int
	var err error
	result, err = ModifyI(optic, OpI(func(index int, focus int) int {
		if index%2 == 0 {
			return focus * 10
		} else {
			return focus
		}
	}), data)
	fmt.Println(result, err)

	//Output:
	//[10 2 30 4] <nil>
}

func ExampleModifyContextI() {

	data := []int{1, 2, 3, 4}

	optic := TraverseSlice[int]()

	ctx, cancel := context.WithCancel(context.Background())

	var result []int
	var err error
	result, err = ModifyContextI(ctx, optic, OpI(func(index int, focus int) int {
		if index%2 == 0 {
			return focus * 10
		} else {
			return focus
		}
	}), data)
	fmt.Println(result, err)

	cancel()

	_, err = ModifyContextI(ctx, optic, OpI(func(index int, focus int) int {
		if index%2 == 0 {
			return focus * 10
		} else {
			return focus
		}
	}), data)
	fmt.Println(err)

	//Output:
	//[10 2 30 4] <nil>
	//context canceled
	//optic error path:
	//	func2
	//	Traverse
}

func ExampleMustModifyI() {

	data := []int{1, 2, 3, 4}

	optic := TraverseSlice[int]()

	var result []int = MustModifyI(optic, OpI(func(index int, focus int) int {
		if index%2 == 0 {
			return focus * 10
		} else {
			return focus
		}
	}), data)
	fmt.Println(result)

	//Output:
	//[10 2 30 4]
}

func ExampleModifyCheckI() {

	data := []int{1, 2, 3, 4}

	//This optic focuses each int in the slice
	optic := TraverseSlice[int]()

	var result []int
	var modified bool
	var err error

	result, modified, err = ModifyCheckI(optic, OpI(func(index int, focus int) int {
		if index%2 == 0 {
			return focus * 10
		} else {
			return focus
		}
	}), data)
	fmt.Println(result, modified, err)

	//When no elements are focused false is returned for modified
	emptyData := []int{}
	_, modified, _ = ModifyCheckI(optic, OpI(func(index int, focus int) int {
		if index%2 == 0 {
			return focus * 10
		} else {
			return focus
		}
	}), emptyData)
	fmt.Println(modified)

	//Output:
	//[10 2 30 4] true <nil>
	//false
}

func ExampleModifyCheckContextI() {

	data := []int{1, 2, 3, 4}

	//This optic focuses each int in the slice
	optic := TraverseSlice[int]()

	ctx, cancel := context.WithCancel(context.Background())

	var result []int
	var modified bool
	var err error

	result, modified, err = ModifyCheckContextI(ctx, optic, OpI(func(index int, focus int) int {
		if index%2 == 0 {
			return focus * 10
		} else {
			return focus
		}
	}), data)
	fmt.Println(result, modified, err)

	cancel()

	_, _, err = ModifyCheckContextI(ctx, optic, OpI(func(index int, focus int) int {
		if index%2 == 0 {
			return focus * 10
		} else {
			return focus
		}
	}), data)
	fmt.Println(err)

	//Output:
	//[10 2 30 4] true <nil>
	//context canceled
	//optic error path:
	//	func2
	//	Traverse
}

func ExampleMustModifyCheckI() {

	data := []int{1, 2, 3, 4}

	//This optic focuses each int in the slice
	optic := TraverseSlice[int]()

	var result []int
	var modified bool
	result, modified = MustModifyCheckI(optic, OpI(func(index int, focus int) int {
		if index%2 == 0 {
			return focus * 10
		} else {
			return focus
		}
	}), data)
	fmt.Println(result, modified)

	_, modified = MustModifyCheckI(optic, OpI(func(index int, focus int) int {
		if index%2 == 0 {
			return focus * 10
		} else {
			return focus
		}
	}), []int{})
	fmt.Println(modified)

	//Output:
	//[10 2 30 4] true
	//false
}

func ExampleGetFirstI() {

	data := map[string]int{
		"alpha": 1,
		"beta":  2,
		"gamma": 3,
		"delta": 4,
	}

	optic := TraverseMap[string, int]()

	index, viewResult, found, err := GetFirstI(optic, data)
	fmt.Println(index, viewResult, found, err)

	_, found, err = GetFirst(optic, map[string]int{})
	fmt.Println(found, err)

	//Output:
	//alpha 1 true <nil>
	//false <nil>
}

func ExampleGetFirstContextI() {

	data := map[string]int{
		"alpha": 1,
		"beta":  2,
		"gamma": 3,
		"delta": 4,
	}

	optic := TraverseMap[string, int]()

	ctx, cancel := context.WithCancel(context.Background())

	index, viewResult, found, err := GetFirstContextI(ctx, optic, data)
	fmt.Println(index, viewResult, found, err)

	cancel()

	_, _, err = GetFirstContext(ctx, optic, data)
	fmt.Println(err)

	//Output:
	//alpha 1 true <nil>
	//context canceled
	//optic error path:
	//	Traverse
}

func ExampleMustGetFirstI() {

	data := map[string]int{
		"alpha": 1,
		"beta":  2,
		"gamma": 3,
		"delta": 4,
	}

	optic := TraverseMap[string, int]()

	index, viewResult, found := MustGetFirstI(optic, data)
	fmt.Println(index, viewResult, found)

	_, _, found = MustGetFirstI(optic, map[string]int{})
	fmt.Println(found)

	//Output:
	//alpha 1 true
	//false
}

func TestFailModifyContextWithTraverseCol(t *testing.T) {

	data := ValCol[int](1, 2, 3)

	result, ok, err := ModifyCheckContext(
		context.Background(),
		TraverseCol[int, int](),
		Mul(2),
		data,
	)

	if !ok || err != nil {
		t.Fatal(result, ok, err)
	}

}
