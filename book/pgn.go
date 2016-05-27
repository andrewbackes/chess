package book

import (
	"errors"
	"github.com/andrewbackes/chess"
	"github.com/andrewbackes/chess/board"
)

// FromPGN creates an opening book from a PGN. 'depth' is the number of plies
// to include in the opening book.
func FromPGN(pgns []*chess.PGN, depth int) (*Book, error) {
	if len(pgns) == 0 {
		return nil, errors.New("no games in pgn")
	}
	book := New()
	for _, pgn := range pgns {
		// skip games where we don't know the opening moves
		if pgn.Tags["FEN"] == "" {
			book.addPgn(pgn, depth)
		}
	}
	return book, nil
}

func (b *Book) addPgn(pgn *chess.PGN, depth int) {
	g := chess.NewGame()
	for d, m := range pgn.Moves {
		if d >= depth {
			return
		}
		if mv, err := g.ParseMove(m); err == nil {
			key := g.Polyglot()
			b.addMove(key, mv)
			status := g.MakeMove(mv)
			if status != chess.InProgress {
				return
			}
		} else {
			return
		}
	}
}

func (b *Book) addMove(key uint64, move board.Move) {
	ml := b.Positions[key]
	for i := range ml {
		if string(ml[i].Move) == string(move) {
			ml[i].Weight++
			b.Positions[key] = ml
			return
		}
	}
	b.Positions[key] = append(b.Positions[key], Move{Move: string(move), Weight: 1})
}
