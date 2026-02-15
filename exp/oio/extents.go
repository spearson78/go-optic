package oio

import (
	"context"
	"errors"
	"hash/crc32"
	"io"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

//TODO: This package is extremely complicated. I'd like to simplify it.
//I'm not sure the idea to base io on a Collection of Extents was a good idea.

type Extent interface {
	Read() ([]byte, error)
	WriteTo(w io.Writer) (n int64, err error)
	Length() int

	SliceFromStart(untilIndex int) (Extent, int)
	SliceToEnd(fromIndex int) (Extent, int)

	Acquire()
	Release() error
}

func ExtentData() Optic[Void, Extent, Extent, []byte, []byte, ReturnOne, ReadWrite, BiDir, Err] {
	return CombiIsoMod[ReadWrite, BiDir, Err, Void, Extent, Extent, []byte, []byte](
		func(ctx context.Context, source Extent) (Void, []byte, error) {
			data, err := source.Read()
			return Void{}, data, JoinCtxErr(ctx, err)
		},
		func(ctx context.Context, fmap func(index Void, focus []byte) ([]byte, error), source Extent) (Extent, error) {

			data, err := source.Read()
			err = JoinCtxErr(ctx, err)
			if err != nil {
				return nil, err
			}

			dataCrc := crc32.ChecksumIEEE(data)

			newData, err := fmap(Void{}, data)
			err = JoinCtxErr(ctx, err)

			if len(data) != len(newData) || dataCrc != crc32.ChecksumIEEE(newData) {
				return bytesExtent(newData), nil
			} else {
				return source, nil
			}
		},
		func(ctx context.Context, focus []byte) (Extent, error) {
			return bytesExtent(focus), nil
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.ExtentData{
				OpticTypeExpr: ot,
			}
		}),
	)
}

func Bytes() Optic[Void, Collection[int64, Extent, Err], Collection[int64, Extent, Err], []byte, []byte, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoE[Collection[int64, Extent, Err], []byte](
		func(ctx context.Context, source Collection[int64, Extent, Err]) ([]byte, error) {

			var ret []byte
			var retErr error

			source.AsIter()(ctx)(func(val ValueIE[int64, Extent]) bool {
				_, focus, err := val.Get()
				if err != nil {
					retErr = err
					return false
				}

				data, err := focus.Read()
				err = JoinCtxErr(ctx, err)
				if err != nil {
					retErr = err
					return false
				}

				ret = append(ret, data...)
				return true
			})

			return ret, retErr
		},
		func(ctx context.Context, focus []byte) (Collection[int64, Extent, Err], error) {
			return ColIE[Err](
				func(ctx context.Context) SeqIE[int64, Extent] {
					return func(yield func(ValueIE[int64, Extent]) bool) {
						yield(ValIE(int64(0), Extent(bytesExtent(focus)), nil))
					}
				},
				nil,
				IxMatchComparable[int64](),
				func(ctx context.Context) (int, error) {
					return 1, nil
				},
			), nil
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Bytes{
				OpticTypeExpr: ot,
			}
		}),
	)
}

type readerAtContext struct {
	refCount int
	closeFnc func() error
	r        io.ReaderAt
}

type readerAtExtent struct {
	ctx *readerAtContext
	pos int64
	len int
}

func (r *readerAtExtent) Read() ([]byte, error) {
	data := make([]byte, r.len)
	_, err := r.ctx.r.ReadAt(data, r.pos)
	return data, err
}

func (r *readerAtExtent) WriteTo(w io.Writer) (int64, error) {
	data := make([]byte, r.len)
	n, err := r.ctx.r.ReadAt(data, r.pos)

	if n > 0 {
		w.Write(data)
		if errors.Is(err, io.EOF) {
			return int64(n), nil
		}
	}
	return int64(n), err
}

func (r *readerAtExtent) SliceFromStart(untilIndex int) (Extent, int) {
	if untilIndex >= r.len {
		return nil, r.len
	}

	r.ctx.refCount = r.ctx.refCount + 1

	return &readerAtExtent{
		ctx: r.ctx,
		pos: r.pos,
		len: untilIndex,
	}, untilIndex
}

func (r *readerAtExtent) SliceToEnd(fromIndex int) (Extent, int) {

	if fromIndex >= r.len {
		return nil, r.len
	}

	r.ctx.refCount = r.ctx.refCount + 1

	return &readerAtExtent{
		ctx: r.ctx,
		pos: r.pos + int64(fromIndex),
		len: r.len - fromIndex,
	}, fromIndex
}

