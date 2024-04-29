package game

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/andrewbackes/chess/fen"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
)

// This is an example of how you might play a game.
func ExampleGame() {
	// Create a new game:
	g := New()
	// Moves can be created based on source and destination squares:
	f3 := move.Move{Source: square.F2, Destination: square.F3}
	g.MakeMove(f3)
	// They can also be created by parsing algebraic notation:
	e5, _ := g.Position().ParseMove("e5")
	g.MakeMove(e5)
	// Or by using piece coordinate notation:
	g4 := move.Parse("g2g4")
	g.MakeMove(g4)
	// Another example of SAN:
	foolsmate, _ := g.Position().ParseMove("Qh4#")
	// Making a move also returns the game status:
	gamestatus, _ := g.MakeMove(foolsmate)
	fmt.Println(gamestatus == WhiteCheckmated)
	// Output: true
}

func ExampleLegalMoves() {
	game, _ := gameFromFEN("8/8/1KP5/3r4/8/8/8/k7 w - - 0 1")
	moves := game.LegalMoves()
	fmt.Println(moves)
	// Output: map[c6c7:{} b6a6:{} b6c7:{} b6b7:{} b6a7:{}]
}

func TestGamePrint(t *testing.T) {
	tc := NewTimeControl(10*time.Minute, 40, 0, true)
	g := NewTimedGame(map[piece.Color]TimeControl{piece.White: tc, piece.Black: tc})
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
	control := map[piece.Color]TimeControl{piece.White: standard, piece.Black: standard}
	NewTimedGame(control)
}

func TestNonexistentMove(t *testing.T) {
	g := New()
	mv := move.Parse("e4e5")
	status, _ := g.MakeMove(mv)
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
	s, _ := g.MakeMove(move.Parse("e1g1"))
	if err != nil || s != WhiteIllegalMove {
		t.Fail()
	}
}

