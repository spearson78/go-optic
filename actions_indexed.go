package optic

import (
	"context"
	"errors"
	"reflect"
)

// The GetI action returns the single indexed focus of the given optic.
//
// GetI only accepts optics that return a single value [ReturnOne].
//
// See:
//   - [GetFirstI] for a version the accepts [ReturnMany] optics
//   - [Get] for a non index aware version.
//   - [GetContextI] for a context aware version
//   - [MustGetI] for a version that panics on error.
func GetI[I, S, A any, RET TReturnOne, RO any, DIR any, ERR any](o OpticRO[I, S, A, RET, RO, DIR, ERR], source S) (I, A, error) {
	return GetContextI(context.Background(), o, source)
}

// The GetContextI action returns the single index and focus of the given optic.
//
// GetContextI only accepts optics that return a single value [ReturnOne]. For a version the accepts any optic see [GetFirstContextI]
//
// See:
//   - [GetFirstContextI] for a version the accepts [ReturnMany] optics
//   - [GetContext] for a non index aware version.
//   - [GetI] for a non context aware version
//   - [MustGetI] for a version that panics on error.
func GetContextI[I, S, A any, RET TReturnOne, RO any, DIR any, ERR any](ctx context.Context, o OpticRO[I, S, A, RET, RO, DIR, ERR], source S) (I, A, error) {
	exprHandler, err := getExprHandler(ctx, o.AsExprHandler())
	if err != nil {
		var i I
		var a A
		return i, a, err
	}

	if exprHandler != nil {

		index, focus, found, err := exprHandler.Get(ctx, o.AsExpr(), source)
		if err != nil {
			var i I
			var a A
			return i, a, err
		}

		if !found {
			var i I
			var a A
			return i, a, ErrEmptyGet
		}

		i, ok := index.(I)
		if !ok {
			var i I
			var a A
			return i, a, &ErrCastOf{
				Of:   reflect.TypeOf(index),
				From: reflect.TypeFor[any](),
				To:   reflect.TypeFor[I](),
			}
		}

		a, ok := focus.(A)
		if !ok {
			var a A
			return i, a, &ErrCastOf{
				Of:   reflect.TypeOf(focus),
				From: reflect.TypeFor[any](),
				To:   reflect.TypeFor[A](),
			}
		}

		return i, a, err
	} else {
		return o.AsGetter()(ctx, source)
	}
}

// The IMustGet action returns the single index and focus of the given optic and panics on error.
//
// IMustGet only accepts optics that return a single value [ReturnOne]. For a version the accepts any optic see [IMustGetFirst]
//
// See:
//   - [IMustGetFirst] for a version the accepts [ReturnMany] optics
//   - [MustGet] for a non index aware version.
//   - [GetI] for a non context aware version
//   - [GetContextI] for a context aware version
func MustGetI[I, S, A any, RET TReturnOne, RO any, DIR any, ERR TPure](o OpticRO[I, S, A, RET, RO, DIR, ERR], source S) (I, A) {
	return Must2(GetI(o, source))
}

// The ModifyI action modifies the source by applying the given indexed operator to every focus.
//
// See:
//   - [ModifyCheckI] for a version that reports if the input optic focused no elements
//   - [Modify] for a non index aware version.
//   - [ModifyContextI] for a context aware version.
//   - [MustModifyI] for a version that panics on error.
func ModifyI[I, S, T, A, B any, RET any, RW TReadWrite, DIR any, RETF TReturnOne, ERR any, ERRFMAP any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], f OperationI[I, A, B, RETF, ERRFMAP], source S) (T, error) {
	return ModifyContextI(context.Background(), o, f, source)
}

// The ModifyContextI action modifies the source by applying the given indexed operator to every focus.
//
// See:
//   - [ModifyCheckContextI] for a version that reports if the input optic focused no elements
//   - [ModifyContext] for a non index aware version.
//   - [ModifyI] for a non context aware version.
//   - [MustModifyI] for a version that panics on error.
func ModifyContextI[I, S, T, A, B any, RET any, RW TReadWrite, DIR any, FRET TReturnOne, ERR any, ERRFMAP any](ctx context.Context, o Optic[I, S, T, A, B, RET, RW, DIR, ERR], f OperationI[I, A, B, FRET, ERRFMAP], source S) (T, error) {
	ret, _, err := ModifyCheckContextI(ctx, o, f, source)
	return ret, err
}

// The MustModifyI action modifies the source by applying the given indexed operator to every focus panicking on error.
//
// See:
//   - [MustModifyCheckI] for a version that reports if the input optic focused no elements
//   - [MustModify] for a non index aware version.
//   - [ModifyI] for a non context aware version.
//   - [ModifyContextI] for a context aware version.
func MustModifyI[I, S, T, A, B any, RET any, RW TReadWrite, DIR any, FRET TReturnOne, ERR TPure, ERRFMAP TPure](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], f OperationI[I, A, B, FRET, ERRFMAP], source S) T {
	return Must(ModifyI(o, f, source))
}

// The ModifyCheckI action modifies the source by applying the given indexed operator to every focus returning false if the input optic focused no elements
//
// See:
//   - [ModifyI] for a version that does not report if the input optic focused no elements
//   - [ModifyCheck] for a non index aware version.
//   - [ModifyCheckContextI] for a context aware version.
//   - [MustModifyCheckI] for a version that panics on error.
func ModifyCheckI[I, S, T, A, B any, RET any, RW TReadWrite, DIR any, FRET TReturnOne, ERR, ERRF any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], fmap OperationI[I, A, B, FRET, ERRF], source S) (T, bool, error) {
	return ModifyCheckContextI(context.Background(), o, fmap, source)
}

