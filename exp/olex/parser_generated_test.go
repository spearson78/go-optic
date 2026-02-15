package olex_test

import (
	"github.com/samber/mo"
	"github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/exp/oio"
	. "github.com/spearson78/go-optic/exp/olex"
)

type lAstBinaryOp[I any, S any, T any, RET any, RW any, DIR any, ERR any] struct {
	optic.Optic[I, S, T, AstBinaryOp, AstBinaryOp, RET, RW, DIR, ERR]
}

func (s *lAstBinaryOp[I, S, T, RET, RW, DIR, ERR]) Span() optic.Optic[I, S, T, Span[oio.LinePosition], Span[oio.LinePosition], RET, RW, optic.UniDir, ERR] {
	return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s, optic.FieldLens(func(x *AstBinaryOp) *Span[oio.LinePosition] {
		return &x.Span
	}))))))
}
func (s *lAstBinaryOp[I, S, T, RET, RW, DIR, ERR]) Left() optic.Optic[I, S, T, AstNode[oio.LinePosition], AstNode[oio.LinePosition], RET, RW, optic.UniDir, ERR] {
	return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s, optic.FieldLens(func(x *AstBinaryOp) *AstNode[oio.LinePosition] {
		return &x.Left
	}))))))
}
func (s *lAstBinaryOp[I, S, T, RET, RW, DIR, ERR]) Op() optic.MakeLensOrdOps[I, S, T, string, RET, RW, optic.UniDir, ERR] {
	return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s, optic.FieldLens(func(x *AstBinaryOp) *string {
		return &x.Op
	})))))))
}
func (s *lAstBinaryOp[I, S, T, RET, RW, DIR, ERR]) Right() optic.Optic[I, S, T, AstNode[oio.LinePosition], AstNode[oio.LinePosition], RET, RW, optic.UniDir, ERR] {
	return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s, optic.FieldLens(func(x *AstBinaryOp) *AstNode[oio.LinePosition] {
		return &x.Right
	}))))))
}

type sAstBinaryOp[I any, S any, T any, RET any, RW any, DIR any, ERR any] struct {
	optic.Optic[I, S, T, optic.Collection[int, AstBinaryOp, optic.Pure], optic.Collection[int, AstBinaryOp, optic.Pure], RET, RW, DIR, ERR]
	o optic.Optic[I, S, T, []AstBinaryOp, []AstBinaryOp, RET, RW, DIR, ERR]
}

func (s *sAstBinaryOp[I, S, T, RET, RW, DIR, ERR]) Traverse() *lAstBinaryOp[int, S, T, optic.ReturnMany, RW, optic.UniDir, ERR] {
	return OAstBinaryOpOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o, optic.TraverseSlice[AstBinaryOp]()))))))
}
func (s *sAstBinaryOp[I, S, T, RET, RW, DIR, ERR]) Nth(index int) *lAstBinaryOp[int, S, T, optic.ReturnMany, RW, optic.UniDir, ERR] {
	return OAstBinaryOpOf(optic.Index(s.Traverse(), index))
}

type mAstBinaryOp[I comparable, J any, S any, T any, RET any, RW any, DIR any, ERR any] struct {
	optic.Optic[J, S, T, optic.Collection[I, AstBinaryOp, optic.Pure], optic.Collection[I, AstBinaryOp, optic.Pure], RET, RW, DIR, ERR]
	o optic.Optic[J, S, T, map[I]AstBinaryOp, map[I]AstBinaryOp, RET, RW, DIR, ERR]
}

func (s *mAstBinaryOp[I, J, S, T, RET, RW, DIR, ERR]) Traverse() *lAstBinaryOp[I, S, T, optic.ReturnMany, RW, optic.UniDir, ERR] {
	return OAstBinaryOpOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o, optic.TraverseMap[I, AstBinaryOp]()))))))
}
func (s *mAstBinaryOp[I, J, S, T, RET, RW, DIR, ERR]) Key(index I) *lAstBinaryOp[I, S, T, optic.ReturnMany, RW, optic.UniDir, ERR] {
	return OAstBinaryOpOf(optic.Index(s.Traverse(), index))
}

type oAstBinaryOp[I any, S any, T any, RET any, RW any, DIR any, ERR any] struct {
	optic.Optic[I, S, T, *AstBinaryOp, *AstBinaryOp, RET, RW, DIR, ERR]
}

