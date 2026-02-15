package optic

import (
	"context"
	"errors"
	"fmt"
	"unsafe"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spearson78/go-optic/expr"
)

const enableComposeIdentityOptimization = true

// Compose returns an [Optic] composed of the 2 input optics.
//
// Composition combines the 2 optics such that the output of the left is connected to the input of the right with the result using the index of the right input.
// The composed optic is compatible with both view and modify actions.
//
// See:
//   - [ComposeLeft] for a version that returns an optic using the left index.
//   - [ComposeBoth] for a version that returns an optic using both indices combined together.
//   - [ComposeI] for a version that returns a custom index.
//   - [Compose3] for a version that takes 3 parameters.
//   - [Compose4] for a version that takes 4 parameters.
//   - [Compose5] for a version that takes 5 parameters.
//   - [Compose6] for a version that takes 6 parameters.
//   - [Compose7] for a version that takes 7 parameters.
//   - [Compose8] for a version that takes 8 parameters.
//   - [Compose9] for a version that takes 9 parameters.
//   - [Compose10] for a version that takes 10 parameters.
//
// If you are implementing a Combinator then it is recommended to use either a condensing combinator ([Ret1],[RetM]) or condensing compose e.g. [ComposeRet1] , [ComposeRetM]
func Compose[I, J, S, T, A, B, C, D, RETI, RETJ any, RWI, RWJ any, DIRI, DIRJ any, ERRI, ERRJ any](left Optic[I, S, T, A, B, RETI, RWI, DIRI, ERRI], right Optic[J, A, B, C, D, RETJ, RWJ, DIRJ, ERRJ]) Optic[J, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]] {

	if enableComposeIdentityOptimization {

		leftIsIdentity := left.OpticType() == expr.OpticTypeIdentity
		rightIsIdentity := right.OpticType() == expr.OpticTypeIdentity

		if leftIsIdentity {

			composeExpr := func() expr.OpticExpression {
				return expr.Compose{
					OpticTypeExpr: expr.NewOpticTypeExpr[J, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]](),
					Left:          left.AsExpr(),
					Right:         right.AsExpr(),
					IxMap:         IxMapRight[I, J](right.AsIxMatch()).AsExpr(),
				}
			}

			if rightIsIdentity {
				return leftAndRightIdentityCompose[J, S, T, A, B, C, D, RETI, RETJ, RWI, RWJ, DIRI, DIRJ, ERRI, ERRJ](composeExpr)
			} else {
				return leftIdentityCompose[J, S, T, A, B, C, D, RETI, RETJ, RWI, RWJ, DIRI, DIRJ, ERRI, ERRJ](composeExpr, right)
			}
		} else {
			if rightIsIdentity {
				composeExpr := func() expr.OpticExpression {
					return expr.Compose{
						OpticTypeExpr: expr.NewOpticTypeExpr[J, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]](),
						Left:          left.AsExpr(),
						Right:         right.AsExpr(),
						IxMap:         IxMapRight[I, J](right.AsIxMatch()).AsExpr(),
					}
				}

				return rightIdentityCompose[I, J, S, T, A, B, C, D, RETI, RETJ, RWI, RWJ, DIRI, DIRJ, ERRI, ERRJ](composeExpr, left)
			} else {
				//Fall through to normal composition logic
			}
		}
	}

	return composeI(
		IxMapRight[I, J](right.AsIxMatch()),
		left,
		right,
	)
}

// ComposeLeft returns an [Optic] composed of the 2 input optics using the left index in the result.
//
// See:
//   - [Compose] for a the general version that returns an optic using the right index.
//   - [ComposeBoth] for a version that returns an optic using both indices combined together.
//   - [ComposeI] for a version that returns a custom index.
//
// If you are implementing a Combinator then it is recommended to use either a condensing combinator ([Ret1],[RetM]) or condensing compose e.g. [ComposeRet1] , [ComposeRetM]
func ComposeLeft[I, J, S, T, A, B, C, D, RETI, RETJ, RWI, RWJ, DIRI, DIRJ, ERRI, ERRJ any](left Optic[I, S, T, A, B, RETI, RWI, DIRI, ERRI], right Optic[J, A, B, C, D, RETJ, RWJ, DIRJ, ERRJ]) Optic[I, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]] {

	if enableComposeIdentityOptimization {

		leftIsIdentity := left.OpticType() == expr.OpticTypeIdentity
		rightIsIdentity := right.OpticType() == expr.OpticTypeIdentity

		if leftIsIdentity {

			composeExpr := func() expr.OpticExpression {
				return expr.Compose{
					OpticTypeExpr: expr.NewOpticTypeExpr[I, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]](),
					Left:          left.AsExpr(),
					Right:         right.AsExpr(),
					IxMap:         IxMapLeft[I, J](left.AsIxMatch()).AsExpr(),
				}
			}

			if rightIsIdentity {
				return leftAndRightIdentityCompose[I, S, T, A, B, C, D, RETI, RETJ, RWI, RWJ, DIRI, DIRJ, ERRI, ERRJ](composeExpr)
			} else {
				return leftIdentityComposeLeft[I, J, S, T, A, B, C, D, RETI, RETJ, RWI, RWJ, DIRI, DIRJ, ERRI, ERRJ](composeExpr, right)
			}
		} else {
			if rightIsIdentity {
				composeExpr := func() expr.OpticExpression {
					return expr.Compose{
						OpticTypeExpr: expr.NewOpticTypeExpr[I, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]](),
						Left:          left.AsExpr(),
						Right:         right.AsExpr(),
						IxMap:         IxMapLeft[I, J](left.AsIxMatch()).AsExpr(),
					}
				}
				return rightIdentityComposeLeft[I, S, T, A, B, C, D, RETI, RETJ, RWI, RWJ, DIRI, DIRJ, ERRI, ERRJ](composeExpr, left)
			} else {
				//Fall through to normal composition logic
			}
		}
	}

	return composeI(
		IxMapLeft[I, J](left.AsIxMatch()),
		left,
		right,
	)
}

