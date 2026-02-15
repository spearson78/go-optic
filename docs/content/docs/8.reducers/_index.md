+++
title = "Reducers"
weight = 8
+++
# Reducers
To support generating aggregated values ,like the sum,mean or max of a set of values, go-optics includes the `Reducer` interface. This enables the `Reduce` combinator to focus the aggregate value.

A Reducer is built from 3 components.
 1. An empty value to begin the reduction from e.g. for Sum this would be 0.
 2. A reduce function that takes the current state and a value and returns the new state. e.g. for Mean this might be the current sum and the number of reduced elements.
 3. An end function that takes the final state and returns the final value. e.g. For Mean this would calculate the final Mean value.
The `Reduce` combinator can then be used to apply this `Reducer` to the foci of an optic.

{{< playground file="/content/docs/8.reducers/examples_test.go" id="reducers_sum" >}}

If the optic focuses no values then ok will be false.
In this example result is the value `10`.

`Reducers` can be created from compatible `ReturnOne` optics using the `AsReducer` constructor.

{{< playground file="/content/docs/8.reducers/examples_test.go" id="reducers_addt2" playgroundid="reducers_playground_addt2" >}}

`0` is the empty value.
`AddT2` is a built in operation that takes a tuple of 2 ints and returns their sum. This is used as the reduce function.

The following Reducers are provided

| Reducer      | Purpose                                                                                |
| ------------ | -------------------------------------------------------------------------------------- |
| Sum          | Sum of all the foci                                                                    |
| Prod         | Multiplies all the foci                                                                |
| MaxReducer   | Returns the maximum of all the foci.                                                   |
| MaxReducerI  | Returns the maximum of all the foci and it's index.                                    |
| MinReducer   | Returns the minimum of all the foci.                                                   |
| MinReducerI  | Returns the minimum of all the foci and it's index.                                    |
| FirstReducer | Returns the first focused value.                                                       |
| Mean         | Returns the mean of all the foci.                                                      |
| Median       | Returns the median of all the foci.                                                    |
| Mode         | Returns the modal value of all the foci.                                               |
| ReducerT2    | Enables 2 reductions to be run at one time. Returning a tuple of the 2 reduced values. |
| ReducerT3    | Enables 3 reductions to be run at one time. Returning a tuple of the 3 reduced values. |