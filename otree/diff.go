package otree

import (
	"context"
	"reflect"

	"github.com/samber/lo"
	"github.com/samber/mo"
	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

func DiffTreeI[I, J, A any, RET TReturnOne, RW, DIR, ERR any, DRET TReturnOne, ERRI TPure](right mo.Option[A], children Optic[J, A, A, Collection[I, A, ERR], Collection[I, A, ERR], RET, RW, DIR, ERR], threshold float64, distance Operation[lo.Tuple2[ValueI[*PathNode[I], A], ValueI[*PathNode[I], A]], float64, DRET, ERR], ixMatch Predicate[lo.Tuple2[I, I], ERRI], filterDiff DiffType, detectPosChange bool) Optic[Diff[*PathNode[I], A], mo.Option[A], mo.Option[A], mo.Option[A], mo.Option[A], ReturnMany, RW, UniDir, ERR] {
	return RetM(Ud(RwR(RwMergeL(EErrR(EErrMergeL(Compose(
		T2Of(
			Identity[mo.Option[A]](),
			IgnoreWrite(Const[mo.Option[A]](right)),
		),
		DiffTreeT2I(children, threshold, distance, ixMatch, filterDiff, detectPosChange),
	)))))))

}

// When a node is added only the added node is reported in the diff not all of it's children recursively.
func DiffTreeT2I[I, J, A any, RET TReturnOne, RW, DIR, ERR any, DRET TReturnOne, ERRI TPure](children Optic[J, A, A, Collection[I, A, ERR], Collection[I, A, ERR], RET, RW, DIR, ERR], threshold float64, distance Operation[lo.Tuple2[ValueI[*PathNode[I], A], ValueI[*PathNode[I], A]], float64, DRET, ERR], ixMatch Predicate[lo.Tuple2[I, I], ERRI], filterDiff DiffType, detectPosChange bool) Optic[Diff[*PathNode[I], A], lo.Tuple2[mo.Option[A], mo.Option[A]], lo.Tuple2[mo.Option[A], mo.Option[A]], mo.Option[A], mo.Option[A], ReturnMany, RW, UniDir, ERR] {
	var iter func(ctx context.Context, afterPath *PathNode[I], beforePath *PathNode[I], node A, origNode A, yield func(ValueIE[Diff[*PathNode[I], A], mo.Option[A]]) bool) bool
	iter = func(ctx context.Context, afterPath *PathNode[I], beforePath *PathNode[I], node A, origNode A, yield func(ValueIE[Diff[*PathNode[I], A], mo.Option[A]]) bool) bool {

		cont := true

		diffCol := DiffColT2I(threshold, OpE(func(ctx context.Context, p lo.Tuple2[ValueI[I, A], ValueI[I, A]]) (float64, error) {
			return distance.AsOpGet()(ctx, lo.T2(ValI(beforePath.Append(p.A.Index()), p.A.Value()), ValI(afterPath.Append(p.B.Index()), p.B.Value())))
		}), ixMatch, DiffNoFilter, detectPosChange)

		afterChildren, err := children.AsOpGet()(ctx, node)
		if err != nil {
			var i Diff[*PathNode[I], A]
			var a mo.Option[A]
			return yield(ValIE(i, a, err))
		}

		beforeChildren, err := children.AsOpGet()(ctx, origNode)
		if err != nil {
			var i Diff[*PathNode[I], A]
			var a mo.Option[A]
			return yield(ValIE(i, a, err))
		}

		diffCol.AsIter()(ctx, lo.T2(
			ColErr(afterChildren),
			ColErr(beforeChildren),
		))(func(val ValueIE[Diff[I, A], mo.Option[A]]) bool {
			index, focus, err := val.Get()
			if err != nil {
				var i Diff[*PathNode[I], A]
				var a mo.Option[A]
				cont = yield(ValIE(i, a, err))
				return cont
			}

			pathDiff := Diff[*PathNode[I], A]{
				BeforePos:   index.BeforePos,
				AfterPos:    index.AfterPos,
				BeforeValue: index.BeforeValue,
				Type:        index.Type,
				Distance:    index.Distance,
			}

			if index.Type == DiffModify || index.Type == DiffAdd {
				pathDiff.AfterIndex = afterPath.Append(index.AfterIndex)
			}

			if index.Type == DiffModify || index.Type == DiffRemove {
				pathDiff.BeforeIndex = beforePath.Append(index.BeforeIndex)
			}

			if (filterDiff & index.Type) == 0 {
				cont = yield(ValIE(pathDiff, focus, nil))
			}

			if cont && (index.Type == DiffModify || index.Type == DiffNone) {
				if val, ok := focus.Get(); ok {
					cont = iter(ctx, pathDiff.AfterIndex, pathDiff.BeforeIndex, val, index.BeforeValue, yield)
				}
			}

			return cont
		})

		return cont
	}

	var modify func(ctx context.Context, afterPath *PathNode[I], beforePath *PathNode[I], node A, origNode A, fmap func(index Diff[*PathNode[I], A], focus mo.Option[A]) (mo.Option[A], error)) (A, error)
	modify = func(ctx context.Context, afterPath *PathNode[I], beforePath *PathNode[I], node A, origNode A, fmap func(index Diff[*PathNode[I], A], focus mo.Option[A]) (mo.Option[A], error)) (A, error) {

		diffCol := DiffColT2I[I, A](threshold, OpE(func(ctx context.Context, p lo.Tuple2[ValueI[I, A], ValueI[I, A]]) (float64, error) {
			return distance.AsOpGet()(ctx, lo.T2(ValI(beforePath.Append(p.A.Index()), p.A.Value()), ValI(afterPath.Append(p.B.Index()), p.B.Value())))
		}), ixMatch, DiffNoFilter, true)

		beforeChildren, err := children.AsOpGet()(ctx, origNode)
		if err != nil {
			var a A
			return a, err
		}

		return children.AsModify()(ctx, func(index J, afterChildren Collection[I, A, ERR]) (Collection[I, A, ERR], error) {

			newChildren, err := diffCol.AsModify()(ctx, func(index Diff[I, A], focus mo.Option[A]) (mo.Option[A], error) {

				pathDiff := Diff[*PathNode[I], A]{
					BeforePos:   index.BeforePos,
					AfterPos:    index.AfterPos,
					BeforeValue: index.BeforeValue,
					Type:        index.Type,
					Distance:    index.Distance,
				}

				if index.Type == DiffModify || index.Type == DiffAdd {
					pathDiff.AfterIndex = afterPath.Append(index.AfterIndex)
				}

				if index.Type == DiffModify || index.Type == DiffRemove {
					pathDiff.BeforeIndex = beforePath.Append(index.BeforeIndex)
				}

				if (filterDiff & index.Type) == 0 {
					focus, err = fmap(pathDiff, focus)
					if err != nil {
						return mo.None[A](), err
					}
				}

				if index.Type == DiffModify || index.Type == DiffNone {
					if val, ok := focus.Get(); ok {
						ret, err := modify(ctx, pathDiff.AfterIndex, pathDiff.BeforeIndex, val, index.BeforeValue, fmap)
						return mo.Some(ret), err
					} else {
						return mo.None[A](), nil
					}
				} else {
					return focus, nil
				}

			}, lo.T2(ColErr(afterChildren), ColErr(beforeChildren)))

			return ColIE[ERR](newChildren.A.AsIter(), newChildren.A.AsIxGet(), newChildren.A.AsIxMatch(), newChildren.A.AsLengthGetter()), err

		}, node)
	}

	return CombiTraversal[ReturnMany, RW, ERR, Diff[*PathNode[I], A], lo.Tuple2[mo.Option[A], mo.Option[A]], lo.Tuple2[mo.Option[A], mo.Option[A]], mo.Option[A], mo.Option[A]](
		func(ctx context.Context, source lo.Tuple2[mo.Option[A], mo.Option[A]]) SeqIE[Diff[*PathNode[I], A], mo.Option[A]] {
			return func(yield func(ValueIE[Diff[*PathNode[I], A], mo.Option[A]]) bool) {

				var path *PathNode[I]
				if a, aok := source.A.Get(); aok {
					if b, bok := source.B.Get(); bok {

						dist, err := distance.AsOpGet()(ctx, lo.T2(ValI(path, a), ValI(path, b)))
						if err != nil {
							yield(ValIE(Diff[*PathNode[I], A]{}, mo.None[A](), err))
							return
						}
						cont := true
						if dist != 0.0 {
							if dist <= threshold {
								if (filterDiff & DiffModify) == 0 {
									cont = yield(
										ValIE(
											Diff[*PathNode[I], A]{
												BeforeIndex: nil,
												AfterIndex:  nil,
												BeforePos:   0,
												AfterPos:    0,
												BeforeValue: b,
												Type:        DiffModify,
												Distance:    dist,
											},
											source.A,
											nil,
										),
									)
								}
							}
						}

						if cont {
							iter(ctx, nil, nil, a, b, yield)
						}

					} else {
						if (filterDiff & DiffAdd) == 0 {
							yield(
								ValIE(
									Diff[*PathNode[I], A]{
										BeforeIndex: nil,
										AfterIndex:  nil,
										BeforePos:   0,
										AfterPos:    0,
										//BeforeValue:
										Type: DiffAdd,
									},
									source.A,
									nil,
								),
							)
						}
					}
				} else {
					if b, bok := source.B.Get(); bok {
						if (filterDiff & DiffRemove) == 0 {
							yield(
								ValIE(
									Diff[*PathNode[I], A]{
										BeforeIndex: nil,
										AfterIndex:  nil,
										BeforePos:   0,
										AfterPos:    0,
										BeforeValue: b,
										Type:        DiffRemove,
									},
									source.A,
									nil,
								),
							)
						}
					} else {
						if (filterDiff & DiffNone) == 0 {
							yield(
								ValIE(
									Diff[*PathNode[I], A]{
										BeforeIndex: nil,
										AfterIndex:  nil,
										BeforePos:   0,
										AfterPos:    0,
										Type:        DiffNone,
										Distance:    0,
									},
									source.A,
									nil,
								),
							)
						}
					}
				}

			}
		},
		nil,
		func(ctx context.Context, fmap func(index Diff[*PathNode[I], A], focus mo.Option[A]) (mo.Option[A], error), source lo.Tuple2[mo.Option[A], mo.Option[A]]) (lo.Tuple2[mo.Option[A], mo.Option[A]], error) {

			diff := Diff[*PathNode[I], A]{
				BeforeIndex: nil,
				AfterIndex:  nil,
				BeforePos:   0,
				AfterPos:    0,
				Type:        DiffNone,
				Distance:    0,
			}
			doFmap := false

			var path *PathNode[I]
			if a, aok := source.A.Get(); aok {
				if b, bok := source.B.Get(); bok {

					dist, err := distance.AsOpGet()(ctx, lo.T2(ValI(path, a), ValI(path, b)))
					if err != nil {
						return lo.Tuple2[mo.Option[A], mo.Option[A]]{}, err
					}
					if dist != 0.0 {
						if dist <= threshold {
							if (filterDiff & DiffModify) == 0 {
								doFmap = true
								diff = Diff[*PathNode[I], A]{
									BeforeIndex: nil,
									AfterIndex:  nil,
									BeforePos:   0,
									AfterPos:    0,
									BeforeValue: b,
									Type:        DiffModify,
									Distance:    dist,
								}

							}

						}
					}

				} else {
					if (filterDiff & DiffAdd) == 0 {
						doFmap = true
						diff = Diff[*PathNode[I], A]{
							BeforeIndex: nil,
							AfterIndex:  nil,
							BeforePos:   0,
							AfterPos:    0,
							//BeforeValue:
							Type: DiffAdd,
						}
					}
				}
			} else {
				if b, bok := source.B.Get(); bok {
					if (filterDiff & DiffRemove) == 0 {
						doFmap = true
						diff = Diff[*PathNode[I], A]{
							BeforeIndex: nil,
							AfterIndex:  nil,
							BeforePos:   0,
							AfterPos:    0,
							BeforeValue: b,
							Type:        DiffRemove,
						}
					}

				} else {
					if (filterDiff & DiffNone) == 0 {
						doFmap = true
						diff = Diff[*PathNode[I], A]{
							BeforeIndex: nil,
							AfterIndex:  nil,
							BeforePos:   0,
							AfterPos:    0,
							Type:        DiffNone,
							Distance:    0,
						}
					}
				}
			}

			if doFmap {
				mapped, err := fmap(diff, source.A)
				if err != nil {
					return lo.Tuple2[mo.Option[A], mo.Option[A]]{}, err
				}

				if m, ok := mapped.Get(); ok {
					if b, ok := source.B.Get(); ok {
						modified, err := modify(ctx, path, path, m, b, fmap)
						if err != nil {
							return lo.Tuple2[mo.Option[A], mo.Option[A]]{}, err
						}
						return lo.T2(mo.Some(modified), source.B), nil
					} else {
						return lo.T2(mapped, source.B), nil
					}
				} else {
					//Root deleted
					return lo.T2(mo.None[A](), source.B), nil
				}
			} else {
				if a, ok := source.A.Get(); ok {
					if b, ok := source.B.Get(); ok {
						modified, err := modify(ctx, path, path, a, b, fmap)
						if err != nil {
							return lo.Tuple2[mo.Option[A], mo.Option[A]]{}, err
						}
						return lo.T2(mo.Some(modified), source.B), nil
					} else {
						return source, nil
					}
				} else {
					return source, nil
				}
			}

		},
		nil,
		func(indexA, indexB Diff[*PathNode[I], A]) bool {

			if indexA.Type != indexB.Type {
				return false
			}

			switch indexA.Type {
			case DiffAdd:

				var a A
				dist := Must(distance.AsOpGet()(context.Background(), lo.T2(ValI(indexA.AfterIndex, a), ValI(indexB.AfterIndex, a))))
				if dist != 0 {
					return false
				}

				if indexA.AfterPos != indexB.AfterPos {
					return false
				}

				return true
			case DiffRemove:

				dist := Must(distance.AsOpGet()(context.Background(), lo.T2(ValI(indexA.BeforeIndex, indexA.BeforeValue), ValI(indexB.BeforeIndex, indexB.BeforeValue))))
				if dist != 0 {
					return false
				}

				if indexA.BeforePos != indexB.BeforePos {
					return false
				}

				return true
			case DiffModify:

				var a A
				dist := Must(distance.AsOpGet()(context.Background(), lo.T2(ValI(indexA.AfterIndex, a), ValI(indexB.AfterIndex, a))))
				if dist != 0 {
					return false
				}

				if indexA.AfterPos != indexB.AfterPos {
					return false
				}

				dist = Must(distance.AsOpGet()(context.Background(), lo.T2(ValI(indexA.BeforeIndex, indexA.BeforeValue), ValI(indexB.BeforeIndex, indexB.BeforeValue))))
				if dist != 0 {
					return false
				}

				if indexA.BeforePos != indexB.BeforePos {
					return false
				}

				if indexA.Distance != indexB.Distance {
					return false
				}

				return true

			default:
				panic("unknown diff type")
			}

		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.DiffTree{
					OpticTypeExpr: ot,
					Children:      children.AsExpr(),
					I:             reflect.TypeFor[I](),
					A:             reflect.TypeFor[A](),
					Threshold:     threshold,
					Distance:      distance.AsExpr(),
					IxMatch:       ixMatch.AsExpr(),
				}
			},
			children,
			distance,
			ixMatch,
		),
	)
}

