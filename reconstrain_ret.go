package optic

// The Ret1 re-constrain converts an Optic with a complex but ReturnOne return type to ReturnOne
//
// RET TReturnOne --> ReturnOne
//
// See:
//   - Ret1 for a re-constrain to  ReturnOne
//   - RetM for a re-constrain to ReturnMany
//   - RetL for a re-constrain to the left side of a CompositionTree
//   - RetR for a re-constrain to the Right side of a CompositionTree
//   - RetSwap for a re-constrain that swaps the components of a CompositionTree
//   - RetSwap for a re-constrain that swaps the components of a CompositionTree
//   - RetSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RetSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RetMerge for a re-constrain that merges identical components of a CompositionTree
//   - RetMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RetMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RetTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RetTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RetTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func Ret1[I, S, T, A, B any, RET TReturnOne, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, T, A, B, ReturnOne, RW, DIR, ERR] {
	return UnsafeReconstrain[ReturnOne, RW, DIR, ERR](o)
}

// The RetM re-constrain converts an Optic with any return type to ReturnMany
//
// RET any --> ReturnMany
//
// See:
//   - Ret1 for a re-constrain to  ReturnOne
//   - RetM for a re-constrain to ReturnMany
//   - RetL for a re-constrain to the left side of a CompositionTree
//   - RetR for a re-constrain to the Right side of a CompositionTree
//   - RetSwap for a re-constrain that swaps the components of a CompositionTree
//   - RetSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RetSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RetMerge for a re-constrain that merges identical components of a CompositionTree
//   - RetMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RetMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RetTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RetTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RetTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RetM[I, S, T, A, B any, RET any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, T, A, B, ReturnMany, RW, DIR, ERR] {
	return UnsafeReconstrain[ReturnMany, RW, DIR, ERR](o)
}

// The RetL re-constrain converts an Optic's return type to the left side of the [CompositionTree] where the right side is ReturnOne.
//
// CompositionTree[RETL any, RETR TReturnOne] --> RETL
//
// See:
//   - Ret1 for a re-constrain to  ReturnOne
//   - RetM for a re-constrain to ReturnMany
//   - RetL for a re-constrain to the left side of a CompositionTree
//   - RetR for a re-constrain to the Right side of a CompositionTree
//   - RetSwap for a re-constrain that swaps the components of a CompositionTree
//   - RetSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RetSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RetMerge for a re-constrain that merges identical components of a CompositionTree
//   - RetMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RetMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RetTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RetTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RetTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RetL[I, S, T, A, B any, RETL any, RETR TReturnOne, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, CompositionTree[RETL, RETR], RW, DIR, ERR]) Optic[I, S, T, A, B, RETL, RW, DIR, ERR] {
	return UnsafeReconstrain[RETL, RW, DIR, ERR](o)
}

// The RetR re-constrain converts an Optic's return type to the right side of the [CompositionTree] where the left side is ReturnOne.
//
// CompositionTree[RETL TReturnOne, RETR any] --> RETR
//
// See:
//   - Ret1 for a re-constrain to  ReturnOne
//   - RetM for a re-constrain to ReturnMany
//   - RetL for a re-constrain to the left side of a CompositionTree
//   - RetR for a re-constrain to the Right side of a CompositionTree
//   - RetSwap for a re-constrain that swaps the components of a CompositionTree
//   - RetSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RetSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RetMerge for a re-constrain that merges identical components of a CompositionTree
//   - RetMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RetMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RetTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RetTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RetTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RetR[I, S, T, A, B any, RETL TReturnOne, RETR any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, CompositionTree[RETL, RETR], RW, DIR, ERR]) Optic[I, S, T, A, B, RETR, RW, DIR, ERR] {
	return UnsafeReconstrain[RETR, RW, DIR, ERR](o)
}

// The RetSwap re-constrain swaps the left and right side of an Optic's return type [CompositionTree].
//
// CompositionTree[RETL, RETR] --> CompositionTree[RETR, RETL]
//
// See:
//   - Ret1 for a re-constrain to  ReturnOne
//   - RetM for a re-constrain to ReturnMany
//   - RetL for a re-constrain to the left side of a CompositionTree
//   - RetR for a re-constrain to the Right side of a CompositionTree
//   - RetSwap for a re-constrain that swaps the components of a CompositionTree
//   - RetSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RetSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RetMerge for a re-constrain that merges identical components of a CompositionTree
//   - RetMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RetMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RetTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RetTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RetTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RetSwap[I, S, T, A, B any, RETL any, RETR any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, CompositionTree[RETL, RETR], RW, DIR, ERR]) Optic[I, S, T, A, B, CompositionTree[RETR, RETL], RW, DIR, ERR] {
	return UnsafeReconstrain[CompositionTree[RETR, RETL], RW, DIR, ERR](o)
}

