package untilwin

import (
	"fmt"
	"time"

	"github.com/dystopium/2048/players"
	"github.com/dystopium/2048/runners"

	"github.com/dystopium/2048/game"
)

// Run plays the game until a game is won
func Run(gg runners.GameGen, play players.Player) {

	g := &game.Game{}
	var numGames uint64
	var numMoves uint64
	start := time.Now()

	for g.State() != game.StateWon {
		g = gg()
		play(g)

		numGames++
		numMoves += g.TotalMoves()

		if numGames%1000 == 0 {
			elapsed := time.Since(start)
			fmt.Printf("Played %v games in %v with an average %v moves to failure\n", numGames, elapsed, numMoves/numGames)
		}
	}

	fmt.Printf("\nWinning took %v games\n", numGames)
	fmt.Printf("\nScore: %v\tMoves: %v\n\n", g.Score(), g.TotalMoves())
	fmt.Println(g)
}
