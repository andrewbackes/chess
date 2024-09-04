// Package move provides a representation for the act of transitioning one chess position to another.
package move

import (
	"time"

	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/square"
)

// Move represents the action that transitions one chess position to another.
type Move struct {
	Source      square.Square `json:"source"`
	Destination square.Square `json:"destination"`
	Promote     piece.Type    `json:"promote,omitempty"`
	Duration    time.Duration `json:"duration,omitempty"`
}

var (
	// Null represents a move not occurring.
	Null = Move{Source: square.NoSquare, Destination: square.NoSquare, Promote: piece.None, Duration: 0}
)

// Parse takes a move in PCN format and return a Move struct.
func Parse(algebraic string) Move {
	if algebraic == "0000" || algebraic == "" {
		return Null
	}
	if len(algebraic) >= 4 {
		from := square.Parse(algebraic[0:2])
		to := square.Parse(algebraic[2:4])
		promote := piece.None
		if len(algebraic) > 4 {
			switch string(algebraic[len(algebraic)-1]) {
			case "Q", "q":
				promote = piece.Queen
			case "R", "r":
				promote = piece.Rook
			case "B", "b":
				promote = piece.Bishop
			case "N", "n":
				promote = piece.Knight
			}
		}
		return Move{
			Source:      from,
			Destination: to,
			Promote:     promote,
		}
	}
	// unexpected input
	return Null
}

// String will return the move in PCN format.
func (m Move) String() string {
	if m.Promote.String() != " " {
		return m.Source.Algebraic() + m.Destination.Algebraic() + m.Promote.String()
	}
	return m.Source.Algebraic() + m.Destination.Algebraic()
}

// From returns the source square of the move.
func (m Move) From() square.Square {
	return m.Source
}

// To returns the destination square of the move.
func (m Move) To() square.Square {
	return m.Destination
}
