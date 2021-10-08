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
	posTestCases := []int{-10, -5, -2, -1}
	for _, testCase := range posTestCases {
		assert.True(t, gm.checkBorderX(testCase), "error in positive testcase when check border X")
	}
	negTestCases := []int{0, 5, 10, 20, 30}
	for _, testCase := range negTestCases {
		assert.False(t, gm.checkBorderX(testCase), " error in negative testcase when check border X")
	}
}

func TestGame_checkBorderY(t *testing.T) {
	posTestCases := []int{-10, -5, -2, -1}
	for _, testCase := range posTestCases {
		assert.True(t, gm.checkBorderX(testCase), "error in positive testcase when check border Y")
	}
	negTestCases := []int{0, 5, 10, 20, 30}
	for _, testCase := range negTestCases {
		assert.False(t, gm.checkBorderX(testCase), " error in negative testcase when check border Y")
	}
}

func TestGame_calcPosOnBoard(t *testing.T) {
}
