package optic

import (
	"context"

	"github.com/samber/mo"
)

func CombiRw[RW any, I, S, T, A, B any, RET any, DIR any, ERR any, ORW TReadWrite](o Optic[I, S, T, A, B, RET, ORW, DIR, ERR]) Optic[I, S, T, A, B, RET, RW, DIR, ERR] {
	return UnsafeReconstrain[RET, RW, DIR, ERR](o)
}

func CombiDir[DIR any, I, S, T, A, B any, RET any, RW any, ERR any, ODIR TBiDir](o Optic[I, S, T, A, B, RET, RW, ODIR, ERR]) Optic[I, S, T, A, B, RET, RW, DIR, ERR] {
	return UnsafeReconstrain[RET, RW, DIR, ERR](o)
}

func CombiEErr[ERR any, I, S, T, A, B any, RET any, RW any, DIR any, OERR TPure](o Optic[I, S, T, A, B, RET, RW, DIR, OERR]) Optic[I, S, T, A, B, RET, RW, DIR, ERR] {
	return UnsafeReconstrain[RET, RW, DIR, ERR](o)
}

func CombiColTypeErr[ERR any, I, S, T, A, B any, OERR TPure](ct CollectionType[I, S, T, A, B, OERR]) CollectionType[I, S, T, A, B, ERR] {
	return CollectionType[I, S, T, A, B, ERR]{
		traverse: CombiEErr[ERR](ct.traverse),
		toCol:    CombiColFocusErr[ERR](ct.toCol),
		fromColB: CombiColSourceErr[ERR](ct.fromColB),
	}
}

func CombiColFocusErr[ERR, I, J, S, T, A, B, RET, RW, DIR any, OERR, OAERR, OBERR TPure](o Optic[J, S, T, Collection[I, A, OAERR], Collection[I, B, OBERR], RET, RW, DIR, OERR]) Optic[J, S, T, Collection[I, A, ERR], Collection[I, B, ERR], RET, RW, DIR, ERR] {
	return unsafeColFocusErr[ERR](o)
}

func CombiColSourceErr[ERR any, I, J, S, T, A, B, RET, RW, DIR any, OERR, OSERR, OTERR TPure](o Optic[J, Collection[I, S, OSERR], Collection[I, T, OTERR], A, B, RET, RW, DIR, OERR]) Optic[J, Collection[I, S, ERR], Collection[I, T, ERR], A, B, RET, RW, DIR, ERR] {
	return unsafeColSourceErr[ERR](o)
}

func CombiColSourceFocusErr[SERR, I, J, S, T, A, B, RET, RW, DIR any, ERR, OSERR, OTERR, OAERR, OBERR TPure](o Optic[J, Collection[I, S, OSERR], Collection[I, T, OTERR], Collection[I, A, OAERR], Collection[I, B, OBERR], RET, RW, DIR, ERR]) Optic[J, Collection[I, S, SERR], Collection[I, T, SERR], Collection[I, A, SERR], Collection[I, B, SERR], RET, RW, DIR, SERR] {
	return unsafeColSourceFocusErr[SERR](o)
}

func CombiTraversal[RET, RW, ERR, I, S, T, A, B any](
	iterate func(ctx context.Context, source S) SeqIE[I, A],
	lengthGetter LengthGetterFunc[S],
	modify func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error),
	ixget func(ctx context.Context, index I, source S) SeqIE[I, A],
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, T, A, B, RET, RW, UniDir, ERR] {
	return Omni[I, S, T, A, B, RET, RW, UniDir, ERR](
		func(ctx context.Context, source S) (index I, ret A, err error) {
			err = ErrEmptyGet

			cont := true
			iterate(ctx, source)(func(val ValueIE[I, A]) bool {
				focusIndex, pa, focusErr := val.Get()
				if !cont {
					panic(yieldAfterBreak)
				}
				index = focusIndex
				err = JoinCtxErr(ctx, focusErr)
				ret = pa
				cont = false
				return false
			})
			return index, ret, err
		},
		func(ctx context.Context, f B, s S) (T, error) {
			return modify(ctx, func(index I, focus A) (B, error) {
				return f, ctx.Err()
			}, s)
		},
		iterate,
		lengthGetter,
		modify,
		ixget,
		ixMatch,
		unsupportedReverseGetter[B, T],
		exprDef,
	)
}

