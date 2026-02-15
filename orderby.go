package optic

import (
	"cmp"
	"context"
	"errors"
	"reflect"

	"github.com/samber/lo"
	"github.com/spearson78/go-optic/expr"
)

// The OrderByPredicate interface is a simplified read only version of the [Optic] interface used for sorting.
// OrderByPredicates should return true if the left value is less than the right value
type OrderByPredicate[A, ERR any] Predicate[lo.Tuple2[A, A], ERR]

// The OrderByPredicateI interface is a simplified read only version of the [Optic] interface used for sorting.
// OrderByIxPredicates should return true if the left value is less than the right value
type OrderByPredicateI[I, A, ERR any] Predicate[lo.Tuple2[ValueI[I, A], ValueI[I, A]], ERR]

// Constructor for an [OrderByPredicate] optic.
// This constructor is a convenience wrapper for [Operator] without an expression parameter.
//
// The less function should return true if valueLeft is less than valueRight.
//
// The following additional constructors are available.
//
//   - [OrderByOp] pure, non indexed
//   - [OrderByOpE] impure, non indexed
//   - [OrderByOpI] pure, index aware
//   - [OrderByOpIE] impure, index aware
//
// See:
//   - [OrderBy] for a combinator that converts an [Operation] to an [OrderByPredicate]
func OrderByOp[A any](less func(valueLeft A, valueRight A) bool) OrderByPredicate[A, Pure] {
	return Operator[lo.Tuple2[A, A], bool](
		func(source lo.Tuple2[A, A]) bool {
			return less(source.A, source.B)
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.OrderBy{
				OpticTypeExpr: ot,
				Optic: expr.GoFuncExpr{
					Func: reflect.ValueOf(less),
				},
			}
		}),
	)
}

// Constructor for an [OrderByPredicate] optic.
// This constructor is a convenience wrapper for [OperatorE] without an expression parameter.
//
// The less function should return true if valueLeft is less than valueRight.
//
// The following additional constructors are available.
//
//   - [OrderByOp] pure, non indexed
//   - [OrderByOpE] impure, non indexed
//   - [OrderByOpI] pure, index aware
//   - [OrderByOpIE] impure, index aware
//
// See:
//   - [OrderBy] for a combinator that converts an [Operation] to an [OrderByPredicate]
func OrderByOpE[A any](less func(ctx context.Context, valueLeft A, valueRight A) (bool, error)) Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool, ReturnOne, ReadOnly, UniDir, Err] {
	return OperatorE[lo.Tuple2[A, A], bool](
		func(ctx context.Context, source lo.Tuple2[A, A]) (bool, error) {
			return less(ctx, source.A, source.B)
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.OrderBy{
				OpticTypeExpr: ot,
				Optic: expr.GoFuncExpr{
					Func: reflect.ValueOf(less),
				},
			}
		}),
	)
}

// Constructor for an [OrderByPredicate] optic.
// This constructor is a convenience wrapper for [Operator] without an expression parameter.
//
// The less function should return true if valueLeft is less than valueRight.
//
// The following additional constructors are available.
//
//   - [OrderByOp] pure, non indexed
//   - [OrderByOpE] impure, non indexed
//   - [OrderByOpI] pure, index aware
//   - [OrderByOpIE] impure, index aware
//
// See:
//   - [OrderBy] for a combinator that converts an [Operation] to an [OrderByPredicate]
func OrderByOpI[I, A any](less func(indexLeft I, valueLeft A, indexRight I, valueRight A) bool) Optic[Void, lo.Tuple2[ValueI[I, A], ValueI[I, A]], lo.Tuple2[ValueI[I, A], ValueI[I, A]], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return Operator[lo.Tuple2[ValueI[I, A], ValueI[I, A]], bool](
		func(source lo.Tuple2[ValueI[I, A], ValueI[I, A]]) bool {
			return less(source.A.index, source.A.value, source.B.index, source.B.value)
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.OrderBy{
				OpticTypeExpr: ot,
				Optic: expr.GoFuncExpr{
					Func: reflect.ValueOf(less),
				},
			}
		}),
	)
}

// Constructor for an [OrderByPredicate] optic.
// This constructor is a convenience wrapper for [OperatorE] without an expression parameter.
//
// The less function should return true if valueLeft is less than valueRight.
//
// The following additional constructors are available.
//   - [OrderByOp] pure, non indexed
//   - [OrderByOpE] impure, non indexed
//   - [OrderByOpI] pure, index aware
//   - [OrderByOpIE] impure, index aware
//
// See:
//   - [OrderBy] for a combinator that converts an [Operation] to an [OrderByPredicate]
func OrderByOpIE[I, A any](less func(ctx context.Context, indexLeft I, valueLeft A, indexRight I, valueRight A) (bool, error)) Optic[Void, lo.Tuple2[ValueI[I, A], ValueI[I, A]], lo.Tuple2[ValueI[I, A], ValueI[I, A]], bool, bool, ReturnOne, ReadOnly, UniDir, Err] {
	return OperatorE[lo.Tuple2[ValueI[I, A], ValueI[I, A]], bool](
		func(ctx context.Context, source lo.Tuple2[ValueI[I, A], ValueI[I, A]]) (bool, error) {
			return less(ctx, source.A.index, source.A.value, source.B.index, source.B.value)
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.OrderBy{
				OpticTypeExpr: ot,
				Optic: expr.GoFuncExpr{
					Func: reflect.ValueOf(less),
				},
			}
		}),
	)
}

