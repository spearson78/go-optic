
//build !makecolops

package optic

//***************************************
//THIS FILE IS AUTO GENERATED DO NOT EDIT
//***************************************

import (

	"github.com/samber/lo"
	"github.com/samber/mo"
)

// The StringCol function converts a String to a [Collection].
func StringCol(source string) Collection[int, rune, Pure] {
	return MustGet(ColTypeToCol(StringColType()), source)
}

// TraverseStringCol returns a [Traversal] over the [Collection] representation of a String.
func TraverseStringCol() Optic[int, Collection[int, rune, Pure], Collection[int, rune, Pure], rune, rune, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseColIEP[int, rune, rune, Pure](TraverseColType(StringColType()).AsIxMatch())
}




// EqString returns a [Predicate] that is satisfied if the elements of the focused String are all equal to (==) the elements of the provided constant String and element [Predicate].
func EqString(right string) Optic[Void, string, string, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	
	return OpT2Bind(EqStringT2(),right)
	
}

// EqStringT2 returns an [BinaryOp] that is satisfied if every element and index in String A and String B in focused order are equal.
//
// See:
//   - [EqString] for a unary version.
func EqStringT2() Optic[Void, lo.Tuple2[string, string], lo.Tuple2[string, string], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return Ret1(Ro(EPure(Compose(
		DelveT2(ColTypeToCol(StringColType())),		
		EqColT2[int, rune, Pure](
			EqT2[rune](),
		),
	))))
}


// The StringOp combinator applies the given collection operation to a String and focuses the modified String.
//
// See:
//   - [StringTraverseOp] for a combinator that efficiently traverses the elements of the modified String.
//   - [StringColOp] for a combinator that focuses the modified [Collection].
func StringOp[I any, RET TReturnOne, RW any, DIR, ERR any](o Optic[I, Collection[int, rune, ERR], Collection[int, rune, ERR], Collection[int, rune, ERR], Collection[int, rune, ERR], RET, RW, DIR, ERR]) Optic[I, string, string, string, string, ReturnOne, RW, DIR, ERR] {
	return ColTypeOp(CombiColTypeErr[ERR](StringColType()), o)
}

// ColToString returns an [Iso] that converts from a String like collection to a String
//
// See:
//  - ColToStringP for a polymorphic version.
func ColToString() Optic[Void, Collection[int, rune, Pure], Collection[int, rune, Pure], string, string, ReturnOne, ReadWrite, BiDir, Pure] {
	return ColToColType(StringColType())
}

// ColToStringP returns a polymorphic [Iso] that converts from a String like collection to a String
//
// See:
//  - ColToString for a no polymorphic version.
func ColToStringP() Optic[Void, Collection[int, rune, Pure], Collection[int, rune, Pure], string, string, ReturnOne, ReadWrite, BiDir, Pure] {
	return ColToColType(StringColType())
}

// DiffString returns a [Traversal] that focuses on the differences between the given and focused Strings.
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffStringI] for an index aware version
//   - [DiffStringT2I] for a version that compares 2 focused Strings.
//   - [Diff] for the index detailing the detected diff.
//   - [Distance] for a convenience constructor for the distance operation.
func DiffString[RET TReturnOne, ERR any](col string,threshold float64, distance Operation[lo.Tuple2[rune, rune], float64, RET, ERR], filterDiff DiffType) Optic[Diff[int, rune], string, string, mo.Option[rune], mo.Option[rune], ReturnMany, ReadWrite, UniDir, ERR] {
	return DiffColType(CombiColTypeErr[ERR](StringColType()), col, threshold, distance, filterDiff)
}

// DiffStringT2 returns a [Traversal] that focuses on the differences between 2 Strings.
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Note: The source tuple should have the modified String in the first position and the reference String in the second position. This is compatible with [T2Dup] returning the modified tree.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffStringT2I] for an index aware version
//   - [DiffString] for a a version that compares against a given String
//   - [Diff] for the index detailing the detected diff.
//   - [Distance] for a convenience constructor for the distance operation.
func DiffStringT2[RET TReturnOne, ERR any](threshold float64, distance Operation[lo.Tuple2[rune, rune], float64, RET, ERR], filterDiff DiffType) Optic[Diff[int, rune], lo.Tuple2[string, string], lo.Tuple2[string, string], mo.Option[rune], mo.Option[rune], ReturnMany, ReadWrite, UniDir, ERR] {
	return DiffColTypeT2(CombiColTypeErr[ERR](StringColType()), threshold, distance, filterDiff)
}


