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
	defaultScore                   = 0
	defaultDifficulty              = 1
	defaultLockDownDelay           = 500
	defaultHeight                  = 20
	defaultBufferHeight            = 20
	defaultWidth                   = 10
	defaultDropSpeedRatio          = 20
	defaultAllowSRS                = true
	defaultAllowGhost              = true
	defaultAllowHardDropOp         = false
	defaultAllowLockDownPeek       = true
	defaultAllowPlayAboveSkyline   = true
	defaultAllowForcedAboveSkyline = true
	defaultAllowTopOut             = false // for now, this must be false
	defaultAllowLockOut            = false
	defaultAllowBlockOut           = false // for now, this must be false
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

const (
	tetriminoShapeO = iota
	tetriminoShapeI
	tetriminoShapeT
	tetriminoShapeL
	tetriminoShapeJ
	tetriminoShapeS
	tetriminoShapeZ
)

// Tetrimino facings (A 1.3)
//   _____________
//   | 0| 1| 2| 3|
//   | 4| 5| 6| 7|
//   | 8| 9| A| B|
//   | C| D| E| F|
//   -------------
var tetriminoShapes = [7][4]int{
	{0x6521, 0x6521, 0x6521, 0x6521}, // O-tetrimino
	{0x7654, 0xEA62, 0xBA98, 0xD951}, // I-tetrimino
	{0x6541, 0x9651, 0x9654, 0x9541}, // T-tetrimino
	{0x6542, 0xA951, 0x8654, 0x9510}, // L-tetrimino
	{0x6540, 0x9521, 0xA654, 0x9851}, // J-tetrimino
	{0x5421, 0xA651, 0x9865, 0x9540}, // S-tetrimino
	{0x6510, 0x9652, 0xA954, 0x8541}, // Z-tetrimino
}

var tSpinCheckShape = 0xA820

const kickWallTableTestCaseNum = 5

var kickWallTableForJLSTZ = [][kickWallTableTestCaseNum][2]int{
	{{0, 0}, {-1, 0}, {-1, +1}, {0, -2}, {-1, -2}}, // 0->R
	{{0, 0}, {+1, 0}, {+1, +1}, {0, -2}, {+1, -2}}, // 0->L

	{{0, 0}, {+1, 0}, {+1, -1}, {0, +2}, {+1, +2}}, // R->2
	{{0, 0}, {+1, 0}, {+1, -1}, {0, +2}, {+1, +2}}, // R->0

	{{0, 0}, {+1, 0}, {+1, +1}, {0, -2}, {+1, -2}}, // 2->L
	{{0, 0}, {-1, 0}, {-1, +1}, {0, -2}, {-1, -2}}, // 2->R

	{{0, 0}, {-1, 0}, {-1, -1}, {0, +2}, {-1, +2}}, // L->0
	{{0, 0}, {-1, 0}, {-1, -1}, {0, +2}, {-1, +2}}, // L->2
}

var kickWallTableForI = [][kickWallTableTestCaseNum][2]int{
	{{0, 0}, {-2, 0}, {+1, 0}, {-2, -1}, {+1, +2}}, // 0->R
	{{0, 0}, {-1, 0}, {+2, 0}, {-1, +2}, {+2, -1}}, // 0->L

	{{0, 0}, {-1, 0}, {+2, 0}, {-1, +2}, {+2, -1}}, // R->2
	{{0, 0}, {+2, 0}, {-1, 0}, {+2, +1}, {-1, -2}}, // R->0

	{{0, 0}, {+2, 0}, {-1, 0}, {+2, +1}, {-1, -2}}, // 2->L
	{{0, 0}, {+1, 0}, {-2, 0}, {+1, -2}, {-2, +1}}, // 2->R

	{{0, 0}, {+1, 0}, {-2, 0}, {+1, -2}, {-2, +1}}, // L->0
	{{0, 0}, {-2, 0}, {+1, 0}, {-2, -1}, {+1, +2}}, // L->2

}

