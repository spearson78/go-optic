package optic

import (
	"context"
	"errors"
	"strings"

	"github.com/samber/mo"
	"github.com/spearson78/go-optic/expr"
)

// OpticErrorPath is an error type that is joined with any error returned from an action.
// It contains the path to the optic where the error occurred.
// It can be extracted from the error using [errors.As]
type OpticErrorPath struct {
	path []func() expr.OpticExpression
}

func (err *OpticErrorPath) Error() string {
	var sb strings.Builder

	sb.WriteString("optic error path:\n")
	for _, v := range err.path {
		sb.WriteString("\t")
		sb.WriteString(v().Short())
		sb.WriteString("\n")
	}

	return sb.String()
}

// OpticError adds the given optic expression to the error if the error is != nil
//
// This should only be called from within an [Omni] optic when writing custom combinators.
func OpticError(err error, expression func() expr.OpticExpression) error {
	if err != nil {
		var opticError *OpticErrorPath
		if errors.As(err, &opticError) {
			opticError.path = append(opticError.path, expression)
			return err
		} else {
			return errors.Join(err, &OpticErrorPath{
				path: []func() expr.OpticExpression{
					expression,
				},
			})
		}
	} else {
		return nil
	}
}

// Error returns an [Iso] that always returns an error.
//
// See:
//   - [ErrorP] for a polymorphic version.
func Error[S, A any](err error) Optic[Void, S, S, A, A, ReturnOne, ReadWrite, BiDir, Err] {
	return ErrorP[S, S, A, A](err)
}

func ErrorI[I, S, A any](err error, ixMatch func(a, b I) bool) Optic[I, S, S, A, A, ReturnOne, ReadWrite, BiDir, Err] {
	return ErrorIP[I, S, S, A, A](err, ixMatch)
}

// Throw returns a [Lens] that returns the source error.
func Throw[A any]() Optic[Void, error, error, A, A, ReturnOne, ReadWrite, UniDir, Err] {
	return throwP[error, A, A]()
}

// Throw returns a polymorphic [Lens] that returns the source error.
func throwP[T, A, B any]() Optic[Void, error, T, A, B, ReturnOne, ReadWrite, UniDir, Err] {
	return CombiLens[ReadWrite, Err, Void, error, T, A, B](
		func(ctx context.Context, source error) (Void, A, error) {
			var a A
			return Void{}, a, source
		},
		func(ctx context.Context, focus B, source error) (T, error) {
			var t T
			return t, source
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Throw{
					OpticTypeExpr: ot,
				}
			},
		),
	)
}

// ErrorP returns a polymorphic [Iso] that always returns an error.
//
// See:
//   - [Error] for a non polymorphic version.
func ErrorP[S, T, A, B any](err error) Optic[Void, S, T, A, B, ReturnOne, ReadWrite, BiDir, Err] {
	return ErrorIP[Void, S, T, A, B](err, IxMatchVoid())
}

func ErrorIP[I, S, T, A, B any](err error, ixMatch func(a, b I) bool) Optic[I, S, T, A, B, ReturnOne, ReadWrite, BiDir, Err] {
	return Omni[I, S, T, A, B, ReturnOne, ReadWrite, BiDir, Err](
		func(ctx context.Context, source S) (I, A, error) {
			var i I
			var a A
			return i, a, err
		},
		func(ctx context.Context, focus B, source S) (T, error) {
			var t T
			return t, err
		},
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(ValueIE[I, A]) bool) {
				var i I
				var a A
				yield(ValIE(i, a, err))
			}
		},
		func(ctx context.Context, source S) (int, error) {
			return 0, err
		},
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {
			var t T
			return t, err
		},
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			return func(yield func(ValueIE[I, A]) bool) {
				var a A
				yield(ValIE(index, a, err))
			}
		},
		ixMatch,
		func(ctx context.Context, focus B) (T, error) {
			var t T
			return t, err
		},
		ExprDef(
			func(t expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Error{
					OpticTypeExpr: t,
					Err:           err,
				}
			},
		),
	)
}

// ErrorIs returns a [Predicate] that focuses true if the error is the target error
//
// See:
//   - [ErrorAs] for an error as version.
func ErrorIs(target error) Optic[Void, error, error, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return CombiGetter[Pure, Void, error, error, bool, bool](
		func(ctx context.Context, source error) (Void, bool, error) {
			return Void{}, errors.Is(source, target), nil
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.ErrorIs{
					OpticTypeExpr: ot,
					Target:        target,
				}
			},
		),
	)
}

// ErrorAs returns a [Prism] that focuses the error as the given type or nothing if the types do not match.
//
// See:
//   - [ErrorIs] for an errors.Is version.
func ErrorAs[A error]() Optic[Void, error, error, A, A, ReturnMany, ReadWrite, BiDir, Pure] {
	return PrismIP[Void, error, error, A, A](
		func(source error) mo.Either[error, ValueI[Void, A]] {
			var target A
			if errors.As(source, &target) {
				return mo.Right[error, ValueI[Void, A]](ValI(Void{}, target))
			} else {
				return mo.Left[error, ValueI[Void, A]](source)
			}
		},
		func(focus A) error {
			return focus
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.ErrorAs{
					OpticTypeExpr: ot,
				}
			},
		),
	)
}

func WithPanic[I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, T, A, B, RET, RW, DIR, Pure] {
	return Omni[I, S, T, A, B, RET, RW, DIR, Pure](
		func(ctx context.Context, source S) (I, A, error) {
			i, a, err := o.AsGetter()(ctx, source)
			if err != nil {
				panic(err)
			}
			return i, a, nil
		},
		func(ctx context.Context, focus B, source S) (T, error) {
			t, err := o.AsSetter()(ctx, focus, source)
			if err != nil {
				panic(err)
			}
			return t, nil
		},
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					index, focus, err := val.Get()
					if err != nil {
						panic(err)
					}
					ret := yield(ValIE(index, focus, nil))
					return ret
				})
			}

		},
		func(ctx context.Context, source S) (int, error) {
			l, err := o.AsLengthGetter()(ctx, source)
			if err != nil {
				panic(err)
			}
			return l, nil

		},
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {
			t, err := o.AsModify()(ctx, func(index I, focus A) (B, error) {
				b, err := fmap(index, focus)
				if err != nil {
					panic(err)
				}
				return b, nil
			}, source)
			return t, err
		},
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				o.AsIxGetter()(ctx, index, source)(func(val ValueIE[I, A]) bool {
					focusIndex, focus, err := val.Get()
					if err != nil {
						panic(err)
					}
					ret := yield(ValIE(focusIndex, focus, nil))
					return ret
				})
			}
		},
		o.AsIxMatch(),
		func(ctx context.Context, focus B) (T, error) {
			t, err := o.AsReverseGetter()(ctx, focus)
			if err != nil {
				panic(err)
			}
			return t, nil
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.WithPanic{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}
