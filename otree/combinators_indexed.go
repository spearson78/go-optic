package otree

import (
	"context"

	"github.com/samber/mo"
	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

// TopDownFilteredI is a [Traversal] that walks the tree in a top down order and filters out branches that don't match the index aware predicate.
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
// - [TopDownFiltered] for a simpler non index aware predicate version.
// - [TopDownMatchI] for a version that performs a match instead of a filter
// - [TopDown] for a simpler non filtering version.
// - [BottomUpFilteredI] for a bottom up version
// - [BreadthFirstFilteredI] for a breadth first version
func TopDownFilteredI[I, A, RET, RW, DIR, ERR any, ERRP any](children Optic[I, A, A, A, A, RET, RW, DIR, ERR], pred PredicateI[*PathNode[I], A, ERRP]) Optic[*PathNode[I], A, A, A, A, ReturnMany, RW, UniDir, CompositionTree[ERR, ERRP]] {
	return CombiTraversal[ReturnMany, RW, CompositionTree[ERR, ERRP], *PathNode[I], A, A, A, A](
		func(ctx context.Context, source A) SeqIE[*PathNode[I], A] {
			return func(yield func(ValueIE[*PathNode[I], A]) bool) {
				iterTdf(ctx, nil, children.AsIter(), pred, source, yield)
			}
		},
		nil,
		func(ctx context.Context, fmap func(index *PathNode[I], focus A) (A, error), source A) (A, error) {
			return modifyTdf(ctx, nil, children.AsModify(), pred, source, fmap)
		},
		func(ctx context.Context, index *PathNode[I], source A) SeqIE[*PathNode[I], A] {
			return func(yield func(val ValueIE[*PathNode[I], A]) bool) {
				resolvePath := MustGet(
					SliceOf(
						Reversed(
							TraversePath[I](),
						),
						MustGet(Length(TraversePath[I]()), index),
					),
					index,
				)
				treeIxGet(ctx, nil, resolvePath, source, children, pred, yield)
			}
		},
		func(indexA, indexB *PathNode[I]) bool {
			return indexA.IxMatch(children.AsIxMatch(), indexB)
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.TopDownFiltered{
					OpticTypeExpr: ot,
					Children:      children.AsExpr(),
					Pred:          pred.AsExpr(),
				}
			},
			children,
			pred,
		),
	)
}

func iterTdf[I, A, ERRP any](ctx context.Context, path *PathNode[I], iterChildren IterFunc[I, A, A], pred PredicateI[*PathNode[I], A, ERRP], tree A, yield func(ValueIE[*PathNode[I], A]) bool) bool {

	match, err := PredGetI(ctx, pred, tree, path)
	if err != nil {
		var a A
		return yield(ValIE(path, a, err))
	}

	if match {
		if yield(ValIE(path, tree, nil)) {
			cont := true
			iterChildren(ctx, tree)(func(val ValueIE[I, A]) bool {
				childIndex, childFocus, childErr := val.Get()
				if childErr != nil {
					var a A
					return yield(ValIE(path, a, childErr))
				}

				cont = iterTdf(ctx, path.Append(childIndex), iterChildren, pred, childFocus, yield)
				return cont
			})
			return cont
		} else {
			return false
		}

	} else {
		return true
	}
}

func modifyTdf[I, A, ERR any](ctx context.Context, path *PathNode[I], modifyChildren ModifyFunc[I, A, A, A, A], pred PredicateI[*PathNode[I], A, ERR], tree A, fmap func(index *PathNode[I], focus A) (A, error)) (A, error) {

	match, err := PredGetI(ctx, pred, tree, path)
	if err != nil {
		var a A
		return a, err
	}

	if match {

		newVal, err := fmap(path, tree)
		if err != nil {
			var a A
			return a, err
		}

		return modifyChildren(
			ctx,
			func(childIndex I, childFocus A) (A, error) {
				return modifyTdf(ctx, path.Append(childIndex), modifyChildren, pred, childFocus, fmap)
			},
			newVal,
		)
	} else {
		return tree, nil
	}
}

