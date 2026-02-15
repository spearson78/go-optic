package main

import (
	"log"
	"os"
	"strings"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/ohtml"
	"github.com/spearson78/go-optic/otree"

	"golang.org/x/net/html"
)

func main() {
	//Check arguments
	if len(os.Args) < 2 {
		log.Fatalf("filename argument required.")
	}

	//Read html file
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("read file failed: %v", err)
	}

	//Parse bytes to to html
	parseHtml := ohtml.Parse()

	//Traverse all the nodes in the html tree recursively
	allNodesTopDown := otree.TopDown(ohtml.O.Node().TraverseChildren())

	//focus on only TextNodes
	filterTextNodes := Filtered(
		allNodesTopDown,
		ohtml.O.Node().NodeType().Eq(html.TextNode),
	)

	//Focus on the data of the text nodes
	nodeData := ohtml.O.Node().Data()

	//Compose all optics together
	composedOptic := Compose3(
		parseHtml,
		filterTextNodes,
		nodeData,
	)

	var metrics Metrics

	//Now that we are focusing on a string we can use strings.ToUpper to convert all text to uppercase
	newHtml, err := Modify(WithMetrics(composedOptic, &metrics), Op(strings.ToUpper), data)
	if err != nil {
		log.Fatalf("transform failed: %v", err)
	}

	log.Printf("Metrics: Iterations=%v Focused=%v Lookups=%v Comparisons=%v", metrics.Access, metrics.Focused, metrics.IxGet, metrics.Custom["comparisons"])

	//Note that in the output only text nodes are affected. the names of html tags have not been modified
	print(string(newHtml))
}
