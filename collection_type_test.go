package optic_test

import (
	"context"
	"fmt"
	"iter"
	"strings"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
)

func ExampleColType() {

	//To create a custom collection we first need an iso that converts
	//from the collection type to a [Collection] and back.
	sliceToCol := CombiIso[ReturnOne, ReadWrite, Pure, []string, []string, Collection[int, string, Pure], Collection[int, string, Pure]](
		//convert slice to collection
		func(ctx context.Context, source []string) (Collection[int, string, Pure], error) {
			return ColI(
				func(yield func(index int, focus string) bool) {
					for i, v := range source {
						if !yield(i, v) {
							return
						}
					}
				},
				func(index int) iter.Seq2[int, string] {
					return func(yield func(index int, focus string) bool) {
						yield(index, source[index])
					}
				},
				IxMatchComparable[int](),
				func() int {
					return len(source)
				},
			), nil
		},
		//convert collection to slice
		func(ctx context.Context, focus Collection[int, string, Pure]) ([]string, error) {
			var ret []string
			for val := range focus.AsIter()(ctx) {
				_, focus, focusErr := val.Get()
				if focusErr != nil {
					return nil, focusErr
				}
				ret = append(ret, focus)
			}
			return ret, nil
		},
		ExprCustom("ExampleCol"),
	)

	//We can now use the iso to create a Collection
	sliceCol := ColType(
		sliceToCol,
		TraverseSlice[string](),
	)

	//We can now convert a col op into a col op for our specific collection

	//Optic[Void, []string, []string, []string, []string, ReturnOne, ReadWrite, UniDir]
	stringReversed := ColTypeOp(sliceCol, ReversedCol[int, string]())

	result := MustGet(stringReversed, []string{"alpha", "beta", "gamma", "delta"})

	fmt.Println(result)
	//Output:
	//[delta gamma beta alpha]
}

func ExampleColTypeP() {

	//We can now use the iso to create a Collection
	colType := ColTypeP(
		ColFocusErr(SliceToColP[string, int]()),
		ColSourceErr(AsReverseGet(SliceToColP[int, string]())),
		EErr(TraverseSliceP[string, int]()),
	)

	op := ColOfP(
		EErr(ComposeLeft(
			TraverseColEP[int, string, int, Err](),
			ParseIntP[int](10, 0),
		)),
	)

	colTypeOp := ColTypeOp(colType, op)

	data := lo.T2(1, []string{"1", "2", "3", "4"})

	var result lo.Tuple2[int, []int]
	result, err := Modify(T2BP[int, []string, []int](), colTypeOp, data)

	fmt.Println(result, err)

	//Output:
	//{1 [1 2 3 4]} <nil>
}

func ExampleTraverseColType() {
	sliceCol := SliceColType[string]()

	traverse := TraverseColType(sliceCol)

	result := MustModify(
		traverse,
		Op(strings.ToUpper),
		[]string{"alpha", "beta", "gamma", "delta"},
	)

	fmt.Println(result)
	//Output:
	//[ALPHA BETA GAMMA DELTA]
}

func ExampleColTypeToCol() {

	colType := ColType(
		SliceToCol[string](),
		TraverseSlice[string](),
	)

	data := []string{"gamma", "beta", "alpha"}

	optic := ColTypeToCol(colType) //[]string --> Collection[int,string,Pure]

	var res []string = MustModify(
		optic,
		ReversedCol[int, string](), //Collection[int,string,Pure]
		data,
	)

	fmt.Println(res)
	//Output:
	//[alpha beta gamma]
}

func ExampleColToColType() {

	colType := ColType(
		SliceToCol[string](),
		TraverseSlice[string](),
	)

	data := []string{"gamma", "beta", "alpha"}

	optic := Compose4(
		ColTypeToCol(colType),        //[]string --> Collection[int,string,Pure]
		FilteredCol[int](Ne("beta")), //Collection[int,string,Pure]
		ReversedCol[int, string](),   //Collection[int,string,Pure]
		ColToColType(colType),        // Collection[int,string,Pure] --> []string
	)

	var res []string = MustGet(
		optic,
		data,
	)

	fmt.Println(res)
	//Output:
	//[alpha gamma]
}

func ExampleColTypeOf() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	optic := ColTypeOf(
		MapColType[int, string](),
		TraverseSlice[string](),
	)

	var viewResult map[int]string = MustGet(optic, data)
	fmt.Println(viewResult)

	//Output:
	//map[0:alpha 1:beta 2:gamma 3:delta]
}

func ExampleColTypeOfP() {

	data := []string{"1", "2", "3", "4"}

	optic := ColTypeOfP(
		MapColTypeP[int, string, int](),
		TraverseSliceP[string, int](),
	)

	var viewResult map[int]string
	viewResult, err := Get(optic, data)
	fmt.Println(viewResult, err)

	var modifyResult []int
	modifyResult, err = Modify(
		optic,
		AsModify(
			TraverseMapP[int, string, int](),
			ParseIntP[int](10, 0),
		),
		data,
	)
	fmt.Println(modifyResult, err)

	//Output:
	//map[0:1 1:2 2:3 3:4] <nil>
	//[1 2 3 4] <nil>
}
