package optic_test

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/samber/lo"

	. "github.com/spearson78/go-optic"
)

func TestMathConsistency(t *testing.T) {

	ValidateOpticTest(t, OpT2ToOptic(T2A[int, int](), MulT2[int](), T2B[int, int]()), lo.T2(2, 10), 0)
	ValidateOpticTest(t, PredT2ToOptic(T2A[bool, bool](), AndT2(), T2B[bool, bool]()), lo.T2(true, false), false)

	ValidateOpticTest(t, Mul(2), 10, 100)
	ValidateOpticTest(t, Div(2), 10, 100)
	ValidateOpticTest(t, Mod(2), 13, 0)
	ValidateOpticTest(t, Add(2), 10, 100)
	ValidateOpticTest(t, Sub(2), 10, 100)
	ValidateOpticTest(t, Pow(2), 10, 100)

	ValidateOpticTest(t, Negate[int](), 10, 100)
	ValidateOpticTest(t, Abs[int](), 10, 100)
	ValidateOpticTest(t, Abs[int](), -10, -100)

	ValidateOpticTest(t, ParseInt[int](10, 32), "10", 100)
	ValidateOpticTestP(t, ParseIntP[int](10, 32), "10", 100,
		func(i int) int { return i }, func(i int) string {
			return strconv.FormatInt(int64(i), 10)
		},
	)

	ValidateOpticTest(t, ParseFloat[float64]('f', 1, 64), "10.5", 100.5)

	ValidateOpticTest(t, Min(5), 1, 0)
	ValidateOpticTest(t, Min(1), 5, 0)

	ValidateOpticTest(t, Max(5), 1, 0)
	ValidateOpticTest(t, Max(1), 5, 0)

	ValidateOpticTest(t, Clamp(5, 7), 1, 0)
	ValidateOpticTest(t, Clamp(5, 7), 6, 0)
	ValidateOpticTest(t, Clamp(5, 7), 10, 0)

}

func TestDividingMultipling(t *testing.T) {
	if r, err := Modify(Compose(Div(10.0), Mul(2.0)), Add(1.0), 30.0); err != nil || r != 35 {
		t.Fatalf("Iso Varieties 10 %v", r)
	}
}

func ExampleBinaryOp() {

	customGreaterThan := BinaryOp(func(left, right int) bool { return left > right }, ">")

	data := []int{1, 2, 3, 4, 5}

	var filtered []int = MustGet(SliceOf(Filtered(TraverseSlice[int](), OpT2Bind(customGreaterThan, 3)), len(data)), data)
	fmt.Println(filtered)

	var overResult []bool = MustModify(TraverseSliceP[int, bool](), OpT2Bind(customGreaterThan, 3), data)
	fmt.Println(overResult)

	//Output:[4 5]
	//[false false false true true]
}

func ExampleBinaryOpE() {

	customDivide := BinaryOpE(func(ctx context.Context, left, right float64) (float64, error) {
		if right == 0 {
			return 0, errors.New("divide by 0")
		}
		return left / right, nil
	}, "/")

	data := []float64{10, 20}

	var overResult []float64
	overResult, err := Modify(TraverseSlice[float64](), OpT2Bind(customDivide, 2), data)
	fmt.Println(overResult, err)

	//Output:[5 10] <nil>
}

func ExampleUnaryOp() {

	//string length is non reversible so an iso cannot be used.
	strLen := UnaryOp(func(val string) int { return len(val) }, "len")

	data := []string{"alpha", "beta", "gamma", "delta", "epsilon", "zeta"}

	var filtered []string = MustGet(SliceOf(Filtered(TraverseSlice[string](), Compose(strLen, Gt(4))), len(data)), data)
	fmt.Println(filtered)

	var strLens []int = MustGet(SliceOf(Compose(TraverseSlice[string](), strLen), len(data)), data)
	fmt.Println(strLens)

	var overResult []int = MustModify(TraverseSliceP[string, int](), strLen, data)
	fmt.Println(overResult)

	//Output:
	//[alpha gamma delta epsilon]
	//[5 4 5 5 7 4]
	//[5 4 5 5 7 4]
}

func ExampleOpT2IToOptic() {

	data := []int{2, 4, 6, 8}

	//This optic focuses each int in the slice
	optic := TraverseSlice[int]()

	//This operator will encode the slice index of the int in the 100s section
	op := OpT2IToOptic(
		OpToOpI[int]( //IxOp enables Identity to be used with indexed operators
			Identity[int](),
		),
		OpT2ToOpT2I[Void, int](AddT2[int]()), //IxOpT2 enables Add to be used with indexed operators
		ComposeLeft(
			ValueIIndex[int, int](),
			Mul(100),
		),
	)

	var result []int = MustModifyI(
		optic,
		op,
		data,
	)
	fmt.Println(result)

	//Output:
	//[2 104 206 308]
}

func ExampleOpT2ToOptic() {

	data := [][]int{
		{1, 1},
		{1, 2},
		{2, 1},
		{2, 2},
	}

	//returns a operator that returns s[0] * s[1] of a slice
	operator := OpT2ToOptic(
		FirstOrDefault(Index(TraverseSlice[int](), 0), 0),
		MulT2[int](),
		FirstOrDefault(Index(TraverseSlice[int](), 1), 0),
	)

	var overResult []int = MustModify(
		TraverseSliceP[[]int, int](),
		operator,
		data,
	)
	fmt.Println(overResult)

	//Output:
	//[1 2 2 4]
}

