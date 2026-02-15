package optic

import (
	"cmp"
	"context"
	"errors"
	"iter"

	"github.com/samber/lo"
	"github.com/spearson78/go-optic/expr"
)

// A Predicate is an [Operation] that return true is the predicate is satisfied.
type Predicate[S any, ERR any] interface {
	AsOpGet() OpGetFunc[S, bool]
	ErrType() ERR
	AsExprHandler() func(ctx context.Context) (ExprHandler, error)
	AsExpr() expr.OpticExpression
}

// A PredicateE is a variant of [Predicate]  where it should be assumed an error is returned.
// This interface can be used to implement fluent style APIs at the cost of no longer being pure.
type PredicateE[S any] interface {
	AsOpGet() OpGetFunc[S, bool]
	AsExprHandler() func(ctx context.Context) (ExprHandler, error)
	AsExpr() expr.OpticExpression
}

// [PredGet] gets the predicate value from an [Operation] or [OperationI] defaulting to false if no value is found.
func PredGet[S any](ctx context.Context, fnc PredicateE[S], source S) (bool, error) {
	ret, err := fnc.AsOpGet()(ctx, source)
	if errors.Is(err, ErrEmptyGet) {
		return false, nil
	}
	return ret, err
}

// True returns a [Predicate] that is always satisfied.
func True[A any]() Optic[Void, A, A, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return Ro(Const[A](true))
}

// False returns a [Predicate] that is never satisfied.
func False[A any]() Optic[Void, A, A, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return Ro(Const[A](false))
}

// Even returns a predicate is satisfied if the given arithmetic value is even.
func Even[A Integer]() Optic[Void, A, A, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return EPure(Ret1(Ro(Ud(Compose(Mod(A(2)), Eq(A(0)))))))
}

// Odd returns a predicate is satisfied if the given arithmetic value is odd.
func Odd[A Integer]() Optic[Void, A, A, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return EPure(Ret1(Ro(Ud(Compose(Mod(A(2)), Eq(A(1)))))))
}

// GtT2 returns an [BinaryOp] that is satisfied if A > B.
//
// See:
//   - [Gt] for a unary version.
//   - [GtOp] for a version that is applied to the focus of 2 [Operation]s.
func GtT2[A cmp.Ordered]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right A) bool { return left > right }, ">")
}

// Gt returns a [Predicate] that is satisfied if the focused value is greater than (>) the provided constant value.
//
// See:
//   - [GtOp] for a [Predicate] that compares 2 focuses instead of a constant.
//   - [GtI] for a version that compares indices.
func Gt[A cmp.Ordered](right A) Optic[Void, A, A, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2Bind(GtT2[A](), right)
}

// GtOp returns a [Predicate] that is satisfied if the left focus is greater than (>) the right focus.
//
// See [Gt] for a predicate that checks constant a value rather than 2 focuses.
func GtOp[S any, A cmp.Ordered, RETL, RETR TReturnOne, LERR, RERR any](left Operation[S, A, RETL, LERR], right Operation[S, A, RETR, RERR]) Optic[Void, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, CompositionTree[LERR, RERR]] {
	return EErrL(OpT2ToOptic(
		left,
		GtT2[A](),
		right,
	))
}

// GteT2 returns an [BinaryOp] that is satisfied if A >= B.
//
// See:
//   - [Gte] for a unary version.
//   - [GteOp] for a version that is applied to the focus of 2 [Operation]s.
func GteT2[A cmp.Ordered]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right A) bool { return left >= right }, ">=")
}

// Gte returns a [Predicate] that is satisfied if the focused value is greater than or equal to (>=) the provided constant value.
//
// See [GteOp] for a [Predicate] that compares 2 focuses instead of a constant.
func Gte[A cmp.Ordered](right A) Optic[Void, A, A, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2Bind(GteT2[A](), right)
}

// GteOp returns a [Predicate] that is satisfied if the left focus is greater than or equals to (>=) the right focus.
//
// See:
//   - [Gt] for a predicate that checks constant a value rather than 2 focuses.
//   - [GteI] for a version that compares indices.
func GteOp[S any, A cmp.Ordered, RETL, RETR TReturnOne, LERR, RERR any](left Operation[S, A, RETL, LERR], right Operation[S, A, RETR, RERR]) Optic[Void, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, CompositionTree[LERR, RERR]] {
	return EErrL(OpT2ToOptic(
		left,
		GteT2[A](),
		right,
	))
}

// LtT2 returns an [BinaryOp] that is satisfied if A < B.
//
// See:
//   - [Lt] for a unary version.
//   - [LtOp] for a version that is applied to the focus of 2 [Operation]s.
func LtT2[A cmp.Ordered]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right A) bool {
		return left < right
	}, "<")
}

// Lt returns a [Predicate] that is satisfied if the focused value is less than (<) the provided constant value.
//
// See:
//   - [LtOp] for a [Predicate] that compares 2 focuses instead of a constant.
//   - [LtI] for a version that compares indices.
func Lt[A cmp.Ordered](right A) Optic[Void, A, A, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2Bind(LtT2[A](), right)
}

