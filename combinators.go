package optic

import (
	"cmp"
	"context"
	"errors"
	"slices"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spearson78/go-optic/expr"
)

// The Length combinator returns an optic that focuses the number of elements in the given optic.
//
// Warning: using a modify action on the return optic will return an error.
func Length[I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[Void, S, S, int, int, ReturnOne, ReadOnly, UniDir, ERR] {
	return CombiGetter[ERR, Void, S, S, int, int](
		func(ctx context.Context, source S) (Void, int, error) {
			ret, err := o.AsLengthGetter()(ctx, source)
			return Void{}, ret, err
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Length{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The Taking combinator returns an optic that focuses on the first amount elements in the given optic.
func Taking[I, S, T, A, RETI, RW, DIR, ERR any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR], amount int) Optic[I, S, T, A, A, ReturnMany, RW, UniDir, ERR] {
	return EErrL(filtered(o, Lt(amount), Const[ValueI[I, A]](true), FilterContinue, FilterStop))
}

// The Dropping combinator returns an optic that skips the first amount elements in the given optic and focuses on the remaining elements.
func Dropping[I, S, T, A, RETI, RW, DIR, ERR any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR], amount int) Optic[I, S, T, A, A, ReturnMany, RW, UniDir, ERR] {
	return EErrL(filtered(o, Gte(amount), Const[ValueI[I, A]](true), FilterYieldAll, FilterContinue))
}

// The Reversed combinator returns an optic that focuses on the elements in the input optic in reverse order.
//
// Warning: the input optic is fully iterated in order to reverse the order.
func Reversed[I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, T, A, B, RET, RW, UniDir, ERR] {
	return CombiTraversal[RET, RW, ERR, I, S, T, A, B](
		func(ctx context.Context, source S) SeqIE[I, A] {
			var stash []ValueIE[I, A]
			o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
				stash = append(stash, val)
				return ctx.Err() == nil
			})

			slices.Reverse(stash)

			return func(yield func(ValueIE[I, A]) bool) {
				for _, v := range stash {
					i, a, err := v.Get()
					if !yield(ValIE(i, a, JoinCtxErr(ctx, err))) {
						break
					}
				}
			}
		},
		o.AsLengthGetter(),
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {
			var err error
			var stash []ValueI[I, A]
			o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
				focusIndex, v, focusErr := val.Get()
				focusErr = JoinCtxErr(ctx, focusErr)
				if focusErr != nil {
					err = focusErr
					return false
				}
				stash = append(stash, ValI(focusIndex, v))
				return true
			})

			err = JoinCtxErr(ctx, err)
			if err != nil {
				var t T
				return t, err
			}

			slices.Reverse(stash)
			stashB := make([]B, len(stash))

			//Fill stashB in reverse so we visit in reverse but can return them modified in the original order
			i := len(stash) - 1
			for _, ixfa := range stash {
				a, err := fmap(ixfa.index, ixfa.value)
				err = JoinCtxErr(ctx, err)
				if err != nil {
					var t T
					return t, err
				}

				stashB[i] = a
				i--
			}

			stash = nil

			i = 0
			ret, err := o.AsModify()(ctx, func(index I, focus A) (B, error) {
				ret := stashB[i]
				i++
				return ret, ctx.Err()
			}, source)

			return ret, JoinCtxErr(ctx, err)
		},
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			var stash []ValueIE[I, A]
			o.AsIxGetter()(ctx, index, source)(func(val ValueIE[I, A]) bool {
				stash = append(stash, val)
				return ctx.Err() == nil
			})

			slices.Reverse(stash)

			return func(yield func(ValueIE[I, A]) bool) {
				for _, v := range stash {
					i, a, err := v.Get()
					if !yield(ValIE(i, a, JoinCtxErr(ctx, err))) {
						break
					}
				}
			}
		},
		o.AsIxMatch(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Reversed{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The TakingWhile combinator focuses on the elements of the given optic until the predicate returns false.
//
// See [TakingWhileI] for an indexed version.
func TakingWhile[I, S, T, A, RETI, RW any, DIR, ERR, ERRP any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR], pred Predicate[A, ERRP]) Optic[I, S, T, A, A, ReturnMany, RW, UniDir, CompositionTree[ERR, ERRP]] {
	return TakingWhileI(o, PredToOpI[I](pred))
}

// The DroppingWhile combinator skips the the elements of the given optic until the predicate returns false and focuses on the remaining elements.
//
// See [DroppingWhileI] for an indexed version.
func DroppingWhile[I, S, T, A, RETI, RW any, DIR any, ERR, ERRP any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR], pred Predicate[A, ERRP]) Optic[I, S, T, A, A, ReturnMany, RW, UniDir, CompositionTree[ERR, ERRP]] {
	return DroppingWhileI(o, PredToOpI[I](pred))
}

// The TrimmingWhile combinator skips leading and trailing elements until an element that does not match the predicate is found.
//
// See [TrimmingWhileI] for an indexed version.
func TrimmingWhile[I, S, T, A, RETI, RW any, DIR any, ERR, ERRP any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR], pred Predicate[A, ERRP]) Optic[I, S, T, A, A, ReturnMany, RW, UniDir, CompositionTree[ERR, ERRP]] {
	return TrimmingWhileI(o, PredToOpI[I](pred))
}

// The Filtered combinator focuses on the elements of the given optic that match the given predicate.
//
// See:
//   - [FilteredBy] for a filtering combinator that yields the source instead of the focus.
//   - [FilteredI] for an indexed version
func Filtered[I, S, T, A, RETI, RW any, DIR any, ERR, ERRP any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR], pred Predicate[A, ERRP]) Optic[I, S, T, A, A, ReturnMany, RW, UniDir, CompositionTree[ERR, ERRP]] {
	return FilteredI(o, PredToOpI[I](pred))
}

