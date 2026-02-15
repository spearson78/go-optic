package optic_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"unicode"

	"github.com/samber/lo"
	"github.com/samber/mo"

	. "github.com/spearson78/go-optic"
)

func TestCombinatorConsistency(t *testing.T) {

	o := TraverseSlice[int]()

	ValidateOpticTestPred(t, Length(o), []int{1, 2, 3, 4, 5}, 0, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, Taking(o, 2), []int{1, 2, 3, 4, 5}, 10, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, Taking(o, 0), []int{1, 2, 3, 4, 5}, 10, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, Dropping(o, 2), []int{1, 2, 3, 4, 5}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, Dropping(o, 2), []int{1, 2, 3, 4, 5}, 10, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, Reversed(o), []int{1, 2, 3, 4, 5}, 10, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, TakingWhile(o, Lte(10)), []int{1, 2, 30, 40, 50}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, TakingWhile(o, False[int]()), []int{1, 2, 30, 40, 50}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, TakingWhile(o, True[int]()), []int{1, 2, 30, 40, 50}, 10, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, DroppingWhile(o, Lt(10)), []int{1, 2, 30, 40, 50}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, DroppingWhile(o, False[int]()), []int{1, 2, 30, 40, 50}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, DroppingWhile(o, True[int]()), []int{1, 2, 30, 40, 50}, 10, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, TrimmingWhile(o, Lt(10)), []int{1, 2, 30, 2, 1}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, TrimmingWhile(o, False[int]()), []int{1, 2, 30, 2, 1}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, TrimmingWhile(o, True[int]()), []int{1, 2, 30, 2, 1}, 10, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, Filtered(o, Lte(10)), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, Filtered(o, False[int]()), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, Filtered(o, True[int]()), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, Element(o, 2), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, Element(o, 100), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, SliceOf(o, 10), []int{1, 2, 30, 40, 5}, []int{10, 20, 3, 4, 50}, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, SliceOfP(o, 10), []int{1, 2, 30, 40, 5}, []int{10, 20, 3, 4, 50}, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, MapOf(o, 5), []int{1, 2, 30, 40, 5}, map[int]int{0: 10, 1: 20, 2: 3, 3: 4, 4: 50}, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, MapOfP(o, 5), []int{1, 2, 30, 40, 5}, map[int]int{0: 10, 1: 20, 2: 3, 3: 4, 4: 50}, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, AtMap[string, int]("alpha"), map[string]int{"alpha": 1, "beta": 2}, mo.Some(10), EqDeepT2[map[string]int]())
	ValidateOpticTestPred(t, AtMap[string, int]("alpha"), map[string]int{"alpha": 1, "beta": 2}, mo.None[int](), EqDeepT2[map[string]int]())
	ValidateOpticTestPred(t, AtMap[string, int]("gamma"), map[string]int{"alpha": 1, "beta": 2}, mo.Some(10), EqDeepT2[map[string]int]())
	ValidateOpticTestPred(t, AtMap[string, int]("gamma"), map[string]int{"alpha": 1, "beta": 2}, mo.None[int](), EqDeepT2[map[string]int]())

	ValidateOpticTestPred(t, AtSlice[int](1, EqT2[int]()), []int{1, 2, 30, 40, 5}, mo.Some(10), EqDeepT2[[]int]())
	ValidateOpticTestPred(t, AtSlice[int](1, EqT2[int]()), []int{1, 2, 30, 40, 5}, mo.None[int](), EqDeepT2[[]int]())
	ValidateOpticTestPred(t, AtSlice[int](100, EqT2[int]()), []int{1, 2, 30, 40, 5}, mo.Some(10), EqDeepT2[[]int]())
	ValidateOpticTestPred(t, AtSlice[int](100, EqT2[int]()), []int{1, 2, 30, 40, 5}, mo.None[int](), EqDeepT2[[]int]())

	ValidateOpticTestPred(t, CoalesceN(Element(o, 10), o), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, CoalesceN(Element(o, 2), o), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, Coalesce(Element(o, 10), o), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, Coalesce(Element(o, 2), o), []int{1, 2, 30, 40, 5}, 10, EqDeepT2[[]int]())

	ValidateOpticTest(t, AsReverseGet(Add(10)), 10, 20)

	ValidateOpticTestPred(t, Ignore(o, Const[error](true)), []int{1, 2, 3}, 10, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, Reduce(o, Sum[int]()), []int{1, 2, 3}, 0, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, MapReduce(o, Mul(10), Sum[int]()), []int{1, 2, 3}, 0, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, First(o), []int{1, 20, 3}, 40, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, FirstOrDefault(o, 10), []int{1, 20, 3}, 0, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, FirstOrDefault(o, 10), []int{}, 0, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, Last(o), []int{1, 20, 3}, 40, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, MaxOf(o, Identity[int]()), []int{1, 20, 3}, 40, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, MaxOf(o, Identity[int]()), []int{}, 1, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, MinOf(o, Identity[int]()), []int{10, 20, 3}, 2, EqDeepT2[[]int]())
	ValidateOpticTestPred(t, MinOf(o, Identity[int]()), []int{}, 2, EqDeepT2[[]int]())

	ValidateOpticTestPred(t, Ordered(o, OrderBy(Identity[int]())), []int{5, 9, 3, 4, 2, 1, 8, 10, 6, 7}, 10, EqDeepT2[[]int]())

	ValidateOpticTest(t, Lookup(o, []int{5, 9, 3, 4, 2, 1, 8, 10, 6, 7}), 5, mo.Some(10))
	ValidateOpticTest(t, Lookup(o, []int{5, 9, 3, 4, 2, 1, 8, 10, 6, 7}), 100, mo.Some(10))
}

