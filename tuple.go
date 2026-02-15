package optic

import (
	"context"
	"reflect"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spearson78/go-optic/expr"
)

// SwappedT2 returns an [Iso] that focuses on the swapped values of a [lo.Tuple2]
func SwappedT2[A, B any]() Optic[Void, lo.Tuple2[A, B], lo.Tuple2[A, B], lo.Tuple2[B, A], lo.Tuple2[B, A], ReturnOne, ReadWrite, BiDir, Pure] {

	return Iso[lo.Tuple2[A, B], lo.Tuple2[B, A]](
		func(source lo.Tuple2[A, B]) lo.Tuple2[B, A] {
			return lo.T2(source.B, source.A)
		},
		func(source lo.Tuple2[B, A]) lo.Tuple2[A, B] {
			return lo.T2(source.B, source.A)
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.SwapTuple{
				OpticTypeExpr: ot,
				A:             reflect.TypeFor[A](),
				B:             reflect.TypeFor[B](),
			}
		}),
	)
}

// BesideT2 returns a [Traversal] that focuses all elements in both provided optics by returning an optic with a combined [lo.Tuple2] source and result types.
func BesideT2[IL, IR, SL, SR, TR, TL, A, B, RETL, RETR any, RWL, RWR TReadWrite, DIRL, DIRR any, LERR, RERR any](left Optic[IL, SL, TL, A, B, RETL, RWL, DIRL, LERR], right Optic[IR, SR, TR, A, B, RETR, RWR, DIRR, RERR]) Optic[mo.Either[IL, IR], lo.Tuple2[SL, SR], lo.Tuple2[TL, TR], A, B, ReturnMany, CompositionTree[RWL, RWR], UniDir, CompositionTree[LERR, RERR]] {
	return CombiTraversal[ReturnMany, CompositionTree[RWL, RWR], CompositionTree[LERR, RERR], mo.Either[IL, IR], lo.Tuple2[SL, SR], lo.Tuple2[TL, TR], A, B](
		func(ctx context.Context, source lo.Tuple2[SL, SR]) SeqIE[mo.Either[IL, IR], A] {
			return func(yield func(ValueIE[mo.Either[IL, IR], A]) bool) {
				cont := true
				left.AsIter()(ctx, source.A)(func(val ValueIE[IL, A]) bool {
					focusIndex, o, focusErr := val.Get()
					cont = yield(ValIE(mo.Left[IL, IR](focusIndex), o, JoinCtxErr(ctx, focusErr)))
					return cont
				})
				if cont {
					right.AsIter()(ctx, source.B)(func(val ValueIE[IR, A]) bool {
						focusIndex, o, focusErr := val.Get()
						cont = yield(ValIE(mo.Right[IL, IR](focusIndex), o, JoinCtxErr(ctx, focusErr)))
						return cont
					})
				}
			}
		},
		func(ctx context.Context, source lo.Tuple2[SL, SR]) (int, error) {
			leftLen, err := left.AsLengthGetter()(ctx, source.A)
			err = JoinCtxErr(ctx, err)
			if err != nil {
				return 0, err
			}

			rightLen, err := left.AsLengthGetter()(ctx, source.A)
			err = JoinCtxErr(ctx, err)
			if err != nil {
				return 0, err
			}

			return leftLen + rightLen, nil
		},
		func(ctx context.Context, fmap func(index mo.Either[IL, IR], focus A) (B, error), source lo.Tuple2[SL, SR]) (lo.Tuple2[TL, TR], error) {
			l, err := left.AsModify()(ctx, func(i IL, focus A) (B, error) {
				ret, err := fmap(mo.Left[IL, IR](i), focus)
				return ret, JoinCtxErr(ctx, err)
			}, source.A)
			err = JoinCtxErr(ctx, err)
			if err != nil {
				var ret lo.Tuple2[TL, TR]
				return ret, err
			}
			r, err := right.AsModify()(ctx, func(i IR, focus A) (B, error) {
				ret, err := fmap(mo.Right[IL, IR](i), focus)
				return ret, JoinCtxErr(ctx, err)
			}, source.B)
			err = JoinCtxErr(ctx, err)
			if err != nil {
				var ret lo.Tuple2[TL, TR]
				return ret, err
			}

			return lo.T2(l, r), ctx.Err()
		},
		func(ctx context.Context, index mo.Either[IL, IR], source lo.Tuple2[SL, SR]) SeqIE[mo.Either[IL, IR], A] {
			return func(yield func(ValueIE[mo.Either[IL, IR], A]) bool) {
				if il, ok := index.Left(); ok {
					left.AsIxGetter()(ctx, il, source.A)(func(val ValueIE[IL, A]) bool {
						focusIndex, o, focusErr := val.Get()
						return yield(ValIE(mo.Left[IL, IR](focusIndex), o, JoinCtxErr(ctx, focusErr)))
					})
				} else {
					right.AsIxGetter()(ctx, index.MustRight(), source.B)(func(val ValueIE[IR, A]) bool {
						focusIndex, o, focusErr := val.Get()
						return yield(ValIE(mo.Right[IL, IR](focusIndex), o, JoinCtxErr(ctx, focusErr)))
					})
				}
			}
		},
		func(indexA, indexB mo.Either[IL, IR]) bool {
			if leftA, ok := indexA.Left(); ok {
				if leftB, ok := indexB.Left(); ok {
					return left.AsIxMatch()(leftA, leftB)
				} else {
					return false
				}
			} else {
				if rightB, ok := indexB.Right(); ok {
					rightA := indexA.MustRight()
					return right.AsIxMatch()(rightA, rightB)
				} else {
					return false
				}
			}
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Beside{
					OpticTypeExpr: ot,
					Optic1:        left.AsExpr(),
					Optic2:        right.AsExpr(),
				}
			},
			left,
			right,
		),
	)
}

