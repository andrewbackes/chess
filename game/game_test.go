package game

import (
	"errors"
	"fmt"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
	"strconv"
	"strings"
	"testing"
	"time"
)

// This is an example of how you might play a game.
func ExampleGame() {
	// Create a new game:
	g := New()
	// Moves can be created based on source and destination squares:
	f3 := move.Move{Source: square.F2, Destination: square.F3}
	g.MakeMove(f3)
	// They can also be created by parsing algebraic notation:
	e5, _ := g.Position.ParseMove("e5")
	g.MakeMove(e5)
	// Or by using piece coordinate notation:
	g4 := move.Parse("g2g4")
	g.MakeMove(g4)
	// Another example of SAN:
	foolsmate, _ := g.Position.ParseMove("Qh4#")
	// Making a move also returns the game status:
	gamestatus := g.MakeMove(foolsmate)
	fmt.Println(gamestatus == WhiteCheckmated)
	// Output: true
}

func ExampleLegalMoves() {
	game, _ := gameFromFEN("8/8/1KP5/3r4/8/8/8/k7 w - - 0 1")
	moves := game.LegalMoves()
	fmt.Println(moves)
	// Output : map[b6b7:{} b6a7:{} c6c7:{} b6a6:{} b6c7:{}]
}

func TestGamePrint(t *testing.T) {
	tc := NewTimeControl(10*time.Minute, 40, 0, true)
	g := NewTimedGame([2]TimeControl{tc, tc})
	got := fmt.Sprint(g)
	expected := `   +---+---+---+---+---+---+---+---+
 8 | r | n | b | q | k | b | n | r |   Active Color:    White
   +---+---+---+---+---+---+---+---+
 7 | p | p | p | p | p | p | p | p |
   +---+---+---+---+---+---+---+---+
 6 |   |   |   |   |   |   |   |   |   En Passant:      None
   +---+---+---+---+---+---+---+---+
 5 |   |   |   |   |   |   |   |   |   Castling Rights: KQkq
   +---+---+---+---+---+---+---+---+
 4 |   |   |   |   |   |   |   |   |   50 Move Rule:    0
   +---+---+---+---+---+---+---+---+
 3 |   |   |   |   |   |   |   |   |
   +---+---+---+---+---+---+---+---+
 2 | P | P | P | P | P | P | P | P |   White's Clock:   10m0s (40 moves)
   +---+---+---+---+---+---+---+---+
 1 | R | N | B | Q | K | B | N | R |   Black's Clock:   10m0s (40 moves)
   +---+---+---+---+---+---+---+---+
     A   B   C   D   E   F   G   H
`
	if got != expected {
		fmt.Print("'", expected, "'\n")
		fmt.Print("'", got, "'\n")
		t.Fail()
	}
}

func TestNewTimedGame(t *testing.T) {
	standard := TimeControl{
		Time:  40 * time.Minute,
		Moves: 40,
	}
	control := [2]TimeControl{standard, standard}
	NewTimedGame(control)
}

func TestNonexistentMove(t *testing.T) {
	g := New()
	mv := move.Parse("e4e5")
	status := g.MakeMove(mv)
	if status != WhiteIllegalMove {
		t.Error("Got: ", status, " Wanted: ", WhiteIllegalMove)
	}
}

func TestActiveColor(t *testing.T) {
	g := New()
	if g.ActiveColor() != piece.White {
		t.Error("it's white to move")
	}
	g.MakeMove(move.Parse("e2e4"))
	if g.ActiveColor() != piece.Black {
		t.Error("it's black to move")
	}
}

func TestIllegalCheck(t *testing.T) {

}

func TestIllegalCastle(t *testing.T) {
	g, err := gameFromFEN("4k3/8/8/8/6r1/8/8/R3K2R w KQ - 0 1")
	s := g.MakeMove(move.Parse("e1g1"))
	if err != nil || s != WhiteIllegalMove {
		t.Fail()
	}
}

