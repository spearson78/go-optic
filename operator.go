package optic

import (
	"context"
	"errors"
	"reflect"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spearson78/go-optic/expr"
)

// The Operation interface is a simplified read only version of the [Optic] interface used for operators and predicates. Provides access to the low level optic implementation functions. These functions are not intended for general user use but are provided to enable efficient combinators and actions to be implemented in external packages.
type Operation[S, A any, RET TReturnOne, ERR any] interface {
	AsOpGet() OpGetFunc[S, A]
	ReturnType() RET
	ErrType() ERR
	AsExprHandler() func(ctx context.Context) (ExprHandler, error)
	AsExpr() expr.OpticExpression
}

// The OperationI interface is a simplified read only version of the [Optic] interface used for operators and predicates. Provides access to the low level optic implementation functions. These functions are not intended for general user use but are provided to enable efficient combinators and actions to be implemented in external packages.
type OperationI[I, S, A any, RET TReturnOne, ERR any] Operation[ValueI[I, S], A, RET, ERR]

// Constructor for an [Operator] optic.
// This constructor is a convenience wrapper for [Operator] without an expression parameter.
//
// The following additional constructors are available.
//
//	. [Operator] for a version that has an expression parameter.
//	- [Op] pure, non indexed,non polymorphic
//	- [OpI] pure, index aware,non polymorphic
//	- [OpE] impure, non indexed,non polymorphic
//	- [OpIE] impure, index aware, non polymorphic
func Op[S, A any](op func(source S) A) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return Operator(op, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.GoFuncExpr{
			OpticTypeExpr: ot,
			Func:          reflect.ValueOf(op),
		}
	}))
}

// Constructor for an [Operator] optic.
// This constructor is a convenience wrapper for [Operator] without an expression parameter.
//
// The following additional constructors are available.
//
//	. [OperatorE] for a version that has an expression parameter.
//	- [Op] pure, non indexed,non polymorphic
//	- [OpI] pure, index aware,non polymorphic
//	- [OpE] impure, non indexed,non polymorphic
//	- [OpIE] impure, index aware, non polymorphic
func OpE[S, A any](op func(ctx context.Context, source S) (A, error)) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, Err] {
	return OperatorE(op, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.GoFuncExpr{
			OpticTypeExpr: ot,
			Func:          reflect.ValueOf(op),
		}
	}))
}

// Constructor for an [Operator] optic.
// This constructor is a convenience wrapper for [Operator] without an expression parameter.
//
// The following additional constructors are available.
//
//	. [OperatorI] for a version that has an expression parameter.
//	- [Op] pure, non indexed,non polymorphic
//	- [OpI] pure, index aware,non polymorphic
//	- [OpE] impure, non indexed,non polymorphic
//	- [OpIE] impure, index aware, non polymorphic
func OpI[I, S, A any](op func(I, S) A) Optic[Void, ValueI[I, S], ValueI[I, S], A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return OperatorI(op, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.GoFuncExpr{
			OpticTypeExpr: ot,
			Func:          reflect.ValueOf(op),
		}
	}))
}

// Constructor for an [Operator] optic.
// This constructor is a convenience wrapper for [Operator] without an expression parameter.
//
// The following additional constructors are available.
//
//	. [OperatorE] for a version that has an expression parameter.
//	- [Op] pure, non indexed,non polymorphic
//	- [OpI] pure, index aware,non polymorphic
//	- [OpE] impure, non indexed,non polymorphic
//	- [OpIE] impure, index aware, non polymorphic
func OpIE[I, S, A any](opFnc func(context.Context, I, S) (A, error)) Optic[Void, ValueI[I, S], ValueI[I, S], A, A, ReturnOne, ReadOnly, UniDir, Err] {
	return OperatorIE[I, S, A, Err](
		opFnc,
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.GoFuncExpr{
				OpticTypeExpr: ot,
				Func:          reflect.ValueOf(opFnc),
			}
		}),
	)

}

// Constructor for an [Operator] optic that takes 2 parameters.
// This constructor is a convenience wrapper for [Operator] without an expression parameter.
//
// The following additional constructors are available.
//
//	. [Op] for a version that takes a single parameter
//	- [OpT2] pure, non indexed,non polymorphic
//	- [OpT2I] pure, index aware,non polymorphic
//	- [OpT2E] impure, non indexed,non polymorphic
//	- [OpT2IE] impure, index aware, non polymorphic
func OpT2[S, T, A any](op func(a S, b T) A) Optic[Void, lo.Tuple2[S, T], lo.Tuple2[S, T], A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return Operator(func(source lo.Tuple2[S, T]) A { return op(source.A, source.B) }, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.GoFuncExpr{
			OpticTypeExpr: ot,
			Func:          reflect.ValueOf(op),
		}
	}))
}

