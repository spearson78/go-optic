package optic_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"

	. "github.com/spearson78/go-optic"
)

func TestCastConsistency(t *testing.T) {
	var data io.Writer = &bytes.Buffer{}
	var newVal *bytes.Buffer = &bytes.Buffer{}
	ValidateOpticTestPred(t, DownCast[io.Writer, *bytes.Buffer](), data, newVal, EqDeepT2[io.Writer]())
}

func ExampleDownCast() {

	data := []any{1, 2, 3}

	intCast := Compose(TraverseSlice[any](), DownCast[any, int]())

	var res []int = MustGet(SliceOf(intCast, 3), data)
	fmt.Println(res)

	badData := []any{1, "two", 3}

	res = MustGet(SliceOf(intCast, 3), badData)
	fmt.Println(res)

	var modifyResult []any = MustModify(intCast, Mul(2), data)
	fmt.Println(modifyResult)

	modifyResult = MustModify(intCast, Mul(2), badData)
	fmt.Println(modifyResult)

	//Output:
	//[1 2 3]
	//[1 3]
	//[2 4 6]
	//[2 two 6]
}

func ExampleDownCastP() {

	data := []any{1, 2, 3}

	intCast := Compose(TraverseSlice[any](), DownCastP[any, int, string]())

	res, err := Get(SliceOfP(intCast, 3), data)
	fmt.Println(res, err)

	badData := []any{1, true, 3}

	res, err = Get(SliceOfP(intCast, 3), badData)
	fmt.Println(res, err)

	var modifyResult []any
	modifyResult = MustModify(intCast, Compose3(Mul(2), FormatInt[int](10), PrependString(StringCol("Number:"))), data)
	fmt.Println(modifyResult)

	modifyResult = MustModify(intCast, Compose3(Mul(2), FormatInt[int](10), PrependString(StringCol("Number:"))), badData)
	fmt.Println(modifyResult)

	//Output:
	//[1 2 3] <nil>
	//[1 3] <nil>
	//[Number:2 Number:4 Number:6]
	//[Number:2 true Number:6]
}

func ExampleIsoCastE() {

	data := []any{1, 2, 3}

	intCast := Compose(TraverseSlice[any](), IsoCastE[any, int]())

	res, err := Get(SliceOf(intCast, 3), data)
	fmt.Println(res, err)

	badData := []any{1, true, 3}

	res, err = Get(SliceOf(intCast, 3), badData)
	fmt.Println(res, err)

	var modifyResult []any
	modifyResult, err = Modify(intCast, Mul(2), data)
	fmt.Println(modifyResult, err)

	modifyResult, err = Modify(intCast, Mul(2), badData)
	fmt.Println(modifyResult, err)

	//Output:
	//[1 2 3] <nil>
	//[] cast failed interface {} -> int for bool
	//optic error path:
	//	Cast
	//	Traverse
	//	SliceOf(Traverse | Cast)
	//
	//[2 4 6] <nil>
	//[] cast failed interface {} -> int for bool
	//optic error path:
	//	Cast
	//	Traverse
}

func ExampleIsoCastEP() {

	data := []any{1, 2, 3}

	intCast := Compose(TraverseSliceP[any, []rune](), IsoCastEP[any, []rune, int, string]())

	res, err := Get(SliceOfP(intCast, 3), data)
	fmt.Println(res, err)

	badData := []any{1, true, 3}

	res, err = Get(SliceOfP(intCast, 3), badData)
	fmt.Println(res, err)

	var modifyResult [][]rune
	modifyResult, err = Modify(intCast, Op(strconv.Itoa), data)
	fmt.Println(modifyResult, err)

	modifyResult, err = Modify(intCast, Op(strconv.Itoa), badData)
	fmt.Println(modifyResult, err)

	//Output:
	//[1 2 3] <nil>
	//[] cast failed interface {} -> int for bool
	//optic error path:
	//	Cast
	//	Traverse
	//	SliceOf(Traverse | Cast)
	//
	//[[49] [50] [51]] <nil>
	//[] cast failed interface {} -> int for bool
	//optic error path:
	//	Cast
	//	Traverse
}

func TestUpCastPanic(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				var errCast *ErrCast
				if errors.As(err, &errCast) {
					return
				}
			}
			panic(r)
		}
	}()

	//WARNING: This is an unsafe UpCast and will cause a panic
	unsafeOptic := Compose(TraverseSlice[any](), UpCast[any, int]()) //Casting any to int may fail. In this case a panic will occur.
	res, err := Get(SliceOf(unsafeOptic, 3), []any{1, 2, 3})
	t.Fatal(fmt.Printf("up cast did not panic: %v : %v", res, err))

}

func ExampleUpCast() {

	data := []int{1, 2, 3}

	optic := Compose(TraverseSlice[int](), UpCast[int, any]()) //Casting int to any always succeeds this is a safe UpCast

	var res []any = MustGet(SliceOf(optic, 3), data)
	fmt.Println(res)

	//WARNING: This is an unsafe UpCast and will cause a panic
	//unsafeOptic := Compose(TraverseSlice[any](), UpCast[any, int]()) //Casting any to int may fail. In this case a panic will occur.
	//var unsafeRes []int = MustGet(SliceOf(unsafeOptic, 3), []any{1, 2, 3})
	//fmt.Println(unsafeRes)

	//Output:
	//[1 2 3]
}

func ExampleIsoCast() {

	data := []string{"alpha", "beta", "gamma"}

	optic := Compose(TraverseSlice[string](), IsoCast[string, []byte]()) //Casting string to []byte and back is always safe.

	var res [][]byte = MustGet(SliceOf(optic, 3), data)
	fmt.Println(res)

	var result []string = MustModify(
		Compose3(
			TraverseSlice[string](),
			IsoCast[string, []byte](),
			TraverseSlice[byte](),
		),
		Add[byte](1),
		data,
	)

	fmt.Println(result)

	//Output:
	//[[97 108 112 104 97] [98 101 116 97] [103 97 109 109 97]]
	//[bmqib cfub hbnnb]
}

func ExampleIsoCastP() {

	data := []byte{'a', 'l', 'p', 'h', 'a'}

	optic := IsoCastP[[]byte, []rune, string, string]() //Casting []byte to string and string to []rune is always safe.

	var res string = MustGet(optic, data)
	fmt.Println(res)

	var modifyRes []rune = MustModify(
		optic,
		Op(strings.ToUpper),
		data,
	)

	fmt.Println(string(modifyRes))

	//Output:
	//alpha
	//ALPHA
}
