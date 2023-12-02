package pgn

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/andrewbackes/chess/fen"
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
)

func decodeToGame(f string) (*game.Game, error) {
	if f == "" {
		return game.New(), nil
	}
	p, err := fen.Decode(f)
	if err != nil {
		return nil, err
	}
	g := game.New()
	g.Positions = append(g.Positions, p)
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
	m, err := g.Position().ParseMove("Qxe5#")
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
		m, err := g.Position().ParseMove(move)
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
	for i, m := range pgn.Moves {
		if move.Parse(m) != game.Positions[i+1].LastMove {
			t.Log(move.Parse(m), "!=", game.Positions[i+1].LastMove)
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
	g.Position().Clear()
	g.Position().QuickPut(piece.New(piece.White, piece.King), square.E1)
	g.Position().QuickPut(piece.New(piece.Black, piece.King), square.E8)
	if g.Result() != "1/2-1/2" {
		t.Fail()
	}
}

func TestEncodeSAN(t *testing.T) {
	testCases := []struct {
		name, FEN string
		moves     []string
		want      *PGN
	}{
		{
			"Empty",
			"", []string(nil),
			&PGN{FirstMoveNum: 1, Tags: map[string]string{"Result": "*"}},
		},
		{
			"StartMove",
			"", []string{"e4", "e5"},
			&PGN{
				FirstMoveNum: 1,
				Tags: map[string]string{
					"Result": "*",
				},
				Moves: []string{"e4", "e5"},
			},
		},
		{
			"TheeMovesMate",
			"", []string{"e4", "e5", "Qh5", "Ke7", "Qxe5#"},
			&PGN{
				FirstMoveNum: 1,
				Tags: map[string]string{
					"Result": "1-0",
				},
				Moves: []string{"e4", "e5", "Qh5", "Ke7", "Qxe5#"},
			},
		},
		{
			"GameWithNearlyAllTypesOfSAN",
			"", []string{"e2e4", "e7e6", "f1c4", "d7d5", "c4d5", "e6d5", "c2c4", "c7c6", "c4d5", "d8a5", "d1b3", "c8g4", "g1f3", "b8d7", "e1g1", "e8c8", "b3b7", "c8b7", "d5d6", "g8h6", "b1c3", "f7f6", "f1e1", "h6f7", "a2a4", "d7e5", "a1a3", "e5f3", "g1h1", "f7e5", "d6d7", "d8c8", "d7c8q", "b7b6", "c8g4", "c6c5", "c3b5", "c5c4", "b2b4", "c4b3", "d2d4", "b3b2", "d4d5", "b2b1n", "d5d6", "e5f7", "d6d7", "b1d2", "a3e3", "d2b3", "e3e2", "b3c5", "e2c2", "c5d3", "g4e6", "b6b7", "d7d8q", "f3e5", "d8c8"},
			&PGN{
				FirstMoveNum: 1,
				Tags: map[string]string{
					"Result": "1-0",
				},
				Moves: []string{"e4", "e6", "Bc4", "d5", "Bxd5", "exd5", "c4", "c6", "cxd5", "Qa5", "Qb3", "Bg4", "Nf3", "Nd7", "O-O", "O-O-O", "Qxb7+", "Kxb7", "d6", "Nh6", "Nc3", "f6", "Re1", "Nf7", "a4", "Nde5", "Ra3", "Nxf3+", "Kh1", "N7e5", "d7", "Rc8", "dxc8=Q+", "Kb6", "Qxg4", "c5", "Nb5", "c4", "b4", "cxb3", "d4", "b2", "d5", "b1=N", "d6", "Nf7", "d7", "Nbd2", "Rae3", "Nb3", "R3e2", "Nc5", "Rc2", "Nd3", "Qe6+", "Kb7", "d8=Q", "Nf3e5", "Qdc8#"},
			},
		},
		{
			"FromFEN-WinningMove",
			"rnbq1bnr/ppppkppp/8/4p2Q/4P3/8/PPPP1PPP/RNB1KBNR w KQ - 1 3",
			[]string{"Qxe5"},
			&PGN{
				Tags: map[string]string{
					"Result": "1-0",
					"FEN":    "rnbq1bnr/ppppkppp/8/4p2Q/4P3/8/PPPP1PPP/RNB1KBNR w KQ - 1 3",
					"Setup":  "1",
				},
				Moves:        []string{"Qxe5#"},
				FirstMoveNum: 3,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Game preparation from FEN and applying moves.
			g, err := decodeToGame(tc.FEN)
			if err != nil {
				t.Fatalf("Can't decode FEN string \"%s\" to a game, due to error: %v", tc.FEN, err)
			}
			for i, move := range tc.moves {
				m, err := g.Position().ParseMove(move)
				if err != nil {
					t.Fatalf("Couldn't parse move #%d: \"%s\", due to error: %v", i, move, err)
				}
				_, err = g.MakeMove(m)
				if err != nil {
					t.Fatalf("Couldn't make move #%d: \"%s\", due to error: %v", i, move, err)
				}
			}

			// Encoding to PGN with moves in SAN.
			enc := EncodeSAN(g)

			// Compare output with wanted result.
			if !reflect.DeepEqual(enc, tc.want) {
				t.Errorf("PGN.Encode(...) failed.\nGot:\n\t%#v,\nwant:\n\t%#v.", enc, tc.want)
			}
		})
	}
}
