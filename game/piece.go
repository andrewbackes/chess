package game

// Piece represents a chess piece.
type Piece struct {
	Color Color
	Type  PieceType
}

// String returns a pretty print version of the piece.
func (P Piece) String() string {
	if P.Type == None {
		return " "
	}
	abbrev := [2][6]string{{"P", "N", "B", "R", "Q", "K"}, {"p", "n", "b", "r", "q", "k"}}
	return abbrev[P.Color][P.Type]
}

// NewPiece returns a new chess piece type.
func NewPiece(c Color, t PieceType) Piece {
	return Piece{
		Color: c,
		Type:  t,
	}
}