func playTestGame(t *testing.T, g *Game, moves []string, expected GameStatus) error {
	for i, san := range moves {
		move, err := g.Position.ParseMove(san)
		if err != nil {
			return err
		}
		s := g.MakeMove(move)
		if (s != InProgress && i+1 < len(moves)) || (i+1 >= len(moves) && s != expected) {
			return errors.New(fmt.Sprint("half-move ", i, " (", san, ") ended with status ", s))
		}
	}
	return nil
}

func TestTimedOut(t *testing.T) {
	tc := TimeControl{
		Time:  40 * time.Minute,
		Moves: 40,
	}
	g := NewTimedGame([2]TimeControl{tc, tc})
	s := g.MakeTimedMove(move.Parse("e2e4"), 41*time.Minute)
	if s != WhiteTimedOut {
		t.Fail()
	}
}

func timedTestGame() *Game {
	tc := TimeControl{Time: 40 * time.Minute, Moves: 2, Increment: 5 * time.Minute, Repeating: true}
	return NewTimedGame([2]TimeControl{tc, tc})
}

func TestTimeIncrement(t *testing.T) {
	g := timedTestGame()
	s := g.MakeTimedMove(move.Parse("e2e4"), 1*time.Minute)
	if s != InProgress {
		t.Error("game should be in progress")
	}
	if g.control[piece.White].clock != 44*time.Minute {
		t.Error("should have 44 min on clock but have", g.control[piece.White].clock)
	}
}

func TestTimeReset(t *testing.T) {
	g := timedTestGame()
	g.MakeTimedMove(move.Parse("e2e4"), 5*time.Minute)
	g.MakeTimedMove(move.Parse("e7e5"), 5*time.Minute)
	g.MakeTimedMove(move.Parse("d2d4"), 5*time.Minute)
	g.MakeTimedMove(move.Parse("d7d5"), 5*time.Minute)
	if g.control[piece.White].movesLeft != g.control[piece.White].Moves {
		t.Error(g.control[piece.White].movesLeft, "!=", g.control[piece.White].Moves)
	}
}

func TestFiftyMoveRule(t *testing.T) {
	fen := "8/8/2B2k2/8/3r1NKp/3N4/8/8 b - - 0 62"
	g, _ := gameFromFEN(fen)
	g.Position.ActiveColor = piece.Black
	moves := []string{"Rd8", "Kxh4", "Rg8", "Be4", "Rg1", "Nh5+", "Ke6", "Ng3", "Kf6", "Kg4", "Ra1", "Bd5", "Ra5", "Bf3", "Ra1", "Kf4", "Ke6", "Nc5+", "Kd6", "Nge4+", "Ke7", "Ke5", "Rf1", "Bg4", "Rg1", "Be6", "Re1", "Bc8", "Rc1", "Kd4", "Rd1", "Nd3", "Kf7", "Ke3", "Ra1", "Kf4", "Ke7", "Nb4", "Rc1", "Nd5+", "Kf7", "Bd7", "Rf1", "Ke5", "Ra1", "Ng5+", "Kg6", "Nf3", "Kg7", "Bg4", "Kg6", "Nf4+", "Kg7", "Nd4", "Re1", "Kf5", "Rc1", "Be2", "Re1", "Bh5", "Ra1", "Nfe6+", "Kh6", "Be8", "Ra8", "Bc6", "Ra1", "Kf6", "Kh7", "Ng5+", "Kh8", "Nde6", "Ra6", "Be8", "Ra8", "Bh5", "Ra1", "Bg6", "Rf1", "Ke7", "Ra1", "Nf7+", "Kg8", "Nh6+", "Kh8", "Nf5", "Ra7", "Kf6", "Ra1", "Ne3", "Re1", "Nd5", "Rg1", "Bf5", "Rf1", "Ndf4", "Ra1", "Ng6+", "Kg8", "Ne7+", "Kh8", "Ng6+"}
	err := playTestGame(t, g, moves, FiftyMoveRule)
	if err != nil {
		t.Error(err)
	}
}