// ComposeBoth returns an [Optic] composed of the 2 input optics using a combined left and right index in the result.
//
// See:
//   - [Compose] for a the general version that returns an optic using the right index.
//   - [ComposeLeft] for a version that returns an optic using th left index.
//   - [ComposeI] for a version that returns a custom index.
//
// If you are implementing a Combinator then it is recommended to use either a condensing combinator ([Ret1],[RetM]) or condensing compose e.g. [ComposeRet1] , [ComposeRetM]
func ComposeBoth[I, J, S, T, A, B, C, D, RETI, RETJ, RWI, RWJ, DIRI, DIRJ, ERRI, ERRJ any](o1 Optic[I, S, T, A, B, RETI, RWI, DIRI, ERRI], o2 Optic[J, A, B, C, D, RETJ, RWJ, DIRJ, ERRJ]) Optic[lo.Tuple2[I, J], S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]] {
	return composeI(
		IxMapBoth[I, J](o1.AsIxMatch(), o2.AsIxMatch()),
		o1,
		o2,
	)
}

type IxMapper[LEFT, RIGHT, MAPPED any, RET TReturnOne, RW any, DIR any, ERR TPure] Optic[MAPPED, lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]], lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]], MAPPED, MAPPED, RET, RW, DIR, ERR]

// ComposeI returns an [Optic] composed of the 2 input optics using a custom combination of left and right indices in the result.
//
// See:
//   - [Compose] for a the general version that returns an optic using the right index.
//   - [ComposeLeft] for a version that returns an optic using the left index.
//   - [ComposeBoth] for a version that returns an optic using both indices combined together as a [lo.Tuple2]
//
// If you are implementing a Combinator then it is recommended to use either a condensing combinator ([Ret1],[RetM]) or condensing compose e.g. [ComposeRet1] , [ComposeRetM]
func ComposeI[I, J, K, S, T, A, B, C, D any, RETI TReturnOne, RETL, RETR, RWI, RWL, RWR, DIRI, DIRL, DIRR, ERRL, ERRR any, ERRI TPure](ixmap IxMapper[I, J, K, RETI, RWI, DIRI, ERRI], o1 Optic[I, S, T, A, B, RETL, RWL, DIRL, ERRL], o2 Optic[J, A, B, C, D, RETR, RWR, DIRR, ERRR]) Optic[K, S, T, C, D, CompositionTree[RETL, RETR], CompositionTree[RWL, RWR], CompositionTree[DIRL, DIRR], CompositionTree[ERRL, ERRR]] {
	return composeI(ixmap, o1, o2)
}

// JoinCtxErr returns the given error joined with ctx.Err()
//
// This function is intended to be used in custom action implementations to enable the reporting of deadline exceeded errors.
func JoinCtxErr(ctx context.Context, err error) error {
	ctxErr := ctx.Err()
	if ctxErr != nil && !errors.Is(err, ctxErr) {
		if err != nil {
			return errors.Join(err, ctxErr)
		} else {
			return ctxErr
		}
	}
	return err
}

func composeI[I, J, K, S, T, A, B, C, D any, RETI TReturnOne, RETL, RETR, RWI, RWL, RWR, DIRI, DIRL, DIRR, ERRL, ERRR any, ERRI TPure](ixmap IxMapper[I, J, K, RETI, RWI, DIRI, ERRI], o1 Optic[I, S, T, A, B, RETL, RWL, DIRL, ERRL], o2 Optic[J, A, B, C, D, RETR, RWR, DIRR, ERRR]) Optic[K, S, T, C, D, CompositionTree[RETL, RETR], CompositionTree[RWL, RWR], CompositionTree[DIRL, DIRR], CompositionTree[ERRL, ERRR]] {
	return composeReconstrainI[CompositionTree[RETL, RETR], CompositionTree[RWL, RWR], CompositionTree[DIRL, DIRR], CompositionTree[ERRL, ERRR]](ixmap, o1, o2)
}

func ixGet2Yield[I, J, K, C any](ctx context.Context, index K, o1Index I, yield func(ValueIE[K, C]) bool, cont *bool, mapix func(ctx context.Context, left I, right J) (K, error), ixmatch func(a, b K) bool, o1Expr func() expr.OpticExpression) func(val ValueIE[J, C]) bool {
	return func(val ValueIE[J, C]) bool {
		o2Index, o2Focus, o2Err := val.Get()
		o2Err = JoinCtxErr(ctx, o2Err)
		k, err := mapix(ctx, o1Index, o2Index)
		err = OpticError(errors.Join(JoinCtxErr(ctx, err), o2Err), o1Expr)
		if err != nil {
			*cont = yield(ValIE(k, o2Focus, err))
			return *cont
		}

		match := ixmatch(index, k)
		if match {
			*cont = yield(ValIE(k, o2Focus, nil))
			return *cont
		}
		return true
	}
}

