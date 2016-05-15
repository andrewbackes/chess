package board

import (
	"fmt"
	"github.com/andrewbackes/chess/piece"
	"strings"
)

// Split breaks the move into its source and destination squares.
func Split(m Move) (Square, Square) {
	alg := string(m)
	from := alg[:2]
	to := alg[2:4]
	return ParseSquare(from), ParseSquare(to)
}

// ParseSquare takes the algebraic notation of a square and returns a Square.
func ParseSquare(alg string) Square {
	alg = strings.ToLower(alg)
	f := ((alg[1] - 48) * 8) - (alg[0] - 96)
	return Square(f)
}

func promotedPiece(m Move) piece.Type {
	alg := string(m)
	if len(alg) > 4 {
		p := make(map[string]piece.Type)
		p = map[string]piece.Type{
			"Q": piece.Queen, "N": piece.Knight, "B": piece.Bishop, "R": piece.Rook,
			"q": piece.Queen, "n": piece.Knight, "b": piece.Bishop, "r": piece.Rook,
		}
		return p[string(alg[len(alg)-1])]
	}
	return piece.None
}

func popcount(b uint64) uint {
	var count uint
	for i := uint(0); i < 64; i++ {
		if (b & (1 << i)) != 0 {
			count++
		}
	}
	return count
}

func bitscan(b uint64) uint {
	for i := uint(0); i < 64; i++ {
		if (b & (1 << i)) != 0 {
			return i
		}
	}
	return 64
}

func getAlg(s Square) string {
	var r string

	file := rune(97 + (7 - (int(s) % 8)))
	rank := rune((int(s) / 8) + 49)

	r = string(file) + string(rank)

	return r
}

func bsf(b uint64) uint {
	for i := uint(0); i < 64; i++ {
		if (b & (1 << i)) != 0 {
			return i
		}
	}
	return 64
}

func bsr(b uint64) uint {
	for i := uint(63); i > 0; i-- {
		if (b & (1 << i)) != 0 {
			return i
		}
	}
	if b&1 != 0 {
		return 0
	}
	return 64
}

func bitprint(x uint64) {
	for i := 7; i >= 0; i-- {
		fmt.Printf("%08b\n", (x >> uint64(8*i) & 255))
	}
}
