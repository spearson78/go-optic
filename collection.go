package optic

import (
	"context"
	"fmt"
	"iter"
	"reflect"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/spearson78/go-optic/expr"
)

// Collection represents a container of multiple indexed values.
type Collection[I, A, ERR any] interface {
	//Equals[Collection[I, A, ERR]]
	//EqualsAny
	AsIter() func(ctx context.Context) SeqIE[I, A]
	AsIxGet() func(ctx context.Context, index I) SeqIE[I, A]
	AsIxMatch() func(a, b I) bool
	AsLengthGetter() func(ctx context.Context) (int, error)
	ErrType() ERR
	AsExpr() expr.SeqExpr
	String() string
}

// Constructor for a non index aware pure [Collection]
//
// An int index is generated automatically.
//
// See:
//   - [ColIE] for an index aware impure version.
//   - [ColE] for an impure version.
//   - [ColI] for an index aware pure version.
//   - [ValCol] for a varargs/slice based version.
//
// If any of the given functions are nil then a default implementation will be provided.
func Col[A any](
	seq iter.Seq[A],
	lengthGetter func() int,
) Collection[int, A, Pure] {

	var lengthgetf func(ctx context.Context) (int, error)
	if lengthGetter != nil {
		lengthgetf = func(ctx context.Context) (int, error) {
			return lengthGetter(), nil
		}
	}

	var reSeq func(ctx context.Context) SeqIE[int, A]
	if seq != nil {
		reSeq = func(ctx context.Context) SeqIE[int, A] {
			return func(yield func(ValueIE[int, A]) bool) {
				i := 0
				seq(func(v A) bool {
					cont := yield(ValIE(i, v, ctx.Err()))
					i++
					return cont
				})
			}
		}
	}

	return ColIE[Pure](
		reSeq,
		nil,
		IxMatchComparable[int](),
		lengthgetf,
	)
}

// Constructor for an impure [Collection]
//
// See:
//   - [ColIE] for an index aware impure version.
//   - [ColI] for an index aware pure version.
//   - [Col] for a non index aware pure version.
//   - [ValColE] for a varargs/slice based version.
//
// If any of the given functions are nil then a default implementation will be provided.
func ColE[A any](
	seq func(ctx context.Context) SeqE[A],
	lengthGetter func(ctx context.Context) (int, error),
) Collection[int, A, Err] {
	return ColIE[Err](
		func(ctx context.Context) SeqIE[int, A] {
			return func(yield func(ValueIE[int, A]) bool) {
				i := 0
				seq(ctx)(func(v ValueE[A]) bool {
					focus, err := v.Get()
					cont := yield(ValIE(i, focus, err))
					i++
					return cont
				})
			}
		},
		nil,
		IxMatchComparable[int](),
		lengthGetter,
	)
}

// Constructor for an impure [Collection]
//
// See:
//   - [ColIE] for an index aware impure version.
//   - [ColE] for an impure version.
//   - [Col] for a non index aware pure version.
//   - [ValColI] for a varargs/slice based version.
//
// If any of the given functions are nil then a default implementation will be provided.
func ColI[I, A any](
	seq iter.Seq2[I, A],
	ixget func(index I) iter.Seq2[I, A],
	ixMatch func(a, b I) bool,
	lengthGetter func() int,
) Collection[I, A, Pure] {

	var seqf func(ctx context.Context) SeqIE[I, A]
	if seq != nil {
		seqf = func(ctx context.Context) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				seq(func(i I, v A) bool {
					return yield(ValIE(i, v, ctx.Err()))
				})
			}
		}
	}

	var ixgetf func(ctx context.Context, index I) SeqIE[I, A]
	if ixget != nil {
		ixgetf = func(ctx context.Context, index I) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				ixget(index)(func(i I, v A) bool {
					return yield(ValIE(i, v, ctx.Err()))
				})
			}
		}
	}

	var lengthgetf func(ctx context.Context) (int, error)
	if lengthGetter != nil {
		lengthgetf = func(ctx context.Context) (int, error) {
			return lengthGetter(), nil
		}
	}

	return ColIE[Pure](
		seqf,
		ixgetf,
		ixMatch,
		lengthgetf,
	)
}

// Constructor for an index aware impure [Collection]
//
// See:
//   - [ColE] for a non index aware impure version.
//   - [ColI] for an index aware pure version.
//   - [Col] for a non index aware pure version.
//   - [ValColIE] for a varargs/slice based version.
//
// If any of the given functions are nil then a default implementation will be provided.
func ColIE[ERR, I, A any](
	seq func(ctx context.Context) SeqIE[I, A],
	ixget func(ctx context.Context, index I) SeqIE[I, A],
	ixMatch func(a, b I) bool,
	lengthGetter func(ctx context.Context) (int, error),
) Collection[I, A, ERR] {

	if seq == nil {
		seq = func(ctx context.Context) SeqIE[I, A] { return func(yield func(val ValueIE[I, A]) bool) {} }
	}

	if ixMatch == nil {
		ixMatch = func(a, b I) bool {
			return false
		}
		ixget = func(ctx context.Context, index I) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {}
		}
	}

	if ixget == nil {
		ixget = func(ctx context.Context, index I) SeqIE[I, A] {
			return func(yield func(ValueIE[I, A]) bool) {
				seq(ctx)(func(val ValueIE[I, A]) bool {
					focusIndex, focus, err := val.Get()
					err = JoinCtxErr(ctx, err)
					if err != nil || ixMatch(index, focusIndex) {
						return yield(ValIE(focusIndex, focus, err))
					}
					return true
				})
			}
		}
	}

	if lengthGetter == nil {
		lengthGetter = func(ctx context.Context) (int, error) {
			i := 0
			var retErr error
			seq(ctx)(func(val ValueIE[I, A]) bool {
				err := val.Error()
				if err != nil {
					retErr = err
					return false
				}
				i++
				return true
			})
			if retErr != nil {
				return 0, retErr
			}
			return i, nil
		}
	}

	return collection[I, A, ERR]{
		seq:          seq,
		ixMatch:      ixMatch,
		ixGet:        ixget,
		lengthGetter: lengthGetter,
	}
}

