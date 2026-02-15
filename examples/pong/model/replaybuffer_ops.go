
//build !makecolops

package model

//***************************************
//THIS FILE IS AUTO GENERATED DO NOT EDIT
//***************************************

import (
	. "github.com/spearson78/go-optic"


	"github.com/samber/lo"
	"github.com/samber/mo"
)

// The ReplayBufferCol function converts a ReplayBuffer to a [Collection].
func ReplayBufferCol[A any](len int, source ReplayBuffer[A]) Collection[int, A, Pure] {
	return MustGet(ColTypeToCol(ReplayBufferColTypeP[A,A](len)), source)
}

// TraverseReplayBufferCol returns a [Traversal] over the [Collection] representation of a ReplayBuffer.
func TraverseReplayBufferCol[A any](len int, ) Optic[int, Collection[int, A, Pure], Collection[int, A, Pure], A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseColIEP[int, A, A, Pure](TraverseColType(ReplayBufferColTypeP[A,A](len)).AsIxMatch())
}


// The ReplayBufferOf combinator focuses on a ReplayBuffer of all the elements in the given optic.
//
// Under modification this collection can be modified and will be rebuilt into the original data structure.
// If the modified ReplayBuffer contains fewer elements the result will use values from the original source.
// If the modified ReplayBuffer contains more elements they will be ignored.
//
// See:
//   - ReplayBufferOfP for a polymorphic version
func ReplayBufferOf[A any, OS, OT, RET, RW, DIR, ERR any](len int, o Optic[int, OS, OT, A, A, RET, RW, DIR, ERR]) Optic[Void, OS, OT, ReplayBuffer[A], ReplayBuffer[A], ReturnOne, RW, UniDir, ERR] {
	return ColTypeOf(CombiColTypeErr[ERR](ReplayBufferColTypeP[A,A](len)), o)
}

// The ReplayBufferOfP polymorphic combinator focuses on a ReplayBuffer of all the elements in the given optic.
//
// Under modification this collection can be modified and will be rebuilt into the original data structure.
// If the modified ReplayBuffer contains fewer elements then [ErrUnsafeMissingElement] will be returned
// If the modified ReplayBuffer contains more elements they will be ignored.
//
// See:
//   - ReplayBufferOf for a non polymorphic version
func ReplayBufferOfP[A,B any, OS, OT, RET, RW, DIR, ERR any](len int, o Optic[int, OS, OT, A, B, RET, RW, DIR, ERR]) Optic[Void, OS, OT, ReplayBuffer[A], ReplayBuffer[B], ReturnOne, RW, UniDir, Err] {
	return ColTypeOfP(CombiColTypeErr[ERR](ReplayBufferColTypeP[A,B](len)), o)
}



// EqReplayBufferT2 returns a [Predicate] that is satisfied if the elements of the focused ReplayBuffer are all equal to (==) the elements of the provided constant ReplayBuffer and element [Predicate].
func EqReplayBuffer[A any, PERR TPure](len int, right ReplayBuffer[A], eq Predicate[lo.Tuple2[A, A], PERR]) Optic[Void, ReplayBuffer[A], ReplayBuffer[A], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	
	return OpT2Bind(EqReplayBufferT2[A](len, eq), right)
	
}

// EqReplayBufferT2 returns a [BinaryOp] that is satisfied if every element and index in ReplayBuffer A and ReplayBuffer B in focused order are equal.
//
// See:
//   - [EqReplayBuffer] for a unary version.
func EqReplayBufferT2[A any, PERR TPure](len int, eq Predicate[lo.Tuple2[A, A], PERR]) Optic[Void, lo.Tuple2[ReplayBuffer[A], ReplayBuffer[A]], lo.Tuple2[ReplayBuffer[A], ReplayBuffer[A]], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return Ret1(Ro(EPure(Compose(
		DelveT2(ColTypeToCol(ReplayBufferColTypeP[A,A](len))),		
		EqColT2[int, A, Pure](
			eq,
		),
	))))
}


// The ReplayBufferOp combinator applies the given collection operation to a ReplayBuffer and focuses the modified ReplayBuffer.
//
// See:
//   - [ReplayBufferTraverseOp] for a combinator that efficiently traverses the elements of the modified ReplayBuffer.
//   - [ReplayBufferColOp] for a combinator that focuses the modified [Collection].
func ReplayBufferOp[A,B any, I any, RET TReturnOne, RW any, DIR, ERR any](len int, o Optic[I, Collection[int, A, ERR], Collection[int, B, ERR], Collection[int, B, ERR], Collection[int, B, ERR], RET, RW, DIR, ERR]) Optic[I, ReplayBuffer[A], ReplayBuffer[B], ReplayBuffer[B], ReplayBuffer[B], ReturnOne, RW, DIR, ERR] {
	return ColTypeOp(CombiColTypeErr[ERR](ReplayBufferColTypeP[A,B](len)), o)
}

