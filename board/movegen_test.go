package board

import (
	"github.com/andrewbackes/chess/piece"
	"testing"
)

func TestRootMoves(t *testing.T) {
	b := New()
	moves := b.Moves(piece.White, nil, [2][2]bool{})
	if len(moves) != 20 {
		t.Log(moves, len(moves))
		t.Log(b)
		t.Fail()
	}
}
