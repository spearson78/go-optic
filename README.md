# Go-optics. Modifying immutable data types in Go.

Go-optics provides tooling and functions to modify deeply nested immutable data structures in go.

```go
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
```

Check out the documentation here [go-optics](https://spearson78.github.io/go-optic/)

There are example programs in the examples directory.

The `makelens` tool which generates optics for your structs is available in the cmd directory.

