package optic

import (
	"cmp"
	"context"
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/samber/lo"
	"github.com/spearson78/go-optic/expr"
)

// BinaryOp returns an [Operation] version of the given operator function.
// This constructor is a convenience wrapper for [Operator] with a simpler expression parameter.
//
// The following additional constructors are available:
//   - [BinaryOpE] impure
func BinaryOp[L, R, V any](op func(left L, right R) V, exprOp string) Optic[Void, lo.Tuple2[L, R], lo.Tuple2[L, R], V, V, ReturnOne, ReadOnly, UniDir, Pure] {
	return CombiGetter[Pure, Void, lo.Tuple2[L, R], lo.Tuple2[L, R], V, V](
		func(ctx context.Context, vals lo.Tuple2[L, R]) (Void, V, error) {
			return Void{}, op(vals.A, vals.B), ctx.Err()
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.BinaryExpr{
				OpticTypeExpr: ot,
				Op:            exprOp,
				L:             reflect.TypeFor[L](),
				R:             reflect.TypeFor[R](),
			}
		}),
	)
}

// BinaryOpE returns an [Operation] version of the given operator function.
// This constructor is a convenience wrapper for [Operator] with a simpler expression parameter.
//
// The following additional constructors are available:
//   - [BinaryOp] pure
func BinaryOpE[L, R, V any](op func(ctx context.Context, left L, right R) (V, error), exprOp string) Optic[Void, lo.Tuple2[L, R], lo.Tuple2[L, R], V, V, ReturnOne, ReadOnly, UniDir, Err] {
	return CombiGetter[Err, Void, lo.Tuple2[L, R], lo.Tuple2[L, R], V, V](
		func(ctx context.Context, vals lo.Tuple2[L, R]) (Void, V, error) {
			r, e := op(ctx, vals.A, vals.B)
			return Void{}, r, JoinCtxErr(ctx, e)
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.BinaryExpr{
				OpticTypeExpr: ot,
				Op:            exprOp,
				L:             reflect.TypeFor[L](),
				R:             reflect.TypeFor[R](),
			}
		}),
	)
}

// UnaryOp returns an [Operation] version of the given unary operator function.
// This constructor is a convenience wrapper for [Operator] with a simpler expression parameter.
//
// This should only be used for 1 way lossy operations.
// If the operation is reversible then use an [Iso] instead.
func UnaryOp[S, A any](op func(S) A, exprOp string) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return Operator[S, A](
		op,
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.UnaryExpr{
				OpticTypeExpr: ot,
				Op:            exprOp,
			}
		}),
	)
}

// MulT2 returns an [BinaryOp] that multiplies 2 values.
//
// See:
//   - [Mul] for a unary version.
//   - [MulOp] for a version that is applied to the focus of 2 [Operation]s.
//   - [Product] for a Reducer variant.
func MulT2[A Arithmetic]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right A) A { return left * right }, "*")
}

// Mul returns an [Iso] that multiplies by a constant value.
//
// See:
//   - [MulT2] for the [BinaryOp] version.
//   - [MulOp] for a version that is applied to the focus of 2 [Operation]s.
//   - [Product] for a Reducer variant.
func Mul[A Arithmetic](right A) Optic[Void, A, A, A, A, ReturnOne, ReadWrite, BiDir, Pure] {
	return EPure(OpT2IsoBind(MulT2[A](), DivT2[A](), right))
}

// MulOp returns an [Operation] expression that multiplies the focuses of 2 [Operation]s.
//
// See:
//   - [MulT2] for the [BinaryOp] version.
//   - [Mul] for the unary version.
func MulOp[S any, A Arithmetic, RETL, RETR TReturnOne, ERRL, ERRR any](left Operation[S, A, RETL, ERRL], right Operation[S, A, RETR, ERRR]) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, CompositionTree[ERRL, ERRR]] {
	return EErrL(OpT2ToOptic(left, MulT2[A](), right))
}

// Product returns a [Reducer] that multiplies all the focused elements.
//
// See [Mul] for the non Reducer version.
func Product[T Arithmetic]() Reduction[T, Pure] {
	return AsReducer(1, MulT2[T]())
}

// DivT2 returns an [BinaryOp] that divides 2 values.
//
// See:
//   - [Div] for a unary version.
//   - [DivOp] for a version that is applied to the focus of 2 [Operation]s.
func DivT2[A Arithmetic]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right A) A { return left / right }, "/")
}