func TestIx(t *testing.T) {

	if r, ok, err := GetFirst(Index(TraverseSlice[string](), 1), []string{"Borg", "Cardassian", "Talaxian"}); !ok || err != nil || r != "Cardassian" {
		t.Fatal("IX 1")
	}
}

func ExampleLength() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	strLenGt4 := Compose(Length(TraverseString()), Gt(4))

	filteredResult := MustGet(SliceOf(Filtered(TraverseSlice[string](), strLenGt4), len(data)), data)
	fmt.Println(filteredResult)

	overResult := MustModify(Filtered(TraverseSlice[string](), strLenGt4), Op(strings.ToUpper), data)
	fmt.Println(overResult)

	//Output:
	//[alpha gamma delta]
	//[ALPHA beta GAMMA DELTA]
}

func ExampleTaking() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	takingOptic := Taking(TraverseSlice[string](), 2)

	listResult := MustGet(SliceOf(takingOptic, len(data)), data)
	fmt.Println(listResult)

	overResult := MustModify(takingOptic, Op(strings.ToUpper), data)
	fmt.Println(overResult)

	//Output:
	//[alpha beta]
	//[ALPHA BETA gamma delta]
}

func ExampleDropping() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	droppingOptic := Dropping(TraverseSlice[string](), 2)

	listResult := MustGet(SliceOf(droppingOptic, len(data)), data)
	fmt.Println(listResult)

	overResult := MustModify(droppingOptic, Op(strings.ToUpper), data)
	fmt.Println(overResult)

	//Output:
	//[gamma delta]
	//[alpha beta GAMMA DELTA]
}

func ExampleReversed() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	//taking 2 reversed is equivalent to taking the last 2 elements in the optic
	reversedOptic := Taking(Reversed(TraverseSlice[string]()), 2)

	listResult := MustGet(SliceOf(reversedOptic, len(data)), data)
	fmt.Println(listResult)

	//Note that the result of over is still in the original collection.
	//However the last 2 elements were converted to upper case as they were focused by the optic.
	overResult := MustModify(reversedOptic, Op(strings.ToUpper), data)
	fmt.Println(overResult)

	//Output:
	//[delta gamma]
	//[alpha beta GAMMA DELTA]
}

func ExampleTakingWhile() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	//focuses on the elements of the slice until the gamma string is found
	optic := TakingWhile(TraverseSlice[string](), Ne("gamma"))

	listResult := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(listResult)

	overResult := MustModify(optic, Op(strings.ToUpper), data)
	fmt.Println(overResult)

	//Output:
	//[alpha beta]
	//[ALPHA BETA gamma delta]
}

func ExampleDroppingWhile() {

	data := []int{1, 2, 30, 4, 5}

	//skips elements until a value >=10 is found.
	optic := DroppingWhile(TraverseSlice[int](), Lt(10))

	listResult := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(listResult)

	overResult := MustModify(optic, Mul(10), data)
	fmt.Println(overResult)

	//Output:
	//[30 4 5]
	//[1 2 300 40 50]
}