// Constructor for an [Operator] optic that takes 2 parameters.
// This constructor is a convenience wrapper for [Operator] without an expression parameter.
//
// The following additional constructors are available.
//
//	. [OpE] for a version that takes a single parameters.
//	- [OpT2] pure, non indexed,non polymorphic
//	- [OpT2I] pure, index aware,non polymorphic
//	- [OpT2E] impure, non indexed,non polymorphic
//	- [OpT2IE] impure, index aware, non polymorphic
func OpT2E[S, T, A any](op func(ctx context.Context, a S, b T) (A, error)) Optic[Void, lo.Tuple2[S, T], lo.Tuple2[S, T], A, A, ReturnOne, ReadOnly, UniDir, Err] {
	return OperatorE(func(ctx context.Context, source lo.Tuple2[S, T]) (A, error) {
		return op(ctx, source.A, source.B)
	}, ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.GoFuncExpr{
			OpticTypeExpr: ot,
			Func:          reflect.ValueOf(op),
		}
	}))
}

// Constructor for an [Operator] optic that takes 2 parameters.
// This constructor is a convenience wrapper for [Operator] without an expression parameter.
//
// The following additional constructors are available.
//
//	. [OpI] for a version that takes a single parameters.
//	- [OpT2] pure, non indexed,non polymorphic
//	- [OpT2I] pure, index aware,non polymorphic
//	- [OpT2E] impure, non indexed,non polymorphic
//	- [OpT2IE] impure, index aware, non polymorphic
func OpT2I[I, J, S, T, A any](op func(ia I, a S, ib J, b T) A) Optic[Void, lo.Tuple2[ValueI[I, S], ValueI[J, T]], lo.Tuple2[ValueI[I, S], ValueI[J, T]], A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return CombiGetter[Pure, Void, lo.Tuple2[ValueI[I, S], ValueI[J, T]], lo.Tuple2[ValueI[I, S], ValueI[J, T]], A, A](
		func(ctx context.Context, source lo.Tuple2[ValueI[I, S], ValueI[J, T]]) (Void, A, error) {
			ret := op(source.A.index, source.A.value, source.B.index, source.B.value)
			return Void{}, ret, nil
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.GoFuncExpr{
				OpticTypeExpr: ot,
				Func:          reflect.ValueOf(op),
			}
		}),
	)
}

// Constructor for an [Operator] optic that takes 2 parameters.
// This constructor is a convenience wrapper for [Operator] without an expression parameter.
//
// The following additional constructors are available.
//
//	. [OpIE] for a version that takes a single parameter.
//	- [OpT2] pure, non indexed,non polymorphic
//	- [OpT2I] pure, index aware,non polymorphic
//	- [OpT2E] impure, non indexed,non polymorphic
//	- [OpT2IE] impure, index aware, non polymorphic
func OpT2IE[I, J, S, T, A any](opFnc func(ctx context.Context, i I, a S, j J, b T) (A, error)) Optic[Void, lo.Tuple2[ValueI[I, S], ValueI[J, T]], lo.Tuple2[ValueI[I, S], ValueI[J, T]], A, A, ReturnOne, ReadOnly, UniDir, Err] {
	return GetterIE[Void, lo.Tuple2[ValueI[I, S], ValueI[J, T]], A](
		func(ctx context.Context, source lo.Tuple2[ValueI[I, S], ValueI[J, T]]) (Void, A, error) {
			ret, err := opFnc(ctx, source.A.index, source.A.value, source.B.index, source.B.value)
			return Void{}, ret, err
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.GoFuncExpr{
				OpticTypeExpr: ot,
				Func:          reflect.ValueOf(opFnc),
			}
		}),
	)
}

// Constructor for an [Operator] optic.
//
// The following additional constructors are available.
//
//	. [Op] for a version that does not require an [expr.OpticExpression],
//	- [Operator] pure, non indexed,non polymorphic
//	- [OperatorI] pure, index aware,non polymorphic
//	- [OperatorE] impure, non indexed,non polymorphic
//	- [OperatorIE] impure, index aware, non polymorphic
func Operator[S, A any](
	opFnc func(source S) A,
	exprDef ExpressionDef,
) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return CombiGetter[Pure, Void, S, S, A, A](
		func(ctx context.Context, source S) (Void, A, error) {
			return Void{}, opFnc(source), ctx.Err()
		},
		IxMatchVoid(),
		exprDef,
	)
}

