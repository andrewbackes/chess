package chess

import (
	"bufio"
	"errors"
	"fmt"
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
	if err := perftSuite(f, d, false); err != nil {
		t.Error(err)
	}
}

/*
func TestPerftSuitePos(t *testing.T) {
	f := "perftsuite.epd"
	d := 4
	tests, _ := loadPerftSuite(f)
	err := checkPerft(tests[1].fen, d, tests[1].nodes[d])
	if err != nil {
		t.Error(err)
	}
}
*/

/*******************************************************************************

	Perft Suite:

*******************************************************************************/

type epdTest struct {
	fen   string
	nodes []uint64
}

func perftSuite(filename string, maxdepth int, failFast bool) error {

	test, err := loadPerftSuite(filename)
	if err != nil {
		return err
	}
	for i, t := range test {

		fmt.Print("FEN ", i+1, ":  \t")
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
	perftNodes := perft(G, depth)
	passed := perftNodes == nodes
	fmt.Print(map[bool]string{
		true:  "pass. ",
		false: "FAIL. ",
	}[passed])
	lapsed := time.Since(start)
	fmt.Print(lapsed, " ")
	if !passed {
		fmt.Println(perftNodes, "!=", nodes)
		err = errors.New("incorrect node count")
	}
	//fmt.Print("\t")
	return err
}

func loadPerftSuite(filename string) ([]epdTest, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)

	var test []epdTest
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, ";")

		var newTest epdTest
		newTest.fen = words[0]
		newTest.nodes = append(newTest.nodes, 1) // depth 0 = 1 node

		for i := 1; i < len(words); i++ {
			n, _ := strconv.ParseUint(strings.Split(words[i], " ")[1], 10, 0)
			newTest.nodes = append(newTest.nodes, n)
		}

		test = append(test, newTest)
	}
	f.Close()

	return test, err
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