func ExampleTrimmingWhile() {

	data := "ooooooooleading and trailing o'sooooooooo"

	optic := TrimmingWhile(TraverseString(), Eq('o'))

	listResult := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(string(listResult))

	overResult := MustModify(optic, Op(unicode.ToUpper), data)
	fmt.Println(overResult)

	//Output:
	//leading and trailing o's
	//ooooooooLEADING AND TRAILING O'Sooooooooo

}

func ExampleFiltered() {

	data := []int{1, 2, 30, 4, 5}

	optic := Filtered(TraverseSlice[int](), Lt(10))

	listResult := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(listResult)

	overResult := MustModify(optic, Add(1), data)
	fmt.Println(overResult)

	//Output:
	//[1 2 4 5]
	//[2 3 30 5 6]
}

func ExampleElement() {

	//       0123456789
	data := "Lorem ipsum"

	optic := Element(TraverseString(), 6)

	viewResult, found := MustGetFirst(optic, data)
	fmt.Println(string(viewResult), found)

	modifyResult := MustModify(optic, Op(unicode.ToUpper), data)
	fmt.Println(modifyResult)

	//Output:
	//i true
	//Lorem Ipsum
}

func ExampleCoalesceN() {

	data := map[int]string{
		1: "alpha",
		2: "beta",
		3: "gamma",
		4: "delta",
	}

	//Polymorphic map traversal with a source type of map[int]string and result type of map[int]int
	index5 := Index(TraverseMap[int, string](), 5)
	index3 := Index(TraverseMap[int, string](), 3)

	failingOptic := CoalesceN(index5, index3)

	//There is no element 5 in the map so the index3 optic is used instead.
	viewResult, found := MustGetFirst(failingOptic, data)
	fmt.Println(viewResult, found) //gamma

	var modifyResult map[int]string = MustModify(failingOptic, Op(strings.ToUpper), data)
	fmt.Println(modifyResult)

	//Once we add an element at index 5 then the index5 will match
	data[5] = "epsilon"

	viewResult, found = MustGetFirst(failingOptic, data)
	fmt.Println(viewResult, found) //epsilon

	modifyResult = MustModify(failingOptic, Op(strings.ToUpper), data)
	fmt.Println(modifyResult)

	//Output:
	//gamma true
	//map[1:alpha 2:beta 3:GAMMA 4:delta]
	//epsilon true
	//map[1:alpha 2:beta 3:gamma 4:delta 5:EPSILON]
}

func ExampleAsReverseGet() {

	celsiusToFahrenheit := Compose(Mul(1.8), Add(32.0))

	var fahrenheit float64 = MustGet(celsiusToFahrenheit, 32)
	fmt.Println(fahrenheit)

	fahrenheitToCelsius := AsReverseGet(celsiusToFahrenheit)

	var celsius float64 = MustGet(fahrenheitToCelsius, 90)
	fmt.Println(celsius)

	//Note: this result is in fahrenheit but the 10 is added in celsius
	var add10CelsiusResult float64 = MustModify(fahrenheitToCelsius, Add(10.0), 90)
	fmt.Println(add10CelsiusResult)

	//Output:
	//89.6
	//32.22222222222222
	//108
}

func ExampleEmbed() {

	//DownCast returns a Prism (ReturnMany) as the downcast may focus nothing if the cast fails.
	anyToString := DownCast[any, string]()

	//Embed reverses the direction of the downCast Prism. The cast from string to any always suceeds
	//So the returned optic is now a Getter (ReturnOne,ReadOnly)
	stringToAny := Embed(anyToString)

	var res any = MustGet(stringToAny, "alpha")
	fmt.Println(res)

	//Output:
	//alpha
}

func TestIgnoreCancel(t *testing.T) {

	data := []string{"1", "2", "three", "4"}

	optic := Compose(TraverseSlice[string](), Ignore(ParseInt[int](10, 32), Const[error](true)))

	ctx, cancel := context.WithCancel(context.Background())

	_, err := ModifyContext(ctx, optic, Op(func(focus int) int {
		cancel()
		return focus
	}), data)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("context.Canceled : %v", err)
	}

}

