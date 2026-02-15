package optic_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/samber/mo"

	. "github.com/spearson78/go-optic"
)

func TestApplicativeConsistency(t *testing.T) {

	o := TraverseSlice[int]()

	ValidateOpticTestPred(t, WithLogging(o), []int{1, 2, 3, 4, 5}, 10, EqDeepT2[[]int]())
	ValidateOpticTestPredP(t, WithOption(o), []int{1, 2, 3, 4, 5}, 10, func(a int) mo.Option[int] { return mo.Some(a) }, func(t mo.Option[[]int]) []int { s, _ := t.Get(); return s }, EqDeepT2[[]int]())

	ValidateOpticTestPredP(t, WithComprehension(o), []int{1, 2, 3, 4, 5}, 10, func(i int) []int { return []int{i, i * 2} }, func(c [][]int) []int { return c[0] }, EqDeepT2[[]int]())
	ValidateOpticTestPredP(t, WithEither[string](o), []int{1, 2, 3, 4, 5}, 10, func(i int) mo.Either[string, int] { return mo.Right[string](i) }, func(t mo.Either[string, []int]) []int { s, _ := t.Right(); return s }, EqDeepT2[[]int]())
	ValidateOpticTestPredP(t, WithValidation[string](o), []int{1, 2, 3, 4, 5}, 10, func(i int) mo.Either[string, int] { return mo.Right[string](i) }, func(t mo.Either[[]string, []int]) []int { s, _ := t.Right(); return s }, EqDeepT2[[]int]())

}

func TestResultApplicative(t *testing.T) {

	testName := "applicative result"
	source := []string{"1", "beta", "2.0", "true"}
	//There is no result applicitave any more I had to roll errors into the standard optics.
	values := TraverseSlice[string]()

	//Compile error due to multiple return from each slice

	//if r := View(TraverseSlice[string](), []string{"a", "b"}); r != "1" {
	//	t.Fatalf("%v 1 : %v", testName, r)
	//}

	//Compile error due to multiple return from each slice

	//if r, err := ReView(values, "test"); err != nil || r != nil { //Embeding into the pure empty slice is the empty slice
	//	t.Fatalf("%v 2 : %v", testName, r)
	//}

	//Short circuits on the first error in and converts the error
	shortErr := errors.New("short")
	callCount := 0 //Verify that short circuit is working, AsResultJoin will combine all the errors together
	if r, err := Modify(values, OpE(func(ctx context.Context, s string) (string, error) {
		callCount++
		if len(s) < 3 {
			return strings.ToUpper(s), nil
		} else {
			return "", shortErr
		}
	}), source); err == nil || err.Error() != `short
optic error path:
	func1
	ValueI[int,string].value
	Traverse
` || callCount != 2 {
		t.Fatalf("%v 3 : %v : %v : callCount %v", testName, r, err, callCount)
	}

	if r, err := Get(SliceOf(Compose(values, ParseInt[int](10, 32)), len(source)), source); err == nil {
		t.Fatalf("%v 4 : %v", testName, r)
	}

	//Like ListOf but returns the first result only. THis leads to a nested Option Result
	//The Option indicates whether the list was empty. Result indicates if the first element was an error or not
	//I initially didn't like this but it makes logical sense
	//We could replace Option with error and in the case of a result return combine them together some how
	if r, ok, err := GetFirst(Compose(values, ParseInt[int](10, 32)), source); !ok || err != nil || r != 1 {
		t.Fatalf("%v 5 : %v", testName, r)
	}

	//If the first element is an error it is returned by FirstOf
	if r, ok, err := GetFirst(Compose(values, ParseInt[int](10, 32)), []string{"bad", "1"}); ok || err == nil {
		t.Fatalf("%v 5a : %v", testName, r)
	}

	//Set doesn't care the parse failed it just writes the new "unparsed" integer value directly into the result array
	if r, err := Set(Compose(values, ParseInt[int](10, 32)), 100, source); err != nil || !reflect.DeepEqual(r, []string{"100", "100", "100", "100"}) {
		t.Fatalf("%v 6 : %v : %v", testName, r, err)
	}

	//Parse error is returned
	if r, found, err := GetFirst(Filtered(Compose(values, ParseInt[int](10, 32)), Gt(3)), source); found || err == nil || err.Error() != `strconv.ParseInt: parsing "beta": invalid syntax
optic error path:
	ParseInt(10,32)
	Traverse
	Filtered(Const(true),ValueI[github.com/spearson78/go-optic.Void,int].value | > 3,FilterContinue,FilterContinue)
` {
		t.Fatalf("%v 6a : %v : %v : %v", testName, r, found, err)
	}

	//Parse error is ignored
	if r, found, err := GetFirst(Filtered(Compose(values, Ignore(ParseInt[int](10, 32), Const[error](true))), Gt(3)), source); found {
		t.Fatalf("%v 6a ignore : %v : %v", testName, r, err)
	}

	//We find a result before the parse error occurs
	if r, found, err := GetFirst(Filtered(Compose(values, ParseInt[int](10, 32)), Eq(1)), source); !found || err != nil || r != 1 {
		t.Fatalf("%v 6b : %v : %v", testName, r, err)
	}

	//Cannot sum due to bad string values
	if r, ok, err := GetFirst(Reduce(Compose(values, ParseInt[int](10, 32)), Sum[int]()), source); err == nil {
		t.Fatalf("%v 7 : %v : %v : %v", testName, r, err, ok)
	}

	if r, ok, err := GetFirst(Reduce(Compose(values, Ignore(ParseInt[int](10, 32), Const[error](true))), Sum[int]()), source); err != nil || r != 1 {
		t.Fatalf("%v 7 : %v : %v : %v", testName, r, err, ok)
	}

	if r, ok, err := GetFirst(Reduce(Compose(values, ParseInt[int](10, 32)), Product[int]()), source); err == nil {
		t.Fatalf("%v 8 : %v : %v", testName, r, ok)
	}

	//The last element exists but cannot be parsed so we get a parse error
	if r, found, err := GetFirst(Last(Compose(values, ParseInt[int](10, 32))), source); err == nil {
		t.Fatalf("%v 9 : %v , %v , %v", testName, r, found, err)
	}

	//The parser failurs are converted to an error result
	if r, ok, err := ModifyCheck(Taking(Compose(values, ParseInt[int](10, 32)), 2), Mul(100), source); !ok || err == nil {
		t.Fatalf("%v 10 : %v , %v , %v", testName, r, ok, err)
	}
}

