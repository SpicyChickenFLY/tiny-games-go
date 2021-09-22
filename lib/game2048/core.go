package game2048

import "math/rand"

const (
	emptyElement = 0
	firstElement = 0
)

// Game implement Game interface
type Game struct {
	score int
	fps   int

	board                 []int
	height, width         int
	up, down, left, right int
	difficult             int
	alive                 bool
	boardFree             int
	lastMoveValid         bool
	lastNewNumberIndex    int
}

// ============== Main Progress ================

// if we assum width=4, height=4
// the g.board is a 1-demension slice like
//  [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15]
// but we can regard it as follow
//  [
//    0,  1,  2,  3,
//    4,  5,  6,  7,
//    8,  9, 10, 11,
//    12, 13, 14, 15
//  ]
// and we can judge the index by following equation:
//  index = width * row_index + column_index
//  e.g. g.board[2][1] = 4 * 2 + 1 = 9 = g.board[9]
//
// so the 4 direction can be defined as integer
//  UP: 	width * -1
//  DOWN: 	width
//  LEFT: 	-1
//  RIGHT:	1
func (g *Game) init(w, h, difficult int) {
	g.score = 0
	g.lastMoveValid = false
	g.fps = 60
	g.difficult = difficult
	g.width, g.height = w, h
	g.up, g.down, g.left, g.right = -1*g.width, g.width, -1, 1
	g.board = make([]int, g.width*g.height)

	// init board with two numbers
	pos1 := rand.Intn(len(g.board))
	pos2 := rand.Intn(len(g.board) - 1)
	if pos2 >= pos1 {
		pos2++
	}
	g.board[pos1] = rand.Intn(g.difficult) + 1
	g.board[pos2] = rand.Intn(g.difficult) + 1
	g.boardFree = len(g.board) - 2
}

func (g *Game) moveOrMergeElement(pos, direction int) {
	newPos := afterMove(pos, direction)

	if g.board[pos] == emptyElement || !g.checkBorderAfterMove(pos, direction) {
		// fmt.Printf("%d cant be moved/merged to %d\n", pos, newPos)
		return // can not be moved, end for next element
	}

	if g.board[newPos] == g.board[pos] { // can be merged
		g.score += g.board[pos]
		g.score += g.board[pos]
		g.board[newPos]++                  // promote the target element
		g.board[newPos] = -g.board[newPos] // mark the target not be merged again
		g.board[pos] = emptyElement
		g.boardFree++
		g.lastMoveValid = true
	}

	if g.board[newPos] == emptyElement { // can be moved
		g.board[newPos] = g.board[pos] // assign the target element
		g.board[pos] = emptyElement
		g.moveOrMergeElement(newPos, direction) // (iteration)
		g.lastMoveValid = true
	}
}

func (g *Game) operate(direction int) {
	g.lastMoveValid = false

	// first loop: move or merge
	if direction == g.up || direction == g.left {
		// fmt.Printf("direction:%d, normal traverse\n", direction)
		for pos := firstElement; pos < len(g.board); pos++ {
			g.moveOrMergeElement(pos, direction)
		}
	} else {
		// fmt.Printf("direction:%d, reverse traverse\n", direction)
		for pos := len(g.board) - 1; pos >= firstElement; pos-- {
			g.moveOrMergeElement(pos, direction)
		}
	}

	// second loop: clean marks, judge alive, generate new number
	newNumIndexCountDown := -1
	if g.lastMoveValid {
		newNumIndexCountDown = rand.Intn(g.boardFree) // choose a free space randomly for new number
	}
	for pos := firstElement; pos < len(g.board); pos++ {
		if g.board[pos] < 0 {
			g.board[pos] = -g.board[pos] // clean this mark by making it positive
		} else if g.board[pos] == emptyElement {
			if newNumIndexCountDown == 0 { // the free pos will be assign for new number
				g.board[pos] = rand.Intn(g.difficult) + 1 // rand.Intn range from [0,n) while we need [1,n], so plus 1
				g.boardFree--                             // free space decreased because of the new number
			}
			if newNumIndexCountDown >= 0 { // wait for next free pos
				newNumIndexCountDown--
			}
		}
	}

	g.checkAlive()
}

func (g *Game) checkAlive() {
	if g.boardFree > 0 { // if not valid (never move or merge this turn), end for next operation
		g.alive = true
	} else {
		for pos := 0; pos < len(g.board)-1; pos++ { // we ignore the last element
			if pos < len(g.board)-g.width && g.board[pos] == g.board[afterMove(pos, g.down)] || // ignore last row while chechk two row
				pos%g.width < g.width-1 && g.board[pos] == g.board[afterMove(pos, g.right)] { // ignore last column while check two column
				g.alive = true
				return
			}
		}
		g.alive = false
	}
}

func (g *Game) checkBorderAfterMove(pos, direction int) bool {
	newPos := afterMove(pos, direction)
	if newPos < firstElement || newPos >= len(g.board) {
		return false // out of vertical border
	}
	if (direction == g.left || direction == g.right) && !g.isSameLine(pos, newPos) {
		return false // out of horizon border
	}
	return true
}

func (g *Game) isSameLine(pos1, pos2 int) bool {
	rowAtPosion1 := pos1 / g.width
	rowAtPosion2 := pos2 / g.width
	return rowAtPosion1 == rowAtPosion2
}

// ============== Utils ====================

func afterMove(i, direction int) int {
	return i + direction
}
