// Package uci is for working with UCI chess engines.
package uci

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/andrewbackes/chess"
	"github.com/andrewbackes/chess/engines"
	"io"
	"time"
)

// Engine represents a chess engine that uses the UCI protocol.
type Engine struct {
	filepath string
	reader   *bufio.Reader
	writer   *bufio.Writer
	output   chan []byte
	input    chan []byte
	stop     chan struct{}
}

const (
	initTimeout = 10 * time.Second
)

func NewEngine(filepath string) (*Engine, error) {
	e := Engine{
		filepath: filepath,
		output:   make(chan []byte, 1024),
		input:    make(chan []byte, 1024),
		stop:     make(chan struct{}),
	}
	r, w, err := engines.Exec(filepath)
	if err != nil {
		return nil, err
	}
	e.reader = r
	e.writer = w
	go sub(e.reader, e.output, e.stop)
	go pub(e.input, e.writer, e.stop)
	err = e.initialize()
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
			fmt.Println(string(message))
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

func should(stop chan struct{}) bool {
	select {
	case <-stop:
		return true
	default:
		// don't block
	}
	return false
}

func (e *Engine) initialize() error {
	return e.sendAndWait([]byte("uci"), "uciok", initTimeout, func([]byte) {})
}

func (e *Engine) sendAndWait(send []byte, expected string, timeout time.Duration, perLine func([]byte)) error {
	e.input <- send
	for {
		select {
		case line := <-e.output:
			perLine(line)
			if string(line) == expected {
				return nil
			}
		case <-e.stop:
			return nil
		case <-time.After(timeout):
			return errors.New("timed out waiting for '" + expected + "' after '" + string(send) + "'")
		}
	}
}

func (e *Engine) isready() bool {
	return e.sendAndWait([]byte("isready"), "readyok", initTimeout, func([]byte) {}) == nil
}

func (e *Engine) Close() error {
	e.input <- []byte("quit")
	close(e.stop)
	return nil
}

func (e *Engine) Stop() error {
	return nil
}

func (e *Engine) Search(*chess.Game) (*engines.SearchInfo, error) {
	return nil, nil
}
