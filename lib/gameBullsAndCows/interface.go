package gameBullsAndCows

import (
	"fmt"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type inputModel struct {
	guess  []int
	cursor int
}

var im = inputModel{
	guess:  make([]int, 0),
	cursor: 0,
}

func ListenToInput(g Game, stopCh chan struct{}) {
	termbox.SetInputMode(termbox.InputEsc)
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventError {
			panic(ev.Err)
		}
		if ev.Type == termbox.EventKey {
			switch ev.Key {
			case termbox.KeyArrowLeft:
				im.cursor--
				if im.cursor < 0 {
					im.cursor = len(im.guess) - 1
				}
			case termbox.KeyArrowRight:
				im.cursor++
				if im.cursor >= len(im.guess) {
					im.cursor = 0
				}
			case termbox.KeyArrowUp:
				im.guess[im.cursor]++
				if im.guess[im.cursor] > 9 {
					im.guess[im.cursor] = 0
				}
			case termbox.KeyArrowDown:
				im.guess[im.cursor]--
				if im.guess[im.cursor] < 0 {
					im.guess[im.cursor] = 9
				}
			case termbox.KeyEnter:
				if g.guess(im.guess) {
					close(stopCh)
					fmt.Println("You win!")
					return
				}
			case termbox.KeyEsc:
				close(stopCh)
				return
			}
		}
	}
}

func render(
	g Game,
	renderFunc func(history []history),
	stopCh <-chan struct{}) {
	for {
		select {
		case <-stopCh:
			return
		default:
			renderFunc(g.historys)
		}
	}
}

func renderToConsole(historys []history) {
	tbprint(
		0, 0, termbox.ColorDefault, termbox.ColorRed,
		"Guess:")
	for i, guessNum := range im.guess {
		if i == im.cursor {
			tbprint(
				0, i*2+6, termbox.ColorDefault, termbox.ColorRed,
				fmt.Sprintf("%d", guessNum))
		} else {
			tbprint(
				0, i*2+6, termbox.ColorDefault, termbox.ColorDefault,
				fmt.Sprintf("%d", guessNum))
		}
	}
	for j, history := range historys {
		tbprint(
			j, 0, termbox.ColorWhite, termbox.ColorDefault,
			fmt.Sprintf("%v", history.guess))
		tbprint(
			j, 20, termbox.ColorGreen, termbox.ColorDefault,
			fmt.Sprintf("Right:%d", history.bulls))
		tbprint(
			j, 30, termbox.ColorYellow, termbox.ColorDefault,
			fmt.Sprintf("Contain:%d", history.cows))
	}

}

// This function is often useful:
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func StartGame() {
	g := Game{}
	g.init(4)
	im = inputModel{make([]int, 4), 0}
	stopCh := make(chan struct{}, 5)
	go render(g, renderToConsole, stopCh)
	ListenToInput(g, stopCh)
}
