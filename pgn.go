package chess

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/andrewbackes/chess/board"
	"io"
	"strconv"
	"strings"
)

// PGN represents a game in Portable Game Notation.
type PGN struct {
	Tags  map[string]string
	Moves []string
}

func EmptyTags() map[string]string {
	tags := make(map[string]string)
	tags["Event"] = ""
	tags["Site"] = ""
	tags["Date"] = ""
	tags["Round"] = ""
	tags["White"] = ""
	tags["Black"] = ""
	tags["Result"] = ""
	return tags
}

// FromPGN returns a Game from a PGN string. The string should only contain one
// game, not a series of games. If you need to load a series of PGN games from a
// file use OpenPGN(filename) instead.
func FromPGN(pgn *PGN) (*Game, error) {
	g := NewGame()
	g.tags = pgn.Tags
	for _, san := range pgn.Moves {
		move, err := g.ParseMove(san)
		if err != nil {
			return nil, err
		}
		g.MakeMove(move)
	}
	return g, nil
}

// NewPGN returns a new blank PGN game.
func NewPGN() *PGN {
	return &PGN{
		Tags: make(map[string]string),
	}
}

// PGN returns the PGN of the game.
func (G *Game) PGN() string {
	var pgn string
	status := G.statusString()
	pgn += G.tagsString(status)
	pgn += fmt.Sprintln("")
	pgn += G.enumerateMoves()
	pgn += status
	pgn += fmt.Sprintln("")
	pgn += fmt.Sprintln("")
	return pgn
}

// tagsString formats the game tags the way PGN does.
func (G *Game) tagsString(status string) string {
	tags := [][]string{
		{"Event", G.tags["Event"]},
		{"Site", G.tags["Site"]},
		{"Date", G.tags["Date"]},
		{"Round", G.tags["Round"]},
		{"White", G.tags["White"]},
		{"Black", G.tags["Black"]},
		{"Result", status},
		/*
			{"WhiteElo", "-"},
			{"BlackElo", "-"},
			{"Time", "-"},
			{"TimeControl", "-"},
		*/
	}
	if G.history.startingFen != "" {
		tags = append(tags, []string{"Setup", "1"})
		tags = append(tags, []string{"FEN", G.history.startingFen})
	}
	var s string
	for _, t := range tags {
		s += fmt.Sprintln("[" + t[0] + " \"" + t[1] + "\"]")
	}
	return s
}

func (G *Game) statusString() string {
	if WhiteWon&G.Status() != 0 {
		return "1-0"
	}
	if BlackWon&G.Status() != 0 {
		return "0-1"
	}
	if Draw&G.Status() != 0 {
		return "1/2-1/2"
	}
	return "*"
}

func (G *Game) enumerateMoves() string {
	moves := ""
	for j, move := range G.history.move {
		if move == board.NullMove {
			// dont print book moves, since the FEN tag would mess it up.
			continue
		}
		if j%2 == 0 {
			moves += strconv.Itoa((j/2)+1) + ". "
		}
		moves += string(move) + " "
	}
	return moves
}

// ParsePGN reads a string containing a single PGN and returns a PGN object.
// To read multiple PGNs from a string use:
//     ReadPGN(strings.NewReader(multiPgnString))
func ParsePGN(pgn string) (*PGN, error) {
	games, err := ReadPGN(strings.NewReader(pgn))
	if err != nil {
		return nil, err
	}
	if len(games) < 1 {
		return nil, errors.New("could not read game")
	}
	return games[0], nil
}

// Line by line reads a pgn file. Turns what is read into a PGNGame struct.
// This is an improvement over LoadPGN() since it goes line by line.
// For huge PGN files, this will work but LoadPGN() will not.
func ReadPGN(file io.Reader) ([]*PGN, error) {
	var GameList []*PGN
	// Read line by line:
	scanner := bufio.NewScanner(file)
	readingmoves := false // flag
	currentGame := NewPGN()
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		if line[0] == '[' {
			if readingmoves {
				// since we are no longer reading moves, we know this is a new game
				readingmoves = false
				// so we need to sort out what to do with the game that we previously read:
				GameList = append(GameList, currentGame)
				currentGame = NewPGN()
			}
			key, value := splitTag(line)
			currentGame.Tags[string(key)] = string(value)
		} else {
			readingmoves = true
			appendMoves(currentGame, line)
		}
	}
	GameList = append(GameList, currentGame)
	err := scanner.Err()
	return GameList, err
}

func appendMoves(game *PGN, line []byte) {
	// example: 1. e2e4 d7d5 2. b1c3 f7f5 {asd asd} 3. a2a3 ;asdasdasdasd"
	l := removeComments(line)
	//l = RemoveNumbering(line)
	moves := strings.Split(string(l), " ")
	for _, m := range moves {
		if strings.Contains(m, ".") {
			m = strings.Split(m, ".")[1]
			m = strings.Trim(m, " ")
		}
		//if isMove(StripAnnotations(m)) {
		if m != "1/2-1/2" && m != "1-0" && m != "0-1" && m != "*" && m != "" {
			game.Moves = append(game.Moves, m)
		}
		//}

	}
}

func removeComments(line []byte) []byte {
	marker := 0
	i := 0
	for {
		if line[i] == ';' {
			if i-1 < 0 {
				return []byte{}
			}
			return line[:i-1]
		}
		if line[i] == '{' {
			marker = i
		}
		if line[i] == '}' {
			if marker-1 < 0 {
				marker = marker + 1
			}
			if i+1 > len(line)-1 {
				line = line[:marker-1]
			} else {
				line = append(line[:marker-1], line[i+1:]...)
			}
			i = marker
		}
		i++
		if i > len(line)-1 {
			break
		}
	}
	return line
}

/*
func removeNumbering(line []byte) []byte {
	marker := 0
	i := 0
	for {
		if line[i] == ' ' {
			marker = i
		}
		if line[i] == '.' {
			if i+1 > len(line)-1 {
				line = line[:marker]
			} else if marker == 0 && i+2 < len(line) {
				line = line[i+2:]
			} else {
				line = append(line[:marker], line[i+1:]...)
			}
			i = marker
		}

		i++
		if i > len(line)-1 {
			break
		}
	}
	return line
}
*/

// takes a tag and returns its key and value components.
// ex: [Event "Testing"] ==> "Event", "Testing"
func splitTag(line []byte) ([]byte, []byte) {

	// Look for a space:
	space := 0
	for i, _ := range line {
		if line[i] == ' ' {
			space = i
			break
		}
	}
	key := line[1:space]
	value := line[space+2 : len(line)-2]

	return key, value
}

/*
// Verifies that the game meets the requirements of the filter.
func satisfiesFilters(game *PGN, filters *[]PGNFilter) bool {
	if filters == nil {
		return true
	}
	for _, f := range *filters {
		v := game.Tags[f.Tag]
		if v == "" || !valuesMatch(v, f.Value) {
			return false
		}
	}
	return true
}


// ex: valuesMatch("2701",">2700") == true
func valuesMatch(value string, constraint string) bool {
	if len(constraint) > 1 {
		c := constraint[1:]
		switch constraint[:1] {
		case ">":
			c_int, _ := strconv.Atoi(c)
			value_int, _ := strconv.Atoi(value)
			return value_int > c_int
		case "<":
			c_int, _ := strconv.Atoi(c)
			value_int, _ := strconv.Atoi(value)
			return value_int < c_int
		case "=":
			return value == c
		}
	}
	return value == constraint
}
*/
