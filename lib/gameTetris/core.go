package gameTetris

import (
	"math"
	"math/rand"
	"time"
)

const (
	emptyTile = 0
	tetriNum  = 4
	rowEmpty  = -2
	rowFull   = -1
)

// operation type (Chapter 4)
const (
	moveLeft = iota
	moveRight
	softDrop
	hardDrop
	rotateClockwise
	rotateCounterClockwise
	hold
)

// Tetrimino facings (A 1.3)
//   _____________
//   | 0| 1| 2| 3|
//   | 4| 5| 6| 7|
//   | 8| 9| A| B|
//   | C| D| E| F|
//   -------------
var tetriminoShapes = [][]int{
	{0x6521, 0x6521, 0x6521, 0x6521}, // O-tetrimino
	{0x7654, 0xEA62, 0xBA98, 0xD951}, // I-tetrimino
	{0x6541, 0x9651, 0x9654, 0x9541}, // T-tetrimino
	{0x6542, 0xA951, 0x8654, 0x9510}, // L-tetrimino
	{0x6540, 0x9521, 0xA654, 0x9851}, // J-tetrimino
	{0x5421, 0xA651, 0x9865, 0x9540}, // S-tetrimino
	{0x6510, 0x9652, 0xA954, 0x8541}, // Z-tetrimino
}

// Game implement Game interface
type Game struct {
	playfield                                     []int
	height, bufferHeight, width                   int
	bag, nextBag                                  []int
	bagIdx                                        int
	stash                                         []int
	tetriminoX, tetriminoY, tetriminoSize         int
	tetriminoIdx, tetriminoDrct, nextTetriminoIdx int
	ghostX, ghostY                                int
	level                                         int
	fallSpeed                                     float64
	lockDelay                                     int
	score, difficult, lockDownDelay               int
	comboCounter                                  int
	dropLine                                      int

	// internal judge flag
	hardDropFlag, softDropFlag    bool
	moveFlag, spaceFlag, landFlag bool
	patternMatchFlag              bool
	tSpinFlag, backToBackFlag     bool

	// game optional enable flag
	allowSRS, allowGhost, allowHardDropOp               bool
	allowLockDownPeek, allowAboveSkyline, allowForcedUp bool
	allowTopOut, allowLockOut, allowBlockOut            bool
	// io utils
	inputCh chan int
}

// ============== Utils ====================

func (g *Game) init(w, h, difficult int) {
	rand.Seed(time.Now().UnixNano())
	g.score = 0
	g.difficult = difficult
	g.level = difficult
	g.width, g.height, g.bufferHeight = w, h, h
	g.playfield = make([]int, g.width*(g.height+g.bufferHeight))
	g.calcFallSpeed()
	g.comboCounter = -1
}

func (g *Game) calcPosOnBoard(posOnShape int) (x, y int) {
	x = g.tetriminoX + posOnShape%tetriNum
	y = g.tetriminoY - posOnShape/tetriNum%tetriNum
	return x, y
}

// calculate the fall speed in current level (unit: Millisecond Per Line)
func (g *Game) calcFallSpeed() {
	g.fallSpeed = math.Pow(0.8-float64(g.level-1)*0.007, float64(g.level-1)) * 1000
}

// calculate the soft drop speed in current level (unit: Millisecond Per Line)
func (g *Game) calcDropSpeed() {
	g.calcFallSpeed()
	g.fallSpeed = g.fallSpeed / 20
}

// get from bag system (A 1.2.1)
func (g *Game) genFromBag() {
	if g.bagIdx == len(g.bag) {
		for i := 0; i < len(tetriminoShapes); i++ {
			g.bag = append(g.bag, i)
		}
		rand.Shuffle(
			len(g.bag),
			func(i, j int) {
				g.bag[i], g.bag[j] = g.bag[j], g.bag[i]
			})
		g.bagIdx = 0
	}
	g.tetriminoIdx = g.bag[g.bagIdx]
	if g.bagIdx == len(g.bag)-1 {
		g.nextTetriminoIdx = g.nextBag[0]
	} else {
		g.nextTetriminoIdx = g.bag[g.bagIdx+1]
	}
}

func (g *Game) checkLanded() bool {
	g.tetriminoY--
	for i := tetriminoShapes[g.tetriminoIdx][g.tetriminoDrct]; i != 0; i >>= 4 {
		x, y := g.calcPosOnBoard(i)
		if g.playfield[x+y*g.width] != 0 {
			g.tetriminoY++
			return false
		}
	}
	g.tetriminoY++
	return true
}

func (g *Game) processInput() {
	for input := range g.inputCh {
		switch input {
		case moveLeft:
			g.move(true)
		case moveRight:
			g.move(false)
		case rotateClockwise:
			g.rotate(true)
		case rotateCounterClockwise:
			g.rotate(false)
		case softDrop:
			g.softDrop()
		case hardDrop:
			g.hardDropFlag = true
		}
	}
}