// The Desc combinator converts am OrderByOperation to operate in the reverse order.
//
// See:
//   - [DescI] for an index aware version.
func Desc[A, ERR any](lessPred OrderByPredicate[A, ERR]) Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return CombiGetter[ERR, Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool](
		func(ctx context.Context, source lo.Tuple2[A, A]) (Void, bool, error) {
			less, err := PredGet(ctx, lessPred, source)
			return Void{}, !less, err
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Desc{
					OpticTypeExpr: ot,

					Optic: lessPred.AsExpr(),
				}
			},
			lessPred,
		),
	)
}

// Constructor for an [OrderByPredicate] optic.
//
// The following additional constructors are available.
//   - [OrderByI] index aware
//
// The cmpBy [Operation] should focus on a [cmp.Ordered] value that will be used to order by.
//
// See:
//   - [OrderByOp] for a constructor that takes an arbitrary function to perform ordering.
//   - [OrderByN] for a combinator enables ordering on multiple criteria
func OrderBy[A any, C cmp.Ordered, RET TReturnOne, ERR any](cmpBy Operation[A, C, RET, ERR]) Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {

	cmpByFnc := cmpBy.AsOpGet()

	return CombiGetter[ERR, Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool](
		func(ctx context.Context, source lo.Tuple2[A, A]) (Void, bool, error) {

			aKey, err := cmpByFnc(ctx, source.A)
			if err != nil {
				return Void{}, false, err
			}

			bKey, err := cmpByFnc(ctx, source.B)

			return Void{}, aKey < bKey, err
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.OrderBy{
					OpticTypeExpr: ot,

					Optic: cmpBy.AsExpr(),
				}
			},
			cmpBy,
		),
	)
}

// Constructor for an [OrderByPredicateI] optic.
//
// The following additional constructors are available.
//   - [OrderBy] non index aware
//
// The cmpBy [OperationI] should focus on a [cmp.Ordered]] value that will be used to order by.
//
// See:
//   - [OrderByOpI] for a constructor that takes an arbitrary function to perform ordering.
//   - [OrderByNI] for a combinator enables ordering on multiple criteria
func OrderByI[I, A any, C cmp.Ordered, RET TReturnOne, ERR any](cmpBy OperationI[I, A, C, RET, ERR]) Optic[Void, lo.Tuple2[ValueI[I, A], ValueI[I, A]], lo.Tuple2[ValueI[I, A], ValueI[I, A]], bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return CombiGetter[ERR, Void, lo.Tuple2[ValueI[I, A], ValueI[I, A]], lo.Tuple2[ValueI[I, A], ValueI[I, A]], bool, bool](
		func(ctx context.Context, source lo.Tuple2[ValueI[I, A], ValueI[I, A]]) (Void, bool, error) {

			aKey, err := cmpBy.AsOpGet()(ctx, source.A)
			if err != nil {
				return Void{}, false, err
			}

			bKey, err := cmpBy.AsOpGet()(ctx, source.B)

			return Void{}, aKey < bKey, err
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.OrderBy{
					OpticTypeExpr: ot,

					Optic: cmpBy.AsExpr(),
				}
			},
			cmpBy,
		),
	)
}

// The OrderByN combinator combines multiple [OrderByPredicate] into a single [OrderByPredicate] applying each order by in sequence.
//
// See:
//   - [OrderBy2], [OrderBy3], [OrderBy4] for versions that allows the ERR to differ between orderBys
func OrderByN[A, ERR any](orderBys ...OrderByPredicate[A, ERR]) Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	naked := MustGet(
		SliceOf(
			Compose(
				TraverseSlice[OrderByPredicate[A, ERR]](),
				UpCast[OrderByPredicate[A, ERR], PredicateE[lo.Tuple2[A, A]]](),
			),
			len(orderBys),
		),
		orderBys,
	)

	return orderByN[ERR](naked...)
}

