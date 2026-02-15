package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/samber/lo"

	. "github.com/spearson78/go-optic"
)

func ExampleCompose() {

	data := []lo.Tuple2[string, int]{
		lo.T2("alpha", 1),
		lo.T2("beta", 2),
		lo.T2("gamma", 3),
		lo.T2("delta", 4),
	}

	//In this example we will fully specify the types of the optics
	//This is is not necessary in normal usage.

	//The left optic we will use in the compose
	var leftOptic Optic[
		int,                      //Index Type
		[]lo.Tuple2[string, int], //Source Type
		[]lo.Tuple2[string, int], //Result Type
		lo.Tuple2[string, int],   //Focus type
		lo.Tuple2[string, int],   //Result focus type
		ReturnMany,               //Return type
		ReadWrite,                //ReadWrite type
		UniDir,                   //Dir type
		Pure,                     //Error type
	] = TraverseSlice[lo.Tuple2[string, int]]()

	//This optic will receive a slice of tuples and focus on the individual tuples.
	//Under modification it will allow the tuple to be edited and combine the edits
	//together into a slice of tuples.

	//The right optic we will use in the compose
	var rightOptic Optic[
		int,                    //Index Type
		lo.Tuple2[string, int], //Source Type
		lo.Tuple2[string, int], //Result Type
		string,                 //Focus type
		string,                 //Result focus type
		ReturnOne,              //Return type
		ReadWrite,              //ReadWrite type
		UniDir,                 //Dir type
		Pure,                   //Error type
	] = T2A[string, int]()

	//This optic will receive a tuple and focus on the first element.
	//Under modification it will allow the first element to be edited
	//and the result is a tuple with a modified first element.

	var composedOptic Optic[
		int,                                    //Index Type
		[]lo.Tuple2[string, int],               //Source Type
		[]lo.Tuple2[string, int],               //Result Type
		string,                                 //Focus type
		string,                                 //Result focus type
		CompositionTree[ReturnMany, ReturnOne], //Return type
		CompositionTree[ReadWrite, ReadWrite],  //ReadWrite type
		CompositionTree[UniDir, UniDir],        //Dir type
		CompositionTree[Pure, Pure],            //Error type
	] = Compose(leftOptic, rightOptic)

	//The composed optic will receive a slice of tuples and focus on the first element of each tuple.
	//Under modification it will allow the first element to be edited
	//and the result is slice of tuples with a modified first element.

	//Note the result index is the Void index from the right optic
	//See ICompose for the ability to define a different index.

	//MustListOf retrieves the focuses and converts them to a slice
	var listOfResult []string = MustGet(SliceOf(composedOptic, len(data)), data)
	fmt.Println(listOfResult) //[alpha beta gamma delta]

	//Over applies the strings.ToUpper function to each focused element and combines
	//them into a new copy of the original data structure.
	//In the result the first elements are in upper case the second elements are unmodified.
	var overResult []lo.Tuple2[string, int] = MustModify(composedOptic, Op(strings.ToUpper), data)
	fmt.Println(overResult) //[{ALPHA 1} {BETA 2} {GAMMA 3} {DELTA 3}]

	//Output:
	//[alpha beta gamma delta]
	//[{ALPHA 1} {BETA 2} {GAMMA 3} {DELTA 4}]
}

func ExampleComposeLeft() {

	data := []lo.Tuple2[string, int]{
		lo.T2("alpha", 1),
		lo.T2("beta", 2),
		lo.T2("gamma", 3),
		lo.T2("delta", 4),
	}

	//In this example we will fully specify the types of the optics
	//This is is not necessary in normal usage.

	//The left optic we will use in the compose
	var leftOptic Optic[
		int,                      //Index Type
		[]lo.Tuple2[string, int], //Source Type
		[]lo.Tuple2[string, int], //Result Type
		lo.Tuple2[string, int],   //Focus type
		lo.Tuple2[string, int],   //Result focus type
		ReturnMany,               //Return type
		ReadWrite,                //ReadWrite type
		UniDir,                   //Dir type
		Pure,                     //Error type
	] = TraverseSlice[lo.Tuple2[string, int]]()

	//This optic will receive a slice of tuples and focus on the individual tuples.
	//Under modification it will allow the tuple to be edited and combine the edits
	//together into a slice of tuples.

	//The right optic we will use in the compose
	var rightOptic Optic[
		int,                    //Index Type
		lo.Tuple2[string, int], //Source Type
		lo.Tuple2[string, int], //Result Type
		string,                 //Focus type
		string,                 //Result focus type
		ReturnOne,              //Return type
		ReadWrite,              //ReadWrite type
		UniDir,                 //Dir type
		Pure,                   //Error type
	] = T2A[string, int]()

	//This optic will receive a tuple and focus on the first element.
	//Under modification it will allow the first element to be edited
	//tand the result is a tuple with a modified first element.

	var composedOptic Optic[
		int,                                    //Index Type
		[]lo.Tuple2[string, int],               //Source Type
		[]lo.Tuple2[string, int],               //Result Type
		string,                                 //Focus type
		string,                                 //Result focus type
		CompositionTree[ReturnMany, ReturnOne], //Return type
		CompositionTree[ReadWrite, ReadWrite],  //ReadWrite type
		CompositionTree[UniDir, UniDir],        //Dir type
		CompositionTree[Pure, Pure],            //Error type
	] = ComposeLeft(leftOptic, rightOptic)

	//The composed optic will receive a slice of tuples and focus on the first element of each tuple.
	//Under modification it will allow the first element to be edited
	//and the result is slice of tuples with a modified first element.

	//Note the result index is the int index from the left optic
	//See ICompose for the ability to define a different index.

	//MustIListOf retrieves the focuses and their indices and converts them to a slice
	//The returned indices are the index into the slice.
	var listOfResult []ValueI[int, string] = MustGet(SliceOf(WithIndex(composedOptic), 4), data)
	fmt.Println(listOfResult) //[0:alpha 1:beta 2:gamma 3:delta]

	//Over applies the custom ToUpper function to each focused element and combines
	//them into a new copy of the original data structure.
	//In the result the first elements with an even index are in upper case the second elements are unmodified.
	var overResult []lo.Tuple2[string, int] = MustModifyI(composedOptic, OpI(func(index int, focus string) string {
		if index%2 == 0 {
			return strings.ToUpper(focus)
		}
		return focus
	}), data)
	fmt.Println(overResult) //[{ALPHA 1} {beta 2} {GAMMA 3} {delta 3}]

	//Output:
	//[0:alpha 1:beta 2:gamma 3:delta]
	//[{ALPHA 1} {beta 2} {GAMMA 3} {delta 4}]
}