func CombiLens[RW, ERR any, I, S, T, A, B any](
	get func(ctx context.Context, source S) (I, A, error),
	set func(ctx context.Context, focus B, source S) (T, error),
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, T, A, B, ReturnOne, RW, UniDir, ERR] {

	ixMatch = ensureIxMatch(ixMatch)

	return Omni[I, S, T, A, B, ReturnOne, RW, UniDir, ERR](
		func(ctx context.Context, source S) (I, A, error) {
			i, focus, err := get(ctx, source)
			return i, focus, err
		},
		set,
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(ValueIE[I, A]) bool) {
				i, r, err := get(ctx, source)
				yield(ValIE(i, r, JoinCtxErr(ctx, err)))
			}
		},
		func(ctx context.Context, source S) (int, error) {
			return 1, nil
		},
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {
			i, a, err := get(ctx, source)
			if err != nil {
				var ft T
				return ft, err
			}

			b, err := fmap(i, a)
			err = JoinCtxErr(ctx, err)
			if err != nil {
				var ft T
				return ft, err
			}
			ret, err := set(ctx, b, source)
			return ret, err
		},
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			return func(yield func(ValueIE[I, A]) bool) {
				retIx, ret, err := get(ctx, source)
				if err != nil {
					var a A
					yield(ValIE(index, a, err))
					return
				}

				match := ixMatch(retIx, index)

				if match {
					yield(ValIE(index, ret, err))
				}
			}
		},
		ixMatch,
		unsupportedReverseGetter[B, T],
		exprDef,
	)
}

func CombiGetMod[RW, ERR any, I, S, T, A, B any](
	get func(ctx context.Context, source S) (I, A, error),
	modify func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error),
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, T, A, B, ReturnOne, RW, UniDir, ERR] {

	ixMatch = ensureIxMatch(ixMatch)

	return Omni[I, S, T, A, B, ReturnOne, RW, UniDir, ERR](
		func(ctx context.Context, source S) (I, A, error) {
			i, focus, err := get(ctx, source)
			return i, focus, err
		},
		func(ctx context.Context, setFocus B, source S) (T, error) {
			ret, err := modify(ctx, func(index I, focus A) (B, error) {
				return setFocus, nil
			}, source)
			return ret, err
		},
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				i, a, err := get(ctx, source)
				yield(ValIE(i, a, err))
			}
		},
		func(ctx context.Context, source S) (int, error) {
			return 1, nil
		},
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {
			ret, err := modify(ctx, func(index I, focus A) (B, error) {
				return fmap(index, focus)
			}, source)
			return ret, err
		},
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			i, v, err := get(ctx, source)
			if err != nil {
				return func(yield func(val ValueIE[I, A]) bool) {
					yield(ValIE(i, v, err))
				}
			}

			match := ixMatch(i, index)
			if match || err != nil {
				return func(yield func(val ValueIE[I, A]) bool) {
					yield(ValIE(i, v, err))
				}
			} else {
				return func(yield func(val ValueIE[I, A]) bool) {}
			}
		},
		ixMatch,
		unsupportedReverseGetter[B, T],
		exprDef,
	)
}

func CombiIsoMod[RW, DIR, ERR any, I, S, T, A, B any](
	get func(ctx context.Context, source S) (I, A, error),
	modify func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error),
	reverse func(ctx context.Context, focus B) (T, error),
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, T, A, B, ReturnOne, RW, DIR, ERR] {

	ixMatch = ensureIxMatch(ixMatch)

	return Omni[I, S, T, A, B, ReturnOne, RW, DIR, ERR](
		func(ctx context.Context, source S) (I, A, error) {
			i, focus, err := get(ctx, source)
			return i, focus, err
		},
		func(ctx context.Context, setFocus B, source S) (T, error) {
			ret, err := modify(ctx, func(index I, focus A) (B, error) {
				return setFocus, nil
			}, source)
			return ret, err
		},
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				i, a, err := get(ctx, source)
				yield(ValIE(i, a, err))
			}
		},
		func(ctx context.Context, source S) (int, error) {
			return 1, nil
		},
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {
			ret, err := modify(ctx, func(index I, focus A) (B, error) {
				return fmap(index, focus)
			}, source)
			return ret, err
		},
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			i, v, err := get(ctx, source)
			if err != nil {
				return func(yield func(val ValueIE[I, A]) bool) {
					yield(ValIE(i, v, err))
				}
			}

			match := ixMatch(i, index)
			if match || err != nil {
				return func(yield func(val ValueIE[I, A]) bool) {
					yield(ValIE(i, v, err))
				}
			} else {
				return func(yield func(val ValueIE[I, A]) bool) {}
			}
		},
		ixMatch,
		reverse,
		exprDef,
	)
}

func CombiGetter[ERR any, I, S, T, A, B any](
	get func(ctx context.Context, source S) (I, A, error),
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, T, A, B, ReturnOne, ReadOnly, UniDir, ERR] {
	return CombiLens[ReadOnly, ERR, I, S, T, A, B](
		get,
		func(ctx context.Context, focus B, source S) (T, error) {
			var t T
			return t, UnsupportedOpticMethod
		},
		ixMatch,
		exprDef,
	)
}

