package optic

import (
	"context"
	"errors"

	"github.com/spearson78/go-optic/expr"
)

var ErrEmptyGet = errors.New("empty getter result")

// Constructor for a Getter optic. Lenses focus on a single element.
//
// The following constructors are available.
//   - [Getter]
//   - [GetterE] error aware
//   - [GetterI] indexed
//   - [GetterIE] indexed,error aware
//
// A getter is constructed from 2 functions
//
//   - get : should return the single value focused by the lens
//   - expr: should return the expression type. See the [expr] package for more information.
func Getter[S, A any](
	getFnc func(source S) A,
	exprDef ExpressionDef,
) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, Pure] {
	return CombiGetter[Pure, Void, S, S, A, A](
		func(ctx context.Context, source S) (Void, A, error) {
			return Void{}, getFnc(source), nil
		},
		IxMatchVoid(),
		exprDef,
	)
}

// Constructor for a Getter optic. Lenses focus on a single element.
//
// The following constructors are available.
//   - [Getter]
//   - [GetterE] error aware
//   - [GetterI] indexed
//   - [GetterIE] indexed,error aware
//
// A getter is constructed from 2 functions
//
//   - get : should return the single value focused by the lens
//   - expr: should return the expression type. See the [expr] package for more information.
func GetterE[S, A any](
	getFnc func(ctx context.Context, source S) (A, error),
	exprDef ExpressionDef,
) Optic[Void, S, S, A, A, ReturnOne, ReadOnly, UniDir, Err] {
	return CombiGetter[Err, Void, S, S, A, A](
		func(ctx context.Context, source S) (Void, A, error) {
			a, err := getFnc(ctx, source)
			return Void{}, a, err
		},
		IxMatchVoid(),
		exprDef,
	)
}

// Constructor for a Getter optic. Lenses focus on a single element.
//
// The following constructors are available.
//   - [Getter]
//   - [GetterE] error aware
//   - [GetterI] indexed
//   - [GetterIE] indexed,error aware
//
// A getter is constructed from 2 functions
//
//   - get : should return the single value focused by the lens
//   - expr: should return the expression type. See the [expr] package for more information.
func GetterI[I, S, A any](
	getFnc func(source S) (I, A),
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, S, A, A, ReturnOne, ReadOnly, UniDir, Pure] {

	ixMatchFnc := ensureSimpleIxMatch(ixMatch)

	return CombiGetter[Pure, I, S, S, A, A](
		func(ctx context.Context, source S) (I, A, error) {
			i, a := getFnc(source)
			return i, a, nil
		},
		ixMatchFnc,
		exprDef,
	)
}

// Constructor for a Getter optic. Lenses focus on a single element.
//
// The following constructors are available.
//   - [Getter]
//   - [GetterE] error aware
//   - [GetterI] indexed
//   - [GetterIE] indexed,error aware
//
// A getter is constructed from 2 functions
//
//   - get : should return the single value focused by the lens
//   - expr: should return the expression type. See the [expr] package for more information.
func GetterIE[I, S, A any](
	getFnc func(ctx context.Context, source S) (I, A, error),
	ixMatch func(indexA, indexB I) bool,
	exprDef ExpressionDef,
) Optic[I, S, S, A, A, ReturnOne, ReadOnly, UniDir, Err] {

	return CombiGetter[Err, I, S, S, A, A](
		getFnc,
		ixMatch,
		exprDef,
	)
}

func rawGetterF[I, S, A, ERR any](
	getFnc func(ctx context.Context, source S) (I, A, error),
	ixMatch func(indexA, indexB I) bool,
	handler func(ctx context.Context) (ExprHandler, error),
	expression func() expr.OpticExpression,
) Optic[I, S, S, A, A, ReturnOne, ReadOnly, UniDir, ERR] {

	ixMatch = ensureIxMatch(ixMatch)

	return UnsafeOmni[I, S, S, A, A, ReturnOne, ReadOnly, UniDir, ERR](
		func(ctx context.Context, source S) (I, A, error) {
			i, a, err := getFnc(ctx, source)
			return i, a, JoinCtxErr(ctx, err)
		},
		func(ctx context.Context, focus A, source S) (S, error) {
			return source, UnsupportedOpticMethod
		},
		func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				i, a, err := getFnc(ctx, source)
				err = JoinCtxErr(ctx, err)
				yield(ValIE(i, a, err))
			}
		},
		func(ctx context.Context, source S) (int, error) {
			return 1, ctx.Err()
		},
		func(ctx context.Context, fmap func(index I, focus A) (A, error), source S) (S, error) {
			return source, UnsupportedOpticMethod
		},
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				i, a, err := getFnc(ctx, source)
				err = JoinCtxErr(ctx, err)
				if err != nil {
					yield(ValIE(i, a, err))
					return
				}

				match := ixMatch(index, i)
				if match {
					yield(ValIE(i, a, JoinCtxErr(ctx, err)))
					return
				}
			}
		},
		ixMatch,
		func(ctx context.Context, focus A) (S, error) {
			var s S
			return s, UnsupportedOpticMethod
		},
		handler,
		expression,
	)

}
