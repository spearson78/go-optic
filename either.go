package optic

import (
	"context"
	"reflect"

	"github.com/samber/mo"
	"github.com/spearson78/go-optic/expr"
)

type ChosenSide bool

const (
	//Indicates the [Chosen] was left.
	ChosenLeft ChosenSide = false
	//Indicates the [Chosen] was right.
	ChosenRight ChosenSide = true
)

// Chosen returns a [Lens] that focuses on the present element of a [mo.Either]. The index indicates which side was present.
func Chosen[A any]() Optic[ChosenSide, mo.Either[A, A], mo.Either[A, A], A, A, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, ChosenSide, mo.Either[A, A], mo.Either[A, A], A, A](
		func(ctx context.Context, source mo.Either[A, A]) (ChosenSide, A, error) {
			if l, ok := source.Left(); ok {
				return ChosenLeft, l, ctx.Err()
			} else {
				return ChosenRight, source.MustRight(), ctx.Err()
			}
		},
		func(ctx context.Context, focus A, source mo.Either[A, A]) (mo.Either[A, A], error) {
			if _, ok := source.Left(); ok {
				return mo.Left[A, A](focus), ctx.Err()
			} else {
				return mo.Right[A, A](focus), ctx.Err()
			}
		},
		IxMatchComparable[ChosenSide](),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Chosen{
				OpticTypeExpr: ot,
			}
		}),
	)
}

// Left returns a [Prism] that matches on the left element of a [mo.Either]
func Left[A, B any]() Optic[ChosenSide, mo.Either[A, B], mo.Either[A, B], A, A, ReturnMany, ReadWrite, BiDir, Pure] {
	return LeftP[A, B, A]()
}

func LeftP[A, B, C any]() Optic[ChosenSide, mo.Either[A, B], mo.Either[C, B], A, C, ReturnMany, ReadWrite, BiDir, Pure] {
	return PrismIP[ChosenSide, mo.Either[A, B], mo.Either[C, B], A, C](
		func(source mo.Either[A, B]) mo.Either[mo.Either[C, B], ValueI[ChosenSide, A]] {
			if a, ok := source.Left(); ok {
				return mo.Right[mo.Either[C, B], ValueI[ChosenSide, A]](ValI(ChosenLeft, a))
			} else {
				return mo.Left[mo.Either[C, B], ValueI[ChosenSide, A]](mo.Right[C, B](source.MustRight()))
			}
		},
		func(focus C) mo.Either[C, B] {
			return mo.Left[C, B](focus)
		},
		IxMatchComparable[ChosenSide](),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Left{
				OpticTypeExpr: ot,
				B:             reflect.TypeFor[B](),
			}
		}),
	)
}

// Right returns a [Prism] that matches on the right element of a [mo.Either]
func Right[A any, B any]() Optic[ChosenSide, mo.Either[A, B], mo.Either[A, B], B, B, ReturnMany, ReadWrite, BiDir, Pure] {
	return RightP[A, B, B]()
}

func RightP[A, B, C any]() Optic[ChosenSide, mo.Either[A, B], mo.Either[A, C], B, C, ReturnMany, ReadWrite, BiDir, Pure] {
	return PrismIP[ChosenSide, mo.Either[A, B], mo.Either[A, C], B, C](
		func(source mo.Either[A, B]) mo.Either[mo.Either[A, C], ValueI[ChosenSide, B]] {
			if b, ok := source.Right(); ok {
				return mo.Right[mo.Either[A, C], ValueI[ChosenSide, B]](ValI(ChosenRight, b))
			} else {
				return mo.Left[mo.Either[A, C], ValueI[ChosenSide, B]](mo.Left[A, C](source.MustLeft()))
			}
		},
		func(focus C) mo.Either[A, C] {
			return mo.Right[A, C](focus)
		},
		IxMatchComparable[ChosenSide](),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Right{
				OpticTypeExpr: ot,
				A:             reflect.TypeFor[A](),
			}
		}),
	)
}

