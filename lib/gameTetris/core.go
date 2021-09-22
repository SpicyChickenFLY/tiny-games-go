package gameTetris

const (
	emptyTile = 0
)

// operations
const (
	moveLeft = iota
	moveRight
	rotate
	accelerate
	decelerate
)

var shapeLib = []int{
	// Two-status-shape
	0xC840, 0x3210, // I-Shape
	0x5410, 0x5410, // O-shape
	0x6510, 0x8541, // filped-Z-shape
	0x5421, 0x9540, // Z-shape
	// Four-status-shape
	0x6541, 0x9651, 0x9654, 0x9541, // T-shape
	0x6542, 0xA951, 0x8654, 0x9510, // 7-shape
	0x6540, 0x9521, 0xA654, 0x9851, // filped-7-shape
}

var shapeColor = []int{}

// Game implement Game interface
type Game struct {
	board                               []int
	height, width, boardSize            int
	shapeX, shapeY, shapeSize, shapeIdx int
	score, difficult, delay             int
	alive                               bool
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
	g.fps = 60
	g.difficult = difficult
	g.delay = difficult
	g.width, g.height, g.boardSize = w, h, w*h
	g.board = make([]int, g.width*g.height)
}

func (g *Game) removeShapeFromBoard() {
	for pos := shapeLib[g.shapeIdx]; pos > 0; pos >>= 4 {
		_, _, pob := g.calPosOnBoard(pos)
		g.board[pob] = emptyTile
	}
}

func (g *Game) addShapeToBoard() {
	color := shapeColor[g.shapeIdx]
	for pos := shapeLib[g.shapeIdx]; pos > 0; pos >>= 4 {
		_, _, pob := g.calPosOnBoard(pos)
		g.board[pob] = color
	}
}

func (g *Game) calPosOnBoard(posOnShape int) (x, y, pos int) {
	x = g.shapeX + posOnShape%g.shapeSize
	y = g.shapeY + posOnShape/g.shapeSize%g.shapeSize
	pos = y*g.width + x
	return x, y, pos
}

func (g *Game) rotate() {

}

func (g *Game) move(isDirectionLeft bool) {
	g.removeShapeFromBoard()

	if isDirectionLeft {
		pos
	} else {

	}
	for pos := shapeLib[g.shapeIdx]; pos > 0; pos >>= 4 {
		x, _, pob := g.calPosOnBoard(pos)
		if x < 0 || x >= g.width || g.board[pob] != 0 {
			return
		}
	}
	g.addShapeToBoard()
}

func (g *Game) drop() {

}

func (g *Game) accelerate() {

}

func (g *Game) operate(operation int) {
	switch operation {
	case moveLeft:
		g.move(true)
	case moveRight:
		g.move(false)
	case rotate:
		g.rotate()
	case accelerate:
		g.accelerate()
	case decelerate:
		g.decelerate()
	}

	BadExpr

}