type collection[I any, A, ERR any] struct {
	seq          func(ctx context.Context) SeqIE[I, A]
	ixGet        func(ctx context.Context, index I) SeqIE[I, A]
	ixMatch      func(a, b I) bool
	lengthGetter func(ctx context.Context) (int, error)
}

func (n collection[I, A, ERR]) AsIter() func(ctx context.Context) SeqIE[I, A] {
	return n.seq
}

func (n collection[I, A, ERR]) AsIxGet() func(ctx context.Context, index I) SeqIE[I, A] {
	return n.ixGet
}

func (n collection[I, A, ERR]) AsIxMatch() func(a, b I) bool {
	return n.ixMatch
}

func (n collection[I, A, ERR]) AsLengthGetter() func(ctx context.Context) (int, error) {
	return n.lengthGetter
}

func (n collection[I, A, ERR]) ErrType() ERR {
	var err ERR
	return err
}

func (n collection[I, A, ERR]) AsExpr() expr.SeqExpr {
	return func(ctx context.Context, yield func(expr.SeqExprValue) bool) {
		n.seq(ctx)(func(val ValueIE[I, A]) bool {
			index, focus, err := val.Get()
			return yield(expr.SeqExprValue{
				Index:    index,
				Value:    focus,
				ValuePtr: &focus,
				Error:    err,
			})
		})
	}
}

