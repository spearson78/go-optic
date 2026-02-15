+++
title = "Indexed Optics"
weight = 11
+++

Indexed optics focus an index value in addition to the value. This typically occurs when traversing a collection.

{{< playground file="/content/docs/11.indexed/examples_test.go" id="traversemap" >}}

Here we have used `SeqIOf` which is the index aware version of `SeqOf`. `SeqIOf` focuses an `iter.Seq2` of the focused index and value.

The result is this sequence of values being printed.

```go
alpha 1
beta 2
delta 4
gamma 3
```
Traversing a map in go optics always iterates the value in a consistent order. The sorting is performed by sorting by the index values.

There is only 1 active index at any time. `Compose` by default maintains the right optic´s index. This can cause the index to be lost.

{{< playground file="/content/docs/11.indexed/examples_test.go" id="lostindex" >}}

`Mul(10)` doesn't have an index and this causes the `TraverseMap`'s index to be lost. The result is this sequence of values being printed.
```go
{} 10
{} 20
{} 40
{} 30
```
In go-optics an empty struct called `Void` is used to represent the lack of an index.

There are 3 built in variants of `Compose` that can retain either the left or both indices.
- `ComposeLeft` retains the left index.
- `ComposeBoth` retains both indices wrapped in a `lo.Tuple2`
- `ComposeI` uses a user provided `IxMapper` to map the index.

## Compose Left
`ComposeLeft` is a drop in replacement for `Compose` that retains the left instead of right index.

{{< playground file="/content/docs/11.indexed/examples_test.go" id="composeleft"  playgroundid="play_composeleft" >}}

This example now focuses the `string` index from the map instead of the `Void` index from the `Mul(10)` operation. The example prints the following values. 

```go
alpha 10
beta 20
delta 40
gamma 30
```
The index focused by `ComposeLeft` is compatible with the `Index` combinator which focuses a value at a given index.

{{< playground file="/content/docs/11.indexed/examples_test.go" id="ixget_composeleft" playgroundid="play_ixget_composeleft" >}}

This prints the following values.

{{< result file="/content/docs/11.indexed/examples_test.go" id="res_ixget_composeleft" >}}

`ComposedLeft` provides `Index`  with an optimised optic that is still able to perform a `data["beta"}` lookup into the map.

## Compose Both
`ComposeBoth` is also a drop in replacement for `Compose` that retains both indices as a `lo.Tuple2`

{{< playground file="/content/docs/11.indexed/examples_test.go" id="composeboth" >}}
In this example we can see that the index is a combination of the slice index and map key.
```go
{0 alpha} 1
{0 beta} 2
{1 delta} 4
{1 gamma} 3
```
`ComposeBoth` also provides an optimized `Index` implementation. We simply pass the 2 index values as a `lo.Tuple2`

{{< playground file="/content/docs/11.indexed/examples_test.go" id="ixget_composeboth" playgroundid="play_ixget_composeboth">}}

This example returns the "gamma" value from the second slice. The `Index` was able to perform both the slice index and map lookup `data[1]["gamma"]`

{{< result file="/content/docs/11.indexed/examples_test.go" id="res_ixget_composeboth" >}}

## ComposeI

`ComposeI` enables custom mapping of the index. The mapping is performed by an `IxMap` optic. By using the `IxMapIso` constructor the composed optic is able to provide an optimized `Index` implementation as it is able to recover the original index values by performing a `ReverseGet` on the `IxMap`

{{< playground file="/content/docs/11.indexed/examples_test.go" id="composei" >}}

This example is equivalent to a `ComposeBoth`

```go
{0 alpha} 1
{0 beta} 2
{1 delta} 4
{1 gamma} 3
```

## Reindexing

It's sometimes necessary to modify the index of an optic for example when building a map. For this purposes the reindexing combinators can be used.

{{< playground file="/content/docs/11.indexed/examples_test.go" id="reindexed1" >}}

This example uses the `ReIndexed` combinator to modify the `int` index of a slice traversal and converts it into a `string` and builds a `map` from the result. The `EqT2` parameter is an `IxMatch` optic. It is used to perform equality checks on the new index during `Index` lookups using the new index. After reindexing the `Index` lookup is performed by filtering on the index using the `IxMatch`

There are several other reindexing combinators.

| Combinator   | Purpose                                    |
|--------------|--------------------------------------------|
| `Indexed`      | replaces the index with an integer counter |
| `SelfIndex`    | uses the focused element as the index      |
| `ReIndexed`    | enables an arbitrary transformation if the index |

By combining the `SelfIndex` and `ReIndexed` combinators an arbitrary property of the focus can be used as the index.

{{< playground file="/content/docs/11.indexed/examples_test.go" id="reindexed2" >}}

This example prints

```go
map[Erika Mustermann:{Erika Mustermann 37} Max Mustermann:{Max Mustermann 42}] <nil>
```

By using the reindexed combinators arbitrary transformations can be applied to the index of an optic.