// LtOp returns a [Predicate] that is satisfied if the left focus is less (<) the right focus.
//
// See:
//   - [Lt] for a predicate that checks constant a value rather than 2 focuses.
func LtOp[S any, A cmp.Ordered, RETL, RETR TReturnOne, LERR, RERR any](left Operation[S, A, RETL, LERR], right Operation[S, A, RETR, RERR]) Optic[Void, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, CompositionTree[LERR, RERR]] {
	return EErrL(OpT2ToOptic(
		left,
		LtT2[A](),
		right,
	))
}

// LteT2 returns an [BinaryOp] that is satisfied if A <= B.
//
// See:
//   - [Lte] for a unary version.
//   - [LteOp] for a version that is applied to the focus of 2 [Operation]s.
func LteT2[A cmp.Ordered]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right A) bool { return left <= right }, "<=")
}

// Lte returns a [Predicate] that is satisfied if the focused value is less than or equal to (<=) the provided constant value.
//
// See:
//   - [LteOp] for a [Predicate] that compares 2 focuses instead of a constant.
//   - [LteI] for a version that compares indices.
func Lte[A cmp.Ordered](right A) Optic[Void, A, A, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2Bind(LteT2[A](), right)
}

// LteOp returns a [Predicate] that is satisfied if the left focus is less than or equal to (<=) the right focus.
//
// See [Lte] for a predicate that checks constant a value rather than 2 focuses.
func LteOp[S any, A cmp.Ordered, RETL, RETR TReturnOne, LERR, RERR any](left Operation[S, A, RETL, LERR], right Operation[S, A, RETR, RERR]) Optic[Void, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, CompositionTree[LERR, RERR]] {
	return EErrL(OpT2ToOptic(
		left,
		LteT2[A](),
		right,
	))
}

// EqT2 returns an [BinaryOp] that is satisfied if A == B.
//
// See:
//   - [Eq] for a unary version.
//   - [EqOp] for a version that is applied to the focus of 2 [Operation]s.
func EqT2[A comparable]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right A) bool {
		return left == right
	}, "==")
}

// The EqT2Of combinator returns an optic that compares that each focused element and it's index are equal and that both traversals focus the same number of elements.
func EqT2Of[I, S, A, RET, RW, DIR, ERR any, PERR TPure](o Optic[I, S, S, A, A, RET, RW, DIR, ERR], eq Predicate[lo.Tuple2[A, A], PERR]) Optic[Void, lo.Tuple2[S, S], lo.Tuple2[S, S], bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {

	return CombiGetter[ERR, Void, lo.Tuple2[S, S], lo.Tuple2[S, S], bool, bool](
		func(ctx context.Context, source lo.Tuple2[S, S]) (Void, bool, error) {
			a := source.A
			b := source.B

			otherNext, otherStop := iter.Pull(iter.Seq[ValueIE[I, A]](o.AsIter()(ctx, b)))
			defer otherStop()
			for focus := range o.AsIter()(ctx, a) {
				if focus.err != nil {
					return Void{}, false, focus.err
				}
				otherFocus, otherOk := otherNext()
				if !otherOk {
					//a is longer than other
					return Void{}, false, nil
				}
				if otherFocus.err != nil {
					return Void{}, false, otherFocus.err
				}

				if !o.AsIxMatch()(focus.index, otherFocus.index) {
					return Void{}, false, nil
				}

				res, err := PredGet(ctx, eq, lo.T2(focus.value, otherFocus.value))
				if err != nil {
					return Void{}, false, err
				}

				if !res {
					return Void{}, false, nil
				}
			}

			_, otherOk := otherNext()
			if otherOk {
				//b is longer than a
				return Void{}, false, nil
			}

			return Void{}, true, nil
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.EqT2Of{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
					Eq:            eq.AsExpr(),
				}
			},
			o,
			eq,
		),
	)
}

// Eq returns a [Predicate] that is satisfied if the focused value is equal to (==) the provided constant value.
//
// See:
//   - [EqOp] for a [Predicate] that compares 2 focuses instead of a constant.
//   - [EqI] for a version that compares indices.
func Eq[A comparable](right A) Optic[Void, A, A, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2Bind(EqT2[A](), right)
}

// EqOp returns a [Predicate] that is satisfied if the left focus is equal to (==) the right focus.
//
// See [Eq] for a predicate that checks constant a value rather than 2 focuses.
func EqOp[S any, A comparable, RETL, RETR TReturnOne, LERR, RERR any](left Operation[S, A, RETL, LERR], right Operation[S, A, RETR, RERR]) Optic[Void, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, CompositionTree[LERR, RERR]] {
	return EErrL(OpT2ToOptic(
		left,
		EqT2[A](),
		right,
	))
}

