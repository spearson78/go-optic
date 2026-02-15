package optic

import (
	"context"

	"github.com/samber/lo"
	"github.com/spearson78/go-optic/expr"
)

// AndT2 returns an [BinaryOp] that performs a logical and
//
// See:
//   - [And] for a unary version.
//   - [AndOp] for a version that is applied to the focus of 2 [Operation]s.
func AndT2() Optic[Void, lo.Tuple2[bool, bool], lo.Tuple2[bool, bool], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right bool) bool {
		return left && right
	}, "&&")
}

// And returns an [Operation] that performs a logical and with a constant value.
//
// See:
//   - [IAnd] for an [OperationI] version
//   - [AndOp] for a version that is applied to the focus of 2 getters.
func And(right bool) Optic[Void, bool, bool, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2Bind(AndT2(), right)
}

// AndOp returns an [Operation] that logically ands the value of the left getter with the right getters value
//
// See:
//   - [And] for the constant version
func AndOp[S, LERR, RERR any](left Predicate[S, LERR], right Predicate[S, RERR]) Optic[Void, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, CompositionTree[LERR, RERR]] {
	//TODO: consider a short circuit optimization.
	return EErrL(PredT2ToOptic(left, AndT2(), right))
}

// OrT2 returns an [BinaryOp] that performs a logical or
//
// See:
//   - [Or] for a unary version.
//   - [OrOp] for a version that is applied to the focus of 2 [Operation]s.
func OrT2() Optic[Void, lo.Tuple2[bool, bool], lo.Tuple2[bool, bool], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right bool) bool { return left || right }, "||")
}

// Or returns an [Operation] that performs a logical or with a constant value.
//
// See:
//   - [IOr] for an [OperationI] version
//   - [OrOp] for a version that is applied to the focus of 2 getters.
func Or(right bool) Optic[Void, bool, bool, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2Bind(OrT2(), right)
}

// OrOp returns an [Operation] that logically ors the value of the left getter with the right getters value
//
// See:
//   - [Or] for the constant version
func OrOp[S, LERR, RERR any](left Predicate[S, LERR], right Predicate[S, RERR]) Optic[Void, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, CompositionTree[LERR, RERR]] {
	return EErrL(PredT2ToOptic(left, OrT2(), right))
}

// Not returns a [Iso] that logically nots the value.
func Not() Optic[Void, bool, bool, bool, bool, ReturnOne, ReadWrite, BiDir, Pure] {
	return Involuted[bool](
		func(b bool) bool {
			return !b
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.UnaryExpr{
				OpticTypeExpr: ot,
				Op:            "not",
			}
		}),
	)
}

// NotOp is a combinator that applies the [Not] operation to a [Predicate].
func NotOp[S, ERR any](pred Predicate[S, ERR]) Optic[Void, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return Ret1(Ro(Ud(EErrL(Compose(PredToOptic(pred), Not())))))
}

func Case[I, S, T, A, B, RET, RW, DIR, ERR, ERRP any](
	condition Predicate[S, ERRP],
	optic Optic[I, S, T, A, B, RET, RW, DIR, ERR],
) switchCase[I, S, T, A, B, RET, RW, DIR, ERR, ERRP] {
	return switchCase[I, S, T, A, B, RET, RW, DIR, ERR, ERRP]{
		condition: condition,
		optic:     optic,
	}
}

// If returns an [Optic] that executes either the onTrue or onFalse operation depending on the result of the condition.
func If[I, S, T, A, B, RETT, RETF, RWT, RWF, DIRT, DIRF, ERRP, ERRT, ERRF any](condition Predicate[S, ERRP], onTrue Optic[I, S, T, A, B, RETT, RWT, DIRT, ERRT], onFalse Optic[I, S, T, A, B, RETF, RWF, DIRF, ERRF]) Optic[I, S, T, A, B, CompositionTree[RETT, RETF], CompositionTree[RWT, RWF], UniDir, CompositionTree[CompositionTree[ERRT, ERRF], ERRP]] {
	return switchN[CompositionTree[RETT, RETF], CompositionTree[RWT, RWF], CompositionTree[CompositionTree[ERRT, ERRF], ERRP]](
		onFalse,
		Case(
			condition,
			onTrue,
		),
	)
}

func Switch[I, S, T, A, B, RET1, RETD, RW1, RWD, DIR1, DIRD, ERRP1, ERR1, ERRD any](
	case1 switchCase[I, S, T, A, B, RET1, RW1, DIR1, ERR1, ERRP1],
	def Optic[I, S, T, A, B, RETD, RWD, DIRD, ERRD],
) Optic[I, S, T, A, B, CompositionTree[RET1, RETD], CompositionTree[RW1, RWD], UniDir, CompositionTree[CompositionTree[ERR1, ERRP1], ERRD]] {
	return switchN[CompositionTree[RET1, RETD], CompositionTree[RW1, RWD], CompositionTree[CompositionTree[ERR1, ERRP1], ERRD]](
		def,
		case1,
	)
}

