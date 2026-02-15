package optic

// The Ro re-constrain converts an Optic with any read/write type to ReadOnly
//
// RW any --> ReadOnly
//
// See:
//   - Ro for a re-constrain to  ReadOnly
//   - Rw for a re-constrain to ReadWrite
//   - RwL for a re-constrain to the left side of a CompositionTree
//   - RwR for a re-constrain to the Right side of a CompositionTree
//   - RwSwap for a re-constrain that swaps the components of a CompositionTree
//   - RwSwap for a re-constrain that swaps the components of a CompositionTree
//   - RwSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RwSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RwMerge for a re-constrain that merges identical components of a CompositionTree
//   - RwMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RwMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RwTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RwTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RwTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func Ro[I, S, T, A, B any, RET any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, T, A, B, RET, ReadOnly, UniDir, ERR] {
	return UnsafeReconstrain[RET, ReadOnly, UniDir, ERR](o)
}

// The Rw re-constrain converts an Optic with a complex but ReadWrite read/write type to ReadWrite
//
// RW TReadWrite --> ReadWrite
//
// See:
//   - Ro for a re-constrain to  ReadOnly
//   - Rw for a re-constrain to ReadWrite
//   - RwL for a re-constrain to the left side of a CompositionTree
//   - RwR for a re-constrain to the Right side of a CompositionTree
//   - RwSwap for a re-constrain that swaps the components of a CompositionTree
//   - RwSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RwSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RwMerge for a re-constrain that merges identical components of a CompositionTree
//   - RwMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RwMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RwTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RwTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RwTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func Rw[I, S, T, A, B any, RET any, RW TReadWrite, DIR any, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, T, A, B, RET, ReadWrite, DIR, ERR] {
	return UnsafeReconstrain[RET, ReadWrite, DIR, ERR](o)
}

// The RwL re-constrain converts an Optic's read/write type to the left side of the [CompositionTree] where the right side is ReadWrite.
//
// CompositionTree[RWL any, RWR TReadOnly] --> RWL
//
// See:
//   - Ro for a re-constrain to  ReadOnly
//   - Rw for a re-constrain to ReadWrite
//   - RwL for a re-constrain to the left side of a CompositionTree
//   - RwR for a re-constrain to the Right side of a CompositionTree
//   - RwSwap for a re-constrain that swaps the components of a CompositionTree
//   - RwSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RwSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RwMerge for a re-constrain that merges identical components of a CompositionTree
//   - RwMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RwMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RwTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RwTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RwTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RwL[I, S, T, A, B any, RET any, RWL any, RWR TReadWrite, DIR any, ERR any](o Optic[I, S, T, A, B, RET, CompositionTree[RWL, RWR], DIR, ERR]) Optic[I, S, T, A, B, RET, RWL, DIR, ERR] {
	return UnsafeReconstrain[RET, RWL, DIR, ERR](o)
}

// The RwR re-constrain converts an Optic's read/write type to the right side of the [CompositionTree] where the left side is ReadWrite.
//
// CompositionTree[RWL TReadWrite, RWR any] --> RWR
//
// See:
//   - Ro for a re-constrain to  ReadOnly
//   - Rw for a re-constrain to ReadWrite
//   - RwL for a re-constrain to the left side of a CompositionTree
//   - RwR for a re-constrain to the Right side of a CompositionTree
//   - RwSwap for a re-constrain that swaps the components of a CompositionTree
//   - RwSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RwSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RwMerge for a re-constrain that merges identical components of a CompositionTree
//   - RwMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RwMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RwTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RwTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RwTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RwR[I, S, T, A, B any, RET any, RWL TReadWrite, RWR any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, CompositionTree[RWL, RWR], DIR, ERR]) Optic[I, S, T, A, B, RET, RWR, DIR, ERR] {
	return UnsafeReconstrain[RET, RWR, DIR, ERR](o)
}

// The RwSwap re-constrain swaps the left and right side of an Optic's read/write type [CompositionTree].
//
// CompositionTree[RWL, RWR] --> CompositionTree[RWR, RWL]
//
// See:
//   - Ro for a re-constrain to  ReadOnly
//   - Rw for a re-constrain to ReadWrite
//   - RwL for a re-constrain to the left side of a CompositionTree
//   - RwR for a re-constrain to the Right side of a CompositionTree
//   - RwSwap for a re-constrain that swaps the components of a CompositionTree
//   - RwSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RwSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RwMerge for a re-constrain that merges identical components of a CompositionTree
//   - RwMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RwMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RwTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RwTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RwTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RwSwap[I, S, T, A, B any, RET any, RWL any, RWR any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, CompositionTree[RWL, RWR], DIR, ERR]) Optic[I, S, T, A, B, RET, CompositionTree[RWR, RWL], DIR, ERR] {
	return UnsafeReconstrain[RET, CompositionTree[RWR, RWL], DIR, ERR](o)
}

