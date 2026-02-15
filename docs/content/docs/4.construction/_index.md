+++
title = "Object Construction"
weight = 4
+++
# Object Construction
Although the primary purpose of go-optics is immutable updates to nested data structures it also provides powerful object graph query abilities. These are used to focus in the correct elements to update but equally these focused elements can be returned directly using the `Get` family of actions. However `Get` can return only a single value. In order to return multiple values they need to be collected together into a container.

{{< playground file="/content/docs/4.construction/examples_test.go" id="construct_sliceof" >}}

This example focuses the values in the slice with a value >= 30
`SliceOf` collects those values into a new slice.

The values focused by an Optic also maintain their index this enables maps to also be constructed.

{{< playground file="/content/docs/4.construction/examples_test.go" id="construct_mapof" >}}

The index value of a traversed slice is the original slice index.This example produces the following map.

{{< result file="/content/docs/4.construction/examples_test.go" id="construct_mapof_result" >}}

Object construction in go-optics also supports modification. In this case the updated values in the updated collection are used to repopulate the original data structure. This is done by position. The index values are ignored during updates.

{{< playground file="/content/docs/4.construction/examples_test.go" id="construct_sliceof_modify" >}}

The slice passed to the modification `Op` is flattened values from the nested slices.
```go
[]string{"alpha", "beta", "gamma", "delta"}
```

The modification `Op` changes the element at index 2.

The values in the return from the op are woven into the original structure. In this case the `result` is

{{< result file="/content/docs/4.construction/examples_test.go" id="construct_sliceof_modify_result" >}}

The element at index 2 was "gamma" and that is the value that now contains the updated value. This demonstrates the power of Optics to update complex nested structures with ease.

If the modified collection contains fewer values then the original values are used for the missing values. If the modified collection contains additional elements these elements are ignored.

The following object construction optics are available.
- `MapOf`
- `ListOf`
- `SliceOf`
- `StringOf`
- `T2Of`
- `ColOf`
- `SeqOf`

Polymorphic versions are available with the suffix P. However the polymorphic versions will return an error if the modified collection has fewer elements than the original.