func ExampleComposeBoth() {

	data := map[string][]int{
		"alpha": []int{1, 2},
		"beta":  []int{3},
		"gamma": []int{4, 5, 6},
		"delta": []int{7, 8},
	}

	//In this example we will fully specify the types of the optics
	//This is is not necessary in normal usage.

	//The left optic we will use in the compose
	var leftOptic Optic[
		string,           //Index Type
		map[string][]int, //Source Type
		map[string][]int, //Result Type
		[]int,            //Focus type
		[]int,            //Result focus type
		ReturnMany,       //Return type
		ReadWrite,        //ReadWrite type
		UniDir,           //Dir type
		Pure,             //Error type
	] = TraverseMap[string, []int]()

	//This optic will recieve a map[string][]int and focus on the []int for each key
	//Under modification it will allow the []int to be edited and combine the edits
	//together into a map[string][]int.

	//The right optic we will use in the compose
	var rightOptic Optic[
		int,        //Index Type
		[]int,      //Source Type
		[]int,      //Result Type
		int,        //Focus type
		int,        //Result focus type
		ReturnMany, //Return type
		ReadWrite,  //ReadWrite type
		UniDir,     //Dir type
		Pure,       //Error type
	] = TraverseSlice[int]()

	//This optic will receive a []int] and focus on the individual ints.
	//Under modification it will allow the ints to be edited
	//and the result is a new []int with the new values.

	var composedOptic Optic[
		lo.Tuple2[string, int],                  //Index Type
		map[string][]int,                        //Source Type
		map[string][]int,                        //Result Type
		int,                                     //Focus type
		int,                                     //Result focus type
		CompositionTree[ReturnMany, ReturnMany], //Return type
		CompositionTree[ReadWrite, ReadWrite],   //ReadWrite type
		CompositionTree[UniDir, UniDir],         //Dir type
		CompositionTree[Pure, Pure],             //Error type
	] = ComposeBoth(leftOptic, rightOptic)

	//The composed optic will receive a map[string][]int and focus on each individual int.
	//Under modification it will allow the int to be edited
	//and the result is map[string][]int with the updated values.

	//Note the result index is a combination of the map string index and the slice int index.
	//See ICompose for the ability to define a different index.

	//MustIListOf retrieves the focuses and their indices and converts them to a slice
	//The returned indices are the index of the map and the slice combines
	var listOfResult []ValueI[lo.Tuple2[string, int], int] = MustGet(SliceOf(WithIndex(composedOptic), 8), data)
	fmt.Println(listOfResult) //[{alpha 0}:1 {alpha 1}:2 {beta 0}:3 {delta 0}:7 {delta 1}:8 {gamma 0}:4 {gamma 1}:5 {gamma 2}:6]

	//Over applies the custom function to each focused element and combines
	//them into a new copy of the original data structure.
	//In the result the int value is multiplied by 100 only if the map key has a length > 4 and the slice index is even.
	var overResult map[string][]int = MustModifyI(composedOptic, OpI(func(index lo.Tuple2[string, int], focus int) int {
		if len(index.A) > 4 && index.B%2 == 0 {
			return focus * 100
		}
		return focus
	}), data)
	fmt.Println(overResult) //map[alpha:[100 2] beta:[3] delta:[700 8] gamma:[400 5 600]]

	//Index uses the composed index to lookup the element at
	//data["gamma"][1]
	//Through the use of IComposeBoth both of these lookups are
	//performed efficiently.
	indexResult, found := MustGetFirst(Index(composedOptic, lo.T2("gamma", 1)), data)
	fmt.Println(indexResult, found) //5 true

	//Output:
	//[{alpha 0}:1 {alpha 1}:2 {beta 0}:3 {delta 0}:7 {delta 1}:8 {gamma 0}:4 {gamma 1}:5 {gamma 2}:6]
	//map[alpha:[100 2] beta:[3] delta:[700 8] gamma:[400 5 600]]
	//5 true
}

