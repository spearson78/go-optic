package main

import (
	"fmt"
	"image"
	"image/color"
	"iter"
	"log"
	"math"

	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	. "github.com/spearson78/go-optic"
	. "github.com/spearson78/go-optic/examples/pong/model"
)

const (
	screenWidth  = 320
	screenHeight = 240

	replayLen = 250
)

var (
	whiteImage = ebiten.NewImage(3, 3).SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

type Game struct {
	//Immutable game state
	pong Pong

	//Action replay
	replay       func() (Pong, bool)
	stop         func()
	replayFrame  int
	replayBuffer ReplayBuffer[Pong]

	//Cpu target y position
	cpuTarget float64
}

func (g *Game) Update() error {

	if g.replay != nil {
		//Action replay is active. Retrieve the next state from the action replay.
		pong, ok := g.replay()
		if ok {
			g.replayFrame++
			g.pong = pong
		} else {
			//End of action repay
			g.stop()
			g.replay = nil
			g.stop = nil
			g.pong = reset(g.pong)
			g.replayBuffer = NewReplayBuffer[Pong]()
		}

		return nil
	} else {
		//No action replay so process user input and cpu target position.
		keyUp := false
		keyDown := false

		var keys [10]ebiten.Key
		for _, k := range inpututil.AppendPressedKeys(keys[:0]) {
			if ebiten.KeyName(k) == "w" {
				keyUp = true
			}

			if ebiten.KeyName(k) == "s" {
				keyDown = true
			}
		}

		//Store the key press in the game state
		g.pong = MustSet(O.Pong().KeyUpPressed(), keyUp, g.pong)
		g.pong = MustSet(O.Pong().KeyDownPressed(), keyDown, g.pong)

		//Update the game state and determine if the cpu has a new target
		var score bool
		var newTarget float64
		g.pong, score, newTarget = update(g.pong, g.cpuTarget)
		if !math.IsNaN(newTarget) {
			g.cpuTarget = newTarget
		}

		//Append the new game state to the action replay buffer
		g.replayBuffer = MustGet(
			AppendReplayBuffer(
				replayLen,
				ValCol(g.pong),
			),
			g.replayBuffer,
		)

		//A player has scored show the action replay
		if score {
			g.replay, g.stop = iter.Pull(
				MustGet(
					SeqOf(
						TraverseReplayBuffer[Pong]()),
					g.replayBuffer,
				),
			)
			g.replayFrame = 0
		}

		return nil
	}
}

func modelToScreen(width float64, height float64) Optic[Void, Point, Point, Point, Point, ReturnOne, ReadWrite, UniDir, Pure] {

	//The model coordinate system has the origin in the center of the playfield with a range from -1 to +1
	//The coordinates of an element are in its centre

	//This converts a game element coordinate for an element to screen coordinate.
	return Ret1(Rw(EPure(PointOf(
		Compose4(
			O.Point().X(),
			Add(1.0),             //Translate to the screen origin, model width is 2 (-1 -> 1)
			Sub(width/2.0),       //Translate to the center of the element
			Mul(screenWidth/2.0), //Scale it to screen width
		),
		Compose4(
			O.Point().Y(),
			Add(1.0),              //Translate to the screen origin, model height is 2 (-1 -> 1)
			Sub(height/2.0),       //Translate to the center of the element
			Mul(screenHeight/2.0), //Scale it to screen height
		),
	))))

}

// Convert ball coordinates to screen coordinates
var ballModelToScreen = Compose(
	O.Pong().BallPos(),
	modelToScreen(ballRadius, ballRadius),
)

// Convert paddle coordinates to screen coordinates
var paddlesModelToScreen = T2Of(
	Compose(
		O.Pong().Paddle1Point(),
		modelToScreen(paddleWidth, paddleHeight),
	),
	Compose(
		O.Pong().Paddle2Point(),
		modelToScreen(paddleWidth, paddleHeight),
	),
)

func (g *Game) drawBall(screen *ebiten.Image) {
	ballScreenPoint := MustGet(
		ballModelToScreen,
		g.pong,
	)
	vector.DrawFilledCircle(screen, float32(MustGet(O.Point().X(), ballScreenPoint)), float32(MustGet(O.Point().Y(), ballScreenPoint)), ballRadius*(screenWidth/2.0), color.White, true)
}

func (g *Game) drawPaddle(screen *ebiten.Image, paddleScreenPoint Point, color color.Color) {
	vector.DrawFilledRect(screen, float32(MustGet(O.Point().X(), paddleScreenPoint)), float32(MustGet(O.Point().Y(), paddleScreenPoint)), paddleWidth*(screenWidth/2), paddleHeight*(screenHeight/2), color, true)
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.drawBall(screen)

	paddlePoints := MustGet(
		paddlesModelToScreen, //Convert both paddle coordinates to screen.
		g.pong,
	)

	//Draw the 2 paddles
	g.drawPaddle(screen, paddlePoints.A, color.White)
	g.drawPaddle(screen, paddlePoints.B, color.White)

	//Draw the scores
	text.Draw(screen, fmt.Sprintf("%v | %v", MustGet(O.Pong().Score1(), g.pong), MustGet(O.Pong().Score2(), g.pong)), bitmapfont.Face, (screenWidth/2)-12, 12, color.White)

	//Flash Action Replay if a replay is playing
	if g.replay != nil {
		if (g.replayFrame>>4)%2 == 0 {
			text.Draw(screen, "Action Replay", bitmapfont.Face, (screenWidth/2)-24, screenHeight/2, color.White)
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Optics Pong Example")
	if err := ebiten.RunGame(&Game{
		pong:         reset(Pong{}),
		replayBuffer: NewReplayBuffer[Pong](),
	}); err != nil {
		log.Fatal(err)
	}
}