func composeReconstrainI[RETC, RWC, DIRC, ERRC, I, J, K, S, T, A, B, C, D any, RETI TReturnOne, RETL, RETR, RWI, RWL, RWR, DIRI, DIRL, DIRR, ERRL, ERRR any, ERRI TPure](ixmap IxMapper[I, J, K, RETI, RWI, DIRI, ERRI], o1 Optic[I, S, T, A, B, RETL, RWL, DIRL, ERRL], o2 Optic[J, A, B, C, D, RETR, RWR, DIRR, ERRR]) Optic[K, S, T, C, D, RETC, RWC, DIRC, ERRC] {

	//WARNING: the optimized functions need to manualls re-add the OpticError that they have optimized away the nested call for.

	composeExpr := func() expr.OpticExpression {
		return expr.Compose{
			OpticTypeExpr: expr.NewOpticTypeExpr[K, S, T, C, D, RETC, RWC, DIRC, ERRC](),
			Left:          o1.AsExpr(),
			Right:         o2.AsExpr(),
			IxMap:         ixmap.AsExpr(),
		}
	}

	o1Get := o1.AsGetter()
	o2Get := o2.AsGetter()

	o1IxGet := o1.AsIxGetter()
	o2IxGet := o2.AsIxGetter()

	o1Iter := o1.AsIter()
	o2Iter := o2.AsIter()

	o1Modify := o1.AsModify()
	o2Modify := o2.AsModify()

	o2LengthGetter := o2.AsLengthGetter()

	o1Setter := o1.AsSetter()
	o2Setter := o2.AsSetter()

	o1ReverseGetter := o1.AsReverseGetter()
	o2ReverseGetter := o2.AsReverseGetter()

	o1Expr := o1.AsExpr

	unmapIx := func(ctx context.Context, index K) (I, bool, J, bool, error) {
		var leftIndex I
		var leftOk = false
		var rightIndex J
		var rightOk = false
		if ixmap.OpticType()&expr.OpticTypeBiDirFlag != 0 {
			mapped, err := ixmap.AsReverseGetter()(ctx, index)
			if err != nil {
				return leftIndex, false, rightIndex, false, err
			}
			leftIndex, leftOk = mapped.A.Get()
			rightIndex, rightOk = mapped.B.Get()
		}
		return leftIndex, leftOk, rightIndex, rightOk, nil
	}

	mapIx := func(ctx context.Context, left I, right J) (K, error) {
		_, k, err := ixmap.AsGetter()(ctx, lo.T2(mo.Some(left), mo.Some(right)))
		return k, err
	}

	ixMatchFnc := ixmap.AsIxMatch()

	leftReturnOne := false
	rightReturnOne := false

	if !DisableOpticTypeOptimizations {
		leftReturnOne = !isType(o1.OpticType(), expr.OpticTypeReturnManyFlag)
		rightReturnOne = !isType(o2.OpticType(), expr.OpticTypeReturnManyFlag)
	}

	var ixgetter func(ctx context.Context, index K, source S) SeqIE[K, C]

	if leftReturnOne {
		if rightReturnOne {
			ixgetter = func(ctx context.Context, index K, source S) SeqIE[K, C] {
				return func(yield func(ValueIE[K, C]) bool) {
					indexI, a, focusErr := o1Get(ctx, source)
					if focusErr != nil {
						var k K
						var c C
						yield(ValIE(k, c, JoinCtxErr(ctx, focusErr)))
						return
					}

					indexJ, c, err := o2Get(ctx, a)
					mergedIndex, errIx := mapIx(ctx, indexI, indexJ)
					err = JoinCtxErr(ctx, errors.Join(err, errIx))
					err = OpticError(err, o1Expr)
					if err != nil {
						yield(ValIE(mergedIndex, c, err))
						return
					}

					if ixMatchFnc(index, mergedIndex) {
						yield(ValIE(mergedIndex, c, err))
						return
					}
				}
			}

		} else {

			ixgetter = func(ctx context.Context, index K, source S) SeqIE[K, C] {
				return func(yield func(ValueIE[K, C]) bool) {
					cont := true
					_, _, rightIndex, rightOk, err := unmapIx(ctx, index)
					if err != nil {
						var c C
						yield(ValIE(index, c, err))
						return
					}

					o1Index, o1Focus, focusErr := o1Get(ctx, source)
					if focusErr != nil {
						var k K
						var c C
						yield(ValIE(k, c, JoinCtxErr(ctx, focusErr)))
						return
					}

					o2Yield := ixGet2Yield(ctx, index, o1Index, yield, &cont, mapIx, ixMatchFnc, o1Expr)

					if rightOk {
						o2IxGet(ctx, rightIndex, o1Focus)(o2Yield)
					} else {
						o2Iter(ctx, o1Focus)(o2Yield)
					}
				}
			}

		}
	} else {
		if rightReturnOne {

			ixgetter = func(ctx context.Context, index K, source S) SeqIE[K, C] {
				return func(yield func(ValueIE[K, C]) bool) {
					cont := true
					leftIndex, leftOk, _, _, err := unmapIx(ctx, index)
					if err != nil {
						var c C
						yield(ValIE(index, c, err))
						return
					}

					if leftOk {
						o1IxGet(ctx, leftIndex, source)(func(val ValueIE[I, A]) bool {
							o1Index, o1Focus, o1Err := val.Get()
							o1Err = JoinCtxErr(ctx, o1Err)
							if o1Err != nil {
								var c C
								cont = yield(ValIE(index, c, o1Err))
								return cont
							}

							indexJ, c, err := o2Get(ctx, o1Focus)
							mergedIndex, errIx := mapIx(ctx, o1Index, indexJ)
							err = JoinCtxErr(ctx, errors.Join(err, errIx))
							err = OpticError(err, o1Expr)
							if err != nil {
								cont = yield(ValIE(mergedIndex, c, err))
								return cont
							}

							if ixMatchFnc(index, mergedIndex) {
								cont = yield(ValIE(mergedIndex, c, err))
								return cont
							}

							return cont
						})
					} else {
						o1Iter(ctx, source)(func(val ValueIE[I, A]) bool {
							o1Index, o1Focus, o1Err := val.Get()
							o1Err = JoinCtxErr(ctx, o1Err)
							if o1Err != nil {
								var c C
								cont = yield(ValIE(index, c, o1Err))
								return cont
							}

							indexJ, c, err := o2Get(ctx, o1Focus)
							mergedIndex, errIx := mapIx(ctx, o1Index, indexJ)
							err = JoinCtxErr(ctx, errors.Join(err, errIx))
							err = OpticError(err, o1Expr)
							if err != nil {
								cont = yield(ValIE(mergedIndex, c, err))
								return cont
							}

							if ixMatchFnc(index, mergedIndex) {
								cont = yield(ValIE(mergedIndex, c, err))
								return cont
							}

							return cont
						})
					}
				}
			}

		} else {

			ixgetter = func(ctx context.Context, index K, source S) SeqIE[K, C] {
				return func(yield func(ValueIE[K, C]) bool) {
					cont := true
					leftIndex, leftOk, rightIndex, rightOk, err := unmapIx(ctx, index)
					if err != nil {
						var c C
						yield(ValIE(index, c, err))
						return
					}

					if leftOk {
						o1IxGet(ctx, leftIndex, source)(func(val ValueIE[I, A]) bool {
							o1Index, o1Focus, o1Err := val.Get()
							o1Err = JoinCtxErr(ctx, o1Err)
							if o1Err != nil {
								var c C
								cont = yield(ValIE(index, c, o1Err))
								return cont
							}

							o2Yield := ixGet2Yield(ctx, index, o1Index, yield, &cont, mapIx, ixMatchFnc, o1Expr)

							if rightOk {
								o2IxGet(ctx, rightIndex, o1Focus)(o2Yield)
							} else {
								o2Iter(ctx, o1Focus)(o2Yield)
							}

							return cont
						})
					} else {
						o1Iter(ctx, source)(func(val ValueIE[I, A]) bool {
							o1Index, o1Focus, o1Err := val.Get()
							o1Err = JoinCtxErr(ctx, o1Err)
							if o1Err != nil {
								var c C
								cont = yield(ValIE(index, c, o1Err))
								return cont
							}

							o2Yield := ixGet2Yield(ctx, index, o1Index, yield, &cont, mapIx, ixMatchFnc, o1Expr)

							if rightOk {
								o2IxGet(ctx, rightIndex, o1Focus)(o2Yield)
							} else {
								o2Iter(ctx, o1Focus)(o2Yield)
							}

							return cont
						})
					}
				}
			}
		}

	}

	var getFnc func(ctx context.Context, s S) (retK K, retC C, retErr error)
	if leftReturnOne {
		getFnc = func(ctx context.Context, s S) (retK K, retC C, retErr error) {
			i, a, err := o1Get(ctx, s)
			if err != nil {
				retErr = err
				return
			}

			//We only have one a value whatever o2get returns thats our result
			j, c, err := o2Get(ctx, a)
			k, errIx := mapIx(ctx, i, j)
			err = JoinCtxErr(ctx, errors.Join(OpticError(err, o1Expr), errIx))
			if err != nil {
				retErr = err
				return
			}

			retK = k
			retC = c
			retErr = nil
			return
		}
	} else {
		if rightReturnOne {

			getFnc = func(ctx context.Context, s S) (retK K, retC C, retErr error) {
				//o2get always returns a result we can pass the result of o1Get into it and will get 1 result back.
				i, a, err := o1Get(ctx, s)
				if err != nil {
					retErr = err
					return
				}

				j, c, err := o2Get(ctx, a)
				k, errIx := mapIx(ctx, i, j)
				err = JoinCtxErr(ctx, errors.Join(OpticError(err, o1Expr), errIx))
				if err != nil {
					retErr = err
					return
				}

				retK = k
				retC = c
				retErr = nil
				return
			}
		} else {
			getFnc = func(ctx context.Context, s S) (retK K, retC C, retErr error) {
				retErr = ErrEmptyGet
				f1Seq := o1Iter(ctx, s)
				f1Seq(func(val ValueIE[I, A]) bool {
					indexI, a, focusErr := val.Get()
					if focusErr != nil {
						retErr = JoinCtxErr(ctx, focusErr)
						return false
					}
					cont := true
					f2Seq := o2Iter(ctx, a)
					f2Seq(func(val ValueIE[J, C]) bool {
						indexJ, c, err := val.Get()
						mergedIndex, errIx := mapIx(ctx, indexI, indexJ)
						retK = mergedIndex
						retC = c
						retErr = OpticError(errors.Join(err, errIx), o1Expr)
						cont = false
						return false
					})
					return cont
				})
				return
			}
		}
	}

	var setFnc func(ctx context.Context, val D, source S) (T, error)
	if leftReturnOne {
		setFnc = func(ctx context.Context, val D, source S) (T, error) {

			//Only 1 value we can get it and avoid allocating an fmap functin
			_, a, err := o1Get(ctx, source)
			if err != nil {
				var t T
				return t, err
			}

			//Whether o2 is ReturnOne or ReturnMany set will set all the values.
			b, err := o2Setter(ctx, val, a)
			if err != nil {
				var t T
				return t, OpticError(err, o1Expr)
			}

			return o1Setter(ctx, b, source)
		}

	} else {
		setFnc = func(ctx context.Context, val D, source S) (T, error) {
			//right prefers setter to modify
			o1ret, err := o1Modify(ctx, func(leftIndex I, leftFocus A) (B, error) {
				//Whether o2 is ReturnOne or ReturnMany set will set all the values.
				o2ret, err := o2Setter(ctx, val, leftFocus)
				return o2ret, JoinCtxErr(ctx, err)
			}, source)
			return o1ret, JoinCtxErr(ctx, err)
		}
	}

	var iterFnc func(ctx context.Context, source S) SeqIE[K, C]
	if leftReturnOne {
		if rightReturnOne {
			iterFnc = func(ctx context.Context, source S) SeqIE[K, C] {
				return func(yield func(ValueIE[K, C]) bool) {

					indexI, a, focusErr := o1Get(ctx, source)

					if focusErr != nil {
						var k K
						var c C
						yield(ValIE(k, c, JoinCtxErr(ctx, focusErr)))
						return
					}

					indexJ, c, err := o2Get(ctx, a)
					mergedIndex, errIx := mapIx(ctx, indexI, indexJ)
					err = JoinCtxErr(ctx, errors.Join(err, errIx))
					err = OpticError(err, o1Expr)
					yield(ValIE(mergedIndex, c, err))
				}
			}
		} else {

			iterFnc = func(ctx context.Context, source S) SeqIE[K, C] {
				return func(yield func(ValueIE[K, C]) bool) {

					indexI, a, focusErr := o1Get(ctx, source)

					if focusErr != nil {
						var k K
						var c C
						yield(ValIE(k, c, JoinCtxErr(ctx, focusErr)))
						return
					}

					f2Seq := o2Iter(ctx, a)
					f2Seq(func(val ValueIE[J, C]) bool {
						indexJ, c, err := val.Get()
						mergedIndex, errIx := mapIx(ctx, indexI, indexJ)
						err = OpticError(errors.Join(err, errIx), o1Expr)
						return yield(ValIE(mergedIndex, c, JoinCtxErr(ctx, err)))
					})
				}
			}

		}
	} else {

		if rightReturnOne {

			iterFnc = func(ctx context.Context, source S) SeqIE[K, C] {
				return func(yield func(ValueIE[K, C]) bool) {

					f1Seq := o1Iter(ctx, source)
					f1Seq(func(val ValueIE[I, A]) bool {
						indexI, a, focusErr := val.Get()
						if focusErr != nil {
							var k K
							var c C
							return yield(ValIE(k, c, JoinCtxErr(ctx, focusErr)))
						}

						indexJ, c, err := o2Get(ctx, a)
						mergedIndex, errIx := mapIx(ctx, indexI, indexJ)
						err = JoinCtxErr(ctx, errors.Join(err, errIx))
						err = OpticError(err, o1Expr)
						return yield(ValIE(mergedIndex, c, err))
					})

				}
			}

		} else {

			iterFnc = func(ctx context.Context, source S) SeqIE[K, C] {
				return func(yield func(ValueIE[K, C]) bool) {

					cont := true
					f1Seq := o1Iter(ctx, source)
					f1Seq(func(val ValueIE[I, A]) bool {
						indexI, a, focusErr := val.Get()
						if focusErr != nil {
							var k K
							var c C
							cont = yield(ValIE(k, c, JoinCtxErr(ctx, focusErr)))
							return cont
						}

						f2Seq := o2Iter(ctx, a)
						f2Seq(func(val ValueIE[J, C]) bool {
							indexJ, c, err := val.Get()
							mergedIndex, errIx := mapIx(ctx, indexI, indexJ)
							err = OpticError(errors.Join(err, errIx), o1Expr)
							cont = yield(ValIE(mergedIndex, c, JoinCtxErr(ctx, err)))
							return cont
						})
						return cont

					})

				}
			}

		}

	}

	var lengthGetter func(ctx context.Context, source S) (length int, retErr error)
	if leftReturnOne {
		lengthGetter = func(ctx context.Context, source S) (length int, retErr error) {
			_, a, err := o1Get(ctx, source)
			if err != nil {
				return 0, err
			}

			o2Len, err := o2LengthGetter(ctx, a)
			err = JoinCtxErr(ctx, OpticError(err, o1Expr))
			return o2Len, err
		}

	} else {
		lengthGetter = func(ctx context.Context, source S) (length int, retErr error) {
			f1Seq := o1Iter(ctx, source)
			f1Seq(func(val ValueIE[I, A]) bool {
				_, a, focusErr := val.Get()
				if focusErr != nil {
					retErr = focusErr
					return false
				}

				o2Len, err := o2LengthGetter(ctx, a)
				err = JoinCtxErr(ctx, err)
				if err != nil {
					retErr = OpticError(err, o1Expr)
					return false
				}
				length += o2Len

				return true
			})

			return
		}
	}

	var modifyFnc func(ctx context.Context, fmap func(indexK K, focus C) (D, error), source S) (T, error)
	if leftReturnOne {
		if rightReturnOne {
			modifyFnc = func(ctx context.Context, fmap func(indexK K, focus C) (D, error), source S) (T, error) {
				k, c, err := getFnc(ctx, source)
				if err != nil {
					var t T
					return t, err
				}

				mapped, err := fmap(k, c)
				if err != nil {
					var t T
					return t, OpticError(OpticError(err, o2.AsExpr), o1Expr)
				}

				ret, err := setFnc(ctx, mapped, source)
				if err != nil {
					var t T
					return t, err
				}
				return ret, err
			}
		} else {
			modifyFnc = func(ctx context.Context, fmap func(indexK K, focus C) (D, error), source S) (T, error) {
				i, a, err := o1Get(ctx, source)
				if err != nil {
					var t T
					return t, err
				}

				b, err := o2Modify(ctx, func(j J, c C) (D, error) {
					k, err := mapIx(ctx, i, j)
					if err != nil {
						var d D
						return d, OpticError(err, o1Expr)
					}

					ret, err := fmap(k, c)
					if err != nil {
						var d D
						return d, OpticError(err, o1Expr)
					}
					return ret, err
				}, a)
				if err != nil {
					var t T
					return t, OpticError(err, o1Expr)
				}

				ret, err := o1Setter(ctx, b, source)
				if err != nil {
					var t T
					return t, err
				}
				return ret, err
			}
		}
	} else {
		if rightReturnOne {

			modifyFnc = func(ctx context.Context, fmap func(indexK K, focus C) (D, error), source S) (T, error) {
				o1ret, err := o1Modify(ctx, func(leftIndex I, leftFocus A) (B, error) {

					j, c, err := o2Get(ctx, leftFocus)
					if err != nil {
						var b B
						return b, err
					}

					k, err := mapIx(ctx, leftIndex, j)
					if err != nil {
						var b B
						return b, OpticError(err, o2.AsExpr)
					}

					d, err := fmap(k, c)
					if err != nil {
						var b B
						return b, OpticError(err, o2.AsExpr)
					}

					ret, err := o2Setter(ctx, d, leftFocus)
					if err != nil {
						var b B
						return b, err
					}
					return ret, err
				}, source)
				return o1ret, JoinCtxErr(ctx, err)
			}

		} else {
			modifyFnc = func(ctx context.Context, fmap func(indexK K, focus C) (D, error), source S) (T, error) {
				o1ret, err := o1Modify(ctx, func(leftIndex I, leftFocus A) (B, error) {
					o2ret, err := o2Modify(ctx, func(rightIndex J, rightFocus C) (D, error) {
						mergedIndex, err := mapIx(ctx, leftIndex, rightIndex)
						if err != nil {
							var d D
							return d, err
						}
						ret, err := fmap(mergedIndex, rightFocus)
						return ret, JoinCtxErr(ctx, err)
					}, leftFocus)
					return o2ret, JoinCtxErr(ctx, err)
				}, source)
				return o1ret, JoinCtxErr(ctx, err)
			}

		}
	}

	return UnsafeOmni[K, S, T, C, D, RETC, RWC, DIRC, ERRC](
		getFnc,
		setFnc,
		iterFnc,
		lengthGetter,
		modifyFnc,
		ixgetter,
		ixMatchFnc,
		func(ctx context.Context, s D) (T, error) {
			left, err := o2ReverseGetter(ctx, s)
			err = JoinCtxErr(ctx, err)
			if err != nil {
				var ft T
				return ft, err
			}
			right, err := o1ReverseGetter(ctx, left)
			err = JoinCtxErr(ctx, err)
			return right, err
		},
		expressionHandlerCompose(o1, o2),
		composeExpr,
	)
}

