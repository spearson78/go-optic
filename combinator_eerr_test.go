package optic_test

import (
	"fmt"

	. "github.com/spearson78/go-optic"
)

func FilteredColIExample[I comparable, J, S, T, A, RET, RW, DIR, ERR any](o Optic[J, S, T, Collection[I, A, ERR], Collection[I, A, ERR], RET, RW, DIR, ERR], index I) Optic[Void, S, T, Collection[I, A, ERR], Collection[I, A, ERR], RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrMerge(Compose(
		o,
		FilteredColI(
			CombiEErr[ERR](EqI[A](index)),
			IxMatchComparable[I](),
		),
	)))))
}

func ExampleCombiEErr() {

	data := []string{"alpha", "beta", "gamma", "delta"}

	optic := FilteredColIExample(
		SliceToCol[string](),
		1,
	)

	res := MustGet(optic, data)
	fmt.Println(res)

	//Output:
	//Col[1:beta]

}