// ColToReplayBuffer returns an [Iso] that converts from a ReplayBuffer like collection to a ReplayBuffer
//
// See:
//  - ColToReplayBufferP for a polymorphic version.
func ColToReplayBuffer[A any](len int, ) Optic[Void, Collection[int, A, Pure], Collection[int, A, Pure], ReplayBuffer[A], ReplayBuffer[A], ReturnOne, ReadWrite, BiDir, Pure] {
	return ColToColType(ReplayBufferColTypeP[A,A](len))
}

// ColToReplayBufferP returns a polymorphic [Iso] that converts from a ReplayBuffer like collection to a ReplayBuffer
//
// See:
//  - ColToReplayBuffer for a no polymorphic version.
func ColToReplayBufferP[A,B any](len int, ) Optic[Void, Collection[int, A, Pure], Collection[int, B, Pure], ReplayBuffer[A], ReplayBuffer[B], ReturnOne, ReadWrite, BiDir, Pure] {
	return ColToColType(ReplayBufferColTypeP[A,B](len))
}

// DiffReplayBuffer returns a [Traversal] that focuses on the differences between the given and focused ReplayBuffers.
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffReplayBufferI] for an index aware version
//   - [DiffReplayBufferT2I] for a version that compares 2 focused ReplayBuffers.
//   - [Diff] for the index detailing the detected diff.
//   - [Distance] for a convenience constructor for the distance operation.
func DiffReplayBuffer[A any, RET TReturnOne, ERR any](len int, col ReplayBuffer[A],threshold float64, distance Operation[lo.Tuple2[A, A], float64, RET, ERR], filterDiff DiffType) Optic[Diff[int, A], ReplayBuffer[A], ReplayBuffer[A], mo.Option[A], mo.Option[A], ReturnMany, ReadWrite, UniDir, ERR] {
	return DiffColType(CombiColTypeErr[ERR](ReplayBufferColTypeP[A,A](len)), col, threshold, distance, filterDiff)
}

// DiffReplayBufferT2 returns a [Traversal] that focuses on the differences between 2 ReplayBuffers.
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Note: The source tuple should have the modified ReplayBuffer in the first position and the reference ReplayBuffer in the second position. This is compatible with [T2Dup] returning the modified tree.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffReplayBufferT2I] for an index aware version
//   - [DiffReplayBuffer] for a a version that compares against a given ReplayBuffer
//   - [Diff] for the index detailing the detected diff.
//   - [Distance] for a convenience constructor for the distance operation.
func DiffReplayBufferT2[A any, RET TReturnOne, ERR any](len int, threshold float64, distance Operation[lo.Tuple2[A, A], float64, RET, ERR], filterDiff DiffType) Optic[Diff[int, A], lo.Tuple2[ReplayBuffer[A], ReplayBuffer[A]], lo.Tuple2[ReplayBuffer[A], ReplayBuffer[A]], mo.Option[A], mo.Option[A], ReturnMany, ReadWrite, UniDir, ERR] {
	return DiffColTypeT2(CombiColTypeErr[ERR](ReplayBufferColTypeP[A,A](len)), threshold, distance, filterDiff)
}


// DiffReplayBufferI returns a [Traversal] that focuses on the differences between a given and focused ReplayBuffer.
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffReplayBuffer] for a non index aware version
//   - [DiffReplayBufferT2] for a version that compares 2 focused ReplayBuffers.
//   - [Diff] for the index detailing the detected diff.
//   - [Distance] for a convenience constructor for the distance operation.
func DiffReplayBufferI[A any, RET TReturnOne, ERR any](len int, col ReplayBuffer[A], threshold float64, distance Operation[lo.Tuple2[ValueI[int, A], ValueI[int, A]], float64, RET, ERR], filterDiff DiffType) Optic[Diff[int, A], ReplayBuffer[A], ReplayBuffer[A], mo.Option[A], mo.Option[A], ReturnMany, ReadWrite, UniDir, ERR] {
	return DiffColTypeI(CombiColTypeErr[ERR](ReplayBufferColTypeP[A,A](len)), col , threshold, distance, filterDiff)
}

// DiffReplayBufferT2I returns a [Traversal] that focuses on the differences between 2 ReplayBuffers.
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Note: The source tuple should have the modified ReplayBuffer in the first position and the reference ReplayBuffer in the second position. This is compatible with [T2Dup] returning the modified tree.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffReplayBufferT2] for a non index aware version
//   - [DiffReplayBufferI] for a a version that compares against a given ReplayBuffer
//   - [Diff] for the index detailing the detected diff.
//   - [Distance] for a convenience constructor for the distance operation.
func DiffReplayBufferT2I[A any, RET TReturnOne, ERR any](len int, threshold float64, distance Operation[lo.Tuple2[ValueI[int, A], ValueI[int, A]], float64, RET, ERR], filterDiff DiffType) Optic[Diff[int, A], lo.Tuple2[ReplayBuffer[A], ReplayBuffer[A]], lo.Tuple2[ReplayBuffer[A], ReplayBuffer[A]], mo.Option[A], mo.Option[A], ReturnMany, ReadWrite, UniDir, ERR] {
	return DiffColTypeT2I(CombiColTypeErr[ERR](ReplayBufferColTypeP[A,A](len)), threshold, distance, filterDiff)
}

