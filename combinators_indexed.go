package optic

import (
	"context"
	"errors"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spearson78/go-optic/expr"
)

// The TakingWhileI combinator focuses on the elements of the given optic until the indexed predicate returns false.
//
// See [TakingWhile] for a non indexed version.
func TakingWhileI[I, S, T, A, RETI, RW any, DIR, ERR, ERRP any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR], pred PredicateI[I, A, ERRP]) Optic[I, S, T, A, A, ReturnMany, RW, UniDir, CompositionTree[ERR, ERRP]] {
	return EErrR(EErrTransR(filtered(o, Const[int](true), pred, FilterContinue, FilterStop)))
}

// The DroppingWhileI combinator skips the the elements of the given optic until the indexed predicate returns false and focuses on the remaining elements.
//
// See [DroppingWhile] for a non indexed version.
func DroppingWhileI[I, S, T, A, RETI, RW any, DIR, ERR, ERRP any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR], pred PredicateI[I, A, ERRP]) Optic[I, S, T, A, A, ReturnMany, RW, UniDir, CompositionTree[ERR, ERRP]] {
	return EErrR(EErrTransR(filtered(o, Const[int](true), NotOp(pred), FilterYieldAll, FilterContinue)))
}

// The TrimmingWhileI combinator skips leading and trailing elements until an element that does not match the indexed predicate is found.
//
// See [TrimmingWhile] for a non indexed version.
func TrimmingWhileI[I, S, T, A, RETI, RW any, DIR, ERR, ERRP any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR], pred PredicateI[I, A, ERRP]) Optic[I, S, T, A, A, ReturnMany, RW, UniDir, CompositionTree[ERR, ERRP]] {
	return EErrSwap(EErrMergeL(EErrTransL(EErrSwapL(Reversed(DroppingWhileI(Reversed(DroppingWhileI(o, pred)), pred))))))
}

// The FilteredI combinator focuses on the elements of the given optic that match the given indexed predicate.
//
// See:
//   - [FilteredBy] for a filtering combinator that yields the source instead of the focus.
//   - [Filtered] for a non indexed version
func FilteredI[I, S, T, A, RETI any, RW any, DIR any, ERR, ERRP any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR], pred PredicateI[I, A, ERRP]) Optic[I, S, T, A, A, ReturnMany, RW, UniDir, CompositionTree[ERR, ERRP]] {
	return EErrR(EErrTransR(filtered(o, Const[int](true), pred, FilterContinue, FilterContinue)))
}

func FilteredIP[I, S, T, A, B, RETI any, RW any, DIR any, ERR, ERRP any](o Optic[I, S, T, A, B, RETI, RW, DIR, ERR], pred PredicateI[I, A, ERRP]) Optic[I, S, T, A, A, ReturnMany, ReadOnly, UniDir, CompositionTree[ERR, ERRP]] {
	return EErrR(EErrTransR(filtered(Polymorphic[T, A](o), Const[int](true), pred, FilterContinue, FilterContinue)))
}

type FilterMode byte

const (
	FilterStop FilterMode = iota
	FilterContinue
	FilterYieldAll
)

