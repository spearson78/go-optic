package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
	"github.com/spearson78/go-optic/internal/playground/data"
	_ "github.com/spearson78/go-optic/internal/playground/data"
	goscript "github.com/spearson78/go-script"
)

func anyfyedPredT2[I, S, T, RET, RW, DIR, ERR any](o optic.Optic[I, lo.Tuple2[S, S], lo.Tuple2[T, T], bool, bool, RET, RW, DIR, ERR]) optic.Optic[any, lo.Tuple2[any, any], lo.Tuple2[any, any], bool, bool, any, any, any, any] {
	return optic.UnsafeOmni[any, lo.Tuple2[any, any], lo.Tuple2[any, any], bool, bool, any, any, any, any](
		func(ctx context.Context, source lo.Tuple2[any, any]) (any, bool, error) {
			t2 := goscript.Coerce[lo.Tuple2[S, S]](nil, source)
			return o.AsGetter()(ctx, lo.T2(t2.A, t2.B))
		},
		func(ctx context.Context, focus bool, source lo.Tuple2[any, any]) (lo.Tuple2[any, any], error) {
			t2 := goscript.Coerce[lo.Tuple2[S, S]](nil, source)
			ret, err := o.AsSetter()(ctx, focus, lo.T2(t2.A, t2.B))
			return lo.T2[any, any](ret.A, ret.B), err
		},
		func(ctx context.Context, source lo.Tuple2[any, any]) optic.SeqIE[any, bool] {
			return func(yield func(optic.ValueIE[any, bool]) bool) {
				t2 := goscript.Coerce[lo.Tuple2[S, S]](nil, source)
				o.AsIter()(ctx, lo.T2(t2.A, t2.B))(func(val optic.ValueIE[I, bool]) bool {
					return yield(optic.ValIE(any(val.Index()), val.Value(), val.Error()))
				})
			}
		},
		func(ctx context.Context, source lo.Tuple2[any, any]) (int, error) {
			t2 := goscript.Coerce[lo.Tuple2[S, S]](nil, source)
			return o.AsLengthGetter()(ctx, lo.T2(t2.A, t2.B))
		},
		func(ctx context.Context, fmap func(index any, focus bool) (bool, error), source lo.Tuple2[any, any]) (lo.Tuple2[any, any], error) {
			t2 := goscript.Coerce[lo.Tuple2[S, S]](nil, source)
			ret, err := o.AsModify()(ctx, func(index I, focus bool) (bool, error) {
				ret, err := fmap(index, focus)
				return ret, err
			}, lo.T2(t2.A, t2.B))
			return lo.T2[any, any](ret.A, ret.B), err
		},
		func(ctx context.Context, index any, source lo.Tuple2[any, any]) optic.SeqIE[any, bool] {
			t2 := goscript.Coerce[lo.Tuple2[S, S]](nil, source)
			return func(yield func(optic.ValueIE[any, bool]) bool) {
				o.AsIxGetter()(ctx, goscript.Coerce[I](nil, index), lo.T2(t2.A, t2.B))(func(val optic.ValueIE[I, bool]) bool {
					return yield(optic.ValIE(any(val.Index()), val.Value(), val.Error()))
				})
			}
		},
		func(anyA, anyB any) bool {
			if a, ok := goscript.SafeCoerce[I](nil, anyA); ok {
				if b, ok := goscript.SafeCoerce[I](nil, anyB); ok {
					ret := o.AsIxMatch()(a, b)
					return ret
				}
			}
			return false
		},
		func(ctx context.Context, focus bool) (lo.Tuple2[any, any], error) {
			ret, err := o.AsReverseGetter()(ctx, focus)
			return lo.T2[any, any](ret.A, ret.B), err
		},
		o.AsExprHandler(),
		o.AsExpr,
	)
}

func Coerce(v reflect.Value, t reflect.Type) reflect.Value {
	r, ok := goscript.ReflectCoerce(nil, v, t)
	if !ok {
		panic(fmt.Errorf("coerce %v to %v failed", v.Type(), t))
	}
	return r
}

