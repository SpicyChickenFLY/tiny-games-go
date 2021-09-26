package gameTetris

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var gm *GameManager

func TestNewGameManager(t *testing.T) {
	gm = NewGameManager()
	assert.NotNil(t, gm, "Failed to instantiate class(GameManager)")
}

func TestGameManager_reload(t *testing.T) {
	gm.reload()
}

func TestGame_checkBorderX(t *testing.T) {
	// testCasesForDefault =
}

func TestGame_checkBorderY(t *testing.T) {}

func TestGame_calcPosOnBoard(t *testing.T) {}
