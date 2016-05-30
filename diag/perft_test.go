package diag

import (
	"github.com/andrewbackes/chess"
	"github.com/andrewbackes/chess/board"
	"github.com/andrewbackes/chess/piece"
	"testing"
)

func TestSimplePerftOutput(t *testing.T) {
	g := chess.NewGame()
	g.Board().Clear()
	g.Board().Put(piece.New(piece.White, piece.Pawn), board.E2)
	c := Perft(g, 1)
	if c != 2 {
		t.Fail()
	}
}

func TestDivideOutput(t *testing.T) {
	g := chess.NewGame()
	g.Board().Clear()
	g.Board().Put(piece.New(piece.White, piece.King), board.A1)
	g.Board().Put(piece.New(piece.Black, piece.King), board.A8)
	m := Divide(g, 1)
	if len(m) != 3 {
		t.Fail()
	}
	if m["a1b1"] != 1 || m["a1b2"] != 1 || m["a1a2"] != 1 {
		t.Fail()
	}
}