func anyfyedReflectP[I, S, T, A, B any](
	o reflect.Value,
	toI func(reflect.Value) I,
	fromI func(I) reflect.Value,
	toS func(reflect.Value) S,
	fromS func(S) reflect.Value,
	toT func(reflect.Value) T,
	fromT func(T) reflect.Value,
	toA func(reflect.Value) A,
	fromA func(A) reflect.Value,
	toB func(reflect.Value) B,
	fromB func(B) reflect.Value,
) optic.Optic[I, S, T, A, B, any, any, any, any] {
	return optic.UnsafeOmni[I, S, T, A, B, any, any, any, any](
		func(ctx context.Context, source S) (I, A, error) {

			fnc := reflect.ValueOf(o.Interface()).MethodByName("AsGetter")
			if fnc.IsZero() {
				asOpGet := reflect.ValueOf(o.Interface()).MethodByName("AsOpGet").Call(nil)[0]

				ret := fnc.Call([]reflect.Value{
					reflect.ValueOf(ctx),
					Coerce(fromS(source), asOpGet.Type().In(1)),
				})

				var err error
				if !ret[1].IsNil() {
					err = ret[1].Interface().(error)
				}

				return toI(reflect.ValueOf(optic.Void{})), toA(ret[0]), err

			} else {

				asGetter := fnc.Call(nil)[0]

				ret := asGetter.Call([]reflect.Value{
					reflect.ValueOf(ctx),
					Coerce(fromS(source), asGetter.Type().In(1)),
				})

				var err error
				if !ret[2].IsNil() {
					err = ret[2].Interface().(error)
				}

				return toI(ret[0]), toA(ret[1]), err
			}
		},
		func(ctx context.Context, focus B, source S) (T, error) {

			asSetter := reflect.ValueOf(o.Interface()).MethodByName("AsSetter").Call(nil)[0]

			ret := asSetter.Call([]reflect.Value{
				reflect.ValueOf(ctx),
				Coerce(fromB(focus), asSetter.Type().In(1)),
				Coerce(fromS(source), asSetter.Type().In(2)),
			})

			var err error
			if !ret[1].IsNil() {
				err = ret[1].Interface().(error)
			}

			return toT(ret[0]), err

		},
		func(ctx context.Context, source S) optic.SeqIE[I, A] {
			return func(yield func(optic.ValueIE[I, A]) bool) {

				asIter := reflect.ValueOf(o.Interface()).MethodByName("AsIter").Call(nil)[0]

				seq := asIter.Call([]reflect.Value{
					reflect.ValueOf(ctx),
					Coerce(fromS(source), asIter.Type().In(1)),
				})[0]

				yieldFnc := reflect.MakeFunc(seq.Type().In(0), func(args []reflect.Value) (results []reflect.Value) {

					ix := args[0].MethodByName("Index").Call(nil)[0]
					val := args[0].MethodByName("Value").Call(nil)[0]
					errV := args[0].MethodByName("Error").Call(nil)[0]

					var err error
					if !errV.IsNil() {
						err = errV.Interface().(error)
					}

					return []reflect.Value{reflect.ValueOf(yield(optic.ValIE(toI(ix), toA(val), err)))}
				})

				seq.Call([]reflect.Value{yieldFnc})
			}
		},
		func(ctx context.Context, source S) (int, error) {

			as := reflect.ValueOf(o.Interface()).MethodByName("AsLengthGetter").Call(nil)[0]

			ret := as.Call([]reflect.Value{
				reflect.ValueOf(ctx),
				Coerce(fromS(source), as.Type().In(1)),
			})

			var err error
			if !ret[1].IsNil() {
				err = ret[1].Interface().(error)
			}

			return int(ret[1].Int()), err

		},
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {

			asModify := reflect.ValueOf(o.Interface()).MethodByName("AsModify")
			asModifyInstance := asModify.Call(nil)[0]
			fmapType := asModifyInstance.Type().In(1)

			ret := asModifyInstance.Call([]reflect.Value{
				reflect.ValueOf(ctx),
				reflect.MakeFunc(
					fmapType,
					func(args []reflect.Value) (results []reflect.Value) {
						b, err := fmap(toI(args[0]), toA(args[1]))

						var errV reflect.Value
						if err != nil {
							errV = reflect.ValueOf(err)
						} else {
							errV = reflect.Zero(reflect.TypeFor[error]())
						}

						coerced := Coerce(fromB(b), fmapType.Out(0))

						return []reflect.Value{coerced, errV}
					},
				),
				Coerce(fromS(source), asModifyInstance.Type().In(2)),
			},
			)

			var err error
			if !ret[1].IsNil() {
				err = ret[1].Interface().(error)
			}

			return toT(ret[0]), err
		},
		func(ctx context.Context, index I, source S) optic.SeqIE[I, A] {
			return func(yield func(optic.ValueIE[I, A]) bool) {

				asIter := reflect.ValueOf(o.Interface()).MethodByName("AsIxGetter").Call(nil)[0]

				seq := asIter.Call([]reflect.Value{
					reflect.ValueOf(ctx),
					Coerce(fromI(index), asIter.Type().In(1)),
					Coerce(fromS(source), asIter.Type().In(2)),
				})[0]

				yieldFnc := reflect.MakeFunc(seq.Type().In(0), func(args []reflect.Value) (results []reflect.Value) {

					ix := args[0].MethodByName("Index").Call(nil)[0]
					val := args[0].MethodByName("Value").Call(nil)[0]
					errV := args[0].MethodByName("Error").Call(nil)[0]

					var err error
					if !errV.IsNil() {
						err = errV.Interface().(error)
					}

					yieldVal := yield(optic.ValIE(toI(ix), toA(val), err))

					return []reflect.Value{reflect.ValueOf(yieldVal)}
				})

				seq.Call([]reflect.Value{yieldFnc})
			}
		},
		func(indexA, indexB I) bool {
			asIxMatch := reflect.ValueOf(o.Interface()).MethodByName("AsIxMatch").Call(nil)[0]

			ret := asIxMatch.Call([]reflect.Value{
				Coerce(fromI(indexA), asIxMatch.Type().In(0)),
				Coerce(fromI(indexB), asIxMatch.Type().In(1)),
			})

			return ret[0].Bool()
		},
		func(ctx context.Context, focus B) (T, error) {
			asReverseGetter := reflect.ValueOf(o.Interface()).MethodByName("AsReverseGetter").Call(nil)[0]

			ret := asReverseGetter.Call([]reflect.Value{
				reflect.ValueOf(ctx),
				Coerce(fromB(focus), asReverseGetter.Type().In(1)),
			})

			var err error
			if !ret[1].IsNil() {
				err = ret[1].Interface().(error)
			}

			return toT(ret[0]), err
		},
		func(ctx context.Context) (optic.ExprHandler, error) {
			asExpr := reflect.ValueOf(o.Interface()).MethodByName("AsExprHandler")

			asExprInstance := asExpr.Call(nil)[0]
			if asExprInstance.IsNil() {
				return nil, nil
			}
			ret := asExprInstance.Call([]reflect.Value{
				reflect.ValueOf(ctx),
			})

			var handler optic.ExprHandler
			if !ret[0].IsNil() {
				handler = ret[0].Interface().(optic.ExprHandler)
			}

			var err error
			if !ret[1].IsNil() {
				err = ret[1].Interface().(error)
			}

			return handler, err
		},
		func() expr.OpticExpression {

			asExpr := reflect.ValueOf(o.Interface()).MethodByName("AsExpr")

			ret := asExpr.Call(nil)

			return ret[0].Interface().(expr.OpticExpression)

		},
	)
}

func anyfyedReflect(o reflect.Value) optic.Optic[any, any, any, any, any, any, any, any, any] {

	return anyfyedReflectP[any, any, any, any, any](
		o,
		reflect.Value.Interface,
		reflect.ValueOf,
		reflect.Value.Interface,
		reflect.ValueOf,
		reflect.Value.Interface,
		reflect.ValueOf,
		reflect.Value.Interface,
		reflect.ValueOf,
		reflect.Value.Interface,
		reflect.ValueOf,
	)
}

func getField(structVal reflect.Value, i int) reflect.Value {
	field := structVal.Field(i)
	writeableField := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
	return writeableField
}

func setField(structVal reflect.Value, i int, setVal reflect.Value) {
	field := structVal.Field(i)
	writeableField := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
	writeableField.Set(setVal)
}