// Div returns an [Iso] that divides by a constant value.
//
// See:
//   - [DivQuotient] for a version that uses a constant quotient.
//   - [DivOp] for a version that is applied to the focus of 2 [Operation]s.
func Div[A Arithmetic](divisor A) Optic[Void, A, A, A, A, ReturnOne, ReadWrite, BiDir, Pure] {
	return EPure(OpT2IsoBind(DivT2[A](), MulT2[A](), divisor))
}

// DivQuotient returns an [Iso] that divides the given quotient by the focused value..
//
// See:
//   - [Div] for a version that uses a constant divisor.
//   - [DivOp] for a version that is applied to the focus of 2 [Operation]s.
func DivQuotient[A Arithmetic](quotient A) Optic[Void, A, A, A, A, ReturnOne, ReadWrite, BiDir, Pure] {
	divQuotient := Ret1(Ro(Ud(EPure(Compose(SwappedT2[A, A](), DivT2[A]())))))
	return EPure(OpT2IsoBind(divQuotient, divQuotient, quotient))
}

// DivOp returns an [Operation] expression that multiplies the focuses of 2 [Operation]s.
//
// See:
//   - [Div] for the unary version.
//   - [DivT2] for the [BinaryOp] version.
func DivOp[S any, A Arithmetic, LRET, RRET TReturnOne, LERR, RERR any](left Operation[S, A, LRET, LERR], right Operation[S, A, RRET, RERR]) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, CompositionTree[LERR, RERR]] {
	return EErrL(OpT2ToOptic(
		left,
		DivT2[A](),
		right,
	))
}

// ModT2 returns an [BinaryOp] that applies A % B.
//
// See:
//   - [Mod] for a unary version.
//   - [ModOp] for a version that is applied to the focus of 2 [Operation]s.
func ModT2[A Integer]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right A) A { return left % right }, "%")
}

// Mod returns a [Iso] that performs a modulo by a constant value.
//
// See:
//   - [ModT2] for the [BinaryOp] version
//   - [ModOp] for a version that is applied to the focus of 2 [Operation]s.
func Mod[A Integer](right A) Optic[Void, A, A, A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2Bind(ModT2[A](), right)
}

// ModOp returns an [Operation] expression that applies a modulo to the focuses of 2 [Operation]s.
//
// See:
//   - [ModT2] for the [BinaryOp] version
//   - [Mod] for the constant version.
func ModOp[S any, A Integer, LRET, RRET TReturnOne, LERR, RERR any](left Operation[S, A, LRET, LERR], right Operation[S, A, RRET, RERR]) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, CompositionTree[LERR, RERR]] {
	return EErrL(OpT2ToOptic(
		left,
		ModT2[A](),
		right,
	))
}

// AddT2 returns an [BinaryOp] that adds 2 values.
//
// See:
//   - [Add] for a unary version.
//   - [AddOp] for a version that is applied to the focus of 2 [Operation]s.
func AddT2[A Arithmetic]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right A) A {
		return left + right
	}, "+")
}

// Add returns a [Iso] that performs an addition by a constant value.
//
// See:
//   - [AddOp] for a version that is applied to the focus of 2 [Operation]s.
//   - [Sum] for a Reducer variant.
func Add[A Arithmetic](right A) Optic[Void, A, A, A, A, ReturnOne, ReadWrite, BiDir, Pure] {
	return EPure(OpT2IsoBind(AddT2[A](), SubT2[A](), right))
}

// AddOp returns an [Operation] expression that adds the focuses of 2 [Operation]s.
//
// See:
//   - [Add] for the constant version.
func AddOp[S any, A Arithmetic, LRET, RRET TReturnOne, LERR, RERR any](left Operation[S, A, LRET, LERR], right Operation[S, A, RRET, RERR]) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, CompositionTree[LERR, RERR]] {
	return EErrL(OpT2ToOptic(
		left,
		AddT2[A](),
		right,
	))
}

// Sum returns a [Reducer] that sums all the focused elements.
//
// See [Add] for the non Reducer version.
func Sum[T Arithmetic]() Reduction[T, Pure] {
	return AsReducer(0, AddT2[T]())
}

