package position

import (
	"github.com/andrewbackes/chess/piece"
	"testing"
)

func TestParseLongCastle(t *testing.T) {

}

func TestParseNoPromo(t *testing.T) {
	p := New()
	p.Clear()
	p.QuickPut(piece.New(piece.White, piece.Pawn), E7)
	p.QuickPut(piece.New(piece.White, piece.King), E1)
	p.QuickPut(piece.New(piece.Black, piece.King), A8)
	move, _ := p.ParseMove("e7e8")
	if move != Move("e7e8q") {
		t.Fail()
	}
}
