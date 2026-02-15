package optic

import (
	"context"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spearson78/go-optic/expr"
)

// Reduction is non polymorphic interface that reduces multiple values into a single result.
//
// For a polymorphic version see [ReductionP]
//
//   - The Empty method should return a value for the initial state of the Reducer. e.g. 0 for sum and 1 for product.
//   - The Reduce method should combine the 2 values and return a new combined result. e.g a+b for sum and a*b for product.
//   - AsExpr should return the expression type. See the [expr] package for more information.
type Reduction[S, ERR any] ReductionP[S, S, S, ERR]

// ReducerP is a polymorphic interface that reduces multiple values into a single result.
//
// For a non polymorphic version see [ReducerP]
// The S type parameter is the internal state that will be appended to.
// The E type parameter is the type of values that will be reduced.
// The R type parameter is the final result type.
//
//   - The Empty method should return an initial state of the Reducer. e.g. 0 for sum and 1 for product.
//   - The Reduce method should combine the previous state with a new value and return a new combined state. e.g a+b for sum and a*b for product.
//   - The End method should extract the final result from the state.
//   - AsExpr should return the expression type. See the [expr] package for more information.
type ReductionP[S, E, R, ERR any] interface {
	Empty(ctx context.Context) (S, error)
	Reduce(ctx context.Context, a S, b E) (S, error)
	End(ctx context.Context, s S) (R, error)
	ErrType() ERR
	AsExpr() expr.ReducerExpression
	AsExprHandler() func(ctx context.Context) (ExprHandler, error)
}

func Reducer[A any](
	empty func() A,
	append func(state, appendVal A) A,
	expr ReducerExpressionDef,
) Reduction[A, Pure] {
	return ReducerP[A, A, A](
		empty,
		append,
		func(state A) A { return state },
		expr,
	)
}

func ReducerP[S, E, R any](
	empty func() S,
	append func(state S, appendVal E) S,
	end func(state S) R,
	expression ReducerExpressionDef,
) ReductionP[S, E, R, Pure] {
	return &reducer[S, E, R, Pure]{
		empty: func(ctx context.Context) (S, error) { return empty(), nil },
		append: func(ctx context.Context, t lo.Tuple2[S, E]) (S, error) {
			next := append(t.A, t.B)
			return next, ctx.Err()
		},
		end: func(ctx context.Context, s S) (R, error) {
			return end(s), ctx.Err()
		},
		expr: func() expr.ReducerExpression {
			return expression.expression(expr.NewReducerTypeExpr[S, E, R]())
		},
		handler: expression.handler,
	}
}

func ReducerE[A any](
	empty func(ctx context.Context) (A, error),
	append func(ctx context.Context, state, appendVal A) (A, error),
	expression ReducerExpressionDef,
) Reduction[A, Err] {
	return ReducerEP[A, A, A](
		empty,
		append,
		func(ctx context.Context, state A) (A, error) { return state, nil },
		expression,
	)
}

func ReducerEP[S, E, R any](
	empty func(ctx context.Context) (S, error),
	append func(ctx context.Context, state S, appendVal E) (S, error),
	end func(ctx context.Context, state S) (R, error),
	expression ReducerExpressionDef,
) ReductionP[S, E, R, Err] {
	return &reducer[S, E, R, Err]{
		empty: empty,
		append: func(ctx context.Context, t lo.Tuple2[S, E]) (S, error) {
			next, err := append(ctx, t.A, t.B)
			return next, JoinCtxErr(ctx, err)
		},
		end: func(ctx context.Context, s S) (R, error) {
			ret, err := end(ctx, s)
			return ret, JoinCtxErr(ctx, err)
		},
		expr: func() expr.ReducerExpression {
			return expression.expression(expr.NewReducerTypeExpr[S, E, R]())
		},
		handler: expression.handler,
	}
}