func Switch2[I, S, T, A, B, RET1, RW1, DIR1, ERRP1, ERR1, RET2, RW2, DIR2, ERRP2, ERR2, RETD, RWD, DIRD, ERRD any](
	case1 switchCase[I, S, T, A, B, RET1, RW1, DIR1, ERR1, ERRP1],
	case2 switchCase[I, S, T, A, B, RET2, RW2, DIR2, ERR2, ERRP2],
	def Optic[I, S, T, A, B, RETD, RWD, DIRD, ERRD],
) Optic[I, S, T, A, B, CompositionTree[CompositionTree[RET1, RET2], RETD], CompositionTree[CompositionTree[RW1, RW2], RWD], UniDir, CompositionTree[CompositionTree[CompositionTree[ERR1, ERRP1], CompositionTree[ERR2, ERRP2]], ERRD]] {
	return switchN[CompositionTree[CompositionTree[RET1, RET2], RETD], CompositionTree[CompositionTree[RW1, RW2], RWD], CompositionTree[CompositionTree[CompositionTree[ERR1, ERRP1], CompositionTree[ERR2, ERRP2]], ERRD]](
		def,
		case1,
		case2,
	)
}

func Switch3[I, S, T, A, B, RET1, RW1, DIR1, ERRP1, ERR1, RET2, RW2, DIR2, ERRP2, ERR2, RET3, RW3, DIR3, ERRP3, ERR3, RETD, RWD, DIRD, ERRD any](
	case1 switchCase[I, S, T, A, B, RET1, RW1, DIR1, ERR1, ERRP1],
	case2 switchCase[I, S, T, A, B, RET2, RW2, DIR2, ERR2, ERRP2],
	case3 switchCase[I, S, T, A, B, RET3, RW3, DIR3, ERR3, ERRP3],
	def Optic[I, S, T, A, B, RETD, RWD, DIRD, ERRD],
) Optic[I, S, T, A, B, CompositionTree[CompositionTree[CompositionTree[RET1, RET2], RET3], RETD], CompositionTree[CompositionTree[CompositionTree[RW1, RW2], RW3], RWD], UniDir, CompositionTree[CompositionTree[CompositionTree[CompositionTree[ERR1, ERRP1], CompositionTree[ERR2, ERRP2]], CompositionTree[ERR3, ERRP3]], ERRD]] {
	return switchN[CompositionTree[CompositionTree[CompositionTree[RET1, RET2], RET3], RETD], CompositionTree[CompositionTree[CompositionTree[RW1, RW2], RW3], RWD], CompositionTree[CompositionTree[CompositionTree[CompositionTree[ERR1, ERRP1], CompositionTree[ERR2, ERRP2]], CompositionTree[ERR3, ERRP3]], ERRD]](
		def,
		case1,
		case2,
		case3,
	)
}

func Switch4[I, S, T, A, B, RET1, RW1, DIR1, ERRP1, ERR1, RET2, RW2, DIR2, ERRP2, ERR2, RET3, RW3, DIR3, ERRP3, ERR3, RET4, RW4, DIR4, ERRP4, ERR4, RETD, RWD, DIRD, ERRD any](
	case1 switchCase[I, S, T, A, B, RET1, RW1, DIR1, ERR1, ERRP1],
	case2 switchCase[I, S, T, A, B, RET2, RW2, DIR2, ERR2, ERRP2],
	case3 switchCase[I, S, T, A, B, RET3, RW3, DIR3, ERR3, ERRP3],
	case4 switchCase[I, S, T, A, B, RET4, RW4, DIR4, ERR4, ERRP4],
	def Optic[I, S, T, A, B, RETD, RWD, DIRD, ERRD],
) Optic[I, S, T, A, B, CompositionTree[CompositionTree[CompositionTree[RET1, RET2], RET3], RETD], CompositionTree[CompositionTree[CompositionTree[RW1, RW2], RW3], RWD], UniDir, CompositionTree[CompositionTree[CompositionTree[CompositionTree[ERR1, ERRP1], CompositionTree[ERR2, ERRP2]], CompositionTree[CompositionTree[ERR3, ERRP3], CompositionTree[ERR4, ERRP4]]], ERRD]] {
	return switchN[CompositionTree[CompositionTree[CompositionTree[RET1, RET2], RET3], RETD], CompositionTree[CompositionTree[CompositionTree[RW1, RW2], RW3], RWD], CompositionTree[CompositionTree[CompositionTree[CompositionTree[ERR1, ERRP1], CompositionTree[ERR2, ERRP2]], CompositionTree[CompositionTree[ERR3, ERRP3], CompositionTree[ERR4, ERRP4]]], ERRD]](
		def,
		case1,
		case2,
		case3,
		case4,
	)
}

