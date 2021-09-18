package game2048

import (
	"fmt"
	"math/rand"
	"time"

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
				inputCh <- keyLeft
			case termbox.KeyArrowDown:
				inputCh <- keyDown
			case termbox.KeyArrowRight:
				inputCh <- keyRight
			case termbox.KeyArrowUp:
				inputCh <- keyUp
			case termbox.KeyEsc:
				inputCh <- keyEsc
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
		panic(err)
	}
	tbprint(0, 0, termbox.ColorDefault, termbox.ColorDefault, fmt.Sprintf("score:%d", score))
	tbprint(0, 1, termbox.ColorDefault, termbox.ColorDefault, "=====================================")

	tbprint(0, 3, termbox.ColorDefault, termbox.ColorDefault, "-------------------------")
	for i := 0; i < height; i++ {
		tbprint(0, i*2+4, termbox.ColorDefault, termbox.ColorDefault, "|")
		for j := 0; j < width; j++ {
			tbprint(6*j+1, i*2+4, colorMap[board[i*width+j]], termbox.ColorBlack, strMap[board[i*width+j]])
			tbprint(6*j+6, i*2+4, termbox.ColorDefault, termbox.ColorDefault, "|")
		}
		tbprint(0, i*2+5, termbox.ColorDefault, termbox.ColorDefault, "-------------------------")
	}
	if err := termbox.Flush(); err != nil {
		panic(err)
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
	go recordLog(logChannel)

	game := Game{}
	score := game.run(name, width, height, difficult, inputChannel, logChannel, render)
	fmt.Println("your final socre is: ", score)
}