func expressionHandlerCompose(e ...asExprHandler) func(ctx context.Context) (ExprHandler, error) {
	return func(ctx context.Context) (ExprHandler, error) {

		var found ExprHandler

		for _, h := range e {
			if h == nil {
				continue
			}

			hh, err := getExprHandler(ctx, h.AsExprHandler())
			if err != nil {
				return nil, err
			}

			if hh != nil {
				if found != nil {
					if found.TypeId() != hh.TypeId() {
						return nil, fmt.Errorf("mismatched expression handler types %v & %v", found.TypeId(), hh.TypeId())
					}
				} else {
					found = hh
				}
			}
		}

		return found, nil
	}
}

func expressionHandlerComposeVarArgs[T asExprHandler](e ...T) func(ctx context.Context) (ExprHandler, error) {
	return func(ctx context.Context) (ExprHandler, error) {

		var found ExprHandler

		for _, h := range e {
			hh, err := getExprHandler(ctx, h.AsExprHandler())
			if err != nil {
				return nil, err
			}

			if hh != nil {
				if found != nil {
					if found.TypeId() != hh.TypeId() {
						return nil, fmt.Errorf("mismatched expression handler types %v & %v", found.TypeId(), hh.TypeId())
					}
				} else {
					found = hh
				}
			}
		}

		return found, nil
	}
}

