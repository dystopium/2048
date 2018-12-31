package game

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"strings"
)

// Game represents the game state
type Game struct {
	board  [][]int
	width  int
	height int
	rnd    *rand.Rand
	score  int
}

func newRand() *rand.Rand {
	b := make([]byte, 8)
	crand.Read(b)
	seed := binary.LittleEndian.Uint64(b)
	return rand.New(rand.NewSource(int64(seed)))
}

// NewGame creates a new game board with two randomly placed tiles
func NewGame(width, height int) *Game {
	board := make([][]int, height)
	for i := range board {
		board[i] = make([]int, width)
	}

	g := &Game{
		board:  board,
		width:  width,
		height: height,
		rnd:    newRand(),
	}

	g.placeNew()
	g.placeNew()

	return g
}

func (g *Game) Score() int {
	return g.score
}

func (g *Game) getPos(row, col int) int {
	return g.board[row][col]
}

func (g *Game) setPos(row, col, newValue int) {
	g.board[row][col] = newValue
}

// placeNew places a new random tile on the board
// 90% chance for a 2, 10% chance for a 4
func (g *Game) placeNew() {
	// first get the random value to place
	newValue := 2

	if g.rnd.Float64() > 0.9 {
		newValue = 4
	}

	for {
		row := g.rnd.Intn(g.height)
		col := g.rnd.Intn(g.width)

		// if the board gets full this might take a few tries
		// if the board is large, this might take a LOT of tries and
		//   it can be optimized to only pick from a list of open squares
		if g.getPos(row, col) == 0 {
			g.setPos(row, col, newValue)
			break
		}
	}
}

// Direction is a strongly typed way to represent a move direction
type Direction byte

const (
	dirInvalid Direction = iota

	// DirUp represents the up direction
	DirUp

	// DirDown represents the down direction
	DirDown

	// DirLeft represents the left direction
	DirLeft

	// DirRight represents the right direction
	DirRight
)

// Move makes a move in the direction specified
// Moves tiles, combines ones that can be combined, and adds a new random tile
func (g *Game) Move(dir Direction) {
	switch dir {
	case DirUp:
		for col := 0; col < g.width; col++ {
			var cur []int

			for row := 0; row < g.height; row++ {
				if g.getPos(row, col) != 0 {
					cur = append(cur, g.getPos(row, col))
				}
			}

			if len(cur) == 0 {
				continue
			}

			new, score := consolidate(cur, g.height)
			g.score += score

			for row := 0; row < g.height; row++ {
				g.setPos(row, col, new[row])
			}
		}

	case DirDown:
		for col := 0; col < g.width; col++ {
			var cur []int

			for row := g.height - 1; row >= 0; row-- {
				if g.getPos(row, col) != 0 {
					cur = append(cur, g.getPos(row, col))
				}
			}

			if len(cur) == 0 {
				continue
			}

			new, score := consolidate(cur, g.height)
			g.score += score

			for row := 0; row < g.height; row++ {
				g.setPos(row, col, new[g.height-1-row])
			}
		}

	case DirLeft:
		for row := 0; row < g.height; row++ {
			var cur []int

			for col := 0; col < g.width; col++ {
				if g.getPos(row, col) != 0 {
					cur = append(cur, g.getPos(row, col))
				}
			}

			if len(cur) == 0 {
				continue
			}

			new, score := consolidate(cur, g.width)
			g.score += score

			for col := 0; col < g.height; col++ {
				g.setPos(row, col, new[col])
			}
		}

	case DirRight:
		for row := 0; row < g.height; row++ {
			var cur []int

			for col := g.width - 1; col >= 0; col-- {
				if g.getPos(row, col) != 0 {
					cur = append(cur, g.getPos(row, col))
				}
			}

			if len(cur) == 0 {
				continue
			}

			new, score := consolidate(cur, g.width)
			g.score += score

			for col := 0; col < g.width; col++ {
				g.setPos(row, col, new[g.width-1-col])
			}
		}
	}

	g.placeNew()
}

// Consolidates pairs of numbers and returns a new, padded slice and the added score
// assumes the slice is more than 0 length
func consolidate(cur []int, padTo int) ([]int, int) {
	if len(cur) == 0 {
		panic("SHORT SLICE")
	}

	var new []int
	var score int

	if len(cur) == 1 {
		new = cur

	} else {
		// combine numbers that are the same by twos
		i := 0
		for {
			if i > len(cur)-1 {
				break
			}

			if i == len(cur)-1 || cur[i] != cur[i+1] {
				new = append(new, cur[i])

			} else {
				new = append(new, cur[i]*2)
				score += cur[i] * 2
				i++
			}

			i++
		}
	}

	for len(new) < padTo {
		new = append(new, 0)
	}

	return new, score
}

func (g *Game) String() string {
	bdr := &strings.Builder{}

	// Write the top border row
	bdr.WriteRune('╔')

	for i := 0; i < g.width; i++ {
		bdr.WriteString("════")

		if i < g.width-1 {
			bdr.WriteRune('╤')
		} else {
			bdr.WriteRune('╗')
		}
	}

	bdr.WriteRune('\n')

	for ridx, row := range g.board {
		// Write out the values in the row
		bdr.WriteRune('║')

		for cidx, col := range row {
			if col > 0 {
				bdr.WriteString(fmt.Sprintf("%4d", col))
			} else {
				bdr.WriteString("    ")
			}

			if cidx < g.width-1 {
				bdr.WriteRune('│')
			} else {
				bdr.WriteRune('║')
			}
		}

		bdr.WriteRune('\n')

		// Write out the next row, which might be a middle separator or the bottom row
		if ridx < g.height-1 {
			bdr.WriteRune('╟')

			for i := 0; i < g.width; i++ {
				bdr.WriteString("────")

				if i < g.width-1 {
					bdr.WriteRune('┼')
				} else {
					bdr.WriteRune('╢')
				}
			}

		} else {
			bdr.WriteRune('╚')

			for i := 0; i < g.width; i++ {
				bdr.WriteString("════")

				if i < g.width-1 {
					bdr.WriteRune('╧')
				} else {
					bdr.WriteRune('╝')
				}
			}
		}

		bdr.WriteRune('\n')

	}

	return bdr.String()
}

// WriteTo writes the string representation of the board to the given io.Writer
func (g *Game) WriteTo(w io.Writer) (int64, error) {
	temp, temp2 := w.Write([]byte(g.String()))
	return int64(temp), temp2
}