// TopDownMatchI is a [Traversal] that walks the tree in a top down order and focuses only branches that match the index aware predicate.
// The top down walk ends at a matching branch.
//
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
// - [TopDownMatch] for a simpler non index aware predicate version.
// - [TopDownFilteredI] for a version that performs a filter instead of a match
// - [TopDown] for a simpler non filtering version.
// - [BottomUpFilteredI] for a bottom up version
// - [BreadthFirstFilteredI] for a breadth first version
func TopDownMatchI[I, S, RET any, RW any, DIR any, ERR, MERR any](children Optic[I, S, S, S, S, RET, RW, DIR, ERR], matcher PredicateI[*PathNode[I], S, MERR]) Optic[*PathNode[I], S, S, S, S, ReturnMany, RW, UniDir, ERR] {

	var recurseChildren func(ctx context.Context, s S, path *PathNode[I], yield func(ValueIE[*PathNode[I], S]) bool) bool
	recurseChildren = func(ctx context.Context, s S, path *PathNode[I], yield func(ValueIE[*PathNode[I], S]) bool) bool {
		ok, err := PredGet(ctx, matcher, ValI(path, s))

		if err != nil {
			var s S
			yield(ValIE(path, s, err))
			return false
		}

		if !ok {
			//No match recurse into children
			cont := true
			children.AsIter()(ctx, s)(func(val ValueIE[I, S]) bool {
				childIndex, focus, focusErr := val.Get()
				focusErr = JoinCtxErr(ctx, focusErr)
				childPath := path.Append(childIndex)
				if focusErr != nil {
					var s S
					yield(ValIE(childPath, s, focusErr))
					return false
				}
				cont = recurseChildren(ctx, focus, childPath, yield)
				return cont
			})
			return cont
		} else {
			return yield(ValIE(path, s, nil))
			//Do not recurse into children on match
		}
	}

	var modifyChildren func(ctx context.Context, source S, path *PathNode[I], fmap func(index *PathNode[I], focus S) (S, error)) (S, error)
	modifyChildren = func(ctx context.Context, source S, path *PathNode[I], fmap func(index *PathNode[I], focus S) (S, error)) (S, error) {
		found, err := PredGet(ctx, matcher, ValI(path, source))
		err = JoinCtxErr(ctx, err)
		if err != nil {
			var s S
			return s, err
		}

		if found {
			//Do not recurse into children on match
			ret, err := fmap(path, source)
			return ret, JoinCtxErr(ctx, err)
		}

		return children.AsModify()(ctx, func(index I, focus S) (S, error) {
			return modifyChildren(ctx, focus, path.Append(index), fmap)
		}, source)
	}

	return CombiTraversal[ReturnMany, RW, ERR, *PathNode[I], S, S, S, S](
		func(ctx context.Context, source S) SeqIE[*PathNode[I], S] {
			return func(yield func(ValueIE[*PathNode[I], S]) bool) {
				recurseChildren(ctx, source, nil, yield)
			}
		},
		nil,
		func(ctx context.Context, fmap func(index *PathNode[I], focus S) (S, error), source S) (S, error) {
			return modifyChildren(ctx, source, nil, fmap)
		},
		nil,
		func(indexA, indexB *PathNode[I]) bool {
			return indexA.IxMatch(children.AsIxMatch(), indexB)
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.TopDownMatch{
					OpticTypeExpr: ot,
					Children:      children.AsExpr(),
					Matcher:       matcher.AsExpr(),
				}
			},
			children,
			matcher,
		),
	)
}

