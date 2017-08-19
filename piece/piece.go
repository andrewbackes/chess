// Package piece provides tools for working with chess pieces.
package piece

import (
	"strings"
)

// Type is a player's piece. Ex: King, Queen, etc.
type Type uint8

// Possible pieces.
const (
	None Type = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

func (T Type) String() string {
	return map[Type]string{
		None:   " ",
		Pawn:   "p",
		Knight: "n",
		Bishop: "b",
		Rook:   "r",
		Queen:  "q",
		King:   "k",
	}[T]
}

// Piece represents a chess piece.
type Piece struct {
	Color Color
	Type  Type
}

// String returns a pretty print version of the piece.
func (P Piece) String() string {
	t := P.Type.String()
	if P.Color == White {
		return strings.ToUpper(t)
	}
	return t
}

// Figurine returns a chess piece icon
func (P Piece) Figurine() string {
	if P.Color == NoColor || P.Type == None {
		return " "
	}
	return string(map[Color]map[Type]rune{
		White: {Pawn: 0x2659, Knight: 0x2658, Bishop: 0x2657, Rook: 0x2656, Queen: 0x2655, King: 0x2654},
		Black: {Pawn: 0x265F, Knight: 0x265E, Bishop: 0x265D, Rook: 0x265C, Queen: 0x265B, King: 0x265A},
	}[P.Color][P.Type])
}

// New returns a new chess piece type.
func New(c Color, t Type) Piece {
	return Piece{
		Color: c,
		Type:  t,
	}
}
