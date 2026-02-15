package optic

import (
	"context"
	"errors"
	"reflect"
	"slices"

	"github.com/samber/lo"
	"github.com/spearson78/go-optic/expr"
)

func reverse[I, A any](ctx context.Context, source SeqIE[I, A]) SeqIE[I, A] {
	return func(yield func(val ValueIE[I, A]) bool) {
		var ret []ValueIE[I, A]
		source(func(val ValueIE[I, A]) bool {
			ret = append(ret, val)
			return true
		})
		slices.Reverse(ret)
		for _, val := range ret {
			if !yield(val) {
				break
			}
		}
	}
}

// ReversedCol returns an [Iso] that reverses the order of a [Collection].
//
// See:
// - [ReversedColP] for a polymorphic version
// - [Reversed] for a non collection version.
func ReversedCol[I, A any]() Optic[Void, Collection[I, A, Pure], Collection[I, A, Pure], Collection[I, A, Pure], Collection[I, A, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return ReversedColP[I, A, A]()
}

// ReversedColP returns an polymorphic [Iso] that reverses the order of an [Collection].
//
// See:
// - [ReversedCol] for a non polymorphic version
// - [Reversed] for a non collection version.
func ReversedColP[I, A, B any]() Optic[Void, Collection[I, A, Pure], Collection[I, B, Pure], Collection[I, A, Pure], Collection[I, B, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return CombiIso[ReadWrite, BiDir, Pure, Collection[I, A, Pure], Collection[I, B, Pure], Collection[I, A, Pure], Collection[I, B, Pure]](
		func(ctx context.Context, source Collection[I, A, Pure]) (Collection[I, A, Pure], error) {
			return ColIE[Pure](
				func(ctx context.Context) SeqIE[I, A] {
					return reverse(ctx, source.AsIter()(ctx))
				},
				func(ctx context.Context, index I) SeqIE[I, A] {
					return reverse(ctx, source.AsIxGet()(ctx, index))
				},
				source.AsIxMatch(),
				source.AsLengthGetter(),
			), nil
		},
		func(ctx context.Context, focus Collection[I, B, Pure]) (Collection[I, B, Pure], error) {
			return ColIE[Pure](
				func(ctx context.Context) SeqIE[I, B] {
					return reverse(ctx, focus.AsIter()(ctx))
				},
				func(ctx context.Context, index I) SeqIE[I, B] {
					return reverse(ctx, focus.AsIxGet()(ctx, index))
				},
				focus.AsIxMatch(),
				focus.AsLengthGetter(),
			), nil
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Reverse{
				OpticTypeExpr: ot,
				I:             reflect.TypeFor[I](),
				A:             reflect.TypeFor[A](),
				B:             reflect.TypeFor[B](),
			}
		}),
	)
}

// FilteredCol returns a [Lens] that removes elements from a [Collection] that match a predicate
//
// See:
// - [FilteredColI] for an index aware version
// - [Filtered] for a non collection version.
func FilteredCol[I comparable, A any, ERR any](pred Predicate[A, ERR]) Optic[Void, Collection[I, A, ERR], Collection[I, A, ERR], Collection[I, A, ERR], Collection[I, A, ERR], ReturnOne, ReadWrite, UniDir, ERR] {
	return FilteredColI(PredToOpI[I](pred), IxMatchComparable[I]())
}

// FilteredColI returns an index aware [Lens] that removes elements from an [Collection] that match a predicate
//
// See:
// - [FilteredCol] for a non index aware version
// - [FilteredI] for a non collection version
func FilteredColI[I, A any, ERR any](pred PredicateI[I, A, ERR], ixMatch func(a, b I) bool) Optic[Void, Collection[I, A, ERR], Collection[I, A, ERR], Collection[I, A, ERR], Collection[I, A, ERR], ReturnOne, ReadWrite, UniDir, ERR] {
	return ColOf(EErrMerge(FilteredI(TraverseColIE[I, A, ERR](ixMatch), pred)))
}

// Append returns a [Lens] that appends additional integer indexed elements to the end of the source [Collection].
//
// See:
//   - [AppendP] for a polymorphic version
//   - [AppendCol] for a version that takes an arbitrary [Collection] as a parameter
//   - [PrependCol] for a prepend version
func Append[A any](toAppend ...A) Optic[Void, Collection[int, A, Pure], Collection[int, A, Pure], Collection[int, A, Pure], Collection[int, A, Pure], ReturnOne, ReadOnly, UniDir, Pure] {
	return AppendCol(ValCol(toAppend...))
}

// AppendCol returns a [Lens] that appends additional elements to the end of the source [Collection].
//
// See:
//   - [AppendColP] for a polymorphic version
//   - [PrependCol] for a prepend version
func AppendCol[I, A, ERR any](toAppend Collection[I, A, ERR]) Optic[Void, Collection[I, A, ERR], Collection[I, A, ERR], Collection[I, A, ERR], Collection[I, A, ERR], ReturnOne, ReadOnly, UniDir, ERR] {
	return Ret1(Ro(Ud(EErrR(EErrMergeL(Compose(
		T2Of(
			Identity[Collection[I, A, ERR]](),
			Const[Collection[I, A, ERR]](toAppend),
		),
		AppendColT2[I, A, ERR](),
	))))))
}

func AppendColT2[I, A, ERR any]() Optic[Void, lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], Collection[I, A, ERR], Collection[I, A, ERR], ReturnOne, ReadOnly, UniDir, ERR] {
	return CombiGetter[ERR, Void, lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], Collection[I, A, ERR], Collection[I, A, ERR]](
		func(ctx context.Context, source lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]]) (Void, Collection[I, A, ERR], error) {
			ret := ColIE[ERR, I, A](
				func(ctx context.Context) SeqIE[I, A] {
					return func(yield func(val ValueIE[I, A]) bool) {
						cont := true
						source.A.AsIter()(ctx)(func(val ValueIE[I, A]) bool {
							index, focus, err := val.Get()
							if !cont {
								panic(yieldAfterBreak)
							}
							cont = yield(ValIE(index, focus, JoinCtxErr(ctx, err)))
							return cont
						})

						if cont {
							source.B.AsIter()(ctx)(func(val ValueIE[I, A]) bool {
								focusIndex, focus, err := val.Get()
								if !cont {
									panic(yieldAfterBreak)
								}
								cont = yield(ValIE(focusIndex, focus, JoinCtxErr(ctx, err)))
								return cont
							})
						}
					}
				},
				func(ctx context.Context, index I) SeqIE[I, A] {
					return func(yield func(val ValueIE[I, A]) bool) {
						cont := true
						source.A.AsIxGet()(ctx, index)(func(val ValueIE[I, A]) bool {
							index, focus, err := val.Get()
							if !cont {
								panic(yieldAfterBreak)
							}
							cont = yield(ValIE(index, focus, JoinCtxErr(ctx, err)))
							return cont
						})

						if cont {
							source.B.AsIxGet()(ctx, index)(func(val ValueIE[I, A]) bool {
								focusIndex, focus, err := val.Get()
								if !cont {
									panic(yieldAfterBreak)
								}
								cont = yield(ValIE(focusIndex, focus, JoinCtxErr(ctx, err)))
								return cont
							})
						}
					}
				},
				func(a, b I) bool {
					return source.A.AsIxMatch()(a, b) || source.B.AsIxMatch()(a, b)
				},
				func(ctx context.Context) (int, error) {
					l, err := source.A.AsLengthGetter()(ctx)
					if err != nil {
						return 0, err
					}
					l2, err := source.B.AsLengthGetter()(ctx)
					if err != nil {
						return 0, err
					}

					return l + l2, nil
				},
			)

			return Void{}, ret, nil
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.AppendCol{
					OpticTypeExpr: ot,
					I:             reflect.TypeFor[I](),
					A:             reflect.TypeFor[A](),
				}
			},
		),
	)
}

