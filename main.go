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
	var width int
	var height int
	var limitPower uint
	var numAdds int

	flag.StringVar(&cpuprofilename, "cpuprofile", "", "File name for a CPU profile")
	flag.StringVar(&playerType, "player", "console", "Player type. One of: console")
	flag.IntVar(&width, "width", 4, "Width of the playing board.")
	flag.IntVar(&height, "height", 4, "Height of the playing board.")
	flag.UintVar(&limitPower, "lim", 11, "Power of 2 to set as the winning number. Default gives 2048.")
	flag.IntVar(&numAdds, "adds", 1, "The number of random values to add after each move")
	flag.Parse()

	if cpuprofilename != "" {
		cpuprofile, err := os.Create(cpuprofilename)
		if err != nil {
			panic(err)
		}

		pprof.StartCPUProfile(cpuprofile)
		defer pprof.StopCPUProfile()
	}

	var player players.Player

	switch playerType {
	case "console":
		player = console.Player{}

	case "random":
		player = random.Player{}
	}

	g := &game.Game{}
	var numGames int
	start := time.Now()

	for g.State() != game.StateWon {
		g = game.NewGame(width, height, limitPower, numAdds)
		player.Play(g)
		//fmt.Println(g.State())

		numGames++

		if numGames%1000 == 0 {
			elapsed := time.Since(start)
			fmt.Printf("Played %v games in %v\n", numGames, elapsed)
		}
	}

	fmt.Printf("\nWinning took %v games\n\n", numGames)
	fmt.Printf("\nScore: %v\tMoves: %v\n\n", g.Score(), g.TotalMoves())
	fmt.Println(g)
}
