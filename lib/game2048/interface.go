package game2048

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	keyEsc = 0

	keyUp    = 1
	keyDown  = 2
	keyLeft  = 4
	keyRight = 8
)

func process(g Game, inputCh chan int) {
	for input := range inputCh {
		switch input {
		case keyUp:
			g.operate(g.up)
		case keyDown:
			g.operate(g.down)
		case keyLeft:
			g.operate(g.left)
		case keyRight:
			g.operate(g.right)
		case keyEsc:
			return
		}
	}
}

func render(
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

func log(logCh chan string, direction int) {

}

func run(
	g Game,
	name string,
	width, height, difficult int,
	inputChannel chan int,
	logChannel chan string,
	renderFunc func(board []int, height, width, score, fps int),
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
	rand.Seed(time.Now().UnixNano())
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	// log.SetOutput(os.Stdout)
	// log.SetLevel(log.InfoLevel)

	inputChannel := make(chan int, 5)
	logChannel := make(chan string, 5)

	go listenToInput(inputChannel)

	game := Game{}
	score := run(game, name, width, height, difficult, inputChannel, logChannel, renderToScreen)
	fmt.Println("your final socre is: ", score)
}
