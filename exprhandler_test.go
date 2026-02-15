package optic_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/spearson78/go-optic/expr"

	. "github.com/spearson78/go-optic"
)

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

func ExampleExprOptic() {

	//This example is a stub demonstrating how an ExprOptic can be used.
	//A complete example is outside the scope of this documentation.

	intsFromDb := ExprOptic[int, string, string, int, int, ReturnMany, ReadWrite, UniDir, Err](
		&ExampleExprHandler{},
		expr.Custom("ExampleExprOptic"),
	)

	optic := Compose(
		Ordered(
			intsFromDb,
			OrderBy(Identity[int]()),
		),
		Mul(10),
	)

	//Note: even though the intsFromDb [ExprOptic] is nested in a compose and ordered it's handler takes control of executing each action completely
	//The expression tree passed to the handler is always the root optic.

	var dbConnectionString = "sqlite:example.db"

	_, err := Get(SeqEOf(optic), dbConnectionString)
	fmt.Println(err)

	_, err = Modify(optic, Mul(10), dbConnectionString)
	fmt.Println(err)

	_, err = Set(optic, 100, dbConnectionString)
	fmt.Println(err)

	//Output:
	//ExampleExprHandler.Get(Optic:SeqOf(Ordered(Custom(ExampleExprOptic),OrderBy(Identity)) | * 10),Source:sqlite:example.db)
	//getter not implemented
	//ExampleExprHandler.Modify(Optic:Ordered(Custom(ExampleExprOptic),OrderBy(Identity)) | * 10,Fmap:ValueI[github.com/spearson78/go-optic.Void,int].value | * 10,Source:sqlite:example.db)
	//modify not implemented
	//ExampleExprHandler.Setter(Optic:Ordered(Custom(ExampleExprOptic),OrderBy(Identity)) | * 10,Source:sqlite:example.db,Val:100)
	//setter not implemented

}
