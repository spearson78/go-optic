
//build !makecolops

package optic

//***************************************
//THIS FILE IS AUTO GENERATED DO NOT EDIT
//***************************************

import (

	"github.com/samber/lo"
	"github.com/samber/mo"
)

// The SliceCol function converts a Slice to a [Collection].
func SliceCol[A any](source []A) Collection[int, A, Pure] {
	return MustGet(ColTypeToCol(SliceColTypeP[A, A]()), source)
}

// TraverseSliceCol returns a [Traversal] over the [Collection] representation of a Slice.
func TraverseSliceCol[A any]() Optic[int, Collection[int, A, Pure], Collection[int, A, Pure], A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseColIEP[int, A, A, Pure](TraverseColType(SliceColTypeP[A, A]()).AsIxMatch())
}




// EqSliceT2 returns a [Predicate] that is satisfied if the elements of the focused Slice are all equal to (==) the elements of the provided constant Slice and element [Predicate].
func EqSlice[A any, PERR TPure](right []A, eq Predicate[lo.Tuple2[A, A], PERR]) Optic[Void, []A, []A, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	
	return OpT2Bind(EqSliceT2[A](eq), right)
	
}

// EqSliceT2 returns a [BinaryOp] that is satisfied if every element and index in Slice A and Slice B in focused order are equal.
//
// See:
//   - [EqSlice] for a unary version.
func EqSliceT2[A any, PERR TPure](eq Predicate[lo.Tuple2[A, A], PERR]) Optic[Void, lo.Tuple2[[]A, []A], lo.Tuple2[[]A, []A], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return Ret1(Ro(EPure(Compose(
		DelveT2(ColTypeToCol(SliceColTypeP[A, A]())),		
		EqColT2[int, A, Pure](
			eq,
		),
	))))
}


// The SliceOp combinator applies the given collection operation to a Slice and focuses the modified Slice.
//
// See:
//   - [SliceTraverseOp] for a combinator that efficiently traverses the elements of the modified Slice.
//   - [SliceColOp] for a combinator that focuses the modified [Collection].
func SliceOp[A, B any, I any, RET TReturnOne, RW any, DIR, ERR any](o Optic[I, Collection[int, A, ERR], Collection[int, B, ERR], Collection[int, B, ERR], Collection[int, B, ERR], RET, RW, DIR, ERR]) Optic[I, []A, []B, []B, []B, ReturnOne, RW, DIR, ERR] {
	return ColTypeOp(CombiColTypeErr[ERR](SliceColTypeP[A, B]()), o)
}

// ColToSlice returns an [Iso] that converts from a Slice like collection to a Slice
//
// See:
//  - ColToSliceP for a polymorphic version.
func ColToSlice[A any]() Optic[Void, Collection[int, A, Pure], Collection[int, A, Pure], []A, []A, ReturnOne, ReadWrite, BiDir, Pure] {
	return ColToColType(SliceColTypeP[A, A]())
}

// ColToSliceP returns a polymorphic [Iso] that converts from a Slice like collection to a Slice
//
// See:
//  - ColToSlice for a no polymorphic version.
func ColToSliceP[A, B any]() Optic[Void, Collection[int, A, Pure], Collection[int, B, Pure], []A, []B, ReturnOne, ReadWrite, BiDir, Pure] {
	return ColToColType(SliceColTypeP[A, B]())
}

// DiffSlice returns a [Traversal] that focuses on the differences between the given and focused Slices.
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffSliceI] for an index aware version
//   - [DiffSliceT2I] for a version that compares 2 focused Slices.
//   - [Diff] for the index detailing the detected diff.
//   - [Distance] for a convenience constructor for the distance operation.
func DiffSlice[A any, RET TReturnOne, ERR any](col []A,threshold float64, distance Operation[lo.Tuple2[A, A], float64, RET, ERR], filterDiff DiffType) Optic[Diff[int, A], []A, []A, mo.Option[A], mo.Option[A], ReturnMany, ReadWrite, UniDir, ERR] {
	return DiffColType(CombiColTypeErr[ERR](SliceColTypeP[A, A]()), col, threshold, distance, filterDiff)
}

// DiffSliceT2 returns a [Traversal] that focuses on the differences between 2 Slices.
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Note: The source tuple should have the modified Slice in the first position and the reference Slice in the second position. This is compatible with [T2Dup] returning the modified tree.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffSliceT2I] for an index aware version
//   - [DiffSlice] for a a version that compares against a given Slice
//   - [Diff] for the index detailing the detected diff.
//   - [Distance] for a convenience constructor for the distance operation.
func DiffSliceT2[A any, RET TReturnOne, ERR any](threshold float64, distance Operation[lo.Tuple2[A, A], float64, RET, ERR], filterDiff DiffType) Optic[Diff[int, A], lo.Tuple2[[]A, []A], lo.Tuple2[[]A, []A], mo.Option[A], mo.Option[A], ReturnMany, ReadWrite, UniDir, ERR] {
	return DiffColTypeT2(CombiColTypeErr[ERR](SliceColTypeP[A, A]()), threshold, distance, filterDiff)
}