func (s *oAstBinaryOp[I, S, T, RET, RW, DIR, ERR]) Span() optic.Optic[I, S, T, Span[oio.LinePosition], Span[oio.LinePosition], optic.ReturnMany, RW, optic.UniDir, ERR] {
	return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s, optic.PtrFieldLens(func(x *AstBinaryOp) *Span[oio.LinePosition] {
		return &x.Span
	}))))))
}
func (s *oAstBinaryOp[I, S, T, RET, RW, DIR, ERR]) Left() optic.Optic[I, S, T, AstNode[oio.LinePosition], AstNode[oio.LinePosition], optic.ReturnMany, RW, optic.UniDir, ERR] {
	return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s, optic.PtrFieldLens(func(x *AstBinaryOp) *AstNode[oio.LinePosition] {
		return &x.Left
	}))))))
}
func (s *oAstBinaryOp[I, S, T, RET, RW, DIR, ERR]) Op() optic.Optic[I, S, T, string, string, optic.ReturnMany, RW, optic.UniDir, ERR] {
	return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s, optic.PtrFieldLens(func(x *AstBinaryOp) *string {
		return &x.Op
	}))))))
}
func (s *oAstBinaryOp[I, S, T, RET, RW, DIR, ERR]) Right() optic.Optic[I, S, T, AstNode[oio.LinePosition], AstNode[oio.LinePosition], optic.ReturnMany, RW, optic.UniDir, ERR] {
	return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s, optic.PtrFieldLens(func(x *AstBinaryOp) *AstNode[oio.LinePosition] {
		return &x.Right
	}))))))
}
func (s *oAstBinaryOp[I, S, T, RET, RW, DIR, ERR]) Some() *lAstBinaryOp[optic.Void, mo.Option[AstBinaryOp], mo.Option[AstBinaryOp], optic.ReturnMany, optic.ReadWrite, optic.BiDir, optic.Pure] {
	return OAstBinaryOpOf(optic.Some[AstBinaryOp]())
}
func (s *oAstBinaryOp[I, S, T, RET, RW, DIR, ERR]) Option() optic.Optic[I, S, T, mo.Option[AstBinaryOp], mo.Option[AstBinaryOp], RET, RW, DIR, ERR] {
	return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s, optic.PtrOption[AstBinaryOp]())))))
}

type lAstNumberLiteral[I any, S any, T any, RET any, RW any, DIR any, ERR any] struct {
	optic.Optic[I, S, T, AstNumberLiteral, AstNumberLiteral, RET, RW, DIR, ERR]
}

func (s *lAstNumberLiteral[I, S, T, RET, RW, DIR, ERR]) Value() optic.MakeLensRealOps[I, S, T, float64, RET, RW, optic.UniDir, ERR] {
	return optic.NewMakeLensRealOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s, optic.FieldLens(func(x *AstNumberLiteral) *float64 {
		return &x.Value
	})))))))
}
func (s *lAstNumberLiteral[I, S, T, RET, RW, DIR, ERR]) Span() optic.Optic[I, S, T, Span[oio.LinePosition], Span[oio.LinePosition], RET, RW, optic.UniDir, ERR] {
	return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s, optic.FieldLens(func(x *AstNumberLiteral) *Span[oio.LinePosition] {
		return &x.Span
	}))))))
}

type sAstNumberLiteral[I any, S any, T any, RET any, RW any, DIR any, ERR any] struct {
	optic.Optic[I, S, T, optic.Collection[int, AstNumberLiteral, optic.Pure], optic.Collection[int, AstNumberLiteral, optic.Pure], RET, RW, DIR, ERR]
	o optic.Optic[I, S, T, []AstNumberLiteral, []AstNumberLiteral, RET, RW, DIR, ERR]
}

func (s *sAstNumberLiteral[I, S, T, RET, RW, DIR, ERR]) Traverse() *lAstNumberLiteral[int, S, T, optic.ReturnMany, RW, optic.UniDir, ERR] {
	return OAstNumberLiteralOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o, optic.TraverseSlice[AstNumberLiteral]()))))))
}
func (s *sAstNumberLiteral[I, S, T, RET, RW, DIR, ERR]) Nth(index int) *lAstNumberLiteral[int, S, T, optic.ReturnMany, RW, optic.UniDir, ERR] {
	return OAstNumberLiteralOf(optic.Index(s.Traverse(), index))
}

type mAstNumberLiteral[I comparable, J any, S any, T any, RET any, RW any, DIR any, ERR any] struct {
	optic.Optic[J, S, T, optic.Collection[I, AstNumberLiteral, optic.Pure], optic.Collection[I, AstNumberLiteral, optic.Pure], RET, RW, DIR, ERR]
	o optic.Optic[J, S, T, map[I]AstNumberLiteral, map[I]AstNumberLiteral, RET, RW, DIR, ERR]
}

