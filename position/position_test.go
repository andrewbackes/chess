package position

import (
	"fmt"
	"testing"

	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
)

func piecesOnSquare(b *Position, s square.Square) int {
	count := 0
	for c := piece.White; c <= piece.Black; c++ {
		for p := piece.Pawn; p <= piece.King; p++ {
			if (b.bitBoard[c][p] & (1 << s)) != 0 {
				count++
			}
		}
	}
	return count
}

func changedBitBoards(before, after *Position) map[piece.Piece]struct{} {
	changed := make(map[piece.Piece]struct{})

	for c := range before.bitBoard {
		for p := range before.bitBoard[piece.Color(c)] {
			if before.bitBoard[piece.Color(c)][p] != after.bitBoard[piece.Color(c)][p] {
				changed[piece.New(piece.Color(c), piece.Type(p))] = struct{}{}
			}
		}
	}
	return changed
}

func TestMovePawn(t *testing.T) {
	beforeMove := New()
	afterMove := beforeMove.MakeMove(move.Parse("e2e4"))
	changed := changedBitBoards(beforeMove, afterMove)
	t.Log("Changed: ", changed)
	if _, c := changed[piece.New(piece.White, piece.Pawn)]; !c || len(changed) != 1 {
		t.Fail()
	}
}

func TestMoveKnight(t *testing.T) {
	beforeMove := New()
	afterMove := beforeMove.MakeMove(move.Parse("b1c3"))
	changed := changedBitBoards(beforeMove, afterMove)
	t.Log("Changed: ", changed)
	if _, c := changed[piece.New(piece.White, piece.Knight)]; !c || len(changed) != 1 {
		t.Fail()
	}
}

func TestPutOnOccSquare(t *testing.T) {
	b := New()
	b.Clear()
	b.QuickPut(piece.New(piece.White, piece.Pawn), square.E2)
	b.Put(piece.New(piece.Black, piece.Queen), square.E2)
	if b.bitBoard[piece.White][piece.Pawn] != 0 {
		t.Fail()
	}
}

func TestFind(t *testing.T) {
	b := New()
	s := b.Find(piece.New(piece.White, piece.King))
	if len(s) != 1 {
		t.Fail()
	}
	if _, ok := s[square.E1]; !ok {
		t.Fail()
	}
}

// printBitBoards is a helper for diagnosing issues.
func (b *Position) printBitBoards() {
	for c := range b.bitBoard {
		for j := range b.bitBoard[c] {
			fmt.Println(piece.New(piece.Color(c), piece.Type(j)))
			fmt.Println(BitBoard(b.bitBoard[c][j]))
		}
	}
}

/*
// TODO(andrewbackes): add more advanced insufficient material checks.
func TestInsufMaterial(t *testing.T) {
	fens := []string{
		"8/8/4kb2/8/8/3K4/8/8 w - - 0 1",
		"8/8/4k3/8/6N1/3K4/8/8 w - - 0 1",
	}
	for _, fen := range fens {
		b, _ := FromFEN(fen)
		if b.InsufficientMaterial() != true {
			t.Fail()
		}
	}
}
*/

func TestKandBvKandOpB(t *testing.T) {
	b := New()
	b.Clear()
	b.QuickPut(piece.New(piece.White, piece.Bishop), square.A1)
	b.QuickPut(piece.New(piece.Black, piece.Bishop), square.B1)
	b.QuickPut(piece.New(piece.White, piece.King), square.A8)
	b.QuickPut(piece.New(piece.Black, piece.King), square.H1)
	if b.InsufficientMaterial() != false {
		t.Fail()
	}
}

/*
func TestBoardPrint(t *testing.T) {
	expected := `+---+---+---+---+---+---+---+---+
| r | n | b | q | k | b | n | r |
+---+---+---+---+---+---+---+---+
| p | p | p | p | p | p | p | p |
+---+---+---+---+---+---+---+---+
|   |   |   |   |   |   |   |   |
+---+---+---+---+---+---+---+---+
|   |   |   |   |   |   |   |   |
+---+---+---+---+---+---+---+---+
|   |   |   |   |   |   |   |   |
+---+---+---+---+---+---+---+---+
|   |   |   |   |   |   |   |   |
+---+---+---+---+---+---+---+---+
| P | P | P | P | P | P | P | P |
+---+---+---+---+---+---+---+---+
| R | N | B | Q | K | B | N | R |
+---+---+---+---+---+---+---+---+
`
	got := New().String()
	if got != expected {
		t.Error(got)
	}
}
*/

