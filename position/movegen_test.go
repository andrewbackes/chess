package position

import (
	"testing"

	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
)

func TestRootMoves(t *testing.T) {
	testCases := []struct {
		name    string
		c       piece.Color
		wantLen int
	}{
		{"White", piece.White, 20},
		{"Black", piece.Black, 20},
		{"NoColor", piece.NoColor, 0},
		{"Color(5)", piece.Color(5), 0},
	}
	p := New()
	for _, tc := range testCases {
		p.ActiveColor = tc.c
		t.Run(tc.name, func(t *testing.T) {
			if moves := p.Moves(); moves == nil || len(moves) != tc.wantLen {
				t.Errorf("len(p.Moves()) for new board with %s as active color is %d, want %d", tc.c, len(moves), tc.wantLen)
			}
		})
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

type testThreatenedPositionTestCases struct {
	name string
	sq   square.Square
	who  piece.Color
	want bool
}

func TestThreatened(t *testing.T) {
	testCases := []struct {
		name      string
		position  testPosition
		testCases []testThreatenedPositionTestCases
	}{
		{"EmptyBoard",
			map[square.Square]piece.Piece{},
			[]testThreatenedPositionTestCases{
				{"B1-White", square.B1, piece.White, false},
				{"B1-Black", square.B1, piece.Black, false},
				{"B2-White", square.B2, piece.White, false},
				{"B2-Black", square.B2, piece.Black, false},
				{"B3-White", square.B3, piece.White, false},
				{"B3-Black", square.B3, piece.Black, false},
				{"B4-White", square.B4, piece.White, false},
				{"B4-Black", square.B4, piece.Black, false},
				{"B5-White", square.B5, piece.White, false},
				{"B5-Black", square.B5, piece.Black, false},
				{"B6-White", square.B6, piece.White, false},
				{"B6-Black", square.B6, piece.Black, false},
				{"B7-White", square.B7, piece.White, false},
				{"B7-Black", square.B7, piece.Black, false},
				{"B8-White", square.B8, piece.White, false},
				{"B8-Black", square.B8, piece.Black, false},
				{"G1-White", square.G1, piece.White, false},
				{"G1-Black", square.G1, piece.Black, false},
				{"G2-White", square.G2, piece.White, false},
				{"G2-Black", square.G2, piece.Black, false},
				{"G3-White", square.G3, piece.White, false},
				{"G3-Black", square.G3, piece.Black, false},
				{"G4-White", square.G4, piece.White, false},
				{"G4-Black", square.G4, piece.Black, false},
				{"G5-White", square.G5, piece.White, false},
				{"G5-Black", square.G5, piece.Black, false},
				{"G6-White", square.G6, piece.White, false},
				{"G6-Black", square.G6, piece.Black, false},
				{"G7-White", square.G7, piece.White, false},
				{"G7-Black", square.G7, piece.Black, false},
				{"G8-White", square.G8, piece.White, false},
				{"G8-Black", square.G8, piece.Black, false},
				{"A1-NoColor", square.A1, piece.NoColor, false},
				{"A1-Color(5)", square.A1, piece.Color(5), false},
				{"H3-NoColor", square.H3, piece.NoColor, false},
				{"H3-Color(5)", square.H3, piece.Color(5), false},
				{"NoSquare-White", square.NoSquare, piece.White, false},
				{"NoSquare-Black", square.NoSquare, piece.Black, false},
				{"NoSquare-NoColor", square.NoSquare, piece.NoColor, false},
				{"NoSquare-Color(5)", square.NoSquare, piece.Color(5), false},
				{"Square(100)-White", square.Square(100), piece.White, false},
				{"Square(100)-Black", square.Square(100), piece.Black, false},
				{"Square(100)-NoColor", square.Square(100), piece.NoColor, false},
				{"Square(100)-Color(5)", square.Square(100), piece.Color(5), false},
			},
		},
		{"InitialBoard",
			map[square.Square]piece.Piece(nil), // Initial board.
			[]testThreatenedPositionTestCases{
				{"B1-White", square.B1, piece.White, true},
				{"B1-Black", square.B1, piece.Black, false},
				{"B2-White", square.B2, piece.White, true},
				{"B2-Black", square.B2, piece.Black, false},
				{"B3-White", square.B3, piece.White, true},
				{"B3-Black", square.B3, piece.Black, false},
				{"B4-White", square.B4, piece.White, false},
				{"B4-Black", square.B4, piece.Black, false},
				{"B5-White", square.B5, piece.White, false},
				{"B5-Black", square.B5, piece.Black, false},
				{"B6-White", square.B6, piece.White, false},
				{"B6-Black", square.B6, piece.Black, true},
				{"B7-White", square.B7, piece.White, false},
				{"B7-Black", square.B7, piece.Black, true},
				{"B8-White", square.B8, piece.White, false},
				{"B8-Black", square.B8, piece.Black, true},
				{"G1-White", square.G1, piece.White, true},
				{"G1-Black", square.G1, piece.Black, false},
				{"G2-White", square.G2, piece.White, true},
				{"G2-Black", square.G2, piece.Black, false},
				{"G3-White", square.G3, piece.White, true},
				{"G3-Black", square.G3, piece.Black, false},
				{"G4-White", square.G4, piece.White, false},
				{"G4-Black", square.G4, piece.Black, false},
				{"G5-White", square.G5, piece.White, false},
				{"G5-Black", square.G5, piece.Black, false},
				{"G6-White", square.G6, piece.White, false},
				{"G6-Black", square.G6, piece.Black, true},
				{"G7-White", square.G7, piece.White, false},
				{"G7-Black", square.G7, piece.Black, true},
				{"G8-White", square.G8, piece.White, false},
				{"G8-Black", square.G8, piece.Black, true},
			},
		},
		{"InitialBoard-NoPawns",
			map[square.Square]piece.Piece{
				square.A1: piece.New(piece.White, piece.Rook),
				square.B1: piece.New(piece.White, piece.Knight),
				square.C1: piece.New(piece.White, piece.Bishop),
				square.D1: piece.New(piece.White, piece.Queen),
				square.E1: piece.New(piece.White, piece.King),
				square.F1: piece.New(piece.White, piece.Bishop),
				square.G1: piece.New(piece.White, piece.Knight),
				square.H1: piece.New(piece.White, piece.Rook),
				square.A8: piece.New(piece.Black, piece.Rook),
				square.B8: piece.New(piece.Black, piece.Knight),
				square.C8: piece.New(piece.Black, piece.Bishop),
				square.D8: piece.New(piece.Black, piece.Queen),
				square.E8: piece.New(piece.Black, piece.King),
				square.F8: piece.New(piece.Black, piece.Bishop),
				square.G8: piece.New(piece.Black, piece.Knight),
				square.H8: piece.New(piece.Black, piece.Rook),
			},
			[]testThreatenedPositionTestCases{
				{"B1-White", square.B1, piece.White, true},
				{"B1-Black", square.B1, piece.Black, false},
				{"B2-White", square.B2, piece.White, true},
				{"B2-Black", square.B2, piece.Black, false},
				{"B3-White", square.B3, piece.White, true},
				{"B3-Black", square.B3, piece.Black, false},
				{"B4-White", square.B4, piece.White, false},
				{"B4-Black", square.B4, piece.Black, true},
				{"B5-White", square.B5, piece.White, true},
				{"B5-Black", square.B5, piece.Black, false},
				{"B6-White", square.B6, piece.White, false},
				{"B6-Black", square.B6, piece.Black, true},
				{"B7-White", square.B7, piece.White, false},
				{"B7-Black", square.B7, piece.Black, true},
				{"B8-White", square.B8, piece.White, false},
				{"B8-Black", square.B8, piece.Black, true},
				{"G1-White", square.G1, piece.White, true},
				{"G1-Black", square.G1, piece.Black, false},
				{"G2-White", square.G2, piece.White, true},
				{"G2-Black", square.G2, piece.Black, false},
				{"G3-White", square.G3, piece.White, false},
				{"G3-Black", square.G3, piece.Black, false},
				{"G4-White", square.G4, piece.White, true},
				{"G4-Black", square.G4, piece.Black, true},
				{"G5-White", square.G5, piece.White, true},
				{"G5-Black", square.G5, piece.Black, true},
				{"G6-White", square.G6, piece.White, false},
				{"G6-Black", square.G6, piece.Black, false},
				{"G7-White", square.G7, piece.White, false},
				{"G7-Black", square.G7, piece.Black, true},
				{"G8-White", square.G8, piece.White, false},
				{"G8-Black", square.G8, piece.Black, true},
			},
		},
	}
	for _, group := range testCases {
		t.Run(group.name, func(t *testing.T) {
			// Get position.
			p, err := testCasePosition(group.position, nil)
			if err != nil {
				t.Fatalf("Position preparation error: %s", err)
			}

			for _, tc := range group.testCases {
				t.Run(tc.name, func(t *testing.T) {
					defer func() {
						if pm := recover(); pm != nil {
							t.Logf("Position:\n%v", p)
							t.Errorf("*Position.Threatened(%v, %v) should be %v, but panicked with: %v", tc.sq, tc.who, tc.want, pm)
						}
					}()
					if th := p.Threatened(tc.sq, tc.who); th != tc.want {
						t.Logf("Position:\n%v", p)
						t.Errorf("*Position.Threatened(%v, %v) = %v, want %v", tc.sq, tc.who, th, tc.want)
					}
				})
			}
		})
	}
}

func BenchmarkPositionThreatened(b *testing.B) {
	// Get position.
	p, err := testCasePosition(testPosition{
		square.A1: piece.New(piece.White, piece.Rook),
		square.B1: piece.New(piece.White, piece.Knight),
		square.C1: piece.New(piece.White, piece.Bishop),
		square.D1: piece.New(piece.White, piece.Queen),
		square.E1: piece.New(piece.White, piece.King),
		square.F1: piece.New(piece.White, piece.Bishop),
		square.G1: piece.New(piece.White, piece.Knight),
		square.H1: piece.New(piece.White, piece.Rook),
		square.A8: piece.New(piece.Black, piece.Rook),
		square.B8: piece.New(piece.Black, piece.Knight),
		square.C8: piece.New(piece.Black, piece.Bishop),
		square.D8: piece.New(piece.Black, piece.Queen),
		square.E8: piece.New(piece.Black, piece.King),
		square.F8: piece.New(piece.Black, piece.Bishop),
		square.G8: piece.New(piece.Black, piece.Knight),
		square.H8: piece.New(piece.Black, piece.Rook),
	}, nil)
	if err != nil {
		b.Fatalf("Position preparation error: %s", err)
	}
	var th bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for sq := square.Square(0); sq < 64; sq += 1 {
			for _, c := range [2]piece.Color{piece.White, piece.Black} {
				th = p.Threatened(sq, c)
			}
		}
	}
	_ = th
}