func ExampleMul() {

	data := []lo.Tuple2[string, int]{
		lo.T2("a", 1),
		lo.T2("b", 2),
		lo.T2("c", 3),
		lo.T2("d", 4),
		lo.T2("e", 5),
	}

	//Focus the second element in a slice of tuples
	optic := Compose(TraverseSlice[lo.Tuple2[string, int]](), T2B[string, int]())

	//Multiply each focused int by 10.
	var overResult []lo.Tuple2[string, int] = MustModify(optic, Mul(10), data)
	fmt.Println(overResult)

	//Getters are optics and can be composed.
	//This composition multiples by 10 then adds 1
	var composedResult []lo.Tuple2[string, int] = MustModify(optic, Compose(Mul(10), Add(1)), data)
	fmt.Println(composedResult)

	//The getter can be composed with the query optic.
	//This has the effect of applying the Add operation in a multipleied by 10 context
	//But when the result is built the value is divided by 10 again to return it to the original context
	//In effect the Add(100) is in effect an Add(10)
	var preComposedResult []lo.Tuple2[string, int] = MustModify(Compose(optic, Mul(10)), Add(100), data)
	fmt.Println(preComposedResult)

	//Output:[{a 10} {b 20} {c 30} {d 40} {e 50}]
	//[{a 11} {b 21} {c 31} {d 41} {e 51}]
	//[{a 11} {b 12} {c 13} {d 14} {e 15}]
}

func ExampleMulT2() {

	data := []int{1, 2, 3, 4}

	//returns a reducer that multiplies all the numbers in a slice
	optic := Reduce(
		TraverseSlice[int](),
		AsReducer(1, MulT2[int]()),
	)

	var res int
	res, ok := MustGetFirst(
		optic,
		data,
	)
	fmt.Println(res, ok)

	//Output:
	//24 true
}

func ExampleMulOp() {

	data := []lo.Tuple2[int, int]{
		lo.T2(1, 6),
		lo.T2(2, 7),
		lo.T2(3, 8),
		lo.T2(4, 9),
		lo.T2(5, 10),
	}

	//2 optics to focus on the first and second elements of the tuple
	leftOptic := T2A[int, int]()
	rightOptic := T2B[int, int]()

	//Create a getter that applies the Mul function to the focuses of both optics
	multiplier := MulOp(leftOptic, rightOptic)

	var singleValueResult int = MustGet(multiplier, lo.T2(2, 10))
	fmt.Println(singleValueResult) // 20

	//The created function is suitable for usage together with polymorphic optics that have a source focus of lo.Tuple2[int, int] but result focus type of int
	var overResult []int = MustModify(TraverseSliceP[lo.Tuple2[int, int], int](), multiplier, data)
	fmt.Println(overResult) // [6 14 24 36 50]

	//Output:20
	//[6 14 24 36 50]
}

func ExampleProduct() {
	data := []int{1, 2, 3, 4}

	product, ok := MustGetFirst(Reduce(TraverseSlice[int](), Product[int]()), data)

	fmt.Println(product, ok)
	//Output:24 true
}

func ExampleDiv() {

	data := []lo.Tuple2[string, int]{
		lo.T2("a", 10),
		lo.T2("b", 20),
		lo.T2("c", 30),
		lo.T2("d", 40),
		lo.T2("e", 50),
	}

	//Focus the second element in a slice of tuples
	optic := Compose(TraverseSlice[lo.Tuple2[string, int]](), T2B[string, int]())

	//Divide each focused int by 10.
	var overResult []lo.Tuple2[string, int] = MustModify(optic, Div(10), data)
	fmt.Println(overResult)

	//Getters are optics and can be composed.
	//This composition Divides by 10 then adds 1
	var composedResult []lo.Tuple2[string, int] = MustModify(optic, Compose(Div(10), Add(10)), data)
	fmt.Println(composedResult)

	//The getter can be composed with the query optic.
	//This has the effect of applying the Add operation in a divided by 10 context
	//But when the result is built the value is multiplied by 10 again to return it to the original context
	//In effect the Add(10) is in effect an Add(100)
	var preComposedResult []lo.Tuple2[string, int] = MustModify(Compose(optic, Div(10)), Add(10), data)
	fmt.Println(preComposedResult)

	//Output:[{a 1} {b 2} {c 3} {d 4} {e 5}]
	//[{a 11} {b 12} {c 13} {d 14} {e 15}]
	//[{a 110} {b 120} {c 130} {d 140} {e 150}]
}

func ExampleDivT2() {

	data := []int{5, 2}

	//returns a reducer that divides 100 sequentially by all the numbers in a slice
	optic := Reduce(
		TraverseSlice[int](),
		AsReducer(100, DivT2[int]()),
	)

	var res int
	res, ok := MustGetFirst(
		optic,
		data,
	)
	fmt.Println(res, ok)

	//Output:
	//10 true
}