func CombiReducer[ERR, S, E, R any](
	empty func(ctx context.Context) (S, error),
	append func(ctx context.Context, state S, appendVal E) (S, error),
	end func(ctx context.Context, state S) (R, error),
	expression ReducerExpressionDef,
) ReductionP[S, E, R, ERR] {
	return &reducer[S, E, R, ERR]{
		empty: empty,
		append: func(ctx context.Context, t lo.Tuple2[S, E]) (S, error) {
			next, err := append(ctx, t.A, t.B)
			return next, JoinCtxErr(ctx, err)
		},
		end: func(ctx context.Context, s S) (R, error) {
			ret, err := end(ctx, s)
			return ret, JoinCtxErr(ctx, err)
		},
		expr: func() expr.ReducerExpression {
			return expression.expression(expr.NewReducerTypeExpr[S, E, R]())
		},
		handler: expression.handler,
	}
}

type ReducerExpressionDef struct {
	handler    func(ctx context.Context) (ExprHandler, error)
	expression func(t expr.ReducerTypeExpr) expr.ReducerExpression
}

func ReducerExprDef(
	expression func(t expr.ReducerTypeExpr) expr.ReducerExpression,
	subOptics ...asExprHandler,
) ReducerExpressionDef {
	return ReducerExpressionDef{
		handler:    expressionHandlerCompose(subOptics...),
		expression: expression,
	}
}

func ReducerExprDefVarArgs[T asExprHandler](
	expression func(t expr.ReducerTypeExpr) expr.ReducerExpression,
	subOptics ...T,
) ReducerExpressionDef {
	return ReducerExpressionDef{
		handler:    expressionHandlerComposeVarArgs(subOptics...),
		expression: expression,
	}
}

func ReducerExprCustom(id string) ReducerExpressionDef {
	return ReducerExpressionDef{
		handler:    nil,
		expression: expr.CustomReducer(id),
	}
}

func ReducerExprTODO(id string) ReducerExpressionDef {
	return ReducerExpressionDef{
		handler:    nil,
		expression: expr.CustomReducer(id),
	}
}

type reducer[S, E, R, ERR any] struct {
	empty   func(context.Context) (S, error)
	append  func(context.Context, lo.Tuple2[S, E]) (S, error)
	end     func(context.Context, S) (R, error)
	expr    func() expr.ReducerExpression
	handler func(ctx context.Context) (ExprHandler, error)
}

func (m *reducer[S, E, R, ERR]) Empty(ctx context.Context) (S, error) {
	return m.empty(ctx)
}
func (m *reducer[S, E, R, ERR]) Reduce(ctx context.Context, a S, b E) (S, error) {
	next, err := m.append(ctx, lo.T2(a, b))
	err = JoinCtxErr(ctx, err)
	return next, err
}

func (m *reducer[S, E, R, ERR]) End(ctx context.Context, v S) (R, error) {
	return m.end(ctx, v)
}

func (m *reducer[S, E, R, ERR]) ErrType() ERR {
	var err ERR
	return err
}

func (m *reducer[S, E, R, ERR]) AsExpr() expr.ReducerExpression {
	return m.expr()
}

func (m *reducer[S, E, R, ERR]) AsExprHandler() func(ctx context.Context) (ExprHandler, error) {
	return m.handler
}

// AsReducer returns a [Reducer] from an [Operation] function. like [Add]
//
//   - empty is the initial value
//   - reduce is the operation function to use as the append operation
//
// See:
//   - AsReducerP for a version that supports end transformations.
func AsReducer[S, E any, RET TReturnOne, ERR any](empty S, reduce Operation[lo.Tuple2[S, E], S, RET, ERR]) ReductionP[S, E, S, ERR] {
	return &reducer[S, E, S, ERR]{
		empty:  func(ctx context.Context) (S, error) { return empty, nil },
		append: reduce.AsOpGet(),
		end: func(ctx context.Context, s S) (S, error) {
			return s, nil
		},
		expr: func() expr.ReducerExpression {
			return expr.AsReducer{
				ReducerTypeExpr: expr.NewReducerTypeExpr[S, E, S](),
				Empty:           Const[Void](empty).AsExpr(),
				Append:          reduce.AsExpr(),
			}
		},
		handler: reduce.AsExprHandler(),
	}
}

