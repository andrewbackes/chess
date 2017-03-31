package position

import (
	"fmt"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
	"testing"
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

func changedbitBoards(before, after *Position) map[piece.Piece]struct{} {
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
	changed := changedbitBoards(beforeMove, afterMove)
	t.Log("Changed: ", changed)
	if _, c := changed[piece.New(piece.White, piece.Pawn)]; !c || len(changed) != 1 {
		t.Fail()
	}
}

func TestMoveKnight(t *testing.T) {
	beforeMove := New()
	afterMove := beforeMove.MakeMove(move.Parse("b1c3"))
	changed := changedbitBoards(beforeMove, afterMove)
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

// printbitBoards is a helper for diagnosing issues.
func (b *Position) printbitBoards() {
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
