package game

import (
	"testing"
)

// TestNewGame just makes sure we can get a new game.
func TestRootMoves(t *testing.T) {
	g := NewGame()
	moves := g.LegalMoves()
	if len(moves) != 20 {
		t.Log(moves, len(moves))
		t.Log(g.board)
		t.Fail()
	}
}

func TestPerftSuite(t *testing.T) {
	f := "perftsuite.epd"
	d := 1
	if testing.Short() {
		d = 1
	}
	if err := PerftSuite(f, d); err != nil {
		t.Error(err)
	}
}
