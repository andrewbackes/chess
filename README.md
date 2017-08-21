[![Go Report Card](https://goreportcard.com/badge/github.com/andrewbackes/chess)](https://goreportcard.com/report/github.com/andrewbackes/chess) [![GoDoc](https://godoc.org/github.com/andrewbackes/chess?status.svg)](https://godoc.org/github.com/andrewbackes/chess) [![Build Status](https://travis-ci.org/andrewbackes/chess.svg?branch=master)](https://travis-ci.org/andrewbackes/chess) [![Coverage Status](https://coveralls.io/repos/github/andrewbackes/chess/badge.svg?branch=master)](https://coveralls.io/github/andrewbackes/chess?branch=master)

# Chess
Multipurpose chess package for Go/Golang.

## What does it do?
This package provides tools for working with chess games. You can:
- Play games
- Play timed games
- Detect checks, checkmates, and draws (stalemate, threefold, 50 move rule, insufficient material)
- Open, save and filter PGN files or strings
- Open EPD files or strings
- Open and save FENs
- Generate legal moves from any position
- and more

For details you can visit the [godoc](https://godoc.org/github.com/andrewbackes/chess)

## How to get it
If you have your GOPATH set in the ([recommended way](https://golang.org/doc/code.html#GOPATH)) then you can use `go get`:

```go get github.com/andrewbackes/chess/<pkg name>```

otherwise you can clone the repo.

## Opinions

The opinion of this package is that a chess game is a sequence of chess positions. Each position holds all of the state neccessary for a completely new game to start from it. A position is just a snapshot in time right *before* a players move. This includes the ability to detect three fold repetition, times left on the clocks, fifty-move rule, etc. This opinion is opposed to the game itself always keeping track of such things.

## Examples

#### Making Moves
```Go
import (
    "fmt"
    "github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/position"
)

func ExampleFoolsMate() {
	// Create a new game:
	g := game.New()
	// Moves can be created based on source and destination squares:
	f3 := position.NewMove(position.F2, position.F3)
	g.MakeMove(f3)
	// They can also be created by parsing algebraic notation:
	e5, _ := g.Position.ParseMove("e5")
	g.MakeMove(e5)
	// Or by using piece coordinate notation:
	g4 := position.Move("g2g4")
	g.MakeMove(g4)
	// Another example of SAN:
	foolsmate, _ := g.Position.ParseMove("Qh4#")
	// Making a move also returns the game status:
	gamestatus := g.MakeMove(foolsmate)
	fmt.Println(gamestatus == game.WhiteCheckmated)
	// Output: true
}
```

#### Setting up a timed game
```Go
import (
    "fmt"
    "github.com/andrewbackes/chess/game"
    "github.com/andrewbackes/chess/position"
    "time"
)

func ExampleTimedGame() {
	forWhite := game.NewTimeControl(5*time.Minute, 40, 0, true)
	forBlack := game.NewTimeControl(1*time.Minute, 40, 0, true)
	tc := [2]game.TimeControl{forWhite, forBlack}
	g := game.NewTimedGame(tc)
	g.MakeTimedMove(position.Move("e2e4"), 1*time.Minute)
}
```

#### Play a timed game in the console
```Go
import (
	"bufio"
	"fmt"
	"github.com/andrewbackes/chess/game"
	"os"
	"time"
)

func ExamplePlayTimedGame() {
	tc := game.NewTimeControl(40*time.Minute, 40, 0, true)
	g := game.NewTimedGame([2]game.TimeControl{tc, tc})
	console := bufio.NewReader(os.Stdin)
	for g.Status() == game.InProgress {
		fmt.Print(g, "\nMove: ")
		start := time.Now()
		input, _ := console.ReadString('\n')
		if move, err := g.Position.ParseMove(input); err == nil {
			g.MakeTimedMove(move, time.Since(start))
		} else {
			fmt.Println("Couldn't understand your move.")
		}
	}
}
```


#### Importing and filtering a PGN
```Go
import (
    "fmt"
    "github.com/andrewbackes/chess/pgn"
    "os"
)

func ExamplePrintGMGames() {
	f, _ := os.Open("myfile.pgn")
	defer f.Close()
	unfiltered, _ := pgn.Read(f)
	filtered := pgn.Filter(unfiltered, pgn.NewTagFilter("WhiteElo>2600"), pgn.NewTagFilter("BlackElo>2600"))
	for _, game := range filtered {
		fmt.Println(game)
	}
}
```

#### Working with FENs
```Go
import (
    "fmt"
    "github.com/andrewbackes/chess/fen"
)

func ExampleSaavedraPositionFEN() {
	decoded, _ := fen.Decode("8/8/1KP5/3r4/8/8/8/k7 w - - 0 1")
	fmt.Println(decoded)
	// Will Output: chess board of the position.
	encoded, _ := fen.Encode(decoded)
	fmt.Println(encoded)
	// Will Output: the inputted FEN
}
```

#### Legal move generation

```Go
import (
    "fmt"
    "github.com/andrewbackes/chess/fen"
	"github.com/andrewbackes/chess/game"
)

func ExampleSaavedraPositionMoves() {
	game, _ := fen.DecodeToGame("8/8/1KP5/3r4/8/8/8/k7 w - - 0 1")
	moves := game.LegalMoves()
	fmt.Println(moves)
	// Will Output: map[b6b7:{} b6a7:{} c6c7:{} b6a6:{} b6c7:{}]
}
```