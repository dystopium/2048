package random

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"

	"github.com/dystopium/2048/game"
	"github.com/dystopium/2048/players"
)

func newRand() *rand.Rand {
	b := make([]byte, 8)
	crand.Read(b)
	seed := binary.LittleEndian.Uint64(b)
	return rand.New(rand.NewSource(int64(seed)))
}

// NewConst returns a new constructor function for console players
func NewConst() players.Const {
	return func() players.Player { return &Player{} }
}

// Player is a human playing at a console window
type Player struct{}

// Play allows the human to play
func (cp Player) Play(g *game.Game) {
	rnd := newRand()

	//fmt.Println(g)

	for g.State() == game.StatePlaying {
		switch rnd.Intn(4) {
		case 0:
			g.Move(game.DirUp)

		case 1:
			g.Move(game.DirDown)

		case 2:
			g.Move(game.DirLeft)

		case 3:
			g.Move(game.DirRight)
		}
	}

	//fmt.Printf("\nScore: %v\tMoves: %v\n\n", g.Score(), g.TotalMoves())
	//fmt.Println(g)

	switch g.State() {
	case game.StateLost:
		//fmt.Println("LOST")

	case game.StateWon:
		//fmt.Println("WON")
	}
}
