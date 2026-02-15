package olex_test

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/gdamore/encoding"
	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/exp/oio"
	. "github.com/spearson78/go-optic/exp/olex"
	"github.com/spearson78/go-optic/oreflect"
	"github.com/spearson78/go-optic/otree"
)

//go:generate ../../makelens -dotImport "github.com/spearson78/go-optic/exp/olex" -import "github.com/spearson78/go-optic/oio" olex_test parser_test.go parser_generated_test.go

type AstNumberLiteral struct {
	Value float64
	Span  Span[oio.LinePosition]
}

func (n AstNumberLiteral) Format(yield func(index Span[oio.LinePosition], token Token) bool) bool {
	return yield(n.Span, NewToken(TokNumber, strconv.FormatFloat(float64(n.Value), 'g', -1, 64)))
}

func (n AstNumberLiteral) String() string {
	return AstNodeToString[oio.LinePosition](n)
}

type AstBinaryOp struct {
	Span Span[oio.LinePosition]

	Left  AstNode[oio.LinePosition]
	Op    string
	Right AstNode[oio.LinePosition]
}

func (n *AstBinaryOp) Format(yield func(index Span[oio.LinePosition], token Token) bool) bool {
	if !yield(n.Span, NewToken(TokOther, "(")) {
		return false
	}

	if !n.Left.Format(yield) {
		return false
	}

	if !yield(n.Span, NewToken(TokOther, n.Op)) {
		return false
	}

	if !n.Right.Format(yield) {
		return false
	}

	if !yield(n.Span, NewToken(TokOther, ")")) {
		return false
	}

	return true
}

func (n *AstBinaryOp) String() string {
	return AstNodeToString(n)
}

type AstVariable struct {
	Name string
	Span Span[oio.LinePosition]
}

func (v AstVariable) Format(yield func(index Span[oio.LinePosition], token Token) bool) bool {
	return yield(v.Span, NewToken(TokIdentifier, v.Name))
}

func (n AstVariable) String() string {
	return AstNodeToString(n)
}

func ExampleParseAst() {

	data := `c=(f-32)/1.8`

	//Lex

	lex := LexI[oio.LinePosition](
		Op(unicode.IsSpace),
		'"',
		Op(unicode.IsDigit),
		OrOp(
			Op(unicode.IsDigit),
			In('e', '.', 'x', '-'),
		),
		OrOp(
			Op(unicode.IsLetter),
			In('_'),
		),
		OrOp(
			OrOp(
				Op(unicode.IsLetter),
				In('_'),
			),
			Op(unicode.IsDigit),
		),
	)

	//parse

	const (
		assign = 70
		sum    = 90
		prod   = 100
	)

	parseCfg := NewAstParserConfig(IxMatchComparable[oio.LinePosition]())
	parseCfg.AddPrefixTypeParselet(TokNumber, func(parse func(l *LexState[oio.LinePosition], currPrecedence int) (AstNode[oio.LinePosition], error), l *LexState[oio.LinePosition]) (AstNode[oio.LinePosition], error) {
		v, err := strconv.ParseFloat(MustGet(LexTokenValue[oio.LinePosition](), l), 64)
		l.Next()
		return AstNumberLiteral{
			v,
			MustGet(LexTokenPosition[oio.LinePosition](), l),
		}, err
	})

	binaryOpParselet := func(parse func(l *LexState[oio.LinePosition], currPrecedence int) (AstNode[oio.LinePosition], error), l *LexState[oio.LinePosition], left AstNode[oio.LinePosition], curPrecedence int) (AstNode[oio.LinePosition], error) {
		res := &AstBinaryOp{
			Span: MustGet(LexTokenPosition[oio.LinePosition](), l),
		}
		res.Op = MustGet(LexTokenValue[oio.LinePosition](), l)
		res.Left = left

		l.Next()

		right, err := parse(l, curPrecedence)
		res.Right = right

		return res, err
	}

	parseCfg.AddInfixParselet(NewToken(TokOther, "="), binaryOpParselet, assign)
	parseCfg.AddInfixParselet(NewToken(TokOther, "+"), binaryOpParselet, sum)
	parseCfg.AddInfixParselet(NewToken(TokOther, "-"), binaryOpParselet, sum)
	parseCfg.AddInfixParselet(NewToken(TokOther, "/"), binaryOpParselet, prod)
	parseCfg.AddInfixParselet(NewToken(TokOther, "*"), binaryOpParselet, prod)

	parseCfg.AddPrefixTypeParselet(TokIdentifier, func(parse func(l *LexState[oio.LinePosition], currPrecedence int) (AstNode[oio.LinePosition], error), l *LexState[oio.LinePosition]) (AstNode[oio.LinePosition], error) {
		node := AstVariable{
			MustGet(LexTokenValue[oio.LinePosition](), l),
			MustGet(LexTokenPosition[oio.LinePosition](), l),
		}
		l.Next()
		return node, nil
	})

	parseCfg.AddPrefixParselet(NewToken(TokOther, "("), func(parse func(l *LexState[oio.LinePosition], currPrecedence int) (AstNode[oio.LinePosition], error), l *LexState[oio.LinePosition]) (AstNode[oio.LinePosition], error) {
		l.Next()

		exp, err := parse(l, 0) //parenthesis reset the precedence
		if err != nil {
			return nil, err
		}

		x := NewToken(TokOther, ")")
		if MustGet(LexCurrent[oio.LinePosition](), l) != x {
			return nil, errors.New("expected )")
		}
		l.Next()

		return exp, nil
	})

	parse := ParseAst(parseCfg)

	//Optic

	parseString := Compose5(
		IsoCast[string, []byte](), //Convert string to []byte
		AsReverseGet(oio.Bytes()), //Convert the []byte to a collection of extents compatible with oio
		oio.DecodeFile(encoding.UTF8),
		lex,
		parse,
	)

	astNode, err := Get(parseString, data)

	fmt.Println(astNode, err)

	code, err := Modify(
		Compose4(
			parseString,
			otree.TopDown(oreflect.TraverseMembers[AstNode[oio.LinePosition], AstNode[oio.LinePosition]]()), //Recursively drill down into all nodes in the AST.
			DownCast[AstNode[oio.LinePosition], AstVariable](),                                              //Identify any variables
			O.AstVariable().Name(), //Focus on their names
		),
		Op(strings.ToUpper), //Convert the variable names to upper case.
		data,
	)
	fmt.Println(code, err)

	//Output:
	//(c=((f-32)/1.8)) <nil>
	//(C=((F-32)/1.8)) <nil>
}
