+++
title = "Polymorphic Optics"
weight = 7
+++
# Polymorphic Optics
The optics we have used so far have been monomorphic. This means that the output source and focus types are the same as the input types.
This can be seen in this visualisation.

![Optic](/poly_1.svg)

The types flowing to the right are identical to this flowing to the left. This is the normal case when working with optics created by [makelens](/go-optic/docs/9.makelens).

It is however possible to convert datatypes within an optic. For these use cases you will need to create user defined optics to perform the conversion.

We can imagine a use case where we are upgrading the rating datatype from int to float64 value.
If we could construct an optic with this structure it would be possible. 

![Optic](/poly_2.svg)

This is in fact possible with polymorphic optics.

```go
package main

import (
	"fmt"
	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

// V1 data types
type BlogPost struct {
	author string
	title string
	content string
	comments []Comment
	ratings []Rating
}
  
type Comment struct {
	author string
	title string
	content string
}
  
type Rating struct {
	author string
	stars int
}
  
// V2 data types
type BlogPostV2 struct {
	author string
	title string
	content string
	comments []Comment
	ratings []RatingV2
}
  
type RatingV2 struct {
	author string
	stars float64
}
  
func ExamplePolymorphism() {

	upgradeBlogRatings := LensP[BlogPost, BlogPostV2, []Rating, []RatingV2](
		func(source BlogPost) []Rating {
			return source.ratings
		},
		func(source BlogPost, focus []RatingV2) BlogPostV2 {
			return BlogPostV2{
				//Retain the fields that don't need an upgrade
				author: source.author,
				title: source.title,
				content: source.content,
				comments: source.comments,
	  
				//Use the upgraded focus for ratings
				ratings: focus,
			}
		},
		ExprCustom("upgradeBlogRatings"),
	)
  
	upgradeTraverseSlice := TraverseSliceP[Rating, RatingV2]()
	  
	upgradeRatingStars := LensP[Rating, RatingV2, int, float64](
		func(source Rating) int {
			return source.stars
		},
		func(source Rating, focus float64) RatingV2 {
			return RatingV2{
				author: source.author,
				stars: focus,
			}
		},
		ExprCustom("upgradeRatingStars"),
	)
  
	upgradeOptic := Compose3(
		upgradeBlogRatings,
		upgradeTraverseSlice,
		upgradeRatingStars,
	)
  
	var result BlogPostV2 = MustModify(
		upgradeOptic, 
		IsoCast[int, float64](),
		BlogPost{
			author: "Max Mustermann",
			title: "Polymorphic Optics",
			content: "Lorem ipsum",
			comments: []Comment{},
			ratings: []Rating{
				{
					author: "Max Mustermann",
					rating: 5,
				},
		},
	})
  
	fmt.Println(result)
  
	//Output:
	//{Max Mustermann Polymorphic Optics Lorem ipsum [] [{Max Mustermann 0.5}]}
}
```
In order to change types we have to use the polymorphic version of a lens created with the ``LensP`` constructor.

During the modify we are given the old version of the source and the new version of the focus. We can then perform the upgrade by copying the old values into the new version making sure to use the updated focus in the new copy.
We did not have to implement the upgrade for the slice as ``TraversSlice`` has a polymorphic version ``TraverseSliceP`` where possible all built in functions have a polymorphic version with the ``P`` postfix.