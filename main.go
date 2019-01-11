package main

import (
	"flag"
	"os"
	"runtime/pprof"

	"github.com/dystopium/2048/game"
	"github.com/dystopium/2048/players"
	"github.com/dystopium/2048/players/console"
	"github.com/dystopium/2048/players/greedy"
	"github.com/dystopium/2048/players/random"
	"github.com/dystopium/2048/runners"
	"github.com/dystopium/2048/runners/multiwin"
	"github.com/dystopium/2048/runners/parallel"
	"github.com/dystopium/2048/runners/single"
	"github.com/dystopium/2048/runners/untilwin"
)

func main() {
	var cpuprofilename string
	var playerType string
	var runnerType string
	var width uint64
	var height uint64
	var limitPower uint64
	var numAdds uint64
	var numWins uint64

	flag.StringVar(&cpuprofilename, "cpuprofile", "", "File name for a CPU profile")
	flag.StringVar(&playerType, "player", "console", "Player type. One of: console, random, greedy")
	flag.StringVar(&runnerType, "runner", "untilwin", "Runner type. One of: single, untilwin, parallel, multiwin")
	flag.Uint64Var(&width, "width", 4, "Width of the playing board.")
	flag.Uint64Var(&height, "height", 4, "Height of the playing board.")
	flag.Uint64Var(&limitPower, "lim", 11, "Power of 2 to set as the winning number. Default gives 2048.")
	flag.Uint64Var(&numAdds, "adds", 1, "The number of random values to add after each move")

	flag.Uint64Var(&numWins, "numwins", 10, "Number of wins to get when using multiwin runner")

	flag.Parse()

	if cpuprofilename != "" {
		cpuprofile, err := os.Create(cpuprofilename)
		if err != nil {
			panic(err)
		}

		pprof.StartCPUProfile(cpuprofile)
		defer pprof.StopCPUProfile()
	}

	var p players.Player

	switch playerType {
	case "console":
		p = console.Play

	case "random":
		p = random.Play

	case "greedy":
		p = greedy.Play
	}

	var runner runners.Runner

	switch runnerType {
	case "single":
		runner = single.Run

	case "untilwin":
		runner = untilwin.Run

	case "parallel":
		runner = parallel.Run

	case "multiwin":
		runner = multiwin.New(numWins)
	}

	gg := func() *game.Game {
		return game.NewGame(width, height, limitPower, numAdds)
	}

	runner(gg, p)
}