func PredT2ToOptic[I, S, R, ORW any, ORET, ODIR any, OERR, LERR, RERR any](left Predicate[S, LERR], op Optic[I, lo.Tuple2[bool, bool], lo.Tuple2[bool, bool], R, R, ORET, ORW, ODIR, OERR], right Predicate[S, RERR]) Optic[I, S, S, R, R, ORET, ReadOnly, UniDir, CompositionTree[CompositionTree[LERR, RERR], OERR]] {
	return RetR(Ro(Ud(Compose(
		T2Of(
			PredToOptic(left),
			PredToOptic(right),
		),
		op,
	))))
}

func InT2[A comparable]() Optic[Void, lo.Tuple2[A, []A], lo.Tuple2[A, []A], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left A, right []A) bool {
		for _, v := range right {
			if left == v {
				return true
			}
		}
		return false
	}, "in")
}

// In returns a [Predicate] that is satisfied if the source value is present in a list of values.
//
// See:
//   - [Contains] for a combinator version
func In[A comparable](right ...A) Optic[Void, A, A, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2Bind(InT2[A](), right)
}

// NeT2 returns an [BinaryOp] that is satisfied if A != B.
//
// See:
//   - [NE] for a unary version.
//   - [NeOp] for a version that is applied to the focus of 2 [Operation]s.
func NeT2[A comparable]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right A) bool { return left != right }, "!=")
}

// Ne returns a [Predicate] that is satisfied if the focused value is not equal to (!=) the provided value.
//
// See:
//   - [NeOp] for a [Predicate] that compares 2 focuses instead of a constant.
//   - [NeI] for a version that compares indices.
func Ne[A comparable](right A) Optic[Void, A, A, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2Bind(NeT2[A](), right)
}

// NeOp returns a [Predicate] that is satisfied if the left focus is not equal to (!=) the right focus.
//
// See [Ne] for a predicate that checks constant a value rather than 2 focuses.
func NeOp[S any, A comparable, RETL, RETR TReturnOne, LERR, RERR any](left Operation[S, A, RETL, LERR], right Operation[S, A, RETR, RERR]) Optic[Void, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, CompositionTree[LERR, RERR]] {
	return EErrL(OpT2ToOptic(
		left,
		NeT2[A](),
		right,
	))
}

// Any returns a [Predicate] that is satisfied if any focus satisfies the passed predicate.
//
// Note: If the optic focuses no elements then Any returns false.
//
// See
//   - [AnyI] for an index aware version.
//   - [All] for a version that is satisified if all focuses satisfy the predicate.
func Any[I, S, T, A, B, RET, RW, DIR, ERR, PERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], pred Predicate[A, PERR]) Optic[I, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, CompositionTree[ERR, PERR]] {
	return AnyI(o, PredToOpI[I](pred))
}

// Contains returns a [Predicate] that is satisfied if any focus satisfies the passed predicate.
//
// Note: If the optic focuses no elements then Contains returns false.
//
// See:
// - [In] for a simple predicate version.
// - [All] for a version that requires all focuses to satisfy he predicate
func Contains[I, S, T any, A comparable, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], val ...A) Optic[I, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return EErrL(Any(o, In(val...)))
}

// All returns a [Predicate] that is satisfied if all focuses satisfy the passed predicate
//
// Note: If the optic focuses no elements then All returns true.
//
// See:
//   - [AllI] for an index aware version.
//   - [Any] for a version that is satisfied if any focus satisfies the predicate.
func All[I, S, T, A, B, RET, RW, DIR, ERR, PERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], pred Predicate[A, PERR]) Optic[Void, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, CompositionTree[ERR, PERR]] {
	return AllI(o, PredToOpI[I](pred))
}

// NotEmpty returns a [Predicate] that is satisfied if the the given optic focuses any elements.
//
// See [Empty] for the inverse version.
func NotEmpty[I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[I, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return EErrL(Any(o, True[A]()))
}

// Empty returns a [Predicate] that is satisfied if the the given optic focuses no elements.
//
// See [NotEmpty] for the inverse version.
func Empty[I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[Void, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return NotOp(NotEmpty(o))
}

func PredToOptic[S, ERR any](pred Predicate[S, ERR]) Optic[Void, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	predFnc := pred.AsOpGet()

	return rawGetterF[Void, S, bool, ERR](
		func(ctx context.Context, source S) (Void, bool, error) {
			v, err := predFnc(ctx, source)
			if errors.Is(err, ErrEmptyGet) {
				return Void{}, false, nil
			}
			return Void{}, v, err
		},
		IxMatchVoid(),
		pred.AsExprHandler(),
		pred.AsExpr,
	)
}

func PredEToOptic[S any](pred PredicateE[S]) Optic[Void, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, Err] {
	predFnc := pred.AsOpGet()

	return rawGetterF[Void, S, bool, Err](
		func(ctx context.Context, source S) (Void, bool, error) {
			v, err := predFnc(ctx, source)
			if errors.Is(err, ErrEmptyGet) {
				return Void{}, false, nil
			}
			return Void{}, v, err
		},
		IxMatchVoid(),
		pred.AsExprHandler(),
		pred.AsExpr,
	)
}
