// Package position is for working with chess positions. It holds the state of
// a chess game at a particular move.
package position

import (
	"fmt"
	"github.com/andrewbackes/chess/piece"
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

// Position represents the state of a game during a player's turn.
type Position struct {
	// bitBoard has one bitBoard per player per color.
	bitBoard [2][6]uint64 //[player][piece]

	FiftyMoveCount uint64
	EnPassant      Square
	CastlingRights [2][2]bool
	ActiveColor    piece.Color
	MoveNumber     int
}

type Simple struct {
	bitBoard       [2][6]uint64
	EnPassant      Square
	CastlingRights [2][2]bool
	ActiveColor    piece.Color
}

const (
	ShortSide, kingSide uint = 0, 0
	LongSide, queenSide uint = 1, 1
)

// New returns a game board in the opening position. If you want
// a blank board, use Clear().
func New() *Position {
	p := &Position{
		bitBoard:       [2][6]uint64{},
		CastlingRights: [2][2]bool{{true, true}, {true, true}},
		EnPassant:      NoSquare,
		MoveNumber:     1,
		ActiveColor:    piece.White,
	}
	p.Reset()
	return p
}

// Simplify a position into one that is comparable. We just need to exclude the
// move count information.
func Simplify(p *Position) Simple {
	return Simple{
		bitBoard:       p.bitBoard,
		EnPassant:      p.EnPassant,
		CastlingRights: p.CastlingRights,
		ActiveColor:    p.ActiveColor,
	}
}

// String puts the Board into a pretty print-able format.
func (p Position) String() (str string) {
	str += "+---+---+---+---+---+---+---+---+\n"
	for i := 1; i <= 64; i++ {
		square := Square(64 - i)
		str += "|"
		noPiece := true
		for c := range p.bitBoard {
			for j := range p.bitBoard[c] {
				if ((1 << square) & p.bitBoard[c][j]) != 0 {
					str += fmt.Sprint(" ", p.OnSquare(square), " ")
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
func (p *Position) Clear() {
	p.bitBoard = [2][6]uint64{}
}

// Reset puts the pieces in the new game position.
func (p *Position) Reset() {
	// puts the pieces in their starting/newgame positions
	for color := uint(0); color < 2; color = color + 1 {
		//Pawns first:
		p.bitBoard[color][piece.Pawn] = 255 << (8 + (color * 8 * 5))
		//Then the rest of the pieces:
		p.bitBoard[color][piece.Knight] = (1 << (B1 + Square(color*8*7))) ^ (1 << (G1 + Square(color*8*7)))
		p.bitBoard[color][piece.Bishop] = (1 << (C1 + Square(color*8*7))) ^ (1 << (F1 + Square(color*8*7)))
		p.bitBoard[color][piece.Rook] = (1 << (A1 + Square(color*8*7))) ^ (1 << (H1 + Square(color*8*7)))
		p.bitBoard[color][piece.Queen] = (1 << (D1 + Square(color*8*7)))
		p.bitBoard[color][piece.King] = (1 << (E1 + Square(color*8*7)))
	}
}

// OnSquare returns the piece that is on the specified square.
func (p *Position) OnSquare(s Square) piece.Piece {
	for c := piece.White; c <= piece.Black; c++ {
		for pc := piece.Pawn; pc <= piece.King; pc++ {
			if (p.bitBoard[c][pc] & (1 << s)) != 0 {
				return piece.New(c, pc)
			}
		}
	}
	return piece.New(piece.Neither, piece.None)
}

// Occupied returns a bitBoard with all of the specified colors pieces.
func (p *Position) occupied(c piece.Color) uint64 {
	var mask uint64
	for pc := piece.Pawn; pc <= piece.King; pc++ {
		if c == piece.BothColors {
			mask |= p.bitBoard[piece.White][pc] | p.bitBoard[piece.Black][pc]
		} else {
			mask |= p.bitBoard[c][pc]
		}
	}
	return mask
}

func (p *Position) decompose(m Move) (from, to Square, movingPiece, capturedPiece piece.Piece) {
	from, to = Split(m)
	movingPiece = p.OnSquare(from)
	capturedPiece = p.OnSquare(to)
	return
}

// MakeMove attempts to make the given move no matter legality or validity.
// It does not change game state such as en passant or castling rights.
// What ever move you specify will attempt to be made. If it is illegal
// or invalid you will get undetermined behavior.
func (p *Position) MakeMove(m Move) {
	from, to, movingPiece, capturedPiece := p.decompose(m)
	p.adjustMoveCounter(movingPiece, capturedPiece)
	p.adjustCastlingRights(movingPiece, from, to)
	p.adjustEnPassant(movingPiece, from, to)
	p.adjustBoard(m, from, to, movingPiece, capturedPiece)
	p.ActiveColor = (p.ActiveColor + 1) % 2
	if p.ActiveColor == piece.White {
		p.MoveNumber++
	}
}

func (p *Position) adjustMoveCounter(movingPiece, capturedPiece piece.Piece) {
	if capturedPiece.Type != piece.None || movingPiece.Type == piece.Pawn {
		p.FiftyMoveCount = 0
	} else {
		p.FiftyMoveCount++
	}
}

func (p *Position) adjustEnPassant(movingPiece piece.Piece, from, to Square) {
	if movingPiece.Type == piece.Pawn {
		p.EnPassant = NoSquare
		if int(from)-int(to) == 16 || int(from)-int(to) == -16 {
			s := Square(int(from) + []int{8, -8}[movingPiece.Color])
			p.EnPassant = s
		}
	} else {
		p.EnPassant = NoSquare
	}
}

func (p *Position) adjustCastlingRights(movingPiece piece.Piece, from, to Square) {
	for side := ShortSide; side <= LongSide; side++ {
		if movingPiece.Type == piece.King || //King moves
			(movingPiece.Type == piece.Rook &&
				from == [2][2]Square{{H1, A1}, {H8, A8}}[movingPiece.Color][side]) {
			p.CastlingRights[movingPiece.Color][side] = false
		}
		if to == [2][2]Square{{H8, A8}, {H1, A1}}[movingPiece.Color][side] {
			p.CastlingRights[[]piece.Color{piece.Black, piece.White}[movingPiece.Color]][side] = false
		}
	}
}

func (p *Position) adjustBoard(m Move, from, to Square, movingPiece, capturedPiece piece.Piece) {
	// Remove captured piece:
	if capturedPiece.Type != piece.None {
		p.bitBoard[capturedPiece.Color][capturedPiece.Type] ^= (1 << to)
	}

	// Move piece:
	p.bitBoard[movingPiece.Color][movingPiece.Type] ^= ((1 << from) | (1 << to))

	// Castle:
	if movingPiece.Type == piece.King {
		if from == (E1+Square(56*uint8(movingPiece.Color))) && (to == G1+Square(56*uint8(movingPiece.Color))) {
			p.bitBoard[movingPiece.Color][piece.Rook] ^= (1 << (H1 + Square(56*movingPiece.Color))) | (1 << (F1 + Square(56*movingPiece.Color)))
		} else if from == E1+Square(56*uint8(movingPiece.Color)) && (to == C1+Square(56*uint8(movingPiece.Color))) {
			p.bitBoard[movingPiece.Color][piece.Rook] ^= (1 << (A1 + Square(56*(movingPiece.Color)))) | (1 << (D1 + Square(56*(movingPiece.Color))))
		}
	}

	if movingPiece.Type == piece.Pawn {
		// Handle en Passant capture:
		// capturedPiece just means the piece on the destination square
		if (int(to)-int(from))%8 != 0 && capturedPiece.Type == piece.None {
			if movingPiece.Color == piece.White {
				p.bitBoard[piece.Black][piece.Pawn] ^= (1 << (to - 8))
			} else if movingPiece.Color == piece.Black {
				p.bitBoard[piece.White][piece.Pawn] ^= (1 << (to + 8))
			}
		}
		// Handle Promotions:
		promotesTo := promotedPiece(m)
		if promotesTo != piece.None {
			p.bitBoard[movingPiece.Color][movingPiece.Type] ^= (1 << to) // remove piece.Pawn
			p.bitBoard[movingPiece.Color][promotesTo] ^= (1 << to)       // add promoted piece
		}
	}
}

// Put places a piece on the square and removes any other piece
// that may be on that square.
func (p *Position) Put(pp piece.Piece, s Square) {
	pc := p.OnSquare(s)
	if pc.Type != piece.None {
		p.bitBoard[pc.Color][pc.Type] ^= (1 << s)
	}
	p.bitBoard[pp.Color][pp.Type] |= (1 << s)
}

// QuickPut places a piece on the square without removing
// any piece that may already be on that square.
func (p *Position) QuickPut(pc piece.Piece, s Square) {
	p.bitBoard[pc.Color][pc.Type] |= (1 << s)
}

// Find returns the squares that hold the specified piece.
func (p *Position) Find(pc piece.Piece) map[Square]struct{} {
	s := make(map[Square]struct{})
	bits := p.bitBoard[pc.Color][pc.Type]
	for bits != 0 {
		sq := bitscan(bits)
		s[Square(sq)] = struct{}{}
		bits ^= (1 << sq)
	}
	return s
}

func (p *Position) InsufficientMaterial() bool {
	/*
		BUG!
		TODO:
		  	-(Any number of additional bishops of either color on the same color of square due to underpromotion do not affect the situation.)
	*/
	loneKing := []bool{
		p.occupied(piece.White)&p.bitBoard[piece.White][piece.King] == p.occupied(piece.White),
		p.occupied(piece.Black)&p.bitBoard[piece.Black][piece.King] == p.occupied(piece.Black)}

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
			if popcount(p.bitBoard[otherColor][piece.Knight]) == 1 {
				mask := p.bitBoard[otherColor][piece.King] | p.bitBoard[otherColor][piece.Knight]
				occuppied := p.occupied(otherColor)
				if occuppied&mask == occuppied {
					return true
				}
			}
			// King vs King & Bishop
			if popcount(p.bitBoard[otherColor][piece.Bishop]) == 1 {
				mask := p.bitBoard[otherColor][piece.King] | p.bitBoard[otherColor][piece.Bishop]
				occuppied := p.occupied(otherColor)
				if occuppied&mask == occuppied {
					return true
				}
			}
		}
		// King vs King & oppoSite bishop
		kingBishopMask := p.bitBoard[color][piece.King] | p.bitBoard[color][piece.Bishop]
		if (p.occupied(color)&kingBishopMask == p.occupied(color)) && (popcount(p.bitBoard[color][piece.Bishop]) == 1) {
			mask := p.bitBoard[otherColor][piece.King] | p.bitBoard[otherColor][piece.Bishop]
			occuppied := p.occupied(otherColor)
			if (occuppied&mask == occuppied) && (popcount(p.bitBoard[otherColor][piece.Bishop]) == 1) {
				color1 := bitscan(p.bitBoard[color][piece.Bishop]) % 2
				color2 := bitscan(p.bitBoard[otherColor][piece.Bishop]) % 2
				if color1 == color2 {
					return true
				}
			}
		}

	}
	return false
}

// Check returns whether or not the specified color is in check.
func (p *Position) Check(color piece.Color) bool {
	opponent := []piece.Color{piece.Black, piece.White}[color]
	kingsq := Square(bitscan(p.bitBoard[color][piece.King]))
	return p.Threatened(kingsq, opponent)
}
