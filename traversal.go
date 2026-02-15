package optic

import (
	"context"
	"iter"

	_ "github.com/spearson78/go-optic/expr"
)

// Constructor for a Traversal optic. Traversals focus on multiple elements.
//
// The following Traversal constructors are available.
//   - [Traversal]
//   - [TraversalP] polymorphic
//   - [TraversalE] error aware
//   - [TraversalEP] polymorphic, error aware
//   - [TraversalI] indexed
//   - [TraversalIP] indexed,polymorphic
//   - [TraversalIE] indexed,error aware
//   - [TraversalIEP] indexed,polymorphic, error aware
//
// A traversal is constructed from 4 functions
//
//   - iter : should iterates over the values in the source.
//   - lengthGetter: should efficiently return the number of elements that will be focused. Pass nil for a default implementation that calls the iter function.
//   - modify : should call fmap on each value in the source and return a new instance of the source type containing the mapped values.
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [Iteration] for [ReadOnly] version.
func Traversal[S, A any](
	iter func(source S) iter.Seq[A],
	lengthGetter func(source S) int,
	modify func(fmap func(focus A) A, source S) S,
	exprDef ExpressionDef,
) Optic[int, S, S, A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraversalP[S, S, A, A](
		iter,
		lengthGetter,
		modify,
		exprDef,
	)
}

// Constructor for a Traversal optic. Traversals focus on multiple elements.
//
// The following Traversal constructors are available.
//   - [Traversal]
//   - [TraversalP] polymorphic
//   - [TraversalE] error aware
//   - [TraversalEP] polymorphic, error aware
//   - [TraversalI] indexed
//   - [TraversalIP] indexed,polymorphic
//   - [TraversalIE] indexed,error aware
//   - [TraversalIEP] indexed,polymorphic, error aware
//
// A traversal is constructed from 4 functions
//
//   - iter : should iterates over the values in the source.
//   - lengthGetter: should efficiently return the number of elements that will be focused. Pass nil for a default implementation that calls the iter function.
//   - modify : should call fmap on each value in the source and return a new instance of the source type containing the mapped values.
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [Iteration] for [ReadOnly] version.
func TraversalP[S, T, A, B any](
	seq func(source S) iter.Seq[A],
	lengthGetter func(source S) int,
	modify func(fmap func(focus A) B, source S) T,
	exprDef ExpressionDef,
) Optic[int, S, T, A, B, ReturnMany, ReadWrite, UniDir, Pure] {

	return CombiTraversal[ReturnMany, ReadWrite, Pure, int, S, T, A, B](
		func(ctx context.Context, source S) SeqIE[int, A] {
			return func(yield func(ValueIE[int, A]) bool) {
				i := 0
				seq(source)(func(v A) bool {
					ret := yield(ValIE(i, v, ctx.Err()))
					i++
					return ret
				})
			}
		},
		maybeUpgradeLengthGetter(lengthGetter),
		func(ctx context.Context, fmap func(index int, focus A) (B, error), source S) (ret T, retErr error) {
			i := 0
			defer handleAbortModify(&retErr)
			ret = modify(func(focus A) B {
				b, err := fmap(i, focus)
				abortModifyError(JoinCtxErr(ctx, err), &retErr)
				i++
				return b
			}, source)
			return
		},
		nil,
		IxMatchComparable[int](),
		exprDef,
	)
}

// Constructor for a Traversal optic. Traversals focus on multiple elements.
//
// The following Traversal constructors are available.
//   - [Traversal]
//   - [TraversalP] polymorphic
//   - [TraversalE] error aware
//   - [TraversalEP] polymorphic, error aware
//   - [TraversalI] indexed
//   - [TraversalIP] indexed,polymorphic
//   - [TraversalIE] indexed,error aware
//   - [TraversalIEP] indexed,polymorphic, error aware
//
// A traversal is constructed from 4 functions
//
//   - iter : should iterates over the values in the source.
//   - lengthGetter: should efficiently return the number of elements that will be focused. Pass nil for a default implementation that calls the iter function.
//   - modify : should call fmap on each value in the source and return a new instance of the source type containing the mapped values.
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [Iteration] for [ReadOnly] version.
func TraversalE[S, A any](
	seq func(ctx context.Context, source S) iter.Seq2[A, error],
	lengthGetter LengthGetterFunc[S],
	modify func(ctx context.Context, fmap func(focus A) (A, error), source S) (S, error),
	exprDef ExpressionDef,
) Optic[int, S, S, A, A, ReturnMany, ReadWrite, UniDir, Err] {
	return TraversalEP[S, S, A, A](
		seq,
		lengthGetter,
		modify,
		exprDef,
	)
}

