package ojson

import (
	"errors"
	"regexp"

	. "github.com/spearson78/go-optic"
)

func Object() Optic[Void, any, any, map[string]any, map[string]any, ReturnMany, ReadWrite, BiDir, Pure] {
	return DownCast[any, map[string]any]()
}

func ObjectE() Optic[Void, any, any, map[string]any, map[string]any, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoCastE[any, map[string]any]()
}

func Array() Optic[Void, any, any, []any, []any, ReturnMany, ReadWrite, BiDir, Pure] {
	return DownCast[any, []any]()
}

func ArrayE() Optic[Void, any, any, []any, []any, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoCastE[any, []any]()
}

func String() Optic[Void, any, any, string, string, ReturnMany, ReadWrite, BiDir, Pure] {
	return DownCast[any, string]()
}

func StringE() Optic[Void, any, any, string, string, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoCastE[any, string]()
}

func Bool() Optic[Void, any, any, bool, bool, ReturnMany, ReadWrite, BiDir, Pure] {
	return DownCast[any, bool]()
}

func BoolE() Optic[Void, any, any, bool, bool, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoCastE[any, bool]()
}

func Float() Optic[Void, any, any, float64, float64, ReturnMany, ReadWrite, BiDir, Pure] {
	return DownCast[any, float64]()
}

func FloatE() Optic[Void, any, any, float64, float64, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoCastE[any, float64]()
}

func Int() Optic[Void, any, any, int64, int64, ReturnMany, ReadWrite, BiDir, Pure] {
	return DownCast[any, int64]()
}

func IntE() Optic[Void, any, any, int64, int64, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoCastE[any, int64]()
}

func Int32() Optic[Void, any, any, int32, int32, ReturnMany, ReadWrite, BiDir, Pure] {
	return DownCast[any, int32]()
}

func Int32E() Optic[Void, any, any, int32, int32, ReturnOne, ReadWrite, BiDir, Err] {
	return IsoCastE[any, int32]()
}

// Key focuses the given index of a JSON object (map[string]any)
// If the source is not a JSON object then nothing is focused
//
// See:
//   - [KeyE] for a version that returns an error if the source is not a JSON object
//   - [Nth] for a version that focuses the elements of a JSON array
func Key(index string) builder[string, Pure] {

	at := ComposeLeft(AtMap[string, any](index), Non[any](nil, EqT2[any]()))

	x := EPure(RetM(Rw(Ud(Compose(DownCast[any, map[string]any](), at)))))

	return builder[string, Pure]{
		Optic: x,
	}
}

// KeyE focuses the given index of a JSON object (map[string]any)
// If the source is not a JSON object then an error is returned.
//
// See:
//   - [KeyE] for a version that returns an error if the source is not a JSON object
//   - [Nth] for a version that focuses the elements of a JSON array
func KeyE(index string) builder[string, Err] {

	at := ComposeLeft(AtMap[string, any](index), Non[any](nil, EqT2[any]()))

	x := EErr(RetM(Rw(Ud(Compose(IsoCastE[any, map[string]any](), at)))))

	return builder[string, Err]{
		Optic: x,
	}
}

// KeyLike focuses the properties of the JSON object where the key matches the given regexp.
func KeyLike(rexpr *regexp.Regexp) builder[string, Pure] {
	keyLike := FilteredI(
		Compose(
			DownCast[any, map[string]any](),
			TraverseMap[string, any](),
		), OpOnIx[any](NotEmpty(MatchString(rexpr, -1))),
	)

	x := EPure(RetM(Rw(Ud(keyLike))))
	return builder[string, Pure]{
		Optic: x,
	}
}

func ValLike(rexpr *regexp.Regexp) builder[string, Pure] {

	keyLike := Filtered(

		Compose(
			DownCast[any, map[string]any](),
			TraverseMap[string, any](),
		), NotEmpty(Compose(
			DownCast[any, string](),
			MatchString(rexpr, -1),
		)),
	)

	x := EPure(RetM(Rw(Ud(keyLike))))
	return builder[string, Pure]{
		Optic: x,
	}
}

func Nth(index int) builder[int, Pure] {
	at := ComposeLeft(AtSlice[any](index, EqDeepT2[any]()), Non[any](nil, EqT2[any]()))
	x := EPure(RetM(Rw(Ud(Compose(DownCast[any, []any](), at)))))
	return builder[int, Pure]{
		Optic: x,
	}
}