// BottomUpFilteredI is a [Traversal] that walks the tree in a bottom up order and filters out branches that don't match the index aware predicate.
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
// - [BottomUpFiltered] for a simpler non index aware predicate version.
// - [BottomUp] for a simpler non filtering version.
// - [TopDownFilteredI] for a top down version
// - [BreadthFirstFilteredI] for a breadth first version
func BottomUpFilteredI[I, A, RET, RW, DIR, ERR any, ERRP any](children Optic[I, A, A, A, A, RET, RW, DIR, ERR], pred PredicateI[*PathNode[I], A, ERRP]) Optic[*PathNode[I], A, A, A, A, ReturnMany, RW, UniDir, CompositionTree[ERR, ERRP]] {

	return CombiTraversal[ReturnMany, RW, CompositionTree[ERR, ERRP], *PathNode[I], A, A, A, A](
		func(ctx context.Context, source A) SeqIE[*PathNode[I], A] {
			return func(yield func(ValueIE[*PathNode[I], A]) bool) {
				iterBuf(ctx, nil, children.AsIter(), pred, source, yield)
			}
		},
		nil,
		func(ctx context.Context, fmap func(index *PathNode[I], focus A) (A, error), source A) (A, error) {
			return modifyBuf(ctx, nil, children.AsModify(), pred, source, fmap)
		},
		func(ctx context.Context, index *PathNode[I], source A) SeqIE[*PathNode[I], A] {
			return func(yield func(ValueIE[*PathNode[I], A]) bool) {
				resolvePath := MustGet(
					SliceOf(
						Reversed(
							TraversePath[I](),
						),
						MustGet(Length(TraversePath[I]()), index),
					),
					index,
				)
				treeIxGet(ctx, nil, resolvePath, source, children, pred, yield)
			}
		},
		func(indexA, indexB *PathNode[I]) bool {
			return indexA.IxMatch(children.AsIxMatch(), indexB)
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.BottomUpFiltered{
					OpticTypeExpr: ot,
					Children:      children.AsExpr(),
					Pred:          pred.AsExpr(),
				}
			},
			children,
			pred,
		),
	)
}

func iterBuf[I, A, ERRP any](ctx context.Context, path *PathNode[I], iterChildren IterFunc[I, A, A], pred PredicateI[*PathNode[I], A, ERRP], tree A, yield func(ValueIE[*PathNode[I], A]) bool) bool {

	match, err := PredGetI(ctx, pred, tree, path)
	if err != nil {
		var a A
		return yield(ValIE(path, a, err))
	}

	if match {
		cont := true
		iterChildren(ctx, tree)(func(val ValueIE[I, A]) bool {
			childIndex, childFocus, childErr := val.Get()
			if childErr != nil {
				var a A
				return yield(ValIE(path, a, childErr))
			}

			cont = iterBuf(ctx, path.Append(childIndex), iterChildren, pred, childFocus, yield)
			return cont
		})

		if cont {
			return yield(ValIE(path, tree, nil))
		} else {
			return false
		}
	} else {
		return true
	}
}

func modifyBuf[I, A, ERR any](ctx context.Context, path *PathNode[I], modifyChildren ModifyFunc[I, A, A, A, A], pred PredicateI[*PathNode[I], A, ERR], tree A, fmap func(index *PathNode[I], focus A) (A, error)) (A, error) {

	match, err := PredGetI(ctx, pred, tree, path)
	if err != nil {
		var a A
		return a, err
	}

	if match {

		newVal, err := modifyChildren(
			ctx,
			func(childIndex I, childFocus A) (A, error) {
				return modifyBuf(ctx, path.Append(childIndex), modifyChildren, pred, childFocus, fmap)
			},
			tree,
		)
		if err != nil {
			var a A
			return a, err
		}

		return fmap(path, newVal)

	} else {
		return tree, nil
	}
}

