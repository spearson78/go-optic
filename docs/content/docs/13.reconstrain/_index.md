+++
title = "Reconstrain"
weight = 13
+++
# Reconstrain

When composing optics their constraints are merged in a `CompositionTree`. 

{{< code file="/content/docs/13.reconstrain/examples_test.go" id="compositiontree" >}}

For normal usage this isn't normally a problem. However If I were to return this optic from a function then I would have to use the full type including the composition tress.

{{< code file="/content/docs/13.reconstrain/examples_test.go" id="celsiusToFahrenheit1" >}}

## Collapsing Reconstraints

The return type constraint is `CompositionTree[ReturnOne, ReturnOne]` we could collapse this back down to `ReturnOne` without affecting the type of the optic. The `Ret1` reconstrain does exactly this. It takes an optic with an arbitrarily complex `ReturnOne` composition path and returns an Optic with `ReturnOne` as it's return type.

{{< code file="/content/docs/13.reconstrain/examples_test.go" id="celsiusToFahrenheit2" >}}

go-optics provides an equivalent reconstrain for each constraint type.

{{< code file="/content/docs/13.reconstrain/examples_test.go" id="celsiusToFahrenheit3" >}}

- `Rw` reconstrains to `ReadWrite`
- `Bd` reconstrains to `BiDir`
- `EPure` reconstrains to `Pure`

There are also equivalent reconstrains for the "opposing" constraints.

- `RetM` reconstarins to `ReturnMany`
- `Ro` reconstrains to `ReadOnly`
- `Ud` reconstrains to `UniDir`
- `EErr` reconstrains to `Err`

This simple reconstrains suffice for a composed optic returned from a function. However when building custom combinators that receive an optic as a parameter. Additional reconstraints are needed.

## Left / Right Reconstraints

{{< code file="/content/docs/13.reconstrain/examples_test.go" id="bytesOf1" >}}

The left side of the composition tree is determined by the optic passed to the combinator.

Here we can use the L (left) and R (right) family of reconstraints.

- `RetL` reconstrains to the left return type as long as the right side is `ReturnOne`.
- `RetR` reconstrains to the right return type as long as the left side is `ReturnOne`.
- `RwL` reconstrains to the left read write type as long as the right side is `ReadWrite`.
- `RwR` reconstrains to the right read write type as long as the left side is `ReadWrite`.
- `DirL` reconstrains to the left direction type as long as the right side is `BiDir`.
- `DirR` reconstrains to the right direction type as long as the left side is `BiDir`.
- `EErrL` reconstrains to the left error type as long as the right side is `Pure`.
- `EErrR` reconstrains to the right error type as long as the left side is `Pure`.

In this case we can use the left hand variants.

{{< code file="/content/docs/13.reconstrain/examples_test.go" id="bytesOf2" >}}

In the case that either the left or right side is not `ReturnOne`,`ReadWrite`,`BiDir` or `Pure` then the simpler `RetM`,`Ro`,`Ud`,`EErr` can be used to reconstrain.

## Merge Reconstraints

When writing a combinator the left and right constraints may both be the same undefined type.

{{< code file="/content/docs/13.reconstrain/examples_test.go" id="makeTuple1" >}}

Here we can see that the constraints on the left and right are identical. We can use the merge family of reconstraints to simplify them.

- RetMerge
- RwMerge
- DirMerge
- EErrMerge

We can additionally use the right family to reconstrains away the left side of the nested composition tree.

{{< code file="/content/docs/13.reconstrain/examples_test.go" id="makeTuple2" >}}

## Advanced Reconstraints

The above reconstraints suffice for most cases. It's typically a good idea to reconstrain as close to composition as possible to avoid the construction of large complex composition trees. This is however sometimes not possible. For these cases go-optics offers 2 solutions.

- Tree Reconstraints
- UnsafeReconstrain

### Tree Reconstraints

Tree reconstraints enable nodes in the composition tree to be switched with one another. The composition tree can then be manipulated until one of the other reconstraints can be used to simplify the node.

Here is an overview of the swap & trans reconstraints.

![Swap & Trans overview](/comptree_ops.svg)

- `Swap`/`SwapL`/`SwapR` swaps the 2 labeled composition tree nodes
- `Merge`/`MergeL`/`MergeR` merges the 2 labeled composition tree nodes
- `Trans`/`TransL`/`TransR` transposes the 2 labeled composition tree nodes across a tree boundary.

`SwapL` and `SwapR` can be used to move a node into position that one of the Tran operations can transpose it to somewhere else in the tree.

`TransL` and `TransR` can be used to move nodes up and down in the tree hierarchy.

`MergeL` and `MergeR` are convenience functions to avoid having to `TransL`/`TransR` nodes up 1 level in the tree just to merge them.

For complex composition trees these functions can become a little unwieldy. They do however provide compile time safety for the constraint types.

Consider this extension of the `lo.Tuple2` example expanded to a `lo.Tuple9`

{{< code file="/content/docs/13.reconstrain/examples_test.go" id="makeTuple9" >}}

The composition trees are deeper than 2 levels so we will have to use `TransL` and `TransR` to lift the nodes into reach of the other reconstrains.

This sequence of calls is able to fully merge the compositions tree.

{{< code file="/content/docs/13.reconstrain/examples_test.go" id="makeTuple9_reconstrain" >}}


### UnsafeReconstrain

`UnsafeReconstrain` will apply the supplied constraints verbatim onto the optic bypassing the typical compile time checks. This can lead to an optic with undefined behavior, including raising  panic, when used. With due care and attention the code is much easier to implement and to read.

{{< code file="/content/docs/13.reconstrain/examples_test.go" id="makeTuple9_unsafereconstrain" >}}

`UnsafeReconstrain` can be made safe by implementing a custom wrapper with the concrete optic type that you are reconstraining.

{{< code file="/content/docs/13.reconstrain/examples_test.go" id="makeTuple9_safereconstrain" >}}

If anything changes in the optics constraints the compile will fail and you can reanalyze the constraints to verify they are still valid. This show the provided reconstraints are implemented.