func TestOption(t *testing.T) {

	testName := "applicative option"
	source := []string{"1", "beta", "2.0", "true"}
	values := WithOption(TraverseSlice[string]())
	isoOptic := WithOption(Mul(10))

	//Compile error due to multiple return from each slice

	//	if r := View(values, source); r != "1" {
	//		t.Fatalf("%v 1 : %v", testName, r)
	//	}

	// Compile error
	//if r, err := ReView(values, mo.Some("test")); err != nil || r.IsPresent() { //Embeding into the pure empty slice is none.
	//	t.Fatalf("%v 2 : r %v err %v", testName, r, err)
	//}

	//Short circuits on the first error in and converts the error
	callCount := 0 //Verify that short circuit is working, AsResultJoin will combine all the errors together
	if r, err := Modify(values, Op(func(s string) mo.Option[string] {
		callCount++
		if len(s) < 3 {
			return mo.Some(strings.ToUpper(s))
		} else {
			return mo.None[string]()
		}
	}), source); err != nil || r.IsPresent() || callCount != 2 {
		t.Fatalf("%v 3 : %v : %v : callCount %v", testName, r, err, callCount)
	}

	if r, err := Get(SliceOfP(values, len(source)), source); err != nil || !reflect.DeepEqual(r, []string{"1", "beta", "2.0", "true"}) {
		t.Fatalf("%v 4 : %v , %v", testName, r, err)
	}

	if r, ok, err := GetFirst(values, source); !ok || err != nil || r != "1" {
		t.Fatalf("%v 5 : %v", testName, r)
	}

	if r, err := Set(values, mo.Some("test"), source); err != nil || !reflect.DeepEqual(r, mo.Some([]string{"test", "test", "test", "test"})) {
		t.Fatalf("%v 6 : %v : %v", testName, r, err)
	}

	if r, err := Set(values, mo.None[string](), source); err != nil || r.IsPresent() {
		t.Fatalf("%v 6 : %v : %v", testName, r, err)
	}

	if r := MustGet(isoOptic, 100); r != 1000 {
		t.Fatalf("%v 7 : %v ", testName, r)
	}

	if r := MustReverseGet(isoOptic, mo.Some(100)); r != mo.Some(10) {
		t.Fatalf("%v 8 : %v ", testName, r)
	}

}

