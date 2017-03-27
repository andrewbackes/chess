package pgn

import (
	"fmt"
	"github.com/andrewbackes/chess/fen"
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
	"strings"
	"testing"
)

func decodeToGame(f string) (*game.Game, error) {
	p, err := fen.Decode(f)
	if err != nil {
		return nil, err
	}
	g := game.New()
	g.Position = p
	g.Tags["FEN"] = f
	g.Tags["Setup"] = "1"
	return g, nil
}

func TestPGNString(t *testing.T) {
	str := `[Event ""]
[Site ""]
[Date ""]
[Round ""]
[White ""]
[Black ""]
[Result "1-0"]
[FEN "rnbq1bnr/ppppkppp/8/4p2Q/4P3/8/PPPP1PPP/RNB1KBNR w KQ - 1 1"]

1. h5e5 1-0

`
	pgn, _ := Parse(str)
	if fmt.Sprint(pgn) != str {
		fmt.Print("'", pgn, "'")
		fmt.Println()
		fmt.Print("'", str, "'")
		t.Fail()
	}
}

func TestPGNnullmoves(t *testing.T) {
	expected := `[Result "1-0"]
[Setup "1"]
[FEN "rnbq1bnr/ppppkppp/8/4p2Q/4P3/8/PPPP1PPP/RNB1KBNR w KQ - 1 3"]

3. h5e5 1-0

`
	expectedAlt := `[Result "1-0"]
[FEN "rnbq1bnr/ppppkppp/8/4p2Q/4P3/8/PPPP1PPP/RNB1KBNR w KQ - 1 3"]
[Setup "1"]

3. h5e5 1-0

`
	test := "rnbq1bnr/ppppkppp/8/4p2Q/4P3/8/PPPP1PPP/RNB1KBNR w KQ - 1 3"
	g, _ := decodeToGame(test)
	m, err := g.Position.ParseMove("Qxe5#")
	if err != nil {
		t.Error("couldnt parse move")
	}
	g.MakeMove(m)
	got := Encode(g).String()
	if got != expected && got != expectedAlt {
		t.Log("wanted:\n", expected)
		t.Log("or:\n", expectedAlt)
		t.Log("got:\n", got)
		t.Fail()
	}
}

func TestPGNoutput(t *testing.T) {
	expected := `[Result "1-0"]

1. e2e4 e7e5 2. d1h5 e8e7 3. h5e5 1-0

`
	g := game.New()
	moves := []string{"e4", "e5", "Qh5", "Ke7", "Qxe5#"}
	for _, move := range moves {
		m, err := g.Position.ParseMove(move)
		if err != nil {
			t.Error("couldnt parse move")
		}
		g.MakeMove(m)
	}
	if Encode(g).String() != expected {
		fmt.Print("expected:\n'", expected, "'")
		fmt.Print("got:\n'", Encode(g), "'")
		t.Fail()
	}
}

func TestReadOnePGN(t *testing.T) {
	input := `[Event "one"]
[Round "1"]
[Result "1-0"]

1. e2e4 e7e5 2. d1h5 e8e7 3. h5e5 1-0
`
	games, err := Read(strings.NewReader(input))
	if err != nil || len(games) != 1 {
		t.Log(games)
		t.Log(err)
		t.Fail()
	}

}

func TestReadTwoPGN(t *testing.T) {
	input := `[Event "one"]
[Round "1"]
[Result "1-0"]

1. e2e4 e7e5 2. d1h5 e8e7 3. h5e5 1-0

[Event "two"]
[Round "2"]
[Result "1/2-1/2"]

1. e2e4 e7e5 2. d1h5 e8e7 1/2-1/2
`
	games, err := Read(strings.NewReader(input))
	if err != nil || len(games) != 2 {
		t.Log(games)
		t.Log(err)
		t.Fail()
		t.Fail()
	}

}

func TestReadThreePGN(t *testing.T) {
	input := `[Event "one"]
[Round "1"]
[Result "1-0"]

1. e2e4 e7e5 2. d1h5 e8e7 3. h5e5 1-0

[Event "one"]
[Round "1"]
[Result "1-0"]

1. e2e4 e7e5 2. d1h5 e8e7 3. h5e5 1-0

[Event "one"]
[Round "1"]
[Result "1-0"]

1. e2e4 e7e5 2. d1h5 e8e7 3. h5e5 1-0
`
	games, err := Read(strings.NewReader(input))
	if err != nil || len(games) != 3 {
		t.Log(games)
		t.Log(err)
		t.Fail()
		t.Fail()
	}
}

func TestReadPGN(t *testing.T) {
	input := `[Event "one"]
[Round "1"]
[Result "1-0"]

1. e2e4 e7e5 2. d1h5 e8e7 3. h5e5 1-0
`
	games, _ := Read(strings.NewReader(input))
	games[0].Tags["Event"] = "one"
	games[0].Tags["Round"] = "1"
	games[0].Tags["Result"] = "1-0"

	moves := []string{"e2e4", "e7e5", "d1h5", "e8e7", "h5e5"}
	for i, m := range games[0].Moves {
		if m != moves[i] {
			t.Log(games[0].Moves)
			t.Fail()
		}
	}
}

func TestFromPGN(t *testing.T) {
	pgn := New()
	pgn.Tags["Event"] = "test"
	pgn.Moves = []string{"e2e4", "e7e5", "d1h5", "e8e7", "h5e5"}
	game, err := Decode(pgn)
	if err != nil {
		t.Error(err)
	}
	moves := game.Moves
	for i, m := range pgn.Moves {
		if move.Parse(m) != moves[i] {
			t.Log(moves)
			t.Fail()
		}
	}
}

func TestStripBracketComments(t *testing.T) {
	moves := "1. e4 d5 { comment here } 2. d4 e5"
	p := removeComments([]byte(moves))
	if string(p) != "1. e4 d5 2. d4 e5" {
		t.Fail()
	}
}

func TestStripColonComments(t *testing.T) {
	moves := "1. e4 d5 ; something here"
	expected := "1. e4 d5"
	p := removeComments([]byte(moves))
	if string(p) != expected {
		t.Log(string(p))
		t.Fail()
	}
}

func TestParsePGN(t *testing.T) {
	input := `[Event "one"]
[Round "1"]
[Result "1-0"]

1. e2e4 e7e5 2. d1h5 e8e7 3. h5e5 1-0

[Event "one"]
[Round "2"]
[Result "1-0"]

1. e2e4 e7e5 2. d1h5 e8e7 3. h5e5 1-0

[Event "one"]
[Round "3"]
[Result "1-0"]

1. e2e4 e7e5 2. d1h5 e8e7 3. h5e5 1-0
`
	pgn, err := Parse(input)
	if err != nil || pgn.Tags["Round"] != "1" {
		t.Fail()
	}
}

func TestStatusStringInDraw(t *testing.T) {
	g := game.New()
	g.Position.Clear()
	g.Position.QuickPut(piece.New(piece.White, piece.King), square.E1)
	g.Position.QuickPut(piece.New(piece.Black, piece.King), square.E8)
	if g.Result() != "1/2-1/2" {
		t.Fail()
	}
}