func ExampleDivOp() {

	data := []lo.Tuple2[int, int]{
		lo.T2(10, 6),
		lo.T2(20, 7),
		lo.T2(30, 8),
		lo.T2(40, 9),
		lo.T2(50, 10),
	}

	//2 optics to focus on the first and second elements of the tuple
	leftOptic := T2A[int, int]()
	rightOptic := T2B[int, int]()

	//Create a getter that applies the Div function to the focuses of both optics
	divider := DivOp(leftOptic, rightOptic)

	var singleValueResult int = MustGet(divider, lo.T2(20, 10))
	fmt.Println(singleValueResult) // 2

	//The created getter is suitable for usage together with polymorphic optics that have a source focus of lo.Tuple2[int, int] but result focus type of int
	var overResult []int = MustModify(TraverseSliceP[lo.Tuple2[int, int], int](), divider, data)
	fmt.Println(overResult) // [1 2 3 4 5]

	//Output:2
	//[1 2 3 4 5]
}

func ExampleDivQuotient() {

	data := []lo.Tuple2[string, int]{
		lo.T2("a", 10),
		lo.T2("b", 20),
		lo.T2("c", 30),
		lo.T2("d", 40),
		lo.T2("e", 50),
	}

	//Focus the second element in a slice of tuples
	optic := Compose(TraverseSlice[lo.Tuple2[string, int]](), T2B[string, int]())

	//Divide 100 by each focused int.
	var modifyResult []lo.Tuple2[string, int] = MustModify(optic, DivQuotient(200), data)
	fmt.Println(modifyResult)

	//Getters are optics and can be composed.
	//This composition Divides by 10 then adds 1
	var composedResult []lo.Tuple2[string, int] = MustModify(optic, Compose(DivQuotient(200), Add(10)), data)
	fmt.Println(composedResult)

	//The getter can be composed with the query optic.
	//This has the effect of applying the Add operation in a divided by 10 context
	//But when the result is built the value is multiplied by 10 again to return it to the original context
	//In effect the Add(10) is in effect an Add(100)
	var preComposedResult []lo.Tuple2[string, int] = MustModify(Compose(optic, Div(10)), Add(10), data)
	fmt.Println(preComposedResult)

	//Output:
	//[{a 20} {b 10} {c 6} {d 5} {e 4}]
	//[{a 30} {b 20} {c 16} {d 15} {e 14}]
	//[{a 110} {b 120} {c 130} {d 140} {e 150}]
}

func ExampleMod() {

	data := []lo.Tuple2[string, int]{
		lo.T2("a", 10),
		lo.T2("b", 21),
		lo.T2("c", 32),
		lo.T2("d", 43),
		lo.T2("e", 55),
	}

	//Focus the second element in a slice of tuples
	optic := Compose(TraverseSlice[lo.Tuple2[string, int]](), T2B[string, int]())

	//Perform modulo 10 on each focused int.
	var overResult []lo.Tuple2[string, int] = MustModify(optic, Mod(10), data)
	fmt.Println(overResult)

	//Getters are optics and can be composed.
	//This composition performs Modulo 10 then adds 10
	var composedResult []lo.Tuple2[string, int] = MustModify(optic, Compose(Mod(10), Add(10)), data)
	fmt.Println(composedResult)

	//Output:[{a 0} {b 1} {c 2} {d 3} {e 5}]
	//[{a 10} {b 11} {c 12} {d 13} {e 15}]
}

func ExampleModOp() {

	data := []lo.Tuple2[int, int]{
		lo.T2(10, 6),
		lo.T2(20, 7),
		lo.T2(30, 8),
		lo.T2(40, 9),
		lo.T2(50, 10),
	}

	//2 optics to focus on the first and second elements of the tuple
	leftOptic := T2A[int, int]()
	rightOptic := T2B[int, int]()

	//Create a getter that applies the Mod function to the focuses of both optics
	moder := ModOp(leftOptic, rightOptic)

	var singleValueResult int = MustGet(moder, lo.T2(20, 7))
	fmt.Println(singleValueResult) // 6

	//The created getter is suitable for usage together with polymorphic optics that have a source focus of lo.Tuple2[int, int] but result focus type of int
	var overResult []int = MustModify(TraverseSliceP[lo.Tuple2[int, int], int](), moder, data)
	fmt.Println(overResult) // [4 6 6 4 0]

	//Output:6
	//[4 6 6 4 0]
}

func ExampleModT2() {

	data := []lo.Tuple2[int, int]{
		lo.T2(5, 2),
		lo.T2(6, 2),
		lo.T2(5, 3),
		lo.T2(6, 3),
	}

	optic := Compose(
		TraverseSlice[lo.Tuple2[int, int]](),
		ModT2[int](),
	)

	var res []int = MustGet(
		SliceOf(
			optic,
			4,
		),
		data,
	)
	fmt.Println(res)

	//Output:
	//[1 0 2 0]
}

