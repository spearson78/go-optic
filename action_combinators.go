package optic

import (
	"context"

	"github.com/spearson78/go-optic/expr"
)

// AsModify is a combinator that returns an [Operation] that applies the op as a modification to the given [Optic]
//
// See:
//   - [AsModifyI] for an index aware version.
//   - [Modify] for an action equivalent.
func AsModify[I, S, T, A, B any, RET any, RW TReadWrite, DIR any, ORET TReturnOne, ERR, OERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], fmap Operation[A, B, ORET, OERR]) Optic[Void, S, S, T, T, ReturnOne, ReadOnly, UniDir, CompositionTree[ERR, OERR]] {
	return AsModifyI(o, OpToOpI[I](fmap))
}

// AsModifyI is an index aware combinator that returns an [Operation] that applies the op as a modification to the given [Optic]
//
// See:
//   - [AsModify] for a non index aware version.
//   - [ModifyI] for an action equivalent.
func AsModifyI[I, S, T, A, B any, RET any, RW TReadWrite, DIR any, ORET TReturnOne, ERR, OERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], fmap OperationI[I, A, B, ORET, OERR]) Optic[Void, S, S, T, T, ReturnOne, ReadOnly, UniDir, CompositionTree[ERR, OERR]] {
	return CombiGetter[CompositionTree[ERR, OERR], Void, S, S, T, T](
		func(ctx context.Context, source S) (Void, T, error) {
			ret, err := o.AsModify()(ctx, func(index I, focus A) (B, error) {
				return fmap.AsOpGet()(ctx, ValI(index, focus))
			}, source)
			return Void{}, ret, err
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Modify{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
					Fmap:          fmap.AsExpr(),
				}
			},
			o,
			fmap,
		),
	)
}

// AsSet is a combinator that returns an [Operation] that focuses the result of setting the given val using the given optic
//
// See:
//   - [Set] for an action equivalent.
func AsSet[I, S, T, A, B any, RET any, RW TReadWrite, DIR any, ORET TReturnOne, ERR, OERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], val Operation[S, B, ORET, OERR]) Optic[Void, S, S, T, T, ReturnOne, ReadOnly, UniDir, CompositionTree[ERR, OERR]] {
	return CombiGetter[CompositionTree[ERR, OERR], Void, S, S, T, T](
		func(ctx context.Context, source S) (Void, T, error) {
			setVal, err := val.AsOpGet()(ctx, source)
			if err != nil {
				var t T
				return Void{}, t, err
			}

			ret, err := o.AsSetter()(ctx, setVal, source)
			return Void{}, ret, err
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Set{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
					Val:           val,
				}
			},
			o,
			val,
		),
	)
}

// The AsReverseGet combinator reverses the direction of an [Iso].
//
// See:
//   - [Embed] for a version that works with [Prism]s
func AsReverseGet[I, S, T, A, B any, RET TReturnOne, RW TReadWrite, DIR TBiDir, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[Void, B, A, T, S, ReturnOne, ReadWrite, BiDir, ERR] {
	return unsafeAsReverseGet[ReadWrite, BiDir](o)
}

// The From combinator reverses the direction of an [Iso] or [Prism] optic.
//
// Warning: reversing a [Prism] may lead to errors being returned from ReverseGet
func unsafeAsReverseGet[RW, DIR any, I, S, T, A, B any, RET any, ORW TReadWrite, ODIR TBiDir, OERR any](o Optic[I, S, T, A, B, RET, ORW, ODIR, OERR]) Optic[Void, B, A, T, S, ReturnOne, RW, DIR, OERR] {

	implRevGetFnc := o.AsReverseGetter()
	implGetFnc := o.AsGetter()

	return Omni[Void, B, A, T, S, ReturnOne, RW, DIR, OERR](
		func(ctx context.Context, source B) (Void, T, error) {
			ret, err := implRevGetFnc(ctx, source)
			return Void{}, ret, err
		},
		func(ctx context.Context, focus S, source B) (A, error) {
			_, ret, err := implGetFnc(ctx, focus)
			return ret, err
		},
		func(ctx context.Context, source B) SeqIE[Void, T] {
			return func(yield func(ValueIE[Void, T]) bool) {
				v, err := implRevGetFnc(ctx, source)
				yield(ValIE(Void{}, v, err))
			}
		},
		func(ctx context.Context, source B) (int, error) {
			return 1, nil
		},
		func(ctx context.Context, fmap func(index Void, focus T) (S, error), source B) (A, error) {
			fa, err := implRevGetFnc(ctx, source)
			if err != nil {
				var t A
				return t, err
			}

			fb, err := fmap(Void{}, fa)
			err = JoinCtxErr(ctx, err)
			if err != nil {
				var t A
				return t, err
			}
			_, ft, err := implGetFnc(ctx, fb)
			return ft, JoinCtxErr(ctx, err)
		},
		func(ctx context.Context, index Void, source B) SeqIE[Void, T] {
			return func(yield func(ValueIE[Void, T]) bool) {
				ret, err := implRevGetFnc(ctx, source)
				err = JoinCtxErr(ctx, err)
				yield(ValIE(index, ret, err))
			}
		},
		IxMatchVoid(),
		func(ctx context.Context, focus S) (A, error) {
			_, ret, err := implGetFnc(ctx, focus)
			return ret, err
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.ReverseGet{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}
