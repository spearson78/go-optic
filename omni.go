package optic

import (
	"context"

	"github.com/spearson78/go-optic/expr"
)

type omniOptic[I, S, T, A, B, RET, RW, DIR, ERR any] struct {
	getter GetterFunc[I, S, A]
	setter SetterFunc[S, T, B]

	iter         IterFunc[I, S, A]
	lengthGetter LengthGetterFunc[S]
	modify       ModifyFunc[I, S, T, A, B]

	ixGetter IxGetterFunc[I, S, A]
	ixMatch  IxMatchFunc[I]

	reverseGetter ReverseGetterFunc[T, B]

	exprHandler func(ctx context.Context) (ExprHandler, error)

	expr      func() expr.OpticExpression
	opticType expr.OpticType
}

// UnsafeOmni returns an optic using the given functions directly.
// UnsafeOmni optics should only be used when creating custom combinators.
//
// Warning: The UnsafeOmni optic does not manage the [OpticErrorPath] this must be done manually using the [OpticError] method.
// Warning: The UnsafeOmni optic does not protect against unsafe iter functions (yield after break protection).
// Warning: The UnsafeOmni optic does not handle context deadlines and cancellation.
//
// See:
//   - [Omni] for a safe version
func UnsafeOmni[I, S, T, A, B, RET, RW, DIR, ERR any](
	getter GetterFunc[I, S, A],
	setter SetterFunc[S, T, B],

	iter IterFunc[I, S, A],
	lengthGetter LengthGetterFunc[S],
	modify ModifyFunc[I, S, T, A, B],

	ixGetter IxGetterFunc[I, S, A],
	ixMatchFnc IxMatchFunc[I],

	reverseGetter ReverseGetterFunc[T, B],

	handler func(ctx context.Context) (ExprHandler, error),
	expression func() expr.OpticExpression,
) Optic[I, S, T, A, B, RET, RW, DIR, ERR] {

	if lengthGetter == nil {
		lengthGetter = func(ctx context.Context, source S) (int, error) {
			count := 0
			var err error
			iter(ctx, source)(func(val ValueIE[I, A]) bool {
				focusErr := val.Error()
				if focusErr != nil {
					err = focusErr
					return false
				}
				count++
				return true
			})
			return count, err
		}
	}

	ixMatchFnc = ensureIxMatch(ixMatchFnc)

	if ixGetter == nil {
		ixGetter = func(ctx context.Context, index I, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				ixGetter(ctx, index, source)(func(val ValueIE[I, A]) bool {
					i, a, err := val.Get()
					if err != nil {
						return yield(ValIE(i, a, err))
					}

					match := ixMatchFnc(i, index)
					if match {
						return yield(ValIE(i, a, err))
					}
					return true
				})
			}

		}
	}

	return omniOptic[I, S, T, A, B, RET, RW, DIR, ERR]{
		getter:        getter,
		setter:        setter,
		iter:          iter,
		lengthGetter:  lengthGetter,
		modify:        modify,
		ixGetter:      ixGetter,
		ixMatch:       ixMatchFnc,
		reverseGetter: reverseGetter,
		exprHandler:   handler,
		expr:          expression,
		opticType:     expr.GetOpticType[RET, RW, DIR](),
	}
}

type ExpressionDef struct {
	handler    func(ctx context.Context) (ExprHandler, error)
	expression func(t expr.OpticTypeExpr) expr.OpticExpression
}

type asExprHandler interface {
	AsExprHandler() func(ctx context.Context) (ExprHandler, error)
}

func ExprDef(
	expression func(t expr.OpticTypeExpr) expr.OpticExpression,
	subOptics ...asExprHandler,
) ExpressionDef {
	return ExpressionDef{
		handler:    expressionHandlerCompose(subOptics...),
		expression: expression,
	}
}

func ExprDefVarArgs[T asExprHandler](
	expression func(t expr.OpticTypeExpr) expr.OpticExpression,
	subOptics ...T,
) ExpressionDef {
	return ExpressionDef{
		handler:    expressionHandlerComposeVarArgs(subOptics...),
		expression: expression,
	}
}

func ExprCustom(id string) ExpressionDef {
	return ExpressionDef{
		handler:    nil,
		expression: expr.Custom(id),
	}
}

func ExprTODO(id string) ExpressionDef {
	return ExpressionDef{
		handler:    nil,
		expression: expr.Custom(id),
	}
}