func ExampleAdd() {

	data := []lo.Tuple2[string, int]{
		lo.T2("a", 1),
		lo.T2("b", 2),
		lo.T2("c", 3),
		lo.T2("d", 4),
		lo.T2("e", 5),
	}

	//Focus the second element in a slice of tuples
	optic := Compose(TraverseSlice[lo.Tuple2[string, int]](), T2B[string, int]())

	//Add 10 to each focused int.
	var overResult []lo.Tuple2[string, int] = MustModify(optic, Add(10), data)
	fmt.Println(overResult)

	//Getters are optics and can be composed.
	//This composition Adds 10 then multiplies by 2
	var composedResult []lo.Tuple2[string, int] = MustModify(optic, Compose(Add(10), Mul(2)), data)
	fmt.Println(composedResult)

	//The getter can be composed with the query optic.
	//This has the effect of applying the Mul operation in a add 10 context
	//But when the result is built the value is subtracted by 10 again to return it to the original context
	var preComposedResult []lo.Tuple2[string, int] = MustModify(Compose(optic, Add(10)), Mul(2), data)
	fmt.Println(preComposedResult)

	//Output:[{a 11} {b 12} {c 13} {d 14} {e 15}]
	//[{a 22} {b 24} {c 26} {d 28} {e 30}]
	//[{a 12} {b 14} {c 16} {d 18} {e 20}]
}

func ExampleAddOp() {

	data := []lo.Tuple2[int, int]{
		lo.T2(1, 6),
		lo.T2(2, 7),
		lo.T2(3, 8),
		lo.T2(4, 9),
		lo.T2(5, 10),
	}

	//2 optics to focus on the first and second elements of the tuple
	leftOptic := T2A[int, int]()
	rightOptic := T2B[int, int]()

	//Create a getter that applies the Add function to the focuses of both optics
	adder := AddOp(leftOptic, rightOptic)

	var singleValueResult int = MustGet(adder, lo.T2(2, 10))
	fmt.Println(singleValueResult) // 12

	//The created function is suitable for usage together with polymorphic optics that have a source focus of lo.Tuple2[int, int] but result focus type of int
	var overResult []int = MustModify(TraverseSliceP[lo.Tuple2[int, int], int](), adder, data)
	fmt.Println(overResult) // [6 14 24 36 50]

	//Output:12
	//[7 9 11 13 15]
}

func ExampleSum() {
	data := []int{1, 2, 3, 4}

	sum, ok := MustGetFirst(Reduce(TraverseSlice[int](), Sum[int]()), data)

	fmt.Println(sum, ok)
	//Output:10 true
}

func ExampleAddT2() {

	data := []int{5, 2, 10, 3}

	//returns a reducer that sums all the numbers in a slice
	optic := Reduce(
		TraverseSlice[int](),
		AsReducer(0, AddT2[int]()),
	)

	var res int
	res, ok := MustGetFirst(
		optic,
		data,
	)
	fmt.Println(res, ok)

	//Output:
	//20 true
}

func ExampleSub() {

	data := []lo.Tuple2[string, int]{
		lo.T2("a", 11),
		lo.T2("b", 22),
		lo.T2("c", 33),
		lo.T2("d", 44),
		lo.T2("e", 55),
	}

	//Focus the second element in a slice of tuples
	optic := Compose(TraverseSlice[lo.Tuple2[string, int]](), T2B[string, int]())

	//Sub 10 to each focused int.
	var overResult []lo.Tuple2[string, int] = MustModify(optic, Sub(10), data)
	fmt.Println(overResult)

	//Getters are optics and can be composed.
	//This composition subtracts 10 then multiplies by 2
	var composedResult []lo.Tuple2[string, int] = MustModify(optic, Compose(Sub(10), Mul(2)), data)
	fmt.Println(composedResult)

	//The getter can be composed with the query optic.
	//This has the effect of applying the Mul operation in a subtract 10 context
	//But when the result is built the value is added by 10 again to return it to the original context
	var preComposedResult []lo.Tuple2[string, int] = MustModify(Compose(optic, Sub(10)), Mul(2), data)
	fmt.Println(preComposedResult)

	//Output:[{a 1} {b 12} {c 23} {d 34} {e 45}]
	//[{a 2} {b 24} {c 46} {d 68} {e 90}]
	//[{a 12} {b 34} {c 56} {d 78} {e 100}]
}

func ExampleSubOp() {

	data := []lo.Tuple2[int, int]{
		lo.T2(6, 1),
		lo.T2(7, 2),
		lo.T2(8, 3),
		lo.T2(9, 4),
		lo.T2(10, 5),
	}

	//2 optics to focus on the first and second elements of the tuple
	leftOptic := T2A[int, int]()
	rightOptic := T2B[int, int]()

	//Create a getter that applies the Sub function to the focuses of both optics
	subber := SubOp(leftOptic, rightOptic)

	var singleValueResult int = MustGet(subber, lo.T2(10, 2))
	fmt.Println(singleValueResult) // 8

	//The created function is suitable for usage together with polymorphic optics that have a source focus of lo.Tuple2[int, int] but result focus type of int
	var overResult []int = MustModify(TraverseSliceP[lo.Tuple2[int, int], int](), subber, data)
	fmt.Println(overResult) // [5 5 5 5 5]

	//Output:8
	//[5 5 5 5 5]
}

func ExampleSubT2() {

	data := []int{5, 2, 10, 3}

	//returns a reducer that subtracts all the numbers in a slice from 100
	optic := Reduce(
		TraverseSlice[int](),
		AsReducer(100, SubT2[int]()),
	)

	var res int
	res, ok := MustGetFirst(
		optic,
		data,
	)
	fmt.Println(res, ok)

	//Output:
	//80 true
}

