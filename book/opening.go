package book

import (
	"errors"
	"github.com/andrewbackes/chess"
)

// Opening is an  opening to a chess game.
type Opening []Entry

// Apply makes the moves in the opening on the game.
func Apply(opening Opening, game *chess.Game) error {
	for _, entry := range opening {
		move, err := game.ParseMove(entry.Move)
		if err != nil {
			return err
		}
		if game.MakeMove(move) != chess.InProgress {
			return errors.New("game ended")
		}
	}
	return nil
}

// RandomOpening picks an opening from the book at random.
func (b *Book) RandomOpening(halfmoves int) (Opening, error) {

	return Opening{}, nil
}