func TestFuncApplicative(t *testing.T) {

	testName := "applicative func"
	source := []string{"1", "beta", "2.0", "true"}
	values := WithFunc[bool](TraverseSlice[string]())
	isoOptic := WithFunc[bool](Mul(10))

	maybeToUpper := func(focus string) func(bool) (string, error) {
		return func(b bool) (string, error) {
			if b {
				return strings.ToUpper(focus), nil
			} else {
				return strings.ToLower(focus), nil
			}
		}
	}

	fnc, err := Modify(values, Op(maybeToUpper), source)
	if err != nil {
		t.Fatalf("%v 3 : %v ", testName, err)
	}

	if r, err := fnc(true); err != nil || !reflect.DeepEqual(r, []string{"1", "BETA", "2.0", "TRUE"}) {
		t.Fatalf("%v 3 : %v ", testName, r)
	}

	if r := MustGet(isoOptic, 100); r != 1000 {
		t.Fatalf("%v 7 : %v ", testName, r)
	}

	if r := MustReverseGet(isoOptic, func(b bool) (int, error) {
		if b {
			return 10, nil
		} else {
			return 100, nil
		}
	}); Must(r(true)) != 1 || Must(r(false)) != 10 {
		t.Fatalf("%v 8 : ", testName)
	}
}

func TestComprehension(t *testing.T) {

	testName := "applicative comprehension"
	source := []string{"alpha", "BETA", "GaMmA"}
	values := WithComprehension(TraverseSlice[string]())
	isoOptic := WithComprehension(Mul(10))

	//Compile error due to multiple return from each slice

	//	if r := View(values, source); r != "1" {
	//		t.Fatalf("%v 1 : %v", testName, r)
	//	}

	// Compile error
	//if r, err := ReView(values, mo.Some("test")); err != nil || r.IsPresent() { //Embeding into the pure empty slice is none.
	//	t.Fatalf("%v 2 : r %v err %v", testName, r, err)
	//}

	//Short circuits on the first error in and converts the error
	callCount := 0 //Verify that short circuit is working, AsResultJoin will combine all the errors together
	if r, err := Modify(values, Op(func(s string) []string {
		return []string{
			strings.ToUpper(s),
			strings.ToLower(s),
		}
	}), source); err != nil || !reflect.DeepEqual(r, [][]string{
		{"ALPHA", "BETA", "GAMMA"},
		{"ALPHA", "BETA", "gamma"},
		{"ALPHA", "beta", "GAMMA"},
		{"ALPHA", "beta", "gamma"},
		{"alpha", "BETA", "GAMMA"},
		{"alpha", "BETA", "gamma"},
		{"alpha", "beta", "GAMMA"},
		{"alpha", "beta", "gamma"},
	}) {
		t.Fatalf("%v 3 : %v : %v : callCount %v", testName, r, err, callCount)
	}

	if r, err := Get(SliceOfP(values, len(source)), source); err != nil || !reflect.DeepEqual(r, []string{"alpha", "BETA", "GaMmA"}) {
		t.Fatalf("%v 4 : %v , %v", testName, r, err)
	}

	if r, ok, err := GetFirst(values, source); !ok || err != nil || r != "alpha" {
		t.Fatalf("%v 5 : %v", testName, r)
	}

	if r, err := Set(values, []string{"A", "B"}, source); err != nil || !reflect.DeepEqual(r, [][]string{
		{"A", "A", "A"},
		{"A", "A", "B"},
		{"A", "B", "A"},
		{"A", "B", "B"},
		{"B", "A", "A"},
		{"B", "A", "B"},
		{"B", "B", "A"},
		{"B", "B", "B"},
	}) {
		t.Fatalf("%v 6 : %v : %v", testName, r, err)
	}

	if r := MustGet(isoOptic, 100); r != 1000 {
		t.Fatalf("%v 7 : %v ", testName, r)
	}

	if r := MustReverseGet(isoOptic, []int{100, 10}); !reflect.DeepEqual(r, []int{10, 1}) {
		t.Fatalf("%v 8 : %v ", testName, r)
	}

}