func TestIgnoreOpticMethods(t *testing.T) {

	sabotage := IsoE[int, int](
		func(ctx context.Context, source int) (int, error) {
			return 0, errors.New("sabotage")
		},
		func(ctx context.Context, focus int) (int, error) {
			return 0, errors.New("sabotage")
		},
		ExprCustom("sabotage"),
	)

	goodOptic := Ignore(ParseInt[int](10, 32), Const[error](true))
	badOptic := Ignore(ComposeLeft(ParseInt[int](10, 32), sabotage), Const[error](true))

	//Fails to compile due to the ReturnMany of the IgnorePrism
	//if focus, err := View(optic, "1"); focus != 1 || err != nil {
	//	t.Fatalf("Ignore View", focus, err)
	//}

	//Setter
	if result, err := Set(goodOptic, 1, "2"); result != "1" || err != nil {
		t.Fatal("Ignore Set Good", result, err)
	}

	if result, err := Set(goodOptic, 1, "bad"); result != "1" || err != nil {
		t.Fatal("Ignore Set Bad Data", result, err)
	}

	if result, err := Set(badOptic, 1, "2"); result != "2" || err != nil {
		t.Fatal("Ignore Set Bad", result, err)
	}

	if result, err := Set(badOptic, 1, "bad"); result != "bad" || err != nil {
		t.Fatal("Ignore Set Bad Both", result, err)
	}

	//Iterate
	if result, err := Get(SliceOf(goodOptic, 1), "2"); !reflect.DeepEqual(result, []int{2}) || err != nil {
		t.Fatal("Ignore Iterate Good", result, err)
	}

	if result, err := Get(SliceOf(goodOptic, 1), "bad"); !reflect.DeepEqual(result, []int{}) || err != nil {
		t.Fatal("Ignore Iterate Bad Data", result, err)
	}

	if result, err := Get(SliceOf(badOptic, 1), "2"); !reflect.DeepEqual(result, []int{}) || err != nil {
		t.Fatal("Ignore Iterate Bad", result, err)
	}

	if result, err := Get(SliceOf(badOptic, 1), "bad"); !reflect.DeepEqual(result, []int{}) || err != nil {
		t.Fatal("Ignore Iterate Bad Both", result, err)
	}

	//LengthGetter

	if result, err := Get(Length(goodOptic), "2"); result != 1 || err != nil {
		t.Fatal("Ignore Length Good", result, err)
	}

	if result, err := Get(Length(goodOptic), "bad"); result != 0 || err != nil {
		t.Fatal("Ignore Length Bad Data", result, err)
	}

	if result, err := Get(Length(badOptic), "2"); result != 0 || err != nil {
		t.Fatal("Ignore Length Bad", result, err)
	}

	if result, err := Get(Length(badOptic), "bad"); result != 0 || err != nil {
		t.Fatal("Ignore Length Bad Both", result, err)
	}

	//Modify

	if result, err := Modify(goodOptic, Mul(10), "2"); result != "20" || err != nil {
		t.Fatal("Ignore Over Good", result, err)
	}

	if result, err := Modify(goodOptic, Mul(10), "bad"); result != "bad" || err != nil {
		t.Fatal("Ignore Over Bad Data", result, err)
	}

	if result, err := Modify(badOptic, Mul(10), "2"); result != "2" || err != nil {
		t.Fatal("Ignore Over Bad", result, err)
	}

	if result, err := Modify(badOptic, Mul(10), "bad"); result != "bad" || err != nil {
		t.Fatal("Ignore Over Bad", result, err)
	}

	//IxGetter

	if result, ok, err := GetFirst(Index(goodOptic, Void{}), "2"); !ok || result != 2 || err != nil {
		t.Fatal("Ignore IxGetter Good", result, err)
	}

	if result, ok, err := GetFirst(Index(goodOptic, Void{}), "bad"); ok || err != nil {
		t.Fatal("Ignore IxGetter Bad Data", result, ok, err)
	}

	if result, ok, err := GetFirst(Index(badOptic, Void{}), "2"); ok || err != nil {
		t.Fatal("Ignore IxGetter Bad", result, ok, err)
	}

	if result, ok, err := GetFirst(Index(badOptic, Void{}), "bad"); ok || err != nil {
		t.Fatal("Ignore IxGetter Bad Both", result, ok, err)
	}

	//ReverseGet errors cannot be ignored.

	if result, err := ReverseGet(goodOptic, 2); result != "2" || err != nil {
		t.Fatal("Ignore ReverseGetter Good", result, err)
	}

	if result, err := ReverseGet(badOptic, 2); err == nil {
		t.Fatal("Ignore ReverseGetter Bad", result, err)
	}

}

func ExampleIgnore() {

	data := []string{"1", "2", "three", "4"}

	optic := Compose(TraverseSlice[string](), Ignore(ParseInt[int](10, 32), Const[error](true)))

	var viewResult []int = MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(viewResult)

	var modifyResult []string = MustModify(optic, Mul(10), data)
	fmt.Println(modifyResult)

	//Output:
	//[1 2 4]
	//[10 20 three 40]
}

