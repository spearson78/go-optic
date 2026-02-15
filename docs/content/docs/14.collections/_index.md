+++
title = "Collections"
weight = 14
+++
# Collections

In go-optics `Collection` provides an immutable wrapper around a concrete collection like a slice or map. Go-optics also offers collection operations that modify these. In conjunction ths enables these generic operations to be performed on arbitrary concrete collection types.

## Collection Conversions

By convention a `___ToCol` `Iso` is used to convert to and from a `Collection`. The following built in collection `Isos` are provided
- `SliceToCol` for slices
- `MapToCol` for maps
- `StringToCol` for strings

For convenience the inverse functions are also provided.
- `ColToSlice` for slices
- `ColToMap` for maps
- `ColToString` for strings

## Operations

The built in collection operations can then be used to modify a `Collection`.

{{< playground file="/content/docs/14.collections/examples_test.go" id="collection_op1">}}

This example prints

{{< result file="/content/docs/14.collections/examples_test.go" id="result_collection_op1">}}

Note that the collection operation was applied and the result converted back to a `[]int`.

The `makelens` tool automatically converts slices and maps to `Collection` to ensure immutability.

{{< playground file="/content/docs/14.collections/examples_test.go" id="makelenscolop">}}

The following collection operations are provided.

- `EqCol` compares a focused collection with a given collection
- `DiffCol` diffs a focused collection with a given collection
- `FilteredCol` removes elements of a collection based on a `Predicate`.
- `AppendCol` adds elements to the end of a collection.
- `PrependCol` adds elements to the start of a collection.
- `SubCol` returns a sub collection with a given start positions and length.
- `ReversedCol` reverses the order of a collection.
- `OrderedCol` sorts the collection according to an `OrderByPredicate`

All operations are compatible with the `Modify` action. For example we can change the collection operation in the first example to `FilteredCol`

{{< playground file="/content/docs/14.collections/examples_test.go" id="collection_op2">}}

This example now prints.

{{< result file="/content/docs/14.collections/examples_test.go" id="result_collection_op2">}}

## Mismatched Collection Error types

The collection type has the following signature.

```go
type Collection[I,A,ERR] interface{...}
```

In addition to the expected index and value type parameters, `I` and `A`, `Collection` has a 3rd type parameter `ERR` to track whether the `Collection` may raise an error during iteration. The `Modify` action needs to iterate the `Collection` this means that the `ERR` of the `Collection` must match the `ERR` of the `Optic`. This causes issues when the error types of the [Optic] and [Collection] don't match.

{{< code file="/content/docs/14.collections/examples_test.go" id="colreconstrain1">}}

In this example the `FilteredCol` operations is an 
`Optic[Void, Collection[int, string, Err], Collection[int, string, Err], Collection[int, string, Err], Collection[int, string, Err], ReturnOne, ReadWrite, UniDir, Err]`

The source and focus collection error types are Err.

The built int collection conversions all return an error type of `Pure` and tying to us them together in a `Modify` action would cause a compile error. Go-optics provides collection error reconstraints to solve this compatibility issue.

{{< code file="/content/docs/14.collections/examples_test.go" id="colreconstrain2">}}

`SliceToCol` has is an `Optic[int,[]string,[]string,Collection[int,string.Pure],Collection[int,string.Pure],ReturnOne,ReadWrite,BiDir,Pure]`

We need some way to change the focused collection error type to `Err`. The `ColFocusErr` provides this function. It takes an optic of type

`Optic[J, S, T, Collection[I, A, OAERR], Collection[I, B, OBERR], RET, RW, DIR, ERR]`

Where `OAERR` and `OERR` and ERR are any type and converts the optic to.

`Optic[J, S, T, Collection[I, A, Err], Collection[I, B, Err], RET, RW, DIR, Err]`

Effectively making the  optic compatible with the `OrderedCol` optic that has a source collection error type of `Err` .

By normalising the source,focus or both to either Pure or Err. Collection optics can be made compatible with one another.

The following collection error type reconstraints are provided.

| Function             | Description                                                  |
| -------------------- | ------------------------------------------------------------ |
| `ColSourcePure`      | Convert the source collection error type to `Pure`           |
| `ColFocusPure`       | Convert the focus collection error type to `Pure`            |
| `ColSourceFocusPure` | Convert the source and focus collection error type to `Pure` |
| `ColSourceErr`       | Convert the source collection error type to `Err`            |
| `ColFocusErr`        | Convert the focus collection error type to `Err`             |
| `ColSourceFocusErr`  | Convert the source and focus collection error type to `Err`  |
|                      |                                                              |
There are 2 additional reconstraints that operate directly on an instance of a `Collection`
- `ColErr` which converts a `Collection[I,A,any]` to a `Collection[I,A,Err]`
- `ColPure`  which converts a `Collection[I,A,TPure]` to a `Collection[I,A,Pure]`

This example shows how to use `ColErr`

{{< code file="/content/docs/14.collections/examples_test.go" id="colreconstrain3">}}

Together these functions enable the collection error types of any `Collection` or collection operation `Optic` to be normalized to `Pure` or `Err` error types.