func getExprHandler(ctx context.Context, f func(ctx context.Context) (ExprHandler, error)) (ExprHandler, error) {
	if f == nil {
		return nil, nil
	}

	return f(ctx)
}

// IxMap returns a simple [IxMapper] where the left and right indexes cannot be recovered from the mapped value
//
// Note: For a more efficient composition when used as an [Traversal] optic then [IxMapIso] should be used.
//
// See: [ComposeI] for an example usage
func IxMap[LEFT, RIGHT, MAPPED any](ixmap func(LEFT, RIGHT) MAPPED) IxMapper[LEFT, RIGHT, MAPPED, ReturnOne, ReadOnly, UniDir, Pure] {
	return CombiGetter[Pure, MAPPED, lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]], lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]], MAPPED, MAPPED](
		func(ctx context.Context, source lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]]) (MAPPED, MAPPED, error) {
			l, lok := source.A.Get()
			r, rok := source.B.Get()
			if !lok || !rok {
				var m MAPPED
				//This error should never occur the options are only used in the return
				//We return an error to avoid a panic in the case the optic action can return an error.
				return m, m, errors.New("ixmap missing index")
			}
			m := ixmap(l, r)
			return m, m, nil
		},
		nil,
		ExprDef(func(t expr.OpticTypeExpr) expr.OpticExpression {
			return expr.IxMap{
				OpticTypeExpr: t,
				Type:          expr.IxMapperCustom,
			}
		}),
	)
}

