
//build !makecolops

package optic

//***************************************
//THIS FILE IS AUTO GENERATED DO NOT EDIT
//***************************************

import (

	"github.com/samber/lo"
	"github.com/samber/mo"
)

// The MapCol function converts a Map to a [Collection].
func MapCol[K comparable, S any](source map[K]S) Collection[K, S, Pure] {
	return MustGet(ColTypeToCol(MapColTypeP[K, S, S]()), source)
}

// TraverseMapCol returns a [Traversal] over the [Collection] representation of a Map.
func TraverseMapCol[K comparable, S any]() Optic[K, Collection[K, S, Pure], Collection[K, S, Pure], S, S, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseColIEP[K, S, S, Pure](TraverseColType(MapColTypeP[K, S, S]()).AsIxMatch())
}




// EqMapT2 returns a [Predicate] that is satisfied if the elements of the focused Map are all equal to (==) the elements of the provided constant Map and element [Predicate].
func EqMap[K comparable, S any, PERR TPure](right map[K]S, eq Predicate[lo.Tuple2[S, S], PERR]) Optic[Void, map[K]S, map[K]S, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	
	return OpT2Bind(EqMapT2[K, ](eq), right)
	
}

// EqMapT2 returns a [BinaryOp] that is satisfied if every element and index in Map A and Map B in focused order are equal.
//
// See:
//   - [EqMap] for a unary version.
func EqMapT2[K comparable, S any, PERR TPure](eq Predicate[lo.Tuple2[S, S], PERR]) Optic[Void, lo.Tuple2[map[K]S, map[K]S], lo.Tuple2[map[K]S, map[K]S], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return Ret1(Ro(EPure(Compose(
		DelveT2(ColTypeToCol(MapColTypeP[K, S, S]())),		
		EqColT2[K, S, Pure](
			eq,
		),
	))))
}


// The MapOp combinator applies the given collection operation to a Map and focuses the modified Map.
//
// See:
//   - [MapTraverseOp] for a combinator that efficiently traverses the elements of the modified Map.
//   - [MapColOp] for a combinator that focuses the modified [Collection].
func MapOp[K comparable, S, T any, I any, RET TReturnOne, RW any, DIR, ERR any](o Optic[I, Collection[K, S, ERR], Collection[K, T, ERR], Collection[K, T, ERR], Collection[K, T, ERR], RET, RW, DIR, ERR]) Optic[I, map[K]S, map[K]T, map[K]T, map[K]T, ReturnOne, RW, DIR, ERR] {
	return ColTypeOp(CombiColTypeErr[ERR](MapColTypeP[K, S, T]()), o)
}

// ColToMap returns an [Iso] that converts from a Map like collection to a Map
//
// See:
//  - ColToMapP for a polymorphic version.
func ColToMap[K comparable, S any]() Optic[Void, Collection[K, S, Pure], Collection[K, S, Pure], map[K]S, map[K]S, ReturnOne, ReadWrite, BiDir, Pure] {
	return ColToColType(MapColTypeP[K, S, S]())
}

// ColToMapP returns a polymorphic [Iso] that converts from a Map like collection to a Map
//
// See:
//  - ColToMap for a no polymorphic version.
func ColToMapP[K comparable, S, T any]() Optic[Void, Collection[K, S, Pure], Collection[K, T, Pure], map[K]S, map[K]T, ReturnOne, ReadWrite, BiDir, Pure] {
	return ColToColType(MapColTypeP[K, S, T]())
}

// DiffMap returns a [Traversal] that focuses on the differences between the given and focused Maps.
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffMapI] for an index aware version
//   - [DiffMapT2I] for a version that compares 2 focused Maps.
//   - [Diff] for the index detailing the detected diff.
//   - [Distance] for a convenience constructor for the distance operation.
func DiffMap[K comparable, S any, RET TReturnOne, ERR any](col map[K]S,threshold float64, distance Operation[lo.Tuple2[S, S], float64, RET, ERR], filterDiff DiffType) Optic[Diff[K, S], map[K]S, map[K]S, mo.Option[S], mo.Option[S], ReturnMany, ReadWrite, UniDir, ERR] {
	return DiffColType(CombiColTypeErr[ERR](MapColTypeP[K, S, S]()), col, threshold, distance, filterDiff)
}

