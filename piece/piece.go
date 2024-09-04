// Package piece provides tools for working with chess pieces.
package piece

import (
	"strings"
)

// Type is a player's piece. Ex: King, Queen, etc.
type Type uint8

const TYPE_COUNT = 7 // Including piece.None. Improves code readability.

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

var typeStrings = [TYPE_COUNT]string{" ", "p", "n", "b", "r", "q", "k"}

func (T Type) String() string {
	if T >= TYPE_COUNT {
		return ""
	}
	return typeStrings[T]
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

var figurines = [COLOR_COUNT][TYPE_COUNT]rune{
	[TYPE_COUNT]rune{' ', 0x2659, 0x2658, 0x2657, 0x2656, 0x2655, 0x2654},
	[TYPE_COUNT]rune{' ', 0x265F, 0x265E, 0x265D, 0x265C, 0x265B, 0x265A},
}

// Figurine returns a chess piece icon
func (P Piece) Figurine() string {
	if P.Color == NoColor || P.Type == None {
		return " "
	}
	if P.Color >= COLOR_COUNT || P.Type >= TYPE_COUNT {
		return "\x00"
	}
	return string(figurines[P.Color][P.Type])
}

// New returns a new chess piece type.
func New(c Color, t Type) Piece {
	return Piece{
		Color: c,
		Type:  t,
	}
}