// BreadthFirstFilteredI is a [Traversal] that walks the tree in a breadth first order and filters out branches that don't match the index aware predicate.
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
// - [TopDownFilteredI] for a top down version
// - [BottomUpFilteredI] for a bottom up version
func BreadthFirstFilteredI[I, A, RET, RW, DIR, ERR any, ERRP any](children Optic[I, A, A, A, A, RET, RW, DIR, ERR], pred PredicateI[*PathNode[I], A, ERRP]) Optic[*PathNode[I], A, A, A, A, ReturnMany, ReadOnly, UniDir, CompositionTree[ERR, ERRP]] {

	return CombiIteration[ReturnMany, CompositionTree[ERR, ERRP], *PathNode[I], A, A, A, A](
		func(ctx context.Context, source A) SeqIE[*PathNode[I], A] {
			return func(yield func(ValueIE[*PathNode[I], A]) bool) {
				iterBff(ctx, nil, children.AsIter(), pred, source, yield)
			}
		},
		nil,
		func(ctx context.Context, index *PathNode[I], source A) SeqIE[*PathNode[I], A] {
			return func(yield func(ValueIE[*PathNode[I], A]) bool) {
				resolvePath := MustGet(
					SliceOf(
						Reversed(
							TraversePath[I](),
						),
						MustGet(Length(TraversePath[I]()), index),
					),
					index,
				)
				treeIxGet(ctx, nil, resolvePath, source, children, pred, yield)
			}
		},
		func(indexA, indexB *PathNode[I]) bool {
			return indexA.IxMatch(children.AsIxMatch(), indexB)
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.BreadthFirstFiltered{
					OpticTypeExpr: ot,
					Children:      children.AsExpr(),
					Pred:          pred.AsExpr(),
				}
			},
			children,
			pred,
		),
	)
}

func iterBff[I, A, ERRP any](ctx context.Context, path *PathNode[I], iterChildren IterFunc[I, A, A], pred PredicateI[*PathNode[I], A, ERRP], tree A, yield func(ValueIE[*PathNode[I], A]) bool) {
	var queue []ValueIE[*PathNode[I], A]
	queue = append(queue, ValIE(path, tree, nil))

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]

		match, err := PredGetI(ctx, pred, v.Value(), v.Index())
		if err != nil {
			var a A
			yield(ValIE((*PathNode[I])(nil), a, err))
			return
		}

		if match {
			if !yield(v) {
				break
			}

			iterChildren(ctx, v.Value())(func(val ValueIE[I, A]) bool {
				childIndex, childFocus, childErr := val.Get()
				queue = append(queue, ValIE(v.Index().Append(childIndex), childFocus, childErr))
				return true
			})
		}
	}
}

func treeIxGet[I, A, RET, RW, DIR, ERR, ERRP any](ctx context.Context, path *PathNode[I], resolvePath []I, source A, children Optic[I, A, A, A, A, RET, RW, DIR, ERR], pred PredicateI[*PathNode[I], A, ERRP], yield func(ValueIE[*PathNode[I], A]) bool) bool {
	if len(resolvePath) == 0 {
		return yield(ValIE(path, source, nil))
	}

	cont := true
	children.AsIxGetter()(ctx, resolvePath[0], source)(func(val ValueIE[I, A]) bool {
		childIndex, childFocus, childErr := val.Get()
		childErr = JoinCtxErr(ctx, childErr)
		if childErr != nil {
			var a A
			return yield(ValIE(path, a, childErr))
		}

		childPath := path.Append(childIndex)

		match, err := PredGetI(ctx, pred, childFocus, childPath)
		if err != nil {
			var a A
			return yield(ValIE(childPath, a, childErr))
		}

		if match {
			cont = treeIxGet(ctx, childPath, resolvePath[1:], childFocus, children, pred, yield)
		}

		return cont
	})
	return cont
}

func rewriteImpl[I, A, RET any, RW TReadWrite, DIR, ERR any, ORET TReturnOne, OERR any](ctx context.Context, path *PathNode[I], tree A, children Optic[I, A, A, A, A, RET, RW, DIR, ERR], op OperationI[*PathNode[I], A, mo.Option[A], ORET, OERR]) (A, bool, error) {

	newNodeOpt, err := op.AsOpGet()(ctx, ValI(path, tree))
	if err != nil {
		var a A
		return a, false, err
	}

	if newNode, ok := newNodeOpt.Get(); ok {
		return newNode, true, nil
	} else {

		modified := false

		var err error
		tree, err = children.AsModify()(ctx, func(childIndex I, childFocus A) (A, error) {
			newChildOpt, childModified, err := rewriteImpl(ctx, path.Append(childIndex), childFocus, children, op)
			modified = modified || childModified
			return newChildOpt, err
		}, tree)

		if err != nil {
			var a A
			return a, false, err
		}

		return tree, modified, nil
	}
}

