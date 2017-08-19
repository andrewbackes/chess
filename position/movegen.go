package position

import (
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/board"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
)

// LegalMoves returns only the legal moves that can be made.
func (p *Position) LegalMoves() map[move.Move]struct{} {
	legalMoves := make(map[move.Move]struct{})
	ml := p.Moves()
	for mv := range ml {
		temp := p.MakeMove(mv)
		if temp.Check(p.ActiveColor) == false {
			legalMoves[mv] = struct{}{}
		}
	}
	return legalMoves
}

// Moves returns all moves that a player can make but ignores legality.
// Moves that put the active color into check are included. Castling moves through
// an attacked square are not included.
func (p *Position) Moves() map[move.Move]struct{} {
	moves := make(map[move.Move]struct{})
	add := func(m move.Move) {
		moves[m] = struct{}{}
	}
	notToMove := piece.Color((p.ActiveColor + 1) % 2)
	p.genPawnMoves(p.ActiveColor, notToMove, p.EnPassant, add)
	p.genKnightMoves(p.ActiveColor, notToMove, add)
	p.genDiagnalMoves(p.ActiveColor, notToMove, add)
	p.genStraightMoves(p.ActiveColor, notToMove, add)
	p.genKingMoves(p.ActiveColor, notToMove, p.CastlingRights, add)
	return moves
}

func (p *Position) genKnightMoves(toMove, notToMove piece.Color, add func(move.Move)) {
	//piece.Knights:
	pieces := p.bitBoard[toMove][piece.Knight]
	for pieces != 0 {
		from := bitscan(pieces)
		destinations := knight_moves[from] &^ p.occupied(toMove)
		for destinations != 0 {
			to := bitscan(destinations)
			add(move.Move{Source: square.Square(from), Destination: square.Square(to), Promote: piece.None})
			destinations ^= (1 << to)
		}
		pieces ^= (1 << from)
	}
}

func (p *Position) genDiagnalMoves(toMove, notToMove piece.Color, add func(move.Move)) {
	// piece.Bishops/piece.Queens:
	pieces := p.bitBoard[toMove][piece.Bishop] | p.bitBoard[toMove][piece.Queen]
	direction := [4][65]uint64{ne, nw, se, sw}
	scan := [4]func(uint64) uint{bsf, bsf, bsr, bsr}
	for pieces != 0 {
		from := bitscan(pieces)
		for i := 0; i < 4; i++ {
			destinations := direction[i][from]
			blockerIndex := scan[i](destinations & p.occupied(piece.BothColors))
			destinations ^= direction[i][blockerIndex]
			destinations &^= p.occupied(toMove)
			for destinations != 0 {
				to := bitscan(destinations)
				add(move.Move{Source: square.Square(from), Destination: square.Square(to), Promote: piece.None})
				destinations ^= (1 << to)
			}
		}
		pieces ^= (1 << from)
	}
}

func (p *Position) genStraightMoves(toMove, notToMove piece.Color, add func(move.Move)) {
	// Rooks/piece.Queens:
	pieces := p.bitBoard[toMove][piece.Rook] | p.bitBoard[toMove][piece.Queen]
	direction := [4][65]uint64{north, west, south, east}
	scan := [4]func(uint64) uint{bsf, bsf, bsr, bsr}
	for pieces != 0 {
		from := bitscan(pieces)
		for i := 0; i < 4; i++ {
			destinations := direction[i][from]
			blockerIndex := scan[i](destinations & p.occupied(piece.BothColors))
			destinations ^= direction[i][blockerIndex]
			destinations &^= p.occupied(toMove)
			for destinations != 0 {
				to := bitscan(destinations)
				add(move.Move{Source: square.Square(from), Destination: square.Square(to), Promote: piece.None})
				destinations ^= (1 << to)
			}
		}
		pieces ^= (1 << from)
	}
}

