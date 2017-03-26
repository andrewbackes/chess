package square

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

func TestNewSquare(t *testing.T) {
	squares := []Square{
		NewSquare(1, 1),
		NewSquare(8, 1),
		NewSquare(1, 8),
		NewSquare(8, 8),
	}
	expected := []Square{
		A1, H1, A8, H8,
	}
	for i := range squares {
		if squares[i] != expected[i] {
			t.Error("got", squares[i], int(squares[i]), "but wanted", expected[i], int(expected[i]))
		}
	}
}
