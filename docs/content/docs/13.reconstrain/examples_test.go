package main

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
)

func TestReconstrainCompositionTree(t *testing.T) {
	//BEGIN compositiontree
	var optic Optic[
		Void,
		float64,
		float64,
		float64,
		float64,
		CompositionTree[ReturnOne, ReturnOne],
		CompositionTree[ReadWrite, ReadWrite],
		CompositionTree[BiDir, BiDir],
		CompositionTree[Pure, Pure],
	]
	optic = Compose(
		Mul(1.8),
		Add(32.0),
	)
	//END compositiontree

	fmt.Println(optic.AsExpr())
}

//BEGIN celsiusToFahrenheit1

func celsiusToFahrenheit1() Optic[
	Void,
	float64,
	float64,
	float64,
	float64,
	CompositionTree[ReturnOne, ReturnOne],
	CompositionTree[ReadWrite, ReadWrite],
	CompositionTree[BiDir, BiDir],
	CompositionTree[Pure, Pure],
] {
	return Compose(
		Mul(1.8),
		Add(32.0),
	)
}

//END celsiusToFahrenheit1

//BEGIN celsiusToFahrenheit2

func celsiusToFahrenheit2() Optic[
	Void,
	float64,
	float64,
	float64,
	float64,
	ReturnOne,
	CompositionTree[ReadWrite, ReadWrite],
	CompositionTree[BiDir, BiDir],
	CompositionTree[Pure, Pure],
] {
	return Ret1(
		Compose(
			Mul(1.8),
			Add(32.0),
		),
	)
}

//END celsiusToFahrenheit2

//BEGIN celsiusToFahrenheit3

func celsiusToFahrenheit3() Optic[
	Void,
	float64,
	float64,
	float64,
	float64,
	ReturnOne,
	ReadWrite,
	BiDir,
	Pure,
] {
	return Ret1(Rw(Bd(EPure(
		Compose(
			Mul(1.8),
			Add(32.0),
		),
	))))
}

//END celsiusToFahrenheit3

//BEGIN bytesOf1

func bytesOf1[I, S any, RET, RW, DIR, ERR any](
	o Optic[I, S, S, string, string, RET, RW, DIR, ERR],
) Optic[
	I,
	S,
	S,
	[]byte,
	[]byte,
	CompositionTree[RET, ReturnOne],
	CompositionTree[RW, ReadWrite],
	CompositionTree[DIR, BiDir],
	CompositionTree[ERR, Pure],
] {
	return ComposeLeft(
		o,
		IsoCast[string, []byte](),
	)
}

//END bytesOf1

//BEGIN bytesOf2

func bytesOf2[I, S any, RET, RW, DIR, ERR any](
	o Optic[I, S, S, string, string, RET, RW, DIR, ERR],
) Optic[
	I,
	S,
	S,
	[]byte,
	[]byte,
	RET,
	RW,
	DIR,
	ERR,
] {
	return RetL(RwL(DirL(EErrL(
		ComposeLeft(
			o,
			IsoCast[string, []byte](),
		),
	))))
}

//END bytesOf2

//BEGIN makeTuple1

func makeTuple1[I, S any, RET, RW, DIR, ERR any](
	o Optic[I, S, S, string, string, RET, RW, DIR, ERR],
) Optic[
	lo.Tuple2[I, I],
	lo.Tuple2[S, S],
	lo.Tuple2[S, S],
	lo.Tuple2[string, string],
	lo.Tuple2[string, string],
	CompositionTree[
		CompositionTree[ReturnOne, RET],
		CompositionTree[ReturnOne, RET],
	],
	CompositionTree[
		CompositionTree[ReadWrite, RW],
		CompositionTree[ReadWrite, RW],
	],
	UniDir,
	CompositionTree[
		CompositionTree[Pure, ERR],
		CompositionTree[Pure, ERR],
	],
] {
	return T2Of(
		Compose(
			T2A[S, S](),
			o,
		),
		Compose(
			T2B[S, S](),
			o,
		),
	)
}

//END makeTuple1

//BEGIN makeTuple2

func makeTuple2[I, S any, RET, RW, DIR, ERR any](
	o Optic[I, S, S, string, string, RET, RW, DIR, ERR],
) Optic[
	lo.Tuple2[I, I],
	lo.Tuple2[S, S],
	lo.Tuple2[S, S],
	lo.Tuple2[string, string],
	lo.Tuple2[string, string],
	RET,
	RW,
	UniDir,
	ERR,
] {
	return RetMerge(RwMerge(EErrMerge(
		T2Of(
			RetR(RwR(EErrR(
				Compose(
					T2A[S, S](),
					o,
				),
			))),
			RetR(RwR(EErrR(
				Compose(
					T2B[S, S](),
					o,
				),
			))),
		),
	)))
}

//END makeTuple2

//BEGIN makeTuple9

