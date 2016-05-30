package chess

import (
	"testing"
	"time"
)

func TestClearTC(t *testing.T) {
	tc := NewTimeControl(time.Minute, 40, time.Second, true)
	tc.Clear()
	if tc.clock != 0 || tc.movesLeft != 0 {
		t.Fail()
	}
}
