package optic

import (
	"context"

	"github.com/samber/mo"
)

// Constructor for a Prism optic. Prisms focus on zero or 1 element.
//
// The following Traversal constructors are available.
//   - [Prism]
//   - [PrismP] polymorphic
//   - [PrismE] error aware
//   - [PrismEP] polymorphic, error aware
//   - [PrismI] indexed
//   - [PrismIP] indexed, polymorphic
//   - [PrismIE] indexed, error aware
//   - [PrismIEP] indexed, polymorphic, error aware
//
// A prism is constructed from 3 functions
//
//   - match : performs a pattern match and returns the match or false.
//   - embed : converts the focus to the source type. This is equivalent to a [ReverseGetterFunc]
//   - expr: should return the expression type. See the [expr] package for more information.
//
// When traversed the prism will return no focuses if the match fails.
// When modified the prism will return the original source if the match fails.
func Prism[S, A any](
	match func(source S) (A, bool),
	embed func(focus A) S,
	exprDef ExpressionDef,
) Optic[Void, S, S, A, A, ReturnMany, ReadWrite, BiDir, Pure] {
	return CombiPrism[ReadWrite, BiDir, Pure](
		func(ctx context.Context, source S) (mo.Either[S, ValueI[Void, A]], error) {

			if right, ok := match(source); ok {
				return mo.Right[S](ValI(Void{}, right)), ctx.Err()
			} else {
				return mo.Left[S, ValueI[Void, A]](source), ctx.Err()
			}
		},
		func(ctx context.Context, focus A) (S, error) {
			return embed(focus), ctx.Err()
		},
		IxMatchVoid(),
		exprDef,
	)
}

// Constructor for a Prism optic. Prisms focus on zero or 1 element.
//
// The following Traversal constructors are available.
//   - [Prism]
//   - [PrismP] polymorphic
//   - [PrismE] error aware
//   - [PrismEP] polymorphic, error aware
//   - [PrismI] indexed
//   - [PrismIP] indexed, polymorphic
//   - [PrismIE] indexed, error aware
//   - [PrismIEP] indexed, polymorphic, error aware
//
// A prism is constructed from 3 functions
//
//   - match : performs a pattern match and returns either the polymorphic result type or the matched value.
//   - embed : converts the focus to the result type. This is equivalent to a [ReverseGetterFunc]
//   - expr: should return the expression type. See the [expr] package for more information.
//
// When traversed the prism will return no focuses if the match fails.
// When modified the prism will return the original source if the match fails.
func PrismP[S, T, A, B any](
	match func(source S) mo.Either[T, A],
	embed func(focus B) T,
	exprDef ExpressionDef,
) Optic[Void, S, T, A, B, ReturnMany, ReadWrite, BiDir, Pure] {
	return CombiPrism[ReadWrite, BiDir, Pure](
		func(ctx context.Context, source S) (mo.Either[T, ValueI[Void, A]], error) {
			e := match(source)
			if right, ok := e.Right(); ok {
				return mo.Right[T](ValI(Void{}, right)), ctx.Err()
			} else {
				return mo.Left[T, ValueI[Void, A]](e.MustLeft()), ctx.Err()
			}
		},
		func(ctx context.Context, focus B) (T, error) {
			return embed(focus), ctx.Err()
		},
		IxMatchVoid(),
		exprDef,
	)
}

// Constructor for a Prism optic. Prisms focus on zero or 1 element.
//
// The following Traversal constructors are available.
//   - [Prism]
//   - [PrismP] polymorphic
//   - [PrismE] error aware
//   - [PrismEP] polymorphic, error aware
//   - [PrismI] indexed
//   - [PrismIP] indexed, polymorphic
//   - [PrismIE] indexed, error aware
//   - [PrismIEP] indexed, polymorphic, error aware
//
// A prism is constructed from 3 functions
//
//   - match : performs a pattern match and returns the match or false.
//   - embed : converts the focus to the source type. This is equivalent to a [ReverseGetterFunc]
//   - expr: should return the expression type. See the [expr] package for more information.
//
// When traversed the prism will return no focuses if the match fails.
// When modified the prism will return the original source if the match fails.
func PrismE[S, A any](
	match func(ctx context.Context, source S) (A, bool, error),
	embed func(ctx context.Context, focus A) (S, error),
	exprDef ExpressionDef,
) Optic[Void, S, S, A, A, ReturnMany, ReadWrite, BiDir, Err] {
	return CombiPrism[ReadWrite, BiDir, Err](
		func(ctx context.Context, source S) (mo.Either[S, ValueI[Void, A]], error) {
			right, ok, err := match(ctx, source)
			if ok {
				return mo.Right[S, ValueI[Void, A]](ValI(Void{}, right)), err
			} else {
				return mo.Left[S, ValueI[Void, A]](source), err
			}
		},
		embed,
		IxMatchVoid(),
		exprDef,
	)
}

