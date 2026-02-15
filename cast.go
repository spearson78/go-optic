package optic

import (
	"context"
	"fmt"
	"reflect"

	"github.com/samber/mo"
	"github.com/spearson78/go-optic/expr"
)

type ErrCast struct {
	From reflect.Type //The expected typed to cast from.
	To   reflect.Type //The type that Of was being cast to.
}

func (e *ErrCast) Error() string {
	return fmt.Sprintf("cannot cast %v to %v", e.From, e.To)
}

// ErrCastOf means that a cast failed.
type ErrCastOf struct {
	Of   reflect.Type //The type of the value that was being cast.
	From reflect.Type //The expected typed to cast from.
	To   reflect.Type //The type that Of was being cast to.
}

func (e *ErrCastOf) Error() string {
	return fmt.Sprintf("cast failed %v -> %v for %v", e.From, e.To, e.Of)
}

func unsafeCast[A, S any](source S) (A, bool) {

	focus, ok := any(source).(A)
	if !ok {
		sourceVal := reflect.ValueOf(source)
		kind := sourceVal.Kind()
		switch kind {
		case 0:
			focus = reflect.Zero(reflect.TypeFor[A]()).Interface().(A)
			ok = true
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
			if sourceVal.IsNil() {
				focus = reflect.Zero(reflect.TypeFor[A]()).Interface().(A)
				ok = true
			} else {
				focusType := reflect.TypeFor[A]()
				if sourceVal.Type().ConvertibleTo(focusType) {
					focusVal := sourceVal.Convert(focusType)
					focus, ok = focusVal.Interface().(A)
				}
			}
		default:
			focusType := reflect.TypeFor[A]()
			if sourceVal.Type().ConvertibleTo(focusType) {
				focusVal := sourceVal.Convert(focusType)
				focus, ok = focusVal.Interface().(A)
			}
		}
	}

	return focus, ok
}

// DownCast returns a [Prism] that performs a safe down cast from S to A and an upcast from A to S
// The down cast may fail in which case no result is focused.
// Warning: If the up cast fails a panic will occur.
//
// See:
//   - [DownCastP] for a polymorphic version
//   - [UpCast] for a read only pure version.
//   - [IsoCast] for a pure isomorphic version.
//   - [IsoCastE] for isomorphic version that returns an error on cast failure.
func DownCast[S, A any]() Optic[Void, S, S, A, A, ReturnMany, ReadWrite, BiDir, Pure] {
	return DownCastP[S, A, A]()
}

// DownCastP returns a polymorphic [Prism] that performs a safe down cast from S to A and an upcast from B to S
// The down cast may fail in which case no result is focused.
// Warning: If the up cast fails a panic will occur.
//
// See:
//   - [DownCast] for a non polymorphic version
//   - [UpCast] for a read only pure version.
//   - [IsoCast] for a pure isomorphic version.
//   - [IsoCastE] for isomorphic version that returns an error on cast failure.
func DownCastP[S, A, B any]() Optic[Void, S, S, A, B, ReturnMany, ReadWrite, BiDir, Pure] {

	ta := reflect.TypeFor[A]()
	ts := reflect.TypeFor[S]()

	if !ta.ConvertibleTo(ts) {
		panic(fmt.Errorf("DownCastP invalid down cast : %w", &ErrCast{
			From: ta,
			To:   ts,
		}))
	}

	tb := reflect.TypeFor[B]()

	if !tb.ConvertibleTo(ts) {
		panic(fmt.Errorf("DownCastP invalid up cast : %w", &ErrCast{
			From: tb,
			To:   ts,
		}))
	}

	return CombiPrism[ReadWrite, BiDir, Pure, Void, S, S, A, B](
		func(ctx context.Context, source S) (mo.Either[S, ValueI[Void, A]], error) {
			focus, ok := unsafeCast[A](source)

			if ok {
				return mo.Right[S, ValueI[Void, A]](ValI(Void{}, focus)), nil
			} else {
				return mo.Left[S, ValueI[Void, A]](source), nil
			}
		},
		func(ctx context.Context, focus B) (S, error) {
			s, ok := unsafeCast[S](focus)
			if !ok {
				return s, &ErrCastOf{
					Of:   reflect.TypeOf(focus),
					From: reflect.TypeFor[B](),
					To:   reflect.TypeFor[S](),
				}
			}
			return s, nil
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Cast{
				OpticTypeExpr: ot,
			}
		}),
	)
}

