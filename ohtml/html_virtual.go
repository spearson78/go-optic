package ohtml

import (
	. "github.com/spearson78/go-optic"
	"golang.org/x/net/html"
)

func (s *lHtml[I, S, T, RET, RW, DIR, ERR]) Children() Optic[Void, S, T, Collection[int, *html.Node, Pure], Collection[int, *html.Node, Pure], RET, RW, UniDir, ERR] {
	return RetL(RwL(Ud(EErrL(Compose(s, ChildrenCol())))))
}

func (s *lHtml[I, S, T, RET, RW, DIR, ERR]) TraverseChildren() *lHtml[int, S, T, ReturnMany, RW, UniDir, ERR] {
	return ONodeOf(RetM(RwL(Ud(EErrL(Compose(s, TraverseChildren()))))))
}

func (s *lHtml[I, S, T, RET, RW, DIR, ERR]) NthChild(index int) *lHtml[int, S, T, ReturnMany, RW, UniDir, ERR] {
	return ONodeOf(RwL(EErrL(Index(Compose(s, TraverseChildren()), index))))
}