func (p *Position) genKingMoves(toMove, notToMove piece.Color, castlingRights map[piece.Color]map[board.Side]bool, add func(move.Move)) {
	pieces := p.bitBoard[toMove][piece.King]
	{
		from := bitscan(pieces)
		destinations := king_moves[from] &^ p.occupied(toMove)
		for destinations != 0 {
			to := bitscan(destinations)
			add(move.Move{Source: square.Square(from), Destination: square.Square(to), Promote: piece.None})
			destinations ^= (1 << to)
		}
		// Castles:
		if castlingRights[toMove][board.ShortSide] == true {
			if square.Square(bsr(east[from]&p.occupied(piece.BothColors))) == []square.Square{square.H1, square.H8}[toMove] {
				if (p.Threatened([]square.Square{square.F1, square.F8}[toMove], notToMove) == false) &&
					(p.Threatened([]square.Square{square.G1, square.G8}[toMove], notToMove) == false) &&
					(p.Threatened([]square.Square{square.E1, square.E8}[toMove], notToMove) == false) {
					add(move.Move{Source: square.Square(from), Destination: []square.Square{square.G1, square.G8}[toMove], Promote: piece.None})
				}
			}
		}
		if castlingRights[toMove][board.LongSide] == true {
			if square.Square(bsf(west[from]&p.occupied(piece.BothColors))) == []square.Square{square.A1, square.A8}[toMove] {
				if (p.Threatened([]square.Square{square.D1, square.D8}[toMove], notToMove) == false) &&
					(p.Threatened([]square.Square{square.C1, square.C8}[toMove], notToMove) == false) &&
					(p.Threatened([]square.Square{square.E1, square.E8}[toMove], notToMove) == false) {
					add(move.Move{Source: square.Square(from), Destination: []square.Square{square.C1, square.C8}[toMove], Promote: piece.None})
				}
			}
		}
	}
}

func (p *Position) genPawnMoves(toMove, notToMove piece.Color, enPassant square.Square, add func(move.Move)) {
	pieces := p.bitBoard[toMove][piece.Pawn] &^ pawns_spawn[notToMove] //&^ = AND_NOT
	for pieces != 0 {
		from := bitscan(pieces)
		//advances:
		advance := pawn_advances[toMove][from] &^ p.occupied(piece.BothColors)
		if advance != 0 {
			to := bitscan(advance)
			add(move.Move{Source: square.Square(from), Destination: square.Square(to), Promote: piece.None})

			advance = pawn_double_advances[toMove][from] &^ p.occupied(piece.BothColors)
			if advance != 0 {
				to = bitscan(advance)
				add(move.Move{Source: square.Square(from), Destination: square.Square(to), Promote: piece.None})
			}
		}
		//captures:
		var enpas uint64
		if enPassant != square.NoSquare {
			enpas = (1 << enPassant)
		}
		captures := pawn_captures[toMove][from] & (p.occupied(notToMove) | enpas)
		for captures != 0 {
			to := bitscan(captures)
			add(move.Move{Source: square.Square(from), Destination: square.Square(to), Promote: piece.None})
			captures ^= (1 << to)
		}

		pieces ^= (1 << from)
	}
	// Promotions:
	pieces = p.bitBoard[toMove][piece.Pawn] & pawns_spawn[notToMove]
	for pieces != 0 {
		from := bitscan(pieces)
		destinations := pawn_advances[toMove][from] &^ p.occupied(piece.BothColors)
		destinations |= pawn_captures[toMove][from] & p.occupied(notToMove)
		for destinations != 0 {
			to := bitscan(destinations)
			p := []piece.Type{piece.Queen, piece.Rook, piece.Bishop, piece.Knight}
			for i := 0; i < 4; i++ {
				add(move.Move{Source: square.Square(from), Destination: square.Square(to), Promote: p[i]})
			}
			destinations ^= (1 << to)
		}
		pieces ^= (1 << from)
	}
}

// Threatened returns whether or not the specified square is under attack
// by the specified color.
func (p *Position) Threatened(square square.Square, byWho piece.Color) bool {
	defender := []piece.Color{piece.Black, piece.White}[byWho]

	// other king attacks:
	if (king_moves[square] & p.bitBoard[byWho][piece.King]) != 0 {
		return true
	}

	// pawn attacks:
	if pawn_captures[defender][square]&p.bitBoard[byWho][piece.Pawn] != 0 {
		return true
	}

	// knight attacks:
	if knight_moves[square]&p.bitBoard[byWho][piece.Knight] != 0 {
		return true
	}
	// diagonal attacks:
	direction := [4][65]uint64{nw, ne, sw, se}
	scan := [4]func(uint64) uint{bsf, bsf, bsr, bsr}
	for i := 0; i < 4; i++ {
		blockerIndex := scan[i](direction[i][square] & p.occupied(piece.BothColors))
		if (1<<blockerIndex)&(p.bitBoard[byWho][piece.Bishop]|p.bitBoard[byWho][piece.Queen]) != 0 {
			return true
		}
	}
	// straight attacks:
	direction = [4][65]uint64{north, west, south, east}
	for i := 0; i < 4; i++ {
		blockerIndex := scan[i](direction[i][square] & p.occupied(piece.BothColors))
		if (1<<blockerIndex)&(p.bitBoard[byWho][piece.Rook]|p.bitBoard[byWho][piece.Queen]) != 0 {
			return true
		}
	}
	return false
}
