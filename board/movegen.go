package board

import (
	"github.com/andrewbackes/chess/piece"
)

// Moves returns all moves that a player can make but ignores legality.
// Moves that put the active color into check are included. Castling moves through
// an attacked square are not included.
func (b *Board) Moves(c piece.Color, enPassant *Square, castlingRights [2][2]bool) map[Move]struct{} {
	moves := make(map[Move]struct{})
	add := func(m Move) {
		moves[m] = struct{}{}
	}
	notToMove := piece.Color((c + 1) % 2)
	b.genPawnMoves(c, notToMove, enPassant, add)
	b.genKnightMoves(c, notToMove, add)
	b.genDiagnalMoves(c, notToMove, add)
	b.genStraightMoves(c, notToMove, add)
	b.genKingMoves(c, notToMove, castlingRights, add)
	return moves
}

func (b *Board) genKnightMoves(toMove, notToMove piece.Color, add func(Move)) {
	//piece.Knights:
	pieces := b.bitBoard[toMove][piece.Knight]
	for pieces != 0 {
		from := bitscan(pieces)
		destinations := knight_moves[from] &^ b.occupied(toMove)
		for destinations != 0 {
			to := bitscan(destinations)
			add(NewMove(Square(from), Square(to)))
			destinations ^= (1 << to)
		}
		pieces ^= (1 << from)
	}
}

func (b *Board) genDiagnalMoves(toMove, notToMove piece.Color, add func(Move)) {
	// piece.Bishops/piece.Queens:
	pieces := b.bitBoard[toMove][piece.Bishop] | b.bitBoard[toMove][piece.Queen]
	direction := [4][65]uint64{ne, nw, se, sw}
	scan := [4]func(uint64) uint{bsf, bsf, bsr, bsr}
	for pieces != 0 {
		from := bitscan(pieces)
		for i := 0; i < 4; i++ {
			destinations := direction[i][from]
			blockerIndex := scan[i](destinations & b.occupied(piece.BothColors))
			destinations ^= direction[i][blockerIndex]
			destinations &^= b.occupied(toMove)
			for destinations != 0 {
				to := bitscan(destinations)
				add(NewMove(Square(from), Square(to)))
				destinations ^= (1 << to)
			}
		}
		pieces ^= (1 << from)
	}
}

func (b *Board) genStraightMoves(toMove, notToMove piece.Color, add func(Move)) {
	// Rooks/piece.Queens:
	pieces := b.bitBoard[toMove][piece.Rook] | b.bitBoard[toMove][piece.Queen]
	direction := [4][65]uint64{north, west, south, east}
	scan := [4]func(uint64) uint{bsf, bsf, bsr, bsr}
	for pieces != 0 {
		from := bitscan(pieces)
		for i := 0; i < 4; i++ {
			destinations := direction[i][from]
			blockerIndex := scan[i](destinations & b.occupied(piece.BothColors))
			destinations ^= direction[i][blockerIndex]
			destinations &^= b.occupied(toMove)
			for destinations != 0 {
				to := bitscan(destinations)
				add(NewMove(Square(from), Square(to)))
				destinations ^= (1 << to)
			}
		}
		pieces ^= (1 << from)
	}
}

func (b *Board) genKingMoves(toMove, notToMove piece.Color, castlingRights [2][2]bool, add func(Move)) {
	pieces := b.bitBoard[toMove][piece.King]
	{
		from := bitscan(pieces)
		destinations := king_moves[from] &^ b.occupied(toMove)
		for destinations != 0 {
			to := bitscan(destinations)
			add(NewMove(Square(from), Square(to)))
			destinations ^= (1 << to)
		}
		// Castles:
		if castlingRights[toMove][shortSide] == true {
			if Square(bsr(east[from]&b.occupied(piece.BothColors))) == []Square{H1, H8}[toMove] {
				if (b.isAttacked([]Square{F1, F8}[toMove], notToMove) == false) &&
					(b.isAttacked([]Square{G1, G8}[toMove], notToMove) == false) &&
					(b.isAttacked([]Square{E1, E8}[toMove], notToMove) == false) {
					add(NewMove(Square(from), []Square{G1, G8}[toMove]))
				}
			}
		}
		if castlingRights[toMove][longSide] == true {
			if Square(bsf(west[from]&b.occupied(piece.BothColors))) == []Square{A1, A8}[toMove] {
				if (b.isAttacked([]Square{D1, D8}[toMove], notToMove) == false) &&
					(b.isAttacked([]Square{C1, C8}[toMove], notToMove) == false) &&
					(b.isAttacked([]Square{E1, E8}[toMove], notToMove) == false) {
					add(NewMove(Square(from), []Square{C1, C8}[toMove]))
				}
			}
		}
	}
}

