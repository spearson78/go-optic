package optic

import (
	"cmp"
)

// MakeLensSlice is used internally by the makelens tool.
type MakeLensSlice[I, S, T, A, RET, RW, DIR, ERR any] struct {
	Optic[I, S, T, Collection[int, A, Pure], Collection[int, A, Pure], RET, RW, UniDir, ERR]
	optic Optic[I, S, T, []A, []A, RET, RW, DIR, ERR]
}

func (s MakeLensSlice[I, S, T, A, RET, RW, DIR, ERR]) Traverse() Optic[int, S, T, A, A, ReturnMany, RW, UniDir, ERR] {
	//No intermediate slice is created here as we were created with an optic that focuses a slice
	return RetM(RwL(Ud(EErrL(Compose(s.optic, TraverseSlice[A]())))))
}

func (s MakeLensSlice[I, S, T, A, RET, RW, DIR, ERR]) Nth(index int) Optic[int, S, T, A, A, ReturnMany, RW, UniDir, ERR] {
	return RwL(EErrL(Index(Compose(s.optic, TraverseSlice[A]()), index)))
}

// NewMakeLensSlice is used internally by the makelens tool.
func NewMakeLensSlice[I, S, T, A, RET any, RW any, DIR any, ERR any](o Optic[I, S, T, []A, []A, RET, RW, DIR, ERR]) MakeLensSlice[I, S, T, A, RET, RW, DIR, ERR] {
	return MakeLensSlice[I, S, T, A, RET, RW, DIR, ERR]{
		Optic: RetL(RwL(Ud(EErrL(ComposeLeft(o, SliceToCol[A]()))))),
		optic: o,
	}
}

// MakeLensMap is used internally by the makelens tool.
type MakeLensMap[I comparable, S, T, A, RET, RW, DIR, ERR any] struct {
	Optic[Void, S, T, Collection[I, A, Pure], Collection[I, A, Pure], RET, RW, UniDir, ERR]
	optic Optic[Void, S, T, map[I]A, map[I]A, RET, RW, DIR, ERR]
}

func (s MakeLensMap[I, S, T, A, RET, RW, DIR, ERR]) Traverse() Optic[I, S, T, A, A, ReturnMany, RW, UniDir, ERR] {
	//No intermediate map is created here as we were created with an optic that focuses a map
	return RetM(RwL(Ud(EErrL(Compose(s.optic, TraverseMap[I, A]())))))
}

func (s MakeLensMap[I, S, T, A, RET, RW, DIR, ERR]) Key(index I) Optic[I, S, T, A, A, ReturnMany, RW, UniDir, ERR] {
	return RwL(EErrL(Index(Compose(s.optic, TraverseMap[I, A]()), index)))
}

// MakeLensMap is used internally by the makelens tool.
func NewMakeLensMap[I comparable, S, T, A, RET any, RW any, DIR any, ERR any](o Optic[Void, S, T, map[I]A, map[I]A, RET, RW, DIR, ERR]) MakeLensMap[I, S, T, A, RET, RW, DIR, ERR] {
	return MakeLensMap[I, S, T, A, RET, RW, DIR, ERR]{
		Optic: Ud(RwL(RetL(EErrL(ComposeLeft(o, MapToCol[I, A]()))))),
		optic: o,
	}
}

// MakeLensCmpOps is used internally by the makelens tool.
type MakeLensCmpOps[I, S, T any, A comparable, RET, RW, DIR, ERR any] struct {
	Optic[I, S, T, A, A, RET, RW, DIR, ERR]
}

// NewMakeLensCmpOps is used internally by the makelens tool.
func NewMakeLensCmpOps[I, S, T any, A comparable, RET, RW, DIR, ERR any](o Optic[I, S, T, A, A, RET, RW, DIR, ERR]) MakeLensCmpOps[I, S, T, A, RET, RW, DIR, ERR] {
	return MakeLensCmpOps[I, S, T, A, RET, RW, DIR, ERR]{
		Optic: o,
	}
}

func (o MakeLensCmpOps[I, S, T, A, RET, RW, DIR, ERR]) Eq(val A) Optic[I, S, T, bool, bool, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, Eq(val))))))
}

func (o MakeLensCmpOps[I, S, T, A, RET, RW, DIR, ERR]) Ne(val A) Optic[I, S, T, bool, bool, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, Ne(val))))))
}

func (o MakeLensCmpOps[I, S, T, A, RET, RW, DIR, ERR]) In(val ...A) Optic[I, S, T, bool, bool, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, In(val...))))))
}

func (o MakeLensCmpOps[I, S, T, A, RET, RW, DIR, ERR]) Contains(val A) Optic[I, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return Contains(o.Optic, val)
}

// MakeLensOrdOps is used internally by the makelens tool.
type MakeLensOrdOps[I, S, T any, A cmp.Ordered, RET, RW, DIR, ERR any] struct {
	Optic[I, S, T, A, A, RET, RW, DIR, ERR]
}

// NewMakeLensOrdOps is used internally by the makelens tool.
func NewMakeLensOrdOps[I, S, T any, A cmp.Ordered, RET, RW, DIR, ERR any](o Optic[I, S, T, A, A, RET, RW, DIR, ERR]) MakeLensOrdOps[I, S, T, A, RET, RW, DIR, ERR] {
	return MakeLensOrdOps[I, S, T, A, RET, RW, DIR, ERR]{
		Optic: o,
	}
}