func FilteredP[I, S, T, A, B, RETI, RW any, DIR any, ERR, ERRP any](o Optic[I, S, T, A, B, RETI, RW, DIR, ERR], pred Predicate[A, ERRP]) Optic[I, S, T, A, A, ReturnMany, ReadOnly, UniDir, CompositionTree[ERR, ERRP]] {
	return FilteredIP(o, PredToOpI[I](pred))
}

// The Element combinator focuses on the element of the given optic at the given 0 based index.
//
// See :
//   - [AtMap] for an optic that focuses an indexed value in a map.
func Element[I, S, T, A, RETI, RW, DIR any, ERR any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR], index int) Optic[I, S, T, A, A, ReturnMany, RW, UniDir, ERR] {
	return EErrL(filtered(o, Eq(index), Const[ValueI[I, A]](true), FilterStop, FilterContinue))
}

var ErrDuplicateKey = errors.New("duplicate key")

type nakedOptic[I, S, T, A, B any] interface {
	// Provides this optics [Traversal] style modification function.
	// This methods behavior is only defined for [ReadWrite] optics.
	AsModify() ModifyFunc[I, S, T, A, B]

	// Provides this optics [Iteration] style view function.
	// This methods behavior is defined for all optic types
	AsIter() IterFunc[I, S, A]

	// Provides this optics [Operation] style view function.
	// This methods behavior is only defined for [ReturnOne] optics.
	AsGetter() GetterFunc[I, S, A]

	// Provides this optics [Lens] style modification function.
	// This methods behavior is only defined for [ReadWrite] optics.
	AsSetter() SetterFunc[S, T, B]

	// Provides this optics [Traversal] style view by index function.
	// This methods behavior is only defined for all optics.
	AsIxGetter() IxGetterFunc[I, S, A]

	AsIxMatch() IxMatchFunc[I]

	// Provides this optics [Iso] or [Prism] style modification function.
	// This methods behavior is only defined for all [BiDir]
	AsReverseGetter() ReverseGetterFunc[T, B]

	// Provides this optics [Iteration] style length getter function.
	// This methods behavior is defined for all optic types
	AsLengthGetter() LengthGetterFunc[S]

	AsOpGet() OpGetFunc[S, A]

	// Provides this optics expression tree representation.
	AsExpr() expr.OpticExpression

	//Provides the custom expression handler. Used by non go backends to execute expressions
	AsExprHandler() func(ctx context.Context) (ExprHandler, error)

	// Provides this optics internal type. This indicates which methods have efficient implementations.
	OpticType() expr.OpticType
}

func Coalesce[I, S, T, A, B, RET1, RW1, DIR1, ERR1, RET2, RW2, DIR2, ERR2 any](o1 Optic[I, S, T, A, B, RET1, RW1, DIR1, ERR1], o2 Optic[I, S, T, A, B, RET2, RW2, DIR2, ERR2]) Optic[I, S, T, A, B, CompositionTree[RET1, RET2], CompositionTree[RW1, RW2], UniDir, CompositionTree[ERR1, ERR2]] {
	return coalesceN[I, S, T, A, B, CompositionTree[RET1, RET2], CompositionTree[RW1, RW2], CompositionTree[ERR1, ERR2]](o1, o2)
}

func Coalesce3[I, S, T, A, B, RET1, RW1, DIR1, ERR1, RET2, RW2, DIR2, ERR2, RET3, RW3, DIR3, ERR3 any](o1 Optic[I, S, T, A, B, RET1, RW1, DIR1, ERR1], o2 Optic[I, S, T, A, B, RET2, RW2, DIR2, ERR2], o3 Optic[I, S, T, A, B, RET3, RW3, DIR3, ERR3]) Optic[I, S, T, A, B, CompositionTree[CompositionTree[RET1, RET2], RET3], CompositionTree[CompositionTree[RW1, RW2], RW3], UniDir, CompositionTree[CompositionTree[ERR1, ERR2], ERR3]] {
	return coalesceN[I, S, T, A, B, CompositionTree[CompositionTree[RET1, RET2], RET3], CompositionTree[CompositionTree[RW1, RW2], RW3], CompositionTree[CompositionTree[ERR1, ERR2], ERR3]](o1, o2, o3)
}

func Coalesce4[I, S, T, A, B, RET1, RW1, DIR1, ERR1, RET2, RW2, DIR2, ERR2, RET3, RW3, DIR3, ERR3, RET4, RW4, DIR4, ERR4 any](o1 Optic[I, S, T, A, B, RET1, RW1, DIR1, ERR1], o2 Optic[I, S, T, A, B, RET2, RW2, DIR2, ERR2], o3 Optic[I, S, T, A, B, RET3, RW3, DIR3, ERR3], o4 Optic[I, S, T, A, B, RET4, RW4, DIR4, ERR4]) Optic[I, S, T, A, B, CompositionTree[CompositionTree[RET1, RET2], CompositionTree[RET3, RET4]], CompositionTree[CompositionTree[RW1, RW2], CompositionTree[RW3, RW4]], UniDir, CompositionTree[CompositionTree[ERR1, ERR2], CompositionTree[ERR3, ERR4]]] {
	return coalesceN[I, S, T, A, B, CompositionTree[CompositionTree[RET1, RET2], CompositionTree[RET3, RET4]], CompositionTree[CompositionTree[RW1, RW2], CompositionTree[RW3, RW4]], CompositionTree[CompositionTree[ERR1, ERR2], CompositionTree[ERR3, ERR4]]](o1, o2, o3, o4)
}

// The CoalesceN combinator focuses on the elements of the first optic that yields at least one element.
func CoalesceN[I, S, T, A, B, RET, RW, DIR, ERR any](o ...Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, T, A, B, RET, RW, UniDir, ERR] {
	naked := MustGet(
		SliceOf(
			Compose(
				TraverseSlice[Optic[I, S, T, A, B, RET, RW, DIR, ERR]](),
				UpCast[Optic[I, S, T, A, B, RET, RW, DIR, ERR], nakedOptic[I, S, T, A, B]](),
			),
			len(o),
		),
		o,
	)

	return coalesceN[I, S, T, A, B, RET, RW, ERR](naked...)

}

