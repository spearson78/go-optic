package optic

import (
	"context"

	"github.com/samber/lo"
)

type heap[T any] struct {
	data []T
	less func(ctx context.Context, a, b T) (bool, error)
}

func newHeap[T any](ctx context.Context, data []T, less func(ctx context.Context, a, b T) (bool, error)) (*heap[T], int, error) {

	heap := &heap[T]{
		data: data,
		less: less,
	}

	cmpCounter := 0

	for i := (len(heap.data) / 2) - 1; i >= 0; i -= 1 {
		cmpCount, _, err := heap.siftDown(ctx, i, len(heap.data))
		cmpCounter += cmpCount
		if err != nil {
			return nil, cmpCounter, err
		}
	}

	return heap, cmpCounter, nil

}

func (h *heap[T]) push(ctx context.Context, x T) (int, error) {
	h.data = append(h.data, x)
	return h.siftUp(ctx, len(h.data)-1)
}

func (h *heap[T]) pop(ctx context.Context) (T, int, error) {
	n := len(h.data) - 1
	h.data[0], h.data[n] = h.data[n], h.data[0]
	cmpCounter, _, err := h.siftDown(ctx, 0, n)
	if err != nil {
		var t T
		return t, cmpCounter, err
	}

	val := h.data[n]
	h.data = h.data[0:n]

	return val, cmpCounter, nil
}

func (h *heap[T]) peek() T {
	return h.data[0]
}

func (h *heap[T]) siftUp(ctx context.Context, j int) (int, error) {

	cmpCounter := 0

	for {
		if ctx.Err() != nil {
			return cmpCounter, ctx.Err()
		}

		i := (j - 1) / 2 // parent
		if i == j {
			break
		}

		cmpCounter++
		less, err := h.less(ctx, h.data[j], h.data[i])
		if err != nil {
			return cmpCounter, err
		}

		if !less {
			break
		}

		h.data[i], h.data[j] = h.data[j], h.data[i]
		j = i
	}

	return cmpCounter, ctx.Err()
}

func (h *heap[T]) siftDown(ctx context.Context, i0, n int) (int, bool, error) {

	cmpCounter := 0

	i := i0
	for {
		if ctx.Err() != nil {
			return cmpCounter, false, ctx.Err()
		}

		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n {
			cmpCounter++
			less, err := h.less(ctx, h.data[j2], h.data[j1])
			if err != nil {
				return cmpCounter, false, err
			}
			if less {
				j = j2 // = 2*i + 2  // right child
			}
		}

		cmpCounter++
		less, err := h.less(ctx, h.data[j], h.data[i])
		if err != nil {
			return cmpCounter, false, err
		}

		if !less {
			break
		}

		h.data[i], h.data[j] = h.data[j], h.data[i]

		i = j
	}
	return cmpCounter, i > i0, ctx.Err()
}

func heapSort[I, A any](ctx context.Context, left SeqIE[I, A], less func(ctx context.Context, a, b ValueI[I, A]) (bool, error)) (SeqIE[I, A], int, []lo.Tuple2[int, ValueI[I, A]], error) {

	const initialHeapCap = 1000

	heapData := make([]lo.Tuple2[int, ValueI[I, A]], 0, initialHeapCap)
	var retErr error
	i := 0

	left(func(val ValueIE[I, A]) bool {
		index, focus, err := val.Get()
		err = JoinCtxErr(ctx, err)
		if err != nil {
			retErr = err
			return false
		}
		heapData = append(heapData, lo.T2(i, ValI(index, focus)))
		i++
		return true
	})
	if retErr != nil {
		return nil, 0, nil, retErr
	}
	if len(heapData) == 0 {
		return func(yield func(val ValueIE[I, A]) bool) {}, 0, nil, nil
	}

	heap, cmpCount, err := newHeap[lo.Tuple2[int, ValueI[I, A]]](
		ctx,
		heapData,
		func(ctx context.Context, a, b lo.Tuple2[int, ValueI[I, A]]) (bool, error) {
			return less(ctx, a.B, b.B)
		},
	)
	IncCustomMetric(ctx, "comparisons", cmpCount)
	if err != nil {
		return nil, 0, nil, err
	}

	return func(yield func(val ValueIE[I, A]) bool) {
		size := len(heap.data)

		cont := true

		for size > 1 {

			val, cmpCount, err := heap.pop(ctx)
			IncCustomMetric(ctx, "comparisons", cmpCount)
			size--

			// 4. yield the sorted value
			cont = yield(ValIE(val.B.index, val.B.value, err))

			if !cont {
				break
			}
		}

		//5. yield the final sorted value
		if cont {
			val := heap.peek()
			yield(ValIE(val.B.index, val.B.value, ctx.Err()))
		}

	}, len(heap.data), heap.data, nil

}
