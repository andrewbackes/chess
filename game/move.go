package game

import (
	"errors"
	"regexp"
	"strings"
)

// NewMove returns a Move object based on the given from & to squares.
func NewMove(from, to Square) Move {
	return Move(getAlg(from) + getAlg(to))
}

// ParseMove takes a move written in standard algebraic notation (SAN)
// to a Pure Coordinate Notation (PCN) Move.
//
// TODO(andrewbackes): ParseMove - What about promotion captures? or ambiguous promotions?
// BUG(andrewbackes): ParseMove - Illegal move: f7g8 (raw: fxg8=Q)
// BUG(andrewbackes): ParseMove - Illegal move: move axb8=Q+
func (G *Game) ParseMove(san string) (Move, error) {

	// Check for null move:
	if san == "0000" {
		return Move(san), nil
	}
	color := G.PlayerToMove()
	// Check for castling:
	if san == "O-O" {
		return Move([]string{"e1g1", "e8g8"}[color]), nil
	}
	if san == "O-O-O" {
		return Move([]string{"e1c1", "e8c8"}[color]), nil
	}

	// Strip uneeded characters:
	san = strings.Replace(san, "-", "", -1)

	// First check to see if it is already in the correct form.
	PCN := "([a-h][1-8])([a-h][1-8])([QBNRqbnr]?)"
	matches, _ := regexp.MatchString(PCN, san)
	if matches {
		parsed := san[:len(san)-1]
		// Some engines dont capitalize the promotion piece:
		parsed += strings.ToLower(san[len(san)-1:])
		// some engines dont specify the promotion piece, assume queen:
		if (parsed[1] == '7' && parsed[3] == '8') || (parsed[1] == '2' && parsed[3] == '1') {
			if len(parsed) <= 4 {
				f, _ := getSquares(Move(parsed))
				p := G.board.OnSquare(f)
				if p.Type == Pawn {
					parsed += "q"
				}
			}
		}
		return Move(parsed), nil
	}

	//	    (piece)    (from)  (from)  (cap) (dest)      (promotion)        (chk  )
	SAN := "([BKNPQR]?)([a-h]?)([0-9]?)([x]?)([a-h][1-8])([=]?[BNPQRbnpqr]?)([+#]?)"
	r, _ := regexp.Compile(SAN)

	matched := r.FindStringSubmatch(san)
	if len(matched) == 0 {
		return Move(san), errors.New("could not parse '" + san + "'")
	}

	piece := matched[1]
	fromFile := matched[2]
	fromRank := matched[3]
	//action := matched[4]      // capture or promote
	destination := matched[5] //or promotion piece if action="="
	//check := matched[6]       //or mate
	promote := strings.Replace(matched[6], "=", "", 1)

	if piece == "" {
		piece = "P"
	}

	origin, err := G.originOfPiece(piece, destination, fromFile, fromRank)
	if err != nil {
		//fmt.Println(err)
		//fmt.Println(G.FEN())
		//fmt.Println(san)
		//G.PrintHUD()
		return Move(san), errors.New("could not find source square of '" + san + "'")
	}

	// Some engines dont tell you to promote to queen, so assume so in that case:
	/*if piece == "P" && ((origin[1] == '7' && destination[1] == '8') || (origin[1] == '2' && destination[1] == '1')) {
		if promote == "" {
			promote = "Q"
		}
	}
	*/
	return Move(origin + destination + strings.ToLower(promote)), nil
}

func (G *Game) originOfPiece(piece, destination, fromFile, fromRank string) (string, error) {
	pieceMap := map[string]PieceType{
		"P": Pawn, "p": Pawn,
		"N": Knight, "n": Knight,
		"B": Bishop, "b": Bishop,
		"R": Rook, "r": Rook,
		"Q": Queen, "q": Queen,
		"K": King, "k": King}

	if fromFile != "" && fromRank != "" {
		return fromFile + fromRank, nil
	}

	// Get all legal moves:
	legalMoves := G.LegalMoves()
	var eligableMoves []Move

	// Grab the legal moves that land on our square:
	for mv := range legalMoves {
		dest := mv[2:4]
		if string(dest) == destination {
			eligableMoves = append(eligableMoves, mv)
		}
	}

	// Get all the squares that have our piece on it from the move list:
	color := G.PlayerToMove()
	var eligableSquares []string
	bits := G.board.bitBoard[color][pieceMap[piece]]
	for bits != 0 {
		bit := bitscan(bits)
		sq := getAlg(Square(bit))
		//verify that its a legal move:
		for _, mv := range eligableMoves {
			if string(mv[:2]) == sq {
				eligableSquares = append(eligableSquares, sq)
				break
			}
		}
		bits ^= (1 << bit)
	}

	// Look for one of the squares that matches the file/rank criteria:
	for _, sq := range eligableSquares {
		if ((sq[0:1] == fromFile) || (fromFile == "")) && ((sq[1:2] == fromRank) || (fromRank == "")) {
			return sq, nil
		}

	}
	//DEBUG:
	/*
		fmt.Println("params: ", piece, destination, fromFile, fromRank)
		fmt.Println("color: ", color)
		fmt.Println("legalMoves:", legalMoves)
		fmt.Println("eligableMoves:", eligableMoves)
		fmt.Println("eligableSquares:", eligableSquares)
	*/
	return "", errors.New("Notation: Can not find source square.")
}
