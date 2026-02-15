package optic_test

import (
	"fmt"
	"strings"
	"unicode"

	. "github.com/spearson78/go-optic"
)

func ExampleIdentity() {

	data := []int{5, 2, 3, 1, 4}

	res := MustGet(
		SliceOf(
			Ordered(
				TraverseSlice[int](),
				OrderBy(
					Identity[int](), //Order by slice value
				),
			),
			len(data),
		),
		data,
	)
	fmt.Println(res)

	//Output:
	//[1 2 3 4 5]
}

func ExampleConst() {

	data := []int{1, 2, 3, 4, 5}

	res := MustGet(
		SliceOf(
			Compose(
				TraverseSlice[int](),
				If(
					Even[int](),
					Const[int]("even"),
					Const[int]("odd"),
				),
			),
			len(data),
		),
		data,
	)

	fmt.Println(res)

	//Output:
	//[odd even odd even odd]
}

func ExampleConstI() {

	data := map[string]int{
		"alpha": 1,
		"beta":  2,
		"gamma": 3,
	}

	res := MustGet(
		MapOfReduced(
			Concat(
				TraverseMap[string, int](),
				ConstI[map[string]int]("delta", 4, EqT2[string]()),
			),
			FirstReducer[int](),
			len(data),
		),
		data,
	)

	fmt.Println(res)

	//Output:
	//map[alpha:1 beta:2 delta:4 gamma:3]
}

func ExampleNothing() {

	data := [][]string{
		{"alpha", "beta"},
		{"gamma"},
		{"delta", "epsilon"},
	}

	optic := Compose(
		TraverseSlice[[]string](),
		If(
			Compose(
				Length(
					TraverseSlice[string](),
				),
				Lt(2),
			),
			ReIndexed(Nothing[[]string, string](), Const[Void](0), EqT2[int]()),
			TraverseSlice[string](),
		),
	)

	res := MustGet(SliceOf(optic, 4), data)
	fmt.Println(res)

	modRes := MustModify(optic, Op(strings.ToUpper), data)
	fmt.Println(modRes)

	//Output:
	//[alpha beta delta epsilon]
	//[[ALPHA BETA] [gamma] [DELTA EPSILON]]

}

func ExampleNothingI() {

	data := [][]string{
		{"alpha", "beta"},
		{"gamma"},
		{"delta", "epsilon"},
	}

	optic := Compose(
		WithIndex(TraverseSlice[[]string]()),
		If(
			OpOnIx[[]string](Odd[int]()),
			Compose(ValueIValue[int, []string](), ReIndexed(Nothing[[]string, string](), Const[Void](0), EqT2[int]())),
			Compose(ValueIValue[int, []string](), TraverseSlice[string]()),
		),
	)

	res := MustGet(SliceOf(optic, 4), data)
	fmt.Println(res)

	modRes := MustModify(optic, Op(strings.ToUpper), data)
	fmt.Println(modRes)

	//Output:
	//[alpha beta delta epsilon]
	//[[ALPHA BETA] [gamma] [DELTA EPSILON]]

}

func ExampleValue() {

	data := "128"

	specialParseInt := Compose(
		If(
			All(
				TraverseString(),
				Op(unicode.IsDigit),
			),
			Identity[string](),
			Value("-1"),
		),
		ParseInt[int](10, 0),
	)

	res, err := Modify(
		specialParseInt,
		Mul(10),
		data,
	)

	fmt.Println(res, err)

	badData := "12°"

	res, err = Modify(
		specialParseInt,
		Mul(10),
		badData,
	)

	fmt.Println(res, err)

	//Output:
	//1280 <nil>
	//-10 <nil>
}
