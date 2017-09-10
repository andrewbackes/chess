package engines

import (
	"bufio"
	"github.com/andrewbackes/chess/game"
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

// Makes sure the engine can go through the initialization process error free.
func TestNewUCIEngine(t *testing.T) {
	if os.Getenv("TEST_ENGINE") == "" {
		t.SkipNow()
	}
	e, err := NewUCIEngine(os.Getenv("TEST_ENGINE"))
	if err != nil {
		t.Fail()
	}
	e.Close()
}

func TestParseInfo(t *testing.T) {
	commands := uciCommands()
	tests := []string{
		"info depth 2 seldepth 5 score cp 100 lowerbound pv e2e4 d7d5",
		"info pv e2e4 d7d5 depth 2 seldepth 5 score cp 100 lowerbound",
		"info pv e2e4 d7d5 depth 2 seldepth 5 score cp 100 lowerbound ",
		"info pv e2e4 d7d5 depth 2 score cp 100 lowerbound seldepth 5",
	}
	expected := map[string]string{
		"depth":      "2",
		"seldepth":   "5",
		"score":      "cp 100",
		"lowerbound": "",
		"pv":         "e2e4 d7d5",
	}
	for _, test := range tests {
		ret := parseInfo(strings.Split(test, " "), commands)
		if len(ret) != 5 {
			t.Log(test)
			t.Log("len(got)=", len(ret), "wanted 5")
			t.Log(ret)
			t.Fail()
		}
		for k, v := range expected {
			if ret[k] != v {
				t.Log(test)
				t.Log("missing", k, "got '", ret[k], "' but wanted '", v, "'")
				t.Fail()
			}
		}
	}
}

type mockWriter struct {
	data []byte
}

func (w *mockWriter) Write(p []byte) (n int, err error) {
	w.data = append(w.data, p...)
	return len(w.data), nil
}

func TestUCIBestMove(t *testing.T) {
	output := []string{
		"uciok\n",
		"info depth 2 seldepth 5 score cp 100 lowerbound pv e2e4 d7d5\n",
		"bestmove e2e4 ponder d7d5\n",
	}
	r := bufio.NewReader(strings.NewReader(strings.Join(output, "")))
	w := bufio.NewWriter(&mockWriter{})
	e, err := newUCIEngine("", r, w)
	if err != nil {
		t.Fail()
	}
	g := game.New()
	sr, err := e.BestMove(g, nil)
	if sr == nil || sr.BestMove != "e2e4" || sr.Ponder != "d7d5" {
		t.Log(sr)
		t.Fail()
	}
}
