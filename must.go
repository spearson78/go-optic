package optic

import "fmt"

// Must panics if err != nil
func Must[T any](t T, err error) T {
	if err != nil {
		panic(fmt.Errorf("must : %w", err))
	}
	return t
}

// Must2 panics if err != nil
func Must2[T any, T2 any](t T, t2 T2, err error) (T, T2) {
	if err != nil {
		panic(fmt.Errorf("must : %w", err))
	}
	return t, t2
}

// MustI panics if err != nil
func MustI[I, T any](i I, val T, err error) (I, T) {
	if err != nil {
		panic(fmt.Errorf("must : %w", err))
	}
	return i, val
}

// Must2I panics if err != nil
func Must2I[I, T, T2 any](i I, val T, val2 T2, err error) (I, T, T2) {
	if err != nil {
		panic(fmt.Errorf("must : %w", err))
	}
	return i, val, val2
}
