package optic

import (
	"cmp"
	"context"
	"errors"

	"github.com/spearson78/go-optic/expr"
)

// An PredicateI is an [OperationI] that return true if the index aware predicate is satisfied.
type PredicateI[I, A, ERR any] Predicate[ValueI[I, A], ERR]

// An PredicateIE is an [PredicateI] where it should be assumed an error is returned.
// This interface can be used to implement fluent style APIs at the cost of no longer being pure.
type PredicateIE[I, A any] PredicateE[ValueI[I, A]]

func PredGetI[I, S, ERR any](ctx context.Context, fnc PredicateI[I, S, ERR], source S, index I) (bool, error) {
	ret, err := fnc.AsOpGet()(ctx, ValI(index, source))
	if errors.Is(err, ErrEmptyGet) {
		return false, nil
	}
	return ret, err
}

// EqI returns a [PredicateI] that is satisfied if the focused index is equal to (==) the provided constant value.
//
// See:
//   - [Eq] for a version that compares the focused value and not the index.
func EqI[A any, P comparable](right P) Optic[Void, ValueI[P, A], ValueI[P, A], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpOnIx[A](Eq(right))
}

// NeI returns a [PredicateI] that is satisfied if the focused index is not equal to (!=) the provided constant value.
//
// See:
//   - [Ne] for a version that compares the focused value and not the index.
func NeI[A any, P comparable](right P) Optic[Void, ValueI[P, A], ValueI[P, A], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpOnIx[A](Ne(right))
}

// GtI returns a [PredicateI] that is satisfied if the focused index is greater than (>) the provided constant value.
//
// See:
//   - [Gt] for a version that compares the focused value and not the index.
func GtI[A any, P cmp.Ordered](right P) Optic[Void, ValueI[P, A], ValueI[P, A], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpOnIx[A](Gt(right))
}

// GteI returns a [PredicateI] that is satisfied if the focused index is greater than or equal to than (>=) the provided constant value.
//
// See:
//   - [Gte] for a version that compares the focused value and not the index.
func GteI[A any, P cmp.Ordered](right P) Optic[Void, ValueI[P, A], ValueI[P, A], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpOnIx[A](Gte(right))
}

// IL returns a [PredicateI] that is satisfied if the focused index is less than (<) the provided constant value.
//
// See:
//   - [Gt] for a version that compares the focused value and not the index.
func LtI[A any, P cmp.Ordered](right P) Optic[Void, ValueI[P, A], ValueI[P, A], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpOnIx[A](Lt(right))
}

// LteI returns a [PredicateI] that is satisfied if the focused index is less than or equal to than (<=) the provided constant value.
//
// See:
//   - [Lte] for a version that compares the focused value and not the index.
func LteI[A any, P cmp.Ordered](right P) Optic[Void, ValueI[P, A], ValueI[P, A], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpOnIx[A](Lte(right))
}

// AnyI returns an [Predicate] that is satisfied if any focus satisfies the given predicate
//
// Note: If the optic focuses no elements then Any returns false.
//
// See:
//   - [Any] for a non index aware version.
//   - [AllI] for a version that requires that all focuses satisfy the predicate.
func AnyI[I, S, T, A, B, RET, RW, DIR, ERR, PERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], pred PredicateI[I, A, PERR]) Optic[I, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, CompositionTree[ERR, PERR]] {

	return CombiGetter[CompositionTree[ERR, PERR], I, S, S, bool, bool](
		func(ctx context.Context, source S) (I, bool, error) {
			var err error
			found := false
			var foundIndex I
			o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
				focusIndex, v, focusErr := val.Get()
				if focusErr != nil {
					err = focusErr
					found = false
					return false
				}
				found, err = PredGetI(ctx, pred, v, focusIndex)
				err = JoinCtxErr(ctx, err)
				if err != nil {
					found = false
					return false
				}
				if found {
					foundIndex = focusIndex
				}
				return !found
			})
			return foundIndex, found, JoinCtxErr(ctx, err)
		},
		o.AsIxMatch(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Any{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
					Pred:          pred.AsExpr(),
				}
			},
			o,
			pred,
		),
	)
}

// AllI returns an [PredicateI] that is satisfied if all focused satisfy the given predicate
//
// Note: If the optic focuses no elements then AllI returns true.
//
// See:
//   - [All] for a non index aware version.
//   - [AnyI] for a version that requires that any focuses satisfies the predicate.
func AllI[I, S, T, A, B, RET, RW, DIR, ERR, PERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], pred PredicateI[I, A, PERR]) Optic[Void, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, CompositionTree[ERR, PERR]] {
	return CombiGetter[CompositionTree[ERR, PERR], Void, S, S, bool, bool](
		func(ctx context.Context, source S) (Void, bool, error) {
			allMatch := true
			var err error
			o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
				focusIndex, v, focusErr := val.Get()
				if focusErr != nil {
					err = focusErr
					allMatch = false
					return false
				}
				ok, err := PredGetI(ctx, pred, v, focusIndex)
				err = JoinCtxErr(ctx, err)
				if err != nil {
					allMatch = false
					return false
				}

				if !ok {
					allMatch = false
					return false
				} else {
					return true
				}
			})
			return Void{}, allMatch, err
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.All{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
					Pred:          pred.AsExpr(),
				}
			},
			o,
			pred,
		),
	)
}