func TestCapture(t *testing.T) {
	b := New()
	b.Clear()
	b.QuickPut(piece.New(piece.White, piece.King), square.E4)
	b.QuickPut(piece.New(piece.Black, piece.King), square.A8)
	b.QuickPut(piece.New(piece.Black, piece.Pawn), square.E5)
	result := b.MakeMove(move.Parse("e4e5"))
	for c := piece.White; c <= piece.Black; c++ {
		for p := piece.Pawn; p < piece.King; p++ {
			if result.bitBoard[c][p] != 0 {
				t.Log(BitBoard(result.bitBoard[c][p]))
				t.Fail()
			}
		}
	}
	if result.bitBoard[piece.White][piece.King] == 0 || result.bitBoard[piece.Black][piece.King] == 0 {
		t.Log("==0")
		t.Fail()
	}
}

func TestShortCastle(t *testing.T) {
	b := New()
	b.Clear()
	b.QuickPut(piece.New(piece.Black, piece.King), square.A8)
	b.QuickPut(piece.New(piece.White, piece.King), square.E1)
	b.QuickPut(piece.New(piece.White, piece.Rook), square.H1)
	result := b.MakeMove(move.Parse("e1g1"))
	if result.bitBoard[piece.White][piece.King] != (1<<square.G1) ||
		result.bitBoard[piece.White][piece.Rook] != (1<<square.F1) {
		t.Fail()
	}
}

func TestLongCastle(t *testing.T) {
	b := New()
	b.Clear()
	b.QuickPut(piece.New(piece.Black, piece.King), square.H8)
	b.QuickPut(piece.New(piece.White, piece.King), square.E1)
	b.QuickPut(piece.New(piece.White, piece.Rook), square.A1)
	result := b.MakeMove(move.Parse("e1c1"))
	if result.bitBoard[piece.White][piece.King] != (1<<square.C1) ||
		result.bitBoard[piece.White][piece.Rook] != (1<<square.D1) {
		fmt.Println(BitBoard(result.bitBoard[piece.White][piece.King]))
		fmt.Println("--Rook:--")
		fmt.Println(BitBoard(result.bitBoard[piece.White][piece.Rook]))
		t.Fail()
	}
}

func TestBitBoardPrint(t *testing.T) {
	b := New()
	expected := `00000000
00000000
00000000
00000000
00000000
00000000
11111111
00000000
`
	got := BitBoard(b.bitBoard[piece.White][piece.Pawn]).String()
	if got != expected {
		t.Error(got)
	}

}

func TestSAN(t *testing.T) {
	for _, group := range testSANGroups {
		t.Run(group.Name, func(t *testing.T) {
			for _, tc := range group.TestCases {
				t.Run(tc.Name, func(t *testing.T) {
					// Get position.
					p, err := testCasePosition(group.Position, tc.positionChangerFunc)
					if err != nil {
						t.Fatalf("Position preparation error: %s", err)
					}

					// Call ParseMove function with test case input on Position.
					s := p.SAN(move.Parse(tc.Move))

					// Compare results with expected values.
					if tc.Want != s {
						t.Errorf("*Position.SAN(%s): got \n%s,\n\twant \n%s\n", tc.Move, s, tc.Want)
					}
				})
			}
		})
	}
}

type testSANGroup struct {
	Name      string
	Position  testPosition
	TestCases []testSANTestCase
}

type testSANTestCase struct {
	Name                string
	positionChangerFunc positionChanger
	Move                string
	Want                string
}

