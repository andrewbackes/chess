package engines

import (
	"bufio"
	"github.com/andrewbackes/chess/fen"
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position"
	"github.com/andrewbackes/chess/position/move"
	"strconv"
	"strings"
	"time"
)

// UCIEngine provides an API for working with UCI chess engines.
type UCIEngine struct {
	filepath     string
	reader       *bufio.Reader
	writer       *bufio.Writer
	output       chan []byte
	input        chan []byte
	stop         chan struct{}
	lastGameUsed *game.Game
}

const (
	initTimeout    = 10 * time.Second
	newGameTimeout = 1 * time.Second
)

// NewUCIEngine execs an engine and allows interaction with the engine through its methods.
func NewUCIEngine(filepath string) (*UCIEngine, error) {
	r, w, err := execEngine(filepath)
	if err != nil {
		return nil, err
	}
	e, err := newUCIEngine(filepath, r, w)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func newUCIEngine(filepath string, reader *bufio.Reader, writer *bufio.Writer) (*UCIEngine, error) {
	e := UCIEngine{
		filepath: filepath,
		output:   make(chan []byte, 1024),
		input:    make(chan []byte, 1024),
		stop:     make(chan struct{}),
	}
	e.reader, e.writer = reader, writer
	go sub(reader, e.output, e.stop)
	go pub(e.input, writer, e.stop)
	err := e.initialize()
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// should is a helper to determing if the channel is passing or not.
func should(stop chan struct{}) bool {
	select {
	case <-stop:
		return true
	default:
	}
	return false
}

func (e *UCIEngine) resetStop() {
	select {
	case <-e.stop:
		e.stop = make(chan struct{})
	default:
	}
}

func (e *UCIEngine) initialize() error {
	_, err := e.sendAndWait([]byte("uci"), "uciok", initTimeout, func([]byte) {})
	return err
}

func (e *UCIEngine) sendAndWait(send []byte, expected string, timeout time.Duration, parse func([]byte)) (time.Duration, error) {
	e.input <- send
	start := time.Now()
	for {
		select {
		case line := <-e.output:
			parse(line)
			if len(line) >= len(expected) && string(line[:len(expected)]) == expected {
				return time.Since(start), nil
			}
		case <-e.stop:
			return time.Since(start), nil
		case <-time.After(timeout):
			return time.Since(start), ErrTimedOut
		}
	}
}

func (e *UCIEngine) isReady() bool {
	_, err := e.sendAndWait([]byte("isready"), "readyok", initTimeout, func([]byte) {})
	return err == nil
}

// Close shuts down the engine.
func (e *UCIEngine) Close() error {
	e.input <- []byte("quit")
	close(e.stop)
	// TODO: need a way to kill the process if it doesnt close on its own.
	return nil
}

// NewGame tells the engine that we will be passing positions and thinking on a new game.
func (e *UCIEngine) NewGame() error {
	e.input <- []byte("ucinewgame")
	e.isReady()
	return nil
}

// Stop sends a command to the engine to stop what it is doing.
func (e *UCIEngine) Stop() error {
	e.input <- []byte("stop\n")
	return nil
}

// SetBoard sets the engines internal board to that of the games.
func (e *UCIEngine) SetBoard(g *game.Game) {
	/*
		    if g != e.lastGameUsed {
				e.NewGame()
			}
	*/
	pos := "startpos"
	if g.Tags["FEN"] != "" && g.Tags["FEN"] != "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" {
		pos = g.Tags["FEN"]
	}
	moves := ""
	if len(g.Positions) > 0 {
		moves = " moves "
		for _, pos := range g.Positions {
			if pos.LastMove != move.Null {
				moves += pos.LastMove.String() + " "
			}
		}
	}
	command := "position " + pos + moves
	e.input <- []byte(command)
}

func (e *UCIEngine) setPosition(p *position.Position) error {
	f, err := fen.Encode(p)
	if err != nil {
		return err
	}
	command := "position fen " + f
	e.input <- []byte(command)
	return nil
}

// Think will do an infinite search on the provided position. It is non-blocking
// and returns a buffered channel where the engine's output is streamed.
func (e *UCIEngine) Think(p *position.Position) (output chan string, err error) {
	e.resetStop()
	err = e.setPosition(p)
	if err != nil {
		return
	}
	output = make(chan string, 2048)
	parse := func(info []byte) {
		if info != nil {
			output <- string(info)
		}
	}
	go e.sendAndWait([]byte("go infinite"), "bestmove ", 8760*time.Hour, parse)
	return
}

func infoToMap(commands map[string]int, line []byte) map[string]string {
	words := strings.Split(string(line), " ")
	if len(words) > 0 {
		if words[0] == "info" {
			if info := parseInfo(words, commands); info != nil {
				return info
			}
		}
	}
	return nil
}

// BestMove tells the engine to return what it things is the best move for the current game.
func (e *UCIEngine) BestMove(g *game.Game, rawOutput chan []byte) (*SearchInfo, error) {
	e.resetStop()
	e.SetBoard(g)
	command := "go"

	s := []string{" wtime ", " btime "}
	for c := piece.White; c <= piece.Black; c++ {
		if g.Clock(c) > 0 {
			command += s[c] + roundToMilliseconds(g.Clock(c))
		}
	}
	m := g.MovesLeft(g.ActiveColor())
	if m > 0 {
		command += " movestogo " + strconv.FormatInt(int64(m), 10)
	}
	timeout := (g.Clock(g.ActiveColor()) * 125) / 100 // 25% buffer on time
	si := SearchInfo{}
	commands := uciCommands()
	parse := func(info []byte) {
		if rawOutput != nil {
			rawOutput <- info
		}
		parseAnalysis(&si, commands, info)
	}
	_, err := e.sendAndWait([]byte(command), "bestmove ", timeout, parse)
	return &si, err
}

func parseAnalysis(si *SearchInfo, commands map[string]int, line []byte) {
	words := strings.Split(string(line), " ")
	if len(words) > 0 {
		if words[0] == "info" {
			if info := parseInfo(words, commands); info != nil {
				si.Analysis = append(si.Analysis, info)
			}
		} else if words[0] == "bestmove" {
			si.BestMove, si.Ponder = parseBestMove(words)
		}
	}
}

func parseInfo(words []string, commands map[string]int) map[string]string {
	info := make(map[string]string)
	for i, word := range words {
		if n, ok := commands[word]; ok {
			if n == -1 {
				end := findEnd(words, i+1, commands)
				info[word] = strings.Join(words[i+1:end], " ")
			} else if len(words) > i+n {
				info[word] = strings.Join(words[i+1:i+n+1], " ")
			}
		}
	}
	if len(info) == 0 {
		return nil
	}
	return info
}

func findEnd(words []string, start int, commands map[string]int) int {
	for i := start; i < len(words); i++ {
		if _, exists := commands[words[i]]; exists {
			return i
		}
	}
	return len(words)
}

func parseBestMove(words []string) (string, string) {
	if len(words) >= 4 && words[2] == "ponder" {
		return words[1], words[3]
	} else if len(words) >= 2 {
		return words[1], ""
	}
	return "", ""
}

// uciCommand returns a map where the keys are the commands and the values are the number of
// words following the command are considered its value.
func uciCommands() map[string]int {
	return map[string]int{
		"depth":          1,
		"seldepth":       1,
		"time":           1,
		"nodes":          1,
		"pv":             -1,
		"multipv":        1,
		"score":          2,
		"currmove":       1,
		"currmovenumber": 1,
		"hashfull":       1,
		"nps":            1,
		"tbhits":         1,
		"cpuload":        1,
		"string":         -1,
		"refutation":     1,
		"currline":       -1,
		"lowerbound":     0,
		"upperbound":     0,
	}
}

func roundToMilliseconds(d time.Duration) string {
	ms := int64(d) / 1000000
	return strconv.FormatInt(ms, 10)
}
