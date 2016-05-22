// Package board provides tools for representing a chess board.
package board

import (
	"errors"
	"fmt"
	"github.com/andrewbackes/chess/piece"
	"strconv"
	"strings"
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

// Board is a representation of a chess Board.
type Board struct {
	// bitBoard has one bitBoard per player per color.
	bitBoard [2][6]uint64 //[player][piece]
}

const (
	ShortSide, kingSide uint = 0, 0
	LongSide, queenSide uint = 1, 1
)

// New returns a game board in the opening position. If you want
// a blank board, use Clear().
func New() Board {
	b := Board{bitBoard: [2][6]uint64{}}
	b.Reset()
	return b
}

// String puts the Board into a pretty print-able format.
func (b Board) String() (str string) {
	str += "+---+---+---+---+---+---+---+---+\n"
	for i := 1; i <= 64; i++ {
		square := Square(64 - i)
		str += "|"
		noPiece := true
		for c := range b.bitBoard {
			for j := range b.bitBoard[c] {
				if ((1 << square) & b.bitBoard[c][j]) != 0 {
					str += fmt.Sprint(" ", b.OnSquare(square), " ")
					noPiece = false
				}
			}
		}
		if noPiece {
			str += "   "
		}
		if square%8 == 0 {
			str += "|\n"
			str += "+---+---+---+---+---+---+---+---+"
			if square < LastSquare {
				str += "\n"
			}
		}
	}
	return
}

// Clear empties the Board.
func (b *Board) Clear() {
	b.bitBoard = [2][6]uint64{}
}

// Reset puts the pieces in the new game position.
func (b *Board) Reset() {
	// puts the pieces in their starting/newgame positions
	for color := uint(0); color < 2; color = color + 1 {
		//Pawns first:
		b.bitBoard[color][piece.Pawn] = 255 << (8 + (color * 8 * 5))
		//Then the rest of the pieces:
		b.bitBoard[color][piece.Knight] = (1 << (B1 + Square(color*8*7))) ^ (1 << (G1 + Square(color*8*7)))
		b.bitBoard[color][piece.Bishop] = (1 << (C1 + Square(color*8*7))) ^ (1 << (F1 + Square(color*8*7)))
		b.bitBoard[color][piece.Rook] = (1 << (A1 + Square(color*8*7))) ^ (1 << (H1 + Square(color*8*7)))
		b.bitBoard[color][piece.Queen] = (1 << (D1 + Square(color*8*7)))
		b.bitBoard[color][piece.King] = (1 << (E1 + Square(color*8*7)))
	}
}

// OnSquare returns the piece that is on the specified square.
func (b *Board) OnSquare(s Square) piece.Piece {
	for c := piece.White; c <= piece.Black; c++ {
		for p := piece.Pawn; p <= piece.King; p++ {
			if (b.bitBoard[c][p] & (1 << s)) != 0 {
				return piece.New(c, p)
			}
		}
	}
	return piece.New(piece.Neither, piece.None)
}

// Occupied returns a bitBoard with all of the specified colors pieces.
func (b *Board) occupied(c piece.Color) uint64 {
	var mask uint64
	for p := piece.Pawn; p <= piece.King; p++ {
		if c == piece.BothColors {
			mask |= b.bitBoard[piece.White][p] | b.bitBoard[piece.Black][p]
		} else {
			mask |= b.bitBoard[c][p]
		}
	}
	return mask
}

// MakeMove attempts to make the given move no matter legality or validity.
// It does not change game state such as en passant or castling rights.
// What ever move you specify will attempt to be made. If it is illegal
// or invalid you will get undetermined behavior.
func (b *Board) MakeMove(m Move) {
	from, to := Split(m)
	movingPiece := b.OnSquare(from)
	capturedPiece := b.OnSquare(to)

	// Remove captured piece:
	if capturedPiece.Type != piece.None {
		b.bitBoard[capturedPiece.Color][capturedPiece.Type] ^= (1 << to)
	}

	// Move piece:
	b.bitBoard[movingPiece.Color][movingPiece.Type] ^= ((1 << from) | (1 << to))

	// Castle:
	if movingPiece.Type == piece.King {
		if from == (E1+Square(56*uint8(movingPiece.Color))) && (to == G1+Square(56*uint8(movingPiece.Color))) {
			b.bitBoard[movingPiece.Color][piece.Rook] ^= (1 << (H1 + Square(56*movingPiece.Color))) | (1 << (F1 + Square(56*movingPiece.Color)))
		} else if from == E1+Square(56*uint8(movingPiece.Color)) && (to == C1+Square(56*uint8(movingPiece.Color))) {
			b.bitBoard[movingPiece.Color][piece.Rook] ^= (1 << (A1 + Square(56*(movingPiece.Color)))) | (1 << (D1 + Square(56*(movingPiece.Color))))
		}
	}

	if movingPiece.Type == piece.Pawn {
		// Handle en Passant capture:
		// capturedPiece just means the piece on the destination square
		if (int(to)-int(from))%8 != 0 && capturedPiece.Type == piece.None {
			if movingPiece.Color == piece.White {
				b.bitBoard[piece.Black][piece.Pawn] ^= (1 << (to - 8))
			} else if movingPiece.Color == piece.Black {
				b.bitBoard[piece.White][piece.Pawn] ^= (1 << (to + 8))
			}
		}
		// Handle Promotions:
		promotesTo := promotedPiece(m)
		if promotesTo != piece.None {
			b.bitBoard[movingPiece.Color][movingPiece.Type] ^= (1 << to) // remove piece.Pawn
			b.bitBoard[movingPiece.Color][promotesTo] ^= (1 << to)       // add promoted piece
		}
	}

}

// GameFromFEN parses the board passed via FEN and returns a board object.
func GameFromFEN(position string) (*Board, error) {
	b := New()
	b.Clear()
	// remove the /'s and replace the numbers with that many spaces
	// so that there is a 1-1 mapping from bytes to squares.
	justBoard := strings.Split(position, " ")[0]
	parsedBoard := strings.Replace(justBoard, "/", "", 9)
	for i := 1; i < 9; i++ {
		parsedBoard = strings.Replace(parsedBoard, strconv.Itoa(i), strings.Repeat(" ", i), -1)
	}
	if len(parsedBoard) < 64 {
		return nil, errors.New("fen: could not parse position")
	}
	p := map[rune]piece.Type{
		'P': piece.Pawn, 'p': piece.Pawn,
		'N': piece.Knight, 'n': piece.Knight,
		'B': piece.Bishop, 'b': piece.Bishop,
		'R': piece.Rook, 'r': piece.Rook,
		'Q': piece.Queen, 'q': piece.Queen,
		'K': piece.King, 'k': piece.King}
	color := map[rune]piece.Color{
		'P': piece.White, 'p': piece.Black,
		'N': piece.White, 'n': piece.Black,
		'B': piece.White, 'b': piece.Black,
		'R': piece.White, 'r': piece.Black,
		'Q': piece.White, 'q': piece.Black,
		'K': piece.White, 'k': piece.Black}
	// adjust the bitboards:
	for pos := 0; pos < len(parsedBoard); pos++ {
		if pos > 64 {
			break
		}
		k := rune(parsedBoard[pos])
		if _, ok := p[k]; ok {
			b.bitBoard[color[k]][p[k]] |= (1 << uint(63-pos))
		}
	}
	return &b, nil
}

// Put places a piece on the square and removes any other piece
// that may be on that square.
func (b *Board) Put(p piece.Piece, s Square) {
	pc := b.OnSquare(s)
	if pc.Type != piece.None {
		b.bitBoard[pc.Color][pc.Type] ^= (1 << s)
	}
	b.bitBoard[p.Color][p.Type] |= (1 << s)
}

// QuickPut places a piece on the square without removing
// any piece that may already be on that square.
func (b *Board) QuickPut(p piece.Piece, s Square) {
	b.bitBoard[p.Color][p.Type] |= (1 << s)
}

// Find returns the squares that hold the specified piece.
func (b *Board) Find(p piece.Piece) map[Square]struct{} {
	s := make(map[Square]struct{})
	bits := b.bitBoard[p.Color][p.Type]
	for bits != 0 {
		sq := bitscan(bits)
		s[Square(sq)] = struct{}{}
		bits ^= (1 << sq)
	}
	return s
}

func (b *Board) InsufficientMaterial() bool {
	/*
		BUG!
		TODO:
		  	-(Any number of additional bishops of either color on the same color of square due to underpromotion do not affect the situation.)
	*/
	loneKing := []bool{
		b.occupied(piece.White)&b.bitBoard[piece.White][piece.King] == b.occupied(piece.White),
		b.occupied(piece.Black)&b.bitBoard[piece.Black][piece.King] == b.occupied(piece.Black)}

	if !loneKing[piece.White] && !loneKing[piece.Black] {
		return false
	}

	for color := piece.White; color <= piece.Black; color++ {
		otherColor := []piece.Color{piece.Black, piece.White}[color]
		if loneKing[color] {
			// King vs King:
			if loneKing[otherColor] {
				return true
			}
			// King vs King & Knight
			if popcount(b.bitBoard[otherColor][piece.Knight]) == 1 {
				mask := b.bitBoard[otherColor][piece.King] | b.bitBoard[otherColor][piece.Knight]
				occuppied := b.occupied(otherColor)
				if occuppied&mask == occuppied {
					return true
				}
			}
			// King vs King & Bishop
			if popcount(b.bitBoard[otherColor][piece.Bishop]) == 1 {
				mask := b.bitBoard[otherColor][piece.King] | b.bitBoard[otherColor][piece.Bishop]
				occuppied := b.occupied(otherColor)
				if occuppied&mask == occuppied {
					return true
				}
			}
		}
		// King vs King & oppoSite bishop
		kingBishopMask := b.bitBoard[color][piece.King] | b.bitBoard[color][piece.Bishop]
		if (b.occupied(color)&kingBishopMask == b.occupied(color)) && (popcount(b.bitBoard[color][piece.Bishop]) == 1) {
			mask := b.bitBoard[otherColor][piece.King] | b.bitBoard[otherColor][piece.Bishop]
			occuppied := b.occupied(otherColor)
			if (occuppied&mask == occuppied) && (popcount(b.bitBoard[otherColor][piece.Bishop]) == 1) {
				color1 := bitscan(b.bitBoard[color][piece.Bishop]) % 2
				color2 := bitscan(b.bitBoard[otherColor][piece.Bishop]) % 2
				if color1 == color2 {
					return true
				}
			}
		}

	}
	return false
}

// Check returns whether or not the specified color is in check.
func (b *Board) Check(color piece.Color) bool {
	opponent := []piece.Color{piece.Black, piece.White}[color]
	kingsq := Square(bitscan(b.bitBoard[color][piece.King]))
	return b.Threatened(kingsq, opponent)
}