// UpCast returns an [Operation] that performs a safe type up cast from S to A.
// The cast is declared safe by the user. If the cast fails a panic will occur.
//
// See:
//   - [DownCast] for a pure version that ignores cast errors.
//   - [IsoCast] for an Iso version.
func UpCast[S, A any]() Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, Pure] {

	ta := reflect.TypeFor[A]()
	ts := reflect.TypeFor[S]()

	if !ts.ConvertibleTo(ta) {
		panic(fmt.Errorf("UpCast invalid up cast : %w", &ErrCast{
			From: ta,
			To:   ts,
		}))
	}

	return CombiGetter[Pure, Void, S, S, A, A](
		func(ctx context.Context, source S) (Void, A, error) {
			focus, ok := unsafeCast[A](source)
			if !ok {
				return Void{}, focus, &ErrCastOf{
					Of:   reflect.TypeOf(source),
					From: reflect.TypeFor[S](),
					To:   reflect.TypeFor[A](),
				}
			}

			return Void{}, focus, ctx.Err()
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Cast{
				OpticTypeExpr: ot,
			}
		}),
	)
}

// IsoCast returns a [Iso] that performs a safe type cast from S to A and back from A to S
// The cast is declared safe by the user. If the cast fails a panic will occur.
//
// See:
//   - [IsoCastP] for a polymorphic version
//   - [DownCast] for a version that ignores cast errors
//   - [UpCast] for an impure version
//
// Warning: if the types are not compatible ReverseGet will return an error.
func IsoCast[S, A any]() Optic[Void, S, S, A, A, ReturnOne, ReadWrite, BiDir, Pure] {
	return IsoCastP[S, S, A, A]()
}

// IsoCastP returns a polymorphic [Iso] that performs a safe type cast from S to A and in reverse from B to T
// The cast is declared safe by the user. If the cast fails a panic will occur.
//
// See:
//   - [IsoCast] for a non polymorphic version
//
// Warning: if the types are not compatible ReverseGet will return an error.
func IsoCastP[S, T, A, B any]() Optic[Void, S, T, A, B, ReturnOne, ReadWrite, BiDir, Pure] {

	return CombiIso[ReadWrite, BiDir, Pure, S, T, A, B](
		func(ctx context.Context, source S) (A, error) {
			focus, ok := unsafeCast[A](source)
			if !ok {
				panic(&ErrCastOf{
					Of:   reflect.TypeOf(source),
					From: reflect.TypeFor[S](),
					To:   reflect.TypeFor[A](),
				})
			}

			return focus, ctx.Err()
		},
		func(ctx context.Context, focus B) (T, error) {
			source, ok := unsafeCast[T](focus)
			if !ok {
				panic(&ErrCastOf{
					Of:   reflect.TypeOf(focus),
					From: reflect.TypeFor[B](),
					To:   reflect.TypeFor[T](),
				})
			}
			return source, nil
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Cast{
				OpticTypeExpr: ot,
			}
		}),
	)
}

// IsoCastE returns a [Iso] that performs an cast from S to A and back from A to S
// If the cast fails an error will occur.
//
// See:
//   - [IsoCastEP] for a polymorphic version
//   - [DownCast] for a version that ignores cast errors
//   - [IsoCast] for an pure version
func IsoCastE[S, A any]() Optic[Void, S, S, A, A, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoCastEP[S, S, A, A]()
}

// IsoCastEP returns a polymorphic [Iso] that performs a cast from S to A and back from B to T
// If the cast fails an error will occur.
//
// See:
//   - [IsoCastE] for a non polymorphic version
//   - [IsoCastP] for a pure version.
func IsoCastEP[S, T, A, B any]() Optic[Void, S, T, A, B, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoEP[S, T, A, B](
		func(ctx context.Context, source S) (A, error) {
			focus, ok := unsafeCast[A](source)
			if !ok {
				return focus, &ErrCastOf{
					Of:   reflect.TypeOf(source),
					From: reflect.TypeFor[S](),
					To:   reflect.TypeFor[A](),
				}
			}
			return focus, nil
		},
		func(ctx context.Context, focus B) (T, error) {
			source, ok := unsafeCast[T](focus)
			if !ok {
				return source, &ErrCastOf{
					Of:   reflect.TypeOf(focus),
					From: reflect.TypeFor[B](),
					To:   reflect.TypeFor[T](),
				}
			}
			return source, nil
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Cast{
				OpticTypeExpr: ot,
			}
		}),
	)
}
