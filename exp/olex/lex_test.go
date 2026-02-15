package olex

import (
	"fmt"
	"unicode"

	"github.com/gdamore/encoding"
	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/exp/oio"
)

func ExampleLex() {

	data := `let a = 1
print("b")`

	optic := LexI[oio.LinePosition](
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

	res, err := Get(
		SliceOf(
			Filtered(
				Compose4(
					AsReverseGet(oio.Bytes()),
					oio.DecodeFile(encoding.UTF8),
					optic,
					TraverseColE[Span[oio.LinePosition], Token, Err](),
				),
				Compose(
					TokenTokenType(),
					Ne(TokWhiteSpace),
				),
			),
			10,
		),
		[]byte(data),
	)

	fmt.Println(res, err)

	modifyRes, err := Modify(
		Compose3(
			Filtered(
				Compose4(
					AsReverseGet(oio.Bytes()),
					oio.DecodeFile(encoding.UTF8),
					optic,
					TraverseColE[Span[oio.LinePosition], Token, Err](),
				),
				Compose(
					TokenTokenType(),
					Eq(TokNumber),
				),
			),
			TokenValue(),
			ParseInt[int](10, 0),
		),
		Add(10),
		[]byte(data),
	)

	fmt.Println(string(modifyRes), err)

	//Output:
	//[{TokIdentifier let} {TokIdentifier a} {TokOther =} {TokNumber 1} {TokIdentifier print} {TokOther (} {TokString b} {TokOther )}] <nil>
	//let a = 11
	//print("b") <nil>

}