func ExampleSubFrom() {

	data := []lo.Tuple2[string, int]{
		lo.T2("a", 11),
		lo.T2("b", 22),
		lo.T2("c", 33),
		lo.T2("d", 44),
		lo.T2("e", 55),
	}

	//Focus the second element in a slice of tuples
	optic := Compose(TraverseSlice[lo.Tuple2[string, int]](), T2B[string, int]())

	//Subtract each focused value from 100.
	var modifyResult []lo.Tuple2[string, int] = MustModify(optic, SubFrom(100), data)
	fmt.Println(modifyResult)

	//Getters are optics and can be composed.
	//This composition subtracts from 100 then multiplies by 2
	var composedResult []lo.Tuple2[string, int] = MustModify(optic, Compose(SubFrom(100), Mul(2)), data)
	fmt.Println(composedResult)

	//The getter can be composed with the query optic.
	//This has the effect of applying the Mul operation in a subtract from 100 context
	//But when the result is built the value is added o again to return it to the original context
	var preComposedResult []lo.Tuple2[string, int] = MustModify(Compose(optic, SubFrom(10)), Mul(2), data)
	fmt.Println(preComposedResult)

	//Output:
	//[{a 89} {b 78} {c 67} {d 56} {e 45}]
	//[{a 178} {b 156} {c 134} {d 112} {e 90}]
	//[{a 12} {b 34} {c 56} {d 78} {e 100}]
}

func ExamplePow() {

	data := []lo.Tuple2[string, float64]{
		lo.T2("a", 1.0),
		lo.T2("b", 2.0),
		lo.T2("c", 3.0),
		lo.T2("d", 4.0),
		lo.T2("e", 5.0),
	}

	//Focus the second element in a slice of tuples
	optic := Compose(TraverseSlice[lo.Tuple2[string, float64]](), T2B[string, float64]())

	//Raise each focused value to the power of 2.
	var overResult []lo.Tuple2[string, float64] = MustModify(optic, Pow(2.0), data)
	fmt.Println(overResult)

	//Getters are optics and can be composed.
	//This composition raises to the power 2 then adds 1
	var composedResult []lo.Tuple2[string, float64] = MustModify(optic, Compose(Pow(2.0), Add(1.0)), data)
	fmt.Println(composedResult)

	//The getter can be composed with the query optic.
	//This has the effect of applying the Add operation in a raise to the power 2 context
	//But when the result is built the value is lowered by root 2 again to return it to the original context
	var preComposedResult []lo.Tuple2[string, float64] = MustModify(Compose(optic, Pow(2.0)), Add(1.0), data)
	fmt.Println(preComposedResult)

	//Output:[{a 1} {b 4} {c 9} {d 16} {e 25}]
	//[{a 2} {b 5} {c 10} {d 17} {e 26}]
	//[{a 1.4142135623730951} {b 2.23606797749979} {c 3.1622776601683795} {d 4.123105625617661} {e 5.0990195135927845}]
}

func ExamplePowOp() {

	data := []lo.Tuple2[float64, float64]{
		lo.T2(1.0, 5.0),
		lo.T2(2.0, 4.0),
		lo.T2(3.0, 3.0),
		lo.T2(4.0, 2.0),
		lo.T2(5.0, 1.0),
	}

	//2 optics to focus on the first and second elements of the tuple
	leftOptic := T2A[float64, float64]()
	rightOptic := T2B[float64, float64]()

	//Create a getter that applies the Pow function to the focuses of both optics
	power := PowOp(leftOptic, rightOptic)

	var singleValueResult float64 = MustGet(power, lo.T2(2.0, 4.0))
	fmt.Println(singleValueResult) // 16

	//The created function is suitable for usage together with polymorphic optics that have a source focus of lo.Tuple2[int, int] but result focus type of int
	var overResult []float64 = MustModify(TraverseSliceP[lo.Tuple2[float64, float64], float64](), power, data)
	fmt.Println(overResult) // [1 16 27 16 5]

	//Output:16
	//[1 16 27 16 5]
}

func ExamplePowT2() {

	data := []lo.Tuple2[float64, float64]{
		lo.T2(5.0, 2.0),
		lo.T2(6.0, 2.0),
		lo.T2(5.0, 3.0),
		lo.T2(6.0, 3.0),
	}

	optic := Compose(
		TraverseSlice[lo.Tuple2[float64, float64]](),
		PowT2(),
	)

	var res []float64 = MustGet(
		SliceOf(
			optic,
			4,
		),
		data,
	)
	fmt.Println(res)

	//Output:
	//[25 36 125 216]
}

