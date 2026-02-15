package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

func TestExpressionOpticErrorPath(t *testing.T) {
	//BEGIN errorpath1
	data := []string{"1", "two", "3"}

	optic := Compose(
		TraverseSlice[string](),
		ParseInt[int](10, 0),
	)

	res, err := Get(SliceOf(optic, 3), data)
	fmt.Println(res, err)
	//END errorpath1

	if !reflect.DeepEqual([]any{res, err.Error()}, []any{
		//BEGIN result_errorpath1
		[]int(nil), `strconv.ParseInt: parsing "two": invalid syntax
optic error path:
	ParseInt(10,0)
	Traverse
	SliceOf(Traverse | ParseInt(10,0))
`,
		//END result_errorpath1
	}) {
		t.Fatal(res, err.Error())
	}
}

func TestExpressionAsExpr(t *testing.T) {
	//BEGIN asexpr1
	optic := Compose(
		TraverseSlice[string](),
		ParseInt[int](10, 0),
	)

	fmt.Println(optic.AsExpr().Short())
	//END asexpr1

	if !reflect.DeepEqual([]any{optic.AsExpr().Short()}, []any{
		//BEGIN result_asexpr1
		`Traverse | ParseInt(10,0)`,
		//END result_asexpr1
	}) {
		t.Fatal(optic.AsExpr().Short())
	}
}

func TestExpressionAsExprStructure(t *testing.T) {

	//BEGIN asexpr2
	e := expr.Compose{
		OpticTypeExpr: expr.NewOpticTypeExpr[Void, []string, []string, int, int, ReturnMany, ReadWrite, UniDir, Err](),
		Left: expr.Traverse{
			OpticTypeExpr: expr.NewOpticTypeExpr[int, []string, []string, string, string, ReturnMany, ReadWrite, UniDir, Pure](),
		},
		Right: expr.ParseInt{
			OpticTypeExpr: expr.NewOpticTypeExpr[Void, string, string, int, int, ReturnMany, ReadWrite, UniDir, Err](),
			Base:          10,
			BitSize:       0,
		},
	}

	fmt.Println(e.Short())
	//END asexpr2

	if !reflect.DeepEqual([]any{e.Short()}, []any{
		//BEGIN result_asexpr2
		`Traverse | ParseInt(10,0)`,
		//END result_asexpr2
	}) {
		t.Fatal(e.Short())
	}

}

// BEGIN example_expr_handler
type ExampleExprHandler struct {
}

// Get implements ExprHandler.
func (e *ExampleExprHandler) Get(ctx context.Context, expr expr.OpticExpression, source any) (any, any, bool, error) {
	fmt.Printf("ExampleExprHandler.Get(Optic:%v,Source:%v)\n", expr.Short(), source)

	//Interpret the expr using the source as a parameter
	//Return the result

	return nil, nil, false, errors.New("getter not implemented")
}

// Modify implements ExprHandler.
func (e *ExampleExprHandler) Modify(ctx context.Context, expr expr.OpticExpression, fmapExpr expr.OpticExpression, fmap func(index any, focus any, focusErr error) (any, error), source any) (any, bool, error) {
	fmt.Printf("ExampleExprHandler.Modify(Optic:%v,Fmap:%v,Source:%v)\n", expr.Short(), fmapExpr.Short(), source)

	//Interpret the expr using the source as a parameter
	//Either map the values using an interpreter over fmapExpr or call the fmap function if the mapping should be executed in go.
	//Return the result

	return nil, false, errors.New("modify not implemented")
}

// ReverseGet implements ExprHandler.
func (e *ExampleExprHandler) ReverseGet(ctx context.Context, expr expr.OpticExpression, focus any) (any, error) {
	fmt.Printf("ExampleExprHandler.ReverseGet(Optic:%v,Focus:%v)\n", expr.Short(), focus)

	//Interpret the expr using the focus as a parameter
	//and return the result

	return nil, errors.New("reverse get not implemented")
}

// Set implements ExprHandler.
func (e *ExampleExprHandler) Set(ctx context.Context, o expr.OpticExpression, source any, val any) (any, error) {
	fmt.Printf("ExampleExprHandler.Setter(Optic:%v,Source:%v,Val:%v)\n", o.Short(), source, val)

	//Interpret the expr using the focus as a parameter
	//Set th value to val
	//Return the result

	return nil, errors.New("setter not implemented")
}

// TypeId implements ExprHandler.
func (e *ExampleExprHandler) TypeId() string {
	//This ID must be unique. Only one ExpressionHandler TypeId is allowed in the entire expression tree.
	return "ExampleExprOptic"
}

//END example_expr_handler

func TestExpressionExprOptic(t *testing.T) {

	//This example is a stub demonstrating how an ExprOptic can be used.
	//A complete example is outside the scope of this documentation.

	//BEGIN exproptic1
	exprOptic := ExprOptic[int, string, string, int, int, ReturnMany, ReadWrite, UniDir, Err](
		&ExampleExprHandler{},
		expr.Custom("ExampleExprOptic"),
	)
	//END exproptic1

	//BEGIN exproptic2
	optic := Compose(
		Ordered(
			exprOptic,
			OrderBy(Identity[int]()),
		),
		Mul(10),
	)

	var dbConnectionString = "sqlite:example.db"

	_, err := Get(SliceOf(optic, 10), dbConnectionString)
	fmt.Println(err)

	_, err = Modify(optic, Mul(10), dbConnectionString)
	fmt.Println(err)

	_, err = Set(optic, 100, dbConnectionString)
	fmt.Println(err)
	//END exproptic2

	//Output:
	//ExampleExprHandler.Get(Optic:SeqOf(Ordered(Custom(ExampleExprOptic),OrderBy(Identity)) | * 10),Source:sqlite:example.db)
	//getter not implemented
	//ExampleExprHandler.Modify(Optic:Ordered(Custom(ExampleExprOptic),OrderBy(Identity)) | * 10,Fmap:ValueI[github.com/spearson78/go-optic.Void,int].value | * 10,Source:sqlite:example.db)
	//modify not implemented
	//ExampleExprHandler.Setter(Optic:Ordered(Custom(ExampleExprOptic),OrderBy(Identity)) | * 10,Source:sqlite:example.db,Val:100)
	//setter not implemented

}