// DiffStringI returns a [Traversal] that focuses on the differences between a given and focused String.
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffString] for a non index aware version
//   - [DiffStringT2] for a version that compares 2 focused Strings.
//   - [Diff] for the index detailing the detected diff.
//   - [Distance] for a convenience constructor for the distance operation.
func DiffStringI[RET TReturnOne, ERR any](col string, threshold float64, distance Operation[lo.Tuple2[ValueI[int, rune], ValueI[int, rune]], float64, RET, ERR], filterDiff DiffType) Optic[Diff[int, rune], string, string, mo.Option[rune], mo.Option[rune], ReturnMany, ReadWrite, UniDir, ERR] {
	return DiffColTypeI(CombiColTypeErr[ERR](StringColType()), col , threshold, distance, filterDiff)
}

// DiffStringT2I returns a [Traversal] that focuses on the differences between 2 Strings.
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Note: The source tuple should have the modified String in the first position and the reference String in the second position. This is compatible with [T2Dup] returning the modified tree.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffStringT2] for a non index aware version
//   - [DiffStringI] for a a version that compares against a given String
//   - [Diff] for the index detailing the detected diff.
//   - [Distance] for a convenience constructor for the distance operation.
func DiffStringT2I[RET TReturnOne, ERR any](threshold float64, distance Operation[lo.Tuple2[ValueI[int, rune], ValueI[int, rune]], float64, RET, ERR], filterDiff DiffType) Optic[Diff[int, rune], lo.Tuple2[string, string], lo.Tuple2[string, string], mo.Option[rune], mo.Option[rune], ReturnMany, ReadWrite, UniDir, ERR] {
	return DiffColTypeT2I(CombiColTypeErr[ERR](StringColType()), threshold, distance, filterDiff)
}

// AppendStringT2 applies the [AppendColT2] operation to a String
func AppendStringT2[I, ERR any]() Optic[Void, lo.Tuple2[string, string], lo.Tuple2[string, string], string, string, ReturnOne, ReadOnly, UniDir, ERR] {
	return AppendColTypeT2(CombiColTypeErr[ERR](StringColType()))
}

// PrependStringT2 applies the [PrependColT2] operation to a String
func PrependStringT2[I, ERR any]() Optic[Void, lo.Tuple2[string, string], lo.Tuple2[string, string], string, string, ReturnOne, ReadOnly, UniDir, ERR] {
	return PrependColTypeT2(CombiColTypeErr[ERR](StringColType()))
}



// FilteredString applies the [Filtered] operation to a String
//
// See:
// - [FilteredStringI] for an index aware version
func FilteredString[ERR any](pred Predicate[rune, ERR]) Optic[Void,string, string, string, string, ReturnOne, ReadWrite, UniDir, ERR] {
	return StringOp(FilteredCol[int](pred))
}




// FilteredStringI applies the [FilteredI] operation to a String
//
// See:
// - [FilteredString] for a non index aware version
func FilteredStringI[ERR any](pred PredicateI[int, rune, ERR]) Optic[Void,string, string, string, string, ReturnOne, ReadWrite, UniDir, ERR] {
	return StringOp(FilteredColI(pred,IxMatchComparable[int]()))
}


			



// AppendString applies the [Append] operation to a String
func AppendString[ERR any](toAppend Collection[int, rune, ERR]) Optic[Void,string, string, string, string, ReturnOne, ReadOnly, UniDir, ERR] {
	return StringOp(AppendCol(toAppend))
}






// PrependString applies the [Prepend] operation to a String
func PrependString[ERR any](toPrepend Collection[int, rune, ERR]) Optic[Void,string, string, string, string, ReturnOne, ReadOnly, UniDir, ERR] {
	return StringOp(PrependCol(toPrepend))
}






// SubString applies the [SubCol] operation to a String
func SubString(start int,length int) Optic[Void,string, string, string, string, ReturnOne, ReadWrite, UniDir, Pure] {
	return StringOp(SubCol[int, rune](start,length))
}






// ReversedString applies the [Reversed] operation to a String
func ReversedString() Optic[Void,string, string, string, string, ReturnOne, ReadWrite, BiDir, Pure] {
	return StringOp(ReversedCol[int, rune]())
}






// OrderedString applies the [Ordered] operation to a String
//
// See:
// - [OrderedStringI] for an index aware version
func OrderedString[ERR any](orderBy OrderByPredicate[rune, ERR]) Optic[Void,string, string, string, string, ReturnOne, ReadWrite, UniDir, ERR] {
	return StringOp(OrderedCol[int](orderBy))
}




// OrderedStringI applies the [OrderedI] operation to a String
//
// See:
// - [OrderedString] for a non index aware version
func OrderedStringI[ERR any](orderBy OrderByPredicateI[int, rune, ERR]) Optic[Void,string, string, string, string, ReturnOne, ReadWrite, UniDir, ERR] {
	return StringOp(OrderedColI(orderBy,IxMatchComparable[int]()))
}


			