// The ModifyCheckContextI action modifies the source by applying the given indexed operator to every focus.
//
// See:
//   - [ModifyContextI] for a version that does not report if the input optic focused no elements
//   - [ModifyCheckContext] for a non index aware version.
//   - [ModifyCheckI] for a non context aware version.
//   - [MustModifyCheckI] for a version that panics on error.
func ModifyCheckContextI[I, S, T, A, B any, RET any, RW TReadWrite, DIR any, FRET TReturnOne, ERR, ERRF any](ctx context.Context, o Optic[I, S, T, A, B, RET, RW, DIR, ERR], fmap OperationI[I, A, B, FRET, ERRF], source S) (T, bool, error) {
	exprHandler, err := getExprHandler(ctx, o.AsExprHandler())
	if err != nil {
		var t T
		return t, false, err
	}

	if exprHandler != nil {

		fmapFnc := fmap.AsOpGet()

		ret, ok, err := exprHandler.Modify(ctx, o.AsExpr(), fmap.AsExpr(), func(index, focus any, focusErr error) (any, error) {
			if focusErr != nil {
				return nil, focusErr
			}

			i, ok := index.(I)
			if !ok {
				return nil, &ErrCastOf{
					Of:   reflect.TypeOf(index),
					From: reflect.TypeFor[any](),
					To:   reflect.TypeFor[I](),
				}
			}

			f, ok := focus.(A)
			if !ok {
				return nil, &ErrCastOf{
					Of:   reflect.TypeOf(focus),
					From: reflect.TypeFor[any](),
					To:   reflect.TypeFor[A](),
				}
			}

			ret, err := fmapFnc(ctx, ValI(i, f))
			return ret, err
		}, source)
		if err != nil {
			var t T
			return t, false, err
		}

		t, ok := ret.(T)
		if !ok {
			return t, false, &ErrCastOf{
				Of:   reflect.TypeOf(ret),
				From: reflect.TypeFor[any](),
				To:   reflect.TypeFor[T](),
			}
		}

		return t, ok, nil
	} else {
		ok := false

		opFnc := fmap.AsOpGet()

		ret, err := o.AsModify()(ctx, func(focusIndex I, focus A) (B, error) {
			ret, err := opFnc(ctx, ValI(focusIndex, focus))
			err = JoinCtxErr(ctx, err)
			if err == nil {
				ok = true
			}
			return ret, err
		}, source)
		return ret, ok, err
	}
}

// The IMustFailOver action modifies the source by applying the given indexed operator to every focus panicking on error.
//
// See:
//   - [MustModifyI] for a version that does not report if the input optic focused no elements
//   - [MustModifyCheck] for a non index aware version.
//   - [ModifyCheckI] for a non context aware version.
//   - [ModifyCheckContextI] for a context aware version.
func MustModifyCheckI[I, S, T, A, B any, RET any, RW TReadWrite, DIR any, FRET TReturnOne, ERR, ERRF TPure](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], fmap OperationI[I, A, B, FRET, ERRF], source S) (T, bool) {
	return Must2(ModifyCheckI(o, fmap, source))
}

// The GetFirstI action returns the first focus and index of the optic or false if there are no elements.
//
// See:
//   - [GetFirst] for a non index aware version.
//   - [GetFirstContextI] for a context aware version.
//   - [MustGetFirstI] for a version that panics on error.
func GetFirstI[I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], source S) (I, A, bool, error) {
	return GetFirstContextI(context.Background(), o, source)
}

// The GetFirstContextI action returns the first focus and index of the optic or false if there are no elements.
//
// See:
//   - [GetFirstContext] for a non index aware version.
//   - [GetFirstI] for a non context aware version.
//   - [MustGetFirstI] for a version that panics on error.
func GetFirstContextI[I, S, T, A, B, RET any, RO any, DIR any, ERR any](ctx context.Context, o Optic[I, S, T, A, B, RET, RO, DIR, ERR], source S) (I, A, bool, error) {
	exprHandler, err := getExprHandler(ctx, o.AsExprHandler())
	if err != nil {
		var i I
		var a A
		return i, a, false, err
	}

	if exprHandler != nil {
		index, focus, found, err := exprHandler.Get(ctx, o.AsExpr(), source)
		if err != nil {
			var i I
			var a A
			return i, a, false, err
		}

		if !found {
			var i I
			var a A
			return i, a, false, nil
		}

		i, ok := index.(I)
		if !ok {
			var i I
			var a A
			return i, a, found, &ErrCastOf{
				Of:   reflect.TypeOf(index),
				From: reflect.TypeFor[any](),
				To:   reflect.TypeFor[I](),
			}
		}

		a, ok := focus.(A)
		if !ok {
			var a A
			return i, a, found, &ErrCastOf{
				Of:   reflect.TypeOf(focus),
				From: reflect.TypeFor[any](),
				To:   reflect.TypeFor[A](),
			}
		}

		return i, a, found, err
	} else {
		i, a, err := o.AsGetter()(ctx, source)
		if errors.Is(err, ErrEmptyGet) {
			return i, a, false, nil
		}
		if err != nil {
			return i, a, false, err
		}
		return i, a, true, nil
	}
}

// The MustGetFirstI action returns the first focus and index of the optic or false if there are no elements.
//
// See:
//   - [MustGetFirst] for a non index aware version.
//   - [GetFirstI] for a non context aware version.
//   - [GetFirstContextI] for a context aware version.
func MustGetFirstI[I, S, T, A, B, RET, RW, DIR any, ERR TPure](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], source S) (I, A, bool) {
	return Must2I(GetFirstI(o, source))
}
