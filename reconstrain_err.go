package optic

// The EPure re-constrain converts an Optic with a complex but Pure error type to Pure
//
// ERR TPure --> Pure
//
// See:
//   - EPure for a re-constrain to Pure
//   - EErr for a re-constrain to Err
//   - EErrL for a re-constrain to the left side of a CompositionTree
//   - EErrR for a re-constrain to the Right side of a CompositionTree
//   - EErrSwap for a re-constrain that swaps the components of a CompositionTree
//   - EErrSwap for a re-constrain that swaps the components of a CompositionTree
//   - EErrSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - EErrSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - EErrMerge for a re-constrain that merges identical components of a CompositionTree
//   - EErrMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - EErrMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - EErrTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - EErrTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - EErrTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func EPure[I, S, T, A, B any, RET any, RW any, DIR any, ERR TPure](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, T, A, B, RET, RW, DIR, Pure] {
	return UnsafeReconstrain[RET, RW, DIR, Pure](o)
}

// The EErr re-constrain converts an Optic with any error type to Err
//
// ERR any --> Err
//
// See:
//   - EPure for a re-constrain to Pure
//   - EErr for a re-constrain to Err
//   - EErrL for a re-constrain to the left side of a CompositionTree
//   - EErrR for a re-constrain to the Right side of a CompositionTree
//   - EErrSwap for a re-constrain that swaps the components of a CompositionTree
//   - EErrSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - EErrSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - EErrMerge for a re-constrain that merges identical components of a CompositionTree
//   - EErrMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - EErrMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - EErrTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - EErrTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - EErrTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func EErr[I, S, T, A, B any, RET any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, T, A, B, RET, RW, DIR, Err] {
	return UnsafeReconstrain[RET, RW, DIR, Err](o)
}

// The EErrL re-constrain converts an Optic's error type to the left side of the [CompositionTree] where the right side is Pure.
//
// CompositionTree[ERRL any, ERRR TPure] --> ERRL
//
// See:
//   - EPure for a re-constrain to Pure
//   - EErr for a re-constrain to Err
//   - EErrL for a re-constrain to the left side of a CompositionTree
//   - EErrR for a re-constrain to the Right side of a CompositionTree
//   - EErrSwap for a re-constrain that swaps the components of a CompositionTree
//   - EErrSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - EErrSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - EErrMerge for a re-constrain that merges identical components of a CompositionTree
//   - EErrMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - EErrMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - EErrTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - EErrTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - EErrTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func EErrL[I, S, T, A, B any, RET any, RW any, DIR any, ERRL any, ERRR TPure](o Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[ERRL, ERRR]]) Optic[I, S, T, A, B, RET, RW, DIR, ERRL] {
	return UnsafeReconstrain[RET, RW, DIR, ERRL](o)
}

// The EErrR re-constrain converts an Optic's error type to the right side of the [CompositionTree] where the left side is Pure.
//
// CompositionTree[ERRL TPure, ERRR any] --> ERRR
//
// See:
//   - EPure for a re-constrain to Pure
//   - EErr for a re-constrain to Err
//   - EErrL for a re-constrain to the left side of a CompositionTree
//   - EErrR for a re-constrain to the Right side of a CompositionTree
//   - EErrSwap for a re-constrain that swaps the components of a CompositionTree
//   - EErrSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - EErrSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - EErrMerge for a re-constrain that merges identical components of a CompositionTree
//   - EErrMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - EErrMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - EErrTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - EErrTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - EErrTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func EErrR[I, S, T, A, B any, RET any, RW any, DIR any, ERRL TPure, ERRR any](o Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[ERRL, ERRR]]) Optic[I, S, T, A, B, RET, RW, DIR, ERRR] {
	return UnsafeReconstrain[RET, RW, DIR, ERRR](o)
}

// The EErrSwap re-constrain swaps the left and right side of an Optic's error type [CompositionTree].
//
// CompositionTree[ERRL, ERRR] --> CompositionTree[ERRR, ERRL]
//
// See:
//   - EPure for a re-constrain to Pure
//   - EErr for a re-constrain to Err
//   - EErrL for a re-constrain to the left side of a CompositionTree
//   - EErrR for a re-constrain to the Right side of a CompositionTree
//   - EErrSwap for a re-constrain that swaps the components of a CompositionTree
//   - EErrSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - EErrSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - EErrMerge for a re-constrain that merges identical components of a CompositionTree
//   - EErrMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - EErrMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - EErrTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - EErrTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - EErrTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func EErrSwap[I, S, T, A, B any, RET any, RW any, DIR any, ERRL any, ERRR any](o Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[ERRL, ERRR]]) Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[ERRR, ERRL]] {
	return UnsafeReconstrain[RET, RW, DIR, CompositionTree[ERRR, ERRL]](o)
}

