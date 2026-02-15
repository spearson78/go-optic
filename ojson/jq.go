package ojson

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/samber/lo"
	"github.com/samber/mo"
	. "github.com/spearson78/go-optic"
	"github.com/spearson78/go-optic/expr"
	"github.com/spearson78/go-optic/otree"
)

func subseq(from, length int) Optic[any, any, any, any, any, ReturnMany, ReadWrite, UniDir, Err] {
	return RetM(Rw(Ud(EErr(ReIndexed(
		Coalesce(
			Compose3(DownCast[any, []any](), SubSlice[any](from, length), IsoCastE[[]any, any]()),
			Compose3(DownCast[any, string](), SubString(from, length), IsoCastE[string, any]()),
		),
		UpCast[Void, any](),
		EqT2[any](),
	)))))
}

// JqSlice returns an Optic that focuses the [SubSlice] or [SubString] of the source JSON object.
func JqSlice(from, length int) builder[any, Err] {
	x := subseq(from, length)
	return builder[any, Err]{
		Optic: x,
	}
}

func (b builder[I, ERR]) JqSlice(from, to int) builder[any, Err] {
	x := subseq(from, to)
	return builder[any, Err]{
		Optic: EErr(RetM(Rw(Ud(Compose(b.Optic, x))))),
	}
}

// JqObjectOf returns on Optic that focuses the [MapOf] the source and casts to [any]
func JqObjectOf[RET any, RW any, DIR any, ERR any](fields Optic[string, any, any, any, any, RET, RW, DIR, ERR], size int) Optic[Void, any, any, any, any, ReturnOne, RW, UniDir, Err] {
	return Ret1(RwL(Ud(EErr(Compose(MapOf(fields, size), IsoCastE[map[string]any, any]())))))
}

func iterObjects[J, ERR any](ctx context.Context, fieldVals []ValueI[string, Collection[J, any, ERR]], obj map[string]any, yield func(ValueIE[Void, any]) bool) bool {

	if len(fieldVals) == 0 {
		return yield(ValIE(Void{}, any(obj), nil))
	}

	fieldName := fieldVals[0].Index()
	fieldVals[0].Value().AsIter()(ctx)(func(val ValueIE[J, any]) bool {
		_, focus, err := val.Get()
		if err != nil {
			return yield(ValIE(Void{}, any(nil), err))
		}

		newObj := make(map[string]any, len(obj)+1)
		for k, v := range obj {
			newObj[k] = v
		}
		newObj[fieldName] = focus
		return iterObjects(ctx, fieldVals[1:], newObj, yield)
	})

	return true
}

// JqObjectsOf returns an [Iteration] that focuses JSON objects built from the given fields. If a field returns a sequence of multiple values then multiple objects will be focused.
func JqObjectsOf[RET any, RW any, DIR any, ERR any, J, T, B any](fields Optic[string, any, T, Collection[J, any, ERR], B, RET, RW, DIR, ERR]) Optic[Void, any, any, any, any, ReturnMany, ReadOnly, UniDir, ERR] {
	return CombiIteration[ReturnMany, ERR, Void, any, any, any, any](
		func(ctx context.Context, source any) SeqIE[Void, any] {
			return func(yield func(ValueIE[Void, any]) bool) {
				var fieldVals []ValueI[string, Collection[J, any, ERR]]

				fields.AsIter()(ctx, source)(func(val ValueIE[string, Collection[J, any, ERR]]) bool {
					index, focus, err := val.Get()
					if err != nil {
						return yield(ValIE(Void{}, any(nil), err))
					}

					fieldVals = append(fieldVals, ValI(index, focus))
					return true
				})

				iterObjects(ctx, fieldVals, map[string]any{}, yield)
			}
		},
		nil,
		nil,
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.JqObjectsOf{
				OpticTypeExpr: ot,
				Fields:        fields.AsExpr(),
			}
		}, fields),
	)
}

