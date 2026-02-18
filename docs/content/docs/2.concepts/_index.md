+++
title = "Concepts"
weight = 2
+++
# Concepts

There are 5 core concepts in go-optics
1. Actions
2. Optics
3. Compose
4. Combinators
5. The Identity Rule

Here is an example of performing an update using go-optics.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="intro" >}}

This example sets the value at slice index 3 to 5. The result value is ```[]int{10,20,30,4,50}```

```MustSet``` is the action `TraverseSlice[int]()` is an ```Optic``` and `Index(3)` is a combinator

Actions determine what operation will be performed on the data structure. There are 2 types of actions read and write. Write actions do not directly modify the given data, they return a modified copy.

The optic determines what data will be read or modified. Optics are described by their source,focus and index types. In this example ```TraverseSlice[int]()``` has
- a source type of `[]int`
- a focus type of `int`. The focus is the elements of the slice.
- an index of `int` which is the position of the focus in the slice.

`TraverseSlice[int]()` focuses on the ```int```values and int indexes within a source slice. 

We can visualize the optic in this way.

![Traverse](/concepts_1.svg)

On the left is the source ``[]int`` and on the right is the focus `int`, The index type is inside the square brackets `[int]`

Combinators modify the behaviour of optics. In this example ```Index(...,3)``` modifies the ```TraverseSlice[int]()``` optic to focus on the slice element at index 3. Combinators return optics enabling the results to be further refined by other combinators.

We can visualize the combinator like this.

![Index](/concepts_2.svg)

Again the source type is on the left. The empty slot in the middle is where the ``TraverseSlice]int()`` plugs in and on the right is the focus type.

![Index](/concepts_3.svg)

Here are some further example optic statements. 

{{< playground file="/content/docs/2.concepts/examples_test.go" id="mustgetfirst" >}}

This example is similar to the first one except we have changed the action from `MustSet` to `MustGetFirst` this actions reads the first focused value instead of updating it. in this case result is the integer value `40` and `ok` is true as we found a value at index 3.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="mustset" >}}

This example uses the  ```TraverseMap[string,int]()``` optic instead of ```TraverseSlice[int]()```. ```TraverseMap[string,int]()``` focuses on the values of a map instead of a slice. This optic is however still compatible with the same ```Index``` combinator used in the ```TraverseSlice[int]()``` examples. 
In this example the map entry with key "alpha" is set to the value 1. The result is a new map with the following value.

{{< result file="/content/docs/2.concepts/examples_test.go" id="mustset_result" >}}

Again the original map is not modified. ```TraverseSlice``` and ```TraverseMap``` immutably update maps and slices by returning new copies with the necessary updates applied.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="mustmodify" >}}


This example uses the `MustModify` action which applies an operation to each focused value. This is different to `MustSet` which sets a fixed value. In this example ```Mul(2)``` is the operation that will be applied to each focused value.
It may seem that ```Mul(2)``` is a new kind of concept but it is in fact an optic. It focuses on a value that is the source value multiplied by 2.
The optic that determines which values will be updated in this example is `Filtered( TraverseSlice[int](), Lt(10) )`. ```Filtered```  is a combinator that modifies an optic to only focus on values that match a predicate. The predicate in this example is `Lt(10)` meaning less than 10.
The `Filtered` optic will only focus on values within the the  ```TraverseSlice[int]()``` that are less than 10. Again predicates like ```Lt(10)``` may look like a new concept but predicates are also optics that focus on a boolean value for each source value. In this case ```Lt(10)``` will focus on true for values less than 10 and false for all other values. 
The ```Filtered``` combinator in this example raises an interesting question. What does it mean to update a filtered slice? The answer is that the ```Mul(2)``` operation is only applied to the values that match the filter. The other values are retained with their original values. The result in this example is the following slice.

{{< result file="/content/docs/2.concepts/examples_test.go" id="mustmodify_result" >}}

Notice that the 30 value is greater than 10 and was therefore not multiplied by 2. This behaviour is due to the Identity Rule. Which states that when applying the `Identity` operation to a Modify action the output should be identical to the input.In the case of filtered that means that non matching elements remain with their original values.