func coalesceN[I, S, T, A, B, RET, RW, ERR any](o ...nakedOptic[I, S, T, A, B]) Optic[I, S, T, A, B, RET, RW, UniDir, ERR] {
	return CombiTraversal[RET, RW, ERR, I, S, T, A, B](
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(ValueIE[I, A]) bool) {
				for _, l := range o {
					ok := false
					l.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
						focusIndex, v, focusErr := val.Get()
						ok = true
						return yield(ValIE(focusIndex, v, JoinCtxErr(ctx, focusErr)))
					})
					if ok {
						break
					}
				}
			}
		},
		func(ctx context.Context, source S) (int, error) {
			for _, l := range o {
				llen, err := l.AsLengthGetter()(ctx, source)
				err = JoinCtxErr(ctx, err)
				if err != nil {
					return 0, err
				}
				if llen > 0 {
					return llen, nil
				}
			}
			return 0, nil
		},
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {
			var lastT T
			for _, l := range o {
				ok := false
				var err error
				lastT, err = l.AsModify()(ctx, func(index I, focus A) (B, error) {
					ok = true
					ret, err := fmap(index, focus)
					return ret, JoinCtxErr(ctx, err)
				}, source)

				err = JoinCtxErr(ctx, err)
				if err != nil {
					var t T
					return t, err
				}
				if ok {
					break
				}
			}
			return lastT, ctx.Err()
		},
		nil,
		func(indexA, indexB I) bool {
			if len(o) > 0 {
				return o[0].AsIxMatch()(indexA, indexB)
			}
			return false
		},
		ExprDefVarArgs(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				var exprs []expr.OpticExpression
				for _, v := range o {
					exprs = append(exprs, v.AsExpr())
				}
				return expr.Coalesce{
					OpticTypeExpr: ot,
					Optics:        exprs,
				}
			},
			o...,
		),
	)
}

// The Embed combinator reverses the direction of a [Prism].
// The result is a [Getter] as a [Prism] may focus no elements.
//
// See:
//   - [AsReverseGet] for a version that will reverse an [Iso] and return an [Iso]
func Embed[I, S, T, A, B any, RET any, RW TReadWrite, DIR TBiDir, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[Void, B, A, T, S, ReturnOne, ReadOnly, UniDir, ERR] {
	return unsafeAsReverseGet[ReadOnly, UniDir](o)
}

func filterIgnoreErr(ctx context.Context, err error, ignorePred func(context.Context, error) (bool, error)) (error, bool) {
	if ctx.Err() != nil {
		return errors.Join(ctx.Err(), err), false
	}

	if err == nil {
		return nil, false
	}

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) || errors.Is(err, yieldAfterBreak) || errors.Is(err, matched) {
		return err, false
	}

	if ignorePred != nil {
		ignore, ferr := ignorePred(ctx, err)
		if errors.Is(ferr, ErrEmptyGet) {
			return err, false
		}

		if ferr != nil {
			return ferr, false
		}

		if ignore {
			return nil, true
		} else {
			return err, false
		}
	} else {
		return nil, true
	}
}

// The Ignore combinator returns an optic that ignores errors that match the given predicate.
// under modification the original values are retained.
//
// Note: Ignore does not ignore the [context.DeadlineExceeded] or [context.Cancelled] errors.
//
// See:
//   - [Catch] for an optic that can handle errors.
func Ignore[I, S, A, RET, RW, DIR, ERR any, PERR TPure](o Optic[I, S, S, A, A, RET, RW, DIR, ERR], pred Predicate[error, PERR]) Optic[I, S, S, A, A, ReturnMany, RW, DIR, Pure] {

	predOp := pred.AsOpGet()

	return Omni[I, S, S, A, A, ReturnMany, RW, DIR, Pure](
		func(ctx context.Context, source S) (I, A, error) {
			i, a, err := o.AsGetter()(ctx, source)
			err, filtered := filterIgnoreErr(ctx, err, predOp)
			if filtered {
				return i, a, ErrEmptyGet
			} else {
				return i, a, err
			}
		},
		func(ctx context.Context, focus A, source S) (S, error) {
			t, err := o.AsSetter()(ctx, focus, source)
			err, filtered := filterIgnoreErr(ctx, err, predOp)
			if filtered {
				return source, nil
			} else {
				return t, err
			}
		},
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					index, focus, err := val.Get()
					err, filtered := filterIgnoreErr(ctx, err, predOp)
					if filtered {
						//Ignore the errors in the output
						return true
					} else {
						return yield(ValIE(index, focus, err))
					}
				})
			}
		},
		func(ctx context.Context, source S) (length int, retErr error) {

			o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
				err := val.Error()
				err, filtered := filterIgnoreErr(ctx, err, predOp)
				if filtered {
					//Ignore the errors in the output
					return true
				} else {
					if err != nil {
						retErr = err
						return false
					}
					length++
					return true
				}
			})

			return
		},
		func(ctx context.Context, fmap func(index I, focus A) (A, error), source S) (S, error) {
			t, err := o.AsModify()(ctx, fmap, source)

			err, filtered := filterIgnoreErr(ctx, err, predOp)
			if filtered {
				return source, nil
			} else {
				return t, err
			}
		},
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				o.AsIxGetter()(ctx, index, source)(func(val ValueIE[I, A]) bool {
					index, focus, err := val.Get()
					err, filtered := filterIgnoreErr(ctx, err, predOp)
					if filtered {
						//Ignore the errors in the output
						return true
					} else {
						return yield(ValIE(index, focus, err))
					}
				})
			}
		},
		o.AsIxMatch(),
		func(ctx context.Context, focus A) (S, error) {
			s, err := o.AsReverseGetter()(ctx, focus)
			err, filtered := filterIgnoreErr(ctx, err, predOp)
			if filtered {
				return s, ErrEmptyGet
			} else {
				return s, err
			}

		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Ignore{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
					Predicate:     pred.AsExpr(),
				}
			},
			o,
			pred,
		),
	)
}

