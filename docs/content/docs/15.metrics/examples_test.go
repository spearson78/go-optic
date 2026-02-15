package main

import (
	"context"
	"fmt"
	"testing"

	. "github.com/spearson78/go-optic"
)

func TestWithMetrics(t *testing.T) {

	//BEGIN withmetrics
	data := []string{"alpha", "beta", "gamma", "delta"}

	sortStringSlice := Ordered(
		TraverseSlice[string](),
		OrderBy[string, string](
			Identity[string](),
		),
	)

	//Attach metrics to the sort
	var m Metrics
	sortStringSlice = WithMetrics(sortStringSlice, &m)

	res := MustGet(
		SliceOf(
			sortStringSlice,
			4,
		),
		data,
	)

	fmt.Println(res)
	fmt.Println(m)
	//END withmetrics

}

func TestIncCustomMetric(t *testing.T) {

	//BEGIN custmetrics
	data := []string{"alpha", "beta", "gamma", "delta"}

	//custom optic that publishes the len(string) as a custom metric.
	stringLen := CombiGetter[Pure, Void, string, string, int, int](
		func(ctx context.Context, source string) (Void, int, error) {
			strLen := len(source)
			IncCustomMetric(ctx, "len", strLen)
			return Void{}, strLen, nil
		},
		IxMatchVoid(),
		ExprCustom("custom string len"),
	)

	//Attach metrics to the stringLen
	var m Metrics
	stringLen = WithMetrics(stringLen, &m)

	res := MustGet(
		SliceOf(
			Compose(
				TraverseSlice[string](),
				stringLen,
			),
			4,
		),
		data,
	)

	fmt.Println(res)
	fmt.Println(m)
	//END custmetrics

	//Output:
	//Col[1:3 2:3 0:2] <nil>
}

func TestWithLogging(t *testing.T) {

	//BEGIN withlogging
	data := []string{"alpha", "beta", "gamma", "delta"}

	sortStringSlice := Ordered(
		TraverseSlice[string](),
		OrderBy[string, string](
			Identity[string](),
		),
	)

	//Attach logging to the sort
	sortStringSlice = WithLogging(sortStringSlice)

	res := MustGet(
		SliceOf(
			sortStringSlice,
			4,
		),
		data,
	)

	fmt.Println("Result:", res)
	//END withlogging

}