func ExampleStop() {

	optic := Stop(
		Compose(
			TraverseSlice[string](),
			ParseInt[int32](10, 32),
		),
		ErrorIs(strconv.ErrSyntax),
	)

	viewResult, err := Get(SliceOf(optic, 3), []string{"1", "two", "3"})
	fmt.Println(viewResult, err)

	viewResult, err = Get(SliceOf(optic, 3), []string{"1", "2147483648"})
	fmt.Println(viewResult, err)

	//Output:
	//[1] <nil>
	//[] strconv.ParseInt: parsing "2147483648": value out of range
	//optic error path:
	//	ParseInt(10,32)
	//	Traverse
	//	Stop(Traverse | ParseInt(10,32),ErrorIs(invalid syntax))
	//	SliceOf(Stop(Traverse | ParseInt(10,32),ErrorIs(invalid syntax)))
}

func ExampleCatch() {

	optic := Catch(
		Compose(
			TraverseSlice[string](),
			ParseInt[int32](10, 32),
		),
		If(
			ErrorIs(strconv.ErrSyntax),
			EErr(Const[error](int32(-1))),
			Ro(Throw[int32]()),
		),
	)

	viewResult, err := Get(SliceOf(optic, 3), []string{"1", "two", "3"})
	fmt.Println(viewResult, err)

	viewResult, err = Get(SliceOf(optic, 3), []string{"1", "2147483648", "3"})
	fmt.Println(viewResult, err)

	//Output:
	//[1 -1 3] <nil>
	//[] strconv.ParseInt: parsing "2147483648": value out of range
	//optic error path:
	//	ParseInt(10,32)
	//	Traverse
	//	Throw()
	//	Switch(Case(ErrorIs(invalid syntax) -> Const(-1))Default(Throw()))
	//	Catch(Traverse | ParseInt(10,32),Switch(Case(ErrorIs(invalid syntax) -> Const(-1))Default(Throw())),Throw())
	//	SliceOf(Catch(Traverse | ParseInt(10,32),Switch(Case(ErrorIs(invalid syntax) -> Const(-1))Default(Throw())),Throw()))
}

func ExampleCatchP() {

	optic := CatchP(
		Compose(
			TraverseSliceP[string, int32](),
			ParseIntP[int32](10, 32),
		),
		Ret1(Ro(EErr(If(
			ErrorIs(strconv.ErrSyntax),
			EErr(Const[error](int32(-1))),
			Ro(Throw[int32]()),
		)))),
		EErr(Const[error, []int32](nil)),
	)

	viewResult, err := Get(SliceOf(optic, 3), []string{"1", "two", "3"})
	fmt.Println(viewResult, err)

	modifyResult, err := Modify(optic, Mul[int32](2), []string{"1", "two", "3"})
	fmt.Println(modifyResult, err)

	//Output:
	//[1 -1 3] <nil>
	//[] <nil>
}

func ExampleReduce() {
	data := []int{1, 2, 3, 4}

	sum, ok := MustGetFirst(Reduce(TraverseSlice[int](), Sum[int]()), data)

	fmt.Println(sum, ok)
	//Output:10 true
}

func ExampleMapReduce() {
	data := []string{"1", "2", "3", "4"}

	sum, ok, err := GetFirst(MapReduce(TraverseSlice[string](), ParseInt[int](10, 32), Sum[int]()), data)

	fmt.Println(sum, ok, err)
	//Output:10 true <nil>
}

func ExampleFirst() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	optic := TraverseSlice[string]()

	var result string
	var found bool
	var err error
	result, found, err = GetFirst(First(optic), data)
	fmt.Println(result, found, err)

	var overResult []string = MustModify(First(optic), Op(strings.ToUpper), data)
	fmt.Println(overResult)

	//Output:
	//alpha true <nil>
	//[ALPHA beta gamma delta]
}

func ExampleFirstOrDefault() {

	optic := TraverseSlice[string]()

	var result string = MustGet(FirstOrDefault(optic, "default"), []string{})
	fmt.Println(result)

	//Output:
	//default
}

