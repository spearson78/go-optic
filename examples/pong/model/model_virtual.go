package model

import (
	"github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic"
)

func (s *lPong[I, S, T, RET, RW, DIR, ERR]) Paddle1Point() optic.Optic[optic.Void, S, T, Point, Point, RET, RW, UniDir, ERR] {
	//in the model only the y location of the paddles are stored. We need a Point for rendering the paddle.
	return RetR(RwR(Ud(EErrR(PointOf(
		IgnoreWrite( //Make the constant value ReadWrite by ignoring the write
			ConstP[S, T, float64, float64](-1+0.02), //The constant x location of paddle 1
		),
		s.Paddle1(),
	)))))
}

func (s *lPong[I, S, T, RET, RW, DIR, ERR]) Paddle2Point() optic.Optic[optic.Void, S, T, Point, Point, RET, RW, UniDir, ERR] {
	//in the model only the y location of the paddles are stored. We need a Point for rendering the paddle.
	return RetR(RwR(Ud(EErrR(PointOf(
		IgnoreWrite( //Make the constant value ReadWrite by ignoring the write
			ConstP[S, T, float64, float64](1-0.02), //The constant x location of paddle 2
		),
		s.Paddle2(),
	)))))
}
