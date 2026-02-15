package optic_test

import (
	"context"
	"fmt"
	"iter"

	. "github.com/spearson78/go-optic"
)

func ExampleWithMetrics() {

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var metrics Metrics

	res, ok := MustGetFirst(
		Reduce(
			WithMetrics(
				TraverseSlice[int](),
				&metrics,
			),
			Sum[int](),
		),
		data,
	)

	fmt.Println(res, ok)
	fmt.Println(metrics)

	//Output:
	//55 true
	//metrics[F:10 A:1 I:0 L:0 Custom:map[]]
}

func ExampleIncCustomMetric() {

	traverseIntSlice := TraversalE[[]int, int](
		func(ctx context.Context, source []int) iter.Seq2[int, error] {
			return func(yield func(int, error) bool) {
				focusCount := 0
				for _, v := range source {
					focusCount++
					if !yield(v, nil) {
						break
					}
				}
				IncCustomMetric(ctx, "custom", focusCount)
			}
		},
		func(ctx context.Context, source []int) (int, error) {
			return len(source), nil
		},
		func(ctx context.Context, fmap func(focus int) (int, error), source []int) ([]int, error) {

			focusCount := 0

			ret := make([]int, 0, len(source))
			for _, oldVal := range source {
				focusCount++
				newVal, err := fmap(oldVal)
				if err != nil {
					IncCustomMetric(ctx, "custom", focusCount)
					return nil, err
				}
				ret = append(ret, newVal)
			}
			IncCustomMetric(ctx, "custom", focusCount)
			return ret, nil
		},
		ExprCustom("ExampleIncCustomMetric"),
	)

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var metrics Metrics

	result, err := Modify(WithMetrics(traverseIntSlice, &metrics), Mul(2), data)
	fmt.Println(result, err)
	fmt.Println(metrics)

	//Output:
	//[2 4 6 8 10 12 14 16 18 20] <nil>
	//metrics[F:10 A:1 I:0 L:0 Custom:map[custom:10]]
}