// PrependCol returns a [Lens] that prepends additional elements to the start of the source [Collection].
//
// See:
//   - [PrependColP] for a polymorphic version
//   - [AppendCol] for an append version
//   - [Prepending] for a non collection version
func PrependCol[I, A, ERR any](toPrepend Collection[I, A, ERR]) Optic[Void, Collection[I, A, ERR], Collection[I, A, ERR], Collection[I, A, ERR], Collection[I, A, ERR], ReturnOne, ReadOnly, UniDir, ERR] {
	return Ret1(Ro(Ud(EErrR(EErrMergeL(Compose(
		T2Of(
			Const[Collection[I, A, ERR]](toPrepend),
			Identity[Collection[I, A, ERR]](),
		),
		AppendColT2[I, A, ERR](),
	))))))
}

func PrependColT2[I, A, ERR any]() Optic[Void, lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], Collection[I, A, ERR], Collection[I, A, ERR], ReturnOne, ReadOnly, UniDir, ERR] {
	return Ret1(Ro(Ud(EErrR(Compose(
		SwappedT2[Collection[I, A, ERR], Collection[I, A, ERR]](),
		AppendColT2[I, A, ERR](),
	)))))
}

// Prepend returns a [Lens] that appends additional integer indexed elements to the start of the source [Collection].
//
// See:
//   - [PrependP] for a polymorphic version
//   - [PrependCol] for a version that takes an arbitrary [Collection] as a parameter
//   - [AppendCol] for an append version
func Prepend[A any](toAppend ...A) Optic[Void, Collection[int, A, Pure], Collection[int, A, Pure], Collection[int, A, Pure], Collection[int, A, Pure], ReturnOne, ReadOnly, UniDir, Pure] {
	return PrependCol(ValCol(toAppend...))
}