// The RwSwapL re-constrain swaps the left and right side of the left side of a nested [CompositionTree].
//
// CompositionTree[CompositionTree[RW1, RW2], RW3] --> CompositionTree[CompositionTree[RW2, RW1], RW3]
//
// See:
//   - Ro for a re-constrain to  ReadOnly
//   - Rw for a re-constrain to ReadWrite
//   - RwL for a re-constrain to the left side of a CompositionTree
//   - RwR for a re-constrain to the Right side of a CompositionTree
//   - RwSwap for a re-constrain that swaps the components of a CompositionTree
//   - RwSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RwSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RwMerge for a re-constrain that merges identical components of a CompositionTree
//   - RwMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RwMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RwTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RwTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RwTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RwSwapL[I, S, T, A, B any, RET any, RW1 any, RW2 any, RW3 any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, CompositionTree[CompositionTree[RW1, RW2], RW3], DIR, ERR]) Optic[I, S, T, A, B, RET, CompositionTree[CompositionTree[RW2, RW1], RW3], DIR, ERR] {
	return UnsafeReconstrain[RET, CompositionTree[CompositionTree[RW2, RW1], RW3], DIR, ERR](o)
}

// The RwSwapR re-constrain swaps the left and right side of the right side of a nested [CompositionTree].
//
// CompositionTree[RW1, CompositionTree[RW2, RW3]] --> CompositionTree[RW1, CompositionTree[RW3, RW2]]
//
// See:
//   - Ro for a re-constrain to  ReadOnly
//   - Rw for a re-constrain to ReadWrite
//   - RwL for a re-constrain to the left side of a CompositionTree
//   - RwR for a re-constrain to the Right side of a CompositionTree
//   - RwSwap for a re-constrain that swaps the components of a CompositionTree
//   - RwSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RwSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RwSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RwMerge for a re-constrain that merges identical components of a CompositionTree
//   - RwMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RwMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RwTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RwTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RwTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RwSwapR[I, S, T, A, B any, RET any, RW1 any, RW2 any, RW3 any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, CompositionTree[RW1, CompositionTree[RW2, RW3]], DIR, ERR]) Optic[I, S, T, A, B, RET, CompositionTree[RW1, CompositionTree[RW3, RW2]], DIR, ERR] {
	return UnsafeReconstrain[RET, CompositionTree[RW1, CompositionTree[RW3, RW2]], DIR, ERR](o)
}

// The RwMerge re-constrain combines identical left and right side read/write types of a composition tree..
//
// CompositionTree[RW, RW] --> RW
//
// See:
//   - Ro for a re-constrain to  ReadOnly
//   - Rw for a re-constrain to ReadWrite
//   - RwL for a re-constrain to the left side of a CompositionTree
//   - RwR for a re-constrain to the Right side of a CompositionTree
//   - RwSwap for a re-constrain that swaps the components of a CompositionTree
//   - RwSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RwSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RwMerge for a re-constrain that merges identical components of a CompositionTree
//   - RwMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RwMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RwTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RwTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RwTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RwMerge[I, S, T, A, B any, RET any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, CompositionTree[RW, RW], DIR, ERR]) Optic[I, S, T, A, B, RET, RW, DIR, ERR] {
	return UnsafeReconstrain[RET, RW, DIR, ERR](o)
}

// The RwMergeL re-constrain combines identical left and right side of the left side of a nested {CompositionTree]
//
// CompositionTree[CompositionTree[RW1, RW1], RW2] --> CompositionTree[RW1, RW2]
//
// See:
//   - Ro for a re-constrain to  ReadOnly
//   - Rw for a re-constrain to ReadWrite
//   - RwL for a re-constrain to the left side of a CompositionTree
//   - RwR for a re-constrain to the Right side of a CompositionTree
//   - RwSwap for a re-constrain that swaps the components of a CompositionTree
//   - RwSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RwSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RwMerge for a re-constrain that merges identical components of a CompositionTree
//   - RwMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RwMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RwTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RwTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RwTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RwMergeL[I, S, T, A, B any, RET any, RW1 any, RW2 any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, CompositionTree[CompositionTree[RW1, RW1], RW2], DIR, ERR]) Optic[I, S, T, A, B, RET, CompositionTree[RW1, RW2], DIR, ERR] {
	return UnsafeReconstrain[RET, CompositionTree[RW1, RW2], DIR, ERR](o)
}

