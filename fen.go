package chess

import (
	"errors"
	"github.com/andrewbackes/chess/board"
	"github.com/andrewbackes/chess/piece"
	"strconv"
	"strings"
)

// FEN returns the fen of the current position of the game.
func (G *Game) FEN() string {

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

// FromFEN creates a game from the provided FEN.
func FromFEN(fen string) (*Game, error) {
	G := NewGame()
	words := strings.Split(fen, " ")
	if len(words) < 6 {
		return nil, errors.New("FEN: incomplete fen")
	}
	if words[1] != "w" && words[1] != "b" {
		return nil, errors.New("FEN: can not determine active player")
	}
	G.board = *board.FromFEN(words[0])
	h, _ := parseMoveHistory(words[1], words[5], words[4])
	G.history = *h
	G.history.castlingRights = parseCastlingRights(words[2])
	G.history.enPassant = parseEnPassantSquare(words[3])
	return G, nil
}
