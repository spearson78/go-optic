package optic

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/samber/lo"
)

func ExampleErrorIs() {

	data := []string{"1", "2", "three", "4"}

	optic := Compose(
		TraverseSlice[string](),
		Ignore(
			ParseInt[int](10, 32),
			ErrorIs(strconv.ErrSyntax),
		),
	)

	var viewResult []int = MustGet(SliceOf(optic, len(data)), data)
	fmt.Println(viewResult)

	var modifyResult []string = MustModify(optic, Mul(10), data)
	fmt.Println(modifyResult)

	//Output:
	//[1 2 4]
	//[10 20 three 40]
}

func ExampleErrorAs() {

	optic := Catch(
		ParseInt[int32](10, 32),
		FirstOrDefault(
			Coalesce(
				Compose(
					ErrorAs[*strconv.NumError](),
					OpE(func(ctx context.Context, source *strconv.NumError) (int32, error) {
						return -1, fmt.Errorf("custom parse error of value %v", source.Num)
					}),
				),
				Throw[int32](),
			),
			-1,
		),
	)

	viewResult, ok, err := GetFirst(optic, "1")
	fmt.Println(viewResult, ok, err)

	viewResult, ok, err = GetFirst(optic, "one")
	fmt.Println(viewResult, ok, err)

	//Output:
	//1 true <nil>
	//-1 false custom parse error of value one
	//optic error path:
	//	func1
	//	ErrorAs(*strconv.NumError)
	//	Coalesce(ErrorAs(*strconv.NumError) | func1,Throw())
	//	ReIndexed(Const({}),Coalesce(ErrorAs(*strconv.NumError) | func1,Throw()))
	//	Polymorphic(ReIndexed(Const({}),Coalesce(ErrorAs(*strconv.NumError) | func1,Throw())))
	//	Coalesce(Polymorphic(ReIndexed(Const({}),Coalesce(ErrorAs(*strconv.NumError) | func1,Throw()))),Const(-1))
	//	Filtered(< 1,Const(true),FilterContinue,FilterStop)
	//	Catch(ParseInt(10,32),Filtered(< 1,Const(true),FilterContinue,FilterStop),Throw())
}

func ExampleError() {

	data := []int{1, 2, -3, 4}

	optic := Compose(
		TraverseSlice[int](),
		If(
			Lt(0),
			Error[int, int](errors.New("<0")),
			Identity[int](),
		),
	)

	res, err := Get(SliceOf(optic, len(data)), data)
	fmt.Println(res, err)

	//Output:
	//[] <0
	//optic error path:
	//	Error(<0)
	//	Switch(Case(< 0 -> Error(<0))Default(Identity))
	//	Traverse
	//	SliceOf(Traverse | Switch(Case(< 0 -> Error(<0))Default(Identity)))
}

func ExampleErrorP() {

	optic := Compose(
		TraverseSlice[any](),
		Coalesce(
			DownCastP[any, string, int](),
			ErrorP[any, any, string, int](errors.New("down cast failed")),
		),
	)

	res, err := Modify(optic, ParseInt[int](10, 0), []any{"1", "2", "3"})
	fmt.Println(res, err)

	res, err = Modify(optic, ParseInt[int](10, 0), []any{"1", false, "3"})
	fmt.Println(res, err)

	//Output:
	//[1 2 3] <nil>
	//[] down cast failed
	//optic error path:
	//	Error(down cast failed)
	//	Coalesce(Cast,Error(down cast failed))
	//	Traverse
}

func ExampleThrow() {

	optic := Catch(
		ParseInt[int32](10, 32),
		If(
			ErrorIs(strconv.ErrSyntax),
			EErr(Const[error](int32(-1))),
			Ro(Throw[int32]()),
		),
	)

	viewResult, ok, err := GetFirst(optic, "1")
	fmt.Println(viewResult, ok, err)

	viewResult, ok, err = GetFirst(optic, "one")
	fmt.Println(viewResult, ok, err)

	viewResult, ok, err = GetFirst(optic, "2147483648")
	fmt.Println(viewResult, ok, err)

	//Output:
	//1 true <nil>
	//-1 true <nil>
	//0 false strconv.ParseInt: parsing "2147483648": value out of range
	//optic error path:
	//	ParseInt(10,32)
	//	Throw()
	//	Switch(Case(ErrorIs(invalid syntax) -> Const(-1))Default(Throw()))
	//	Catch(ParseInt(10,32),Switch(Case(ErrorIs(invalid syntax) -> Const(-1))Default(Throw())),Throw())
}

func TestErrorReporting(t *testing.T) {

	data := []lo.Tuple2[string, string]{lo.T2("one", "1"), lo.T2("1", "2")}

	_, err := Modify(Compose4(Index(TraverseSlice[lo.Tuple2[string, string]](), 0), T2A[string, string](), ParseInt[int](10, 32), Mul(10)), OpE(func(ctx context.Context, i int) (int, error) {
		return 0, errors.New("sabotage")
	}), data)

	if err.Error() != `strconv.ParseInt: parsing "one": invalid syntax
optic error path:
	ParseInt(10,32)
	TupleElement(0)
	Traverse
	Index(0)
` {
		t.Fatal(err)
	}

	_, err = Modify(Compose4(Index(TraverseSlice[lo.Tuple2[string, string]](), 1), T2A[string, string](), ParseInt[int](10, 32), Mul(10)), OpE(func(ctx context.Context, i int) (int, error) {
		return 0, errors.New("sabotage")
	}), data)

	if err.Error() != `sabotage
optic error path:
	func2
	ValueI[github.com/spearson78/go-optic.Void,int].value
	ParseInt(10,32) | * 10
	TupleElement(0)
	Traverse
	Index(1)
` {
		t.Fatal(err)
	}

}

func TestSafeYield(t *testing.T) {

	iteration := Iteration[int, int](
		func(limit int) iter.Seq[int] {
			return func(yield func(focus int) bool) {
				for i := 0; i < limit; i++ {
					yield(i)
				}
			}
		},
		nil,
		ExprCustom("iteration"),
	)

	defer func() {
		if r := recover(); r != nil {
			if r == yieldAfterBreak {
				//Ignore and let the error be returned
			} else {
				panic(r)
			}
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	GetContext(ctx, Length(iteration), math.MaxInt)

	t.Fatal("yieldAfterBreak panics expected")

}

func TestViewContextDeadline(t *testing.T) {

	iteration := Iteration[int, int](
		func(limit int) iter.Seq[int] {
			return func(yield func(focus int) bool) {
				for i := 0; i < limit; i++ {
					if !yield(i) {
						break
					}
				}
			}
		},
		nil,
		ExprCustom("iteration"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := GetContext(ctx, Length(iteration), math.MaxInt)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal("expected deadline error")
	}

}

func TestAppendingSafeYield(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			if r == yieldAfterBreak {
				//Ignore and let the error be returned
			} else {
				panic(r)
			}
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	ModifyContext(ctx, Identity[[]int](), AppendSlice(Col(
		func(yield func(focus int) bool) {
			i := 0
			for {
				yield(i)
				i++
			}

		},
		func() int {
			return math.MaxInt
		},
	)), []int{1, 2, 3})

	t.Fatal("yieldAfterBreak panics expected")

}

func TestAppendingDeadline(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := ModifyContext(ctx, Identity[[]int](), AppendSlice(
		Col(
			func(yield func(focus int) bool) {
				i := 0
				for yield(i) {
					i++
				}
			},
			func() int {
				return math.MaxInt
			},
		),
	), []int{1, 2, 3})

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal("expected deadline error")
	}

}