// Constructor for an [Operator] optic.
//
// The following additional constructors are available.
//
//	. [OpI] for a version that does not require an [expr.OpticExpression],
//	- [Operator] pure, non indexed,non polymorphic
//	- [OperatorI] pure, index aware,non polymorphic
//	- [OperatorE] impure, non indexed,non polymorphic
//	- [OperatorIE] impure, index aware, non polymorphic
func OperatorI[I, S, A any](
	opFnc func(I, S) A,
	exprDef ExpressionDef,
) Optic[Void, ValueI[I, S], ValueI[I, S], A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return CombiGetter[Pure, Void, ValueI[I, S], ValueI[I, S], A, A](
		func(ctx context.Context, source ValueI[I, S]) (Void, A, error) {
			a := opFnc(source.index, source.value)
			return Void{}, a, nil
		},
		IxMatchVoid(),
		exprDef,
	)
}

// Constructor for an [Operator] optic.
//
// The following additional constructors are available.
//
//	. [OpE] for a version that does not require an [expr.OpticExpression],
//	- [Operator] pure, non indexed,non polymorphic
//	- [OperatorI] pure, index aware,non polymorphic
//	- [OperatorE] impure, non indexed,non polymorphic
//	- [OperatorIE] impure, index aware, non polymorphic
func OperatorE[S, A any](
	opFnc func(ctx context.Context, source S) (A, error),
	exprDef ExpressionDef,
) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, Err] {
	return CombiGetter[Err, Void, S, S, A, A](
		func(ctx context.Context, source S) (Void, A, error) {
			ret, err := opFnc(ctx, source)
			return Void{}, ret, err
		},
		IxMatchVoid(),
		exprDef,
	)
}

// Constructor for an [Operator] optic.
//
// The following additional constructors are available.
//
//	. [OpIE] for a version that does not require an [expr.OpticExpression],
//	- [Operator] pure, non indexed,non polymorphic
//	- [OperatorI] pure, index aware,non polymorphic
//	- [OperatorE] impure, non indexed,non polymorphic
//	- [OperatorIE] impure, index aware, non polymorphic
func OperatorIE[I, S, A, ERR any](
	opFnc func(ctx context.Context, index I, source S) (A, error),
	exprDef ExpressionDef,
) Optic[Void, ValueI[I, S], ValueI[I, S], A, A, ReturnOne, ReadOnly, UniDir, ERR] {
	return CombiGetter[ERR, Void, ValueI[I, S], ValueI[I, S], A, A](
		func(ctx context.Context, source ValueI[I, S]) (Void, A, error) {
			a, err := opFnc(ctx, source.index, source.value)
			return Void{}, a, err
		},
		IxMatchVoid(),
		exprDef,
	)
}

// OpToOpI is a combinator that converts an [Operator] to an [OperatorI] that ignores the index when performing the operation.
//
// See:
//   - [OpOnIx] for a combinator that converts an [Operator] to operate on the index instead of the focus
func OpToOpI[I, S, A any, RET TReturnOne, ERR any](fnc Operation[S, A, RET, ERR]) Optic[Void, ValueI[I, S], ValueI[I, S], A, A, ReturnOne, ReadOnly, UniDir, ERR] {
	return Ret1(Ro(Ud(EErrR(Compose(
		ValueIValue[I, S](),
		OpToOptic(fnc),
	)))))
}

// PredT2ToOpT2I is a combinator that converts a 2 parameter [Predicate] to a 2 parameter [PredicateI] that ignores the index when performing the operation.
//
// See:
//   - [PredToOpI] for a single parameter version
func PredT2ToOpT2I[I, J, S, T, ERR any](fnc Predicate[lo.Tuple2[S, T], ERR]) Optic[Void, lo.Tuple2[ValueI[I, S], ValueI[J, T]], lo.Tuple2[ValueI[I, S], ValueI[J, T]], bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	opGetFnc := fnc.AsOpGet()
	return rawGetterF[Void, lo.Tuple2[ValueI[I, S], ValueI[J, T]], bool, ERR](
		func(ctx context.Context, source lo.Tuple2[ValueI[I, S], ValueI[J, T]]) (Void, bool, error) {
			ret, err := opGetFnc(ctx, lo.T2(source.A.value, source.B.value))
			if errors.Is(err, ErrEmptyGet) {
				return Void{}, false, nil
			}
			return Void{}, ret, err
		},
		IxMatchVoid(),
		fnc.AsExprHandler(),
		fnc.AsExpr,
	)
}