func (r *readerAtExtent) Acquire() {
	r.ctx.refCount = r.ctx.refCount + 1
}

func (r *readerAtExtent) Release() error {
	r.ctx.refCount = r.ctx.refCount - 1
	if r.ctx.refCount == 0 {
		return r.ctx.closeFnc()
	}
	return nil
}

func (r *readerAtExtent) Length() int {
	return r.len
}

type bytesExtent []byte

func (r bytesExtent) Read() ([]byte, error) {
	return []byte(r), nil
}

func (r bytesExtent) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(r)
	return int64(n), err
}

func (r bytesExtent) SliceFromStart(untilIndex int) (Extent, int) {
	if untilIndex >= len(r) {
		return nil, len(r)
	}
	return bytesExtent(r[:untilIndex]), 0
}

func (r bytesExtent) SliceToEnd(fromIndex int) (Extent, int) {
	if fromIndex >= len(r) {
		return nil, len(r)
	}
	return bytesExtent(r[fromIndex:]), 0
}

func (r bytesExtent) Acquire() {
}

func (r bytesExtent) Release() error {
	return nil
}

func (r bytesExtent) Length() int {
	return len(r)
}

type multiExtent []Extent

func newMultiExtent(e []Extent) Extent {

	if len(e) == 1 {
		return e[0]
	} else {
		return multiExtent(e)
	}

}

func (r multiExtent) WriteTo(w io.Writer) (int64, error) {

	totalWritten := int64(0)

	for _, v := range r {
		n, err := v.WriteTo(w)
		totalWritten += int64(n)
		if err != nil {
			return totalWritten, err
		}
	}
	return totalWritten, nil
}

func (r multiExtent) Read() ([]byte, error) {

	ret := make([]byte, 0, r.Length())

	for _, v := range r {
		d, err := v.Read()
		if err != nil {
			return nil, err
		}
		ret = append(ret, d...)
	}

	return ret, nil
}

func (r multiExtent) Release() error {
	var err error
	for _, v := range r {
		errors.Join(err, v.Release())
	}
	return err
}

func (r multiExtent) Acquire() {
	for _, v := range r {
		v.Acquire()
	}
}

func (r multiExtent) SliceFromStart(untilIndex int) (Extent, int) {

	var extents []Extent

	totalConsumed := 0

	for _, v := range r {

		e, consumed := v.SliceFromStart(untilIndex)
		totalConsumed += consumed
		untilIndex -= consumed
		if e == nil {
			extents = append(extents, v)
		} else {
			extents = append(extents, e)
		}

		if untilIndex == 0 {
			break
		}
	}

	return multiExtent(extents), totalConsumed

}

func (r multiExtent) SliceToEnd(fromIndex int) (Extent, int) {

	var extents []Extent

	totalConsumed := 0

	for _, v := range r {

		if fromIndex == 0 {
			extents = append(extents, v)
		} else {
			e, consumed := v.SliceToEnd(fromIndex)
			totalConsumed += consumed
			fromIndex -= consumed
			if e != nil {
				extents = append(extents, e)
			}
		}

	}

	return multiExtent(extents), totalConsumed

}

func (r multiExtent) Length() int {
	l := 0
	for _, v := range r {
		l += v.Length()
	}
	return l
}

type mmapContext struct {
	refCount int
	closeFnc func() error
}

type mmapExtent struct {
	ctx  *mmapContext
	data []byte
}

func (r mmapExtent) Read() ([]byte, error) {
	return r.data, nil
}

func (r mmapExtent) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(r.data)
	return int64(n), err
}

func (r mmapExtent) SliceFromStart(untilIndex int) (Extent, int) {
	if untilIndex >= len(r.data) {
		return nil, len(r.data)
	}

	r.ctx.refCount++

	return mmapExtent{
		ctx:  r.ctx,
		data: r.data[:untilIndex],
	}, 0
}

func (r mmapExtent) SliceToEnd(fromIndex int) (Extent, int) {
	if fromIndex >= len(r.data) {
		return nil, len(r.data)
	}

	r.ctx.refCount++

	return mmapExtent{
		ctx:  r.ctx,
		data: r.data[fromIndex:],
	}, 0
}

func (r mmapExtent) Acquire() {
	r.ctx.refCount++
}

func (r mmapExtent) Release() error {
	r.ctx.refCount--
	if r.ctx.refCount == 0 {
		return r.ctx.closeFnc()
	}
	return nil
}

func (r mmapExtent) Length() int {
	return len(r.data)
}