// IxMapIso returns a [IxMapper] where the left and right may be recovered from the mapped value
// This enables more efficient operation when the composed optic is used as an [Traversal]
//
// See: [IxMap] for a simpler inefficient version.
func IxMapIso[LEFT, RIGHT, MAPPED any](
	ixmap func(left LEFT, right RIGHT) MAPPED,
	ixmatch func(a, b MAPPED) bool,
	unmap func(mapped MAPPED) (LEFT, bool, RIGHT, bool),
	exprDef ExpressionDef,
) IxMapper[LEFT, RIGHT, MAPPED, ReturnOne, ReadWrite, BiDir, Pure] {

	ixmatch = ensureIxMatch(ixmatch)

	getter := func(ctx context.Context, source lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]]) (MAPPED, error) {
		l, lok := source.A.Get()
		r, rok := source.B.Get()
		if !lok || !rok {
			var m MAPPED
			return m, errors.New("ixmap both missing index")
		}
		return ixmap(l, r), nil
	}

	reverse := func(ctx context.Context, focus MAPPED) (lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]], error) {

		l, lok, r, rok := unmap(focus)

		return lo.T2(mo.TupleToOption(l, lok), mo.TupleToOption(r, rok)), nil
	}

	return Omni[MAPPED, lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]], lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]], MAPPED, MAPPED, ReturnOne, ReadWrite, BiDir, Pure](
		func(ctx context.Context, source lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]]) (MAPPED, MAPPED, error) {
			ret, err := getter(ctx, source)
			return ret, ret, err
		},
		func(ctx context.Context, focus MAPPED, source lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]]) (lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]], error) {
			ret, err := reverse(ctx, focus)
			return ret, err
		},
		func(ctx context.Context, source lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]]) SeqIE[MAPPED, MAPPED] {
			return func(yield func(ValueIE[MAPPED, MAPPED]) bool) {
				v, err := getter(ctx, source)
				yield(ValIE(v, v, err))
			}
		},
		func(ctx context.Context, source lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]]) (int, error) {
			return 1, nil
		},
		func(ctx context.Context, fmap func(index MAPPED, focus MAPPED) (MAPPED, error), source lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]]) (lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]], error) {
			fa, err := getter(ctx, source)
			if err != nil {
				var t lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]]
				return t, err
			}

			fb, err := fmap(fa, fa)
			err = JoinCtxErr(ctx, err)
			if err != nil {
				var t lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]]
				return t, err
			}
			ft, err := reverse(ctx, fb)
			return ft, err
		},
		func(ctx context.Context, index MAPPED, source lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]]) SeqIE[MAPPED, MAPPED] {
			return func(yield func(ValueIE[MAPPED, MAPPED]) bool) {
				ret, err := getter(ctx, source)
				yield(ValIE(index, ret, err))
			}
		},
		ixmatch,
		func(ctx context.Context, focus MAPPED) (lo.Tuple2[mo.Option[LEFT], mo.Option[RIGHT]], error) {
			ret, err := reverse(ctx, focus)
			return ret, err
		},
		exprDef,
	)
}

func IxMapLeft[I, J any](ixmatch func(I, I) bool) IxMapper[I, J, I, ReturnOne, ReadWrite, BiDir, Pure] {
	return IxMapIso[I, J, I](
		func(left I, right J) I {
			return left
		},
		ixmatch,
		func(mapped I) (I, bool, J, bool) {
			var j J
			return mapped, true, j, false
		},
		ExprDef(func(t expr.OpticTypeExpr) expr.OpticExpression {
			return expr.IxMap{
				OpticTypeExpr: t,
				Type:          expr.IxMapperLeft,
			}
		}),
	)
}