func (o MakeLensOrdOps[I, S, T, A, RET, RW, DIR, ERR]) Eq(val A) Optic[I, S, T, bool, bool, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, Eq(val))))))
}

func (o MakeLensOrdOps[I, S, T, A, RET, RW, DIR, ERR]) Ne(val A) Optic[I, S, T, bool, bool, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, Ne(val))))))
}

func (o MakeLensOrdOps[I, S, T, A, RET, RW, DIR, ERR]) In(val ...A) Optic[I, S, T, bool, bool, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, In(val...))))))
}

func (o MakeLensOrdOps[I, S, T, A, RET, RW, DIR, ERR]) Contains(val A) Optic[I, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return Contains(o.Optic, val)
}

func (o MakeLensOrdOps[I, S, T, A, RET, RW, DIR, ERR]) Gt(val A) Optic[I, S, T, bool, bool, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, Gt(val))))))
}

func (o MakeLensOrdOps[I, S, T, A, RET, RW, DIR, ERR]) Gte(val A) Optic[I, S, T, bool, bool, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, Gte(val))))))
}

func (o MakeLensOrdOps[I, S, T, A, RET, RW, DIR, ERR]) Lt(val A) Optic[I, S, T, bool, bool, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, Lt(val))))))
}

func (o MakeLensOrdOps[I, S, T, A, RET, RW, DIR, ERR]) Lte(val A) Optic[I, S, T, bool, bool, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, Lte(val))))))
}

// MakeLensRealOps is used internally by the makelens tool.
type MakeLensRealOps[I, S, T any, A Real, RET, RW, DIR, ERR any] struct {
	Optic[I, S, T, A, A, RET, RW, DIR, ERR]
}

// NewMakeLensRealOps is used internally by the makelens tool.
func NewMakeLensRealOps[I, S, T any, A Real, RET, RW, DIR, ERR any](o Optic[I, S, T, A, A, RET, RW, DIR, ERR]) MakeLensRealOps[I, S, T, A, RET, RW, DIR, ERR] {
	return MakeLensRealOps[I, S, T, A, RET, RW, DIR, ERR]{
		Optic: o,
	}
}

func (o MakeLensRealOps[I, S, T, A, RET, RW, DIR, ERR]) Eq(val A) Optic[I, S, T, bool, bool, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, Eq(val))))))
}

func (o MakeLensRealOps[I, S, T, A, RET, RW, DIR, ERR]) Ne(val A) Optic[I, S, T, bool, bool, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, Ne(val))))))
}

func (o MakeLensRealOps[I, S, T, A, RET, RW, DIR, ERR]) In(val ...A) Optic[I, S, T, bool, bool, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, In(val...))))))
}

func (o MakeLensRealOps[I, S, T, A, RET, RW, DIR, ERR]) Contains(val A) Optic[I, S, S, bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return Contains(o.Optic, val)
}

func (o MakeLensRealOps[I, S, T, A, RET, RW, DIR, ERR]) Gt(val A) Optic[I, S, T, bool, bool, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, Gt(val))))))
}

func (o MakeLensRealOps[I, S, T, A, RET, RW, DIR, ERR]) Gte(val A) Optic[I, S, T, bool, bool, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, Gte(val))))))
}

func (o MakeLensRealOps[I, S, T, A, RET, RW, DIR, ERR]) Lt(val A) Optic[I, S, T, bool, bool, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, Lt(val))))))
}

func (o MakeLensRealOps[I, S, T, A, RET, RW, DIR, ERR]) Lte(val A) Optic[I, S, T, bool, bool, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, Lte(val))))))
}

func (o MakeLensRealOps[I, S, T, A, RET, RW, DIR, ERR]) Mul(val A) Optic[I, S, T, A, A, RET, RW, DIR, ERR] {
	return RetL(RwL(DirL(EErrL(ComposeLeft(o.Optic, Mul(val))))))
}

func (o MakeLensRealOps[I, S, T, A, RET, RW, DIR, ERR]) Div(val A) Optic[I, S, T, A, A, RET, RW, DIR, ERR] {
	return RetL(RwL(DirL(EErrL(ComposeLeft(o.Optic, Div(val))))))
}

func (o MakeLensRealOps[I, S, T, A, RET, RW, DIR, ERR]) Add(val A) Optic[I, S, T, A, A, RET, RW, DIR, ERR] {
	return RetL(RwL(DirL(EErrL(ComposeLeft(o.Optic, Add(val))))))
}

func (o MakeLensRealOps[I, S, T, A, RET, RW, DIR, ERR]) Sub(val A) Optic[I, S, T, A, A, RET, RW, DIR, ERR] {
	return RetL(RwL(DirL(EErrL(ComposeLeft(o.Optic, Sub(val))))))
}

func (o MakeLensRealOps[I, S, T, A, RET, RW, DIR, ERR]) Negate() Optic[I, S, T, A, A, RET, RW, DIR, ERR] {
	return RetL(RwL(DirL(EErrL(ComposeLeft(o.Optic, Negate[A]())))))
}

func (o MakeLensRealOps[I, S, T, A, RET, RW, DIR, ERR]) Min(val A) Optic[I, S, T, A, A, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, Min(val))))))
}

func (o MakeLensRealOps[I, S, T, A, RET, RW, DIR, ERR]) Max(val A) Optic[I, S, T, A, A, RET, ReadOnly, UniDir, ERR] {
	return RetL(Ro(Ud(EErrL(ComposeLeft(o.Optic, Max(val))))))
}
