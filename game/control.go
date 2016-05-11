package game

import (
	"time"
)

type TimeControl struct {
	Time      time.Duration
	BonusTime time.Duration
	Moves     int64
	Repeating bool
	clock     time.Duration
	movesLeft int64
}

func (t *TimeControl) Reset() {
	t.clock = t.Time
	t.movesLeft = t.Moves
}
