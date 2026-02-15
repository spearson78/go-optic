package model

import (
	"container/list"
	"context"
	"reflect"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

type ReplayBuffer[A any] struct {
	buf *list.List
}

func NewReplayBuffer[A any]() ReplayBuffer[A] {
	return ReplayBuffer[A]{
		buf: list.New(),
	}
}

func TraverseReplayBuffer[A any]() Optic[int, ReplayBuffer[A], ReplayBuffer[A], A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	l := FieldLens(func(s *ReplayBuffer[A]) **list.List { return &s.buf })
	return traverseReplayBufferP(l)
}

func TraverseReplayBufferP[A, B any]() Optic[int, ReplayBuffer[A], ReplayBuffer[B], A, B, ReturnMany, ReadWrite, UniDir, Pure] {

	l := FieldLensP[ReplayBuffer[A], ReplayBuffer[B], *list.List, *list.List](
		func(source *ReplayBuffer[A]) **list.List {
			return &source.buf
		},
		func(focus *list.List, source ReplayBuffer[A]) ReplayBuffer[B] {
			return ReplayBuffer[B]{
				buf: focus,
			}
		},
	)

	return traverseReplayBufferP(l)
}

func traverseReplayBufferP[A, B any](l Optic[Void, ReplayBuffer[A], ReplayBuffer[B], *list.List, *list.List, ReturnOne, ReadWrite, UniDir, Pure]) Optic[int, ReplayBuffer[A], ReplayBuffer[B], A, B, ReturnMany, ReadWrite, UniDir, Pure] {
	return CombiTraversal[ReturnMany, ReadWrite, Pure, int, ReplayBuffer[A], ReplayBuffer[B], A, B](
		func(ctx context.Context, source ReplayBuffer[A]) SeqIE[int, A] {
			return func(yield func(ValueIE[int, A]) bool) {
				i := 0
				for element := source.buf.Front(); element != nil; element = element.Next() {
					if !yield(ValIE(i, element.Value.(A), ctx.Err())) {
						i++
						break
					}
					i++
				}
			}
		},
		func(ctx context.Context, source ReplayBuffer[A]) (int, error) {
			return source.buf.Len(), nil
		},
		func(ctx context.Context, fmap func(index int, focus A) (B, error), source ReplayBuffer[A]) (ReplayBuffer[B], error) {
			ret := NewReplayBuffer[B]()

			i := 0
			for element := source.buf.Front(); element != nil; element = element.Next() {
				b, err := fmap(i, element.Value.(A))
				i++
				err = JoinCtxErr(ctx, err)
				if err != nil {
					return ReplayBuffer[B]{}, err
				}
				ret.buf.PushBack(b)
			}

			return ret, ctx.Err()
		},
		func(ctx context.Context, index int, source ReplayBuffer[A]) SeqIE[int, A] {
			return func(yield func(ValueIE[int, A]) bool) {
				if index >= 0 && index < source.buf.Len() {
					i := 0
					for element := source.buf.Front(); element != nil; element = element.Next() {
						if i == index {
							yield(ValIE(index, element.Value.(A), ctx.Err()))
							return
						}
					}
				}
			}
		},
		IxMatchComparable[int](),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Traverse{
				OpticTypeExpr: ot,
			}
		}),
	)
}

func ReplayBufferToCol[A any](len int) Optic[Void, ReplayBuffer[A], ReplayBuffer[A], Collection[int, A, Pure], Collection[int, A, Pure], ReturnOne, ReadWrite, BiDir, Pure] {

	return ReplayBufferToColP[A, A](len)
}

func ReplayBufferToColP[A, B any](len int) Optic[Void, ReplayBuffer[A], ReplayBuffer[B], Collection[int, A, Pure], Collection[int, B, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return CombiIso[ReadWrite, BiDir, Pure, ReplayBuffer[A], ReplayBuffer[B], Collection[int, A, Pure], Collection[int, B, Pure]](
		func(ctx context.Context, source ReplayBuffer[A]) (Collection[int, A, Pure], error) {
			return ColI(
				func(yield func(index int, focus A) bool) {
					i := 0
					for element := source.buf.Front(); element != nil; element = element.Next() {
						if !yield(i, element.Value.(A)) {
							break
						}
					}
				},
				nil,
				IxMatchComparable[int](),
				func() int {
					return source.buf.Len()
				},
			), nil
		},
		func(ctx context.Context, focus Collection[int, B, Pure]) (ReplayBuffer[B], error) {

			ret := NewReplayBuffer[B]()
			var retErr error

			focus.AsIter()(ctx)(func(val ValueIE[int, B]) bool {
				_, focus, focusErr := val.Get()
				if focusErr != nil {
					retErr = focusErr
					return false
				}

				ret.buf.PushBack(focus)

				if ret.buf.Len() > len {
					ret.buf.Remove(ret.buf.Front())
				}

				return true
			})

			return ret, retErr

		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.ToCol{
				OpticTypeExpr: ot,
				I:             reflect.TypeFor[int](),
				A:             reflect.TypeFor[A](),
				B:             reflect.TypeFor[B](),
			}
		}),
	)
}

//go:generate ../../../makecolops replaybuffer_ops.go "model" "ReplayBuffer" "A any" "A,B any" "len int" "ReplayBufferColTypeP[A,A](len)" "ReplayBufferColTypeP[A,B](len)" "int" "ReplayBuffer[A]" "ReplayBuffer[B]" "A" "B"
func ReplayBufferColTypeP[A, B any](len int) CollectionType[int, ReplayBuffer[A], ReplayBuffer[B], A, B, Pure] {
	return ColTypeP(
		ReplayBufferToColP[A, B](len),
		AsReverseGet(ReplayBufferToColP[B, A](len)),
		TraverseReplayBufferP[A, B](),
	)
}
