package game2048

import (
	"log"
	"os"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var colorMap = []termbox.Attribute{
	termbox.ColorRed,
	termbox.ColorRed,
	termbox.ColorYellow,
	termbox.ColorGreen,
	termbox.ColorCyan,
	termbox.ColorBlue,
	termbox.ColorMagenta,
	termbox.ColorDarkGray,
	termbox.ColorBlack,
	termbox.ColorBlack,
	termbox.ColorBlack,
	termbox.ColorBlack,
	termbox.ColorBlack,
	termbox.ColorBlack,
	termbox.ColorBlack,
	termbox.ColorBlack,
	termbox.ColorBlack,
}

var strMap = []string{
	"     ",
	"  2  ",
	"  4  ",
	"  8  ",
	"  16 ",
	"  32 ",
	"  64 ",
	" 128 ",
	" 256 ",
	" 512 ",
	" 1024",
	" 2048",
	" 4096",
	" 8192",
	"16384",
	"32768",
	"65536",
}

//  =================== Utils ===================
func listenToInput(inputCh chan int) {
	termbox.SetInputMode(termbox.InputEsc)
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowLeft:
				inputCh <- KeyLeft
			case termbox.KeyArrowDown:
				inputCh <- KeyDown
			case termbox.KeyArrowRight:
				inputCh <- KeyRight
			case termbox.KeyArrowUp:
				inputCh <- KeyUp
			case termbox.KeyEsc:
				inputCh <- KeyEsc
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func recordLog(logCh chan string) {
	for _ = range logCh {
		// log.Info(logStr)
	}
}

func render(board []int, height, width, score, fps int) {
	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	tbprint(0, 0, termbox.ColorWhite, termbox.ColorBlack, "")
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			tbprint(6*j+1, i+1, colorMap[board[i*width+j]], termbox.ColorBlack, strMap[board[i*width+j]])
		}
	}
	if err := termbox.Flush(); err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	time.Sleep(time.Duration(1000/fps) * time.Millisecond)
}

// This function is often useful:
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func Run(name string, width int, height int, difficult int) {
	// log.SetOutput(os.Stdout)
	// log.SetLevel(log.InfoLevel)

	inputChannel := make(chan int, 5)
	logChannel := make(chan string, 5)

	go listenToInput(inputChannel)
	go recordLog(logChannel)

	game := Game{}
	game.run(name, width, height, difficult, inputChannel, logChannel, render)
}
