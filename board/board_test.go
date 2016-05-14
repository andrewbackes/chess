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

// printbitBoards is a helper for diagnosing issues.
func (b *Board) printbitBoards() {
	for c := range b.bitBoard {
		for j := range b.bitBoard[c] {
			fmt.Println(piece.New(piece.Color(c), piece.Type(j)))
			bitprint(b.bitBoard[c][j])
		}
	}
}
