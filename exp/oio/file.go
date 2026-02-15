package oio

import (
	"bufio"
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/edsrzf/mmap-go"
	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

type FileInfo interface {
	fs.FileInfo
	FullPath() string
	modified() *modifiedFileInfo
}

type fileInfo struct {
	fs.FileInfo
	fullPath string
}

func (f *fileInfo) FullPath() string {
	return f.fullPath
}

func (f *fileInfo) modified() *modifiedFileInfo {
	return nil
}

type modifiedFileInfo struct {
	origInfo     fs.FileInfo
	origFullPath string

	name    string      // base name of the file
	mode    fs.FileMode // file mode bits
	modTime time.Time   // modification time
}

func (f *modifiedFileInfo) FullPath() string {
	return f.origFullPath
}

func (f *modifiedFileInfo) Name() string {
	return f.name
}

func (f *modifiedFileInfo) Size() int64 {
	return f.origInfo.Size()
}

func (f *modifiedFileInfo) Mode() fs.FileMode {
	return f.mode
}

func (f *modifiedFileInfo) ModTime() time.Time {
	return f.modTime
}

func (f *modifiedFileInfo) IsDir() bool {
	return f.origInfo.IsDir()
}

func (f *modifiedFileInfo) Sys() any {
	return f.origInfo.Sys()
}

func (f *modifiedFileInfo) modified() *modifiedFileInfo {
	return f
}

func Stat() Optic[Void, string, string, FileInfo, FileInfo, ReturnOne, ReadOnly, UniDir, Err] {
	return CombiGetter[Err, Void, string, string, FileInfo, FileInfo](
		func(ctx context.Context, source string) (Void, FileInfo, error) {
			source = filepath.Clean(source)
			fi, err := os.Stat(source)
			return Void{}, &fileInfo{
				FileInfo: fi,
				fullPath: source,
			}, err
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.StatExpr{
				OpticTypeExpr: ot,
			}
		}),
	)
}

func TraverseFileInfo() Optic[string, FileInfo, FileInfo, FileInfo, FileInfo, ReturnMany, ReadOnly, UniDir, Err] {

	return IterationIE[string, FileInfo, FileInfo](
		func(ctx context.Context, source FileInfo) SeqIE[string, FileInfo] {
			return func(yield func(ValueIE[string, FileInfo]) bool) {

				first := true

				walkPath := source.FullPath()

				err := filepath.Walk(walkPath, func(childPath string, d fs.FileInfo, err error) error {
					if err != nil {
						yield(ValIE(childPath, FileInfo(nil), err))
						return err
					}

					if first {
						first = false
						return nil
					}

					child := &fileInfo{
						FileInfo: d,
						fullPath: childPath,
					}

					if !yield(ValIE(childPath[len(walkPath)+1:], FileInfo(child), err)) {
						return filepath.SkipAll
					}

					if d.IsDir() {
						return filepath.SkipDir
					}

					return nil
				})

				if err != nil {
					yield(ValIE(source.FullPath(), FileInfo(nil), err))
				}

			}
		},
		nil,
		func(ctx context.Context, index string, source FileInfo) SeqIE[string, FileInfo] {
			return func(yield func(ValueIE[string, FileInfo]) bool) {
				fullPath := path.Join(source.FullPath(), index)
				d, err := os.Stat(fullPath)
				if err != nil {
					if errors.Is(err, os.ErrNotExist) {
						//Nop
						return
					}

					yield(ValIE(index, FileInfo(nil), err))
				}

				child := &fileInfo{
					FileInfo: d,
					fullPath: fullPath,
				}

				yield(ValIE(index, FileInfo(child), nil))
			}
		},
		IxMatchComparable[string](),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.TraverseDirExpr{
				OpticTypeExpr: ot,
			}
		}),
	)
}

func FileInfoFullPath() Optic[Void, FileInfo, FileInfo, string, string, ReturnOne, ReadOnly, UniDir, Pure] {
	return MethodGetter(FileInfo.FullPath)
}

func FileInfoName() Optic[Void, FileInfo, FileInfo, string, string, ReturnOne, ReadOnly, UniDir, Pure] {
	return MethodGetter[FileInfo, string](
		FileInfo.Name,
	)
}

func FileInfoSize() Optic[Void, FileInfo, FileInfo, int64, int64, ReturnOne, ReadOnly, UniDir, Pure] {
	return MethodGetter(FileInfo.Size)
}

