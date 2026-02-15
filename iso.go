package optic

import (
	"context"
)

// Constructor for the Involuted case of the [Iso] optic.
// This constructor is pure and non polymorphic
//
// An involuted is an iso where the conversion to and from the result is the same. e.g. negation of a integer.
//
// The following constructors are available.
//   - [Involuted]
//   - [InvolutedE] error aware.
//
// An involuted iso is constructed from 2 functions
//
//   - getter : should convert the source to the result.
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [Iso] for a version where the conversion to and from the result is different.
func Involuted[S any](
	getter func(source S) S,
	exprDef ExpressionDef,
) Optic[Void, S, S, S, S, ReturnOne, ReadWrite, BiDir, Pure] {
	return Iso[S](
		getter,
		getter,
		exprDef,
	)
}

// Constructor for the Involuted case of the [Iso] optic.
// This constructor is pure and non polymorphic
//
// An involuted is an iso where the conversion to and from the result is the same. e.g. negation of a integer.
//
// The following constructors are available.
//   - [Involuted]
//   - [InvolutedE] error aware.
//
// An involuted iso is constructed from 2 functions
//
//   - getter : should convert the source to the result.
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [IsoE] for a version where the conversion to and from the result is different.
func InvolutedE[S any](
	getter func(ctx context.Context, source S) (S, error),
	exprDef ExpressionDef,
) Optic[Void, S, S, S, S, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoE[S](
		getter,
		getter,
		exprDef,
	)
}

// Constructor for an Iso optic. Iso's perform a lossless conversion between 2 types.
// This constructor is pure and non polymorphic
//
// The following constructors are available.
//   - [Iso]
//   - [IsoP] polymorphic
//   - [IsoE] error aware
//   - [IsoEP] polymorphic, error aware
//
// An iso is constructed from 3 functions
//
//   - getter : should convert the source type to the result type.
//   - reverse : should convert the result type back to the source type
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [Involuted] for a version where the getter and reverse are identical.
func Iso[S, A any](
	getter func(source S) A,
	reverse func(focus A) S,
	exprDef ExpressionDef,
) Optic[Void, S, S, A, A, ReturnOne, ReadWrite, BiDir, Pure] {
	return CombiIso[ReadWrite, BiDir, Pure, S, S, A, A](
		func(ctx context.Context, v S) (A, error) {
			return getter(v), nil
		},
		func(ctx context.Context, v A) (S, error) {
			return reverse(v), nil
		},
		exprDef,
	)
}

// Constructor for an Iso optic. Iso's perform a lossless conversion between 2 types.
// This constructor is impure , non indexed and non polymorphic
//
// The following constructors are available.
//   - [Iso]
//   - [IsoP] polymorphic
//   - [IsoE] error aware
//   - [IsoEP] polymorphic, error aware

// An iso is constructed from 3 functions
//
//   - getter : should convert the source type to the result type.
//   - reverse : should convert the result type back to the source type
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [InvolutedR] for a version where the getter and reverse are identical.
func IsoE[S, A any](
	getter func(ctx context.Context, source S) (A, error),
	reverse func(ctx context.Context, focus A) (S, error),
	exprDef ExpressionDef,
) Optic[Void, S, S, A, A, ReturnOne, ReadWrite, BiDir, Err] {
	return CombiIso[ReadWrite, BiDir, Err, S, S, A, A](
		getter,
		reverse,
		exprDef,
	)
}

// Constructor for an Iso optic. Iso's perform a lossless conversion between 2 types.
// This constructor is pure , non indexed and polymorphic
//
// The following constructors are available.
//   - [Iso]
//   - [IsoP] polymorphic
//   - [IsoE] error aware
//   - [IsoEP] polymorphic, error aware
//
// An iso is constructed from 3 functions
//
//   - getter : should convert the source type to the result type.
//   - reverse : should convert the result type back to the source type
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [Involuted] for a version where the getter and reverse are identical.
func IsoP[S, T, A, B any](
	getter func(source S) A,
	reverse func(focus B) T,
	exprDef ExpressionDef,
) Optic[Void, S, T, A, B, ReturnOne, ReadWrite, BiDir, Pure] {
	return CombiIso[ReadWrite, BiDir, Pure, S, T, A, B](
		func(ctx context.Context, v S) (A, error) {
			return getter(v), nil
		},
		func(ctx context.Context, v B) (T, error) {
			return reverse(v), nil
		},
		exprDef,
	)
}

// Constructor for an Iso optic. Iso's perform a lossless conversion between 2 types.
// This constructor is impure , non indexed and polymorphic
//
// The following constructors are available.
//   - [Iso]
//   - [IsoP] polymorphic
//   - [IsoE] error aware
//   - [IsoEP] polymorphic, error aware
//
// An iso is constructed from 3 functions
//
//   - getter : should convert the source type to the result type.
//   - reverse : should convert the result type back to the source type
//   - expr: should return the expression type. See the [expr] package for more information.
//
// See:
//   - [InvolutedR] for a version where the getter and reverse are identical.
func IsoEP[S, T, A, B any](
	getter func(ctx context.Context, source S) (A, error),
	reverse func(ctx context.Context, focus B) (T, error),
	exprDef ExpressionDef,
) Optic[Void, S, T, A, B, ReturnOne, ReadWrite, BiDir, Err] {

	return CombiIso[ReadWrite, BiDir, Err, S, T, A, B](
		getter,
		reverse,
		exprDef,
	)
}
