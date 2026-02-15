package otree

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"reflect"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

type PathNode[I any] struct {
	parent *PathNode[I]
	value  I
}

func (p *PathNode[I]) Append(i I) *PathNode[I] {
	return &PathNode[I]{
		parent: p,
		value:  i,
	}
}

func (p *PathNode[I]) Parent() *PathNode[I] {
	if p == nil {
		return nil
	}

	return p.parent
}

func (p *PathNode[I]) Value() I {
	return p.value
}

func (p *PathNode[I]) Slice() []I {
	return MustGet(
		SliceOf(
			Reversed(
				TraversePath[I](),
			),
			MustGet(Length(TraversePath[I]()), p),
		),
		p,
	)
}

func (p *PathNode[I]) String() string {
	return fmt.Sprintf("%v", p.Slice())
}

func (p *PathNode[I]) IxMatch(ixMatch func(indexA, indexB I) bool, other *PathNode[I]) bool {
	if p == nil {
		if other == nil {
			return true
		} else {
			return false
		}
	} else {
		if other == nil {
			return false
		} else {
			match := ixMatch(p.value, other.value)
			if !match {
				return match
			}

			return p.parent.IxMatch(ixMatch, other.parent)
		}
	}
}

func Path[I any](indices ...I) *PathNode[I] {
	var p *PathNode[I]
	for _, v := range indices {
		p = p.Append(v)
	}
	return p
}

func PathValue[I any]() Optic[Void, *PathNode[I], *PathNode[I], I, I, ReturnMany, ReadWrite, UniDir, Pure] {
	return PtrFieldLens(func(source *PathNode[I]) *I { return &source.value })
}

func TraversePath[I any]() Optic[int, *PathNode[I], *PathNode[I], I, I, ReturnMany, ReadOnly, UniDir, Pure] {
	return Iteration[*PathNode[I], I](
		func(source *PathNode[I]) iter.Seq[I] {
			return func(yield func(I) bool) {
				for i := source; i != nil; i = i.parent {
					if !yield(i.value) {
						return
					}
				}
			}
		},
		func(source *PathNode[I]) int {
			len := 0
			for i := source; i != nil; i = i.parent {
				len++
			}
			return len
		},
		ExprDef(func(t expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Traverse{
				OpticTypeExpr: t,
			}
		}),
	)
}

// TreeNode is a data structure representing a tree.
// A TreeNode has a value and a list of child trees.
//
// See:
//   - [Tree] for the constructor for this type.
type TreeNode[I, T, ERR any] struct {
	value    T
	children Collection[I, TreeNode[I, T, ERR], ERR]
}

// Tree is the constructor for [TreeNode]
//
// See:
//   - [Tree] for a version that supports arbitrary index types.
func Tree[T any](value T, children ...Collection[int, TreeNode[int, T, Pure], Pure]) TreeNode[int, T, Pure] {
	return TreeI(value, IxMatchComparable[int](), children...)
}

func TreeI[I any, T any](value T, ixMatch func(a, b I) bool, children ...Collection[I, TreeNode[I, T, Pure], Pure]) TreeNode[I, T, Pure] {
	return treeIE(value, ixMatch, children...)
}

func TreeE[T any](value T, children ...Collection[int, TreeNode[int, T, Err], Err]) TreeNode[int, T, Err] {
	return treeIE(value, IxMatchComparable[int](), children...)
}

// TreeIE is a constructor for [TreeNode]
//
// See:
//   - [Tree] for a simpler version that supports only comparable index types.
func TreeIE[I any, T any, ERR any](value T, ixMatch func(a, b I) bool, children ...Collection[I, TreeNode[I, T, ERR], ERR]) TreeNode[I, T, ERR] {
	return treeIE(value, ixMatch, children...)
}