func ExampleRoot() {

	data := []lo.Tuple2[string, float64]{
		lo.T2("a", 1.0),
		lo.T2("b", 2.0),
		lo.T2("c", 3.0),
		lo.T2("d", 4.0),
		lo.T2("e", 5.0),
	}

	//Focus the second element in a slice of tuples
	optic := Compose(TraverseSlice[lo.Tuple2[string, float64]](), T2B[string, float64]())

	//Square root of each focused value.
	var overResult []lo.Tuple2[string, float64] = MustModify(optic, Root(2.0), data)
	fmt.Println(overResult)

	//Getters are optics and can be composed.
	var composedResult []lo.Tuple2[string, float64] = MustModify(optic, Compose(Root(2.0), Add(1.0)), data)
	fmt.Println(composedResult)

	//The getter can be composed with the query optic.
	//This has the effect of applying the Add operation in a root 2 context
	//But when the result is built the value is raised by power 2 again to return it to the original context
	var preComposedResult []lo.Tuple2[string, float64] = MustModify(Compose(optic, Root(2.0)), Add(1.0), data)
	fmt.Println(preComposedResult)

	//Output:
	//[{a 1} {b 1.4142135623730951} {c 1.7320508075688772} {d 2} {e 2.23606797749979}]
	//[{a 2} {b 2.414213562373095} {c 2.732050807568877} {d 3} {e 3.23606797749979}]
	//[{a 4} {b 5.82842712474619} {c 7.464101615137754} {d 9} {e 10.47213595499958}]
}

func ExampleRootOp() {

	data := []lo.Tuple2[float64, float64]{
		lo.T2(1.0, 5.0),
		lo.T2(2.0, 4.0),
		lo.T2(3.0, 3.0),
		lo.T2(4.0, 2.0),
		lo.T2(5.0, 1.0),
	}

	//2 optics to focus on the first and second elements of the tuple
	leftOptic := T2A[float64, float64]()
	rightOptic := T2B[float64, float64]()

	//Create a getter that applies the Root function to the focuses of both optics
	root := RootOp(leftOptic, rightOptic)

	var singleValueResult float64 = MustGet(root, lo.T2(2.0, 4.0))
	fmt.Println(singleValueResult) // 16

	//The created function is suitable for usage together with polymorphic optics that have a source focus of lo.Tuple2[int, int] but result focus type of int
	var overResult []float64 = MustModify(TraverseSliceP[lo.Tuple2[float64, float64], float64](), root, data)
	fmt.Println(overResult) // [1 16 27 16 5]

	//Output:
	//1.189207115002721
	//[1 1.189207115002721 1.4422495703074085 2 5]
}

func ExampleRootT2() {

	data := []lo.Tuple2[float64, float64]{
		lo.T2(5.0, 2.0),
		lo.T2(6.0, 2.0),
		lo.T2(5.0, 3.0),
		lo.T2(6.0, 3.0),
	}

	optic := Compose(
		TraverseSlice[lo.Tuple2[float64, float64]](),
		RootT2(),
	)

	var res []float64 = MustGet(
		SliceOf(
			optic,
			4,
		),
		data,
	)
	fmt.Println(res)

	//Output:
	//[2.23606797749979 2.449489742783178 1.7099759466766968 1.8171205928321397]
}

func ExampleNegate() {

	data := []lo.Tuple2[string, int]{
		lo.T2("a", 11),
		lo.T2("b", 22),
		lo.T2("c", 33),
		lo.T2("d", 44),
		lo.T2("e", 55),
	}

	//Focus the second element in a slice of tuples
	optic := Compose(TraverseSlice[lo.Tuple2[string, int]](), T2B[string, int]())

	//Negate each focused int.
	var overResult []lo.Tuple2[string, int] = MustModify(optic, Negate[int](), data)
	fmt.Println(overResult)

	//Getters are optics and can be composed.
	//This composition negates then Adds 2
	var composedResult []lo.Tuple2[string, int] = MustModify(optic, Compose(Negate[int](), Add(2)), data)
	fmt.Println(composedResult)

	//The getter can be composed with the query optic.
	//This has the effect of applying the Add operation in a negated context
	//But when the result is built the value is negated again to return it to the original context.
	var preComposedResult []lo.Tuple2[string, int] = MustModify(Compose(optic, Negate[int]()), Add(2), data)
	fmt.Println(preComposedResult)

	//Output:[{a -11} {b -22} {c -33} {d -44} {e -55}]
	//[{a -9} {b -20} {c -31} {d -42} {e -53}]
	//[{a 9} {b 20} {c 31} {d 42} {e 53}]
}

func ExampleAbs() {

	data := []lo.Tuple2[string, int]{
		lo.T2("a", 1),
		lo.T2("b", -2),
		lo.T2("c", 3),
		lo.T2("d", -4),
		lo.T2("e", 5),
	}

	//Focus the second element in a slice of tuples
	optic := Compose(TraverseSlice[lo.Tuple2[string, int]](), T2B[string, int]())

	//Abs each focused int.
	var overResult []lo.Tuple2[string, int] = MustModify(optic, Abs[int](), data)
	fmt.Println(overResult)

	//Getters are optics and can be composed.
	//This composition negates then Adds 2
	var composedResult []lo.Tuple2[string, int] = MustModify(optic, Compose(Abs[int](), Add(2)), data)
	fmt.Println(composedResult)

	//The getter can be composed with the query optic.
	var preComposedResult []int = MustGet(SliceOf(Compose(optic, Abs[int]()), 5), data)
	fmt.Println(preComposedResult)

	//Output:
	//[{a 1} {b 2} {c 3} {d 4} {e 5}]
	//[{a 3} {b 4} {c 5} {d 6} {e 7}]
	//[1 2 3 4 5]
}