// AppendReplayBufferT2 applies the [AppendColT2] operation to a ReplayBuffer
func AppendReplayBufferT2[A any, I, ERR any](len int, ) Optic[Void, lo.Tuple2[ReplayBuffer[A], ReplayBuffer[A]], lo.Tuple2[ReplayBuffer[A], ReplayBuffer[A]], ReplayBuffer[A], ReplayBuffer[A], ReturnOne, ReadOnly, UniDir, ERR] {
	return AppendColTypeT2(CombiColTypeErr[ERR](ReplayBufferColTypeP[A,A](len)))
}

// PrependReplayBufferT2 applies the [PrependColT2] operation to a ReplayBuffer
func PrependReplayBufferT2[A any, I, ERR any](len int, ) Optic[Void, lo.Tuple2[ReplayBuffer[A], ReplayBuffer[A]], lo.Tuple2[ReplayBuffer[A], ReplayBuffer[A]], ReplayBuffer[A], ReplayBuffer[A], ReturnOne, ReadOnly, UniDir, ERR] {
	return PrependColTypeT2(CombiColTypeErr[ERR](ReplayBufferColTypeP[A,A](len)))
}



// FilteredReplayBuffer applies the [Filtered] operation to a ReplayBuffer
//
// See:
// - [FilteredReplayBufferI] for an index aware version
func FilteredReplayBuffer[A any, ERR any](len int, pred Predicate[A, ERR]) Optic[Void,ReplayBuffer[A], ReplayBuffer[A], ReplayBuffer[A], ReplayBuffer[A], ReturnOne, ReadWrite, UniDir, ERR] {
	return ReplayBufferOp(len, FilteredCol[int](pred))
}




// FilteredReplayBufferI applies the [FilteredI] operation to a ReplayBuffer
//
// See:
// - [FilteredReplayBuffer] for a non index aware version
func FilteredReplayBufferI[A any, ERR any](len int, pred PredicateI[int, A, ERR]) Optic[Void,ReplayBuffer[A], ReplayBuffer[A], ReplayBuffer[A], ReplayBuffer[A], ReturnOne, ReadWrite, UniDir, ERR] {
	return ReplayBufferOp(len, FilteredColI(pred,IxMatchComparable[int]()))
}


			



// AppendReplayBuffer applies the [Append] operation to a ReplayBuffer
func AppendReplayBuffer[A any, ERR any](len int, toAppend Collection[int, A, ERR]) Optic[Void,ReplayBuffer[A], ReplayBuffer[A], ReplayBuffer[A], ReplayBuffer[A], ReturnOne, ReadOnly, UniDir, ERR] {
	return ReplayBufferOp(len, AppendCol(toAppend))
}






// PrependReplayBuffer applies the [Prepend] operation to a ReplayBuffer
func PrependReplayBuffer[A any, ERR any](len int, toPrepend Collection[int, A, ERR]) Optic[Void,ReplayBuffer[A], ReplayBuffer[A], ReplayBuffer[A], ReplayBuffer[A], ReturnOne, ReadOnly, UniDir, ERR] {
	return ReplayBufferOp(len, PrependCol(toPrepend))
}






// SubReplayBuffer applies the [SubCol] operation to a ReplayBuffer
func SubReplayBuffer[A any](len int, start int,length int) Optic[Void,ReplayBuffer[A], ReplayBuffer[A], ReplayBuffer[A], ReplayBuffer[A], ReturnOne, ReadWrite, UniDir, Pure] {
	return ReplayBufferOp(len, SubCol[int, A](start,length))
}






// ReversedReplayBuffer applies the [Reversed] operation to a ReplayBuffer
func ReversedReplayBuffer[A any](len int) Optic[Void,ReplayBuffer[A], ReplayBuffer[A], ReplayBuffer[A], ReplayBuffer[A], ReturnOne, ReadWrite, BiDir, Pure] {
	return ReplayBufferOp(len, ReversedCol[int, A]())
}






// OrderedReplayBuffer applies the [Ordered] operation to a ReplayBuffer
//
// See:
// - [OrderedReplayBufferI] for an index aware version
func OrderedReplayBuffer[A any, ERR any](len int, orderBy OrderByPredicate[A, ERR]) Optic[Void,ReplayBuffer[A], ReplayBuffer[A], ReplayBuffer[A], ReplayBuffer[A], ReturnOne, ReadWrite, UniDir, ERR] {
	return ReplayBufferOp(len, OrderedCol[int](orderBy))
}




// OrderedReplayBufferI applies the [OrderedI] operation to a ReplayBuffer
//
// See:
// - [OrderedReplayBuffer] for a non index aware version
func OrderedReplayBufferI[A any, ERR any](len int, orderBy OrderByPredicateI[int, A, ERR]) Optic[Void,ReplayBuffer[A], ReplayBuffer[A], ReplayBuffer[A], ReplayBuffer[A], ReturnOne, ReadWrite, UniDir, ERR] {
	return ReplayBufferOp(len, OrderedColI(orderBy,IxMatchComparable[int]()))
}


			


