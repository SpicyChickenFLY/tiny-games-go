package gameTetris

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

const (
	keyEsc = 0

	keyUp    = 1
	keyDown  = 2
	keyLeft  = 4
	keyRight = 8
)

func renderFrame(
	g Game,
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

func run(
	g Game,
	width, height, difficult int,
	inputChannel chan int,
	renderFunc func(playfield []int, height, width, score, fps int),
) (score int) {
	g.init(4, 4, 2)

	stopRenderCh := make(chan struct{})
	go render(g, renderFunc, stopRenderCh)
	defer close(stopRenderCh)

	process(g, inputChannel)

	return g.score
}

// Run is the entrance of game 2048 in cmd
func Run(name string, width int, height int, difficult int) {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	inputChannel := make(chan int, 5)
	go listenToInput(inputChannel)

	game := Game{}

	score := run(game, width, height, difficult, inputChannel, renderToScreen)
	fmt.Println("your final socre is: ", score)
}