func ExampleParseInt() {

	data := []string{"1", "2", "3"}

	var composedResult []string
	composedResult, err := Modify(
		Compose(
			TraverseSlice[string](),
			ParseInt[int](10, 32),
		),
		Mul(2),
		data,
	)
	fmt.Println(composedResult, err)

	//Output:
	//[2 4 6] <nil>
}

func ExampleParseIntP() {

	data := []string{"1", "2", "3"}

	var composedResult []int
	composedResult, err := Modify(
		Compose(
			TraverseSliceP[string, int](),
			ParseIntP[int](10, 32),
		),
		Mul(2),
		data,
	)
	fmt.Println(composedResult, err)

	//Output:
	//[2 4 6] <nil>
}

func ExampleFormatInt() {

	data := []int{1, 2, 3}

	var composedResult []string = MustGet(
		SliceOf(
			Compose(
				TraverseSlice[int](),
				FormatInt[int](10),
			),
			len(data),
		),
		data,
	)
	fmt.Println(composedResult)

	//Output:
	//[1 2 3]
}

func ExampleParseFloat() {

	data := []string{"1.1", "2.2", "3.3"}

	var composedResult []string
	composedResult, err := Modify(
		Compose(
			TraverseSlice[string](),
			ParseFloat[float64]('f', -1, 64),
		),
		Mul(2.0),
		data,
	)
	fmt.Println(composedResult, err)

	//Output:
	//[2.2 4.4 6.6] <nil>
}

func ExampleParseFloatP() {

	data := []string{"1.1", "2.2", "3.3"}

	var composedResult []float64
	composedResult, err := Modify(
		Compose(
			TraverseSliceP[string, float64](),
			ParseFloatP[float64](64),
		),
		Mul(2.0),
		data,
	)
	fmt.Println(composedResult, err)

	//Output:
	//[2.2 4.4 6.6] <nil>
}

func ExampleFormatFloat() {

	data := []float64{1.1, 2.2, 3.3}

	var composedResult []string = MustGet(
		SliceOf(
			Compose(
				TraverseSlice[float64](),
				FormatFloat[float64]('f', -1, 64),
			),
			len(data),
		),
		data,
	)
	fmt.Println(composedResult)

	//Output:
	//[1.1 2.2 3.3]
}

func ExampleMean() {
	data := []float64{1, 1, 1, 20, 40, 45, 50}

	mean, ok := MustGetFirst(Reduce(TraverseSlice[float64](), Mean[float64]()), data)

	fmt.Println(mean, ok)
	//Output:22.571428571428573 true

}

func ExampleMedian() {

	data := []float64{1, 2, 3, 20, 40, 45, 50}

	median, ok := MustGetFirst(Reduce(TraverseSlice[float64](), Median[float64]()), data)

	fmt.Println(median, ok)
	//Output:20 true

}

func ExampleMode() {

	data := []float64{1, 1, 1, 20, 40, 45, 50}

	median, ok := MustGetFirst(Reduce(TraverseSlice[float64](), Mode[float64]()), data)

	fmt.Println(median, ok)
	//Output:1 true

}

func ExampleMin() {

	data := []lo.Tuple2[string, int]{
		lo.T2("a", 1),
		lo.T2("b", 2),
		lo.T2("c", 3),
		lo.T2("d", 4),
		lo.T2("e", 5),
	}

	//Focus the second element in a slice of tuples
	optic := Compose(TraverseSlice[lo.Tuple2[string, int]](), T2B[string, int]())

	//Update each element so the value is 3 or less
	var overResult []lo.Tuple2[string, int] = MustModify(optic, Min(3), data)
	fmt.Println(overResult)

	//Getters are optics and can be composed.
	//This composition applies a minimum of 3 then adds 1
	var composedResult []lo.Tuple2[string, int] = MustModify(optic, Compose(Min(3), Add(1)), data)
	fmt.Println(composedResult)

	//Output:[{a 1} {b 2} {c 3} {d 3} {e 3}]
	//[{a 2} {b 3} {c 4} {d 4} {e 4}]
}

func ExampleMinOp() {

	data := []lo.Tuple2[int, int]{
		lo.T2(1, 6),
		lo.T2(7, 2),
		lo.T2(3, 8),
		lo.T2(9, 4),
		lo.T2(5, 10),
	}

	//2 optics to focus on the first and second elements of the tuple
	leftOptic := T2A[int, int]()
	rightOptic := T2B[int, int]()

	//Create a getter that applies the Mul function to the focuses of both optics
	minimum := MinOp(leftOptic, rightOptic)

	var singleValueResult int = MustGet(minimum, lo.T2(2, 10))
	fmt.Println(singleValueResult) // 2

	//The created function is suitable for usage together with polymorphic optics that have a source focus of lo.Tuple2[int, int] but result focus type of int
	var overResult []int = MustModify(TraverseSliceP[lo.Tuple2[int, int], int](), minimum, data)
	fmt.Println(overResult) // [1 2 3 4 5]

	//Output:2
	//[1 2 3 4 5]
}

func ExampleMinReducer() {
	data := []int{10, 20, 3, 40}

	min, ok := MustGetFirst(Reduce(TraverseSlice[int](), MinReducer[int]()), data)

	fmt.Println(min, ok)

	//If there are no elements to  calculate he minimum of then false is returned
	_, ok = MustGetFirst(Reduce(TraverseSlice[int](), MinReducer[int]()), []int{})

	fmt.Println(ok)
	//Output:3 true
	//false
}

