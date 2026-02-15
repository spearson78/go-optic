package otree

import (
	"github.com/samber/mo"
	"github.com/spearson78/go-optic"
)

type loTree[X, I, J any, S any, T any, RET any, RW any, DIR any, ERR any] struct {
	optic.Optic[J, S, T, mo.Option[TreeNode[I, X, ERR]], mo.Option[TreeNode[I, X, ERR]], RET, RW, DIR, ERR]
	ixMatch func(a, b I) bool
}

func (s *loTree[X, I, J, S, T, RET, RW, DIR, ERR]) Value() optic.Optic[J, S, T, X, X, optic.ReturnMany, RW, optic.UniDir, ERR] {

	x := optic.FieldLens(func(x *TreeNode[I, X, ERR]) *X {
		return &x.value
	})
	y := optic.RetM(optic.Rw(optic.Ud(optic.EPure(optic.ComposeLeft(optic.Some[TreeNode[I, X, ERR]](), x)))))

	z := optic.ComposeLeft(s, y)

	a := optic.EErrL(optic.RetM(optic.RwL(optic.Ud(z))))

	return a
}
func (s *loTree[X, I, J, S, T, RET, RW, DIR, ERR]) Children() *sTree[X, I, J, S, T, optic.ReturnMany, RW, optic.UniDir, ERR] {

	x := optic.FieldLens(func(x *TreeNode[I, X, ERR]) *optic.Collection[I, TreeNode[I, X, ERR], ERR] {
		return &x.children
	})

	y := optic.RetM(optic.Rw(optic.Ud(optic.EPure(optic.ComposeLeft(optic.Some[TreeNode[I, X, ERR]](), x)))))

	z := optic.ComposeLeft(s, y)

	a := optic.EErrL(optic.RetM(optic.RwL(optic.Ud(z))))

	return &sTree[X, I, J, S, T, optic.ReturnMany, RW, optic.UniDir, ERR]{
		Optic:   a,
		ixMatch: s.ixMatch,
	}
}

type lTree[X, I, J any, S any, T any, RET any, RW any, DIR any, ERR any] struct {
	optic.Optic[J, S, T, TreeNode[I, X, ERR], TreeNode[I, X, ERR], RET, RW, DIR, ERR]
	ixMatch func(a, b I) bool
}

func (s *lTree[X, I, J, S, T, RET, RW, DIR, ERR]) Value() optic.Optic[J, S, T, X, X, RET, RW, optic.UniDir, ERR] {

	x := optic.FieldLens(func(x *TreeNode[I, X, ERR]) *X {
		return &x.value
	})

	z := optic.ComposeLeft(s, x)

	a := optic.EErrL(optic.RetL(optic.RwL(optic.Ud(z))))

	return a
}
func (s *lTree[X, I, J, S, T, RET, RW, DIR, ERR]) Children() *sTree[X, I, J, S, T, RET, RW, optic.UniDir, ERR] {

	x := optic.FieldLens(func(x *TreeNode[I, X, ERR]) *optic.Collection[I, TreeNode[I, X, ERR], ERR] {
		return &x.children
	})

	z := optic.ComposeLeft(s, x)

	a := optic.EErrL(optic.RetL(optic.RwL(optic.Ud(z))))

	return &sTree[X, I, J, S, T, RET, RW, optic.UniDir, ERR]{
		Optic:   a,
		ixMatch: s.ixMatch,
	}
}

func (s *loTree[X, I, J, S, T, RET, RW, DIR, ERR]) TopDown() optic.Optic[*PathNode[I], S, T, mo.Option[TreeNode[I, X, ERR]], mo.Option[TreeNode[I, X, ERR]], optic.ReturnMany, RW, optic.UniDir, optic.Err] {
	x := optic.Compose(optic.Compose(s.Optic, optic.Some[TreeNode[I, X, ERR]]()), optic.ComposeLeft(TopDown(TraverseTreeChildrenI[I, X, ERR](s.ixMatch)), optic.AsReverseGet(optic.SomeE[TreeNode[I, X, ERR]]())))
	return optic.RetM(optic.RwL(optic.RwL(optic.Ud(optic.EErr(x)))))
}

func (s *loTree[X, I, J, S, T, RET, RW, DIR, ERR]) BottomUp() optic.Optic[*PathNode[I], S, T, mo.Option[TreeNode[I, X, ERR]], mo.Option[TreeNode[I, X, ERR]], optic.ReturnMany, RW, optic.UniDir, optic.Err] {
	x := optic.Compose(optic.Compose(s.Optic, optic.Some[TreeNode[I, X, ERR]]()), optic.ComposeLeft(BottomUp(TraverseTreeChildrenI[I, X, ERR](s.ixMatch)), optic.AsReverseGet(optic.SomeE[TreeNode[I, X, ERR]]())))
	return optic.RetM(optic.RwL(optic.RwL(optic.Ud(optic.EErr(x)))))
}

