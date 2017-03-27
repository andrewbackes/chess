// Package position is for working with chess positions. It holds the state of
// a chess game at a particular move.
package position

import (
	"fmt"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/move"
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

// Position represents the state of a game during a player's turn.
type Position struct {
	// bitBoard has one bitBoard per player per color.
	bitBoard       map[piece.Color]map[piece.Type]uint64
	FiftyMoveCount uint64        `json:"fiftyMoveCount,omitempty" bson:"fiftyMoveCount,omitempty"`
	EnPassant      square.Square `json:"enPassant,omitempty" bson:"enPassant,omitempty"`
	CastlingRights [2][2]bool    `json:"castlingRights" bson:"castlingRights"`
	ActiveColor    piece.Color   `json:"activeColor" bson:"activeColor"`
	MoveNumber     int           `json:"moveNumber" bson:"moveNumber"`
}

type Simple struct {
	bitBoard       map[piece.Color]map[piece.Type]uint64
	EnPassant      square.Square
	CastlingRights [2][2]bool
	ActiveColor    piece.Color
}

const (
	ShortSide, kingSide uint = 0, 0
	LongSide, queenSide uint = 1, 1
)

func newBitboards() map[piece.Color]map[piece.Type]uint64 {
	m := make(map[piece.Color]map[piece.Type]uint64)
	b := make(map[piece.Type]uint64)
	w := make(map[piece.Type]uint64)
	m[piece.White] = w
	m[piece.Black] = b
	return m
}

// New returns a game board in the opening position. If you want
// a blank board, use Clear().
func New() *Position {
	p := &Position{
		bitBoard:       newBitboards(),
		CastlingRights: [2][2]bool{{true, true}, {true, true}},
		EnPassant:      square.NoSquare,
		MoveNumber:     1,
		ActiveColor:    piece.White,
	}
	p.Reset()
	return p
}

// Copy a position.
func Copy(p *Position) *Position {
	n := &Position{
		bitBoard:       newBitboards(),
		EnPassant:      p.EnPassant,
		CastlingRights: p.CastlingRights,
		ActiveColor:    p.ActiveColor,
		MoveNumber:     p.MoveNumber,
	}
	for k, v := range p.bitBoard[piece.White] {
		n.bitBoard[piece.White][k] = v
	}
	for k, v := range p.bitBoard[piece.Black] {
		n.bitBoard[piece.Black][k] = v
	}
	return n
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
		sq := square.Square(64 - i)
		str += "|"
		noPiece := true
		for c := range p.bitBoard {
			for j := range p.bitBoard[c] {
				if ((1 << sq) & p.bitBoard[c][j]) != 0 {
					str += fmt.Sprint(" ", p.OnSquare(sq), " ")
					noPiece = false
				}
			}
		}
		if noPiece {
			str += "   "
		}
		if sq%8 == 0 {
			str += "|\n"
			str += "+---+---+---+---+---+---+---+---+"
			if sq < square.LastSquare {
				str += "\n"
			}
		}
	}
	return
}

// Clear empties the Board.
func (p *Position) Clear() {
	p.bitBoard = newBitboards()
}

// Reset puts the pieces in the new game position.
func (p *Position) Reset() {
	// puts the pieces in their starting/newgame positions
	for color := piece.Color(0); color < 2; color = color + 1 {
		//Pawns first:
		p.bitBoard[color][piece.Pawn] = 255 << (8 + (color * 8 * 5))
		//Then the rest of the pieces:
		p.bitBoard[color][piece.Knight] = (1 << (square.B1 + square.Square(color*8*7))) ^ (1 << (square.G1 + square.Square(color*8*7)))
		p.bitBoard[color][piece.Bishop] = (1 << (square.C1 + square.Square(color*8*7))) ^ (1 << (square.F1 + square.Square(color*8*7)))
		p.bitBoard[color][piece.Rook] = (1 << (square.A1 + square.Square(color*8*7))) ^ (1 << (square.H1 + square.Square(color*8*7)))
		p.bitBoard[color][piece.Queen] = (1 << (square.D1 + square.Square(color*8*7)))
		p.bitBoard[color][piece.King] = (1 << (square.E1 + square.Square(color*8*7)))
	}
}