// Constructor for a Traversal optic. Traversals focus on multiple elements.
//
// The following Traversal constructors are available.
//   - [Traversal]
//   - [TraversalP] polymorphic
//   - [TraversalE] error aware
//   - [TraversalEP] polymorphic, error aware
//   - [TraversalI] indexed
//   - [TraversalIP] indexed,polymorphic
//   - [TraversalIE] indexed,error aware
//   - [TraversalIEP] indexed,polymorphic, error aware
//
// A traversal is constructed from 4 functions
//
//   - iter : should iterates over the values in the source.
//   - lengthGetter: should efficiently return the number of elements that will be focused. Pass nil for a default implementation that calls the iter function.
//   - modify : should call fmap on each value in the source and return a new instance of the source type containing the mapped values.
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [Iteration] for [ReadOnly] version.
func TraversalEP[S, T, A, B any](
	seq func(ctx context.Context, source S) iter.Seq2[A, error],
	lengthGetter LengthGetterFunc[S],
	modify func(ctx context.Context, fmap func(focus A) (B, error), source S) (T, error),
	exprDef ExpressionDef,
) Optic[int, S, T, A, B, ReturnMany, ReadWrite, UniDir, Err] {

	return CombiTraversal[ReturnMany, ReadWrite, Err, int, S, T, A, B](
		func(ctx context.Context, source S) SeqIE[int, A] {
			return func(yield func(ValueIE[int, A]) bool) {
				i := 0
				seq(ctx, source)(func(v A, err error) bool {
					ret := yield(ValIE(i, v, JoinCtxErr(ctx, err)))
					i++
					return ret
				})
			}
		},
		lengthGetter,
		func(ctx context.Context, fmap func(index int, focus A) (B, error), source S) (T, error) {
			i := 0
			return modify(ctx, func(focus A) (B, error) {
				b, err := fmap(i, focus)
				i++
				return b, err
			}, source)
		},
		nil,
		IxMatchComparable[int](),
		exprDef,
	)
}

// Constructor for a Traversal optic. Traversals focus on multiple elements.
//
// The following Traversal constructors are available.
//   - [Traversal]
//   - [TraversalP] polymorphic
//   - [TraversalE] error aware
//   - [TraversalEP] polymorphic, error aware
//   - [TraversalI] indexed
//   - [TraversalIP] indexed,polymorphic
//   - [TraversalIE] indexed,error aware
//   - [TraversalIEP] indexed,polymorphic, error aware
//
// An indexed traversal is constructed from 4 functions
//
//   - iter : should iterates over the values in the source.
//   - lengthGetter: should efficiently return the number of elements that will be focused. Pass nil for a default implementation that calls the iter function.
//   - modify : should call fmap on each value in the source and return a new instance of the source type containing the mapped values.
//   - ixGet : should iterate over the values in the source with the given index.
//   - ixMatch : should return true if the 2 index values are equal. Pass nil for a default implementation that calls the iter function.
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [Iteration] for [ReadOnly] version.
func TraversalI[I, S, A any](
	iter func(source S) SeqI[I, A],
	lengthGetter func(source S) int,
	modify func(fmap func(index I, focus A) A, source S) S,
	ixget func(source S, index I) iter.Seq2[I, A],
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, S, A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraversalIP[I, S, S, A, A](
		iter,
		lengthGetter,
		modify,
		ixget,
		ixMatch,
		exprDef,
	)
}

// Constructor for a Traversal optic. Traversals focus on multiple elements.
//
// The following Traversal constructors are available.
//   - [Traversal]
//   - [TraversalP] polymorphic
//   - [TraversalE] error aware
//   - [TraversalEP] polymorphic, error aware
//   - [TraversalI] indexed
//   - [TraversalIP] indexed,polymorphic
//   - [TraversalIE] indexed,error aware
//   - [TraversalIEP] indexed,polymorphic, error aware
//
// An indexed traversal is constructed from 4 functions
//
//   - iter : should iterates over the values in the source.
//   - lengthGetter: should efficiently return the number of elements that will be focused. Pass nil for a default implementation that calls the iter function.
//   - modify : should call fmap on each value in the source and return a new instance of the source type containing the mapped values.
//   - ixGet : should iterate over the values in the source with the given index.
//   - ixMatch : should return true if the 2 index values are equal. Pass nil for a default implementation that calls the iter function.
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [Iteration] for [ReadOnly] version.
func TraversalIP[I, S, T, A, B any](
	seq func(source S) SeqI[I, A],
	lengthGetter func(source S) int,
	modify func(fmap func(index I, focus A) B, source S) T,
	ixget func(source S, index I) iter.Seq2[I, A],
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, T, A, B, ReturnMany, ReadWrite, UniDir, Pure] {

	return CombiTraversal[ReturnMany, ReadWrite, Pure, I, S, T, A, B](
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(ValueIE[I, A]) bool) {
				seq(source)(func(i I, v A) bool {
					return yield(ValIE(i, v, ctx.Err()))
				})
			}
		},
		maybeUpgradeLengthGetter(lengthGetter),
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (ret T, retErr error) {
			defer handleAbortModify(&retErr)
			ret = modify(func(index I, focus A) B {
				b, err := fmap(index, focus)
				abortModifyError(JoinCtxErr(ctx, err), &retErr)
				return b
			}, source)
			return
		},
		maybeUpgradeIxGet(ixget),
		ixMatch,
		exprDef,
	)
}