func treeIE[I any, T any, ERR any](value T, ixMatch func(a, b I) bool, children ...Collection[I, TreeNode[I, T, ERR], ERR]) TreeNode[I, T, ERR] {
	switch len(children) {
	case 0:
		return TreeNode[I, T, ERR]{
			value: value,
			children: ColIE[ERR, I, TreeNode[I, T, ERR]](
				func(ctx context.Context) SeqIE[I, TreeNode[I, T, ERR]] {
					return func(yield func(ValueIE[I, TreeNode[I, T, ERR]]) bool) {}
				},
				func(ctx context.Context, index I) SeqIE[I, TreeNode[I, T, ERR]] {
					return func(yield func(ValueIE[I, TreeNode[I, T, ERR]]) bool) {}
				},
				ixMatch,
				func(ctx context.Context) (int, error) {
					return 0, nil
				},
			),
		}
	case 1:
		return TreeNode[I, T, ERR]{
			value:    value,
			children: children[0],
		}
	default:
		return TreeNode[I, T, ERR]{
			value: value,
			children: ColIE[ERR, I, TreeNode[I, T, ERR]](
				func(ctx context.Context) SeqIE[I, TreeNode[I, T, ERR]] {
					return func(yield func(ValueIE[I, TreeNode[I, T, ERR]]) bool) {
						cont := true
						for _, c := range children {
							c.AsIter()(ctx)(func(val ValueIE[I, TreeNode[I, T, ERR]]) bool {
								cont = yield(val)
								return cont
							})
							if !cont {
								break
							}
						}
					}
				},
				func(ctx context.Context, index I) SeqIE[I, TreeNode[I, T, ERR]] {
					return func(yield func(ValueIE[I, TreeNode[I, T, ERR]]) bool) {
						cont := true
						for _, c := range children {
							c.AsIxGet()(ctx, index)(func(val ValueIE[I, TreeNode[I, T, ERR]]) bool {
								cont = yield(val)
								return cont
							})
							if !cont {
								break
							}
						}
					}
				},
				ixMatch,
				func(ctx context.Context) (int, error) {
					l := 0
					for _, col := range children {
						colLen, err := col.AsLengthGetter()(ctx)
						if err != nil {
							return 0, err
						}
						l = l + colLen
					}
					return l, nil
				},
			),
		}
	}
}

func (r TreeNode[I, T, ERR]) String() string {
	return fmt.Sprintf("{%v : %v}", r.value, r.children)
}

func (r *TreeNode[I, T, ERR]) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Value    T
		Children Collection[I, TreeNode[I, T, ERR], ERR]
	}{
		Value:    r.value,
		Children: r.children,
	})
}

func (r *TreeNode[I, T, ERR]) UnmarshalJSON(data []byte) error {

	if r.children != nil {
		return errors.New("UnmarshalJson: TreeNode[I,T] is immutable")
	}

	u := struct {
		Value    T
		Children Collection[I, TreeNode[I, T, ERR], ERR]
	}{}

	err := json.Unmarshal(data, &u)

	r.value = u.Value
	r.children = u.Children

	return err
}

func EqPath[I any, ERR TPure](right *PathNode[I], eq Predicate[lo.Tuple2[I, I], ERR]) Optic[Void, *PathNode[I], *PathNode[I], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2Bind(EqPathT2(eq), right)
}

func EqPathT2[I any, ERR TPure](eq Predicate[lo.Tuple2[I, I], ERR]) Optic[Void, lo.Tuple2[*PathNode[I], *PathNode[I]], lo.Tuple2[*PathNode[I], *PathNode[I]], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return Ret1(Ro(Ud(EPure(
		AndOp(
			Compose(
				DelveT2(Length(TraversePath[I]())),
				EqT2[int](),
			),
			All(
				DelveT2(TraversePath[I]()),
				eq,
			),
		),
	))))
}

func EqTreeT2[I comparable, A, ERR any, PERR TPure](eq Predicate[lo.Tuple2[A, A], PERR]) Optic[Void, lo.Tuple2[TreeNode[I, A, ERR], TreeNode[I, A, ERR]], lo.Tuple2[TreeNode[I, A, ERR], TreeNode[I, A, ERR]], bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return EqTreeT2I[I, A, ERR](EqT2[I](), eq)
}

func EqTree[I comparable, A, ERR any, PERR TPure](right TreeNode[I, A, ERR], eq Predicate[lo.Tuple2[A, A], PERR]) Optic[Void, TreeNode[I, A, ERR], TreeNode[I, A, ERR], bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return EqTreeI(right, EqT2[I](), eq)
}

