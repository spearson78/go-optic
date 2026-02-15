+++
title = "Metrics & Logging"
weight = 15
+++
# Metric & Logging

Go-optics provides the ability to collect metrics and to log the focused elements during the execution of an action.

## Metrics

Metrics can be collected for an optic by using the `WithMetrics` combinator. The signature of the optic is not altered but it will now report metrics to the given `Metrics` object. The following metrics are collected.

| Metric     | Meaning                                   |
| Focused    | how many elements were focused.           |
| Access     | how many times the source was accessed.   |
| LengthGet  | how many LengthGets were performed.       |
| IxGet      | how many index lookups were performed.    |

{{< playground file="/content/docs/15.metrics/examples_test.go" id="withmetrics">}}

When printing the metrics their names are abbreviated to their initial letter.

In this example we see that 4 elements were focused, The source was accessed once, to iterate over the slice, no index lookups and no length gets were performed.

`Ordered` also publishes a custom metric `"comparisons"`. This indicates how many comparisons were performed during the sort.

### Custom Metrics

Custom metrics can be published from your on own custom optics using the `IncCustomMetric`

{{< playground file="/content/docs/15.metrics/examples_test.go" id="custmetrics">}}

The custom optic increments the `len` metric for each string that it processes. The final metric is the total string length that was processed.

## Logging

Additional logging can be enabled for an individual optic using the `WithLogging` combinator. Like the `WithMetrics` combinator `WithLogging` doesn't change the signature of the optic.

{{< playground file="/content/docs/15.metrics/examples_test.go" id="withlogging">}}

Note: in the playground the logging will be written to rhe browser log. 

