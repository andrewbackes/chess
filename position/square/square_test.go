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

func TestNew(t *testing.T) {
	squares := []Square{
		New(1, 1),
		New(8, 1),
		New(1, 8),
		New(8, 8),
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