// The RetSwapL re-constrain swaps the left and right side of the left side of a nested [CompositionTree].
//
// CompositionTree[CompositionTree[RET1, RET2], RET3] --> CompositionTree[CompositionTree[RET2, RET1], RET3]
//
// See:
//   - Ret1 for a re-constrain to  ReturnOne
//   - RetM for a re-constrain to ReturnMany
//   - RetL for a re-constrain to the left side of a CompositionTree
//   - RetR for a re-constrain to the Right side of a CompositionTree
//   - RetSwap for a re-constrain that swaps the components of a CompositionTree
//   - RetSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RetSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RetMerge for a re-constrain that merges identical components of a CompositionTree
//   - RetMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RetMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RetTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RetTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RetTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RetSwapL[I, S, T, A, B any, RET1 any, RET2 any, RET3 any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, CompositionTree[CompositionTree[RET1, RET2], RET3], RW, DIR, ERR]) Optic[I, S, T, A, B, CompositionTree[CompositionTree[RET2, RET1], RET3], RW, DIR, ERR] {
	return UnsafeReconstrain[CompositionTree[CompositionTree[RET2, RET1], RET3], RW, DIR, ERR](o)
}

// The RetSwapR re-constrain swaps the left and right side of the right side of a nested [CompositionTree].
//
// CompositionTree[RET1, CompositionTree[RET2, RET3]] --> CompositionTree[RET1, CompositionTree[RET3, RET2]]
//
// See:
//   - Ret1 for a re-constrain to  ReturnOne
//   - RetM for a re-constrain to ReturnMany
//   - RetL for a re-constrain to the left side of a CompositionTree
//   - RetR for a re-constrain to the Right side of a CompositionTree
//   - RetSwap for a re-constrain that swaps the components of a CompositionTree
//   - RetSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RetSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RetMerge for a re-constrain that merges identical components of a CompositionTree
//   - RetMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RetMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RetTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RetTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RetTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RetSwapR[I, S, T, A, B any, RET1 any, RET2 any, RET3 any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, CompositionTree[RET1, CompositionTree[RET2, RET3]], RW, DIR, ERR]) Optic[I, S, T, A, B, CompositionTree[RET1, CompositionTree[RET3, RET2]], RW, DIR, ERR] {
	return UnsafeReconstrain[CompositionTree[RET1, CompositionTree[RET3, RET2]], RW, DIR, ERR](o)
}

// The RetMerge re-constrain combines identical left and right side return types of a composition tree..
//
// CompositionTree[RET, RET] --> RET
//
// See:
//   - Ret1 for a re-constrain to  ReturnOne
//   - RetM for a re-constrain to ReturnMany
//   - RetL for a re-constrain to the left side of a CompositionTree
//   - RetR for a re-constrain to the Right side of a CompositionTree
//   - RetSwap for a re-constrain that swaps the components of a CompositionTree
//   - RetSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RetSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RetMerge for a re-constrain that merges identical components of a CompositionTree
//   - RetMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RetMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RetTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RetTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RetTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RetMerge[I, S, T, A, B any, RET any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, CompositionTree[RET, RET], RW, DIR, ERR]) Optic[I, S, T, A, B, RET, RW, DIR, ERR] {
	return UnsafeReconstrain[RET, RW, DIR, ERR](o)
}

// The RetMergeL re-constrain combines identical left and right side of the left side of a nested {CompositionTree]
//
// CompositionTree[CompositionTree[RET1, RET1], RET2] --> CompositionTree[RET1, RET2]
//
// See:
//   - Ret1 for a re-constrain to  ReturnOne
//   - RetM for a re-constrain to ReturnMany
//   - RetL for a re-constrain to the left side of a CompositionTree
//   - RetR for a re-constrain to the Right side of a CompositionTree
//   - RetSwap for a re-constrain that swaps the components of a CompositionTree
//   - RetSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RetSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RetMerge for a re-constrain that merges identical components of a CompositionTree
//   - RetMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RetMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RetTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RetTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RetTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RetMergeL[I, S, T, A, B any, RET1 any, RET2 any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, CompositionTree[CompositionTree[RET1, RET1], RET2], RW, DIR, ERR]) Optic[I, S, T, A, B, CompositionTree[RET1, RET2], RW, DIR, ERR] {
	return UnsafeReconstrain[CompositionTree[RET1, RET2], RW, DIR, ERR](o)
}

// The RetMergeR re-constrain combines identical left and right side of the right side of a nested {CompositionTree]
//
// CompositionTree[RET1, CompositionTree[RET2, RET2]] --> CompositionTree[RET1, RET2]
//
// See:
//   - Ret1 for a re-constrain to  ReturnOne
//   - RetM for a re-constrain to ReturnMany
//   - RetL for a re-constrain to the left side of a CompositionTree
//   - RetR for a re-constrain to the Right side of a CompositionTree
//   - RetSwap for a re-constrain that swaps the components of a CompositionTree
//   - RetSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RetSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RetMerge for a re-constrain that merges identical components of a CompositionTree
//   - RetMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RetMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RetTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RetTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RetTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RetMergeR[I, S, T, A, B any, RET1 any, RET2 any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, CompositionTree[RET1, CompositionTree[RET2, RET2]], RW, DIR, ERR]) Optic[I, S, T, A, B, CompositionTree[RET1, RET2], RW, DIR, ERR] {
	return UnsafeReconstrain[CompositionTree[RET1, RET2], RW, DIR, ERR](o)
}