func pickMapPath(path []any, from map[string]any, into map[string]any) error {
	if len(path) == 0 {
		return nil
	}

	key, ok := path[0].(string)

	if !ok {
		return errors.New("Pick expected string path segment")
	}

	if f, ok := from[key]; ok {
		i, ok := into[key]
		if !ok {
			i = make(map[string]any)
		}

		switch t := f.(type) {
		case map[string]any:
			err := pickMapPath(path[1:], t, i.(map[string]any))
			if err != nil {
				return err
			}
			into[key] = i
		case []any:
			ret, err := pickSlicePath(path[1:], t, i.([]any))
			if err != nil {
				return err
			}
			into[key] = ret
		default:
			into[key] = t
		}
	} else {
		into[key] = nil
	}

	return nil
}

func pickSlicePath(path []any, from []any, into []any) ([]any, error) {
	if len(path) == 0 {
		return into, nil
	}

	key, ok := path[0].(int)

	if !ok {
		return nil, errors.New("Pick expected int path segment")
	}

	if key < len(from) {
		f := from[key]

		if key >= cap(into)-1 {
			clone := make([]any, key+1)
			copy(clone, into)
			into = clone
		}

		if key >= len(into)-1 {
			into = into[:key+1]
		}

		i := into[key]

		switch t := f.(type) {
		case map[string]any:
			err := pickMapPath(path[1:], t, i.(map[string]any))
			if err != nil {
				return nil, err
			}
			into[key] = i
		case []any:
			ret, err := pickSlicePath(path[1:], t, i.([]any))
			if err != nil {
				return nil, err
			}
			into[key] = ret
		default:
			into[key] = t
		}

		return into, nil
	} else {
		if key >= cap(into) {
			clone := make([]any, key)
			copy(clone, into)
			into = clone
		}

		return into, nil
	}

}

// JqPick returns an [Operation] that focuses a JSON object constructed from the leaf elements of the paths.
func JqPick(paths ...*otree.PathNode[any]) Optic[Void, any, any, any, any, ReturnOne, ReadOnly, UniDir, Err] {
	return CombiGetter[Err, Void, any, any, any, any](
		func(ctx context.Context, source any) (Void, any, error) {

			switch t := source.(type) {
			case map[string]any:
				into := make(map[string]any)
				for _, path := range paths {
					err := pickMapPath(path.Slice(), t, into)
					if err != nil {
						return Void{}, nil, err
					}
				}
				return Void{}, into, nil

			case []any:
				var into []any
				for _, path := range paths {
					ret, err := pickSlicePath(path.Slice(), t, into)
					if err != nil {
						return Void{}, nil, err
					}
					into = ret
				}
				return Void{}, into, nil
			default:
				return Void{}, nil, errors.New("Pick expected map or slice")
			}

		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {

			flatPaths := MustGet(
				SliceOf(
					Compose(
						TraverseSlice[*otree.PathNode[any]](),
						SliceOf(
							Reversed(
								otree.TraversePath[any](),
							),
							5,
						),
					),
					5,
				),
				paths,
			)

			return expr.JqPick{
				OpticTypeExpr: ot,
				Paths:         flatPaths,
			}
		}),
	)
}

func getCmpTypeRank(a any) (int, error) {
	if a == nil {
		return 10, nil
	}

	switch a.(type) {
	case bool:
		return 9, nil
	case float64:
		return 8, nil
	case string:
		return 7, nil
	case []any:
		return 6, nil
	case map[string]any:
		return 5, nil
	default:
		return 0, errors.New("CompareJson unsupported type")
	}

}