func SwitchN[I, S, T, A, B, RET, RW, DIR, ERR, ERRP any](def Optic[I, S, T, A, B, RET, RW, DIR, ERR], cases ...switchCase[I, S, T, A, B, RET, RW, DIR, ERR, ERRP]) Optic[I, S, T, A, B, RET, RW, UniDir, CompositionTree[ERR, ERRP]] {
	naked := MustGet(
		SliceOf(
			Compose(
				TraverseSlice[switchCase[I, S, T, A, B, RET, RW, DIR, ERR, ERRP]](),
				UpCast[switchCase[I, S, T, A, B, RET, RW, DIR, ERR, ERRP], nakedCase[I, S, T, A, B]](),
			),
			len(cases),
		),
		cases,
	)

	return switchN[RET, RW, CompositionTree[ERR, ERRP]](def, naked...)
}

type switchCase[I, S, T, A, B, RET, RW, DIR, ERR, ERRP any] struct {
	condition Predicate[S, ERRP]
	optic     Optic[I, S, T, A, B, RET, RW, DIR, ERR]
}

func (c switchCase[I, S, T, A, B, RET, RW, DIR, ERR, ERRP]) getCondition() nakedPredicate[S] {
	return c.condition
}

func (c switchCase[I, S, T, A, B, RET, RW, DIR, ERR, ERRP]) getOptic() nakedOptic[I, S, T, A, B] {
	return c.optic
}

type nakedPredicate[S any] interface {
	AsOpGet() OpGetFunc[S, bool]
	AsExprHandler() func(ctx context.Context) (ExprHandler, error)
	AsExpr() expr.OpticExpression
}

type nakedCase[I, S, T, A, B any] interface {
	getCondition() nakedPredicate[S]
	getOptic() nakedOptic[I, S, T, A, B]
}

func switchN[RET, RW, ERR, I, S, T, A, B any](def nakedOptic[I, S, T, A, B], whens ...nakedCase[I, S, T, A, B]) Optic[I, S, T, A, B, RET, RW, UniDir, ERR] {

	selectBranch := func(ctx context.Context, source S) (nakedOptic[I, S, T, A, B], error) {
		for _, when := range whens {

			match, err := PredGet(ctx, when.getCondition(), source)
			if err != nil {
				return nil, err
			}

			if match {
				return when.getOptic(), nil
			}
		}
		return def, nil
	}

	exprHandlers := make([]asExprHandler, 0, (2*len(whens))+1)
	exprHandlers = append(exprHandlers, def)
	for _, o := range whens {
		exprHandlers = append(exprHandlers, o.getCondition())
		exprHandlers = append(exprHandlers, o.getOptic())

	}

	return Omni[I, S, T, A, B, RET, RW, UniDir, ERR](
		func(ctx context.Context, source S) (I, A, error) {
			branch, err := selectBranch(ctx, source)
			if err != nil {
				var i I
				var a A
				return i, a, err
			}

			return branch.AsGetter()(ctx, source)

		},
		func(ctx context.Context, focus B, source S) (T, error) {
			branch, err := selectBranch(ctx, source)
			if err != nil {
				var t T
				return t, err
			}

			return branch.AsSetter()(ctx, focus, source)
		},
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(ValueIE[I, A]) bool) {
				branch, err := selectBranch(ctx, source)
				if err != nil {
					var i I
					var a A
					yield(ValIE(i, a, JoinCtxErr(ctx, err)))
					return
				}

				branch.AsIter()(ctx, source)(yield)
			}
		},
		func(ctx context.Context, source S) (int, error) {
			branch, err := selectBranch(ctx, source)
			if err != nil {
				return 0, err
			}

			return branch.AsLengthGetter()(ctx, source)
		},
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {

			branch, err := selectBranch(ctx, source)
			if err != nil {
				var t T
				return t, err
			}

			return branch.AsModify()(ctx, fmap, source)
		},
		nil, //We can't skip directly to the index as an earlier case may match the predicate.
		func(indexA, indexB I) bool {
			return def.AsIxMatch()(indexA, indexB)
		},
		unsupportedReverseGetter[B, T],
		ExprDefVarArgs(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {

				whenExprs := make([]expr.Case, 0, len(whens))

				for _, o := range whens {
					whenExprs = append(whenExprs, expr.Case{
						Condition: o.getCondition().AsExpr(),
						Default:   o.getOptic().AsExpr(),
					})
				}

				return expr.Switch{
					OpticTypeExpr: ot,
					Whens:         whenExprs,
					Default:       def.AsExpr(),
				}
			},
			exprHandlers...,
		),
	)
}
