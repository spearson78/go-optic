package optic_test

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	. "github.com/spearson78/go-optic"
)

func ExampleGetModIE() {
	//This lens focuses on the B element of a tuple
	tupleB := GetModIEP[Void, lo.Tuple2[int, int], lo.Tuple2[int, int], int, int](
		func(ctx context.Context, source lo.Tuple2[int, int]) (Void, int, error) {
			return Void{}, source.B, nil
		},
		func(ctx context.Context, fmap func(index Void, focus int) (int, error), source lo.Tuple2[int, int]) (lo.Tuple2[int, int], error) {

			focus := source.B
			newFocus, err := fmap(Void{}, focus)
			if err != nil {
				return source, err
			}

			//LensM provides access to the old and new value at the same time to allow modifications to be detected.
			if newFocus != focus {
				source.B = newFocus
			}

			return source, err
		},
		IxMatchVoid(),
		ExprCustom("ExampleGetModIE"),
	)

	viewResult, err := Get(tupleB, lo.T2(10, 20))
	fmt.Println(viewResult, err)

	overResult, err := Modify(tupleB, Mul(2), lo.T2(10, 20))
	fmt.Println(overResult, err)

	//Output:20 <nil>
	//{10 40} <nil>
}