// DiffMapT2 returns a [Traversal] that focuses on the differences between 2 Maps.
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Note: The source tuple should have the modified Map in the first position and the reference Map in the second position. This is compatible with [T2Dup] returning the modified tree.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffMapT2I] for an index aware version
//   - [DiffMap] for a a version that compares against a given Map
//   - [Diff] for the index detailing the detected diff.
//   - [Distance] for a convenience constructor for the distance operation.
func DiffMapT2[K comparable, S any, RET TReturnOne, ERR any](threshold float64, distance Operation[lo.Tuple2[S, S], float64, RET, ERR], filterDiff DiffType) Optic[Diff[K, S], lo.Tuple2[map[K]S, map[K]S], lo.Tuple2[map[K]S, map[K]S], mo.Option[S], mo.Option[S], ReturnMany, ReadWrite, UniDir, ERR] {
	return DiffColTypeT2(CombiColTypeErr[ERR](MapColTypeP[K, S, S]()), threshold, distance, filterDiff)
}


// DiffMapI returns a [Traversal] that focuses on the differences between a given and focused Map.
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffMap] for a non index aware version
//   - [DiffMapT2] for a version that compares 2 focused Maps.
//   - [Diff] for the index detailing the detected diff.
//   - [Distance] for a convenience constructor for the distance operation.
func DiffMapI[K comparable, S any, RET TReturnOne, ERR any](col map[K]S, threshold float64, distance Operation[lo.Tuple2[ValueI[K, S], ValueI[K, S]], float64, RET, ERR], filterDiff DiffType) Optic[Diff[K, S], map[K]S, map[K]S, mo.Option[S], mo.Option[S], ReturnMany, ReadWrite, UniDir, ERR] {
	return DiffColTypeI(CombiColTypeErr[ERR](MapColTypeP[K, S, S]()), col , threshold, distance, filterDiff)
}

// DiffMapT2I returns a [Traversal] that focuses on the differences between 2 Maps.
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Note: The source tuple should have the modified Map in the first position and the reference Map in the second position. This is compatible with [T2Dup] returning the modified tree.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffMapT2] for a non index aware version
//   - [DiffMapI] for a a version that compares against a given Map
//   - [Diff] for the index detailing the detected diff.
//   - [Distance] for a convenience constructor for the distance operation.
func DiffMapT2I[K comparable, S any, RET TReturnOne, ERR any](threshold float64, distance Operation[lo.Tuple2[ValueI[K, S], ValueI[K, S]], float64, RET, ERR], filterDiff DiffType) Optic[Diff[K, S], lo.Tuple2[map[K]S, map[K]S], lo.Tuple2[map[K]S, map[K]S], mo.Option[S], mo.Option[S], ReturnMany, ReadWrite, UniDir, ERR] {
	return DiffColTypeT2I(CombiColTypeErr[ERR](MapColTypeP[K, S, S]()), threshold, distance, filterDiff)
}

// AppendMapT2 applies the [AppendColT2] operation to a Map
func AppendMapT2[K comparable, S any, I, ERR any]() Optic[Void, lo.Tuple2[map[K]S, map[K]S], lo.Tuple2[map[K]S, map[K]S], map[K]S, map[K]S, ReturnOne, ReadOnly, UniDir, ERR] {
	return AppendColTypeT2(CombiColTypeErr[ERR](MapColTypeP[K, S, S]()))
}

// PrependMapT2 applies the [PrependColT2] operation to a Map
func PrependMapT2[K comparable, S any, I, ERR any]() Optic[Void, lo.Tuple2[map[K]S, map[K]S], lo.Tuple2[map[K]S, map[K]S], map[K]S, map[K]S, ReturnOne, ReadOnly, UniDir, ERR] {
	return PrependColTypeT2(CombiColTypeErr[ERR](MapColTypeP[K, S, S]()))
}



// FilteredMap applies the [Filtered] operation to a Map
//
// See:
// - [FilteredMapI] for an index aware version
func FilteredMap[K comparable, S any, ERR any](pred Predicate[S, ERR]) Optic[Void,map[K]S, map[K]S, map[K]S, map[K]S, ReturnOne, ReadWrite, UniDir, ERR] {
	return MapOp(FilteredCol[K](pred))
}




// FilteredMapI applies the [FilteredI] operation to a Map
//
// See:
// - [FilteredMap] for a non index aware version
func FilteredMapI[K comparable, S any, ERR any](pred PredicateI[K, S, ERR]) Optic[Void,map[K]S, map[K]S, map[K]S, map[K]S, ReturnOne, ReadWrite, UniDir, ERR] {
	return MapOp(FilteredColI(pred,IxMatchComparable[K]()))
}


			



// AppendMap applies the [Append] operation to a Map
func AppendMap[K comparable, S any, ERR any](toAppend Collection[K, S, ERR]) Optic[Void,map[K]S, map[K]S, map[K]S, map[K]S, ReturnOne, ReadOnly, UniDir, ERR] {
	return MapOp(AppendCol(toAppend))
}





