package gameTetris

import (
	"fmt"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var colorMap = []termbox.Attribute{
	termbox.ColorWhite,
	termbox.ColorRed,
	termbox.ColorYellow,
	termbox.ColorGreen,
	termbox.ColorCyan,
	termbox.ColorBlue,
	termbox.ColorMagenta,
	termbox.ColorDarkGray,
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

func RenderToScreen(playfield []int, height, width int) {
	counter := 0
	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		panic(err)
	}
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			tbprint(width-j, height-i, termbox.ColorRed, termbox.ColorDefault, fmt.Sprint(playfield[i*width+j]))
		}

	}
	tbprint(30, 1, termbox.ColorRed, termbox.ColorDefault, fmt.Sprint(counter))
	if err := termbox.Flush(); err != nil {
		panic(err)
	}
	counter++
}

// This function is often useful:
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}
