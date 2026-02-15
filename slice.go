package optic

import (
	"context"
	"errors"
	"iter"
	"reflect"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spearson78/go-optic/expr"
)

// TraverseSlice returns a [Traversal] that focuses on the elements of a slice.
//
// See: [TraverseSliceP] for a polymorphic version
func TraverseSlice[A any]() Optic[int, []A, []A, A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseSliceP[A, A]()
}

// TraverseSliceP returns a polymorphic [Traversal] that focuses on the elements of a slice.
//
// See: [TraverseSlice] for a non polymorphic version
func TraverseSliceP[A, B any]() Optic[int, []A, []B, A, B, ReturnMany, ReadWrite, UniDir, Pure] {
	return CombiTraversal[ReturnMany, ReadWrite, Pure, int, []A, []B, A, B](
		func(ctx context.Context, source []A) SeqIE[int, A] {
			return func(yield func(ValueIE[int, A]) bool) {
				for i, a := range source {
					if !yield(ValIE(i, a, ctx.Err())) {
						return
					}
				}
			}
		},
		func(ctx context.Context, source []A) (int, error) {
			return len(source), nil
		},
		func(ctx context.Context, fmap func(index int, focus A) (B, error), source []A) ([]B, error) {
			ret := make([]B, len(source))
			for i, a := range source {
				b, err := fmap(i, a)
				if err != nil {
					return nil, err
				}
				ret[i] = b
			}

			return ret, nil
		},
		func(ctx context.Context, index int, source []A) SeqIE[int, A] {
			return func(yield func(ValueIE[int, A]) bool) {
				if index >= 0 && index < len(source) {
					yield(ValIE(index, source[index], nil))
				}
			}
		},
		IxMatchComparable[int](),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Traverse{
				OpticTypeExpr: ot,
			}
		}),
	)
}

// AtSlice returns a [Lens] that focuses the index slot in a slice.
// When viewing the none option means either the slice was too small or the value was the uninitialized value.
// When modifying the slice will expand if necessary and be filled with uninitialized value as needed.
//
// See:
//   - [AtSliceT2] for a version that can use a dynamic index.
//   - [Element] for an optic that focuses by index position.
//   - [IndexSlice] for a version without the optional focus
func AtSlice[V, ERR any](index int, equals Predicate[lo.Tuple2[V, V], ERR]) Optic[int, []V, []V, mo.Option[V], mo.Option[V], ReturnOne, ReadWrite, UniDir, ERR] {
	return Ret1(Rw(Ud((EErrR(Compose(
		T2Of(
			Identity[[]V](),
			IgnoreWrite(Const[[]V](index)),
		),
		AtSliceT2[V](equals),
	))))))
}

// AtSliceT2 returns a [Lens] that focuses the tuple.B index slot in a slice tuple.A.
// When viewing the none option means either the slice was too small or the value was the uninitialized value.
// When modifying the slice will expand if necessary and be filled with uninitialized value as needed.
//
// See:
//   - [AtSlice] for a version that uses a fixed index..
//   - [Element] for an optic that focuses by index position.
//   - [IndexSlice] for a version without the optional focus
func AtSliceT2[V, ERR any](equals Predicate[lo.Tuple2[V, V], ERR]) Optic[int, lo.Tuple2[[]V, int], lo.Tuple2[[]V, int], mo.Option[V], mo.Option[V], ReturnOne, ReadWrite, UniDir, ERR] {
	return CombiGetMod[ReadWrite, ERR, int, lo.Tuple2[[]V, int], lo.Tuple2[[]V, int], mo.Option[V], mo.Option[V]](
		func(ctx context.Context, sourceT2 lo.Tuple2[[]V, int]) (int, mo.Option[V], error) {
			source := sourceT2.A
			index := sourceT2.B

			if len(source) == 0 {
				return index, mo.None[V](), nil
			}

			if index < 0 {
				index = len(source) + index
			}

			if index >= len(source) || index < 0 {
				return index, mo.None[V](), nil
			}

			var defaultVal V
			match, err := PredGet(ctx, equals, lo.T2(defaultVal, source[index]))
			if err != nil {
				return index, mo.None[V](), err
			}

			if match {
				return index, mo.None[V](), nil
			}

			return index, mo.Some(source[index]), nil

		},
		func(ctx context.Context, fmap func(index int, focus mo.Option[V]) (mo.Option[V], error), sourceT2 lo.Tuple2[[]V, int]) (lo.Tuple2[[]V, int], error) {
			source := sourceT2.A
			index := sourceT2.B

			if index < 0 {
				index = len(source) + index
			}

			if index >= 0 {

				if index >= len(source)-1 {
					newVal, err := fmap(index, mo.None[V]())
					err = JoinCtxErr(ctx, err)
					if err != nil {
						return lo.T2([]V(nil), index), err
					}

					if v, ok := newVal.Get(); ok {
						ret := make([]V, index+1)
						copy(ret, source)
						ret[index] = v

						return lo.T2(ret, index), nil
					} else {
						return sourceT2, nil
					}
				} else {

					var defaultVal V
					match, err := PredGet(ctx, equals, lo.T2(defaultVal, source[index]))
					if err != nil {
						return lo.T2([]V(nil), index), err
					}

					var newVal mo.Option[V]
					if match {
						newVal, err = fmap(index, mo.None[V]())
					} else {
						newVal, err = fmap(index, mo.Some(source[index]))
					}

					err = JoinCtxErr(ctx, err)
					if err != nil {
						return lo.T2([]V(nil), index), err
					}

					if v, ok := newVal.Get(); ok {
						ret := append([]V(nil), source...)
						ret[index] = v
						return lo.T2(ret, index), nil
					} else {

						var v V
						ret := append([]V(nil), source...)
						ret[index] = v

						return lo.T2(ret, index), nil
					}
				}
			} else {
				newVal, err := fmap(index, mo.None[V]())
				err = JoinCtxErr(ctx, err)
				if err != nil {
					return lo.T2([]V(nil), index), err
				}

				ret := make([]V, -index+len(source))
				copy(ret[-index:], source)

				if v, ok := newVal.Get(); ok {
					ret[0] = v
				}

				return lo.T2(ret, index), nil
			}
		},
		IxMatchComparable[int](),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.AtT2{
					OpticTypeExpr: ot,
					V:             reflect.TypeFor[V](),
					Equals:        equals.AsExpr(),
				}
			},
			equals,
		),
	)
}

