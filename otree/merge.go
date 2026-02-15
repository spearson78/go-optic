package otree

import (
	"context"

	"github.com/samber/lo"
	"github.com/samber/mo"
	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

func MergeTree[I, J, A, RET, SRET any, RW any, SRW ReadWrite, DIR, SDIR, ERR, SERR any](right A, children Optic[I, A, A, A, A, RET, RW, DIR, ERR], childrenCol Optic[J, A, A, Collection[I, A, SERR], Collection[I, A, SERR], SRET, SRW, SDIR, SERR]) Optic[Void, mo.Option[A], mo.Option[A], mo.Option[A], mo.Option[A], ReturnOne, ReadOnly, UniDir, CompositionTree[ERR, SERR]] {
	return Ret1(Ro(EErrR(Compose(
		T2Of(
			Identity[mo.Option[A]](),
			Const[mo.Option[A]](mo.Some(right)),
		),
		MergeTreeT2(children, childrenCol),
	))))
}

func MergeTreeT2[I, J, A, RET, SRET any, RW any, SRW ReadWrite, DIR, SDIR, ERR, SERR any](children Optic[I, A, A, A, A, RET, RW, DIR, ERR], childrenCol Optic[J, A, A, Collection[I, A, SERR], Collection[I, A, SERR], SRET, SRW, SDIR, SERR]) Optic[Void, lo.Tuple2[mo.Option[A], mo.Option[A]], lo.Tuple2[mo.Option[A], mo.Option[A]], mo.Option[A], mo.Option[A], ReturnOne, ReadOnly, UniDir, CompositionTree[ERR, SERR]] {

	var mergeGet func(ctx context.Context, left, right A) (A, error)
	mergeGet = func(ctx context.Context, left, right A) (A, error) {
		return childrenCol.AsModify()(ctx, func(index J, focus Collection[I, A, SERR]) (Collection[I, A, SERR], error) {
			return ColIE[SERR](
				func(ctx context.Context) SeqIE[I, A] {
					return func(yield func(val ValueIE[I, A]) bool) {

						children.AsIter()(ctx, left)(func(val ValueIE[I, A]) bool {
							index, lVal, err := val.Get()
							if err != nil {
								var a A
								return yield(ValIE(index, a, err))
							}
							rFound := false
							var rErr error

							children.AsIxGetter()(ctx, index, right)(func(val ValueIE[I, A]) bool {
								err := val.Error()
								rFound = true
								rErr = err
								return false
							})

							if rErr != nil {
								var a A
								return yield(ValIE(index, a, rErr))
							}

							if !rFound {
								return yield(ValIE(index, lVal, rErr))
							}

							return true
						})

						focus.AsIter()(ctx)(func(val ValueIE[I, A]) bool {
							index, rVal, err := val.Get()
							if err != nil {
								var a A
								return yield(ValIE(index, a, err))
							}

							var lVal A
							lFound := false
							var lErr error
							children.AsIxGetter()(ctx, index, left)(func(val ValueIE[I, A]) bool {
								_, focus, err := val.Get()
								lVal = focus
								lErr = err
								lFound = true
								return false
							})

							if lErr != nil {
								var a A
								return yield(ValIE(index, a, err))
							}

							if lFound {
								mergedChild, err := mergeGet(ctx, lVal, rVal)
								if err != nil {
									var a A
									return yield(ValIE(index, a, err))
								}
								return yield(ValIE(index, mergedChild, nil))
							} else {
								return yield(ValIE(index, rVal, nil))
							}

						})
					}
				},
				nil,
				children.AsIxMatch(),
				nil,
			), nil
		}, right)
	}

	return CombiGetter[CompositionTree[ERR, SERR], Void, lo.Tuple2[mo.Option[A], mo.Option[A]], lo.Tuple2[mo.Option[A], mo.Option[A]], mo.Option[A], mo.Option[A]](
		func(ctx context.Context, source lo.Tuple2[mo.Option[A], mo.Option[A]]) (Void, mo.Option[A], error) {

			if a, ok := source.A.Get(); ok {
				if b, ok := source.B.Get(); ok {
					ret, err := mergeGet(ctx, a, b)
					return Void{}, mo.Some(ret), err
				} else {
					return Void{}, source.A, nil
				}
			} else {
				return Void{}, source.B, nil
			}

		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.TreeMerge{
					OpticTypeExpr: ot,
					Children:      children.AsExpr(),
					ChildrenSeq:   childrenCol.AsExpr(),
				}
			},
			children,
			childrenCol,
		),
	)
}
