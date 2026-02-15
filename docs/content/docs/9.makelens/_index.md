+++
title = "makelens Tool"
weight = 9
+++
# MakeLens
The makelens tool generates the boilerplate code for creating and composing `lenses` for structs.
To ensure immutability it is advised to move your structs to a separate package and make all the fields unexported. This prevents any accidental mutation. The generated lenses should be the only access your struct fields.

The makelens tool takes 3 parameters
  1. the package name to generate
  2. the go file to scan for types
  3. the filename to generate
  
```bash
makelens types types.go types_generated.go
```

The tool generates several exported symbols with the prefix `O` 
In this tutorial we will consider the following structs
```go
type Comment struct {
	Title   string
	Content string
}

type Rating struct {
	Author string
	Stars  int
}

type BlogPost struct {
	Content  string
	Comments []Comment
	Ratings  []Rating
}
```
Note that all the fields are unexported.
Running make lens will generate the following exported Symbols.


| Symbol                 | Usage                                                                                                     |
| ---------------------- | --------------------------------------------------------------------------------------------------------- |
| `O.BlogPost()`         | Provides access to the `BlogPost` lenses & composition from with a source of `BlogPost`.                  |
| `O.Comment()`          | Provides access to the `Comment` lenses & composition from with a source of `Comment`.                    |
| `O.Rating()`           | Provides access to the `Rating` lenses & composition from with a source of `Rating`                       |
| `OBlogPostOf(optic)`   | Provides access to the `BlogPost` lenses & composition from an arbitrary optic that focuses a `BlogPost`. |
| `OCommentOf(optic)`    | Provides access to the `Comment` lenses & composition from an arbitrary optic that focuses a `Comment`.   |
| `ORatingOf(optic)`     | Provides access to the `Comment` lenses & composition from an arbitrary optic that focuses a `Comment`.   |

For each struct a lens is generated for each field e.g. `O.BlogPost().Content()` provides access to the `BlogPost.content` field. Additionally for comparable types access to basic predicates is also provided.

{{< playground file="/content/docs/9.makelens/examples_test.go" id="simple_make_lens" playgroundid="simple_playground_make_lens">}}

Returns a `Predicate` that returns true if a `BlogPost.author` is `"Max Mustermann"`

When the field has a collection type (slice or map) then the generate field lens will not provide direct access to the slice or map but rather an immutable `Collection` wrapper. This ensures the underlying map or slice is immutable. The content of the `Collection` can then be modified using the provided collection operations. e.g.

{{< playground file="/content/docs/9.makelens/examples_test.go" id="collection_make_lens" playgroundid="collection_playground_make_lens">}}

Will return a new `BlogPost` with only comments titled `"Second Comment"` remaining.

For collections a `Traverse()` method is provided that performs the relevant `TraverseSlice`or `TraverseMap` to focus the contained elements, after traversing the nested fields are accessible via the generated code. e.g.

{{< playground file="/content/docs/9.makelens/examples_test.go" id="makelens_traverse" playgroundid="makelens_playground_traverse">}}

Will add 1 to every rating for the given blog post.

For slices an `Nth` function is provided for convenient access to the element at a given index.

{{< playground file="/content/docs/9.makelens/examples_test.go" id="makelens_nth" playgroundid="makelens_playground_nth">}}


is the equivalent of 
```go
thirdComment := blogPost.Comments[1]
```

For maps a `Key` function is provided to retrieve the element with the given key.

See the [Collections](/docs/14.collections) section for more information on `Collection`s

The `O____Of` functions are used when the source object has a different root e.g.

{{< playground file="/content/docs/9.makelens/examples_test.go" id="makelens_ofrom" playgroundid="makelens_playground_ofrom">}}


Will add 1 to every rating in every blog post in the given slice.

The set of lenses for a given struct can be extended by adding additional methods to the generated builder struct.

For `BlogPost` the generated builder is
```go
type lBlogPost[I any,S any,T any,RET any,RW any,DIR any,ERR any] struct {
	optic.Optic[I,S,T,BlogPost,BlogPost,RET,RW,DIR,ERR]
}
```

This is a thin wrapper around the underlying optic to access the `BlogPost` 

Additional "virtual" fields can be added by adding a method like this.
```go
import . "github.com/spearson78/go-optic"

func (s *lBlogPost[I, S, T, RET, RW, DIR, ERR]) MeanRating() Optic[Void, S, S, int, int, ReturnMany, ReadOnly, UniDir, CompositionTree[ERR, Pure]] {
	return Reduce(
		s.Ratings().Traverse().Stars(),
		Mean[int](),
	)
}
```

This will provide a `MeanRating` that can be accessed like this.

{{< playground file="/content/docs/9.makelens/examples_test.go" id="makelens_virtual" playgroundid="makelens_playground_virtual">}}
