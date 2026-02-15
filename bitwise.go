package optic

import (
	"github.com/samber/lo"
	"github.com/spearson78/go-optic/expr"
)

// BitAndT2 returns an [BinaryOp] that applies a bitwise and to the given values.
//
// See:
//   - [BitAnd] for a unary version.
func BitAndT2[A Integer]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right A) A { return left & right }, "&")
}

// BitAnd returns an [Operation] that applies a bitwise and to the given value.
//
// See:
//   - [BitOr] for the bitwise or operator
//   - [BitXor] for the bitwise exclusive or operator
//   - [BitNot] for the bitwise not operator
func BitAnd[A Integer](right A) Optic[Void, A, A, A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2Bind(BitAndT2[A](), right)
}

// BitOrT2 returns an [BinaryOp] that applies a bitwise or to the given values.
//
// See:
//   - [BitOr] for a unary version.
func BitOrT2[A Integer]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right A) A { return left | right }, "|")
}

// BitOr returns an [Operation] that applies a bitwise or to the given value.
//
// See:
//   - [BitAnd] for the bitwise and operator
//   - [BitXor] for the bitwise exclusive or operator
//   - [BitNot] for the bitwise not operator
func BitOr[A Integer](right A) Optic[Void, A, A, A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2Bind(BitOrT2[A](), right)
}

// BitXorT2 returns an [BinaryOp] that applies a bitwise xor to the given values.
//
// See:
//   - [BitXor] for a unary version.
func BitXorT2[A Integer]() Optic[Void, lo.Tuple2[A, A], lo.Tuple2[A, A], A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return BinaryOp(func(left, right A) A { return left ^ right }, "^")
}

// BitXor returns an [Operation] that applies a bitwise xor to the given value.
//
// See:
//   - [BitAnd] for the bitwise and operator
//   - [BitOr] for he bitwise or operator
//   - [BitNot] for the bitwise not operator
func BitXor[A Integer](right A) Optic[Void, A, A, A, A, ReturnOne, ReadWrite, BiDir, Pure] {
	return EErrMerge(OpT2IsoBind(BitXorT2[A](), BitXorT2[A](), right))
}

// BitNot returns an [Operation] that applies a bitwise not.
//
// See:
//   - [BitAnd] for the bitwise and operator
//   - [BitOr] for he bitwise or operator
//   - [BitXor] for the bitwise exclusive or operator
func BitNot[A Integer]() Optic[Void, A, A, A, A, ReturnOne, ReadWrite, BiDir, Pure] {
	return Involuted(
		func(x A) A { return ^x },
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.UnaryExpr{
				OpticTypeExpr: ot,
				Op:            "^",
			}
		}),
	)
}