func CoerceT[T any](v reflect.Value) T {
	r, ok := goscript.ReflectCoerce(nil, v, reflect.TypeFor[T]())
	if !ok {
		panic(fmt.Errorf("coerce %v to %v failed", v.Type(), reflect.TypeFor[T]()))
	}
	return r.Interface().(T)
}

func coerceOptic[I, S, T, A, B any](o reflect.Value) optic.Optic[I, S, T, A, B, any, any, any, any] {

	return anyfyedReflectP[I, S, T, A, B](
		o,
		CoerceT[I],
		func(i I) reflect.Value {
			return reflect.ValueOf(i)
		},
		CoerceT[S],
		func(i S) reflect.Value {
			return reflect.ValueOf(i)
		},
		CoerceT[T],
		func(i T) reflect.Value {
			return reflect.ValueOf(i)
		},
		CoerceT[A],
		func(b A) reflect.Value {
			return reflect.ValueOf(b)
		},
		CoerceT[B],
		func(b B) reflect.Value {
			return reflect.ValueOf(b)
		},
	)
}

func anyfyedReflectReduction(
	o reflect.Value,
) optic.ReductionP[any, any, any, any] {
	return optic.CombiReducer[any, any, any, any](
		func(ctx context.Context) (any, error) {
			fnc := reflect.ValueOf(o.Interface()).MethodByName("Empty")
			ret := fnc.Call([]reflect.Value{
				reflect.ValueOf(ctx),
			})

			var err error
			if !ret[1].IsNil() {
				err = ret[1].Interface().(error)
			}

			return ret[0].Interface(), err
		},
		func(ctx context.Context, state, appendVal any) (any, error) {
			fnc := reflect.ValueOf(o.Interface()).MethodByName("Reduce")
			ret := fnc.Call([]reflect.Value{
				reflect.ValueOf(ctx),
				Coerce(reflect.ValueOf(state), fnc.Type().In(1)),
				Coerce(reflect.ValueOf(appendVal), fnc.Type().In(2)),
			})

			var err error
			if !ret[1].IsNil() {
				err = ret[1].Interface().(error)
			}

			return ret[0].Interface(), err
		},
		func(ctx context.Context, state any) (any, error) {
			fnc := reflect.ValueOf(o.Interface()).MethodByName("End")
			ret := fnc.Call([]reflect.Value{
				reflect.ValueOf(ctx),
				Coerce(reflect.ValueOf(state), fnc.Type().In(1)),
			})

			var err error
			if !ret[1].IsNil() {
				err = ret[1].Interface().(error)
			}

			return ret[0].Interface(), err
		},
		optic.ReducerExprDef(
			func(t expr.ReducerTypeExpr) expr.ReducerExpression {
				fnc := reflect.ValueOf(o.Interface()).MethodByName("AsExpr")
				ret := fnc.Call(nil)

				return ret[0].Interface().(expr.ReducerExpression)
			},
		),
	)
}

func downCast(typeParams []reflect.Type) reflect.Value {

	//When using DownCast[any,any]() the cast always succeeds and all values are allowed through the cast.

	return reflect.ValueOf(func() any {

		return optic.PrismE[any, any](
			func(ctx context.Context, source any) (any, bool, error) {
				rs := reflect.ValueOf(source)
				val, ok := goscript.ReflectCoerce(nil, rs, typeParams[1])
				return val.Interface(), ok, nil
			},
			func(ctx context.Context, focus any) (any, error) {
				rf := reflect.ValueOf(focus)
				val, ok := goscript.ReflectCoerce(nil, rf, typeParams[0])
				if ok {
					return val.Interface(), nil
				} else {
					return rf, &optic.ErrCast{
						From: typeParams[1],
						To:   typeParams[0],
					}
				}
			},
			optic.ExprCustom("scriptDownCast"),
		)
	})
}

func orderBy(typeParams []reflect.Type) reflect.Value {

	return reflect.ValueOf(func(o optic.Optic[any, any, any, any, any, any, any, any, any]) any {
		switch typeParams[0].Kind() {
		case reflect.Int:
			return anyfyedPredT2(optic.OrderBy(
				optic.FirstOrError(
					optic.Compose(
						o,
						optic.DownCast[any, int](),
					),
					errors.New("order by cast failed"),
				),
			))
		case reflect.Float64:
			return anyfyedPredT2(optic.OrderBy(
				optic.FirstOrError(
					optic.Compose(
						o,
						optic.DownCast[any, float64](),
					),
					errors.New("order by cast failed"),
				),
			))
		case reflect.String:
			return anyfyedPredT2(optic.OrderBy(
				optic.FirstOrError(
					optic.Compose(
						o,
						optic.DownCast[any, string](),
					),
					errors.New("order by cast failed"),
				),
			))
		default:
			panic("orderBy unknown type")
		}
	})
}

func wrapOrderByI[I, A any, C cmp.Ordered]() reflect.Value {

	return reflect.ValueOf(func(o optic.Optic[any, any, any, any, any, any, any, any, any]) any {

		x := optic.FirstOrError(
			optic.Compose3(
				optic.UpCast[optic.ValueI[any, any], any](),
				o,
				optic.DownCast[any, C](),
			),
			errors.New("order by cast failed"),
		)

		z := optic.OrderByI(
			x,
		)

		return anyfyedPredT2(z)
	})

}

func orderByI(typeParams []reflect.Type) reflect.Value {
	switch typeParams[2].Kind() {
	case reflect.Int:
		return wrapOrderByI[any, any, int]()
	case reflect.Float64:
		return wrapOrderByI[any, any, float64]()
	case reflect.String:
		return wrapOrderByI[any, any, string]()
	default:
		panic(fmt.Errorf("orderBy unknown I type %v", typeParams[0].Kind()))
	}
}

func maxOf(typeParams []reflect.Type) reflect.Value {

	return reflect.ValueOf(func(o optic.Optic[any, any, any, any, any, any, any, any, any], cmpBy optic.Operation[any, any, any, any]) any {
		switch typeParams[0].Kind() {
		case reflect.Int:

			c := optic.Compose(
				optic.OpToOptic(cmpBy),
				optic.IsoCast[any, int](),
			)

			m := optic.MaxOf(o, c)

			return any(m)
		case reflect.Float64:
			c := optic.Compose(
				optic.OpToOptic(cmpBy),
				optic.IsoCast[any, float64](),
			)

			m := optic.MaxOf(o, c)

			return any(m)
		case reflect.String:
			c := optic.Compose(
				optic.OpToOptic(cmpBy),
				optic.IsoCast[any, string](),
			)

			m := optic.MaxOf(o, c)

			return any(m)

		default:
			panic("orderBy unknown type")
		}
	})
}

