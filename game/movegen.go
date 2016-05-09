package game

func (G *Game) LegalMoves() map[Move]struct{} {
	legalMoves := make(map[Move]struct{})
	ml := G.PsuedoLegalMoves()
	toMove := G.PlayerToMove()
	for mv := range ml {
		temp := *G
		temp.MakeMove(mv)
		if temp.isInCheck(toMove) == false {
			legalMoves[mv] = struct{}{}
		}
	}
	return legalMoves
}

func (G *Game) isInCheck(toMove Color) bool {
	notToMove := []Color{Black, White}[toMove]
	kingsq := bitscan(G.board.BitBoard[toMove][King])
	return G.isAttacked(kingsq, notToMove)
}

func (G *Game) PsuedoLegalMoves() map[Move]struct{} {
	moves := make(map[Move]struct{})
	add := func(m Move) {
		moves[m] = struct{}{}
	}
	toMove := G.PlayerToMove()
	notToMove := []Color{Black, White}[toMove] // too bad !toMove doesnt work =(
	G.genPawnMoves(toMove, notToMove, add)
	G.genKnightMoves(toMove, notToMove, add)
	G.genDiagnalMoves(toMove, notToMove, add)
	G.genStraightMoves(toMove, notToMove, add)
	G.genKingMoves(toMove, notToMove, add)
	return moves
}

func (G *Game) genKnightMoves(toMove, notToMove Color, add func(Move)) {
	//Knights:
	pieces := G.board.BitBoard[toMove][Knight]
	for pieces != 0 {
		from := bitscan(pieces)
		destinations := knight_moves[from] &^ G.board.Occupied(toMove)
		for destinations != 0 {
			to := bitscan(destinations)
			add(NewMove(Square(from), Square(to)))
			destinations ^= (1 << to)
		}
		pieces ^= (1 << from)
	}
}

func (G *Game) genDiagnalMoves(toMove, notToMove Color, add func(Move)) {
	// Bishops/Queens:
	pieces := G.board.BitBoard[toMove][Bishop] | G.board.BitBoard[toMove][Queen]
	direction := [4][65]uint64{ne, nw, se, sw}
	scan := [4]func(uint64) uint{BSF, BSF, BSR, BSR}
	for pieces != 0 {
		from := bitscan(pieces)
		for i := 0; i < 4; i++ {
			destinations := direction[i][from]
			blockerIndex := scan[i](destinations & G.board.Occupied(Both))
			destinations ^= direction[i][blockerIndex]
			destinations &^= G.board.Occupied(toMove)
			for destinations != 0 {
				to := bitscan(destinations)
				add(NewMove(Square(from), Square(to)))
				destinations ^= (1 << to)
			}
		}
		pieces ^= (1 << from)
	}
}

func (G *Game) genStraightMoves(toMove, notToMove Color, add func(Move)) {
	// Rooks/Queens:
	pieces := G.board.BitBoard[toMove][Rook] | G.board.BitBoard[toMove][Queen]
	direction := [4][65]uint64{north, west, south, east}
	scan := [4]func(uint64) uint{BSF, BSF, BSR, BSR}
	for pieces != 0 {
		from := bitscan(pieces)
		for i := 0; i < 4; i++ {
			destinations := direction[i][from]
			blockerIndex := scan[i](destinations & G.board.Occupied(Both))
			destinations ^= direction[i][blockerIndex]
			destinations &^= G.board.Occupied(toMove)
			for destinations != 0 {
				to := bitscan(destinations)
				add(NewMove(Square(from), Square(to)))
				destinations ^= (1 << to)
			}
		}
		pieces ^= (1 << from)
	}
}

