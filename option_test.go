package optic_test

import (
	"fmt"
	"testing"

	"github.com/samber/mo"

	. "github.com/spearson78/go-optic"
)

func TestOptionConsistency(t *testing.T) {

	ValidateOpticTest(t, Some[int](), mo.Some(5), 10)
	ValidateOpticTest(t, Some[int](), mo.None[int](), 10)

	ValidateOpticTest(t, None[int](), mo.Some(5), Void{})
	ValidateOpticTest(t, None[int](), mo.None[int](), Void{})

	ValidateOpticTest(t, Non[int](10, EqT2[int]()), mo.Some(5), 1)
	ValidateOpticTest(t, Non[int](10, EqT2[int]()), mo.Some(5), 10)
	ValidateOpticTest(t, Non[int](10, EqT2[int]()), mo.None[int](), 1)
	ValidateOpticTest(t, Non[int](10, EqT2[int]()), mo.None[int](), 10)

	ptrEqual := Compose(
		DelveT2(PtrOption[int]()),
		EqT2[mo.Option[int]](),
	)

	i := 5
	ValidateOpticTestPred(t, PtrOption[int](), &i, mo.Some(5), ptrEqual)
	ValidateOpticTestPred(t, PtrOption[int](), &i, mo.None[int](), ptrEqual)
	ValidateOpticTestPred(t, PtrOption[int](), nil, mo.Some(5), ptrEqual)
	ValidateOpticTestPred(t, PtrOption[int](), nil, mo.None[int](), ptrEqual)

}
func ExamplePresent() {

	data := []mo.Option[int]{
		mo.Some(1),
		mo.None[int](),
		mo.Some(2),
	}

	//View actions like MustToSliceOf unwrap the option if possible and ignore the None elements.
	var result []bool = MustGet(
		SliceOf(
			Compose(
				TraverseSlice[mo.Option[int]](),
				Present[int](),
			),
			len(data)),
		data,
	)
	fmt.Println(result)

	//Output:
	//[true false true]

}

func ExampleSome() {

	data := []mo.Option[int]{
		mo.Some(1),
		mo.None[int](),
		mo.Some(2),
	}

	//View actions like MustListOf unwrap the option if possible and ignore the None elements.
	var result []int = MustGet(SliceOf(Compose(TraverseSlice[mo.Option[int]](), Some[int]()), len(data)), data)
	fmt.Println(result)

	//Modification actions like over cause Some to skip over out the None elements and updates only the values that are present.
	//The composed lense work together to return the original data structure []mo.Option[int] with the applied mapping.
	var overResult []mo.Option[int] = MustModify(Compose(TraverseSlice[mo.Option[int]](), Some[int]()), Mul(10), data)
	fmt.Println(overResult)

	//Output:[1 2]
	//[{true 10} {false 0} {true 20}]
}

func ExampleSomeP() {

	var matchBool mo.Either[mo.Option[bool], int] = MustGet(
		Matching(
			SomeP[int, bool](),
		),
		mo.None[int](),
	)
	fmt.Println(matchBool.MustLeft().IsPresent())

	matchBool = MustGet(
		Matching(
			SomeP[int, bool](),
		),
		mo.Some[int](1),
	)
	fmt.Println(matchBool.MustRight())

	//Output
	//false
	//1

}

func ExampleTraverseOption() {

	data := []mo.Option[int]{
		mo.Some(1),
		mo.None[int](),
		mo.Some(2),
	}

	//View actions like MustToSliceOf unwrap the option if possible and ignore the None elements.
	var result []int = MustGet(SliceOf(Compose(TraverseSlice[mo.Option[int]](), TraverseOption[int]()), len(data)), data)
	fmt.Println(result)

	//Modification actions like over cause TraverseOption to skip over out the None elements and updates only the values that are present.
	//The composed lense work together to return the original data structure []mo.Option[int] with the applied mapping.
	var overResult []mo.Option[int] = MustModify(Compose(TraverseSlice[mo.Option[int]](), TraverseOption[int]()), Mul(10), data)
	fmt.Println(overResult)

	//Output:[1 2]
	//[{true 10} {false 0} {true 20}]
}

func ExampleNone() {

	data := []mo.Option[int]{
		mo.Some(1),
		mo.None[int](),
		mo.Some(2),
	}

	//View actions like MustView ignore the Some elements.
	count := MustGet(Length(Compose(TraverseSlice[mo.Option[int]](), None[int]())), data)
	fmt.Println(count)

	//There are no sensible transformations of the Void type so modification actions make little sense for None.
	var overResult []mo.Option[int] = MustModify(Compose(TraverseSlice[mo.Option[int]](), None[int]()), Identity[Void](), data)
	fmt.Println(overResult)

	//Output:1
	//[{true 1} {false 0} {true 2}]
}

func ExampleNon() {

	data := []mo.Option[int]{
		mo.Some(1),
		mo.None[int](),
		mo.Some(2),
	}

	//View actions like MustListOf convert the none elements to the default value
	var results []int = MustGet(SliceOf(Compose(TraverseSlice[mo.Option[int]](), Non[int](-1, EqT2[int]())), len(data)), data)
	fmt.Println(results)

	//Unlike Some and None we now have default value we can operate on.
	var overResult []mo.Option[int] = MustModify(Compose(TraverseSlice[mo.Option[int]](), Non[int](-1, EqT2[int]())), Mul(10), data)
	fmt.Println(overResult)

	//Warning: setting the value to the default value will cause it to be represented as the None value in the result.
	var setDefaultResult []mo.Option[int] = MustSet(Compose(TraverseSlice[mo.Option[int]](), Non[int](-1, EqT2[int]())), -1, data)
	fmt.Println(setDefaultResult)

	//Output:[1 -1 2]
	//[{true 10} {true -10} {true 20}]
	//[{false 0} {false 0} {false 0}]

}

func ExamplePtrOption() {

	one, three := 1, 3

	data := []*int{
		&one,
		nil,
		&three,
	}

	optic := Compose(TraverseSlice[*int](), PtrOption[int]())

	//View actions like MustListOf convert the none elements to the default value
	var results []mo.Option[int] = MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(results)

	//Some filters out the None values and allows us to operate on the populate values
	var overResult []*int = MustModify(Compose(optic, Some[int]()), Mul(10), data)
	fmt.Println(*overResult[0], overResult[1], *overResult[2])

	//Output:
	//[{true 1} {false 0} {true 3}]
	//10 <nil> 30
}
