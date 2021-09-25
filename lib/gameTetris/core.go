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

const (
	defaultDifficulty   = 1
	defaultHeight       = 20
	defaultBufferHeight = 20
	defaultWidth        = 10
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

// GameManager implement GameManager interface
type GameManager struct {
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
	score, difficulty, lockDownDelay              int
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

func NewGameManager() *GameManager {
	return &GameManager{}
}

// ============== Utils ====================

func (gm *GameManager) init() {
	rand.Seed(time.Now().UnixNano())
	gm.score = 0
	gm.difficulty, gm.level = defaultDifficulty, defaultDifficulty
	gm.width, gm.height, gm.bufferHeight = defaultWidth, defaultHeight, defaultBufferHeight
	gm.playfield = make([]int, gm.width*(gm.height+gm.bufferHeight))
	gm.calcFallSpeed()
	gm.comboCounter = -1
}

func (gm *GameManager) calcPosOnBoard(posOnShape int) (x, y int) {
	x = gm.tetriminoX + posOnShape%tetriNum
	y = gm.tetriminoY - posOnShape/tetriNum%tetriNum
	return x, y
}

// calculate the fall speed in current level (unit: Millisecond Per Line)
func (gm *GameManager) calcFallSpeed() {
	gm.fallSpeed = math.Pow(0.8-float64(gm.level-1)*0.007, float64(gm.level-1)) * 1000
}

// calculate the soft drop speed in current level (unit: Millisecond Per Line)
func (gm *GameManager) calcDropSpeed() {
	gm.calcFallSpeed()
	gm.fallSpeed = gm.fallSpeed / 20
}

// get from bag system (A 1.2.1)
func (gm *GameManager) genFromBag() {
	if gm.bagIdx == len(gm.bag) {
		for i := 0; i < len(tetriminoShapes); i++ {
			gm.bag = append(gm.bag, i)
		}
		rand.Shuffle(
			len(gm.bag),
			func(i, j int) {
				gm.bag[i], gm.bag[j] = gm.bag[j], gm.bag[i]
			})
		gm.bagIdx = 0
	}
	gm.tetriminoIdx = gm.bag[gm.bagIdx]
	if gm.bagIdx == len(gm.bag)-1 {
		gm.nextTetriminoIdx = gm.nextBag[0]
	} else {
		gm.nextTetriminoIdx = gm.bag[gm.bagIdx+1]
	}
}

func (gm *GameManager) checkCollision() bool {
	for i := tetriminoShapes[gm.tetriminoIdx][gm.tetriminoDrct]; i != 0; i >>= 4 {
		x, y := gm.calcPosOnBoard(i)
		if gm.playfield[x+y*gm.width] != 0 {
			return false
		}
	}
	return false //
}

func (gm *GameManager) checkLanding() bool {
	gm.tetriminoY--
	for i := tetriminoShapes[gm.tetriminoIdx][gm.tetriminoDrct]; i != 0; i >>= 4 {
		x, y := gm.calcPosOnBoard(i)
		if y >= 0 && gm.playfield[x+y*gm.width] != 0 {
			gm.tetriminoY++
			return false
		}
	}
	gm.tetriminoY++
	return true
}

func (gm *GameManager) processInput() {
	for input := range gm.inputCh {
		switch input {
		case moveLeft:
			gm.move(true)
		case moveRight:
			gm.move(false)
		case rotateClockwise:
			gm.rotate(true)
		case rotateCounterClockwise:
			gm.rotate(false)
		case softDrop:
			gm.softDrop()
		case hardDrop:
			gm.hardDropFlag = true
		}
	}
}

// =============== Basic Operation =================
func (gm *GameManager) move(isDrctLeft bool) {
	if isDrctLeft {
		gm.tetriminoX--
	} else {
		gm.tetriminoX++
	}
	for i := tetriminoShapes[gm.tetriminoIdx][gm.tetriminoDrct]; i != 0; i >>= 4 {
		x, y := gm.calcPosOnBoard(i)
		if x < 0 || x >= gm.width || gm.playfield[x+y*gm.width] != 0 {
			if isDrctLeft {
				gm.tetriminoX++
			} else {
				gm.tetriminoX--
			}
			return
		}
	}
}

func (gm *GameManager) rotate(isClockwise bool) {
	if gm.allowSRS {
		gm.superRotate(isClockwise)
	} else {
		gm.classiscRotate(isClockwise)
	}
}

func (gm *GameManager) classiscRotate(isClockwise bool) {}

func (gm *GameManager) superRotate(isClockwise bool) {}

func (gm *GameManager) softDrop() {
	if gm.softDropFlag {
		gm.calcFallSpeed()
	} else {
		gm.calcDropSpeed()
	}
}

func (gm *GameManager) hardDrop() {
	gm.tetriminoX, gm.tetriminoY = gm.ghostX, gm.ghostY
	gm.hardDropFlag = true
}

// ============ Running Flowchart ================

// generration Phase (A 1.2.1)
// return value is the gameOverFlag
func (gm *GameManager) generationPhase() bool {
	// Random Generation
	gm.genFromBag()
	// Generation of Tetriminos
	// TODO: delay time?
	// Starting Location and Orirntation
	gm.tetriminoX = (gm.width - tetriNum) / 2
	gm.tetriminoY = gm.height + 1
	return gm.checkCollision()
}

func (gm *GameManager) fallingPhase() {
	// FIXME: should not block in here
	if !gm.checkLanding() {
		startTime, endTime := time.Now(), time.Now()
		for endTime.Sub(startTime) < time.Duration(gm.fallSpeed)*time.Millisecond {
			gm.processInput()
			if gm.hardDropFlag && !gm.allowHardDropOp {
				return
			}
		}
	}
}

//
// return value is the gameOverFlag
func (gm *GameManager) lockPhase() bool {
	// FIXME: should not block in here
	startTime, endTime := time.Now(), time.Now()
	for endTime.Sub(startTime) < time.Duration(gm.lockDownDelay)*time.Millisecond {
		gm.processInput()
		if gm.hardDropFlag && !gm.allowHardDropOp {
			return gm.checkCollision()
		}
	}
	return gm.checkCollision()
}

func (gm *GameManager) patternPhase() {
	gm.patternMatchFlag = false
	for i := 0; i < len(gm.playfield); i++ {
		// TODO: judge T-spin
		count := gm.width
		if i%gm.width == 0 {
			if count == 0 {
				gm.patternMatchFlag = true
				gm.playfield[i-gm.width] = rowFull
			} else if count == gm.width {
				gm.playfield[i-gm.width] = rowEmpty
				return // above this line is all empty
			}
			count = gm.width
		}
		if gm.playfield[i] > 0 {
			count--
		}
	}
}

// this phase is for more variants and is not used for now
func (gm *GameManager) iteratePhase() {}

// this phase consumes no apparent game time
func (gm *GameManager) animatePhase() {}

func (gm *GameManager) elimatePhase() {
	clearLineCount := 0
	for rowIdx := 0; rowIdx < gm.height+gm.bufferHeight; rowIdx++ {
		rowHeaderPos := rowIdx * gm.width
		if gm.playfield[rowHeaderPos] == rowFull {
			clearLineCount++
		} else if gm.playfield[rowHeaderPos] == rowEmpty {
			return // no need to move empty line
		} else if clearLineCount > 0 {
			for colIdx := rowHeaderPos; colIdx < rowHeaderPos+gm.width; colIdx++ {
				gm.playfield[colIdx] = gm.playfield[colIdx-clearLineCount*gm.width]
			}
		}
	}
	// TODO: GameManager Statistics
}

func (gm *GameManager) completionPhase() {
	// update information
	// level up condition
}

// Tetris engine flowchart
func (gm *GameManager) loopFlow() {
	for gm.generationPhase() {
		for true {
			gm.fallingPhase()
			if (gm.hardDropFlag && !gm.allowHardDropOp) ||
				(gm.moveFlag && !gm.spaceFlag && !gm.landFlag) ||
				!gm.moveFlag {
				break
			}
			for gm.moveFlag && !gm.spaceFlag && gm.landFlag {
				if !gm.lockPhase() {
					return
				}
			}

		}
		gm.patternPhase()
		if gm.patternMatchFlag {
			// Mark block for destruction
		}
		gm.iteratePhase()
		gm.animatePhase()
		gm.elimatePhase()
		gm.completionPhase()
	}
}

// ============= Export Function ===============

func (gm *GameManager) GetSetups() {

}

// Setup game optional settings
func (gm *GameManager) Setup() {}

// Run the game
// GameManager Over condition occurs in Generation Phase and Lock Phase
func (gm *GameManager) Run() {
	gm.init()
	// Tetris engine flowchart
	gm.loopFlow()
	// GameManager Over Events
}
