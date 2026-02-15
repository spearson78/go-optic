package optic

import (
	"context"
	"reflect"
)

// The Get action returns the single focus of the given optic.
//
// Get only accepts optics that return a single value [ReturnOne].
//
// See:
//   - [GetFirst] for a version the accepts [ReturnMany] optics
//   - [GetI] for an index aware version.
//   - [GetContext] for a context aware version
//   - [MustGet] for a version that panics on error.
func Get[I, S, T, A, B any, RET TReturnOne, RO any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, RO, DIR, ERR], source S) (A, error) {
	return GetContext(context.Background(), o, source)
}

// The GetContext action returns the single focus of the given optic.
//
// GetContext only accepts optics that return a single value [ReturnOne]. For a version the accepts any optic see [GetFirstContext]
//
// See:
//   - [GetFirstContext] for a version the accepts [ReturnMany] optics
//   - [GetContextI] for an index aware version.
//   - [Get] for a non context aware version
//   - [MustGet] for a version that panics on error.
func GetContext[I, S, T, A, B any, RET TReturnOne, RO any, DIR any, ERR any](ctx context.Context, o Optic[I, S, T, A, B, RET, RO, DIR, ERR], source S) (A, error) {
	_, ret, err := GetContextI(ctx, o, source)
	return ret, err
}

// The MustGet action returns the single focus of the given optic and panics on error.
//
// MustGet only accepts optics that return a single value [ReturnOne]. For a version the accepts any optic see [MustGetFirst]
//
// See:
//   - [MustGetFirst] for a version the accepts [ReturnMany] optics
//   - [MustGetI] for an index aware version.
//   - [Get] for a non context aware version
//   - [GetContext] for a context aware version
func MustGet[I, S, T, A, B any, RET TReturnOne, RO any, DIR any, ERR TPure](o Optic[I, S, T, A, B, RET, RO, DIR, ERR], source S) A {
	return Must(Get(o, source))
}

// The ReverseGet action retrieves the single reverse get focus of the given optic.
//
// See:
//   - [ReverseGetContext] for a context aware version
//   - [MustReverseGet] for a version that panics on error.
func ReverseGet[I, S, T, A, B any, RET any, RO any, DIR TBiDir, ERR any](o Optic[I, S, T, A, B, RET, RO, DIR, ERR], focus B) (T, error) {
	return ReverseGetContext(context.Background(), o, focus)
}

// The ReverseGetContext action retrieves the single reverse get focus of the given optic.
//
// See:
//   - [ReverseGet] for a non context aware version
//   - [MustReverseGet] for a version that panics on error.
func ReverseGetContext[I, S, T, A, B any, RET any, RO any, DIR TBiDir, ERR any](ctx context.Context, o Optic[I, S, T, A, B, RET, RO, DIR, ERR], focus B) (T, error) {
	exprHandler, err := getExprHandler(ctx, o.AsExprHandler())
	if err != nil {
		var t T
		return t, err
	}

	if exprHandler != nil {
		source, err := exprHandler.ReverseGet(ctx, o.AsExpr(), focus)

		t, ok := source.(T)
		if !ok {
			var t T
			return t, &ErrCastOf{
				Of:   reflect.TypeOf(source),
				From: reflect.TypeFor[any](),
				To:   reflect.TypeFor[T](),
			}
		}

		return t, err
	} else {
		return o.AsReverseGetter()(ctx, focus)
	}
}

// The MustReverseGet action retrieves the single reverse get focus of the given optic and  panics on error
//
// See:
//   - [ReverseGet] for a non context aware version
//   - [ReverseGetContext] for a context aware version
func MustReverseGet[I, S, T, A, B, RET any, RO any, DIR TBiDir, ERR TPure](o Optic[I, S, T, A, B, RET, RO, DIR, ERR], focus B) T {
	return Must(ReverseGet(o, focus))
}

// The Modify action modifies the source by applying the given operator to every focus.
//
// See:
//   - [ModifyCheck] for a version that reports if the input optic focused no elements
//   - [ModifyI] for an index aware version.
//   - [ModifyContext] for a context aware version.
//   - [MustModify] for a version that panics on error.
func Modify[I, S, T, A, B any, RET any, RW TReadWrite, DIR any, RETF TReturnOne, ERR, ERRFMAP any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], op Operation[A, B, RETF, ERRFMAP], source S) (T, error) {
	return ModifyI(o, OpToOpI[I](op), source)
}

// The ModifyContext action modifies the source by applying the given operator to every focus.
//
// See:
//   - [ModifyCheckContext] for a version that reports if the input optic focused no elements
//   - [ModifyContextI] for an index aware version.
//   - [Modify] for a non context aware version.
//   - [MustModify] for a version that panics on error.
func ModifyContext[I, S, T, A, B any, RET any, RW TReadWrite, DIR any, FRET TReturnOne, ERR, ERRFMAP any](ctx context.Context, o Optic[I, S, T, A, B, RET, RW, DIR, ERR], f Operation[A, B, FRET, ERRFMAP], source S) (T, error) {
	return ModifyContextI(ctx, o, OpToOpI[I](f), source)
}

// The MustModify action modifies the source by applying the given operator to every focus panicking on error.
//
// See:
//   - [MustModifyCheck] for a version that panics on error.
//   - [IMustOver] for an index aware version.
//   - [Modify] for a non context aware version.
//   - [ModifyContext] for a context aware version.
func MustModify[I, S, T, A, B any, RET any, RW TReadWrite, DIR any, RETF TReturnOne, ERR, ERRFMAP TPure](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], f Operation[A, B, RETF, ERRFMAP], source S) T {
	return MustModifyI(o, OpToOpI[I](f), source)
}

