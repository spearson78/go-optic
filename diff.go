package optic

import (
	"context"
	"fmt"
	"math"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spearson78/go-optic/expr"
)

type DiffType int

const diffUnknown DiffType = 0
const DiffNoFilter DiffType = 0

const (
	DiffAdd DiffType = 1 << iota
	DiffRemove
	DiffModify
	DiffNone
)

func (d DiffType) String() string {
	switch d {
	case diffUnknown:
		return "Unknown"
	case DiffAdd:
		return "Add"
	case DiffRemove:
		return "Remove"
	case DiffModify:
		return "Modify"
	case DiffNone:
		return "None"
	default:
		return "Unknown"
	}
}

type Diff[I, A any] struct {
	BeforeIndex I
	AfterIndex  I
	BeforePos   int
	AfterPos    int
	BeforeValue A
	Type        DiffType
	Distance    float64
}

func (d Diff[I, A]) String() string {

	switch d.Type {
	case DiffAdd:
		return fmt.Sprintf("DiffAdd(Index=%v)", d.AfterIndex)
	case DiffRemove:
		return fmt.Sprintf("DiffRemove(Index=%v,Value=%v)", d.BeforeIndex, d.BeforeValue)
	case DiffNone:
		return fmt.Sprintf("DiffNone(Pos=%v->%v,Index=%v->%v,BeforeValue=%v,Dist=%v)", d.BeforePos, d.AfterPos, d.BeforeIndex, d.AfterIndex, d.BeforeValue, d.Distance)
	case DiffModify:
		return fmt.Sprintf("DiffModify(Pos=%v->%v,Index=%v->%v,BeforeValue=%v,Dist=%v)", d.BeforePos, d.AfterPos, d.BeforeIndex, d.AfterIndex, d.BeforeValue, d.Distance)

	default:
		return "DiffUnknown"
	}

}

// Constructor for a distance [Operation] for use with [DiffCol]
//
// The function should return the edit distance between a and b
// See:
//   - [DistanceI] for an index aware version
//   - [DistanceE] for an impure version.
//   - [DistanceIE] for an index aware, impure version.
func Distance[A any](fnc func(a, b A) float64) Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], float64, float64, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2(fnc)
}

// Constructor for an index aware distance [Operation] for use with [DiffColI]
//
// The function should return the edit distance between a and b
// See:
//   - [Distance] for a pure version
//   - [DistanceE] for an impure version.
//   - [DistanceIE] for an index aware, impure version.
func DistanceI[I, A any](fnc func(ia I, a A, ib I, b A) float64) Optic[Void, lo.Tuple2[ValueI[I, A], ValueI[I, A]], lo.Tuple2[ValueI[I, A], ValueI[I, A]], float64, float64, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2I(fnc)
}

// Constructor for an index aware distance [Operation] for use with [DiffColI]
//
// The function should return the edit distance between a and b
// See:
//   - [Distance] for a pure version
//   - [DistanceE] for an impure version.
//   - [DistanceIE] for an index aware, impure version.
func DistanceE[I, A any](fnc func(ctx context.Context, a A, b A) (float64, error)) Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], float64, float64, ReturnOne, ReadOnly, UniDir, Err] {
	return OpT2E(fnc)
}

// Constructor for an index aware distance [Operation] for use with [DiffColI]
//
// The function should return the edit distance between a and b
// See:
//   - [Distance] for a pure version
//   - [DistanceE] for an impure version.
//   - [DistanceIE] for an index aware, impure version.
func DistanceIE[I, A any](fnc func(ctx context.Context, ia I, a A, ib I, b A) (float64, error)) Optic[Void, lo.Tuple2[ValueI[I, A], ValueI[I, A]], lo.Tuple2[ValueI[I, A], ValueI[I, A]], float64, float64, ReturnOne, ReadOnly, UniDir, Err] {
	return OpT2IE(fnc)
}