func cmpJson(ctx context.Context, aVal, bVal any) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	aRank, err := getCmpTypeRank(aVal)
	if err != nil {
		return false, err
	}

	bRank, err := getCmpTypeRank(bVal)
	if err != nil {
		return false, err
	}

	if aRank != bRank {
		return aRank > bRank, nil
	}

	if aVal == nil {
		return false, nil
	}

	switch a := aVal.(type) {
	case bool:
		b := bVal.(bool)

		if a {
			if b {
				return false, nil
			} else {
				return true, nil
			}
		} else {
			if b {
				return true, nil
			} else {
				return false, nil
			}
		}
	case float64:
		b := bVal.(float64)
		return a < b, nil
	case string:
		b := bVal.(string)
		return a < b, nil
	case []any:
		b := bVal.([]any)

		la := len(a)
		lb := len(b)
		l := max(la, lb)

		for i := 0; i < l; i++ {
			if i <= la-1 {
				if i <= lb-1 {
					cmp, err := cmpJson(ctx, aVal, bVal)
					if err != nil {
						return false, err
					}
					if cmp {
						return true, nil
					}
				} else {
					return false, err
				}
			} else {
				if i <= lb-1 {
					return true, nil
				} else {
					return false, nil
				}
			}
		}

		return false, nil
	case map[string]any:
		b := bVal.(map[string]any)

		sortedKeys := Compose3(WithIndex(TraverseMap[string, any]()), ValueIIndex[string, any](), UpCast[string, any]())

		aKeys, err := GetContext(ctx, SliceOf(sortedKeys, len(a)), a)
		if err != nil {
			return false, err
		}
		bKeys, err := GetContext(ctx, SliceOf(sortedKeys, len(b)), b)
		if err != nil {
			return false, err
		}

		eq, err := GetContext(ctx, EqDeepT2[[]any](), lo.T2(aKeys, bKeys))
		if err != nil {
			return false, err
		}

		if eq {
			for _, k := range aKeys {
				kStr := k.(string)
				cmp, err := cmpJson(ctx, a[kStr], b[kStr])
				if err != nil {
					return false, err
				}
				if cmp {
					return true, nil
				}
			}
			return false, nil
		} else {
			return cmpJson(ctx, aKeys, bKeys)
		}
	default:
		return false, errors.New("CompareJson unsupported type")
	}

}

// JqOrder returns an [OrderByPredicate] that implements the jq order by rules.
// See:
//   - https://jqlang.org/manual/#sort-sort_by
//   - JqOrderBy for a version that orders by a sub property.
func JqOrder() Optic[Void, lo.Tuple2[any, any], lo.Tuple2[any, any], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return CombiGetter[Pure, Void, lo.Tuple2[any, any], lo.Tuple2[any, any], bool, bool](
		func(ctx context.Context, source lo.Tuple2[any, any]) (Void, bool, error) {
			cmp, err := cmpJson(ctx, source.A, source.B)
			return Void{}, cmp, err
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.JqOrderBy{
				OpticTypeExpr: ot,
			}
		}),
	)
}

// JqOrderBy returns an [OrderByPredicate] that implements the jq order by rules on the source using the focus of the given optic as the order by key.
// See:
//   - https://jqlang.org/manual/#sort-sort_by
//   - JqOrder for a version that orders by the source.
func JqOrderBy[I, RET, RW, DIR, ERR any](o Optic[I, any, any, any, any, RET, RW, DIR, ERR]) Optic[Void, lo.Tuple2[any, any], lo.Tuple2[any, any], bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return Ret1(Ro(Ud(EErrL(Compose(
		EErrMerge(T2Of(
			FirstOrDefault(EErrR(Compose(T2A[any, any](), o)), nil),
			FirstOrDefault(EErrR(Compose(T2B[any, any](), o)), nil),
		)),
		JqOrder(),
	)))))
}

func containsJson(left any, right any) bool {

	switch rt := right.(type) {
	case string:
		lt, ok := left.(string)
		if !ok {
			return false
		}

		return strings.Contains(lt, rt)
	case []any:
		lt, ok := left.([]any)
		if !ok {
			return false
		}

		for _, rVal := range rt {
			match := false
			for _, lVal := range lt {
				if containsJson(lVal, rVal) {
					match = true
					break
				}
			}
			if !match {
				return false
			}
		}

		return true

	case map[string]any:
		lt, ok := left.(map[string]any)
		if !ok {
			return false
		}

		for rKey, rVal := range rt {

			lVal, ok := lt[rKey]
			if !ok {
				return false
			}

			if !containsJson(lVal, rVal) {
				return false
			}
		}

		return true

	default:
		return MustGet(EqDeepT2[any](), lo.T2(left, right))
	}

}