These examples only touch the surface of what optics are capable of but demonstrate the core concepts of actions, optics, combinators and the identity rule.

## Actions
Actions determine the operation that will be performed using an optic. Actions are provided according to the following naming convention.

|                 | Pure        | Error Aware | Context Aware  |
| --------------- | ----------- | ----------- | -------------- |
| **Non Indexed** | MustAction  | Action      | ActionContext  |
| **Indexed**     | MustActionI | ActionI     | ActionContextI |

The non-index/pure forms provide a simpler interface when the use case does not require the use of indexes and the optic is pure i.e. will not return an error.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="actions_mustmodify" >}}

The `TraverseMap[string,int]` focuses on the values within a map.  `MustModify` is an action that applies the given operation ,`Mul(2)`, to each focused element. In this example the result is.

{{< result file="/content/docs/2.concepts/examples_test.go" id="actions_mustmodify_result" >}}

`MustModifyI` is an indexed form of `MustModify` that provides access to the index during modification.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="actions_mustmodifyi" >}}


`OpI` is a specialized `Getter` constructor that operates on an index & value instead of the just the focused value. The result of this example is

{{< result file="/content/docs/2.concepts/examples_test.go" id="actions_mustmodifyi_result" >}}

`MustModifyI` can only modify the focused value. In order to modify an index you will need to use the `ReIndexed` combinator.

Indexes are also relevant to read actions.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="actions_mustgetfirst" >}}


In this example the `MustGetFirst` action is a non indexed action that will return the first focused value. The result in this example is

{{< result file="/content/docs/2.concepts/examples_test.go" id="actions_mustgetfirst_result" >}}

The true values indicates that a value was found.

``MustGetFirstI`` is the indexed form of ``MustGetFirst`` if we drop it into the example above.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="actions_mustgetfirsti" >}}


We see that the return now includes the index of the found element. In this example the result is

{{< result file="/content/docs/2.concepts/examples_test.go" id="actions_mustgetfirsti_result" >}}

So far all the actions we have used have been prefixed with `Must` These are non error aware and only accept pure Optics that will not return an error.

The error and context aware actions accept any optic and have an extra error return value.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="actions_erroraware" >}}

In this example we parse a string to int and multiply by 2. This can obviously fail if the string is not a number.

The results of this example is

{{< result file="/content/docs/2.concepts/examples_test.go" id="actions_erroraware_result" >}}

The context aware versions have an additional context parameter.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="actions_contextaware" >}}

Context aware actions are cancel & deadline aware and will stop processing and return an error if the context is canceled or the deadline has expired.

## Optics
There are 6 different kinds of optics in the go optics library.
1. `Getter`
2. `Lens`
3. `Iteration`
4. `Traversal`
5. `Iso`
6. `Prism`

Their basic behaviour is determined by combining 4 different properties.
1. Return type
2. Read only
3. Direction
4. Error raising

Return type determines whether an optic returns exactly 1 result or may return 0 or more results.
An example of a return 1 optic is a `Getter` e.g. ``Gt(1)`` this clearly always returns exactly 1 value either true or false. ``TraverseSlice[int]()`` will focus on the number of elements in the slice, which may of course be empty.

Read only determines whether an optic supports modifications, again `Getters` are read only. You cannot set a value into ``Eq(1)``. ``TraverseSlice[int]()`` however does support modifications by returning a modified version of the slice.

Direction determines whether the effect of an optic can be reversed, again `Getters` can´t be reversed. It's not possible to take a true value and pass it to ``Gt(1)`` and determine what the source value was. ``TraverseSlice[int)()`` is also not bidirectional as it focuses on multiple values and only single values can be reversed. ``Add`` is a bidirectional optic as te reverse an addition by subtracting.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="optics_intro_getter" >}}

Will add `5` to `10` the and return the integer `15` as a result.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="optics_intro_reversegetter" >}}

Will reverse add/subtract the integer value `5` from `10` and return `5`.

The following table demonstrates the behaviour of each of the optic types.

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

Read only bidirectional optics do not exist as the reverse operation is equivalent to modification.