// The DistancePercent combinator returns a modified distance optice that scales the returned value by the result of length.
func DistancePercent[I, J, S, T, T2, B, C any, RETO, RETL TReturnOne, RWO, RWL any, DIRO, DIRL any, ERRO, ERRL TPure](distance Optic[I, lo.Tuple2[S, S], T, float64, B, RETO, RWO, DIRO, ERRO], length Optic[J, S, T2, int, C, RETL, RWL, DIRL, ERRL]) Optic[Void, lo.Tuple2[S, S], lo.Tuple2[S, S], float64, float64, ReturnOne, ReadOnly, UniDir, Pure] {
	x := EPure(DivOp(
		distance,
		EErrL(Compose(
			EErrMerge(MaxOp(
				EErrR(Compose(T2AP[S, S, T2](), length)),
				EErrR(Compose(T2BP[S, S, T2](), length)),
			)),
			IsoCast[int, float64](),
		)),
	))

	return x
}

// The DistancePercent combinator returns a modified distance optic that scales the returned value by the result of length.
func DistancePercentI[I, J, K, S, T, T2, B, C any, RETO, RETL TReturnOne, RWO, RWL any, DIRO, DIRL any, ERRO, ERRL any](distance Optic[K, lo.Tuple2[ValueI[I, S], ValueI[I, S]], T, float64, B, RETO, RWO, DIRO, ERRO], length Optic[J, S, T2, int, C, RETL, RWL, DIRL, ERRL]) Optic[Void, lo.Tuple2[ValueI[I, S], ValueI[I, S]], lo.Tuple2[ValueI[I, S], ValueI[I, S]], float64, float64, ReturnOne, ReadOnly, UniDir, CompositionTree[ERRO, ERRL]] {
	x := DivOp(
		distance,
		EErrL(Compose(
			EErrMerge(MaxOp(
				EErrR(Compose(EPure(Compose(T2AP[ValueI[I, S], ValueI[I, S], T2](), Polymorphic[T2, T2](ValueIValue[I, S]()))), length)),
				EErrR(Compose(EPure(Compose(T2BP[ValueI[I, S], ValueI[I, S], T2](), Polymorphic[T2, T2](ValueIValue[I, S]()))), length)),
			)),
			IsoCast[int, float64](),
		)),
	)
	return x
}

func DiffColI[I, A any, RET TReturnOne, ERR any, ERRI TPure](
	right Collection[I, A, ERR],
	threshold float64,
	distance Operation[lo.Tuple2[ValueI[I, A], ValueI[I, A]], float64, RET, ERR],
	ixMatch Predicate[lo.Tuple2[I, I], ERRI],
	filterDiffType DiffType,
	detectPosChange bool,
) Optic[Diff[I, A], Collection[I, A, ERR], Collection[I, A, ERR], mo.Option[A], mo.Option[A], ReturnMany, ReadWrite, UniDir, ERR] {
	return RetM(Rw(Ud(EErrR(EErrMergeL(Compose(
		T2Of(
			Identity[Collection[I, A, ERR]](),
			IgnoreWrite(Const[Collection[I, A, ERR]](right)),
		),
		DiffColT2I(threshold, distance, ixMatch, filterDiffType, detectPosChange),
	))))))
}

// DiffColT2I returns an [Traversal] that focuses on the differences between 2 [Collection]s.
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Note: The source tuple should have the modified [Collection] in the first position and the reference [Collection] in the second position. This is compatible with [T2Dup] returning the modified tree.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffColT2] for a non index aware version
//   - [Diff] for the index detailing the detected diff.
//   - [DistanceI] for a convenience constructor for the distance operation.
func DiffColT2I[I, A any, RET TReturnOne, ERR any, ERRI TPure](
	threshold float64,
	distance Operation[lo.Tuple2[ValueI[I, A], ValueI[I, A]], float64, RET, ERR],
	ixMatch Predicate[lo.Tuple2[I, I], ERRI],
	filterDiffType DiffType,
	detectPosChange bool,
) Optic[Diff[I, A], lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], mo.Option[A], mo.Option[A], ReturnMany, ReadWrite, UniDir, ERR] {
	return diffColT2IE(threshold, distance, ixMatch, filterDiffType, detectPosChange)
}

