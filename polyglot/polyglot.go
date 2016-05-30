package polyglot

import (
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position"
)

// Encode returns the polyglot hash of the current game position. For more
// info you can check out http://hardy.uhasselt.be/Toga/book_format.html
func Encode(p *position.Position) (hash uint64) {
	// pieces:
	for s := position.Square(0); s <= position.LastSquare; s++ {
		if pp := p.OnSquare(s); pp.Type != piece.None {
			pc := pieceToPG(pp)
			file, row := indexToFR(int(s))
			index := 64*pc + 8*row + file
			hash ^= randomPiece[index]
		}
	}

	// castles:
	if p.CastlingRights[piece.White][position.ShortSide] {
		hash ^= randomCastle[0]
	}
	if p.CastlingRights[piece.White][position.LongSide] {
		hash ^= randomCastle[1]
	}
	if p.CastlingRights[piece.Black][position.ShortSide] {
		hash ^= randomCastle[2]
	}
	if p.CastlingRights[piece.Black][position.LongSide] {
		hash ^= randomCastle[3]
	}

	// enpassant:
	if p.EnPassant != position.NoSquare {
		file, _ := indexToFR(int(p.EnPassant))
		hashit := false
		rank := []uint{5, 4}[p.ActiveColor]
		if file > 0 {
			pp := p.OnSquare(position.NewSquare(uint(file), rank))
			if pp.Color == p.ActiveColor && pp.Type == piece.Pawn {
				hashit = true
			}
		}
		if file < 7 {
			pp := p.OnSquare(position.NewSquare(uint(file+2), rank)) //+2 b/c NewSquare uses 1-8 for file, not 0-7
			if pp.Color == p.ActiveColor && pp.Type == piece.Pawn {
				hashit = true
			}
		}
		if hashit {
			hash ^= randomEnpassant[file]
		}
	}

	// turn:
	if p.ActiveColor == piece.White {
		hash ^= randomTurn[0]
	}
	return
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
		return (int(p.Type) * 2) + 1
	}
	return (int(p.Type) * 2)
}