// The Stop combinator returns an optic that stops on errors that match the given predicate.
//
// See:
//   - [Catch] for an optic that can handle errors.
//   - [Ignore] for an optic that ignores errors.
func Stop[I, S, A, RET, RW, DIR, ERR any, PERR TPure](o Optic[I, S, S, A, A, RET, RW, DIR, ERR], stopPred Predicate[error, PERR]) Optic[I, S, S, A, A, ReturnMany, ReadOnly, UniDir, Pure] {

	return Omni[I, S, S, A, A, ReturnMany, ReadOnly, UniDir, Pure](
		func(ctx context.Context, source S) (I, A, error) {
			i, a, err := o.AsGetter()(ctx, source)
			isBreak := Must(PredGet(ctx, stopPred, err))
			if isBreak {
				return i, a, ErrEmptyGet
			} else {
				return i, a, err
			}
		},
		func(ctx context.Context, focus A, source S) (S, error) {
			var s S
			return s, UnsupportedOpticMethod
		},
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					index, focus, err := val.Get()
					isBreak := Must(PredGet(ctx, stopPred, err))
					if isBreak {
						//Stop as soon as an error occurs.
						return false
					} else {
						return yield(ValIE(index, focus, err))
					}
				})
			}
		},
		func(ctx context.Context, source S) (length int, retErr error) {

			o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
				err := val.Error()
				isBreak := Must(PredGet(ctx, stopPred, err))
				if isBreak {
					//Stop as soon as an error occurs.
					return false
				} else {
					if err != nil {
						retErr = err
						return false
					}
					length++
					return true
				}
			})

			return
		},
		func(ctx context.Context, fmap func(index I, focus A) (A, error), source S) (S, error) {
			var s S
			return s, UnsupportedOpticMethod
		},
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				o.AsIxGetter()(ctx, index, source)(func(val ValueIE[I, A]) bool {
					index, focus, err := val.Get()
					isBreak := Must(PredGet(ctx, stopPred, err))
					if isBreak {
						//Stop as soon as an error occurs.
						return false
					} else {
						return yield(ValIE(index, focus, err))
					}
				})
			}
		},
		o.AsIxMatch(),
		func(ctx context.Context, focus A) (S, error) {
			var s S
			return s, UnsupportedOpticMethod
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Stop{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
					Predicate:     stopPred.AsExpr(),
				}
			},
			o,
			stopPred,
		),
	)

}

func mayCatchErr(ctx context.Context, err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) || errors.Is(err, yieldAfterBreak) || errors.Is(err, matched) {
		return false
	}

	return true
}

// The Catch combinator returns an optic that enables an error to be converted to a focus or to remain in error..
//
// Note: Catch does not catch the [context.DeadlineExceeded] or [context.Cancelled] errors.
// See:
//   - [CatchP] for a read write polymorphic version
//   - [Ignore] for an optic that can ignore errors.
func Catch[I, S, T, A, B, RET, RW, DIR, ERR any, CT any, CB any, CRET TReturnOne, CRW any, CDIR any, CERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], catch Optic[I, error, CT, A, CB, CRET, CRW, CDIR, CERR]) Optic[I, S, T, A, B, RET, ReadOnly, UniDir, ERR] {
	return Ro(CatchP(o, catch, Throw[T]()))
}

// The CatchP combinator returns an optic that enables an error to be converted to a focus or to remain in error..
//
// Note: Catch does not catch the [context.DeadlineExceeded] or [context.Cancelled] errors.
// See:
//   - [Catch] for a read only non polymorphic version
//   - [Ignore] for an optic that can ignore errors.
func CatchP[I, S, T, A, B, RET, RW, DIR, ERR any, CAT any, CAB any, CARET, CBRET TReturnOne, CARW any, CADIR any, CAERR, CBERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], catchA Optic[I, error, CAT, A, CAB, CARET, CARW, CADIR, CAERR], catchT Operation[error, T, CBRET, CBERR]) Optic[I, S, T, A, B, RET, RW, DIR, ERR] {

	return Omni[I, S, T, A, B, RET, RW, DIR, ERR](
		func(ctx context.Context, source S) (I, A, error) {
			i, a, err := o.AsGetter()(ctx, source)
			if mayCatchErr(ctx, err) {
				return catchA.AsGetter()(ctx, err)
			} else {
				return i, a, err
			}
		},
		func(ctx context.Context, focus B, source S) (T, error) {
			t, err := o.AsSetter()(ctx, focus, source)
			if mayCatchErr(ctx, err) {
				return catchT.AsOpGet()(ctx, err)
			} else {
				return t, err
			}
		},
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					index, focus, err := val.Get()
					if mayCatchErr(ctx, err) {
						return yield(ValIE(catchA.AsGetter()(ctx, err)))
					} else {
						return yield(ValIE(index, focus, err))
					}
				})
			}
		},
		nil,
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {
			t, err := o.AsModify()(ctx, fmap, source)
			if mayCatchErr(ctx, err) {
				return catchT.AsOpGet()(ctx, err)
			} else {
				return t, err
			}
		},
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				o.AsIxGetter()(ctx, index, source)(func(val ValueIE[I, A]) bool {
					index, focus, err := val.Get()
					if mayCatchErr(ctx, err) {
						cont := true
						catchA.AsIxGetter()(ctx, index, err)(func(val ValueIE[I, A]) bool {
							index, focus, err := val.Get()
							cont := yield(ValIE(index, focus, err))
							return cont
						})
						return cont
					} else {
						return yield(ValIE(index, focus, err))
					}
				})
			}
		},
		o.AsIxMatch(),
		func(ctx context.Context, focus B) (T, error) {
			t, err := o.AsReverseGetter()(ctx, focus)
			if mayCatchErr(ctx, err) {
				return catchT.AsOpGet()(ctx, err)
			} else {
				return t, err
			}
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Catch{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
					CatchA:        catchA.AsExpr(),
					CatchT:        catchT.AsExpr(),
				}
			},
			o,
			catchA,
			catchT,
		),
	)
}