func filtered[I, S, T, A, RETI any, RW any, DIR any, ERR, ERRP, ERRI any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR], posPred Predicate[int, ERRP], pred PredicateI[I, A, ERRI], matchMode FilterMode, noMatchMode FilterMode) Optic[I, S, T, A, A, ReturnMany, RW, UniDir, CompositionTree[ERR, CompositionTree[ERRP, ERRI]]] {
	return CombiTraversal[ReturnMany, RW, CompositionTree[ERR, CompositionTree[ERRP, ERRI]], I, S, T, A, A](
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(ValueIE[I, A]) bool) {
				pos := 0
				yieldAll := false
				posPredRes := false
				posPredResPrimed := false

				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					index, focus, focusErr := val.Get()
					focusErr = JoinCtxErr(ctx, focusErr)
					if focusErr != nil {
						return yield(ValIE(index, focus, focusErr))
					}

					if yieldAll {
						return yield(ValIE(index, focus, focusErr))
					}

					if !posPredResPrimed {
						var err error
						posPredResPrimed = false
						posPredRes, err = PredGet(ctx, posPred, pos)
						pos++
						err = JoinCtxErr(ctx, err)
						if err != nil {
							return yield(ValIE(index, focus, err))
						}
					}

					if posPredRes {
						predRes, err := PredGetI(ctx, pred, focus, index)
						err = JoinCtxErr(ctx, err)
						if err != nil {
							return yield(ValIE(index, focus, err))
						}

						if predRes {

							cont := true

							switch matchMode {
							case FilterContinue:
								cont = yield(ValIE(index, focus, ctx.Err()))
							case FilterStop:
								yield(ValIE(index, focus, ctx.Err()))
								cont = false
							case FilterYieldAll:
								yieldAll = true
								cont = yield(ValIE(index, focus, ctx.Err()))
							default:
								panic("unknown filter mode")
							}

							if cont {

								if noMatchMode == FilterStop {
									posPredResPrimed = true
									posPredRes, err = PredGet(ctx, posPred, pos)
									pos++
									err = JoinCtxErr(ctx, err)
									if err != nil {
										return yield(ValIE(index, focus, err))
									}

									cont = posPredRes
								}
							}

							return cont
						}
					}

					//predicate did not match
					switch noMatchMode {
					case FilterContinue:
						return true
					case FilterStop:
						return false
					case FilterYieldAll:
						yieldAll = true
						return yield(ValIE(index, focus, ctx.Err()))
					default:
						panic("unknown filter mode")
					}

				})
			}
		},
		nil,
		func(ctx context.Context, fmap func(index I, focus A) (A, error), source S) (T, error) {
			pos := 0
			stopped := false
			fmapAll := false
			return o.AsModify()(ctx, func(focusIndex I, focus A) (A, error) {
				if fmapAll {
					pos++
					ret, err := fmap(focusIndex, focus)
					return ret, JoinCtxErr(ctx, err)
				}

				if stopped {
					return focus, ctx.Err()
				}
				predRes, err := PredGet(ctx, posPred, pos)
				pos++
				err = JoinCtxErr(ctx, err)
				if err != nil {
					var fa A
					return fa, err
				}

				if predRes {
					ixPredRes, err := PredGetI(ctx, pred, focus, focusIndex)
					err = JoinCtxErr(ctx, err)
					if err != nil {
						var fa A
						return fa, err
					}
					if ixPredRes {

						switch matchMode {
						case FilterContinue:
							//No op
						case FilterStop:
							stopped = true
						case FilterYieldAll:
							fmapAll = true
						default:
							panic("unknown filter mode")
						}

						ret, err := fmap(focusIndex, focus)

						return ret, err
					}
				}

				//Predicate did not match

				switch noMatchMode {
				case FilterContinue:
					return focus, err
				case FilterStop:
					stopped = true
					return focus, err
				case FilterYieldAll:
					fmapAll = true
					return fmap(focusIndex, focus)
				default:
					panic("unknown filter mode")
				}

			}, source)
		},
		nil, //We need to calculate the pos so IxGet has to fallback to the iterator.
		o.AsIxMatch(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Filtered{
					OpticTypeExpr: ot,
					PosPred:       posPred.AsExpr(),
					Pred:          pred.AsExpr(),
					Optic:         o.AsExpr(),
					NoMatchMode:   expr.FilterMode(noMatchMode),
					MatchMode:     expr.FilterMode(matchMode),
				}
			},
			posPred,
			pred,
			o,
		),
	)
}