func NthE(index int) builder[int, Err] {
	at := ComposeLeft(AtSlice[any](index, EqDeepT2[any]()), Non[any](nil, EqT2[any]()))
	x := EErr(RetM(Rw(Ud(Compose(IsoCastE[any, []any](), at)))))
	return builder[int, Err]{
		Optic: x,
	}
}

func traverse() Optic[any, any, any, any, any, ReturnMany, ReadWrite, UniDir, Pure] {
	return RetM(Rw(Ud(EPure(Coalesce(
		Compose(DownCast[any, []any](), ReIndexed(TraverseSlice[any](), UpCast[int, any](), EqT2[any]())),
		Compose(DownCast[any, map[string]any](), ReIndexed(TraverseMap[string, any](), UpCast[string, any](), EqT2[any]())),
	)))))
}

func Traverse() builder[any, Pure] {
	t := traverse()
	return builder[any, Pure]{
		Optic: t,
	}
}

var ErrTraverse = errors.New("ojson Traverse unknown type")

func traverseE() Optic[any, any, any, any, any, ReturnMany, ReadWrite, UniDir, Err] {
	return RetM(Rw(EErr(If(
		OrOp(
			NotEmpty(DownCast[any, []any]()),
			NotEmpty(DownCast[any, map[string]any]()),
		),
		EErr(traverse()),
		ReIndexed(Error[any, any](ErrTraverse), UpCast[Void, any](), EqT2[any]()),
	))))
}

// Note: returns an error if the source is not traversable.
func TraverseE() builder[any, Err] {
	t := traverseE()
	return builder[any, Err]{
		Optic: t,
	}
}

type builder[I, ERR any] struct {
	Optic[I, any, any, any, any, ReturnMany, ReadWrite, UniDir, ERR]
}

func (b builder[I, ERR]) Key(index string) builder[string, ERR] {
	k := Key(index)
	c := EErrL(RetM(Rw(Ud(Compose(b.Optic, k)))))
	return builder[string, ERR]{
		Optic: c,
	}
}

func (b builder[I, ERR]) KeyE(index string) builder[string, Err] {
	k := KeyE(index)
	c := EErr(RetM(Rw(Ud(Compose(b.Optic, k)))))
	return builder[string, Err]{
		Optic: c,
	}
}

func (b builder[I, ERR]) KeyLike(rexpr *regexp.Regexp) builder[string, ERR] {
	k := KeyLike(rexpr)
	c := EErrL(RetM(Rw(Ud(Compose(b.Optic, k)))))
	return builder[string, ERR]{
		Optic: c,
	}
}

func (b builder[I, ERR]) Eq(right any) Optic[Void, any, any, bool, bool, ReturnMany, ReadOnly, UniDir, ERR] {
	k := Eq[any](right)
	c := EErrL(RetM(Ro(Ud(Compose(b.Optic, k)))))
	return c
}

func (b builder[I, ERR]) Like(right *regexp.Regexp) Optic[MatchIndex, any, any, bool, bool, ReturnOne, ReadOnly, UniDir, ERR] {
	return EErrL(EErrL(NotEmpty(
		Compose3(
			b.Optic,
			DownCast[any, string](),
			MatchString(right, -1),
		),
	)))
}

func (b builder[I, ERR]) Gt(right float64) Optic[Void, any, any, bool, bool, ReturnMany, ReadOnly, UniDir, ERR] {
	k := Compose(DownCast[any, float64](), Gt(right))
	c := EErrL(RetM(Ro(Ud(Compose(b.Optic, k)))))
	return c

}

func (b builder[I, ERR]) Gte(right float64) Optic[Void, any, any, bool, bool, ReturnMany, ReadOnly, UniDir, ERR] {
	k := Compose(DownCast[any, float64](), Gte(right))
	c := EErrL(RetM(Ro(Ud(Compose(b.Optic, k)))))
	return c

}

func (b builder[I, ERR]) Lt(right float64) Optic[Void, any, any, bool, bool, ReturnMany, ReadOnly, UniDir, ERR] {
	k := Compose(DownCast[any, float64](), Lt(right))
	c := EErrL(RetM(Ro(Ud(Compose(b.Optic, k)))))
	return c

}

func (b builder[I, ERR]) Lte(right float64) Optic[Void, any, any, bool, bool, ReturnMany, ReadOnly, UniDir, ERR] {
	k := Compose(DownCast[any, float64](), Lte(right))
	c := EErrL(RetM(Ro(Ud(Compose(b.Optic, k)))))
	return c

}