// Constructor for a Prism optic. Prisms focus on zero or 1 element.
//
// The following Traversal constructors are available.
//   - [Prism]
//   - [PrismP] polymorphic
//   - [PrismE] error aware
//   - [PrismEP] polymorphic, error aware
//   - [PrismI] indexed
//   - [PrismIP] indexed, polymorphic
//   - [PrismIE] indexed, error aware
//   - [PrismIEP] indexed, polymorphic, error aware
//
// A prism is constructed from 3 functions
//
//   - match : performs a pattern match and returns either the polymorphic result type or the matched value.
//   - embed : converts the focus to the result type. This is equivalent to a [ReverseGetterFunc]
//   - expr: should return the expression type. See the [expr] package for more information.
//
// When traversed the prism will return no focuses if the match fails.
// When modified the prism will return the original source if the match fails.
func PrismEP[S, T, A, B any](
	match func(ctx context.Context, source S) (mo.Either[T, A], error),
	embed func(ctx context.Context, focus B) (T, error),
	exprDef ExpressionDef,
) Optic[Void, S, T, A, B, ReturnMany, ReadWrite, BiDir, Err] {
	return CombiPrism[ReadWrite, BiDir, Err, Void, S, T, A, B](
		func(ctx context.Context, source S) (mo.Either[T, ValueI[Void, A]], error) {
			e, err := match(ctx, source)
			if err != nil {
				return mo.Either[T, ValueI[Void, A]]{}, err
			}

			if right, ok := e.Right(); ok {
				return mo.Right[T](ValI(Void{}, right)), ctx.Err()
			} else {
				return mo.Left[T, ValueI[Void, A]](e.MustLeft()), ctx.Err()
			}
		},
		embed,
		IxMatchVoid(),
		exprDef,
	)
}

// Constructor for a Prism optic. Prisms focus on zero or 1 element.
//
// The following Traversal constructors are available.
//   - [Prism]
//   - [PrismP] polymorphic
//   - [PrismE] error aware
//   - [PrismEP] polymorphic, error aware
//   - [PrismI] indexed
//   - [PrismIP] indexed, polymorphic
//   - [PrismIE] indexed, error aware
//   - [PrismIEP] indexed, polymorphic, error aware
//
// A prism is constructed from 3 functions
//
//   - match : performs a pattern match and returns the match or false.
//   - embed : converts the focus to the source type. This is equivalent to a [ReverseGetterFunc]
//   - expr: should return the expression type. See the [expr] package for more information.
//
// When traversed the prism will return no focuses if the match fails.
// When modified the prism will return the original source if the match fails.
func PrismI[I, S, A any](
	match func(source S) (I, A, bool),
	embed func(focus A) S,
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, S, A, A, ReturnMany, ReadWrite, BiDir, Pure] {

	ixMatchFnc := ensureSimpleIxMatch(ixMatch)

	return CombiPrism[ReadWrite, BiDir, Pure, I, S, S, A, A](
		func(ctx context.Context, source S) (mo.Either[S, ValueI[I, A]], error) {

			if ix, right, ok := match(source); ok {
				return mo.Right[S](ValI(ix, right)), ctx.Err()
			} else {
				return mo.Left[S, ValueI[I, A]](source), ctx.Err()
			}
		},
		func(ctx context.Context, focus A) (S, error) {
			return embed(focus), ctx.Err()
		},
		ixMatchFnc,
		exprDef,
	)
}

