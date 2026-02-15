package optic

import (
	"github.com/samber/lo"
	"github.com/samber/mo"
)

// ColTypeOp returns an optic that applies the given collection operation optic to the given collection.
// The source and focus are the underlying collection type
func ColTypeOp[I, J, S, T, A, B any, RET TReturnOne, RW any, DIR any, ERR any](s CollectionType[I, S, T, A, B, ERR], o Optic[J, Collection[I, A, ERR], Collection[I, B, ERR], Collection[I, B, ERR], Collection[I, B, ERR], RET, RW, DIR, ERR]) Optic[J, S, T, T, T, ReturnOne, RW, DIR, ERR] {
	return EErrMerge(Ret1(RwL(DirL(ComposeLeft(EErrMerge(RwR(DirR(Compose(s.toCol, o)))), s.fromColB)))))
}

func colTypeT2ToCol[I, S, A any, ERR any](s CollectionType[I, S, S, A, A, ERR]) Optic[Void, lo.Tuple2[S, S], lo.Tuple2[S, S], lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], ReturnOne, ReadWrite, UniDir, ERR] {

	a := EErrR(Compose(T2A[S, S](), s.toCol))
	b := EErrR(Compose(T2B[S, S](), s.toCol))

	t2 := Ret1(EErrMerge(Rw(T2Of(a, b))))

	voidT2 := EErrL(ReIndexed(t2, Const[lo.Tuple2[Void, Void]](Void{}), EqT2[Void]()))

	return voidT2

}

func DiffColType[I, S, A any, RET TReturnOne, ERR any](s CollectionType[I, S, S, A, A, ERR], right S, threshold float64, distance Operation[lo.Tuple2[A, A], float64, RET, ERR], filterDiff DiffType) Optic[Diff[I, A], S, S, mo.Option[A], mo.Option[A], ReturnMany, ReadWrite, UniDir, ERR] {
	return RetM(Rw(Ud(EErrR(EErrMergeL(Compose(
		T2Of(
			Identity[S](),
			IgnoreWrite(Const[S](right)),
		),
		DiffColTypeT2(s, threshold, distance, filterDiff),
	))))))
}

// DiffColTypeT2 returns a [Traversal] that focuses on the differences between 2 [CollectionTypes].
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Note: The source tuple should have the modified [Collection] in the first position and the reference [Collection] in the second position. This is compatible with [T2Dup] returning the modified tree.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffColTypeE] for an impure version
//   - [DiffColTypeT2I] for an index aware version
//   - [Diff] for the index detailing the detected diff.
//   - [Distance] for a convenience constructor for the distance operation.
func DiffColTypeT2[I, S, A any, RET TReturnOne, ERR any](s CollectionType[I, S, S, A, A, ERR], threshold float64, distance Operation[lo.Tuple2[A, A], float64, RET, ERR], filterDiff DiffType) Optic[Diff[I, A], lo.Tuple2[S, S], lo.Tuple2[S, S], mo.Option[A], mo.Option[A], ReturnMany, ReadWrite, UniDir, ERR] {
	t2Col := colTypeT2ToCol(s)

	ixMatch := AsIxMatchT2(s.traverse)

	d := DiffColT2(threshold, distance, ixMatch, filterDiff, true)
	ret := Compose(t2Col, d)

	return EErrMerge(RetM(Rw(Ud(ret))))
}

func DiffColTypeI[I, S, A any, RET TReturnOne, ERR any](s CollectionType[I, S, S, A, A, ERR], right S, threshold float64, distance Operation[lo.Tuple2[ValueI[I, A], ValueI[I, A]], float64, RET, ERR], filterDiff DiffType) Optic[Diff[I, A], S, S, mo.Option[A], mo.Option[A], ReturnMany, ReadWrite, UniDir, ERR] {
	return RetM(Rw(Ud(EErrR(EErrMergeL(Compose(
		T2Of(
			Identity[S](),
			IgnoreWrite(Const[S](right)),
		),
		DiffColTypeT2I(s, threshold, distance, filterDiff),
	))))))
}

// DiffColTypeT2I returns an [Traversal] that focuses on the differences between 2 [CollectionTypes]s.
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Note: The source tuple should have the modified [Collection] in the first position and the reference [Collection] in the second position. This is compatible with [T2Dup] returning the modified tree.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffColTypeT2] for a non index aware version
//   - [Diff] for the index detailing the detected diff.
//   - [DistanceI] for a convenience constructor for the distance operation.
func DiffColTypeT2I[I, S, A any, RET TReturnOne, ERR any](s CollectionType[I, S, S, A, A, ERR], threshold float64, distance Operation[lo.Tuple2[ValueI[I, A], ValueI[I, A]], float64, RET, ERR], filterDiff DiffType) Optic[Diff[I, A], lo.Tuple2[S, S], lo.Tuple2[S, S], mo.Option[A], mo.Option[A], ReturnMany, ReadWrite, UniDir, ERR] {
	t2Col := colTypeT2ToCol(s)
	ixMatch := AsIxMatchT2(s.traverse)
	d := DiffColT2I[I](threshold, distance, ixMatch, filterDiff, true)
	return EErrMerge(RetM(Rw(Ud(Compose(t2Col, d)))))
}

func AppendColTypeT2[I, S, A any, ERR any](s CollectionType[I, S, S, A, A, ERR]) Optic[Void, lo.Tuple2[S, S], lo.Tuple2[S, S], S, S, ReturnOne, ReadOnly, UniDir, ERR] {
	a := Compose3(
		colTypeT2ToCol(s),
		AppendColT2[I, A, ERR](),
		ColToColType(s),
	)
	c := Ret1(Ro(Ud(EErrMerge(EErrMergeL(a)))))
	return c
}

func PrependColTypeT2[I, S, A any, ERR any](s CollectionType[I, S, S, A, A, ERR]) Optic[Void, lo.Tuple2[S, S], lo.Tuple2[S, S], S, S, ReturnOne, ReadOnly, UniDir, ERR] {
	a := Compose3(
		colTypeT2ToCol(s),
		PrependColT2[I, A, ERR](),
		ColToColType(s),
	)
	c := Ret1(Ro(Ud(EErrMerge(EErrMergeL(a)))))
	return c
}
