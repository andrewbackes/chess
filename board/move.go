package board

// Move is a chess move on the board.
type Move string

// NewMove returns a Move object based on the given from & to squares.
func NewMove(from, to Square) Move {
	return Move(getAlg(from) + getAlg(to))
}
