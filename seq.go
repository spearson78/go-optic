package optic

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"reflect"
	"strings"

	"github.com/samber/lo"
	"github.com/spearson78/go-optic/expr"
)

// Push style iterator for index and value
type SeqI[I, A any] iter.Seq2[I, A]

// Push style iterator for error and value
//
// Note: error is first to avoid it being ignored.
type SeqE[A any] iter.Seq[ValueE[A]]

// Push style iterator for index,value and error.
type SeqIE[I any, A any] iter.Seq[ValueIE[I, A]]

func (s SeqI[I, A]) String() string {
	var sb strings.Builder

	first := true
	sb.WriteString("Seq[")

	if s != nil {

		for index, focus := range s {
			if !first {
				sb.WriteString(" ")
			}
			first = false

			sb.WriteString(fmt.Sprintf("%v:%v", index, focus))
		}
	} else {
		sb.WriteString("NIL")
	}
	sb.WriteString("]")

	return sb.String()
}

func (s SeqE[A]) String() string {
	var sb strings.Builder

	first := true
	sb.WriteString("Seq[")

	if s != nil {

		for val := range s {
			focus, err := val.Get()
			if !first {
				sb.WriteString(" ")
			}
			first = false

			if err != nil {
				sb.WriteString(fmt.Sprintf("<%v>", err))
			} else {
				sb.WriteString(fmt.Sprintf("%v", focus))
			}

		}
	} else {
		sb.WriteString("NIL")
	}
	sb.WriteString("]")

	return sb.String()
}

func (s SeqIE[I, A]) String() string {
	var sb strings.Builder

	first := true
	sb.WriteString("Seq[")

	if s != nil {

		s(func(val ValueIE[I, A]) bool {
			index, focus, err := val.Get()
			if !first {
				sb.WriteString(" ")
			}
			first = false

			if err != nil {
				sb.WriteString(fmt.Sprintf("<%v>", err))
			} else {
				sb.WriteString(fmt.Sprintf("%v:%v", index, focus))
			}
			return true
		})
	} else {
		sb.WriteString("NIL")
	}
	sb.WriteString("]")

	return sb.String()
}

func (s SeqIE[I, A]) MarshalJSON() ([]byte, error) {
	var slice []ValueI[I, A]
	var retErr error

	s(func(val ValueIE[I, A]) bool {
		index, focus, err := val.Get()
		if err != nil {
			retErr = err
			return false
		}
		slice = append(slice, ValI(index, focus))
		return true
	})

	if retErr != nil {
		return nil, retErr
	}

	return json.Marshal(slice)
}

func (s *SeqIE[I, A]) UnmarshalJSON(data []byte) error {

	var slice []ValueI[I, A]
	err := json.Unmarshal(data, &slice)
	if err != nil {
		return err
	}

	*s = func(yield func(val ValueIE[I, A]) bool) {
		for _, v := range slice {
			if !yield(ValIE(v.index, v.value, nil)) {
				break
			}
		}
	}

	return nil
}

// PullI converts the “push-style” [SeqI] into a “pull-style” iterator accessed by the two functions next and stop.
//
// See [iter.Pull2] for more documentation.
func PullI[I, A any](iterFnc SeqI[I, A]) (next func() (index I, focus A, ok bool), stop func()) {
	return iter.Pull2(iter.Seq2[I, A](iterFnc))
}

// PullE converts the “push-style” [SeqI] into a “pull-style” iterator accessed by the two functions next and stop.
//
// See [iter.Pull2] for more documentation.
func PullE[I, A any](iterFnc SeqE[A]) (next func() (focus A, ok bool, err error), stop func()) {
	n, s := iter.Pull(iter.Seq[ValueE[A]](iterFnc))
	return func() (focus A, ok bool, err error) {
		val, ok := n()
		return val.Value(), ok, val.Error()
	}, s

}

// PullIE converts the “push-style” [SeqIE] into a “pull-style” iterator accessed by the two functions next and stop.
//
// See [iter.Pull] for more documentation.
func PullIE[I, A any](iterFnc SeqIE[I, A]) (next func() (index I, focus A, err error, ok bool), stop func()) {
	iterAdapter := func(yield func(lo.Tuple3[I, A, error]) bool) {
		iterFnc(func(val ValueIE[I, A]) bool {
			index, focus, err := val.Get()
			return yield(lo.T3(index, focus, err))
		})
	}

	inext, istop := iter.Pull(iterAdapter)

	return func() (I, A, error, bool) {
		v, ok := inext()
		return v.A, v.B, v.C, ok

	}, istop
}

