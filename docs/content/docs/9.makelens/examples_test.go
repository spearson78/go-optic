package makelens

import (
	"fmt"
	"testing"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/internal/playground/data"
)

func TestSimpleMakeLens(t *testing.T) {

	//BEGIN simple_playground_make_lens
	result := MustGet(
		//BEGIN simple_make_lens
		O.BlogPost().Content().Eq("Blog Content"),
		//END simple_make_lens
		BlogPost{
			Content: "Blog Content",
			Comments: []Comment{
				Comment{
					Title:   "First Comment",
					Content: "Comment 1 content",
				},
				Comment{
					Title:   "Second Comment",
					Content: "Comment 2 content",
				},
			},
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  5,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  3,
				},
			},
		},
	)
	fmt.Println(result)
	//END simple_playground_make_lens
}

func TestCollectionMakeLens(t *testing.T) {

	//BEGIN collection_playground_make_lens
	result := MustModify(
		//BEGIN collection_make_lens
		O.BlogPost().Comments(),
		FilteredCol[int](
			O.Comment().Title().Eq("Second Comment"),
		),
		//END collection_make_lens
		BlogPost{
			Content: "Blog Content",
			Comments: []Comment{
				Comment{
					Title:   "First Comment",
					Content: "Comment 1 content",
				},
				Comment{
					Title:   "Second Comment",
					Content: "Comment 2 content",
				},
			},
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  5,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  3,
				},
			},
		},
	)
	fmt.Println(result)
	//END collection_playground_make_lens
}

func TestMakeLensTraverse(t *testing.T) {

	//BEGIN makelens_playground_traverse
	//BEGIN makelens_traverse
	result := MustModify(
		O.BlogPost().Ratings().Traverse().Stars(),
		Add(1),
		BlogPost{
			//END makelens_traverse
			Content: "Blog Content",
			Comments: []Comment{
				Comment{
					Title:   "First Comment",
					Content: "Comment 1 content",
				},
				Comment{
					Title:   "Second Comment",
					Content: "Comment 2 content",
				},
			},
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  5,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  3,
				},
			},
		},
	)
	fmt.Println(result)
	//END makelens_playground_traverse

}

func TestMakeLensNth(t *testing.T) {

	//BEGIN makelens_playground_nth
	//BEGIN makelens_nth
	result, ok := MustGetFirst(
		O.BlogPost().Comments().Nth(1),
		BlogPost{
			//END makelens_nth
			Content: "Blog Content",
			Comments: []Comment{
				Comment{
					Title:   "First Comment",
					Content: "Comment 1 content",
				},
				Comment{
					Title:   "Second Comment",
					Content: "Comment 2 content",
				},
			},
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  5,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  3,
				},
			},
		},
	)
	fmt.Println(result, ok)
	//END makelens_playground_nth
}

func TestMakeLensOFrom(t *testing.T) {

	//BEGIN makelens_playground_ofrom
	//BEGIN makelens_ofrom
	result := MustModify(
		OBlogPostOf(
			TraverseSlice[BlogPost](),
		).Ratings().Traverse().Stars(),
		Add(1),
		[]BlogPost{
			//END makelens_ofrom
			BlogPost{
				Content: "First Blog",
				Ratings: []Rating{
					Rating{
						Author: "Max Mustermann",
						Stars:  5,
					},
					Rating{
						Author: "Erika Mustermann",
						Stars:  3,
					},
				},
			},
			BlogPost{
				Content: "Second Blog",
				Ratings: []Rating{
					Rating{
						Author: "Max Mustermann",
						Stars:  2,
					},
				},
			},
		},
	)
	fmt.Println(result)
	//END makelens_playground_ofrom

}

func TestMakeLensVirtual(t *testing.T) {

	//BEGIN makelens_playground_virtual
	//BEGIN makelens_virtual
	result, ok := MustGetFirst(
		O.BlogPost().MeanRating(),
		BlogPost{
			//END makelens_virtual
			Content: "First Blog",
			Ratings: []Rating{
				Rating{
					Author: "Max Mustermann",
					Stars:  5,
				},
				Rating{
					Author: "Erika Mustermann",
					Stars:  3,
				},
			},
		},
	)
	fmt.Println(result, ok)
	//END makelens_playground_virtual

}