func (b *Board) genPawnMoves(toMove, notToMove piece.Color, enPassant *Square, add func(Move)) {
	pieces := b.bitBoard[toMove][piece.Pawn] &^ pawns_spawn[notToMove] //&^ = AND_NOT
	for pieces != 0 {
		from := bitscan(pieces)
		//advances:
		advance := pawn_advances[toMove][from] &^ b.occupied(piece.BothColors)
		if advance != 0 {
			to := bitscan(advance)
			add(NewMove(Square(from), Square(to)))

			advance = pawn_double_advances[toMove][from] &^ b.occupied(piece.BothColors)
			if advance != 0 {
				to = bitscan(advance)
				add(NewMove(Square(from), Square(to)))
			}
		}
		//captures:
		var enpas uint64
		if enPassant != nil {
			enpas = (1 << *enPassant)
		}
		captures := pawn_captures[toMove][from] & (b.occupied(notToMove) | enpas)
		for captures != 0 {
			to := bitscan(captures)
			add(NewMove(Square(from), Square(to)))
			captures ^= (1 << to)
		}

		pieces ^= (1 << from)
	}
	// Promotions:
	pieces = b.bitBoard[toMove][piece.Pawn] & pawns_spawn[notToMove]
	for pieces != 0 {
		from := bitscan(pieces)
		destinations := pawn_advances[toMove][from] &^ b.occupied(piece.BothColors)
		destinations |= pawn_captures[toMove][from] & b.occupied(notToMove)
		for destinations != 0 {
			to := bitscan(destinations)
			p := []string{"q", "r", "b", "n"}
			for i := 0; i < 4; i++ {
				mv := string(NewMove(Square(from), Square(to)))
				mv += p[i]
				add(Move(mv))
			}
			destinations ^= (1 << to)
		}
		pieces ^= (1 << from)
	}
}

func (b *Board) isAttacked(square Square, byWho piece.Color) bool {
	defender := []piece.Color{piece.Black, piece.White}[byWho]

	// other king attacks:
	if (king_moves[square] & b.bitBoard[byWho][piece.King]) != 0 {
		return true
	}

	// pawn attacks:
	if pawn_captures[defender][square]&b.bitBoard[byWho][piece.Pawn] != 0 {
		return true
	}

	// knight attacks:
	if knight_moves[square]&b.bitBoard[byWho][piece.Knight] != 0 {
		return true
	}
	// diagonal attacks:
	direction := [4][65]uint64{nw, ne, sw, se}
	scan := [4]func(uint64) uint{bsf, bsf, bsr, bsr}
	for i := 0; i < 4; i++ {
		blockerIndex := scan[i](direction[i][square] & b.occupied(piece.BothColors))
		if (1<<blockerIndex)&(b.bitBoard[byWho][piece.Bishop]|b.bitBoard[byWho][piece.Queen]) != 0 {
			return true
		}
	}
	// straight attacks:
	direction = [4][65]uint64{north, west, south, east}
	for i := 0; i < 4; i++ {
		blockerIndex := scan[i](direction[i][square] & b.occupied(piece.BothColors))
		if (1<<blockerIndex)&(b.bitBoard[byWho][piece.Rook]|b.bitBoard[byWho][piece.Queen]) != 0 {
			return true
		}
	}
	return false
}
