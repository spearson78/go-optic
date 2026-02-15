package ohtml

import (
	"bytes"
	"context"
	"iter"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
	"golang.org/x/net/html"
)

func Parse() Optic[Void, []byte, []byte, *html.Node, *html.Node, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoEP[[]byte, []byte, *html.Node, *html.Node](
		func(ctx context.Context, b []byte) (*html.Node, error) {
			node, err := html.Parse(bytes.NewReader(b))
			return node, JoinCtxErr(ctx, err)
		},
		func(ctx context.Context, n *html.Node) ([]byte, error) {
			buf := new(bytes.Buffer)
			err := html.Render(buf, n)
			r := buf.Bytes()
			return r, JoinCtxErr(ctx, err)
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.ParseHtml{
				OpticTypeExpr: ot,
			}
		}),
	)
}

func ChildrenCol() Optic[Void, *html.Node, *html.Node, Collection[int, *html.Node, Pure], Collection[int, *html.Node, Pure], ReturnOne, ReadWrite, UniDir, Pure] {

	return CombiLens[ReadWrite, Pure, Void, *html.Node, *html.Node, Collection[int, *html.Node, Pure], Collection[int, *html.Node, Pure]](
		func(ctx context.Context, source *html.Node) (Void, Collection[int, *html.Node, Pure], error) {
			return Void{}, Col[*html.Node](
				func(yield func(*html.Node) bool) {
					child := source.FirstChild
					for child != nil {
						if !yield(child) {
							break
						}
						child = child.NextSibling
					}
				},
				nil,
			), nil
		},
		func(ctx context.Context, focus Collection[int, *html.Node, Pure], source *html.Node) (*html.Node, error) {

			ret := *source
			var lastChild *html.Node

			first := true
			for valIE := range focus.AsIter()(ctx) {
				_, v, err := valIE.Get()
				if err != nil {
					return nil, err
				}
				cloneChild := *v

				if first {
					ret.FirstChild = &cloneChild
					first = false
				} else {
					lastChild.NextSibling = &cloneChild
					cloneChild.PrevSibling = lastChild
				}

				lastChild = &cloneChild
			}

			if lastChild != nil {
				lastChild.NextSibling = nil
			}
			ret.LastChild = lastChild

			return &ret, nil
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.HtmlChildren{
				OpticTypeExpr: ot,
			}
		}),
	)

}

func TraverseChildren() Optic[int, *html.Node, *html.Node, *html.Node, *html.Node, ReturnMany, ReadWrite, UniDir, Pure] {
	return Traversal[*html.Node, *html.Node](
		func(source *html.Node) iter.Seq[*html.Node] {
			return func(yield func(*html.Node) bool) {
				child := source.FirstChild
				for child != nil {
					if !yield(child) {
						break
					}
					child = child.NextSibling
				}
			}
		},
		nil,
		func(fmap func(focus *html.Node) *html.Node, source *html.Node) *html.Node {
			ret := *source
			var lastChild *html.Node

			child := source.FirstChild
			for child != nil {
				cloneChild := *child
				newChild := fmap(&cloneChild)

				if lastChild == nil {
					ret.FirstChild = newChild
				} else {
					lastChild.NextSibling = newChild
					newChild.PrevSibling = lastChild
				}

				lastChild = newChild
				child = child.NextSibling
			}

			if lastChild != nil {
				lastChild.NextSibling = nil
			}
			ret.LastChild = lastChild

			return &ret
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Traverse{
				OpticTypeExpr: ot,
			}
		}),
	)

}
