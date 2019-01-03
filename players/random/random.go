package random

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"

	"github.com/dystopium/2048/game"
)

func newRand() *rand.Rand {
	b := make([]byte, 8)
	crand.Read(b)
	seed := binary.LittleEndian.Uint64(b)
	return rand.New(rand.NewSource(int64(seed)))
}

// Play allows the human to play
func Play(g *game.Game) {
	rnd := newRand()

	var currnd uint64

	for g.State() == game.StatePlaying {
		if currnd == 0 {
			currnd = rnd.Uint64()
		}

		dir := currnd & 0x3
		currnd >>= 2

		switch dir {
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
}
