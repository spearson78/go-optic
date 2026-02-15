package oio

import (
	"unsafe"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

// BytesString returns an Iso that performs an unsafe casts from []byte to a string
//
// WARNING: This cast is unsafe. The []byte MUST NOT BE ALTERED.
// See:
//   - [unsafe.String]
//   - [unsafe.StringData]
//   - [unsafe.Slice]
func BytesString() Optic[Void, []byte, []byte, string, string, ReturnOne, ReadWrite, BiDir, Pure] {
	return Iso[[]byte, string](
		func(focus []byte) string {
			p := unsafe.SliceData(focus)
			return unsafe.String(p, len(focus))
		},
		func(source string) []byte {
			p := unsafe.StringData(source)
			b := unsafe.Slice(p, len(source))
			return b
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Cast{
				OpticTypeExpr: ot,
			}
		}),
	)
}
