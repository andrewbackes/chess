package square

import (
	"strings"
)

// Square on the board.
type Square uint16

// Squares on the board.
const (
	H1 Square = iota
	G1
	F1
	E1
	D1
	C1
	B1
	A1
	H2
	G2
	F2
	E2
	D2
	C2
	B2
	A2
	H3
	G3
	F3
	E3
	D3
	C3
	B3
	A3
	H4
	G4
	F4
	E4
	D4
	C4
	B4
	A4
	H5
	G5
	F5
	E5
	D5
	C5
	B5
	A5
	H6
	G6
	F6
	E6
	D6
	C6
	B6
	A6
	H7
	G7
	F7
	E7
	D7
	C7
	B7
	A7
	H8
	G8
	F8
	E8
	D8
	C8
	B8
	A8
	NoSquare
)

// LastSquare is the end of the board.
const LastSquare Square = A8

func (s Square) String() string {
	var r string
	file := rune(97 + (7 - (int(s) % 8)))
	rank := rune((int(s) / 8) + 49)
	r = string(file) + string(rank)
	return r
}

// Mask returns a mask to use with bitboards.
func (s Square) Mask() uint64 {
	return mask[s]
}

// New returns a new square based on file and rank.
func New(file, rank uint) Square {
	return Square(((rank - 1) * 8) + (8 - file))
}

// Algebraic returns the file/rank notation of the square.
func (s Square) Algebraic() string {
	var r string
	file := rune(97 + (7 - (int(s) % 8)))
	rank := rune((int(s) / 8) + 49)
	r = string(file) + string(rank)
	return r
}

// Parse takes the algebraic notation of a square and returns a Square.
func Parse(alg string) Square {
	alg = strings.ToLower(alg)
	f := ((alg[1] - 48) * 8) - (alg[0] - 96)
	return Square(f)
}
