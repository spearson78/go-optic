package optic

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"
	"github.com/spearson78/go-optic/expr"
)

func seqToSlice[I, A any](i SeqIE[I, A]) []ValueIE[I, A] {
	var ret []ValueIE[I, A]
	i(func(val ValueIE[I, A]) bool {
		index, focus, err := val.Get()
		ret = append(ret, ValIE(index, focus, err))
		return true
	})

	return ret
}

func eqValIESlice[I, A any](a, b []ValueIE[I, A]) bool {
	if len(a) != len(b) {
		return false
	}
	for i, aVal := range a {
		bVal := b[i]
		if !eqValIE(aVal, bVal) {
			return false
		}
	}

	return true
}

func eqValIE[I, A any](a, b ValueIE[I, A]) bool {

	if !MustGet(EqDeepT2[I](), lo.T2(a.index, b.index)) {
		return false
	}

	if !MustGet(EqDeepT2[A](), lo.T2(a.value, b.value)) {
		return false
	}

	if a.Error() != nil {
		if b.Error() != nil {
			return a.Error().Error() == b.Error().Error()
		} else {
			return false
		}
	} else {
		if b.Error() != nil {
			return false
		} else {
			//nil==nil
		}
	}

	return true

}

func testLensLaw[I, S, T, A, B, RET, RW, DIR, ERR any](ctx context.Context, o Optic[I, S, T, A, B, RET, RW, DIR, ERR], initialSource S, newValue A, aToB func(A) B, tToS func(T) S) error {

	newS, err := o.AsSetter()(ctx, aToB(newValue), initialSource)
	if err != nil {
		return fmt.Errorf("set failed %v", err)
	}
	_, newA, err := o.AsGetter()(ctx, tToS(newS))
	if err != nil {
		return fmt.Errorf("get failed %v", err)
	}

	if !MustGet(EqDeepT2[A](), lo.T2(newValue, newA)) {
		return fmt.Errorf("set-get failed expected: %v actual: %v", newValue, newA)
	}

	_, initialA, err := o.AsGetter()(ctx, initialSource)
	if err != nil {
		return fmt.Errorf("get failed %v", err)
	}

	getSetS, err := o.AsSetter()(ctx, aToB(initialA), initialSource)
	if err != nil {
		return fmt.Errorf("get-set failed %v", err)
	}

	if !MustGet(EqDeepT2[S](), lo.T2(initialSource, tToS(getSetS))) {
		return fmt.Errorf("get-set failed expected: %v actual: %v", initialSource, getSetS)
	}

	doubleSetS, err := o.AsSetter()(ctx, aToB(newValue), tToS(newS))
	if err != nil {
		return fmt.Errorf("set-set failed %v", err)
	}

	if !MustGet(EqDeepT2[T](), lo.T2(newS, doubleSetS)) {
		return fmt.Errorf("set-set failed expected: %v actual: %v", newS, doubleSetS)
	}

	return nil
}

func testIsoLaw[I, S, T, A, B, RET, RW, DIR, ERR any](ctx context.Context, o Optic[I, S, T, A, B, RET, RW, DIR, ERR], initialSource S, newValue A, aToB func(A) B, tToS func(T) S) error {

	_, a, err := o.AsGetter()(ctx, initialSource)
	if err != nil {
		return fmt.Errorf("AsGetter failed : %w", err)
	}

	b := aToB(a)

	t, err := o.AsReverseGetter()(ctx, b)
	if err != nil {
		return fmt.Errorf("AsReverseGetter failed : %w", err)
	}

	s := tToS(t)

	if !MustGet(EqDeepT2[S](), lo.T2(initialSource, s)) {
		return fmt.Errorf("iso reversibility failed failed expected: %v got: %v", initialSource, s)
	}

	newB := aToB(newValue)

	newT, err := o.AsReverseGetter()(ctx, newB)
	if err != nil {
		return fmt.Errorf("AsReverseGetter failed : %w", err)
	}

	newS := tToS(newT)

	_, newA, err := o.AsGetter()(ctx, newS)
	if err != nil {
		return fmt.Errorf("AsGetter failed : %w", err)
	}

	if !MustGet(EqDeepT2[A](), lo.T2(newValue, newA)) {
		return fmt.Errorf("iso reverse reversibility failed failed expected: %v got: %v", newValue, newA)
	}

	return nil
}

func testTraversalLaw[I, S, T, A, B, RET, RW, DIR, ERR any](ctx context.Context, o Optic[I, S, T, A, B, RET, RW, DIR, ERR], initialSource S, aToB func(A) B, tToS func(T) S, eq func(a, b S) (bool, error)) error {

	modifyPure, err := o.AsModify()(ctx, func(index I, focus A) (B, error) {
		return aToB(focus), nil
	}, initialSource)
	if err != nil {
		return fmt.Errorf("modify failed %v", err)
	}

	ok, err := eq(initialSource, tToS(modifyPure))
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("modify-pure failed expected: %v actual: %v", initialSource, modifyPure)
	}

	return nil
}

