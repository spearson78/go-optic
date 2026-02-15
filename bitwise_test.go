package optic_test

import (
	"fmt"

	. "github.com/spearson78/go-optic"
)

func ExampleBitAnd() {

	data := []int{
		0b000,
		0b001,
		0b010,
		0b011,
		0b100,
		0b101,
		0b110,
		0b111,
	}

	optic := TraverseSlice[int]()

	var viewResult []int = MustGet(SliceOf(Filtered(optic, Compose(BitAnd(0b010), Ne(0))), len(data)), data)
	fmt.Println(viewResult)

	var overResult []int = MustModify(optic, BitAnd(0b010), data)
	fmt.Println(overResult)

	// Getters are optics and can be composed.
	// This composition ands with 0b010 and then ors with 0b001
	var composedResult []int = MustModify(optic, Compose(BitAnd(0b010), BitOr(0b001)), data)
	fmt.Println(composedResult)

	//Output:
	//[2 3 6 7]
	//[0 0 2 2 0 0 2 2]
	//[1 1 3 3 1 1 3 3]
}

func ExampleBitAndT2() {

	data := []int{
		0b1101,
		0b1011,
		0b1111,
	}

	optic := Reduce(
		TraverseSlice[int](),
		AsReducer(
			0b1111,
			BitAndT2[int](),
		),
	)

	viewResult, ok := MustGetFirst(optic, data)
	fmt.Println(viewResult, ok)

	//Output:
	//9 true
}

func ExampleBitOr() {

	data := []int{
		0b000,
		0b001,
		0b010,
		0b011,
		0b100,
		0b101,
		0b110,
		0b111,
	}

	optic := TraverseSlice[int]()

	var overResult []int = MustModify(optic, BitOr(0b010), data)
	fmt.Println(overResult)

	// Getters are optics and can be composed.
	// This composition ands with 0b010 and then ors with 0b001
	var composedResult []int = MustModify(optic, Compose(BitAnd(0b010), BitOr(0b001)), data)
	fmt.Println(composedResult)

	//Output:
	//[2 3 2 3 6 7 6 7]
	//[1 1 3 3 1 1 3 3]
}

func ExampleBitOrT2() {

	data := []int{
		0b0001,
		0b0010,
		0b0100,
	}

	optic := Reduce(
		TraverseSlice[int](),
		AsReducer(
			0b0000,
			BitOrT2[int](),
		),
	)

	viewResult, ok := MustGetFirst(optic, data)
	fmt.Println(viewResult, ok)

	//Output:
	//7 true
}

func ExampleBitXor() {

	data := []int{
		0b000,
		0b001,
		0b010,
		0b011,
		0b100,
		0b101,
		0b110,
		0b111,
	}

	optic := TraverseSlice[int]()

	var overResult []int = MustModify(optic, BitXor(0b010), data)
	fmt.Println(overResult)

	// Getters are optics and can be composed.
	// This composition ands with 0b010 and then xors with 0b010
	var composedResult []int = MustModify(optic, Compose(BitAnd(0b010), BitXor(0b010)), data)
	fmt.Println(composedResult)

	//Output:
	//[2 3 0 1 6 7 4 5]
	//[2 2 0 0 2 2 0 0]
}

func ExampleBitXorT2() {

	data := []int{
		0b0101,
		0b1010,
		0b0100,
	}

	optic := Reduce(
		TraverseSlice[int](),
		AsReducer(
			0b0000,
			BitXorT2[int](),
		),
	)

	viewResult, ok := MustGetFirst(optic, data)
	fmt.Println(viewResult, ok)

	//Output:
	//11 true
}

func ExampleBitNot() {

	data := []int{
		0b000,
		0b001,
		0b010,
		0b011,
		0b100,
		0b101,
		0b110,
		0b111,
	}

	optic := TraverseSlice[int]()

	var overResult []int = MustModify(optic, BitNot[int](), data)
	fmt.Println(overResult)

	// Getters are optics and can be composed.
	// This composition ands with 0b010 and then nots the result
	var composedResult []int = MustModify(optic, Compose(BitAnd(0b010), BitNot[int]()), data)
	fmt.Println(composedResult)

	//Output:
	//[-1 -2 -3 -4 -5 -6 -7 -8]
	//[-1 -1 -3 -3 -1 -1 -3 -3]
}
