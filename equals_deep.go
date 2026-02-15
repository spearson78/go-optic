package optic

import (
	"context"
	"reflect"

	"github.com/samber/lo"
	"github.com/spearson78/go-optic/expr"
)

// EqDeepT2 returns an [BinaryOp] that is satisfied if A is [reflect.DeepEqual] to B.
//
// See:
//   - [EqDeep] for a unary version.
//   - [EqDeepOp] for a version that is applied to the focus of 2 [Operation]s.
func EqDeepT2[A any]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	//EqDeepT2Of is not needed we can pass EqDeepT2 into EqT2Of
	return CombiGetter[Pure, Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool](
		func(ctx context.Context, source lo.Tuple2[A, A]) (Void, bool, error) {
			ret := reflect.DeepEqual(source.A, source.B)
			return Void{}, ret, nil
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.EqDeep{
				OpticTypeExpr: ot,
			}
		}),
	)
}

// EqDeep returns a [Predicate] that is satisfied if the focused value is [reflect.DeepEqual] to (==) the provided constant value.
//
// See:
//   - [DeepEqOp] for a [Predicate] that compares 2 focuses instead of a constant.
func EqDeep[A any](right A) Optic[Void, A, A, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2Bind(EqDeepT2[A](), right)
}

// EqDeepOp returns a [Predicate] that is satisfied if the left focus is [reflect.DeepEqual] to (==) the right focus.
//
// See [EqDeep] for a predicate that checks constant a value rather than 2 focuses.
func EqDeepOp[S any, A any, RETL, RETR TReturnOne, LERR, RERR any](left Operation[S, A, RETL, LERR], right Operation[S, A, RETR, RERR]) Optic[Void, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, CompositionTree[LERR, RERR]] {
	return EErrL(OpT2ToOptic(
		left,
		EqDeepT2[A](),
		right,
	))
}
