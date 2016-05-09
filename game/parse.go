package game

import (
	"errors"
	"strconv"
	"strings"
)

func ParseMove(s string) Move {
	// Todo: convert the different notation types.
	return Move(s)
}

func parseMoveHistory(activeColor, moveCount, fiftyMoveCount string) (*GameHistory, error) {
	h := GameHistory{}
	unknownMove := Move("")
	fullMoves, err := strconv.ParseUint(moveCount, 10, 0)
	if err != nil {
		return nil, errors.New("FEN: could not parse move count")
	}
	halfMoves := ((fullMoves - 1) * 2) + map[string]uint64{"w": 0, "b": 1}[activeColor]
	for i := uint64(0); i < halfMoves; i++ {
		h.move = append(h.move, unknownMove)
	}
	fmc, err := strconv.ParseUint(fiftyMoveCount, 10, 0)
	if err != nil {
		return nil, errors.New("FEN: could not parse fifty move rule count")
	}
	// Since internally we store half moves:
	h.fiftyMoveCount = (fmc * 2) + map[string]uint64{"w": 0, "b": 1}[activeColor]
	return &h, nil
}

func parseEnPassantSquare(sq string) *Square {
	if sq != "-" {
		s := toSquare(sq)
		return &s
	}
	return nil
}

func parseCastlingRights(KQkq string) [2][2]bool {
	return [2][2]bool{
		{strings.Contains(KQkq, "K"), strings.Contains(KQkq, "Q")},
		{strings.Contains(KQkq, "k"), strings.Contains(KQkq, "q")}}
}
