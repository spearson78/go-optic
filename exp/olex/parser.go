package olex

import (
	"context"
	"fmt"

	"github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

type LexState[I any] struct {
	next func() (Span[I], Token, error, bool)

	curIndex Span[I]
	current  Token
	err      error
}

func LexCurrent[I any]() Optic[Void, *LexState[I], *LexState[I], Token, Token, ReturnOne, ReadWrite, UniDir, Pure] {
	return WithPanic(PtrFieldLensE(func(source *LexState[I]) *Token { return &source.current }))
}

func LexTokenType[I any]() Optic[Void, *LexState[I], *LexState[I], TokenType, TokenType, ReturnOne, ReadWrite, UniDir, Pure] {
	return Ret1(Rw(Ud(EPure(Compose(LexCurrent[I](), TokenTokenType())))))
}

func LexTokenValue[I any]() Optic[Void, *LexState[I], *LexState[I], string, string, ReturnOne, ReadWrite, UniDir, Pure] {
	return Ret1(Rw(Ud(EPure(Compose(LexCurrent[I](), TokenValue())))))
}

func LexTokenPosition[I any]() Optic[Void, *LexState[I], *LexState[I], Span[I], Span[I], ReturnOne, ReadWrite, UniDir, Pure] {
	return WithPanic(PtrFieldLensE(func(source *LexState[I]) *Span[I] { return &source.curIndex }))
}

func (l *LexState[I]) Next() {
	index, t, err, ok := l.next()
	if err != nil {
		l.current = Token{tokType: TokError, value: err.Error()}
		l.err = err
		return
	}
	if !ok {
		l.current = Token{tokType: TokEnd}
		return
	}
	l.curIndex = index
	l.current = t
}

type prefixParslet[I any] func(parse func(l *LexState[I], currPrecedence int) (AstNode[I], error), l *LexState[I]) (AstNode[I], error)
type infixParslet[I any] func(parse func(l *LexState[I], currPrecedence int) (AstNode[I], error), l *LexState[I], left AstNode[I], curPrecedence int) (AstNode[I], error)

type astParserConfig[I any] struct {
	prefixValueParslets map[Token]prefixParslet[I]
	prefixWildParslets  map[TokenType]prefixParslet[I]
	infixValueParslets  map[Token]infixParslet[I]
	infixWildParslets   map[TokenType]infixParslet[I]
	valuePrecedence     map[Token]int
	wildPrecedence      map[TokenType]int
	ixMatch             func(a, b I) bool
}

func NewAstParserConfig[I any](ixmatch func(a, b I) bool) *astParserConfig[I] {
	p := &astParserConfig[I]{}

	p.ixMatch = ixmatch

	p.prefixValueParslets = make(map[Token]prefixParslet[I])
	p.infixValueParslets = make(map[Token]infixParslet[I])
	p.prefixWildParslets = make(map[TokenType]prefixParslet[I])
	p.infixWildParslets = make(map[TokenType]infixParslet[I])
	p.valuePrecedence = make(map[Token]int)
	p.wildPrecedence = make(map[TokenType]int)

	return p
}

func (p *astParserConfig[I]) AddPrefixParselet(t Token, parslet prefixParslet[I]) {
	p.prefixValueParslets[t] = parslet
}

func (p *astParserConfig[I]) AddPrefixTypeParselet(t TokenType, parslet prefixParslet[I]) {
	p.prefixWildParslets[t] = parslet
}

func (p *astParserConfig[I]) AddInfixParselet(t Token, parslet infixParslet[I], precedence int) {
	p.infixValueParslets[t] = parslet
	p.valuePrecedence[t] = precedence
}

func (p *astParserConfig[I]) AddInfixTypeParselet(t TokenType, parslet infixParslet[I], precedence int) {
	p.infixWildParslets[t] = parslet
	p.wildPrecedence[t] = precedence
}
func (p *astParserConfig[I]) getPrecedence(t Token) int {
	precedence, ok := p.valuePrecedence[t]

	if !ok {
		precedence, ok = p.wildPrecedence[t.tokType]
		if !ok {
			precedence = -1
		}
	}

	return precedence
}

func (p *astParserConfig[I]) getPrefixParslet(t Token) (prefixParslet[I], error) {
	prefixParslet, ok := p.prefixValueParslets[t]
	if !ok {
		prefixParslet, ok = p.prefixWildParslets[t.tokType]
		if !ok {
			return nil, fmt.Errorf("Prefix Parslet not found %v", t)
		}
	}

	return prefixParslet, nil
}

func (p *astParserConfig[I]) getInfixParslet(t Token) (infixParslet[I], error) {

	infixParslet, ok := p.infixValueParslets[t]
	if !ok {
		infixParslet, ok = p.infixWildParslets[t.tokType]
		if !ok {
			return nil, fmt.Errorf("Infix Parslet not found %v", t)
		}
	}

	return infixParslet, nil
}

func (p *astParserConfig[I]) parsletParse(l *LexState[I], currPrecedence int) (AstNode[I], error) {

	//Consume whitespace before a prefix parslet
	if l.current.tokType == TokWhiteSpace {
		l.Next()
	}

	prefixParslet, err := p.getPrefixParslet(l.current)
	if err != nil {
		return nil, err
	}

	expr, err := prefixParslet(p.parsletParse, l)
	if err != nil {
		return nil, err
	}

	nextPrec := p.getPrecedence(l.current)
	for l.current.tokType != TokEnd && currPrecedence < nextPrec {

		infixParslet, err := p.getInfixParslet(l.current)
		if err != nil {
			return nil, err
		}

		expr, err = infixParslet(p.parsletParse, l, expr, nextPrec)
		if err != nil {
			return nil, err
		}

		nextPrec = p.getPrecedence(l.current)
	}

	return expr, l.err
}

type ParseAstExpr struct {
	expr.OpticTypeExpr
	Get func(ctx context.Context, source any) (any, error)
	Set func(ctx context.Context, focus any) (any, error)
}

func (e ParseAstExpr) Short() string {
	return "ParseAst"
}

func (e ParseAstExpr) String() string {
	return "ParseAst"
}

func ParseAst[I any](config *astParserConfig[I]) Optic[Void, Collection[Span[I], Token, optic.Err], Collection[Span[I], Token, optic.Err], AstNode[I], AstNode[I], ReturnOne, ReadWrite, BiDir, Err] {

	get := func(ctx context.Context, source Collection[Span[I], Token, Err]) (AstNode[I], error) {

		next, stop := PullIE(source.AsIter()(ctx))
		defer stop()

		l := &LexState[I]{
			next: next,
		}

		l.Next()

		return config.parsletParse(l, -1)
	}

	set := func(ctx context.Context, focus AstNode[I]) (Collection[Span[I], Token, Err], error) {

		return ColIE[Err, Span[I], Token](
			func(ctx context.Context) SeqIE[Span[I], Token] {
				return func(yield func(val ValueIE[Span[I], Token]) bool) {
					focus.Format(func(index Span[I], focus Token) bool {
						cont := yield(ValIE(index, focus, nil))
						return cont
					})
				}
			},
			nil,
			func(a, b Span[I]) bool {
				return config.ixMatch(a.start, b.start) && config.ixMatch(a.end, b.end)
			},
			nil,
		), nil

	}

	return CombiIso[ReadWrite, BiDir, Err, Collection[Span[I], Token, optic.Err], Collection[Span[I], Token, optic.Err], AstNode[I], AstNode[I]](
		get,
		set,
		ExprDef(
			func(t expr.OpticTypeExpr) expr.OpticExpression {
				return ParseAstExpr{
					OpticTypeExpr: t,
					Get: func(ctx context.Context, source any) (any, error) {
						return get(ctx, source.(Collection[Span[I], Token, Err]))
					},
					Set: func(ctx context.Context, focus any) (any, error) {
						return set(ctx, focus.(AstNode[I]))
					},
				}
			},
		),
	)

}