// PredToOpI is a combinator that converts a [Predicate] to an [OperatorI] that ignores the index when performing the operation.
//
// See:
//   - [OpOnIx] for a combinator that converts an [Operator] to operate on the index instead of the focus
func PredToOpI[I, S, ERR any](fnc Predicate[S, ERR]) Optic[Void, ValueI[I, S], ValueI[I, S], bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return Ret1(Ro(Ud(EErrR(Compose(
		ValueIValue[I, S](),
		PredToOptic(fnc),
	)))))
}

func OpT2ToOpT2I[J, S, A any, RET TReturnOne, ERR any](op Operation[lo.Tuple2[S, S], A, RET, ERR]) Optic[Void, lo.Tuple2[ValueI[J, S], ValueI[J, S]], lo.Tuple2[ValueI[J, S], ValueI[J, S]], A, A, ReturnOne, ReadOnly, UniDir, ERR] {
	opGetFnc := op.AsOpGet()
	return rawGetterF[Void, lo.Tuple2[ValueI[J, S], ValueI[J, S]], A, ERR](
		func(ctx context.Context, source lo.Tuple2[ValueI[J, S], ValueI[J, S]]) (Void, A, error) {
			ret, err := opGetFnc(ctx, lo.T2(source.A.value, source.B.value))
			return Void{}, ret, err
		},
		IxMatchVoid(),
		op.AsExprHandler(),
		op.AsExpr,
	)
}

