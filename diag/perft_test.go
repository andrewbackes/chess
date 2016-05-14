package diag

import (
	"github.com/andrewbackes/chess/game"
	"testing"
)

func TestSimplePerftOutput(t *testing.T) {
	g := game.NewGame()
	g.Board().Clear()
	g.Board().Put(game.NewPiece(game.White, game.Pawn), game.E2)
	c := Perft(g, 1)
	if c != 2 {
		t.Fail()
	}
}

func TestDivideOutput(t *testing.T) {
	g := game.NewGame()
	g.Board().Clear()
	g.Board().Put(game.NewPiece(game.White, game.King), game.A1)
	g.Board().Put(game.NewPiece(game.Black, game.King), game.A8)
	m := Divide(g, 1)
	if len(m) != 3 {
		t.Fail()
	}
	if m["a1b1"] != 1 || m["a1b2"] != 1 || m["a1a2"] != 1 {
		t.Fail()
	}
}