Error awareness determines whether the Optic can return an error. The `Getter` `Gt(1)` is pure it never returns an error. `ParseInt` however may raise an error if passed a non numeric string.

The optics library includes compile time support to prevent optics being used with incompatible actions. e.g.

```go
//Compile error : ReadOnly does not satisfy comparable
Set( Eq(1) , false , 1 )
```
This clearly makes no sense to try and set `1 == 1` to false.

```go
//Compile error : ReturnMany does not satisfy comparable
Get( TraverseSlice[int]() , []int{ 1 , 2 , 3 } )
```
Traversing a slice focuses multiple values so we cannot view the single focused value. Get could try to return the first focused value but  if the slice is empty then there is no value to return.

```go
//Compile error : UniDir does not satisfy comparable
ReverseGet( Gt(1) , false )
```
`Gt` is unidirectional. We can't determine which integer to return for a false result from a greater than 1 operation.

```go
//Compile error : Err does not satisfy comparable
MustGet( ParseInt(10,32) , "1" )
```

`ParseInt` might return an error so it can't be used with Must actions.

When constructing optics several variants of the constructor are provided. 

| Variant | Capabilities                                        |
| ------- | --------------------------------------------------- |
| Base    | Non index aware, Non error raising, Non polymorphic |
| BaseI   | Index aware                                         |
| BaseE   | Error raising                                       |
| BaseP   | Polymorphic                                         |
| BaseIE  | Index aware, polymorphic.                           |
| BaseEP  | Error raising, polymorphic                          |
| BaseIEP | Index aware, error raising, polymorphic.            |

Index aware constructors require additional parameters and return values.
Non error raising constructors are required to create Pure optics. 
Polymorphic constructors require additional type parameters.

### Getters
`Getters` return exactly one result, cannot be written to and are unidirectional.  They are primarily used for computed values that cannot be reversed.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="optics_getter" >}}

The ``Op()`` family of constructors is provided to make it easy to wrap a go function into a `getter`

{{< playground file="/content/docs/2.concepts/examples_test.go" id="optics_getter_op" >}}

Existing functions can also be wrapped.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="optics_getter_op_existing" >}}

This example returns

{{< result file="/content/docs/2.concepts/examples_test.go" id="optics_getter_op_existing_result" >}}

### Lenses
Lenses return exactly one result, can be written to and are unidirectional. They are most often used to provide access the fields of a struct using the built in `FieldLens`

{{< playground file="/content/docs/2.concepts/examples_test.go" id="optics_lenses_fieldlens" >}}

Lenses to access the fields of a struct are usually auto-generated using the [makelens](/go-optic/docs/9.makelens) tool

{{< playground file="/content/docs/2.concepts/examples_data_test.go" id="optics_lens_makelens" >}}

The [makelens](/go-optic/docs/9.makelens) tool generated the `data.O.BlogPost().Content()` optic providing access to the `BlogPost.content` field.

### Isos
`Isos` return exactly one result and are bidirectional. Writing to an `Iso` is equivalent to the reverse operation. `Isos` are named after isomorphism as the conversion from source to focus and back should be lossless so the two types are isomorphic to one another.

In go-optics the unary mathematical operations are implemented as `Isos` This has the property of providing the inverse conversion automatically.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="optics_isos_celsius" >}}

Here we were able to `compose` a `Mul` with an `Add` and using `ReverseGet` we can perform the reverse conversion.

The `AsReverseGet` combinator is able to reverse the direction of an `Iso`

{{< playground file="/content/docs/2.concepts/examples_test.go" id="optics_isos_from" >}}

### Traversals & Iterations
Traversals return 0 or more results, can be written to and are unidirectional. Traversals are used to iterate over and immutably update the contents of a container.

In addition to the built in traversals user defined traversals can be constructed using the ``Traversal()`` family of constructors. This enables new container data types to be supported by go optics.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="optics_traversal_seqof" >}}

`SeqOf` is a `Combinator` that focuses an `iter.Seq` of the foci of an optic.
`TraverseSlice` is a `Traversal` that focuses the elements of a slice.

A `Traversal` can also be written to which will return a copy of the original data structure with the modification applied.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="optics_traversal_modify" >}}

