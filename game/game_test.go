package game

import (
	"testing"
	"time"
)

// TestNewGame just makes sure we can get a new game.
func TestNewGame(t *testing.T) {
	NewGame()
	//t.Log(g.board.String())
}

func TestNewTimedGame(t *testing.T) {
	standard := TimeControl{
		Time:  40 * time.Minute,
		Moves: 40,
	}
	control := [2]TimeControl{standard, standard}
	NewTimedGame(control)
}