func minOf(typeParams []reflect.Type) reflect.Value {

	return reflect.ValueOf(func(o optic.Optic[any, any, any, any, any, any, any, any, any], cmpBy optic.Operation[any, any, any, any]) any {
		switch typeParams[0].Kind() {
		case reflect.Int:

			c := optic.Compose(
				optic.OpToOptic(cmpBy),
				optic.IsoCast[any, int](),
			)

			m := optic.MinOf(o, c)

			return any(m)
		case reflect.Float64:
			c := optic.Compose(
				optic.OpToOptic(cmpBy),
				optic.IsoCast[any, float64](),
			)

			m := optic.MinOf(o, c)

			return any(m)
		case reflect.String:
			c := optic.Compose(
				optic.OpToOptic(cmpBy),
				optic.IsoCast[any, string](),
			)

			m := optic.MinOf(o, c)

			return any(m)

		default:
			panic("orderBy unknown type")
		}
	})
}

func fieldLens(fref func(source any) any) optic.Optic[any, any, any, any, any, any, any, any, any] {

	return optic.UnsafeOmni[any, any, any, any, any, any, any, any, any](
		func(ctx context.Context, source any) (any, any, error) {
			fieldPtr := fref(&source)
			vField := reflect.ValueOf(fieldPtr)
			return optic.Void{}, vField.Elem().Interface(), nil
		},
		func(ctx context.Context, focus, source any) (any, error) {
			vSource := reflect.ValueOf(source)
			vPtr := reflect.New(vSource.Type())
			vPtr.Elem().Set(vSource)
			fieldPtr := fref(vPtr.Interface())

			vField := reflect.ValueOf(fieldPtr)
			vField.Elem().Set(reflect.ValueOf(focus))
			return vPtr.Elem().Interface(), nil
		},
		func(ctx context.Context, source any) optic.SeqIE[any, any] {
			return func(yield func(optic.ValueIE[any, any]) bool) {
				fieldPtr := fref(&source)
				vField := reflect.ValueOf(fieldPtr)
				yield(optic.ValIE[any, any](optic.Void{}, vField.Elem().Interface(), nil))
			}
		},
		func(ctx context.Context, source any) (int, error) {
			return 1, nil
		},
		func(ctx context.Context, fmap func(index any, focus any) (any, error), source any) (any, error) {
			vSource := reflect.ValueOf(source)
			vPtr := reflect.New(vSource.Type())
			vPtr.Elem().Set(vSource)
			fieldPtr := fref(vPtr.Interface())

			vField := reflect.ValueOf(fieldPtr)

			newVal, err := fmap(optic.Void{}, vField.Elem().Interface())

			coerced, ok := goscript.ReflectCoerce(nil, reflect.ValueOf(newVal), vField.Type().Elem())
			if !ok {
				panic("field lens coerce failed")
			}

			vField.Elem().Set(coerced)
			return vPtr.Elem().Interface(), err
		},
		func(ctx context.Context, index, source any) optic.SeqIE[any, any] {
			return func(yield func(optic.ValueIE[any, any]) bool) {
				fieldPtr := fref(&source)
				vField := reflect.ValueOf(fieldPtr)
				yield(optic.ValIE[any, any](optic.Void{}, vField.Elem().Interface(), nil))
			}
		},
		func(anyA, anyB any) bool {
			return true
		},
		func(ctx context.Context, focus any) (any, error) {
			panic("field lens reverse get not implemented")
		},
		nil,
		func() expr.OpticExpression {
			return expr.FieldLens{}
		},
	)

}

func anyfyedCollection[ERR any](v reflect.Value) optic.Collection[any, any, ERR] {

	return optic.ColIE[ERR, any, any](
		func(ctx context.Context) optic.SeqIE[any, any] {
			return func(yield func(optic.ValueIE[any, any]) bool) {

				asIter := v.MethodByName("AsIter")

				seq := asIter.Call(nil)[0].Call([]reflect.Value{
					reflect.ValueOf(ctx),
				})[0]

				yieldFnc := reflect.MakeFunc(seq.Type().In(0), func(args []reflect.Value) (results []reflect.Value) {

					ix := args[0].MethodByName("Index").Call(nil)[0]
					val := args[0].MethodByName("Value").Call(nil)[0]
					errV := args[0].MethodByName("Error").Call(nil)[0]

					var err error
					if !errV.IsNil() {
						err = errV.Interface().(error)
					}

					return []reflect.Value{reflect.ValueOf(yield(optic.ValIE(ix.Interface(), val.Interface(), err)))}
				})

				seq.Call([]reflect.Value{yieldFnc})
			}
		},
		func(ctx context.Context, index any) optic.SeqIE[any, any] {
			return func(yield func(optic.ValueIE[any, any]) bool) {

				asIter := v.MethodByName("AsIxGet")

				seq := asIter.Call(nil)[0].Call([]reflect.Value{
					reflect.ValueOf(ctx),
					reflect.ValueOf(index),
				})[0]

				yieldFnc := reflect.MakeFunc(seq.Type().In(0), func(args []reflect.Value) (results []reflect.Value) {

					ix := args[0].MethodByName("Index").Call(nil)[0]
					val := args[0].MethodByName("Value").Call(nil)[0]
					errV := args[0].MethodByName("Error").Call(nil)[0]

					var err error
					if !errV.IsNil() {
						err = errV.Interface().(error)
					}

					return []reflect.Value{reflect.ValueOf(yield(optic.ValIE(ix.Interface(), val.Interface(), err)))}
				})

				seq.Call([]reflect.Value{yieldFnc})
			}
		},
		func(a, b any) bool {
			asIxMatch := v.MethodByName("AsIxMatch")

			ret := asIxMatch.Call(nil)[0].Call([]reflect.Value{
				reflect.ValueOf(a),
				reflect.ValueOf(b),
			})

			return ret[0].Bool()
		},
		func(ctx context.Context) (int, error) {
			as := v.MethodByName("AsLengthGetter")

			ret := as.Call(nil)[0].Call([]reflect.Value{
				reflect.ValueOf(ctx),
			})

			var err error
			if !ret[1].IsNil() {
				err = ret[1].Interface().(error)
			}

			return int(ret[1].Int()), err
		},
	)

}

