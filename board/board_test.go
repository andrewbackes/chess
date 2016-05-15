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
			bitprint(b.bitBoard[c][j])
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