func diffColT2IE[I, A any, RET TReturnOne, ERR any, ERRI TPure](
	threshold float64,
	distance Operation[lo.Tuple2[ValueI[I, A], ValueI[I, A]], float64, RET, ERR],
	ixMatch Predicate[lo.Tuple2[I, I], ERRI],
	filterDiffType DiffType,
	detectPosChange bool,
) Optic[Diff[I, A], lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], mo.Option[A], mo.Option[A], ReturnMany, ReadWrite, UniDir, ERR] {

	diffIxMatch := func(indexA, indexB Diff[I, A]) bool {

		if indexA.Type != indexB.Type {
			return false
		}

		switch indexA.Type {
		case DiffAdd:
			match := Must(PredGet(context.Background(), ixMatch, lo.T2(indexA.AfterIndex, indexB.AfterIndex)))
			if !match {
				return false
			}

			if detectPosChange && (indexA.AfterPos != indexB.AfterPos) {
				return false
			}

			return true
		case DiffRemove:

			match := Must(PredGet(context.Background(), ixMatch, lo.T2(indexA.BeforeIndex, indexB.BeforeIndex)))
			if !match {
				return false
			}

			dist := Must(distance.AsOpGet()(context.Background(), lo.T2(ValI(indexA.BeforeIndex, indexA.BeforeValue), ValI(indexB.BeforeIndex, indexB.BeforeValue))))
			if dist != 0 {
				return dist == 0
			}

			if detectPosChange && (indexA.BeforePos != indexB.BeforePos) {
				return false
			}

			return true
		case DiffModify:

			match := Must(PredGet(context.Background(), ixMatch, lo.T2(indexA.AfterIndex, indexB.AfterIndex)))
			if !match {
				return false
			}

			if detectPosChange && (indexA.AfterPos != indexB.AfterPos) {
				return false
			}

			match = Must(PredGet(context.Background(), ixMatch, lo.T2(indexA.BeforeIndex, indexB.BeforeIndex)))
			if !match {
				return false
			}

			dist := Must(distance.AsOpGet()(context.Background(), lo.T2(ValI(indexA.BeforeIndex, indexA.BeforeValue), ValI(indexB.BeforeIndex, indexB.BeforeValue))))
			if dist != 0 {
				return dist == 0
			}

			if detectPosChange && (indexA.BeforePos != indexB.BeforePos) {
				return false
			}

			if indexA.Distance != indexB.Distance {
				return false
			}

			return true

		default:
			panic("unknown diff type")
		}
	}

	diffSeq := func(ctx context.Context, source lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]]) (Collection[Diff[I, A], mo.Option[A], ERR], error) {

		ixMatchFnc := PredToIxMatch(ixMatch)

		beforeSlice, err := GetContext(ctx, SliceOf(WithIndex(TraverseColIE[I, A, ERR](ixMatchFnc)), 10), source.B)
		if err != nil {
			return nil, err
		}

		afterSlice, err := GetContext(ctx, SliceOf(WithIndex(TraverseColIE[I, A, ERR](ixMatchFnc)), 10), source.A)
		if err != nil {
			return nil, err
		}

		distances := make([]lo.Tuple3[int, int, float64], 0, len(beforeSlice)*len(afterSlice))
		for beforePos := 0; beforePos < len(beforeSlice); beforePos++ {
			for afterPos := 0; afterPos < len(afterSlice); afterPos++ {

				d, err := distance.AsOpGet()(ctx, lo.T2(beforeSlice[beforePos], afterSlice[afterPos]))
				if err != nil {
					return nil, err
				}

				distances = append(distances,
					lo.T3(
						beforePos,
						afterPos,
						math.Abs(d),
					),
				)
			}
		}

		return ColIE[ERR](
			func(ctx context.Context) SeqIE[Diff[I, A], mo.Option[A]] {
				return func(yield func(ValueIE[Diff[I, A], mo.Option[A]]) bool) {

					emittedBeforeIndexes := make([]DiffType, len(beforeSlice))
					emittedAfterIndexes := make([]Diff[I, A], len(afterSlice))
					emitCount := 0

					seq := MustGet(SeqOf(Ordered(TraverseSlice[lo.Tuple3[int, int, float64]](), OrderBy(T3C[int, int, float64]()))), distances)
					for matchedElements := range seq {
						//matchedElements.A beforePos
						//matchedElements.B afterPos
						//matchedElements.C absoluteDistance
						if ctx.Err() != nil {
							v := ValIE(Diff[I, A]{}, mo.None[A](), ctx.Err())
							if !yield(v) {
								return
							}
						}

						if diffType := emittedAfterIndexes[matchedElements.B]; diffType.Type != diffUnknown {
							//We already emitted this element
							continue
						}

						if emittedBeforeIndexes[matchedElements.A] != diffUnknown {
							//We already emitted this element
							continue
						}

						if matchedElements.C != 0.0 {
							//similar enough to declare an edit
							if matchedElements.C <= threshold {
								i := Diff[I, A]{
									BeforeIndex: beforeSlice[matchedElements.A].index,
									AfterIndex:  afterSlice[matchedElements.B].index,
									BeforePos:   matchedElements.A,
									AfterPos:    matchedElements.B,
									BeforeValue: beforeSlice[matchedElements.A].value,
									Type:        DiffModify,
									Distance:    matchedElements.C,
								}

								emittedBeforeIndexes[matchedElements.A] = DiffModify
								emittedAfterIndexes[matchedElements.B] = i

								emitCount++

								if emitCount == len(afterSlice) {
									//We have emitted all of the original indexes
									break
								}

							} else {
								//The elements are now too different to be consider a modification
								break
							}
						} else {

							diffType := DiffModify
							isModify := (detectPosChange && matchedElements.A != matchedElements.B)
							if !isModify {
								diffType = DiffNone
							}

							i := Diff[I, A]{
								BeforeIndex: beforeSlice[matchedElements.A].index,
								AfterIndex:  afterSlice[matchedElements.B].index,
								BeforePos:   matchedElements.A,
								AfterPos:    matchedElements.B,
								BeforeValue: beforeSlice[matchedElements.A].value,
								Type:        diffType,
								Distance:    matchedElements.C,
							}

							emittedBeforeIndexes[matchedElements.A] = diffType
							emittedAfterIndexes[matchedElements.B] = i
							emitCount++

							if emitCount == len(afterSlice) {
								//We have emitted all of the original indexes
								break
							}

						}
					}

					maxLen := len(beforeSlice)
					if len(afterSlice) > maxLen {
						maxLen = len(afterSlice)
					}

					for i := 0; i < maxLen; i++ {

						if i < len(beforeSlice) {
							if emittedBeforeIndexes[i] == diffUnknown {
								diff := Diff[I, A]{
									BeforeIndex: beforeSlice[i].index,
									BeforePos:   i,
									BeforeValue: beforeSlice[i].value,
									Type:        DiffRemove,
									Distance:    0.0,
								}

								if (filterDiffType & DiffRemove) == 0 {
									if !yield(ValIE(diff, mo.None[A](), nil)) {
										return
									}
								}
							}
						}

						if i < len(afterSlice) {
							switch emittedAfterIndexes[i].Type {
							case diffUnknown:
								diff := Diff[I, A]{
									AfterIndex: afterSlice[i].index,
									AfterPos:   i,
									Type:       DiffAdd,
								}

								if (filterDiffType & DiffAdd) == 0 {
									if !yield(ValIE(diff, mo.Some(afterSlice[i].value), nil)) {
										return
									}
								}

							case DiffModify:

								if (filterDiffType & DiffModify) == 0 {
									if !yield(ValIE(emittedAfterIndexes[i], mo.Some(afterSlice[i].value), nil)) {
										return
									}
								}

							default:

								if (filterDiffType & DiffNone) == 0 {

									if !yield(ValIE(emittedAfterIndexes[i], mo.Some(afterSlice[i].value), nil)) {
										return
									}
								}
							}
						}

					}
				}
			},
			nil,
			diffIxMatch,
			nil,
		), nil

	}

	return CombiTraversal[ReturnMany, ReadWrite, ERR, Diff[I, A], lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], mo.Option[A], mo.Option[A]](
		func(ctx context.Context, source lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]]) SeqIE[Diff[I, A], mo.Option[A]] {
			seq, err := diffSeq(ctx, source)
			if err != nil {
				return func(yield func(ValueIE[Diff[I, A], mo.Option[A]]) bool) {
					var i Diff[I, A]
					yield(ValIE(i, mo.None[A](), err))
				}
			}

			return seq.AsIter()(ctx)
		},
		nil,
		func(ctx context.Context, fmap func(index Diff[I, A], focus mo.Option[A]) (mo.Option[A], error), source lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]]) (lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], error) {
			seq, err := diffSeq(ctx, source)
			if err != nil {
				return lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]]{}, err
			}

			var retErr error
			var afterSlice []ValueIE[I, A]

			seq.AsIter()(ctx)(func(val ValueIE[Diff[I, A], mo.Option[A]]) bool {
				index, focus, err := val.Get()
				if err != nil {
					retErr = err
					return false
				}

				yieldIndex := index.AfterIndex
				if index.Type == DiffRemove {
					yieldIndex = index.BeforeIndex
				}

				if (filterDiffType & index.Type) != 0 {
					//No fmap
					if val, ok := focus.Get(); ok {
						afterSlice = append(afterSlice, ValIE(yieldIndex, val, nil))
					}
					return true
				} else {

					newOpt, err := fmap(index, focus)
					if newVal, ok := newOpt.Get(); ok || err != nil {
						//Element was resurrected
						if err != nil {
							retErr = err
							return false
						}
						//Under modification diff may be pure
						afterSlice = append(afterSlice, ValIE(yieldIndex, newVal, nil))
					} else {
						//Element is removed
					}

					return true
				}
			})

			if retErr != nil {
				return lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]]{}, retErr
			}

			return lo.T2[Collection[I, A, ERR], Collection[I, A, ERR]](ValColIE[I, A, ERR](
				func(a, b I) bool {
					return Must(PredGet(context.Background(), ixMatch, lo.T2(a, b)))
				},
				afterSlice...,
			), source.B), nil

		},
		nil,
		diffIxMatch,
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.DiffCol{
					OpticTypeExpr: ot,
					Threshold:     threshold,
					Distance:      distance.AsExpr(),
					IxMatch:       ixMatch.AsExpr(),
				}
			},
			distance,
			ixMatch,
		),
	)
}