// OrderedCol returns a collection operation that sorts an [Collection] according to the given [OrderByPredicate]
//
// See:
//   - [SortP] for a polymorphic version
//   - [OrderedColI] for an index aware version
//   - [Ordered] for a non collection operation version.
func OrderedCol[I comparable, A any, ERR any](orderBy OrderByPredicate[A, ERR]) Optic[Void, Collection[I, A, ERR], Collection[I, A, ERR], Collection[I, A, ERR], Collection[I, A, ERR], ReturnOne, ReadWrite, UniDir, ERR] {
	return OrderedColI(PredT2ToOpT2I[I, I](orderBy), IxMatchComparable[I]())
}

// OrderedColI returns an index aware collection operation that sorts an [Collection] according to the given [OrderByPredicateI]
//
// See:
//   - [OrderedCol] for a non index aware version
//   - [OrderedI] for a non collection operation version.
func OrderedColI[I, A any, ERR any](orderBy OrderByPredicateI[I, A, ERR], ixMatch func(a, b I) bool) Optic[Void, Collection[I, A, ERR], Collection[I, A, ERR], Collection[I, A, ERR], Collection[I, A, ERR], ReturnOne, ReadWrite, UniDir, ERR] {
	return ColOf(EErrMerge(OrderedI(TraverseColIE[I, A, ERR](ixMatch), orderBy)))
}

