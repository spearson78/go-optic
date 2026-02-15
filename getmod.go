package optic

import "context"

// Constructor for a [Lens] like optic with a [Traversal] style modify function.
// The modify function provides access to the old and new focus values at the same time. This can be used to detect which modifications were made to the focus.
//
// The following constructors are available.
//   - [GetMod]
//   - [GetModP] polymorphic
//   - [GetModE] error aware
//   - [GetModEP] polymorphic, error aware
//   - [GetModI] indexed
//   - [GetModIP] indexed,polymorphic
//   - [GetModIE] indexed,error aware
//   - [GetModIEP] indexed,polymorphic, error aware
//
// An GetMod is constructed from 3 functions
//
//   - get     : should return the single value focused by the lens
//   - set     : should return a new result type with the new given focused value.
//   - expr    : should return the expression type. See the [expr] package for more information.
func GetMod[S, A any](
	get func(S) A,
	modify func(fmap func(focus A) A, source S) S,
	exprDef ExpressionDef,
) Optic[Void, S, S, A, A, ReturnOne, ReadWrite, UniDir, Pure] {
	return GetModP[S, S, A, A](get, modify, exprDef)
}

// Constructor for a [Lens] like optic with a [Traversal] style modify function.
// The modify function provides access to the old and new focus values at the same time. This can be used to detect which modifications were made to the focus.
//
// The following constructors are available.
//   - [GetMod]
//   - [GetModP] polymorphic
//   - [GetModE] error aware
//   - [GetModEP] polymorphic, error aware
//   - [GetModI] indexed
//   - [GetModIP] indexed,polymorphic
//   - [GetModIE] indexed,error aware
//   - [GetModIEP] indexed,polymorphic, error aware
//
// An GetMod is constructed from 3 functions
//
//   - get     : should return the single value focused by the lens
//   - set     : should return a new result type with the new given focused value.
//   - expr    : should return the expression type. See the [expr] package for more information.
func GetModP[S, T, A, B any](
	get func(S) A,
	modify func(fmap func(focus A) B, source S) T,
	exprDef ExpressionDef,
) Optic[Void, S, T, A, B, ReturnOne, ReadWrite, UniDir, Pure] {

	return CombiGetMod[ReadWrite, Pure, Void, S, T, A, B](
		func(ctx context.Context, source S) (Void, A, error) {
			a := get(source)
			return Void{}, a, nil
		},
		func(ctx context.Context, fmap func(index Void, focus A) (B, error), source S) (ret T, retErr error) {
			defer handleAbortModify(&retErr)
			ret = modify(func(focus A) B {
				b, err := fmap(Void{}, focus)
				abortModifyError(JoinCtxErr(ctx, err), &retErr)
				return b
			}, source)
			return
		},
		IxMatchVoid(),
		exprDef,
	)
}

// Constructor for a [Lens] like optic with a [Traversal] style modify function.
// The modify function provides access to the old and new focus values at the same time. This can be used to detect which modifications were made to the focus.
//
// The following constructors are available.
//   - [GetMod]
//   - [GetModP] polymorphic
//   - [GetModE] error aware
//   - [GetModEP] polymorphic, error aware
//   - [GetModI] indexed
//   - [GetModIP] indexed,polymorphic
//   - [GetModIE] indexed,error aware
//   - [GetModIEP] indexed,polymorphic, error aware
//
// An GetMod is constructed from 3 functions
//
//   - get     : should return the single value focused by the lens
//   - set     : should return a new result type with the new given focused value.
//   - expr    : should return the expression type. See the [expr] package for more information.
func GetModE[S, A any](
	get func(ctx context.Context, source S) (A, error),
	modify func(ctx context.Context, fmap func(focus A) (A, error), source S) (S, error),
	exprDef ExpressionDef,
) Optic[Void, S, S, A, A, ReturnOne, ReadWrite, UniDir, Err] {
	return GetModEP(get, modify, exprDef)
}

// Constructor for a [Lens] like optic with a [Traversal] style modify function.
// The modify function provides access to the old and new focus values at the same time. This can be used to detect which modifications were made to the focus.
//
// The following constructors are available.
//   - [GetMod]
//   - [GetModP] polymorphic
//   - [GetModE] error aware
//   - [GetModEP] polymorphic, error aware
//   - [GetModI] indexed
//   - [GetModIP] indexed,polymorphic
//   - [GetModIE] indexed,error aware
//   - [GetModIEP] indexed,polymorphic, error aware
//
// An GetMod is constructed from 3 functions
//
//   - get     : should return the single value focused by the lens
//   - set     : should return a new result type with the new given focused value.
//   - expr    : should return the expression type. See the [expr] package for more information.
func GetModEP[S, T, A, B any](
	get func(ctx context.Context, source S) (A, error),
	modify func(ctx context.Context, fmap func(focus A) (B, error), source S) (T, error),
	exprDef ExpressionDef,
) Optic[Void, S, T, A, B, ReturnOne, ReadWrite, UniDir, Err] {

	return CombiGetMod[ReadWrite, Err, Void, S, T, A, B](
		func(ctx context.Context, source S) (Void, A, error) {
			a, err := get(ctx, source)
			return Void{}, a, err
		},
		func(ctx context.Context, fmap func(index Void, focus A) (B, error), source S) (T, error) {
			return modify(ctx, func(focus A) (B, error) {
				return fmap(Void{}, focus)
			}, source)
		},
		IxMatchVoid(),
		exprDef,
	)
}