func TestWithEitherApplicative(t *testing.T) {

	testName := "applicative option"
	source := []string{"1", "beta", "2.0", "true"}
	values := WithEither[string](TraverseSlice[string]())
	isoOptic := WithEither[string](Mul(10))

	//Compile error due to multiple return from each slice

	//	if r := View(values, source); r != "1" {
	//		t.Fatalf("%v 1 : %v", testName, r)
	//	}

	// Compile error
	//if r, err := ReView(values, mo.Some("test")); err != nil || r.IsPresent() { //Embeding into the pure empty slice is none.
	//	t.Fatalf("%v 2 : r %v err %v", testName, r, err)
	//}

	//Short circuits on the first error in and converts the error
	if r, err := Modify(values, Op(func(s string) mo.Either[string, string] {
		if len(s) < 3 {
			return mo.Right[string, string](s)
		} else {
			return mo.Left[string, string]("string len >=3")
		}
	}), source); err != nil || r.IsRight() {
		t.Fatalf("%v 3 : %v : %v ", testName, r, err)
	}

	if r, err := Get(SliceOfP(values, len(source)), source); err != nil || !reflect.DeepEqual(r, []string{"1", "beta", "2.0", "true"}) {
		t.Fatalf("%v 4 : %v , %v", testName, r, err)
	}

	if r, ok, err := GetFirst(values, source); !ok || err != nil || r != "1" {
		t.Fatalf("%v 5 : %v", testName, r)
	}

	if r, err := Set(values, mo.Right[string, string]("test"), source); err != nil || !reflect.DeepEqual(r, mo.Right[string, []string]([]string{"test", "test", "test", "test"})) {
		t.Fatalf("%v 6 : %v : %v", testName, r, err)
	}

	if r, err := Set(values, mo.Left[string, string]("error"), source); err != nil || !reflect.DeepEqual(r, mo.Left[string, []string]("error")) {
		t.Fatalf("%v 6 : %v : %v", testName, r, err)
	}

	if r := MustGet(isoOptic, 100); r != 1000 {
		t.Fatalf("%v 7 : %v ", testName, r)
	}

	if r := MustReverseGet(isoOptic, mo.Left[string, int]("error")); r != mo.Left[string, int]("error") {
		t.Fatalf("%v 8 : %v ", testName, r)
	}

	if r := MustReverseGet(isoOptic, mo.Right[string, int](100)); r != mo.Right[string, int](10) {
		t.Fatalf("%v 8 : %v ", testName, r)
	}

}