// GameManager implement GameManager interface
type GameManager struct {
	// game optional variables(can be modified before game start)
	difficulty, lockDownDelay                  int
	height, bufferHeight, width, stashQueueCap int
	dropSpeedRatio                             float64
	allowSRS, allowGhost, allowHardDropOp      bool
	allowLockDownPeek, allowPlayAboveSkyline   bool
	allowForcedAboveSkyline                    bool
	allowTopOut, allowLockOut, allowBlockOut   bool

	// game internal variables(can not be modified by user)
	playfield                                            []int
	bag, nextBag                                         []int
	stashQueue                                           []int
	tetriminoX, tetriminoY, ghostX, ghostY               int
	tetriminoSpawnX, tetriminoSpawnY                     int
	tetriminoIdx, tetriminoDrct, nextTetriminoIdx        int
	score, level, bagIdx, comboCounter, dropLine, lastOp int
	fallSpeed                                            float64
	hardDropFlag, softDropFlag                           bool
	moveFlag, landFlag, lockDownTimerResetFlag           bool
	patternMatchFlag, tSpinFlag, backToBackFlag          bool

	// io utils
	inputCh  chan int
	renderer func(playfield, next []int, height, width, score int)
}

// NewGameManager return *GameManager
func NewGameManager(inputCh chan int, renderer func(playfield, next []int, height, width, score int)) *GameManager {
	return &GameManager{
		difficulty:              defaultDifficulty,
		lockDownDelay:           defaultLockDownDelay,
		height:                  defaultHeight,
		bufferHeight:            defaultBufferHeight,
		width:                   defaultWidth,
		dropSpeedRatio:          defaultDropSpeedRatio,
		allowSRS:                defaultAllowSRS,
		allowGhost:              defaultAllowGhost,
		allowHardDropOp:         defaultAllowHardDropOp,
		allowLockDownPeek:       defaultAllowLockDownPeek,
		allowPlayAboveSkyline:   defaultAllowPlayAboveSkyline,
		allowForcedAboveSkyline: defaultAllowForcedAboveSkyline,
		allowTopOut:             defaultAllowTopOut,
		allowLockOut:            defaultAllowLockOut,
		allowBlockOut:           defaultAllowBlockOut,
		inputCh:                 inputCh,
		renderer:                renderer,
	}
}

// ================ Utils ====================

func (gm *GameManager) reload() {
	rand.Seed(time.Now().UnixNano())
	gm.playfield = make([]int, gm.width*(gm.height+gm.bufferHeight))
	gm.tetriminoSpawnX = (gm.width - tetriNum) / 2
	gm.tetriminoSpawnY = gm.height
	gm.bag = make([]int, len(tetriminoShapes))
	gm.nextBag = make([]int, len(tetriminoShapes))
	// use bag system twice to ensure current bag is shuffled
	gm.bagIdx = len(gm.bag)
	gm.useBagSystem()
	gm.bagIdx = len(gm.bag)
	gm.useBagSystem()
	gm.stashQueue = make([]int, gm.stashQueueCap)
	gm.score = defaultScore
	gm.level = gm.difficulty
	gm.comboCounter = -1
	gm.dropLine = 0
	gm.calcFallSpeed()
}

