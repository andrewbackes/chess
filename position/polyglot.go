package position

import (
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/board"
	"github.com/andrewbackes/chess/position/square"
)

// Hash is a polyglot encoding of the given position.
type Hash uint64

// Encode returns the polyglot hash of the current game position. For more
// info you can check out http://hardy.uhasselt.be/Toga/book_format.html
func (p *Position) Polyglot() Hash {
	var hash uint64
	// pieces:
	for s := square.Square(0); s <= square.LastSquare; s++ {
		if pp := p.OnSquare(s); pp.Type != piece.None {
			pc := pieceToPG(pp)
			file, row := indexToFR(int(s))
			index := 64*pc + 8*row + file
			hash ^= randomPiece[index]
		}
	}

	// castles:
	if p.GetCastlingRights()[piece.White][board.ShortSide] {
		hash ^= randomCastle[0]
	}
	if p.GetCastlingRights()[piece.White][board.LongSide] {
		hash ^= randomCastle[1]
	}
	if p.GetCastlingRights()[piece.Black][board.ShortSide] {
		hash ^= randomCastle[2]
	}
	if p.GetCastlingRights()[piece.Black][board.LongSide] {
		hash ^= randomCastle[3]
	}

	// enpassant:
	if p.GetEnPassant() != square.NoSquare {
		file, _ := indexToFR(int(p.GetEnPassant()))
		hashit := false
		rank := []uint{5, 4}[p.GetActiveColor()]
		if file > 0 {
			pp := p.OnSquare(square.New(uint(file), rank))
			if pp.Color == p.GetActiveColor() && pp.Type == piece.Pawn {
				hashit = true
			}
		}
		if file < 7 {
			pp := p.OnSquare(square.New(uint(file+2), rank)) //+2 b/c NewSquare uses 1-8 for file, not 0-7
			if pp.Color == p.GetActiveColor() && pp.Type == piece.Pawn {
				hashit = true
			}
		}
		if hashit {
			hash ^= randomEnpassant[file]
		}
	}

	// turn:
	if p.GetActiveColor() == piece.White {
		hash ^= randomTurn[0]
	}
	return Hash(hash)
}

func indexToFR(index int) (file int, row int) {
	// 0  --> h1 --> 7,0
	// 7  --> a1 --> 0,0 (row,file)
	// 63 --> a8 --> 0,7
	// 56 --> h8 --> 7,7
	row = index / 8
	file = 7 - (index % 8)
	return
}

func pieceToPG(p piece.Piece) int {
	/*
		"kind_of_piece" is encoded as follows
		black pawn    0
		white pawn    1
		black knight  2
		white knight  3
		black bishop  4
		white bishop  5
		black rook    6
		white rook    7
		black queen   8
		white queen   9
		black king   10
		white king   11
	*/
	if p.Color == piece.White {
		return ((int(p.Type) - 1) * 2) + 1
	}
	return ((int(p.Type) - 1) * 2)
}
