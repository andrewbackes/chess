// Package board holds the chess board object and methods for interacting with it.
// Bitboards, 'mailbox' view of the board, etc.
package board

import (
	"fmt"
)

// Piece represents a chess piece.
type Piece struct {
	Color Color
	Type  PieceType
}

func piece(c Color, t PieceType) Piece {
	return Piece{
		Color: c,
		Type:  t,
	}
}

// Board is a chess board.
type Board struct {
	// PieceBB has one bitboard per player per color.
	// So
	PieceBB [2][6]uint64 //[player][piece]
}

// String puts the board into a pretty print-able format.
func (B *Board) String() (str string) {
	abbrev := [2][6]string{{"P", "N", "B", "R", "Q", "K"}, {"p", "n", "b", "r", "q", "k"}}

	str += fmt.Sprintln("+---+---+---+---+---+---+---+---+")
	for i := 1; i <= 64; i++ {
		square := uint(64 - i)

		str += fmt.Sprint("|")
		blankSquare := true
		for j := Pawn; j <= King; j = j + 1 {
			for color := Color(0); color <= Black; color++ {
				if ((1 << square) & B.PieceBB[color][j]) != 0 {
					str += fmt.Sprint(" ", abbrev[color][j], " ")
					blankSquare = false
				}
			}
		}
		if blankSquare == true {
			str += fmt.Sprint("   ")
		}
		if square%8 == 0 {
			str += fmt.Sprintln("|")
			str += fmt.Sprintln("+---+---+---+---+---+---+---+---+")
		}
	}
	return
}

// Clear empties the board.
func (B *Board) Clear() {
	B.PieceBB = [2][6]uint64{}
}

// Reset puts the pieces in the new game position.
func (B *Board) Reset() {
	// puts the pieces in their starting/newgame positions
	for color := uint(0); color < 2; color = color + 1 {
		//Pawns first:
		B.PieceBB[color][Pawn] = 255 << (8 + (color * 8 * 5))
		//Then the rest of the pieces:
		B.PieceBB[color][Knight] = (1 << (B1 + (color * 8 * 7))) ^ (1 << (G1 + (color * 8 * 7)))
		B.PieceBB[color][Bishop] = (1 << (C1 + (color * 8 * 7))) ^ (1 << (F1 + (color * 8 * 7)))
		B.PieceBB[color][Rook] = (1 << (A1 + (color * 8 * 7))) ^ (1 << (H1 + (color * 8 * 7)))
		B.PieceBB[color][Queen] = (1 << (D1 + (color * 8 * 7)))
		B.PieceBB[color][King] = (1 << (E1 + (color * 8 * 7)))
	}
}

// OnSquare returns the piece that is on the specified square.
func (B *Board) OnSquare(s Square) Piece {
	for c := White; c <= Black; c++ {
		for p := Pawn; p <= King; p++ {
			if (B.PieceBB[c][p] & (1 << s)) != 0 {
				return piece(c, p)
			}
		}
	}
	return piece(Neither, None)
}

// Occupied returns a bitboard with all of the specified colors pieces.
func (B *Board) Occupied(c Color) uint64 {
	var mask uint64
	for p := Pawn; p <= King; p++ {
		if c == Both {
			mask |= B.PieceBB[White][p] | B.PieceBB[Black][p]
		} else {
			mask |= B.PieceBB[c][p]
		}
	}
	return mask
}
