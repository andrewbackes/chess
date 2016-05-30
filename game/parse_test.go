package chess

import (
	"github.com/andrewbackes/chess/board"
	"github.com/andrewbackes/chess/piece"
	"testing"
)

func TestParseLongCastle(t *testing.T) {

}

func TestParseNoPromo(t *testing.T) {
	g := NewGame()
	g.board.Clear()
	g.board.QuickPut(piece.New(piece.White, piece.Pawn), board.E7)
	g.board.QuickPut(piece.New(piece.White, piece.King), board.E1)
	g.board.QuickPut(piece.New(piece.Black, piece.King), board.A8)
	move, _ := g.ParseMove("e7e8")
	if move != board.Move("e7e8q") {
		t.Fail()
	}
}
