package oio

import (
	"bytes"
	"context"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

func SplitFile(splitOn byte) Optic[Void, Collection[int64, Extent, Err], Collection[int64, Extent, Err], Collection[int64, Extent, Err], Collection[int64, Extent, Err], ReturnOne, ReadWrite, UniDir, Err] {

	splitIter := func(ctx context.Context, source Collection[int64, Extent, Err]) (Void, Collection[int64, Extent, Err], error) {
		return Void{}, ColIE[Err](
			func(ctx context.Context) SeqIE[int64, Extent] {
				return func(yield func(ValueIE[int64, Extent]) bool) {

					buf := make([]Extent, 0, 2)
					i := int64(0)

					cont := true

					source.AsIter()(ctx)(
						func(val ValueIE[int64, Extent]) bool {
							_, focus, err := val.Get()
							if err != nil {
								return yield(ValIE(i, Extent(nil), err))
							}

							data, err := focus.Read()

							for {

								splitIndex := bytes.IndexByte(data, splitOn)
								if splitIndex == -1 {
									focus.Acquire()
									buf = append(buf, focus)
									return true
								} else {

									start, _ := focus.SliceFromStart(splitIndex)
									if start != nil {
										buf = append(buf, start)
									}
									ext := newMultiExtent(buf)
									cont = yield(ValIE(i, ext, nil))
									err = ext.Release()
									if !cont {
										return false
									}

									if err != nil {
										cont = yield(ValIE(i, Extent(nil), err))
										if !cont {
											return false
										}
									}

									i++
									buf = nil

									if len(data) > splitIndex {
										data = data[splitIndex+1:]
										nextFocus, _ := focus.SliceToEnd(splitIndex + 1)
										if nextFocus == nil {
											return true
										}
										focus = nextFocus
									} else {
										return true
									}
								}
							}
						},
					)

					if cont && len(buf) > 0 {
						ext := newMultiExtent(buf)
						cont = yield(ValIE(i, ext, nil))
						err := ext.Release()
						if cont && err != nil {
							yield(ValIE(i, Extent(nil), err))
						}
					}
				}
			},
			nil,
			IxMatchComparable[int64](),
			nil,
		), nil
	}

	return GetModIEP[Void, Collection[int64, Extent, Err], Collection[int64, Extent, Err], Collection[int64, Extent, Err], Collection[int64, Extent, Err]](
		splitIter,
		func(ctx context.Context, fmap func(index Void, focus Collection[int64, Extent, Err]) (Collection[int64, Extent, Err], error), source Collection[int64, Extent, Err]) (Collection[int64, Extent, Err], error) {

			_, origSeq, err := splitIter(ctx, source)
			if err != nil {
				return nil, err
			}

			newSeq, err := fmap(Void{}, origSeq)
			err = JoinCtxErr(ctx, err)

			return ColIE[Err](
				func(ctx context.Context) SeqIE[int64, Extent] {
					return func(yield func(ValueIE[int64, Extent]) bool) {
						cont := true
						i := int64(0)

						newSeq.AsIter()(ctx)(func(val ValueIE[int64, Extent]) bool {
							_, focus, err := val.Get()
							if i != 0 {
								cont = yield(ValIE(i, Extent(bytesExtent([]byte{splitOn})), err))
								i++
								if !cont {
									return false
								}
							}
							cont = yield(ValIE(i, focus, err))
							i++
							return cont
						})
					}
				},
				nil,
				IxMatchComparable[int64](),
				nil,
			), err

		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.SplitFileExpr{
				OpticTypeExpr: ot,
				SplitOn:       splitOn,
			}
		}),
	)

}
