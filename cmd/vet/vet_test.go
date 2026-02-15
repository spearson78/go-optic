package main

import (
	"context"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

func VetLength[I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[Void, S, S, int, int, ReturnOne, ReadOnly, UniDir, Err] {
	return OperatorE[S, int](
		func(ctx context.Context, source S) (int, error) {
			return o.AsLengthGetter()(ctx, source)
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Length{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}
