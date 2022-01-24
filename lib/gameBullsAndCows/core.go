package gameBullsAndCows

import (
	"math/rand"
)

type history struct {
	guess []int
	bulls int
	cows  int
}

type Game struct {
	numLen   int
	secret   []int
	historys []history
}

func NewGame(numLen int) Game {
	g := Game{}
	g.init(numLen)
	return g
}

func (g *Game) init(numLen int) {
	g.numLen = numLen
	g.secret = make([]int, 0)
	for i := 0; i < numLen; i++ {
		g.secret = append(g.secret, rand.Intn(10))
	}
}

func (g *Game) Guess(guess []int) (bool, int, int) {
	bulls, cows := g.getHint(guess)
	g.historys = append(g.historys, history{guess, bulls, cows})
	return bulls == g.numLen, bulls, cows
}

func (g *Game) guess(guess []int) bool {
	bulls, cows := g.getHint(guess)
	g.historys = append(g.historys, history{guess, bulls, cows})
	return bulls == g.numLen
}

func (g *Game) getHint(guess []int) (bulls, cows int) {
	l := len(g.secret)
	cnt := [10]int{}
	for i := 0; i < l; i++ {
		if g.secret[i] == guess[i] {
			bulls++
		} else {
			cnt[g.secret[i]]++
			cnt[guess[i]]--
		}
	}
	sum := 0
	for _, v := range cnt {
		if v <= 0 {
			sum -= v
		} else {
			sum += v
		}
	}
	cows = l - bulls - sum/2
	return
}
