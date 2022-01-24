package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/SpicyChickenFLY/tiny-games-go/lib/gameBullsAndCows"
)

func main() {
	// if err := termbox.Init(); err != nil {
	// 	panic(err)
	// }
	// termbox.HideCursor()

	// gameBullsAndCows.NewGame()

	g := gameBullsAndCows.NewGame(4)
	inputReader := bufio.NewReader(os.Stdin) //创建一个读取器，并将其与标准输入绑定。
	for {
		var guessStr string
		guessStr, _ = inputReader.ReadString('\n') //读取器对象提供一个方法 ReadString(delim byte) ，该方法从输入中读取内容，直到碰到 delim 指定的字符，然后将读取到的内容连同 delim 字符一起放到缓冲区。
		guessStr = strings.ReplaceAll(guessStr, "\n", "")
		if len(guessStr) != 5 {
			fmt.Println("wrong input")
			continue
		}
		guess := make([]int, 4)
		for i := 0; i < 4; i++ {
			guess[i], _ = strconv.Atoi(string(guessStr[i]))
		}
		result, bulls, cows := g.Guess(guess)
		fmt.Printf("%v correct:%d, contain%d\n", guess, bulls, cows)
		if result {
			return
		}
	}
}
