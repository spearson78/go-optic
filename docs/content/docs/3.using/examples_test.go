package using

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/internal/playground/data"
)

func TestUsingConcepts(t *testing.T) {
	//BEGIN using_playground_modify

	//BEGIN using_fieldlens
	blogPostRatings := FieldLens(func(source *BlogPost) *[]Rating {
		return &source.Ratings
	})

	ratingAuthor := FieldLens(func(source *Rating) *string {
		return &source.Author
	})

	ratingStars := FieldLens(func(source *Rating) *int {
		return &source.Stars
	})
	//END using_fieldlens

	//BEGIN using_traverse
	traverseRatings := TraverseSlice[Rating]()
	//END using_traverse

	//BEGIN using_eq
	eqMustermann := Eq("Max Mustermann")
	//END using_eq

	//BEGIN using_add1
	incOne := Add(1)
	//END using_add1

	//BEGIN using_compose
	// for _,rating := range blogPost.ratings
	traverseBlogPostRatings := Compose(blogPostRatings, traverseRatings)

	// rating.author == "Max Mustermann"
	ratingAuthorEqMaxMustermann := Compose(ratingAuthor, eqMustermann)

	// if rating.author == "Max Mustermann"
	ifRatingAuthorMustermann := Filtered(traverseBlogPostRatings, ratingAuthorEqMaxMustermann)

	// rating.stars = ....
	matchingRatings := Compose(
		ifRatingAuthorMustermann,
		ratingStars,
	)
	//END using_compose

	//BEGIN using_modify
	result := MustModify(
		matchingRatings,
		incOne,
		BlogPost{
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  0,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  0,
				},
			},
		},
	)
	fmt.Println(result)
	//END using_modify
	//END using_playground_modify

}

func TestUsingConceptsMakeLens(t *testing.T) {

	//BEGIN using_makelens
	result := MustModify(
		Compose(
			Filtered(
				O.BlogPost().Ratings().Traverse(),
				O.Rating().Author().Eq("Max Mustermann"),
			),
			O.Rating().Stars(),
		),
		Add(1),
		BlogPost{
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  0,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  0,
				},
			},
		},
	)
	fmt.Println(result)
	//END using_makelens

}

func TestUsingConceptsQuery(t *testing.T) {
	//BEGIN using_query
	result := MustGet(
		SliceOf(
			Filtered(
				O.BlogPost().Ratings().Traverse(),
				O.Rating().Author().Eq("Max Mustermann"),
			),
			10, //Initial slice length
		),
		BlogPost{
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  0,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  0,
				},
			},
		},
	)
	fmt.Println(result)
	//END using_query
}

func TestUsingFilterIndexed(t *testing.T) {

	//BEGIN using_filter_indexed
	result := MustModify(
		Compose(
			FilteredI(
				O.BlogPost().Ratings().Traverse(),
				OpOnIx[Rating](Even[int]()),
			),
			O.Rating().Stars(),
		),
		Add(1),
		BlogPost{
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  0,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  0,
				},
			},
		},
	)
	fmt.Println(result)
	//END using_filter_indexed

}

func TestUsingConceptsSimplePredicate(t *testing.T) {
	//BEGIN using_simple_predicate
	result := MustModify(
		Compose(
			Filtered(
				O.BlogPost().Ratings().Traverse(),
				O.Rating().Author().Eq("Max Mustermann"),
			),
			O.Rating().Stars(),
		),
		Add(1),
		BlogPost{
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  0,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  0,
				},
			},
		},
	)
	fmt.Println(result)
	//END using_simple_predicate
}

func TestUsingConceptsAndOp(t *testing.T) {
	//BEGIN using_applyand_predicate
	result := MustModify(
		Compose(
			Filtered(
				O.BlogPost().Ratings().Traverse(),
				AndOp(
					O.Rating().Author().Eq("Max Mustermann"),
					O.Rating().Stars().Lt(5),
				),
			),
			O.Rating().Stars(),
		),
		Add(1),
		BlogPost{
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  0,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  0,
				},
			},
		},
	)
	fmt.Println(result)
	//END using_applyand_predicate
}

func TestUsingReordered(t *testing.T) {

	//BEGIN using_reordered
	data := []int{3, 2, 5, 4, 1}

	ordered := Ordered(
		TraverseSlice[int](),
		OrderBy[int](Identity[int]()),
	)

	result := MustGet(
		SliceOf(
			ordered,
			len(data),
		),
		data,
	)

	fmt.Println(result)
	//END using_reordered

	if !reflect.DeepEqual([]any{result}, []any{
		//BEGIN using_reordered_result
		[]int{1, 2, 3, 4, 5},
		//END using_reordered_result
	},
	) {
		t.Fatal(result)
	}
}

func TestUsingReorderedModify(t *testing.T) {

	//BEGIN playground_using_reordered_modify
	data := []int{3, 2, 5, 4, 1}

	ordered := Ordered(
		TraverseSlice[int](),
		OrderBy[int](Identity[int]()),
	)

	//BEGIN using_reordered_modify
	result := MustModify(
		ordered,
		Mul(2),
		data,
	)
	//END using_reordered_modify

	fmt.Println(result)
	//END playground_using_reordered_modify

	if !reflect.DeepEqual([]any{result}, []any{
		//BEGIN using_reordered_modify_result
		[]int{6, 4, 10, 8, 2},
		//END using_reordered_modify_result
	},
	) {
		t.Fatal(result)
	}
}

func TestUsingReorderedModifyTaking(t *testing.T) {

	//BEGIN playground_using_reordered_modify_taking
	data := []int{3, 2, 5, 4, 1}

	ordered := Ordered(
		TraverseSlice[int](),
		OrderBy[int](Identity[int]()),
	)

	//BEGIN using_reordered_modify_taking
	result := MustModify(
		Taking(
			ordered,
			2,
		),
		Mul(10),
		data,
	)
	//END using_reordered_modify_taking

	fmt.Println(result)
	//END playground_using_reordered_modify_taking

	if !reflect.DeepEqual([]any{result}, []any{
		//BEGIN using_reordered_modify_result_taking
		[]int{3, 20, 5, 4, 10},
		//END using_reordered_modify_result_taking
	},
	) {
		t.Fatal(result)
	}
}

func TestUsingCollectionsFilteredCol(t *testing.T) {

	//BEGIN using_collection_filteredcol
	result := MustModify(
		O.BlogPost().Ratings(),
		FilteredCol[int](
			O.Rating().Author().Eq("Max Mustermann"),
		),
		BlogPost{
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  0,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  0,
				},
			},
		},
	)

	fmt.Println(result)
	//END using_collection_filteredcol

}

func TestUsingCollectionsReversedSlice(t *testing.T) {

	//BEGIN using_collection_reversed_slice
	result := MustModify(
		Identity[[]int](),
		ReversedSlice[int](),
		[]int{1, 2, 3},
	)

	fmt.Println(result)
	//END using_collection_reversed_slice

	if !reflect.DeepEqual(result,
		//BEGIN using_collection_reversed_slice_result
		[]int{3, 2, 1},
		//END using_collection_reversed_slice_result
	) {
		t.Fatal(result)
	}

}
