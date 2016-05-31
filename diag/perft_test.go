package diag

import (
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position"
	"testing"
)

func TestSimplePerftOutput(t *testing.T) {
	p := position.New()
	p.Clear()

	p.Put(piece.New(piece.White, piece.Pawn), position.E2)
	c := Perft(p, 1)
	if c != 2 {
		t.Fail()
	}
}

func TestDivideOutput(t *testing.T) {
	p := position.New()
	p.Clear()
	p.Put(piece.New(piece.White, piece.King), position.A1)
	p.Put(piece.New(piece.Black, piece.King), position.A8)
	m := Divide(p, 1)
	if len(m) != 3 {
		t.Fail()
	}
	if m["a1b1"] != 1 || m["a1b2"] != 1 || m["a1a2"] != 1 {
		t.Fail()
	}
}
