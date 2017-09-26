// Package engines provides an API for UCI/WinBoard chess engines.
package engines

import (
	"bufio"
	"errors"
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/position"
	"io"
	"os/exec"
	"path/filepath"
)

var (
	ErrTimedOut = errors.New("timed out")
)

// Engine is an interface for using different types of engines (UCI or WinBoard)
type Engine interface {

	// Shutdown stops the engine's executable.
	Close() error

	// NewGame signals to the engine that the next search will be on a new game.
	NewGame() error

	// Search finds the best move for the game.
	BestMove(g *game.Game, rawOutput chan []byte) (*SearchInfo, error)

	// Think starts the engine deciding the best move for the given position.
	// It is non-blocking and the output is put into the output channel.
	Think(*position.Position) (output chan string, err error)

	// Stop tells the engine to stop doing what ever its doing.
	Stop() error
}

// SearchInfo is used to transfer information from an engine after searching.
type SearchInfo struct {
	BestMove string
	Ponder   string
	// keys: depth, seldepth, score, lowerbound, upperbound, time, nodes, pv
	Analysis []map[string]string
}

// Exec executes the engine executable and wires up the input and output as Readers and Writers.
func execEngine(enginePath string) (*bufio.Reader, *bufio.Writer, error) {
	fullpath, _ := filepath.Abs(enginePath)
	cmd := exec.Command(fullpath)
	cmd.Dir, _ = filepath.Abs(filepath.Dir(enginePath))

	// Setup the pipes to communicate with the engine:
	StdinPipe, errIn := cmd.StdinPipe()
	if errIn != nil {
		return nil, nil, errors.New("can not establish inward pipe")
	}
	StdoutPipe, errOut := cmd.StdoutPipe()
	if errOut != nil {
		return nil, nil, errors.New("can not establish outward pipe")
	}
	r, w := bufio.NewReader(StdoutPipe), bufio.NewWriter(StdinPipe)

	if err := cmd.Start(); err != nil {
		return nil, nil, errors.New("couldnt execute " + enginePath + " - " + err.Error())
	}

	// Setup up for when the engine exits:
	go func() {
		cmd.Wait()
		//TODO: add some confirmation that the engine has terminated correctly.
	}()

	return r, w, nil
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

func IsUCI(filename string) bool {
	return false
}

func IsWinboard(filename string) bool {
	return false
}
