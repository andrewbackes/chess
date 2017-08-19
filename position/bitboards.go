package position

import (
	"fmt"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/square"
)

// BitBoard is a 64 bit integer where each bit represents a square on
// the chess board. In this implementation, Bit 0 is H1 and bit 63 is A8.
// Since you can only turn a bit on or off, you will need a bitboard for
// each piece type and color to represent a while chess board full of pieces.
// The advantage to this type of chess board representation is that when you
// combine them with pregenerated masks for piece movement at a given square,
// the action of intersecting, bisecting and excluding them becomes a
// bitwise operation.
//
// If you are new to bitboards, Robert Hyatt has a great write-up which
// is more than worth the read:
//   https://www.cis.uab.edu/hyatt/bitmaps.html
type BitBoard uint64

func (b BitBoard) String() string {
	s := ""
	for i := 7; i >= 0; i-- {
		s += fmt.Sprintf("%08b\n", (b >> uint64(8*i) & 255))
	}
	return s
}

type BitBoards map[piece.Color]map[piece.Type]uint64

func (b BitBoards) MailBox() string {
	r := make([]byte, 64, 64)
	for i := uint16(0); i < 64; i++ {
		r[i] = []byte(b.OnSquare(square.Square(i)).String())[0]
	}
	return string(r)
}

func newBitboards() map[piece.Color]map[piece.Type]uint64 {
	m := make(map[piece.Color]map[piece.Type]uint64)
	b := make(map[piece.Type]uint64)
	w := make(map[piece.Type]uint64)
	m[piece.White] = w
	m[piece.Black] = b
	return m
}

// OnSquare returns the piece that is on the specified square.
func (b BitBoards) OnSquare(s square.Square) piece.Piece {
	for c := piece.White; c <= piece.Black; c++ {
		for pc := piece.Pawn; pc <= piece.King; pc++ {
			if (b[c][pc] & (1 << s)) != 0 {
				return piece.New(c, pc)
			}
		}
	}
	return piece.New(piece.Neither, piece.None)
}
