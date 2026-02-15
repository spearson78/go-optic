package model

import (
	"github.com/samber/lo"
	. "github.com/spearson78/go-optic"
)

//go:generate ../../../makelens model model.go model_generated.go

// Immutable data type for storing the game state
type Pong struct {
	ballPos   Point
	ballSpeed Vector
	paddle1   float64 //y location of paddle 1
	paddle2   float64 //y location of paddle 1
	score1    int
	score2    int

	keyUpPressed   bool
	keyDownPressed bool
}

type Point struct {
	x float64
	y float64
}

// Convert a [Point] to a [lo.Tuple2]
func PointToT2() Optic[Void, Point, Point, lo.Tuple2[float64, float64], lo.Tuple2[float64, float64], ReturnOne, ReadWrite, BiDir, Pure] {
	//Converting the Point to a T2 enables us to use T2Of to construct a Point.
	return Iso[Point, lo.Tuple2[float64, float64]](
		func(source Point) lo.Tuple2[float64, float64] {
			return lo.T2(source.x, source.y)
		},
		func(focus lo.Tuple2[float64, float64]) Point {
			return Point{
				x: focus.A,
				y: focus.B,
			}
		},
		ExprCustom("PointToT2"),
	)
}

// Construct [Point]s from an x and y optic.
func PointOf[I0 any, I1 any, S, T any, RET0 any, RW0 any, DIR0 any, ERR0 any, RET1 any, RW1 any, DIR1 any, ERR1 any](x Optic[I0, S, S, float64, float64, RET0, RW0, DIR0, ERR0], y Optic[I1, S, T, float64, float64, RET1, RW1, DIR1, ERR1]) Optic[Void, S, T, Point, Point, CompositionTree[RET0, RET1], CompositionTree[RW0, RW1], UniDir, CompositionTree[ERR0, ERR1]] {
	return RetL(RwL(Ud(EErrL(Compose(
		//Create a T2 of the x and y coordinate optics.
		T2Of(
			x,
			y,
		),
		//Convert the T2 to a Point
		AsReverseGet(PointToT2()),
	)))))
}

type Vector struct {
	u float64
	v float64
}

// Convert a [Vector] to a [lo.Tuple2]
func VectorToT2() Optic[Void, Vector, Vector, lo.Tuple2[float64, float64], lo.Tuple2[float64, float64], ReturnOne, ReadWrite, BiDir, Pure] {
	//Converting the Vector to a T2 enables us to use T2Of to construct a Vector.
	return Iso[Vector, lo.Tuple2[float64, float64]](
		func(source Vector) lo.Tuple2[float64, float64] {
			return lo.T2(source.u, source.v)
		},
		func(focus lo.Tuple2[float64, float64]) Vector {
			return Vector{
				u: focus.A,
				v: focus.B,
			}
		},
		ExprCustom("VectorToT2"),
	)
}

// Construct [Vector]s from a u and v optic.
func VectorOf[I0 any, I1 any, S, T any, RET0 any, RW0 any, DIR0 any, ERR0 any, RET1 any, RW1 any, DIR1 any, ERR1 any](u Optic[I0, S, S, float64, float64, RET0, RW0, DIR0, ERR0], v Optic[I1, S, T, float64, float64, RET1, RW1, DIR1, ERR1]) Optic[Void, S, T, Vector, Vector, CompositionTree[RET0, RET1], CompositionTree[RW0, RW1], UniDir, CompositionTree[ERR0, ERR1]] {
	return RetL(RwL(Ud(EErrL(Compose(
		//Create a T2 of the u and v component optics.
		T2Of(
			u,
			v,
		),
		//Convert the T2 to a Vector
		AsReverseGet(VectorToT2()),
	)))))
}

func NewVector(u float64, v float64) Vector {
	return Vector{
		u: u,
		v: v,
	}
}
