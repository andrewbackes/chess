// Package chess contains examples for the other packages in this repo.
package chess

import (
	"bufio"
	"fmt"
	"github.com/andrewbackes/chess/fen"
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/pgn"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
	"os"
	"time"
)

func ExampleFoolsMate() {
	// Create a new game:
	g := game.New()
	// Moves can be created based on source and destination squares:
	f3 := &move.Move{Source: square.F2, Destination: square.F3}
	g.MakeMove(f3)
	// They can also be created by parsing algebraic notation:
	e5, _ := g.Position.ParseMove("e5")
	g.MakeMove(e5)
	// Or by using piece coordinate notation:
	g4 := move.Parse("g2g4")
	g.MakeMove(g4)
	// Another example of SAN:
	foolsmate, _ := g.Position.ParseMove("Qh4#")
	// Making a move also returns the game status:
	gamestatus := g.MakeMove(foolsmate)
	fmt.Println(gamestatus == game.WhiteCheckmated)
	// Output: true
}

func ExampleTimedGame() {
	forWhite := game.NewTimeControl(5*time.Minute, 40, 0, true)
	forBlack := game.NewTimeControl(1*time.Minute, 40, 0, true)
	tc := [2]game.TimeControl{forWhite, forBlack}
	g := game.NewTimedGame(tc)
	g.MakeTimedMove(move.Parse("e2e4"), 1*time.Minute)
}

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

func ExamplePrintGMGames() {
	f, _ := os.Open("myfile.pgn")
	defer f.Close()
	unfiltered, _ := pgn.Read(f)
	filtered := pgn.Filter(unfiltered, pgn.NewTagFilter("WhiteElo>2600"), pgn.NewTagFilter("BlackElo>2600"))
	for _, game := range filtered {
		fmt.Println(game)
	}
}

func ExampleSaavedraPositionFEN() {
	decoded, _ := fen.Decode("8/8/1KP5/3r4/8/8/8/k7 w - - 0 1")
	fmt.Println(decoded)
	// Will Output: chess board of the position.
	encoded, _ := fen.Encode(decoded)
	fmt.Println(encoded)
	// Will Output: the inputted FEN
}
