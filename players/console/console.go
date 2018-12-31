package console

import (
	"bufio"
	"fmt"
	"os"

	"github.com/dystopium/2048/game"
)

// Player is a human playing at a console window
type Player struct{}

// Play allows the human to play
func (cp Player) Play(g *game.Game) {
	input := bufio.NewReader(os.Stdin)

	fmt.Println("Use IJKL for Up Left Down Right")

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
