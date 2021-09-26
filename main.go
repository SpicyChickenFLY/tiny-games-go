package main

import (
	"github.com/SpicyChickenFLY/tiny-games-go/lib/gameTetris"
	"github.com/nsf/termbox-go"
)

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	termbox.HideCursor()

	inputCh := make(chan int, 5)
	go gameTetris.ListenToInput(inputCh)

	game := gameTetris.NewGameManager(inputCh, gameTetris.RenderToScreen)
	game.NewGame()
}