func CombiPrism[RW, DIR, ERR any, I, S, T, A, B any](
	match func(ctx context.Context, source S) (mo.Either[T, ValueI[I, A]], error),
	embed func(ctx context.Context, focus B) (T, error),
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, T, A, B, ReturnMany, RW, DIR, ERR] {

	if ixMatch == nil {
		ixMatch = ensureIxMatch(ixMatch)
	}

	return Omni[I, S, T, A, B, ReturnMany, RW, DIR, ERR](
		func(ctx context.Context, source S) (index I, ret A, err error) {
			either, err := match(ctx, source)
			if err != nil {
				var i I
				var r A
				return i, r, err
			}
			if a, ok := either.Right(); ok {
				return a.index, a.value, nil
			} else {
				var i I
				var r A
				return i, r, ErrEmptyGet
			}
		},
		func(ctx context.Context, focus B, source S) (T, error) {
			either, err := match(ctx, source)
			if err != nil {
				var t T
				return t, err
			}
			if either.IsRight() {
				ret, err := embed(ctx, focus)
				return ret, err
			} else {
				t := either.MustLeft()
				return t, nil
			}
		},
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(ValueIE[I, A]) bool) {
				either, err := match(ctx, source)
				if err != nil {
					var i I
					var a A
					yield(ValIE(i, a, err))
					return
				}
				if a, ok := either.Right(); ok {
					yield(ValIE(a.index, a.value, nil))
					return
				}
			}
		},
		func(ctx context.Context, source S) (int, error) {
			either, err := match(ctx, source)
			if either.IsRight() {
				return 1, err
			} else {
				return 0, err
			}
		},
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {
			e, err := match(ctx, source)
			if err != nil {
				var t T
				return t, err
			}

			if fa, ok := e.Right(); ok {
				fb, err := fmap(fa.index, fa.value)
				err = JoinCtxErr(ctx, err)
				if err != nil {
					var ft T
					return ft, err
				}

				ft, err := embed(ctx, fb)
				return ft, err
			} else {
				t := e.MustLeft()
				return t, err
			}
		},
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			return func(yield func(ValueIE[I, A]) bool) {
				either, err := match(ctx, source)
				if err != nil {
					var a A
					yield(ValIE(index, a, err))
					return
				}

				if a, ok := either.Right(); ok {
					match := ixMatch(index, a.index)
					if match {
						yield(ValIE(a.index, a.value, err))
					}
					return
				}
			}
		},
		ixMatch,
		embed,
		exprDef,
	)
}

func CombiIteration[RET, ERR, I, S, T, A, B any](
	iterate func(ctx context.Context, source S) SeqIE[I, A],
	lengthGetter LengthGetterFunc[S],
	ixget func(ctx context.Context, index I, source S) SeqIE[I, A],
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, T, A, B, RET, ReadOnly, UniDir, ERR] {

	return CombiTraversal[RET, ReadOnly, ERR, I, S, T, A, B](
		iterate,
		lengthGetter,
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {
			var t T
			return t, UnsupportedOpticMethod
		},
		ixget,
		ixMatch,
		exprDef,
	)
}

func CombiIso[RW, DIR, ERR, S, T, A, B any](
	getter func(ctx context.Context, source S) (A, error),
	reverse func(ctx context.Context, focus B) (T, error),
	exprDef ExpressionDef,
) Optic[Void, S, T, A, B, ReturnOne, RW, BiDir, ERR] {

	return Omni[Void, S, T, A, B, ReturnOne, RW, BiDir, ERR](
		func(ctx context.Context, source S) (Void, A, error) {
			ret, err := getter(ctx, source)
			return Void{}, ret, err
		},
		func(ctx context.Context, v B, s S) (T, error) {
			ret, err := reverse(ctx, v)
			return ret, err
		},
		func(ctx context.Context, source S) SeqIE[Void, A] {
			return func(yield func(ValueIE[Void, A]) bool) {
				v, err := getter(ctx, source)
				yield(ValIE(Void{}, v, err))
			}
		},
		func(ctx context.Context, source S) (int, error) {
			return 1, nil
		},
		func(ctx context.Context, fmap func(index Void, focus A) (B, error), source S) (T, error) {
			fa, err := getter(ctx, source)
			if err != nil {
				var t T
				return t, err
			}

			fb, err := fmap(Void{}, fa)
			err = JoinCtxErr(ctx, err)
			if err != nil {
				var t T
				return t, err
			}
			ft, err := reverse(ctx, fb)
			return ft, err
		},
		func(ctx context.Context, index Void, source S) SeqIE[Void, A] {
			return func(yield func(ValueIE[Void, A]) bool) {
				ret, err := getter(ctx, source)
				yield(ValIE(index, ret, err))
			}
		},
		IxMatchVoid(),
		func(ctx context.Context, focus B) (T, error) {
			ret, err := reverse(ctx, focus)
			return ret, err
		},
		exprDef,
	)
}