func deanyfyedCollection[I, E, ERR any](
	v reflect.Value,
) optic.Collection[I, E, ERR] {

	return optic.ColIE[ERR, I, E](
		func(ctx context.Context) optic.SeqIE[I, E] {
			return func(yield func(optic.ValueIE[I, E]) bool) {

				asIter := v.MethodByName("AsIter")

				seq := asIter.Call(nil)[0].Call([]reflect.Value{
					reflect.ValueOf(ctx),
				})[0]

				yieldFnc := reflect.MakeFunc(seq.Type().In(0), func(args []reflect.Value) (results []reflect.Value) {

					ix := args[0].MethodByName("Index").Call(nil)[0]
					val := args[0].MethodByName("Value").Call(nil)[0]
					errV := args[0].MethodByName("Error").Call(nil)[0]

					var err error
					if !errV.IsNil() {
						err = errV.Interface().(error)
					}

					return []reflect.Value{reflect.ValueOf(yield(optic.ValIE(Coerce(ix, reflect.TypeFor[I]()).Interface().(I), Coerce(val, reflect.TypeFor[E]()).Interface().(E), err)))}
				})

				seq.Call([]reflect.Value{yieldFnc})
			}
		},
		func(ctx context.Context, index I) optic.SeqIE[I, E] {
			return func(yield func(optic.ValueIE[I, E]) bool) {

				asIter := v.MethodByName("AsIxGet")

				seq := asIter.Call(nil)[0].Call([]reflect.Value{
					reflect.ValueOf(ctx),
					reflect.ValueOf(index),
				})[0]

				yieldFnc := reflect.MakeFunc(seq.Type().In(0), func(args []reflect.Value) (results []reflect.Value) {

					ix := args[0].MethodByName("Index").Call(nil)[0]
					val := args[0].MethodByName("Value").Call(nil)[0]
					errV := args[0].MethodByName("Error").Call(nil)[0]

					var err error
					if !errV.IsNil() {
						err = errV.Interface().(error)
					}

					return []reflect.Value{reflect.ValueOf(yield(optic.ValIE(Coerce(ix, reflect.TypeFor[I]()).Interface().(I), Coerce(val, reflect.TypeFor[E]()).Interface().(E), err)))}
				})

				seq.Call([]reflect.Value{yieldFnc})
			}
		},
		func(a, b I) bool {
			asIxMatch := v.MethodByName("AsIxMatch")

			ret := asIxMatch.Call(nil)[0].Call([]reflect.Value{
				reflect.ValueOf(a),
				reflect.ValueOf(b),
			})

			return ret[0].Bool()
		},
		func(ctx context.Context) (int, error) {
			as := v.MethodByName("AsLengthGetter")

			ret := as.Call(nil)[0].Call([]reflect.Value{
				reflect.ValueOf(ctx),
			})

			var err error
			if !ret[1].IsNil() {
				err = ret[1].Interface().(error)
			}

			return int(ret[1].Int()), err
		},
	)

}

func ixMapIso(
	ixmap func(left any, right any) any,
	ixmatch func(a, b any) bool,
	unmap func(mapped any) (any, bool, any, bool),
	exprDef optic.ExpressionDef,
) optic.Optic[any, lo.Tuple2[mo.Option[any], mo.Option[any]], lo.Tuple2[mo.Option[any], mo.Option[any]], any, any, any, any, any, any] {
	return optic.UnsafeReconstrain[any, any, any, any](optic.IxMapIso(ixmap, ixmatch, unmap, exprDef))
}