// The Reduce combinator returns an [Iteration] that focuses the result of applying the [Reducer] to the given optic.
//
// If the given optic is empty then no reduced value is focused.
func Reduce[I, S, T, A, B, RET, RW, DIR, ERR, SA, RA, ERRREDUCE any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], m ReductionP[SA, A, RA, ERRREDUCE]) Optic[Void, S, S, RA, RA, RET, ReadOnly, UniDir, CompositionTree[ERR, ERRREDUCE]] {
	return EErrL(MapReduce(o, Identity[A](), m))
}

// The MapReduce combinator returns an [Iteration] that focuses the result of applying the [Reducer] to the given optic after applying the mapper operator.
func MapReduce[I, S, T, A, B, RET, RW, DIR, ERR, SM, M, RM, ERRMAP, ERRREDUCE any, MRET TReturnOne](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], mapper Operation[A, M, MRET, ERRMAP], m ReductionP[SM, M, RM, ERRREDUCE]) Optic[Void, S, S, RM, RM, RET, ReadOnly, UniDir, CompositionTree[CompositionTree[ERR, ERRREDUCE], ERRMAP]] {
	return CombiIteration[RET, CompositionTree[CompositionTree[ERR, ERRREDUCE], ERRMAP], Void, S, S, RM, RM](
		func(ctx context.Context, source S) SeqIE[Void, RM] {
			return func(yield func(ValueIE[Void, RM]) bool) {
				ret, err := m.Empty(ctx)
				if err != nil {
					var rm RM
					yield(ValIE(Void{}, rm, err))
					return
				}
				var retErr error
				found := false

				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					_, a, focusErr := val.Get()
					retErr = JoinCtxErr(ctx, focusErr)
					if retErr != nil {
						return false
					}

					var mapped M
					mapped, retErr = mapper.AsOpGet()(ctx, a)
					if retErr != nil {
						return false
					}

					ret, retErr = m.Reduce(ctx, ret, mapped)
					retErr = JoinCtxErr(ctx, focusErr)
					if retErr != nil {
						return false
					}

					found = true
					return true
				})

				if retErr != nil {
					var v RM
					yield(ValIE(Void{}, v, retErr))
					return
				}

				if found {
					v, err := m.End(ctx, ret)
					yield(ValIE(Void{}, v, err))
				}
			}
		},
		nil,
		nil,
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.MapReduce{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
					Reducer:       m.AsExpr(),
					Mapper:        mapper.AsExpr(),
				}
			},
			o,
			m,
			mapper,
		),
	)
}

// The First combinator returns an optic that focuses the first focused element.
func First[I, S, T, A, RET, RW, DIR, ERR any](o Optic[I, S, T, A, A, RET, RW, DIR, ERR]) Optic[I, S, T, A, A, RET, RW, UniDir, ERR] {
	//If o is ReturnOne then Taking 1 will always return one.
	return UnsafeReconstrain[RET, RW, UniDir, ERR](Taking(o, 1))
}

// The FirstOrDefault combinator returns an optic that focuses the first focused element or the default value if no elements are focused.
// See:
//   - FirstOrDefaultI for an index aware version
//   - FirstOrError for a version that returns an error
func FirstOrDefault[I, S, T, A, RET, RW, DIR, ERR any](o Optic[I, S, T, A, A, RET, RW, DIR, ERR], defaultVal A) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, ERR] {
	return FirstOrDefaultI(
		EErrL(ReIndexed(o, Const[I](Void{}), Const[lo.Tuple2[Void, Void]](true))),
		Void{},
		defaultVal,
	)
}

// The FirstOrError combinator returns an optic that focuses the first focused element or returns the given error if no elements are focused.
// See:
//   - [FirstOrDefault] for a version that returns a default value
func FirstOrError[I, S, A, RET, RW, DIR, ERR any](o Optic[I, S, S, A, A, RET, RW, DIR, ERR], err error) Optic[I, S, S, A, A, ReturnOne, RW, UniDir, Err] {
	return UnsafeReconstrain[ReturnOne, RW, UniDir, Err](Taking(
		Coalesce(
			o,
			ErrorI[I, S, A](err, o.AsIxMatch()),
		),
		1,
	))
}