// Omni returns an optic using the given functions with added protection for unsafe iteration functions (yield after break) and automatic generation of the [OpticError] path
// Omni optics should only be used when creating custom combinators.
//
// Note:
// If the getter function is unable to return a result due to an empty iteration it MUST return an [ErrEmptyGet]
//
// See:
//   - [UnsafeOmni] for an unsafe version
//   - [CombiTraversal] for a simpler optic designed for [ReturnMany] combinators.
func Omni[I, S, T, A, B, RET, RW, DIR, ERR any](
	getter GetterFunc[I, S, A],
	setter SetterFunc[S, T, B],

	iter IterFunc[I, S, A],
	lengthGetter LengthGetterFunc[S],
	modify ModifyFunc[I, S, T, A, B],

	ixGetter IxGetterFunc[I, S, A],
	ixMatchFnc IxMatchFunc[I],

	reverseGetter ReverseGetterFunc[T, B],

	exprDef ExpressionDef,
) Optic[I, S, T, A, B, RET, RW, DIR, ERR] {

	exprFnc := func() expr.OpticExpression {
		return exprDef.expression(expr.NewOpticTypeExpr[I, S, T, A, B, RET, RW, DIR, ERR]())
	}

	if disableSafeYieldAndOpticError {
		return UnsafeOmni[I, S, T, A, B, RET, RW, DIR, ERR](
			getter,
			setter,
			iter,
			lengthGetter,
			modify,
			ixGetter,
			ixMatchFnc,
			reverseGetter,
			exprDef.handler,
			exprFnc,
		)
	}

	omniLengthGetter := ensureLengthGetter(lengthGetter, iter, exprFnc)
	ixMatchFnc = ensureIxMatch(ixMatchFnc)
	omniIxGet := ensureIxGetter(ixGetter, iter, ixMatchFnc, exprFnc)

	return omniOptic[I, S, T, A, B, RET, RW, DIR, ERR]{
		getter: func(ctx context.Context, source S) (I, A, error) {
			i, a, err := getter(ctx, source)
			err = OpticError(JoinCtxErr(ctx, err), exprFnc)
			return i, a, err
		},
		setter: func(ctx context.Context, focus B, source S) (T, error) {
			t, err := setter(ctx, focus, source)
			err = OpticError(JoinCtxErr(ctx, err), exprFnc)
			return t, err
		},
		iter: func(ctx context.Context, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				cont := true
				iter(ctx, source)(func(val ValueIE[I, A]) bool {
					i, a, err := val.Get()
					if !cont {
						panic(yieldAfterBreak)
					}
					err = OpticError(JoinCtxErr(ctx, err), exprFnc)
					cont = yield(ValIE(i, a, err))
					return cont
				})
			}
		},
		lengthGetter: omniLengthGetter,
		modify: func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {
			t, err := modify(ctx, func(index I, focus A) (B, error) {
				b, err := fmap(index, focus)
				err = JoinCtxErr(ctx, err)
				return b, err
			}, source)
			err = OpticError(JoinCtxErr(ctx, err), exprFnc)
			return t, err
		},
		ixGetter: omniIxGet,
		ixMatch:  ixMatchFnc,
		reverseGetter: func(ctx context.Context, focus B) (T, error) {
			t, err := reverseGetter(ctx, focus)
			err = OpticError(JoinCtxErr(ctx, err), exprFnc)
			return t, err
		},
		exprHandler: exprDef.handler,
		expr:        exprFnc,
		opticType:   expr.GetOpticType[RET, RW, DIR](),
	}
}

func (o omniOptic[I, S, T, A, B, RET, RW, DIR, ERR]) AsModify() ModifyFunc[I, S, T, A, B] {
	return o.modify
}

func (o omniOptic[I, S, T, A, B, RET, RW, DIR, ERR]) AsIter() IterFunc[I, S, A] {
	return o.iter
}

func (o omniOptic[I, S, T, A, B, RET, RW, DIR, ERR]) AsGetter() GetterFunc[I, S, A] {
	return o.getter
}

func (o omniOptic[I, S, T, A, B, RET, RW, DIR, ERR]) AsSetter() SetterFunc[S, T, B] {
	return o.setter
}

func (o omniOptic[I, S, T, A, B, RET, RW, DIR, ERR]) AsIxGetter() IxGetterFunc[I, S, A] {
	return o.ixGetter
}

func (o omniOptic[I, S, T, A, B, RET, RW, DIR, ERR]) AsIxMatch() IxMatchFunc[I] {
	return o.ixMatch
}