// SubCol returns a collection operation that removes elements from the start and end of a [Collection]].
//
// length may be negative in which case it will remove length elements from the end of the [Collection]
//
// See:
// - [SubColP] for a non polymorphic version
// - [FilteredCol] for a version that can remove arbitrary elements
func SubCol[I, A any](start int, length int) Optic[Void, Collection[I, A, Pure], Collection[I, A, Pure], Collection[I, A, Pure], Collection[I, A, Pure], ReturnOne, ReadWrite, UniDir, Pure] {

	expr := ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
		return expr.SubCol{
			OpticTypeExpr: ot,
			Start:         start,
			Length:        length,
		}
	})

	//Avoid corner cases in the more complex combinations

	if start >= 0 {
		if length >= 0 {

			end := (start + length) - 1

			return CombiLens[ReadWrite, Pure, Void, Collection[I, A, Pure], Collection[I, A, Pure], Collection[I, A, Pure], Collection[I, A, Pure]](
				func(ctx context.Context, source Collection[I, A, Pure]) (Void, Collection[I, A, Pure], error) {
					return Void{}, ColIE[Pure](
						func(ctx context.Context) SeqIE[I, A] {

							return func(yield func(val ValueIE[I, A]) bool) {
								i := 0
								fnc := func(val ValueIE[I, A]) bool {
									err := val.Error()
									if err != nil {
										return yield(val)
									}

									isInRange := i >= start && i <= end
									i++
									if isInRange {

										return yield(val)
									}

									return true
								}
								source.AsIter()(ctx)(fnc)
							}
						},
						nil,
						source.AsIxMatch(),
						nil,
					), nil
				},
				func(ctx context.Context, focus, source Collection[I, A, Pure]) (Collection[I, A, Pure], error) {
					return ColIE[Pure](
						func(ctx context.Context) SeqIE[I, A] {
							spliced := false

							return func(yield func(val ValueIE[I, A]) bool) {
								i := 0
								cont := true

								fnc := func(val ValueIE[I, A]) bool {
									err := val.Error()
									if err != nil {
										return yield(val)
									}

									isInRange := i >= start && i <= end
									i++
									if isInRange {
										if !spliced {
											spliced = true
											focus.AsIter()(ctx)(func(val ValueIE[I, A]) bool {
												cont = yield(val)
												return cont
											})
											return cont
										} else {
											//Values were replaced with those spliced in above
											return true
										}
									} else {
										cont = yield(val)
										return cont
									}
								}
								source.AsIter()(ctx)(fnc)

							}
						},
						nil,
						source.AsIxMatch(),
						nil,
					), nil
				},
				IxMatchVoid(),
				expr,
			)

		} else {

			return CombiLens[ReadWrite, Pure, Void, Collection[I, A, Pure], Collection[I, A, Pure], Collection[I, A, Pure], Collection[I, A, Pure]](
				func(ctx context.Context, source Collection[I, A, Pure]) (Void, Collection[I, A, Pure], error) {
					return Void{}, ColIE[Pure](
						func(ctx context.Context) SeqIE[I, A] {
							return func(yield func(ValueIE[I, A]) bool) {

								buf := make([]ValueIE[I, A], -length)
								i := 0
								source.AsIter()(ctx)(func(val ValueIE[I, A]) bool {
									err := val.Error()
									if err != nil {
										return yield(val)
									}

									if i < start {
										i++
										return true
									}

									if i >= start+(-length) {
										//Length is negative
										val := buf[(i+length)%-length]
										if !yield(val) {
											return false
										}
									}

									buf[i%-length] = val

									i++

									return true
								})
							}
						},
						nil,
						source.AsIxMatch(),
						nil,
					), nil
				},
				func(ctx context.Context, focus, source Collection[I, A, Pure]) (Collection[I, A, Pure], error) {
					return ColIE[Pure](
						func(ctx context.Context) SeqIE[I, A] {
							cont := true
							splicePrimed := false
							splicePos := 0
							return func(yield func(val ValueIE[I, A]) bool) {

								buf := make([]ValueIE[I, A], -length)
								i := -1
								source.AsIter()(ctx)(func(val ValueIE[I, A]) bool {
									i++
									err := val.Error()
									if err != nil {
										cont = yield(val)
										return cont
									}

									yieldVal := buf[i%-length]
									buf[i%-length] = val

									if i < -length {
										//Still filling the buffer
										i++
										return true
									}

									//Buffer is filled
									if i >= start {
										if !splicePrimed {
											splicePrimed = true
											splicePos = i + length
										}

										if i == splicePos {
											//We don't need any more original values
											return false
										} else {
											//yield the delayed values
											cont = yield(yieldVal)
											return cont
										}
									}

									return true
								})

								if cont {
									focus.AsIter()(ctx)(func(val ValueIE[I, A]) bool {
										cont = yield(val)
										return cont
									})
								}

								//Now spool out the buffer
								spoolCount := min(-length, i)
								for o := range spoolCount {
									val := buf[(i-o)%-length]
									cont = yield(val)
									if !cont {
										break
									}
								}
							}

						},
						nil,
						source.AsIxMatch(),
						nil,
					), nil
				},
				IxMatchVoid(),
				expr,
			)
		}
	} else {

		if length < 0 {
			length = -start + length
			if length < 0 {
				length = 0
			}
		}

		if length > 0 {

			length := min(-start, length)

			return CombiLens[ReadWrite, Pure, Void, Collection[I, A, Pure], Collection[I, A, Pure], Collection[I, A, Pure], Collection[I, A, Pure]](
				func(ctx context.Context, source Collection[I, A, Pure]) (Void, Collection[I, A, Pure], error) {
					return Void{}, ColIE[Pure](
						func(ctx context.Context) SeqIE[I, A] {
							return func(yield func(val ValueIE[I, A]) bool) {
								found := 0
								buf := make([]ValueIE[I, A], length)
								i := 0
								source.AsIter()(ctx)(func(val ValueIE[I, A]) bool {
									err := val.Error()
									if err != nil {
										yield(val)
										//The error prevented us from reaching the end of the sequence
										//Yield the error and bail
										found = 0
										return false
									}

									buf[i%length] = val
									i++
									found++
									found = min(length, found)
									return true
								})

								for j := 0; j < found; j++ {
									val := buf[(i+j)%length]
									if !yield(val) {
										return
									}
								}
							}
						},
						nil,
						source.AsIxMatch(),
						nil,
					), nil
				},
				func(ctx context.Context, focus, source Collection[I, A, Pure]) (Collection[I, A, Pure], error) {
					return ColIE[Pure](
						func(ctx context.Context) SeqIE[I, A] {
							return func(yield func(val ValueIE[I, A]) bool) {
								cont := true
								found := 0
								buf := make([]ValueIE[I, A], length)
								i := -1
								source.AsIter()(ctx)(func(val ValueIE[I, A]) bool {
									i++
									found++
									found = min(length, found)

									err := val.Error()
									if err != nil {
										yield(val)
										//The error prevented us from reaching the end of the sequence
										//Yield the error and bail
										found = 0
										return false
									}

									yieldVal := buf[i%length]
									buf[i%length] = val

									if i > length {
										cont = yield(yieldVal)
										return cont
									}

									return true
								})

								if cont {
									for j := 0; j < found; j++ {
										val := buf[(i+j)%length]
										if !yield(val) {
											return
										}
									}
								}
							}
						},
						nil,
						source.AsIxMatch(),
						nil,
					), nil
				},
				IxMatchVoid(),
				expr,
			)

		} else {
			return CombiLens[ReadWrite, Pure, Void, Collection[I, A, Pure], Collection[I, A, Pure], Collection[I, A, Pure], Collection[I, A, Pure]](
				func(ctx context.Context, source Collection[I, A, Pure]) (Void, Collection[I, A, Pure], error) {
					return Void{}, ValColIE[I, A, Pure](source.AsIxMatch()), nil
				},
				func(ctx context.Context, focus, source Collection[I, A, Pure]) (Collection[I, A, Pure], error) {
					return focus, nil
				},
				IxMatchVoid(),
				expr,
			)
		}

	}
}

