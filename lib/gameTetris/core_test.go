package gameTetris

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var gm *GameManager = nil

func TestNewGameManager(t *testing.T) {
	gm = NewGameManager()
	assert.NotNil(t, gm)
}

func TestGame_init(t *testing.T) {
	// gm.init()
}

func TestGame_calcPosOnBoard(t *testing.T) {

}