// Constructor for a Traversal optic. Traversals focus on multiple elements.
//
// The following Traversal constructors are available.
//   - [Traversal]
//   - [TraversalP] polymorphic
//   - [TraversalE] error aware
//   - [TraversalEP] polymorphic, error aware
//   - [TraversalI] indexed
//   - [TraversalIP] indexed,polymorphic
//   - [TraversalIE] indexed,error aware
//   - [TraversalIEP] indexed,polymorphic, error aware
//
// An indexed traversal is constructed from 4 functions
//
//   - iter : should iterates over the values in the source.
//   - lengthGetter: should efficiently return the number of elements that will be focused. Pass nil for a default implementation that calls the iter function.
//   - modify : should call fmap on each value in the source and return a new instance of the source type containing the mapped values.
//   - ixGet : should iterate over the values in the source with the given index.
//   - ixMatch : should return true if the 2 index values are equal. Pass nil for a default implementation that calls the iter function.
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [Iteration] for [ReadOnly] version.
func TraversalIE[I, S, A any](
	seq func(ctx context.Context, source S) SeqIE[I, A],
	lengthGetter LengthGetterFunc[S],
	modify func(ctx context.Context, fmap func(index I, focus A) (A, error), source S) (S, error),
	ixget func(ctx context.Context, index I, source S) SeqIE[I, A],
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, S, A, A, ReturnMany, ReadWrite, UniDir, Err] {
	return TraversalIEP[I, S, S, A, A](
		seq,
		lengthGetter,
		modify,
		ixget,
		ixMatch,
		exprDef,
	)
}

// Constructor for a Traversal optic. Traversals focus on multiple elements.
//
// The following Traversal constructors are available.
//   - [Traversal]
//   - [TraversalP] polymorphic
//   - [TraversalE] error aware
//   - [TraversalEP] polymorphic, error aware
//   - [TraversalI] indexed
//   - [TraversalIP] indexed,polymorphic
//   - [TraversalIE] indexed,error aware
//   - [TraversalIEP] indexed,polymorphic, error aware
//
// An indexed traversal is constructed from 4 functions
//
//   - iter : should iterates over the values in the source.
//   - lengthGetter: should efficiently return the number of elements that will be focused. Pass nil for a default implementation that calls the iter function.
//   - modify : should call fmap on each value in the source and return a new instance of the source type containing the mapped values.
//   - ixGet : should iterate over the values in the source with the given index.
//   - ixMatch : should return true if the 2 index values are equal. Pass nil for a default implementation that calls the iter function.
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [Iteration] for [ReadOnly] version.
func TraversalIEP[I, S, T, A, B any](
	seq func(ctx context.Context, source S) SeqIE[I, A],
	lengthGetter LengthGetterFunc[S],
	modify func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error),
	ixget func(ctx context.Context, index I, source S) SeqIE[I, A],
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, T, A, B, ReturnMany, ReadWrite, UniDir, Err] {

	return CombiTraversal[ReturnMany, ReadWrite, Err, I, S, T, A, B](
		seq,
		lengthGetter,
		modify,
		ixget,
		ixMatch,
		exprDef,
	)
}

func maybeUpgradeIxGet[I, S, A any](ixGet func(source S, index I) iter.Seq2[I, A]) IxGetterFunc[I, S, A] {
	if ixGet != nil {
		return func(ctx context.Context, index I, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				ixGet(source, index)(func(index I, focus A) bool {
					return yield(ValIE(index, focus, ctx.Err()))
				})
			}
		}
	}

	return nil
}
