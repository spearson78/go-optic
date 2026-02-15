+++
title = "Optic Expressions"
weight = 12
+++
# Optic Expressions
When constructing a custom optic like a `Getter`  you are asked to provide an `OpticExpression` for typical use cases  `ExprCustom` can be used to provide a named `OpticExpression`

{{< playground file="/content/docs/10.custom/examples_test.go" id="getter1" >}}

This section of the documentation will explain the meaning of these expressions and what are they used for.

## Optic Error Path
When an error occurs during execution of an optic action an `Optic Error Path` is included. This is a kind of stack trace indicating where in the composed optic the error occurred.

{{< playground file="/content/docs/12.expression/examples_test.go" id="errorpath1" >}}

This example fails parsing the value `"two"` and outputs the following.

{{< result file="/content/docs/12.expression/examples_test.go" id="result_errorpath1" >}}

Notice the error includes an "optic error path" section. The top most entry in the path is the optic where the error occurred, in this case `ParseInt(10,0)`,
The next line indicates that optic from where `ParseInt` received it's value to parse, i this case a `Traverse`. The final line indicates that `SliceOf` is collecting the parsed results.
Every built in optic has a unique `OpticExpression` that not only captures the type of the optic but also all the parameters passed into it when constructed, notice that `ParseInt` captured it's parameters.

## AsExpr()

The user can access the `OpticExpression` by calling `AsExpr()` on any `Optic`

{{< playground file="/content/docs/12.expression/examples_test.go" id="asexpr1" >}}

In this example the optic expression is displayed as.

{{< result file="/content/docs/12.expression/examples_test.go" id="result_asexpr1" >}}

The actual in memory structure of this expression is actually this.

{{< playground file="/content/docs/12.expression/examples_test.go" id="asexpr2" >}}

The `OpticTypeExpr` captures the type parameters of the `Optic`. The `Compose` expression contains the expressions for the left and right side of the composition. `Traverse` requires no additional parameters other than the `OpticTypeExpr`. `ParseInt` captures the `Base` and `BitSize` parameters.

In this way the `OpticExpression` provides a complete reflection like runtime accessible representation of the `Optic`

## Expression Handlers
Go-optics supports overriding the default execution of an action by using the `ExprOptic`.

{{< code file="/content/docs/12.expression/examples_test.go" id="exproptic1" >}}

The `ExampleExprHandler` implements the `ExprHandler` interface.
```go
type ExprHandler interface {
	TypeId() string

	Modify(
		ctx context.Context,
		o expr.OpticExpression,
		fmapExpr expr.OpticExpression,
		fmap func(index any, focus any, focusErr error) (any, error),
		source any,
	) (any, bool, error)

	Set(
		ctx context.Context,
		o expr.OpticExpression,
		source any,
		val any,
	) (any, error)

	Get(
		ctx context.Context,
		expr expr.OpticExpression,
		source any,
	) (index any, value any, found bool, err error)

	ReverseGet(
		ctx context.Context,
		expr expr.OpticExpression,
		focus any,
	) (any, error)

}
```

`TypeId` should return a unique id for your `ExprHandler`

The other methods represent the built in actions. The `ExprOptic`  is an `Optic` and is fully compatible with other `Optics`

{{< code file="/content/docs/12.expression/examples_test.go" id="exproptic2" >}}

Here we can see that the `ExprOptic` has been passed to the `Ordered` combinator which in turn has been composed with the `Mul` optic.

During execution of an action go-optics will call into the `ExprHandler` and pass it the top level `OpticExpression`. In this way the `ExprHandler`  takes over complete execution of the action.


