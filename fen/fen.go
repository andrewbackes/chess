package chess

import (
	"errors"
	"github.com/andrewbackes/chess/board"
	"github.com/andrewbackes/chess/piece"
	"strconv"
	"strings"
)

// Encode returns the fen of the current position of the game.
func Encode(G *Game) string {

	pc := [][]string{
		{"P", "N", "B", "R", "Q", "K", " "},
		{"p", "n", "b", "r", "q", "k", " "},
		{" ", " ", " ", " ", " ", " ", " "}}

	var boardstr string
	// put what is on each square into a squence (including blanks):
	for i := int(63); i >= 0; i-- {
		p := G.board.OnSquare(board.Square(i))
		boardstr += pc[p.Color][p.Type]
		if i%8 == 0 && i > 0 {
			boardstr += "/"
		}
	}
	// replace groups of spaces with numbers instead
	for i := 8; i > 0; i-- {
		boardstr = strings.Replace(boardstr, strings.Repeat(" ", i), strconv.Itoa(i), -1)
	}
	// Player to move:
	turn := []string{"w", "b"}[G.ActiveColor()]
	// Castling Rights:
	var rights string
	castles := [][]string{{"K", "Q"}, {"k", "q"}}
	for c := piece.White; c <= piece.Black; c++ {
		for side := board.ShortSide; side <= board.LongSide; side++ {
			if G.history.castlingRights[c][side] {
				rights += castles[c][side]
			}
		}
	}
	if rights == "" {
		rights = "-"
	}
	// en Passant:
	var enPas string
	if G.history.enPassant != nil {
		enPas = board.Square(*G.history.enPassant).String()
	} else {
		enPas = "-"
	}
	// Moves and 50 move rule
	fifty := strconv.Itoa(int(G.history.fiftyMoveCount / 2))
	move := strconv.Itoa(int(len(G.history.move)/2) + 1)
	// all together:
	fen := boardstr + " " + turn + " " + rights + " " + enPas + " " + fifty + " " + move
	return fen
}

// Decode creates a game from the provided FEN.
func Decode(fen string) (*Game, error) {
	G := NewGame()
	words := strings.Split(fen, " ")
	if len(words) < 4 {
		return nil, errors.New("FEN: incomplete fen")
	}
	if words[1] != "w" && words[1] != "b" {
		return nil, errors.New("FEN: can not determine active player")
	}
	b, err := board.FromFEN(words[0])
	if err != nil {
		return nil, err
	}
	G.board = *b
	if len(words) >= 6 {
		h, _ := parseMoveHistory(words[1], words[5], words[4])
		G.history = *h
	} else if strings.ToLower(words[1]) == "b" {
		// add a null move since we want it to be black's turn.
		var m []board.Move
		m = append(m, board.NullMove)
		G.history.move = m
	}
	G.history.castlingRights = parseCastlingRights(words[2])
	G.history.enPassant = parseEnPassantSquare(words[3])
	G.history.startingFen = fen

	return G, nil
}

func parseMoveHistory(activeColor, moveCount, fiftyMoveCount string) (*gameHistory, error) {
	h := gameHistory{}
	fullMoves, err := strconv.ParseUint(moveCount, 10, 0)
	if err != nil {
		return nil, errors.New("FEN: could not parse move count")
	}
	halfMoves := ((fullMoves - 1) * 2) + map[string]uint64{"w": 0, "b": 1}[activeColor]
	for i := uint64(0); i < halfMoves; i++ {
		h.move = append(h.move, board.NullMove)
	}
	fmc, err := strconv.ParseUint(fiftyMoveCount, 10, 0)
	if err != nil {
		return nil, errors.New("FEN: could not parse fifty move rule count")
	}
	// Since internally we store half moves:
	h.fiftyMoveCount = (fmc * 2) + map[string]uint64{"w": 0, "b": 1}[activeColor]
	return &h, nil
}

func parseEnPassantSquare(sq string) *board.Square {
	if sq != "-" {
		s := board.ParseSquare(sq)
		return &s
	}
	return nil
}

func parseCastlingRights(KQkq string) [2][2]bool {
	return [2][2]bool{
		{strings.Contains(KQkq, "K"), strings.Contains(KQkq, "Q")},
		{strings.Contains(KQkq, "k"), strings.Contains(KQkq, "q")}}
}

/*
// GameFromFEN parses the board passed via FEN and returns a board object.
func FromFEN(position string) (*Board, error) {
	b := New()
	b.Clear()
	// remove the /'s and replace the numbers with that many spaces
	// so that there is a 1-1 mapping from bytes to squares.
	justBoard := strings.Split(position, " ")[0]
	parsedBoard := strings.Replace(justBoard, "/", "", 9)
	for i := 1; i < 9; i++ {
		parsedBoard = strings.Replace(parsedBoard, strconv.Itoa(i), strings.Repeat(" ", i), -1)
	}
	if len(parsedBoard) < 64 {
		return nil, errors.New("fen: could not parse position")
	}
	p := map[rune]piece.Type{
		'P': piece.Pawn, 'p': piece.Pawn,
		'N': piece.Knight, 'n': piece.Knight,
		'B': piece.Bishop, 'b': piece.Bishop,
		'R': piece.Rook, 'r': piece.Rook,
		'Q': piece.Queen, 'q': piece.Queen,
		'K': piece.King, 'k': piece.King}
	color := map[rune]piece.Color{
		'P': piece.White, 'p': piece.Black,
		'N': piece.White, 'n': piece.Black,
		'B': piece.White, 'b': piece.Black,
		'R': piece.White, 'r': piece.Black,
		'Q': piece.White, 'q': piece.Black,
		'K': piece.White, 'k': piece.Black}
	// adjust the bitboards:
	for pos := 0; pos < len(parsedBoard); pos++ {
		if pos > 64 {
			break
		}
		k := rune(parsedBoard[pos])
		if _, ok := p[k]; ok {
			b.bitBoard[color[k]][p[k]] |= (1 << uint(63-pos))
		}
	}
	return &b, nil
}
*/
