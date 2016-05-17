package chess

import (
	"testing"
)

func TestPGNnullmoves(t *testing.T) {
	expected := `[Event ""]
[Site ""]
[Date ""]
[Round ""]
[White ""]
[Black ""]
[Result "1-0"]
[Setup "1"]
[FEN "rnbq1bnr/ppppkppp/8/4p2Q/4P3/8/PPPP1PPP/RNB1KBNR w KQ - 1 3"]

3. h5e5 1-0

`
	g, _ := FromFEN("rnbq1bnr/ppppkppp/8/4p2Q/4P3/8/PPPP1PPP/RNB1KBNR w KQ - 1 3")
	m, err := g.ParseMove("Qxe5#")
	if err != nil {
		t.Error("couldnt parse move")
	}
	g.MakeMove(m)
	if g.PGN() != expected {
		t.Log(g.PGN())
		t.Fail()
	}
}

func TestPGNoutput(t *testing.T) {
	expected := `[Event ""]
[Site ""]
[Date ""]
[Round ""]
[White ""]
[Black ""]
[Result "1-0"]

1. e2e4 e7e5 2. d1h5 e8e7 3. h5e5 1-0

`
	g := NewGame()
	moves := []string{"e4", "e5", "Qh5", "Ke7", "Qxe5#"}
	for _, move := range moves {
		m, err := g.ParseMove(move)
		if err != nil {
			t.Error("couldnt parse move")
		}
		g.MakeMove(m)
	}
	if g.PGN() != expected {
		t.Fail()
	}
}