// The SeqOf combinator focuses on a [iter.Seq] of all the elements in the given optic.
//
// Under modification this seq can be modified using standard operations and will be rebuilt into the original data structure.
// If the modified seq contains fewer elements the result will use values from the original source.
// If the modified seq contains more elements they will be ignored.
//
// See:
//   - [SeqEOfP] for a polymorphic version.
//   - [ISeqOf] for an indexed version.
//   - [SeqIEOf] for a [SeqIE] version.
func SeqOf[I, S, T, A, RETI, RW, DIR any, ERR TPure](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR]) Optic[Void, S, T, iter.Seq[A], iter.Seq[A], ReturnOne, RW, UniDir, ERR] {
	return CombiLens[RW, ERR, Void, S, T, iter.Seq[A], iter.Seq[A]](
		func(ctx context.Context, source S) (Void, iter.Seq[A], error) {
			ret := func(yield func(A) bool) {
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					_, a, focusErr := val.Get()
					if focusErr != nil {
						panic(focusErr) //ERR is TPure
					}
					return yield(a)
				})
			}
			return Void{}, ret, nil
		},
		func(ctx context.Context, setVal iter.Seq[A], source S) (T, error) {

			next, stop := iter.Pull(setVal)
			defer stop()
			ret, err := o.AsModify()(ctx, func(index I, focus A) (A, error) {

				val, ok := next()
				if !ok {
					return focus, ctx.Err()
				} else {
					return val, ctx.Err()
				}
			}, source)

			return ret, JoinCtxErr(ctx, err)
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.SeqOf{
					OpticTypeExpr: ot,
					A:             reflect.TypeFor[A](),
					B:             reflect.TypeFor[A](),
					I:             reflect.TypeFor[Void](),
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The SeqOfP polymorphic combinator focuses on a [iter.Seq] of all the elements in the given optic.
//
// Under modification this seq can be modified using standard operations and will be rebuilt into the original data structure.
// If the modified seq contains fewer elements then [ErrUnsafeMissingElement] will be returned
// If the modified seq contains more elements they will be ignored.
//
// See:
//   - [SeqEOfP] for a polymorphic version.
//   - [ISeqOf] for an indexed version.
//   - [SeqIEOf] for a [SeqIE] version.
func SeqOfP[I, S, T, A, B, RETI, RW, DIR any, ERR TPure](o Optic[I, S, T, A, B, RETI, RW, DIR, ERR]) Optic[Void, S, T, iter.Seq[A], iter.Seq[B], ReturnOne, RW, UniDir, Err] {
	return CombiLens[RW, Err, Void, S, T, iter.Seq[A], iter.Seq[B]](
		func(ctx context.Context, source S) (Void, iter.Seq[A], error) {
			ret := func(yield func(A) bool) {
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					_, a, focusErr := val.Get()
					if focusErr != nil {
						panic(focusErr) //ERR is TPure
					}
					return yield(a)
				})
			}
			return Void{}, ret, nil
		},
		func(ctx context.Context, setVal iter.Seq[B], source S) (T, error) {

			next, stop := iter.Pull(setVal)
			defer stop()
			ret, err := o.AsModify()(ctx, func(index I, focus A) (B, error) {

				val, ok := next()
				if !ok {
					var b B
					return b, ErrUnsafeMissingElement
				} else {
					return val, ctx.Err()
				}
			}, source)

			return ret, JoinCtxErr(ctx, err)
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.SeqOf{
					OpticTypeExpr: ot,
					A:             reflect.TypeFor[A](),
					B:             reflect.TypeFor[B](),
					I:             reflect.TypeFor[Void](),
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The SeqIOf combinator focuses on a [iter.Seq] of all the elements in the given optic.
//
// Under modification this seq can be modified using standard operations and will be rebuilt into the original data structure.
// If the modified seq contains fewer elements the result will use values from the original source.
// If the modified seq contains more elements they will be ignored.
//
// See:
//   - [SeqEOfP] for a polymorphic version.
//   - [ISeqOf] for an indexed version.
//   - [SeqIEOf] for a [SeqIE] version.
func SeqIOf[I, S, T, A, RETI, RW, DIR any, ERR TPure](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR]) Optic[Void, S, T, SeqI[I, A], SeqI[I, A], ReturnOne, RW, UniDir, ERR] {

	return CombiLens[RW, ERR, Void, S, T, SeqI[I, A], SeqI[I, A]](
		func(ctx context.Context, source S) (Void, SeqI[I, A], error) {
			ret := func(yield func(I, A) bool) {
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					focusIndex, a, focusErr := val.Get()
					if focusErr != nil {
						panic(focusErr) //ERR is TPure
					}
					return yield(focusIndex, a)
				})
			}
			return Void{}, ret, nil
		},
		func(ctx context.Context, setVal SeqI[I, A], source S) (T, error) {

			next, stop := iter.Pull2(iter.Seq2[I, A](setVal))
			defer stop()
			ret, err := o.AsModify()(ctx, func(index I, focus A) (A, error) {

				_, val, ok := next()
				if !ok {
					return focus, ctx.Err()
				} else {
					return val, ctx.Err()
				}
			}, source)

			return ret, JoinCtxErr(ctx, err)
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.SeqOf{
					OpticTypeExpr: ot,
					A:             reflect.TypeFor[A](),
					B:             reflect.TypeFor[A](),
					I:             reflect.TypeFor[I](),
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The SeqIOfP combinator focuses on a [iter.Seq] of all the elements in the given optic.
//
// Under modification this seq can be modified using standard operations and will be rebuilt into the original data structure.
// If the modified seq contains fewer elements the result will use values from the original source.
// If the modified seq contains more elements they will be ignored.
//
// See:
//   - [SeqEOfP] for a polymorphic version.
//   - [ISeqOf] for an indexed version.
//   - [SeqIEOf] for a [SeqIE] version.
func SeqIOfP[I, S, T, A, B, RETI, RW, DIR any, ERR TPure](o Optic[I, S, T, A, B, RETI, RW, DIR, ERR]) Optic[Void, S, T, SeqI[I, A], SeqI[I, B], ReturnOne, RW, UniDir, Err] {

	return CombiLens[RW, Err, Void, S, T, SeqI[I, A], SeqI[I, B]](
		func(ctx context.Context, source S) (Void, SeqI[I, A], error) {
			ret := func(yield func(I, A) bool) {
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					focusIndex, a, focusErr := val.Get()
					if focusErr != nil {
						panic(focusErr) //ERR is TPure
					}
					return yield(focusIndex, a)
				})
			}
			return Void{}, ret, nil
		},
		func(ctx context.Context, setVal SeqI[I, B], source S) (T, error) {

			next, stop := iter.Pull2(iter.Seq2[I, B](setVal))
			defer stop()
			ret, err := o.AsModify()(ctx, func(index I, focus A) (B, error) {

				_, val, ok := next()
				if !ok {
					var b B
					return b, ErrUnsafeMissingElement
				} else {
					return val, ctx.Err()
				}
			}, source)

			return ret, JoinCtxErr(ctx, err)
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.SeqOf{
					OpticTypeExpr: ot,
					A:             reflect.TypeFor[A](),
					I:             reflect.TypeFor[I](),
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The SeqEOf combinator focuses on a [iter.Seq2] of all the elements in the given optic.
//
// Under modification this seq can be modified using standard operations and will be rebuilt into the original data structure.
// If the modified seq contains fewer elements the result will use values from the original source.
// If the modified seq contains more elements they will be ignored.
//
// See:
//   - [SeqEOfP] for a polymorphic version.
//   - [ISeqOf] for an indexed version.
//   - [SeqIEOf] for a [SeqIE] version.
func SeqEOf[I, S, T, A, RETI, RW, DIR, ERR any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR]) Optic[Void, S, T, SeqE[A], SeqE[A], ReturnOne, RW, UniDir, ERR] {
	return CombiLens[RW, ERR, Void, S, T, SeqE[A], SeqE[A]](
		func(ctx context.Context, source S) (Void, SeqE[A], error) {
			ret := func(yield func(val ValueE[A]) bool) {
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					_, a, focusErr := val.Get()
					return yield(ValE(a, JoinCtxErr(ctx, focusErr)))
				})
			}
			return Void{}, ret, nil
		},
		func(ctx context.Context, setVal SeqE[A], source S) (T, error) {

			next, stop := iter.Pull(iter.Seq[ValueE[A]](setVal))
			defer stop()
			ret, err := o.AsModify()(ctx, func(index I, focus A) (A, error) {
				v, ok := next()
				val, err := v.Get()
				if !ok {
					return focus, ctx.Err()
				} else {
					return val, JoinCtxErr(ctx, err)
				}
			}, source)

			return ret, JoinCtxErr(ctx, err)
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.SeqOf{
					OpticTypeExpr: ot,
					A:             reflect.TypeFor[A](),
					B:             reflect.TypeFor[A](),
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The SeqEOfP combinator focuses on a polymorphic [iter.Seq2] of all the elements in the given optic.
//
// Under modification this seq can be modified using standard operations and will be rebuilt into the original data structure.
// If the modified seq contains fewer elements then [ErrUnsafeMissingElement] will be returned
// If the modified seq contains more elements they will be ignored.
//
// See:
//   - [SeqEOf] for a safe non polymorphic version.
//   - [ISeqOfP] for an indexed version.
func SeqEOfP[I, S, T, A, B, RETI, RW, DIR, ERR any](o Optic[I, S, T, A, B, RETI, RW, DIR, ERR]) Optic[Void, S, T, SeqE[A], SeqE[B], ReturnOne, RW, UniDir, Err] {
	return CombiLens[RW, Err, Void, S, T, SeqE[A], SeqE[B]](
		func(ctx context.Context, source S) (Void, SeqE[A], error) {
			ret := func(yield func(ValueE[A]) bool) {
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					_, a, focusErr := val.Get()
					return yield(ValE(a, JoinCtxErr(ctx, focusErr)))
				})
			}
			return Void{}, ret, nil
		},
		func(ctx context.Context, focus SeqE[B], source S) (T, error) {

			next, stop := iter.Pull(iter.Seq[ValueE[B]](focus))
			defer stop()
			ret, err := o.AsModify()(ctx, func(index I, focus A) (B, error) {

				v, ok := next()
				val, err := v.Get()
				if !ok {
					var b B
					return b, JoinCtxErr(ctx, ErrUnsafeMissingElement)
				} else {
					return val, JoinCtxErr(ctx, err)
				}
			}, source)

			return ret, JoinCtxErr(ctx, err)
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.SeqOf{
					OpticTypeExpr: ot,
					A:             reflect.TypeFor[A](),
					B:             reflect.TypeFor[B](),
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The SeqIEOf combinator focuses on a [SeqIE] of all the elements in the given optic.
//
// Under modification this seq can be modified using standard operations and will be rebuilt into the original data structure.
// If the modified seq contains fewer elements the result will use values from the original source.
// If the modified seq contains more elements they will be ignored.
//
// See:
//   - [SeqIEOfP] for a polymorphic version.
func SeqIEOf[I, S, T, A, RETI, RW, DIR, ERR any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR]) Optic[Void, S, T, SeqIE[I, A], SeqIE[I, A], ReturnOne, RW, UniDir, ERR] {
	return CombiLens[RW, ERR, Void, S, T, SeqIE[I, A], SeqIE[I, A]](
		func(ctx context.Context, source S) (Void, SeqIE[I, A], error) {
			ret := o.AsIter()(ctx, source)
			return Void{}, ret, nil
		},
		func(ctx context.Context, focus SeqIE[I, A], source S) (T, error) {

			next, stop := PullIE(focus)
			defer stop()
			ret, err := o.AsModify()(ctx, func(index I, focus A) (A, error) {

				_, val, err, ok := next()
				if !ok {
					return focus, ctx.Err()
				} else {
					return val, JoinCtxErr(ctx, err)
				}
			}, source)

			return ret, JoinCtxErr(ctx, err)
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.SeqOf{
					OpticTypeExpr: ot,
					A:             reflect.TypeFor[A](),
					I:             reflect.TypeFor[I](),
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The SeqIEOfP combinator focuses on a polymorphic [SeqIE] of all the elements in the given optic.
//
// Under modification this seq can be modified using standard operations and will be rebuilt into the original data structure.
// If the modified seq contains fewer elements then [ErrUnsafeMissingElement] will be returned
// If the modified seq contains more elements they will be ignored.
//
// See:
//   - [SeqIEOf] for a safe non polymorphic version.
func SeqIEOfP[I, S, T, A, B, RETI, RW, DIR, ERR any](o Optic[I, S, T, A, B, RETI, RW, DIR, ERR]) Optic[Void, S, T, SeqIE[I, A], SeqIE[I, B], ReturnOne, RW, UniDir, Err] {

	return CombiLens[RW, Err, Void, S, T, SeqIE[I, A], SeqIE[I, B]](
		func(ctx context.Context, source S) (Void, SeqIE[I, A], error) {
			ret := o.AsIter()(ctx, source)
			return Void{}, ret, nil
		},
		func(ctx context.Context, focus SeqIE[I, B], source S) (T, error) {

			next, stop := PullIE(focus)
			defer stop()
			ret, err := o.AsModify()(ctx, func(index I, focus A) (B, error) {

				_, val, err, ok := next()
				if !ok {
					var b B
					return b, JoinCtxErr(ctx, ErrUnsafeMissingElement)
				} else {
					return val, JoinCtxErr(ctx, err)
				}
			}, source)

			return ret, JoinCtxErr(ctx, err)
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.SeqOf{
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
