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