func (b builder[I, ERR]) Nth(index int) builder[int, ERR] {
	k := Nth(index)
	c := EErrL(RetM(Rw(Ud(Compose(b.Optic, k)))))
	return builder[int, ERR]{
		Optic: c,
	}
}

func (b builder[I, ERR]) NthE(index int) builder[int, Err] {
	k := NthE(index)
	c := EErr(RetM(Rw(Ud(Compose(b.Optic, k)))))
	return builder[int, Err]{
		Optic: c,
	}
}

func (b builder[I, ERR]) Object() Optic[Void, any, any, map[string]any, map[string]any, ReturnMany, ReadWrite, UniDir, ERR] {
	return EErrL(RetM(Rw(Ud(Compose(b.Optic, Object())))))
}

func (b builder[I, ERR]) ObjectE() Optic[Void, any, any, map[string]any, map[string]any, ReturnMany, ReadWrite, UniDir, Err] {
	return EErr(RetM(Rw(Ud(Compose(b.Optic, ObjectE())))))
}

func (b builder[I, ERR]) String() Optic[Void, any, any, string, string, ReturnMany, ReadWrite, UniDir, ERR] {
	return EErrL(RetM(Rw(Ud(Compose(b.Optic, String())))))
}

func (b builder[I, ERR]) StringE() Optic[Void, any, any, string, string, ReturnMany, ReadWrite, UniDir, Err] {
	return EErr(RetM(Rw(Ud(Compose(b.Optic, StringE())))))
}

func (b builder[I, ERR]) Bool() Optic[Void, any, any, bool, bool, ReturnMany, ReadWrite, UniDir, ERR] {
	return EErrL(RetM(Rw(Ud(Compose(b.Optic, Bool())))))
}

func (b builder[I, ERR]) BoolE() Optic[Void, any, any, bool, bool, ReturnMany, ReadWrite, UniDir, Err] {
	return EErr(RetM(Rw(Ud(Compose(b.Optic, BoolE())))))
}

func (b builder[I, ERR]) Int() Optic[Void, any, any, int64, int64, ReturnMany, ReadWrite, UniDir, ERR] {
	return EErrL(RetM(Rw(Ud(Compose(b.Optic, Int())))))
}

func (b builder[I, ERR]) IntE() Optic[Void, any, any, int64, int64, ReturnMany, ReadWrite, UniDir, Err] {
	return EErr(RetM(Rw(Ud(Compose(b.Optic, IntE())))))
}

func (b builder[I, ERR]) Int32() Optic[Void, any, any, int32, int32, ReturnMany, ReadWrite, UniDir, ERR] {
	return EErrL(RetM(Rw(Ud(Compose(b.Optic, Int32())))))
}

func (b builder[I, ERR]) Int32E() Optic[Void, any, any, int32, int32, ReturnMany, ReadWrite, UniDir, Err] {
	return EErr(RetM(Rw(Ud(Compose(b.Optic, Int32E())))))
}

func (b builder[I, ERR]) Float() Optic[Void, any, any, float64, float64, ReturnMany, ReadWrite, UniDir, ERR] {
	return EErrL(RetM(Rw(Ud(Compose(b.Optic, Float())))))
}

func (b builder[I, ERR]) FloatE() Optic[Void, any, any, float64, float64, ReturnMany, ReadWrite, UniDir, Err] {
	return EErr(RetM(Rw(Ud(Compose(b.Optic, FloatE())))))
}

func (b builder[I, ERR]) Array() Optic[Void, any, any, []any, []any, ReturnMany, ReadWrite, UniDir, ERR] {
	return EErrL(RetM(Rw(Ud(Compose(b.Optic, Array())))))
}

func (b builder[I, ERR]) ArrayE() Optic[Void, any, any, []any, []any, ReturnMany, ReadWrite, UniDir, Err] {
	return EErr(RetM(Rw(Ud(Compose(b.Optic, ArrayE())))))
}

func (b builder[I, ERR]) Traverse() builder[any, ERR] {
	return builder[any, ERR]{
		Optic: EErrL(RetM(Rw(Ud(Compose(b.Optic, traverse()))))),
	}
}

func (b builder[I, ERR]) TraverseE() builder[any, Err] {
	return builder[any, Err]{
		Optic: EErr(RetM(Rw(Ud(Compose(b.Optic, traverseE()))))),
	}
}
