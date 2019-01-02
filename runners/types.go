package runners

import (
	"github.com/dystopium/2048/game"
	"github.com/dystopium/2048/players"
)

// GameGen creates a new game each time it is called
// Ideally it's closed over the config needed
type GameGen func() *game.Game

// Runner is a thing that will run games however it sees fit
type Runner func(GameGen, players.Player)