// SubT2 returns an [BinaryOp] that subtracts 2 values.
//
// See:
//   - [Sub] for a unary version.
//   - [SubOp] for a version that is applied to the focus of 2 [Operation]s.
func SubT2[A Arithmetic]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right A) A { return left - right }, "-")
}

// Sub returns a [Iso] that performs a subtraction by a constant value.
//
// See:
//   - [SubFrom] for a version that subtracts from a constant value.
//   - [SubOp] for a version that is applied to the focus of 2 [Operation]s.
func Sub[A Arithmetic](right A) Optic[Void, A, A, A, A, ReturnOne, ReadWrite, BiDir, Pure] {
	return EPure(OpT2IsoBind(SubT2[A](), AddT2[A](), right))
}

// SubFrom returns a [Iso] that performs a subtraction from a constant value.
//
// See:
//   - [Sub] for a version that subtracts a constant value.
//   - [SubOp] for a version that is applied to the focus of 2 [Operation]s.
func SubFrom[A Arithmetic](left A) Optic[Void, A, A, A, A, ReturnOne, ReadWrite, BiDir, Pure] {
	subFrom := Ret1(Ro(Ud(EPure(Compose(SwappedT2[A, A](), SubT2[A]())))))
	inverse := Ret1(Ro(Ud(EPure(Compose(SubT2[A](), Negate[A]())))))
	return EPure(OpT2IsoBind(subFrom, inverse, left))
}

// SubOp returns an [Operation] expression that subtracts the focuses of 2 [Operation]s.
//
// See:
//   - [Add] for the constant version.
func SubOp[S any, A Arithmetic, LRET, RRET TReturnOne, LERR, RERR any](left Operation[S, A, LRET, LERR], right Operation[S, A, RRET, RERR]) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, CompositionTree[LERR, RERR]] {
	return EErrL(OpT2ToOptic(
		left,
		SubT2[A](),
		right,
	))
}

// PowT2 returns an [BinaryOp] that raises A to the power B.
//
// See:
//   - [Pow] for a unary version.
//   - [PowOp] for a version that is applied to the focus of 2 [Operation]s.
func PowT2() Optic[Void, lo.Tuple2[float64, float64], lo.Tuple2[float64, float64], float64, float64, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right float64) float64 { return math.Pow(left, right) }, "**")
}

// Pow returns a [Iso] that raises a value to a constant power.
//
// See:
//   - [PowOp] for a version that is applied to the focus of 2 [Operation]s.
func Pow(right float64) Optic[Void, float64, float64, float64, float64, ReturnOne, ReadWrite, BiDir, Pure] {
	return EPure(OpT2IsoBind(PowT2(), RootT2(), right))
}

// PowOp returns an [Operation] expression that raises the left focus to the power of the right focus.
//
// See:
//   - [Pow] for the constant version.
func PowOp[S any, LRET, RRET TReturnOne, LERR, RERR any](left Operation[S, float64, LRET, LERR], right Operation[S, float64, RRET, RERR]) Optic[Void, S, S, float64, float64, ReturnOne, ReadOnly, UniDir, CompositionTree[LERR, RERR]] {
	return EErrL(OpT2ToOptic(
		left,
		PowT2(),
		right,
	))
}

// RootT2 returns an [BinaryOp] that returns root B of A.
//
// See:
//   - [Root] for a unary version.
//   - [RootOp] for a version that is applied to the focus of 2 [Operation]s.
func RootT2() Optic[Void, lo.Tuple2[float64, float64], lo.Tuple2[float64, float64], float64, float64, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right float64) float64 { return math.Pow(left, 1.0/right) }, "root")
}

// Root returns a [Iso] that returns a root of a value.
//
// See:
//   - [RootOp] for a version that is applied to the focus of 2 [Operation]s.
func Root(right float64) Optic[Void, float64, float64, float64, float64, ReturnOne, ReadWrite, BiDir, Pure] {
	return EPure(OpT2IsoBind(RootT2(), PowT2(), right))
}

// RootOp returns an [Operation] expression that returns the right root of the left focus.
//
// See:
//   - [Root] for the constant version.
func RootOp[S any, LRET, RRET TReturnOne, LERR, RERR any](left Operation[S, float64, LRET, LERR], right Operation[S, float64, RRET, RERR]) Optic[Void, S, S, float64, float64, ReturnOne, ReadOnly, UniDir, CompositionTree[LERR, RERR]] {
	return EErrL(OpT2ToOptic(
		left,
		RootT2(),
		right,
	))
}

