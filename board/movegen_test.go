package game

/*
import (
	"bufio"
	"errors"
	"fmt"
	"github.com/andrewbackes/chess/diag"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)
*/

import (
	"github.com/andrewbackes/chess/diag"
	"testing"
)

func TestCallingDiag(t *testing.T) {
	g := NewGame()
	diag.Divide(g, 1)
}

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

/*
func TestPerftSuite(t *testing.T) {
	f := "perftsuite.epd"
	d := 3
	if strings.ToLower(os.Getenv("TEST_FULL_PERFT_SUITE")) == "true" {
		d = 6
	}
	if testing.Short() {
		d = 1
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
	if err := perftSuite(f, d, false); err != nil {
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

/*******************************************************************************

	Perft Suite:

*******************************************************************************/

/*
func perftSuite(filename string, maxdepth int, failFast bool) error {

	test, err := loadPerftSuite(filename)
	if err != nil {
		return err
	}
	for i, t := range test {

		fmt.Print("FEN ", i+1, ": \n")
		for depth, nodes := range t.nodes {
			if depth > maxdepth {
				break
			}
			er := checkPerft(t.fen, depth, nodes)
			if er != nil {
				err = er
			}
			if er != nil && failFast {
				return err
			}
		}
		fmt.Print("\n")
	}
	return err
}

func checkPerft(fen string, depth int, nodes uint64) error {
	G, err := FromFEN(fen)
	if err != nil {
		return err
	}
	start := time.Now()
	fmt.Print("\tD", depth, ": ")
	perftNodes, checks, castles, mates, captures, promotions, enpassant := perft(G, depth)
	passed := perftNodes == nodes
	fmt.Print(map[bool]string{
		true:  "pass. ",
		false: "FAIL. ",
	}[passed])
	lapsed := time.Since(start)
	fmt.Print(lapsed, " ")
	if !passed {
		fmt.Println(perftNodes, "!=", nodes)
		fmt.Print(
			"\t\tchecks:      \t", checks,
			"\n\t\tcastles:   \t", castles,
			"\n\t\tmates:     \t", mates,
			"\n\t\tcaptures:  \t", captures,
			"\n\t\tpromotions:\t", promotions,
			"\n\t\tenpassant: \t", enpassant)
		err = errors.New("incorrect node count")
	}
	fmt.Println()
	return err
}
*/
