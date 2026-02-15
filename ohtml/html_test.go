package ohtml_test

import (
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/ohtml"
	"github.com/spearson78/go-optic/otree"
	"golang.org/x/net/html"
)

func TestHtml(t *testing.T) {

	htmlFile, openErr := os.Open("test.html")
	if openErr != nil {
		t.Fatal(openErr)
	}
	htmlData, readErr := io.ReadAll(htmlFile)
	if readErr != nil {
		t.Fatal(readErr)
	}

	node, err := Get(Parse(), htmlData)
	if err != nil {
		t.Fatal(err)
	}

	topDown := otree.TopDown(O.Node().TraverseChildren())

	if res, err := Get(SliceOf(topDown, 9), node); err != nil || len(res) != 9 {
		t.Fatalf("TopDown %v", len(res))
	}

	bottomUp := otree.BottomUp(O.Node().TraverseChildren())

	newHtml, err := Modify(
		Compose(
			Filtered(
				bottomUp,
				O.Node().NodeType().Eq(html.TextNode),
			),
			O.Node().Data(),
		),
		Op(strings.ToUpper),
		node,
	)

	newHtmlData, err := ReverseGet(Parse(), newHtml)
	if err != nil {
		t.Fatal(err)
	}
	newHtmlStr := string(newHtmlData)
	if newHtmlStr != `<!DOCTYPE html><html><head></head><body><h1>MY FIRST HEADING</h1><p>MY FIRST PARAGRAPH.</p></body></html>` {
		t.Fatalf("NewHtml : '%v'", newHtmlStr)
	}

	origHtmlData, err := ReverseGet(Parse(), node)
	if err != nil {
		t.Fatal(err)
	}
	origHtmlStr := string(origHtmlData)
	if origHtmlStr != `<!DOCTYPE html><html><head></head><body><h1>My First Heading</h1><p>My first paragraph.</p></body></html>` {
		t.Fatalf("OrigHtml : '%v'", origHtmlStr)
	}

	if res, ok, err := GetFirst(O.Node().NthChild(1).NthChild(1).NthChild(1).NthChild(0).Data(), node); err != nil || !ok || res != "My first paragraph." {
		t.Fatalf("Builder %v", res)
	}

	if res, err := Get(SliceOf(Compose(Filtered(topDown, O.Node().NodeType().Eq(html.TextNode)), O.Node().Data()), 2), node); err != nil || !reflect.DeepEqual(res, []string{"My First Heading", "My first paragraph."}) {
		t.Fatalf("Filtered %v", res)
	}

}
