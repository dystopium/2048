package multiwin

import (
	"fmt"
	"time"

	"github.com/dystopium/2048/players"
	"github.com/dystopium/2048/runners"

	"github.com/dystopium/2048/game"
)

// New creates a new runner that plays until numWins games have been won
func New(numWins uint64) runners.Runner {
	return func(gg runners.GameGen, play players.Player) {
		var totalGamesToWin uint64
		var totalMovesToWin uint64
		var totalWinningScore uint64

		for i := uint64(0); i < numWins; i++ {
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

			//fmt.Printf("\nWinning took %v games\n", numGames)
			//fmt.Printf("\nScore: %v\tMoves: %v\n\n", g.Score(), g.TotalMoves())
			//fmt.Println(g)

			totalGamesToWin += numGames
			totalMovesToWin += g.TotalMoves()
			totalWinningScore += g.Score()
		}

		fmt.Printf("\nAverage %v games to win\n", float64(totalGamesToWin)/float64(numWins))
		fmt.Printf("\nAverage %v moves in winning games\n", float64(totalMovesToWin)/float64(numWins))
		fmt.Printf("\nAverage %v score in winning games\n", float64(totalWinningScore)/float64(numWins))
	}
}
