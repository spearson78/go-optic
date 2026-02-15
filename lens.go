package optic

import (
	"context"
)

// Constructor for a Lens optic. Lenses focus on a single element.
//
// The following Lens constructors are available.
//   - [Lens]
//   - [LensP] polymorphic
//   - [LensE] error aware
//   - [LensEP] polymorphic, error aware
//   - [LensI] indexed
//   - [LensIP] indexed,polymorphic
//   - [LensIE] indexed,error aware
//   - [LensIEP] indexed,polymorphic, error aware
//
// A lens is constructed from 3 functions
//
//   - get : should return the single value focused by the lens
//   - set : should return a new result type with the new given focused value.
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [GetModIEP] for a lens like optic that uses a modify function instead of a setter.
func Lens[S, A any](
	get func(source S) A,
	set func(focus A, source S) S,
	exprDef ExpressionDef,
) Optic[Void, S, S, A, A, ReturnOne, ReadWrite, UniDir, Pure] {
	return LensP[S, S, A, A](
		get,
		set,
		exprDef,
	)
}

// Constructor for a Lens optic. Lenses focus on a single element.
//
// The following Lens constructors are available.
//   - [Lens]
//   - [LensP] polymorphic
//   - [LensE] error aware
//   - [LensEP] polymorphic, error aware
//   - [LensI] indexed
//   - [LensIP] indexed,polymorphic
//   - [LensIE] indexed,error aware
//   - [LensIEP] indexed,polymorphic, error aware
//
// A lens is constructed from 3 functions
//
//   - get : should return the single value focused by the lens
//   - set : should return a new result type with the new given focused value.
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [GetModIEP] for a lens like optic that uses a modify function instead of a setter.
func LensP[S, T, A, B any](
	get func(source S) A,
	set func(focus B, source S) T,
	exprDef ExpressionDef,
) Optic[Void, S, T, A, B, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, Void, S, T, A, B](
		func(ctx context.Context, s S) (Void, A, error) {
			iv := get(s)
			return Void{}, iv, ctx.Err()
		},
		func(ctx context.Context, v B, s S) (T, error) {
			return set(v, s), ctx.Err()
		},
		IxMatchVoid(),
		exprDef,
	)
}

// Constructor for a Lens optic. Lenses focus on a single element.
//
// The following Lens constructors are available.
//   - [Lens]
//   - [LensP] polymorphic
//   - [LensE] error aware
//   - [LensEP] polymorphic, error aware
//   - [LensI] indexed
//   - [LensIP] indexed,polymorphic
//   - [LensIE] indexed,error aware
//   - [LensIEP] indexed,polymorphic, error aware
//
// A lens is constructed from 3 functions
//
//   - get : should return the single value focused by the lens
//   - set : should return a new result type with the new given focused value.
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [GetModIEP] for a lens like optic that uses a modify function instead of a setter.
func LensE[S, A any](
	get func(ctx context.Context, source S) (A, error),
	set func(ctx context.Context, focus A, source S) (S, error),
	exprDef ExpressionDef,
) Optic[Void, S, S, A, A, ReturnOne, ReadWrite, UniDir, Err] {
	return CombiLens[ReadWrite, Err, Void, S, S, A, A](
		func(ctx context.Context, s S) (Void, A, error) {
			a, err := get(ctx, s)
			return Void{}, a, JoinCtxErr(ctx, err)
		},
		func(ctx context.Context, focus A, source S) (S, error) {
			ret, err := set(ctx, focus, source)
			return ret, JoinCtxErr(ctx, err)
		},
		IxMatchVoid(),
		exprDef,
	)
}

// Constructor for a Lens optic. Lenses focus on a single element.
//
// The following Lens constructors are available.
//   - [Lens]
//   - [LensP] polymorphic
//   - [LensE] error aware
//   - [LensEP] polymorphic, error aware
//   - [LensI] indexed
//   - [LensIP] indexed,polymorphic
//   - [LensIE] indexed,error aware
//   - [LensIEP] indexed,polymorphic, error aware
//
// A lens is constructed from 3 functions
//
//   - get : should return the single value focused by the lens
//   - set : should return a new result type with the new given focused value.
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [GetModIEP] for a lens like optic that uses a modify function instead of a setter.
func LensEP[S, T, A, B any](
	get func(ctx context.Context, source S) (A, error),
	set func(ctx context.Context, focus B, source S) (T, error),
	exprDef ExpressionDef,
) Optic[Void, S, T, A, B, ReturnOne, ReadWrite, UniDir, Err] {
	return CombiLens[ReadWrite, Err, Void, S, T, A, B](
		func(ctx context.Context, s S) (Void, A, error) {
			a, err := get(ctx, s)
			return Void{}, a, JoinCtxErr(ctx, err)
		},
		func(ctx context.Context, focus B, source S) (T, error) {
			ret, err := set(ctx, focus, source)
			return ret, JoinCtxErr(ctx, err)
		},
		IxMatchVoid(),
		exprDef,
	)
}

