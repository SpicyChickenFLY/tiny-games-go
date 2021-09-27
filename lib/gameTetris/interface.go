package gameTetris

import (
	"fmt"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var colorMap = []termbox.Attribute{
	termbox.ColorBlack,
	termbox.ColorWhite,
	termbox.ColorCyan,
	termbox.ColorMagenta,
	termbox.ColorBlue,
	termbox.ColorYellow,
	termbox.ColorGreen,
	termbox.ColorRed,
}

//  =================== Utils ===================
func ListenToInput(inputCh chan int) {
	termbox.SetInputMode(termbox.InputEsc)
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowLeft:
				inputCh <- moveLeft
			case termbox.KeyArrowRight:
				inputCh <- moveRight
			case termbox.KeyArrowUp:
				inputCh <- rotateClockwise
			case termbox.KeyArrowDown:
				inputCh <- softDrop
			case termbox.KeySpace:
				inputCh <- hardDrop

			case termbox.KeyEsc:
				panic("bye")
			}
			switch ev.Ch {
			case 'x':
				inputCh <- rotateClockwise
			case 'z':
				inputCh <- rotateCounterClockwise
			}

		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func RenderToScreen(playfield, next []int, height, width, score int) {
	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		panic(err)
	}
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if playfield[i*width+j] > 0 {
				// tbprint(j*2, height-i, colorMap[playfield[i*width+j]], termbox.ColorDefault, fmt.Sprint(playfield[i*width+j]))
				tbprint(j*2, height-i, colorMap[playfield[i*width+j]], termbox.ColorBlack, "⬛")
			} else {
				// tbprint(j*2, height-i, termbox.ColorBlack, colorMap[playfield[i*width+j]*-1], fmt.Sprint(playfield[i*width+j]))
				tbprint(j*2, height-i, termbox.ColorBlack, colorMap[playfield[i*width+j]*-1], "⬛")
			}

		}

	}
	tbprint(22, 5, termbox.ColorRed, termbox.ColorDefault, fmt.Sprintf("score: %8d", score))
	if err := termbox.Flush(); err != nil {
		panic(err)
	}
}

// This function is often useful:
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}