// Negate returns a [Iso] that negates the value.
func Negate[A Arithmetic]() Optic[Void, A, A, A, A, ReturnOne, ReadWrite, BiDir, Pure] {
	return Involuted(
		func(x A) A { return -x },
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.UnaryExpr{
				OpticTypeExpr: ot,
				Op:            "negate",
			}
		}),
	)
}

// Abs returns a [Operation] that focuses on the absolute value.
func Abs[A Real]() Optic[Void, A, A, A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return Operator(
		func(x A) A {
			if x < 0 {
				return -x
			} else {
				return x
			}
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.UnaryExpr{
				OpticTypeExpr: ot,
				Op:            "abs",
			}
		}),
	)
}

// ParseInt returns an [Iso] that parses the value to an int.
// An error is returned if the string cannot be parsed
//
// The base and bitSize parameters correspond to the [strconv.ParseInt] parameters with the same name.
// See:
//   - [ParseIntP] for a polymorphic version.
//   - [FormatInt] for a pure readonly formating version
func ParseInt[T Real](base int, bitSize int) Optic[Void, string, string, T, T, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoEP[string, string, T, T](
		func(ctx context.Context, s string) (T, error) {
			i, err := strconv.ParseInt(s, base, bitSize)
			return T(i), err
		},
		func(ctx context.Context, f T) (string, error) {
			return strconv.FormatInt(int64(f), base), nil
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.ParseInt{
				OpticTypeExpr: ot,
				Base:          base,
				BitSize:       bitSize,
			}
		}),
	)
}

// FormatInt returns an [Operation] that formats an int to a string.
//
// The base parameters correspond to the [strconv.FormatInt] parameter with the same name.
// See:
//   - [ParseInt] for an [Iso] version.
func FormatInt[T Real](base int) Optic[Void, T, T, string, string, ReturnOne, ReadOnly, UniDir, Pure] {
	return CombiGetter[Pure, Void, T, T, string, string](
		func(ctx context.Context, source T) (Void, string, error) {
			return Void{}, strconv.FormatInt(int64(source), base), nil
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.FormatInt{
				OpticTypeExpr: ot,
				Base:          base,
			}
		}),
	)
}

// ParseIntP returns a polymorphic [Iso] that parses the value to an int.
// An error is returned if the string cannot be parsed
//
// The base and bitSize parameters correspond to the [strconv.ParseInt] parameters with the same name.
// See:
//   - [ParseInt] for a non polymorphic version.
func ParseIntP[T Real](base int, bitSize int) Optic[Void, string, T, T, T, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoEP[string, T, T, T](
		func(ctx context.Context, s string) (T, error) {
			i, err := strconv.ParseInt(s, base, bitSize)
			return T(i), JoinCtxErr(ctx, err)
		},
		func(ctx context.Context, f T) (T, error) {
			return f, ctx.Err()
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.ParseInt{
				OpticTypeExpr: ot,
				Base:          base,
				BitSize:       bitSize,
			}
		}),
	)
}

// ParseFloat returns a [Iso] that parses the value to an float.
//
// The fmt and prec parameters corresponds to the [strconv.FormatFloat] parameters with the same name.
// The bitSize parameter corresponds to the [strconv.ParseFloat] parameter with the same name.
//
// See:
//   - [FormatFloat] for a pure readonly formating version
func ParseFloat[T Real](fmt byte, prec int, bitSize int) Optic[Void, string, string, T, T, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoEP[string, string, T, T](
		func(ctx context.Context, s string) (T, error) {
			i, err := strconv.ParseFloat(s, bitSize)
			return T(i), err
		},
		func(ctx context.Context, f T) (string, error) {
			return strconv.FormatFloat(float64(f), fmt, prec, bitSize), nil
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.ParseFloat{
				OpticTypeExpr: ot,
				Fmt:           fmt,
				Prec:          prec,
				BitSize:       bitSize,
			}
		}),
	)
}

// ParseFloatP returns a polymorphic [Iso] that parses the value to an float.
//
// The base and bitSize parameters correspond to the [strconv.ParseFloat] parameters with the same name.
func ParseFloatP[T Real](bitSize int) Optic[Void, string, T, T, T, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoEP[string, T, T, T](
		func(ctx context.Context, source string) (T, error) {
			i, err := strconv.ParseFloat(source, bitSize)
			return T(i), JoinCtxErr(ctx, err)
		},
		func(ctx context.Context, f T) (T, error) {
			return f, ctx.Err()
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.ParseFloatP{
				OpticTypeExpr: ot,
				BitSize:       bitSize,
			}
		}),
	)
}