// OnSquare returns the piece that is on the specified square.
func (p *Position) OnSquare(s square.Square) piece.Piece {
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

func (p *Position) decompose(m *move.Move) (from, to square.Square, movingPiece, capturedPiece piece.Piece) {
	return m.From(), m.To(), p.OnSquare(m.From()), p.OnSquare(m.To())
}

// MakeMove attempts to make the given move no matter legality or validity.
// It does not change game state such as en passant or castling rights.
// What ever move you specify will attempt to be made. If it is illegal
// or invalid you will get undetermined behavior.
func (p *Position) MakeMove(m *move.Move) {
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

func (p *Position) adjustEnPassant(movingPiece piece.Piece, from, to square.Square) {
	if movingPiece.Type == piece.Pawn {
		p.EnPassant = square.NoSquare
		if int(from)-int(to) == 16 || int(from)-int(to) == -16 {
			s := square.Square(int(from) + []int{8, -8}[movingPiece.Color])
			p.EnPassant = s
		}
	} else {
		p.EnPassant = square.NoSquare
	}
}

func (p *Position) adjustCastlingRights(movingPiece piece.Piece, from, to square.Square) {
	for side := ShortSide; side <= LongSide; side++ {
		if movingPiece.Type == piece.King || //King moves
			(movingPiece.Type == piece.Rook &&
				from == [2][2]square.Square{{square.H1, square.A1}, {square.H8, square.A8}}[movingPiece.Color][side]) {
			p.CastlingRights[movingPiece.Color][side] = false
		}
		if to == [2][2]square.Square{{square.H8, square.A8}, {square.H1, square.A1}}[movingPiece.Color][side] {
			p.CastlingRights[[]piece.Color{piece.Black, piece.White}[movingPiece.Color]][side] = false
		}
	}
}

func (p *Position) adjustBoard(m *move.Move, from, to square.Square, movingPiece, capturedPiece piece.Piece) {
	// Remove captured piece:
	if capturedPiece.Type != piece.None {
		p.bitBoard[capturedPiece.Color][capturedPiece.Type] ^= (1 << to)
	}

	// Move piece:
	p.bitBoard[movingPiece.Color][movingPiece.Type] ^= ((1 << from) | (1 << to))

	// Castle:
	if movingPiece.Type == piece.King {
		if from == (square.E1+square.Square(56*uint8(movingPiece.Color))) && (to == square.G1+square.Square(56*uint8(movingPiece.Color))) {
			p.bitBoard[movingPiece.Color][piece.Rook] ^= (1 << (square.H1 + square.Square(56*movingPiece.Color))) | (1 << (square.F1 + square.Square(56*movingPiece.Color)))
		} else if from == square.E1+square.Square(56*uint8(movingPiece.Color)) && (to == square.C1+square.Square(56*uint8(movingPiece.Color))) {
			p.bitBoard[movingPiece.Color][piece.Rook] ^= (1 << (square.A1 + square.Square(56*(movingPiece.Color)))) | (1 << (square.D1 + square.Square(56*(movingPiece.Color))))
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
		if m.Promote != piece.None {
			p.bitBoard[movingPiece.Color][movingPiece.Type] ^= (1 << to) // remove piece.Pawn
			p.bitBoard[movingPiece.Color][m.Promote] ^= (1 << to)        // add promoted piece
		}
	}
}

// Put places a piece on the square and removes any other piece
// that may be on that square.
func (p *Position) Put(pp piece.Piece, s square.Square) {
	pc := p.OnSquare(s)
	if pc.Type != piece.None {
		p.bitBoard[pc.Color][pc.Type] ^= (1 << s)
	}
	p.bitBoard[pp.Color][pp.Type] |= (1 << s)
}

// QuickPut places a piece on the square without removing
// any piece that may already be on that square.
func (p *Position) QuickPut(pc piece.Piece, s square.Square) {
	p.bitBoard[pc.Color][pc.Type] |= (1 << s)
}

// Find returns the squares that hold the specified piece.
func (p *Position) Find(pc piece.Piece) map[square.Square]struct{} {
	s := make(map[square.Square]struct{})
	bits := p.bitBoard[pc.Color][pc.Type]
	for bits != 0 {
		sq := bitscan(bits)
		s[square.Square(sq)] = struct{}{}
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
	kingsq := square.Square(bitscan(p.bitBoard[color][piece.King]))
	return p.Threatened(kingsq, opponent)
}
