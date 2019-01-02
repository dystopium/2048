package single

import (
	"fmt"

	"github.com/dystopium/2048/game"

	"github.com/dystopium/2048/players"
	"github.com/dystopium/2048/runners"
)

// Run will run a single game until it wins or loses
func Run(gg runners.GameGen, pc players.Const) {
	player := pc()
	g := gg()

	player.Play(g)

	switch g.State() {
	case game.StateWon:
		fmt.Println("\nYOU WON!")

	case game.StateLost:
		fmt.Println("\nYOU LOST!")
	}
	fmt.Printf("\nScore: %v\tMoves: %v\n\n", g.Score(), g.TotalMoves())
	fmt.Println(g)
}