// The EErrSwapL re-constrain swaps the left and right side of the left side of a nested [CompositionTree].
//
// CompositionTree[CompositionTree[ERR1, ERR2], ERR3] --> CompositionTree[CompositionTree[ERR2, ERR1], ERR3]
//
// See:
//   - EPure for a re-constrain to Pure
//   - EErr for a re-constrain to Err
//   - EErrL for a re-constrain to the left side of a CompositionTree
//   - EErrR for a re-constrain to the Right side of a CompositionTree
//   - EErrSwap for a re-constrain that swaps the components of a CompositionTree
//   - EErrSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - EErrSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - EErrMerge for a re-constrain that merges identical components of a CompositionTree
//   - EErrMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - EErrMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - EErrTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - EErrTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - EErrTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func EErrSwapL[I, S, T, A, B any, RET any, RW any, DIR any, ERR1 any, ERR2 any, ERR3 any](o Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[CompositionTree[ERR1, ERR2], ERR3]]) Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[CompositionTree[ERR2, ERR1], ERR3]] {
	return UnsafeReconstrain[RET, RW, DIR, CompositionTree[CompositionTree[ERR2, ERR1], ERR3]](o)
}

// The EErrSwapR re-constrain swaps the left and right side of the right side of a nested [CompositionTree].
//
// CompositionTree[ERR1, CompositionTree[ERR2, ERR3]]] --> CompositionTree[ERR1, CompositionTree[ERR3, ERR2]]
//
// See:
//   - EPure for a re-constrain to Pure
//   - EErr for a re-constrain to Err
//   - EErrL for a re-constrain to the left side of a CompositionTree
//   - EErrR for a re-constrain to the Right side of a CompositionTree
//   - EErrSwap for a re-constrain that swaps the components of a CompositionTree
//   - EErrSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - EErrSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - EErrMerge for a re-constrain that merges identical components of a CompositionTree
//   - EErrMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - EErrMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - EErrTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - EErrTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - EErrTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func EErrSwapR[I, S, T, A, B any, RET any, RW any, DIR any, ERR1 any, ERR2 any, ERR3 any](o Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[ERR1, CompositionTree[ERR2, ERR3]]]) Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[ERR1, CompositionTree[ERR3, ERR2]]] {
	return UnsafeReconstrain[RET, RW, DIR, CompositionTree[ERR1, CompositionTree[ERR3, ERR2]]](o)
}

// The EErrMerge re-constrain combines identical left and right side error types of a composition tree..
//
// CompositionTree[ERR, ERR] --> ERR
//
// See:
//   - EPure for a re-constrain to Pure
//   - EErr for a re-constrain to Err
//   - EErrL for a re-constrain to the left side of a CompositionTree
//   - EErrR for a re-constrain to the Right side of a CompositionTree
//   - EErrSwap for a re-constrain that swaps the components of a CompositionTree
//   - EErrSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - EErrSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - EErrMerge for a re-constrain that merges identical components of a CompositionTree
//   - EErrMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - EErrMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - EErrTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - EErrTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - EErrTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func EErrMerge[I, S, T, A, B any, RET any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[ERR, ERR]]) Optic[I, S, T, A, B, RET, RW, DIR, ERR] {
	return UnsafeReconstrain[RET, RW, DIR, ERR](o)
}

// The EErrMergeL re-constrain combines identical left and right side of the left side of a nested {CompositionTree]
//
// CompositionTree[CompositionTree[ERR1, ERR1], ERR2] --> CompositionTree[ERR1, ERR2]
//
// See:
//   - EPure for a re-constrain to Pure
//   - EErr for a re-constrain to Err
//   - EErrL for a re-constrain to the left side of a CompositionTree
//   - EErrR for a re-constrain to the Right side of a CompositionTree
//   - EErrSwap for a re-constrain that swaps the components of a CompositionTree
//   - EErrSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - EErrSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - EErrMerge for a re-constrain that merges identical components of a CompositionTree
//   - EErrMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - EErrMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - EErrTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - EErrTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - EErrTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func EErrMergeL[I, S, T, A, B any, RET any, RW any, DIR any, ERR1 any, ERR2 any](o Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[CompositionTree[ERR1, ERR1], ERR2]]) Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[ERR1, ERR2]] {
	return UnsafeReconstrain[RET, RW, DIR, CompositionTree[ERR1, ERR2]](o)
}