var testSANGroups = []testSANGroup{
	{"Misc",
		map[square.Square]piece.Piece{
			square.E1: piece.New(piece.White, piece.King),
			square.F6: piece.New(piece.White, piece.Queen),
			square.A1: piece.New(piece.White, piece.Rook),
			square.D5: piece.New(piece.White, piece.Bishop),
			square.E4: piece.New(piece.White, piece.Pawn),

			square.E8: piece.New(piece.Black, piece.King),
			square.F3: piece.New(piece.Black, piece.Queen),
			square.H8: piece.New(piece.Black, piece.Rook),
			square.D4: piece.New(piece.Black, piece.Bishop),
			square.E5: piece.New(piece.White, piece.Pawn),

			// . . . . k . . r 8
			// . . . . . . . . 7
			// . . . . . Q . . 6
			// . . . B p . . . 5
			// . . . b P . . . 4
			// . . . . . q . . 3
			// . . . . . . . . 2
			// R . . . K . . . 1
			// a b c d e f g h
		},
		[]testSANTestCase{
			{"IllegalMove-White", active(piece.White), "e1e3", ""},
			{"IllegalMove-Black", active(piece.Black), "e8e6", ""},
			{"Move-White", active(piece.White), "e1d2", "Kd2"},
			{"Move-Black", active(piece.Black), "e8d7", "Kd7"},
			{"Move-Check-White", active(piece.White), "a1a8", "Ra8+"},
			{"Move-Check-Black", active(piece.Black), "h8h1", "Rh1+"},
			{"Move-Mate-White", active(piece.White), "d5c6", "Bc6#"},
			{"Move-Mate-Black", active(piece.Black), "d4c3", "Bc3#"},
		},
	},
	{"Promo",
		map[square.Square]piece.Piece{
			square.E1: piece.New(piece.White, piece.King),
			square.C1: piece.New(piece.White, piece.Rook),
			square.C7: piece.New(piece.White, piece.Rook),
			square.G2: piece.New(piece.White, piece.Rook),
			square.H1: piece.New(piece.White, piece.Rook),
			square.B7: piece.New(piece.White, piece.Pawn),
			square.D2: piece.New(piece.White, piece.Pawn),
			square.H7: piece.New(piece.White, piece.Pawn),

			square.E8: piece.New(piece.Black, piece.King),
			square.C2: piece.New(piece.Black, piece.Rook),
			square.C8: piece.New(piece.Black, piece.Rook),
			square.G7: piece.New(piece.Black, piece.Rook),
			square.H8: piece.New(piece.Black, piece.Rook),
			square.B2: piece.New(piece.Black, piece.Pawn),
			square.D7: piece.New(piece.Black, piece.Pawn),
			square.H2: piece.New(piece.Black, piece.Pawn),

			// . . r . k . . r 8
			// . P R p . . r P 7
			// . . . . . . . . 6
			// . . . . . . . . 5
			// . . . . . . . . 4
			// . . . . . . . . 3
			// . p r P . . R p 2
			// . . R . K . . R 1
			// a b c d e f g h
		},
		[]testSANTestCase{
			{"PawnMove-NoPromo-White", active(piece.White), "d2d3", "d3"},
			{"PawnMove-NoPromo-Black", active(piece.Black), "d7d6", "d6"},
			{"RookMove-NoPromo-White", active(piece.White), "g2g3", "Rg3"},
			{"RookMove-NoPromo-Black", active(piece.Black), "g7g6", "Rg6"},
			{"OpponentHomeRankMove-PawnMove-PromoQueen-White", active(piece.White), "b7b8q", "b8=Q"},
			{"OpponentHomeRankMove-PawnMove-PromoQueen-Black", active(piece.Black), "b2b1q", "b1=Q"},
			{"OpponentHomeRankMove-PawnMove-PromoRook-White", active(piece.White), "b7b8r", "b8=R"},
			{"OpponentHomeRankMove-PawnMove-PromoRook-Black", active(piece.Black), "b2b1r", "b1=R"},
			{"OpponentHomeRankMove-PawnMove-PromoBishop-White", active(piece.White), "b7b8b", "b8=B"},
			{"OpponentHomeRankMove-PawnMove-PromoBishop-Black", active(piece.Black), "b2b1b", "b1=B"},
			{"OpponentHomeRankMove-PawnMove-PromoKnight-White", active(piece.White), "b7b8n", "b8=N"},
			{"OpponentHomeRankMove-PawnMove-PromoKnight-Black", active(piece.Black), "b2b1n", "b1=N"},
			{"OpponentHomeRankMove-RookCaptureMove-NoPromo-Check-White", active(piece.White), "c7c8", "Rxc8+"},
			{"OpponentHomeRankMove-RookCaptureMove-NoPromo-Check-Black", active(piece.Black), "c2c1", "Rxc1+"},
			{"OpponentHomeRankMove-PawnCaptureMove-PromoQueen-Check-White", active(piece.White), "b7c8q", "bxc8=Q+"},
			{"OpponentHomeRankMove-PawnCaptureMove-PromoQueen-Check-Black", active(piece.Black), "b2c1q", "bxc1=Q+"},
			{"OpponentHomeRankMove-PawnCaptureMove-PromoRook-Check-White", active(piece.White), "b7c8r", "bxc8=R+"},
			{"OpponentHomeRankMove-PawnCaptureMove-PromoRook-Check-Black", active(piece.Black), "b2c1r", "bxc1=R+"},
			{"OpponentHomeRankMove-PawnCaptureMove-PromoBishop-White", active(piece.White), "b7c8b", "bxc8=B"},
			{"OpponentHomeRankMove-PawnCaptureMove-PromoBishop-Black", active(piece.Black), "b2c1b", "bxc1=B"},
			{"OpponentHomeRankMove-PawnCaptureMove-PromoKnight-White", active(piece.White), "b7c8n", "bxc8=N"},
			{"OpponentHomeRankMove-PawnCaptureMove-PromoKnight-Black", active(piece.Black), "b2c1n", "bxc1=N"},
		},
	},
	{"Pawn",
		map[square.Square]piece.Piece{
			square.G1: piece.New(piece.White, piece.King),
			square.F6: piece.New(piece.White, piece.Queen),
			square.D6: piece.New(piece.White, piece.Rook),
			square.B1: piece.New(piece.White, piece.Knight),
			square.B5: piece.New(piece.White, piece.Pawn),
			square.C4: piece.New(piece.White, piece.Pawn),
			square.E2: piece.New(piece.White, piece.Pawn),
			square.F2: piece.New(piece.White, piece.Pawn),
			square.G2: piece.New(piece.White, piece.Pawn),
			square.H3: piece.New(piece.White, piece.Pawn),
			square.H7: piece.New(piece.White, piece.Pawn),

			square.E8: piece.New(piece.Black, piece.King),
			square.F3: piece.New(piece.Black, piece.Queen),
			square.D3: piece.New(piece.Black, piece.Rook),
			square.G8: piece.New(piece.Black, piece.Knight),
			square.A2: piece.New(piece.Black, piece.Pawn),
			square.B4: piece.New(piece.Black, piece.Pawn),
			square.C5: piece.New(piece.Black, piece.Pawn),
			square.F7: piece.New(piece.Black, piece.Pawn),
			square.E7: piece.New(piece.Black, piece.Pawn),
			square.G7: piece.New(piece.Black, piece.Pawn),
			square.H6: piece.New(piece.Black, piece.Pawn),

			// . . . . k . n . 8
			// . . . . p p p P 7
			// . . . R . Q . p 6
			// . P p . . . . . 5
			// . p P . . . . . 4
			// . . . r . q . P 3
			// p . . . P P P . 2
			// . N . . . . K . 1
			// a b c d e f g h
		},
		[]testSANTestCase{
			{"Move-White", active(piece.White), "e2e3", "e3"},
			{"Move-Black", active(piece.Black), "e7e6", "e6"},
			{"LongMove-White", active(piece.White), "e2e4", "e4"},
			{"LongMove-Black", active(piece.Black), "e7e5", "e5"},
			{"CaptureMove-White", active(piece.White), "e2d3", "exd3"},
			{"CaptureMove-Black", active(piece.Black), "e7d6", "exd6"},
			{"CaptureMove-Ambiguous-White", active(piece.White), "g2f3", "gxf3"},
			{"CaptureMove-Ambiguous-Black", active(piece.Black), "e7f6", "exf6"},
			{"PromoMove-PromoQueen-White", active(piece.White), "h7h8q", "h8=Q"},
			{"PromoMove-PromoQueen-Black", active(piece.Black), "a2a1q", "a1=Q"},
			{"PromoMove-PromoRook-White", active(piece.White), "h7h8r", "h8=R"},
			{"PromoMove-PromoRook-Black", active(piece.Black), "a2a1r", "a1=R"},
			{"PromoMove-PromoBishop-White", active(piece.White), "h7h8b", "h8=B"},
			{"PromoMove-PromoBishop-Black", active(piece.Black), "a2a1b", "a1=B"},
			{"PromoMove-PromoKnight-White", active(piece.White), "h7h8n", "h8=N"},
			{"PromoMove-PromoKnight-Black", active(piece.Black), "a2a1n", "a1=N"},
			{"CapturePromoMove-PromoBishop-White", active(piece.White), "h7g8b", "hxg8=B"},
			{"CapturePromoMove-PromoKnight-Black", active(piece.Black), "a2b1n", "axb1=N"},
			{"CapturePromoMove-PromoQueen-Mate-White", active(piece.White), "h7g8q", "hxg8=Q#"},
			{"CapturePromoMove-PromoQueen-Check-Black", active(piece.Black), "a2b1q", "axb1=Q+"},
			{"CaptureMove-EnPassantOnC6-White", multi(active(piece.White), enPassant(square.C6)), "b5c6", "bxc6"},
			{"CaptureMove-EnPassantOnC3-Black", multi(active(piece.Black), enPassant(square.C3)), "b4c3", "bxc3"},
		},
	},
	{"Knight",
		map[square.Square]piece.Piece{
			square.H1: piece.New(piece.White, piece.King),
			square.A4: piece.New(piece.White, piece.Knight),
			square.B1: piece.New(piece.White, piece.Knight),
			square.B3: piece.New(piece.White, piece.Knight),
			square.B5: piece.New(piece.White, piece.Knight),
			square.B7: piece.New(piece.White, piece.Knight),
			square.C4: piece.New(piece.White, piece.Knight),
			square.D1: piece.New(piece.White, piece.Knight),
			square.D3: piece.New(piece.White, piece.Knight),
			square.D5: piece.New(piece.White, piece.Knight),
			square.G8: piece.New(piece.White, piece.Knight),
			square.H3: piece.New(piece.White, piece.Knight),

			square.H8: piece.New(piece.Black, piece.King),
			square.A5: piece.New(piece.Black, piece.Knight),
			square.B2: piece.New(piece.Black, piece.Knight),
			square.B4: piece.New(piece.Black, piece.Knight),
			square.B6: piece.New(piece.Black, piece.Knight),
			square.B8: piece.New(piece.Black, piece.Knight),
			square.C5: piece.New(piece.Black, piece.Knight),
			square.D2: piece.New(piece.Black, piece.Knight),
			square.D4: piece.New(piece.Black, piece.Knight),
			square.D8: piece.New(piece.Black, piece.Knight),
			square.G1: piece.New(piece.Black, piece.Knight),
			square.H6: piece.New(piece.Black, piece.Knight),

			// . n . n . . N k 8
			// . N . . . . . . 7
			// . n . . . . . n 6
			// n N n N . . . . 5
			// N n N n . . . . 4
			// . N . N . . . N 3
			// . n . n . . . . 2
			// . N . N . . n K 1
			// a b c d e f g h
		},
		[]testSANTestCase{
			{"Move-White", active(piece.White), "h3g5", "Ng5"},
			{"Move-Black", active(piece.Black), "h6g4", "Ng4"},
			{"AmbiguousMove-SpecifyFile-White", active(piece.White), "b5c7", "Nbc7"},
			{"AmbiguousMove-SpecifyFile-Black", active(piece.Black), "b4c2", "Nbc2"},
			{"AmbiguousMove-SpecifyRank-White", active(piece.White), "b1a3", "N1a3"},
			{"AmbiguousMove-SpecifyRank-Black", active(piece.Black), "b8a6", "N8a6"},
			{"AmbiguousMove-SpecifyFileRank-White", active(piece.White), "d1c3", "Nd1c3"},
			{"AmbiguousMove-SpecifyFileRank-Black", active(piece.Black), "d4c6", "Nd4c6"},
			{"CaptureMove-White", active(piece.White), "h3g1", "Nxg1"},
			{"CaptureMove-Black", active(piece.Black), "h6g8", "Nxg8"},
			{"AmbiguousCaptureMove-SpecifyFile-White", active(piece.White), "c4a5", "Ncxa5"},
			{"AmbiguousCaptureMove-SpecifyFile-Black", active(piece.Black), "c5a4", "Ncxa4"},
			{"AmbiguousCaptureMove-SpecifyRank-White", active(piece.White), "b3a5", "N3xa5"},
			{"AmbiguousCaptureMove-SpecifyRank-Black", active(piece.Black), "b2a4", "N2xa4"},
			{"AmbiguousCaptureMove-SpecifyFileRank-White", active(piece.White), "b3c5", "Nb3xc5"},
			{"AmbiguousCaptureMove-SpecifyFileRank-Black", active(piece.Black), "b2c4", "Nb2xc4"},
		},
	},
	{"Bishop",
		map[square.Square]piece.Piece{
			square.E1: piece.New(piece.White, piece.King),
			square.A1: piece.New(piece.White, piece.Bishop),
			square.A3: piece.New(piece.White, piece.Bishop),
			square.B5: piece.New(piece.White, piece.Bishop),
			square.B7: piece.New(piece.White, piece.Bishop),
			square.C1: piece.New(piece.White, piece.Bishop),
			square.C3: piece.New(piece.White, piece.Bishop),
			square.E7: piece.New(piece.White, piece.Bishop),
			square.F1: piece.New(piece.White, piece.Bishop),
			square.F3: piece.New(piece.White, piece.Bishop),
			square.H1: piece.New(piece.White, piece.Bishop),
			square.H3: piece.New(piece.White, piece.Bishop),

			square.E8: piece.New(piece.Black, piece.King),
			square.A6: piece.New(piece.Black, piece.Bishop),
			square.A8: piece.New(piece.Black, piece.Bishop),
			square.B2: piece.New(piece.Black, piece.Bishop),
			square.B4: piece.New(piece.Black, piece.Bishop),
			square.C6: piece.New(piece.Black, piece.Bishop),
			square.C8: piece.New(piece.Black, piece.Bishop),
			square.E2: piece.New(piece.Black, piece.Bishop),
			square.F6: piece.New(piece.Black, piece.Bishop),
			square.F8: piece.New(piece.Black, piece.Bishop),
			square.H6: piece.New(piece.Black, piece.Bishop),
			square.H8: piece.New(piece.Black, piece.Bishop),

			// b . b . k b . b 8
			// . B . . B . . . 7
			// b . b . . b . b 6
			// . B . . . . . . 5
			// . b . . . . . . 4
			// B . B . . B . B 3
			// . b . . b . . . 2
			// B . B . K B . B 1
			// a b c d e f g h
		},
		[]testSANTestCase{
			{"Move-White", active(piece.White), "c1e3", "Be3"},
			{"Move-Black", active(piece.Black), "c8e6", "Be6"},
			{"AmbiguousMove-SpecifyFile-White", active(piece.White), "f3g4", "Bfg4"},
			{"AmbiguousMove-SpecifyFile-Black", active(piece.Black), "f6g5", "Bfg5"},
			{"AmbiguousMove-SpecifyRank-White", active(piece.White), "c1d2", "B1d2"},
			{"AmbiguousMove-SpecifyRank-Black", active(piece.Black), "c8d7", "B8d7"},
			{"AmbiguousMove-SpecifyFileRank-White", active(piece.White), "h1g2", "Bh1g2"},
			{"AmbiguousMove-SpecifyFileRank-Black", active(piece.Black), "h8g7", "Bh8g7"},
			{"CaptureMove-White", active(piece.White), "e7f6", "Bxf6"},
			{"CaptureMove-Black", active(piece.Black), "e2f3", "Bxf3"},
			{"AmbiguousCaptureMove-SpecifyFile-White", active(piece.White), "a3b4", "Baxb4"},
			{"AmbiguousCaptureMove-SpecifyFile-Black", active(piece.Black), "a6b5", "Baxb5"},
			{"AmbiguousCaptureMove-SpecifyRank-White", active(piece.White), "f1e2", "B1xe2"},
			{"AmbiguousCaptureMove-SpecifyRank-Black", active(piece.Black), "f8e7", "B8xe7"},
			{"AmbiguousCaptureMove-SpecifyFileRank-White", active(piece.White), "a1b2", "Ba1xb2"},
			{"AmbiguousCaptureMove-SpecifyFileRank-Black", active(piece.Black), "a8b7", "Ba8xb7"},
		},
	},
	{"Rook",
		map[square.Square]piece.Piece{
			square.D1: piece.New(piece.White, piece.King),
			square.A2: piece.New(piece.White, piece.Rook),
			square.A6: piece.New(piece.White, piece.Rook),
			square.B1: piece.New(piece.White, piece.Rook),
			square.B3: piece.New(piece.White, piece.Rook),
			square.B5: piece.New(piece.White, piece.Rook),
			square.B7: piece.New(piece.White, piece.Rook),
			square.C2: piece.New(piece.White, piece.Rook),
			square.G2: piece.New(piece.White, piece.Rook),
			square.G4: piece.New(piece.White, piece.Rook),
			square.H3: piece.New(piece.White, piece.Rook),

			square.D8: piece.New(piece.Black, piece.King),
			square.A3: piece.New(piece.Black, piece.Rook),
			square.A7: piece.New(piece.Black, piece.Rook),
			square.B2: piece.New(piece.Black, piece.Rook),
			square.B4: piece.New(piece.Black, piece.Rook),
			square.B6: piece.New(piece.Black, piece.Rook),
			square.B8: piece.New(piece.Black, piece.Rook),
			square.C7: piece.New(piece.Black, piece.Rook),
			square.G5: piece.New(piece.Black, piece.Rook),
			square.G7: piece.New(piece.Black, piece.Rook),
			square.H6: piece.New(piece.Black, piece.Rook),

			// . r . k . . . . 8
			// r R r . . . r . 7
			// R r . . . . . r 6
			// . R . . . . r . 5
			// . r . . . . R . 4
			// r R . . . . . R 3
			// R r R . . . R . 2
			// . R . K . . . . 1
			// a b c d e f g h
		},
		[]testSANTestCase{
			{"Move-White", active(piece.White), "g2g1", "Rg1"},
			{"Move-Black", active(piece.Black), "g7g8", "Rg8"},
			{"AmbiguousMove-SpecifyFile-White", active(piece.White), "c2e2", "Rce2"},
			{"AmbiguousMove-SpecifyFile-Black", active(piece.Black), "c7e7", "Rce7"},
			{"AmbiguousMove-SpecifyRank-White", active(piece.White), "g2g3", "R2g3"},
			{"AmbiguousMove-SpecifyRank-Black", active(piece.Black), "g7g6", "R7g6"},
			{"CaptureMove-White", active(piece.White), "h3h6", "Rxh6"},
			{"CaptureMove-Black", active(piece.Black), "h6h3", "Rxh3"},
			{"AmbiguousCaptureMove-SpecifyFile-White", active(piece.White), "b3a3", "Rbxa3"},
			{"AmbiguousCaptureMove-SpecifyFile-Black", active(piece.Black), "b6a6", "Rbxa6"},
			{"AmbiguousCaptureMove-SpecifyRank-White", active(piece.White), "b3b4", "R3xb4"},
			{"AmbiguousCaptureMove-SpecifyRank-Black", active(piece.Black), "b4b5", "R4xb5"},
		},
	},
	{"Queen",
		map[square.Square]piece.Piece{
			square.E1: piece.New(piece.White, piece.King),
			square.A1: piece.New(piece.White, piece.Queen),
			square.A2: piece.New(piece.White, piece.Queen),
			square.B1: piece.New(piece.White, piece.Queen),
			square.G3: piece.New(piece.White, piece.Queen),
			square.G4: piece.New(piece.White, piece.Queen),
			square.H4: piece.New(piece.White, piece.Queen),
			square.A6: piece.New(piece.White, piece.Pawn),
			square.F6: piece.New(piece.White, piece.Pawn),
			square.H6: piece.New(piece.White, piece.Pawn),

			square.E8: piece.New(piece.Black, piece.King),
			square.A7: piece.New(piece.Black, piece.Queen),
			square.A8: piece.New(piece.Black, piece.Queen),
			square.B8: piece.New(piece.Black, piece.Queen),
			square.G5: piece.New(piece.Black, piece.Queen),
			square.G6: piece.New(piece.Black, piece.Queen),
			square.H5: piece.New(piece.Black, piece.Queen),
			square.A3: piece.New(piece.Black, piece.Pawn),
			square.F3: piece.New(piece.Black, piece.Pawn),
			square.H3: piece.New(piece.Black, piece.Pawn),

			// q q . . k . . . 8
			// q . . . . . . . 7
			// P . . . . P q P 6
			// . . . . . . q q 5
			// . . . . . . Q Q 4
			// p . . . . p Q p 3
			// Q . . . . . . . 2
			// Q Q . . K . . . 1
			// a b c d e f g h
		},
		[]testSANTestCase{
			{"Move-White", active(piece.White), "b1c1", "Qc1"},
			{"Move-Black", active(piece.Black), "b8c8", "Qc8"},
			{"AmbiguousMove-SpecifyFile-White", active(piece.White), "b1b2", "Qbb2"},
			{"AmbiguousMove-SpecifyFile-Black", active(piece.Black), "b8b7", "Qbb7"},
			{"AmbiguousMove-SpecifyRank-White", active(piece.White), "a2b2", "Q2b2"},
			{"AmbiguousMove-SpecifyRank-Black", active(piece.Black), "a7b7", "Q7b7"},
			{"AmbiguousMove-SpecifyFileRank-White", active(piece.White), "a1b2", "Qa1b2"},
			{"AmbiguousMove-SpecifyFileRank-Black", active(piece.Black), "a8b7", "Qa8b7"},
			{"CaptureMove-White", active(piece.White), "a2a3", "Qxa3"},
			{"CaptureMove-Black", active(piece.Black), "a7a6", "Qxa6"},
			{"AmbiguousCaptureMove-SpecifyFile-White", active(piece.White), "h4h3", "Qhxh3"},
			{"AmbiguousCaptureMove-SpecifyFile-Black", active(piece.Black), "h5h6", "Qhxh6"},
			{"AmbiguousCaptureMove-SpecifyRank-White", active(piece.White), "g3h3", "Q3xh3"},
			{"AmbiguousCaptureMove-SpecifyRank-Black", active(piece.Black), "g6h6", "Q6xh6"},
			{"AmbiguousCaptureMove-SpecifyFileRank-White", active(piece.White), "g4h3", "Qg4xh3"},
			{"AmbiguousCaptureMove-SpecifyFileRank-Black", active(piece.Black), "g5h6", "Qg5xh6"},
		},
	},
	{"King",
		map[square.Square]piece.Piece{
			square.E1: piece.New(piece.White, piece.King),
			square.A1: piece.New(piece.White, piece.Rook),
			square.H1: piece.New(piece.White, piece.Rook),

			square.E8: piece.New(piece.Black, piece.King),
			square.A8: piece.New(piece.Black, piece.Rook),
			square.H8: piece.New(piece.Black, piece.Rook),

			// r . . . k . . r 8
			// . . . . . . . . 7
			// . . . . . . . . 6
			// . . . . . . . . 5
			// . . . . . . . . 4
			// . . . . . . . . 3
			// . . . . . . . . 2
			// R . . . K . . R 1
			// a b c d e f g h
		},
		[]testSANTestCase{
			{"Castling-KingSide-White", active(piece.White), "e1g1", "O-O"},
			{"Castling-KingSide-Black", active(piece.Black), "e8g8", "O-O"},
			{"Castling-QueenSide-White", active(piece.White), "e1c1", "O-O-O"},
			{"Castling-QueenSide-Black", active(piece.Black), "e8c8", "O-O-O"},
			{"Move-White", active(piece.White), "e1d2", "Kd2"},
			{"Move-Black", active(piece.Black), "e8d7", "Kd7"},
			// Note: To capture, a new piece is added to the board, using the pos positionChangerFunc.
			{"CaptureMove-White", multi(pos(square.D2, piece.New(piece.Black, piece.Rook)), active(piece.White)), "e1d2", "Kxd2"},
			{"CaptureMove-Black", multi(pos(square.D7, piece.New(piece.White, piece.Rook)), active(piece.Black)), "e8d7", "Kxd7"},
		},
	},
}