func TestWithValidationApplicative(t *testing.T) {

	testName := "applicative option"
	source := []string{"1", "beta", "2.0", "true"}
	values := WithValidation[string](TraverseSlice[string]())
	isoOptic := WithValidation[string](Mul(10))

	//Compile error due to multiple return from each slice

	//	if r := View(values, source); r != "1" {
	//		t.Fatalf("%v 1 : %v", testName, r)
	//	}

	// Compile error
	//if r, err := ReView(values, mo.Some("test")); err != nil || r.IsPresent() { //Embeding into the pure empty slice is none.
	//	t.Fatalf("%v 2 : r %v err %v", testName, r, err)
	//}

	//Short circuits on the first error in and converts the error
	if r, err := Modify(values, Op(func(s string) mo.Either[string, string] {
		if len(s) < 4 {
			return mo.Left[string, string]("len(" + s + ")<4")
		} else {
			return mo.Right[string, string](s)
		}
	}), source); err != nil || r.IsRight() || len(r.MustLeft()) != 2 {
		t.Fatalf("%v 3 : %v : %v ", testName, r, err)
	}

	if r, err := Get(SliceOfP(values, len(source)), source); err != nil || !reflect.DeepEqual(r, []string{"1", "beta", "2.0", "true"}) {
		t.Fatalf("%v 4 : %v , %v", testName, r, err)
	}

	if r, ok, err := GetFirst(values, source); !ok || err != nil || r != "1" {
		t.Fatalf("%v 5 : %v", testName, r)
	}

	if r, err := Set(values, mo.Right[string, string]("test"), source); err != nil || !reflect.DeepEqual(r, mo.Right[[]string, []string]([]string{"test", "test", "test", "test"})) {
		t.Fatalf("%v 6 : %v : %v", testName, r, err)
	}

	if r, err := Set(values, mo.Left[string, string]("error"), source); err != nil || !reflect.DeepEqual(r, mo.Left[[]string, []string]([]string{"error", "error", "error", "error"})) {
		t.Fatalf("%v 6 : %v : %v", testName, r, err)
	}

	if r := MustGet(isoOptic, 100); r != 1000 {
		t.Fatalf("%v 7 : %v ", testName, r)
	}

	if r := MustReverseGet(isoOptic, mo.Left[string, int]("error")); !reflect.DeepEqual(r, mo.Left[[]string, int]([]string{"error"})) {
		t.Fatalf("%v 8 : %v ", testName, r)
	}

	if r := MustReverseGet(isoOptic, mo.Right[string, int](100)); !reflect.DeepEqual(r, mo.Right[[]string, int](10)) {
		t.Fatalf("%v 8 : %v ", testName, r)
	}

}

func TestAsFunc(t *testing.T) {

	fieldA := FieldLens(func(source *TestStruct) *string { return &source.A })

	fieldAFunc := WithFunc[bool](fieldA)

	testStruct := TestStruct{
		A: "Hello World",
	}

	uplowFunc, err := Modify(fieldAFunc, Op(func(focus string) func(bool) (string, error) {
		return func(b bool) (string, error) {
			if b {
				return strings.ToUpper(focus), nil
			} else {
				return strings.ToLower(focus), nil
			}
		}
	}), testStruct)
	if err != nil {
		t.Fatalf("TraverseOf AsFunc %v", err)
	}

	upper, err := uplowFunc(true)
	if err != nil {
		t.Fatalf("TraverseOf AsFunc upper %v", err)
	}
	lower, err := uplowFunc(false)
	if err != nil {
		t.Fatalf("TraverseOf AsFunc lower %v", err)
	}

	if upper.A != "HELLO WORLD" {
		t.Fatalf("upper.A : %v", upper.A)
	}

	if lower.A != "hello world" {
		t.Fatalf("upper.A : %v", lower.A)
	}

}

func TestWithEither(t *testing.T) {

	fieldB := FieldLens(func(source *TestStruct) *any { return &source.B })

	fieldAFunc := WithEither[error](fieldB)

	testStruct := TestStruct{
		B: errors.New("Hello World"),
	}

	withEither, err := Modify(fieldAFunc, Op(func(focus any) mo.Either[error, any] {
		if stringer, ok := focus.(error); ok {
			return mo.Left[error, any](stringer)
		} else {
			return mo.Right[error, any](focus)
		}
	}), testStruct)

	if err != nil || withEither.IsRight() {
		t.Fatalf("RIGHT %v", withEither)
	}
}

func ExampleWithLogging() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	result := MustModify(
		WithLogging(
			TraverseSlice[string](),
		),
		Op(strings.ToUpper),
		data,
	)

	fmt.Println(result)

	//Output:
	//[ALPHA BETA GAMMA DELTA]

}