// FormatFloat returns an [Operation] that formats a float to a string.
//
// The fmt,prec and bitSize parameters correspond to the [strconv.FormatFloat] parameters with the same name.
// See:
//   - [ParseFloat] for an [Iso] version.
func FormatFloat[T Real](fmt byte, prec int, bitSize int) Optic[Void, T, T, string, string, ReturnOne, ReadOnly, UniDir, Pure] {
	return CombiGetter[Pure, Void, T, T, string, string](
		func(ctx context.Context, source T) (Void, string, error) {
			return Void{}, strconv.FormatFloat(float64(source), fmt, prec, bitSize), nil
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.FormatFloat{
				OpticTypeExpr: ot,
				Fmt:           fmt,
				Prec:          prec,
				BitSize:       bitSize,
			}
		}),
	)
}

// Mean returns a [Reducer] that calculates the mean value.
func Mean[T Real]() ReductionP[lo.Tuple2[T, int], T, T, Pure] {

	reshape := Ret1(Rw(Ud(EPure(T2Of(
		T2Of(
			Compose(
				T2A[lo.Tuple2[T, int], T](),
				T2A[T, int](),
			),
			T2B[lo.Tuple2[T, int], T](),
		),
		Compose(
			T2A[lo.Tuple2[T, int], T](),
			T2B[T, int](),
		),
	)))))

	fmap := T2Of(
		Compose(
			T2A[lo.Tuple2[T, T], int](),
			T2Of(
				AddT2[T](),
				T2B[T, T](),
			),
		),
		Compose(
			T2B[lo.Tuple2[T, T], int](),
			Add(1),
		),
	)

	reduce := EPure(Compose(
		AsModify(reshape, fmap),
		T2A[lo.Tuple2[T, int], T](),
	))

	end := Compose(
		T2Of(
			T2A[T, int](),
			Compose(
				T2B[T, int](),
				IsoCast[int, T](),
			),
		),
		DivT2[T](),
	)

	var t T
	return AsReducerP(
		Const[Void](lo.T2(t, 0)),
		reduce,
		end,
	)
}

type medianState[T cmp.Ordered] struct {
	maxHeap *heap[T]
	minHeap *heap[T]
}

// Median returns a [Reducer] that calculates the median value.
func Median[T Real]() ReductionP[medianState[T], T, T, Pure] {
	return CombiReducer[Pure, medianState[T], T, T](
		func(ctx context.Context) (medianState[T], error) {

			maxHeap, _, _ := newHeap[T](context.Background(), nil, func(ctx context.Context, a T, b T) (bool, error) {
				return a < b, nil
			})

			minHeap, _, _ := newHeap[T](context.Background(), nil, func(ctx context.Context, a T, b T) (bool, error) {
				return a > b, nil
			})

			return medianState[T]{
				minHeap: minHeap,
				maxHeap: maxHeap,
			}, nil
		},
		func(ctx context.Context, state medianState[T], appendVal T) (medianState[T], error) {
			state.maxHeap.push(ctx, appendVal)

			temp, cmpCounter1, err := state.maxHeap.pop(ctx)
			if err != nil {
				IncCustomMetric(ctx, "comparisons", cmpCounter1)
				return state, err
			}

			cmpCounter2, err := state.minHeap.push(ctx, temp)
			if err != nil {
				IncCustomMetric(ctx, "comparisons", cmpCounter1+cmpCounter2)
				return state, err
			}

			if len(state.minHeap.data) > len(state.maxHeap.data) {
				temp, cmpCounter3, err := state.minHeap.pop(ctx)
				if err != nil {
					IncCustomMetric(ctx, "comparisons", cmpCounter1+cmpCounter2+cmpCounter3)
					return state, err
				}

				cmpCounter4, err := state.maxHeap.push(ctx, temp)
				IncCustomMetric(ctx, "comparisons", cmpCounter1+cmpCounter2+cmpCounter3+cmpCounter4)
				return state, err

			} else {
				IncCustomMetric(ctx, "comparisons", cmpCounter1+cmpCounter2)
				return state, nil
			}
		},
		func(ctx context.Context, state medianState[T]) (T, error) {
			v := state.maxHeap.peek()
			if len(state.maxHeap.data) != len(state.minHeap.data) {
				return v, nil
			}
			v2 := state.minHeap.peek()

			return (v + v2) / T(2), nil
		},
		ReducerExprDef(
			func(t expr.ReducerTypeExpr) expr.ReducerExpression {
				return expr.MedianReducer{
					ReducerTypeExpr: t,
				}
			},
		),
	)
}

