// Package fen is for encoding/decoding chess positions into fen format.
// You can also decode straight into a playable game.
package fen

import (
	"errors"
	"fmt"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position"
	"github.com/andrewbackes/chess/position/board"
	"github.com/andrewbackes/chess/position/reader"
	"github.com/andrewbackes/chess/position/square"
	"strconv"
	"strings"
)

// Encode will take a position
func Encode(p reader.PositionReader) (fen string, err error) {

	var boardstr string
	// put what is on each square into a squence (including blanks):
	for i := int(63); i >= 0; i-- {
		sq := p.OnSquare(square.Square(i))
		boardstr += sq.String()
		if i%8 == 0 && i > 0 {
			boardstr += "/"
		}
	}
	// replace groups of spaces with numbers instead
	for i := 8; i > 0; i-- {
		boardstr = strings.Replace(boardstr, strings.Repeat(" ", i), strconv.Itoa(i), -1)
	}
	// Player to move:
	turn := []string{"w", "b"}[p.GetActiveColor()]
	// Castling Rights:
	var rights string
	castles := [][]string{{"K", "Q"}, {"k", "q"}}
	for c := piece.White; c <= piece.Black; c++ {
		for side := board.ShortSide; side <= board.LongSide; side++ {
			if p.GetCastlingRights()[c][side] {
				rights += castles[c][side]
			}
		}
	}
	if rights == "" {
		rights = "-"
	}
	// en Passant:
	enPas := "-"
	if p.GetEnPassant() != square.NoSquare {
		enPas = fmt.Sprint(p.GetEnPassant())
	}
	// Moves and 50 move rule
	fifty := strconv.Itoa(int(p.GetFiftyMoveCount() / 2))
	move := strconv.Itoa(p.GetMoveNumber())
	// all together:
	fen = boardstr + " " + turn + " " + rights + " " + enPas + " " + fifty + " " + move
	return fen, nil
}

// Decode creates a game position from the provided FEN.
func Decode(fen string) (*position.Position, error) {
	words := strings.Split(fen, " ")
	if len(words) < 4 {
		return nil, errors.New("FEN: incomplete fen")
	}
	if words[1] != "w" && words[1] != "b" {
		return nil, errors.New("FEN: can not determine active player")
	}
	p, err := parseBoard(words[0])
	if err != nil {
		return nil, err
	}

	if len(words) >= 6 {
		appendMoveHistory(words[1], words[5], words[4], p)
	}
	if strings.ToLower(words[1]) == "b" {
		p.ActiveColor = piece.Black
	}
	p.CastlingRights = parseCastlingRights(words[2])
	p.EnPassant = parseEnPassantSquare(words[3])
	if len(words) >= 6 {
		p.MoveNumber, _ = strconv.Atoi(words[5])
	}
	return p, nil
}

/*
// DecodeToGame converts a fen string into a game and sets the appropriate tags.
func DecodeToGame(fen string) (*game.Game, error) {
	p, err := Decode(fen)
	if err != nil {
		return nil, err
	}
	g := game.New()
	g.Position = p
	g.Tags["FEN"] = fen
	g.Tags["Setup"] = "1"
	return g, nil
}
*/

func appendMoveHistory(activeColor, moveCount, fiftyMoveCount string, pos *position.Position) error {

	fullMoves, err := strconv.ParseUint(moveCount, 10, 0)
	if err != nil {
		return errors.New("FEN: could not parse move count")
	}
	//halfMoves := ((fullMoves - 1) * 2) + map[string]uint64{"w": 0, "b": 1}[activeColor]
	pos.MoveNumber = int(fullMoves)
	fmc, err := strconv.ParseUint(fiftyMoveCount, 10, 0)
	if err != nil {
		return errors.New("FEN: could not parse fifty move rule count")
	}
	// Since internally we store half moves:
	pos.FiftyMoveCount = (fmc * 2) + map[string]uint64{"w": 0, "b": 1}[activeColor]
	return nil
}

func parseEnPassantSquare(sq string) square.Square {
	if sq != "-" {
		return square.Parse(sq)
	}
	return square.NoSquare
}

func parseCastlingRights(KQkq string) map[piece.Color]map[board.Side]bool {
	return map[piece.Color]map[board.Side]bool{
		piece.White: {board.ShortSide: strings.Contains(KQkq, "K"), board.LongSide: strings.Contains(KQkq, "Q")},
		piece.Black: {board.ShortSide: strings.Contains(KQkq, "k"), board.LongSide: strings.Contains(KQkq, "q")}}
}

// GameFromFEN parses the board passed via FEN and returns a board object.
func parseBoard(board string) (*position.Position, error) {
	p := position.New()
	p.Clear()
	// remove the /'s and replace the numbers with that many spaces
	// so that there is a 1-1 mapping from bytes to squares.
	justBoard := strings.Split(board, " ")[0]
	parsedBoard := strings.Replace(justBoard, "/", "", 9)
	for i := 1; i < 9; i++ {
		parsedBoard = strings.Replace(parsedBoard, strconv.Itoa(i), strings.Repeat(" ", i), -1)
	}
	if len(parsedBoard) < 64 {
		return nil, errors.New("fen: incomplete position")
	}
	pc := map[rune]piece.Type{
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
		if _, ok := pc[k]; ok {
			p.Put(piece.New(color[k], pc[k]), square.Square(63-pos))
		}
	}
	return p, nil
}
