package ojson

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
)

type bytesReaderCloser struct {
	*bytes.Reader
}

func (*bytesReaderCloser) Close() error {
	return nil
}

func Parse[A any]() Optic[Void, []byte, []byte, A, A, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoEP[[]byte, []byte, A, A](
		func(ctx context.Context, source []byte) (A, error) {
			var ret A
			err := json.Unmarshal(source, &ret)
			return ret, JoinCtxErr(ctx, err)
		},
		func(ctx context.Context, focus A) ([]byte, error) {
			ret, err := json.Marshal(&focus)
			return ret, JoinCtxErr(ctx, err)
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.ParseJson{
				OpticTypeExpr: ot,
			}
		}),
	)
}

func ParseString[A any]() Optic[Void, string, string, A, A, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoEP[string, string, A, A](
		func(ctx context.Context, source string) (A, error) {
			var ret A
			err := json.Unmarshal([]byte(source), &ret)
			return ret, JoinCtxErr(ctx, err)
		},
		func(ctx context.Context, focus A) (string, error) {
			ret, err := json.Marshal(&focus)
			return string(ret), JoinCtxErr(ctx, err)
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.ParseJson{
				OpticTypeExpr: ot,
			}
		}),
	)
}

func ParseReader[A any]() Optic[Void, io.Reader, io.Reader, A, A, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoEP[io.Reader, io.Reader, A, A](
		func(ctx context.Context, source io.Reader) (A, error) {
			var ret A
			err := json.NewDecoder(source).Decode(&ret)
			return ret, JoinCtxErr(ctx, err)
		},
		func(ctx context.Context, focus A) (io.Reader, error) {
			newData, err := json.Marshal(&focus)
			var ret io.ReadCloser
			ret = &bytesReaderCloser{
				Reader: bytes.NewReader(newData),
			}
			return ret, JoinCtxErr(ctx, err)
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.ParseJson{
				OpticTypeExpr: ot,
			}
		}),
	)
}