func DiffTree[I, J, A any, RET TReturnOne, RW any, DIR any, ERR any, DRET TReturnOne, ERRI TPure](right mo.Option[A], children Optic[J, A, A, Collection[I, A, ERR], Collection[I, A, ERR], RET, RW, DIR, ERR], threshold float64, distance Operation[lo.Tuple2[A, A], float64, DRET, ERR], ixMatch Predicate[lo.Tuple2[I, I], ERRI], filterDiff DiffType, detectPosChange bool) Optic[Diff[*PathNode[I], A], mo.Option[A], mo.Option[A], mo.Option[A], mo.Option[A], ReturnMany, RW, UniDir, ERR] {
	return RetM(Ud(RwR(RwMergeL(EErrR(EErrMergeL(Compose(
		T2Of(
			Identity[mo.Option[A]](),
			IgnoreWrite(Const[mo.Option[A]](right)),
		),
		DiffTreeT2(children, threshold, distance, ixMatch, filterDiff, detectPosChange),
	)))))))
}

func DiffTreeT2[I, J, A any, RET TReturnOne, RW any, DIR any, ERR any, DRET TReturnOne, ERRI TPure](children Optic[J, A, A, Collection[I, A, ERR], Collection[I, A, ERR], RET, RW, DIR, ERR], threshold float64, distance Operation[lo.Tuple2[A, A], float64, DRET, ERR], ixMatch Predicate[lo.Tuple2[I, I], ERRI], filterDiff DiffType, detectPosChange bool) Optic[Diff[*PathNode[I], A], lo.Tuple2[mo.Option[A], mo.Option[A]], lo.Tuple2[mo.Option[A], mo.Option[A]], mo.Option[A], mo.Option[A], ReturnMany, RW, UniDir, ERR] {

	idistance := EErrR(Compose(
		DelveT2(ValueIValue[*PathNode[I], A]()),
		OpToOptic(distance),
	))

	return DiffTreeT2I(children, threshold, idistance, ixMatch, filterDiff, detectPosChange)
}
