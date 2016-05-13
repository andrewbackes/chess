// Package game plays chess games.
package game

import (
	"fmt"
	"strconv"
	"strings"
)

func NewMove(from, to Square) Move {
	return Move(getAlg(from) + getAlg(to))
}

// Piece represents a chess piece.
type Piece struct {
	Color Color
	Type  PieceType
}

func (P Piece) String() string {
	if P.Type == None {
		return " "
	}
	abbrev := [2][6]string{{"P", "N", "B", "R", "Q", "K"}, {"p", "n", "b", "r", "q", "k"}}
	return abbrev[P.Color][P.Type]
}

// NewPiece returns a new chess piece type.
func NewPiece(c Color, t PieceType) Piece {
	return Piece{
		Color: c,
		Type:  t,
	}
}

// Board is a chess Board.
type Board struct {
	// BitBoard has one bitBoard per player per color.
	BitBoard [2][6]uint64 //[player][piece]
}

func NewBoard() Board {
	b := Board{BitBoard: [2][6]uint64{}}
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
		for c := range b.BitBoard {
			for j := range b.BitBoard[c] {
				if ((1 << square) & b.BitBoard[c][j]) != 0 {
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
	b.BitBoard = [2][6]uint64{}
}

// Reset puts the pieces in the new game position.
func (b *Board) Reset() {
	// puts the pieces in their starting/newgame positions
	for color := uint(0); color < 2; color = color + 1 {
		//Pawns first:
		b.BitBoard[color][Pawn] = 255 << (8 + (color * 8 * 5))
		//Then the rest of the pieces:
		b.BitBoard[color][Knight] = (1 << (B1 + (color * 8 * 7))) ^ (1 << (G1 + (color * 8 * 7)))
		b.BitBoard[color][Bishop] = (1 << (C1 + (color * 8 * 7))) ^ (1 << (F1 + (color * 8 * 7)))
		b.BitBoard[color][Rook] = (1 << (A1 + (color * 8 * 7))) ^ (1 << (H1 + (color * 8 * 7)))
		b.BitBoard[color][Queen] = (1 << (D1 + (color * 8 * 7)))
		b.BitBoard[color][King] = (1 << (E1 + (color * 8 * 7)))
	}
}

// OnSquare returns the piece that is on the specified square.
func (b *Board) OnSquare(s Square) Piece {
	for c := White; c <= Black; c++ {
		for p := Pawn; p <= King; p++ {
			if (b.BitBoard[c][p] & (1 << s)) != 0 {
				return NewPiece(c, p)
			}
		}
	}
	return NewPiece(Neither, None)
}

// Occupied returns a bitBoard with all of the specified colors pieces.
func (b *Board) Occupied(c Color) uint64 {
	var mask uint64
	for p := Pawn; p <= King; p++ {
		if c == Both {
			mask |= b.BitBoard[White][p] | b.BitBoard[Black][p]
		} else {
			mask |= b.BitBoard[c][p]
		}
	}
	return mask
}

func (b *Board) MakeMove(m Move) {
	from, to := getSquares(m)
	movingPiece := b.OnSquare(from)
	capturedPiece := b.OnSquare(to)

	// Remove captured piece:
	if capturedPiece.Type != None {
		b.BitBoard[capturedPiece.Color][capturedPiece.Type] ^= (1 << to)
	}

	// Move piece:
	b.BitBoard[movingPiece.Color][movingPiece.Type] ^= ((1 << from) | (1 << to))

	// Castle:
	if movingPiece.Type == King {
		if from == Square(E1+56*uint8(movingPiece.Color)) && (to == Square(G1+56*uint8(movingPiece.Color))) {
			b.BitBoard[movingPiece.Color][Rook] ^= (1 << (H1 + 56*uint8(movingPiece.Color))) | (1 << (F1 + 56*uint8(movingPiece.Color)))
		} else if from == Square(E1+56*uint8(movingPiece.Color)) && to == Square(C1+56*uint8(movingPiece.Color)) {
			b.BitBoard[movingPiece.Color][Rook] ^= (1 << (A1 + 56*uint8(movingPiece.Color))) | (1 << (D1 + 56*uint8(movingPiece.Color)))
		}
	}

	if movingPiece.Type == Pawn {
		// Handle en Passant capture:
		// capturedPiece just means the piece on the destination square
		if (int(to)-int(from))%8 != 0 && capturedPiece.Type == None {
			if movingPiece.Color == White {
				b.BitBoard[Black][Pawn] ^= (1 << (to - 8))
			} else if movingPiece.Color == Black {
				b.BitBoard[White][Pawn] ^= (1 << (to + 8))
			}
		}
		// Handle Promotions:
		promotesTo := promotedPiece(m)
		if promotesTo != NoPiece {
			b.BitBoard[movingPiece.Color][movingPiece.Type] ^= (1 << to) // remove Pawn
			b.BitBoard[movingPiece.Color][promotesTo] ^= (1 << to)       // add promoted piece
		}
	}

}

func parseBoard(position string) *Board {
	b := NewBoard()
	b.Clear()
	// remove the /'s and replace the numbers with that many spaces
	// so that there is a 1-1 mapping from bytes to squares.
	parsedBoard := strings.Replace(position, "/", "", 9)
	for i := 1; i < 9; i++ {
		parsedBoard = strings.Replace(parsedBoard, strconv.Itoa(i), strings.Repeat(" ", i), -1)
	}
	piece := map[string]PieceType{
		"P": Pawn, "p": Pawn,
		"N": Knight, "n": Knight,
		"B": Bishop, "b": Bishop,
		"R": Rook, "r": Rook,
		"Q": Queen, "q": Queen,
		"K": King, "k": King}
	color := map[string]Color{
		"P": White, "p": Black,
		"N": White, "n": Black,
		"B": White, "b": Black,
		"R": White, "r": Black,
		"Q": White, "q": Black,
		"K": White, "k": Black}
	// adjust the bitboards:
	for pos := 0; pos < len(parsedBoard); pos++ {
		k := parsedBoard[pos:(pos + 1)]
		if _, ok := piece[k]; ok {
			b.BitBoard[color[k]][piece[k]] |= (1 << uint(63-pos))
		}
	}
	return &b
}

func (b *Board) put(p Piece, s Square) {
	b.BitBoard[p.Color][p.Type] |= (1 << s)
}

func (b *Board) printBitBoards() {
	for c := range b.BitBoard {
		for j := range b.BitBoard[c] {
			fmt.Println(NewPiece(Color(c), PieceType(j)))
			bitprint(b.BitBoard[c][j])
		}
	}
}
