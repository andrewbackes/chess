package chess

import (
	//"os"
	"errors"
	"fmt"
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
	Position   string
	Operations []Operation
}

func (e EPD) String() string {
	return fmt.Sprint("Position:   ", e.Position, "\nOperations: ", e.Operations)
}

// ParseEPD turns a string representation of an epd into an object.
func ParseEPD(epd string) (*EPD, error) {
	// POS 2 3 4 a b; c d; e f;
	s := strings.Split(epd, " ")
	if len(s) <= 4 {
		return &EPD{Position: epd, Operations: nil}, nil
	}
	posStr := strings.Join(s[:4], " ")
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
	return &EPD{Position: posStr, Operations: opers}, nil
}

/*
// FromEPD returns a game based on the epd provided.
func FromEPD(epd EPD) *Game {
	g, err := FromFEN(epd.Position)
	return g, err
}

// OpenEPD loads a file with new line delimited epd's into a slice of Games.
func OpenEPD(f os.File) []*Game {
	return nil
}

func openEPD(f *os.File) ([]epdTest, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)

	var test []epdTest
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, ";")

		var newTest epdTest
		newTest.fen = words[0]
		newTest.nodes = append(newTest.nodes, 1) // depth 0 = 1 node

		for i := 1; i < len(words); i++ {
			n, _ := strconv.ParseUint(strings.Split(words[i], " ")[1], 10, 0)
			newTest.nodes = append(newTest.nodes, n)
		}

		test = append(test, newTest)
	}
	f.Close()

	return test, err
}
*/