// The Index combinator returns an optic that focuses on the given index in the given optic.
//
// Note: multiple elements may share the same index.
// Note: Index is not able to add or remove elements.
//
// See:
//   - [Element] for an optic that focuses on a element by position.
//   - [AtMap] for an optic that is able to add/remove elements from a map.
//   - [Lookup] for a version that looks up an index in a fixed source.
func Index[I, S, T, A, RETI any, RW, DIR, ERR any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR], index I) Optic[I, S, T, A, A, ReturnMany, RW, UniDir, ERR] {
	return CombiTraversal[ReturnMany, RW, ERR, I, S, T, A, A](
		func(ctx context.Context, source S) SeqIE[I, A] {
			return o.AsIxGetter()(ctx, index, source)
		},
		nil,
		func(ctx context.Context, fmap func(index I, focus A) (A, error), source S) (T, error) {
			return o.AsModify()(ctx, func(focusIndex I, focus A) (A, error) {
				match := o.AsIxMatch()(focusIndex, index)

				if match {
					ret, err := fmap(focusIndex, focus)
					return ret, JoinCtxErr(ctx, err)
				} else {
					return focus, ctx.Err()
				}
			}, source)
		},
		func(ctx context.Context, ixGetIndex I, source S) SeqIE[I, A] {
			match := o.AsIxMatch()(ixGetIndex, index)

			if match {
				return o.AsIxGetter()(ctx, index, source)
			} else {
				return func(yield func(val ValueIE[I, A]) bool) {}
			}
		},
		o.AsIxMatch(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Index{
					OpticTypeExpr: ot,
					Index:         index,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The Lookup combinator returns an optic that focuses on the indexed element within the given optic and source.
//
// See:
//   - [Index] for a version that looks up a fixed index.
func Lookup[I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], source S) Optic[I, I, I, mo.Option[A], mo.Option[A], ReturnOne, ReadOnly, UniDir, ERR] {
	return CombiGetter[ERR, I, I, I, mo.Option[A], mo.Option[A]](
		func(ctx context.Context, focus I) (I, mo.Option[A], error) {

			var ret A
			var retErr error
			found := false

			o.AsIxGetter()(ctx, focus, source)(func(val ValueIE[I, A]) bool {
				_, focus, err := val.Get()
				ret = focus
				retErr = err
				found = true
				return false
			})

			return focus, mo.TupleToOption(ret, found), JoinCtxErr(ctx, retErr)
		},
		o.AsIxMatch(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Lookup{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
					Source:        source,
				}
			},
			o,
		),
	)
}

const indexTunnelKey = "github.com/spearson78/go-optics/indexTunnel"

type indexTunnel struct {
	index any
	ok    bool
}

// The WithIndex combinator returns an optic that lifts the index value into a combined focus with the original focused value..
//
// Note: Modifying the ValueIE index is not supported. The original index will be retained.
func WithIndex[I, S, T, A, B, RETI, RW, DIR, ERR any](o Optic[I, S, T, A, B, RETI, RW, DIR, ERR]) Optic[I, S, T, ValueI[I, A], ValueI[I, B], RETI, RW, UniDir, ERR] {
	return CombiTraversal[RETI, RW, ERR, I, S, T, ValueI[I, A], ValueI[I, B]](
		func(ctx context.Context, source S) SeqIE[I, ValueI[I, A]] {
			return func(yield func(val ValueIE[I, ValueI[I, A]]) bool) {
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					index, focus, err := val.Get()
					return yield(ValIE(index, ValI(index, focus), err))
				})
			}
		},
		o.AsLengthGetter(),
		func(ctx context.Context, fmap func(index I, focus ValueI[I, A]) (ValueI[I, B], error), source S) (T, error) {
			ixTunnel, _ := ctx.Value(indexTunnelKey).(*indexTunnel)

			return o.AsModify()(ctx, func(index I, focus A) (B, error) {
				ixval, err := fmap(index, ValI(index, focus))
				if ixTunnel != nil {
					ixTunnel.index = ixval.index
					ixTunnel.ok = true
				}
				return ixval.value, JoinCtxErr(ctx, err)
			}, source)
		},
		func(ctx context.Context, index I, source S) SeqIE[I, ValueI[I, A]] {
			return func(yield func(val ValueIE[I, ValueI[I, A]]) bool) {
				o.AsIxGetter()(ctx, index, source)(func(val ValueIE[I, A]) bool {
					index, focus, err := val.Get()
					return yield(ValIE(index, ValI(index, focus), err))
				})
			}
		},
		o.AsIxMatch(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.WithIndex{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The Indices combinator focuses on elements in the given optic where the predicate matches the elements index.
//
// See:
//   - [FilteredI] for a more general filter that has access to both the index and value
//   - [Filtered] for a non indexed version
func Indices[I, S, T, A, RET, RW, DIR, ERR, ERRP any](o Optic[I, S, T, A, A, RET, RW, DIR, ERR], pred Predicate[I, ERRP]) Optic[I, S, T, A, A, ReturnMany, RW, UniDir, CompositionTree[ERR, ERRP]] {
	return FilteredI(o, PredOnIx[A](pred))
}

// The Indexing combinator returns an optic where the index is replaced with an integer position.
//
// See [ReIndexed] for a version that replaces the index with an arbitrary value.
func Indexing[I, RETI, S, T, A, B, RW, DIR, ERR any](o Optic[I, S, T, A, B, RETI, RW, DIR, ERR]) Optic[int, S, T, A, B, RETI, RW, UniDir, ERR] {
	return CombiTraversal[RETI, RW, ERR, int, S, T, A, B](
		func(ctx context.Context, source S) SeqIE[int, A] {
			return func(yield func(ValueIE[int, A]) bool) {
				i := 0
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					_, iv, focusErr := val.Get()
					cont := yield(ValIE(i, iv, JoinCtxErr(ctx, focusErr)))
					i++
					return cont
				})
			}
		},
		o.AsLengthGetter(),
		func(ctx context.Context, fmap func(focusIndex int, focus A) (B, error), source S) (T, error) {
			i := 0
			return o.AsModify()(ctx, func(focusIndex I, focus A) (B, error) {
				ret, err := fmap(i, focus)
				i++
				return ret, JoinCtxErr(ctx, err)
			}, source)
		},
		nil,
		IxMatchComparable[int](),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Indexing{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The ReIndexed combinator returns an optic where the index is replaced with an arbitrarily mapped value.
// The ixmap should ideally be an [Iso] to enable the [AsIndex] combinator to operate efficiently.
//
// Note: multiple values may share the same index.
//
// See [Indexing] for a version that replaces the index with an integer position.
func ReIndexed[I, J, K, S, T, A, B, RET, RW, DIR, ERR any, RETI TReturnOne, RWI any, DIRI any, ERRI any, ERRP TPure](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], ixmap Optic[K, I, I, J, J, RETI, RWI, DIRI, ERRI], ixMatch Predicate[lo.Tuple2[J, J], ERRP]) Optic[J, S, T, A, B, RET, RW, DIR, CompositionTree[ERR, ERRI]] {

	return Omni[J, S, T, A, B, RET, RW, DIR, CompositionTree[ERR, ERRI]](
		func(ctx context.Context, source S) (J, A, error) {
			i, ret, err := o.AsGetter()(ctx, source)
			j, ixErr := ixmap.AsOpGet()(ctx, i)
			return j, ret, errors.Join(err, ixErr)
		},
		o.AsSetter(),
		func(ctx context.Context, source S) SeqIE[J, A] {
			return func(yield func(ValueIE[J, A]) bool) {
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					index, focus, err := val.Get()
					j, ixErr := ixmap.AsOpGet()(ctx, index)
					err = errors.Join(err, ixErr)
					return yield(ValIE(j, focus, err))
				})
			}
		},
		o.AsLengthGetter(),
		func(ctx context.Context, fmap func(index J, focus A) (B, error), source S) (T, error) {
			return o.AsModify()(ctx, func(focusIndex I, focus A) (B, error) {
				res, ixErr := ixmap.AsOpGet()(ctx, focusIndex)
				if ixErr != nil {
					var b B
					return b, ixErr
				}
				ret, err := fmap(res, focus)
				return ret, err
			}, source)
		},
		func(ctx context.Context, getIndex J, source S) SeqIE[J, A] {
			if ixmap.OpticType()&expr.OpticTypeBiDirFlag != 0 {
				i, err := ixmap.AsReverseGetter()(ctx, getIndex)
				if err != nil {
					return func(yield func(ValueIE[J, A]) bool) {
						var a A
						yield(ValIE(getIndex, a, err))
					}
				}

				return func(yield func(ValueIE[J, A]) bool) {
					o.AsIxGetter()(ctx, i, source)(func(val ValueIE[I, A]) bool {
						_, focus, err := val.Get()
						return yield(ValIE(getIndex, focus, err))
					})
				}
			} else {
				return func(yield func(ValueIE[J, A]) bool) {

					o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
						index, focus, err := val.Get()
						j, ixErr := ixmap.AsOpGet()(ctx, index)
						err = errors.Join(err, ixErr)
						if err != nil {
							return yield(ValIE(j, focus, err))
						}

						match, err := PredGet(ctx, ixMatch, lo.T2(getIndex, j))
						if match || err != nil {
							return yield(ValIE(j, focus, err))
						}

						return true
					})
				}
			}
		},
		func(indexA, indexB J) bool {
			return Must(PredGet(context.Background(), ixMatch, lo.T2(indexA, indexB)))
		},
		o.AsReverseGetter(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.ReIndexed{
					OpticTypeExpr: ot,
					IxMap:         ixmap.AsExpr(),
					IxMatch:       ixMatch.AsExpr(),
					Optic:         o.AsExpr(),
				}
			},
			ixmap,
			ixMatch,
			o,
		),
	)
}