func ExampleLast() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	optic := TraverseSlice[string]()

	var result string
	var found bool
	var err error
	result, found, err = GetFirst(Last(optic), data)
	fmt.Println(result, found, err)

	var overResult []string = MustModify(Last(optic), Op(strings.ToUpper), data)
	fmt.Println(overResult)

	//Output:
	//delta true <nil>
	//[alpha beta gamma DELTA]
}

func ExampleMaxOf() {
	data := []lo.Tuple2[string, int]{
		lo.T2("alpha", 1),
		lo.T2("beta", 2),
		lo.T2("gamma", 30),
		lo.T2("delta", 4),
	}

	optic := TraverseSlice[lo.Tuple2[string, int]]()
	cmp := T2B[string, int]()

	var index int
	var maxTuple lo.Tuple2[string, int]
	var ok bool
	index, maxTuple, ok = MustGetFirstI(MaxOf(optic, cmp), data)

	fmt.Println(index, maxTuple, ok)

	var overResult []lo.Tuple2[string, int] = MustModify(Compose(MaxOf(optic, cmp), T2B[string, int]()), Mul(10), data)
	fmt.Println(overResult)

	//Output:
	//2 {gamma 30} true
	//[{alpha 1} {beta 2} {gamma 300} {delta 4}]
}

func ExampleMinOf() {
	data := []lo.Tuple2[string, int]{
		lo.T2("alpha", 10),
		lo.T2("beta", 20),
		lo.T2("gamma", 3),
		lo.T2("delta", 40),
	}

	optic := TraverseSlice[lo.Tuple2[string, int]]()
	cmp := T2B[string, int]()

	var index int
	var minTuple lo.Tuple2[string, int]
	var ok bool
	index, minTuple, ok = MustGetFirstI(MinOf(optic, cmp), data)

	fmt.Println(index, minTuple, ok)

	var overResult []lo.Tuple2[string, int] = MustModify(Compose(MinOf(optic, cmp), T2B[string, int]()), Mul(10), data)
	fmt.Println(overResult)

	//Output:2 {gamma 3} true
	//[{alpha 10} {beta 20} {gamma 30} {delta 40}]
}

func ExampleOrdered() {

	data := []int{9, 8, 3, 7, 2, 5, 1, 4, 10, 6}

	optic := Taking(
		Ordered(
			TraverseSlice[int](),
			OrderBy(Identity[int]()),
		),
		5,
	)

	var getResult []int = MustGet(SliceOf(optic, 5), data)
	fmt.Println(getResult)

	var modifyResult []int = MustModify(optic, Mul(100), data)
	fmt.Println(modifyResult)

	//Output:
	//[1 2 3 4 5]
	//[9 8 300 7 200 500 100 400 10 6]
}

func ExampleLookup() {

	vowels := map[rune]bool{
		'a': true,
		'e': true,
		'i': true,
		'o': true,
		'u': true,
	}

	inVowels := Compose(Lookup(TraverseMap[rune, bool](), vowels), Non(false, EqT2[bool]()))

	index, vowel, found := MustGetFirstI(Filtered(TraverseString(), inVowels), "hello world")
	fmt.Println(index, string(vowel), found)
	//Output:
	//1 e true
}

func ExampleMatching() {

	var result mo.Either[any, string] = MustGet(
		Matching(
			DownCast[any, string](),
		),
		"alpha",
	)
	fmt.Println(result.MustRight()) //The cast matches so we get the cast value on the right

	result = MustGet(
		Matching(
			DownCast[any, string](),
		),
		1.2,
	)
	fmt.Println(result.MustLeft()) //The cast does not match so we get the original value on the left

	//Output:
	//alpha
	//1.2
}

func ExampleForEach() {

	lookupData := map[string]string{
		"alpha": "1",
		"beta":  "2",
		"gamma": "3",
	}

	lookupNames := Compose(Lookup(TraverseMap[string, string](), lookupData), Some[string]())

	optic := ForEach(
		EErr(lookupNames),
		ParseInt[int](10, 32),
	)

	var res Collection[Void, string, Err]
	res, err := Modify(optic, Mul(10), "alpha")
	fmt.Println(res, err)
	//Output:
	//Col[{}:10] <nil>

}

func ExampleGrouped() {

	data := []int{1, 2, 3, 4, 5, 6, 4, 3, 2, 1, 4}

	res, err := Get(
		SliceOf(
			Grouped(
				SelfIndex(TraverseSlice[int](), EqT2[int]()),
				FirstReducer[int](),
			),
			len(data),
		),
		data,
	)

	fmt.Println(res, err)

	//Output:
	//[1 2 3 4 5 6] <nil>

}