// The EErrMergeR re-constrain combines identical left and right side of the right side of a nested {CompositionTree]
//
// CompositionTree[ERR1, CompositionTree[ERR2, ERR2]] --> CompositionTree[ERR1, ERR2]
//
// See:
//   - EPure for a re-constrain to Pure
//   - EErr for a re-constrain to Err
//   - EErrL for a re-constrain to the left side of a CompositionTree
//   - EErrR for a re-constrain to the Right side of a CompositionTree
//   - EErrSwap for a re-constrain that swaps the components of a CompositionTree
//   - EErrSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - EErrSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - EErrMerge for a re-constrain that merges identical components of a CompositionTree
//   - EErrMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - EErrMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - EErrTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - EErrTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - EErrTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func EErrMergeR[I, S, T, A, B any, RET any, RW any, DIR any, ERR1 any, ERR2 any](o Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[ERR1, CompositionTree[ERR2, ERR2]]]) Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[ERR1, ERR2]] {
	return UnsafeReconstrain[RET, RW, DIR, CompositionTree[ERR1, ERR2]](o)
}

// The EErrTrans re-constrain transposes the middle error type of a nested [CompositionTree].
//
// CompositionTree[CompositionTree[ERR1, ERR2], CompositionTree[ERR3, ERR4] --> CompositionTree[CompositionTree[ERR1, ERR3], CompositionTree[ERR2, ERR4]
//
// See:
//   - EPure for a re-constrain to Pure
//   - EErr for a re-constrain to Err
//   - EErrL for a re-constrain to the left side of a CompositionTree
//   - EErrR for a re-constrain to the Right side of a CompositionTree
//   - EErrSwap for a re-constrain that swaps the components of a CompositionTree
//   - EErrSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - EErrSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - EErrMerge for a re-constrain that merges identical components of a CompositionTree
//   - EErrMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - EErrMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - EErrTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - EErrTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - EErrTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func EErrTrans[I, S, T, A, B any, RET any, RW any, DIR any, ERR1 any, ERR2 any, ERR3 any, ERR4 any](o Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[CompositionTree[ERR1, ERR2], CompositionTree[ERR3, ERR4]]]) Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[CompositionTree[ERR1, ERR3], CompositionTree[ERR2, ERR4]]] {
	return UnsafeReconstrain[RET, RW, DIR, CompositionTree[CompositionTree[ERR1, ERR3], CompositionTree[ERR2, ERR4]]](o)
}

// The EErrTransL re-constrain transposes the right element with the right child of a nested composition tree to the left.
//
// CompositionTree[CompositionTree[ERR1, ERR2], ERR3] --> CompositionTree[CompositionTree[ERR1, ERR3], ERR2]
// This has the effect of removing a level of nested from the ERR2 path of the composition tree.
//
// See:
//   - EPure for a re-constrain to Pure
//   - EErr for a re-constrain to Err
//   - EErrL for a re-constrain to the left side of a CompositionTree
//   - EErrR for a re-constrain to the Right side of a CompositionTree
//   - EErrSwap for a re-constrain that swaps the components of a CompositionTree
//   - EErrSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - EErrSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - EErrMerge for a re-constrain that merges identical components of a CompositionTree
//   - EErrMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - EErrMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - EErrTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - EErrTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - EErrTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func EErrTransL[I, S, T, A, B any, RET any, RW any, DIR any, ERR1 any, ERR2 any, ERR3 any](o Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[CompositionTree[ERR1, ERR2], ERR3]]) Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[CompositionTree[ERR1, ERR3], ERR2]] {
	return UnsafeReconstrain[RET, RW, DIR, CompositionTree[CompositionTree[ERR1, ERR3], ERR2]](o)
}

// The EErrTransR re-constrain transposes the left element with the left child of a nested composition tree to the right.
//
// CompositionTree[ERR1, CompositionTree[ERR2, ERR3]] --> CompositionTree[ERR2, CompositionTree[ERR1, ERR3]]
// This has the effect of removing a level of nested from the ERR2 path of the composition tree.
//
// See:
//   - EPure for a re-constrain to Pure
//   - EErr for a re-constrain to Err
//   - EErrL for a re-constrain to the left side of a CompositionTree
//   - EErrR for a re-constrain to the Right side of a CompositionTree
//   - EErrSwap for a re-constrain that swaps the components of a CompositionTree
//   - EErrSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - EErrSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - EErrMerge for a re-constrain that merges identical components of a CompositionTree
//   - EErrMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - EErrMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - EErrTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - EErrTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - EErrTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func EErrTransR[I, S, T, A, B any, RET any, RW any, DIR any, ERR1 any, ERR2 any, ERR3 any](o Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[ERR1, CompositionTree[ERR2, ERR3]]]) Optic[I, S, T, A, B, RET, RW, DIR, CompositionTree[ERR2, CompositionTree[ERR1, ERR3]]] {
	return UnsafeReconstrain[RET, RW, DIR, CompositionTree[ERR2, CompositionTree[ERR1, ERR3]]](o)
}