// The SelfIndex combinator replaces the index of the given optic with the source value.
//
// See [ReIndexed] for a version that replaces the index with an arbitrary value.
func SelfIndex[I, S, T, A, B, RET, RW, DIR, ERR any, ERRP TPure](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], ixMatch Predicate[lo.Tuple2[A, A], ERRP]) Optic[A, S, T, A, B, RET, RW, DIR, ERR] {

	return Omni[A, S, T, A, B, RET, RW, DIR, ERR](
		func(ctx context.Context, source S) (A, A, error) {
			_, a, err := o.AsGetter()(ctx, source)
			return a, a, err
		},
		func(ctx context.Context, focus B, source S) (T, error) {
			ret, err := o.AsSetter()(ctx, focus, source)
			return ret, err
		},
		func(ctx context.Context, source S) SeqIE[A, A] {
			return func(yield func(ValueIE[A, A]) bool) {
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					_, focus, err := val.Get()
					return yield(ValIE(focus, focus, err))
				})
			}
		},
		o.AsLengthGetter(),
		func(ctx context.Context, fmap func(index A, focus A) (B, error), source S) (T, error) {
			return o.AsModify()(ctx, func(index I, focus A) (B, error) {
				ret, err := fmap(focus, focus)
				return ret, err
			}, source)
		},
		func(ctx context.Context, index A, source S) SeqIE[A, A] {
			return func(yield func(ValueIE[A, A]) bool) {
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					_, focus, err := val.Get()
					if err != nil {
						return yield(ValIE(focus, focus, err))
					}

					match, err := PredGet(ctx, ixMatch, lo.T2(index, focus))
					if err != nil || match {
						return yield(ValIE(focus, focus, err))
					}
					return true
				})
			}
		},
		func(indexA, indexB A) bool {
			return Must(PredGet(context.Background(), ixMatch, lo.T2(indexA, indexB)))
		},
		o.AsReverseGetter(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.SelfIndex{
					OpticTypeExpr: ot,
					IxMatch:       ixMatch.AsExpr(),
					Optic:         o.AsExpr(),
				}
			},
			o,
			ixMatch,
		),
	)
}