func ExampleComposeI() {

	//This example demonstrates the custom index mapping
	//for an example on composition in general check the Compose function
	data := [][]int{
		{1, 2, 3},
		{4, 5},
		{6, 7, 8},
	}

	//Combine the 2 indices into an array rather than the
	//[lo.Tuple2] returned by IComposeBoth
	ixMap := IxMap(func(left int, right int) [2]int {
		return [2]int{left, right}
	})

	//Note: the IxMap function does not support ReverseGet
	//to recover the left and right indexes this will
	//result in poor performance when using the Index() combinator.

	//The composed optic focuses on the individual ints in the nested slices.
	composedOptic := ComposeI(ixMap, TraverseSlice[[]int](), TraverseSlice[int]())

	//Note that the index in the result is [2]int the result of our ixMap function
	var listOfResult []ValueI[[2]int, int] = MustGet(SliceOf(WithIndex(composedOptic), 8), data)
	fmt.Println(listOfResult) //[[0 0]:1 [0 1]:2 [0 2]:3 [1 0]:4 [1 1]:5 [2 0]:6 [2 1]:7 [2 2]:8]

	//Over applies the custom function to each focused element and combines
	//them into a new copy of the original data structure.
	//In the result the int value is multiplied by 100 only if both slice indices are even.
	var overResult [][]int = MustModifyI(composedOptic, OpI(func(index [2]int, focus int) int {
		if index[0]%2 == 0 && index[1]%2 == 0 {
			return focus * 100
		}
		return focus
	}), data)
	fmt.Println(overResult) //[[100 2 300] [4 5] [600 7 800]]

	//Index uses the composed index to lookup the element at
	//data[1][1]
	//Note: even though the IxMap has an inefficient implementation
	//it still returns the correct results
	//For an efficient implementation use IComposeBoth or implement a custom IxMapper
	indexResult, found := MustGetFirst(Index(composedOptic, [2]int{1, 1}), data)
	fmt.Println(indexResult, found) //5 true

	//Output:
	//[[0 0]:1 [0 1]:2 [0 2]:3 [1 0]:4 [1 1]:5 [2 0]:6 [2 1]:7 [2 2]:8]
	//[[100 2 300] [4 5] [600 7 800]]
	//5 true
}

func ExampleIxMap() {

	//This example demonstrates the custom index mapping
	//for an example on composition in general check the Compose function
	data := [][]int{
		{1, 2, 3},
		{4, 5},
		{6, 7, 8},
	}

	//Combine the 2 indices into an array rather than the
	//[lo.Tuple2] returned by ComposeBoth
	//Even though The mapping is technically lossless IxMap does not support
	//recovering the original indexes and [Traversal] operations will be inefficient
	//See: [IxMapF] for an efficient version.
	ixMap := IxMap(
		//ixmap
		func(left int, right int) [2]int {
			return [2]int{left, right}
		},
	)

	//The composed optic focuses on the individual ints in the nested slices.
	composedOptic := ComposeI(ixMap, TraverseSlice[[]int](), TraverseSlice[int]())

	//Note that the index in the result is [2]int the result of our ixMap function
	var listOfResult []ValueI[[2]int, int] = MustGet(SliceOf(WithIndex(composedOptic), 8), data)
	fmt.Println(listOfResult) //[[0 0]:1 [0 1]:2 [0 2]:3 [1 0]:4 [1 1]:5 [2 0]:6 [2 1]:7 [2 2]:8]

	//Over applies the custom function to each focused element and combines
	//them into a new copy of the original data structure.
	//In the result the int value is multiplied by 100 only if both slice indices are even.
	var overResult [][]int = MustModifyI(composedOptic, OpI(func(index [2]int, focus int) int {
		if index[0]%2 == 0 && index[1]%2 == 0 {
			return focus * 100
		}
		return focus
	}), data)
	fmt.Println(overResult) //[[100 2 300] [4 5] [600 7 800]]

	//Index uses the composed index to lookup the element at
	//data[1][1]
	//Note: the Index lookup will be inefficient as we have not used IxMapF's unmap support.
	indexResult, found := MustGetFirst(Index(composedOptic, [2]int{1, 1}), data)
	fmt.Println(indexResult, found) //5 true

	//Output:
	//[[0 0]:1 [0 1]:2 [0 2]:3 [1 0]:4 [1 1]:5 [2 0]:6 [2 1]:7 [2 2]:8]
	//[[100 2 300] [4 5] [600 7 800]]
	//5 true

}

