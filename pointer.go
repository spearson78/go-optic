package optic

import (
	"context"
	"errors"

	"github.com/samber/mo"
	"github.com/spearson78/go-optic/expr"
)

func TraversePtr[S any]() Optic[Void, *S, *S, S, S, ReturnMany, ReadWrite, BiDir, Pure] {
	return TraversePtrP[S, S]()
}

func TraversePtrP[S, T any]() Optic[Void, *S, *T, S, T, ReturnMany, ReadWrite, BiDir, Pure] {
	return PrismP[*S, *T, S, T](
		func(source *S) mo.Either[*T, S] {
			if source != nil {
				return mo.Right[*T, S](*source)
			} else {
				return mo.Left[*T, S](nil)
			}
		},
		func(focus T) *T {
			return &focus
		},
		ExprDef(
			func(t expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Traverse{
					OpticTypeExpr: t,
				}
			},
		),
	)
}

var ErrPtrNil = errors.New("nil passed to TraversePtrE / TraversePtrEP")

func TraversePtrE[S any]() Optic[Void, *S, *S, S, S, ReturnMany, ReadWrite, BiDir, Err] {
	return TraversePtrEP[S, S]()
}

func TraversePtrEP[S, T any]() Optic[Void, *S, *T, S, T, ReturnMany, ReadWrite, BiDir, Err] {
	return PrismEP[*S, *T, S, T](
		func(ctx context.Context, source *S) (mo.Either[*T, S], error) {
			if source != nil {
				return mo.Right[*T, S](*source), nil
			} else {
				return mo.Left[*T, S](nil), ErrPtrNil
			}
		},
		func(ctx context.Context, focus T) (*T, error) {
			return &focus, nil
		},
		ExprDef(
			func(t expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Traverse{
					OpticTypeExpr: t,
				}
			},
		),
	)
}
