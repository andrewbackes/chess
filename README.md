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

#### Setting up a timed game
```
import (
    "fmt"
    "github.com/andrewbackes/chess"
    "github.com/andrewbackes/chess/board"
    "time"
)

func TimedGame() {
    forWhite := chess.NewTimeControl(5 * time.Minute, 40, 0, true)
    forBlack := chess.NewTimeControl(1 * time.Minute, 40, 0, true)
    tc := [2]chess.TimeControl{forWhite, forBlack}
    game := chess.NewTimedGame(tc)
    game.MakeTimedMove(board.Move("e2e4"), 1*time.Minute)
}
```

#### Play a timed game in the console
```
import (
	"bufio"
	"fmt"
	"github.com/andrewbackes/chess"
	"os"
	"time"
)

func main() {
	tc := chess.NewTimeControl(40*time.Minute, 40, 0, true)
	game := chess.NewTimedGame([2]chess.TimeControl{tc, tc})
	console := bufio.NewReader(os.Stdin)
	for game.Status() == chess.InProgress {
		fmt.Print(game, "\nMove: ")
		start := time.Now()
		input, _ := console.ReadString('\n')
		if move, err := game.ParseMove(input); err == nil {
			game.MakeTimedMove(move, time.Since(start))
		} else {
			fmt.Println("Couldn't understand your move.")
		}
	}
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
    pgn := chess.ReadPGN(f)
    filtered := chess.FilterPGNs(pgn, NewTagFilter("WhiteElo>2600"), NewTagFilter("BlackElo>2600"))
	for _, game := range filtered {
		fmt.Println(game)
	} 
}
```

#### Working with FENs
```
import (
    "fmt"
    "github.com/andrewbackes/chess"
)

func SaavedraPosition() {
    game, _ := chess.FromFEN("8/8/1KP5/3r4/8/8/8/k7 w - - 0 1")
    fmt.Println(game)
    // Output: chess board of the position.
    fmt.Println(game.FEN())
    // Output: the inputted FEN
}
```

#### Legal move generation

```
import (
    "fmt"
    "github.com/andrewbackes/chess"
)

func SaavedraPosition() {
    game, _ := chess.FromFEN("8/8/1KP5/3r4/8/8/8/k7 w - - 0 1")
    moves := game.LegalMoves()
    fmt.Println(moves)
    // Output: map[b6b7:{} b6a7:{} c6c7:{} b6a6:{} b6c7:{}]
}
```