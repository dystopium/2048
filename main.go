package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/dystopium/2048/game"
)

func main() {
	var gameType string
	var width int
	var height int

	flag.StringVar(&gameType, "type", "manual", "Game type. 'manual' to play in the console")
	flag.IntVar(&width, "width", 4, "Width of the playing board.")
	flag.IntVar(&height, "height", 4, "Height of the playing board.")
	flag.Parse()

	input := bufio.NewReader(os.Stdin)

	fmt.Println("Use IJKL for Up Left Down Right")

	g := game.NewGame(width, height)
	print := true

	for {
		if print {
			fmt.Printf("\nScore: %v\tMoves: %v\n\n", g.Score(), g.TotalMoves())
			fmt.Println(g)
		}
		print = true

		if g.State() != game.StatePlaying {
			break
		}

		cmd, _, err := input.ReadRune()
		if err != nil {
			panic(err)
		}

		switch cmd {
		case 'i':
			fallthrough
		case 'I':
			g.Move(game.DirUp)

		case 'k':
			fallthrough
		case 'K':
			g.Move(game.DirDown)

		case 'j':
			fallthrough
		case 'J':
			g.Move(game.DirLeft)

		case 'l':
			fallthrough
		case 'L':
			g.Move(game.DirRight)

		case '\n':
			print = false
		}
	}

	switch g.State() {
	case game.StateLost:
		fmt.Println("YOU LOST")

	case game.StateWon:
		fmt.Println("YOU WON")
	}
}