// get from bag system (A 1.2.1)
func (gm *GameManager) useBagSystem() {
	gm.bagIdx++
	if gm.bagIdx >= len(gm.bag) {
		copy(gm.bag, gm.nextBag)
		for i := 0; i < len(gm.nextBag); i++ {
			gm.nextBag[i] = i
		}
		rand.Shuffle(
			len(gm.nextBag),
			func(i, j int) {
				gm.nextBag[i], gm.nextBag[j] = gm.nextBag[j], gm.nextBag[i]
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

func (gm *GameManager) lockDown() {
	for i := tetriminoShapes[gm.tetriminoIdx][gm.tetriminoDrct]; i != 0; i >>= 4 {
		x, y := gm.calcMinoPosOnBoard(i)
		gm.playfield[x+y*gm.width] = gm.tetriminoIdx + 1
	}
}

func (gm *GameManager) calcMinoPosOnBoard(posOnShape int) (x, y int) {
	x = gm.tetriminoX + posOnShape%tetriNum
	y = gm.tetriminoY - posOnShape/tetriNum%tetriNum
	return x, y
}

func (gm *GameManager) calcGhostMinoPosOnBoard(posOnShape int) (x, y int) {
	x = gm.ghostX + posOnShape%tetriNum
	y = gm.ghostY - posOnShape/tetriNum%tetriNum
	return x, y
}

func (gm *GameManager) calcGhostPos() {
	gm.ghostX = gm.tetriminoX
	gm.ghostY = gm.tetriminoY
	gm.checkLanding()
	if gm.landFlag {
		return
	}
	for true {
		gm.ghostY--
		for i := tetriminoShapes[gm.tetriminoIdx][gm.tetriminoDrct]; i != 0; i >>= 4 {
			x, y := gm.calcGhostMinoPosOnBoard(i)
			if y < 0 || gm.playfield[x+y*gm.width] != 0 {
				gm.ghostY++
				return
			}
		}
	}
}

// calculate the fall speed in current level (unit: Millisecond Per Line)
func (gm *GameManager) calcFallSpeed() {
	// TODO: can we modify these ratio?
	gm.fallSpeed = math.Pow(0.8-float64(gm.level-1)*0.007, float64(gm.level-1)) * 1000
}

// calculate the soft drop speed in current level (unit: Millisecond Per Line)
func (gm *GameManager) calcDropSpeed() {
	gm.calcFallSpeed()
	gm.fallSpeed = gm.fallSpeed / gm.dropSpeedRatio
}
func (gm *GameManager) checkBorderX(x int) bool {
	return x >= 0 && x < gm.width
}

func (gm *GameManager) checkBorderY(y int) bool {
	return y >= 0 && y < gm.height+gm.bufferHeight
}

func (gm *GameManager) checkNoCollision() bool {
	for i := tetriminoShapes[gm.tetriminoIdx][gm.tetriminoDrct]; i != 0; i >>= 4 {
		x, y := gm.calcMinoPosOnBoard(i)
		if gm.checkBorderX(x) && gm.checkBorderY(y) && gm.playfield[x+y*gm.width] != 0 {
			return false
		}
	}
	return true
}

func (gm *GameManager) checkLockOut() bool {
	return false
}

func (gm *GameManager) checkLanding() {
	gm.tetriminoY--
	for i := tetriminoShapes[gm.tetriminoIdx][gm.tetriminoDrct]; i != 0; i >>= 4 {
		x, y := gm.calcMinoPosOnBoard(i)
		if y < 0 || gm.playfield[x+y*gm.width] != 0 {
			gm.tetriminoY++
			gm.landFlag = true
			return
		}
	}
	gm.tetriminoY++
	gm.landFlag = false
}

func (gm *GameManager) checkTSpin() {

}

// =============== Basic Operation =================
func (gm *GameManager) move(opCode int) {
	if opCode == moveLeft {
		gm.tetriminoX--
	} else {
		gm.tetriminoX++
	}
	for i := tetriminoShapes[gm.tetriminoIdx][gm.tetriminoDrct]; i != 0; i >>= 4 {
		x, y := gm.calcMinoPosOnBoard(i)
		if !gm.checkBorderX(x) || !gm.checkBorderY(y) || gm.playfield[x+y*gm.width] != 0 {
			if opCode == moveLeft {
				gm.tetriminoX++
			} else {
				gm.tetriminoX--
			}
			return
		}
	}
	gm.calcGhostPos()
	gm.lastOp = opCode
	gm.moveFlag = true
}

func (gm *GameManager) rotate(opCode int) {
	srcDrct := gm.tetriminoDrct
	// Per rotate
	if opCode == rotateClockwise {
		gm.tetriminoDrct++
		if gm.tetriminoDrct >= len(tetriminoShapes[gm.tetriminoIdx]) {
			gm.tetriminoDrct = 0
		}
	} else {
		gm.tetriminoDrct--
		if gm.tetriminoDrct < 0 {
			gm.tetriminoDrct = len(tetriminoShapes[gm.tetriminoIdx]) - 1
		}
	}
	// check if blocked
	var testCaseNum int
	if gm.allowSRS {
		testCaseNum = kickWallTableTestCaseNum
	} else {
		testCaseNum = 1
	}
	var testCases [][kickWallTableTestCaseNum][2]int
	if gm.tetriminoIdx == tetriminoShapeI {
		testCases = kickWallTableForI
	} else {
		testCases = kickWallTableForJLSTZ
	}
	var directionOffset int
	if opCode == rotateClockwise {
		directionOffset = 0
	} else {
		directionOffset = 1
	}
	for i := 0; i < testCaseNum; i++ {
		testCase := testCases[srcDrct*2+directionOffset][i]
		offsetX, offsetY := testCase[0], testCase[1]
		blockFlag := false
		for i := tetriminoShapes[gm.tetriminoIdx][gm.tetriminoDrct]; i != 0; i >>= 4 {
			x, y := gm.calcMinoPosOnBoard(i)
			x += offsetX
			y += offsetY
			if !gm.checkBorderX(x) || !gm.checkBorderY(y) || gm.playfield[x+y*gm.width] != 0 {
				blockFlag = true
				break
			}
		}
		if blockFlag {
			continue
		}
		// pass test case
		gm.tetriminoX += offsetX
		gm.tetriminoY += offsetY
		gm.calcGhostPos()
		gm.lastOp = opCode
		gm.moveFlag = true
		return
	}
	// not pass any case, rollback
	if opCode != rotateClockwise {
		gm.tetriminoDrct++
		if gm.tetriminoDrct >= len(tetriminoShapes[gm.tetriminoIdx]) {
			gm.tetriminoDrct = 0
		}
	} else {
		gm.tetriminoDrct--
		if gm.tetriminoDrct < 0 {
			gm.tetriminoDrct = len(tetriminoShapes[gm.tetriminoIdx]) - 1
		}
	}
}

func (gm *GameManager) softDrop() {
	if gm.softDropFlag {
		gm.calcFallSpeed()
	} else {
		gm.calcDropSpeed()
		gm.lastOp = softDrop
	}
	gm.softDropFlag = !gm.softDropFlag
}

func (gm *GameManager) hardDrop() {
	gm.lastOp = hardDrop
	gm.tetriminoX, gm.tetriminoY = gm.ghostX, gm.ghostY
	gm.hardDropFlag = true
}

// ============ Running Flowchart ================

// Generration Phase (A 1.2.1)
func (gm *GameManager) generationPhase() bool {
	// Random Generation of Tetriminos
	gm.useBagSystem()
	// Starting Location and Orirntation
	gm.tetriminoX = gm.tetriminoSpawnX
	gm.tetriminoY = gm.tetriminoSpawnY
	gm.tetriminoDrct = 0
	gm.calcGhostPos()
	gm.moveFlag = true
	gm.renderOutput()
	return gm.checkNoCollision() || gm.allowBlockOut
}

func (gm *GameManager) fallingPhase() {
	gm.checkLanding()
	for !gm.landFlag {
		gm.hardDropFlag = false
		startTime, endTime := time.Now(), time.Now()
		for endTime.Sub(startTime) < time.Duration(gm.fallSpeed)*time.Millisecond {
			gm.processInput()
			gm.calcGhostPos()
			if gm.hardDropFlag && !gm.allowHardDropOp {
				return
			}
			endTime = time.Now()
		}
		gm.tetriminoY--
		gm.renderOutput()
		gm.checkLanding()
	}
}

// Lock Phase (A 1.2.1)
func (gm *GameManager) lockPhase() bool {
	startTime, endTime := time.Now(), time.Now()
	for !gm.hardDropFlag || gm.allowHardDropOp {
		if endTime.Sub(startTime) >= time.Duration(gm.lockDownDelay)*time.Millisecond {
			break
		}
		gm.processInput()
		gm.calcGhostPos()
		if gm.moveFlag && gm.landFlag && gm.lockDownTimerResetFlag {
			startTime = time.Now()
		}
		endTime = time.Now()
	}
	if gm.moveFlag && !gm.landFlag {
		return true
	}
	if !gm.checkNoCollision() || (gm.checkLockOut() && !gm.allowLockOut) {
		return false
	}
	gm.lockDown() // Lock down this tetrimino
	return true
}

// Pattern Phase (A 1.2.1)
func (gm *GameManager) patternPhase() {
	gm.patternMatchFlag = false
	count := gm.width
	for i := 0; i < len(gm.playfield); i++ {
		// TODO: judge T-spin
		if i%gm.width == 0 { // row head
			if count == 0 {
				gm.patternMatchFlag = true
				gm.playfield[i-gm.width] = rowFull
			} else if count == gm.width && i > 0 {
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

// this phase is for more variants (not used for now)
func (gm *GameManager) iteratePhase() {}

// this phase consumes no apparent game time (gm.renderOutput())
func (gm *GameManager) animatePhase() {}

func (gm *GameManager) elimatePhase() {
	clearLineCount := 0
	for rowIdx := 0; rowIdx < gm.height+gm.bufferHeight; rowIdx++ {
		rowHeaderPos := rowIdx * gm.width
		if gm.playfield[rowHeaderPos] == rowFull {
			clearLineCount++
		} else if gm.playfield[rowHeaderPos] == rowEmpty {
			gm.playfield[rowHeaderPos] = 0
			break // no need to move empty line
		}
		if clearLineCount > 0 {
			for colIdx := rowHeaderPos; colIdx < rowHeaderPos+gm.width; colIdx++ {
				gm.playfield[colIdx] = gm.playfield[colIdx+clearLineCount*gm.width]
				if gm.playfield[colIdx] < 0 {
					gm.playfield[colIdx] = 0
				}
			}
		}
	}
	// TODO: GameManager Statistics

	gm.renderOutput()
}

func (gm *GameManager) completionPhase() {
	// update information
	// level up condition
}

// Tetris engine flowchart
func (gm *GameManager) loopFlow() {
	for gm.generationPhase() {
		for gm.moveFlag && !gm.landFlag {
			gm.fallingPhase()
			if !gm.lockPhase() {
				return
			}
		}
		gm.patternPhase()
		gm.iteratePhase()
		gm.animatePhase()
		gm.elimatePhase()
		gm.completionPhase()
	}
}

// =================== IO ========================

func (gm *GameManager) processInput() {
	select {
	case input := <-gm.inputCh:
		switch input {
		case moveLeft, moveRight:
			gm.move(input)
		case rotateClockwise, rotateCounterClockwise:
			gm.rotate(input)
		case softDrop:
			gm.softDrop()
		case hardDrop:
			gm.hardDrop()
		}
	default:

	}
	gm.renderOutput()

}

func (gm *GameManager) renderOutput() {
	playfield := make([]int, gm.width*gm.height)
	for i := 0; i < len(playfield); i++ {
		playfield[i] = gm.playfield[i]
	}
	for i := tetriminoShapes[gm.tetriminoIdx][gm.tetriminoDrct]; i != 0; i >>= 4 {
		x, y := gm.calcGhostMinoPosOnBoard(i)
		if x+y*gm.width < len(playfield) {
			playfield[x+y*gm.width] = (gm.tetriminoIdx + 1) * -1
		}
		x, y = gm.calcMinoPosOnBoard(i)
		if x+y*gm.width < len(playfield) {
			playfield[x+y*gm.width] = gm.tetriminoIdx + 1
		}
	}
	// nextTetrimin := make([][]int, 1)

	gm.renderer(playfield, gm.height, gm.width, gm.score)
}

// ============= Export Function ===============

// GetSetups of game manager
func (gm *GameManager) GetSetups() {}

// Setup game optional settings
func (gm *GameManager) Setup() {}

// RestoreDefaultSetup for game manager
func (gm *GameManager) RestoreDefaultSetup() {}

// NewGame start for caller
// GameManager Over condition occurs in Generation Phase and Lock Phase
func (gm *GameManager) NewGame() {
	gm.reload()
	// Tetris engine flowchart
	gm.loopFlow()
	// GameManager Over Events
}
