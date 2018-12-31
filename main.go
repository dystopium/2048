package main

import (
	"flag"

	"github.com/dystopium/2048/game"
	"github.com/dystopium/2048/players"
	"github.com/dystopium/2048/players/console"
)

func main() {
	var playerType string
	var width int
	var height int

	flag.StringVar(&playerType, "player", "console", "Player type. One of: console")
	flag.IntVar(&width, "width", 4, "Width of the playing board.")
	flag.IntVar(&height, "height", 4, "Height of the playing board.")
	flag.Parse()

	var player players.Player

	switch playerType {
	case "console":
		player = console.Player{}
	}

	g := game.NewGame(width, height)

	player.Play(g)
}