func ExampleMinT2() {

	data := []int{5, 2}

	//returns a reducer that takes the min of all the numbers in a slice
	optic := Reduce(
		TraverseSlice[int](),
		AsReducer(math.MaxInt, MinT2[int]()),
	)

	var res int
	res, ok := MustGetFirst(
		optic,
		data,
	)
	fmt.Println(res, ok)

	//Output:
	//2 true
}

func ExampleMaxT2() {

	data := []int{5, 2}

	//returns a reducer that takes the max of all the numbers in a slice
	optic := Reduce(
		TraverseSlice[int](),
		AsReducer(math.MinInt, MaxT2[int]()),
	)

	var res int
	res, ok := MustGetFirst(
		optic,
		data,
	)
	fmt.Println(res, ok)

	//Output:
	//5 true
}

func TestMinReducerAliasTypes(t *testing.T) {

	type intLike int

	data := []intLike{10, 20, 3, 40}

	if min, ok := MustGetFirst(Reduce(TraverseSlice[intLike](), MinReducer[intLike]()), data); !ok || min != 3 {
		t.Fatalf("min intlike %v %v", min, ok)
	}

	//If there are no elements to  calculate he minimum of then false is returned
	if _, ok := MustGetFirst(Reduce(TraverseSlice[intLike](), MinReducer[intLike]()), []intLike{}); ok {
		t.Fatalf("min intlike empty %v", ok)
	}
}

func TestMaxReducerAliasTypes(t *testing.T) {

	type intLike int

	data := []intLike{10, 20, 3, 40}

	if max, ok := MustGetFirst(Reduce(TraverseSlice[intLike](), MaxReducer[intLike]()), data); !ok || max != 40 {
		t.Fatalf("min intlike %v %v", max, ok)
	}

	//If there are no elements to  calculate he minimum of then false is returned
	if _, ok := MustGetFirst(Reduce(TraverseSlice[intLike](), MaxReducer[intLike]()), []intLike{}); ok {
		t.Fatalf("min intlike empty %v", ok)
	}
}

func ExampleMax() {

	data := []lo.Tuple2[string, int]{
		lo.T2("a", 1),
		lo.T2("b", 2),
		lo.T2("c", 3),
		lo.T2("d", 4),
		lo.T2("e", 5),
	}

	//Focus the second element in a slice of tuples
	optic := Compose(TraverseSlice[lo.Tuple2[string, int]](), T2B[string, int]())

	//Update each element so the value is 3 or greater
	var overResult []lo.Tuple2[string, int] = MustModify(optic, Max(3), data)
	fmt.Println(overResult)

	//Getters are optics and can be composed.
	//This composition applies a maximum of 3 then adds 1
	var composedResult []lo.Tuple2[string, int] = MustModify(optic, Compose(Max(3), Add(1)), data)
	fmt.Println(composedResult)

	//Output:[{a 3} {b 3} {c 3} {d 4} {e 5}]
	//[{a 4} {b 4} {c 4} {d 5} {e 6}]
}

func ExampleMaxOp() {

	data := []lo.Tuple2[int, int]{
		lo.T2(1, 6),
		lo.T2(7, 2),
		lo.T2(3, 8),
		lo.T2(9, 4),
		lo.T2(5, 10),
	}

	//2 optics to focus on the first and second elements of the tuple
	leftOptic := T2A[int, int]()
	rightOptic := T2B[int, int]()

	//Create a getter that applies the Mul function to the focuses of both optics
	minimum := MaxOp(leftOptic, rightOptic)

	var singleValueResult int = MustGet(minimum, lo.T2(2, 10))
	fmt.Println(singleValueResult) // 10

	//The created function is suitable for usage together with polymorphic optics that have a source focus of lo.Tuple2[int, int] but result focus type of int
	var overResult []int = MustModify(TraverseSliceP[lo.Tuple2[int, int], int](), minimum, data)
	fmt.Println(overResult) // [6 7 8 9 10]

	//Output:10
	//[6 7 8 9 10]
}

func ExampleMaxReducer() {
	data := []int{1, 2, 30, 4}

	product, ok := MustGetFirst(Reduce(TraverseSlice[int](), MaxReducer[int]()), data)

	fmt.Println(product, ok)

	//If there are no elements to  calculate the maximum of then false is returned
	_, ok = MustGetFirst(Reduce(TraverseSlice[int](), MaxReducer[int]()), []int{})

	fmt.Println(ok)
	//Output:30 true
	//false
}

func ExampleClamp() {

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var result []int = MustModify(
		TraverseSlice[int](),
		Clamp(3, 6),
		data,
	)
	fmt.Println(result)

	//Output:
	//[3 3 3 4 5 6 6 6 6 6]
}

func TestCelsiusToFahrenheit(t *testing.T) {
	celsiusToFahrenheit := Compose(
		Mul(1.8),
		Add(32.0),
	)

	celsiusResult := MustModify(
		celsiusToFahrenheit,
		Add(1.0),
		32,
	)

	if !(math.Abs(celsiusResult-32.55555555555555) <= 0.000000000001) {
		t.Fatal(celsiusResult)
	}
}
