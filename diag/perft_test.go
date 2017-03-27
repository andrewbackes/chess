package diag

import (
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
	"testing"
)

func TestSimplePerftOutput(t *testing.T) {
	p := position.New()
	p.Clear()

	p.Put(piece.New(piece.White, piece.Pawn), square.E2)
	c := Perft(p, 1)
	if c != 2 {
		t.Fail()
	}
}

func TestDivideOutput(t *testing.T) {
	p := position.New()
	p.Clear()
	p.Put(piece.New(piece.White, piece.King), square.A1)
	p.Put(piece.New(piece.Black, piece.King), square.A8)
	m := Divide(p, 1)
	if len(m) != 3 {
		t.Fail()
	}
	if m[*move.Parse("a1b1")] != 1 || m[*move.Parse("a1b2")] != 1 || m[*move.Parse("a1a2")] != 1 {
		t.Fail()
	}
}