// The Last combinator returns an optic that focuses the last focused element.
func Last[I, S, T, A, RET, RW, DIR, ERR any](o Optic[I, S, T, A, A, RET, RW, DIR, ERR]) Optic[I, S, T, A, A, RET, RW, UniDir, ERR] {
	return CombiTraversal[RET, RW, ERR, I, S, T, A, A](
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				var index I
				var val A
				found := false
				var err error

				o.AsIter()(ctx, source)(func(valie ValueIE[I, A]) bool {
					focusIndex, v, focusErr := valie.Get()
					err = JoinCtxErr(ctx, focusErr)
					if err != nil {
						return false
					}
					index = focusIndex
					val = v
					found = true
					return true
				})

				if found {
					yield(ValIE(index, val, err))
				}
			}
		},
		nil,
		func(ctx context.Context, fmap func(index I, focus A) (A, error), source S) (T, error) {

			l, err := o.AsLengthGetter()(ctx, source)
			if err != nil {
				var t T
				return t, err
			}

			l -= 1
			i := 0
			return o.AsModify()(ctx, func(index I, focus A) (A, error) {
				if i == l {
					return fmap(index, focus)
				}

				i++
				return focus, ctx.Err()
			}, source)
		},
		nil,
		o.AsIxMatch(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Last{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The LastOrDefault combinator returns an optic that focuses the last focused element or the default value if no elements are focused.
// See:
//   - [LastOrDefaultI] for an index aware version
//   - [LastOrError] for a version that returns an error
func LastOrDefault[I, S, T, A, RET, RW, DIR, ERR any](o Optic[I, S, T, A, A, RET, RW, DIR, ERR], defaultVal A) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, ERR] {
	return LastOrDefaultI(
		EErrL(ReIndexed(o, Const[I](Void{}), Const[lo.Tuple2[Void, Void]](true))),
		Void{},
		defaultVal,
	)
}

// The LastOrError combinator returns an optic that focuses the last focused element or returns the given error if no elements are focused.
// See:
//   - [LastOrDefault] for a version that returns a default value
func LastOrError[I, S, A, RET, RW, DIR, ERR any](o Optic[I, S, S, A, A, RET, RW, DIR, ERR], err error) Optic[I, S, S, A, A, ReturnOne, RW, UniDir, Err] {
	return UnsafeReconstrain[ReturnOne, RW, UniDir, Err](
		Coalesce(
			Last(o),
			ErrorI[I, S, A](err, o.AsIxMatch()),
		),
	)
}

// The MaxOf combinator returns an optic that focuses the element with the maximum value as focused by the cmpBy optic
func MaxOf[C cmp.Ordered, I, S, T, A any, RETI any, RETJ TReturnOne, RWI, DIRI any, ERR, ERRCMP any](o Optic[I, S, T, A, A, RETI, RWI, DIRI, ERR], cmpBy Operation[A, C, RETJ, ERRCMP]) Optic[I, S, T, A, A, RETI, RWI, UniDir, CompositionTree[ERR, ERRCMP]] {
	return First(Ordered(o, Desc(OrderBy(cmpBy))))
}

// The MinOf combinator returns an optic that focuses the element with the minimum value as focused by the cmpBy optic
func MinOf[C cmp.Ordered, I, S, T, A any, RETI any, RETJ TReturnOne, RWI, DIRI any, ERR, ERRCMP any](o Optic[I, S, T, A, A, RETI, RWI, DIRI, ERR], cmpBy Operation[A, C, RETJ, ERRCMP]) Optic[I, S, T, A, A, RETI, RWI, UniDir, CompositionTree[ERR, ERRCMP]] {
	return First(Ordered(o, OrderBy(cmpBy)))
}

// The Ordered combinator focuses on the elements of the given optic in the order defined by the orderBy predicate
//
// Note: under modification the elements are focused in order but the result retains the original order.
//
// See:
//   - [OrderedCol] for a operation that re-orders a collection.
//   - [OrderedI] for an index aware version.
func Ordered[I, S, T, A, B any, RET, RW, DIR, ERR, PERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], orderBy OrderByPredicate[A, PERR]) Optic[I, S, T, A, B, RET, RW, UniDir, CompositionTree[ERR, PERR]] {
	return OrderedI(o, PredT2ToOpT2I[I, I](orderBy))
}

var matched = errors.New("matched")

