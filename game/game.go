package game

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"strings"
)

type State byte

const (
	stateInvalid State = iota
	StatePlaying
	StateWon
	StateLost
)

func (s State) String() string {
	switch s {
	case StatePlaying:
		return "Playing"
	case StateWon:
		return "Won"
	case StateLost:
		return "Lost"
	default:
		return "?invalid state?"
	}
}

// Game represents the game state
type Game struct {
	board  [][]int
	width  int
	height int
	limit  int
	rnd    *rand.Rand
	score  int
	moves  int
	state  State
}

func newRand() *rand.Rand {
	b := make([]byte, 8)
	crand.Read(b)
	seed := binary.LittleEndian.Uint64(b)
	return rand.New(rand.NewSource(int64(seed)))
}

// NewGame creates a new game board with two randomly placed tiles
// limitPower is the power of 2 required to win
func NewGame(width, height int, limitPower uint) *Game {
	if limitPower > 13 {
		panic("Limit powers greater than 13 are not supported")
	}

	board := make([][]int, height)
	for i := range board {
		board[i] = make([]int, width)
	}

	g := &Game{
		board:  board,
		width:  width,
		height: height,
		limit:  1 << limitPower,
		rnd:    newRand(),
		state:  StatePlaying,
	}

	g.placeNew()
	g.placeNew()

	return g
}

// TotalMoves reports the cumulative number of moves performed in the game
func (g *Game) TotalMoves() int {
	return g.moves
}

// Score reports the current game score
func (g *Game) Score() int {
	return g.score
}

func (g *Game) State() State {
	return g.state
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
	// sneaky sneaky
	if g.state != StatePlaying {
		return
	}

	old := make([][]int, g.height)
	for i := range old {
		old[i] = make([]int, g.width)
	}

	for row := 0; row < g.height; row++ {
		for col := 0; col < g.width; col++ {
			old[row][col] = g.board[row][col]
		}
	}

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

			for col := 0; col < g.width; col++ {
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

	if !same(old, g.board) {
		g.moves++
	}

	g.setWonOrLost()

	// if old and new boards are the same, don't place new
	if g.state == StatePlaying && !same(old, g.board) {
		g.placeNew()
	}
}

func same(old, new [][]int) bool {
	for row := 0; row < len(old); row++ {
		for col := 0; col < len(old[0]); col++ {
			if old[row][col] != new[row][col] {
				return false
			}
		}
	}

	return true
}

func (g *Game) setWonOrLost() {
	if isWin(g.board, g.limit) {
		g.state = StateWon
	} else if isLose(g.board) {
		g.state = StateLost
	}
}

func isWin(board [][]int, winValue int) bool {
	for _, row := range board {
		for _, val := range row {
			if val == winValue {
				return true
			}
		}
	}

	return false
}

func isLose(board [][]int) bool {
	height := len(board)
	width := len(board[0])

	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			// blank space means there's still a move
			if board[row][col] == 0 {
				return false
			}

			// up
			if row > 1 &&
				board[row][col] == board[row-1][col] {
				return false
			}

			// down
			if row < height-1 &&
				board[row][col] == board[row+1][col] {
				return false
			}

			// left
			if col > 1 &&
				board[row][col] == board[row][col-1] {
				return false
			}

			// right
			if col < width-1 &&
				board[row][col] == board[row][col+1] {
				return false
			}
		}
	}

	return true
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
