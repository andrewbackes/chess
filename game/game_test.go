package game

import (
	"testing"
	"time"
)

// TestNewGame just makes sure we can get a new game.
func TestNewGame(t *testing.T) {
	NewGame()
}

func TestNewTimedGame(t *testing.T) {
	standard := TimeControl{
		Time:  40 * time.Minute,
		Moves: 40,
	}
	control := [2]TimeControl{standard, standard}
	NewTimedGame(control)
}

func TestNonexistentMove(t *testing.T) {
	g := NewGame()
	mv := Move("e4e5")
	status := g.MakeMove(mv)
	if status != WhiteIllegalMove {
		t.Error("Got: ", status, " Wanted: ", WhiteIllegalMove)
	}
}
