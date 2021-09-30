package gameTetris

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	
}

var gm *GameManager

func TestNewGameManager(t *testing.T) {
	gm = NewGameManager()
	assert.NotNil(t, gm, "Failed to instantiate class(GameManager)")
}

func TestGameManager_reload(t *testing.T) {
	gm.reload()
}

func TestGame_checkBorderX(t *testing.T) {
	testCases := make([]interface{}, 0)
	10, -1, 0, 1, 4, 5, 10, 20, 30, 100
	for _, testCase := range testCases {
		gm.checkBorderX(testCase)
	}
}

func TestGame_checkBorderY(t *testing.T) {}

func TestGame_calcPosOnBoard(t *testing.T) {}
