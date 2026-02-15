package olex

import (
	"context"
	"fmt"
	"strings"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

//go:generate go run golang.org/x/tools/cmd/stringer@latest --type TokenType
type TokenType int

const (
	TokUnknown TokenType = iota
	TokWhiteSpace
	TokString
	TokNumber
	TokIdentifier
	TokOther
	TokEnd
	TokError
)

type Token struct {
	tokType TokenType
	value   string
}

func NewToken(t TokenType, v string) Token {
	return Token{
		tokType: t,
		value:   v,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("{%v %v}", t.tokType.String(), t.value)
}

type Span[I any] struct {
	start I
	end   I
}

func SpanStart[I any]() Optic[Void, Span[I], Span[I], I, I, ReturnOne, ReadWrite, UniDir, Pure] {
	return FieldLens(func(source *Span[I]) *I { return &source.start })
}

func SpanEnd[I any]() Optic[Void, Span[I], Span[I], I, I, ReturnOne, ReadWrite, UniDir, Pure] {
	return FieldLens(func(source *Span[I]) *I { return &source.end })
}

func (p Span[I]) String() string {
	return fmt.Sprintf("Span{%v:%v}", p.start, p.end)
}

func TokenTokenType() Optic[Void, Token, Token, TokenType, TokenType, ReturnOne, ReadWrite, UniDir, Pure] {
	return FieldLens(func(source *Token) *TokenType { return &source.tokType })
}

func TokenValue() Optic[Void, Token, Token, string, string, ReturnOne, ReadWrite, UniDir, Pure] {
	return FieldLens(func(source *Token) *string { return &source.value })
}

func predGetYield[I any](ctx context.Context, pred PredicateE[rune], r rune, yield func(val ValueIE[I, Token]) bool) (bool, bool) {
	ret, err := PredGet(ctx, pred, r)
	if err != nil {
		var i I
		return false, yield(ValIE(i, Token{}, err))
	}
	return ret, true
}

func Lex(
	isWhiteSpace PredicateE[rune],
	quoteChar rune,
	isNumberPrefix PredicateE[rune],
	isNumber PredicateE[rune],
	isIdentifierPrefix PredicateE[rune],
	isIdentifier PredicateE[rune],
) Optic[Void, Collection[int64, rune, Err], Collection[int64, rune, Err], Collection[Span[int64], Token, Err], Collection[Span[int64], Token, Err], ReturnOne, ReadWrite, BiDir, Err] {
	return LexI[int64](isWhiteSpace, quoteChar, isNumberPrefix, isNumber, isIdentifierPrefix, isIdentifier)
}

type LexExpr struct {
	expr.OpticTypeExpr
	IsWhiteSpace       expr.OpticExpression
	QuoteChar          rune
	IsNumberPrefix     expr.OpticExpression
	IsNumber           expr.OpticExpression
	IsIdentifierPrefix expr.OpticExpression
	IsIdentifier       expr.OpticExpression
}

func (e LexExpr) Short() string {
	return "Lex"
}

func (e LexExpr) String() string {
	return "Lex"
}

func LexI[I any](
	isWhiteSpace PredicateE[rune],
	quoteChar rune,
	isNumberPrefix PredicateE[rune],
	isNumber PredicateE[rune],
	isIdentifierPrefix PredicateE[rune],
	isIdentifier PredicateE[rune],
) Optic[Void, Collection[I, rune, Err], Collection[I, rune, Err], Collection[Span[I], Token, Err], Collection[Span[I], Token, Err], ReturnOne, ReadWrite, BiDir, Err] {
	return CombiIso[ReadWrite, BiDir, Err, Collection[I, rune, Err], Collection[I, rune, Err], Collection[Span[I], Token, Err], Collection[Span[I], Token, Err]](
		func(ctx context.Context, source Collection[I, rune, Err]) (Collection[Span[I], Token, Err], error) {
			return ColIE[Err, Span[I], Token](
				func(ctx context.Context) SeqIE[Span[I], Token] {

					var curType TokenType
					var curValue strings.Builder
					var tokStartIndex I
					var curIndex I

					return func(yield func(ValueIE[Span[I], Token]) bool) {

						handleRune := func(index I, focus rune) bool {
							space, cont := predGetYield(ctx, isWhiteSpace, focus, yield)
							if !cont {
								return false
							}
							if space {
								tokStartIndex = index
								curType = TokWhiteSpace
								curValue.WriteRune(focus)
								return true
							}

							if focus == quoteChar {
								tokStartIndex = index
								curType = TokString
								return true
							}

							num, cont := predGetYield(ctx, isNumberPrefix, focus, yield)
							if !cont {
								return false
							}
							if num {
								tokStartIndex = index
								curType = TokNumber
								curValue.WriteRune(focus)
								return true
							}

							id, cont := predGetYield(ctx, isIdentifierPrefix, focus, yield)
							if !cont {
								return false
							}
							if id {
								tokStartIndex = index
								curType = TokIdentifier
								curValue.WriteRune(focus)
								return true
							}

							curType = TokUnknown
							tokStartIndex = index
							return yield(ValIE(Span[I]{
								index,
								index,
							}, Token{
								tokType: TokOther,
								value:   string([]rune{focus}),
							}, nil))
						}

						source.AsIter()(ctx)(func(val ValueIE[I, rune]) bool {
							index, focus, err := val.Get()
							if err != nil {
								return yield(ValIE(Span[I]{index, index}, Token{}, err))
							}
							curIndex = index
							switch curType {
							case TokUnknown:
								return handleRune(index, focus)
							case TokWhiteSpace:
								space, cont := predGetYield(ctx, isWhiteSpace, focus, yield)
								if !cont {
									return false
								}

								if space {
									curValue.WriteRune(focus)
									return true
								} else {
									if !yield(ValIE(Span[I]{
										tokStartIndex,
										index,
									}, Token{
										tokType: TokWhiteSpace,
										value:   curValue.String(),
									}, nil)) {
										return false
									}

									curValue.Reset()
									return handleRune(index, focus)
								}

							case TokString:
								if focus == quoteChar {
									if !yield(ValIE(Span[I]{
										tokStartIndex,
										index,
									}, Token{
										tokType: TokString,
										value:   curValue.String(),
									}, nil)) {
										return false
									}

									curValue.Reset()
									curType = TokUnknown
									return true
								} else {
									curValue.WriteRune(focus)
									return true
								}
							case TokNumber:

								num, cont := predGetYield(ctx, isNumber, focus, yield)
								if !cont {
									return false
								}
								if !num {
									if !yield(ValIE(Span[I]{
										tokStartIndex,
										index,
									}, Token{
										tokType: TokNumber,
										value:   curValue.String(),
									}, nil)) {
										return false
									}

									curValue.Reset()
									return handleRune(index, focus)
								} else {
									curValue.WriteRune(focus)
									return true
								}

							case TokIdentifier:

								num, cont := predGetYield(ctx, isIdentifier, focus, yield)
								if !cont {
									return false
								}
								if !num {
									if !yield(ValIE(Span[I]{
										tokStartIndex,
										index,
									}, Token{
										tokType: TokIdentifier,
										value:   curValue.String(),
									}, nil)) {
										return false
									}

									curValue.Reset()
									return handleRune(index, focus)
								} else {
									curValue.WriteRune(focus)
									return true
								}
							default:
								panic("invalid state")
							}
						})

						if curType != TokUnknown {
							yield(ValIE(Span[I]{
								tokStartIndex,
								curIndex,
							}, Token{
								tokType: curType,
								value:   curValue.String(),
							}, nil))
						}

					}
				},
				nil,
				func(a, b Span[I]) bool {
					return source.AsIxMatch()(a.start, b.start) && source.AsIxMatch()(a.end, b.end)
				},
				nil,
			), nil

		},
		func(ctx context.Context, focus Collection[Span[I], Token, Err]) (Collection[I, rune, Err], error) {
			return ColIE[Err, I, rune](
				func(ctx context.Context) SeqIE[I, rune] {
					return func(yield func(ValueIE[I, rune]) bool) {
						focus.AsIter()(ctx)(func(val ValueIE[Span[I], Token]) bool {
							index, focus, err := val.Get()
							if err != nil {
								return yield(ValIE(index.start, rune(0), err))
							}

							if focus.tokType == TokString {

								if !yield(ValIE(index.start, quoteChar, nil)) {
									return false
								}

								for _, r := range []rune(focus.value) {
									if !yield(ValIE(index.start, r, nil)) {
										return false
									}
								}

								if !yield(ValIE(index.start, quoteChar, nil)) {
									return false
								}
								return true
							} else {
								for _, r := range []rune(focus.value) {
									if !yield(ValIE(index.start, r, nil)) {
										return false
									}
								}
								return true
							}

						})
					}
				},
				nil,
				func(a, b I) bool {
					return focus.AsIxMatch()(Span[I]{a, a}, Span[I]{b, b})
				},
				nil,
			), nil
		},
		ExprDef(
			func(t expr.OpticTypeExpr) expr.OpticExpression {
				return LexExpr{
					OpticTypeExpr:      t,
					IsWhiteSpace:       isWhiteSpace.AsExpr(),
					QuoteChar:          quoteChar,
					IsNumberPrefix:     isNumberPrefix.AsExpr(),
					IsNumber:           isNumber.AsExpr(),
					IsIdentifierPrefix: isIdentifierPrefix.AsExpr(),
					IsIdentifier:       isIdentifier.AsExpr(),
				}
			},
			isWhiteSpace,
			isNumberPrefix,
			isNumber,
			isIdentifierPrefix,
			isIdentifier,
		),
	)

}
