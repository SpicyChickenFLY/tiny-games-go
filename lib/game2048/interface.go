package game2048

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	KeyEsc = 0

	KeyUp    = 1
	KeyDown  = 2
	KeyLeft  = 4
	KeyRight = 8
)

func (g *Game) process(inputCh chan int, logCh chan string) {
	for input := range inputCh {
		switch input {
		case KeyUp:
			g.operate(g.up)
			g.log(logCh, g.up)
		case KeyDown:
			g.operate(g.down)
			g.log(logCh, g.down)
		case KeyLeft:
			g.operate(g.left)
			g.log(logCh, g.left)
		case KeyRight:
			g.operate(g.right)
			g.log(logCh, g.right)
		case KeyEsc:
			return
		}
	}
}

func (g *Game) render(
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

func (g *Game) log(logCh chan string, direction int) {

}

func (g *Game) run(
	name string,
	width, height, difficult int,
	inputChannel chan int,
	logChannel chan string,
	renderFunc func(board []int, height, width, score, fps int),
) (score int) {

	rand.Seed(time.Now().UnixNano())
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	g.init(4, 4, 2)

	stopRenderCh := make(chan struct{})
	go g.render(renderFunc, stopRenderCh)
	defer close(stopRenderCh)

	go g.process(inputChannel, logChannel)

	return g.score
}
