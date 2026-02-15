package main

import (
	"context"

	"github.com/spearson78/go-optic"
)

func traverseTNIter[A any](ctx context.Context, elements ...*A) optic.SeqIE[int, A] {
	return func(yield func(optic.ValueIE[int, A]) bool) {
		for i, v := range elements {
			if !yield(optic.ValIE(i, *v, ctx.Err())) {
				break
			}
		}
	}
}

func traverseTNModify[A, B any](ctx context.Context, fmap func(index int, focus A) (B, error), elements ...*A) ([]B, error) {
	ret := make([]B, 0, len(elements))

	for i, v := range elements {
		b, err := fmap(i, *v)
		err = optic.JoinCtxErr(ctx, err)
		if err != nil {
			return nil, err
		}
		ret = append(ret, b)
	}

	return ret, nil
}

func traverseTNIxGet[A any](index int, elements ...*A) optic.SeqIE[int, A] {
	if index < 0 || index > len(elements)-1 {
		return func(yield func(optic.ValueIE[int, A]) bool) {}
	}
	val := *elements[index]
	return func(yield func(optic.ValueIE[int, A]) bool) {
		yield(optic.ValIE(index, val, nil))
	}
}