func (s collection[I, A, ERR]) String() string {
	var sb strings.Builder

	first := true
	sb.WriteString("Col[")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	s.seq(ctx)(func(val ValueIE[I, A]) bool {
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

	sb.WriteString("]")

	return sb.String()
}

// EqCol returns a [Predicate] that is satisfied if the elements of the focused [Collection] are all equal to (==) the elements of the provided constant [Collection] and element [Predicate].
func EqCol[I comparable, A any, ERR any, PERR TPure](right Collection[I, A, ERR], eq Predicate[lo.Tuple2[A, A], PERR]) Optic[Void, Collection[I, A, ERR], Collection[I, A, ERR], bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return OpT2Bind(EqColT2[I, A, ERR](eq), right)
}

// EqColT2 returns an [BinaryOp] that is satisfied if every element and index in Collection A and Collection B in focused order are equal.
//
// See:
//   - [EqCol] for a unary version.
//   - [EqColT2I] for a version that supports arbitrary index types
func EqColT2[I comparable, A, ERR any, PERR TPure](eq Predicate[lo.Tuple2[A, A], PERR]) Optic[Void, lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return EqColT2I[I, A, ERR](EqT2[I](), eq)
}

// EqColT2 returns an index aware [BinaryOp] that is satisfied if Collection A == Collection B.
//
// See:
//   - [EqColT2] for a version that simpler version that supports only comparable index types.
func EqColT2I[I any, A, ERR any, IERR, PERR TPure](ixMatch Predicate[lo.Tuple2[I, I], IERR], eq Predicate[lo.Tuple2[A, A], PERR]) Optic[Void, lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], lo.Tuple2[Collection[I, A, ERR], Collection[I, A, ERR]], bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return EqT2Of(
		TraverseColIE[I, A, ERR](PredToIxMatch(ixMatch)),
		eq,
	)
}

func unsafeColSourceErr[SERR, I, J, S, T, A, B, RET, RW, DIR, ERR, OSERR, OTERR any](o Optic[J, Collection[I, S, OSERR], Collection[I, T, OTERR], A, B, RET, RW, DIR, ERR]) Optic[J, Collection[I, S, SERR], Collection[I, T, SERR], A, B, RET, RW, DIR, SERR] {

	return Omni[J, Collection[I, S, SERR], Collection[I, T, SERR], A, B, RET, RW, DIR, SERR](
		func(ctx context.Context, source Collection[I, S, SERR]) (J, A, error) {
			return o.AsGetter()(ctx, ColIE[OSERR](source.AsIter(), source.AsIxGet(), source.AsIxMatch(), source.AsLengthGetter()))
		},
		func(ctx context.Context, focus B, source Collection[I, S, SERR]) (Collection[I, T, SERR], error) {
			ret, err := o.AsSetter()(ctx, focus, ColIE[OSERR](source.AsIter(), source.AsIxGet(), source.AsIxMatch(), source.AsLengthGetter()))
			if err != nil {
				return ValColIE[I, T, SERR](nil), err
			}
			return ColIE[SERR](ret.AsIter(), ret.AsIxGet(), ret.AsIxMatch(), ret.AsLengthGetter()), err
		},
		func(ctx context.Context, source Collection[I, S, SERR]) SeqIE[J, A] {
			return o.AsIter()(ctx, ColIE[OSERR](source.AsIter(), source.AsIxGet(), source.AsIxMatch(), source.AsLengthGetter()))
		},
		func(ctx context.Context, source Collection[I, S, SERR]) (int, error) {
			return o.AsLengthGetter()(ctx, ColIE[OSERR](source.AsIter(), source.AsIxGet(), source.AsIxMatch(), source.AsLengthGetter()))
		},
		func(ctx context.Context, fmap func(index J, focus A) (B, error), source Collection[I, S, SERR]) (Collection[I, T, SERR], error) {
			ret, err := o.AsModify()(ctx, fmap, ColIE[OSERR](source.AsIter(), source.AsIxGet(), source.AsIxMatch(), source.AsLengthGetter()))
			if err != nil {
				return ValColIE[I, T, SERR](nil), err
			}
			return ColIE[SERR](ret.AsIter(), ret.AsIxGet(), ret.AsIxMatch(), ret.AsLengthGetter()), err
		},
		func(ctx context.Context, index J, source Collection[I, S, SERR]) SeqIE[J, A] {
			return o.AsIxGetter()(ctx, index, ColIE[OSERR](source.AsIter(), source.AsIxGet(), source.AsIxMatch(), source.AsLengthGetter()))
		},
		o.AsIxMatch(),
		func(ctx context.Context, focus B) (Collection[I, T, SERR], error) {
			ret, err := o.AsReverseGetter()(ctx, focus)
			if err != nil {
				return ValColIE[I, T, SERR](nil), err
			}
			return ColIE[SERR](ret.AsIter(), ret.AsIxGet(), ret.AsIxMatch(), ret.AsLengthGetter()), err
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.ColSourceErr{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

func unsafeColFocusErr[ERR, I, J, S, T, A, B, RET, RW, DIR, OERR, OAERR, OBERR any](o Optic[J, S, T, Collection[I, A, OAERR], Collection[I, B, OBERR], RET, RW, DIR, OERR]) Optic[J, S, T, Collection[I, A, ERR], Collection[I, B, ERR], RET, RW, DIR, ERR] {

	return Omni[J, S, T, Collection[I, A, ERR], Collection[I, B, ERR], RET, RW, DIR, ERR](
		func(ctx context.Context, source S) (J, Collection[I, A, ERR], error) {
			i, ret, err := o.AsGetter()(ctx, source)
			return i, ColIE[ERR](ret.AsIter(), ret.AsIxGet(), ret.AsIxMatch(), ret.AsLengthGetter()), err
		},
		func(ctx context.Context, focus Collection[I, B, ERR], source S) (T, error) {
			return o.AsSetter()(ctx, ColIE[OBERR](focus.AsIter(), focus.AsIxGet(), focus.AsIxMatch(), focus.AsLengthGetter()), source)
		},
		func(ctx context.Context, source S) SeqIE[J, Collection[I, A, ERR]] {
			return func(yield func(ValueIE[J, Collection[I, A, ERR]]) bool) {
				o.AsIter()(ctx, source)(func(val ValueIE[J, Collection[I, A, OAERR]]) bool {
					index, focus, err := val.Get()
					return yield(ValIE(index, ColIE[ERR](focus.AsIter(), focus.AsIxGet(), focus.AsIxMatch(), focus.AsLengthGetter()), err))
				})
			}
		},
		o.AsLengthGetter(),
		func(ctx context.Context, fmap func(index J, focus Collection[I, A, ERR]) (Collection[I, B, ERR], error), source S) (T, error) {
			return o.AsModify()(ctx, func(index J, focus Collection[I, A, OAERR]) (Collection[I, B, OBERR], error) {
				ret, err := fmap(index, ColIE[ERR](focus.AsIter(), focus.AsIxGet(), focus.AsIxMatch(), focus.AsLengthGetter()))
				return ColIE[OBERR](ret.AsIter(), ret.AsIxGet(), ret.AsIxMatch(), ret.AsLengthGetter()), err
			}, source)
		},
		func(ctx context.Context, index J, source S) SeqIE[J, Collection[I, A, ERR]] {
			return func(yield func(ValueIE[J, Collection[I, A, ERR]]) bool) {
				o.AsIxGetter()(ctx, index, source)(func(val ValueIE[J, Collection[I, A, OAERR]]) bool {
					index, focus, err := val.Get()
					return yield(ValIE(index, ColIE[ERR](focus.AsIter(), focus.AsIxGet(), focus.AsIxMatch(), focus.AsLengthGetter()), err))
				})
			}
		},
		o.AsIxMatch(),
		func(ctx context.Context, focus Collection[I, B, ERR]) (T, error) {
			return o.AsReverseGetter()(ctx, ColIE[OBERR](focus.AsIter(), focus.AsIxGet(), focus.AsIxMatch(), focus.AsLengthGetter()))
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.ColFocusErr{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

func unsafeColSourceFocusErr[ERR, I, J, S, T, A, B, RET, RW, DIR, OERR any, OSERR, OTERR, OAERR, OBERR any](o Optic[J, Collection[I, S, OSERR], Collection[I, T, OTERR], Collection[I, A, OAERR], Collection[I, B, OBERR], RET, RW, DIR, OERR]) Optic[J, Collection[I, S, ERR], Collection[I, T, ERR], Collection[I, A, ERR], Collection[I, B, ERR], RET, RW, DIR, ERR] {

	return Omni[J, Collection[I, S, ERR], Collection[I, T, ERR], Collection[I, A, ERR], Collection[I, B, ERR], RET, RW, DIR, ERR](
		func(ctx context.Context, source Collection[I, S, ERR]) (J, Collection[I, A, ERR], error) {
			i, ret, err := o.AsGetter()(ctx, ColIE[OSERR](source.AsIter(), source.AsIxGet(), source.AsIxMatch(), source.AsLengthGetter()))
			return i, ColIE[ERR](ret.AsIter(), ret.AsIxGet(), ret.AsIxMatch(), ret.AsLengthGetter()), err
		},
		func(ctx context.Context, focus Collection[I, B, ERR], source Collection[I, S, ERR]) (Collection[I, T, ERR], error) {
			ret, err := o.AsSetter()(ctx, ColIE[OBERR](focus.AsIter(), focus.AsIxGet(), focus.AsIxMatch(), focus.AsLengthGetter()), ColIE[OSERR](source.AsIter(), source.AsIxGet(), source.AsIxMatch(), source.AsLengthGetter()))
			if err != nil {
				return ValColIE[I, T, ERR](nil), err
			}
			return ColIE[ERR](ret.AsIter(), ret.AsIxGet(), ret.AsIxMatch(), ret.AsLengthGetter()), err
		},
		func(ctx context.Context, source Collection[I, S, ERR]) SeqIE[J, Collection[I, A, ERR]] {
			return func(yield func(ValueIE[J, Collection[I, A, ERR]]) bool) {
				o.AsIter()(ctx, ColIE[OSERR](source.AsIter(), source.AsIxGet(), source.AsIxMatch(), source.AsLengthGetter()))(func(val ValueIE[J, Collection[I, A, OAERR]]) bool {
					index, focus, err := val.Get()
					return yield(ValIE(index, ColIE[ERR](focus.AsIter(), focus.AsIxGet(), focus.AsIxMatch(), focus.AsLengthGetter()), err))
				})
			}
		},
		func(ctx context.Context, source Collection[I, S, ERR]) (int, error) {
			return o.AsLengthGetter()(ctx, ColIE[OSERR](source.AsIter(), source.AsIxGet(), source.AsIxMatch(), source.AsLengthGetter()))
		},
		func(ctx context.Context, fmap func(index J, focus Collection[I, A, ERR]) (Collection[I, B, ERR], error), source Collection[I, S, ERR]) (Collection[I, T, ERR], error) {
			ret, err := o.AsModify()(ctx, func(index J, focus Collection[I, A, OAERR]) (Collection[I, B, OBERR], error) {
				ret, err := fmap(index, ColIE[ERR](focus.AsIter(), focus.AsIxGet(), focus.AsIxMatch(), focus.AsLengthGetter()))
				return ColIE[OBERR](ret.AsIter(), ret.AsIxGet(), ret.AsIxMatch(), ret.AsLengthGetter()), err
			}, ColIE[OSERR](source.AsIter(), source.AsIxGet(), source.AsIxMatch(), source.AsLengthGetter()))
			return ColIE[ERR](ret.AsIter(), ret.AsIxGet(), ret.AsIxMatch(), ret.AsLengthGetter()), err
		},
		func(ctx context.Context, index J, source Collection[I, S, ERR]) SeqIE[J, Collection[I, A, ERR]] {
			return func(yield func(ValueIE[J, Collection[I, A, ERR]]) bool) {
				o.AsIxGetter()(ctx, index, ColIE[OSERR](source.AsIter(), source.AsIxGet(), source.AsIxMatch(), source.AsLengthGetter()))(func(val ValueIE[J, Collection[I, A, OAERR]]) bool {
					index, focus, err := val.Get()
					return yield(ValIE(index, ColIE[ERR](focus.AsIter(), focus.AsIxGet(), focus.AsIxMatch(), focus.AsLengthGetter()), err))
				})
			}
		},
		o.AsIxMatch(),
		func(ctx context.Context, focus Collection[I, B, ERR]) (Collection[I, T, ERR], error) {
			ret, err := o.AsReverseGetter()(ctx, ColIE[OBERR](focus.AsIter(), focus.AsIxGet(), focus.AsIxMatch(), focus.AsLengthGetter()))
			return ColIE[ERR](ret.AsIter(), ret.AsIxGet(), ret.AsIxMatch(), ret.AsLengthGetter()), err
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.ColSourceFocusErr{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// ColSourceFocusPure returns the input [Optic] with the SERR as the source and focus [Collection] ERR.
//
// This function is intended when combinining optics that use a mix of [Pure] and [Err] [Collection]s.
//
// See:
//   - [ColSourceFocusPureP] for a polymorphic version
//   - [ColSourcePure] for a version that only modifies the source [Collection] ERR.
//   - [ColFocusPure] for a version that only modifies the focus [Collection] ERR
//   - [ColSourceFocusErr] for a version that modifies any [Collection] to Err
func ColSourceFocusPure[I, J, S, T, A, B, RET, RW, DIR, ERR any, OSERR, OTERR, OAERR, OBERR TPure](o Optic[J, Collection[I, S, OSERR], Collection[I, T, OTERR], Collection[I, A, OAERR], Collection[I, B, OBERR], RET, RW, DIR, ERR]) Optic[J, Collection[I, S, Pure], Collection[I, T, Pure], Collection[I, A, Pure], Collection[I, B, Pure], RET, RW, DIR, Pure] {
	return unsafeColSourceFocusErr[Pure](o)
}

// ColSourceFocusErr returns the input [Optic] with Err as the source [Collection] SERR.
//
// This function is intended when combinining optics that use a mix of [Pure] and [Err] [Collection]s.
//
// See:
//   - [ColSourceFocus] for a version that modifes pure source and focus [Collection] ERR.
//   - [ColSourceErrP] for a polymorphic version
//   - [ColSourceFocusPure] for a version that modifies both the source and focus [Collection] ERR.
//   - [ColFocusPure] for a version that only modifies the focus [Collection] ERR
func ColSourceFocusErr[I, J, S, T, A, B, RET, RW, DIR, ERR any, OSERR, OTERR, OAERR, OBERR any](o Optic[J, Collection[I, S, OSERR], Collection[I, T, OTERR], Collection[I, A, OAERR], Collection[I, B, OBERR], RET, RW, DIR, ERR]) Optic[J, Collection[I, S, Err], Collection[I, T, Err], Collection[I, A, Err], Collection[I, B, Err], RET, RW, DIR, Err] {
	return unsafeColSourceFocusErr[Err](o)
}

// ColSourcePure returns the pure input [Optic] with the SERR as the source [Collection] ERR.
//
// This function is intended when combinining optics that use a mix of [Pure] and [Err] [Collection]s.
//
// See:
//   - [ColSourceErrP] for a polymorphic version
//   - [ColSourceFocusPure] for a version that modifies both the source and focus [Collection] ERR.
//   - [ColFocusPure] for a version that only modifies the focus [Collection] ERR
//   - [ColSourceErr] for a version that modifies any source [Collection] to Err
func ColSourcePure[I, J, S, T, A, B, RET, RW, DIR any, ERR TPure, OSERR TPure](o Optic[J, Collection[I, S, OSERR], Collection[I, T, OSERR], A, B, RET, RW, DIR, ERR]) Optic[J, Collection[I, S, Pure], Collection[I, T, Pure], A, B, RET, RW, DIR, Pure] {
	return unsafeColSourceErr[Pure](o)
}

// ColSourceErr returns the input [Optic] with Err as the source [Collection] SERR.
//
// This function is intended when combinining optics that use a mix of [Pure] and [Err] [Collection]s.
//
// See:
//   - [ColSourceErrP] for a polymorphic version
//   - [ColSourceFocusPure] for a version that modifies both the source and focus [Collection] ERR.
//   - [ColFocusPure] for a version that only modifies the focus [Collection] ERR
//   - [ColSourcePure] for a version that modifies pure source [Collection] to an arbitrary ERR
func ColSourceErr[I, J, S, T, A, B, RET, RW, DIR, ERR any, SERR, TERR any](o Optic[J, Collection[I, S, SERR], Collection[I, T, TERR], A, B, RET, RW, DIR, ERR]) Optic[J, Collection[I, S, Err], Collection[I, T, Err], A, B, RET, RW, DIR, Err] {
	return unsafeColSourceErr[Err](o)
}

// ColFocusPure returns the pure input [Optic] with the given focus [Collection] ERR.
//
// This function is intended when combinining optics that use a mix of [Pure] and [Err] [Collection]s.
//
// See:
//   - [ColFocusPure] for a non polymorphic version
//   - [ColSourceFocusPureP] for a version that modifies both the source and focus [Collection] ERR.
//   - [ColSourceErrP] for a version that only modifies the source [Collection] ERR
//   - [ColSourceErr] for a version that modifies any source [Collection] to Err
func ColFocusPure[I, J, S, T, A, B, RET, RW, DIR, ERR any, OAERR, OBERR TPure](o Optic[J, S, T, Collection[I, A, OAERR], Collection[I, B, OBERR], RET, RW, DIR, ERR]) Optic[J, S, T, Collection[I, A, Pure], Collection[I, B, Pure], RET, RW, DIR, Pure] {
	return unsafeColFocusErr[Pure](o)
}

// ColFocusErr returns the input [Optic] with Err as the focus [Collection] SERR.
//
// This function is intended when combinining optics that use a mix of [Pure] and [Err] [Collection]s.
//
// See:
//   - [ColFocusPure] for a version that modifes pure source and focus [Collection] ERR.
//   - [ColFocusErrP] for a polymorphic version
//   - [ColSourceFocusPure] for a version that modifies both the source and focus [Collection] ERR.
func ColFocusErr[I, J, S, T, A, B, RET, RW, DIR, ERR any, OAERR, OBERR any](o Optic[J, S, T, Collection[I, A, OAERR], Collection[I, B, OBERR], RET, RW, DIR, ERR]) Optic[J, S, T, Collection[I, A, Err], Collection[I, B, Err], RET, RW, DIR, Err] {
	return unsafeColFocusErr[Err](o)
}

// ColPure modifies a pure [Collection] to a [Collection] with an arbitrary ERR type.
func ColPure[I, A any, ERR TPure](s Collection[I, A, ERR]) Collection[I, A, Pure] {
	return collection[I, A, Pure]{
		seq:          s.AsIter(),
		ixGet:        s.AsIxGet(),
		ixMatch:      s.AsIxMatch(),
		lengthGetter: s.AsLengthGetter(),
	}
}

// ColErr modifies any [Collection]s to Err
func ColErr[I, A any, ERR any](s Collection[I, A, ERR]) Collection[I, A, Err] {
	return collection[I, A, Err]{
		seq:          s.AsIter(),
		ixGet:        s.AsIxGet(),
		ixMatch:      s.AsIxMatch(),
		lengthGetter: s.AsLengthGetter(),
	}
}

// ValCol returns a [Collection] for the given values with an int index type.
//
// See:
//   - [ValCol] for a version with an arbitrary numeric index type.
//   - [ValColI] for a version with a custom index
//   - [ValColE] for a version with error support
//   - [ValColIE] for a version with custom index and error support
func ValCol[A any](values ...A) Collection[int, A, Pure] {
	return ColI(
		func(yield func(int, A) bool) {
			for i, v := range values {
				if !yield(i, v) {
					break
				}
			}
		},
		func(index int) iter.Seq2[int, A] {
			return func(yield func(index int, focus A) bool) {
				if index >= 0 && index < len(values) {
					yield(index, values[index])
				}
			}
		},
		IxMatchComparable[int](),
		func() int {
			return len(values)
		},
	)
}

// ColOfValues returns an [Collection] for the given index and values
//
// See:
//   - [ValCol] for a version with a automatic numeric index type
//   - [ValColE] for a version with error support
//   - [ValColIE] for a version with custom index and error support
func ValColI[I any, A any](ixMatch func(a, b I) bool, values ...ValueI[I, A]) Collection[I, A, Pure] {
	return ColI(
		func(yield func(I, A) bool) {
			for _, v := range values {
				if !yield(v.index, v.value) {
					break
				}
			}
		},
		nil,
		ixMatch,
		func() int {
			return len(values)
		},
	)
}

// ValColE returns a [Collection] for the given values and erros with an int index type.
//
// See:
//   - [ValCol] for a version with a automatic numeric index type
//   - [ValColI] for a version with a custom index
//   - [ValColE] for a version with error support
//   - [ValColIE] for a version with custom index and error support
func ValColE[A any](values ...ValueE[A]) Collection[int, A, Err] {
	return ColE(
		func(ctx context.Context) SeqE[A] {
			return func(yield func(val ValueE[A]) bool) {
				for _, v := range values {
					res, err := v.Get()
					err = JoinCtxErr(ctx, err)
					if !yield(ValE(res, err)) {
						break
					}
				}
			}
		},
		nil,
	)
}

// ValColIE returns an [Collection] for the given index,value and errors values.
//
// See:
//   - [ValCol] for a simpler version with an int index
//   - [ValCol] for a version with a automatic numeric index type
//   - [ValColI] for a version with a custom index
//   - [ValColE] for a version with error support
func ValColIE[I any, A, ERR any](ixMatch func(a, b I) bool, values ...ValueIE[I, A]) Collection[I, A, ERR] {
	return ColIE[ERR](
		func(ctx context.Context) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				for _, v := range values {
					if !yield(ValIE(v.index, v.value, JoinCtxErr(ctx, v.Error()))) {
						break
					}
				}
			}
		},
		nil,
		ixMatch,
		nil,
	)
}

// RangeCol return a [Collection] that yields the values between from and to inclusive.
func RangeCol[A Real](from, to A) Collection[int, A, Pure] {
	return ColIE[Pure](
		func(ctx context.Context) SeqIE[int, A] {
			return func(yield func(ValueIE[int, A]) bool) {
				index := 0
				for i := from; i <= to; i++ {
					if !yield(ValIE(index, i, ctx.Err())) {
						break
					}
					index++
				}
			}
		},
		func(ctx context.Context, index int) SeqIE[int, A] {
			if A(index) >= from && A(index) <= to {
				return func(yield func(ValueIE[int, A]) bool) {
					yield(ValIE(index-int(from), A(index), nil))
				}
			} else {
				return func(yield func(ValueIE[int, A]) bool) {}
			}
		},
		IxMatchComparable[int](),
		func(ctx context.Context) (int, error) {
			return int((to - from) + 1), nil
		},
	)
}

// TraverseCol returns an [Traversal] that focuses the elements of a collection.
//
// See:
//   - [TraverseColI] for a version that supports arbitrary index types.
//   - [TraverseColE] for an impure version
//   - [TraverseColP] for a polymorphic version
//   - [TraverseColEP] for an impure,polymorphic version
//   - [TraverseColIEP] for am impure,polymorphic version that supports arbitrary index types.
func TraverseCol[I comparable, A any]() Optic[I, Collection[I, A, Pure], Collection[I, A, Pure], A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseColI[I, A](IxMatchComparable[I]())
}

// TraverseColI returns an [Traversal] that focuses the elements of a collection.
//
// See:
//   - [TraverseCol] for a simpler version that only supports comparable index types.
//   - [TraverseColE] for an impure version
//   - [TraverseColP] for a polymorphic version
//   - [TraverseColEP] for an impure,polymorphic version
//   - [TraverseColIEP] for am impure,polymorphic version that supports arbitrary index types.
func TraverseColI[I, A any](ixmatch func(a, b I) bool) Optic[I, Collection[I, A, Pure], Collection[I, A, Pure], A, A, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseColIEP[I, A, A, Pure](ixmatch)
}

func TraverseColIE[I, A, ERR any](ixmatch func(a, b I) bool) Optic[I, Collection[I, A, ERR], Collection[I, A, ERR], A, A, ReturnMany, ReadWrite, UniDir, ERR] {
	return TraverseColIEP[I, A, A, ERR](ixmatch)
}

// TraverseColE returns an [Traversal] that focuses the elements of a collection.
//
// The modifyRetainErrors controls whether errors in the collection are maintained. true will return a collection with error entries, false will return the first error.
//
// See:
//   - [TraverseCol] for a simpler version that only supports comparable index types.
//   - [TraverseColI] for a version that supports arbitrary index types.
//   - [TraverseColP] for a polymorphic version
//   - [TraverseColEP] for an impure,polymorphic version
//   - [TraverseColIEP] for am impure,polymorphic version that supports arbitrary index types.
func TraverseColE[I comparable, A any, ERR any]() Optic[I, Collection[I, A, ERR], Collection[I, A, ERR], A, A, ReturnMany, ReadWrite, UniDir, ERR] {
	return TraverseColIEP[I, A, A, ERR](IxMatchComparable[I]())
}

// TraverseColP returns an [Traversal] that focuses the elements of a collection.
//
// See:
//   - [TraverseCol] for a simpler version that only supports comparable index types.
//   - [TraverseColI] for a version that supports arbitrary index types.
//   - [TraverseColE] for an impure version
//   - [TraverseColEP] for an impure,polymorphic version
//   - [TraverseColIEP] for am impure,polymorphic version that supports arbitrary index types.
func TraverseColP[I comparable, A, B any]() Optic[I, Collection[I, A, Pure], Collection[I, B, Pure], A, B, ReturnMany, ReadWrite, UniDir, Pure] {
	return TraverseColIEP[I, A, B, Pure](IxMatchComparable[I]())
}

// TraverseColP returns an [Traversal] that focuses the elements of a collection.
//
// See:
//   - [TraverseCol] for a simpler version that only supports comparable index types.
//   - [TraverseColI] for a version that supports arbitrary index types.
//   - [TraverseColE] for an impure version
//   - [TraverseColEP] for an impure,polymorphic version
//   - [TraverseColIEP] for am impure,polymorphic version that supports arbitrary index types.
func TraverseColEP[I comparable, A, B any, ERR any]() Optic[I, Collection[I, A, ERR], Collection[I, B, ERR], A, B, ReturnMany, ReadWrite, UniDir, ERR] {
	return TraverseColIEP[I, A, B, ERR](IxMatchComparable[I]())
}

// TraverseColIEP returns an [Traversal] that focuses the elements of a collection.
//
// under Modification the returned collection is lazily evaluated. See
//
// See:
//   - [TraverseCol] for a simpler version that only supports comparable index types.
//   - [TraverseColI] for a version that supports arbitrary index types.
//   - [TraverseColE] for an impure version
//   - [TraverseColP] for a polymorphic version
func TraverseColIEP[I, A, B, ERR any](ixmatch func(a, b I) bool) Optic[I, Collection[I, A, ERR], Collection[I, B, ERR], A, B, ReturnMany, ReadWrite, UniDir, ERR] {
	return CombiTraversal[ReturnMany, ReadWrite, ERR, I, Collection[I, A, ERR], Collection[I, B, ERR], A, B](
		func(ctx context.Context, source Collection[I, A, ERR]) SeqIE[I, A] {
			if source == nil {
				return func(yield func(val ValueIE[I, A]) bool) {}
			}

			return source.AsIter()(ctx)
		},
		func(ctx context.Context, source Collection[I, A, ERR]) (int, error) {
			return source.AsLengthGetter()(ctx)
		},
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source Collection[I, A, ERR]) (Collection[I, B, ERR], error) {
			if source == nil {
				return ValColIE[I, B, ERR](ixmatch), nil
			}

			indexTunnel, _ := ctx.Value(indexTunnelKey).(*indexTunnel)

			col := ColIE[ERR, I, B](
				func(ctx context.Context) SeqIE[I, B] {
					return func(yield func(ValueIE[I, B]) bool) {
						source.AsIter()(ctx)(func(val ValueIE[I, A]) bool {
							index, focus, err := val.Get()
							if err != nil {
								var b B
								return yield(ValIE(index, b, err))
							}
							b, err := fmap(index, focus)
							if indexTunnel != nil && indexTunnel.ok {
								index = indexTunnel.index.(I)
								indexTunnel.ok = false
							}
							return yield(ValIE(index, b, err))
						})

					}
				},
				func(ctx context.Context, index I) SeqIE[I, B] {
					return func(yield func(ValueIE[I, B]) bool) {
						source.AsIxGet()(ctx, index)(func(val ValueIE[I, A]) bool {
							index, focus, err := val.Get()
							if err != nil {
								var b B
								return yield(ValIE(index, b, err))
							}
							b, err := fmap(index, focus)
							return yield(ValIE(index, b, err))
						})
					}
				},
				ixmatch,
				source.AsLengthGetter(),
			)

			if indexTunnel != nil {
				return col, nil
			} else {
				return materializeCol(ctx, col)
			}

		},
		func(ctx context.Context, index I, source Collection[I, A, ERR]) SeqIE[I, A] {
			if source == nil {
				return func(yield func(val ValueIE[I, A]) bool) {}
			}

			return source.AsIxGet()(ctx, index)
		},
		ixmatch,
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Traverse{
				OpticTypeExpr: ot,
			}
		}),
	)
}

func materializeCol[I, A, ERR any](ctx context.Context, source Collection[I, A, ERR]) (Collection[I, A, ERR], error) {
	l, err := source.AsLengthGetter()(ctx)
	if err != nil {
		return nil, err
	}
	var retErr error
	ret := make([]ValueIE[I, A], 0, l)

	source.AsIter()(ctx)(func(val ValueIE[I, A]) bool {
		index, focus, err := val.Get()
		if err != nil {
			retErr = err
			return false
		}
		ret = append(ret, ValIE(index, focus, nil))
		return true
	})

	if retErr != nil {
		return nil, retErr
	}

	return ValColIE[I, A, ERR](source.AsIxMatch(), ret...), nil
}

type asIxMatch[I any] interface {
	AsIxMatch() IxMatchFunc[I]
}

// The ColOf combinator focuses on a [Collection] of all the elements in the given optic.
//
// Under modification this collection can be modified and will be rebuilt into the original data structure.
// If the modified collection contains fewer elements the result will use values from the original source.
// If the modified collection contains more elements they will be ignored.
//
// See:
//   - [ColOfP] for a polymorphic version.
//   - [ColyTypeOf] for a version that operates on concrete collection types
func ColOf[I, S, T, A, RETI, RW, DIR, ERR any](o Optic[I, S, T, A, A, RETI, RW, DIR, ERR]) Optic[Void, S, T, Collection[I, A, ERR], Collection[I, A, ERR], ReturnOne, RW, UniDir, ERR] {

	//WithIndex,TraverseColIEP and the Combinator below work together to support modifying indexes during a modify()
	op := WithIndex(o)

	return CombiLens[RW, ERR, Void, S, T, Collection[I, A, ERR], Collection[I, A, ERR]](
		func(ctx context.Context, source S) (Void, Collection[I, A, ERR], error) {
			ret := ColIE[ERR, I, A](
				func(ctx context.Context) SeqIE[I, A] {
					return func(yield func(ValueIE[I, A]) bool) {
						op.AsIter()(ctx, source)(func(val ValueIE[I, ValueI[I, A]]) bool {
							return yield(ValIE(val.value.index, val.value.value, val.err))
						})
					}
				},
				func(ctx context.Context, index I) SeqIE[I, A] {
					return func(yield func(ValueIE[I, A]) bool) {
						op.AsIxGetter()(ctx, index, source)(func(val ValueIE[I, ValueI[I, A]]) bool {
							return yield(ValIE(val.value.index, val.value.value, val.err))
						})
					}
				},
				o.AsIxMatch(),
				func(ctx context.Context) (int, error) {
					return op.AsLengthGetter()(ctx, source)
				},
			)
			return Void{}, ret, nil
		},
		func(ctx context.Context, focus Collection[I, A, ERR], source S) (T, error) {

			ctx = context.WithValue(ctx, indexTunnelKey, &indexTunnel{})
			mappedNext, _ := iter.Pull(iter.Seq[ValueIE[I, A]](focus.AsIter()(ctx)))
			ret, err := op.AsModify()(ctx, func(index I, focus ValueI[I, A]) (ValueI[I, A], error) {
				mapVal, ok := mappedNext()
				if !ok {
					return ValI(focus.index, focus.value), nil
				} else {
					return ValI(mapVal.index, mapVal.value), mapVal.err
				}
			}, source)

			return ret, err
		},
		IxMatchVoid(),
		ExprDef(
			func(t expr.OpticTypeExpr) expr.OpticExpression {
				return expr.CollectionOf{
					OpticTypeExpr: t,
					I:             reflect.TypeFor[I](),
					A:             reflect.TypeFor[A](),
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// The ColOfP combinator focuses on a polymorphic [Collection] of all the elements in the given optic.
//
// Under modification this collection can be modified and will be rebuilt into the original data structure.
// If the modified collection contains fewer elements then [ErrUnsafeMissingElement] will be returned
// If the modified collection contains more elements they will be ignored.
//
// See:
//   - [ColOf] for a safe non polymorphic version.
//   - [ColyTypeOfP] for a version that operates on concrete collection types
func ColOfP[I, S, T, A, B, RETI, RW, DIR, ERR any](o Optic[I, S, T, A, B, RETI, RW, DIR, ERR]) Optic[Void, S, T, Collection[I, A, Err], Collection[I, B, Err], ReturnOne, RW, UniDir, Err] {

	//WithIndex,TraverseColIEP and the Combinator below work together to support modifying indexes during a modify()
	op := WithIndex(o)
	return CombiLens[RW, Err, Void, S, T, Collection[I, A, Err], Collection[I, B, Err]](
		func(ctx context.Context, source S) (Void, Collection[I, A, Err], error) {
			ret := ColIE[Err, I, A](
				func(ctx context.Context) SeqIE[I, A] {
					return func(yield func(ValueIE[I, A]) bool) {
						op.AsIter()(ctx, source)(func(val ValueIE[I, ValueI[I, A]]) bool {
							return yield(ValIE(val.value.index, val.value.value, val.err))
						})
					}
				},
				func(ctx context.Context, index I) SeqIE[I, A] {
					return func(yield func(ValueIE[I, A]) bool) {
						op.AsIxGetter()(ctx, index, source)(func(val ValueIE[I, ValueI[I, A]]) bool {
							return yield(ValIE(val.value.index, val.value.value, val.err))
						})
					}
				},
				o.AsIxMatch(),
				func(ctx context.Context) (int, error) {
					return op.AsLengthGetter()(ctx, source)
				},
			)
			return Void{}, ret, nil
		},
		func(ctx context.Context, focus Collection[I, B, Err], source S) (T, error) {

			ctx = context.WithValue(ctx, indexTunnelKey, &indexTunnel{})
			mappedNext, _ := iter.Pull(iter.Seq[ValueIE[I, B]](focus.AsIter()(ctx)))
			ret, err := op.AsModify()(ctx, func(index I, focus ValueI[I, A]) (ValueI[I, B], error) {
				mapVal, ok := mappedNext()
				if !ok {
					return ValueI[I, B]{}, ErrUnsafeMissingElement
				} else {
					return ValI(mapVal.index, mapVal.value), mapVal.err
				}
			}, source)

			return ret, err
		},
		IxMatchVoid(),
		ExprDef(
			func(t expr.OpticTypeExpr) expr.OpticExpression {
				return expr.CollectionOf{
					OpticTypeExpr: t,
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
