[![Go Report Card](https://goreportcard.com/badge/github.com/andrewbackes/chess)](https://goreportcard.com/report/github.com/andrewbackes/chess) [![GoDoc](https://godoc.org/github.com/andrewbackes/chess?status.svg)](https://godoc.org/github.com/andrewbackes/chess) [![Build Status](https://travis-ci.org/andrewbackes/chess.svg?branch=master)](https://travis-ci.org/andrewbackes/chess) [![Coverage Status](https://coveralls.io/repos/github/andrewbackes/chess/badge.svg?branch=master)](https://coveralls.io/github/andrewbackes/chess?branch=master)

# Chess
Multipurpose chess package for Go/Golang.

## What does it do?
This package provides tools for working with chess games. You can:
- Play games
- Play timed games
- Detect checks, checkmates, and draws (stalemate, threefold, 50 move rule, insufficient material)
- Open PGN files or strings and filter them
- Open EPD files or strings
- Load FENs
- Generate legal moves from any position.
- and more

For details you can visit the [godoc](https://godoc.org/github.com/andrewbackes/chess)

## How to get it
If you have your GOPATH set in the recommended way ([golang.org](https://golang.org/doc/code.html#GOPATH)):

```go get github.com/andrewbackes/chess```

otherwise you can clone the repo.

## Examples

#### Playing a game
```
import (
    "fmt"
    "github.com/andrewbackes/chess"
	"github.com/andrewbackes/chess/board"
	"github.com/andrewbackes/chess/piece"
)

func FoolsMate() {
	// Create a new game:
	g := chess.NewGame()
	// Moves can be created based on source and destination squares:
	f3 := board.NewMove(board.F2, board.F3)
	g.MakeMove(f3)
	// They can also be created by parsing algebraic notation:
	e5, _ := g.ParseMove("e5")
	g.MakeMove(e5)
	// Or by using piece coordinate notation:
	g4 := board.Move("g2g4")
	g.MakeMove(g4)
	// Another example of SAN:
	foolsmate, _ := g.ParseMove("Qh4#")
	// Making a move also returns the game status:
	gamestatus := g.MakeMove(foolsmate)
	fmt.Println(gamestatus == chess.WhiteCheckmated)
	// Output: true
}
```

#### Importing and filtering a PGN
```
import (
    "fmt"
    "github.com/andrewbackes/chess"
    "os"
)

func PrintGrandmasterGames() {
    f, _ := os.Open("myfile.pgn")
    pgns := chess.ReadPGN(f)
    filtered := chess.FilterPGNs(pgns, NewTagFilter("WhiteElo>2600"), NewTagFilter("BlackElo>2600"))
	for _, pgn := range filtered {
		fmt.Println(pgn)
	} 
}
```