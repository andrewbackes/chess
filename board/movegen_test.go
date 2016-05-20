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

func TestCheck(t *testing.T) {
	whiteChecked := []string{"rnb1kbnr/pppp1ppp/8/4p3/4P1q1/2N5/PPPPKPPP/R1BQ1BNR w kq - 4 4"}
	for _, check := range whiteChecked {
		b, _ := FromFEN(check)
		if !b.Check(piece.White) {
			t.Fail()
		}
	}
}

func TestGenCastles(t *testing.T) {
	// TODO
}

func TestPromotion(t *testing.T) {
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
