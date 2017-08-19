package pgn

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/andrewbackes/chess/game"
	"io"
	"strings"
)

// PGN represents a game in Portable Game Notation.
type PGN struct {
	Tags         map[string]string
	Moves        []string
	FirstMoveNum int
}

func (p PGN) String() string {
	s := ""
	ordering := []string{"Event", "Site", "Date", "Round", "White", "Black", "Result"}
	for _, t := range ordering {
		if v, ok := p.Tags[t]; ok {
			s += fmt.Sprint("[", t, " ", "\"", v, "\"]\n")
		}
	}
	alreadyPrinted := func(t string) bool {
		for _, v := range ordering {
			if v == t {
				return true
			}
		}
		return false
	}
	for t, v := range p.Tags {
		if !alreadyPrinted(t) {
			s += fmt.Sprint("[", t, " ", "\"", v, "\"]\n")
		}
	}
	s += fmt.Sprintln()
	for i, m := range p.Moves {
		if i%2 == 0 {
			s += fmt.Sprint(p.FirstMoveNum+(i/2), ". ")
		}
		s += fmt.Sprint(m, " ")
	}
	s += fmt.Sprintln(p.Tags["Result"])
	s += fmt.Sprintln()
	return s
}

// MarshalText allows PGN to implement the TextMarshaler interface.
func (p *PGN) MarshalText() (text []byte, err error) {
	text = []byte(p.String())
	err = nil
	return
}

// UnmarshalText allows PGN to implement the TextUnmarshaler interface.
func (p *PGN) UnmarshalText(text []byte) error {
	pgn, err := Parse(string(text))
	if err != nil {
		return err
	}
	p.Tags = pgn.Tags
	p.Moves = pgn.Moves
	return nil
}

// Decode returns a Game from a PGN struct. To load a PGN string ParsePGN()
// or use ReadPGN() to load it from a file.
func Decode(pgn *PGN) (*game.Game, error) {
	g := game.New()
	g.Tags = pgn.Tags
	for _, san := range pgn.Moves {
		move, err := g.Position().ParseMove(san)
		if err != nil {
			return nil, err
		}
		g.MakeMove(move)
	}
	return g, nil
}

// New returns a new blank PGN game.
func New() *PGN {
	return &PGN{
		Tags:         make(map[string]string),
		FirstMoveNum: 1,
	}
}

//Encode returns the PGN of the game.
// If you want to see it as a string then you can use:
// 		G.PGN().String()
// or:
//		G.PGN().UnmarshalText()
func Encode(G *game.Game) *PGN {
	pgn := New()
	//G.appendTags()
	pgn.Tags = G.Tags
	pgn.Tags["Result"] = G.Result()
	/*
		firstRealMove := 0

		for i, move := range G.Moves {
			if move != position.NullMove {
				firstRealMove = i
				break
			}
		}
		pgn.FirstMoveNum = firstRealMove/2 + 1
	*/
	//firstRealMove := G.Position().MoveNumber - (len(G.Positions) / 2)
	//pgn.FirstMoveNum = firstRealMove
	pgn.FirstMoveNum = G.Positions[0].MoveNumber
	for i := 1; i < len(G.Positions); i++ {
		pgn.Moves = append(pgn.Moves, G.Positions[i].LastMove.String())
	}
	return pgn
}

// Parse reads a string containing a single PGN and returns a PGN object.
// To read multiple PGNs from a string use:
//     Read(strings.NewReader(multiPgnString))
func Parse(pgn string) (*PGN, error) {
	games, err := Read(strings.NewReader(pgn))
	if err != nil {
		return nil, err
	}
	if len(games) < 1 {
		return nil, errors.New("could not read game")
	}
	return games[0], nil
}

// Read reads the passed file line by line. Each game that is read is loaded
// into a PGN. If you want to load a PGN fron a string you can use ParsePGN(your_pgn_str)
func Read(file io.Reader) ([]*PGN, error) {
	var GameList []*PGN
	// Read line by line:
	scanner := bufio.NewScanner(file)
	readingmoves := false // flag
	currentGame := New()
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
				currentGame = New()
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
// minimum tag: [ T "" ]
func splitTag(line []byte) ([]byte, []byte) {
	if len(line) < 6 {
		return nil, nil
	}
	// Look for a space:
	space := 0
	for i := range line {
		if line[i] == ' ' {
			space = i
			break
		}
	}
	if space == 0 {
		return nil, nil
	}
	key := line[1:space]
	var value []byte
	if len(line)-2 > 0 && space+2 < len(line) {
		value = line[space+2 : len(line)-2]
	}

	return key, value
}
