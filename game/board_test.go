package game

import (
	"testing"
)

func piecesOnSquare(b *Board, s Square) int {
	count := 0
	for c := White; c <= Black; c++ {
		for p := Pawn; p <= King; p++ {
			if (b.bitBoard[c][p] & (1 << s)) != 0 {
				count++
			}
		}
	}
	return count
}

func changedbitBoards(before, after *Board) map[Piece]struct{} {
	changed := make(map[Piece]struct{})

	for c := range before.bitBoard {
		for p := range before.bitBoard[Color(c)] {
			if before.bitBoard[Color(c)][p] != after.bitBoard[Color(c)][p] {
				changed[NewPiece(Color(c), PieceType(p))] = struct{}{}
			}
		}
	}
	return changed
}

func TestMovePawn(t *testing.T) {
	beforeMove := NewBoard()
	afterMove := NewBoard()
	afterMove.MakeMove("e2e4")
	changed := changedbitBoards(&beforeMove, &afterMove)
	t.Log("Changed: ", changed)
	if _, c := changed[NewPiece(White, Pawn)]; !c || len(changed) != 1 {
		t.Fail()
	}
}

func TestMoveKnight(t *testing.T) {
	beforeMove := NewBoard()
	afterMove := NewBoard()
	afterMove.MakeMove("b1c3")
	changed := changedbitBoards(&beforeMove, &afterMove)
	t.Log("Changed: ", changed)
	if _, c := changed[NewPiece(White, Knight)]; !c || len(changed) != 1 {
		t.Fail()
	}
}

func TestPutOnOccSquare(t *testing.T) {
	b := NewBoard()
	b.Clear()
	b.QuickPut(NewPiece(White, Pawn), E2)
	b.Put(NewPiece(Black, Queen), E2)
	if b.bitBoard[White][Pawn] != 0 {
		t.Fail()
	}
}