func IxMapRight[I, J any](ixmatch func(J, J) bool) IxMapper[I, J, J, ReturnOne, ReadWrite, BiDir, Pure] {
	return IxMapIso[I, J, J](
		func(left I, right J) J {
			return right
		},
		ixmatch,
		func(mapped J) (I, bool, J, bool) {
			var i I
			return i, false, mapped, true
		},
		ExprDef(func(t expr.OpticTypeExpr) expr.OpticExpression {
			return expr.IxMap{
				OpticTypeExpr: t,
				Type:          expr.IxMapperRight,
			}
		}),
	)
}

func IxMapBoth[I, J any](
	ixmatchI func(I, I) bool,
	ixmatchJ func(J, J) bool,
) IxMapper[I, J, lo.Tuple2[I, J], ReturnOne, ReadWrite, BiDir, Pure] {
	return IxMapIso[I, J, lo.Tuple2[I, J]](
		func(left I, right J) lo.Tuple2[I, J] {
			return lo.T2(left, right)
		},
		func(t1, t2 lo.Tuple2[I, J]) bool {
			leftMatch := ixmatchI(t1.A, t2.A)
			if !leftMatch {
				return false
			}
			return ixmatchJ(t1.B, t2.B)
		},
		func(mapped lo.Tuple2[I, J]) (I, bool, J, bool) {
			return mapped.A, true, mapped.B, true
		},
		ExprDef(func(t expr.OpticTypeExpr) expr.OpticExpression {
			return expr.IxMap{
				OpticTypeExpr: t,
				Type:          expr.IxMapperBoth,
			}
		}),
	)
}

func rightIdentityComposeLeft[I, S, T, A, B, C, D, RETI, RETJ, RWI, RWJ, DIRI, DIRJ, ERRI, ERRJ any](composeExpr func() expr.OpticExpression, o1 Optic[I, S, T, A, B, RETI, RWI, DIRI, ERRI]) Optic[I, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]] {

	return UnsafeOmni[I, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]](
		func(ctx context.Context, source S) (I, C, error) {
			ix, ret, err := o1.AsGetter()(ctx, source)
			return ix, *(*C)(unsafe.Pointer(&ret)), err
		},
		func(ctx context.Context, focus D, source S) (T, error) {
			return o1.AsSetter()(ctx, *(*B)(unsafe.Pointer(&focus)), source)
		},
		func(ctx context.Context, source S) SeqIE[I, C] {
			return func(yield func(ValueIE[I, C]) bool) {
				ret := o1.AsIter()(ctx, source)
				ret(func(val ValueIE[I, A]) bool {
					return yield(ValIE[I, C](val.index, *(*C)(unsafe.Pointer(&val.value)), val.err))
				})
			}
		},
		o1.AsLengthGetter(),
		func(ctx context.Context, fmap func(index I, focus C) (D, error), source S) (T, error) {
			return o1.AsModify()(ctx, func(index I, focus A) (B, error) {
				ret, err := fmap(index, *(*C)(unsafe.Pointer(&focus)))
				return *(*B)(unsafe.Pointer(&ret)), err
			}, source)
		},
		func(ctx context.Context, index I, source S) SeqIE[I, C] {
			return func(yield func(ValueIE[I, C]) bool) {
				ret := o1.AsIxGetter()(ctx, index, source)
				ret(func(val ValueIE[I, A]) bool {
					return yield(ValIE[I, C](val.index, *(*C)(unsafe.Pointer(&val.value)), val.err))
				})
			}
		},
		o1.AsIxMatch(),
		func(ctx context.Context, focus D) (T, error) {
			return o1.AsReverseGetter()(ctx, *(*B)(unsafe.Pointer(&focus)))
		},
		o1.AsExprHandler(),
		composeExpr,
	)
}

func leftIdentityComposeLeft[I, J, S, T, A, B, C, D, RETI, RETJ, RWI, RWJ, DIRI, DIRJ, ERRI, ERRJ any](composeExpr func() expr.OpticExpression, o2 Optic[J, A, B, C, D, RETJ, RWJ, DIRJ, ERRJ]) Optic[I, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]] {
	return UnsafeOmni[I, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]](
		func(ctx context.Context, source S) (I, C, error) {
			var i I //I is Void
			_, ret, err := o2.AsGetter()(ctx, *(*A)(unsafe.Pointer(&source)))
			return i, *(*C)(unsafe.Pointer(&ret)), err
		},
		func(ctx context.Context, focus D, source S) (T, error) {
			b, err := o2.AsSetter()(ctx, focus, *(*A)(unsafe.Pointer(&source)))
			return *(*T)(unsafe.Pointer(&b)), err
		},
		func(ctx context.Context, source S) SeqIE[I, C] {
			return func(yield func(ValueIE[I, C]) bool) {
				o2.AsIter()(ctx, *(*A)(unsafe.Pointer(&source)))(func(focus ValueIE[J, C]) bool {
					var i I //I is Void
					return yield(ValIE(i, focus.value, focus.err))
				})
			}
		},
		func(ctx context.Context, source S) (int, error) {
			return o2.AsLengthGetter()(ctx, *(*A)(unsafe.Pointer(&source)))
		},
		func(ctx context.Context, fmap func(index I, focus C) (D, error), source S) (T, error) {
			b, err := o2.AsModify()(ctx, func(index J, focus C) (D, error) {
				var i I //I is Void
				return fmap(i, focus)
			}, *(*A)(unsafe.Pointer(&source)))

			return *(*T)(unsafe.Pointer(&b)), err
		},
		func(ctx context.Context, index I, source S) SeqIE[I, C] {
			return func(yield func(ValueIE[I, C]) bool) {
				o2.AsIter()(ctx, *(*A)(unsafe.Pointer(&source)))(func(focus ValueIE[J, C]) bool {
					var i I //I is Void so all right elements are matched
					return yield(ValIE(i, focus.value, focus.err))
				})
			}
		},
		func(indexA, indexB I) bool {
			//I is Void
			return true
		},
		func(ctx context.Context, focus D) (T, error) {
			b, err := o2.AsReverseGetter()(ctx, focus)
			return *(*T)(unsafe.Pointer(&b)), err
		},
		o2.AsExprHandler(),
		composeExpr,
	)
}

