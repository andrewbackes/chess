package game

import (
	"fmt"
	"os"
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
	d := 5
	if os.Getenv("TEST_FULL_PERFT_SUITE") != "" {
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
	if err := PerftSuite(f, d, true); err != nil {
		t.Error(err)
	}
}
