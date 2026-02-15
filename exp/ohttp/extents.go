package ohttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/exp/oio"
	"github.com/spearson78/go-optic/expr"
)

func DefaultReadCloser() Optic[Void, string, string, io.ReadCloser, io.ReadCloser, ReturnOne, ReadWrite, UniDir, Err] {
	return ReadCloser[ReturnOne, Pure](nil, nil)
}

type Rest struct {
	expr.OpticTypeExpr
	Client     *http.Client
	Middleware expr.OpticExpression
}

func (e Rest) Short() string {
	return fmt.Sprintf("Rest(%v)", e.Middleware)
}

func (e Rest) String() string {
	return fmt.Sprintf("Rest(%v)", e.Middleware)
}

func ReadCloser[RET TReturnOne, ERR any](client *http.Client, middleware Operation[*http.Request, *http.Request, RET, ERR]) Optic[Void, string, string, io.ReadCloser, io.ReadCloser, ReturnOne, ReadWrite, UniDir, Err] {

	if client == nil {
		client = http.DefaultClient
	}

	return CombiLens[ReadWrite, Err, Void, string, string, io.ReadCloser, io.ReadCloser](
		func(ctx context.Context, source string) (Void, io.ReadCloser, error) {

			req, err := http.NewRequestWithContext(ctx, "GET", source, nil)
			if err != nil {
				return Void{}, nil, err
			}

			if middleware != nil {
				req, err = middleware.AsOpGet()(ctx, req)
				if err != nil {
					return Void{}, nil, err
				}
			}

			res, err := client.Do(req)
			if err != nil {
				return Void{}, nil, err
			}

			if res.StatusCode != http.StatusOK {
				return Void{}, nil, HttpError(res.StatusCode)
			}

			return Void{}, res.Body, nil
		},
		func(ctx context.Context, focus io.ReadCloser, source string) (string, error) {
			defer focus.Close()
			req, err := http.NewRequestWithContext(ctx, "POST", source, focus)
			if err != nil {
				return source, err
			}
			if middleware != nil {
				req, err = middleware.AsOpGet()(ctx, req)
				if err != nil {
					return source, err
				}
			}
			res, err := client.Do(req)
			if err != nil {
				return source, err
			}
			res.Body.Close()

			if res.StatusCode != http.StatusOK {
				return source, HttpError(res.StatusCode)
			}

			return source, nil
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {

				var m expr.OpticExpression
				if middleware != nil {
					m = middleware.AsExpr()
				}

				return Rest{
					OpticTypeExpr: ot,
					Client:        client,
					Middleware:    m,
				}
			},
			middleware,
		),
	)
}

func DefaultBytes() Optic[Void, string, string, []byte, []byte, ReturnOne, ReadWrite, UniDir, Err] {
	return Bytes[ReturnOne, Pure](nil, nil)
}

func Bytes[RET TReturnOne, ERR any](client *http.Client, middleware Operation[*http.Request, *http.Request, RET, ERR]) Optic[Void, string, string, []byte, []byte, ReturnOne, ReadWrite, UniDir, Err] {

	if client == nil {
		client = http.DefaultClient
	}

	return CombiLens[ReadWrite, Err, Void, string, string, []byte, []byte](
		func(ctx context.Context, source string) (Void, []byte, error) {

			req, err := http.NewRequestWithContext(ctx, "GET", source, nil)
			if err != nil {
				return Void{}, nil, err
			}
			if middleware != nil {
				req, err = middleware.AsOpGet()(ctx, req)
				if err != nil {
					return Void{}, nil, err
				}
			}
			res, err := client.Do(req)
			if err != nil {
				return Void{}, nil, err
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				return Void{}, nil, HttpError(res.StatusCode)
			}

			ret, err := io.ReadAll(res.Body)

			return Void{}, ret, err

		},
		func(ctx context.Context, focus []byte, source string) (string, error) {

			req, err := http.NewRequestWithContext(ctx, "POST", source, bytes.NewReader(focus))
			if err != nil {
				return source, err
			}
			if middleware != nil {
				req, err = middleware.AsOpGet()(ctx, req)
				if err != nil {
					return source, err
				}
			}
			res, err := client.Do(req)
			if err != nil {
				return source, err
			}
			res.Body.Close()

			if res.StatusCode != http.StatusOK {
				return source, HttpError(res.StatusCode)
			}

			return source, nil
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {

				var m expr.OpticExpression
				if middleware != nil {
					m = middleware.AsExpr()
				}

				return Rest{
					OpticTypeExpr: ot,
					Client:        client,
					Middleware:    m,
				}
			},
			middleware,
		),
	)
}

func DefaultExtents() Optic[Void, string, string, Collection[int64, oio.Extent, Err], Collection[int64, oio.Extent, Err], ReturnOne, ReadWrite, UniDir, Err] {
	return Extents[ReturnOne, Pure](nil, 8192, nil)
}

func Extents[RET TReturnOne, ERR any](client *http.Client, bufSize int, middleware Operation[*http.Request, *http.Request, RET, ERR]) Optic[Void, string, string, Collection[int64, oio.Extent, Err], Collection[int64, oio.Extent, Err], ReturnOne, ReadWrite, UniDir, Err] {

	return CombiLens[ReadWrite, Err, Void, string, string, Collection[int64, oio.Extent, Err], Collection[int64, oio.Extent, Err]](
		func(ctx context.Context, source string) (Void, Collection[int64, oio.Extent, Err], error) {

			req, err := http.NewRequestWithContext(ctx, "GET", source, nil)
			if err != nil {
				return Void{}, nil, err
			}
			if middleware != nil {
				req, err = middleware.AsOpGet()(ctx, req)
				if err != nil {
					return Void{}, nil, err
				}
			}
			res, err := client.Do(req)
			if err != nil {
				return Void{}, nil, err
			}

			if res.StatusCode != http.StatusOK {
				return Void{}, nil, HttpError(res.StatusCode)
			}

			ret, err := oio.ReadCloser(bufSize).AsReverseGetter()(ctx, res.Body)
			return Void{}, ret, err

		},
		func(ctx context.Context, focus Collection[int64, oio.Extent, Err], source string) (string, error) {

			readCloser, err := Get(oio.ReadCloser(bufSize), focus)
			if err != nil {
				return source, err
			}
			defer readCloser.Close()

			req, err := http.NewRequestWithContext(ctx, "POST", source, readCloser)
			if err != nil {
				return source, err
			}
			if middleware != nil {
				req, err = middleware.AsOpGet()(ctx, req)
				if err != nil {
					return source, err
				}
			}
			res, err := client.Do(req)
			if err != nil {
				return source, err
			}
			res.Body.Close()

			if res.StatusCode != http.StatusOK {
				return source, HttpError(res.StatusCode)
			}

			return source, nil
		},
		IxMatchVoid(),
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {

				var m expr.OpticExpression
				if middleware != nil {
					m = middleware.AsExpr()
				}

				return Rest{
					OpticTypeExpr: ot,
					Client:        client,
					Middleware:    m,
				}
			},
			middleware,
		),
	)

}