// JqContainsT2 returns an [BinaryOp] that focuses the result of the jq contains function
//
// See:
//   - https://jqlang.org/manual/#contains
//   - [JqContains] for a right unary version.
//   - [JqInside] for a left unary version.
func JqContainsT2() Optic[Void, lo.Tuple2[any, any], lo.Tuple2[any, any], bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return CombiGetter[Pure, Void, lo.Tuple2[any, any], lo.Tuple2[any, any], bool, bool](
		func(ctx context.Context, source lo.Tuple2[any, any]) (Void, bool, error) {
			return Void{}, containsJson(source.A, source.B), nil
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.JqContains{
				OpticTypeExpr: ot,
			}
		}),
	)
}

// JqContains returns an [Operation] that focuses the result of the jq contains function
//
// See:
//   - https://jqlang.org/manual/#contains
//   - [JqInside] for a left unary version.
func JqContains(right any) Optic[Void, any, any, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2Bind(JqContainsT2(), right)
}

// JqInside returns an [Operation] that focuses the result of the jq inside function
//
// See:
//   - https://jqlang.org/manual/#inside
//   - [JqContains] for a right unary version.
func JqInside(left any) Optic[Void, any, any, bool, bool, ReturnOne, ReadOnly, UniDir, Pure] {
	return OpT2Bind(Ret1(Ro(Ud(EPure(Compose(SwappedT2[any, any](), JqContainsT2()))))), left)
}

// JqIndices returns an [Operation] that focuses the result of the jq indices function
//
// Note: This function will panic if the regexp fails to compile.
//
// See:
//   - https://jqlang.org/manual/#indices
func JqIndices(b any) Optic[Void, any, any, any, any, ReturnMany, ReadOnly, UniDir, Pure] {

	switch bt := b.(type) {
	case string:

		x := RetM(Ro(Ud(EPure(
			Compose(
				Index(
					Compose4(
						WithIndex(
							Compose(
								DownCast[any, string](),
								MatchString(regexp.MustCompile(bt), -1),
							),
						),
						ValueIIndex[MatchIndex, string](),
						MatchIndexOffsets(),
						TraverseSlice[int](),
					),
					0,
				),
				UpCast[int, any](),
			)))))

		return x

	case []any:

		x := RetM(Ro(Ud(EPure(Compose4(
			DownCast[any, []any](),
			WithIndex(
				Filtered(
					Indexing(otree.TopDown(SliceTail[any]())),
					Compose3(
						SliceToCol[any](),
						SubCol[int, any](0, len(bt)),
						EqCol(ValCol(bt...), EqDeepT2[any]()),
					),
				),
			),
			ValueIIndex[int, []any](),
			UpCast[int, any](),
		)))))

		return x

	default:
		x := RetM(Ro(Ud(EPure(Compose(
			WithIndex(
				Filtered(
					Traverse(),
					Eq[any](bt),
				),
			),
			ValueIIndex[any, any](),
		)))))

		return x
	}

}

type breakError string

func (e breakError) Error() string {
	return string(e)
}

// JqLabel returns an Optic that will stop when the corresponding [JqBreak] is focused.
func JqLabel[I, S, A, RET, RW, DIR, ERR any](o Optic[I, S, S, A, A, RET, RW, DIR, ERR], name string) Optic[I, S, S, A, A, ReturnMany, ReadOnly, UniDir, Pure] {
	isBreak :=
		Compose(
			ErrorAs[breakError](),
			Eq(breakError(name)),
		)

	return Stop(o, isBreak)
}

// JqBreak returns an Optic that will break out to the named [JqLabel]
func JqBreak[S, A any](name string) Optic[Void, S, S, A, A, ReturnOne, ReadWrite, BiDir, Err] {
	return Error[S, A](breakError(name))
}

