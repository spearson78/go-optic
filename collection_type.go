package optic

import (
	"reflect"

	"github.com/spearson78/go-optic/expr"
)

// CollectionType contains [Iso]s to convert between [Collection] and an arbitrary collection type like slice or map.
// This enables collection operations like [FilteredCol],[Append] and [OrderedCol] to be applied to any CollectionType type.
type CollectionType[I, S, T, A, B any, ERR any] struct {
	toCol    Optic[Void, S, T, Collection[I, A, ERR], Collection[I, B, ERR], ReturnOne, ReadWrite, BiDir, ERR]
	fromColP Optic[Void, Collection[I, A, ERR], Collection[I, B, ERR], S, T, ReturnOne, ReadWrite, BiDir, ERR]
	fromColB Optic[Void, Collection[I, B, ERR], Collection[I, B, ERR], T, T, ReturnOne, ReadWrite, BiDir, ERR]
	traverse Optic[I, S, T, A, B, ReturnMany, ReadWrite, UniDir, ERR]
}

// ColType returns a [CollectionType] using the given [Iso] to convert to and from a [Collection]
//
// See:
//   - [ColTypeP] for a polymorphic version
func ColType[I, S, A any, TRET any, RET TReturnOne, TRW, RW TReadWrite, TDIR any, DIR TBiDir, ERR any](
	toCol Optic[Void, S, S, Collection[I, A, ERR], Collection[I, A, ERR], RET, RW, DIR, ERR],
	traverse Optic[I, S, S, A, A, TRET, TRW, TDIR, ERR],
) CollectionType[I, S, S, A, A, ERR] {
	fromCol := AsReverseGet(toCol)
	return ColTypeP(toCol, fromCol, traverse)
}

// ColTypeP returns a [CollectionType] using the given [Iso]s to convert to and from a [Collection]
//
// See:
//   - [ColType] for a non polymorphic version.
func ColTypeP[I, S, T, A, B any, TRET any, RETI, RETJ TReturnOne, TRW, RWI, RWJ TReadWrite, TDIR any, DIRI, DIRJ TBiDir, ERR any](
	toCol Optic[Void, S, T, Collection[I, A, ERR], Collection[I, B, ERR], RETI, RWI, DIRI, ERR],
	fromColP Optic[Void, Collection[I, A, ERR], Collection[I, B, ERR], S, T, RETJ, RWJ, DIRJ, ERR],
	traverse Optic[I, S, T, A, B, TRET, TRW, TDIR, ERR],
) CollectionType[I, S, T, A, B, ERR] {

	fromColB := CombiIso[ReadWrite, BiDir, ERR, Collection[I, B, ERR], Collection[I, B, ERR], T, T](
		toCol.AsReverseGetter(),
		fromColP.AsReverseGetter(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.ReverseGet{
					OpticTypeExpr: ot,
					Optic: expr.ToCol{
						OpticTypeExpr: expr.NewOpticTypeExpr[Void, T, T, Collection[I, B, ERR], Collection[I, B, ERR], ReturnOne, ReadWrite, BiDir, ERR](),
						I:             reflect.TypeFor[I](),
						A:             reflect.TypeFor[B](),
						B:             reflect.TypeFor[A](),
					},
				}
			},
			toCol,
			fromColP,
		),
	)

	return CollectionType[I, S, T, A, B, ERR]{
		traverse: RetM(Rw(Ud(traverse))),
		toCol:    Ret1(Rw(Bd(toCol))),
		fromColP: Ret1(Rw(Bd(fromColP))),
		fromColB: Ret1(Rw(Bd(fromColB))),
	}
}

// TraverseColType returns an optic that focuses the elements of the source collection.
func TraverseColType[I, S, T, A, B any, ERR any](s CollectionType[I, S, T, A, B, ERR]) Optic[I, S, T, A, B, ReturnMany, ReadWrite, UniDir, ERR] {
	return s.traverse
}

// ColTypeToCol returns an [Iso] that converts the concrete collection type ([]string]) to a [Collection]
//
// This provides easy access to [Collection] operations like [Reversed] without needing to use [makecolop] to generate the specialized collection operations.
//
// See:
//   - [ColToColType] for the reverse
func ColTypeToCol[I, S, T, A, B any, ERR any](s CollectionType[I, S, T, A, B, ERR]) Optic[Void, S, T, Collection[I, A, ERR], Collection[I, B, ERR], ReturnOne, ReadWrite, BiDir, ERR] {
	return s.toCol
}

// ColToColType returns an [Iso] that converts a [Collection] to the concrete collection type ([]string])
//
// This provides easy access to [Collection] operations like [Reversed] without needing to use [makecolop] to generate the specialized collection operations.
//
// See:
//   - [ColTypeToCol] for the reverse
func ColToColType[I, S, T, A, B any, ERR any](s CollectionType[I, S, T, A, B, ERR]) Optic[Void, Collection[I, A, ERR], Collection[I, B, ERR], S, T, ReturnOne, ReadWrite, BiDir, ERR] {
	return s.fromColP
}

// The ColTypeOf combinator focuses on the concrete collection of all the elements in the given optic.
//
// Under modification this collection can be modified and will be rebuilt into the original data structure.
// If the modified collection contains fewer elements the result will use values from the original source.
// If the modified collection contains more elements they will be ignored.
//
// See:
//   - [ColTypeOf] for a [Collection] version
//   - [ColOf] for a [Collection] version
func ColTypeOf[I, SC, TC, ERR, A, S, T, ORET, ORW, ODIR any](ct CollectionType[I, SC, TC, A, A, ERR], o Optic[I, S, T, A, A, ORET, ORW, ODIR, ERR]) Optic[Void, S, T, TC, TC, ReturnOne, ORW, UniDir, ERR] {
	return Ret1(RwL(Ud(EErrMerge(Compose(ColOf(o), ct.fromColB)))))
}

// The ColTypeOfP combinator focuses on a polymorphic concrete collection of all the elements in the given optic.
//
// Under modification this collection can be modified and will be rebuilt into the original data structure.
// If the modified collection contains fewer elements then [ErrUnsafeMissingElement] will be returned
// If the modified collection contains more elements they will be ignored.
//
// See:
//   - [ColOf] for a safe non polymorphic version.
//   - [ColyTypeOfP] for a version that operates on concrete collection types
func ColTypeOfP[I, SC, TC, ERR, A, B, S, T, ORET, ORW, ODIR any](ct CollectionType[I, SC, TC, A, B, ERR], o Optic[I, S, T, A, B, ORET, ORW, ODIR, ERR]) Optic[Void, S, T, SC, TC, ReturnOne, ORW, UniDir, Err] {
	return Ret1(RwL(Ud(EErrMerge(Compose(ColOfP(o), ColSourceErr(ct.fromColP))))))
}