func makeTuple9[I, S any, RET, RW, DIR, ERR any](
	o Optic[I, S, S, string, string, RET, RW, DIR, ERR],
) Optic[
	lo.Tuple9[I, I, I, I, I, I, I, I, I],
	lo.Tuple9[S, S, S, S, S, S, S, S, S],
	lo.Tuple9[S, S, S, S, S, S, S, S, S],
	lo.Tuple9[string, string, string, string, string, string, string, string, string],
	lo.Tuple9[string, string, string, string, string, string, string, string, string],
	CompositionTree[
		CompositionTree[
			CompositionTree[
				CompositionTree[RET, RET],
				RET,
			],
			CompositionTree[RET, RET]],
		CompositionTree[
			CompositionTree[RET, RET],
			CompositionTree[RET, RET],
		],
	],
	CompositionTree[
		CompositionTree[
			CompositionTree[
				CompositionTree[RW, RW],
				RW],
			CompositionTree[RW, RW],
		],
		CompositionTree[
			CompositionTree[RW, RW],
			CompositionTree[RW, RW],
		],
	],
	UniDir,
	CompositionTree[
		CompositionTree[
			CompositionTree[
				CompositionTree[ERR, ERR],
				ERR,
			],
			CompositionTree[ERR, ERR],
		],
		CompositionTree[
			CompositionTree[ERR, ERR],
			CompositionTree[ERR, ERR],
		],
	],
] {
	return T9Of(
		RetR(RwR(EErrR(
			Compose(
				T9A[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9B[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9C[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9D[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9E[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9F[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9G[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9H[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9I[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
	)
}

//END makeTuple9

func makeTuple9_reconstrain[I, S any, RET, RW, DIR, ERR any](
	o Optic[I, S, S, string, string, RET, RW, DIR, ERR],
) Optic[
	lo.Tuple9[I, I, I, I, I, I, I, I, I],
	lo.Tuple9[S, S, S, S, S, S, S, S, S],
	lo.Tuple9[S, S, S, S, S, S, S, S, S],
	lo.Tuple9[string, string, string, string, string, string, string, string, string],
	lo.Tuple9[string, string, string, string, string, string, string, string, string],
	RET,
	RW,
	UniDir,
	ERR,
] {

	makeTuple := T9Of(
		RetR(RwR(EErrR(
			Compose(
				T9A[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9B[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9C[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9D[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9E[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9F[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9G[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9H[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9I[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
	)

	//BEGIN makeTuple9_reconstrain
	retMerge := RetMerge(RetMerge(RetTransR(RetMergeL(RetTransL(RetSwapL(RetMergeR(RetTransL(RetMergeR(RetMergeR(makeTuple))))))))))
	rwMerge := RwMerge(RwMerge(RwTransR(RwMergeL(RwTransL(RwSwapL(RwMergeR(RwTransL(RwMergeR(RwMergeR(retMerge))))))))))
	errMerge := EErrMerge(EErrMerge(EErrTransR(EErrMergeL(EErrTransL(EErrSwapL(EErrMergeR(EErrTransL(EErrMergeR(EErrMergeR(rwMerge))))))))))
	//END makeTuple9_reconstrain

	return errMerge
}

func makeTuple9_unsafereconstrain[I, S any, RET, RW, DIR, ERR any](
	o Optic[I, S, S, string, string, RET, RW, DIR, ERR],
) Optic[
	lo.Tuple9[I, I, I, I, I, I, I, I, I],
	lo.Tuple9[S, S, S, S, S, S, S, S, S],
	lo.Tuple9[S, S, S, S, S, S, S, S, S],
	lo.Tuple9[string, string, string, string, string, string, string, string, string],
	lo.Tuple9[string, string, string, string, string, string, string, string, string],
	RET,
	RW,
	UniDir,
	ERR,
] {

	makeTuple := T9Of(
		RetR(RwR(EErrR(
			Compose(
				T9A[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9B[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9C[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9D[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9E[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9F[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9G[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9H[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
		RetR(RwR(EErrR(
			Compose(
				T9I[S, S, S, S, S, S, S, S, S](),
				o,
			),
		))),
	)

	//BEGIN makeTuple9_unsafereconstrain
	return UnsafeReconstrain[RET, RW, UniDir, ERR](makeTuple)
	//END makeTuple9_unsafereconstrain

}

//BEGIN makeTuple9_safereconstrain

func reconstrainMakeTuple9[I, S, RET, RW, ERR any](
	o Optic[
		lo.Tuple9[I, I, I, I, I, I, I, I, I],
		lo.Tuple9[S, S, S, S, S, S, S, S, S],
		lo.Tuple9[S, S, S, S, S, S, S, S, S],
		lo.Tuple9[string, string, string, string, string, string, string, string, string],
		lo.Tuple9[string, string, string, string, string, string, string, string, string],
		CompositionTree[
			CompositionTree[
				CompositionTree[
					CompositionTree[RET, RET],
					RET,
				],
				CompositionTree[RET, RET]],
			CompositionTree[
				CompositionTree[RET, RET],
				CompositionTree[RET, RET],
			],
		],
		CompositionTree[
			CompositionTree[
				CompositionTree[
					CompositionTree[RW, RW],
					RW],
				CompositionTree[RW, RW],
			],
			CompositionTree[
				CompositionTree[RW, RW],
				CompositionTree[RW, RW],
			],
		],
		UniDir,
		CompositionTree[
			CompositionTree[
				CompositionTree[
					CompositionTree[ERR, ERR],
					ERR,
				],
				CompositionTree[ERR, ERR],
			],
			CompositionTree[
				CompositionTree[ERR, ERR],
				CompositionTree[ERR, ERR],
			],
		],
	]) Optic[
	lo.Tuple9[I, I, I, I, I, I, I, I, I],
	lo.Tuple9[S, S, S, S, S, S, S, S, S],
	lo.Tuple9[S, S, S, S, S, S, S, S, S],
	lo.Tuple9[string, string, string, string, string, string, string, string, string],
	lo.Tuple9[string, string, string, string, string, string, string, string, string],
	RET,
	RW,
	UniDir,
	ERR,
] {
	return UnsafeReconstrain[RET, RW, UniDir, ERR](o)
}

//END makeTuple9_safereconstrain
