package optic_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	. "github.com/spearson78/go-optic"
)

// This example combinator returns an optic that focuses the first focused element.
func FirstExample[I, S, T, A, RET, RW, DIR, ERR any](o Optic[I, S, T, A, A, RET, RW, DIR, ERR]) Optic[I, S, T, A, A, RET, RW, UniDir, ERR] {
	return CombiTraversal[RET, RW, ERR, I, S, T, A, A](
		//iterate
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					focusIndex, v, focusErr := val.Get()
					yield(ValIE(focusIndex, v, focusErr))
					return false //End iteration after the first element
				})
			}
		},
		//length getter
		nil, //nil means auto generate a length getter based on iterate
		//modify,
		func(ctx context.Context, fmap func(index I, focus A) (A, error), source S) (T, error) {
			first := true

			return o.AsModify()(ctx, func(index I, focus A) (A, error) {
				if first {
					first = false
					return fmap(index, focus) //Map only the first element
				} else {
					return focus, nil //Otherwise pass through all elements unaltered
				}
			}, source)
		},
		//ixget
		nil, //nil means auto generate an ix get based on iterate
		//ixmatch
		IxMatchDeep[I](),
		ExprCustom("FirstExample"),
	)
}

func ExampleCombiTraversal() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	optic := FirstExample(TraverseSlice[string]())

	result := MustModify(optic, Op(strings.ToUpper), data)
	fmt.Println(result)

	//Output:
	//[ALPHA beta gamma delta]
}

func TestReconstrain5(t *testing.T) {

	type Err1 struct{}
	type Err2 struct{}
	type Err3 struct{}
	type Err4 struct{}
	type Err5 struct{}

	i := Identity[int]()

	c5 := Compose5(
		CombiEErr[Err1](i),
		CombiEErr[Err2](i),
		CombiEErr[Err3](i),
		CombiEErr[Err4](i),
		CombiEErr[Err5](i),
	)

	//This sequence of calls moves Err5 to be next to Err1.
	s0 := EErrSwapR(c5)
	s1 := EErrTrans(s0)
	s2 := EErrSwapL(s1)
	s3 := EErrTransL(s2)
	s4 := EErrSwapR(s3)
	s5 := EErrSwapL(s4)
	s6 := EErrTrans(s5)
	s7 := EErrSwap(s6)
	var s8 Optic[Void, int, int, int, int, CompositionTree[CompositionTree[CompositionTree[ReturnOne, ReturnOne], ReturnOne], CompositionTree[ReturnOne, ReturnOne]], CompositionTree[CompositionTree[CompositionTree[ReadWrite, ReadWrite], ReadWrite], CompositionTree[ReadWrite, ReadWrite]], CompositionTree[CompositionTree[CompositionTree[BiDir, BiDir], BiDir], CompositionTree[BiDir, BiDir]], CompositionTree[CompositionTree[Err1, Err5], CompositionTree[CompositionTree[Err3, Err4], Err2]]] = EErrSwapL(s7)

	t.Log(s8)

	//This sequence of calls moves Err4 to be next to Err1.
	x1 := EErrTrans(c5)
	x2 := EErrSwapL(x1)
	x3 := EErrTransL(x2)
	x4 := EErrSwapR(x3)
	x5 := EErrSwapL(x4)
	x6 := EErrTrans(x5)
	x7 := EErrSwap(x6)
	var x8 Optic[Void, int, int, int, int, CompositionTree[CompositionTree[CompositionTree[ReturnOne, ReturnOne], ReturnOne], CompositionTree[ReturnOne, ReturnOne]], CompositionTree[CompositionTree[CompositionTree[ReadWrite, ReadWrite], ReadWrite], CompositionTree[ReadWrite, ReadWrite]], CompositionTree[CompositionTree[CompositionTree[BiDir, BiDir], BiDir], CompositionTree[BiDir, BiDir]], CompositionTree[CompositionTree[Err1, Err4], CompositionTree[CompositionTree[Err3, Err5], Err2]]] = EErrSwapL(x7)

	t.Log(x8)

	//This sequence of calls swaps Err5 and Err1.
	swn1 := EErrTrans(c5)
	sw0 := EErrSwap(swn1)
	sw1 := EErrSwapL(sw0)
	sw2 := EErrTrans(sw1)
	sw3 := EErrTransL(sw2)
	sw4 := EErrSwapL(sw3)
	sw5 := EErrTrans(sw4)
	sw6 := EErrSwapL(sw5)
	sw7 := EErrTransL(sw6)
	sw8 := EErrSwapL(sw7)
	sw9 := EErrTrans(sw8)
	sw10 := EErrSwap(sw9)
	sw11 := EErrSwapL(sw10)
	var sw12 Optic[Void, int, int, int, int, CompositionTree[CompositionTree[CompositionTree[ReturnOne, ReturnOne], ReturnOne], CompositionTree[ReturnOne, ReturnOne]], CompositionTree[CompositionTree[CompositionTree[ReadWrite, ReadWrite], ReadWrite], CompositionTree[ReadWrite, ReadWrite]], CompositionTree[CompositionTree[CompositionTree[BiDir, BiDir], BiDir], CompositionTree[BiDir, BiDir]], CompositionTree[CompositionTree[CompositionTree[Err5, Err2], Err3], CompositionTree[Err4, Err1]]] = EErrSwap(sw11)

	t.Log(sw12)

}
