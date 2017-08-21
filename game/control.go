package game

import (
	"time"
)

// TimeControl represents the time contraints for a chess game.
// TODO(andrewbackes): new TimeControl from parsed string. ex: 40/40
type TimeControl struct {
	Time      time.Duration `json:"time,omitempty"`
	Increment time.Duration `json:"increment,omitempty"`
	Moves     int           `json:"moves,omitempty"`
	Repeating bool          `json:"repeating,omitempty"`
}

// NewTimeControl creates a time control where 'time' is the time per control,
// 'moves' is the number of moves allotted for that time control, 'inc' is the amount
// of time added after each move, and 'repeating' is whether the time control starts
// over once 'moves' has been met.
func NewTimeControl(time time.Duration, moves int, inc time.Duration, repeating bool) TimeControl {
	return TimeControl{
		Time:      time,
		Moves:     moves,
		Increment: inc,
		Repeating: repeating,
	}
}

/*
// Reset adds the time control to the clock current clock value.
func (t *TimeControl) Reset() {
	t.clock += t.Time
	t.movesLeft = t.Moves
}

// Clear sets the clock to zero.
func (t *TimeControl) Clear() {
	t.clock = 0
	t.movesLeft = 0
}
*/
