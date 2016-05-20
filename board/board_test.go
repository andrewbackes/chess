package board

import (
	"fmt"
	"github.com/andrewbackes/chess/piece"
	"testing"
)

func piecesOnSquare(b *Board, s Square) int {
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

func changedbitBoards(before, after *Board) map[piece.Piece]struct{} {
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
	afterMove := New()
	afterMove.MakeMove("e2e4")
	changed := changedbitBoards(&beforeMove, &afterMove)
	t.Log("Changed: ", changed)
	if _, c := changed[piece.New(piece.White, piece.Pawn)]; !c || len(changed) != 1 {
		t.Fail()
	}
}

func TestMoveKnight(t *testing.T) {
	beforeMove := New()
	afterMove := New()
	afterMove.MakeMove("b1c3")
	changed := changedbitBoards(&beforeMove, &afterMove)
	t.Log("Changed: ", changed)
	if _, c := changed[piece.New(piece.White, piece.Knight)]; !c || len(changed) != 1 {
		t.Fail()
	}
}

func TestPutOnOccSquare(t *testing.T) {
	b := New()
	b.Clear()
	b.QuickPut(piece.New(piece.White, piece.Pawn), E2)
	b.Put(piece.New(piece.Black, piece.Queen), E2)
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
	if _, ok := s[E1]; !ok {
		t.Fail()
	}
}

// printbitBoards is a helper for diagnosing issues.
func (b *Board) printbitBoards() {
	for c := range b.bitBoard {
		for j := range b.bitBoard[c] {
			fmt.Println(piece.New(piece.Color(c), piece.Type(j)))
			fmt.Println(BitBoard(b.bitBoard[c][j]))
		}
	}
}

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

func TestCapture(t *testing.T) {
	b := New()
	b.Clear()
	b.QuickPut(piece.New(piece.White, piece.King), E4)
	b.QuickPut(piece.New(piece.Black, piece.King), A8)
	b.QuickPut(piece.New(piece.Black, piece.Pawn), E5)
	b.MakeMove(Move("e4e5"))
	for c := piece.White; c <= piece.Black; c++ {
		for p := piece.Pawn; p < piece.King; p++ {
			if b.bitBoard[c][p] != 0 {
				t.Fail()
			}
		}
	}
	if b.bitBoard[piece.White][piece.King] == 0 || b.bitBoard[piece.Black][piece.King] == 0 {
		t.Fail()
	}
}

func TestShortCastle(t *testing.T) {
	b := New()
	b.Clear()
	b.QuickPut(piece.New(piece.Black, piece.King), A8)
	b.QuickPut(piece.New(piece.White, piece.King), E1)
	b.QuickPut(piece.New(piece.White, piece.Rook), H1)
	b.MakeMove(Move("e1g1"))
	if b.bitBoard[piece.White][piece.King] != (1<<G1) ||
		b.bitBoard[piece.White][piece.Rook] != (1<<F1) {
		t.Fail()
	}
}

func TestLongCastle(t *testing.T) {
	b := New()
	b.Clear()
	b.QuickPut(piece.New(piece.Black, piece.King), H8)
	b.QuickPut(piece.New(piece.White, piece.King), E1)
	b.QuickPut(piece.New(piece.White, piece.Rook), A1)
	b.MakeMove(Move("e1c1"))
	if b.bitBoard[piece.White][piece.King] != (1<<C1) ||
		b.bitBoard[piece.White][piece.Rook] != (1<<D1) {
		fmt.Println(BitBoard(b.bitBoard[piece.White][piece.King]))
		fmt.Println("--Rook:--")
		fmt.Println(BitBoard(b.bitBoard[piece.White][piece.Rook]))
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