func (s *mAstNumberLiteral[I, J, S, T, RET, RW, DIR, ERR]) Traverse() *lAstNumberLiteral[I, S, T, optic.ReturnMany, RW, optic.UniDir, ERR] {
	return OAstNumberLiteralOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o, optic.TraverseMap[I, AstNumberLiteral]()))))))
}
func (s *mAstNumberLiteral[I, J, S, T, RET, RW, DIR, ERR]) Key(index I) *lAstNumberLiteral[I, S, T, optic.ReturnMany, RW, optic.UniDir, ERR] {
	return OAstNumberLiteralOf(optic.Index(s.Traverse(), index))
}

type oAstNumberLiteral[I any, S any, T any, RET any, RW any, DIR any, ERR any] struct {
	optic.Optic[I, S, T, *AstNumberLiteral, *AstNumberLiteral, RET, RW, DIR, ERR]
}

func (s *oAstNumberLiteral[I, S, T, RET, RW, DIR, ERR]) Value() optic.Optic[I, S, T, float64, float64, optic.ReturnMany, RW, optic.UniDir, ERR] {
	return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s, optic.PtrFieldLens(func(x *AstNumberLiteral) *float64 {
		return &x.Value
	}))))))
}
func (s *oAstNumberLiteral[I, S, T, RET, RW, DIR, ERR]) Span() optic.Optic[I, S, T, Span[oio.LinePosition], Span[oio.LinePosition], optic.ReturnMany, RW, optic.UniDir, ERR] {
	return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s, optic.PtrFieldLens(func(x *AstNumberLiteral) *Span[oio.LinePosition] {
		return &x.Span
	}))))))
}
func (s *oAstNumberLiteral[I, S, T, RET, RW, DIR, ERR]) Some() *lAstNumberLiteral[optic.Void, mo.Option[AstNumberLiteral], mo.Option[AstNumberLiteral], optic.ReturnMany, optic.ReadWrite, optic.BiDir, optic.Pure] {
	return OAstNumberLiteralOf(optic.Some[AstNumberLiteral]())
}
func (s *oAstNumberLiteral[I, S, T, RET, RW, DIR, ERR]) Option() optic.Optic[I, S, T, mo.Option[AstNumberLiteral], mo.Option[AstNumberLiteral], RET, RW, DIR, ERR] {
	return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s, optic.PtrOption[AstNumberLiteral]())))))
}

type lAstVariable[I any, S any, T any, RET any, RW any, DIR any, ERR any] struct {
	optic.Optic[I, S, T, AstVariable, AstVariable, RET, RW, DIR, ERR]
}

func (s *lAstVariable[I, S, T, RET, RW, DIR, ERR]) Name() optic.MakeLensOrdOps[I, S, T, string, RET, RW, optic.UniDir, ERR] {
	return optic.NewMakeLensOrdOps(optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s, optic.FieldLens(func(x *AstVariable) *string {
		return &x.Name
	})))))))
}
func (s *lAstVariable[I, S, T, RET, RW, DIR, ERR]) Span() optic.Optic[I, S, T, Span[oio.LinePosition], Span[oio.LinePosition], RET, RW, optic.UniDir, ERR] {
	return optic.RetL(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s, optic.FieldLens(func(x *AstVariable) *Span[oio.LinePosition] {
		return &x.Span
	}))))))
}

type sAstVariable[I any, S any, T any, RET any, RW any, DIR any, ERR any] struct {
	optic.Optic[I, S, T, optic.Collection[int, AstVariable, optic.Pure], optic.Collection[int, AstVariable, optic.Pure], RET, RW, DIR, ERR]
	o optic.Optic[I, S, T, []AstVariable, []AstVariable, RET, RW, DIR, ERR]
}

func (s *sAstVariable[I, S, T, RET, RW, DIR, ERR]) Traverse() *lAstVariable[int, S, T, optic.ReturnMany, RW, optic.UniDir, ERR] {
	return OAstVariableOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o, optic.TraverseSlice[AstVariable]()))))))
}
func (s *sAstVariable[I, S, T, RET, RW, DIR, ERR]) Nth(index int) *lAstVariable[int, S, T, optic.ReturnMany, RW, optic.UniDir, ERR] {
	return OAstVariableOf(optic.Index(s.Traverse(), index))
}

type mAstVariable[I comparable, J any, S any, T any, RET any, RW any, DIR any, ERR any] struct {
	optic.Optic[J, S, T, optic.Collection[I, AstVariable, optic.Pure], optic.Collection[I, AstVariable, optic.Pure], RET, RW, DIR, ERR]
	o optic.Optic[J, S, T, map[I]AstVariable, map[I]AstVariable, RET, RW, DIR, ERR]
}