func DiffCol[I, A any, RET TReturnOne, ERR any, ERRI TPure](
	right Collection[I, A, ERR],
	threshold float64,
	distance Operation[lo.Tuple2[A, A], float64, RET, ERR],
	ixMatch Predicate[lo.Tuple2[I, I], ERRI],
	filterDiffType DiffType,
	detectPosChange bool,
) Optic[Diff[I, A], Collection[I, A, ERR], Collection[I, A, ERR], mo.Option[A], mo.Option[A], ReturnMany, ReadWrite, UniDir, ERR] {
	return RetM(Rw(Ud(EErrR(EErrMergeL(Compose(
		T2Of(
			Identity[Collection[I, A, ERR]](),
			IgnoreWrite(Const[Collection[I, A, ERR]](right)),
		),
		DiffColT2(threshold, distance, ixMatch, filterDiffType, detectPosChange),
	))))))
}

// DiffColT2 returns an [Traversal] that focuses on the differences between 2 [Collection].
// The distance [Operation] is used to determine if an element has either changed position or been modified.
// Elements are considered same if they have the smallest distance compared to all other elements and the distance is lower than the threshold.
//
// Note: The source tuple should have the modified [Collection] in the first position and the reference [Collection] in the second position. This is compatible with [T2Dup] returning the modified tree.
//
// Details of the diff are encoded in the index of the focused elements.
//
// See:
//   - [DiffColE] for an impure version
//   - [DiffColT2I] for an index aware version
//   - [Diff] for the index detailing the detected diff.
//   - [Distance] for a convenience constructor for the distance operation.
func DiffColT2[I, A any, RET TReturnOne, ERR any, ERRI TPure](threshold float64, distance Operation[lo.Tuple2[A, A], float64, RET, ERR], ixMatch Predicate[lo.Tuple2[I, I], ERRI], filterDiffType DiffType, detectPosChange bool) Optic[Diff[I, A], lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], mo.Option[A], mo.Option[A], ReturnMany, ReadWrite, UniDir, ERR] {
	distFnc := distance.AsOpGet()

	idist := rawGetterF[Void, lo.Tuple2[ValueI[I, A], ValueI[I, A]], float64, ERR](
		func(ctx context.Context, source lo.Tuple2[ValueI[I, A], ValueI[I, A]]) (Void, float64, error) {
			ret, err := distFnc(ctx, lo.T2(source.A.value, source.B.value))
			return Void{}, ret, err
		},
		IxMatchVoid(),
		distance.AsExprHandler(),
		distance.AsExpr,
	)

	return DiffColT2I(threshold, idist, ixMatch, filterDiffType, detectPosChange)
}