// The ReIndexedCol combinator focuses on a [Collection] where the comparable index has been remapped using the ixmap optic.
//
// See:
//   - [ReIndexedColP] for a polymorphic version
//   - [ReIndexedColI] for a version that supports a non comparable index
//   - [ReIndexedColIP] for a polymorphic version, non comparable index version.
func ReIndexedCol[A any, I, J comparable, L any, IRET TReturnOne, IRW any, IDIR any, IERR TPure](ixmap Optic[L, I, I, J, J, IRET, IRW, IDIR, IERR]) Optic[Void, Collection[I, A, Pure], Collection[I, A, Pure], Collection[J, A, Pure], Collection[J, A, Pure], ReturnOne, CompositionTree[IRW, IDIR], IDIR, Pure] {
	return ReIndexedColP[A, A](ixmap)
}

// The ReIndexedColI combinator focuses on a [Collection] where the index has been remapped using the ixmap optic.
//
// See:
//   - [ReIndexedCol] for a simpler version that supports only comparable indices
//   - [ReIndexedColP] for a polymorphic version that supports only comparable indices
//   - [ReIndexedColIP] for a polymorphic version, non comparable index version.
func ReIndexedColI[A, I, J, L any, IRET TReturnOne, IRW any, IDIR any, IERR, MIERR, MJERR TPure](ixmap Optic[L, I, I, J, J, IRET, IRW, IDIR, IERR], ixmatchi Predicate[lo.Tuple2[I, I], MIERR], ixmatchj Predicate[lo.Tuple2[J, J], MJERR]) Optic[Void, Collection[I, A, Pure], Collection[I, A, Pure], Collection[J, A, Pure], Collection[J, A, Pure], ReturnOne, CompositionTree[IRW, IDIR], IDIR, Pure] {
	return ReIndexedColIP[A, A](ixmap, ixmatchi, ixmatchj)
}