func (G *Game) genKingMoves(toMove, notToMove Color, add func(Move)) {
	pieces := G.board.BitBoard[toMove][King]
	{
		from := bitscan(pieces)
		destinations := king_moves[from] &^ G.board.Occupied(toMove)
		for destinations != 0 {
			to := bitscan(destinations)
			add(NewMove(Square(from), Square(to)))
			destinations ^= (1 << to)
		}
		// Castles:
		if G.history.castlingRights[toMove][ShortSide] == true {
			if BSR(east[from]&G.board.Occupied(Both)) == []uint{H1, H8}[toMove] {
				if (G.isAttacked([]uint{F1, F8}[toMove], notToMove) == false) &&
					(G.isAttacked([]uint{G1, G8}[toMove], notToMove) == false) &&
					(G.isAttacked([]uint{E1, E8}[toMove], notToMove) == false) {
					add(NewMove(Square(from), Square([]uint{G1, G8}[toMove])))
				}
			}
		}
		if G.history.castlingRights[toMove][LongSide] == true {
			if BSF(west[from]&G.board.Occupied(Both)) == []uint{A1, A8}[toMove] {
				if (G.isAttacked([]uint{D1, D8}[toMove], notToMove) == false) &&
					(G.isAttacked([]uint{C1, C8}[toMove], notToMove) == false) &&
					(G.isAttacked([]uint{E1, E8}[toMove], notToMove) == false) {
					add(NewMove(Square(from), Square([]uint{C1, C8}[toMove])))
				}
			}
		}
	}
}

func (G *Game) genPawnMoves(toMove, notToMove Color, add func(Move)) {
	pieces := G.board.BitBoard[toMove][Pawn] &^ pawns_spawn[notToMove] //&^ = AND_NOT
	for pieces != 0 {
		from := bitscan(pieces)

		//advances:
		advance := pawn_advances[toMove][from] &^ G.board.Occupied(BothColors)
		if advance != 0 {
			to := bitscan(advance)
			add(NewMove(Square(from), Square(to)))

			advance = pawn_double_advances[toMove][from] &^ G.board.Occupied(BothColors)
			if advance != 0 {
				to = bitscan(advance)
				add(NewMove(Square(from), Square(to)))
			}
		}

		//captures:
		var enpas uint64
		if G.history.enPassant != nil {
			enpas = (1 << *G.history.enPassant)
		}
		captures := pawn_captures[toMove][from] & (G.board.Occupied(notToMove) | enpas)
		for captures != 0 {
			to := bitscan(captures)
			add(NewMove(Square(from), Square(to)))
			captures ^= (1 << to)
		}

		pieces ^= (1 << from)
	}
	// Promotions:
	pieces = G.board.BitBoard[toMove][Pawn] & pawns_spawn[notToMove]
	for pieces != 0 {
		from := bitscan(pieces)
		destinations := pawn_advances[toMove][from] &^ G.board.Occupied(Both)
		destinations |= pawn_captures[toMove][from] & G.board.Occupied(notToMove)
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

func (G *Game) isAttacked(square uint, byWho Color) bool {
	defender := []Color{Black, White}[byWho]

	// other king attacks:
	if (king_moves[square] & G.board.BitBoard[byWho][King]) != 0 {
		return true
	}

	// pawn attacks:
	if pawn_captures[defender][square]&G.board.BitBoard[byWho][Pawn] != 0 {
		return true
	}

	// knight attacks:
	if knight_moves[square]&G.board.BitBoard[byWho][Knight] != 0 {
		return true
	}
	// diagonal attacks:
	direction := [4][65]uint64{nw, ne, sw, se}
	scan := [4]func(uint64) uint{BSF, BSF, BSR, BSR}
	for i := 0; i < 4; i++ {
		blockerIndex := scan[i](direction[i][square] & G.board.Occupied(Both))
		if (1<<blockerIndex)&(G.board.BitBoard[byWho][Bishop]|G.board.BitBoard[byWho][Queen]) != 0 {
			return true
		}
	}
	// straight attacks:
	direction = [4][65]uint64{north, west, south, east}
	for i := 0; i < 4; i++ {
		blockerIndex := scan[i](direction[i][square] & G.board.Occupied(Both))
		if (1<<blockerIndex)&(G.board.BitBoard[byWho][Rook]|G.board.BitBoard[byWho][Queen]) != 0 {
			return true
		}
	}
	return false
}