// =============== Basic Operation =================
func (g *Game) move(isDrctLeft bool) {
	if isDrctLeft {
		g.tetriminoX--
	} else {
		g.tetriminoX++
	}
	for i := tetriminoShapes[g.tetriminoIdx][g.tetriminoDrct]; i != 0; i >>= 4 {
		x, y := g.calcPosOnBoard(i)
		if x < 0 || x >= g.width || g.playfield[x+y*g.width] != 0 {
			if isDrctLeft {
				g.tetriminoX++
			} else {
				g.tetriminoX--
			}
			return
		}
	}
}

func (g *Game) rotate(isClockwise bool) {
	if g.allowSRS {
		g.superRotate(isClockwise)
	} else {
		g.classiscRotate(isClockwise)
	}
}

func (g *Game) classiscRotate(isClockwise bool) {}

func (g *Game) superRotate(isClockwise bool) {}

func (g *Game) softDrop() {
	if g.softDropFlag {
		g.calcFallSpeed()
	} else {
		g.calcDropSpeed()
	}
}

func (g *Game) hardDrop() {
	g.tetriminoX, g.tetriminoY = g.ghostX, g.ghostY
	g.hardDropFlag = true
}

// ============ Running Flowchart ================

// generration Phase (A 1.2.1)
// return value is the gameOverFlag
func (g *Game) generationPhase() bool {
	// Random Generation
	g.genFromBag()
	// Generation of Tetriminos
	// TODO: delay time?
	// Starting Location and Orirntation
	g.tetriminoX = (g.width - tetriNum) / 2
	g.tetriminoY = g.height + 1
}

func (g *Game) fallingPhase() {
	// FIXME: should not block in here
	if !g.checkLanded() {
		startTime, endTime := time.Now(), time.Now()
		for endTime.Sub(startTime) < time.Duration(g.fallSpeed)*time.Millisecond {
			g.processInput()
			if g.hardDropFlag && !g.allowHardDropOp {
				return
			}
		}
	}
	g.lockPhase()
}

//
// return value is the gameOverFlag
func (g *Game) lockPhase() {
	// FIXME: should not block in here
	startTime, endTime := time.Now(), time.Now()
	for endTime.Sub(startTime) < time.Duration(g.lockDownDelay)*time.Millisecond {
		g.processInput()
		if g.hardDropFlag && !g.allowHardDropOp {
			return
		}
	}
	if g.moveFlag {
		if g.spaceFlag {
			g.fallingPhase()
		} else if g.landFlag {
			g.lockPhase()
		}
	}
}

func (g *Game) patternPhase() {
	g.patternMatchFlag = false
	for i := 0; i < len(g.playfield); i++ {
		// TODO: judge T-spin
		count := g.width
		if i%g.width == 0 {
			if count == 0 {
				g.patternMatchFlag = true
				g.playfield[i-g.width] = rowFull
			} else if count == g.width {
				g.playfield[i-g.width] = rowEmpty
				return // above this line is all empty
			}
			count = g.width
		}
		if g.playfield[i] > 0 {
			count--
		}
	}
}

// this phase is for more variants and is not used for now
func (g *Game) iteratePhase() {}

// this phase consumes no apparent game time
func (g *Game) animatePhase() {}

func (g *Game) elimatePhase() {
	clearLineCount := 0
	for rowIdx := 0; rowIdx < g.height+g.bufferHeight; rowIdx++ {
		rowHeaderPos := rowIdx * g.width
		if g.playfield[rowHeaderPos] == rowFull {
			clearLineCount++
		} else if g.playfield[rowHeaderPos] == rowEmpty {
			return // no need to move empty line
		} else if clearLineCount > 0 {
			for colIdx := rowHeaderPos; colIdx < rowHeaderPos+g.width; colIdx++ {
				g.playfield[colIdx] = g.playfield[colIdx-clearLineCount*g.width]
			}
		}
	}
	// TODO: Game Statistics

}

func (g *Game) completionPhase() {
	// update information
	// level up condition
}

// Run the game
// Game Over condition occurs in Generation Phase and Lock Phase
func (g *Game) Run(w, h, difficult int) {
	g.init(w, h, difficult)
	// Tetris engine flowchart
	for true {
		g.generationPhase()
		for true {
			g.fallingPhase()
			if g.hardDropFlag && !g.allowHardDropOp {
				break
			}
			for g.moveFlag && !g.spaceFlag && g.landFlag {
				g.lockPhase()
			}
			if !g.moveFlag || (g.moveFlag && !g.spaceFlag && !g.landFlag) {
				break
			}
		}
		g.patternPhase()
		if g.patternMatchFlag {
			// Mark block for destruction
		}
		g.iteratePhase()
		g.animatePhase()
		g.elimatePhase()
		g.completionPhase()
	}

	// Game Over Events

}