// The RetTrans re-constrain transposes the middle return type of a nested [CompositionTree].
//
// CompositionTree[CompositionTree[ERR1, ERR2], CompositionTree[ERR3, ERR4] --> CompositionTree[CompositionTree[ERR1, ERR3], CompositionTree[ERR2, ERR4]
//
// See:
//   - Ret1 for a re-constrain to  ReturnOne
//   - RetM for a re-constrain to ReturnMany
//   - RetL for a re-constrain to the left side of a CompositionTree
//   - RetR for a re-constrain to the Right side of a CompositionTree
//   - RetSwap for a re-constrain that swaps the components of a CompositionTree
//   - RetSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RetSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RetMerge for a re-constrain that merges identical components of a CompositionTree
//   - RetMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RetMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RetTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RetTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RetTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RetTrans[I, S, T, A, B any, RET1 any, RET2 any, RET3 any, RET4 any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, CompositionTree[CompositionTree[RET1, RET2], CompositionTree[RET3, RET4]], RW, DIR, ERR]) Optic[I, S, T, A, B, CompositionTree[CompositionTree[RET1, RET3], CompositionTree[RET2, RET4]], RW, DIR, ERR] {
	return UnsafeReconstrain[CompositionTree[CompositionTree[RET1, RET3], CompositionTree[RET2, RET4]], RW, DIR, ERR](o)
}

// The RetTransL re-constrain transposes the right element with the right child of a nested composition tree to the left.
//
// CompositionTree[CompositionTree[RET1, RET2], RET3] --> CompositionTree[CompositionTree[RET1, RET3], RET2]
// This has the effect of removing a level of nested from the RET2 path of the composition tree.
//
// See:
//   - Ret1 for a re-constrain to  ReturnOne
//   - RetM for a re-constrain to ReturnMany
//   - RetL for a re-constrain to the left side of a CompositionTree
//   - RetR for a re-constrain to the Right side of a CompositionTree
//   - RetSwap for a re-constrain that swaps the components of a CompositionTree
//   - RetSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RetSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RetMerge for a re-constrain that merges identical components of a CompositionTree
//   - RetMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RetMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RetTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RetTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RetTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RetTransL[I, S, T, A, B any, RET1 any, RET2 any, RET3 any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, CompositionTree[CompositionTree[RET1, RET2], RET3], RW, DIR, ERR]) Optic[I, S, T, A, B, CompositionTree[CompositionTree[RET1, RET3], RET2], RW, DIR, ERR] {
	return UnsafeReconstrain[CompositionTree[CompositionTree[RET1, RET3], RET2], RW, DIR, ERR](o)
}

// The RetTransR re-constrain transposes the left element with the left child of a nested composition tree to the right.
//
// CompositionTree[RET1, CompositionTree[RET2, RET3]] --> CompositionTree[RET2, CompositionTree[RET1, RET3]]
// This has the effect of removing a level of nested from the RET2 path of the composition tree.
//
// See:
//   - Ret1 for a re-constrain to  ReturnOne
//   - RetM for a re-constrain to ReturnMany
//   - RetL for a re-constrain to the left side of a CompositionTree
//   - RetR for a re-constrain to the Right side of a CompositionTree
//   - RetSwap for a re-constrain that swaps the components of a CompositionTree
//   - RetSwapL for a re-constrain that swaps the left components of a nested CompositionTree
//   - RetSwapR for a re-constrain that swaps the right components of a nested CompositionTree
//   - RetMerge for a re-constrain that merges identical components of a CompositionTree
//   - RetMergeL for a re-constrain that merges the identical left components of a nested CompositionTree
//   - RetMergeR for a re-constrain that merges the identical right components of a nested CompositionTree
//   - RetTrans for a re-constrain that swaps the middle elements of a nested CompositionTree.
//   - RetTransL for a re-constrain that swaps the right components into the left nested CompositionTree
//   - RetTransR for a re-constrain that swaps the left components into the right nested CompositionTree
func RetTransR[I, S, T, A, B any, RET1 any, RET2 any, RET3 any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, CompositionTree[RET1, CompositionTree[RET2, RET3]], RW, DIR, ERR]) Optic[I, S, T, A, B, CompositionTree[RET2, CompositionTree[RET1, RET3]], RW, DIR, ERR] {
	return UnsafeReconstrain[CompositionTree[RET2, CompositionTree[RET1, RET3]], RW, DIR, ERR](o)
}
