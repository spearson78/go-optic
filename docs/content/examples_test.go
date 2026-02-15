package main

import (
	"fmt"
	"testing"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/internal/playground/data"
)

func TestOverview(t *testing.T) {
	//BEGIN overview

	//Create a nested immutable data structure.
	data := NewBlogPost(
		"BlogPost",
		[]Comment{
			NewComment("First Comment", "First comment content"),
			NewComment("Second Comment", "Second comment content"),
		},
	)

	//Create an optic to focus the first comments content.
	optic := O.BlogPost().Comments().Nth(0).Content()

	//Set a new value to the focused comment.
	newData := MustSet(
		optic,
		"Updated comment content",
		data,
	)

	fmt.Println(newData)
	//Output:
	//{BlogPost [{First Comment Updated comment content} {Second Comment Second comment content}] []}

	//END overview
}
