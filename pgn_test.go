package chess

import (
	"fmt"
	"testing"
)

func TestPGNoutput(t *testing.T) {

}

func TestPGNnullmoves(t *testing.T) {
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
	fmt.Println(g.PGN())
	fmt.Println(g.FEN())
}
