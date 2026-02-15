+++
title = "Custom Optics"
weight = 10
+++
# Custom Optics

## Optic Types

When implementing a custom optic it is important to first identify the correct type of Optic.
The following table lists the available optics and their properties.

| Return Type | Read Type  | Direction Type | Optic Type      |
| ----------- | ---------- | -------------- | --------------- |
| Return One  | Read Only  | Uni            | Getter          |
| Return One  | Read Only  | Bi             | --------------- |
| Return One  | Read Write | Uni            | Lens            |
| Return One  | Read Write | Bi             | Iso             |
| Return Many | Read Only  | Uni            | Iteration       |
| Return Many | Read Only  | Bi             | --------------- |
| Return Many | Read Write | Uni            | Traversal       |
| Return Many | Read Write | Bi             | Prism           |


To determine the correct return type you need to identify how many elements your new optic will focus. `ReturnOne` means that the optic always focuses exactly 1 element or returns an error. `ReturnMany` means the optic can focus 0 or more elements.

To determine the read type you need to identify if your optic will support updating the source. If updates are possible then a `ReadWrite` optic should be selected.

To determine the direction type you need to determine if the effect of reading your optic can be reversed without loss. If the operation is reversible a `BiDir` optic should be selected.

It is important to identify the correct optic type as this determines the functions you will need to implement.

## Optic Type Variants

Custom optics are constructed using a method named after the optic type.

{{< playground file="/content/docs/10.custom/examples_test.go" id="getter1" >}}

The `Getter` optic requires 2 parameters the get function and an `OpticExpression`. `OpticExpressions` are covered in another section of the documentation. In these examples we will use the `ExprCustom` helper to create a custom expression.  See [Optic Expressions](/docs/12.expression) for more information.

There are several variants of the `Getter` constructor available.

| Variant Postfix | Capabilities                                        |
| --------------- | --------------------------------------------------- |
|                 | Non index aware, Non error raising, Non polymorphic |
| I               | Index aware                                         |
| E               | Error raising                                       |
| IE              | Index aware, Error raising.                         |
Index aware getters return an index value in addition to the value.

{{< playground file="/content/docs/10.custom/examples_test.go" id="getterI" >}}

In this example we create an index aware getter for the last element of a slice. That returns the index of the last element. In addition to the get function we now need to provide an `ixMatch` function that is able to compare 2 index values for equality. This is needed to support the `Index` combinator that focuses on an element in an optic with a specific index. The `GetterI` constructor provides an implementation for this using the `ixMatch` function. Index aware optics are also compatible with the non-index aware actions so we must use `MustGetI` to retrieve the index. If we were to use `MustGet` instead we would retrieve only the value of the last element of the slice.

The non error returning constructors provide a safe way to create pure optics that cannot return an error.

Error raising getters receive a context and are able to return an error. 

{{< playground file="/content/docs/10.custom/examples_test.go" id="getterE" >}}