func ExampleIxMapIso() {

	//This example demonstrates the custom index mapping
	//for an example on composition in general check the Compose function
	data := [][]int{
		{1, 2, 3},
		{4, 5},
		{6, 7, 8},
	}

	//Combine the 2 indices into an array rather than the
	//[lo.Tuple2] returned by ComposeBoth
	//The mapping is lossless so we can recover the left and right components
	//This enables indexed operations to be performed in an efficient way.
	ixMap := IxMapIso[int, int, [2]int](
		//ixmap
		func(left int, right int) [2]int {
			return [2]int{left, right}
		},
		//ixmatch
		func(i1, i2 [2]int) bool {
			return i1[0] == i2[0] && i1[1] == i2[1]
		},
		//unmap
		func(mapped [2]int) (int, bool, int, bool) {
			return mapped[0], true, mapped[1], true
		},
		ExprCustom("ExampleIxMapF"),
	)

	//The composed optic focuses on the individual ints in the nested slices.
	composedOptic := ComposeI(ixMap, TraverseSlice[[]int](), TraverseSlice[int]())

	//Note that the index in the result is [2]int the result of our ixMap function
	var listOfResult []ValueI[[2]int, int] = MustGet(SliceOf(WithIndex(composedOptic), 8), data)
	fmt.Println(listOfResult) //[[0 0]:1 [0 1]:2 [0 2]:3 [1 0]:4 [1 1]:5 [2 0]:6 [2 1]:7 [2 2]:8]

	//Over applies the custom function to each focused element and combines
	//them into a new copy of the original data structure.
	//In the result the int value is multiplied by 100 only if both slice indices are even.
	var overResult [][]int = MustModifyI(composedOptic, OpI(func(index [2]int, focus int) int {
		if index[0]%2 == 0 && index[1]%2 == 0 {
			return focus * 100
		}
		return focus
	}), data)
	fmt.Println(overResult) //[[100 2 300] [4 5] [600 7 800]]

	//Index uses the composed index to lookup the element at
	//data[1][1]
	//Note: the given unmap function will be used to perform an efficient lookup.
	indexResult, found := MustGetFirst(Index(composedOptic, [2]int{1, 1}), data)
	fmt.Println(indexResult, found) //5 true

	//Output:
	//[[0 0]:1 [0 1]:2 [0 2]:3 [1 0]:4 [1 1]:5 [2 0]:6 [2 1]:7 [2 2]:8]
	//[[100 2 300] [4 5] [600 7 800]]
	//5 true
}

func ExampleIxMapLeft() {

	data := map[string]lo.Tuple2[int, string]{
		"a": lo.T2(10, "alpha"),
		"b": lo.T2(11, "beta"),
		"c": lo.T2(12, "gamma"),
	}

	optic := ComposeI(
		IxMapLeft[string, int](IxMatchComparable[string]()),
		TraverseMap[string, lo.Tuple2[int, string]](),
		T2B[int, string](),
	)

	res := MustGet(SeqIOf(optic), data)
	fmt.Println(res)

	//Output:
	//Seq[a:alpha b:beta c:gamma]

}

func ExampleIxMapRight() {

	data := map[string]lo.Tuple2[int, string]{
		"a": lo.T2(10, "alpha"),
		"b": lo.T2(11, "beta"),
		"c": lo.T2(12, "gamma"),
	}

	optic := ComposeI(
		IxMapRight[string, int](IxMatchComparable[int]()),
		TraverseMap[string, lo.Tuple2[int, string]](),
		T2B[int, string](),
	)

	res := MustGet(SeqIOf(optic), data)
	fmt.Println(res)

	//Output:
	//Seq[1:alpha 1:beta 1:gamma]

}

func ExampleIxMapBoth() {

	data := map[string]lo.Tuple2[int, string]{
		"a": lo.T2(10, "alpha"),
		"b": lo.T2(11, "beta"),
		"c": lo.T2(12, "gamma"),
	}

	optic := ComposeI(
		IxMapBoth[string, int](IxMatchComparable[string](), IxMatchComparable[int]()),
		TraverseMap[string, lo.Tuple2[int, string]](),
		T2B[int, string](),
	)

	res := MustGet(SeqIOf(optic), data)
	fmt.Println(res)

	//Output:
	//Seq[{a 1}:alpha {b 1}:beta {c 1}:gamma]

}

func TestReturnManyGetterComposition(t *testing.T) {

	if _, res, err := Compose(
		TraverseSlice[int](),
		Add(5),
	).AsGetter()(
		context.Background(),
		[]int{10, 20, 30},
	); res != 15 || err != nil {
		//The first value from the slice is taken and 5 is added this seems ok.
		t.Fatal("1", res, err)
	}

	if _, res, err := AsIxGet(
		ComposeLeft(
			TraverseSlice[int](),
			Add(5),
		),
	).AsGetter()(
		context.Background(),
		ValI(1, []int{10, 20, 30}),
	); res.MustGet() != 25 || err != nil {
		t.Fatal("1.a", res, err)
	}

	if _, res, err := Compose(
		WithIndex(Reversed(TraverseSlice[int]())),
		OpOnIx[int](Add(5)),
	).AsGetter()(
		context.Background(),
		[]int{10, 20, 30},
	); res != 7 || err != nil {
		//OpGet should return the first iteration result of the backward slice
		//Then we add five to it
		t.Fatal("2", res, err)
	}

	if _, res, err := AsIxGet(Compose(
		TraverseSlice[[]int](),
		TraverseSlice[int](),
	)).AsGetter()(
		context.Background(),
		ValI(1, [][]int{{10, 20}, {30, 40}}),
	); res.MustGet() != 20 || err != nil {
		//left opget returns the first slice
		//right opget uses the index we pass
		t.Fatal("3", res, err)
	}

	if _, res, err := AsIxGet(ComposeLeft(
		TraverseSlice[[]int](),
		TraverseSlice[int](),
	)).AsGetter()(
		context.Background(),
		ValI(1, [][]int{{10, 20}, {30, 40}}),
	); res.MustGet() != 30 || err != nil {
		//Left now uses the index we pass
		//Right returns the first element
		t.Fatal("4", res, err)
	}

	if _, res, err := AsIxGet(ComposeBoth(
		TraverseSlice[[]int](),
		TraverseSlice[int](),
	)).AsGetter()(
		context.Background(),
		ValI(lo.T2(1, 1), [][]int{{10, 20}, {30, 40}}),
	); res.MustGet() != 40 || err != nil {
		//ComposeBoth us access to both indexes
		t.Fatal("5", res, err)
	}

	if _, res, err := Compose(
		Add(5),
		Mul(2),
	).AsGetter()(
		context.Background(),
		10,
	); res != 30 || err != nil {
		//Both params are Void this works as expected
		t.Fatal("6", res, err)
	}

	if _, res, err := Compose(
		OpOnIx[int](Add(5)),
		Mul(2),
	).AsGetter()(
		context.Background(),
		ValI(10, 20),
	); res != 30 || err != nil {
		//OpIx[int] uses the index we pass instead the value
		t.Fatal("7", res, err)
	}

	if _, res, err := Compose(
		WithIndex(ReIndexed(OpOnIx[int](Add(5)), Op(func(index Void) int { return 1 }), EqT2[int]())),
		OpOnIx[int](Mul(2)),
	).AsGetter()(
		context.Background(),
		ValI(10, 20),
	); res != 2 || err != nil {
		//OpIx2 works better that OpIx for the left side as the source forces an Index to be passed via the source
		t.Fatal("8", res, err)
	}

	if _, res, err := Compose3(
		AsIxGet(TraverseSlice[int]()),
		Non(0, EqT2[int]()),
		Add(5),
	).AsGetter()(
		context.Background(),
		ValI(1, []int{10, 20, 30}),
	); res != 25 || err != nil {
		//OpGet should return the first iteration result of the backward slice
		//Then we add five to it
		t.Fatal("9", res, err)
	}
}

