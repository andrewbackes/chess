package book

import (
	"errors"
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/pgn"
	"github.com/andrewbackes/chess/polyglot"
	"github.com/andrewbackes/chess/position"
)

// FromPGN creates an opening book from a PGN. 'depth' is the number of plies
// to include in the opening book.
func FromPGN(pgns []*pgn.PGN, depth int) (*Book, error) {
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

func (b *Book) addPgn(pgn *pgn.PGN, depth int) {
	g := game.New()
	for d, m := range pgn.Moves {
		if d >= depth {
			return
		}
		if mv, err := g.Position.ParseMove(m); err == nil {
			key := polyglot.Encode(g.Position)
			b.addMove(key, mv)
			status := g.MakeMove(mv)
			if status != game.InProgress {
				return
			}
		} else {
			return
		}
	}
}

func (b *Book) addMove(key uint64, move position.Move) {
	ml := b.Positions[key]
	for i := range ml {
		if string(ml[i].Move) == string(move) {
			ml[i].Weight++
			b.Positions[key] = ml
			return
		}
	}
	b.Positions[key] = append(b.Positions[key], Entry{Move: string(move), Weight: 1})
}
