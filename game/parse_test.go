package game

import (
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position"
	"testing"
)

func TestParseLongCastle(t *testing.T) {

}

func TestParseNoPromo(t *testing.T) {
	g := New()
	g.Position.Clear()
	g.Position.QuickPut(piece.New(piece.White, piece.Pawn), position.E7)
	g.Position.QuickPut(piece.New(piece.White, piece.King), position.E1)
	g.Position.QuickPut(piece.New(piece.Black, piece.King), position.A8)
	move, _ := g.ParseMove("e7e8")
	if move != position.Move("e7e8q") {
		t.Fail()
	}
}
