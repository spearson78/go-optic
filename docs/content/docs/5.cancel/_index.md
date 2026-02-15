+++
title = "Deadlines & Cancellation"
weight = 5
+++
# Deadlines & Cancellation
All actions in go optics support context deadlines and cancellation. The context is checked periodically and on deadline expiry or cancellation will return an error.
This is useful when processing potentially unlimited data structures to prevent denial of service.

{{< playground file="/content/docs/5.cancel/examples_test.go" id="cancel" >}}

Here we have attempt to create a slice from an infinite sequence. Once the deadline has expired the yield will return false and the deadline error will be returned.   

If you need access to the context in your user defined optic or operation in order to honour the deadline then be sure to use an error aware constructor for the optic you are using. These end with a suffix of `E` 
1. ``OpE``
2. ``OperationE``
3. ``LensE``
4. ``IterationE``
5. ``TraversalE``
6. ``IsoE``
7. ``PrismE``