func init() {

	optic.DisableOpticTypeOptimizations = true
	goscript.RegisterTransformer(func(fromVal reflect.Value, toType reflect.Type) reflect.Value {

		if toType.Name() == "OpticRO[interface {},interface {},interface {},interface {},interface {},interface {},interface {}]" {
			return reflect.ValueOf(anyfyedReflect(fromVal).(optic.OpticRO[any, any, any, any, any, any, any]))
		}

		if toType.Name() == "Tuple2[float64,float64]" {

			a := fromVal.FieldByName("A")
			b := fromVal.FieldByName("B")

			afloat := Coerce(a, reflect.TypeFor[float64]())
			bfloat := Coerce(b, reflect.TypeFor[float64]())

			return reflect.ValueOf(lo.T2(afloat.Interface().(float64), bfloat.Interface().(float64)))
		}

		if toType.Name() == "Operation[github.com/samber/lo.Tuple2[interface {},interface {}],interface {},interface {},interface {}]" {

			x := anyfyedReflectP[optic.Void, lo.Tuple2[any, any], lo.Tuple2[any, any], any, any](
				fromVal,
				func(v reflect.Value) optic.Void {
					return optic.Void{}
				},
				func(v optic.Void) reflect.Value {
					return reflect.ValueOf(v)
				},
				func(v reflect.Value) lo.Tuple2[any, any] {
					return lo.T2(
						v.Field(0).Interface(),
						v.Field(1).Interface(),
					)
				},
				func(t lo.Tuple2[any, any]) reflect.Value {
					opticExpr := fromVal.MethodByName("AsExpr").Call(nil)[0].Interface().(expr.OpticExpression)
					sType := opticExpr.OpticS()
					iType := sType.Field(0).Type
					vType := sType.Field(1).Type
					retPtr := reflect.New(sType)
					setField(retPtr.Elem(), 0, Coerce(reflect.ValueOf(t.A), iType))
					setField(retPtr.Elem(), 1, Coerce(reflect.ValueOf(t.B), vType))

					return retPtr.Elem()
				},
				func(v reflect.Value) lo.Tuple2[any, any] {
					return lo.T2(
						v.Field(0).Interface(),
						v.Field(1).Interface(),
					)
				},
				func(t lo.Tuple2[any, any]) reflect.Value {
					opticExpr := fromVal.MethodByName("AsExpr").Call(nil)[0].Interface().(expr.OpticExpression)
					sType := opticExpr.OpticS()
					iType := sType.Field(0).Type
					vType := sType.Field(1).Type
					retPtr := reflect.New(sType)
					setField(retPtr.Elem(), 0, Coerce(reflect.ValueOf(t.A), iType))
					setField(retPtr.Elem(), 1, Coerce(reflect.ValueOf(t.B), vType))

					return retPtr.Elem()
				},
				reflect.Value.Interface,
				reflect.ValueOf,
				reflect.Value.Interface,
				reflect.ValueOf,
			)

			return reflect.ValueOf(x.(optic.Operation[lo.Tuple2[any, any], any, any, any]))

		}

		if toType.Name() == "Predicate[github.com/samber/lo.Tuple2[interface {},interface {}],interface {}]" {

			x := anyfyedReflectP[optic.Void, lo.Tuple2[any, any], lo.Tuple2[any, any], bool, bool](
				fromVal,
				func(v reflect.Value) optic.Void {
					return optic.Void{}
				},
				func(v optic.Void) reflect.Value {
					return reflect.ValueOf(v)
				},
				func(v reflect.Value) lo.Tuple2[any, any] {
					return lo.T2(
						v.Field(0).Interface(),
						v.Field(1).Interface(),
					)
				},
				func(t lo.Tuple2[any, any]) reflect.Value {
					opticExpr := fromVal.MethodByName("AsExpr").Call(nil)[0].Interface().(expr.OpticExpression)
					sType := opticExpr.OpticS()
					iType := sType.Field(0).Type
					vType := sType.Field(1).Type
					retPtr := reflect.New(sType)
					setField(retPtr.Elem(), 0, Coerce(reflect.ValueOf(t.A), iType))
					setField(retPtr.Elem(), 1, Coerce(reflect.ValueOf(t.B), vType))

					return retPtr.Elem()
				},
				func(v reflect.Value) lo.Tuple2[any, any] {
					return lo.T2(
						v.Field(0).Interface(),
						v.Field(1).Interface(),
					)
				},
				func(t lo.Tuple2[any, any]) reflect.Value {
					opticExpr := fromVal.MethodByName("AsExpr").Call(nil)[0].Interface().(expr.OpticExpression)
					sType := opticExpr.OpticS()
					iType := sType.Field(0).Type
					vType := sType.Field(1).Type
					retPtr := reflect.New(sType)
					setField(retPtr.Elem(), 0, Coerce(reflect.ValueOf(t.A), iType))
					setField(retPtr.Elem(), 1, Coerce(reflect.ValueOf(t.B), vType))

					return retPtr.Elem()
				},
				func(v reflect.Value) bool {
					return v.Bool()
				},
				func(b bool) reflect.Value {
					return reflect.ValueOf(b)
				},
				func(v reflect.Value) bool {
					return v.Bool()
				},
				func(b bool) reflect.Value {
					return reflect.ValueOf(b)
				},
			)

			return reflect.ValueOf(x.(optic.Predicate[lo.Tuple2[any, any], any]))

		}

		if toType.Name() == "Collection[interface {},interface {},github.com/spearson78/go-optic.Pure]" {
			return reflect.ValueOf(anyfyedCollection[optic.Pure](fromVal))
		}

		if toType.Name() == "Collection[int,github.com/spearson78/go-optic/internal/playground/data.Rating,github.com/spearson78/go-optic.Pure]" {
			return reflect.ValueOf(deanyfyedCollection[int, data.Rating, optic.Pure](fromVal))
		}

		if toType.Name() == "Collection[int,github.com/spearson78/go-optic/internal/playground/data.Comment,github.com/spearson78/go-optic.Pure]" {
			return reflect.ValueOf(deanyfyedCollection[int, data.Comment, optic.Pure](fromVal))
		}

		if toType.Name() == "Collection[interface {},interface {},interface {}]" {
			return reflect.ValueOf(anyfyedCollection[any](fromVal))
		}

		if toType.Name() == "Optic[interface {},interface {},interface {},interface {},interface {},interface {},interface {},interface {},interface {}]" {
			return reflect.ValueOf(anyfyedReflect(fromVal))
		}

		if toType.Name() == "Operation[interface {},interface {},interface {},interface {}]" {
			return reflect.ValueOf(anyfyedReflect(fromVal))
		}

		if toType.Name() == "PredicateI[interface {},interface {},interface {}]" {
			x := anyfyedReflectP[optic.Void, optic.ValueI[any, any], optic.ValueI[any, any], bool, bool](
				fromVal,
				func(v reflect.Value) optic.Void {
					return optic.Void{}
				},
				func(v optic.Void) reflect.Value {
					return reflect.ValueOf(v)
				},
				func(v reflect.Value) optic.ValueI[any, any] {
					return optic.ValI(
						v.Field(0).Interface(),
						v.Field(1).Interface(),
					)
				},
				func(vi optic.ValueI[any, any]) reflect.Value {
					opticExpr := fromVal.MethodByName("AsExpr").Call(nil)[0].Interface().(expr.OpticExpression)
					sType := opticExpr.OpticS()
					iType := sType.Field(0).Type
					vType := sType.Field(1).Type
					retPtr := reflect.New(sType)
					setField(retPtr.Elem(), 0, Coerce(reflect.ValueOf(vi.Index()), iType))
					setField(retPtr.Elem(), 1, Coerce(reflect.ValueOf(vi.Value()), vType))

					return retPtr.Elem()
				},
				func(v reflect.Value) optic.ValueI[any, any] {
					return optic.ValI(
						v.Field(0).Interface(),
						v.Field(1).Interface(),
					)
				},
				func(vi optic.ValueI[any, any]) reflect.Value {
					opticExpr := fromVal.MethodByName("AsExpr").Call(nil)[0].Interface().(expr.OpticExpression)
					sType := opticExpr.OpticS()
					iType := sType.Field(0).Type
					vType := sType.Field(1).Type
					retPtr := reflect.New(sType)
					setField(retPtr.Elem(), 0, Coerce(reflect.ValueOf(vi.Index()), iType))
					setField(retPtr.Elem(), 1, Coerce(reflect.ValueOf(vi.Value()), vType))

					return retPtr.Elem()
				},
				func(v reflect.Value) bool {
					return Coerce(v, reflect.TypeFor[bool]()).Interface().(bool)
				},
				func(b bool) reflect.Value {
					return reflect.ValueOf(b)
				},
				func(v reflect.Value) bool {
					return Coerce(v, reflect.TypeFor[bool]()).Interface().(bool)
				},
				func(b bool) reflect.Value {
					return reflect.ValueOf(b)
				},
			)

			return reflect.ValueOf(x.(optic.PredicateI[any, any, any]))
		}

		if toType.Name() == "Predicate[interface {},interface {}]" {
			x := anyfyedReflectP[optic.Void, any, any, bool, bool](
				fromVal,
				func(v reflect.Value) optic.Void {
					return optic.Void{}
				},
				func(v optic.Void) reflect.Value {
					return reflect.ValueOf(v)
				},
				reflect.Value.Interface,
				reflect.ValueOf,
				reflect.Value.Interface,
				reflect.ValueOf,
				func(v reflect.Value) bool {
					return Coerce(v, reflect.TypeFor[bool]()).Interface().(bool)
				},
				func(b bool) reflect.Value {
					return reflect.ValueOf(b)
				},
				func(v reflect.Value) bool {
					return Coerce(v, reflect.TypeFor[bool]()).Interface().(bool)
				},
				func(b bool) reflect.Value {
					return reflect.ValueOf(b)
				},
			)

			return reflect.ValueOf(x.(optic.Predicate[any, any]))
		}

		if toType.Name() == "OperationI[interface {},interface {},interface {},interface {},interface {}]" {

			x := anyfyedReflectP[optic.Void, optic.ValueI[any, any], optic.ValueI[any, any], any, any](
				fromVal,
				func(v reflect.Value) optic.Void {
					return optic.Void{}
				},
				func(v optic.Void) reflect.Value {
					return reflect.ValueOf(v)
				},
				func(v reflect.Value) optic.ValueI[any, any] {
					return optic.ValI(
						v.Field(0).Interface(),
						v.Field(1).Interface(),
					)
				},
				func(vi optic.ValueI[any, any]) reflect.Value {
					opticExpr := fromVal.MethodByName("AsExpr").Call(nil)[0].Interface().(expr.OpticExpression)
					sType := opticExpr.OpticS()
					iType := sType.Field(0).Type
					vType := sType.Field(1).Type
					retPtr := reflect.New(sType)
					setField(retPtr.Elem(), 0, Coerce(reflect.ValueOf(vi.Index()), iType))
					setField(retPtr.Elem(), 1, Coerce(reflect.ValueOf(vi.Value()), vType))

					return retPtr.Elem()
				},
				func(v reflect.Value) optic.ValueI[any, any] {
					return optic.ValI(
						v.Field(0).Interface(),
						v.Field(1).Interface(),
					)
				},
				func(vi optic.ValueI[any, any]) reflect.Value {
					opticExpr := fromVal.MethodByName("AsExpr").Call(nil)[0].Interface().(expr.OpticExpression)
					sType := opticExpr.OpticS()
					iType := sType.Field(0).Type
					vType := sType.Field(1).Type
					retPtr := reflect.New(sType)
					setField(retPtr.Elem(), 0, Coerce(reflect.ValueOf(vi.Index()), iType))
					setField(retPtr.Elem(), 1, Coerce(reflect.ValueOf(vi.Value()), vType))

					return retPtr.Elem()
				},
				reflect.Value.Interface,
				reflect.ValueOf,
				reflect.Value.Interface,
				reflect.ValueOf,
			)

			return reflect.ValueOf(x.(optic.OperationI[any, any, any, any, any]))

		}

		if toType.Name() == "Optic[interface {},interface {},interface {},github.com/spearson78/go-optic/internal/playground/data.BlogPost,github.com/spearson78/go-optic/internal/playground/data.BlogPost,interface {},interface {},interface {},interface {}]" {
			return reflect.ValueOf(coerceOptic[any, any, any, data.BlogPost, data.BlogPost](fromVal))
		}

		if toType.Name() == "ReductionP[interface {},interface {},interface {},interface {}]" {
			return reflect.ValueOf(anyfyedReflectReduction(fromVal))
		}

		if toType.Name() == "Reduction[interface {},interface {},interface {}]" {
			return reflect.ValueOf(anyfyedReflectReduction(fromVal))
		}

		if toType.Name() == "Operation[interface {},float64,interface {},interface {}]" {
			return reflect.ValueOf(coerceOptic[any, any, any, float64, float64](fromVal))
		}

		if toType.Name() == "Operation[interface {},string,interface {},interface {}]" {
			return reflect.ValueOf(coerceOptic[any, any, any, string, string](fromVal))
		}

		if toType.Name() == "OrderByPredicateI[interface {},interface {},interface {}]" {

			x := anyfyedReflectP[optic.Void, lo.Tuple2[optic.ValueI[any, any], optic.ValueI[any, any]], lo.Tuple2[optic.ValueI[any, any], optic.ValueI[any, any]], bool, bool](
				fromVal,
				func(v reflect.Value) optic.Void {
					return optic.Void{}
				},
				func(v optic.Void) reflect.Value {
					return reflect.ValueOf(v)
				},
				func(v reflect.Value) lo.Tuple2[optic.ValueI[any, any], optic.ValueI[any, any]] {
					return lo.T2(
						optic.ValI(
							getField(v.Field(0), 0).Interface(),
							getField(v.Field(0), 1).Interface(),
						),
						optic.ValI(
							getField(v.Field(1), 0).Interface(),
							getField(v.Field(1), 1).Interface(),
						),
					)
				},
				func(t lo.Tuple2[optic.ValueI[any, any], optic.ValueI[any, any]]) reflect.Value {
					return reflect.ValueOf(
						lo.T2(
							any(t.A),
							any(t.B),
						),
					)
				},
				func(v reflect.Value) lo.Tuple2[optic.ValueI[any, any], optic.ValueI[any, any]] {
					return lo.T2(
						optic.ValI(
							getField(v.Field(0), 0).Interface(),
							getField(v.Field(0), 1).Interface(),
						),
						optic.ValI(
							getField(v.Field(1), 0).Interface(),
							getField(v.Field(1), 1).Interface(),
						),
					)
				},
				func(t lo.Tuple2[optic.ValueI[any, any], optic.ValueI[any, any]]) reflect.Value {
					return reflect.ValueOf(
						lo.T2(
							any(t.A),
							any(t.B),
						),
					)
				},
				func(v reflect.Value) bool {
					return v.Bool()
				},
				func(b bool) reflect.Value {
					return reflect.ValueOf(b)
				},
				func(v reflect.Value) bool {
					return v.Bool()
				},
				func(b bool) reflect.Value {
					return reflect.ValueOf(b)
				},
			)

			return reflect.ValueOf(x.(optic.OrderByPredicateI[interface{}, interface{}, interface{}]))

		}

		if toType.Name() == "Tuple2[github.com/spearson78/go-optic.ValueI[interface {},interface {}],github.com/spearson78/go-optic.ValueI[interface {},interface {}]]" {
			return reflect.ValueOf(lo.T2(
				optic.ValI(
					fromVal.Field(0).Elem().MethodByName("Index").Call(nil)[0].Interface(),
					fromVal.Field(0).Elem().MethodByName("Value").Call(nil)[0].Interface(),
				),
				optic.ValI(
					fromVal.Field(1).Elem().MethodByName("Index").Call(nil)[0].Interface(),
					fromVal.Field(1).Elem().MethodByName("Value").Call(nil)[0].Interface(),
				),
			))
		}

		if toType.Name() == "Collection[int,interface {},github.com/spearson78/go-optic.Pure]" {
			return reflect.ValueOf(
				optic.ColIE[optic.Pure, int, any](
					func(ctx context.Context) optic.SeqIE[int, any] {
						return func(yield func(optic.ValueIE[int, any]) bool) {
							seq := fromVal.MethodByName("AsIter").Call(nil)[0].Call([]reflect.Value{reflect.ValueOf(ctx)})[0]

							yieldFnc := reflect.MakeFunc(seq.Type().In(0), func(args []reflect.Value) (results []reflect.Value) {

								ix := args[0].MethodByName("Index").Call(nil)[0]
								val := args[0].MethodByName("Value").Call(nil)[0]
								errV := args[0].MethodByName("Error").Call(nil)[0]

								var err error
								if !errV.IsNil() {
									err = errV.Interface().(error)
								}

								coerceIx, ok := goscript.ReflectCoerce(nil, ix, reflect.TypeFor[int]())
								if !ok {
									coerceIx = reflect.ValueOf(0)
									err = errors.Join(fmt.Errorf("coerce failed %v to int", ix.Kind()))
								}

								yieldVal := optic.ValIE(int(coerceIx.Int()), val.Interface(), err)

								return []reflect.Value{reflect.ValueOf(yield(yieldVal))}
							})

							seq.Call([]reflect.Value{yieldFnc})
						}
					},
					func(ctx context.Context, index int) optic.SeqIE[int, any] {
						return func(yield func(optic.ValueIE[int, any]) bool) {
							seq := fromVal.MethodByName("AsIxGet").Call(nil)[0].Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(index)})[0]

							yieldFnc := reflect.MakeFunc(seq.Type().In(0), func(args []reflect.Value) (results []reflect.Value) {

								ix := args[0].MethodByName("Index").Call(nil)[0]
								val := args[0].MethodByName("Value").Call(nil)[0]
								errV := args[0].MethodByName("Error").Call(nil)[0]

								var err error
								if !errV.IsNil() {
									err = errV.Interface().(error)
								}

								coerceIx, ok := goscript.ReflectCoerce(nil, ix, reflect.TypeFor[int]())
								if !ok {
									coerceIx = reflect.ValueOf(0)
									err = errors.Join(fmt.Errorf("coerce failed %v to int", ix.Kind()))
								}

								yieldVal := optic.ValIE(int(coerceIx.Int()), val.Interface(), err)

								return []reflect.Value{reflect.ValueOf(yield(yieldVal))}
							})

							seq.Call([]reflect.Value{yieldFnc})
						}
					},
					optic.IxMatchComparable[int](),
					func(ctx context.Context) (int, error) {
						res := fromVal.MethodByName("AsLengthGetter").Call(nil)[0].Call([]reflect.Value{reflect.ValueOf(ctx)})
						var err error
						if !res[1].IsNil() {
							err = res[1].Interface().(error)
						}

						return int(res[0].Int()), err
					},
				),
			)
		}

		return fromVal
	})

	goscript.RegisterEraser(func(fromVal reflect.Value) reflect.Value {

		if fromVal.Type().Name() == "Optic[interface {},github.com/samber/lo.Tuple2[github.com/samber/mo.Option[interface {}],github.com/samber/mo.Option[interface {}]],github.com/samber/lo.Tuple2[github.com/samber/mo.Option[interface {}],github.com/samber/mo.Option[interface {}]],interface {},interface {},interface {},interface {},interface {},interface {}]" {
			//IxMapper exclusion it's already anyfyed
			return fromVal
		}

		if strings.HasPrefix(fromVal.Type().Name(), "Optic[") && fromVal.Type().Name() != "Optic[interface {},interface {},interface {},interface {},interface {},interface {},interface {},interface {},interface {}]" {
			x := anyfyedReflect(fromVal)
			v := reflect.ValueOf(x).Convert(reflect.TypeFor[optic.Optic[any, any, any, any, any, any, any, any, any]]())
			return v
		}

		return fromVal
	})
}

/*
func main() {
	err := goscript.Run(context.Background(), `
package main

import (
	"fmt"
	"strings"

	. "github.com/spearson78/go-optic"
)

type Person struct {
	Name    string
	Age     int
	Hobbies []string
}

func main() {

	nameField := FieldLens(func(source *Person) *string { return &source.Name })
	ageField := FieldLens(func(source *Person) *int { return &source.Age })
	hobbiesField := FieldLens(func(source *Person) *[]string { return &source.Hobbies })

	data := Person{
		Name:    "Max Mustermann",
		Age:     46,
		//Hobbies: []string{"eating", "sleeping"},
	}

	name := MustGet(nameField, data)
	age := MustGet(ageField, data)
	//hobbies := MustGet(hobbiesField, data)

	fmt.Println(name)
	fmt.Println(age) //, hobbies)

	fmt.Println(data)
	olderPerson := MustSet(ageField, 47, data)
	fmt.Println(olderPerson)
	fmt.Println(data)

	//Output:
	//Max Mustermann 46 [eating sleeping]
	//{Max Mustermann 47 [eating sleeping]}
	//{Max Mustermann 46 [EATING SLEEPING]}
}
`)

	if err != nil {
		log.Fatal(err)
	}

}

*/