func FileInfoMode() Optic[Void, FileInfo, FileInfo, os.FileMode, os.FileMode, ReturnOne, ReadOnly, UniDir, Pure] {
	return MethodGetter[FileInfo, os.FileMode](
		FileInfo.Mode,
	)
}

func FileInfoModTime() Optic[Void, FileInfo, FileInfo, time.Time, time.Time, ReturnOne, ReadOnly, UniDir, Pure] {
	return MethodGetter[FileInfo, time.Time](
		FileInfo.ModTime,
	)
}

func FileInfoIsDir() Optic[Void, FileInfo, FileInfo, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return MethodGetter(FileInfo.IsDir)
}

func FileInfoSys() Optic[Void, FileInfo, FileInfo, any, any, ReturnOne, ReadOnly, UniDir, Pure] {
	return MethodGetter(FileInfo.Sys)
}

func FileBytes(perm os.FileMode, bufSize int, mode ReadMode) Optic[Void, string, string, []byte, []byte, ReturnOne, ReadWrite, UniDir, Err] {
	return Ret1(Rw(Ud(EErr(Compose(TraverseFileOverwrite(perm, bufSize, mode), Bytes())))))
}

type extentReader struct {
	err  error
	buf  []byte
	next func() (int64, Extent, error, bool)
	stop func()
}

func (e *extentReader) Close() error {
	e.stop()
	return nil
}

func (e *extentReader) Read(p []byte) (n int, err error) {

	if e.err != nil {
		return 0, e.err
	}

	if len(e.buf) == 0 {
		_, ext, err, ok := e.next()
		if !ok {
			e.stop()
			e.err = io.EOF
			return 0, io.EOF
		}
		if err != nil {
			e.stop()
			e.err = err
			return 0, err
		}

		data, err := ext.Read()

		if err != nil {
			e.stop()
			e.err = err
			return 0, err
		}

		e.buf = data
	}

	if len(p) >= len(e.buf) {
		copy(p, e.buf)
		l := len(e.buf)
		e.buf = nil
		return l, nil
	} else {
		copy(p, e.buf[:len(p)])
		e.buf = e.buf[len(p):]
		return len(p), nil
	}
}