// The ModifyCheck action modifies the source by applying the given operator to every focus returning false if the input optic focused no elements
//
// See:
//   - [Modify] for a version that does not report if the input optic focused no elements
//   - [ModifyCheckI] for an index aware version.
//   - [ModifyCheckContext] for a context aware version.
//   - [MustModifyCheck] for a version that panics on error.
func ModifyCheck[I, S, T, A, B any, RET any, RW TReadWrite, DIR any, FRET TReturnOne, ERR, ERRFMAP any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], fmap Operation[A, B, FRET, ERRFMAP], source S) (T, bool, error) {
	return ModifyCheckContext(context.Background(), o, fmap, source)
}

// The ModifyCheckContext action modifies the source by applying the given operator to every focus returning false if the input optic focused no elements.
//
// See:
//   - [ModifyContext] for a version that does not report if the input optic focused no elements
//   - [ModifyCheckContextI] for an index aware version.
//   - [ModifyCheck] for a non context aware version.
//   - [MustModifyCheck] for a version that panics on error.
func ModifyCheckContext[I, S, T, A, B any, RET any, RW TReadWrite, DIR any, FRET TReturnOne, ERR, ERRFMAP any](ctx context.Context, o Optic[I, S, T, A, B, RET, RW, DIR, ERR], fmap Operation[A, B, FRET, ERRFMAP], source S) (T, bool, error) {
	return ModifyCheckContextI(ctx, o, OpToOpI[I](fmap), source)
}

// The MustModifyCheck action modifies the source by applying the given operator to every focus panicking on error and returning false if the input optic focused no elements
//
// See:
//   - [MustModify] for a version that does not report if the input optic focused no elements
//   - [IMustFailOver] for an index aware version.
//   - [ModifyCheck] for a non context aware version.
//   - [ModifyCheckContext] for a context aware version.
func MustModifyCheck[I, S, T, A, B any, RET any, RW TReadWrite, DIR any, FRET TReturnOne, ERR, ERRFMAP TPure](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], fmap Operation[A, B, FRET, ERRFMAP], source S) (T, bool) {
	return Must2(ModifyCheck(o, fmap, source))
}

// The GetFirst action returns the first focus of the optic or false if there are no elements.
//
// See:
//   - [GetFirstI] for an index aware version.
//   - [GetFirstContext] for a context aware version.
//   - [MustGetFirst] for a version that panics on error.
func GetFirst[I, S, T, A, B, RET any, RO, DIR, ERR any](o Optic[I, S, T, A, B, RET, RO, DIR, ERR], source S) (A, bool, error) {
	return GetFirstContext(context.Background(), o, source)
}

// The GetFirstContext action returns the first focus of the optic or false if there are no elements.
//
// See:
//   - [GetFirstContextI] for an index aware version.
//   - [GetFirst] for a non context aware version.
//   - [MustGetFirst] for a version that panics on error.
func GetFirstContext[I, S, T, A, B, RET any, RO, DIR any, ERR any](ctx context.Context, o Optic[I, S, T, A, B, RET, RO, DIR, ERR], source S) (A, bool, error) {
	_, res, ok, err := GetFirstContextI(ctx, o, source)
	return res, ok, JoinCtxErr(ctx, err)
}

// The MustGetFirst action returns the first focus of the optic or false if there are no elements.
//
// See:
//   - [IMustPreView] for an index aware version.
//   - [GetFirst] for a non context aware version.
//   - [GetFirstContext] for a context aware version.
func MustGetFirst[I, S, T, A, B, RET any, RO, DIR any, ERR TPure](o Optic[I, S, T, A, B, RET, RO, DIR, ERR], source S) (A, bool) {
	return Must2(GetFirst(o, source))
}

// The Set action sets the focused elements of the optic to the given value.
//
// See:
//   - [Modify] for a version that applies a mapping function to each focus.
//   - [SetContext] for a context aware version.
//   - [MustSet] for a version that panics on error.
func Set[I, S, T, A, B, RET any, RW TReadWrite, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], val B, source S) (T, error) {
	return SetContext(context.Background(), o, val, source)
}

// The SetContext action sets the focused elements of the optic to the given value.
//
// See:
//   - [ModifyContext] for a version that applies a mapping function to each focus.
//   - [Set] for a non context aware version.
//   - [MustSet] for a version that panics on error.
func SetContext[I, S, T, A, B any, RET any, RW TReadWrite, DIR, ERR any](ctx context.Context, o Optic[I, S, T, A, B, RET, RW, DIR, ERR], val B, source S) (T, error) {
	exprHandler, err := getExprHandler(ctx, o.AsExprHandler())
	if err != nil {
		var t T
		return t, err
	}

	if exprHandler != nil {

		ret, err := exprHandler.Set(ctx, o.AsExpr(), source, val)
		if err != nil {
			var t T
			return t, err
		}

		t, ok := ret.(T)
		if !ok {
			return t, &ErrCastOf{
				Of:   reflect.TypeOf(ret),
				From: reflect.TypeFor[any](),
				To:   reflect.TypeFor[T](),
			}
		}

		return t, nil
	} else {
		ret, err := o.AsSetter()(ctx, val, source)
		return ret, err
	}
}

// The MustSet action sets the focused elements of the optic to the given value and panics on error
//
// See:
//   - [MustModify] for a version that applies a mapping function to each focus.
//   - [Set] for a non context aware version.
//   - [SetContext] for a context aware version.
func MustSet[I, S, T, A, B, RET any, RW TReadWrite, DIR any, ERR TPure](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], val B, source S) T {
	return Must(Set(o, val, source))
}