// The OrderBy2 combinator combines 2 [OrderByPredicate] into a single [OrderByPredicate] applying each order by in sequence.
//
// See:
//   - [OrderBy2], [OrderBy3], [OrderBy4] for other versions that allows the ERR to differ between orderBys
//   - [OrderByN] for a versions that accepts an arbitrary number of orderbys but requires they all have the same ERR
func OrderBy2[A, ERR1, ERR2 any](orderBy1 OrderByPredicate[A, ERR1], orderBy2 OrderByPredicate[A, ERR2]) Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool, ReturnOne, ReadOnly, UniDir, CompositionTree[ERR1, ERR2]] {
	return orderByN[CompositionTree[ERR1, ERR2]](PredicateE[lo.Tuple2[A, A]](orderBy1), PredicateE[lo.Tuple2[A, A]](orderBy2))
}

// The OrderBy3 combinator combines 3 [OrderByPredicate] into a single [OrderByPredicate] applying each order by in sequence.
//
// See:
//   - [OrderBy2], [OrderBy3], [OrderBy4] for other versions that allows the ERR to differ between orderBys
//   - [OrderByN] for a versions that accepts an arbitrary number of orderbys but requires they all have the same ERR
func OrderBy3[A, ERR1, ERR2, ERR3 any](orderBy1 OrderByPredicate[A, ERR1], orderBy2 OrderByPredicate[A, ERR2], orderBy3 OrderByPredicate[A, ERR3]) Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool, ReturnOne, ReadOnly, UniDir, CompositionTree[CompositionTree[ERR1, ERR2], ERR3]] {
	return orderByN[CompositionTree[CompositionTree[ERR1, ERR2], ERR3]](PredicateE[lo.Tuple2[A, A]](orderBy1), PredicateE[lo.Tuple2[A, A]](orderBy2), PredicateE[lo.Tuple2[A, A]](orderBy3))
}

// The OrderBy4 combinator combines 4 [OrderByPredicate] into a single [OrderByPredicate] applying each order by in sequence.
//
// See:
//   - [OrderBy2], [OrderBy3], [OrderBy4] for other versions that allows the ERR to differ between orderBys
//   - [OrderByN] for a versions that accepts an arbitrary number of orderbys but requires they all have the same ERR
func OrderBy4[A, ERR1, ERR2, ERR3, ERR4 any](orderBy1 OrderByPredicate[A, ERR1], orderBy2 OrderByPredicate[A, ERR2], orderBy3 OrderByPredicate[A, ERR3], orderBy4 OrderByPredicate[A, ERR4]) Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool, ReturnOne, ReadOnly, UniDir, CompositionTree[CompositionTree[ERR1, ERR2], CompositionTree[ERR3, ERR4]]] {
	return orderByN[CompositionTree[CompositionTree[ERR1, ERR2], CompositionTree[ERR3, ERR4]]](PredicateE[lo.Tuple2[A, A]](orderBy1), PredicateE[lo.Tuple2[A, A]](orderBy2), PredicateE[lo.Tuple2[A, A]](orderBy3), PredicateE[lo.Tuple2[A, A]](orderBy4))
}

func orderByN[ERR, A any](orderBys ...PredicateE[lo.Tuple2[A, A]]) Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {

	return CombiGetter[ERR, Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool](
		func(ctx context.Context, source lo.Tuple2[A, A]) (Void, bool, error) {

			for _, orderBy := range orderBys {
				less, err := orderBy.AsOpGet()(ctx, source)
				if less || err != nil {
					if errors.Is(ErrEmptyGet, err) {
						less = false
						err = nil
					}
					return Void{}, true, err
				}

				swapped := lo.T2(source.B, source.A)
				more, err := orderBy.AsOpGet()(ctx, swapped)
				if more || err != nil {
					if errors.Is(ErrEmptyGet, err) {
						more = true
						err = nil
					}
					return Void{}, !more, err
				}

			}

			//Equal
			return Void{}, false, nil
		},
		IxMatchVoid(),
		ExprDefVarArgs(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				var exprs []expr.OpticExpression

				for _, v := range orderBys {
					exprs = append(exprs, v.AsExpr())
				}

				return expr.OrderByN{
					OpticTypeExpr: ot,

					Optics: exprs,
				}
			},
			orderBys...,
		),
	)
}

// The DescI combinator converts an [OrderByPredicateI] to operate in the reverse order.
//
// See:
//   - [Desc] for a non index aware version.
func DescI[I, A, ERR any](less OrderByPredicateI[I, A, ERR]) Optic[Void, lo.Tuple2[ValueI[I, A], ValueI[I, A]], lo.Tuple2[ValueI[I, A], ValueI[I, A]], bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return CombiGetter[ERR, Void, lo.Tuple2[ValueI[I, A], ValueI[I, A]], lo.Tuple2[ValueI[I, A], ValueI[I, A]], bool, bool](
		func(ctx context.Context, source lo.Tuple2[ValueI[I, A], ValueI[I, A]]) (Void, bool, error) {
			less, err := PredGet(ctx, less, source)
			return Void{}, !less, err
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Desc{
					OpticTypeExpr: ot,

					Optic: less.AsExpr(),
				}
			},
			less,
		),
	)
}
