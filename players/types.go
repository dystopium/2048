package players

import "github.com/dystopium/2048/game"

// Player represents something that can attempt to win the game
type Player interface {
	Play(*game.Game)
}

// Const is a fucntion that returns a player
type Const func() Player