`Mul(2)` is an `Iso` that multiplies the value by 2.

The result of the example  is

{{< result file="/content/docs/2.concepts/examples_test.go" id="optics_traversal_modify_result" >}}

A read only `Traversal` is called an `Iteration`. Traversals are also index aware. The ``Index`` combinator can be used to focus an element at a given index.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="optics_traversal_index" >}}

`TraverseMap]string,int]()` is a `Traversal` that focuses the elements of a `map[string}int`.

The result of the example is

{{< result file="/content/docs/2.concepts/examples_test.go" id="optics_traversal_index_result" >}}


### Prisms
Prisms return 0 or 1 results, can be written to and are bidirectional. Prisms are most often used to perform type safe conversions from a super type to a sub type. If the cast fails then 0 results are returned. If the cast succeeds then the cast value is returned.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="optics_prism_get" >}}

`SliceOf` is a combinator that creates a slice from the focused elements.
`TraverseSlice[any]()` is a `Traversal` that focuses the elements of a `[]any`
`DownCase[any,int]()` is a `Prism` that focuses an `int` if the case from any succeeds.

This example returns.

{{< result file="/content/docs/2.concepts/examples_test.go" id="optics_prism_get_result" >}}

Notice that the string `"two"` is missing from the result.

`Prisms` can also be written to.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="optics_prism_modify" >}}

`Mul(2)` is an `Iso` that multiples a value by 2.

The result of this example is

{{< result file="/content/docs/2.concepts/examples_test.go" id="optics_prism_modify_result" >}}

The integer values have been multiplied by 2. The string was ignored.

## Compose

Compose enables optics to focus deeply into nested structures by combining Optics together in a chain.

For a simple BlogPost data structure  

{{< code file="/content/docs/2.concepts/examples_test.go" id="optics_compose_types" >}}

We can define the following optics

{{< code file="/content/docs/2.concepts/examples_test.go" id="optics_compose_fieldlens" >}}


`FieldLenses` focus the given field of the parent struct.

The `blogComments` lens focuses the `BlogPost.comments` field

![Optic](/concepts_7.svg)

The `commentTitle` lens focuses the `Comment.title` field

![Optic](/concepts_8.svg)

We can use `Compose` to combine these lenses to focus the `Comments` nested within a `BlogPost`

{{< code file="/content/docs/2.concepts/examples_test.go" id="optics_compose_fieldlens_composed" >}}

In this example we have chained together 3 optics of differing types.
1. ``blogComments`` is a lens with source type ``BlogPost`` that focuses on the ``Comments`` field  which is of type ``[]Comment``
2. ``TraverseSlice[Comment]`` is a traversal with source type ``[]Comment`` that focuses on each  ``Comment`` in the slice.
3. ``commentTitle`` is a lens with source type ``Comment`` that focuses on the ``title`` field which is of type ``string``

We can visualise this chaining in this way.

![Blog comment's title](/concepts_4.svg)

Notice that the focus of each optic matches the source type of the next optic in the composed chain. Index types however do not need to match `TraverseSlice` has  an int index however `comment.Title` is un-indexed.

The resulting ``blogCommentTitles`` optic has a source type of ``BlogPost`` and a focus type of ``string`` where that string is the title field of a ``Comment``

This composed `Optic` is fully compatible with the built in `Actions`

{{< playground file="/content/docs/2.concepts/examples_test.go" id="optics_compose_get" typesid="optics_compose_types" playgroundid="optics_compose_fieldlens_playground" >}}

This will iterate over the title of each comment in a `BlogPost`

The composed optic is still write capable.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="optics_compose_modify" typesid="optics_compose_types" playgroundid="optics_compose_fieldlens_modify_playground" >}}


The effect of this ``MustModify`` action is to convert all comment titles to upper case.

In the visualisation the connecting arrows indicate a flow to the right but for modifiable optics the arrows flow in both directions. We can visualise the ``MustModify`` in this way.

