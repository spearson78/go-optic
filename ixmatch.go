package optic

import (
	"context"
	"reflect"

	"github.com/samber/lo"
	"github.com/spearson78/go-optic/expr"
)

func PredToIxMatch[I any, ERR TPure](ixMatch Predicate[lo.Tuple2[I, I], ERR]) func(a, b I) bool {
	return func(a, b I) bool {
		ret, err := PredGet(context.Background(), ixMatch, lo.T2(a, b))
		if err != nil {
			panic(err) //ixMatch is pure
		}
		return ret
	}
}

func AsIxMatchT2[I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR]) Optic[Void, lo.Tuple2[I, I], lo.Tuple2[I, I], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return CombiGetter[Pure, Void, lo.Tuple2[I, I], lo.Tuple2[I, I], bool, bool](
		func(ctx context.Context, source lo.Tuple2[I, I]) (Void, bool, error) {
			ret := o.AsIxMatch()(source.A, source.B)
			return Void{}, ret, nil
		},
		IxMatchVoid(),
		ExprDef(
			func(t expr.OpticTypeExpr) expr.OpticExpression {
				return expr.IxMatch{
					OpticTypeExpr: t,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

func IxMatchVoid() func(indexA, indexB Void) bool {
	return func(indexA, indexB Void) bool {
		return true
	}
}

func IxMatchComparable[I comparable]() func(indexA, indexB I) bool {
	return func(indexA, indexB I) bool {
		return indexA == indexB
	}
}

func IxMatchDeep[I any]() func(indexA, indexB I) bool {
	return func(indexA, indexB I) bool {
		return reflect.DeepEqual(indexA, indexB)
	}
}

func ensureSimpleIxMatch[I any](fnc func(indexA, indexB I) bool) func(indexA, indexB I) bool {
	if fnc != nil {
		return func(indexA, indexB I) bool {
			return fnc(indexA, indexB)
		}
	} else {
		return IxMatchDeep[I]()
	}
}