// RewriteOpI constructs an index aware rewrite operation.
//
// See:
//   - [RewriteI] for more details.
//   - [RewriteOp] for a non index aware version
//   - [RewriteOpI] for an index aware version
//   - [RewriteOpE] for an error raising version
//   - [RewriteOpIE] for an index aware error raising version
func RewriteOpI[I, A any](rewriteHandler func(index *PathNode[I], node A) (A, bool)) Optic[Void, ValueI[*PathNode[I], A], ValueI[*PathNode[I], A], mo.Option[A], mo.Option[A], ReturnOne, ReadOnly, UniDir, Pure] {
	return Op[ValueI[*PathNode[I], A], mo.Option[A]](
		func(iv ValueI[*PathNode[I], A]) mo.Option[A] {
			return mo.TupleToOption(rewriteHandler(iv.Index(), iv.Value()))
		},
	)
}

// RewriteOpIE constructs an index aware error raising rewrite operation.
//
// See:
//   - [RewriteOp] for a non index aware version
//   - [RewriteOpI] for an index aware version
//   - [RewriteOpE] for an error raising version
//   - [RewriteOpIE] for an index aware error raising version
func RewriteOpIE[I, A any](rewriteHandler func(ctx context.Context, index *PathNode[I], node A) (A, bool, error)) Optic[Void, ValueI[*PathNode[I], A], ValueI[*PathNode[I], A], mo.Option[A], mo.Option[A], ReturnOne, ReadOnly, UniDir, Err] {
	return OpE[ValueI[*PathNode[I], A], mo.Option[A]](
		func(ctx context.Context, source ValueI[*PathNode[I], A]) (mo.Option[A], error) {
			newNode, modified, err := rewriteHandler(ctx, source.Index(), source.Value())
			if err != nil {
				return mo.None[A](), err
			}
			return mo.TupleToOption(newNode, modified), nil
		},
	)
}

// The RewriteI combinator focuses a rewritten version of the source tree.
//
// The op parameter can be constructed using the [RewriteOp] family of constructors.
//
// The tree is rewritten by applying the op to very node in the tree. If the op returns a some option then the value replaces the node in the tree. If the op returns none then the node remains unaltered.
// If any node was replaced then the whole process is run again until the op returns none for all nodes in the tree.
//
// See:
//   - [Rewrite] for a simpler non index non polymorphic version
//   - [RewriteI] for an index aware version
func RewriteI[I, A, RET any, RW TReadWrite, DIR, ERR any, ORET TReturnOne, OERR any](children Optic[I, A, A, A, A, RET, RW, DIR, ERR], op OperationI[*PathNode[I], A, mo.Option[A], ORET, OERR]) Optic[Void, A, A, A, A, ReturnOne, ReadOnly, UniDir, CompositionTree[ERR, OERR]] {

	return CombiGetter[CompositionTree[ERR, OERR], Void, A, A, A, A](
		func(ctx context.Context, source A) (Void, A, error) {

			for {
				var modified bool
				var err error
				source, modified, err = rewriteImpl(ctx, nil, source, children, op)
				err = JoinCtxErr(ctx, err)
				if err != nil {
					var a A
					return Void{}, a, err
				}

				if !modified {
					break
				}
			}

			return Void{}, source, ctx.Err()

		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Rewrite{
					OpticTypeExpr: ot,
					Children:      children.AsExpr(),
					RewriteOp:     op.AsExpr(),
				}
			},
			children,
			op,
		),
	)
}
