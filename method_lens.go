package optic

import (
	"context"
	"reflect"
	"runtime"

	"github.com/spearson78/go-optic/expr"
)

func MethodGetter[S, A any](fref func(S) A) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return Getter[S, A](
		fref,
		getMethodGetterExpr(fref),
	)
}

func MethodGetterI[I, S, A any](fref func(S) (I, A), ixMatch func(indexA, indexB I) bool) Optic[I, S, S, A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return GetterI[I, S, A](
		fref,
		ixMatch,
		getMethodGetterExpr(fref),
	)
}

func MethodGetterE[S, A any](fref func(S) (A, error)) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, Err] {
	return GetterE[S, A](
		func(ctx context.Context, source S) (A, error) {
			return fref(source)
		},
		getMethodGetterExpr(fref),
	)
}

func MethodGetterContext[S, A any](fref func(S, context.Context) (A, error)) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, Err] {
	return GetterE[S, A](
		func(ctx context.Context, source S) (A, error) {
			return fref(source, ctx)
		},
		getMethodGetterExpr(fref),
	)
}

func MethodGetterIE[I, S, A any](fref func(S, context.Context) (I, A, error), ixMatch func(indexA, indexB I) bool) Optic[I, S, S, A, A, ReturnOne, ReadOnly, UniDir, Err] {
	return CombiGetter[Err, I, S, S, A, A](
		func(ctx context.Context, source S) (I, A, error) {
			return fref(source, ctx)
		},
		ixMatch,
		getMethodGetterExpr(fref),
	)
}

func getMethodGetterExpr(fref any) ExpressionDef {
	return ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {

		name := runtime.FuncForPC(reflect.ValueOf(fref).Pointer()).Name()
		method := reflect.TypeOf(fref)

		return expr.MethodGetter{
			OpticTypeExpr: ot,
			MethodName:    name,
			Method:        method,
		}
	})
}

func MethodLens[S, A any](get func(S) A, set func(S, A) S) Optic[Void, S, S, A, A, ReturnOne, ReadWrite, UniDir, Pure] {
	return Lens[S, A](
		get,
		func(focus A, source S) S {
			return set(source, focus)
		},
		getMethodLensExpr(get, set),
	)
}

func MethodLensI[I, S, A any](get func(S) (I, A), set func(S, A) S, ixMatch func(indexA, indexB I) bool) Optic[I, S, S, A, A, ReturnOne, ReadWrite, UniDir, Pure] {
	return LensI[I, S, A](
		get,
		func(focus A, source S) S {
			return set(source, focus)
		},
		ixMatch,
		getMethodLensExpr(get, set),
	)
}

func MethodLensP[S, T, A, B any](get func(S) A, set func(S, B) T) Optic[Void, S, T, A, B, ReturnOne, ReadWrite, UniDir, Pure] {
	return LensP[S, T, A, B](
		get,
		func(focus B, source S) T {
			return set(source, focus)
		},
		getMethodLensExpr(get, set),
	)
}

func MethodLensE[S, A any](get func(S) (A, error), set func(S, A) (S, error)) Optic[Void, S, S, A, A, ReturnOne, ReadWrite, UniDir, Err] {
	return LensE[S, A](
		func(ctx context.Context, source S) (A, error) {
			return get(source)
		},
		func(ctx context.Context, focus A, source S) (S, error) {
			return set(source, focus)
		},
		getMethodLensExpr(get, set),
	)
}

func MethodLensContext[S, A any](get func(S, context.Context) (A, error), set func(S, context.Context, A) (S, error)) Optic[Void, S, S, A, A, ReturnOne, ReadWrite, UniDir, Err] {
	return LensE[S, A](
		func(ctx context.Context, source S) (A, error) {
			return get(source, ctx)
		},
		func(ctx context.Context, focus A, source S) (S, error) {
			return set(source, ctx, focus)
		},
		getMethodLensExpr(get, set),
	)
}

func MethodLensIEP[I, S, T, A, B any](get func(S, context.Context) (I, A, error), set func(S, context.Context, B) (T, error), ixMatch func(indexA, indexB I) bool) Optic[I, S, T, A, B, ReturnOne, ReadWrite, UniDir, Err] {
	return CombiLens[ReadWrite, Err, I, S, T, A, B](
		func(ctx context.Context, source S) (I, A, error) {
			return get(source, ctx)
		},
		func(ctx context.Context, focus B, source S) (T, error) {
			return set(source, ctx, focus)
		},
		ixMatch,
		getMethodLensExpr(get, set),
	)
}

func getMethodLensExpr(get any, set any) ExpressionDef {
	return ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {

		getName := runtime.FuncForPC(reflect.ValueOf(get).Pointer()).Name()
		getMethod := reflect.TypeOf(get)

		setName := runtime.FuncForPC(reflect.ValueOf(set).Pointer()).Name()
		setMethod := reflect.TypeOf(set)

		return expr.MethodLens{
			OpticTypeExpr: ot,
			GetterName:    getName,
			Getter:        getMethod,
			SetterName:    setName,
			Setter:        setMethod,
		}
	})
}