func (o omniOptic[I, S, T, A, B, RET, RW, DIR, ERR]) AsOpGet() OpGetFunc[S, A] {
	return func(ctx context.Context, source S) (A, error) {
		_, a, err := o.getter(ctx, source)
		return a, err
	}
}

func (o omniOptic[I, S, T, A, B, RET, RW, DIR, ERR]) AsReverseGetter() ReverseGetterFunc[T, B] {
	return o.reverseGetter
}

func (o omniOptic[I, S, T, A, B, RET, RW, DIR, ERR]) AsLengthGetter() LengthGetterFunc[S] {
	return o.lengthGetter
}

func (o omniOptic[I, S, T, A, B, RET, RW, DIR, ERR]) AsExprHandler() func(ctx context.Context) (ExprHandler, error) {
	return o.exprHandler
}

func (o omniOptic[I, S, T, A, B, RET, RW, DIR, ERR]) AsExpr() expr.OpticExpression {
	return o.expr()
}

func (o omniOptic[I, S, T, A, B, RET, RW, DIR, ERR]) OpticType() expr.OpticType {
	return o.opticType
}

func (o omniOptic[I, S, T, A, B, RET, RW, DIR, ERR]) ReturnType() RET {
	var ret RET
	return ret
}

func (o omniOptic[I, S, T, A, B, RET, RW, DIR, ERR]) ReadWriteType() RW {
	var ret RW
	return ret
}

func (o omniOptic[I, S, T, A, B, RET, RW, DIR, ERR]) DirType() DIR {
	var ret DIR
	return ret
}

func (o omniOptic[I, S, T, A, B, RET, RW, DIR, ERR]) ErrType() ERR {
	var ret ERR
	return ret
}

func ensureLengthGetter[I, S, A any](lengthGetter LengthGetterFunc[S], iterate func(ctx context.Context, source S) SeqIE[I, A], expr func() expr.OpticExpression) LengthGetterFunc[S] {
	if lengthGetter != nil {
		return func(ctx context.Context, source S) (int, error) {
			l, err := lengthGetter(ctx, source)
			err = OpticError(JoinCtxErr(ctx, err), expr)
			return l, err
		}
	}

	return func(ctx context.Context, source S) (int, error) {
		var err error
		count := 0
		cont := true
		iterate(ctx, source)(func(val ValueIE[I, A]) bool {
			if !cont {
				panic(yieldAfterBreak)
			}
			focusErr := val.Error()

			err = OpticError(JoinCtxErr(ctx, focusErr), expr)
			if err != nil {
				cont = false
				return false
			}
			count++
			return true
		})
		return count, JoinCtxErr(ctx, err)
	}

}

func ensureIxGetter[I, S, A any](ixGetter IxGetterFunc[I, S, A], iterate func(ctx context.Context, source S) SeqIE[I, A], ixMatch func(indexA, indexB I) bool, expr func() expr.OpticExpression) IxGetterFunc[I, S, A] {
	if ixGetter != nil {
		return func(ctx context.Context, index I, source S) SeqIE[I, A] {
			return func(yield func(val ValueIE[I, A]) bool) {
				cont := true
				ixGetter(ctx, index, source)(func(val ValueIE[I, A]) bool {
					i, a, err := val.Get()
					if !cont {
						panic(yieldAfterBreak)
					}
					err = OpticError(JoinCtxErr(ctx, err), expr)
					cont = yield(ValIE(i, a, err))
					return cont
				})
			}
		}
	}

	return func(ctx context.Context, index I, source S) SeqIE[I, A] {
		return func(yield func(val ValueIE[I, A]) bool) {
			cont := true
			iterate(ctx, source)(func(val ValueIE[I, A]) bool {
				focusIndex, pa, focusErr := val.Get()
				if !cont {
					panic(yieldAfterBreak)
				}

				focusErr = OpticError(JoinCtxErr(ctx, focusErr), expr)
				if focusErr != nil {
					cont = yield(ValIE(index, pa, focusErr))
					return cont
				}

				match := ixMatch(focusIndex, index)
				if match {
					cont = yield(ValIE(index, pa, ctx.Err()))
					return cont
				}

				return cont
			})
		}
	}
}

func ensureIxMatch[I any](ixMatchFnc func(indexA, indexB I) bool) func(indexA, indexB I) bool {

	if ixMatchFnc == nil {
		ixMatchFnc = IxMatchDeep[I]()
	}

	return ixMatchFnc

}
