package game

import (
	"fmt"
	"strings"
)

func toSquare(alg string) Square {
	alg = strings.ToLower(alg)
	f := ((alg[1] - 48) * 8) - (alg[0] - 96)
	return Square(f)
}

func getSquares(m Move) (Square, Square) {
	alg := string(m)
	from := alg[:2]
	to := alg[2:4]
	return toSquare(from), toSquare(to)
}

func promotedPiece(m Move) PieceType {
	alg := string(m)
	if len(alg) > 4 {
		p := make(map[string]PieceType)
		p = map[string]PieceType{"Q": Queen, "N": Knight, "B": Bishop, "R": Rook,
			"q": Queen, "n": Knight, "b": Bishop, "r": Rook}
		return p[string(alg[len(alg)-1])]
	}
	return None
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
