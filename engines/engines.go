// Package engines provides tools for working with chess engines.
package engines

import (
	"bufio"
	"errors"
	"github.com/andrewbackes/chess"
	"github.com/andrewbackes/chess/board"
	"io"
	"os/exec"
	"path/filepath"
)

// Engine is an interface for using different types of engines (UCI or WinBoard)
type Engine interface {
	Start()
	Search(*chess.Game) (*SearchInfo, error)
}

// SearchInfo holds the data that the engine returned during a search.
type SearchInfo struct {
	BestMove board.Move
	Pv       string
}

// Exec executes the engine executable and wires up the input and output as Readers and Writers
func Exec(enginePath string) (io.Writer, io.Reader, error) {
	fullpath, _ := filepath.Abs(enginePath)
	cmd := exec.Command(fullpath)
	cmd.Dir, _ = filepath.Abs(filepath.Dir(enginePath))

	// Setup the pipes to communicate with the engine:
	StdinPipe, errIn := cmd.StdinPipe()
	if errIn != nil {
		E.LogError("Initializing Engine:" + errIn.Error())
		return nil, nil, errors.New("Error Initializing Engine: can not establish inward pipe.")
	}
	StdoutPipe, errOut := cmd.StdoutPipe()
	if errOut != nil {
		E.LogError("Initializing Engine:" + errOut.Error())
		return nil, nilerrors.New("Error Initializing Engine: can not establish outward pipe.")
	}
	E.writer, E.reader = bufio.NewWriter(StdinPipe), bufio.NewReader(StdoutPipe)

	// Start the engine:
	started := make(chan struct{})
	errChan := make(chan error)
	go func() {
		// Question: Does this force the engine to run in its own thread?
		if err := cmd.Start(); err != nil {
			errChan <- err
			return
			//return errors.New("Error executing " + E.Path + " - " + err.Error())
		}
		close(started)
	}()
	select {
	case <-started:
	case e := <-errChan:
		return errors.New("Error starting engine:" + e.Error())
	}

	// Get the engine ready:
	if err := E.Initialize(); err != nil {
		E.LogError("Initializing Engine: " + err.Error())
		//E.Shutdown()
		//cmd.Process.Kill()
		//return err
	}

	//E.NewGame()

	// Setup up for when the engine exits:
	go func() {
		cmd.Wait()
		//TODO: add some confirmation that the engine has terminated correctly.
	}()

	return nil
}