// The OrderedI combinator focuses on the elements of the given optic in the order defined by the orderBy predicate
//
// Note: under modification the elements are focused in order but the result retains the original order.
//
// See:
//   - [OrderedColI] for a operation that re-orders a collection.
//   - [Ordered] for a non index aware version.
func OrderedI[I, S, T, A, B any, RET, RW, DIR, ERR, PERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], orderBy OrderByPredicateI[I, A, PERR]) Optic[I, S, T, A, B, RET, RW, UniDir, CompositionTree[ERR, PERR]] {

	orderByFnc := orderBy.AsOpGet()

	return CombiTraversal[RET, RW, CompositionTree[ERR, PERR], I, S, T, A, B](
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(ValueIE[I, A]) bool) {

				i := o.AsIter()(ctx, source)
				seq, _, _, err := heapSort(ctx, i, func(ctx context.Context, a, b ValueI[I, A]) (bool, error) {
					ret, err := orderByFnc(ctx, lo.T2(a, b))
					if errors.Is(err, ErrEmptyGet) {
						return false, nil
					}
					return ret, err
				})
				if err != nil {
					var i I
					var a A
					yield(ValIE(i, a, err))
					return
				}
				seq(yield)
			}
		},
		o.AsLengthGetter(),
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {

			si := o.AsIter()(ctx, source)
			sortedSeq, size, heap, err := heapSort(ctx, si, func(ctx context.Context, a, b ValueI[I, A]) (bool, error) {
				ret, err := orderByFnc(ctx, lo.T2(a, b))
				if errors.Is(err, ErrEmptyGet) {
					return false, nil
				}
				return ret, err
			})
			if err != nil {
				var t T
				return t, err
			}

			origSequence := make([]ValueIE[I, B], size)

			i := size - 1
			sortedSeq(func(val ValueIE[I, A]) bool {
				index, focus, err := val.Get()
				if err != nil {
					var b B
					origSequence[heap[i].A] = ValIE(index, b, err)
				}

				b, err := fmap(index, focus)
				origSequence[heap[i].A] = ValIE(index, b, err)

				i--
				return true
			})

			i = 0
			return o.AsModify()(ctx, func(index I, focus A) (B, error) {
				ret, err := origSequence[i].Value(), origSequence[i].Error()
				i++
				return ret, err
			}, source)

		},
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			return func(yield func(ValueIE[I, A]) bool) {
				i := o.AsIxGetter()(ctx, index, source)
				seq, _, _, err := heapSort(ctx, i, func(ctx context.Context, a, b ValueI[I, A]) (bool, error) {
					ret, err := orderByFnc(ctx, lo.T2(a, b))
					if errors.Is(err, ErrEmptyGet) {
						return false, nil
					}
					return ret, err
				})
				if err != nil {
					var i I
					var a A
					yield(ValIE(i, a, err))
					return
				}
				seq(yield)
			}
		},
		o.AsIxMatch(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Ordered{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
					OrderBy:       orderBy.AsExpr(),
				}
			},
			o,
			orderBy,
		),
	)
}

