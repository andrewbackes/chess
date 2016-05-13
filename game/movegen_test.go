package game


import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// TestNewGame just makes sure we can get a new game.
func TestRootMoves(t *testing.T) {
	g := NewGame()
	moves := g.LegalMoves()
	if len(moves) != 20 {
		t.Log(moves, len(moves))
		t.Log(g.board)
		t.Fail()
	}
}


func TestPerftSuite(t *testing.T) {
	f := "perftsuite.epd"
	d := 3
	if strings.ToLower(os.Getenv("TEST_FULL_PERFT_SUITE")) == "true" {
		d = 6
	}
	if testing.Short() {
		d = 3
	} else {
		t := time.NewTicker(time.Minute)
		stop := make(chan struct{})
		defer func() {
			t.Stop()
			close(stop)
		}()
		go func() {
			for {
				select {
				case <-t.C:
					fmt.Print(".")
				case <-stop:
					break
				}
			}
		}()
	}
	if err := PerftSuite(f, d, false); err != nil {
		t.Error(err)
	}
}
*/
/*
func TestPerftSuitePos(t *testing.T) {
	f := "perftsuite.epd"
	d := 4
	tests, _ := LoadPerftSuite(f)
	err := CheckPerft(tests[1].fen, d, tests[1].nodes[d])
	if err != nil {
		t.Error(err)
	}
}
*/
/*
d4:
    e2c4 : 84835
    e5c4 : 77751
d3 (e2c4):
    c7c5

    d5c6

    missing: e7c5
*/
/*
func TestWithDivide(t *testing.T) {

	fen := "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
	g, _ := FromFEN(fen)
	g.MakeMove(Move("e2c4"))
	g.MakeMove(Move("c7c5"))
	g.MakeMove(Move("d5c6"))
	divide(g, 1)
	fmt.Print(g.board)
	fmt.Println(g.FEN())
}
*/
