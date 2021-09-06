package game2048

const (
	keyEsc = 0

	keyUp    = 1
	keyDown  = 2
	keyLeft  = 4
	keyRight = 8
)

func (g *Game) process(inputCh chan int, logCh chan string) {
	for input := range inputCh {
		switch input {
		case keyUp:
			g.operate(g.up)
			g.log(logCh, g.up)
		case keyDown:
			g.operate(g.down)
			g.log(logCh, g.down)
		case keyLeft:
			g.operate(g.left)
			g.log(logCh, g.left)
		case keyRight:
			g.operate(g.right)
			g.log(logCh, g.right)
		case keyEsc:
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
	g.init(4, 4, 2)

	stopRenderCh := make(chan struct{})
	go g.render(renderFunc, stopRenderCh)
	defer close(stopRenderCh)

	g.process(inputChannel, logChannel)

	return g.score
}
