package main

import (
	"math"
	"math/rand"

	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/examples/pong/model"
)

//This example is based on https://github.com/ekmett/lens/blob/master/examples/Pong.hs
//
// Original Copyright   :  (C) 2012 Edward Kmett, Niklas Haas
// Original License     :  BSD-style (see the file https://github.com/ekmett/lens/blob/master/examples/LICENSE)
//-----------------------------------------------------------------------------

const (
	ballRadius      = 0.02
	speedIncrease   = 1.05
	losingAccuracy  = 0.98
	winningAccuracy = 0.1
	initialSpeed    = 0.006
	paddleWidth     = 0.02
	paddleHeight    = 0.3
	paddleSpeed     = 0.02
)

func accuracy(pong Pong) float64 {

	score1 := MustGet(
		O.Pong().Score1(),
		pong,
	)

	score2 := MustGet(
		O.Pong().Score1(),
		pong,
	)

	return MustGet(
		//Clamp the result to a range
		Clamp(winningAccuracy, losingAccuracy),
		//As the player loses make the computer more inaccurate
		(float64(score2-score1)*-0.27)+0.98,
	)

}

// Game update logic
func update(pong Pong, cpuTarget float64) (Pong, bool, float64) {
	//Update the paddle positions provided the cpu with it's target position
	pong = updatePaddles(pong, cpuTarget)
	//Update the pall position
	pong = updateBall(pong)

	//Check for collisions with a paddle or an out of bounds causing a player to score.
	pong, score, newTarget := paddleBounce(pong)

	return pong, score, newTarget
}

const edge = 1.0 - (ballRadius + paddleWidth)

func updateBall(pong Pong) Pong {

	//Make sure it doesn't leave the playing area
	clampBallRadius := Clamp(ballRadius-1, 1-ballRadius)

	//Update the ball position
	pong = MustModify(
		O.Pong().BallPos(),
		PointOf(
			//Update x pos with U component of ball speed
			Compose3(
				O.Point().X(),
				Add(
					MustGet(
						O.Pong().BallSpeed().U(),
						pong,
					),
				),
				clampBallRadius,
			),
			//Update y pos with V component of ball speed
			Compose3(
				O.Point().Y(),
				Add(
					MustGet(
						O.Pong().BallSpeed().V(),
						pong,
					),
				),
				clampBallRadius,
			),
		),
		pong,
	)

	//Check for collisions with the top or bottom
	if MustGet(
		Compose(
			O.Pong().BallPos().Y(),
			Abs[float64](),
		),
		pong,
	) >= edge {
		pong = MustModify(
			O.Pong().BallSpeed().V(),
			Negate[float64](),
			pong,
		)
	}

	return pong
}

func updatePaddles(pong Pong, cpuTarget float64) Pong {

	paddleMovement := paddleSpeed

	//Update player paddles position
	if MustGet(O.Pong().KeyUpPressed(), pong) {
		pong = MustModify(
			O.Pong().Paddle1(),
			Compose(
				Sub(paddleMovement),
				Max(-1+(paddleHeight/2)),
			),
			pong,
		)
	} else if MustGet(O.Pong().KeyDownPressed(), pong) {
		pong = MustModify(
			O.Pong().Paddle1(),
			Compose(
				Add(paddleMovement),
				Min(1-(paddleHeight/2)),
			),
			pong,
		)
	}

	//Move the CPU's paddle towards the optimal position
	pong = MustModify(
		O.Pong().Paddle2(),
		Compose(
			Op(func(paddle float64) float64 {

				dist := cpuTarget - paddle

				if math.Abs(dist) > paddleHeight/3.0 {
					if dist > 0 {
						paddle += paddleMovement
					} else if dist < 0 {
						paddle -= paddleMovement
					}
				}

				return paddle
			}),
			//Ensure the paddle does not leave the playfield.
			Clamp(-1+(paddleHeight/2), 1-(paddleHeight/2)),
		),
		pong,
	)

	return pong
}

func paddleBounce(pong Pong) (Pong, bool, float64) {

	handlePaddleBounce := func(paddle float64, opponentScore Optic[Void, Pong, Pong, int, int, ReturnOne, ReadWrite, UniDir, Pure]) (Pong, bool) {
		ballPosY := MustGet(O.Pong().BallPos().Y(), pong)

		if ballPosY >= paddle-paddleHeight/2 && ballPosY <= paddle+paddleHeight/2 {
			//Ball has hit the paddle
			return MustModify(
				O.Pong().BallSpeed(),
				VectorOf(
					Compose3(
						O.Vector().U(),
						Negate[float64](),  //Invert the balls x direction
						Mul(speedIncrease), //And increase the speed of the game
					),
					Compose3(
						O.Vector().V(),
						Op(func(v float64) float64 {
							return v + 0.02*(ballPosY-paddle) //Set the y direction based on the location of the paddle hit.
						}),
						Mul(speedIncrease), //And increase the speed of the game
					),
				),
				pong,
			), false
		} else {
			//Ball has left the playfield increase the opponents score.
			pong = MustModify(
				opponentScore,
				Add(1),
				pong,
			)
			return pong, true
		}
	}

	target := math.NaN()
	score := false
	if MustGet(
		//Ball has exited the right side of the playfield
		O.Pong().BallPos().X().Gte(edge),
		pong,
	) {
		pong, score = handlePaddleBounce(
			MustGet(O.Pong().Paddle2(), pong), //Right side paddle position
			O.Pong().Score1(),                 //Opponents score to increase
		)
	} else if MustGet(
		//Ball has exited the left side of the playfield
		O.Pong().BallPos().X().Lte(-edge),
		pong,
	) {
		pong, score = handlePaddleBounce(
			MustGet(O.Pong().Paddle1(), pong), //Left side paddle position
			O.Pong().Score2(),                 //Opponents score to increase
		)

		if !score {
			//The ball bounced of the user's paddle calculate the new cpu target position.

			futurePong := pong
			//Update a new branch of the game forward until the ball exits the right edge.
			for MustGet(O.Pong().BallPos().X().Lt(edge), futurePong) {
				futurePong = updateBall(futurePong)
			}

			//Determine a random error to the position based on the score difference.
			acc := rand.Float64() * (1 - accuracy(pong))

			//Apply the random error to the final ball position
			optimum := MustGet(O.Pong().BallPos().Y(), futurePong)
			if rand.Float64() >= 0.5 {
				target = optimum + (acc * paddleHeight / 2)
			} else {
				target = optimum - (acc * paddleHeight / 2)
			}
		}
	}

	return pong, score, target
}

func reset(pong Pong) Pong {
	//Reset the game state
	pong = MustSet(O.Pong().BallPos(), Point{}, pong)
	pong = MustSet(O.Pong().BallSpeed(), NewVector(-initialSpeed, -initialSpeed), pong)

	return pong
}
