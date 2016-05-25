// Package uci is for working with UCI chess engines.
package uci

import (
	"bufio"
	"errors"
	"github.com/andrewbackes/chess"
	"github.com/andrewbackes/chess/board"
	"github.com/andrewbackes/chess/engines"
	"github.com/andrewbackes/chess/piece"
	"io"
	"strconv"
	"strings"
	"time"
)

// Engine represents a chess engine that uses the UCI protocol.
type Engine struct {
	filepath     string
	reader       *bufio.Reader
	writer       *bufio.Writer
	output       chan []byte
	input        chan []byte
	stop         chan struct{}
	lastGameUsed *chess.Game
}

const (
	initTimeout    = 10 * time.Second
	newGameTimeout = 1 * time.Second
)

func NewEngine(filepath string) (*Engine, error) {
	r, w, err := engines.Exec(filepath)
	if err != nil {
		return nil, err
	}
	e, err := newEngine(filepath, r, w)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func newEngine(filepath string, reader *bufio.Reader, writer *bufio.Writer) (*Engine, error) {
	e := Engine{
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

func pub(source chan []byte, dest *bufio.Writer, stop chan struct{}) {
	for {
		select {
		case message := <-source:
			dest.Write(message) // TODO: error handling
			dest.WriteByte('\n')
			dest.Flush()
		case <-stop:
			return
		}
	}
}

func sub(source *bufio.Reader, dest chan []byte, stop chan struct{}) {
	for {
		line, err := source.ReadBytes('\n')
		if err == io.EOF {
			return
		}
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			dest <- line[:len(line)-2]
		} else if len(line) >= 1 && line[len(line)-1] == '\n' {
			dest <- line[:len(line)-1]
		} else {
			dest <- line
		}
		if should(stop) {
			return
		}
	}
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

func (e *Engine) initialize() error {
	_, err := e.sendAndWait([]byte("uci"), "uciok", initTimeout, func([]byte) {})
	return err
}

func (e *Engine) sendAndWait(send []byte, expected string, timeout time.Duration, parse func([]byte)) (time.Duration, error) {
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
			return time.Since(start), errors.New("timed out waiting for '" + expected + "' after '" + string(send) + "'")
		}
	}
}

func (e *Engine) isReady() bool {
	_, err := e.sendAndWait([]byte("isready"), "readyok", initTimeout, func([]byte) {})
	return err == nil
}

// Close shuts down the engine.
func (e *Engine) Close() error {
	e.input <- []byte("quit")
	close(e.stop)
	// TODO: need a way to kill the process if it doesnt close on its own.
	return nil
}

// NewGame tells the engine that we will be passing positions and thinking on a new game.
func (e *Engine) NewGame() {
	e.input <- []byte("ucinewgame")
	e.isReady()
}

// Stop sends a command to the engine to stop what it is doing.
func (e *Engine) Stop() error {
	e.input <- []byte("stop\n")
	return nil
}

// SetBoard sets the engines internal board to that of the games.
func (e *Engine) SetBoard(g *chess.Game) {
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
	if len(g.MoveHistory()) > 0 {
		moves = " moves "
		for _, move := range g.MoveHistory() {
			if move != board.NullMove {
				moves += string(move) + " "
			}
		}
	}
	command := "position " + pos + moves
	e.input <- []byte(command)
}

// BestMove tells the engine to return what it things is the best move for the current game.
func (e *Engine) BestMove(g *chess.Game) (*engines.SearchResult, error) {
	command := "go"

	s := []string{" wtime ", " btime "}
	for c := piece.White; c <= piece.Black; c++ {
		if g.Clock(c) > 0 {
			command += s[c] + roundToMilliseconds(g.Clock(c))
		}
	}
	m := g.MovesLeft(g.ActiveColor())
	if m > 0 {
		command += " movestogo " + strconv.FormatInt(m, 10)
	}
	timeout := (g.Clock(g.ActiveColor()) * 125) / 100 // 25% buffer on time
	si := engines.SearchResult{}
	commands := map[string]int{}
	parse := func(info []byte) {
		parseAnalysis(&si, commands, info)
	}
	e.sendAndWait([]byte(command), "bestmove ", timeout, parse)
	return &si, nil
}

func parseAnalysis(si *engines.SearchResult, commands map[string]int, line []byte) {
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
