package greedy

import (
	crand "crypto/rand"
	"encoding/binary"
	"math"
	"math/rand"

	"github.com/dystopium/2048/game"
)

func newRand() *rand.Rand {
	b := make([]byte, 8)
	crand.Read(b)
	seed := binary.LittleEndian.Uint64(b)
	return rand.New(rand.NewSource(int64(seed)))
}

// Play allows the human to play
func Play(g *game.Game) {
	rnd := newRand()

	for g.State() == game.StatePlaying {
		var dirs []game.Direction

		curBoard := g.Board()
		height := uint64(len(curBoard))
		width := uint64(len(curBoard[0]))
		var maxCombines uint64

		//fmt.Println(g)

		// up
		var totalCombines uint64
		for col := uint64(0); col < width; col++ {
			for row := height - 1; row > 0 && row != math.MaxUint64; row-- {
				rownext := row - 1

				// find the next non-blank square
				for rownext != math.MaxUint64 && curBoard[rownext][col] == 0 {
					rownext--
				}

				if rownext != math.MaxUint64 && curBoard[row][col] == curBoard[rownext][col] && curBoard[row][col] != 0 {
					totalCombines++
					row--
				}
			}
		}

		//fmt.Println("up", totalCombines)

		maxCombines = totalCombines

		if totalCombines > 0 {
			dirs = append(dirs, game.DirUp)
		}

		// down
		totalCombines = 0
		for col := uint64(0); col < width; col++ {
			for row := uint64(0); row < height-2; row++ {
				rownext := row + 1

				// find the next non-blank square
				for rownext < height && curBoard[rownext][col] == 0 {
					rownext++
				}

				if rownext < height && curBoard[row][col] == curBoard[rownext][col] && curBoard[row][col] != 0 {
					totalCombines++
					row++
				}
			}
		}

		//fmt.Println("down", totalCombines)

		if totalCombines > 0 && totalCombines == maxCombines {
			dirs = append(dirs, game.DirDown)
		} else if totalCombines > maxCombines {
			maxCombines = totalCombines
			dirs = []game.Direction{game.DirDown}
		}

		// left
		totalCombines = 0
		for row := uint64(0); row < height; row++ {
			for col := width - 1; col > 0 && col != math.MaxUint64; col-- {
				colnext := col - 1

				// find the next non-blank square
				for colnext != math.MaxUint64 && curBoard[row][colnext] == 0 {
					colnext--
				}

				if colnext != math.MaxUint64 && curBoard[row][col] == curBoard[row][colnext] && curBoard[row][col] != 0 {
					totalCombines++
					col--
				}
			}
		}

		//fmt.Println("left", totalCombines)

		if totalCombines > 0 && totalCombines == maxCombines {
			dirs = append(dirs, game.DirLeft)
		} else if totalCombines > maxCombines {
			maxCombines = totalCombines
			dirs = []game.Direction{game.DirLeft}
		}

		// right
		totalCombines = 0
		for row := uint64(0); row < height; row++ {
			for col := uint64(0); col < width-2; col++ {
				colnext := col + 1

				// find the next non-blank square
				for colnext < width && curBoard[row][colnext] == 0 {
					colnext++
				}

				if colnext < width && curBoard[row][col] == curBoard[row][col+1] && curBoard[row][col] != 0 {
					totalCombines++
					col++
				}
			}
		}

		//fmt.Println("right", totalCombines)

		if totalCombines > 0 && totalCombines == maxCombines {
			dirs = append(dirs, game.DirRight)
		} else if totalCombines > maxCombines {
			maxCombines = totalCombines
			dirs = []game.Direction{game.DirRight}
		}

		//fmt.Println(dirs)

		if len(dirs) == 0 {
			dirs = []game.Direction{game.DirUp, game.DirDown, game.DirLeft, game.DirRight}
		}

		dir := dirs[rnd.Intn(len(dirs))]

		//fmt.Println(dir)
		g.Move(dir)
	}
}