// The ReIndexedColP polymorphic combinator focuses on a [Collection] where the index has been remapped using the ixmap optic.
//
// See:
//   - [ReIndexedCol] for a non polymorphic version
//   - [ReIndexedColI] for a non polymorphic version that supports non comparable indices
//   - [ReIndexedColIP] for a polymorphic version, non comparable index version.
func ReIndexedColP[A, B any, I, J comparable, L any, IRET TReturnOne, IRW any, IDIR any, IERR TPure](ixmap Optic[L, I, I, J, J, IRET, IRW, IDIR, IERR]) Optic[Void, Collection[I, A, Pure], Collection[I, B, Pure], Collection[J, A, Pure], Collection[J, B, Pure], ReturnOne, CompositionTree[IRW, IDIR], IDIR, Pure] {
	return ReIndexedColIP[A, B](ixmap, EqT2[I](), EqT2[J]())
}

// The ReIndexedColIP polymorphic combinator focuses on a [Collection] where the index has been remapped using the ixmap optic.
//
// See:
//   - [ReIndexedCol] for simpler a non polymorphic version that supports only comparable indices
//   - [ReIndexedColI] for a non polymorphic version.
//   - [ReIndexedColP] for a comparable index version.
func ReIndexedColIP[A, B, I, J, L any, IRET TReturnOne, IRW any, IDIR any, IERR, MIERR, MJERR TPure](ixmap Optic[L, I, I, J, J, IRET, IRW, IDIR, IERR], ixmatchi Predicate[lo.Tuple2[I, I], MIERR], ixmatchj Predicate[lo.Tuple2[J, J], MJERR]) Optic[Void, Collection[I, A, Pure], Collection[I, B, Pure], Collection[J, A, Pure], Collection[J, B, Pure], ReturnOne, CompositionTree[IRW, IDIR], IDIR, Pure] {
	x := CombiIso[CompositionTree[IRW, IDIR], IDIR, Pure, Collection[I, A, Pure], Collection[I, B, Pure], Collection[J, A, Pure], Collection[J, B, Pure]](
		func(ctx context.Context, source Collection[I, A, Pure]) (Collection[J, A, Pure], error) {

			return ColIE[Pure](
				func(ctx context.Context) SeqIE[J, A] {
					return func(yield func(ValueIE[J, A]) bool) {
						source.AsIter()(ctx)(func(val ValueIE[I, A]) bool {
							index, focus, err := val.Get()
							j, mapErr := ixmap.AsOpGet()(ctx, index)
							return yield(ValIE(j, focus, errors.Join(err, mapErr)))
						})
					}
				},
				nil,
				PredToIxMatch(ixmatchj),
				source.AsLengthGetter(),
			), nil
		},
		func(ctx context.Context, focus Collection[J, B, Pure]) (Collection[I, B, Pure], error) {
			return ColIE[Pure](
				func(ctx context.Context) SeqIE[I, B] {
					return func(yield func(ValueIE[I, B]) bool) {
						focus.AsIter()(ctx)(func(val ValueIE[J, B]) bool {
							index, focus, err := val.Get()
							i, mapErr := ixmap.AsReverseGetter()(ctx, index)
							return yield(ValIE(i, focus, errors.Join(err, mapErr)))
						})
					}
				},
				nil,
				PredToIxMatch(ixmatchi),
				focus.AsLengthGetter(),
			), nil
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.ReIndexedCol{
					OpticTypeExpr: ot,
					IxMap:         ixmap.AsExpr(),
					IxMatchI:      ixmatchi.AsExpr(),
					IxMatchJ:      ixmatchj.AsExpr(),
				}
			},
			ixmap,
			ixmatchi,
			ixmatchj,
		),
	)

	return CombiDir[IDIR](x)
}