// The SliceOf combinator focuses on a slice of all the elements in the given optic.
//
// Under modification this slice can be modified using standard operations and will be rebuilt into the original data structure.
// If the modified slice contains fewer elements the result will use values from the original source.
// If the modified slice contains more elements they will be ignored.
//
// See:
//   - [SliceOfP] for a polymorphic version.
func SliceOf[I, S, T, A, RETI, RW, DIR, ERR any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR], size int) Optic[Void, S, T, []A, []A, ReturnOne, RW, UniDir, ERR] {
	return CombiLens[RW, ERR, Void, S, T, []A, []A](
		func(ctx context.Context, source S) (Void, []A, error) {
			ret := make([]A, 0, size)
			var err error
			o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
				_, a, focusErr := val.Get()
				err = JoinCtxErr(ctx, focusErr)
				if err != nil {
					return false
				}
				ret = append(ret, a)
				return true
			})

			if err != nil {
				return Void{}, nil, err
			}

			return Void{}, ret, err
		},
		func(ctx context.Context, va []A, vs S) (T, error) {
			i := 0
			l := len(va)
			ret, err := o.AsModify()(ctx, func(index I, focus A) (A, error) {
				if i >= l {
					return focus, ctx.Err()
				} else {
					ret := va[i]
					i++
					return ret, ctx.Err()
				}
			}, vs)

			return ret, JoinCtxErr(ctx, err)
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.SliceOf{
					OpticTypeExpr: ot,
					I:             reflect.TypeFor[I](),
					A:             reflect.TypeFor[A](),
					B:             reflect.TypeFor[A](),
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

var ErrUnsafeMissingElement = errors.New("not enough elements were supplied")

// The SliceOfP combinator focuses on a polymorphic slice of all the elements in the given optic.
//
// Under modification this slice can be modified using standard operations and will be rebuilt into the original data structure.
// If the modified slice contains fewer elements then [ErrUnsafeMissingElement] will be returned
// If the modified slice contains more elements they will be ignored.
//
// See:
//   - [SliceOf] for a safe non polymorphic version.
func SliceOfP[I, S, T, A, B, RET any, RW any, DIR any, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], size int) Optic[Void, S, T, []A, []B, ReturnOne, RW, UniDir, Err] {
	return CombiLens[RW, Err, Void, S, T, []A, []B](
		func(ctx context.Context, source S) (Void, []A, error) {
			ret := make([]A, 0, size)
			var err error
			o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
				_, a, focusErr := val.Get()
				err = JoinCtxErr(ctx, focusErr)
				if err != nil {
					return false
				}
				ret = append(ret, a)
				return true
			})

			if err != nil {
				return Void{}, nil, err
			}

			return Void{}, ret, err
		},
		func(ctx context.Context, b []B, source S) (T, error) {
			i := 0

			l := len(b)
			ret, err := o.AsModify()(ctx, func(index I, focus A) (B, error) {
				if i >= l {
					var b B
					return b, JoinCtxErr(ctx, ErrUnsafeMissingElement)
				} else {
					rb := b[i]
					i++
					return rb, ctx.Err()
				}
			}, source)

			err = JoinCtxErr(ctx, err)
			return ret, err
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.SliceOf{
					OpticTypeExpr: ot,
					I:             reflect.TypeFor[I](),
					A:             reflect.TypeFor[A](),
					B:             reflect.TypeFor[B](),
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// SliceToCol returns an [Iso] that converts a slice to an [Collection]
func SliceToCol[T any]() Optic[Void, []T, []T, Collection[int, T, Pure], Collection[int, T, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return SliceToColP[T, T]()
}

// SliceToColP returns a polymorphic [Iso] that converts a slice to an [Collection]
func SliceToColP[A, B any]() Optic[Void, []A, []B, Collection[int, A, Pure], Collection[int, B, Pure], ReturnOne, ReadWrite, BiDir, Pure] {
	return CombiIso[ReadWrite, BiDir, Pure, []A, []B, Collection[int, A, Pure], Collection[int, B, Pure]](
		func(ctx context.Context, source []A) (Collection[int, A, Pure], error) {
			return ColI(
				func(yield func(index int, focus A) bool) {
					for i, v := range source {
						if !yield(i, v) {
							return
						}
					}
				},
				func(index int) iter.Seq2[int, A] {
					return func(yield func(index int, focus A) bool) {
						if index >= 0 && index < len(source) {
							yield(index, source[index])
						}
					}
				},
				IxMatchComparable[int](),
				func() int {
					return len(source)
				},
			), nil
		},
		func(ctx context.Context, focus Collection[int, B, Pure]) ([]B, error) {
			var ret []B
			var retErr error

			i := 0
			focus.AsIter()(ctx)(func(val ValueIE[int, B]) bool {
				_, focus, focusErr := val.Get()
				if focusErr != nil {
					retErr = focusErr
					return false
				}
				i++
				ret = append(ret, focus)
				return true
			})
			return ret, retErr
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.ToCol{
				OpticTypeExpr: ot,
				I:             reflect.TypeFor[int](),
				A:             reflect.TypeFor[A](),
				B:             reflect.TypeFor[A](),
			}
		}),
	)
}

//go:generate ./makecolops -nocolof slice_ops.go "optic" "Slice" "A any" "A, B any" "" "SliceColTypeP[A, A]()" "SliceColTypeP[A, B]()" "int" "[]A" "[]B" "A" "B"

// SliceColType returns a [CollectionType] for slices
func SliceColType[S any]() CollectionType[int, []S, []S, S, S, Pure] {
	return SliceColTypeP[S, S]()
}

// SliceColTypeP returns a polymorphic [CollectionType] for slices.
func SliceColTypeP[S, T any]() CollectionType[int, []S, []T, S, T, Pure] {
	return ColTypeP(
		SliceToColP[S, T](),
		AsReverseGet(SliceToColP[T, S]()),
		TraverseSliceP[S, T](),
	)
}

// AppendSliceReducer retuens a reducer that appends the elements to a slice.
func AppendSliceReducer[T any](cap int) ReductionP[[]T, T, []T, Pure] {
	return ReducerP[[]T, T, []T](
		func() []T {
			if cap == 0 {
				return nil
			}
			return make([]T, 0, cap)
		},
		func(state []T, appendVal T) []T {
			return append(state, appendVal)
		},
		func(state []T) []T {
			return state
		},
		ReducerExprDef(
			func(t expr.ReducerTypeExpr) expr.ReducerExpression {
				return expr.AppendSliceReducerExpr{
					ReducerTypeExpr: t,
				}
			},
		),
	)
}

// SliceTail focuses a slice containing all but the first element of the source or nothing if the slice has 0 or 1 elements.
func SliceTail[A any]() Optic[Void, []A, []A, []A, []A, ReturnMany, ReadWrite, UniDir, Pure] {
	return CombiTraversal[ReturnMany, ReadWrite, Pure, Void, []A, []A, []A, []A](
		func(ctx context.Context, source []A) SeqIE[Void, []A] {
			return func(yield func(val ValueIE[Void, []A]) bool) {
				if len(source) == 0 {
					return
				}

				if len(source) == 1 {
					return
				}

				yield(ValIE(Void{}, source[1:], nil))
			}
		},
		func(ctx context.Context, source []A) (int, error) {
			if len(source) == 0 {
				return 0, nil
			}

			if len(source) == 1 {
				return 0, nil
			}

			return 1, nil
		},
		func(ctx context.Context, fmap func(index Void, focus []A) ([]A, error), source []A) ([]A, error) {
			modifiedTail, err := fmap(Void{}, source[1:])
			if err != nil {
				return nil, err
			}

			ret := make([]A, len(source))
			ret[0] = source[0]
			n := copy(ret[1:], modifiedTail)
			return ret[:n+1], nil

		},
		nil,
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.SliceChildren{
					OpticTypeExpr: ot,
					A:             reflect.TypeFor[A](),
				}
			},
		),
	)
}

func MakeSlice[S any, V any](size, cap int) Optic[Void, S, S, []V, []V, ReturnOne, ReadOnly, UniDir, Pure] {
	return Getter[S, []V](
		func(source S) []V {
			return make([]V, size, cap)
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Make{
				OpticTypeExpr: ot,
				Size:          []int{size, cap},
			}
		}),
	)
}
