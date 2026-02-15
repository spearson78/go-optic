package optic

import (
	"context"
	"errors"
	"reflect"

	"github.com/spearson78/go-optic/expr"
)

type withVariableKey string

type sentinilT struct{}

var sentinil = any((*sentinilT)(nil))

// The WithVar combinator creates a named variable that is accessible inside the withVar optic.
//
// See:
//   - [Var] for a n optic that focuses the value of the variable.
func WithVar[I, S, T, A, B any, RET any, RW, DIR, ERR any, V any, VRET TReturnOne, VERR any](withVar Optic[I, S, T, A, B, RET, RW, DIR, ERR], name string, value Operation[S, V, VRET, VERR]) Optic[I, S, T, A, B, RET, RW, UniDir, CompositionTree[ERR, VERR]] {
	return Omni[I, S, T, A, B, RET, RW, UniDir, CompositionTree[ERR, VERR]](
		func(ctx context.Context, source S) (I, A, error) {

			val, err := value.AsOpGet()(ctx, source)
			if err != nil {
				var i I
				var a A
				return i, a, err
			}

			av := any(val)
			if av == nil {
				av = sentinil
			}

			varCtx := context.WithValue(ctx, withVariableKey(name), av)

			return withVar.AsGetter()(varCtx, source)
		},
		func(ctx context.Context, focus B, source S) (T, error) {
			val, err := value.AsOpGet()(ctx, source)
			if err != nil {
				var t T
				return t, err
			}

			av := any(val)
			if av == nil {
				av = sentinil
			}

			varCtx := context.WithValue(ctx, withVariableKey(name), av)

			return withVar.AsSetter()(varCtx, focus, source)
		},
		func(ctx context.Context, source S) SeqIE[I, A] {
			val, err := value.AsOpGet()(ctx, source)
			if err != nil {
				return func(yield func(val ValueIE[I, A]) bool) {
					var i I
					var a A
					yield(ValIE(i, a, err))
				}
			}

			av := any(val)
			if av == nil {
				av = sentinil
			}

			varCtx := context.WithValue(ctx, withVariableKey(name), av)

			return withVar.AsIter()(varCtx, source)
		},
		func(ctx context.Context, source S) (int, error) {
			val, err := value.AsOpGet()(ctx, source)
			if err != nil {
				return 0, err
			}

			av := any(val)
			if av == nil {
				av = sentinil
			}

			varCtx := context.WithValue(ctx, withVariableKey(name), av)

			return withVar.AsLengthGetter()(varCtx, source)
		},
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {
			val, err := value.AsOpGet()(ctx, source)
			if err != nil {
				var t T
				return t, err
			}

			av := any(val)
			if av == nil {
				av = sentinil
			}

			varCtx := context.WithValue(ctx, withVariableKey(name), av)

			return withVar.AsModify()(varCtx, fmap, source)
		},
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			val, err := value.AsOpGet()(ctx, source)
			if err != nil {
				return func(yield func(val ValueIE[I, A]) bool) {
					var a A
					yield(ValIE(index, a, err))
				}
			}

			av := any(val)
			if av == nil {
				av = sentinil
			}

			varCtx := context.WithValue(ctx, withVariableKey(name), av)

			return withVar.AsIxGetter()(varCtx, index, source)
		},
		withVar.AsIxMatch(),
		unsupportedReverseGetter[B, T],
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.WithVar{
					OpticTypeExpr: ot,
					WithVar:       withVar.AsExpr(),
					Name:          name,
					Value:         value.AsExpr(),
				}
			},
			withVar,
			value,
		),
	)
}

var ErrVariableNotDefined = errors.New("variable not defined")

// Var is an [Operator] that focuses the value of the variable.
//
// Note: if the variable is not defined an [ErrVariableNotDefined] will be returned.
func Var[S, A any](name string) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, Err] {

	return OperatorE[S, A](
		func(ctx context.Context, source S) (A, error) {
			val := ctx.Value(withVariableKey(name))
			if val == nil {
				var a A
				return a, ErrVariableNotDefined
			}

			_, isSentinil := val.(*sentinilT)
			if isSentinil {
				var a A
				return a, nil
			}

			a, ok := val.(A)
			if !ok {
				return a, &ErrCast{
					From: reflect.TypeOf(a),
					To:   reflect.TypeFor[A](),
				}
			}

			return a, nil
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Var{
				OpticTypeExpr: ot,
				Name:          name,
			}
		}),
	)
}
