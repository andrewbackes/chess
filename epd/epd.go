// Package epd is for working with Extended Position Description.
// You can decode and/or open edp files. Combine this with the engines
// package to really get some cool stuff going on.
package epd

import (
	//"os"
	"bufio"
	"errors"
	"fmt"
	"github.com/andrewbackes/chess/fen"
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/position"
	"io"
	"strings"
)

// Operation is an opcode/operand pair.
//
// Opcode mnemonics:
//     acn - analysis count nodes
//     acs - analysis count seconds
//     am - avoid move(s)
//     bm - best move(s)
//     c0 - comment (primary, also c1 though c9)
//     ce - centipawn evaluation
//     dm - direct mate fullmove count
//     draw_accept - accept a draw offer
//     draw_claim - claim a draw
//     draw_offer - offer a draw
//     draw_reject - reject a draw offer
//     eco - Encyclopedia of Chess Openings opening code
//     fmvn - fullmove number
//     hmvc - halfmove clock
//     id - position identification
//     nic - _New In Chess_ opening code
//     noop - no operation
//     pm - predicted move
//     pv - predicted variation
//     rc - repetition count
//     resign - game resignation
//     sm - supplied move
//     tcgs - telecommunication game selector
//     tcri - telecommunication receiver identification
//     tcsi - telecommunication sender identification
//     v0 - variation name (primary, also v1 though v9)
type Operation struct {
	Code    string
	Operand string
}

// EPD is an Extended Position Description. Position is a FEN like representation
// of the board. Operations are the operations to perform on that position.
type EPD struct {
	Position   *position.Position
	Operations []Operation
}

func (e EPD) String() string {
	return fmt.Sprint("Position:   ", e.Position, "\nOperations: ", e.Operations)
}

// Decode turns a string representation of an epd into an object.
func Decode(epd string) (*EPD, error) {
	s := strings.Split(epd, " ")
	posStr := strings.Join(s[:4], " ")
	p, err := fen.Decode(posStr)
	if len(s) <= 4 {
		return &EPD{Position: p, Operations: nil}, err
	}
	opsStr := strings.TrimRight(strings.Join(s[4:], " "), ";")
	ops := strings.Split(opsStr, ";")
	var opers []Operation
	for _, op := range ops {
		pair := strings.Split(strings.TrimSpace(op), " ")
		if len(pair) != 2 {
			return nil, errors.New("epd: could not parse operation")
		}
		opers = append(opers, Operation{Code: pair[0], Operand: pair[1]})
	}
	return &EPD{Position: p, Operations: opers}, nil
}

// ToGame returns a game based on the position in the EPD provided.
func (e EPD) ToGame() *game.Game {
	g := game.New()
	g.Position = e.Position
	return g
}

// Open loads a file with multiple EPD's. Each EPD needs to be on its own line.
func Open(file io.Reader) ([]*EPD, error) {
	scanner := bufio.NewScanner(file)
	var ret []*EPD
	for scanner.Scan() {
		line := scanner.Text()
		epd, err := Decode(line)
		if err != nil {
			return nil, err
		}
		ret = append(ret, epd)
	}
	return ret, nil
}
