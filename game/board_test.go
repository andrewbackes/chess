package game

import (
	"testing"
)

func piecesOnSquare(b *Board, s Square) int {
	count := 0
	for c := White; c <= Black; c++ {
		for p := Pawn; p <= King; p++ {
			if (b.BitBoard[c][p] & (1 << s)) != 0 {
				count++
			}
		}
	}
	return count
}

func changedBitBoards(before, after *Board) map[Piece]struct{} {
	changed := make(map[Piece]struct{})

	for c := range before.BitBoard {
		for p := range before.BitBoard[Color(c)] {
			if before.BitBoard[Color(c)][p] != after.BitBoard[Color(c)][p] {
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
	changed := changedBitBoards(&beforeMove, &afterMove)
	t.Log("Changed: ", changed)
	if _, c := changed[NewPiece(White, Pawn)]; !c || len(changed) != 1 {
		t.Fail()
	}
}

func TestMoveKnight(t *testing.T) {
	beforeMove := NewBoard()
	afterMove := NewBoard()
	afterMove.MakeMove("b1c3")
	changed := changedBitBoards(&beforeMove, &afterMove)
	t.Log("Changed: ", changed)
	if _, c := changed[NewPiece(White, Knight)]; !c || len(changed) != 1 {
		t.Fail()
	}
}