func TestIdentityOptimization(t *testing.T) {

	data := []string{"alpha", "beta"}
	newVal := "gamma"

	ValidateOpticTestPred(t, Compose(Identity[[]string](), TraverseSlice[string]()), data, newVal, EqDeepT2[[]string]())
	ValidateOpticTestPred(t, Compose(TraverseSlice[string](), Identity[string]()), data, newVal, EqDeepT2[[]string]())
	ValidateOpticTest(t, Compose(Identity[string](), Identity[string]()), "data", "newVal")

	ValidateOpticTestPred(t, ComposeLeft(Identity[[]string](), TraverseSlice[string]()), data, newVal, EqDeepT2[[]string]())
	ValidateOpticTestPred(t, ComposeLeft(TraverseSlice[string](), Identity[string]()), data, newVal, EqDeepT2[[]string]())
	ValidateOpticTest(t, ComposeLeft(Identity[string](), Identity[string]()), "data", "newVal")

}

func TestComposeOptimizationsRet1Ret1(t *testing.T) {

	left := T2A[string, string]()
	right := ParseInt[int](10, 0)
	data := lo.T2("1", "B")
	errData := lo.T2("A", "B")

	lRet1rRet1 := Compose(left, right)
	lRet1rRet1cLeft := ComposeLeft(left, right)
	lRet1rRet1cBoth := ComposeBoth(left, right)
	lRet1rRet1cNone := Indexing(ComposeBoth(left, right))
	lRet1rRetErr := Compose(left, Error[string, int](errors.New("sabotage")))

	//Get
	gres, err := Get(lRet1rRet1, data)
	if err != nil || gres != 1 {
		t.Fatal(gres, err)
	}

	gres, err = Get(lRet1rRet1, errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	TupleElement(0)
` {
		t.Fatal(gres, err)
	}

	gres, ok, err := GetFirst(lRet1rRet1, data)
	if !ok || err != nil || gres != 1 {
		t.Fatal(gres, ok, err)
	}

	gres, ok, err = GetFirst(lRet1rRet1, errData)
	if ok || err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	TupleElement(0)
` {
		t.Fatal(gres, ok, err)
	}

	//Set
	sres, err := Set(lRet1rRet1, 2, data)
	if err != nil || !reflect.DeepEqual(sres, lo.T2("2", "B")) {
		t.Fatal(sres, err)
	}

	sres, err = Set(lRet1rRetErr, 2, data)
	if ok || err == nil || err.Error() != `sabotage
optic error path:
	Error(sabotage)
	TupleElement(0)
` {
		t.Fatal(sres, ok, err)
	}

	//Iterate
	ires, err := Get(SliceOf(lRet1rRet1, 1), data)
	if err != nil || !reflect.DeepEqual(ires, []int{1}) {
		t.Fatal(ires, err)
	}

	ires, err = Get(SliceOf(lRet1rRet1, 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	TupleElement(0)
	SliceOf(TupleElement(0) | ParseInt(10,0))
` {
		t.Fatal(gres, ok, err)
	}

	//LengthGetter
	lres, err := Get(Length(lRet1rRet1), data)
	if err != nil || lres != 1 {
		t.Fatal(lres, err)
	}

	lres, err = Get(Length(lRet1rRetErr), errData)
	if err == nil || err.Error() != `sabotage
optic error path:
	Error(sabotage)
	TupleElement(0)
	Length(TupleElement(0) | Error(sabotage))
` {
		t.Fatal(lres, ok, err)
	}

	//Modify
	mres, err := Modify(lRet1rRet1, Add(1), data)
	if err != nil || !reflect.DeepEqual(mres, lo.T2("2", "B")) {
		t.Fatal(mres, err)
	}

	mres, err = Modify(lRet1rRet1, Add(1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	TupleElement(0)
` {
		t.Fatal(mres, err)
	}

	//IxGetter

	ixres, err := Get(SliceOf(Index(lRet1rRet1, Void{}), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(lRet1rRet1cLeft, 0), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(lRet1rRet1cBoth, lo.T2(0, Void{})), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(lRet1rRet1cNone, 0), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(lRet1rRet1, Void{}), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	TupleElement(0)
	Index({})
	SliceOf(Index({}))
` {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(lRet1rRet1cLeft, 0), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	TupleElement(0)
	Index(0)
	SliceOf(Index(0))
` {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(lRet1rRet1cBoth, lo.T2(0, Void{})), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	TupleElement(0)
	Index({0 {}})
	SliceOf(Index({0 {}}))
` {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(lRet1rRet1cNone, 0), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	TupleElement(0)
	Indexing(TupleElement(0) | ParseInt(10,0))
	Index(0)
	SliceOf(Index(0))
` {
		t.Fatal(ixres, err)
	}
}

func TestComposeOptimizationsRetMRet1(t *testing.T) {

	left := Filtered(TraverseSlice[string](), Ne("X"))
	right := ParseInt[int](10, 0)
	data := []string{"X", "1", "2"}
	errData := []string{"A", "2"}

	optic := Compose(left, right)
	opticLeft := ComposeLeft(left, right)
	opticBoth := ComposeBoth(left, right)
	opticNone := Indexing(ComposeBoth(left, right))
	opticErr := Compose(left, Error[string, int](errors.New("sabotage")))

	//Get

	gres, ok, err := GetFirst(optic, data)
	if !ok || err != nil || gres != 1 {
		t.Fatal(gres, ok, err)
	}

	gres, ok, err = GetFirst(optic, errData)
	if ok || err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
` {
		t.Fatal(gres, ok, err)
	}

	//Set
	sres, err := Set(optic, 3, data)
	if err != nil || !reflect.DeepEqual(sres, []string{"X", "3", "3"}) {
		t.Fatal(sres, err)
	}

	sres, err = Set(opticErr, 3, data)
	if ok || err == nil || err.Error() != `sabotage
optic error path:
	Error(sabotage)
	Traverse
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
` {
		t.Fatal(sres, ok, err)
	}

	//Iterate
	ires, err := Get(SliceOf(optic, 2), data)
	if err != nil || !reflect.DeepEqual(ires, []int{1, 2}) {
		t.Fatal(ires, err)
	}

	ires, err = Get(SliceOf(optic, 2), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	SliceOf(Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue) | ParseInt(10,0))
` {
		t.Fatal(gres, ok, err)
	}

	//LengthGetter
	lres, err := Get(Length(optic), data)
	if err != nil || lres != 2 {
		t.Fatal(lres, err)
	}

	lres, err = Get(Length(opticErr), errData)
	if err == nil || err.Error() != `sabotage
optic error path:
	Error(sabotage)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	Length(Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue) | Error(sabotage))
` {
		t.Fatal(lres, ok, err)
	}

	//Modify
	mres, err := Modify(optic, Add(1), data)
	if err != nil || !reflect.DeepEqual(mres, []string{"X", "2", "3"}) {
		t.Fatal(mres, err)
	}

	mres, err = Modify(optic, Add(1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Traverse
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
` {
		t.Fatal(mres, err)
	}

	//IxGetter

	ixres, err := Get(SliceOf(Index(optic, Void{}), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1, 2}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticLeft, 1), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticBoth, lo.T2(1, Void{})), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticNone, 0), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(optic, Void{}), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	Index({})
	SliceOf(Index({}))
` {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticLeft, 0), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	Index(0)
	SliceOf(Index(0))
` {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticBoth, lo.T2(0, Void{})), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	Index({0 {}})
	SliceOf(Index({0 {}}))
` {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticNone, 0), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	Indexing(Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue) | ParseInt(10,0))
	Index(0)
	SliceOf(Index(0))
` {
		t.Fatal(ixres, err)
	}
}

func TestComposeOptimizationsRet1RetM(t *testing.T) {

	left := T2A[[]string, string]()
	right := ComposeLeft(Filtered(TraverseSlice[string](), Ne("X")), ParseInt[int](10, 0))
	data := lo.T2([]string{"X", "1", "2"}, "B")
	errData := lo.T2([]string{"X", "A", "2"}, "B")

	optic := Compose(left, right)
	opticLeft := ComposeLeft(left, right)
	opticBoth := ComposeBoth(left, right)
	opticNone := Indexing(ComposeBoth(left, right))
	opticErr := Compose(left, Compose(Filtered(TraverseSlice[string](), Ne("X")), Error[string, int](errors.New("sabotage"))))

	//Get

	gres, ok, err := GetFirst(optic, data)
	if !ok || err != nil || gres != 1 {
		t.Fatal(gres, ok, err)
	}

	gres, ok, err = GetFirst(optic, errData)
	if ok || err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	TupleElement(0)
` {
		t.Fatal(gres, ok, err)
	}

	//Set
	sres, err := Set(optic, 3, data)
	if err != nil || !reflect.DeepEqual(sres, lo.T2([]string{"X", "3", "3"}, "B")) {
		t.Fatal(sres, err)
	}

	sres, err = Set(opticErr, 3, data)
	if ok || err == nil || err.Error() != `sabotage
optic error path:
	Error(sabotage)
	Traverse
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	TupleElement(0)
` {
		t.Fatal(sres, ok, err)
	}

	//Iterate
	ires, err := Get(SliceOf(optic, 2), data)
	if err != nil || !reflect.DeepEqual(ires, []int{1, 2}) {
		t.Fatal(ires, err)
	}

	ires, err = Get(SliceOf(optic, 2), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	TupleElement(0)
	SliceOf(TupleElement(0) | Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue) | ParseInt(10,0))
` {
		t.Fatal(gres, ok, err)
	}

	//LengthGetter
	lres, err := Get(Length(optic), data)
	if err != nil || lres != 2 {
		t.Fatal(lres, err)
	}

	lres, err = Get(Length(opticErr), errData)
	if err == nil || err.Error() != `sabotage
optic error path:
	Error(sabotage)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	TupleElement(0)
	Length(TupleElement(0) | Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue) | Error(sabotage))
` {
		t.Fatal(lres, ok, err)
	}

	//Modify
	mres, err := Modify(optic, Add(1), data)
	if err != nil || !reflect.DeepEqual(mres, lo.T2([]string{"X", "2", "3"}, "B")) {
		t.Fatal(mres, err)
	}

	mres, err = Modify(optic, Add(1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Traverse
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	TupleElement(0)
` {
		t.Fatal(mres, err)
	}

	//IxGetter

	ixres, err := Get(SliceOf(Index(optic, 1), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticLeft, 0), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1, 2}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticBoth, lo.T2(0, 1)), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticNone, 0), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(optic, 1), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	TupleElement(0)
	Index(1)
	SliceOf(Index(1))
` {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticLeft, 0), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	TupleElement(0)
	Index(0)
	SliceOf(Index(0))
` {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticBoth, lo.T2(0, 1)), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	TupleElement(0)
	Index({0 1})
	SliceOf(Index({0 1}))
` {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticNone, 0), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	TupleElement(0)
	Indexing(TupleElement(0) | Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue) | ParseInt(10,0))
	Index(0)
	SliceOf(Index(0))
` {
		t.Fatal(ixres, err)
	}
}

func TestComposeOptimizationsRetMRetM(t *testing.T) {

	left := TraverseSlice[[]string]()
	right := ComposeLeft(Filtered(TraverseSlice[string](), Ne("X")), ParseInt[int](10, 0))
	data := [][]string{
		{"X", "1", "2"},
		{"3"},
	}
	errData := [][]string{
		{"X", "A", "2"},
		{"3"},
	}

	optic := Compose(left, right)
	opticLeft := ComposeLeft(left, right)
	opticBoth := ComposeBoth(left, right)
	opticNone := Indexing(ComposeBoth(left, right))
	opticErr := Compose(left, Compose(Filtered(TraverseSlice[string](), Ne("X")), Error[string, int](errors.New("sabotage"))))

	//Get

	gres, ok, err := GetFirst(optic, data)
	if !ok || err != nil || gres != 1 {
		t.Fatal(gres, ok, err)
	}

	gres, ok, err = GetFirst(optic, errData)
	if ok || err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	Traverse
` {
		t.Fatal(gres, ok, err)
	}

	//Set
	sres, err := Set(optic, 3, data)
	if err != nil || !reflect.DeepEqual(sres, [][]string{{"X", "3", "3"}, {"3"}}) {
		t.Fatal(sres, err)
	}

	sres, err = Set(opticErr, 3, data)
	if ok || err == nil || err.Error() != `sabotage
optic error path:
	Error(sabotage)
	Traverse
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	Traverse
` {
		t.Fatal(sres, ok, err)
	}

	//Iterate
	ires, err := Get(SliceOf(optic, 3), data)
	if err != nil || !reflect.DeepEqual(ires, []int{1, 2, 3}) {
		t.Fatal(ires, err)
	}

	ires, err = Get(SliceOf(optic, 3), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	Traverse
	SliceOf(Traverse | Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue) | ParseInt(10,0))
` {
		t.Fatal(gres, ok, err)
	}

	//LengthGetter
	lres, err := Get(Length(optic), data)
	if err != nil || lres != 3 {
		t.Fatal(lres, err)
	}

	lres, err = Get(Length(opticErr), errData)
	if err == nil || err.Error() != `sabotage
optic error path:
	Error(sabotage)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	Traverse
	Length(Traverse | Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue) | Error(sabotage))
` {
		t.Fatal(lres, ok, err)
	}

	//Modify
	mres, err := Modify(optic, Add(1), data)
	if err != nil || !reflect.DeepEqual(mres, [][]string{{"X", "2", "3"}, {"4"}}) {
		t.Fatal(mres, err)
	}

	mres, err = Modify(optic, Add(1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Traverse
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	Traverse
` {
		t.Fatal(mres, err)
	}

	//IxGetter

	ixres, err := Get(SliceOf(Index(optic, 1), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticLeft, 0), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1, 2}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticBoth, lo.T2(0, 1)), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticNone, 0), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(optic, 1), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	Traverse
	Index(1)
	SliceOf(Index(1))
` {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticLeft, 0), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	Traverse
	Index(0)
	SliceOf(Index(0))
` {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticBoth, lo.T2(0, 1)), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	Traverse
	Index({0 1})
	SliceOf(Index({0 1}))
` {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticNone, 0), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue)
	Traverse
	Indexing(Traverse | Filtered(Const(true),ValueI[int,string].value | != X,FilterContinue,FilterContinue) | ParseInt(10,0))
	Index(0)
	SliceOf(Index(0))
` {
		t.Fatal(ixres, err)
	}
}

func TestComposeOptimizationsFilteredByRetMRetMFilteredBy(t *testing.T) {

	left := ComposeLeft(
		TraverseSlice[string](),
		Filtered(Identity[string](), Ne("X")),
	)
	right := ParseInt[int](10, 0)
	data := []string{"X", "1", "2"}
	errData := []string{"X", "A", "2"}

	optic := Compose(left, right)
	opticLeft := ComposeLeft(left, right)
	opticBoth := ComposeBoth(left, right)
	opticNone := Indexing(ComposeBoth(left, right))
	opticErr := Compose(left, Error[string, int](errors.New("sabotage")))

	//Get

	gres, ok, err := GetFirst(optic, data)
	if !ok || err != nil || gres != 1 {
		t.Fatal(gres, ok, err)
	}

	gres, ok, err = GetFirst(optic, errData)
	if ok || err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Traverse | Filtered(Const(true),ValueI[github.com/spearson78/go-optic.Void,string].value | != X,FilterContinue,FilterContinue)
` {
		t.Fatal(gres, ok, err)
	}

	//Set
	sres, err := Set(optic, 3, data)
	if err != nil || !reflect.DeepEqual(sres, []string{"X", "3", "3"}) {
		t.Fatal(sres, err)
	}

	sres, err = Set(opticErr, 3, data)
	if ok || err == nil || err.Error() != `sabotage
optic error path:
	Error(sabotage)
	Filtered(Const(true),ValueI[github.com/spearson78/go-optic.Void,string].value | != X,FilterContinue,FilterContinue)
	Traverse
` {
		t.Fatal(sres, ok, err)
	}

	//Iterate
	ires, err := Get(SliceOf(optic, 3), data)
	if err != nil || !reflect.DeepEqual(ires, []int{1, 2}) {
		t.Fatal(ires, err)
	}

	ires, err = Get(SliceOf(optic, 3), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Traverse | Filtered(Const(true),ValueI[github.com/spearson78/go-optic.Void,string].value | != X,FilterContinue,FilterContinue)
	SliceOf(Traverse | Filtered(Const(true),ValueI[github.com/spearson78/go-optic.Void,string].value | != X,FilterContinue,FilterContinue) | ParseInt(10,0))
` {
		t.Fatal(gres, ok, err)
	}

	//LengthGetter
	lres, err := Get(Length(optic), data)
	if err != nil || lres != 2 {
		t.Fatal(lres, err)
	}

	lres, err = Get(Length(opticErr), errData)
	if err == nil || err.Error() != `sabotage
optic error path:
	Error(sabotage)
	Traverse | Filtered(Const(true),ValueI[github.com/spearson78/go-optic.Void,string].value | != X,FilterContinue,FilterContinue)
	Length(Traverse | Filtered(Const(true),ValueI[github.com/spearson78/go-optic.Void,string].value | != X,FilterContinue,FilterContinue) | Error(sabotage))
` {
		t.Fatal(lres, ok, err)
	}

	//Modify
	mres, err := Modify(optic, Add(1), data)
	if err != nil || !reflect.DeepEqual(mres, []string{"X", "2", "3"}) {
		t.Fatal(mres, err)
	}

	mres, err = Modify(optic, Add(1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Filtered(Const(true),ValueI[github.com/spearson78/go-optic.Void,string].value | != X,FilterContinue,FilterContinue)
	Traverse
` {
		t.Fatal(mres, err)
	}

	//IxGetter

	ixres, err := Get(SliceOf(Index(optic, Void{}), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1, 2}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticLeft, 1), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticBoth, lo.T2(1, Void{})), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticNone, 0), 1), data)
	if err != nil || !reflect.DeepEqual(ixres, []int{1}) {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(optic, Void{}), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Traverse | Filtered(Const(true),ValueI[github.com/spearson78/go-optic.Void,string].value | != X,FilterContinue,FilterContinue)
	Index({})
	SliceOf(Index({}))
` {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticLeft, 1), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Traverse | Filtered(Const(true),ValueI[github.com/spearson78/go-optic.Void,string].value | != X,FilterContinue,FilterContinue)
	Index(1)
	SliceOf(Index(1))
` {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticBoth, lo.T2(1, Void{})), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Traverse | Filtered(Const(true),ValueI[github.com/spearson78/go-optic.Void,string].value | != X,FilterContinue,FilterContinue)
	Index({1 {}})
	SliceOf(Index({1 {}}))
` {
		t.Fatal(ixres, err)
	}

	ixres, err = Get(SliceOf(Index(opticNone, 0), 1), errData)
	if err == nil || err.Error() != `strconv.ParseInt: parsing "A": invalid syntax
optic error path:
	ParseInt(10,0)
	Traverse | Filtered(Const(true),ValueI[github.com/spearson78/go-optic.Void,string].value | != X,FilterContinue,FilterContinue)
	Indexing(Traverse | Filtered(Const(true),ValueI[github.com/spearson78/go-optic.Void,string].value | != X,FilterContinue,FilterContinue) | ParseInt(10,0))
	Index(0)
	SliceOf(Index(0))
` {
		t.Fatal(ixres, err)
	}
}
