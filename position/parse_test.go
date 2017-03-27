package position

import (
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
	"testing"
)

func TestParseLongCastle(t *testing.T) {

}

func TestParseNoPromo(t *testing.T) {
	p := New()
	p.Clear()
	p.QuickPut(piece.New(piece.White, piece.Pawn), square.E7)
	p.QuickPut(piece.New(piece.White, piece.King), square.E1)
	p.QuickPut(piece.New(piece.Black, piece.King), square.A8)
	m, _ := p.ParseMove("e7e8")
	if m != move.Parse("e7e8q") {
		t.Fail()
	}
}
