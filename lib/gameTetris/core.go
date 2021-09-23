package gameTetris

import (
	"math"
	"math/rand"
	"time"
)

const (
	emptyTile = 0
	tetriNum  = 4
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
	stash                                         []int
	bag, nextBag                                  []int
	bagIdx                                        int
	tetriminoX, tetriminoY, tetriminoSize         int
	tetriminoIdx, tetriminoDrct, nextTetriminoIdx int
	ghostX, ghostY                                int
	level                                         int
	fallSpeed                                     float64
	lockDelay                                     int
	score, difficult, lockDownDelay               int
	alive                                         bool
	hardDropFlag, softDropFlag                    bool
	movedFlag, spaceFlag, landFlag                bool
	patternMatchFlag                              bool
	allowSRS, allowGhost, allowHardDropOp         bool
	inputCh                                       chan int
}

// ============== Main Progress ================

func (g *Game) init(w, h, difficult int) {
	rand.Seed(time.Now().UnixNano())
	g.score = 0
	g.difficult = difficult
	g.level = difficult
	g.width, g.height, g.bufferHeight = w, h, h
	g.playfield = make([]int, g.width*(g.height+g.bufferHeight))
	g.calcFallSpeed()

}

func (g *Game) calcPosOnBoard(posOnShape int) (x, y int) {
	x = g.tetriminoX + posOnShape%tetriNum
	y = g.tetriminoY - posOnShape/tetriNum%tetriNum
	return x, y
}

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

func (g *Game) superRotate(isClockwise bool) {

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
	return true
}

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

func (g *Game) accelerate() {}

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
		g.nextTetriminoIdx = tetriminoShapes[g.nextBag[0]][0]
	} else {
		g.nextTetriminoIdx = tetriminoShapes[g.bag[g.bagIdx+1]][0]
	}
}

// generration Phase (A 1.2.1)
func (g *Game) generationPhase() {
	// Random Generation
	g.genFromBag()
	// Generation of Tetriminos
	// TODO: delay time?
	// Starting Location and Orirntation
	g.tetriminoX = (g.width - tetriNum) / 2
	g.tetriminoY = g.height + 1
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

func (g *Game) fallingPhase() {
	if !g.checkLanded() {
		startTime, endTime := time.Now(), time.Now()
		for endTime.Sub(startTime) < time.Duration(g.fallSpeed)*time.Millisecond {
			g.processInput()
			if g.hardDropFlag && !g.allowHardDropOp {
				return
			}
		}
	}

	g.tetriminoY--

	g.lockPhase()
}

func (g *Game) lockPhase() {
	startTime, endTime := time.Now(), time.Now()
	for endTime.Sub(startTime) < time.Duration(g.lockDownDelay)*time.Millisecond {
		g.processInput()
		if g.hardDropFlag && !g.allowHardDropOp {
			return
		}
	}
	if g.movedFlag {
		if g.spaceFlag {
			g.fallingPhase()
		} else if g.landFlag {
			g.lockPhase()
		}
	}
}

func (g *Game) patternPhase() {
	if g.patternMatchFlag {
		// Mark Block for Destruction
	}
	g.iteratePhase()
}

func (g *Game) iteratePhase() {

}

func (g *Game) animatePhase() {}

func (g *Game) detectLineClear() bool {
	return false // FIXEME: not done
}

func (g *Game) elimatePhase() {
	if g.detectLineClear() {
		g.elimatePhase()
	}
	g.completionPhase()
}

func (g *Game) completionPhase() {}

// calculate the fall speed in current level (unit: Millisecond Per Line)
func (g *Game) calcFallSpeed() {
	g.fallSpeed = math.Pow(0.8-float64(g.level-1)*0.007, float64(g.level-1)) * 1000
}

func (g *Game) calcDropSpeed() {
	g.calcFallSpeed()
	g.fallSpeed = g.fallSpeed / 20
}

// Run the game
func (g *Game) Run() {
	// Tetris engine flowchart
	for true {
		g.generationPhase()
		for true {
			g.fallingPhase()
			if g.hardDropFlag && !g.allowHardDropOp {
				break
			}
			for g.movedFlag && !g.spaceFlag && g.landFlag {
				g.lockPhase()
			}
			if !g.movedFlag || (g.movedFlag && !g.spaceFlag && !g.landFlag) {
				break
			}
		}
		g.patternPhase()
		if g.patternMatchFlag {
			// Mark block for destruction
		}
		g.iteratePhase()
		g.animatePhase()
		for g.detectLineClear() {
			g.elimatePhase()
		}
		g.completionPhase()
	}
}