// Constructor for a Prism optic. Prisms focus on zero or 1 element.
//
// The following Traversal constructors are available.
//   - [Prism]
//   - [PrismP] polymorphic
//   - [PrismE] error aware
//   - [PrismEP] polymorphic, error aware
//   - [PrismI] indexed
//   - [PrismIP] indexed, polymorphic
//   - [PrismIE] indexed, error aware
//   - [PrismIEP] indexed, polymorphic, error aware
//
// A prism is constructed from 3 functions
//
//   - match : performs a pattern match and returns either the polymorphic result type or the matched value.
//   - embed : converts the focus to the result type. This is equivalent to a [ReverseGetterFunc]
//   - expr: should return the expression type. See the [expr] package for more information.
//
// When traversed the prism will return no focuses if the match fails.
// When modified the prism will return the original source if the match fails.
func PrismIP[I, S, T, A, B any](
	match func(source S) mo.Either[T, ValueI[I, A]],
	embed func(focus B) T,
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, T, A, B, ReturnMany, ReadWrite, BiDir, Pure] {
	return CombiPrism[ReadWrite, BiDir, Pure](
		func(ctx context.Context, source S) (mo.Either[T, ValueI[I, A]], error) {
			e := match(source)
			if right, ok := e.Right(); ok {
				return mo.Right[T](right), ctx.Err()
			} else {
				return mo.Left[T, ValueI[I, A]](e.MustLeft()), ctx.Err()
			}
		},
		func(ctx context.Context, focus B) (T, error) {
			return embed(focus), ctx.Err()
		},
		ixMatch,
		exprDef,
	)
}

// Constructor for a Prism optic. Prisms focus on zero or 1 element.
//
// The following Traversal constructors are available.
//   - [Prism]
//   - [PrismP] polymorphic
//   - [PrismE] error aware
//   - [PrismEP] polymorphic, error aware
//   - [PrismI] indexed
//   - [PrismIP] indexed, polymorphic
//   - [PrismIE] indexed, error aware
//   - [PrismIEP] indexed, polymorphic, error aware
//
// A prism is constructed from 3 functions
//
//   - match : performs a pattern match and returns the match or false.
//   - embed : converts the focus to the source type. This is equivalent to a [ReverseGetterFunc]
//   - expr: should return the expression type. See the [expr] package for more information.
//
// When traversed the prism will return no focuses if the match fails.
// When modified the prism will return the original source if the match fails.
func PrismIE[I, S, A any](
	match func(ctx context.Context, source S) (I, A, bool, error),
	embed func(ctx context.Context, focus A) (S, error),
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, S, A, A, ReturnMany, ReadWrite, BiDir, Err] {
	return CombiPrism[ReadWrite, BiDir, Err](
		func(ctx context.Context, source S) (mo.Either[S, ValueI[I, A]], error) {
			i, right, ok, err := match(ctx, source)
			if ok {
				return mo.Right[S, ValueI[I, A]](ValI(i, right)), err
			} else {
				return mo.Left[S, ValueI[I, A]](source), err
			}
		},
		embed,
		ixMatch,
		exprDef,
	)
}

// Constructor for a Prism optic. Prisms focus on zero or 1 element.
//
// The following Traversal constructors are available.
//   - [Prism]
//   - [PrismP] polymorphic
//   - [PrismE] error aware
//   - [PrismEP] polymorphic, error aware
//   - [PrismI] indexed
//   - [PrismIP] indexed, polymorphic
//   - [PrismIE] indexed, error aware
//   - [PrismIEP] indexed, polymorphic, error aware
//
// A prism is constructed from 3 functions
//
//   - match : performs a pattern match and returns either the polymorphic result type or the matched value.
//   - embed : converts the focus to the result type. This is equivalent to a [ReverseGetterFunc]
//   - expr: should return the expression type. See the [expr] package for more information.
//
// When traversed the prism will return no focuses if the match fails.
// When modified the prism will return the original source if the match fails.
func PrismIEP[I, S, T, A, B any](
	match func(ctx context.Context, source S) (mo.Either[T, ValueI[I, A]], error),
	embed func(ctx context.Context, focus B) (T, error),
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, T, A, B, ReturnMany, ReadWrite, BiDir, Err] {
	return CombiPrism[ReadWrite, BiDir, Err](
		match,
		embed,
		ixMatch,
		exprDef,
	)
}