// Constructor for a Lens optic. Lenses focus on a single element.
//
// The following Lens constructors are available.
//   - [Lens]
//   - [LensP] polymorphic
//   - [LensE] error aware
//   - [LensEP] polymorphic, error aware
//   - [LensI] indexed
//   - [LensIP] indexed,polymorphic
//   - [LensIE] indexed,error aware
//   - [LensIEP] indexed,polymorphic, error aware
//
// An indexed lens is constructed from 4 functions
//
//   - get     : should return the single value focused by the lens
//   - set     : should return a new result type with the new given focused value.
//   - ixMatch : should return true if the 2 index values are equal
//   - expr    : should return the expression type. See the [expr] package for more information.
//
// See:
//   - [GetModIEP] for a lens like optic that uses a modify function instead of a setter.
func LensI[I, S, A any](
	get func(source S) (I, A),
	set func(focus A, source S) S,
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, S, A, A, ReturnOne, ReadWrite, UniDir, Pure] {

	ixMatchFnc := ensureSimpleIxMatch(ixMatch)

	return CombiLens[ReadWrite, Pure, I, S, S, A, A](
		func(ctx context.Context, s S) (I, A, error) {
			i, v := get(s)
			return i, v, ctx.Err()
		},
		func(ctx context.Context, v A, s S) (S, error) {
			return set(v, s), ctx.Err()
		},
		ixMatchFnc,
		exprDef,
	)
}

// Constructor for a Lens optic. Lenses focus on a single element.
//
// The following Lens constructors are available.
//   - [Lens]
//   - [LensP] polymorphic
//   - [LensE] error aware
//   - [LensEP] polymorphic, error aware
//   - [LensI] indexed
//   - [LensIP] indexed,polymorphic
//   - [LensIE] indexed,error aware
//   - [LensIEP] indexed,polymorphic, error aware
//
// An indexed lens is constructed from 4 functions
//
//   - get     : should return the single value focused by the lens
//   - set     : should return a new result type with the new given focused value.
//   - ixMatch : should return true if the 2 index values are equal
//   - expr    : should return the expression type. See the [expr] package for more information.
//
// See:
//   - [GetModIEP] for a lens like optic that uses a modify function instead of a setter.
func LensIP[I, S, T, A, B any](
	get func(source S) (I, A),
	set func(focus B, source S) T,
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, T, A, B, ReturnOne, ReadWrite, UniDir, Pure] {
	return CombiLens[ReadWrite, Pure, I, S, T, A, B](
		func(ctx context.Context, s S) (I, A, error) {
			i, iv := get(s)
			return i, iv, ctx.Err()
		},
		func(ctx context.Context, v B, s S) (T, error) {
			return set(v, s), ctx.Err()
		},
		ixMatch,
		exprDef,
	)
}

// Constructor for a Lens optic. Lenses focus on a single element.
//
// The following Lens constructors are available.
//   - [Lens]
//   - [LensP] polymorphic
//   - [LensE] error aware
//   - [LensEP] polymorphic, error aware
//   - [LensI] indexed
//   - [LensIP] indexed,polymorphic
//   - [LensIE] indexed,error aware
//   - [LensIEP] indexed,polymorphic, error aware
//
// An indexed lens is constructed from 4 functions
//
//   - get     : should return the single value focused by the lens
//   - set     : should return a new result type with the new given focused value.
//   - ixMatch : should return true if the 2 index values are equal
//   - expr    : should return the expression type. See the [expr] package for more information.
//
// See:
//   - [GetModIEP] for a lens like optic that uses a modify function instead of a setter.
func LensIE[I, S, A any](
	get func(ctx context.Context, source S) (I, A, error),
	set func(ctx context.Context, focus A, source S) (S, error),
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, S, A, A, ReturnOne, ReadWrite, UniDir, Err] {
	return CombiLens[ReadWrite, Err, I, S, S, A, A](
		get,
		set,
		ixMatch,
		exprDef,
	)
}

// Constructor for a Lens optic. Lenses focus on a single element.
//
// The following Lens constructors are available.
//   - [Lens]
//   - [LensP] polymorphic
//   - [LensE] error aware
//   - [LensEP] polymorphic, error aware
//   - [LensI] indexed
//   - [LensIP] indexed,polymorphic
//   - [LensIE] indexed,error aware
//   - [LensIEP] indexed,polymorphic, error aware
//
// An indexed lens is constructed from 4 functions
//
//   - get     : should return the single value focused by the lens
//   - set     : should return a new result type with the new given focused value.
//   - ixMatch : should return true if the 2 index values are equal
//   - expr    : should return the expression type. See the [expr] package for more information.
//
// See:
//   - [GetModIEP] for a lens like optic that uses a modify function instead of a setter.
func LensIEP[I any, S, T, A, B any](
	get func(ctx context.Context, source S) (I, A, error),
	set func(ctx context.Context, focus B, source S) (T, error),
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, T, A, B, ReturnOne, ReadWrite, UniDir, Err] {
	return CombiLens[ReadWrite, Err](get, set, ixMatch, exprDef)
}
