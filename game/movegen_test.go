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