// The RwMergeR re-constrain combines identical left and right side of the right side of a nested {CompositionTree]
//
// CompositionTree[RW1, CompositionTree[RW2, RW2]] --> CompositionTree[RW1, RW2]
//
// See:
//   - Ro for a re-constrain to  ReadOnly
//   - Rw for a re-constrain to ReadWrite
//   - RwL for a re-constrain to the left side of a CompositionTree
//   - RwR for a re-constrain to the Right side of a CompositionTree
//   - RwSwap for a re-constrain that swaps the components of a CompositionTree
//   - RwSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RwSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RwMerge for a re-constrain that merges identical components of a CompositionTree
//   - RwMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RwMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RwMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RwTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RwTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RwTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RwMergeR[I, S, T, A, B any, RET any, RW1 any, RW2 any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, CompositionTree[RW1, CompositionTree[RW2, RW2]], DIR, ERR]) Optic[I, S, T, A, B, RET, CompositionTree[RW1, RW2], DIR, ERR] {
	return UnsafeReconstrain[RET, CompositionTree[RW1, RW2], DIR, ERR](o)
}

// The RwTrans re-constrain transposes the middle read/write type of a nested [CompositionTree].
//
// CompositionTree[CompositionTree[ERR1, ERR2], CompositionTree[ERR3, ERR4] --> CompositionTree[CompositionTree[ERR1, ERR3], CompositionTree[ERR2, ERR4]
//
// See:
//   - Ro for a re-constrain to  ReadOnly
//   - Rw for a re-constrain to ReadWrite
//   - RwL for a re-constrain to the left side of a CompositionTree
//   - RwR for a re-constrain to the Right side of a CompositionTree
//   - RwSwap for a re-constrain that swaps the components of a CompositionTree
//   - RwSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RwSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RwMerge for a re-constrain that merges identical components of a CompositionTree
//   - RwMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RwMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RwTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RwTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RwTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RwTrans[I, S, T, A, B any, RET any, RW1 any, RW2 any, RW3 any, RW4 any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, CompositionTree[CompositionTree[RW1, RW2], CompositionTree[RW3, RW4]], DIR, ERR]) Optic[I, S, T, A, B, RET, CompositionTree[CompositionTree[RW1, RW3], CompositionTree[RW2, RW4]], DIR, ERR] {
	return UnsafeReconstrain[RET, CompositionTree[CompositionTree[RW1, RW3], CompositionTree[RW2, RW4]], DIR, ERR](o)
}

// The RwTransL re-constrain transposes the right element with the right child of a nested composition tree to the left.
//
// CompositionTree[CompositionTree[RW1, RW2], RW3] --> CompositionTree[CompositionTree[RW1, RW3], RW2]
// This has the effect of removing a level of nested from the RW2 path of the composition tree.
//
// See:
//   - Ro for a re-constrain to  ReadOnly
//   - Rw for a re-constrain to ReadWrite
//   - RwL for a re-constrain to the left side of a CompositionTree
//   - RwR for a re-constrain to the Right side of a CompositionTree
//   - RwSwap for a re-constrain that swaps the components of a CompositionTree
//   - RwSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RwSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RwMerge for a re-constrain that merges identical components of a CompositionTree
//   - RwMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RwMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RwTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RwTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RwTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RwTransL[I, S, T, A, B any, RET any, RW1 any, RW2 any, RW3 any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, CompositionTree[CompositionTree[RW1, RW2], RW3], DIR, ERR]) Optic[I, S, T, A, B, RET, CompositionTree[CompositionTree[RW1, RW3], RW2], DIR, ERR] {
	return UnsafeReconstrain[RET, CompositionTree[CompositionTree[RW1, RW3], RW2], DIR, ERR](o)
}

// The RwTransR re-constrain transposes the left element with the left child of a nested composition tree to the right.
//
// CompositionTree[RW1, CompositionTree[RW2, RW3]] --> CompositionTree[RW2, CompositionTree[RW1, RW3]]
// This has the effect of removing a level of nested from the RW2 path of the composition tree.
//
// See:
//   - Ro for a re-constrain to  ReadOnly
//   - Rw for a re-constrain to ReadWrite
//   - RwL for a re-constrain to the left side of a CompositionTree
//   - RwR for a re-constrain to the Right side of a CompositionTree
//   - RwSwap for a re-constrain that swaps the components of a CompositionTree
//   - RwSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RwSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RwMerge for a re-constrain that merges identical components of a CompositionTree
//   - RwMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RwMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RwMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RwTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RwTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RwTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RwTransR[I, S, T, A, B any, RET any, RW1 any, RW2 any, RW3 any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, CompositionTree[RW1, CompositionTree[RW2, RW3]], DIR, ERR]) Optic[I, S, T, A, B, RET, CompositionTree[RW2, CompositionTree[RW1, RW3]], DIR, ERR] {
	return UnsafeReconstrain[RET, CompositionTree[RW2, CompositionTree[RW1, RW3]], DIR, ERR](o)
}