func EqTreeT2I[I, A, ERR any, IERR, PERR TPure](ixMatch Predicate[lo.Tuple2[I, I], IERR], eq Predicate[lo.Tuple2[A, A], PERR]) Optic[Void, lo.Tuple2[TreeNode[I, A, ERR], TreeNode[I, A, ERR]], lo.Tuple2[TreeNode[I, A, ERR], TreeNode[I, A, ERR]], bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	//This looks  like it doesn't compare the indexes but EqT2Of applies uses the ixmatch func for the indices and verifies that the two TopDowns performed on the tuple provided to the output optic are of equal length and that they produce equal values in lock step with one another.
	return EErrL(EqT2Of(
		ComposeLeft(
			TopDown(
				TraverseTreeChildrenI[I, A, ERR](
					PredToIxMatch(ixMatch),
				),
			),
			TreeValue[I, A, ERR](),
		),
		eq,
	))

}

func EqTreeI[I, A, ERR any, PERR, IERR TPure](right TreeNode[I, A, ERR], ixMatch Predicate[lo.Tuple2[I, I], PERR], eq Predicate[lo.Tuple2[A, A], IERR]) Optic[optic.Void, TreeNode[I, A, ERR], TreeNode[I, A, ERR], bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return OpT2Bind(EqTreeT2I[I, A, ERR](ixMatch, eq), right)
}

// TreeChildren returns a [Lens] that focuses the direct children of a [TreeNode]
//
// See:
//   - [TreeChildrenP] for a polymorphic version.
func TreeValue[I, A, ERR any]() Optic[Void, TreeNode[I, A, ERR], TreeNode[I, A, ERR], A, A, ReturnOne, ReadWrite, UniDir, Pure] {
	return FieldLens(func(source *TreeNode[I, A, ERR]) *A { return &source.value })
}

// TreeChildren returns a [Lens] that focuses the direct children of a [TreeNode]
//
// See:
//   - [TreeChildrenP] for a polymorphic version.
func TreeChildren[I, A, ERR any]() Optic[Void, TreeNode[I, A, ERR], TreeNode[I, A, ERR], Collection[I, TreeNode[I, A, ERR], ERR], Collection[I, TreeNode[I, A, ERR], ERR], ReturnOne, ReadWrite, UniDir, Pure] {
	//TODO: consider I,E and IE versions to simplify the basic usage.
	return FieldLens(func(source *TreeNode[I, A, ERR]) *Collection[I, TreeNode[I, A, ERR], ERR] { return &source.children })
}

func TraverseTreeChildren[I comparable, A, ERR any]() Optic[I, TreeNode[I, A, ERR], TreeNode[I, A, ERR], TreeNode[I, A, ERR], TreeNode[I, A, ERR], ReturnMany, ReadWrite, UniDir, ERR] {
	return TraverseTreeChildrenI[I, A, ERR](IxMatchComparable[I]())
}

func TraverseTreeChildrenI[I, A, ERR any](ixmatch func(a, b I) bool) Optic[I, TreeNode[I, A, ERR], TreeNode[I, A, ERR], TreeNode[I, A, ERR], TreeNode[I, A, ERR], ReturnMany, ReadWrite, UniDir, ERR] {
	return EErrR(RetM(Rw(Ud(Compose(TreeChildren[I, A, ERR](), TraverseColIEP[I, TreeNode[I, A, ERR], TreeNode[I, A, ERR], ERR](ixmatch))))))
}