// The Matching combinator returns an optic that focuses on either the matching value of the given [Prism] or the source transformed to the result type.
func Matching[I, S, T, A, B, RET any, RW TReadWrite, DIR TBiDir, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, S, mo.Either[T, A], mo.Either[T, A], ReturnOne, ReadOnly, UniDir, ERR] {
	return CombiGetter[ERR, I, S, S, mo.Either[T, A], mo.Either[T, A]](
		func(ctx context.Context, source S) (I, mo.Either[T, A], error) {

			var a A
			var i I

			t, err := o.AsModify()(ctx, func(index I, focus A) (B, error) {
				i = index
				a = focus
				var b B
				return b, matched
			}, source)

			if errors.Is(err, matched) {
				return i, mo.Right[T, A](a), nil
			} else {
				return i, mo.Left[T, A](t), err
			}

		},
		o.AsIxMatch(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Matching{
					OpticTypeExpr: ot,
					Match:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The ForEach combinator applies op to each element focused by forEach, focusing on the [Collection] of results.
//
// This can be useful to apply impure modifications to a [ReadOnly] forEach optic.
// If forEach is [ReadWrite] then [Compose] should be used instead.
// See:
//   - IForEach for a version that supports index mapping.
func ForEach[I, J, S, T, A, B, C, D any, RET, RETM any, RW any, RWM any, DIR, DIRM any, SERR any](forEach Optic[I, S, T, A, B, RET, RW, DIR, SERR], op Optic[J, A, B, C, D, RETM, RWM, DIRM, SERR]) Optic[J, S, Collection[I, B, SERR], C, D, CompositionTree[RET, RETM], RWM, UniDir, SERR] {
	return ForEachI(IxMapRight[I, J](op.AsIxMatch()), forEach, op)
}

// The Polymorphic combinator enables any optic to be made polymorphic at the cost of becoming [ReadOnly]
func Polymorphic[TP, BP, I, S, T, A, B any, RET any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, TP, A, BP, RET, ReadOnly, UniDir, ERR] {
	return Omni[I, S, TP, A, BP, RET, ReadOnly, UniDir, ERR](
		o.AsGetter(),
		func(ctx context.Context, focus BP, source S) (TP, error) {
			var tp TP
			return tp, UnsupportedOpticMethod
		},
		o.AsIter(),
		o.AsLengthGetter(),
		func(ctx context.Context, fmap func(index I, focus A) (BP, error), source S) (TP, error) {
			var tp TP
			return tp, UnsupportedOpticMethod
		},
		o.AsIxGetter(),
		o.AsIxMatch(),
		func(ctx context.Context, focus BP) (TP, error) {
			var tp TP
			return tp, UnsupportedOpticMethod
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Polymorphic{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

func EditDistance[I, S, T, A, B any, RET any, RW any, DIR any, ERR, ERRP any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], allowedEdits EditType, equal Predicate[lo.Tuple2[A, A], ERRP], size int) Optic[Void, lo.Tuple2[S, S], lo.Tuple2[S, S], float64, float64, ReturnOne, ReadOnly, UniDir, CompositionTree[ERR, ERRP]] {
	return EErrR(EErrTransR(EditDistanceI(o, allowedEdits, Const[lo.Tuple2[I, I]](true), equal, size)))
}

func Concat[I, S, A any, RET1, RW1, DIR1, ERR1, RET2, RW2, DIR2, ERR2 any](o1 Optic[I, S, S, A, A, RET1, RW1, DIR1, ERR1], o2 Optic[I, S, S, A, A, RET2, RW2, DIR2, ERR2]) Optic[I, S, S, A, A, ReturnMany, CompositionTree[RW1, RW2], UniDir, CompositionTree[ERR1, ERR2]] {
	return concatN[I, S, A, CompositionTree[RW1, RW2], CompositionTree[ERR1, ERR2]](o1, o2)
}

func Concat3[I, S, A any, RET1, RW1, DIR1, ERR1, RET2, RW2, DIR2, ERR2, RET3, RW3, DIR3, ERR3 any](o1 Optic[I, S, S, A, A, RET1, RW1, DIR1, ERR1], o2 Optic[I, S, S, A, A, RET2, RW2, DIR2, ERR2], o3 Optic[I, S, S, A, A, RET3, RW3, DIR3, ERR3]) Optic[I, S, S, A, A, ReturnMany, CompositionTree[CompositionTree[RW1, RW2], RW3], UniDir, CompositionTree[CompositionTree[ERR1, ERR2], ERR3]] {
	return concatN[I, S, A, CompositionTree[CompositionTree[RW1, RW2], RW3], CompositionTree[CompositionTree[ERR1, ERR2], ERR3]](o1, o2, o3)
}

func Concat4[I, S, A any, RET1, RW1, DIR1, ERR1, RET2, RW2, DIR2, ERR2, RET3, RW3, DIR3, ERR3, RET4, RW4, DIR4, ERR4 any](o1 Optic[I, S, S, A, A, RET1, RW1, DIR1, ERR1], o2 Optic[I, S, S, A, A, RET2, RW2, DIR2, ERR2], o3 Optic[I, S, S, A, A, RET3, RW3, DIR3, ERR3], o4 Optic[I, S, S, A, A, RET4, RW4, DIR4, ERR4]) Optic[I, S, S, A, A, ReturnMany, CompositionTree[CompositionTree[RW1, RW2], CompositionTree[RW3, RW4]], UniDir, CompositionTree[CompositionTree[ERR1, ERR2], CompositionTree[ERR3, ERR4]]] {
	return concatN[I, S, A, CompositionTree[CompositionTree[RW1, RW2], CompositionTree[RW3, RW4]], CompositionTree[CompositionTree[ERR1, ERR2], CompositionTree[ERR3, ERR4]]](o1, o2, o3, o4)
}

// The ConcatN combinator returns all the focuses of the given optics in sequence.
func ConcatN[I, S, A any, RET any, RW any, DIR any, ERR any](o ...Optic[I, S, S, A, A, RET, RW, DIR, ERR]) Optic[I, S, S, A, A, ReturnMany, RW, UniDir, ERR] {
	naked := MustGet(
		SliceOf(
			Compose(
				TraverseSlice[Optic[I, S, S, A, A, RET, RW, DIR, ERR]](),
				UpCast[Optic[I, S, S, A, A, RET, RW, DIR, ERR], nakedOptic[I, S, S, A, A]](),
			),
			len(o),
		),
		o,
	)

	return concatN[I, S, A, RW, ERR](naked...)
}

type nakedOpticRO[I, S, A any] interface {
	AsIter() IterFunc[I, S, A]
	AsGetter() GetterFunc[I, S, A]
	AsIxGetter() IxGetterFunc[I, S, A]
	AsIxMatch() IxMatchFunc[I]
	AsLengthGetter() LengthGetterFunc[S]
	AsOpGet() OpGetFunc[S, A]
	AsExpr() expr.OpticExpression
	AsExprHandler() func(ctx context.Context) (ExprHandler, error)
	OpticType() expr.OpticType
}

func concatN[I, S, A any, RW any, ERR any](o ...nakedOptic[I, S, S, A, A]) Optic[I, S, S, A, A, ReturnMany, RW, UniDir, ERR] {
	return CombiTraversal[ReturnMany, RW, ERR, I, S, S, A, A](
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				cont := true

				for _, optic := range o {
					optic.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
						index, focus, err := val.Get()
						cont = yield(ValIE(index, focus, err))
						return cont
					})
					if !cont {
						return
					}
				}
			}
		},
		func(ctx context.Context, source S) (int, error) {
			l := 0
			for _, optic := range o {
				res, err := optic.AsLengthGetter()(ctx, source)
				if err != nil {
					return 0, err
				}
				l += res
			}

			return l, nil
		},
		func(ctx context.Context, fmap func(index I, focus A) (A, error), source S) (S, error) {
			for _, optic := range o {
				var err error
				source, err = optic.AsModify()(ctx, fmap, source)
				if err != nil {
					var s S
					return s, err
				}
			}
			return source, nil
		},
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				cont := true

				for _, optic := range o {
					optic.AsIxGetter()(ctx, index, source)(func(val ValueIE[I, A]) bool {
						index, focus, err := val.Get()
						cont = yield(ValIE(index, focus, err))
						return cont
					})
					if !cont {
						return
					}
				}
			}
		},
		func(indexA, indexB I) bool {
			ret := false
			for _, optic := range o {
				res := optic.AsIxMatch()(indexA, indexB)
				if !res {
					return false
				}
				ret = true
			}
			return ret
		},
		ExprDefVarArgs(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {

				var exprs []expr.OpticExpression
				for _, v := range o {
					exprs = append(exprs, v.AsExpr())
				}

				return expr.Concat{
					OpticTypeExpr: ot,
					Optics:        exprs,
				}
			},
			o...,
		),
	)
}

// The Grouped combinator reduces all focuses with the same index using the given reducer.
func Grouped[I, S, T, A, B, RET, RW, DIR, ERR any, RS, RR, RERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], reducer ReductionP[RS, A, RR, RERR]) Optic[I, S, S, RR, RR, RET, ReadOnly, UniDir, CompositionTree[ERR, RERR]] {

	return CombiIteration[RET, CompositionTree[ERR, RERR], I, S, S, RR, RR](
		func(ctx context.Context, source S) SeqIE[I, RR] {
			return func(yield func(ValueIE[I, RR]) bool) {
				var grouped []lo.Tuple2[I, RS]

				ixMatch := o.AsIxMatch()

				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					index, focus, err := val.Get()
					if err != nil {
						var i I
						var rr RR
						return yield(ValIE(i, rr, err))
					}

					matched := false
					for i, _ := range grouped {
						grouping := &grouped[i]

						match := ixMatch(index, grouping.A)

						if match {
							matched = true
							reduced, err := reducer.Reduce(ctx, grouping.B, focus)
							if err != nil {
								var i I
								var rr RR
								return yield(ValIE(i, rr, err))
							}
							grouping.B = reduced
						}
					}

					if !matched {
						empty, err := reducer.Empty(ctx)
						if err != nil {
							var i I
							var rr RR
							yield(ValIE(i, rr, err))
							return false
						}
						reduced, err := reducer.Reduce(ctx, empty, focus)
						if err != nil {
							var i I
							var rr RR
							return yield(ValIE(i, rr, err))
						}
						grouped = append(grouped, lo.T2(index, reduced))
					}

					return true
				})

				for _, v := range grouped {
					rr, err := reducer.End(ctx, v.B)
					err = JoinCtxErr(ctx, err)
					if !yield(ValIE(v.A, rr, err)) {
						return
					}
				}
			}
		},
		nil,
		func(ctx context.Context, index I, source S) SeqIE[I, RR] {
			return func(yield func(ValueIE[I, RR]) bool) {
				grouping, err := reducer.Empty(ctx)
				if err != nil {
					var i I
					var rr RR
					yield(ValIE(i, rr, err))
					return
				}
				matched := false

				o.AsIxGetter()(ctx, index, source)(func(val ValueIE[I, A]) bool {
					_, focus, err := val.Get()
					if err != nil {
						var i I
						var rr RR
						return yield(ValIE(i, rr, err))
					}
					matched = true

					reduced, err := reducer.Reduce(ctx, grouping, focus)
					if err != nil {
						var i I
						var rr RR
						return yield(ValIE(i, rr, err))
					}
					grouping = reduced

					return true
				})

				if matched {
					rr, err := reducer.End(ctx, grouping)
					err = JoinCtxErr(ctx, err)
					yield(ValIE(index, rr, err))
				}
			}
		},
		o.AsIxMatch(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Grouped{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
					Reducer:       reducer.AsExpr(),
				}
			},
			o,
			reducer,
		),
	)

}

