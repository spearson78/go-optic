package optic

// The Bd re-constrain converts an Optic with a complex but BiDir direction type to BiDir
//
// DIR TBiDir --> BiDir
//
// See:
//   - Bd for a re-constrain to BiDir
//   - Ud for a re-constrain to UniDir
//   - DirL for a re-constrain to the left side of a CompositionTree
//   - DirR for a re-constrain to the Right side of a CompositionTree
//   - DirSwap for a re-constrain that swaps the components of a CompositionTree
//   - DirSwap for a re-constrain that swaps the components of a CompositionTree
//   - DirSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - DirSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - DirMerge for a re-constrain that merges identical components of a CompositionTree
//   - DirMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - DirMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - DirTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - DirTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - DirTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func Bd[I, S, T, A, B any, RET any, RW TReadWrite, DIR TBiDir, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, T, A, B, RET, RW, BiDir, ERR] {
	return UnsafeReconstrain[RET, RW, BiDir, ERR](o)
}

// The Ud re-constrain converts an Optic with any direction type to UniDir
//
// DIR any --> UniDir
//
// See:
//   - Bd for a re-constrain to BiDir
//   - Ud for a re-constrain to UniDir
//   - DirL for a re-constrain to the left side of a CompositionTree
//   - DirR for a re-constrain to the Right side of a CompositionTree
//   - DirSwap for a re-constrain that swaps the components of a CompositionTree
//   - DirSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - DirSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - DirMerge for a re-constrain that merges identical components of a CompositionTree
//   - DirMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - DirMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - DirTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - DirTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - DirTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func Ud[I, S, T, A, B any, RET any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, T, A, B, RET, RW, UniDir, ERR] {
	return UnsafeReconstrain[RET, RW, UniDir, ERR](o)
}

// The DirL re-constrain converts an Optic's direction type to the left side of the [CompositionTree] where the right side is BiDir.
//
// CompositionTree[DIRL any, DIRR TBiDir] --> DIRL
//
// See:
//   - Bd for a re-constrain to BiDir
//   - Ud for a re-constrain to UniDir
//   - DirL for a re-constrain to the left side of a CompositionTree
//   - DirR for a re-constrain to the Right side of a CompositionTree
//   - DirSwap for a re-constrain that swaps the components of a CompositionTree
//   - DirSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - DirSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - DirMerge for a re-constrain that merges identical components of a CompositionTree
//   - DirMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - DirMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - DirTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - DirTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - DirTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func DirL[I, S, T, A, B any, RET any, RW any, DIRL any, DIRR TBiDir, ERR any](o Optic[I, S, T, A, B, RET, RW, CompositionTree[DIRL, DIRR], ERR]) Optic[I, S, T, A, B, RET, RW, DIRL, ERR] {
	return UnsafeReconstrain[RET, RW, DIRL, ERR](o)
}

// The DirR re-constrain converts an Optic's direction type to the right side of the [CompositionTree] where the left side is BiDir.
//
// CompositionTree[DIRL TBiDir, DIRR any] --> DIRR
//
// See:
//   - Bd for a re-constrain to BiDir
//   - Ud for a re-constrain to UniDir
//   - DirL for a re-constrain to the left side of a CompositionTree
//   - DirR for a re-constrain to the Right side of a CompositionTree
//   - DirSwap for a re-constrain that swaps the components of a CompositionTree
//   - DirSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - DirSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - DirMerge for a re-constrain that merges identical components of a CompositionTree
//   - DirMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - DirMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - DirTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - DirTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - DirTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func DirR[I, S, T, A, B any, RET any, RW any, DIRL TBiDir, DIRR any, ERR any](o Optic[I, S, T, A, B, RET, RW, CompositionTree[DIRL, DIRR], ERR]) Optic[I, S, T, A, B, RET, RW, DIRR, ERR] {
	return UnsafeReconstrain[RET, RW, DIRR, ERR](o)
}

// The DirSwap re-constrain swaps the left and right side of an Optic's direction type [CompositionTree].
//
// CompositionTree[DIRL, DIRR] --> CompositionTree[DIRR, DIRL]
//
// See:
//   - Bd for a re-constrain to BiDir
//   - Ud for a re-constrain to UniDir
//   - DirL for a re-constrain to the left side of a CompositionTree
//   - DirR for a re-constrain to the Right side of a CompositionTree
//   - DirSwap for a re-constrain that swaps the components of a CompositionTree
//   - DirSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - DirSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - DirMerge for a re-constrain that merges identical components of a CompositionTree
//   - DirMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - DirMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - DirTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - DirTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - DirTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func DirSwap[I, S, T, A, B any, RET any, RW any, DIRL any, DIRR any, ERR any](o Optic[I, S, T, A, B, RET, RW, CompositionTree[DIRL, DIRR], ERR]) Optic[I, S, T, A, B, RET, RW, CompositionTree[DIRR, DIRL], ERR] {
	return UnsafeReconstrain[RET, RW, CompositionTree[DIRR, DIRL], ERR](o)
}

