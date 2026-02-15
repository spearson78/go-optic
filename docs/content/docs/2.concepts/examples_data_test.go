package main_test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/internal/playground/data"
)

func TestConceptsOpticsComposeFieldLensMakeLens(t *testing.T) {
	//BEGIN optics_compose_makelens
	for commentTitle := range MustGet(
		SeqOf(
			O.BlogPost().Comments().Traverse().Title(),
		),
		NewBlogPost(
			"Content",
			[]Comment{
				NewComment("First Comment", "My comment"),
				NewComment("Second Comment", "Another comment"),
			},
		),
	) {
		fmt.Println(commentTitle)
	}
	//END optics_compose_makelens
}

func TestConceptsOpticsComposeLeftMakeLens(t *testing.T) {
	//BEGIN optics_compose_left_makelens
	index, result, ok := MustGetFirstI(
		O.BlogPost().Comments().Traverse().Title(),
		NewBlogPost(
			"Content",
			[]Comment{
				NewComment("First Comment", "My comment"),
				NewComment("Second Comment", "Another comment"),
			},
		),
	)
	fmt.Println(index, result, ok)
	//END optics_compose_left_makelens

	if !reflect.DeepEqual([]any{index, result, ok}, []any{
		//BEGIN optics_compose_left_makelens_result
		0, "First Comment", true,
		//END optics_compose_left_makelens_result
	},
	) {
		t.Fatal(index, result, ok)
	}
}

func TestConceptsOpticsLens(t *testing.T) {
	//BEGIN optics_lens_makelens
	result := MustSet(
		O.BlogPost().Content(),
		"New Content",
		NewBlogPost(
			"Content",
			nil,
		),
	)
	fmt.Println(result)
	//END optics_lens_makelens
}
