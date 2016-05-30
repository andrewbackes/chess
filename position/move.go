package position

// Move is a chess move on the board.
type Move string

const (
	NullMove Move = "0000"
)

// NewMove returns a Move object based on the given from & to squares.
func NewMove(from, to Square) Move {
	return Move(getAlg(from) + getAlg(to))
}
