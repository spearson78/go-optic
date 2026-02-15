package oio

import (
	"context"
	"errors"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

const defaultBufSize = 8192

type LinePosition struct {
	line int
	col  int
}

func DecodeFile(enc encoding.Encoding) Optic[Void, Collection[int64, Extent, Err], Collection[int64, Extent, Err], Collection[LinePosition, rune, Err], Collection[LinePosition, rune, Err], ReturnOne, ReadWrite, BiDir, Err] {
	return CombiIso[ReadWrite, BiDir, Err, Collection[int64, Extent, Err], Collection[int64, Extent, Err], Collection[LinePosition, rune, Err], Collection[LinePosition, rune, Err]](
		func(ctx context.Context, source Collection[int64, Extent, Err]) (Collection[LinePosition, rune, Err], error) {
			return ColIE[Err, LinePosition, rune](
				func(ctx context.Context) SeqIE[LinePosition, rune] {
					line := 0
					col := 0
					return func(yield func(ValueIE[LinePosition, rune]) bool) {

						d := enc.NewDecoder()
						dstBuf := make([]byte, defaultBufSize)
						srcBuf := make([]byte, 0, defaultBufSize)
						source.AsIter()(ctx)(func(val ValueIE[int64, Extent]) bool {
							_, focus, err := val.Get()
							if err != nil {
								return yield(ValIE(LinePosition{line: line, col: col}, rune(0), err))
							}

							newData, err := focus.Read()
							if err != nil {
								return yield(ValIE(LinePosition{line: line, col: col}, rune(0), err))
							}

							for _, dataByte := range newData {
								srcBuf = append(srcBuf, dataByte)

								nDst, nSrc, err := d.Transform(dstBuf, srcBuf, false)
								if nDst > 0 {
									for _, r := range []rune(string(dstBuf[0:nDst])) {
										if !yield(ValIE(LinePosition{line: line, col: col}, r, nil)) {
											return false
										}

										if r == '\n' {
											col = 0
											line++
										} else {
											col++
										}
									}
								}

								if nSrc > 0 {
									srcBuf = srcBuf[0:0]
								}

								if err == nil || errors.Is(err, transform.ErrShortSrc) {
									//Process next byte
									continue
								} else {
									return yield(ValIE(LinePosition{line: line, col: col}, rune(0), err))
								}
							}

							//Read more data
							return true
						})
					}
				},
				nil,
				func(a, b LinePosition) bool {
					return a.line == b.line && a.col == b.col
				},
				nil,
			), nil
		},
		func(ctx context.Context, focus Collection[LinePosition, rune, Err]) (Collection[int64, Extent, Err], error) {
			return ColIE[Err, int64, Extent](
				func(ctx context.Context) SeqIE[int64, Extent] {
					i := int64(0)
					e := enc.NewEncoder()
					srcBuf := make([]byte, 0, defaultBufSize/4)
					dstBuf := make([]byte, defaultBufSize)
					return func(yield func(ValueIE[int64, Extent]) bool) {
						focus.AsIter()(ctx)(func(val ValueIE[LinePosition, rune]) bool {
							_, focus, err := val.Get()
							if err != nil {
								return yield(ValIE(i, Extent(nil), err))
							}

							if len(srcBuf) >= defaultBufSize-2 {
								nDst, nSrc, err := e.Transform(dstBuf, srcBuf, false)

								if nDst > 0 {

									ext, err := Get(AsReverseGet(ExtentData()), dstBuf[0:nDst])
									if !yield(ValIE(i, ext, err)) {
										return false
									}
									i++
								}

								if nSrc > 0 {
									copy(srcBuf, srcBuf[0:nSrc])
									srcBuf = srcBuf[0:nSrc]
								}

								if err == nil {
									//Process next rune
									return true
								} else {
									return yield(ValIE(i, Extent(nil), err))
								}

							} else {
								srcBuf = append(srcBuf, []byte(string(focus))...)
								return true
							}
						})

						if len(srcBuf) > 0 {

							nDst, _, err := e.Transform(dstBuf, srcBuf, false)

							if nDst > 0 {

								ext, err := Get(AsReverseGet(ExtentData()), dstBuf[0:nDst])
								if !yield(ValIE(int64(0), ext, err)) {
									return
								}
							}

							if err == nil {
								return
							} else {
								yield(ValIE(int64(0), Extent(nil), err))
								return
							}

						}
					}
				},
				nil,
				IxMatchComparable[int64](),
				nil,
			), nil
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.DecodeFileExpr{
				OpticTypeExpr: ot,
				Encoding:      enc,
			}
		}),
	)
}
