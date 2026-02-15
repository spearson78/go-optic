package ohttp

import (
	"context"
	"encoding/json"
	"net/http"

	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/exp/oio"
	"github.com/spearson78/go-optic/expr"
)

func DefaultRestJSON[A any]() Optic[Void, string, string, A, A, ReturnOne, ReadWrite, UniDir, Err] {
	return RestJSON[A, ReturnOne, Pure](nil, 8192, nil)
}

// TODO: move to ojson when oio is completed
func ParseJsonExtents[A any](bufSize int) Optic[Void, Collection[int64, oio.Extent, Err], Collection[int64, oio.Extent, Err], A, A, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoEP[Collection[int64, oio.Extent, Err], Collection[int64, oio.Extent, Err], A, A](
		func(ctx context.Context, source Collection[int64, oio.Extent, Err]) (A, error) {
			_, r, err := oio.ReadCloser(bufSize).AsGetter()(ctx, source)
			var ret A
			if err != nil {
				return ret, err
			}
			err = json.NewDecoder(r).Decode(&ret)
			r.Close()
			return ret, JoinCtxErr(ctx, err)
		},
		func(ctx context.Context, focus A) (Collection[int64, oio.Extent, Err], error) {
			newData, err := json.Marshal(&focus)
			if err != nil {
				return nil, err
			}

			return ReverseGet(oio.Bytes(), newData)
		},
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.ParseJson{
				OpticTypeExpr: ot,
			}
		}),
	)
}

func RestJSON[A any, RET TReturnOne, ERR any](client *http.Client, bufSize int, middleware Operation[*http.Request, *http.Request, RET, ERR]) Optic[Void, string, string, A, A, ReturnOne, ReadWrite, UniDir, Err] {
	return Ret1(Rw(Ud(EErr(Compose(Extents(client, bufSize, middleware), ParseJsonExtents[A](bufSize))))))
}
