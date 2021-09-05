package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const (
	keyEsc = 0

	keyUp    = 1
	keyDown  = 2
	keyLeft  = 4
	keyRight = 8
)

var colorMap = []termbox.Attribute{
	termbox.ColorRed,
	termbox.ColorRed,
	termbox.ColorYellow,
	termbox.ColorGreen,
	termbox.ColorCyan,
	termbox.ColorBlue,
	termbox.ColorMagenta,
	termbox.ColorDarkGray,
	termbox.ColorBlack,
	termbox.ColorBlack,
	termbox.ColorBlack,
	termbox.ColorBlack,
	termbox.ColorBlack,
	termbox.ColorBlack,
	termbox.ColorBlack,
	termbox.ColorBlack,
	termbox.ColorBlack,
}

var strMap = []string{
	"     ",
	"  2  ",
	"  4  ",
	"  8  ",
	"  16 ",
	"  32 ",
	"  64 ",
	" 128 ",
	" 256 ",
	" 512 ",
	" 1024",
	" 2048",
	" 4096",
	" 8192",
	"16384",
	"32768",
	"65536",
}

const (
	emptyElement = 0
	firstElement = 0
)

// Game2048 implement Game interface
type Game2048 struct {
	score int
	fps   int

	board                 []int
	height, width         int
	up, down, left, right int
	difficult             int
	alive                 bool
	boardFree             int
	lastMoveValid         bool
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
func (g *Game2048) init(w, h, difficult int) {
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

func (g *Game2048) moveOrMergeElement(pos, direction int) {
	newPos := pos + direction

	if g.board[pos] == emptyElement || !g.checkBorderAfterMove(pos, direction) {
		// fmt.Printf("%d cant be moved/merged to %d\n", pos, newPos)
		return // can not be moved, end for next element
	}

	if g.board[newPos] == g.board[pos] { // can be merged
		// fmt.Printf("merge:%d->%d\n", pos, newPos)
		// g.Score += g.board[pos]
		g.board[newPos]++                  // promote the target element
		g.board[newPos] = -g.board[newPos] // mark the target not be merged again
		g.board[pos] = emptyElement
		g.boardFree++
		g.lastMoveValid = true
	}

	if g.board[newPos] == emptyElement { // can be moved
		// fmt.Printf("move:%d->%d\n", pos, newPos)
		g.board[newPos] = g.board[pos] // assign the target element
		g.board[pos] = emptyElement
		g.moveOrMergeElement(newPos, direction) // (iteration)
		g.lastMoveValid = true
	}

}

func (g *Game2048) operate(direction int) {
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
				newNumIndexCountDown = -1                 // make count down negatively to avoid generating another new number
			} else if newNumIndexCountDown > 0 { // wait for next free pos
				newNumIndexCountDown--
			}
		}
	}

	g.checkAlive()
}

func (g *Game2048) checkAlive() {
	if !g.lastMoveValid || g.boardFree > 0 { // if not valid (never move or merge this turn), end for next operation
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
	g.lastMoveValid = false
}

func (g *Game2048) checkBorderAfterMove(pos, direction int) bool {
	newPos := afterMove(pos, direction)
	if newPos < firstElement || newPos >= len(g.board) {
		return false // out of vertical border
	}
	if (direction == g.left || direction == g.right) && !g.isSameLine(pos, newPos) {
		return false // out of horizon border
	}
	return true
}

func (g *Game2048) isSameLine(pos1, pos2 int) bool {
	rowAtPosion1 := pos1 / g.width
	rowAtPosion2 := pos2 / g.width
	return rowAtPosion1 == rowAtPosion2
}

// ============== Utils ====================

func afterMove(i, direction int) int {
	return i + direction
}

// ============== Implements ================
// Start implement Game interface to start game
func (g *Game2048) Run(
	name string,
	width, height, difficult int,
	inputChannel chan int,
	renderFunc func(board []int, height, width, score, fps int)) {
	g.init(4, 4, 2)
	stopRenderCh := make(chan struct{})
	go g.Render(renderFunc, stopRenderCh)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go g.Process(&wg, inputChannel)
	wg.Wait()
	close(stopRenderCh)
}

//
func (g *Game2048) Process(wg *sync.WaitGroup, inputCh chan int) {
	for x := range inputCh {
		switch x {
		case keyUp:
			g.operate(g.up)
		case keyDown:
			g.operate(g.down)
		case keyLeft:
			g.operate(g.left)
		case keyRight:
			g.operate(g.right)
		case keyEsc:
			wg.Done()
			return
		}
	}
}

func (g *Game2048) Render(
	renderFunc func(board []int, height, width, score, fps int),
	stopCh <-chan struct{}) {

	for {
		select {
		case <-stopCh:
			return
		default:
			renderFunc(g.board, g.height, g.width, g.score, g.fps)
		}

	}
}

//  =================== Utils ===================
func listenToInput(input chan int) {
	termbox.SetInputMode(termbox.InputEsc)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowLeft:
				input <- keyLeft
			case termbox.KeyArrowDown:
				input <- keyDown
			case termbox.KeyArrowRight:
				input <- keyRight
			case termbox.KeyArrowUp:
				input <- keyUp
			case termbox.KeyEsc:
				input <- keyEsc
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func render(board []int, height, width, score, fps int) {
	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		panic(err)
	}

	// termbox.SetCursor(0, 0)
	tbprint(0, 0, termbox.ColorWhite, termbox.ColorBlack, "")

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			// fmt.Printf("(%2d,%2d) color[%2d] %s ", (1+width)*j, (1+height)*i, i*width+j, strMap[board[i*width+j]])
			tbprint(6*j+1, i+1, colorMap[board[i*width+j]], termbox.ColorBlack, strMap[board[i*width+j]])
		}
		// fmt.Println()
	}

	if err := termbox.Flush(); err != nil {
		panic(err)
	}

	// time.Sleep(time.Duration(1) * time.Second)
	time.Sleep(time.Duration(1000/fps) * time.Millisecond)
}

// This function is often useful:
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	inputChannel := make(chan int, 5)
	go listenToInput(inputChannel)

	g := Game2048{}
	g.Run("chow", 4, 4, 2, inputChannel, render)
}