There is an additional variant postfix `P` for polymorphic these are covered in the [polymorphic optics](#polymorphic-optics) section.
### Getters
`Getters` return exactly one result, cannot be written to and are unidirectional.  They are can be used to provide read-only access to a field or for a computed value.

We already saw some examples of creating a `Getter` in the previous section.

As `getters` are most often used as predicates or as the operation to apply to a focused value a shorthand ``Op()`` family of constructors is provided to avoid the need to specify the optic expression (`ExprCustom`)

{{< playground file="/content/docs/10.custom/examples_test.go" id="getter_op" >}}

### Lenses
Lenses return exactly one result, can be written to and are unidirectional. They are most often used to provide access the fields of a struct.
 
Lenses are constructed from a pair of `getter` and `setter` functions.

{{< playground file="/content/docs/10.custom/examples_test.go" typesid="lens_types" id="lens_customfieldlens" >}}

Lenses to access the fields of a struct are usually auto-generated using the [makelens](/docs/9.makelens) tool which uses the built in `FieldLens` optic.

### Isos
`Isos` return exactly one result and are bidirectional. Writing to an `Iso` is equivalent to the reverse operation. `Isos` are named after isomorphism as the conversion from source to focus and back should be lossless so the two types are isomorphic to one another. They are most often used to convert datatypes.

Isos are constructed from a pair of `getter` and `reverse getter` functions.

{{< playground file="/content/docs/10.custom/examples_test.go" id="custom_iso" >}}

### Iterations
Iterations return 0 or more results are read only and are unidirectional. Iterations are most often used to iterate over the contents of a container.

User defined iterations can be constructed using the ``Iteration()`` family of constructors. Iteration in go optics is based on the range func.

{{< playground file="/content/docs/10.custom/examples_test.go" id="custom_iteration" >}}

The `lengthGetter` function may be nil in which case go-optics will provide a default implementation.

The ``MustGetFirst`` action returns the first focused value. In this example the result is the int value 1. Iterations in go optics support early exit by returning false from the yield function.  In this example the `MustGetFirst` action will cause yield to return false and exit the loop after 1 iteration.

### Traversals
Traversals return 0 or more results, can be written to and are unidirectional. Traversals are used to iterate over and immutably update the contents of a container.

Custom traversals can be constructed using the ``Traversal()`` family of constructors. This enables new container data types to be supported by go optics.

{{< playground file="/content/docs/10.custom/examples_test.go" id="custom_traversal" >}}

It is imperative that the modify function does not mutate the source. This is required as the purpose of go-optics is to support immutability.
The modify function needs to call fmap on each focused element. The fmap function is provided by go-optics and represents the operation that the user is performing. In this example the fmap will be a function that multiples by 2 (``Mul(2)``)

The `lengthGetter` function may be nil in which case go-optics will provide a default implementation.

When constructing traversals it is advisable to use the index aware version of the constructor. This enables optimized index lookups.

{{< playground file="/content/docs/10.custom/examples_test.go" id="custom_traversali" >}}


The ``Index`` combinator will use the index getter function to efficiently access the indexed element.

The `Ix Match` function needs to compare 2 index values. This is required in the case that an efficient index lookup cannot be performed and index values need to be compared individually.

### Prisms
Prisms return 0 or 1 results, can be written to and are bidirectional. Prisms are most often used to perform type safe conversions from a super type to a sub type. If the cast fails then 0 results are returned. If the cast succeeds then the cast value is returned.
Custom Prisms can be constructed using the ``Prism()`` family of constructors.

{{< playground file="/content/docs/10.custom/examples_test.go" id="custom_prism" >}}


A prism is constructed from 2 functions, a match function and an embed function. The match function checks if the source can be converted to the focus type and returns the result of the conversion and a bool to indicate if the match was successful. The embed function reverses the cast and embeds the focus in the super type source.

In this example if the io.Writer passed to the ``MustOver`` action is a ``bytes.Buffer`` then it will be grown by 100 bytes.

## Polymorphic Optics

A polymorphic optic is able the S and T types and A and B types are completely independent.

Polymorphic optics are created using `P` variants of the base optic constructor.
Here is the complete list of all possible constructor post fixes.

| Variant Postfix | Capabilities                                        |
| --------------- | --------------------------------------------------- |
|                 | Non index aware, Non error raising, Non polymorphic |
| I               | Index aware                                         |
| E               | Error raising                                       |
| P               | Polymorphic                                         |
| IE              | Index aware, Error raising.                         |
| IP              | Index aware, Polymorphic                            |
| EP              | Error raising, Polymorphic                          |
| IEP             | Index aware, Error raising, Polymorphic.            |

To demonstrate will extend our `parseInt` `Getter` to be a polymorphic `Iso` using the `IsoEP` constructor.

{{< playground file="/content/docs/10.custom/examples_test.go" id="isoep" >}}

In this example we specify the S,T,A and B types individually.  We construct an optic with this structure.


![Polymorphic parse int](/custom_isoep.svg)
The return source type for the `Iso` is `int64` instead of the expected `string` for a non polymorphic optic.
This has the effect of enabling optics to transform to different data types under modification. We could use this for example to convert a slice of `string` to a slice of `int`
For maximum re-usability you should strive to implement polymorphic optics where possible. 