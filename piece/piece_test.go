package piece

import (
	"fmt"
	"testing"
)

func TestPiecePrint(t *testing.T) {
	p := New(White, Pawn)
	if fmt.Sprint(p) != "P" {
		t.Fail()
	}
}
