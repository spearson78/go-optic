package optic_test

import (
	"context"
	"fmt"
	"strings"

	. "github.com/spearson78/go-optic"
)

func WithExampleLogging[I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, T, A, B, RET, RW, DIR, ERR] {

	return Omni[I, S, T, A, B, RET, RW, DIR, ERR](
		func(ctx context.Context, source S) (I, A, error) {
			i, a, err := o.AsGetter()(ctx, source)
			fmt.Println("Get", source, i, a, err)
			return i, a, err
		},
		func(ctx context.Context, focus B, source S) (T, error) {
			t, err := o.AsSetter()(ctx, focus, source)
			fmt.Println("Set", source, t, err)
			return t, err
		},
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					index, focus, err := val.Get()
					ret := yield(ValIE(index, focus, err))
					fmt.Println("Iterate", source, index, focus, err)
					return ret
				})
			}
		},
		func(ctx context.Context, source S) (int, error) {
			l, err := o.AsLengthGetter()(ctx, source)
			fmt.Println("GetLength", source, l, err)
			return l, err

		},
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {
			t, err := o.AsModify()(ctx, func(index I, focus A) (B, error) {
				b, err := fmap(index, focus)
				fmt.Println("Modify", source, index, focus, b, err)
				return b, err
			}, source)
			return t, err
		},
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				o.AsIxGetter()(ctx, index, source)(func(val ValueIE[I, A]) bool {
					focusIndex, focus, err := val.Get()
					ret := yield(ValIE(focusIndex, focus, err))
					fmt.Println("IxGet", source, index, focusIndex, focus, err, ret)
					return ret
				})
			}
		},
		func(indexA, indexB I) bool {
			match := o.AsIxMatch()(indexA, indexB)
			fmt.Println("IxMatch", indexA, indexB, match)
			return match
		},
		func(ctx context.Context, focus B) (T, error) {
			t, err := o.AsReverseGetter()(ctx, focus)
			fmt.Println("ReverseGet", focus, t, err)
			return t, err
		},
		ExprCustom("WithExampleLogging"),
	)
}

func ExampleOmni() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	optic := WithExampleLogging(TraverseSlice[string]())

	result := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println("Result:", result)

	result = MustModify(optic, Op(strings.ToUpper), data)
	fmt.Println("Result:", result)

	//Output:
	//Iterate [alpha beta gamma delta] 0 alpha <nil>
	//Iterate [alpha beta gamma delta] 1 beta <nil>
	//Iterate [alpha beta gamma delta] 2 gamma <nil>
	//Iterate [alpha beta gamma delta] 3 delta <nil>
	//Result: [alpha beta gamma delta]
	//Modify [alpha beta gamma delta] 0 alpha ALPHA <nil>
	//Modify [alpha beta gamma delta] 1 beta BETA <nil>
	//Modify [alpha beta gamma delta] 2 gamma GAMMA <nil>
	//Modify [alpha beta gamma delta] 3 delta DELTA <nil>
	//Result: [ALPHA BETA GAMMA DELTA]

}