![Blog comment's title](/concepts_5.svg)

We feed a blogPost into the left hand source. The ``[]Comment`` is extracted and fed to ``TraverseSlice`` which feeds each ``Comment`` to ``commentTitle`` which in turn extracts the comment title. This is then passed to ``ToUpper`` which returns an upper case string. This result is then fed back through the whole flow in reverse order.
``commentTitle`` packs the upper case title back into a new copy of the ``Comment`` struct with the updated title. This is then in turn packed into a new ``[]Comment`` slice by the ``TraverseSlice`` and is then in turn packed into a new copy of the original blogPost by ``blogComments``. In this way the optic has focused on the exact field we wanted to modify and then rebuilt a new copy of the original structure on the return path without modifying the original. This is how go optics achieves immutable data type support.

Note that the index value only flows to the right. Index values cannot be updated during an update of a `Traversal`

For this use case of delving into a nested structure this compose syntax is verbose. It is advised to use the [makelens](/go-optic/docs/9.makelens) tool to generate helpers that make this use case much more intuitive. The above example can be simplified to the following when using [makelens](/go-optic/docs/9.makelens)

{{< playground file="/content/docs/2.concepts/examples_data_test.go" id="optics_compose_makelens" >}}

behind the scenes the ``O.BlogPost().Comments().Traverse().Title()`` is performing all the necessary composition.

When using ``Compose`` the type of the return optic is determined by the types of the input optics. 

|               | Getter    | Lens      | Iteration | Traversal | Iso       | Prism     |
| ------------- | --------- | --------- | --------- | --------- | --------- | --------- |
| **Getter**    | Getter    | Getter    | Iteration | Iteration | Getter    | Iteration |
| **Lens**      | Getter    | Lens      | Iteration | Traversal | Lens      | Traversal |
| **Iteration** | Iteration | Iteration | Iteration | Iteration | Iteration | Iteration  |
| **Traversal** | Iteration | Traversal | Iteration | Traversal | Traversal | Traversal |
| **Iso**       | Getter    | Lens      | Iteration | Traversal | Iso       | Prism     |
| **Prism**     | Iteration | Traversal | Iteration | Traversal | Prism     | Prism     |

The important thing to note from this table is that it is complete. Every optic type can be composed with every other optic type and a correctly functioning optic will be returned.

`Compose` by default retains the index of the right most optic. This makes sense when composing with a `Traversal` as traversals typically are indexed.

However in this example.

{{< playground file="/content/docs/2.concepts/examples_test.go" typesid="optics_compose_types" id="optics_compose_right" >}}


After the `TraverseSlice` we composed a `FieldLens` which does not have an index. 


![Optic](/concepts_9.svg)
We can however maintain this index in the resulting `Optic` by using `ComposeLeft`

{{< playground file="/content/docs/2.concepts/examples_test.go" typesid="optics_compose_types" id="optics_compose_left" >}}

![Optic](/concepts_10.svg)

The `ComposeLeft` has retained the `[int]` index of the `TraverseSlice`. As a convenience the `Optics` generated by [makelens](/go-optic/docs/9.makelens) use `ComposeLeft` when the right side of the `Compose` is un-indexed.

{{< playground file="/content/docs/2.concepts/examples_data_test.go" id="optics_compose_left_makelens" >}}

## Combinators
Combinators are functions that take an optic as a parameter and return a new optic with modified behaviour. The optic returned by a combinator may not be the same type (Lens,Traversal,..) as the input optic.

Go-optics provides a wide variety of other combinators that modify the behavior of another optic.

We have already used several combinators
- `Compose` is in fact a combinator that combines multiple `Optics` together.
- `Index` returns an Optic that focuses the elements with a given index.
- `Filtered` returns an Optic that focuses the elements that match a given predicate.
- `AsReverseGet` returns an OPtic that reverses the direction of an `Iso`
- `SeqOf` returns an Optic that focuses an `iter.Seq` of the focused elements.
- `SliceOf` returns an Optic that focuses a slice of the focused elements.

An important `combinator` is `Filtered`

{{< playground file="/content/docs/2.concepts/examples_test.go" id="optics_combinators" >}}

- `Filtered` is a combinator that focuses only elements that meet the `Predicate`
- `TraverseSlice[int]()` is a `Traversal` that focuses the elements of a `[]int`

`Filtered` uses the types from `TraverseSlice[int]()` as it's own source and focus types, the index is retained as well. The int focus type of `TraverseSlice[int]` defines the type of the `Predicate`.

In this case the Predicate is `AndOp` which is a combinator that applies a boolean and to the focus of it's 2 input optics.
 - `Gt(10)` is a `Predicate` that returns true if the value is greater than 10.
 - `Lt(40)` is a `Predicate` that returns true if the value is less than than 40.

The modify operation can also be a combined optic. In this case we composed `Add(10)` with `Mul(2)` creating a combined effect of `(focus+10)*2`

The result of this example is

{{< result file="/content/docs/2.concepts/examples_test.go" id="optics_combinators_result" >}}

Every value >10 and <40 had the operation `(focus+10)*2` applied.

One of the most important concepts in go-optics is that predicates and modify operations are also optics and can be combined together using combinators.

## The Identity Rule
When reading from an Optic the set of focused elements will be returned. If filtering is involved then the filtered results will be missing from the read result. If re-ordering has occurred then the results will be delivered in the order defined by the optic.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="identity_rule" >}}

