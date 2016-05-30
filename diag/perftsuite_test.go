package diag

import (
	"errors"
	"fmt"
	"github.com/andrewbackes/chess/position"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestPerftSuite(t *testing.T) {
	f := "perftsuite.epd"
	d := 3
	if strings.ToLower(os.Getenv("TEST_FULL_PERFT_SUITE")) == "true" {
		d = 6
	}
	if testing.Short() {
		d = 2
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
	if err := perftSuite(f, d, true); err != nil {
		t.Error(err)
	}
}

/*
func TestPerftSuitePos(t *testing.T) {
	edp, _ := ParseEPD("4k3/8/8/8/8/8/8/4K2R b K - 0 1 ;D1 5 ;D2 75 ;D3 459 ;D4 8290 ;D5 47635 ;D6 899442")
	g, _ := GameFromEPD(*edp)
	fmt.Println(edp)
	fmt.Println(g.ActiveColor())
	fmt.Println(g.LegalMoves())

		g.MakeMove("c2c3")
		g.MakeMove("d7d6")
		m := divide(g, 1)
		fmt.Println(g.board.Moves(g.ActiveColor(), g.history.enPassant, g.history.castlingRights))
		for k, v := range m {
			fmt.Println(k, ":", v)
		}
		fmt.Println(len(m))

	//err := checkPerft(edp.Position, 3, 8902)
	//if err != nil {
	//	t.Fail()
	//}

}
*/
func divide(G *Game, depth int) map[position.Move]uint64 {
	div := make(map[position.Move]uint64)
	//fmt.Println("Depth", depth)
	var nodes, moveCount uint64
	ml := G.LegalMoves()
	toMove := G.ActiveColor()
	for mv := range ml {
		temp := *G
		temp.MakeMove(mv)

		if temp.Check(toMove) == false {
			//Count it for mate:
			moveCount++
			n := perft(&temp, depth-1)
			div[mv] = n
			nodes += n
		}
	}
	return div
}

/*******************************************************************************

	Perft Suite:

*******************************************************************************/

func perftSuite(filename string, maxdepth int, failFast bool) error {

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	tests, err := ReadEPD(f)
	if err != nil {
		return err
	}
	for i, test := range tests {
		fmt.Print("EPD ", i+1, ":  ")
		//fmt.Println(test)
		for depth, op := range test.Operations {
			//fmt.Println(op)
			if depth > maxdepth {
				break
			}
			nodes, er := strconv.ParseUint(op.Operand, 10, 0)
			if er != nil {
				return er
			}
			er = checkPerft(test.Position, depth, nodes)
			if er != nil {
				err = er
				if failFast {
					return err
				}
			}
		}
		fmt.Print("\n")
	}
	return err
}

func checkPerft(fen string, depth int, nodes uint64) error {
	G, err := GameFromFEN(fen)
	//fmt.Println(G.board)
	if err != nil {
		return err
	}
	start := time.Now()
	fmt.Print("\tD", depth, ": ")
	perftNodes := perft(G, depth)
	passed := perftNodes == nodes
	fmt.Print(map[bool]string{
		true:  "pass. ",
		false: "FAIL. ",
	}[passed])
	lapsed := time.Since(start)
	fmt.Print(lapsed, " ")
	if !passed {
		fmt.Println("got", perftNodes, "wanted", nodes)
		err = errors.New("incorrect node count")
	}
	//fmt.Print("\t")
	return err
}

func perft(g *Game, depth int) uint64 {
	if depth == 0 {
		return 1
	}
	toMove := g.ActiveColor()
	var nodes uint64
	ml := g.LegalMoves()
	for mv := range ml {
		temp := *g
		temp.QuickMove(mv)
		if temp.Check(toMove) == false {
			nodes += perft(&temp, depth-1)
		}
	}
	return nodes
}
