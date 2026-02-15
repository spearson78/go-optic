package otree

import (
	"context"
	"reflect"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

// WithChildPath returns a Traversal that reindexes a [Collection] to a [PathNode]
//
// See:
//   - [WithChildPathI] for a version that supports arbitrary index types
func WithChildPath[I comparable, S, A any, RET any, RW, DIR, ERR any](o Optic[*PathNode[I], S, S, Collection[I, A, ERR], Collection[I, A, ERR], RET, RW, DIR, ERR]) Optic[*PathNode[I], S, S, Collection[*PathNode[I], A, ERR], Collection[*PathNode[I], A, ERR], RET, RW, UniDir, ERR] {
	return WithChildPathI(o, IxMatchComparable[I]())
}

// WithChildPath returns a Traversal that reindexes a [Collection] to a [PathNode]
//
// See:
//   - [WithChildPath] for a simpler version that supports only comparable index types.
func WithChildPathI[I, S, A any, RET any, RW, DIR, ERR any](o Optic[*PathNode[I], S, S, Collection[I, A, ERR], Collection[I, A, ERR], RET, RW, DIR, ERR], ixMatch func(a, b I) bool) Optic[*PathNode[I], S, S, Collection[*PathNode[I], A, ERR], Collection[*PathNode[I], A, ERR], RET, RW, UniDir, ERR] {

	withPath := func(path *PathNode[I], col Collection[I, A, ERR]) Collection[*PathNode[I], A, ERR] {
		return ColIE[ERR, *PathNode[I], A](
			func(ctx context.Context) SeqIE[*PathNode[I], A] {
				return func(yield func(ValueIE[*PathNode[I], A]) bool) {
					col.AsIter()(ctx)(func(val ValueIE[I, A]) bool {
						i, focus, err := val.Get()
						return yield(ValIE(path.Append(i), focus, err))
					})
				}
			},
			func(ctx context.Context, index *PathNode[I]) SeqIE[*PathNode[I], A] {
				return func(yield func(ValueIE[*PathNode[I], A]) bool) {
					col.AsIxGet()(ctx, index.value)(func(val ValueIE[I, A]) bool {
						i, focus, err := val.Get()
						return yield(ValIE(path.Append(i), focus, err))
					})
				}
			},
			func(a, b *PathNode[I]) bool {
				return MustGet(EqPathT2(OpT2(col.AsIxMatch())), lo.T2(a, b))
			},
			col.AsLengthGetter(),
		)
	}

	return CombiTraversal[RET, RW, ERR, *PathNode[I], S, S, Collection[*PathNode[I], A, ERR], Collection[*PathNode[I], A, ERR]](
		func(ctx context.Context, source S) SeqIE[*PathNode[I], Collection[*PathNode[I], A, ERR]] {
			return func(yield func(ValueIE[*PathNode[I], Collection[*PathNode[I], A, ERR]]) bool) {
				o.AsIter()(ctx, source)(func(val ValueIE[*PathNode[I], Collection[I, A, ERR]]) bool {
					path, col, err := val.Get()
					if err != nil {
						return yield(ValIE[*PathNode[I], Collection[*PathNode[I], A, ERR]](path, nil, err))
					}
					return yield(ValIE(path, withPath(path, col), nil))

				})
			}
		},
		func(ctx context.Context, source S) (int, error) {
			return o.AsLengthGetter()(ctx, source)
		},
		func(ctx context.Context, fmap func(index *PathNode[I], focus Collection[*PathNode[I], A, ERR]) (Collection[*PathNode[I], A, ERR], error), source S) (S, error) {
			return o.AsModify()(ctx, func(index *PathNode[I], focus Collection[I, A, ERR]) (Collection[I, A, ERR], error) {

				mapped, err := fmap(index, withPath(index, focus))
				if err != nil {
					return nil, err
				}

				col := ColIE[ERR, I, A](
					func(ctx context.Context) SeqIE[I, A] {
						return func(yield func(ValueIE[I, A]) bool) {
							mapped.AsIter()(ctx)(func(val ValueIE[*PathNode[I], A]) bool {
								path, focus, err := val.Get()
								if path == nil {
									var i I
									return yield(ValIE(i, focus, err))
								} else {
									return yield(ValIE(path.Value(), focus, err))
								}
							})
						}
					},
					nil,
					ixMatch,
					focus.AsLengthGetter(),
				)

				return col, nil

			}, source)
		},
		func(ctx context.Context, index *PathNode[I], source S) SeqIE[*PathNode[I], Collection[*PathNode[I], A, ERR]] {
			return func(yield func(ValueIE[*PathNode[I], Collection[*PathNode[I], A, ERR]]) bool) {
				o.AsIxGetter()(ctx, index, source)(func(val ValueIE[*PathNode[I], Collection[I, A, ERR]]) bool {
					path, col, err := val.Get()
					if err != nil {
						return yield(ValIE[*PathNode[I], Collection[*PathNode[I], A, ERR]](path, nil, err))
					}
					return yield(ValIE(path, withPath(path, col), nil))

				})
			}
		},
		func(a, b *PathNode[I]) bool {
			return MustGet(EqPathT2(OpT2(ixMatch)), lo.T2(a, b))
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.WithChildPath{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
					I:             reflect.TypeFor[I](),
					A:             reflect.TypeFor[A](),
				}
			},
			o,
		),
	)
}