When using ``MustModify`` the values focused by the filter are updated and the non matching values are passed through unmodified into the result.

{{< result file="/content/docs/2.concepts/examples_test.go" id="identity_rule_result" >}}

The identity rule states that passing ``Identity()`` as the operation to a modify action should return the original data structure back . ``Identity()`` is an operation that returns the source directly back as the focus.

If we change `Mul(2)` to `Identity()`

{{< playground file="/content/docs/2.concepts/examples_test.go" id="identity_rule_identity" >}}

Then due to the Identity Rule the result must be

{{< result file="/content/docs/2.concepts/examples_test.go" id="identity_rule_identity_result" >}}

This can only be true if the filtered combinator retains the unmatched values in the output.
This applies to all built in combinators provided by go optics. 

The same logic applies to re-ordering

{{< playground file="/content/docs/2.concepts/examples_test.go" id="identity_rule_identity_reorder" >}}

Reading the re-ordered data returns a sorted slice.

{{< result file="/content/docs/2.concepts/examples_test.go" id="identity_rule_identity_reorder_result" >}}

However under modification the data is returned in the original order.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="identity_rule_identity_reorder_modify" >}}

{{< result file="/content/docs/2.concepts/examples_test.go" id="identity_rule_identity_reorder_modify_result" >}}

This may seem to make `Ordered` pointless. However the elements were focused in order. In conjunction with other combinators the effect of this ordered focus can be used effectively.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="identity_rule_identity_reorder_modify_taking" >}}

`Taking(...,3)` focuses on the first `3` elements. In this case this is the first 3 ordered elements. Which are then multiplied by 2 and then restored to their original positions in the slice.

{{< result file="/content/docs/2.concepts/examples_test.go" id="identity_rule_identity_reorder_modify_taking_result" >}}

The net effect is that the lowest 3 elements are modified in place.

Another impact of the Identity Rule is that under modification arithmetic operations are applied the modification is made and the the arithmetic operations are reversed.

{{< playground file="/content/docs/2.concepts/examples_test.go" id="identity_rule_context" >}}

The result of this example is

{{< result file="/content/docs/2.concepts/examples_test.go" id="identity_rule_context_result" >}}

The `Add(1.0)` was an increase of 1 Fahrenheit which corresponds to an ~0.5 increase in Celsius. Note that we only specified the conversion to Fahrenheit. `Mul` and `Add` are `Isos` and automatically provide the `ReverseGetter` to reverse the conversion.

This behaviour is a result of the identity rule as if the `Add(1)` were replaced with identity the result must be the original Celsius value.

This may seem counter-intuitive initially but this can be used to create "virtual" fields or conversions. Suppose our underlying data is stored in celsius we can provide an `Optic` that performs the conversion to fahrenheit.

{{< code file="/content/docs/2.concepts/examples_test.go" id="identity_rule_virtual" >}}

Users are now free to work in both fahrenheit or celsius by selecting the correct `Optic`