func leftIdentityCompose[J, S, T, A, B, C, D, RETI, RETJ any, RWI, RWJ any, DIRI, DIRJ any, ERRI, ERRJ any](composeExpr func() expr.OpticExpression, right Optic[J, A, B, C, D, RETJ, RWJ, DIRJ, ERRJ]) Optic[J, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]] {
	return UnsafeOmni[J, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]](
		func(ctx context.Context, source S) (J, C, error) {
			return right.AsGetter()(ctx, *(*A)(unsafe.Pointer(&source)))
		},
		func(ctx context.Context, focus D, source S) (T, error) {
			ret, err := right.AsSetter()(ctx, focus, *(*A)(unsafe.Pointer(&source)))
			return *(*T)(unsafe.Pointer(&ret)), err
		},
		func(ctx context.Context, source S) SeqIE[J, C] {
			return right.AsIter()(ctx, *(*A)(unsafe.Pointer(&source)))
		},
		func(ctx context.Context, source S) (int, error) {
			return right.AsLengthGetter()(ctx, *(*A)(unsafe.Pointer(&source)))
		},
		func(ctx context.Context, fmap func(index J, focus C) (D, error), source S) (T, error) {
			ret, err := right.AsModify()(ctx, fmap, *(*A)(unsafe.Pointer((&source))))
			return *(*T)(unsafe.Pointer(&ret)), err
		},
		func(ctx context.Context, index J, source S) SeqIE[J, C] {
			return right.AsIxGetter()(ctx, index, *(*A)(unsafe.Pointer(&source)))
		},
		right.AsIxMatch(),
		func(ctx context.Context, focus D) (T, error) {
			ret, err := right.AsReverseGetter()(ctx, focus)
			return *(*T)(unsafe.Pointer(&ret)), err
		},
		right.AsExprHandler(),
		composeExpr,
	)
}

func rightIdentityCompose[I, J, S, T, A, B, C, D, RETI, RETJ any, RWI, RWJ any, DIRI, DIRJ any, ERRI, ERRJ any](composeExpr func() expr.OpticExpression, left Optic[I, S, T, A, B, RETI, RWI, DIRI, ERRI]) Optic[J, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]] {
	return UnsafeOmni[J, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]](
		func(ctx context.Context, source S) (J, C, error) {
			//We know that J is Void
			var j J
			_, a, err := left.AsGetter()(ctx, source)
			return j, *(*C)(unsafe.Pointer(&a)), err
		},
		func(ctx context.Context, focus D, source S) (T, error) {
			return left.AsSetter()(ctx, *(*B)(unsafe.Pointer(&focus)), source)
		},
		func(ctx context.Context, source S) SeqIE[J, C] {
			return func(yield func(ValueIE[J, C]) bool) {
				left.AsIter()(ctx, source)(func(focus ValueIE[I, A]) bool {
					var j J
					return yield(ValIE(j, *(*C)(unsafe.Pointer(&focus.value)), focus.err))
				})
			}
		},
		left.AsLengthGetter(),
		func(ctx context.Context, fmap func(index J, focus C) (D, error), source S) (T, error) {

			return left.AsModify()(ctx, func(index I, focus A) (B, error) {
				var j J
				d, err := fmap(j, *(*C)(unsafe.Pointer(&focus)))
				return *(*B)(unsafe.Pointer(&d)), err
			}, source)
		},
		func(ctx context.Context, index J, source S) SeqIE[J, C] {
			return func(yield func(ValueIE[J, C]) bool) {
				left.AsIter()(ctx, source)(func(focus ValueIE[I, A]) bool {
					//J is Void and matches every value
					var j J
					return yield(ValIE(j, *(*C)(unsafe.Pointer(&focus.value)), focus.err))
				})
			}
		},
		func(indexA, indexB J) bool {
			//J is void
			return true
		},
		func(ctx context.Context, focus D) (T, error) {
			return left.AsReverseGetter()(ctx, *(*B)(unsafe.Pointer(&focus)))
		},
		left.AsExprHandler(),
		composeExpr,
	)
}

func leftAndRightIdentityCompose[J, S, T, A, B, C, D, RETI, RETJ any, RWI, RWJ any, DIRI, DIRJ any, ERRI, ERRJ any](composeExpr func() expr.OpticExpression) Optic[J, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]] {
	omni := UnsafeOmni[J, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]](
		func(ctx context.Context, source S) (J, C, error) {
			//We know that J is Void
			var j J
			return j, *(*C)(unsafe.Pointer(&source)), nil
		},
		func(ctx context.Context, focus D, source S) (T, error) {
			return *(*T)(unsafe.Pointer(&focus)), nil
		},
		func(ctx context.Context, source S) SeqIE[J, C] {
			return func(yield func(ValueIE[J, C]) bool) {
				var j J
				yield(ValIE(j, *(*C)(unsafe.Pointer(&source)), nil))
			}
		},
		func(ctx context.Context, source S) (int, error) {
			return 1, nil
		},
		func(ctx context.Context, fmap func(index J, focus C) (D, error), source S) (T, error) {
			var j J
			d, err := fmap(j, *(*C)(unsafe.Pointer(&source)))
			return *(*T)(unsafe.Pointer(&d)), err
		},
		func(ctx context.Context, index J, source S) SeqIE[J, C] {
			return func(yield func(ValueIE[J, C]) bool) {
				//J is Void and matches every value
				var j J
				yield(ValIE(j, *(*C)(unsafe.Pointer(&source)), nil))
			}
		},
		func(indexA, indexB J) bool {
			//J is void
			return true
		},
		func(ctx context.Context, focus D) (T, error) {
			return *(*T)(unsafe.Pointer(&focus)), nil
		},
		nil,
		composeExpr,
	).(omniOptic[J, S, T, C, D, CompositionTree[RETI, RETJ], CompositionTree[RWI, RWJ], CompositionTree[DIRI, DIRJ], CompositionTree[ERRI, ERRJ]])

	omni.opticType |= expr.OpticTypeIdentity

	return omni
}