// The FirstOrDefaultI combinator returns an optic that focuses the first focused element or the default value and index if no elements are focused.
// See:
//   - [FirstOrDefault] for a non index aware version
func FirstOrDefaultI[I, S, T, A, RET, RW, DIR, ERR any](o Optic[I, S, T, A, A, RET, RW, DIR, ERR], defaultIndex I, defaultVal A) Optic[I, S, S, A, A, ReturnOne, ReadOnly, UniDir, ERR] {
	return UnsafeReconstrain[ReturnOne, ReadOnly, UniDir, ERR](Taking(
		Coalesce(
			Polymorphic[S, A](o),
			//TODO: I don't capture the IxMatch Predicate in the ConstIP so why not just use an ixmatch func instead of the predicate. I ran into a similar situation somewhere ele in the code.
			//I created an ixmatcht2 branch with the reverse of this. I tries to always use a Predicate. I think it can work but got messy.
			ConstI[S](defaultIndex, defaultVal, Op(func(source lo.Tuple2[I, I]) bool {
				return o.AsIxMatch()(source.A, source.B)
			})),
		),
		1,
	))
}

func LastOrDefaultI[I, S, T, A, RET, RW, DIR, ERR any](o Optic[I, S, T, A, A, RET, RW, DIR, ERR], defaultIndex I, defaultVal A) Optic[I, S, S, A, A, ReturnOne, ReadOnly, UniDir, ERR] {
	return UnsafeReconstrain[ReturnOne, ReadOnly, UniDir, ERR](
		Coalesce(
			Polymorphic[S, A](Last(o)),
			//TODO: I don't capture the IxMatch Predicate in the ConstIP so why not just use an ixmatch func instead of the predicate. I ran into a similar situation somewhere ele in the code.
			//I created an ixmatcht2 branch with the reverse of this. I tries to always use a Predicate. I think it can work but got messy.
			ConstI[S](defaultIndex, defaultVal, Op(func(source lo.Tuple2[I, I]) bool {
				return o.AsIxMatch()(source.A, source.B)
			})),
		),
	)
}