// Warning the returned reader MUST be fully read to prevent leaks.
func Reader(bufSize int) Optic[Void, Collection[int64, Extent, Err], Collection[int64, Extent, Err], io.Reader, io.Reader, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoE[Collection[int64, Extent, Err], io.Reader](
		func(ctx context.Context, source Collection[int64, Extent, Err]) (io.Reader, error) {
			next, stop := PullIE(source.AsIter()(ctx))
			return &extentReader{
				next: next,
				stop: stop,
			}, nil
		},
		func(ctx context.Context, focus io.Reader) (Collection[int64, Extent, Err], error) {
			return ColIE[Err](
				func(ctx context.Context) SeqIE[int64, Extent] {
					return func(yield func(ValueIE[int64, Extent]) bool) {
						buf := make([]byte, bufSize)
						i := int64(0)

						eof := false
						for !eof {
							n, err := io.ReadFull(focus, buf)
							if errors.Is(err, io.ErrUnexpectedEOF) {
								eof = true
								err = nil
							}
							if err != nil {
								if !yield(ValIE(i, Extent(nil), err)) {
									break
								}
							} else {
								if n > 0 {
									if !yield(ValIE(i, Extent(bytesExtent(buf[:n])), nil)) {
										break
									}
								}
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
			return expr.Reader{
				OpticTypeExpr: ot,
			}
		}),
	)
}

func ReadCloser(bufSize int) Optic[Void, Collection[int64, Extent, Err], Collection[int64, Extent, Err], io.ReadCloser, io.ReadCloser, ReturnOne, ReadWrite, BiDir, Err] {
	return Ret1(Rw(Bd(EErr(Compose(Reader(bufSize), IsoCast[io.Reader, io.ReadCloser]())))))
}

func TraverseFileOverwrite(perm os.FileMode, bufSize int, mode ReadMode) Optic[Void, string, string, Collection[int64, Extent, Err], Collection[int64, Extent, Err], ReturnOne, ReadWrite, UniDir, Err] {
	return Ret1(Rw(Ud(EErr(Compose(OverWrite(), TraverseFile(perm, bufSize, mode))))))
}

func OverWrite() Optic[Void, string, string, FileNames, FileNames, ReturnOne, ReadWrite, BiDir, Pure] {
	return Iso[string, FileNames](
		func(source string) FileNames {
			return FileNames{
				source,
				source,
			}
		},
		func(focus FileNames) string {
			return focus.OutFile
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Overwrite{
				OpticTypeExpr: ot,
			}
		}),
	)
}

type FileNames struct {
	InFile  string
	OutFile string
}

type ReadMode int

const (
	ReadBuffer ReadMode = iota
	ReadAt
	ReadAll
	ReadMMap
)

func getReadFile(bufSize int, mode ReadMode) func(ctx context.Context, source FileNames) (Void, Collection[int64, Extent, Err], error) {

	switch mode {
	case ReadAll:
		return func(ctx context.Context, source FileNames) (Void, Collection[int64, Extent, Err], error) {

			return Void{}, ColIE[Err](
				func(ctx context.Context) SeqIE[int64, Extent] {
					return func(yield func(ValueIE[int64, Extent]) bool) {
						f, err := os.Open(source.InFile)
						if err != nil {
							yield(ValIE(int64(0), Extent(nil), err))
							return
						}

						data, err := io.ReadAll(f)
						f.Close()

						yield(ValIE(int64(0), Extent(bytesExtent(data)), err))
					}
				},
				nil,
				IxMatchComparable[int64](),
				func(ctx context.Context) (int, error) {
					return 1, nil
				},
			), nil
		}
	case ReadMMap:
		return func(ctx context.Context, source FileNames) (Void, Collection[int64, Extent, Err], error) {

			return Void{}, ColIE[Err](
				func(ctx context.Context) SeqIE[int64, Extent] {
					return func(yield func(ValueIE[int64, Extent]) bool) {

						map_file, err := os.Open(source.InFile)
						if err != nil {
							yield(ValIE(int64(0), Extent(nil), err))
							return
						}

						mapped, err := mmap.Map(map_file, mmap.RDONLY, 0)
						if err != nil {
							map_file.Close()
							yield(ValIE(int64(0), Extent(nil), err))
							return
						}

						ctx := &mmapContext{
							refCount: 1,
							closeFnc: func() error {
								return errors.Join(mapped.Unmap(), map_file.Close())
							},
						}

						yield(ValIE(int64(0), Extent(&mmapExtent{
							ctx:  ctx,
							data: mapped,
						}), nil))
					}
				},
				nil,
				IxMatchComparable[int64](),
				func(ctx context.Context) (int, error) {
					return 1, nil
				},
			), nil
		}
	case ReadBuffer:

		return func(ctx context.Context, source FileNames) (Void, Collection[int64, Extent, Err], error) {

			return Void{}, ColIE[Err](
				func(ctx context.Context) SeqIE[int64, Extent] {
					return func(yield func(ValueIE[int64, Extent]) bool) {
						f, err := os.Open(source.InFile)
						if err != nil {
							yield(ValIE(int64(0), Extent(nil), err))
							return
						}
						defer f.Close()

						eof := false
						for !eof {
							b := make([]byte, bufSize)
							n, err := io.ReadFull(f, b)
							b = b[:n]
							if err != nil {
								if errors.Is(err, io.ErrUnexpectedEOF) {
									eof = true
									err = nil
								}
								yield(ValIE(int64(0), Extent(bytesExtent(b)), err))
								return

							}
							if !yield(ValIE(int64(0), Extent(bytesExtent(b)), err)) {
								break
							}
						}
					}
				},
				nil,
				IxMatchComparable[int64](),
				func(ctx context.Context) (int, error) {
					fileInfo, err := os.Stat(source.InFile)
					if err != nil {
						return 0, err
					}
					//Ceiling for the final extent.
					return int((fileInfo.Size() + int64(bufSize) - int64(1)) / int64(bufSize)), nil
				},
			), nil
		}

	case ReadAt:
		return func(ctx context.Context, source FileNames) (Void, Collection[int64, Extent, Err], error) {

			return Void{}, ColIE[Err](
				func(ctx context.Context) SeqIE[int64, Extent] {
					return func(yield func(ValueIE[int64, Extent]) bool) {
						file, err := os.Open(source.InFile)
						err = JoinCtxErr(ctx, err)
						if err != nil {
							yield(ValIE(int64(0), Extent(nil), err))
							return
						}

						info, err := file.Stat()
						err = JoinCtxErr(ctx, err)
						if err != nil {
							yield(ValIE(int64(0), Extent(nil), err))
							return
						}

						size := info.Size()
						extCount := size / int64(bufSize)
						extRemain := int(size % int64(bufSize))

						cont := true

						extContext := readerAtContext{
							refCount: 1,
							closeFnc: file.Close,
							r:        file,
						}

						for i := int64(0); i < extCount; i++ {
							ext := &readerAtExtent{
								ctx: &extContext,
								pos: i * int64(bufSize),
								len: bufSize,
							}

							ext.Acquire()
							cont = yield(ValIE(i, Extent(ext), nil))
							err = ext.Release()
							if !cont {
								break
							}

							if err != nil {
								cont = yield(ValIE(i, Extent(nil), err))
							}
							if !cont {
								break
							}
						}

						if cont && extRemain != 0 {
							ext := &readerAtExtent{
								ctx: &extContext,
								pos: extCount * int64(bufSize),
								len: extRemain,
							}
							ext.Acquire()
							cont = yield(ValIE(extCount, Extent(ext), nil))
							err = ext.Release()
							if cont && err != nil {
								yield(ValIE(extCount, Extent(nil), err))
							}
						}

						extContext.refCount = extContext.refCount - 1
						if extContext.refCount == 0 {
							extContext.closeFnc()
						}
					}
				},
				nil,
				IxMatchComparable[int64](),
				func(ctx context.Context) (int, error) {
					fileInfo, err := os.Stat(source.InFile)
					if err != nil {
						return 0, err
					}
					//Ceiling for the final extent.
					return int((fileInfo.Size() + int64(bufSize) - int64(1)) / int64(bufSize)), nil
				},
			), nil
		}
	default:
		panic("TraverseFile unknown read mode")
	}

}

func TraverseFile(perm os.FileMode, bufSize int, mode ReadMode) Optic[Void, FileNames, FileNames, Collection[int64, Extent, Err], Collection[int64, Extent, Err], ReturnOne, ReadWrite, UniDir, Err] {

	readFile := getReadFile(bufSize, mode)

	return GetModIEP[Void, FileNames, FileNames, Collection[int64, Extent, Err], Collection[int64, Extent, Err]](
		readFile,
		func(ctx context.Context, fmap func(index Void, focus Collection[int64, Extent, Err]) (Collection[int64, Extent, Err], error), source FileNames) (FileNames, error) {

			_, fileSeq, err := readFile(ctx, source)
			if err != nil {
				return source, err
			}

			newFileSeq, err := fmap(Void{}, fileSeq)
			err = JoinCtxErr(ctx, err)
			if err != nil {
				return source, err
			}

			writer, err := os.OpenFile(source.OutFile+".part", os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
			err = JoinCtxErr(ctx, err)
			if err != nil {
				return source, err
			}

			bufWriter := bufio.NewWriterSize(writer, bufSize)

			var retErr error

			newFileSeq.AsIter()(ctx)(func(val ValueIE[int64, Extent]) bool {
				_, focus, err := val.Get()
				if err != nil {
					retErr = err
					return false
				}

				_, err = focus.WriteTo(bufWriter)
				err = JoinCtxErr(ctx, err)
				if err != nil {
					retErr = err
					return false
				}

				return true
			})

			if retErr != nil {
				writer.Close()
				os.Remove(source.OutFile + ".part")
				return source, retErr
			}

			err = bufWriter.Flush()
			if err != nil {
				writer.Close()
				os.Remove(source.OutFile + ".part")
				return source, err
			}

			return source, errors.Join(writer.Close(), os.Rename(source.OutFile+".part", source.OutFile))
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.TraverseFile{
				OpticTypeExpr: ot,
				Perm:          perm,
				BufSize:       bufSize,
			}
		}),
	)
}

func IterateFile(perm os.FileMode, bufSize int, mode ReadMode) Optic[Void, string, string, Collection[int64, Extent, Err], Collection[int64, Extent, Err], ReturnOne, ReadOnly, UniDir, Err] {

	readFile := getReadFile(bufSize, mode)

	return CombiGetter[Err, Void, string, string, Collection[int64, Extent, Err], Collection[int64, Extent, Err]](
		func(ctx context.Context, source string) (Void, Collection[int64, Extent, Err], error) {
			_, col, err := readFile(ctx, FileNames{InFile: source, OutFile: source})
			return Void{}, col, err
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.TraverseFile{
				OpticTypeExpr: ot,
				Perm:          perm,
				BufSize:       bufSize,
			}
		}),
	)
}