// DiffSliceI returns a [Traversal] that focuses on the differences between a given and focused Slice.
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffSlice] for a non index aware version
//   - [DiffSliceT2] for a version that compares 2 focused Slices.
//   - [Diff] for the index detailing the detected diff.
//   - [Distance] for a convenience constructor for the distance operation.
func DiffSliceI[A any, RET TReturnOne, ERR any](col []A, threshold float64, distance Operation[lo.Tuple2[ValueI[int, A], ValueI[int, A]], float64, RET, ERR], filterDiff DiffType) Optic[Diff[int, A], []A, []A, mo.Option[A], mo.Option[A], ReturnMany, ReadWrite, UniDir, ERR] {
	return DiffColTypeI(CombiColTypeErr[ERR](SliceColTypeP[A, A]()), col , threshold, distance, filterDiff)
}

// DiffSliceT2I returns a [Traversal] that focuses on the differences between 2 Slices.
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Note: The source tuple should have the modified Slice in the first position and the reference Slice in the second position. This is compatible with [T2Dup] returning the modified tree.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffSliceT2] for a non index aware version
//   - [DiffSliceI] for a a version that compares against a given Slice
//   - [Diff] for the index detailing the detected diff.
//   - [Distance] for a convenience constructor for the distance operation.
func DiffSliceT2I[A any, RET TReturnOne, ERR any](threshold float64, distance Operation[lo.Tuple2[ValueI[int, A], ValueI[int, A]], float64, RET, ERR], filterDiff DiffType) Optic[Diff[int, A], lo.Tuple2[[]A, []A], lo.Tuple2[[]A, []A], mo.Option[A], mo.Option[A], ReturnMany, ReadWrite, UniDir, ERR] {
	return DiffColTypeT2I(CombiColTypeErr[ERR](SliceColTypeP[A, A]()), threshold, distance, filterDiff)
}

// AppendSliceT2 applies the [AppendColT2] operation to a Slice
func AppendSliceT2[A any, I, ERR any]() Optic[Void, lo.Tuple2[[]A, []A], lo.Tuple2[[]A, []A], []A, []A, ReturnOne, ReadOnly, UniDir, ERR] {
	return AppendColTypeT2(CombiColTypeErr[ERR](SliceColTypeP[A, A]()))
}

// PrependSliceT2 applies the [PrependColT2] operation to a Slice
func PrependSliceT2[A any, I, ERR any]() Optic[Void, lo.Tuple2[[]A, []A], lo.Tuple2[[]A, []A], []A, []A, ReturnOne, ReadOnly, UniDir, ERR] {
	return PrependColTypeT2(CombiColTypeErr[ERR](SliceColTypeP[A, A]()))
}



// FilteredSlice applies the [Filtered] operation to a Slice
//
// See:
// - [FilteredSliceI] for an index aware version
func FilteredSlice[A any, ERR any](pred Predicate[A, ERR]) Optic[Void,[]A, []A, []A, []A, ReturnOne, ReadWrite, UniDir, ERR] {
	return SliceOp(FilteredCol[int](pred))
}




// FilteredSliceI applies the [FilteredI] operation to a Slice
//
// See:
// - [FilteredSlice] for a non index aware version
func FilteredSliceI[A any, ERR any](pred PredicateI[int, A, ERR]) Optic[Void,[]A, []A, []A, []A, ReturnOne, ReadWrite, UniDir, ERR] {
	return SliceOp(FilteredColI(pred,IxMatchComparable[int]()))
}


			



// AppendSlice applies the [Append] operation to a Slice
func AppendSlice[A any, ERR any](toAppend Collection[int, A, ERR]) Optic[Void,[]A, []A, []A, []A, ReturnOne, ReadOnly, UniDir, ERR] {
	return SliceOp(AppendCol(toAppend))
}






// PrependSlice applies the [Prepend] operation to a Slice
func PrependSlice[A any, ERR any](toPrepend Collection[int, A, ERR]) Optic[Void,[]A, []A, []A, []A, ReturnOne, ReadOnly, UniDir, ERR] {
	return SliceOp(PrependCol(toPrepend))
}






// SubSlice applies the [SubCol] operation to a Slice
func SubSlice[A any](start int,length int) Optic[Void,[]A, []A, []A, []A, ReturnOne, ReadWrite, UniDir, Pure] {
	return SliceOp(SubCol[int, A](start,length))
}






// ReversedSlice applies the [Reversed] operation to a Slice
func ReversedSlice[A any]() Optic[Void,[]A, []A, []A, []A, ReturnOne, ReadWrite, BiDir, Pure] {
	return SliceOp(ReversedCol[int, A]())
}






// OrderedSlice applies the [Ordered] operation to a Slice
//
// See:
// - [OrderedSliceI] for an index aware version
func OrderedSlice[A any, ERR any](orderBy OrderByPredicate[A, ERR]) Optic[Void,[]A, []A, []A, []A, ReturnOne, ReadWrite, UniDir, ERR] {
	return SliceOp(OrderedCol[int](orderBy))
}




// OrderedSliceI applies the [OrderedI] operation to a Slice
//
// See:
// - [OrderedSlice] for a non index aware version
func OrderedSliceI[A any, ERR any](orderBy OrderByPredicateI[int, A, ERR]) Optic[Void,[]A, []A, []A, []A, ReturnOne, ReadWrite, UniDir, ERR] {
	return SliceOp(OrderedColI(orderBy,IxMatchComparable[int]()))
}


			