// The ForEachI combinator applies op to each element focused by forEach, focusing on the [Collection] of results using a mapped index
//
// This can be useful to apply impure modifications to a [ReadOnly] forEach optic.
// If forEach is [ReadWrite] then [Compose] should be used instead.
// See:
//   - ForEach for a non index mapping version.
func ForEachI[I, J, K, S, T, A, B, C, D any, RETI TReturnOne, RET, RETM any, RWI, RW any, RWM any, DIRI, DIR, DIRM any, ERRI TPure, SERR any](ixmap IxMapper[I, J, K, RETI, RWI, DIRI, ERRI], forEach Optic[I, S, T, A, B, RET, RW, DIR, SERR], op Optic[J, A, B, C, D, RETM, RWM, DIRM, SERR]) Optic[K, S, Collection[I, B, SERR], C, D, CompositionTree[RET, RETM], RWM, UniDir, SERR] {

	unmapIx := func(ctx context.Context, index K) (I, bool, J, bool, error) {
		var leftIndex I
		var leftOk = false
		var rightIndex J
		var rightOk = false
		if ixmap.OpticType()&expr.OpticTypeBiDirFlag != 0 {
			mapped, err := ixmap.AsReverseGetter()(ctx, index)
			if err != nil {
				return leftIndex, false, rightIndex, false, err
			}
			leftIndex, leftOk = mapped.A.Get()
			rightIndex, rightOk = mapped.B.Get()
		}
		return leftIndex, leftOk, rightIndex, rightOk, nil
	}

	mapIx := func(ctx context.Context, left I, right J) (K, error) {
		_, k, err := ixmap.AsGetter()(ctx, lo.T2(mo.Some(left), mo.Some(right)))
		return k, err
	}

	composedNestedOptic := ComposeLeft(forEach, op)

	return CombiTraversal[CompositionTree[RET, RETM], RWM, SERR, K, S, Collection[I, B, SERR], C, D](
		func(ctx context.Context, source S) SeqIE[K, C] {
			return func(yield func(ValueIE[K, C]) bool) {
				forEach.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					indexI, focusA, err := val.Get()
					if err != nil {
						var k K
						var c C
						return yield(ValIE(k, c, err))
					}

					cont := true

					op.AsIter()(ctx, focusA)(func(val ValueIE[J, C]) bool {
						indexJ, focusC, err := val.Get()
						indexK, errIx := mapIx(ctx, indexI, indexJ)
						cont = yield(ValIE(indexK, focusC, errors.Join(err, errIx)))
						return cont

					})
					return cont
				})
			}
		},
		composedNestedOptic.AsLengthGetter(),
		func(ctx context.Context, fmap func(index K, focus C) (D, error), source S) (Collection[I, B, SERR], error) {

			var modified []ValueIE[I, B]
			var retErr error

			forEach.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
				indexI, focus, err := val.Get()
				if err != nil {
					retErr = err
					return false
				}

				ret, err := op.AsModify()(ctx, func(indexJ J, modifyFocus C) (D, error) {
					indexK, err := mapIx(ctx, indexI, indexJ)
					if err != nil {
						var d D
						return d, err
					}
					mapped, err := fmap(indexK, modifyFocus)
					return mapped, err
				}, focus)

				if err != nil {
					retErr = err
					return false
				}

				modified = append(modified, ValIE(indexI, ret, nil))
				return true
			})

			if retErr != nil {
				return ValColIE[I, B, SERR](forEach.AsIxMatch()), retErr
			}

			return ValColIE[I, B, SERR](forEach.AsIxMatch(), modified...), nil
		},
		func(ctx context.Context, index K, source S) SeqIE[K, C] {
			return func(yield func(ValueIE[K, C]) bool) {
				indexI, iOk, indexJ, jOk, err := unmapIx(ctx, index)
				if err != nil {
					var c C
					yield(ValIE(index, c, err))
					return
				}

				if iOk {
					forEach.AsIxGetter()(ctx, indexI, source)(func(val ValueIE[I, A]) bool {
						indexI, focusA, err := val.Get()
						if err != nil {
							var k K
							var c C
							return yield(ValIE(k, c, err))
						}

						cont := true

						if jOk {
							op.AsIxGetter()(ctx, indexJ, focusA)(func(val ValueIE[J, C]) bool {
								_, focusC, err := val.Get()
								cont = yield(ValIE(index, focusC, err))
								return cont
							})
						} else {
							op.AsIter()(ctx, focusA)(func(val ValueIE[J, C]) bool {
								indexJ, focusC, err := val.Get()
								k, err := mapIx(ctx, indexI, indexJ)
								if err != nil {
									var c C
									cont = yield(ValIE(index, c, err))
									return cont
								}

								match := ixmap.AsIxMatch()(index, k)
								if match {
									cont = yield(ValIE(index, focusC, nil))
								}
								return cont
							})
						}

						return cont
					})
				} else {
					forEach.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
						indexI, focusA, err := val.Get()
						if err != nil {
							var k K
							var c C
							return yield(ValIE(k, c, err))
						}

						cont := true

						if jOk {
							op.AsIxGetter()(ctx, indexJ, focusA)(func(val ValueIE[J, C]) bool {
								indexJ, focusC, err := val.Get()
								k, err := mapIx(ctx, indexI, indexJ)
								if err != nil {
									var c C
									cont = yield(ValIE(index, c, err))
									return cont
								}

								match := ixmap.AsIxMatch()(index, k)
								if match {
									cont = yield(ValIE(index, focusC, nil))
								}
								return cont
							})
						} else {
							op.AsIter()(ctx, focusA)(func(val ValueIE[J, C]) bool {
								indexJ, focusC, err := val.Get()
								k, err := mapIx(ctx, indexI, indexJ)
								if err != nil {
									var c C
									cont = yield(ValIE(index, c, err))
									return cont
								}

								match := ixmap.AsIxMatch()(index, k)
								if match {
									cont = yield(ValIE(index, focusC, nil))
								}
								return cont
							})
						}

						return cont
					})
				}

			}

		},
		ixmap.AsIxMatch(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.ForEach{
					OpticTypeExpr: ot,
					ForEach:       forEach.AsExpr(),
					Op:            op.AsExpr(),
				}
			},
			forEach,
			op,
		),
	)
}