func TestEnPassantMove(t *testing.T) {
	fen := "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
	g, _ := gameFromFEN(fen)
	g.QuickMove(move.Parse("e2c4"))
	g.QuickMove(move.Parse("c7c5"))
	moves := g.LegalMoves()
	if _, ok := moves[move.Parse("d5c6")]; !ok {
		t.Error("missing legal en passant d5c6")
	}
	g.QuickMove(move.Parse("d5c6"))
	if g.Position.OnSquare(square.C5).Type != piece.None {
		t.Error("en passant pawn not captured")
	}
}

func TestThreeFold(t *testing.T) {
	moves := []string{"Nf3", "d6", "d4", "g6", "c4", "Bg7", "Nc3", "Nf6", "e4", "O-O", "Bd3", "Na6", "a3", "c5", "d5", "e6", "O-O", "exd5", "cxd5", "Nc7", "Be3", "Bg4", "h3", "Bxf3", "Qxf3", "Nd7", "Bf4", "Ne5", "Bxe5", "Bxe5", "Rfe1", "a6", "Qd1", "b5", "Qd2", "Qh4", "Ne2", "f5", "f4", "fxe4", "Bxe4", "Bxf4", "Qd3", "Be5", "Rf1", "c4", "Qc2", "Rae8", "Rae1", "Rxf1+", "Rxf1", "Bxb2", "Rf4", "Qe1+", "Rf1", "Qh4", "Rf4", "Qe1+", "Rf1", "Qh4"}
	g := New()
	err := playTestGame(t, g, moves, Threefold)
	if err != nil {
		t.Error(err)
	}
}

func TestStalemate(t *testing.T) {
	fen := "K7/8/k7/1r6/8/8/8/8 w - - 0 1"
	g, _ := gameFromFEN(fen)
	if g.Status() != Stalemate {
		t.Fail()
	}
}

func gameFromFEN(fen string) (*Game, error) {
	g := New()
	p, err := fromFEN(fen)
	g.Position = p
	return g, err
}

func fromFEN(board string) (*position.Position, error) {
	b := position.New()
	b.Clear()
	// remove the /'s and replace the numbers with that many spaces
	// so that there is a 1-1 mapping from bytes to squares.
	justBoard := strings.Split(board, " ")[0]
	parsedBoard := strings.Replace(justBoard, "/", "", 9)
	for i := 1; i < 9; i++ {
		parsedBoard = strings.Replace(parsedBoard, strconv.Itoa(i), strings.Repeat(" ", i), -1)
	}
	if len(parsedBoard) < 64 {
		return nil, errors.New("fen: could not parse position")
	}
	p := map[rune]piece.Type{
		'P': piece.Pawn, 'p': piece.Pawn,
		'N': piece.Knight, 'n': piece.Knight,
		'B': piece.Bishop, 'b': piece.Bishop,
		'R': piece.Rook, 'r': piece.Rook,
		'Q': piece.Queen, 'q': piece.Queen,
		'K': piece.King, 'k': piece.King}
	color := map[rune]piece.Color{
		'P': piece.White, 'p': piece.Black,
		'N': piece.White, 'n': piece.Black,
		'B': piece.White, 'b': piece.Black,
		'R': piece.White, 'r': piece.Black,
		'Q': piece.White, 'q': piece.Black,
		'K': piece.White, 'k': piece.Black}
	// adjust the bitboards:
	for pos := 0; pos < len(parsedBoard); pos++ {
		if pos > 64 {
			break
		}
		k := rune(parsedBoard[pos])
		if _, ok := p[k]; ok {
			b.Put(piece.New(color[k], p[k]), square.Square(63-pos))
			//b.bitBoard[color[k]][p[k]] |= (1 << uint(63-pos))
		}
	}
	return b, nil
}

func TestGameResultInProgress(t *testing.T) {
	g := New()
	if g.Result() != "*" {
		t.Fail()
	}
}
