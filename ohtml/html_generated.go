package ohtml

import (
	"github.com/spearson78/go-optic"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

//TODO: makelens doesn't support recursive structures like html.Node

type lHtml[I any, S any, T any, RET any, RW any, DIR any, ERR any] struct {
	optic.Optic[I, S, T, *html.Node, *html.Node, RET, RW, DIR, ERR]
}

func (s *lHtml[I, S, T, RET, RW, DIR, ERR]) NodeType() optic.MakeLensCmpOps[I, S, T, html.NodeType, RET, RW, optic.UniDir, ERR] {
	return optic.NewMakeLensCmpOps(optic.EErrL(optic.RetL(optic.RwL(optic.Ud(optic.ComposeLeft(s, optic.FieldLens(func(x **html.Node) *html.NodeType {
		return &(*x).Type
	})))))))
}
func (s *lHtml[I, S, T, RET, RW, DIR, ERR]) DataAtom() optic.MakeLensCmpOps[I, S, T, atom.Atom, RET, RW, optic.UniDir, ERR] {
	return optic.NewMakeLensCmpOps(optic.EErrL(optic.RetL(optic.RwL(optic.Ud(optic.ComposeLeft(s, optic.FieldLens(func(x **html.Node) *atom.Atom {
		return &(*x).DataAtom
	})))))))
}
func (s *lHtml[I, S, T, RET, RW, DIR, ERR]) Data() optic.MakeLensOrdOps[I, S, T, string, RET, RW, optic.UniDir, ERR] {
	return optic.NewMakeLensOrdOps(optic.EErrL(optic.RetL(optic.RwL(optic.Ud(optic.ComposeLeft(s, optic.FieldLens(func(x **html.Node) *string {
		return &(*x).Data
	})))))))
}
func (s *lHtml[I, S, T, RET, RW, DIR, ERR]) Namespace() optic.MakeLensOrdOps[I, S, T, string, RET, RW, optic.UniDir, ERR] {
	return optic.NewMakeLensOrdOps(optic.EErrL(optic.RetL(optic.RwL(optic.Ud(optic.ComposeLeft(s, optic.FieldLens(func(x **html.Node) *string {
		return &(*x).Namespace
	})))))))
}
func (s *lHtml[I, S, T, RET, RW, DIR, ERR]) Attr() *sAttribute[I, S, T, RET, RW, optic.UniDir, ERR] {
	return &sAttribute[I, S, T, RET, RW, optic.UniDir, ERR]{
		Optic: optic.EErrL(optic.RetL(optic.RwL(optic.Ud(optic.ComposeLeft(optic.EErrL(optic.RetL(optic.RwL(optic.Ud(optic.ComposeLeft(s, optic.FieldLens(func(x **html.Node) *[]html.Attribute {
			return &(*x).Attr
		})))))), optic.SliceToCol[html.Attribute]()))))),
		o: optic.EErrL(optic.RetL(optic.RwL(optic.Ud(optic.ComposeLeft(s, optic.FieldLens(func(x **html.Node) *[]html.Attribute {
			return &(*x).Attr
		})))))),
	}
}

type lAttribute[I any, S any, T any, RET any, RW any, DIR any, ERR any] struct {
	optic.Optic[I, S, T, html.Attribute, html.Attribute, RET, RW, DIR, ERR]
}

func (s *lAttribute[I, S, T, RET, RW, DIR, ERR]) Namespace() optic.MakeLensOrdOps[I, S, T, string, RET, RW, optic.UniDir, ERR] {
	return optic.NewMakeLensOrdOps(optic.EErrL(optic.RetL(optic.RwL(optic.Ud(optic.ComposeLeft(s, optic.FieldLens(func(x *html.Attribute) *string {
		return &x.Namespace
	})))))))
}

func (s *lAttribute[I, S, T, RET, RW, DIR, ERR]) Key() optic.MakeLensOrdOps[I, S, T, string, RET, RW, optic.UniDir, ERR] {
	return optic.NewMakeLensOrdOps(optic.EErrL(optic.RetL(optic.RwL(optic.Ud(optic.ComposeLeft(s, optic.FieldLens(func(x *html.Attribute) *string {
		return &x.Key
	})))))))
}

func (s *lAttribute[I, S, T, RET, RW, DIR, ERR]) Val() optic.MakeLensOrdOps[I, S, T, string, RET, RW, optic.UniDir, ERR] {
	return optic.NewMakeLensOrdOps(optic.EErrL(optic.RetL(optic.RwL(optic.Ud(optic.ComposeLeft(s, optic.FieldLens(func(x *html.Attribute) *string {
		return &x.Val
	})))))))
}

type sAttribute[I any, S any, T any, RET any, RW any, DIR any, ERR any] struct {
	optic.Optic[I, S, T, optic.Collection[int, html.Attribute, optic.Pure], optic.Collection[int, html.Attribute, optic.Pure], RET, RW, DIR, ERR]
	o optic.Optic[I, S, T, []html.Attribute, []html.Attribute, RET, RW, DIR, ERR]
}

func (s *sAttribute[I, S, T, RET, RW, DIR, ERR]) Traverse() *lAttribute[int, S, T, optic.ReturnMany, optic.CompositionTree[RW, optic.ReadWrite], optic.UniDir, ERR] {
	return OAttributeOf(optic.RetM(optic.Ud(optic.EErrL(optic.Compose(s.o, optic.TraverseSlice[html.Attribute]())))))
}
func (s *sAttribute[I, S, T, RET, RW, DIR, ERR]) Nth(index int) *lAttribute[int, S, T, optic.ReturnMany, optic.CompositionTree[RW, optic.ReadWrite], optic.UniDir, ERR] {
	return OAttributeOf(optic.Index(s.Traverse(), index))
}

type o struct {
}

func (s *o) Node() *lHtml[optic.Void, *html.Node, *html.Node, optic.ReturnOne, optic.ReadWrite, optic.BiDir, optic.Pure] {
	return ONodeOf(optic.Identity[*html.Node]())
}
func (s *o) Attribute() *lAttribute[optic.Void, html.Attribute, html.Attribute, optic.ReturnOne, optic.ReadWrite, optic.BiDir, optic.Pure] {
	return OAttributeOf(optic.Identity[html.Attribute]())
}
func ONodeOf[I any, S any, T any, RET any, RW any, DIR any, ERR any](l optic.Optic[I, S, T, *html.Node, *html.Node, RET, RW, DIR, ERR]) *lHtml[I, S, T, RET, RW, DIR, ERR] {
	return &lHtml[I, S, T, RET, RW, DIR, ERR]{
		Optic: l,
	}
}
func OAttributeOf[I any, S any, T any, RET any, RW any, DIR any, ERR any](l optic.Optic[I, S, T, html.Attribute, html.Attribute, RET, RW, DIR, ERR]) *lAttribute[I, S, T, RET, RW, DIR, ERR] {
	return &lAttribute[I, S, T, RET, RW, DIR, ERR]{
		Optic: l,
	}
}

var O = o{}
