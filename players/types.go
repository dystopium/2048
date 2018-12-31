package players

import "github.com/dystopium/2048/game"

// Player represents something that can attempt to win the game
type Player interface {
	Play(*game.Game)
}