func ExampleWithOption() {

	goodData := []string{"1", "2", "3", "4"}
	badData := []string{"1", "2", "three", "4"}

	//TraverseSliceP normally converts from []string to []int
	//The WithOption modifies this to mo.Option[[]int]
	baseOptic := TraverseSliceP[string, int]()
	optionOptic := WithOption(baseOptic)

	//View actions are identical to the non applicative version of the optic.

	//Modification actions have optional return values.

	var goodOverResult mo.Option[[]int] = MustModify(optionOptic, Op(func(focus string) mo.Option[int] {
		i, err := strconv.ParseInt(focus, 10, 32)
		//Each return is also now an option
		//The first none option will cause the final result to short circuit to none
		return mo.TupleToOption(int(i), err == nil)
	}), goodData)
	fmt.Println(goodOverResult.Get())

	var badOverResult mo.Option[[]int] = MustModify(optionOptic, Op(func(focus string) mo.Option[int] {
		i, err := strconv.ParseInt(focus, 10, 32)
		//Each return is also now an option
		//The first none option will cause the final result to short circuit to none
		return mo.TupleToOption(int(i), err == nil)
	}), badData)
	fmt.Println(badOverResult.Get())

	//WithOption also works with composed optics

	composedData := []lo.Tuple2[string, string]{
		lo.T2("1", "alpha"),
		lo.T2("2", "beta"),
		lo.T2("3", "gamma"),
		lo.T2("4", "delta"),
	}

	composedOptic := Compose(
		TraverseSliceP[lo.Tuple2[string, string], lo.Tuple2[int, string]](),
		T2AP[string, string, int](),
	)

	composedOptionOptic := WithOption(composedOptic)

	var goodComposedResult mo.Option[[]lo.Tuple2[int, string]] = MustModify(composedOptionOptic, Op(func(focus string) mo.Option[int] {
		i, err := strconv.ParseInt(focus, 10, 32)
		//Each return is also now an option
		//The first none option will cause the final result to short circuit to none
		return mo.TupleToOption(int(i), err == nil)
	}), composedData)
	fmt.Println(goodComposedResult.Get())

	//Output:
	//[1 2 3 4] true
	//[] false
	//[{1 alpha} {2 beta} {3 gamma} {4 delta}] true
}

func ExampleWithFunc() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	//The WithFunc modifies this the optic to return a Func(string) []string
	optic := WithFunc[string](TraverseSlice[string]())

	//View actions are identical to the non applicative version of the optic.

	//Modification actions have a func as the return value

	var overResult func(string) ([]string, error) = MustModify(optic, Op(func(focus string) func(string) (string, error) {
		return func(param string) (string, error) {
			return focus + "-" + param, nil
		}
	}), data)

	//Calling the return funcion passes the given parameter into the over action.
	fmt.Println(overResult("quadrant"))
	fmt.Println(overResult("test"))

	//WithFunc also works with composed optics

	composedData := []lo.Tuple2[string, string]{
		lo.T2("1", "alpha"),
		lo.T2("2", "beta"),
		lo.T2("3", "gamma"),
		lo.T2("4", "delta"),
	}

	composedOptic := Compose(
		TraverseSlice[lo.Tuple2[string, string]](),
		T2B[string, string](),
	)

	composedFuncOptic := WithFunc[string](composedOptic)

	var composedResult func(string) ([]lo.Tuple2[string, string], error) = MustModify(composedFuncOptic, Op(func(focus string) func(string) (string, error) {
		return func(param string) (string, error) {
			return focus + "-" + param, nil
		}
	}), composedData)
	fmt.Println(composedResult("quadrant"))

	//Output:
	//[alpha-quadrant beta-quadrant gamma-quadrant delta-quadrant] <nil>
	//[alpha-test beta-test gamma-test delta-test] <nil>
	//[{1 alpha-quadrant} {2 beta-quadrant} {3 gamma-quadrant} {4 delta-quadrant}] <nil>
}