func traverseTNIter[A any](ctx context.Context, elements ...*A) SeqIE[int, A] {
	return func(yield func(ValueIE[int, A]) bool) {
		for i, v := range elements {
			if !yield(ValIE(i, *v, nil)) {
				break
			}
		}
	}
}

func traverseTNModify[A, B any](ctx context.Context, fmap func(index int, focus A) (B, error), elements ...*A) ([]B, error) {
	ret := make([]B, len(elements))

	for i, v := range elements {
		b, err := fmap(i, *v)
		err = JoinCtxErr(ctx, err)
		if err != nil {
			return ret, err
		}
		ret[i] = b
	}

	return ret, nil
}

func traverseTNIxGet[A any](index int, elements ...*A) SeqIE[int, A] {
	if index < 0 || index > len(elements)-1 {
		return func(yield func(ValueIE[int, A]) bool) {}
	}
	val := *elements[index]
	return func(yield func(ValueIE[int, A]) bool) {
		yield(ValIE(index, val, nil))
	}
}

func tnColGetter[A any](ctx context.Context, elements ...*A) Collection[int, A, Pure] {
	return ColIE[Pure, int, A](
		func(ctx context.Context) SeqIE[int, A] {
			return traverseTNIter(ctx, elements...)
		},
		func(ctx context.Context, index int) SeqIE[int, A] {
			return func(yield func(ValueIE[int, A]) bool) {
				if index >= 0 && index < len(elements) {
					yield(ValIE(index, *elements[index], nil))
				}
			}
		},
		IxMatchComparable[int](),
		func(ctx context.Context) (int, error) {
			return len(elements), nil
		},
	)
}

func tnColReverseGet[A any](ctx context.Context, focus Collection[int, A, Pure], elements ...*A) error {
	var retErr error
	i := 0
	focus.AsIter()(ctx)(func(val ValueIE[int, A]) bool {
		_, focus, err := val.Get()
		if err != nil {
			retErr = err
			return false
		}
		*elements[i] = focus
		i++
		return i <= len(elements)
	})
	return retErr
}
