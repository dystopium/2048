package main

import (
	"flag"
	"fmt"

	"github.com/dystopium/2048/game"
	"github.com/dystopium/2048/players"
	"github.com/dystopium/2048/players/console"
	"github.com/dystopium/2048/players/random"
)

func main() {
	var playerType string
	var width int
	var height int
	var limitPower uint

	flag.StringVar(&playerType, "player", "console", "Player type. One of: console")
	flag.IntVar(&width, "width", 4, "Width of the playing board.")
	flag.IntVar(&height, "height", 4, "Height of the playing board.")
	flag.UintVar(&limitPower, "lim", 11, "Power of 2 to set as the winning number. Default gives 2048.")
	flag.Parse()

	var player players.Player

	switch playerType {
	case "console":
		player = console.Player{}

	case "random":
		player = random.Player{}
	}

	g := &game.Game{}
	var numGames int

	for g.State() != game.StateWon {
		g = game.NewGame(width, height, limitPower)
		player.Play(g)
		//fmt.Println(g.State())

		numGames++

		if numGames%1000 == 0 {
			fmt.Println(numGames)
		}
	}

	fmt.Printf("\nWinning took %v games\n\n", numGames)
	fmt.Printf("\nScore: %v\tMoves: %v\n\n", g.Score(), g.TotalMoves())
	fmt.Println(g)
}
