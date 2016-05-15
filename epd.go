package chess

import (
	"os"
)

// FromEPD returns a game based on the epd provided.
func FromEPD(epd string) *Game {
	return NewGame()
}

// OpenEPD loads a file with new line delimited epd's into a slice of Games.
func OpenEPD(f os.File) []*Game {
	return nil
}
