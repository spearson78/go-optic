package main

import (
	"fmt"
	"testing"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/internal/playground/data"
)

//BEGIN intro_boilerplate

// Update the given comment in the blog with the given content and return a new immutable BlogPost
func UpdateCommentContent(source BlogPost, commentIndex int, newContent string) BlogPost {
	//We can't modify the existing comments so we have to clone them
	var updatedComments []Comment
	updatedComments = append(updatedComments, source.Comments...)

	//Comments are stored by value so we can update our copy
	updatedComments[commentIndex].Content = newContent

	//BlogPost was also passed by value so we can update its Comments directly
	source.Comments = updatedComments

	return source
}

//END intro_boilerplate

func TestUpdateCommentContent(t *testing.T) {
	//BEGIN intro_boilerplate_playground
	data := NewBlogPost(
		"BlogPost",
		[]Comment{
			NewComment("First Comment", "comment content"),
			NewComment("Second Comment", "comment content"),
		},
	)

	result := UpdateCommentContent(data, 1, "Updated Content")

	fmt.Println(result)
	//END intro_boilerplate_playground
}

//BEGIN intro_optic_update

// Update the given comment in the blog with the given content and return a new immutable BlogPost
func OpticalUpdateCommentContent(blogPost BlogPost, commentIndex int, newContent string) BlogPost {
	return MustSet(
		O.BlogPost().Comments().Nth(commentIndex).Content(),
		newContent,
		blogPost,
	)
}

//END intro_optic_update

func TestOpticalUpdateCommentContent(t *testing.T) {
	//BEGIN intro_optic_update_playground
	data := NewBlogPost(
		"BlogPost",
		[]Comment{
			NewComment("First Comment", "comment content"),
			NewComment("Second Comment", "comment content"),
		},
	)

	result := OpticalUpdateCommentContent(data, 1, "Updated Content")

	fmt.Println(result)
	//END intro_optic_update_playground
}

func OpticalRead(blogPost BlogPost, commentIndex int) {
	//BEGIN intro_optic_read
	content, found := MustGetFirst(
		O.BlogPost().Comments().Nth(commentIndex).Content(),
		blogPost,
	)
	//END intro_optic_read
	fmt.Println(content, found)
}

func TestOpticalRead(t *testing.T) {
	//BEGIN intro_optic_read_playground
	data := NewBlogPost(
		"BlogPost",
		[]Comment{
			NewComment("First Comment", "comment content"),
			NewComment("Second Comment", "comment content"),
		},
	)

	result, ok := MustGetFirst(
		O.BlogPost().Comments().Nth(0).Content(),
		data,
	)

	fmt.Println(result, ok)
	//END intro_optic_read_playground
}