// The DirSwapL re-constrain swaps the left and right side of the left side of a nested [CompositionTree].
//
// CompositionTree[CompositionTree[DIR1, DIR2], DIR3] --> CompositionTree[CompositionTree[DIR2, DIR1], DIR3]
//
// See:
//   - Bd for a re-constrain to BiDir
//   - Ud for a re-constrain to UniDir
//   - DirL for a re-constrain to the left side of a CompositionTree
//   - DirR for a re-constrain to the Right side of a CompositionTree
//   - DirSwap for a re-constrain that swaps the components of a CompositionTree
//   - DirSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - DirSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - DirMerge for a re-constrain that merges identical components of a CompositionTree
//   - DirMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - DirMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - DirTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - DirTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - DirTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func DirSwapL[I, S, T, A, B any, RET any, RW any, DIR1 any, DIR2 any, DIR3 any, ERR any](o Optic[I, S, T, A, B, RET, RW, CompositionTree[CompositionTree[DIR1, DIR2], DIR3], ERR]) Optic[I, S, T, A, B, RET, RW, CompositionTree[CompositionTree[DIR2, DIR1], DIR3], ERR] {
	return UnsafeReconstrain[RET, RW, CompositionTree[CompositionTree[DIR2, DIR1], DIR3], ERR](o)
}

// The DirSwapR re-constrain swaps the left and right side of the right side of a nested [CompositionTree].
//
// CompositionTree[DIR1, CompositionTree[DIR2, DIR3]] --> CompositionTree[DIR1, CompositionTree[DIR3, DIR2]]
//
// See:
//   - Bd for a re-constrain to BiDir
//   - Ud for a re-constrain to UniDir
//   - DirL for a re-constrain to the left side of a CompositionTree
//   - DirR for a re-constrain to the Right side of a CompositionTree
//   - DirSwap for a re-constrain that swaps the components of a CompositionTree
//   - DirSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - DirSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - DirSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - DirSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - DirMerge for a re-constrain that merges identical components of a CompositionTree
//   - DirMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - DirMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - DirTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - DirTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - DirTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func DirSwapR[I, S, T, A, B any, RET any, RW any, DIR1 any, DIR2 any, DIR3 any, ERR any](o Optic[I, S, T, A, B, RET, RW, CompositionTree[DIR1, CompositionTree[DIR2, DIR3]], ERR]) Optic[I, S, T, A, B, RET, RW, CompositionTree[DIR1, CompositionTree[DIR3, DIR3]], ERR] {
	return UnsafeReconstrain[RET, RW, CompositionTree[DIR1, CompositionTree[DIR3, DIR3]], ERR](o)
}

// The DirMerge re-constrain combines identical left and right side direction types of a composition tree..
//
// CompositionTree[DIR, DIR] --> DIR
//
// See:
//   - Bd for a re-constrain to BiDir
//   - Ud for a re-constrain to UniDir
//   - DirL for a re-constrain to the left side of a CompositionTree
//   - DirR for a re-constrain to the Right side of a CompositionTree
//   - DirSwap for a re-constrain that swaps the components of a CompositionTree
//   - DirSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - DirSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - DirSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - DirMerge for a re-constrain that merges identical components of a CompositionTree
//   - DirMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - DirMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - DirTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - DirTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - DirTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func DirMerge[I, S, T, A, B any, RET any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, RW, CompositionTree[DIR, DIR], ERR]) Optic[I, S, T, A, B, RET, RW, DIR, ERR] {
	return UnsafeReconstrain[RET, RW, DIR, ERR](o)
}

// The DirMergeL re-constrain combines identical left and right side of the left side of a nested {CompositionTree]
//
// CompositionTree[CompositionTree[DIR1, DIR1], DIR2] --> CompositionTree[DIR1, DIR2]
//
// See:
//   - Bd for a re-constrain to BiDir
//   - Ud for a re-constrain to UniDir
//   - DirL for a re-constrain to the left side of a CompositionTree
//   - DirR for a re-constrain to the Right side of a CompositionTree
//   - DirSwap for a re-constrain that swaps the components of a CompositionTree
//   - DirSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - DirSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - DirMerge for a re-constrain that merges identical components of a CompositionTree
//   - DirMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - DirMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - DirTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - DirTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - DirTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func DirMergeL[I, S, T, A, B any, RET any, RW any, DIR1 any, DIR2 any, ERR any](o Optic[I, S, T, A, B, RET, RW, CompositionTree[CompositionTree[DIR1, DIR1], DIR2], ERR]) Optic[I, S, T, A, B, RET, RW, CompositionTree[DIR1, DIR2], ERR] {
	return UnsafeReconstrain[RET, RW, CompositionTree[DIR1, DIR2], ERR](o)
}