func playTestGame(t *testing.T, g *Game, moves []string, expected GameStatus) error {
	for i, san := range moves {
		move, err := g.Position().ParseMove(san)
		if err != nil {
			return err
		}
		s, _ := g.MakeMove(move)
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
	g := NewTimedGame(map[piece.Color]TimeControl{piece.White: tc, piece.Black: tc})
	t.Log(g)
	m := move.Parse("e2e4")
	m.Duration = 41 * time.Minute
	s, e := g.MakeMove(m)
	if s != WhiteTimedOut {
		t.Log(g)
		t.Log(g.LegalMoves())
		t.Log(s, e)
		t.Fail()
	}
}

func timedTestGame() *Game {
	tc := TimeControl{Time: 40 * time.Minute, Moves: 2, Increment: 5 * time.Minute, Repeating: true}
	return NewTimedGame(map[piece.Color]TimeControl{piece.White: tc, piece.Black: tc})
}

func TestTimeIncrement(t *testing.T) {
	g := timedTestGame()
	m := move.Parse("e2e4")
	m.Duration = 1 * time.Minute
	s, _ := g.MakeMove(m)
	if s != InProgress {
		t.Error("game should be in progress")
	}
	if g.Position().Clocks[piece.White] != 44*time.Minute {
		t.Error("should have 44 min on clock but have", g.Position().Clocks[piece.White])
	}
}

func TestTimeReset(t *testing.T) {
	g := timedTestGame()
	timedMove := func(s string, t time.Duration) move.Move {
		m := move.Parse(s)
		m.Duration = t
		return m
	}
	g.MakeMove(timedMove("e2e4", 5*time.Minute))
	g.MakeMove(timedMove("e7e5", 5*time.Minute))
	g.MakeMove(timedMove("d2d4", 5*time.Minute))
	g.MakeMove(timedMove("d7d5", 5*time.Minute))
	if g.Position().MovesLeft[piece.White] != g.control[piece.White].Moves {
		t.Error(g.Position().MovesLeft[piece.White], "!=", g.control[piece.White].Moves)
	}
}

func TestFiftyMoveRule(t *testing.T) {
	fen := "8/8/2B2k2/8/3r1NKp/3N4/8/8 b - - 0 62"
	g, _ := gameFromFEN(fen)
	g.Position().ActiveColor = piece.Black
	moves := []string{"Rd8", "Kxh4", "Rg8", "Be4", "Rg1", "Nh5+", "Ke6", "Ng3", "Kf6", "Kg4", "Ra1", "Bd5", "Ra5", "Bf3", "Ra1", "Kf4", "Ke6", "Nc5+", "Kd6", "Nge4+", "Ke7", "Ke5", "Rf1", "Bg4", "Rg1", "Be6", "Re1", "Bc8", "Rc1", "Kd4", "Rd1", "Nd3", "Kf7", "Ke3", "Ra1", "Kf4", "Ke7", "Nb4", "Rc1", "Nd5+", "Kf7", "Bd7", "Rf1", "Ke5", "Ra1", "Ng5+", "Kg6", "Nf3", "Kg7", "Bg4", "Kg6", "Nf4+", "Kg7", "Nd4", "Re1", "Kf5", "Rc1", "Be2", "Re1", "Bh5", "Ra1", "Nfe6+", "Kh6", "Be8", "Ra8", "Bc6", "Ra1", "Kf6", "Kh7", "Ng5+", "Kh8", "Nde6", "Ra6", "Be8", "Ra8", "Bh5", "Ra1", "Bg6", "Rf1", "Ke7", "Ra1", "Nf7+", "Kg8", "Nh6+", "Kh8", "Nf5", "Ra7", "Kf6", "Ra1", "Ne3", "Re1", "Nd5", "Rg1", "Bf5", "Rf1", "Ndf4", "Ra1", "Ng6+", "Kg8", "Ne7+", "Kh8", "Ng6+"}
	err := playTestGame(t, g, moves, FiftyMoveRule)
	if err != nil {
		t.Log(g)
		t.Log(g.LegalMoves())
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
	if g.Position().OnSquare(square.C5).Type != piece.None {
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

func TestGameResultInProgress(t *testing.T) {
	g := New()
	if g.Result() != "*" {
		t.Fail()
	}
}

func gameFromFEN(FEN string) (*Game, error) {
	p, err := fen.Decode(FEN)
	if err != nil {
		return nil, err
	}
	return &Game{
		control: nil,
		Tags: map[string]string{
			"FEN":   FEN,
			"Setup": "1",
		},
		Positions: []*position.Position{p},
	}, nil
}

// BenchmarkGame benchmarks full games of various lengths using moves as move.Move struct and strings in PCN and SAN notation.
func BenchmarkGame(b *testing.B) {
	benchCases := []struct {
		Name     string
		PCNMoves []string
	}{
		{
			"2-move-fools-mate",
			[]string{"f3", "e5", "g4", "Qh4#"},
		},
		{
			"Random-game-#1338_half-moves-11",
			[]string{"e4", "Na6", "a3", "Nf6", "Nc3", "Rg8", "Bc4", "Nxe4", "Qh5", "Nec5", "Bxf7#"},
		},
		{
			"Random-game-#6589_half-moves-25",
			[]string{"c4", "f5", "Qc2", "g5", "h3", "h6", "h4", "g4", "b3", "d5", "cxd5", "a6", "Qc4", "c5", "h5", "Qc7", "Qxa6", "Kf7", "Qd3", "Kg7", "Qe3", "Qd8", "Qe6", "g3", "Qg6#"},
		},
		{
			"Random-game-#8159_half-moves-50",
			[]string{"b3", "c5", "Nf3", "d5", "Ng5", "Na6", "Nxf7", "Kd7", "Nc3", "Qb6", "Nxd5", "g6", "g3", "Qd6", "e3", "Nc7", "a3", "Ne8", "Be2", "Qc7", "h4", "b5", "Bf3", "Qf4", "Rb1", "b4", "h5", "Bg7", "Bg2", "Qe5", "c4", "Qxc3", "Be4", "e5", "Rh2", "Nef6", "Qe2", "Qxb3", "a4", "Ne8", "Qf3", "Qxc4", "Nh6", "Bxh6", "Qf6", "Ba6", "Ra1", "Qb5", "Qg5", "Qe2#"},
		},
		{
			"Random-game-#2537_half-moves-100",
			[]string{"c4", "f5", "h4", "Na6", "b4", "h5", "g4", "b6", "Nh3", "c6", "d3", "Nf6", "a3", "Rb8", "Rh2", "Rh7", "Ng5", "Bb7", "b5", "Bc8", "Nh3", "hxg4", "Ng1", "Qc7", "d4", "Nh5", "Rh1", "Ra8", "Ra2", "Rb8", "bxa6", "Qg3", "Kd2", "Qf4+", "Kd3", "d5", "Nc3", "c5", "Qd2", "Qh2", "Qb2", "Nf4+", "Kd2", "b5", "h5", "Qxf2", "Nh3", "e5", "cxd5", "e4", "Qb3", "Bb7", "Qxb5+", "Kf7", "Rb2", "e3+", "Kc2", "Re8", "axb7", "Re5", "Rb1", "Rxh5", "Ng5+", "Kf6", "Qc6+", "Kxg5", "Nb5", "c4", "Rg1", "Qxg1", "Nxa7", "Ng2", "Qb5", "Bb4", "Bb2", "Kg6", "d6", "Rd5", "Qc6", "Ba5", "b8=R", "Rh4", "Qc8", "Qf2", "Bxg2", "Qxe2+", "Kc1", "Qxb2+", "Kd1", "Qc3", "R1b7", "Qxa3", "Qd8", "Qc3", "Bf3", "Kh5", "Rc8", "gxf3", "Rc6", "Rhxd4#"},
		},
		{
			"Random-game-#177_half-moves-250",
			[]string{"h3", "a6", "a4", "b6", "e3", "h5", "a5", "d5", "Ke2", "bxa5", "Ra4", "g5", "Rxa5", "Nh6", "d3", "Bf5", "Ra1", "Ng4", "h4", "Ne5", "Ra4", "e6", "Rb4", "Nbc6", "Rc4", "Ng4", "g3", "Nxf2", "Rg4", "Rh7", "Kf3", "Nb8", "Bd2", "Bd6", "Rb4", "g4+", "Kxf2", "Kf8", "Bc1", "Ra7", "Rb7", "Ke7", "Qf3", "Be5", "Qxg4", "Bd4", "b4", "Kd6", "b5", "Ke7", "Ke2", "Bh8", "Bg2", "Bd4", "Qg6", "a5", "c3", "Rg7", "c4", "Qe8", "Rh2", "Qf8", "Bh3", "Rxg6", "Bf1", "Nc6", "bxc6", "Qc8", "Bh3", "Bh8", "Nd2", "d4", "Kf1", "Rxg3", "Rb6", "Be4", "Rb7", "e5", "Kf2", "Rf3+", "Ke1", "Qf5", "Ndxf3", "Bf6", "Rb8", "dxe3", "Bxf5", "Ra8", "Nh3", "Bg7", "Kd1", "Kf6", "Bh7", "Bxh7", "Rd2", "Ra6", "Rb5", "e2+", "Kc2", "e1=Q", "c5", "Rb6", "Nfg5", "Ke7", "Kb1", "Qxd2", "Bxd2", "Bf6", "Ka2", "Bh8", "Be3", "Bxd3", "Rxa5", "Rxc6", "Nf4", "exf4", "Ra6", "Bc3", "Kb3", "Bc4+", "Kc2", "Kf6", "Ne4+", "Kg7", "Ng3", "Ba1", "Nxh5+", "Kg8", "Kc1", "Bb3", "Ra5", "fxe3", "Ra3", "Rf6", "Ra8+", "Kh7", "Ra2", "Bc3", "Rf2", "e2", "Rxf6", "Bd4", "Rh6+", "Kg8", "Ra6", "Ba2", "Rg6+", "Kh7", "Rg7+", "Kh8", "Rg5", "e1=R+", "Kc2", "f6", "Rg3", "Re6", "Rd3", "Re2+", "Rd2", "Rf2", "Nf4", "Ba1", "Nd5", "Bb2", "Nc3", "Bxc3", "Kxc3", "Rh2", "Re2", "Bc4", "Kd2", "Bb3", "Ke1", "Bd1", "Re5", "c6", "Rh5+", "Kg7", "Kxd1", "Rc2", "Rh7+", "Kg6", "h5+", "Kxh7", "Kxc2", "Kg7", "h6+", "Kh7", "Kb2", "f5", "Kb1", "Kg8", "Kc2", "Kh8", "Kb3", "Kg8", "Kb2", "Kh8", "Kb3", "Kh7", "Ka2", "Kxh6", "Ka1", "Kg7", "Kb2", "f4", "Kc1", "Kh6", "Kb1", "Kg7", "Kc2", "Kf6", "Kc3", "Kf5", "Kc4", "Ke5", "Kb4", "f3", "Ka4", "Ke6", "Ka5", "Ke5", "Ka6", "Kf6", "Kb7", "Kg6", "Kxc6", "Kg7", "Kd5", "Kf8", "Ke6", "Ke8", "Kf6", "Kd8", "Kg5", "Kc8", "Kg4", "Kc7", "Kg3", "f2", "Kg4", "f1=B", "Kf5", "Bh3+", "Kf4", "Bf5", "Ke3", "Bg4", "c6", "Kxc6"},
		},
		{
			"Random-game-#577_half-moves-500",
			[]string{"b4", "c6", "e3", "e6", "Qe2", "Na6", "d4", "Nc7", "Bb2", "a6", "a3", "g6", "h4", "Bc5", "e4", "e5", "Nd2", "f6", "Qd1", "Qe7", "Ke2", "g5", "f3", "f5", "hxg5", "h6", "Rh2", "Bxd4", "Rh4", "Rb8", "Rf4", "b6", "Ke1", "Qc5", "Rc1", "Qd6", "Bc4", "Nb5", "Bd5", "Qe6", "Bc4", "hxg5", "Nb1", "Qf6", "Rxf5", "Kf8", "Bxd4", "Nc3", "Bxe5", "Rb7", "Bh2", "Qxf5", "Ba2", "Qxf3", "Qd6+", "Ne7", "Nh3", "Qf7", "Be6", "Qf5", "Kd2", "Na4", "Nc3", "Ke8", "Rh1", "Qd5+", "Qxd5", "Rxh3", "Nb5", "Rg3", "Bf5", "Nc5", "Rd1", "Rxg2+", "Ke1", "d6", "Bh7", "Rc7", "Qxg5", "Ng8", "Qf6", "Rg1+", "Kf2", "Rcg7", "Bxd6", "Rd7", "Rc1", "Nd3+", "cxd3", "Rh1", "Qf4", "Rxd6", "Nxd6+", "Kd7", "d4", "Kc7", "Qf7+", "Kb8", "Bf5", "Bxf5", "exf5", "Nh6", "Qe7", "Ka8", "Nb7", "Rh4", "Na5", "Nxf5", "Rd1", "Kb8", "Rb1", "Nh6", "Qh7", "Nf5", "Kg1", "b5", "Qa7+", "Kxa7", "d5", "Rg4+", "Kf1", "Nh6", "Nb7", "Rc4", "Kf2", "Ng8", "Kg1", "Rc2", "Rf1", "Rc1", "Kh1", "cxd5", "Na5", "Kb6", "a4", "Rb1", "axb5", "Ra1", "Rb1", "Nf6", "Rc1", "Ng8", "Kg1", "Rxa5", "Rc2", "Ra3", "Re2", "Ra4", "Re5", "Ra2", "Re7", "Kxb5", "Kh1", "Nh6", "Rb7+", "Kc4", "Rg7", "Kc3", "Rd7", "Kb3", "Rb7", "Ka3", "Ra7", "Ra1+", "Kg2", "Ka4", "Rf7", "Nf5", "b5", "Ra2+", "Kf1", "Nd6", "Rf5", "Rf2+", "Ke1", "axb5", "Rf3", "Rb2", "Rh3", "Ne4", "Kf1", "Rc2", "Rh4", "Kb3", "Rh3+", "Ka2", "Rd3", "Rc6", "Rd1", "Ng5", "Kg2", "Nh7", "Re1", "Rf6", "Rc1", "Ka3", "Rc4", "Ka2", "Rb4", "Ka3", "Rf4", "Ka2", "Rxf6", "Ng5", "Ra6+", "Kb3", "Kg3", "Kb4", "Kg2", "Kb3", "Kh2", "Ne4", "Rg6", "Ka4", "Rg1", "Nf2", "Re1", "Ka3", "Re6", "Kb4", "Rg6", "Nh3", "Rg4+", "Kc5", "Ra4", "Nf4", "Re4", "b4", "Rd4", "Ne6", "Kg2", "Kd6", "Rg4", "b3", "Rd4", "Nd8", "Kf3", "Nf7", "Kg2", "Ke7", "Kh2", "Nh8", "Rxd5", "Kf6", "Rd1", "Nf7", "Rb1", "Nh8", "Rd1", "b2", "Kg3", "b1=B", "Rh1", "Kg5", "Rd1", "Nf7", "Rd2", "Bf5", "Rb2", "Nd6", "Rb1", "Kh6", "Rb2", "Bc2", "Rb5", "Bf5", "Kf4", "Bh7", "Rb2", "Kg6", "Kg4", "Nf5", "Rg2", "Ng3", "Rc2", "Bg8", "Rd2", "Bd5", "Rd3", "Nh5", "Rc3", "Bb3", "Kh3", "Ba2", "Rc2", "Ng7", "Rd2", "Kh6", "Rc2", "Kh5", "Rf2", "Kg6", "Re2", "Be6+", "Kh2", "Kf5", "Rd2", "Kf6", "Rb2", "Kg6", "Rc2", "Kf6", "Rd2", "Kg6", "Rd4", "Kh6", "Re4", "Bh3", "Re8", "Be6", "Re7", "Bb3", "Re1", "Ba4", "Kh3", "Bc2", "Rh1", "Kg5", "Kg2", "Kg4", "Rc1", "Bb3", "Rc5", "Bf7", "Rg5+", "Kf4", "Kh1", "Nh5", "Rg3", "Be8", "Rc3", "Bd7", "Rc1", "Be6", "Rc4+", "Kf5", "Ra4", "Ng7", "Ra7", "Kg6", "Ra1", "Kf5", "Kg2", "Kg5", "Rf1", "Nf5", "Kh1", "Nh4", "Kg1", "Ng2", "Kxg2", "Bb3", "Rf7", "Kh6", "Rg7", "Ba4", "Rf7", "Kg6", "Kf3", "Be8", "Ra7", "Bb5", "Ra2", "Be8", "Kf4", "Kh7", "Ra6", "Bb5", "Ke5", "Bc6", "Kd4", "Be4", "Ra2", "Bg2", "Rb2", "Kg8", "Rb6", "Kg7", "Kc4", "Kf7", "Rc6", "Ke8", "Kb3", "Bf3", "Rc7", "Bg4", "Kc2", "Be2", "Rc6", "Bc4", "Rh6", "Bb3+", "Kb1", "Bd5", "Rg6", "Ba8", "Kc1", "Bf3", "Ra6", "Bb7", "Ra7", "Ba8", "Ra3", "Bb7", "Kc2", "Bh1", "Kc3", "Kf7", "Kb4", "Kg7", "Rc3", "Kh6", "Kb3", "Kg7", "Rh3", "Kf7", "Rxh1", "Ke7", "Ka3", "Kf8", "Rh5", "Ke8", "Kb2", "Kd8", "Rh3", "Kc8", "Rh7", "Kd8", "Kc2", "Ke8", "Rb7", "Kd8", "Kd3", "Kc8", "Rb8+", "Kd7", "Rb6", "Ke8", "Rb4", "Kd8", "Ke3", "Ke7", "Rd4", "Kf8", "Rd6", "Ke7", "Rd7+", "Ke8", "Ra7", "Kf8", "Kf3", "Kg8", "Ra4", "Kg7", "Re4", "Kh8", "Ke2", "Kg7", "Kf2", "Kh7", "Rh4+", "Kg6", "Rh7", "Kf6", "Rb7", "Kg6", "Kg1", "Kg5", "Ra7", "Kg6", "Rb7", "Kh5", "Rb3", "Kg4", "Rb8", "Kf3", "Rh8", "Kf4", "Re8", "Kg5", "Rf8", "Kh4", "Rf4+", "Kh5", "Rf2", "Kg4", "Kh1", "Kh5", "Rf4", "Kg5", "Rd4", "Kg6", "Kg2", "Kg5", "Kg1", "Kg6", "Rd6+", "Kg7", "Rh6", "Kxh6"},
		},
		{
			"Random-game-#1077_half-moves-724",
			[]string{"b4", "a5", "c3", "b6", "e4", "Ra7", "g3", "b5", "Qg4", "c5", "Bh3", "h6", "Qe2", "g5", "Bxd7+", "Bxd7", "c4", "g4", "a4", "Bf5", "bxc5", "f6", "d4", "Bd7", "axb5", "e6", "Nh3", "Ra6", "Qf1", "Rc6", "f4", "Qe7", "Nd2", "Kd8", "Ba3", "Rb6", "Qf2", "Kc7", "Nb1", "Qg7", "Qe3", "Na6", "Ng5", "Qg6", "Kf2", "Qxe4", "Qe2", "Bd6", "d5", "h5", "Kf1", "Qf3+", "Qxf3", "Be5", "Ke2", "Kc8", "Kf2", "Bb8", "d6", "f5", "c6", "Nf6", "Qc3", "Rh7", "Ke2", "Rf7", "Qd2", "Re7", "Qc1", "Rxb5", "h4", "Rb2+", "Kd3", "Ne8", "Qc3", "Nc5+", "Ke3", "Na6", "Bxb2", "Bxc6", "d7+", "Kc7", "Rd1", "Bb5", "Nh7", "Nd6", "cxb5+", "Nc5", "Ra4", "Nc8", "Qc4", "Kb6", "dxc8=R", "Nd3", "Qc2", "Nc5", "Kd2", "Nb7", "Qd3", "Rf7", "Kc3", "Rd7", "Qe3+", "Rd4", "Bc1", "Ka7", "Ba3", "Nc5", "Rd8", "Nd7", "Bb2", "Nc5", "Ra3", "Rxd8", "b6+", "Ka6", "Rh1", "Nd3", "Qd2", "Rc8+", "Kb3", "Kb7", "Re1", "Rc6", "Ng5", "Rc1", "Ba1", "Rc7", "bxc7", "Kb6", "Bh8", "Ka7", "Rd1", "Kb7", "Qe2", "Ba7", "Qh2", "Ka8", "Ra1", "Ne5", "Ka4", "Bg1", "Qc2", "Nd7", "Qc4", "Ka7", "Qd3", "Bb6", "c8=N+", "Kb7", "Nd6+", "Kc6", "Qa6", "Nc5+", "Ka3", "Nb7", "Kb3", "e5", "Rxa5", "e4", "Re5", "Nc5+", "Kc4", "Kc7", "Rg1", "Nb7", "Nb5+", "Kd8", "Qa7", "Ba5", "Qd4+", "Nd6+", "Kc5", "Bb6+", "Kc6", "Kc8", "Qd1", "Bxg1", "Na7+", "Kd8", "Qxg1", "Nc4", "Na3", "e3", "Nf3", "Nd6", "Nd2", "e2", "Qc5", "Nf7", "Re4", "e1=R", "Nc2", "Re3", "Nc4", "Rb3", "Qb6+", "Rxb6+", "Kc5", "Nh6", "Rd4+", "Rd6", "Rd5", "Kc7", "Rd3", "Rd4", "Nb6", "Rd6", "Bd4", "Re6", "Be3", "Re7", "Nbc8", "Re8", "Kb4", "Rxc8", "Na1", "Nf7", "Rc3+", "Kd6", "Kb3", "Ke6", "Ka2", "Ne5", "Rc7", "Nd7", "Bd4", "Ra8", "Bg1", "Rh8", "Rb7", "Kf7", "Bc5", "Rh7", "Nc8", "Kg6", "Rb8", "Nf8", "Rb6+", "Kg7", "Rb1", "Kf6", "Ne7", "Rf7", "Ba7", "Nd7", "Nc8", "Rf8", "Ne7", "Rh8", "Be3", "Ne5", "Bc1", "Nd7", "Kb2", "Rd8", "Ng8+", "Rxg8", "Be3", "Kg7", "Re1", "Nf6", "Ka3", "Rf8", "Nc2", "Ne4", "Bb6", "Kh6", "Rxe4", "Kg7", "Bg1", "Rf7", "Re7", "Kg8", "Re1", "Kg7", "Bb6", "Re7", "Ka2", "Re2", "Ba5", "Rh2", "Rc1", "Kh8", "Bc7", "Rh3", "Bd8", "Kh7", "Nb4", "Rh2+", "Rc2", "Rf2", "Be7", "Kg6", "Bf8", "Rxf4", "Kb1", "Rf2", "Be7", "Rh2", "Bf6", "f4", "Rc8", "Rf2", "Rc6", "Ra2", "Nxa2", "Kf5", "Bh8", "Ke4", "Rb6", "Kf5", "Rb7", "f3", "Ka1", "f2", "Nc3", "Kg6", "Ne2", "Kh6", "Bg7+", "Kh7", "Rb8", "f1=N", "Bb2", "Ne3", "Ba3", "Nd1", "Rb1", "Nb2", "Nc1", "Kg6", "Bf8", "Nc4", "Bb4", "Na3", "Na2", "Kg7", "Ba5", "Kh6", "Rb4", "Kg6", "Nc1", "Nc2+", "Kb2", "Na3", "Kc3", "Kg7", "Nb3", "Nc4", "Nc1", "Ne5", "Bd8", "Kh7", "Bc7", "Kg6", "Bb8", "Nd7", "Rc4", "Kf7", "Kd4", "Nf8", "Nb3", "Kg6", "Ke5", "Nd7+", "Kd4", "Kh7", "Rc1", "Kg6", "Bf4", "Kf7", "Kd5", "Nf8", "Rc2", "Kf6", "Rc3", "Kf5", "Be3", "Ne6", "Bh6", "Nf8", "Bg7", "Kg6", "Kc4", "Kxg7", "Na5", "Kh6", "Nb7", "Nd7", "Kb3", "Ne5", "Ka3", "Nf3", "Ka4", "Ng5", "Rc1", "Kg7", "Re1", "Kg8", "Re4", "Kf7", "Re1", "Nh7", "Rf1+", "Nf6", "Kb5", "Kf8", "Rh1", "Nd7", "Nd8", "Ne5", "Nb7", "Nd7", "Rd1", "Kg7", "Re1", "Kf7", "Nd8+", "Kg6", "Re8", "Nc5", "Re4", "Kf6", "Ka5", "Kg7", "Kb5", "Kf8", "Ra4", "Nd3", "Ra8", "Ke7", "Ra2", "Kxd8", "Rh2", "Ne5", "Ra2", "Nd3", "Kc6", "Nc5", "Kd5", "Kc7", "Ra6", "Nxa6", "Kc4", "Kd8", "Kb3", "Nc7", "Kb4", "Ke8", "Kb3", "Nb5", "Ka2", "Nc3+", "Ka3", "Nd1", "Kb3", "Nf2", "Kb4", "Kf8", "Ka3", "Nd1", "Ka2", "Kf7", "Kb1", "Kf8", "Kc2", "Nc3", "Kc1", "Kg7", "Kc2", "Nd5", "Kd1", "Kg6", "Kd2", "Nf4", "Kc2", "Kh7", "Kd2", "Kh6", "Ke3", "Kg7", "Kd4", "Kf8", "Ke3", "Nh3", "Kd2", "Kg7", "Kc1", "Nf2", "Kb1", "Kf7", "Kc1", "Nh3", "Kb1", "Ke8", "Ka2", "Kd7", "Kb1", "Kd8", "Kc2", "Ng1", "Kc3", "Ke8", "Kc4", "Nh3", "Kb5", "Nf4", "gxf4", "Kf8", "Ka4", "Kf7", "Ka3", "Kg8", "Kb2", "Kf7", "Kc2", "Kg7", "Kd2", "Kf7", "f5", "Kf6", "Ke1", "Ke5", "Kd1", "g3", "Ke2", "Ke4", "Ke1", "Ke3", "Kd1", "Kd4", "Kc2", "Ke5", "Kb3", "Kd5", "Ka2", "Ke5", "Ka1", "g2", "Kb1", "Kf4", "Kc2", "g1=Q", "Kd3", "Kf3", "Kc2", "Qf1", "f6", "Qa1", "f7", "Kg2", "f8=R", "Qa3", "Kd1", "Kh3", "Rc8", "Kg2", "Rf8", "Qa7", "Rf5", "Qg7", "Rf1", "Qb7", "Rh1", "Qh7", "Kd2", "Kf3", "Ra1", "Qh8", "Kc2", "Qb2+", "Kxb2", "Ke3", "Rc1", "Kd4", "Rg1", "Kc5", "Rg3", "Kd6", "Kb1", "Kc5", "Rg6", "Kb4", "Rg7", "Ka3", "Rg8", "Kb3", "Kc1", "Kb4", "Rg7", "Ka3", "Ra7+", "Kb3", "Kd1", "Kc4", "Kd2", "Kd4", "Rg7", "Ke4", "Kc3", "Kf3", "Kd3", "Kf2", "Rb7", "Ke1", "Kc4", "Ke2", "Rb8", "Kf3", "Ra8", "Ke4", "Rf8", "Ke5", "Ra8", "Kf5", "Ra1", "Ke6", "Ra7", "Ke5", "Rb7", "Ke6", "Re7+", "Kxe7", "Kc5", "Kf8", "Kb4", "Ke8", "Ka5", "Kf7", "Kb5", "Kg7", "Kc5", "Kh7", "Kd5", "Kh8", "Ke6", "Kg8", "Kf6", "Kf8", "Kg6", "Kg8", "Kf5", "Kh7", "Ke5", "Kg6", "Ke6", "Kh7", "Kd5", "Kg6", "Ke5", "Kh7", "Ke4", "Kh8", "Kd5", "Kg7", "Ke6", "Kf8", "Kd5", "Kf7", "Ke5", "Ke7", "Kd5", "Ke8", "Kd4", "Kd8", "Ke4", "Kc7", "Kd4", "Kb8", "Kc3", "Kc7", "Kb4", "Kb6", "Ka3", "Kc7", "Kb3", "Kd6", "Kc4", "Kc7", "Kb4", "Kd6", "Kb3", "Kc5", "Ka3", "Kc4", "Ka4", "Kd4", "Ka5", "Kd5", "Ka4", "Ke6", "Kb3", "Ke5", "Ka4", "Kd5", "Kb3", "Kd6", "Kc3", "Kc6", "Kc2", "Kd6", "Kc3", "Ke5", "Kd3", "Kd6", "Ke3", "Kd7", "Kf3", "Kc6", "Kf2", "Kd6", "Ke2", "Kc7", "Kd1", "Kb6", "Kc1", "Ka6", "Kd1", "Ka7", "Kc2", "Kb7", "Kc1", "Kc8"},
		},
	}
	for i, bc := range benchCases {
		if testing.Short() && i%2 == 0 {
			// Skip every even benchmark if short flag.
			continue
		}
		b.Run(bc.Name, func(b *testing.B) {
			// Check if benchmark game is valid.
			g := New()
			for n := 0; n < len(bc.PCNMoves); n++ {
				m, err := g.Position().ParseMove(bc.PCNMoves[n])
				if err != nil {
					b.Fatalf("Error parsing move #%d: %v", n+1, err)
				}
				_, err = g.MakeMove(m)
				if err != nil {
					b.Fatalf("Error making move #%d: %v", n+1, err)
				}
			}
			if g.Status() == InProgress {
				b.Fatal("Game status after last benchmark move should not be InProgress")
			}
			b.Logf("GameStatus after %d half-moves: %v", len(g.Positions)-1, g.Status())

			// Generate slices for PCN and struct moves benchmarks.
			PCNMoves := make([]string, 0, len(bc.PCNMoves))
			StructMoves := make([]move.Move, 0, len(bc.PCNMoves))
			for _, p := range g.Positions[1:] {
				PCNMoves = append(PCNMoves, p.LastMove.String())
				StructMoves = append(StructMoves, p.LastMove)
			}

			b.Run("SAN-Moves", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					g := New()
					for n := 0; n < len(bc.PCNMoves); n++ {
						m, err := g.Position().ParseMove(bc.PCNMoves[n])
						if err != nil {
							b.Fatalf("Error parsing move #%d: %v", n+1, err)
						}
						_, err = g.MakeMove(m)
						if err != nil {
							b.Fatalf("Error making move #%d: %v", n+1, err)
						}
					}
				}
			})

			if testing.Short() {
				// Skip other benchmarks if short flag.
				return
			}

			b.Run("PCN-Moves", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					g := New()
					for n := 0; n < len(PCNMoves); n++ {
						_, err := g.MakeMove(move.Parse(PCNMoves[n]))
						if err != nil {
							b.Fatalf("Error making move #%d: %v", n+1, err)
						}
					}
				}
			})

			b.Run("Struct-Moves", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					g := New()
					for n := 0; n < len(StructMoves); n++ {
						_, err := g.MakeMove(StructMoves[n])
						if err != nil {
							b.Fatalf("Error making move #%d: %v", n+1, err)
						}
					}
				}
			})
		})
	}
}
