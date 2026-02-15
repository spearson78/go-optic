package optic

import (
	"context"
	"errors"
	"iter"
)

var yieldAfterBreak = errors.New("yield called after break")

// Constructor for an Iteration optic. Iterations focus on multiple elements.
//
// The following constructors are available.
//   - [Iteration]
//   - [IterationE] error aware
//   - [IterationI] indexed
//   - [IterationIE] indexed,error aware
//
// An Iteration is constructed from 3 functions
//
//   - iter : should iterates over the values in the source. yielding errors for individual values as needed.
//   - lengthGetter: should efficiently return the number of elements that will be focused. Pass nil for a default implementation that calls the iter function.
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [Traversal] for a version that returns a [ReadWrite] optic.
func Iteration[S, A any](
	seq func(source S) iter.Seq[A],
	lengthGetter func(source S) int,
	exprDef ExpressionDef,
) Optic[int, S, S, A, A, ReturnMany, ReadOnly, UniDir, Pure] {
	return CombiIteration[ReturnMany, Pure, int, S, S, A, A](
		func(ctx context.Context, source S) SeqIE[int, A] {
			return func(yield func(ValueIE[int, A]) bool) {
				i := 0
				seq(source)(func(focus A) bool {
					ret := yield(ValIE(i, focus, nil))
					i++
					return ret
				})
			}
		},
		maybeUpgradeLengthGetter(lengthGetter),
		nil,
		IxMatchComparable[int](),
		exprDef,
	)
}

// Constructor for an Iteration optic. Iterations focus on multiple elements.
//
// The following constructors are available.
//   - [Iteration]
//   - [IterationE] error aware
//   - [IterationI] indexed
//   - [IterationIE] indexed,error aware
//
// An Iteration is constructed from 3 functions
//
//   - iter : should iterates over the values in the source. yielding errors for individual values as needed.
//   - lengthGetter: should efficiently return the number of elements that will be focused. Pass nil for a default implementation that calls the iter function.
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [TraversalE] for a version that returns a [ReadWrite] optic.
func IterationE[S, A any](
	seq func(ctx context.Context, source S) iter.Seq2[A, error],
	lengthGetter LengthGetterFunc[S],
	exprDef ExpressionDef,
) Optic[int, S, S, A, A, ReturnMany, ReadOnly, UniDir, Err] {
	return CombiIteration[ReturnMany, Err, int, S, S, A, A](
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
		nil,
		IxMatchComparable[int](),
		exprDef,
	)
}

// Constructor for an Iteration optic. Iterations focus on multiple elements.
//
// The following constructors are available.
//   - [Iteration]
//   - [IterationE] error aware
//   - [IterationI] indexed
//   - [IterationIE] indexed,error aware

// An indexed Iteration is constructed from 5 functions
//
//   - iter : should iterates over the values in the source. yielding errors for individual values as needed.
//   - lengthGetter: should efficiently return the number of elements that will be focused. Pass nil for a default implementation that calls the iter function.
//   - ixGet : should iterate over the values in the source with the given index.
//   - ixMatch : should return true if the 2 index values are equal. Pass nil for a default implementation that calls the iter function.
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [TraversalI] for a version that returns a [ReadWrite] optic.
func IterationI[I, S, A any](
	seq func(source S) iter.Seq2[I, A],
	lengthGetter func(source S) int,
	ixget func(source S, index I) iter.Seq2[I, A],
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, S, A, A, ReturnMany, ReadOnly, UniDir, Pure] {

	ixMatchFnc := ensureSimpleIxMatch(ixMatch)

	return CombiIteration[ReturnMany, Pure, I, S, S, A, A](
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(ValueIE[I, A]) bool) {
				seq(source)(func(index I, v A) bool {
					return yield(ValIE(index, v, ctx.Err()))
				})
			}
		},
		maybeUpgradeLengthGetter(lengthGetter),
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			return func(yield func(ValueIE[I, A]) bool) {
				ixget(source, index)(func(index I, v A) bool {
					return yield(ValIE(index, v, ctx.Err()))
				})
			}
		},
		ixMatchFnc,
		exprDef,
	)
}

// Constructor for an Iteration optic. Iterations focus on multiple elements.
//
// The following constructors are available.
//   - [Iteration]
//   - [IterationE] error aware
//   - [IterationI] indexed
//   - [IterationIE] indexed,error aware
//
// An indexed Iteration is constructed from 5 functions
//
//   - iter : should iterates over the values in the source. yielding errors for individual values as needed.
//   - lengthGetter: should efficiently return the number of elements that will be focused. Pass nil for a default implementation that calls the iter function.
//   - ixGet : should iterate over the values in the source with the given index.
//   - ixMatch : should return true if the 2 index values are equal. Pass nil for a default implementation that calls the iter function.
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [TraversalIE] for a version that returns a [ReadWrite] optic.
func IterationIE[I, S, A any](
	iterate func(ctx context.Context, source S) SeqIE[I, A],
	lengthGetter LengthGetterFunc[S],
	ixget func(ctx context.Context, index I, source S) SeqIE[I, A],
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, S, A, A, ReturnMany, ReadOnly, UniDir, Err] {

	return CombiIteration[ReturnMany, Err, I, S, S, A, A](
		iterate,
		lengthGetter,
		ixget,
		ixMatch,
		exprDef,
	)
}