type sTree[X, I, J any, S any, T any, RET any, RW any, DIR any, ERR any] struct {
	optic.Optic[J, S, T, optic.Collection[I, TreeNode[I, X, ERR], ERR], optic.Collection[I, TreeNode[I, X, ERR], ERR], RET, RW, DIR, ERR]
	ixMatch func(a, b I) bool
}

func (s *sTree[X, I, J, S, T, RET, RW, DIR, ERR]) Traverse() *lTree[X, I, I, S, T, optic.ReturnMany, RW, optic.UniDir, ERR] {

	x := optic.EErrMerge(optic.Compose(s.Optic, optic.TraverseColIEP[I, TreeNode[I, X, ERR], TreeNode[I, X, ERR], ERR](s.ixMatch)))

	j := optic.RetM(optic.RwL(optic.Ud(x)))

	c := OTreeFromI(j, s.ixMatch)

	return c
}
func (s *sTree[X, I, J, S, T, RET, RW, DIR, ERR]) Nth(index I) *lTree[X, I, I, S, T, optic.ReturnMany, RW, optic.UniDir, ERR] {

	x := optic.Index(s.Traverse(), index)

	return OTreeFromI(x, s.ixMatch)
}

func OTree[I comparable, X, ERR any]() *lTree[X, I, optic.Void, TreeNode[I, X, ERR], TreeNode[I, X, ERR], optic.ReturnOne, optic.ReadWrite, optic.BiDir, ERR] {
	return OTreeI[I, X, ERR](optic.IxMatchComparable[I]())

}

func OTreeOpt[I comparable, X, ERR any]() *loTree[X, I, optic.Void, mo.Option[TreeNode[I, X, ERR]], mo.Option[TreeNode[I, X, ERR]], optic.ReturnOne, optic.ReadWrite, optic.BiDir, ERR] {
	return OTreeOptI[I, X, ERR](optic.IxMatchComparable[I]())

}

func OTreeFrom[X any, I comparable, J any, S any, T any, RET any, RW any, DIR any, ERR any](l optic.Optic[J, S, T, TreeNode[I, X, ERR], TreeNode[I, X, ERR], RET, RW, DIR, ERR]) *lTree[X, I, J, S, T, RET, RW, DIR, ERR] {
	return OTreeFromI(l, optic.IxMatchComparable[I]())
}

func OTreeOptFrom[X any, I comparable, J any, S any, T any, RET any, RW any, DIR any, ERR any](l optic.Optic[J, S, T, mo.Option[TreeNode[I, X, ERR]], mo.Option[TreeNode[I, X, ERR]], RET, RW, DIR, ERR]) *loTree[X, I, J, S, T, RET, RW, DIR, ERR] {
	return OTreeOptFromI(l, optic.IxMatchComparable[I]())
}

func OTreeI[I any, X, ERR any](ixmatch func(a, b I) bool) *lTree[X, I, optic.Void, TreeNode[I, X, ERR], TreeNode[I, X, ERR], optic.ReturnOne, optic.ReadWrite, optic.BiDir, ERR] {
	return &lTree[X, I, optic.Void, TreeNode[I, X, ERR], TreeNode[I, X, ERR], optic.ReturnOne, optic.ReadWrite, optic.BiDir, ERR]{
		Optic:   optic.CombiEErr[ERR](optic.Identity[TreeNode[I, X, ERR]]()),
		ixMatch: ixmatch,
	}
}

func OTreeOptI[I any, X, ERR any](ixmatch func(a, b I) bool) *loTree[X, I, optic.Void, mo.Option[TreeNode[I, X, ERR]], mo.Option[TreeNode[I, X, ERR]], optic.ReturnOne, optic.ReadWrite, optic.BiDir, ERR] {
	return &loTree[X, I, optic.Void, mo.Option[TreeNode[I, X, ERR]], mo.Option[TreeNode[I, X, ERR]], optic.ReturnOne, optic.ReadWrite, optic.BiDir, ERR]{
		Optic:   optic.CombiEErr[ERR](optic.Identity[mo.Option[TreeNode[I, X, ERR]]]()),
		ixMatch: ixmatch,
	}
}

func OTreeFromI[X any, I any, J any, S any, T any, RET any, RW any, DIR any, ERR any](l optic.Optic[J, S, T, TreeNode[I, X, ERR], TreeNode[I, X, ERR], RET, RW, DIR, ERR], ixmatch func(a, b I) bool) *lTree[X, I, J, S, T, RET, RW, DIR, ERR] {
	return &lTree[X, I, J, S, T, RET, RW, DIR, ERR]{
		Optic:   l,
		ixMatch: ixmatch,
	}
}

func OTreeOptFromI[X any, I any, J any, S any, T any, RET any, RW any, DIR any, ERR any](l optic.Optic[J, S, T, mo.Option[TreeNode[I, X, ERR]], mo.Option[TreeNode[I, X, ERR]], RET, RW, DIR, ERR], ixmatch func(a, b I) bool) *loTree[X, I, J, S, T, RET, RW, DIR, ERR] {
	return &loTree[X, I, J, S, T, RET, RW, DIR, ERR]{
		Optic:   l,
		ixMatch: ixmatch,
	}
}