func ExampleFirstOrError() {

	data := map[string]int{
		"alpha": 0,
		"gamma": 2,
		"delta": 3,
	}

	optic := FirstOrError(
		Index(
			TraverseMap[string, int](),
			"beta",
		),
		errors.New("key not found"),
	)

	res, err := Get(optic, data)
	fmt.Println(res, err)

	//Output:
	//0 key not found
	//optic error path:
	//	Error(key not found)
	//	Coalesce(Index(beta),Error(key not found))
	//	Filtered(< 1,Const(true),FilterContinue,FilterStop)

}

func ExamplePolymorphic() {

	data := []ValueI[int, string]{
		ValI(0, "1"),
		ValI(1, "2"),
		ValI(3, "3"),
	}

	optic := Compose3(
		TraverseSliceP[ValueI[int, string], ValueI[int, int]](),
		Polymorphic[ValueI[int, int], int](ValueIValue[int, string]()),
		ParseIntP[int](10, 0),
	)

	res, err := Get(SliceOf(optic, 3), data)
	fmt.Println(res, err)

	//Output:
	//[1 2 3] <nil>

}

func ExampleEditDistance() {

	data := []lo.Tuple2[string, string]{
		lo.T2("alpha", "alpha"),
		lo.T2("alpha", "Alpha"),
		lo.T2("alpha", "lapha"),
	}

	optic := Compose(
		TraverseSlice[lo.Tuple2[string, string]](),
		EditDistance(TraverseString(), EditOSA, EqT2[rune](), 10),
	)

	res := MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(res)

	//Output:
	//[0 1 1]

}

func ExampleConcat() {

	data := map[string]int{
		"alpha": 1,
		"beta":  2,
	}

	data2 := map[string]int{
		"gamma": 3,
		"delta": 4,
	}

	res := MustGet(
		MapOfReduced(
			Concat(
				Compose(
					T2A[map[string]int, map[string]int](),
					TraverseMap[string, int](),
				),
				Compose(
					T2B[map[string]int, map[string]int](),
					TraverseMap[string, int](),
				),
			),
			FirstReducer[int](),
			len(data)+len(data2),
		),
		lo.T2(data, data2),
	)

	fmt.Println(res)

	//Output:
	//map[alpha:1 beta:2 delta:4 gamma:3]

}

func TestOrdered(t *testing.T) {

	data := []int{278, 9, 3, 98, 3, 4, 690, 72, 8, 9761, 20, 93, 5, 9, 67, 28, 73, 479, 2, 5, 89}

	heapSort := Ordered[int](TraverseSlice[int](), OrderBy(Identity[int]()))

	if lowest, ok := MustGetFirst(heapSort, data); !ok || lowest != 2 {
		t.Fatal("MustPreView", ok, lowest)
	}

	if sorted := MustGet(SliceOf(heapSort, len(data)), data); !reflect.DeepEqual(sorted, []int{2, 3, 3, 4, 5, 5, 8, 9, 9, 20, 28, 67, 72, 73, 89, 93, 98, 278, 479, 690, 9761}) {
		t.Fatal("MustToSliceOf", sorted)
	}

	if sorted := MustGet(SliceOf(Taking(heapSort, 1), len(data)), data); !reflect.DeepEqual(sorted, []int{2}) {
		t.Fatal("Taking(1)", sorted)
	}

	if sorted := MustGet(SliceOf(Taking(heapSort, 2), len(data)), data); !reflect.DeepEqual(sorted, []int{2, 3}) {
		t.Fatal("Taking(2)", sorted)
	}

	if sorted := MustGet(SliceOf(Taking(heapSort, 3), len(data)), data); !reflect.DeepEqual(sorted, []int{2, 3, 3}) {
		t.Fatal("Taking(3)", sorted)
	}

	if sorted := MustGet(SliceOf(Taking(heapSort, 10), len(data)), data); !reflect.DeepEqual(sorted, []int{2, 3, 3, 4, 5, 5, 8, 9, 9, 20}) {
		t.Fatal("Taking(10)", sorted)
	}

	if sorted := MustGet(SliceOf(Taking(heapSort, 18), len(data)), data); !reflect.DeepEqual(sorted, []int{2, 3, 3, 4, 5, 5, 8, 9, 9, 20, 28, 67, 72, 73, 89, 93, 98, 278}) {
		t.Fatal("Taking(18)", sorted)
	}

	if sorted := MustGet(SliceOf(Taking(heapSort, 19), len(data)), data); !reflect.DeepEqual(sorted, []int{2, 3, 3, 4, 5, 5, 8, 9, 9, 20, 28, 67, 72, 73, 89, 93, 98, 278, 479}) {
		t.Fatal("Taking(19)", sorted)
	}

	if sorted := MustGet(SliceOf(Taking(heapSort, 20), len(data)), data); !reflect.DeepEqual(sorted, []int{2, 3, 3, 4, 5, 5, 8, 9, 9, 20, 28, 67, 72, 73, 89, 93, 98, 278, 479, 690}) {
		t.Fatal("Taking(20)", sorted)
	}

	if sorted := MustGet(SliceOf(Taking(heapSort, 21), len(data)), data); !reflect.DeepEqual(sorted, []int{2, 3, 3, 4, 5, 5, 8, 9, 9, 20, 28, 67, 72, 73, 89, 93, 98, 278, 479, 690, 9761}) {
		t.Fatal("Taking(21)", sorted)
	}

	if sorted := MustGet(SliceOf(Taking(heapSort, 22), len(data)), data); !reflect.DeepEqual(sorted, []int{2, 3, 3, 4, 5, 5, 8, 9, 9, 20, 28, 67, 72, 73, 89, 93, 98, 278, 479, 690, 9761}) {
		t.Fatal("Taking(22)", sorted)
	}

}