// BesideEither returns a [Traversal] that focuses the elements in either the left or right optics depending whether the source is left or right.
func BesideEither[IL, IR, SL, TL, SR, TR, A any, B any, RETL, RETR any, RWL, RWR any, DIRL, DIRR any, LERR, RERR any](left Optic[IL, SL, TL, A, B, RETL, RWL, DIRL, LERR], right Optic[IR, SR, TR, A, B, RETR, RWR, DIRR, RERR]) Optic[mo.Either[IL, IR], mo.Either[SL, SR], mo.Either[TL, TR], A, B, CompositionTree[RETL, RETR], CompositionTree[RWL, RWR], UniDir, CompositionTree[LERR, RERR]] {
	return CombiTraversal[CompositionTree[RETL, RETR], CompositionTree[RWL, RWR], CompositionTree[LERR, RERR], mo.Either[IL, IR], mo.Either[SL, SR], mo.Either[TL, TR], A, B](
		func(ctx context.Context, source mo.Either[SL, SR]) SeqIE[mo.Either[IL, IR], A] {
			return func(yield func(ValueIE[mo.Either[IL, IR], A]) bool) {
				if l, ok := source.Left(); ok {
					left.AsIter()(ctx, l)(func(val ValueIE[IL, A]) bool {
						focusIndex, o, focusErr := val.Get()
						return yield(ValIE(mo.Left[IL, IR](focusIndex), o, JoinCtxErr(ctx, focusErr)))
					})
				} else {
					right.AsIter()(ctx, source.MustRight())(func(val ValueIE[IR, A]) bool {
						focusIndex, o, focusErr := val.Get()
						return yield(ValIE(mo.Right[IL, IR](focusIndex), o, JoinCtxErr(ctx, focusErr)))
					})
				}
			}
		},
		func(ctx context.Context, source mo.Either[SL, SR]) (int, error) {
			if l, ok := source.Left(); ok {
				return left.AsLengthGetter()(ctx, l)
			} else {
				return right.AsLengthGetter()(ctx, source.MustRight())
			}
		},
		func(ctx context.Context, fmap func(index mo.Either[IL, IR], focus A) (B, error), source mo.Either[SL, SR]) (mo.Either[TL, TR], error) {
			if l, ok := source.Left(); ok {
				rl, err := left.AsModify()(ctx, func(i IL, o A) (B, error) {
					ret, err := fmap(mo.Left[IL, IR](i), o)
					return ret, JoinCtxErr(ctx, err)
				}, l)
				err = JoinCtxErr(ctx, err)
				if err != nil {
					var ret mo.Either[TL, TR]
					return ret, err
				}

				return mo.Left[TL, TR](rl), ctx.Err()
			} else {
				rr, err := right.AsModify()(ctx, func(i IR, o A) (B, error) {
					ret, err := fmap(mo.Right[IL, IR](i), o)
					return ret, JoinCtxErr(ctx, err)
				}, source.MustRight())
				err = JoinCtxErr(ctx, err)
				if err != nil {
					var ret mo.Either[TL, TR]
					return ret, err
				}

				return mo.Right[TL, TR](rr), ctx.Err()
			}
		},
		func(ctx context.Context, index mo.Either[IL, IR], source mo.Either[SL, SR]) SeqIE[mo.Either[IL, IR], A] {
			return func(yield func(ValueIE[mo.Either[IL, IR], A]) bool) {
				if s, sLeft := source.Left(); sLeft {
					if i, iLeft := index.Left(); iLeft {
						left.AsIxGetter()(ctx, i, s)(func(val ValueIE[IL, A]) bool {
							index, focus, err := val.Get()
							return yield(ValIE(mo.Left[IL, IR](index), focus, err))
						})
					}
				} else {
					s := source.MustRight()
					if i, iRight := index.Right(); iRight {
						right.AsIxGetter()(ctx, i, s)(func(val ValueIE[IR, A]) bool {
							index, focus, err := val.Get()
							return yield(ValIE(mo.Right[IL, IR](index), focus, err))
						})
					}
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
				return expr.BesideEither{
					OpticTypeExpr: ot,
					OnLeft:        left.AsExpr(),
					OnRight:       right.AsExpr(),
				}
			},
			left,
			right,
		),
	)
}
