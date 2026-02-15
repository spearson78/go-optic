package optic

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
)

func ValidateOpticTest[I any, S comparable, A any, RET any, RW, DIR, ERR any](t *testing.T, o Optic[I, S, S, A, A, RET, RW, DIR, ERR], source S, newValue A) {
	ValidateOpticTestPred(t, o, source, newValue, EqT2[S]())
}

func ValidateOpticTestPred[I, S, A any, RET any, RW, DIR, ERR any, PERR any](t *testing.T, o Optic[I, S, S, A, A, RET, RW, DIR, ERR], source S, newValue A, eq Predicate[lo.Tuple2[S, S], PERR]) {
	t.Run(fmt.Sprintf("%v source %v newVal %v", o.AsExpr().Short(), source, newValue), func(t *testing.T) {
		err := ValidateOpticPred(o, source, newValue, eq)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func ValidateOpticTestP[I any, S comparable, T, A, B any, RET any, RW, DIR, ERR any](t *testing.T, o Optic[I, S, T, A, B, RET, RW, DIR, ERR], source S, newValue A, aToB func(A) B, tToS func(T) S) {
	ValidateOpticTestPredP(t, o, source, newValue, aToB, tToS, EqT2[S]())
}

func ValidateOpticTestPredP[I, S, T, A, B any, RET any, RW, DIR, ERR any, PERR any](t *testing.T, o Optic[I, S, T, A, B, RET, RW, DIR, ERR], source S, newValue A, aToB func(A) B, tToS func(T) S, eq Predicate[lo.Tuple2[S, S], PERR]) {
	t.Run(fmt.Sprintf("%v source %v newVal %v", o.AsExpr().Short(), source, newValue), func(t *testing.T) {
		err := ValidateOpticPredP(o, source, newValue, aToB, tToS, eq)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func ValidateOpticFuzz[I any, S comparable, A any, RET any, RW, DIR, ERR any, PERR TPure](f *testing.F, o Optic[I, S, S, A, A, RET, RW, DIR, ERR]) {
	f.Fuzz(func(t *testing.T, initialSource S, newValue A) {
		ValidateOpticTest(t, o, initialSource, newValue)
	})
}