// The DirMergeR re-constrain combines identical left and right side of the right side of a nested {CompositionTree]
//
// CompositionTree[DIR1, CompositionTree[DIR2, DIR2]] --> CompositionTree[DIR1, DIR2]
//
// See:
//   - Bd for a re-constrain to BiDir
//   - Ud for a re-constrain to UniDir
//   - DirL for a re-constrain to the left side of a CompositionTree
//   - DirR for a re-constrain to the Right side of a CompositionTree
//   - DirSwap for a re-constrain that swaps the components of a CompositionTree
//   - DirSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - DirSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - DirMerge for a re-constrain that merges identical components of a CompositionTree
//   - DirMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - DirMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - DirTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - DirTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - DirTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func DirMergeR[I, S, T, A, B any, RET any, RW any, DIR1 any, DIR2 any, ERR any](o Optic[I, S, T, A, B, RET, RW, CompositionTree[DIR1, CompositionTree[DIR2, DIR2]], ERR]) Optic[I, S, T, A, B, RET, RW, CompositionTree[DIR1, DIR2], ERR] {
	return UnsafeReconstrain[RET, RW, CompositionTree[DIR1, DIR2], ERR](o)
}

// The DirTrans re-constrain transposes the middle direction type of a nested [CompositionTree].
//
// CompositionTree[CompositionTree[ERR1, ERR2], CompositionTree[ERR3, ERR4] --> CompositionTree[CompositionTree[ERR1, ERR3], CompositionTree[ERR2, ERR4]
//
// See:
//   - Bd for a re-constrain to BiDir
//   - Ud for a re-constrain to UniDir
//   - DirL for a re-constrain to the left side of a CompositionTree
//   - DirR for a re-constrain to the Right side of a CompositionTree
//   - DirSwap for a re-constrain that swaps the components of a CompositionTree
//   - DirSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - DirSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - DirMerge for a re-constrain that merges identical components of a CompositionTree
//   - DirMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - DirMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - DirTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - DirTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - DirTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func DirTrans[I, S, T, A, B any, RET any, RW any, DIR1 any, DIR2 any, DIR3 any, DIR4 any, ERR any](o Optic[I, S, T, A, B, RET, RW, CompositionTree[CompositionTree[DIR1, DIR2], CompositionTree[DIR3, DIR4]], ERR]) Optic[I, S, T, A, B, RET, RW, CompositionTree[CompositionTree[DIR1, DIR3], CompositionTree[DIR2, DIR4]], ERR] {
	return UnsafeReconstrain[RET, RW, CompositionTree[CompositionTree[DIR1, DIR3], CompositionTree[DIR2, DIR4]], ERR](o)
}

// The DirTransL re-constrain transposes the right element with the right child of a nested composition tree to the left.
//
// CompositionTree[CompositionTree[DIR1, DIR2], DIR3] --> CompositionTree[CompositionTree[DIR1, DIR3], DIR2]
// This has the effect of removing a level of nested from the DIR2 path of the composition tree.
//
// See:
//   - Bd for a re-constrain to BiDir
//   - Ud for a re-constrain to UniDir
//   - DirL for a re-constrain to the left side of a CompositionTree
//   - DirR for a re-constrain to the Right side of a CompositionTree
//   - DirSwap for a re-constrain that swaps the components of a CompositionTree
//   - DirSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - DirSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - DirMerge for a re-constrain that merges identical components of a CompositionTree
//   - DirMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - DirMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - DirTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - DirTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - DirTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func DirTransL[I, S, T, A, B any, RET any, RW any, DIR1 any, DIR2 any, DIR3 any, ERR any](o Optic[I, S, T, A, B, RET, RW, CompositionTree[CompositionTree[DIR1, DIR2], DIR3], ERR]) Optic[I, S, T, A, B, RET, RW, CompositionTree[CompositionTree[DIR1, DIR3], DIR2], ERR] {
	return UnsafeReconstrain[RET, RW, CompositionTree[CompositionTree[DIR1, DIR3], DIR2], ERR](o)
}

// The DirTransR re-constrain transposes the left element with the left child of a nested composition tree to the right.
//
// CompositionTree[DIR1, CompositionTree[DIR2, DIR3]] --> CompositionTree[DIR2, CompositionTree[DIR1, DIR3]]
// This has the effect of removing a level of nested from the DIR2 path of the composition tree.
//
// See:
//   - Bd for a re-constrain to BiDir
//   - Ud for a re-constrain to UniDir
//   - DirL for a re-constrain to the left side of a CompositionTree
//   - DirR for a re-constrain to the Right side of a CompositionTree
//   - DirSwap for a re-constrain that swaps the components of a CompositionTree
//   - DirSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - DirSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - DirSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - DirMerge for a re-constrain that merges identical components of a CompositionTree
//   - DirMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - DirMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - DirTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - DirTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - DirTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func DirTransR[I, S, T, A, B any, RET any, RW any, DIR1 any, DIR2 any, DIR3 any, ERR any](o Optic[I, S, T, A, B, RET, RW, CompositionTree[DIR1, CompositionTree[DIR2, DIR3]], ERR]) Optic[I, S, T, A, B, RET, RW, CompositionTree[DIR2, CompositionTree[DIR1, DIR3]], ERR] {
	return UnsafeReconstrain[RET, RW, CompositionTree[DIR2, CompositionTree[DIR1, DIR3]], ERR](o)
}