// Mode returns a [Reducer] that calculates the modal value.
func Mode[T cmp.Ordered]() ReductionP[map[T]int, T, T, Pure] {

	optic := Compose(
		AtMapT2[T, int](),
		Non[int](0, EqT2[int]()),
	)

	fmap := Add(1)

	reduce := EPure(Compose(
		AsModify(optic, fmap),
		T2A[map[T]int, T](),
	))

	var t T
	end := Compose(
		WithIndex(
			FirstOrDefaultI(
				MaxOf(
					TraverseMap[T, int](),
					Identity[int](),
				),
				t,
				0,
			),
		),
		ValueIIndex[T, int](),
	)

	return AsReducerP(
		MakeMap[Void, T, int](10),
		reduce,
		end,
	)
}

func minVal[T cmp.Ordered]() (r T) {
	switch x := any(&r).(type) {
	case *int:
		*x = math.MinInt
	case *int8:
		*x = math.MinInt8
	case *int16:
		*x = math.MinInt16
	case *int32:
		*x = math.MinInt32
	case *int64:
		*x = math.MinInt64
	case *uint:
		*x = 0
	case *uint8:
		*x = 0
	case *uint16:
		*x = 0
	case *uint32:
		*x = 0
	case *uint64:
		*x = 0
	case *float32:
		*x = -math.MaxFloat32
	case *float64:
		*x = -math.MaxFloat64
	case *string:
		*x = ""
	default:
		//Use reflection for alias types
		rptr := reflect.ValueOf(&r)

		switch rptr.Elem().Kind() {
		case reflect.Int:
			rptr.Elem().SetInt(math.MinInt)
		case reflect.Int8:
			rptr.Elem().SetInt(math.MinInt8)
		case reflect.Int16:
			rptr.Elem().SetInt(math.MinInt16)
		case reflect.Int32:
			rptr.Elem().SetInt(math.MinInt32)
		case reflect.Int64:
			rptr.Elem().SetInt(math.MinInt64)
		case reflect.Uint:
			rptr.Elem().SetUint(0)
		case reflect.Uint8:
			rptr.Elem().SetUint(0)
		case reflect.Uint16:
			rptr.Elem().SetUint(0)
		case reflect.Uint32:
			rptr.Elem().SetUint(0)
		case reflect.Uint64:
			rptr.Elem().SetUint(0)
		case reflect.Uintptr:
			rptr.Elem().SetUint(0)
		case reflect.Float32:
			rptr.Elem().SetFloat(-math.MaxFloat32)
		case reflect.Float64:
			rptr.Elem().SetFloat(-math.MaxFloat64)
		case reflect.String:
			rptr.Elem().SetString("")
		default:
			panic(fmt.Errorf("minVal unknown type : %T", r))
		}
	}
	return
}

