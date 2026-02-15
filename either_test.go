package optic_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/samber/mo"

	. "github.com/spearson78/go-optic"
)

func TestEitherConsistency(t *testing.T) {

	ValidateOpticTest(t, Chosen[int](), mo.Left[int, int](5), 10)
	ValidateOpticTest(t, Chosen[int](), mo.Right[int, int](5), 10)

	ValidateOpticTest(t, Left[int, string](), mo.Left[int, string](5), 10)
	ValidateOpticTest(t, Left[int, string](), mo.Right[int, string]("test"), 10)

	ValidateOpticTest(t, Right[int, string](), mo.Left[int, string](5), "update")
	ValidateOpticTest(t, Right[int, string](), mo.Right[int, string]("test"), "update")

	ValidateOpticTest(t, BesideEither(Identity[int](), Identity[int]()), mo.Right[int, int](1), 10)
	ValidateOpticTest(t, BesideEither(Identity[int](), Identity[int]()), mo.Left[int, int](1), 10)

}

func ExampleChosen() {

	chosenLeft := MustGet(Chosen[int](), mo.Left[int, int](1))
	fmt.Println(chosenLeft)

	chosenRight := MustGet(Chosen[int](), mo.Right[int, int](1))
	fmt.Println(chosenRight)

	//As the left value is chosen. a new Left is returned with the set value
	setChosen := MustSet(Chosen[int](), 2, mo.Left[int, int](1))
	fmt.Println(setChosen)

	//Output:
	//1
	//1
	//{true 2 0}
}

func ExampleLeft() {

	left, matched := MustGetFirst(Left[int, int](), mo.Left[int, int](1))
	fmt.Println(left, matched)

	right, matched := MustGetFirst(Left[int, int](), mo.Right[int, int](1))
	fmt.Println(right, matched)

	//Setting Left on a Right value returns the original Left value
	setLeft := MustSet(Left[int, int](), 2, mo.Right[int, int](1))
	fmt.Println(setLeft)

	//Output:
	//1 true
	//0 false
	//{false 0 1}
}

func ExampleRight() {

	left, matched := MustGetFirst(Right[int, int](), mo.Left[int, int](1))
	fmt.Println(left, matched)

	right, matched := MustGetFirst(Right[int, int](), mo.Right[int, int](1))
	fmt.Println(right, matched)

	//Setting Right on a Left value returns the original Left value
	setRight := MustSet(Right[int, int](), 2, mo.Left[int, int](1))
	fmt.Println(setRight)

	//Output:
	//0 false
	//1 true
	//{true 1 0}
}

func ExampleLeftP() {

	left, matched := MustGetFirst(LeftP[int, int, string](), mo.Left[int, int](1))
	fmt.Println(left, matched)

	right, matched := MustGetFirst(Left[int, int](), mo.Right[int, int](1))
	fmt.Println(right, matched)

	var setLeft mo.Either[string, int] = MustModify(
		LeftP[int, int, string](),
		FormatInt[int](10),
		mo.Left[int, int](1),
	)
	fmt.Println(setLeft)

	//Output:
	//1 true
	//0 false
	//{true 1 0}
}

func ExampleRightP() {

	left, matched := MustGetFirst(RightP[int, int, string](), mo.Left[int, int](1))
	fmt.Println(left, matched)

	right, matched := MustGetFirst(RightP[int, int, string](), mo.Right[int, int](1))
	fmt.Println(right, matched)

	setRight := MustModify(
		RightP[int, int, string](),
		FormatInt[int](10),
		mo.Right[int, int](1),
	)
	fmt.Println(setRight)

	//Output:
	//0 false
	//1 true
	//{false 0 1}
}

func ExampleBesideEither() {

	data := []mo.Either[lo.Tuple2[string, int], lo.Tuple2[int, string]]{
		mo.Left[lo.Tuple2[string, int], lo.Tuple2[int, string]](lo.T2("alpha", 1)),
		mo.Right[lo.Tuple2[string, int], lo.Tuple2[int, string]](lo.T2(2, "beta")),
		mo.Left[lo.Tuple2[string, int], lo.Tuple2[int, string]](lo.T2("gamma", 3)),
		mo.Right[lo.Tuple2[string, int], lo.Tuple2[int, string]](lo.T2(4, "delta")),
	}

	leftOptic := T2A[string, int]()
	rightOptic := T2B[int, string]()

	optic := Compose(
		TraverseSlice[mo.Either[lo.Tuple2[string, int], lo.Tuple2[int, string]]](),
		BesideEither(leftOptic, rightOptic),
	)

	var viewResult []string = MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(viewResult)

	overResult := MustModify(optic, Op(strings.ToUpper), data)
	fmt.Println(overResult)

	//Output:
	//[alpha beta gamma delta]
	//[{true {ALPHA 1} {0 }} {false { 0} {2 BETA}} {true {GAMMA 3} {0 }} {false { 0} {4 DELTA}}]
}