func GetModI[I, S, A any](
	get func(source S) (I, A),
	modify func(fmap func(index I, focus A) A, source S) S,
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, S, A, A, ReturnOne, ReadWrite, UniDir, Err] {
	return GetModIP(get, modify, ixMatch, exprDef)
}

// Constructor for a [Lens] like optic with a [Traversal] style modify function.
// The modify function provides access to the old and new focus values at the same time. This can be used to detect which modifications were made to the focus.
//
// The following constructors are available.
//   - [GetMod]
//   - [GetModP] polymorphic
//   - [GetModE] error aware
//   - [GetModEP] polymorphic, error aware
//   - [GetModI] indexed
//   - [GetModIP] indexed,polymorphic
//   - [GetModIE] indexed,error aware
//   - [GetModIEP] indexed,polymorphic, error aware
//
// An indexed GetMod is constructed from 4 functions
//
//   - get     : should return the single value focused by the lens
//   - set     : should return a new result type with the new given focused value.
//   - ixMatch : should return true if the 2 index values are equal
//   - expr    : should return the expression type. See the [expr] package for more information.
func GetModIP[I, S, T, A, B any](
	get func(source S) (I, A),
	modify func(fmap func(index I, focus A) B, source S) T,
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, T, A, B, ReturnOne, ReadWrite, UniDir, Err] {

	return CombiGetMod[ReadWrite, Err, I, S, T, A, B](
		func(ctx context.Context, source S) (I, A, error) {
			i, a := get(source)
			return i, a, nil
		},
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (ret T, retErr error) {
			defer handleAbortModify(&retErr)
			ret = modify(func(index I, focus A) B {
				b, err := fmap(index, focus)
				abortModifyError(JoinCtxErr(ctx, err), &retErr)
				return b
			}, source)
			return
		},
		ixMatch,
		exprDef,
	)
}

// Constructor for a [Lens] like optic with a [Traversal] style modify function.
// The modify function provides access to the old and new focus values at the same time. This can be used to detect which modifications were made to the focus.
//
// The following constructors are available.
//   - [GetMod]
//   - [GetModP] polymorphic
//   - [GetModE] error aware
//   - [GetModEP] polymorphic, error aware
//   - [GetModI] indexed
//   - [GetModIP] indexed,polymorphic
//   - [GetModIE] indexed,error aware
//   - [GetModIEP] indexed,polymorphic, error aware
//
// An indexed GetMod is constructed from 4 functions
//
//   - get     : should return the single value focused by the lens
//   - set     : should return a new result type with the new given focused value.
//   - ixMatch : should return true if the 2 index values are equal
//   - expr    : should return the expression type. See the [expr] package for more information.
func GetModIE[I, S, A any](
	get func(ctx context.Context, source S) (I, A, error),
	modify func(ctx context.Context, fmap func(index I, focus A) (A, error), source S) (S, error),
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, S, A, A, ReturnOne, ReadWrite, UniDir, Err] {
	return GetModIEP(get, modify, ixMatch, exprDef)
}

// Constructor for a [Lens] like optic with a [Traversal] style modify function.
// The modify function provides access to the old and new focus values at the same time. This can be used to detect which modifications were made to the focus.
//
// The following constructors are available.
//   - [GetMod]
//   - [GetModP] polymorphic
//   - [GetModE] error aware
//   - [GetModEP] polymorphic, error aware
//   - [GetModI] indexed
//   - [GetModIP] indexed,polymorphic
//   - [GetModIE] indexed,error aware
//   - [GetModIEP] indexed,polymorphic, error aware
//
// An indexed GetMod is constructed from 4 functions
//
//   - get     : should return the single value focused by the lens
//   - set     : should return a new result type with the new given focused value.
//   - ixMatch : should return true if the 2 index values are equal
//   - expr    : should return the expression type. See the [expr] package for more information.
func GetModIEP[I, S, T, A, B any](
	get func(ctx context.Context, source S) (I, A, error),
	modify func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error),
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, T, A, B, ReturnOne, ReadWrite, UniDir, Err] {

	return CombiGetMod[ReadWrite, Err, I, S, T, A, B](
		get,
		modify,
		ixMatch,
		exprDef,
	)
}
