package optic

import (
	"context"
	"log"
	"reflect"

	"github.com/samber/mo"
	"github.com/spearson78/go-optic/expr"
)

func WithLogging[I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, T, A, B, RET, RW, DIR, ERR] {
	return UnsafeOmni[I, S, T, A, B, RET, RW, DIR, ERR](
		func(ctx context.Context, source S) (I, A, error) {
			i, a, err := o.AsGetter()(ctx, source)
			log.Printf("%v.AsGetter(%v) -> (%v,%v,%v)", o.AsExpr().Short(), source, i, a, err)
			return i, a, err
		},
		func(ctx context.Context, focus B, source S) (T, error) {
			t, err := o.AsSetter()(ctx, focus, source)
			log.Printf("%v.AsSetter()(%v) -> (%v,%v)", o.AsExpr().Short(), focus, source, t)
			return t, err
		},
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				log.Printf("%v.AsIter()(%v) begin", o.AsExpr().Short(), source)
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					index, focus, err := val.Get()
					ret := yield(val)
					log.Printf("%v.AsIter()(%v) -> yield(%v,%v,%v) -> %v", o.AsExpr().Short(), source, index, focus, err, ret)
					return ret
				})
				log.Printf("%v.AsIter()(%v) end", o.AsExpr().Short(), source)
			}

		},
		func(ctx context.Context, source S) (int, error) {
			l, err := o.AsLengthGetter()(ctx, source)
			log.Printf("%v.AsLengthGetter()(%v) -> (%v,%v)", o.AsExpr().Short(), source, l, err)
			return l, err

		},
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {
			log.Printf("%v.AsModify()(%v) begin", o.AsExpr().Short(), source)
			t, err := o.AsModify()(ctx, func(index I, focus A) (B, error) {
				b, err := fmap(index, focus)
				log.Printf("%v.AsModify()(%v) fmap(%v,%v) (%v,%v)", o.AsExpr().Short(), source, index, focus, b, err)
				return b, err
			}, source)
			log.Printf("%v.AsModify()(%v) -> (%v,%v)", o.AsExpr().Short(), source, t, err)
			return t, err
		},
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				log.Printf("%v.AsIxGetter()(%v) begin", o.AsExpr().Short(), source)
				o.AsIxGetter()(ctx, index, source)(func(val ValueIE[I, A]) bool {
					focusIndex, focus, err := val.Get()
					ret := yield(val)
					log.Printf("%v.AsIxGetter()(%v,%v) -> yield(%v,%v,%v) -> %v", o.AsExpr().Short(), index, source, focusIndex, focus, err, ret)
					return ret
				})
				log.Printf("%v.AsIxGetter()(%v) end", o.AsExpr().Short(), source)
			}
		},
		func(indexA, indexB I) bool {
			ret := o.AsIxMatch()(indexA, indexB)
			log.Printf("%v.AsIxMatch()(%v,%v) -> (%v)", o.AsExpr().Short(), indexA, indexB, ret)
			return ret
		},
		func(ctx context.Context, focus B) (T, error) {
			t, err := o.AsReverseGetter()(ctx, focus)
			log.Printf("%v.AsReverseGetter()(%v) -> (%v,%v)", o.AsExpr().Short(), focus, t, err)
			return t, err
		},
		o.AsExprHandler(),
		o.AsExpr,
	)
}

