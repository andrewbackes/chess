package position

import (
	"errors"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
	"regexp"
	"strings"
)

// ParseMove transforms a move written in standard algebraic notation (SAN)
// to a move written in Pure Coordinate Notation (PCN).
//
// TODO(andrewbackes): ParseMove - What about promotion captures? or ambiguous promotions?
// BUG(andrewbackes): ParseMove - Illegal move: f7g8 (raw: fxg8=Q)
// BUG(andrewbackes): ParseMove - Illegal move: move axb8=Q+
func (p *Position) ParseMove(san string) (move.Move, error) {

	// Check for null move:
	if san == "0000" {
		return move.Parse(san), nil
	}
	color := p.ActiveColor
	// Check for castling:
	if san == "O-O" {
		return move.Parse([]string{"e1g1", "e8g8"}[color]), nil
	}
	if san == "O-O-O" {
		return move.Parse([]string{"e1c1", "e8c8"}[color]), nil
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
				f := move.Parse(parsed).From()
				p := p.OnSquare(f)
				if p.Type == piece.Pawn {
					parsed += "q"
				}
			}
		}
		return move.Parse(parsed), nil
	}

	//	    (piece)    (from)  (from)  (cap) (dest)      (promotion)        (chk  )
	SAN := "([BKNPQR]?)([a-h]?)([0-9]?)([x]?)([a-h][1-8])([=]?[BNPQRbnpqr]?)([+#]?)"
	r, _ := regexp.Compile(SAN)

	matched := r.FindStringSubmatch(san)
	if len(matched) == 0 {
		return move.Parse(san), errors.New("could not parse '" + san + "'")
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

	origin, err := p.originOfPiece(piece, color, destination, fromFile, fromRank)
	if err != nil {
		return move.Parse(san), errors.New("could not find source square of '" + san + "'")
	}

	// Some engines dont tell you to promote to queen, so assume so in that case:
	//if piece == "P" && ((origin[1] == '7' && destination[1] == '8') || (origin[1] == '2' && destination[1] == '1')) {
	//	if promote == "" {
	//		promote = "Q"
	//	}
	//}

	return move.Parse(origin + destination + strings.ToLower(promote)), nil
}

func (p *Position) originOfPiece(pc string, color piece.Color, destination, fromFile, fromRank string) (string, error) {
	pieceMap := map[string]piece.Type{
		"P": piece.Pawn, "p": piece.Pawn,
		"N": piece.Knight, "n": piece.Knight,
		"B": piece.Bishop, "b": piece.Bishop,
		"R": piece.Rook, "r": piece.Rook,
		"Q": piece.Queen, "q": piece.Queen,
		"K": piece.King, "k": piece.King,
	}

	if fromFile != "" && fromRank != "" {
		return fromFile + fromRank, nil
	}

	// Get all legal moves:
	legalMoves := p.LegalMoves()
	var eligableMoves []move.Move

	// Grab the legal moves that land on our square:
	for mv := range legalMoves {
		if mv.To() == square.Parse(destination) {
			eligableMoves = append(eligableMoves, mv)
		}
	}

	// Get all the squares that have our piece on it from the move list:
	var eligableSquares []string
	squares := p.Find(piece.New(color, pieceMap[pc]))
	for sq := range squares {
		for _, mv := range eligableMoves {
			if mv.From() == sq {
				eligableSquares = append(eligableSquares, sq.String())
				break
			}
		}
	}

	// Look for one of the squares that matches the file/rank criteria:
	for _, sq := range eligableSquares {
		if ((sq[0:1] == fromFile) || (fromFile == "")) && ((sq[1:2] == fromRank) || (fromRank == "")) {
			return sq, nil
		}

	}
	//DEBUG:
	//fmt.Println("params: ", piece, destination, fromFile, fromRank)
	//fmt.Println("color: ", color)
	//fmt.Println("legalMoves:", legalMoves)
	//fmt.Println("eligableMoves:", eligableMoves)
	//fmt.Println("eligableSquares:", eligableSquares)

	return "", errors.New("Notation: Can not find source square.")
}