func Prefixed[I, J, S, A any, RET TReturnOne, RW TReadWrite, DIR, ERR any, PERR any](path Optic[I, S, S, Collection[J, A, ERR], Collection[J, A, ERR], RET, RW, DIR, ERR], valMatch Predicate[lo.Tuple2[A, A], PERR], ix J, val A) Optic[Void, S, S, S, S, ReturnMany, ReadWrite, BiDir, CompositionTree[ERR, PERR]] {

	return CombiPrism[ReadWrite, BiDir, CompositionTree[ERR, PERR], Void, S, S, S, S](
		func(ctx context.Context, source S) (mo.Either[S, ValueI[Void, S]], error) {
			pathCol, err := path.AsOpGet()(ctx, source)
			if err != nil {
				return mo.Left[S, ValueI[Void, S]](source), err
			}

			firstPath, ok, err := GetFirstContext(ctx, TraverseColIE[J, A, ERR](pathCol.AsIxMatch()), pathCol)
			if err != nil {
				return mo.Left[S, ValueI[Void, S]](source), err
			}

			pathMatch := false
			if ok {
				pathMatch, err = PredGet(ctx, valMatch, lo.T2(firstPath, val))
				if err != nil {
					return mo.Left[S, ValueI[Void, S]](source), err
				}
			}

			if pathMatch {

				x := ColOf(Dropping(TraverseColIE[J, A, ERR](pathCol.AsIxMatch()), 1))

				updatedPath, err := ModifyContext(
					ctx,
					path,
					x,
					source,
				)
				if err != nil {
					return mo.Left[S, ValueI[Void, S]](source), err
				}

				return mo.Right[S, ValueI[Void, S]](ValI(Void{}, updatedPath)), nil

			} else {
				return mo.Left[S, ValueI[Void, S]](source), nil
			}
		},
		func(ctx context.Context, focus S) (S, error) {
			pathCol, err := path.AsOpGet()(ctx, focus)
			if err != nil {
				var s S
				return s, err
			}

			ret, err := ModifyContext(
				ctx,
				path,
				PrependCol(
					ValColIE[J, A, ERR](pathCol.AsIxMatch(), ValIE(ix, val, nil)),
				),
				focus,
			)

			return ret, err
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Prefixed{
					OpticTypeExpr: ot,
					Path:          path.AsExpr(),
					Val:           val,
					ValMatch:      valMatch.AsExpr(),
				}
			},
			path,
			valMatch,
		),
	)
}
