package uci

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestSub(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("a\nb\nc\nd"))
	c := make(chan []byte, 10)
	s := make(chan struct{})
	// just make sure we dont block. should get eof and bail out.
	sub(r, c, s)
}

func readAll(c chan []byte) []string {
	var lines []string
	for {
		select {
		case line := <-c:
			lines = append(lines, string(line))
		default:
			return lines
		}
	}
}

func TestNewEngine(t *testing.T) {
	fmt.Println()
	if os.Getenv("TEST_ENGINE") == "" {
		t.SkipNow()
	}
	e, err := NewEngine(os.Getenv("TEST_ENGINE"))
	if err != nil {
		t.Fail()
	}
	e.Close()
}
