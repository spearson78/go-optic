package optic

import (
	"context"
	"errors"
	"reflect"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spearson78/go-optic/expr"
)

// [Present] returns a [Predicate] that focuses true if the option is present.
func Present[A any]() Optic[Void, mo.Option[A], mo.Option[A], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return NotEmpty(TraverseOption[A]())
}

// Some returns a [Prism] that matches on the Some value of a [mo.Option]. This has the effect of traversing the option ignoring None values.
//
// This function is synonymous with [TraverseOption]
//
// See:
//   - [SomeP] for a polymorphic version
//   - [SomeE] for an impure version
//   - [SomeEP] for a polymorphic impure version
//   - [None] for a version matches on the None value.
//   - [Non] for an optic that provides default values for None values.
func Some[A any]() Optic[Void, mo.Option[A], mo.Option[A], A, A, ReturnMany, ReadWrite, BiDir, Pure] {
	return SomeP[A, A]()
}

// SomeE returns an [Iso] that focuses the value of a [mo.Option] or returns [ErrOptionNone] for none values.
//
// See:
//   - [Some] for a pure prism version.
//   - [SomeP] for a polymorphic version
//   - [SomeEP] for a polymorphic impure version
//   - [None] for a version matches on the None value.
//   - [Non] for an optic that provides default values for None values.
func SomeE[A any]() Optic[Void, mo.Option[A], mo.Option[A], A, A, ReturnOne, ReadWrite, BiDir, Err] {
	return SomeEP[A, A]()
}

// Some returns a polymorphic [Prism] that matches on the Some value of a [mo.Option]. This has the effect of traversing the option ignoring None values.
//
// See:
//   - [Some] for none a polymorphic version
//   - [SomeE] for a non polymorphic impure version
//   - [SomeEP] for a polymorphic impure version
//   - [None] for a version matches on the None value.
//   - [Non] for an optic that provides default values for None values.
func SomeP[A, B any]() Optic[Void, mo.Option[A], mo.Option[B], A, B, ReturnMany, ReadWrite, BiDir, Pure] {
	return PrismP[mo.Option[A], mo.Option[B], A, B](
		func(source mo.Option[A]) mo.Either[mo.Option[B], A] {
			a, ok := source.Get()
			if ok {
				return mo.Right[mo.Option[B], A](a)
			} else {
				return mo.Left[mo.Option[B], A](mo.None[B]())
			}
		},
		func(focus B) mo.Option[B] {
			return mo.Some(focus)
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Some{
				OpticTypeExpr: ot,
			}
		}),
	)
}

var ErrOptionNone = errors.New("none passed to SomeE / SomeEP")

// SomeEP returns a polymorphic [Iso] that focuses the value of a [mo.Option] or returns [ErrOptionNone] for none values.
//
// See:
//   - [Some] for a pure prism version.
//   - [SomeP] for a pure polymorphic version
//   - [SomeE] for a non polymorphic impure version
//   - [None] for a version matches on the None value.
//   - [Non] for an optic that provides default values for None values.
func SomeEP[A, B any]() Optic[Void, mo.Option[A], mo.Option[B], A, B, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoEP[mo.Option[A], mo.Option[B], A, B](
		func(ctx context.Context, source mo.Option[A]) (A, error) {
			a, ok := source.Get()
			if ok {
				return a, nil
			} else {
				return a, ErrOptionNone
			}
		},
		func(ctx context.Context, focus B) (mo.Option[B], error) {
			return mo.Some(focus), nil
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Some{
				OpticTypeExpr: ot,
			}
		}),
	)
}

// TraverseOption returns a [Prism] that matches on the Some value of a [mo.Option]. This has the effect of traversing the option ignoring None values.
//
// This function is synonymous with [Some]
//
// See:
//   - [None] for a version matches on the None value.
//   - [Non] for an optic that provides default values for None values.
func TraverseOption[A any]() Optic[Void, mo.Option[A], mo.Option[A], A, A, ReturnMany, ReadWrite, BiDir, Pure] {
	return Some[A]()
}