// ChildrenCol focuses the children of a JSON object as a [Collection].
func ChildrenCol() Optic[Void, any, any, Collection[any, any, Err], Collection[any, any, Err], ReturnOne, ReadWrite, UniDir, Err] {
	return CombiLens[ReadWrite, Err, Void, any, any, Collection[any, any, Err], Collection[any, any, Err]](
		func(ctx context.Context, source any) (Void, Collection[any, any, Err], error) {
			switch t := source.(type) {
			case map[string]any:
				return Void{}, ColIE[Err](
					func(ctx context.Context) SeqIE[any, any] {
						return func(yield func(ValueIE[any, any]) bool) {
							for k, v := range t {
								if !yield(ValIE(any(k), v, nil)) {
									return
								}
							}
						}
					},
					func(ctx context.Context, index any) SeqIE[any, any] {
						if i, ok := index.(string); ok {
							return func(yield func(ValueIE[any, any]) bool) {
								if v, ok := t[i]; ok {
									yield(ValIE(any(index), v, nil))
								}
							}
						} else {
							return func(yield func(ValueIE[any, any]) bool) {}
						}
					},
					IxMatchComparable[any](),
					func(ctx context.Context) (int, error) {
						return len(t), nil
					},
				), nil
			case []any:
				return Void{}, ColIE[Err](
					func(ctx context.Context) SeqIE[any, any] {
						return func(yield func(ValueIE[any, any]) bool) {
							for k, v := range t {
								if !yield(ValIE(any(k), v, nil)) {
									return
								}
							}
						}
					},
					func(ctx context.Context, index any) SeqIE[any, any] {
						if i, ok := index.(int); ok {
							return func(yield func(ValueIE[any, any]) bool) {
								if i >= 0 && i < len(t) {
									yield(ValIE(index, t[i], nil))
								}
							}
						} else {
							return func(yield func(ValueIE[any, any]) bool) {}
						}
					},
					IxMatchComparable[any](),
					func(ctx context.Context) (int, error) {
						return len(t), nil
					},
				), nil
			default:
				return Void{}, ColIE[Err](
					func(ctx context.Context) SeqIE[any, any] { return func(yield func(ValueIE[any, any]) bool) {} },
					func(ctx context.Context, index any) SeqIE[any, any] {
						return func(yield func(ValueIE[any, any]) bool) {}
					},
					IxMatchComparable[any](),
					func(ctx context.Context) (int, error) { return 0, nil },
				), nil
			}
		},
		func(ctx context.Context, focus Collection[any, any, Err], source any) (any, error) {

			switch t := source.(type) {
			case map[string]any:

				var retErr error
				ret := make(map[string]any, len(t))
				focus.AsIter()(ctx)(func(val ValueIE[any, any]) bool {
					index, focus, err := val.Get()
					if err != nil {
						retErr = err
						return false
					}
					if strIndex, ok := index.(string); ok {
						ret[strIndex] = focus
						return true
					} else {
						retErr = errors.New("expected string index")
						return false
					}
				})
				if retErr != nil {
					return nil, retErr
				}
				return ret, nil

			case []any:
				ret := make([]any, 0, len(t))
				var retErr error
				focus.AsIter()(ctx)(func(val ValueIE[any, any]) bool {
					_, focus, err := val.Get()
					if err != nil {
						retErr = err
						return false
					}
					ret = append(ret, focus)
					return true
				})
				if retErr != nil {
					return nil, retErr
				}
				return ret, nil
			default:
				return source, nil
			}
		},
		IxMatchVoid(),
		ExprDef(func(ot expr.OpticTypeExpr) expr.OpticExpression {
			return expr.Children{
				OpticTypeExpr: ot,
			}
		}),
	)
}

// JqMergeT2 focuses on the merged version of 2 JSON elements (object or slice)
func JqMergeT2() Optic[Void, lo.Tuple2[mo.Option[any], mo.Option[any]], lo.Tuple2[mo.Option[any], mo.Option[any]], mo.Option[any], mo.Option[any], ReturnOne, ReadOnly, UniDir, Err] {
	return EErr(otree.MergeTreeT2(
		Traverse(),
		ChildrenCol(),
	))
}
