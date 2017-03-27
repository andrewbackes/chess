package position

import (
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
	"testing"
)

func TestRootMoves(t *testing.T) {
	b := New()
	moves := b.Moves()
	if len(moves) != 20 {
		t.Log(moves, len(moves))
		t.Log(b)
		t.Fail()
	}
}

/*
func TestCheck(t *testing.T) {
	whiteChecked := []string{"rnb1kbnr/pppp1ppp/8/4p3/4P1q1/2N5/PPPPKPPP/R1BQ1BNR w kq - 4 4"}
	for _, check := range whiteChecked {
		b, _ := FromFEN(check)
		if !b.Check(piece.White) {
			t.Fail()
		}
	}
}
*/

func TestGenPromotion(t *testing.T) {
	b := New()
	b.Clear()
	b.QuickPut(piece.New(piece.White, piece.Pawn), square.E7)
	moves := b.LegalMoves()
	expected := []string{"e7e8q", "e7e8r", "e7e8b", "e7e8n"}
	if len(moves) != len(expected) {
		t.Fail()
	}
	for _, exp := range expected {
		if _, ok := moves[move.Parse(exp)]; !ok {
			t.Fail()
		}
	}
}

func TestGenCastles(t *testing.T) {
	b := New()
	b.Clear()
	b.QuickPut(piece.New(piece.White, piece.King), square.E1)
	b.QuickPut(piece.New(piece.White, piece.Rook), square.H1)
	b.QuickPut(piece.New(piece.White, piece.Rook), square.A1)
	moves := b.LegalMoves()
	expected := []string{"e1g1", "e1c1"}
	for _, exp := range expected {
		if _, ok := moves[move.Parse(exp)]; !ok {
			t.Fail()
		}
	}
}

func TestGenPromotionCap(t *testing.T) {
	// TODO
}

func TestGenKingCaps(t *testing.T) {
	// TODO
}

func TestGenDiagCaps(t *testing.T) {
	// TODO
}

func TestEnPassant(t *testing.T) {
	// TODO
}