// None returns a [Prism] that matches on the None state of a [mo.Option]. A void focus type is returned as the none value has no valid representation.
//
// See:
//   - [Non] for an optic that provides default values for None values.
//   - [Some] for an optic that matches on Some values.
func None[A any]() Optic[Void, mo.Option[A], mo.Option[A], Void, Void, ReturnMany, ReadWrite, BiDir, Pure] {
	return Prism[mo.Option[A], Void](
		func(source mo.Option[A]) (Void, bool) {
			if source.IsAbsent() {
				return Void{}, true
			} else {
				return Void{}, false
			}
		},
		func(focus Void) mo.Option[A] {
			return mo.None[A]()
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.None{
				OpticTypeExpr: ot,
			}
		}),
	)
}

// Non returns an [Iso] that converts a None value of a [mo.Option] to the given default value.
//
// Warning: this has the effect that under modification the default value is converted to a None [mo.Option] value.
//
//   - See [Some] for an optic that filters out None values.
//   - See [None] for an optic that filters out Some values.
func Non[A, ERR any](defaultValue A, equal Predicate[lo.Tuple2[A, A], ERR]) Optic[Void, mo.Option[A], mo.Option[A], A, A, ReturnOne, ReadWrite, BiDir, ERR] {
	return CombiIso[ReadWrite, BiDir, ERR, mo.Option[A], mo.Option[A], A, A](
		func(ctx context.Context, source mo.Option[A]) (A, error) {
			return source.OrElse(defaultValue), nil
		},
		func(ctx context.Context, focus A) (mo.Option[A], error) {
			eq, err := PredGet(ctx, equal, lo.T2(focus, defaultValue))
			if err != nil {
				return mo.None[A](), err
			}

			if eq {
				return mo.None[A](), nil
			} else {
				return mo.Some(focus), nil
			}
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.Non{
					OpticTypeExpr: ot,
					Default:       defaultValue,
					Equal:         equal.AsExpr(),
				}
			},
			equal,
		),
	)
}

// PtrOption returns a [Iso] that converts a pointer into a [mo.Option]. Nil is converted to the None option.
//
// See:
//   - [PtrOptionP] for a polymorphic version.
func PtrOption[A any]() Optic[Void, *A, *A, mo.Option[A], mo.Option[A], ReturnOne, ReadWrite, BiDir, Pure] {
	return PtrOptionP[A, A]()
}

// PtrOptionP returns a polymorphic [Iso] that converts a pointer into a [mo.Option]. Nil is converted to the None option.
//
// See:
//   - [PtrOption] for a non polymorphic version.
func PtrOptionP[A, B any]() Optic[Void, *A, *B, mo.Option[A], mo.Option[B], ReturnOne, ReadWrite, BiDir, Pure] {

	return IsoP[*A, *B, mo.Option[A], mo.Option[B]](
		func(source *A) mo.Option[A] {
			if source == nil {
				return mo.None[A]()
			} else {
				return mo.Some(*source)
			}
		},
		func(focus mo.Option[B]) *B {
			val, ok := focus.Get()
			if !ok {
				return nil
			}
			return &val
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.PtrOption{
				OpticTypeExpr: ot,
				A:             reflect.TypeFor[A](),
				B:             reflect.TypeFor[B](),
			}
		}),
	)
}

// OptionOfFirst focuses the first element wrapped in a [mo.Option] or [mo.None] if no elements are focused.
func OptionOfFirst[I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, T, mo.Option[A], mo.Option[B], ReturnOne, ReadOnly, UniDir, ERR] {
	return CombiGetter[ERR, I, S, T, mo.Option[A], mo.Option[B]](
		func(ctx context.Context, source S) (I, mo.Option[A], error) {
			i, focus, err := o.AsGetter()(ctx, source)
			if errors.Is(err, ErrEmptyGet) {
				return i, mo.None[A](), nil
			}
			if err != nil {
				return i, mo.None[A](), err
			}
			return i, mo.Some(focus), err
		},
		o.AsIxMatch(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.OptionOfFirst{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)

}
