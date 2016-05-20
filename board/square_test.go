package board

import (
	"testing"
)

func TestSquareToString(t *testing.T) {
	if Square(A8).String() != "a8" ||
		Square(H1).String() != "h1" ||
		Square(D4).String() != "d4" {
		t.Fail()
	}
}