func testPrismLaw[I, S, T, A, B, RET, RW, DIR, ERR any](ctx context.Context, o Optic[I, S, T, A, B, RET, RW, DIR, ERR], initialSource S, aToB func(A) B, tToS func(T) S) error {

	var preview A
	var ok bool
	var err error
	o.AsIter()(ctx, initialSource)(func(val ValueIE[I, A]) bool {
		_, focus, focusErr := val.Get()
		preview = focus
		err = focusErr
		ok = true
		return false
	})

	if err != nil {
		return fmt.Errorf("iterate failed %v", err)
	}

	if ok {

		review, err := o.AsReverseGetter()(ctx, aToB(preview))
		if err != nil {
			return fmt.Errorf("reverse get failed %v", err)
		}

		if !MustGet(EqDeepT2[S](), lo.T2(initialSource, tToS(review))) {
			return fmt.Errorf("preview-review failed expected: %v actual: %v", initialSource, review)
		}

		var preview2 A
		ok = false
		o.AsIter()(ctx, tToS(review))(func(val ValueIE[I, A]) bool {
			_, focus, focusErr := val.Get()
			preview2 = focus
			err = focusErr
			ok = true
			return false
		})

		if !ok || err != nil {
			return fmt.Errorf("review-preview iterate failed %v", err)
		}

		if !MustGet(EqDeepT2[A](), lo.T2(preview, preview2)) {
			return fmt.Errorf("review-preview failed expected: %v actual: %v", initialSource, review)
		}

	}

	return nil
}

func ValidateOptic[I any, S comparable, A any, RET any, RW, DIR, ERR any](o Optic[I, S, S, A, A, RET, RW, DIR, ERR], source S, newValue A) error {
	return ValidateOpticP(o, source, newValue, func(a A) A { return a }, func(s S) S { return s })
}

func ValidateOpticPred[I, S, A any, RET any, RW, DIR, ERR any, PERR any](o Optic[I, S, S, A, A, RET, RW, DIR, ERR], source S, newValue A, eq Predicate[lo.Tuple2[S, S], PERR]) error {
	return ValidateOpticPredP(o, source, newValue, func(a A) A { return a }, func(s S) S { return s }, eq)
}

func ValidateOpticP[I any, S comparable, T, A, B any, RET any, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], source S, newValue A, aToB func(A) B, tToS func(T) S) error {
	return ValidateOpticPredP(o, source, newValue, aToB, tToS, EqT2[S]())
}

func ValidateOpticPredP[I any, S, T, A, B any, RET any, RW, DIR, ERR any, PERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], source S, newValue A, aToB func(A) B, tToS func(T) S, eq Predicate[lo.Tuple2[S, S], PERR]) error {

	ctx := context.Background()

	eqOp := func(a, b S) (bool, error) {
		return PredGet(ctx, eq, lo.T2(a, b))
	}

	if isType(o.OpticType(), expr.OpticTypeReturnManyFlag) {
		if isType(o.OpticType(), expr.OpticTypeReadWriteFlag) {
			if isType(o.OpticType(), expr.OpticTypeBiDirFlag) {
				//Prism
				err := testPrismLaw(ctx, o, source, aToB, tToS)
				if err != nil {
					return err
				}
			} else {
				//Traversal
				err := testTraversalLaw(ctx, o, source, aToB, tToS, eqOp)
				if err != nil {
					return err
				}
			}
		} else {
			if isType(o.OpticType(), expr.OpticTypeBiDirFlag) {
				//Illegal
				return errors.New("optic is ReadOnly and BiDir. BiDir optics must be ReadWrite")
			} else {
				//Iteration
			}
		}
	} else {
		if isType(o.OpticType(), expr.OpticTypeReadWriteFlag) {
			if isType(o.OpticType(), expr.OpticTypeBiDirFlag) {
				//Iso
				err := testIsoLaw(ctx, o, source, newValue, aToB, tToS)
				if err != nil {
					return err
				}
			} else {
				//Lens
				err := testLensLaw(ctx, o, source, newValue, aToB, tToS)
				if err != nil {
					return err
				}
			}
		} else {
			if isType(o.OpticType(), expr.OpticTypeBiDirFlag) {
				// Illegal
				return errors.New("optic is ReadOnly and BiDir. BiDir optics must be ReadWrite")
			} else {
				//Getter / Operation
			}
		}
	}

	err := testConsistencyP(ctx, o, source, newValue, aToB, tToS, eqOp)
	if err != nil {
		return err
	}

	return nil

}