// The WithOption combinator returns an applicative optic enabling the return of optional values.
// under modification
//   - if any focus is None then the result is short circuited to None.
//   - if the source is empty then the result is converted to None.
func WithOption[I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, mo.Option[T], A, mo.Option[B], RET, RW, DIR, ERR] {

	return Omni[I, S, mo.Option[T], A, mo.Option[B], RET, RW, DIR, ERR](
		o.AsGetter(),
		func(ctx context.Context, focus mo.Option[B], source S) (mo.Option[T], error) {
			f, ok := focus.Get()
			if ok {
				t, err := o.AsSetter()(ctx, f, source)
				return mo.TupleToOption(t, err == nil), err
			} else {
				return mo.None[T](), nil
			}
		},
		o.AsIter(),
		o.AsLengthGetter(),
		func(ctx context.Context, fmap func(index I, focus A) (mo.Option[B], error), source S) (ret mo.Option[T], retErr error) {
			defer handleAbortModify(&retErr)

			var res T
			ok := false
			res, retErr = o.AsModify()(ctx, func(focusIndex I, focus A) (B, error) {
				ob, err := fmap(focusIndex, focus)
				var b B
				err = JoinCtxErr(ctx, err)
				if err != nil {
					return b, err
				}
				b, ok = ob.Get()
				if !ok {
					ret = mo.None[T]()
					abortModify(&retErr)
				}
				return b, nil
			}, source)

			ret = mo.TupleToOption(res, ok && retErr == nil)

			return
		},
		o.AsIxGetter(),
		o.AsIxMatch(),
		func(ctx context.Context, focus mo.Option[B]) (mo.Option[T], error) {
			if f, ok := focus.Get(); ok {
				t, err := o.AsReverseGetter()(ctx, f)
				return mo.TupleToOption(t, err == nil), err
			} else {
				return mo.None[T](), nil
			}
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.WithOption{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The WithFunc combinator returns an applicative optic that enables the passing of a parameter into the optic.
func WithFunc[P, I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, func(P) (T, error), A, func(P) (B, error), RET, RW, DIR, ERR] {

	return Omni[I, S, func(P) (T, error), A, func(P) (B, error), RET, RW, DIR, ERR](
		o.AsGetter(),
		func(ctx context.Context, focus func(P) (B, error), source S) (func(P) (T, error), error) {
			return func(p P) (T, error) {
				b, err := focus(p)
				err = JoinCtxErr(ctx, err)
				if err != nil {
					var t T
					return t, err
				}
				ret, err := o.AsSetter()(ctx, b, source)
				return ret, err
			}, nil
		},
		o.AsIter(),
		o.AsLengthGetter(),
		func(ctx context.Context, fmap func(index I, focus A) (func(P) (B, error), error), source S) (func(P) (T, error), error) {
			return func(p P) (T, error) {
				res, err := o.AsModify()(ctx, func(focusIndex I, focus A) (B, error) {
					fncb, err := fmap(focusIndex, focus)
					var b B
					err = JoinCtxErr(ctx, err)
					if err != nil {
						return b, err
					}
					b, err = fncb(p)
					err = JoinCtxErr(ctx, err)
					if err != nil {
						return b, err
					}
					return b, nil
				}, source)

				return res, err
			}, nil
		},
		o.AsIxGetter(),
		o.AsIxMatch(),
		func(ctx context.Context, focus func(P) (B, error)) (func(P) (T, error), error) {
			return func(p P) (T, error) {
				b, err := focus(p)
				err = JoinCtxErr(ctx, err)
				if err != nil {
					var t T
					return t, err
				}
				ret, err := o.AsReverseGetter()(ctx, b)
				return ret, err
			}, nil
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.WithFunc{
					OpticTypeExpr: ot,
					P:             reflect.TypeFor[P](),
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The WithComprehension combinator returns an applicative optic that modifies the given optic to perform a comprehension over the slice of return focuses.
func WithComprehension[I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, []T, A, []B, RET, RW, DIR, ERR] {

	return Omni[I, S, []T, A, []B, RET, RW, DIR, ERR](
		o.AsGetter(),
		func(ctx context.Context, setFocus []B, source S) ([]T, error) {

			atobLen, err := o.AsLengthGetter()(ctx, source)
			err = JoinCtxErr(ctx, err)
			if err != nil {
				return nil, err
			}

			if atobLen == 0 {
				return nil, nil
			}

			atobStat := make([]int, atobLen)

			var ret []T

			for {

				i := 0
				t, err := o.AsModify()(ctx, func(focusIndex I, focus A) (B, error) {
					idB := setFocus[atobStat[i]]
					i++
					return idB, nil
				}, source)
				err = JoinCtxErr(ctx, err)
				if err != nil {
					return nil, err
				}
				ret = append(ret, t)

				incPos := len(atobStat) - 1
				for {
					if err := ctx.Err(); err != nil {
						return nil, err
					}

					atobStat[incPos] = atobStat[incPos] + 1
					if atobStat[incPos] == len(setFocus) {
						atobStat[incPos] = 0
						incPos--
						if incPos == -1 {
							return ret, ctx.Err()
						}
					} else {
						break
					}
				}
			}
		},
		o.AsIter(),
		o.AsLengthGetter(),
		func(ctx context.Context, fmap func(index I, focus A) ([]B, error), source S) ([]T, error) {
			var atob [][]B

			var retErr error
			o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
				index, focus, err := val.Get()
				res, err := fmap(index, focus)
				err = JoinCtxErr(ctx, retErr)
				if err != nil {
					retErr = err
					return false
				}
				atob = append(atob, res)
				return true

			})
			if retErr != nil {
				return nil, retErr
			}

			if len(atob) == 0 {
				return nil, nil
			}

			atobStat := make([]int, len(atob))

			var ret []T

			for {

				i := 0
				t, err := o.AsModify()(ctx, func(focusIndex I, focus A) (B, error) {
					idB := atob[i][atobStat[i]]
					i++
					return idB, nil
				}, source)
				if err != nil {
					return nil, err
				}
				ret = append(ret, t)

				incPos := len(atobStat) - 1
				for {
					if err := ctx.Err(); err != nil {
						return nil, err
					}

					atobStat[incPos] = atobStat[incPos] + 1
					if atobStat[incPos] == len(atob[incPos]) {
						atobStat[incPos] = 0
						incPos--
						if incPos == -1 {
							return ret, nil
						}
					} else {
						break
					}
				}
			}
		},
		o.AsIxGetter(),
		o.AsIxMatch(),
		func(ctx context.Context, focus []B) ([]T, error) {
			ret := make([]T, 0, len(focus))
			for _, b := range focus {
				t, err := o.AsReverseGetter()(ctx, b)
				if err != nil {
					return nil, err
				}
				ret = append(ret, t)
			}
			return ret, nil
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.WithComprehension{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The WithEither combinator returns an applicative optic enabling the return of an Either value.
// under modification
//   - if any focus is left then the result is short circuited to this left value.
//   - if the source is empty then the result is converted to the default left value.
func WithEither[P, I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, mo.Either[P, T], A, mo.Either[P, B], RET, RW, DIR, ERR] {

	return Omni[I, S, mo.Either[P, T], A, mo.Either[P, B], RET, RW, DIR, ERR](
		o.AsGetter(),
		func(ctx context.Context, focus mo.Either[P, B], source S) (mo.Either[P, T], error) {

			if l, isLeft := focus.Left(); isLeft {
				return mo.Left[P, T](l), nil
			} else {
				t, err := o.AsSetter()(ctx, focus.MustRight(), source)
				return mo.Right[P, T](t), err
			}
		},
		o.AsIter(),
		o.AsLengthGetter(),
		func(ctx context.Context, fmap func(index I, focus A) (mo.Either[P, B], error), source S) (ret mo.Either[P, T], retErr error) {
			defer handleAbortModify(&retErr)

			res, err := o.AsModify()(ctx, func(focusIndex I, focus A) (B, error) {
				ob, err := fmap(focusIndex, focus)
				var b B
				err = JoinCtxErr(ctx, err)
				if err != nil {
					return b, err
				}
				b, ok := ob.Right()
				if !ok {
					ret = mo.Left[P, T](ob.MustLeft())
					abortModify(&retErr)
				}

				return b, nil
			}, source)

			ret = mo.Right[P, T](res)
			retErr = err

			return
		},
		o.AsIxGetter(),
		o.AsIxMatch(),
		func(ctx context.Context, focus mo.Either[P, B]) (mo.Either[P, T], error) {
			if l, isLeft := focus.Left(); isLeft {
				return mo.Left[P, T](l), nil
			} else {
				t, err := o.AsReverseGetter()(ctx, focus.MustRight())
				return mo.Right[P, T](t), err
			}
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.WithEither{
					OpticTypeExpr: ot,
					P:             reflect.TypeFor[P](),
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The WithValidation combinator returns an applicative optic able to report multiple validation errors.
func WithValidation[V, I, S, T, A, B, RET, RW any, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, mo.Either[[]V, T], A, mo.Either[V, B], RET, RW, DIR, ERR] {

	return Omni[I, S, mo.Either[[]V, T], A, mo.Either[V, B], RET, RW, DIR, ERR](
		o.AsGetter(),
		func(ctx context.Context, focus mo.Either[V, B], source S) (mo.Either[[]V, T], error) {
			if l, isLeft := focus.Left(); isLeft {
				sLen, err := o.AsLengthGetter()(ctx, source)
				if err != nil {
					return mo.Left[[]V, T](nil), err
				}

				ret := make([]V, sLen)
				for i := range ret {
					ret[i] = l
				}

				return mo.Left[[]V, T](ret), nil
			} else {
				t, err := o.AsSetter()(ctx, focus.MustRight(), source)
				return mo.Right[[]V, T](t), err
			}
		},
		o.AsIter(),
		o.AsLengthGetter(),
		func(ctx context.Context, fmap func(index I, focus A) (mo.Either[V, B], error), source S) (mo.Either[[]V, T], error) {
			var validationErrors []V

			res, err := o.AsModify()(ctx, func(focusIndex I, focus A) (B, error) {
				validation, err := fmap(focusIndex, focus)
				var b B
				err = JoinCtxErr(ctx, err)
				if err != nil {
					return b, err
				}
				if b, ok := validation.Right(); ok {
					return b, nil
				} else {
					validationErrors = append(validationErrors, validation.MustLeft())
					return b, nil
				}
			}, source)

			if len(validationErrors) > 0 {
				return mo.Left[[]V, T](validationErrors), err
			}

			return mo.Right[[]V, T](res), err
		},
		o.AsIxGetter(),
		o.AsIxMatch(),
		func(ctx context.Context, focus mo.Either[V, B]) (mo.Either[[]V, T], error) {
			if l, isLeft := focus.Left(); isLeft {
				return mo.Left[[]V, T]([]V{l}), nil
			} else {
				t, err := o.AsReverseGetter()(ctx, focus.MustRight())
				return mo.Right[[]V, T](t), err
			}
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.WithValidation{
					OpticTypeExpr: ot,
					V:             reflect.TypeFor[V](),
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}