func (s *mAstVariable[I, J, S, T, RET, RW, DIR, ERR]) Traverse() *lAstVariable[I, S, T, optic.ReturnMany, RW, optic.UniDir, ERR] {
	return OAstVariableOf(optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.Compose(s.o, optic.TraverseMap[I, AstVariable]()))))))
}
func (s *mAstVariable[I, J, S, T, RET, RW, DIR, ERR]) Key(index I) *lAstVariable[I, S, T, optic.ReturnMany, RW, optic.UniDir, ERR] {
	return OAstVariableOf(optic.Index(s.Traverse(), index))
}

type oAstVariable[I any, S any, T any, RET any, RW any, DIR any, ERR any] struct {
	optic.Optic[I, S, T, *AstVariable, *AstVariable, RET, RW, DIR, ERR]
}

func (s *oAstVariable[I, S, T, RET, RW, DIR, ERR]) Name() optic.Optic[I, S, T, string, string, optic.ReturnMany, RW, optic.UniDir, ERR] {
	return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s, optic.PtrFieldLens(func(x *AstVariable) *string {
		return &x.Name
	}))))))
}
func (s *oAstVariable[I, S, T, RET, RW, DIR, ERR]) Span() optic.Optic[I, S, T, Span[oio.LinePosition], Span[oio.LinePosition], optic.ReturnMany, RW, optic.UniDir, ERR] {
	return optic.RetM(optic.RwL(optic.Ud(optic.EErrL(optic.ComposeLeft(s, optic.PtrFieldLens(func(x *AstVariable) *Span[oio.LinePosition] {
		return &x.Span
	}))))))
}
func (s *oAstVariable[I, S, T, RET, RW, DIR, ERR]) Some() *lAstVariable[optic.Void, mo.Option[AstVariable], mo.Option[AstVariable], optic.ReturnMany, optic.ReadWrite, optic.BiDir, optic.Pure] {
	return OAstVariableOf(optic.Some[AstVariable]())
}
func (s *oAstVariable[I, S, T, RET, RW, DIR, ERR]) Option() optic.Optic[I, S, T, mo.Option[AstVariable], mo.Option[AstVariable], RET, RW, DIR, ERR] {
	return optic.RetL(optic.RwL(optic.DirL(optic.EErrL(optic.ComposeLeft(s, optic.PtrOption[AstVariable]())))))
}

type o struct {
}

func (s *o) AstBinaryOp() *lAstBinaryOp[optic.Void, AstBinaryOp, AstBinaryOp, optic.ReturnOne, optic.ReadWrite, optic.BiDir, optic.Pure] {
	return OAstBinaryOpOf[optic.Void, AstBinaryOp, AstBinaryOp, optic.ReturnOne](optic.Identity[AstBinaryOp]())
}
func (s *o) AstNumberLiteral() *lAstNumberLiteral[optic.Void, AstNumberLiteral, AstNumberLiteral, optic.ReturnOne, optic.ReadWrite, optic.BiDir, optic.Pure] {
	return OAstNumberLiteralOf[optic.Void, AstNumberLiteral, AstNumberLiteral, optic.ReturnOne](optic.Identity[AstNumberLiteral]())
}
func (s *o) AstVariable() *lAstVariable[optic.Void, AstVariable, AstVariable, optic.ReturnOne, optic.ReadWrite, optic.BiDir, optic.Pure] {
	return OAstVariableOf[optic.Void, AstVariable, AstVariable, optic.ReturnOne](optic.Identity[AstVariable]())
}
func OAstBinaryOpOf[I any, S any, T any, RET any, RW any, DIR any, ERR any](l optic.Optic[I, S, T, AstBinaryOp, AstBinaryOp, RET, RW, DIR, ERR]) *lAstBinaryOp[I, S, T, RET, RW, DIR, ERR] {
	return &lAstBinaryOp[I, S, T, RET, RW, DIR, ERR]{
		Optic: l,
	}
}
func OAstNumberLiteralOf[I any, S any, T any, RET any, RW any, DIR any, ERR any](l optic.Optic[I, S, T, AstNumberLiteral, AstNumberLiteral, RET, RW, DIR, ERR]) *lAstNumberLiteral[I, S, T, RET, RW, DIR, ERR] {
	return &lAstNumberLiteral[I, S, T, RET, RW, DIR, ERR]{
		Optic: l,
	}
}
func OAstVariableOf[I any, S any, T any, RET any, RW any, DIR any, ERR any](l optic.Optic[I, S, T, AstVariable, AstVariable, RET, RW, DIR, ERR]) *lAstVariable[I, S, T, RET, RW, DIR, ERR] {
	return &lAstVariable[I, S, T, RET, RW, DIR, ERR]{
		Optic: l,
	}
}

var O = o{}