func TestOrderByOver(t *testing.T) {

	data := []int{278, 9, 3, 98, 3, 4, 690, 72, 8, 9761, 20, 93, 5, 9, 67, 28, 73, 479, 2, 5, 89}

	heapSort := Ordered(TraverseSlice[int](), OrderBy(Identity[int]()))

	if sorted := MustModify(Taking(heapSort, 2), Mul(2), data); !reflect.DeepEqual(sorted, []int{278, 9, 3, 98, 6, 4, 690, 72, 8, 9761, 20, 93, 5, 9, 67, 28, 73, 479, 4, 5, 89}) {
		t.Fatal("Taking(2)", sorted)
	}
}

func TestOrderByIdentity(t *testing.T) {

	data := []int{3, 2, 5, 4, 1}

	ordered := Ordered(
		TraverseSlice[int](),
		OrderBy(Identity[int]()),
	)

	viewResult := MustGet(SliceOf(ordered, len(data)), data)

	if !reflect.DeepEqual(viewResult, []int{1, 2, 3, 4, 5}) {
		t.Fatal(viewResult)
	}

	overResult := MustModify(ordered, Mul(2), data)

	if !reflect.DeepEqual(overResult, []int{6, 4, 10, 8, 2}) {
		t.Fatal(overResult)
	}

	takingResult := MustModify(Taking(ordered, 2), Mul(10), data)

	if !reflect.DeepEqual(takingResult, []int{3, 20, 5, 4, 10}) {
		t.Fatal(takingResult)
	}

	//Desc

	descResult := MustGet(
		SliceOf(
			Ordered(
				TraverseSlice[int](),
				Desc(OrderBy(Identity[int]())),
			),
			len(data),
		),
		data,
	)

	if !reflect.DeepEqual(descResult, []int{5, 4, 3, 2, 1}) {
		t.Fatal(descResult)
	}
}

func TestLevenshteinDistance(t *testing.T) {

	if ret := MustGet(EditDistance(TraverseString(), EditLevenshtein, EqT2[rune](), 10), lo.T2("hello", "hallo")); ret != 1 {

		t.Fatal("hello : hallo", ret)
	}

	if ret := MustGet(EditDistance(TraverseString(), EditLevenshtein, EqT2[rune](), 10), lo.T2("hello", "helo")); ret != 1 {
		t.Fatal("hello : helo", ret)
	}

	if ret := MustGet(EditDistance(TraverseString(), EditLevenshtein, EqT2[rune](), 10), lo.T2("hello", "halo")); ret != 2 {
		t.Fatal("hello : halo", ret)
	}

	if ret := MustGet(EditDistance(TraverseString(), EditLevenshtein, EqT2[rune](), 10), lo.T2("hello", "halol")); ret != 3 {
		t.Fatal("hello : halol EditLevenshtein", ret)
	}

	if ret := MustGet(EditDistance(TraverseString(), EditOSA, EqT2[rune](), 10), lo.T2("hello", "halol")); ret != 2 {
		t.Fatal("hello : halol EditRestricted", ret)
	}

}
