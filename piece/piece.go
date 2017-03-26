// Package piece provides tools for working with chess pieces.
package piece

// Color is the color of a piece or square.
type Color uint8

// Possible colors of pieces.
const (
	White      Color = 0
	Black      Color = 1
	Both       Color = 2
	BothColors Color = 2
	Neither    Color = 2
	NoColor    Color = 2
)

// Type is a player's piece. Ex: King, Queen, etc.
type Type uint8

// Possible pieces.
const (
	Pawn Type = iota
	Knight
	Bishop
	Rook
	Queen
	King
	None
)

func (T Type) String() string {
	if T == None {
		return " "
	}
	abbrev := [6]string{"p", "n", "b", "r", "q", "k"}
	return abbrev[T]
}

// Piece represents a chess piece.
type Piece struct {
	Color Color
	Type  Type
}

// String returns a pretty print version of the piece.
func (P Piece) String() string {
	if P.Type == None {
		return " "
	}
	abbrev := [2][6]string{{"P", "N", "B", "R", "Q", "K"}, {"p", "n", "b", "r", "q", "k"}}
	return abbrev[P.Color][P.Type]
}

// Figurine returns a chess piece icon
func (P Piece) Figurine() string {
	if P.Color == NoColor || P.Type == None {
		return " "
	}
	figurines := [][]rune{
		{0x2659, 0x2658, 0x2657, 0x2656, 0x2655, 0x2654},
		{0x265F, 0x265E, 0x265D, 0x265C, 0x265B, 0x265A},
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