// AsReducerP returns a [Reducer] from an [Operation] function. like [Add]
//
//   - empty is the initial value
//   - reduce is the operation function to use as the append operation
//   - end is the operation to transform the state to the final result type.
func AsReducerP[S, E, R any, EMRET, RET, ERET TReturnOne, EMERR, ERR any, EERR TPure](empty Operation[Void, S, EMRET, EMERR], reduce Operation[lo.Tuple2[S, E], S, RET, ERR], end Operation[S, R, ERET, EERR]) ReductionP[S, E, R, ERR] {
	return &reducer[S, E, R, ERR]{
		empty:  func(ctx context.Context) (S, error) { return empty.AsOpGet()(ctx, Void{}) },
		append: reduce.AsOpGet(),
		end:    end.AsOpGet(),
		expr: func() expr.ReducerExpression {
			return expr.AsReducer{
				ReducerTypeExpr: expr.NewReducerTypeExpr[S, E, S](),
				Empty:           empty.AsExpr(),
				Append:          reduce.AsExpr(),
				End:             end.AsExpr(),
			}
		},
		handler: reduce.AsExprHandler(),
	}
}

// ReducerT2 returns a composite [ReducerP] enabling 2 aggregations to be executed in a single pass.
//
// For a 3 value variant see [ReducerT3]
func ReducerT2[SA, A, RA, SB, B, RB any, ERRA, ERRB any](ma ReductionP[SA, A, RA, ERRA], mb ReductionP[SB, B, RB, ERRB]) ReductionP[lo.Tuple2[SA, SB], lo.Tuple2[A, B], lo.Tuple2[RA, RB], CompositionTree[ERRA, ERRB]] {
	return CombiReducer[CompositionTree[ERRA, ERRB], lo.Tuple2[SA, SB], lo.Tuple2[A, B], lo.Tuple2[RA, RB]](
		func(ctx context.Context) (lo.Tuple2[SA, SB], error) {
			a, err := ma.Empty(ctx)
			if err != nil {
				return lo.Tuple2[SA, SB]{}, err
			}
			b, err := mb.Empty(ctx)
			if err != nil {
				return lo.Tuple2[SA, SB]{}, err
			}
			return lo.T2(a, b), nil
		},
		func(ctx context.Context, state lo.Tuple2[SA, SB], appendVal lo.Tuple2[A, B]) (lo.Tuple2[SA, SB], error) {
			na, err := ma.Reduce(ctx, state.A, appendVal.A)
			if err != nil {
				var ret lo.Tuple2[SA, SB]
				return ret, err
			}
			nb, err := mb.Reduce(ctx, state.B, appendVal.B)
			return lo.T2(na, nb), JoinCtxErr(ctx, err)
		},
		func(ctx context.Context, state lo.Tuple2[SA, SB]) (lo.Tuple2[RA, RB], error) {
			a, err := ma.End(ctx, state.A)
			if err != nil {
				return lo.Tuple2[RA, RB]{}, err
			}
			b, err := mb.End(ctx, state.B)
			if err != nil {
				return lo.Tuple2[RA, RB]{}, err
			}

			return lo.T2(a, b), nil
		},
		ReducerExprDef(
			func(t expr.ReducerTypeExpr) expr.ReducerExpression {
				return expr.ReducerN{
					ReducerTypeExpr: t,
					Reducers:        []expr.ReducerExpression{ma.AsExpr(), mb.AsExpr()},
				}
			},
			ma,
			mb,
		),
	)
}