func TraverseTreeChildrenP[I, A, B any, ERR any, PERR TPure](ixMatch Predicate[lo.Tuple2[I, I], PERR]) Optic[*PathNode[I], TreeNode[I, A, ERR], TreeNode[I, B, ERR], A, B, ReturnMany, ReadWrite, UniDir, ERR] {
	//TopDown and BottomUp don't support polymorphism.
	var iterRecursive func(ctx context.Context, path *PathNode[I], source TreeNode[I, A, ERR], yield func(ValueIE[*PathNode[I], A]) bool) bool
	iterRecursive = func(ctx context.Context, path *PathNode[I], source TreeNode[I, A, ERR], yield func(ValueIE[*PathNode[I], A]) bool) bool {
		if !yield(ValIE(path, source.value, nil)) {
			return false
		}

		cont := true
		source.children.AsIter()(ctx)(func(val ValueIE[I, TreeNode[I, A, ERR]]) bool {
			index, focus, err := val.Get()
			if err != nil {
				var a A
				cont = yield(ValIE(path, a, err))
				return cont
			}

			cont = iterRecursive(ctx, path.Append(index), focus, yield)
			return cont
		})

		return cont
	}

	var modifyRecursive func(ctx context.Context, path *PathNode[I], source TreeNode[I, A, ERR], fmap func(index *PathNode[I], focus A) (B, error)) (TreeNode[I, B, ERR], error)
	modifyRecursive = func(ctx context.Context, path *PathNode[I], source TreeNode[I, A, ERR], fmap func(index *PathNode[I], focus A) (B, error)) (TreeNode[I, B, ERR], error) {
		var ret TreeNode[I, B, ERR]

		b, err := fmap(path, source.value)
		if err != nil {
			return ret, err
		}

		ret.value = b

		//We can't evaluate the transformed seq lazily. Rewrite requires that the fmap be called and it detects that a modification occurred.
		var retErr error
		var bChildren []ValueIE[I, TreeNode[I, B, ERR]]

		source.children.AsIter()(ctx)(func(val ValueIE[I, TreeNode[I, A, ERR]]) bool {
			index, focus, focusErr := val.Get()
			if focusErr != nil {
				retErr = focusErr
				return false
			}

			b, err := modifyRecursive(ctx, path.Append(index), focus, fmap)
			if err != nil {
				retErr = err
				return false
			}

			bChildren = append(bChildren, ValIE(index, b, nil))
			return true
		})

		if retErr != nil {
			return ret, retErr
		}

		ret.children = ValColIE[I, TreeNode[I, B, ERR], ERR](
			func(a, b I) bool {
				return Must(PredGet(ctx, ixMatch, lo.T2(a, b)))
			},
			bChildren...,
		)
		return ret, nil

	}

	var ixGetRecursive func(ctx context.Context, path *PathNode[I], index []I, source TreeNode[I, A, ERR], yield func(ValueIE[*PathNode[I], A]) bool) bool
	ixGetRecursive = func(ctx context.Context, path *PathNode[I], index []I, source TreeNode[I, A, ERR], yield func(ValueIE[*PathNode[I], A]) bool) bool {
		if len(index) == 0 {
			return yield(ValIE(path, source.value, nil))
		}

		cont := true
		childIndex := index[1:]
		source.children.AsIxGet()(ctx, index[0])(func(val ValueIE[I, TreeNode[I, A, ERR]]) bool {
			focusIndex, focus, err := val.Get()
			if err != nil {
				var a A
				cont = yield(ValIE(path, a, err))
				return cont
			}

			cont = ixGetRecursive(ctx, path.Append(focusIndex), childIndex, focus, yield)
			return cont
		})

		return cont
	}

	return CombiTraversal[ReturnMany, ReadWrite, ERR, *PathNode[I], TreeNode[I, A, ERR], TreeNode[I, B, ERR], A, B](
		func(ctx context.Context, source TreeNode[I, A, ERR]) SeqIE[*PathNode[I], A] {
			return func(yield func(ValueIE[*PathNode[I], A]) bool) {
				iterRecursive(ctx, nil, source, yield)
			}
		},
		nil,
		func(ctx context.Context, fmap func(index *PathNode[I], focus A) (B, error), source TreeNode[I, A, ERR]) (TreeNode[I, B, ERR], error) {
			return modifyRecursive(ctx, nil, source, fmap)
		},
		func(ctx context.Context, index *PathNode[I], source TreeNode[I, A, ERR]) SeqIE[*PathNode[I], A] {
			if index == nil {
				return func(yield func(ValueIE[*PathNode[I], A]) bool) {}
			}
			indexSlice := index.Slice()

			return func(yield func(ValueIE[*PathNode[I], A]) bool) {
				ixGetRecursive(ctx, nil, indexSlice, source, yield)
			}

		},
		func(indexA, indexB *PathNode[I]) bool {
			return indexA.IxMatch(func(indexA, indexB I) bool {
				return Must(PredGet(context.Background(), ixMatch, lo.T2(indexA, indexB)))
			}, indexB)
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.TraverseTreeChildrenExpr{
				OpticTypeExpr: ot,
				IxMatch:       ixMatch.AsExpr(),
				I:             reflect.TypeFor[I](),
				A:             reflect.TypeFor[A](),
			}
		}),
	)

}