func ExampleWithComprehension() {

	data := []string{"alpha", "beta"}

	//The WithComprehension modifies this the optic to return an [][]string
	optic := WithComprehension(TraverseSlice[string]())

	//View actions are identical to the non applicative version of the optic.

	//Modification actions have a slice as the return value
	var overResult [][]string = MustModify(optic, Op(func(focus string) []string {
		return []string{
			focus + "-quadrant",
			focus + "-test",
		}
	}), data)

	//In the result every focused value is combined with each value in the return slice in every possible combination.
	fmt.Println(overResult)

	//WithComprehension also works with composed optics

	composedData := []lo.Tuple2[string, string]{
		lo.T2("1", "alpha"),
		lo.T2("2", "beta"),
	}

	composedOptic := Compose(
		TraverseSlice[lo.Tuple2[string, string]](),
		T2B[string, string](),
	)

	composedFuncOptic := WithComprehension(composedOptic)

	var composedResult [][]lo.Tuple2[string, string] = MustModify(composedFuncOptic, Op(func(focus string) []string {
		return []string{
			focus + "-quadrant",
			focus + "-test",
		}
	}), composedData)
	fmt.Println(composedResult)

	//Output:
	//[[alpha-quadrant beta-quadrant] [alpha-quadrant beta-test] [alpha-test beta-quadrant] [alpha-test beta-test]]
	//[[{1 alpha-quadrant} {2 beta-quadrant}] [{1 alpha-quadrant} {2 beta-test}] [{1 alpha-test} {2 beta-quadrant}] [{1 alpha-test} {2 beta-test}]]
}

func ExampleWithEither() {

	goodData := []string{"1", "2", "3", "4"}
	badData := []string{"1", "2", "three", "four"}

	//TraverseSliceP normally converts from []string to []int
	//The AWithEither modifies this to mo.Either[string,[]int]
	baseOptic := TraverseSliceP[string, int]()
	optionOptic := WithEither[string](baseOptic)

	//View actions are identical to the non applicative version of the optic.

	//Modification actions have an either return value.
	var goodOverResult mo.Either[string, []int] = MustModify(optionOptic, Op(func(focus string) mo.Either[string, int] {
		i, err := strconv.ParseInt(focus, 10, 32)
		if err != nil {
			//We can use the left value to report the invalid string
			return mo.Left[string, int](focus)
		} else {
			return mo.Right[string, int](int(i))
		}
	}), goodData)
	fmt.Println(goodOverResult)

	var badOverResult mo.Either[string, []int] = MustModify(optionOptic, Op(func(focus string) mo.Either[string, int] {
		i, err := strconv.ParseInt(focus, 10, 32)
		if err != nil {
			//We can use the left value to report the invalid string
			return mo.Left[string, int](focus)
		} else {
			return mo.Right[string, int](int(i))
		}
	}), badData)
	fmt.Println(badOverResult)

	//Output:
	//{false  [1 2 3 4]}
	//{true three []}
}

func ExampleWithValidation() {

	goodData := []string{"1", "2", "3", "4"}
	badData := []string{"1", "2", "three", "four"}

	//TraverseSliceP normally converts from []string to []int
	//The WithValidation modifies this to mo.Either[[]string,[]int]
	baseOptic := TraverseSliceP[string, int]()
	optionOptic := WithValidation[string](baseOptic)

	//View actions are identical to the non applicative version of the optic.

	//Modification actions have an either return value.
	var goodOverResult mo.Either[[]string, []int] = MustModify(optionOptic, Op(func(focus string) mo.Either[string, int] {
		i, err := strconv.ParseInt(focus, 10, 32)
		if err != nil {
			//We can use the left value to report the invalid string
			return mo.Left[string, int](focus)
		} else {
			return mo.Right[string, int](int(i))
		}
	}), goodData)
	fmt.Println(goodOverResult)

	var badOverResult mo.Either[[]string, []int] = MustModify(optionOptic, Op(func(focus string) mo.Either[string, int] {
		i, err := strconv.ParseInt(focus, 10, 32)
		if err != nil {
			//We can use the left value to report the invalid string
			return mo.Left[string, int](focus)
		} else {
			return mo.Right[string, int](int(i))
		}
	}), badData)
	fmt.Println(badOverResult)

	//Output:
	//{false [] [1 2 3 4]}
	//{true [three four] []}
}

func TestWithOptionAndTraverseCol(t *testing.T) {

	data := ValCol[int](1, 2, 3)

	result, ok, err := ModifyCheckContext(
		context.Background(),
		WithOption(TraverseCol[int, int]()),
		Error[int, mo.Option[int]](errors.New("sabotage")),
		data,
	)

	if err == nil {
		t.Fatal(result, ok, err)
	}

}
