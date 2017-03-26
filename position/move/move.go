// Package move provides a representation for the act of transitioning one chess position to another.
package move

import (
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/square"
	"time"
)

// Move represents the action that transitions one chess position to another.
type Move struct {
	Source      square.Square
	Destination square.Square
	Promote     piece.Type
	Duration    time.Duration
}

var (
	// Null represents a move not occurring.
	Null = Move{square.NoSquare, square.NoSquare, piece.None, 0}
)

// Parse takes a move in PCN format and return a Move struct.
func Parse(algebraic string) *Move {
	if len(algebraic) >= 4 {
		from := square.Parse(algebraic[0:2])
		to := square.Parse(algebraic[2:4])
		promote := piece.None
		if len(algebraic) > 4 {
			p := make(map[string]piece.Type)
			p = map[string]piece.Type{
				"Q": piece.Queen, "N": piece.Knight, "B": piece.Bishop, "R": piece.Rook,
				"q": piece.Queen, "n": piece.Knight, "b": piece.Bishop, "r": piece.Rook,
			}
			promote = p[string(algebraic[len(algebraic)-1])]
		}
		return &Move{
			Source:      from,
			Destination: to,
			Promote:     promote,
		}
	}
	// unexpected input
	return &Null
}

// String will return the move in PCN format.
func (m Move) String() string {
	return m.Source.Algebraic() + m.Destination.Algebraic() + m.Promote.String()
}

// From returns the source square of the move.
func (m *Move) From() square.Square {
	return m.Source
}

// To returns the destination square of the move.
func (m *Move) To() square.Square {
	return m.Destination
}
