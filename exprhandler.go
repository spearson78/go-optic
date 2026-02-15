package optic

import (
	"context"
	"errors"

	"github.com/spearson78/go-optic/expr"
)

var ignoredExprHandlerErr = errors.New("ExprHandler was ignored")

// The ExprHandler interface allows optics to take control of the execution of the complete expression tree.
//
// This can be used for example to translate optics to SQL statements.
type ExprHandler interface {
	TypeId() string

	Modify(ctx context.Context, o expr.OpticExpression, fmapExpr expr.OpticExpression, fmap func(index any, focus any, focusErr error) (any, error), source any) (any, bool, error)
	Set(ctx context.Context, o expr.OpticExpression, source any, val any) (any, error)
	Get(ctx context.Context, expr expr.OpticExpression, source any) (index any, value any, found bool, err error)
	ReverseGet(ctx context.Context, expr expr.OpticExpression, focus any) (any, error)
}

// ExprOptic returns an optic that will pass the complete expression tree to the given handler.
//
// This can be used for example to translate optics to SQL statements.
//
// Warning: the RET,RW,DIR constraints MUST match the abilities of the handler otherwise runtime errors will be reported instead of being reported at compile time.
//
// An ExprOptic is constructed from 2 parameters
//   - exprHandler is an [ExprHandler] that will be called when this optic is used in any action.
//   - expr: should return the expression type. See the [expr] package for more information.
func ExprOptic[I, S, T, A, B, RET, RW, DIR, ERR any](
	exprHandler ExprHandler,
	expression func(t expr.OpticTypeExpr) expr.OpticExpression,
) Optic[I, S, T, A, B, RET, RW, DIR, ERR] {

	return Omni[I, S, T, A, B, RET, RW, DIR, ERR](
		func(ctx context.Context, source S) (I, A, error) {
			var i I
			var a A
			return i, a, ignoredExprHandlerErr
		},
		func(ctx context.Context, focus B, source S) (T, error) {
			var t T
			return t, ignoredExprHandlerErr
		},
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				var i I
				var a A
				yield(ValIE(i, a, ignoredExprHandlerErr))
			}
		},
		func(ctx context.Context, source S) (int, error) {
			return 0, ignoredExprHandlerErr
		},
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {
			var t T
			return t, ignoredExprHandlerErr
		},
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				var i I
				var a A
				yield(ValIE(i, a, ignoredExprHandlerErr))
			}
		},
		func(indexA, indexB I) bool {
			panic(ignoredExprHandlerErr)
		},
		func(ctx context.Context, focus B) (T, error) {
			var t T
			return t, ignoredExprHandlerErr
		},
		ExpressionDef{
			handler:    func(ctx context.Context) (ExprHandler, error) { return exprHandler, nil },
			expression: expression,
		},
	)
}