// AsIxGet is a combinator that converts a [Traversal] to an [OperationI] that operates on the combined index and source.
//
// See:
//   - [OpToOpI] for a combinator that converts an [Operator] and ignores the index
//   - [OpOnIx] for a version that operates on on the index only.
func AsIxGet[I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[Void, ValueI[I, S], ValueI[I, S], mo.Option[A], mo.Option[A], ReturnOne, ReadOnly, UniDir, ERR] {

	return CombiGetter[ERR, Void, ValueI[I, S], ValueI[I, S], mo.Option[A], mo.Option[A]](
		func(ctx context.Context, source ValueI[I, S]) (Void, mo.Option[A], error) {

			var ret A
			var retErr error
			found := false

			o.AsIxGetter()(ctx, source.index, source.value)(func(val ValueIE[I, A]) bool {
				_, focus, err := val.Get()
				ret = focus
				retErr = err
				found = true
				return false
			})

			return Void{}, mo.TupleToOption(ret, found), JoinCtxErr(ctx, retErr)
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.OpIxGet{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)

}

// OpOnIx is a combinator that converts an [Operator] to an [IxOperator} that operates on the index instead of the focus.
//
// See:
//   - [OpToOpI] for a combinator that converts an [Operator] and ignores the index
func OpOnIx[A, I, R any, RET TReturnOne, ERR any](o Operation[I, R, RET, ERR]) Optic[Void, ValueI[I, A], ValueI[I, A], R, R, ReturnOne, ReadOnly, UniDir, ERR] {
	return Ret1(Ro(Ud(EErrR(Compose(
		ValueIIndex[I, A](),
		OpToOptic(o),
	)))))
}

// PredOnIx is a combinator that converts an [Operator] to an [IxOperator} that operates on the index instead of the focus.
func PredOnIx[A, I, ERR any](o Predicate[I, ERR]) Optic[Void, ValueI[I, A], ValueI[I, A], bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return Ret1(Ro(Ud(EErrR(Compose(
		ValueIIndex[I, A](),
		PredToOptic(o),
	)))))
}

// OpT2IToOptic returns an [Operation] that applies the given [OperationI] fnc on the values focused by the left and right optics.
func OpT2IToOptic[I, S, A, R, LRW, RRW, ORW any, LRET, RRET, ORET TReturnOne, LDIR, RDIR, ODIR any, LERR, RERR, OERR any](left Optic[I, S, S, A, A, LRET, LRW, LDIR, LERR], op Optic[Void, lo.Tuple2[ValueI[I, A], ValueI[I, A]], lo.Tuple2[ValueI[I, A], ValueI[I, A]], R, R, ORET, ORW, ODIR, OERR], right Optic[I, S, S, A, A, RRET, RRW, RDIR, RERR]) Optic[Void, S, S, R, R, ORET, ReadOnly, UniDir, CompositionTree[CompositionTree[LERR, RERR], OERR]] {
	return RetR(Ro(Ud(Compose(
		T2Of(
			WithIndex(left),
			WithIndex(right),
		),
		op,
	))))
}

func OpT2ToOptic[I, S, A, R, ORW any, LRET, RRET TReturnOne, ORET any, ODIR any, LERR, RERR, OERR any](left Operation[S, A, LRET, LERR], op Optic[I, lo.Tuple2[A, A], lo.Tuple2[A, A], R, R, ORET, ORW, ODIR, OERR], right Operation[S, A, RRET, RERR]) Optic[I, S, S, R, R, ORET, ReadOnly, UniDir, CompositionTree[CompositionTree[LERR, RERR], OERR]] {
	return RetR(Ro(Ud(Compose(
		T2Of(
			OpToOptic(left),
			OpToOptic(right),
		),
		op,
	))))
}

func OpT2Bind[A, B, R any, RET TReturnOne, ERR any](op Operation[lo.Tuple2[A, B], R, RET, ERR], rightVal B) Optic[Void, A, A, R, R, ReturnOne, ReadOnly, UniDir, ERR] {
	return CombiGetter[ERR, Void, A, A, R, R](
		func(ctx context.Context, leftVal A) (Void, R, error) {
			ret, err := op.AsOpGet()(ctx, lo.T2(leftVal, rightVal))
			return Void{}, ret, err
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.OpT2BindExpr{
				OpticTypeExpr: ot,
				Op:            op.AsExpr(),
				S:             reflect.TypeFor[A](),
				RightValue:    rightVal,
			}

		},
			op,
		),
	)
}

func OpT2IsoBind[A any, RET TReturnOne, ERR any, RRET TReturnOne, RERR any](op Operation[lo.Tuple2[A, A], A, RET, ERR], revOp Operation[lo.Tuple2[A, A], A, RRET, RERR], rightVal A) Optic[Void, A, A, A, A, ReturnOne, ReadWrite, BiDir, CompositionTree[ERR, RERR]] {
	return CombiIso[ReadWrite, BiDir, CompositionTree[ERR, RERR], A, A, A, A](
		func(ctx context.Context, leftVal A) (A, error) {
			ret, err := op.AsOpGet()(ctx, lo.T2(leftVal, rightVal))
			return ret, err
		},
		func(ctx context.Context, leftVal A) (A, error) {
			ret, err := revOp.AsOpGet()(ctx, lo.T2(leftVal, rightVal))
			return ret, err
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.IsoOpT2BindExpr{
				OpticTypeExpr: ot,
				Op:            op.AsExpr(),
				InvOp:         revOp.AsExpr(),
				S:             reflect.TypeFor[A](),
				Right:         rightVal,
			}

		},
			op,
			revOp,
		),
	)
}

func OpToOptic[S, A any, RET TReturnOne, ERR any](o Operation[S, A, RET, ERR]) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, ERR] {
	fnc := o.AsOpGet()
	return rawGetterF[Void, S, A, ERR](
		func(ctx context.Context, source S) (Void, A, error) {
			v, err := fnc(ctx, source)
			return Void{}, v, err
		},
		IxMatchVoid(),
		o.AsExprHandler(),
		o.AsExpr,
	)
}

func OpIToOptic[I, S, A any, RET TReturnOne, ERR any, ERRP any](o OperationI[I, S, A, RET, ERR], ixMatch Predicate[lo.Tuple2[I, I], ERRP]) Optic[I, ValueI[I, S], ValueI[I, S], A, A, ReturnOne, ReadOnly, UniDir, CompositionTree[ERR, ERRP]] {
	fnc := o.AsOpGet()
	return rawGetterF[I, ValueI[I, S], A, CompositionTree[ERR, ERRP]](
		func(ctx context.Context, source ValueI[I, S]) (I, A, error) {
			v, err := fnc(ctx, source)
			return source.index, v, err
		},
		func(indexA, indexB I) bool {
			return Must(PredGet(context.Background(), ixMatch, lo.T2(indexA, indexB)))
		},
		o.AsExprHandler(),
		o.AsExpr,
	)
}