type EditType expr.EditType

const (
	EditInsert     = EditType(expr.EditInsert)
	EditDelete     = EditType(expr.EditDelete)
	EditSubstitute = EditType(expr.EditSubstitute)
	EditTranspose  = EditType(expr.EditTranspose)

	EditLevenshtein = EditType(EditInsert | EditDelete | EditSubstitute)
	EditOSA         = EditType(EditLevenshtein | EditTranspose)
	EditLCS         = EditType(EditInsert | EditDelete)

	EditAny = EditType(0xFF)
)

func EditDistanceI[I, S, T, A, B any, RET any, RW any, DIR any, ERR any, ERRI any, ERRP any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], allowedEdits EditType, ixMatch Predicate[lo.Tuple2[I, I], ERRI], equal Predicate[lo.Tuple2[A, A], ERRP], size int) Optic[Void, lo.Tuple2[S, S], lo.Tuple2[S, S], float64, float64, ReturnOne, ReadOnly, UniDir, CompositionTree[ERR, CompositionTree[ERRI, ERRP]]] {
	return CombiGetter[CompositionTree[ERR, CompositionTree[ERRI, ERRP]], Void, lo.Tuple2[S, S], lo.Tuple2[S, S], float64, float64](
		func(ctx context.Context, source lo.Tuple2[S, S]) (Void, float64, error) {

			a, err := GetContext(ctx, SliceOfP(WithIndex(o), size), source.A)
			if err != nil {
				return Void{}, 0, err
			}

			b, err := GetContext(ctx, SliceOfP(WithIndex(o), size), source.B)
			if err != nil {
				return Void{}, 0, err
			}

			distance := make([][]int, len(a)+1)
			for i := range distance {
				distance[i] = make([]int, len(b)+1)
			}

			for i := 0; i <= len(a); i++ {
				distance[i][0] = i
			}
			for j := 0; j <= len(b); j++ {
				distance[0][j] = j
			}

			for i := 1; i <= len(a); i++ {
				for j := 1; j <= len(b); j++ {
					cost := 0
					if (allowedEdits & EditSubstitute) == EditSubstitute {

						eq, err := PredGet(ctx, equal, lo.T2(a[i-1].value, b[j-1].value))
						if err != nil {
							return Void{}, 0, err
						}

						if !eq {
							cost = 1
						} else {
							eq, err := PredGet(ctx, ixMatch, lo.T2(a[i-1].index, b[j-1].index))
							if err != nil {
								return Void{}, 0, err
							}

							if !eq {
								cost = 1
							}
						}
					}
					cost = distance[i-1][j-1] + cost

					if (allowedEdits & EditDelete) == EditDelete {
						cost = min(cost, distance[i-1][j]+1)
					}

					if (allowedEdits & EditInsert) == EditInsert {
						cost = min(cost, distance[i][j-1]+1)
					}

					if (allowedEdits&EditTranspose) == EditTranspose && i > 1 && j > 1 {

						t1, err := PredGet(ctx, equal, lo.T2(a[i-1].value, b[j-2].value))
						if err != nil {
							return Void{}, 0, err
						}

						t2, err := PredGet(ctx, equal, lo.T2(a[i-2].value, b[j-1].value))
						if err != nil {
							return Void{}, 0, err
						}

						t1i, err := PredGet(ctx, ixMatch, lo.T2(a[i-1].index, b[j-2].index))
						if err != nil {
							return Void{}, 0, err
						}

						t2i, err := PredGet(ctx, ixMatch, lo.T2(a[i-2].index, b[j-1].index))
						if err != nil {
							return Void{}, 0, err
						}

						if t1 && t1i && t2 && t2i {
							cost = min(cost, distance[i-2][j-2]+1)
						}
					}

					distance[i][j] = cost
				}
			}

			return Void{}, float64(distance[len(a)][len(b)]), nil
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.EditDistance{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
					AllowedEdits:  expr.EditType(allowedEdits),
				}
			},
			o,
			ixMatch,
			equal,
		),
	)
}