func reIndexTree[A, I, J, L any, IRET TReturnOne, IRW any, IDIR any, ERR any](ctx context.Context, ixmap Optic[L, I, I, J, J, IRET, IRW, IDIR, ERR], ixmatch func(J, J) bool, source TreeNode[I, A, ERR]) TreeNode[J, A, ERR] {
	return TreeNode[J, A, ERR]{
		value: source.value,
		children: ColIE[ERR](
			func(ctx context.Context) SeqIE[J, TreeNode[J, A, ERR]] {
				return func(yield func(val ValueIE[J, TreeNode[J, A, ERR]]) bool) {
					source.children.AsIter()(ctx)(func(val ValueIE[I, TreeNode[I, A, ERR]]) bool {
						index, focus, err := val.Get()
						if err != nil {
							var j J
							var t TreeNode[J, A, ERR]
							return yield(ValIE(j, t, err))
						}
						j, err := ixmap.AsOpGet()(ctx, index)
						return yield(ValIE(j, reIndexTree(ctx, ixmap, ixmatch, focus), err))
					})
				}
			},
			nil,
			ixmatch,
			source.children.AsLengthGetter(),
		),
	}
}

func reverseReIndexTree[A, I, J, L any, IRET TReturnOne, IRW any, IDIR any, ERR any](ctx context.Context, ixmap Optic[L, I, I, J, J, IRET, IRW, IDIR, ERR], ixmatch func(I, I) bool, source TreeNode[J, A, ERR]) TreeNode[I, A, ERR] {
	return TreeNode[I, A, ERR]{
		value: source.value,
		children: ColIE[ERR](
			func(ctx context.Context) SeqIE[I, TreeNode[I, A, ERR]] {
				return func(yield func(val ValueIE[I, TreeNode[I, A, ERR]]) bool) {
					source.children.AsIter()(ctx)(func(val ValueIE[J, TreeNode[J, A, ERR]]) bool {
						index, focus, err := val.Get()
						j, err := ixmap.AsReverseGetter()(ctx, index)
						return yield(ValIE(j, reverseReIndexTree(ctx, ixmap, ixmatch, focus), err))
					})
				}
			},
			nil,
			ixmatch,
			source.children.AsLengthGetter(),
		),
	}
}

func ReIndexedTree[A any, I comparable, J comparable, L any, IRET TReturnOne, IRW any, IDIR any, ERR any](ixmap Optic[L, I, I, J, J, IRET, IRW, IDIR, ERR]) Optic[Void, mo.Option[TreeNode[I, A, ERR]], mo.Option[TreeNode[I, A, ERR]], mo.Option[TreeNode[J, A, ERR]], mo.Option[TreeNode[J, A, ERR]], ReturnOne, IRW, IDIR, ERR] {
	return ReIndexedTreeP[A, A](ixmap, IxMatchComparable[I](), IxMatchComparable[J]())
}

func ReIndexedTreeP[A, B, I, J, L any, IRET TReturnOne, IRW any, IDIR any, ERR any](ixmap Optic[L, I, I, J, J, IRET, IRW, IDIR, ERR], ixmatch func(I, I) bool, ixmatchj func(J, J) bool) Optic[Void, mo.Option[TreeNode[I, A, ERR]], mo.Option[TreeNode[I, B, ERR]], mo.Option[TreeNode[J, A, ERR]], mo.Option[TreeNode[J, B, ERR]], ReturnOne, IRW, IDIR, ERR] {
	x := CombiIso[ReadWrite, BiDir, ERR, mo.Option[TreeNode[I, A, ERR]], mo.Option[TreeNode[I, B, ERR]], mo.Option[TreeNode[J, A, ERR]], mo.Option[TreeNode[J, B, ERR]]](
		func(ctx context.Context, source mo.Option[TreeNode[I, A, ERR]]) (mo.Option[TreeNode[J, A, ERR]], error) {
			if node, ok := source.Get(); ok {
				return mo.Some(reIndexTree(ctx, ixmap, ixmatchj, node)), nil
			} else {
				return mo.None[TreeNode[J, A, ERR]](), nil
			}
		},
		func(ctx context.Context, focus mo.Option[TreeNode[J, B, ERR]]) (mo.Option[TreeNode[I, B, ERR]], error) {
			if node, ok := focus.Get(); ok {
				return mo.Some(reverseReIndexTree(ctx, ixmap, ixmatch, node)), nil
			} else {
				return mo.None[TreeNode[I, B, ERR]](), nil
			}
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.ReIndexedTree{
					OpticTypeExpr: ot,
					I:             reflect.TypeFor[I](),
					A:             reflect.TypeFor[A](),
					B:             reflect.TypeFor[B](),
					IxMap:         ixmap.AsExpr(),
				}
			},
			ixmap,
		),
	)

	return CombiRw[IRW](CombiDir[IDIR](x))
}