func maxVal[T Real]() (r T) {

	switch x := any(&r).(type) {
	case *int:
		*x = math.MaxInt
	case *int8:
		*x = math.MaxInt8
	case *int16:
		*x = math.MaxInt16
	case *int32:
		*x = math.MaxInt32
	case *int64:
		*x = math.MaxInt64
	case *uint:
		*x = math.MaxUint
	case *uint8:
		*x = math.MaxUint8
	case *uint16:
		*x = math.MaxUint16
	case *uint32:
		*x = math.MaxUint32
	case *uint64:
		*x = math.MaxUint64
	case *float32:
		*x = math.MaxFloat32
	case *float64:
		*x = math.MaxFloat64
	default:
		//Use reflection for alias types
		rptr := reflect.ValueOf(&r)

		switch rptr.Elem().Kind() {
		case reflect.Int:
			rptr.Elem().SetInt(math.MaxInt)
		case reflect.Int8:
			rptr.Elem().SetInt(math.MaxInt8)
		case reflect.Int16:
			rptr.Elem().SetInt(math.MaxInt16)
		case reflect.Int32:
			rptr.Elem().SetInt(math.MaxInt32)
		case reflect.Int64:
			rptr.Elem().SetInt(math.MaxInt64)
		case reflect.Uint:
			rptr.Elem().SetUint(math.MaxUint)
		case reflect.Uint8:
			rptr.Elem().SetUint(math.MaxUint8)
		case reflect.Uint16:
			rptr.Elem().SetUint(math.MaxUint16)
		case reflect.Uint32:
			rptr.Elem().SetUint(math.MaxUint32)
		case reflect.Uint64:
			rptr.Elem().SetUint(math.MaxUint64)
		case reflect.Uintptr:
			rptr.Elem().SetUint(math.MaxUint64)
		case reflect.Float32:
			rptr.Elem().SetFloat(math.MaxFloat32)
		case reflect.Float64:
			rptr.Elem().SetFloat(math.MaxFloat64)
		default:
			panic(fmt.Errorf("maxVal unknown type : %T", r))
		}
	}
	return
}

// MinT2 returns an [BinaryOp] that returns the minimum.
//
// See:
//   - [Min] for a unary version.
//   - [MinOp] for a version that is applied to the focus of 2 [Operation]s.
func MinT2[A cmp.Ordered]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right A) A { return min(left, right) }, "min")
}

// Min returns an [Operation] that returns the minimum value
//
// See:
//   - [MinOp] for a version that is applied to the focus of 2 [Operation]s.
//   - [MinReducer] for a Reducer variant.
//   - [MinI] for a version that returns also return the index.
func Min[A cmp.Ordered](right A) Optic[Void, A, A, A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2Bind(MinT2[A](), right)
}

// MinOp returns an [Operation] expression that returns the minimum of the 2 [Operation]s.
//
// See:
//   - [Min] for the constant version.
func MinOp[S any, A cmp.Ordered, LRET, RRET TReturnOne, LERR, RERR any](left Operation[S, A, LRET, LERR], right Operation[S, A, RRET, RERR]) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, CompositionTree[LERR, RERR]] {
	return EErrL(OpT2ToOptic(
		left,
		MinT2[A](),
		right,
	))
}

// MinReducer returns a [Reducer] that returns the minimum value.
//
// See:
//   - [Min] for the non Reducer version.
//   - [MinReducerI] for a version that also returns the index.
func MinReducer[T Real]() Reduction[T, Pure] {
	return AsReducer(maxVal[T](), MinT2[T]())
}

// MaxT2 returns an [BinaryOp] that returns the maximum.
//
// See:
//   - [Max] for a unary version.
//   - [MaxOp] for a version that is applied to the focus of 2 [Operation]s.
func MaxT2[A cmp.Ordered]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right A) A { return max(left, right) }, "max")
}

// Max returns an [Operation] that returns the maximum value
//
// See:
//   - [MaxOp] for a version that is applied to the focus of 2 [Operation]s.
//   - [MaxReducer] for a Reducer variant.
//   - [MaxI] for a version that also returns the index.
func Max[A cmp.Ordered](right A) Optic[Void, A, A, A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2Bind(MaxT2[A](), right)
}

// MaxOp returns an [Operation] expression that returns the maximum of the values focused by 2 [Operation]s.
//
// See:
//   - [Max] for the constant version.
//   - [MaxReducer] for a Reducer variant.
func MaxOp[S any, A cmp.Ordered, LRET, RRET TReturnOne, LERR, RERR any](left Operation[S, A, LRET, LERR], right Operation[S, A, RRET, RERR]) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, CompositionTree[LERR, RERR]] {
	return EErrL(OpT2ToOptic(
		left,
		MaxT2[A](),
		right,
	))
}

// MaxReducer returns a [Reducer] that returns the maximum of all focused elements.
//
// See
//   - [Max] for the non Reducer version.
//   - [MaxReducerI] for a version that also returns the index.
func MaxReducer[T cmp.Ordered]() Reduction[T, Pure] {
	return AsReducer(minVal[T](), MaxT2[T]())
}

// Clamp returns an [Operation] that constrains the focused value between min and max.
func Clamp[A cmp.Ordered](min A, max A) Optic[Void, A, A, A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return EPure(Ret1(Ro(Ud(Compose(Max(min), Min(max))))))
}