func testConsistencyP[I, S, T, A, B any, RET any, RW, DIR, ERR any](ctx context.Context, o Optic[I, S, T, A, B, RET, RW, DIR, ERR], source S, newValue A, aToB func(A) B, tToS func(T) S, eq func(a, b S) (bool, error)) error {

	var iterRes []ValueIE[I, A]

	o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
		index, focus, err := val.Get()
		iterRes = append(iterRes, ValIE(index, focus, err))
		return true
	})

	if res, err := o.AsLengthGetter()(ctx, source); err != nil || res != len(iterRes) {
		return fmt.Errorf("consistency check AsLengthGetter Expected:%v Got:%v Err:%v", len(iterRes), res, err)
	}

	if len(iterRes) == 0 {
		return nil
	}

	for i := 0; i < 10; i++ {
		var iterRes2 []ValueIE[I, A]
		o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
			index, focus, err := val.Get()
			iterRes2 = append(iterRes2, ValIE(index, focus, err))
			return true
		})
		if !eqValIESlice(iterRes, iterRes2) {
			return fmt.Errorf("consistency check iteration order changed expected: %v got: %v on iteration %v", iterRes, iterRes2, i)
		}
	}

	//Getter should always return the first iterated value
	if _, res, err := o.AsGetter()(ctx, source); err != iterRes[0].Error() || !MustGet(EqDeepT2[A](), lo.T2(iterRes[0].value, res)) {
		return fmt.Errorf("consistency check AsGetter Expected: %v Got: %v Err: %v", iterRes[0].value, res, err)
	}

	if isType(o.OpticType(), expr.OpticTypeBiDirFlag) {
		if res, err := o.AsReverseGetter()(ctx, aToB(iterRes[0].value)); err != nil || !MustGet(EqDeepT2[S](), lo.T2(source, tToS(res))) {
			return fmt.Errorf("consistency check AsReverseGet Expected: %v Got: %v Err: %v ReverseGet of: %v", source, res, err, iterRes[0].value)
		}
	}

	var groupedIndexes []lo.Tuple2[I, []ValueIE[I, A]]

	for i := 0; i < len(iterRes); i++ {

		grouped := false
		for j, _ := range groupedIndexes {
			v := &groupedIndexes[j]
			if MustGet(EqDeepT2[I](), lo.T2(v.A, iterRes[i].index)) {
				v.B = append(v.B, iterRes[i])
				grouped = true
				break
			}
		}

		if !grouped {
			groupedIndexes = append(groupedIndexes, lo.T2(iterRes[i].index, []ValueIE[I, A]{iterRes[i]}))
		}
	}

	for i := 0; i < len(groupedIndexes); i++ {
		if res := seqToSlice(o.AsIxGetter()(ctx, groupedIndexes[i].A, source)); !eqValIESlice(groupedIndexes[i].B, res) {
			return fmt.Errorf("consistency check AsIxGetter Expected: %v Got: %v ChosenIndex: %v", groupedIndexes[i].B, res, groupedIndexes[i].A)
		}
	}

	if isType(o.OpticType(), expr.OpticTypeReadWriteFlag) {

		//The identity rule
		modifyPure, err := o.AsModify()(ctx, func(index I, focus A) (B, error) {
			return aToB(focus), nil
		}, source)
		if err != nil {
			return fmt.Errorf("modify failed %v", err)
		}

		ok, err := eq(source, tToS(modifyPure))
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("identity rule failed expected: %v actual: %v", source, modifyPure)
		}

		modifiedCount := 0

		modified, err := o.AsModify()(ctx, func(index I, focus A) (B, error) {
			modifiedCount++
			return aToB(newValue), nil
		}, source)
		if err != nil {
			return fmt.Errorf("consistency check AsModify Err: %v", err)
		}

		modifyIterCount := 0
		var iterErr error

		o.AsIter()(ctx, tToS(modified))(func(val ValueIE[I, A]) bool {
			_, focus, err := val.Get()
			if err != nil {
				iterErr = err
				return false
			}
			modifyIterCount++
			if !MustGet(EqDeepT2[A](), lo.T2(focus, newValue)) {
				iterErr = fmt.Errorf("consistency check AsModify failed to update all elements iterated value: %v expected %v", focus, newValue)
				return false
			}
			return true
		})

		if iterErr != nil {
			return iterErr
		}

		if modifiedCount != len(iterRes) || modifiedCount != modifyIterCount {
			return fmt.Errorf("consistency check AsModify failed to update all elements modified : %v originally iterated: %v modified iterated: %v", modifiedCount, len(iterRes), modifyIterCount)
		}

		set, err := o.AsSetter()(ctx, aToB(newValue), source)
		if err != nil {
			return fmt.Errorf("consistency check AsSetter Err: %v", err)
		}

		//THe modify is equivalent to a set
		ok, err = eq(tToS(modified), tToS(set))
		if err != nil {
			return err
		}

		if !ok {
			return fmt.Errorf("consistency check AsModify AsSetter mismatch modified: %v set: %v", modified, set)
		}

	}

	return nil
}
