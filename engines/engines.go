// Package engines provides tools for working with chess engines.
package engines

import (
	"bufio"
	"errors"
	"github.com/andrewbackes/chess"
	"os/exec"
	"path/filepath"
)

// Engine is an interface for using different types of engines (UCI or WinBoard)
type Engine interface {

	// Shutdown stops the engine's executable.
	Close() error

	// NewGame signals to the engine that the next search will be on a new game.
	NewGame() error

	// Search finds the best move for the game.
	BestMove(*chess.Game) (*SearchResult, error)

	// Stop tells the engine to stop doing what ever its doing.
	Stop() error
}

// SearchResult is used to transfer information from an engine after searching.
type SearchResult struct {
	BestMove string
	Ponder   string
	// keys: depth, seldepth, score, lowerbound, upperbound, time, nodes, pv
	Analysis []map[string]string
}

// Exec executes the engine executable and wires up the input and output as Readers and Writers.
func Exec(enginePath string) (*bufio.Reader, *bufio.Writer, error) {
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
