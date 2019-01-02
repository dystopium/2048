package parallel

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/dystopium/2048/players"
	"github.com/dystopium/2048/runners"

	"github.com/dystopium/2048/game"
)

// Run runs games in as many cores as possible
func Run(gg runners.GameGen, pc players.Const) {

	// times 2 for no good reason
	numCores := runtime.NumCPU()

	stop := make(chan *sync.WaitGroup)
	results := make(chan *game.Game, numCores*2)
	wg := &sync.WaitGroup{}

	for i := 0; i < numCores; i++ {
		wg.Add(1)

		go func(gg runners.GameGen, pc players.Const) {
			for {
				player := pc()
				g := gg()

				player.Play(g)

				// either send the result
				// or accept the signal to stop looping
				// This guarantees that if the results channel is full,
				// it can still complete, it'll just throw out the result
				select {
				case results <- g:
				case wg := <-stop:
					wg.Done()
					return
				}
			}
		}(gg, pc)
	}

	g := &game.Game{}
	var numGames uint64
	var numMoves uint64
	start := time.Now()

	for g.State() != game.StateWon {
		g = <-results

		numGames++
		numMoves += g.TotalMoves()

		if numGames%10000 == 0 {
			elapsed := time.Since(start)
			fmt.Printf("Played %v games in %v with an average %v moves to failure\n", numGames, elapsed, numMoves/numGames)
		}
	}

	// After the win, shut down the workers
	for i := 0; i < numCores; i++ {
		stop <- wg
	}
	wg.Wait()

	fmt.Printf("\nWinning took %v games\n", numGames)
	fmt.Printf("\nScore: %v\tMoves: %v\n\n", g.Score(), g.TotalMoves())
	fmt.Println(g)
}