// ReducerT3 returns a composite [ReducerP] enabling 3 aggregations to be executed in a single pass.
//
// For a 2 value variant see [ReducerT2]
func ReducerT3[SA, A, RA, SB, B, RB, SC, C, RC any, ERRA, ERRB, ERRC any](ma ReductionP[SA, A, RA, ERRA], mb ReductionP[SB, B, RB, ERRB], mc ReductionP[SC, C, RC, ERRC]) ReductionP[lo.Tuple3[SA, SB, SC], lo.Tuple3[A, B, C], lo.Tuple3[RA, RB, RC], CompositionTree[CompositionTree[ERRA, ERRB], ERRC]] {
	return CombiReducer[CompositionTree[CompositionTree[ERRA, ERRB], ERRC], lo.Tuple3[SA, SB, SC], lo.Tuple3[A, B, C], lo.Tuple3[RA, RB, RC]](
		func(ctx context.Context) (lo.Tuple3[SA, SB, SC], error) {
			a, err := ma.Empty(ctx)
			if err != nil {
				return lo.Tuple3[SA, SB, SC]{}, err
			}
			b, err := mb.Empty(ctx)
			if err != nil {
				return lo.Tuple3[SA, SB, SC]{}, err
			}
			c, err := mc.Empty(ctx)
			if err != nil {
				return lo.Tuple3[SA, SB, SC]{}, err
			}
			return lo.T3(a, b, c), nil
		},
		func(ctx context.Context, state lo.Tuple3[SA, SB, SC], appendVal lo.Tuple3[A, B, C]) (lo.Tuple3[SA, SB, SC], error) {
			na, err := ma.Reduce(ctx, state.A, appendVal.A)
			if err != nil {
				var ret lo.Tuple3[SA, SB, SC]
				return ret, err
			}
			nb, err := mb.Reduce(ctx, state.B, appendVal.B)
			if err != nil {
				var ret lo.Tuple3[SA, SB, SC]
				return ret, err
			}
			nc, err := mc.Reduce(ctx, state.C, appendVal.C)
			if err != nil {
				var ret lo.Tuple3[SA, SB, SC]
				return ret, err
			}

			return lo.T3(na, nb, nc), JoinCtxErr(ctx, err)
		},
		func(ctx context.Context, state lo.Tuple3[SA, SB, SC]) (lo.Tuple3[RA, RB, RC], error) {
			a, err := ma.End(ctx, state.A)
			if err != nil {
				return lo.Tuple3[RA, RB, RC]{}, err
			}
			b, err := mb.End(ctx, state.B)
			if err != nil {
				return lo.Tuple3[RA, RB, RC]{}, err
			}
			c, err := mc.End(ctx, state.C)
			if err != nil {
				return lo.Tuple3[RA, RB, RC]{}, err
			}

			return lo.T3(a, b, c), nil
		},
		ReducerExprDef(
			func(t expr.ReducerTypeExpr) expr.ReducerExpression {
				return expr.ReducerN{
					ReducerTypeExpr: t,
					Reducers:        []expr.ReducerExpression{ma.AsExpr(), mb.AsExpr(), mc.AsExpr()},
				}
			},
			ma,
			mb,
			mc,
		),
	)
}

// FirstReducer returns a [Reducer] that reduces to the first element.
func FirstReducer[T any]() ReductionP[mo.Option[T], T, T, Pure] {
	var t T
	return AsReducerP(
		Const[Void](mo.None[T]()),
		UnsafeReconstrain[ReturnOne, ReadOnly, UniDir, Pure](
			Coalesce(
				Filtered(
					T2A[mo.Option[T], T](),
					Present[T](),
				),
				ComposeLeft(
					T2B[mo.Option[T], T](),
					Embed(Some[T]()),
				),
			),
		),
		FirstOrDefault(
			Some[T](),
			t,
		),
	)
}

// LastReducer returns a [Reducer] that reduces to the last element.
func LastReducer[T any]() ReductionP[mo.Option[T], T, T, Pure] {

	var t T
	return AsReducerP(
		Const[Void](mo.None[T]()),
		UnsafeReconstrain[ReturnOne, ReadOnly, UniDir, Pure](
			Compose(
				T2B[mo.Option[T], T](),
				Embed(Some[T]()),
			),
		),
		FirstOrDefault(
			Some[T](),
			t,
		),
	)
}
