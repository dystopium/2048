package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/dystopium/2048/game"
	"github.com/dystopium/2048/players"
	"github.com/dystopium/2048/players/console"
	"github.com/dystopium/2048/players/random"
)

func main() {
	var cpuprofilename string
	var playerType string
	var width uint64
	var height uint64
	var limitPower uint64
	var numAdds uint64

	flag.StringVar(&cpuprofilename, "cpuprofile", "", "File name for a CPU profile")
	flag.StringVar(&playerType, "player", "console", "Player type. One of: console")
	flag.Uint64Var(&width, "width", 4, "Width of the playing board.")
	flag.Uint64Var(&height, "height", 4, "Height of the playing board.")
	flag.Uint64Var(&limitPower, "lim", 11, "Power of 2 to set as the winning number. Default gives 2048.")
	flag.Uint64Var(&numAdds, "adds", 1, "The number of random values to add after each move")
	flag.Parse()

	if cpuprofilename != "" {
		cpuprofile, err := os.Create(cpuprofilename)
		if err != nil {
			panic(err)
		}

		pprof.StartCPUProfile(cpuprofile)
		defer pprof.StopCPUProfile()
	}

	var pc players.Const

	switch playerType {
	case "console":
		pc = console.NewConst()

	case "random":
		pc = random.NewConst()
	}

	g := &game.Game{}
	var numGames uint64
	var numMoves uint64
	start := time.Now()

	for g.State() != game.StateWon {
		player := pc()
		g = game.NewGame(width, height, limitPower, numAdds)
		player.Play(g)
		//fmt.Println(g.State())

		numGames++
		numMoves += g.TotalMoves()

		if numGames%1000 == 0 {
			elapsed := time.Since(start)
			fmt.Printf("Played %v games in %v with an average %v moves to failure\n", numGames, elapsed, numMoves/numGames)
		}
	}

	fmt.Printf("\nWinning took %v games\n\n", numGames)
	fmt.Printf("\nScore: %v\tMoves: %v\n\n", g.Score(), g.TotalMoves())
	fmt.Println(g)
}
