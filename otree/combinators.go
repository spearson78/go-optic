package otree

import (
	"context"

	"github.com/samber/mo"
	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

// The ResolvePath combinator returns an optic that focuses the child element wth the given path.
func ResolvePath[I, A, RET, RW, DIR, ERR any](children Optic[I, A, A, A, A, RET, RW, DIR, ERR], resolvePath *PathNode[I]) Optic[*PathNode[I], A, A, A, A, ReturnMany, ReadWrite, UniDir, ERR] {

	var modify func(ctx context.Context, path *PathNode[I], fmap func(index *PathNode[I], focus A) (A, error), source A) (A, error)
	modify = func(ctx context.Context, path *PathNode[I], fmap func(index *PathNode[I], focus A) (A, error), source A) (A, error) {

		match := resolvePath.IxMatch(children.AsIxMatch(), path)

		if match {
			newValue, err := fmap(path, source)
			return newValue, err
		}

		return children.AsModify()(ctx, func(childIndex I, childFocus A) (A, error) {
			return modify(ctx, path.Append(childIndex), fmap, childFocus)
		}, source)
	}

	getter := func(ctx context.Context, source A) SeqIE[*PathNode[I], A] {
		return func(yield func(ValueIE[*PathNode[I], A]) bool) {

			pathSlice := MustGet(
				SliceOf(
					Reversed(
						TraversePath[I](),
					),
					MustGet(Length(TraversePath[I]()), resolvePath),
				),
				resolvePath,
			)
			treeIxGet(ctx, nil, pathSlice, source, children, True[ValueI[*PathNode[I], A]](), yield)
		}
	}

	return CombiTraversal[ReturnMany, ReadWrite, ERR, *PathNode[I], A, A, A, A](
		getter,
		nil,
		func(ctx context.Context, fmap func(index *PathNode[I], focus A) (A, error), source A) (A, error) {
			return modify(ctx, nil, fmap, source)
		},
		func(ctx context.Context, index *PathNode[I], source A) SeqIE[*PathNode[I], A] {
			match := resolvePath.IxMatch(children.AsIxMatch(), index)

			if match {
				return getter(ctx, source)
			} else {
				return func(yield func(ValueIE[*PathNode[I], A]) bool) {}
			}
		},
		func(indexA, indexB *PathNode[I]) bool {
			return indexA.IxMatch(children.AsIxMatch(), indexB)
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {

				path := MustGet(
					SliceOf(
						Reversed(
							Compose(
								TraversePath[I](),
								IsoCast[I, any](),
							),
						),
						MustGet(Length(TraversePath[I]()), resolvePath),
					),
					resolvePath,
				)

				return expr.ResolvePath{

					OpticTypeExpr: ot,
					Children:      children.AsExpr(),
					Path:          path,
				}
			},
			children,
		),
	)

}

// RewriteOp constructs a rewrite operation.
//
// See:
//   - [Rewrite] for more details.
//   - [RewriteOp] for a non index aware version
//   - [RewriteOpI] for an index aware version
//   - [RewriteOpE] for an error raising version
//   - [RewriteOpIE] for an index aware error raising version
func RewriteOp[A any](rewriteHandler func(node A) (A, bool)) Optic[Void, A, A, mo.Option[A], mo.Option[A], ReturnOne, ReadOnly, UniDir, Pure] {
	return Op[A, mo.Option[A]](
		func(tn A) mo.Option[A] {
			newNode, modified := rewriteHandler(tn)
			if modified {
				return mo.Some(newNode)
			} else {
				return mo.None[A]()
			}
		},
	)
}

// RewriteOpE constructs an error raising rewrite operation.
//
// See:
//   - [Rewrite] for more details.
//   - [RewriteOp] for a non index aware version
//   - [RewriteOpI] for an index aware version
//   - [RewriteOpE] for an error raising version
//   - [RewriteOpIE] for an index aware error raising version
func RewriteOpE[A any](rewriteHandler func(ctx context.Context, node A) (A, bool, error)) Optic[Void, A, A, mo.Option[A], mo.Option[A], ReturnOne, ReadOnly, UniDir, Err] {
	return OpE[A, mo.Option[A]](
		func(ctx context.Context, source A) (mo.Option[A], error) {

			newNode, modified, err := rewriteHandler(ctx, source)
			if err != nil {
				return mo.None[A](), err
			}
			if modified {
				return mo.Some(newNode), nil
			} else {
				return mo.None[A](), nil
			}

		},
	)
}

// The Rewrite combinator focuses a rewritten version of the source tree.
//
// The op parameter can be constructed using the [RewriteOp] family of constructors.
//
// The tree is rewritten by applying the op to very node in the tree. If the op returns a some option then the value replaces the node in the tree. If the op returns none then the node remains unaltered.
// If any node was replaced then the whole process is run again until the op returns none for all nodes in the tree.
//
// See:
//   - [Rewrite] for a simpler non index non polymorphic version
//   - [RewriteI] for an index aware version
func Rewrite[I, A, RET any, RW TReadWrite, DIR, ERR any, ORET TReturnOne, OERR any](children Optic[I, A, A, A, A, RET, RW, DIR, ERR], op Operation[A, mo.Option[A], ORET, OERR]) Optic[Void, A, A, A, A, ReturnOne, ReadOnly, UniDir, CompositionTree[ERR, OERR]] {
	return RewriteI(children, OpToOpI[*PathNode[I]](op))
}

// TopDownFiltered is a [Traversal] that walks the tree in a top down order and filters out branches that don't match the predicate.
// During the traversal the child index is used to construct a path represented as a [PathNode]
//
// Note: the predicate is also applied to the root node which may result in no tree elements being focused.
//
// In a tree with this structure
//
//	+-A
//	  +-B
//	  | +-C
//	  +-D
//
// The elements will be focused in the order A,B,C,D
// See:
// - [TopDownFilteredI] for an index aware predicate version.
// - [TopDownMatch] for a version that performs a match instead of filter..
// - [TopDown] for a simpler non filtering version.
// - [BottomUpFiltered] for a bottom up version
// - [BreadthFirstFiltered] for a breadth first version
func TopDownFiltered[I, A, RET, RW, DIR, ERR, ERRP any](children Optic[I, A, A, A, A, RET, RW, DIR, ERR], pred Predicate[A, ERRP]) Optic[*PathNode[I], A, A, A, A, ReturnMany, RW, UniDir, CompositionTree[ERR, ERRP]] {
	return TopDownFilteredI(children, PredToOpI[*PathNode[I]](pred))
}

// TopDownMatch is a [Traversal] that walks the tree in a top down order and focuses only branches that match the predicate.
// The top down walk ends at a matching branch.
//
// During the traversal the child index is used to construct a path represented as a [PathNode]
//
// In a tree with this structure
//
//	+-A
//	  +-B
//	  | +-C
//	  +-D
//
// The elements will be focused in the order A,B,C,D
// See:
// - [TopDownMatchI] for an index aware version
// - [TopDownFilteredI] for an version that filters instead of matches.
// - [TopDown] for a simpler non filtering version.
// - [BottomUpFiltered] for a bottom up version
// - [BreadthFirstFiltered] for a breadth first version
func TopDownMatch[I, S, RET any, RW any, DIR any, ERR, MERR any](children Optic[I, S, S, S, S, RET, RW, DIR, ERR], matcher Predicate[S, MERR]) Optic[*PathNode[I], S, S, S, S, ReturnMany, RW, UniDir, ERR] {
	ixMatcher := RetR(RwR(Ud(EErrR(Compose(ValueIValue[*PathNode[I], S](), PredToOptic(matcher))))))
	return TopDownMatchI(children, ixMatcher)
}

// TopDown is a [Traversal] that walks the tree in a top down order.
// During the traversal the child index is used to construct a path represented as a [PathNode]
//
// In a tree with this structure
//
//	+-A
//	  +-B
//	  | +-C
//	  +-D
//
// The elements will be focused in the order A,B,C,D
// See:
// - [TopDownFiltered] for a filtering version.
// - [BottomUp] for a bottom up version
// - [BreadthFirst] for a breadth first version
func TopDown[I, A any, RET, RW, DIR, ERR any](children Optic[I, A, A, A, A, RET, RW, DIR, ERR]) Optic[*PathNode[I], A, A, A, A, ReturnMany, RW, UniDir, ERR] {
	return EErrL(TopDownFilteredI(children, True[ValueI[*PathNode[I], A]]()))
}

// BottomUpFiltered is a [Traversal] that walks the tree in a bottom up order and filters out branches that don't match the predicate.
// During the traversal the child index is used to construct a path represented as a [PathNode]
//
// Note: the predicate is applied applied in top down order including the root node which may result in no tree elements being focused.
//
// In a tree with this structure
//
//	+-A
//	  +-B
//	  | +-C
//	  +-D
//
// The elements will be focused in the order C,B,D,A
// See:
// - [BottomUpFilteredI] for an index aware predicate version.
// - [BottomUp] for a simpler non filtering version.
// - [TopDownFiltered] for a top down version
// - [BreadthFirstFiltered] for a breadth first version
func BottomUpFiltered[I, A, RET, RW, DIR, ERR, ERRP any](children Optic[I, A, A, A, A, RET, RW, DIR, ERR], pred Predicate[A, ERRP]) Optic[*PathNode[I], A, A, A, A, ReturnMany, RW, UniDir, CompositionTree[ERR, ERRP]] {
	return BottomUpFilteredI(children, PredToOpI[*PathNode[I]](pred))
}

// BottomUp is a [Traversal] that walks the tree in a bottom up order.
// During the traversal the child index is used to construct a path represented as a [PathNode]
//
// In a tree with this structure
//
//	+-A
//	  +-B
//	  | +-C
//	  +-D
//
// The elements will be focused in the order C,B,D,A
// See:
// - [BottomUpFiltered] for a filtering version
// - [TopDown] for a top down version
// - [BreadthFirst] for a breadth first version
func BottomUp[I, A any, RET, RW, DIR, ERR any](children Optic[I, A, A, A, A, RET, RW, DIR, ERR]) Optic[*PathNode[I], A, A, A, A, ReturnMany, RW, UniDir, ERR] {
	return EErrL(BottomUpFilteredI(children, True[ValueI[*PathNode[I], A]]()))
}

// BreadthFirstFiltered is a [Traversal] that walks the tree in a breadth first order and filters out branches that don't match the predicate.
// During the traversal the child index is used to construct a path represented as a [PathNode]
//
// Note: the predicate is applied applied in top down order including the root node which may result in no tree elements being focused.
//
// In a tree with this structure
//
//	+-A
//	  +-B
//	  | +-C
//	  +-D
//
// The elements will be focused in the order A,B,D,C
// See:
// - [BreadthFirstFiltered] for a simpler non index aware predicate version.
// - [BreadthFirst] for a simpler non filtering version.
// - [TopDownFiltered] for a top down version
// - [BottomUpFiltered] for a bottom up version
func BreadthFirstFiltered[I, A, RET, RW, DIR, ERR, ERRP any](children Optic[I, A, A, A, A, RET, RW, DIR, ERR], pred Predicate[A, ERRP]) Optic[*PathNode[I], A, A, A, A, ReturnMany, ReadOnly, UniDir, CompositionTree[ERR, ERRP]] {
	return BreadthFirstFilteredI(children, PredToOpI[*PathNode[I]](pred))
}

// BreadthFirst is a [Traversal] that walks the tree in a breadth first order.
// During the traversal the child index is used to construct a path represented as a [PathNode]
//
// In a tree with this structure
//
//	+-A
//	  +-B
//	  | +-C
//	  +-D
//
// The elements will be focused in the order A,B,D,C
// See:
// - [BreadthFirstFiltered] for a filtering version
// - [TopDown] for a top down version
// - [BottomUp] for a bottom up version
func BreadthFirst[I, A any, RET, RW, DIR, ERR any](children Optic[I, A, A, A, A, RET, RW, DIR, ERR]) Optic[*PathNode[I], A, A, A, A, ReturnMany, ReadOnly, UniDir, ERR] {
	return EErrL(BreadthFirstFilteredI(children, True[ValueI[*PathNode[I], A]]()))
}
