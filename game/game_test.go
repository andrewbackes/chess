package game

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestNewTimedGame(t *testing.T) {
	standard := TimeControl{
		Time:  40 * time.Minute,
		Moves: 40,
	}
	control := [2]TimeControl{standard, standard}
	NewTimedGame(control)
}

func TestNonexistentMove(t *testing.T) {
	g := NewGame()
	mv := Move("e4e5")
	status := g.MakeMove(mv)
	if status != WhiteIllegalMove {
		t.Error("Got: ", status, " Wanted: ", WhiteIllegalMove)
	}
}

func TestPlayerToMove(t *testing.T) {
	g := NewGame()
	if g.PlayerToMove() != White {
		t.Error("it's white to move")
	}
	g.MakeMove("e2e4")
	if g.PlayerToMove() != Black {
		t.Error("it's black to move")
	}
}

func TestIllegalCheck(t *testing.T) {

}

func TestIllegalCastle(t *testing.T) {
	g, err := FromFEN("4k3/8/8/8/6r1/8/8/R3K2R w KQ - 0 1")
	s := g.MakeMove(Move("e1g1"))
	if err != nil || s != WhiteIllegalMove {
		t.Fail()
	}
}

func TestThreeFold(t *testing.T) {
	moves := []string{"Nf3", "d6", "d4", "g6", "c4", "Bg7", "Nc3", "Nf6", "e4", "O-O", "Bd3", "Na6", "a3", "c5", "d5", "e6", "O-O", "exd5", "cxd5", "Nc7", "Be3", "Bg4", "h3", "Bxf3", "Qxf3", "Nd7", "Bf4", "Ne5", "Bxe5", "Bxe5", "Rfe1", "a6", "Qd1", "b5", "Qd2", "Qh4", "Ne2", "f5", "f4", "fxe4", "Bxe4", "Bxf4", "Qd3", "Be5", "Rf1", "c4", "Qc2", "Rae8", "Rae1", "Rxf1+", "Rxf1", "Bxb2", "Rf4", "Qe1+", "Rf1", "Qh4", "Rf4", "Qe1+", "Rf1", "Qh4"}
	g := NewGame()
	err := playTestGame(t, g, moves, Threefold)
	if err != nil {
		t.Error(err)
	}
}

func playTestGame(t *testing.T, g *Game, moves []string, expected GameStatus) error {
	for i, san := range moves {
		move, err := g.ParseMove(san)
		if err != nil {
			return err
		}
		s := g.MakeMove(move)
		if (s != InProgress && i+1 < len(moves)) || (i+1 >= len(moves) && s != expected) {
			for _, fen := range g.history.fen {
				t.Log(fen)
			}
			return errors.New(fmt.Sprint("half-move ", i, " (", san, ") ended with status ", s))
		}
	}
	return nil
}

func TestStalemate(t *testing.T) {
	fen := "K7/8/k7/1r6/8/8/8/8 w - - 0 1"
	g, _ := FromFEN(fen)
	if g.gameStatus() != Stalemate {
		t.Fail()
	}
}

func TestTimedOut(t *testing.T) {
	tc := TimeControl{
		Time:  40 * time.Minute,
		Moves: 40,
	}
	g := NewTimedGame([2]TimeControl{tc, tc})
	s := g.MakeTimedMove(Move("e2e4"), 41*time.Minute)
	if s != WhiteTimedOut {
		t.Fail()
	}
}

func TestInsufMaterial(t *testing.T) {

}

func TestFiftyMoveRule(t *testing.T) {
	fen := "8/8/2B2k2/8/3r1NKp/3N4/8/8 b - - 0 62"
	g, _ := FromFEN(fen)
	moves := []string{"Rd8", "Kxh4", "Rg8", "Be4", "Rg1", "Nh5+", "Ke6", "Ng3", "Kf6", "Kg4", "Ra1", "Bd5", "Ra5", "Bf3", "Ra1", "Kf4", "Ke6", "Nc5+", "Kd6", "Nge4+", "Ke7", "Ke5", "Rf1", "Bg4", "Rg1", "Be6", "Re1", "Bc8", "Rc1", "Kd4", "Rd1", "Nd3", "Kf7", "Ke3", "Ra1", "Kf4", "Ke7", "Nb4", "Rc1", "Nd5+", "Kf7", "Bd7", "Rf1", "Ke5", "Ra1", "Ng5+", "Kg6", "Nf3", "Kg7", "Bg4", "Kg6", "Nf4+", "Kg7", "Nd4", "Re1", "Kf5", "Rc1", "Be2", "Re1", "Bh5", "Ra1", "Nfe6+", "Kh6", "Be8", "Ra8", "Bc6", "Ra1", "Kf6", "Kh7", "Ng5+", "Kh8", "Nde6", "Ra6", "Be8", "Ra8", "Bh5", "Ra1", "Bg6", "Rf1", "Ke7", "Ra1", "Nf7+", "Kg8", "Nh6+", "Kh8", "Nf5", "Ra7", "Kf6", "Ra1", "Ne3", "Re1", "Nd5", "Rg1", "Bf5", "Rf1", "Ndf4", "Ra1", "Ng6+", "Kg8", "Ne7+", "Kh8", "Ng6+"}
	err := playTestGame(t, g, moves, FiftyMoveRule)
	if err != nil {
		t.Error(err)
	}
}

func TestEnPassantMove(t *testing.T) {
	fen := "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
	g, _ := FromFEN(fen)
	g.QuickMove(Move("e2c4"))
	g.QuickMove(Move("c7c5"))
	moves := g.LegalMoves()
	if _, ok := moves["d5c6"]; !ok {
		t.Error("missing legal en passant d5c6")
	}
	g.QuickMove(Move("d5c6"))
	if g.board.OnSquare(C5).Type != None {
		t.Error("en passant pawn not captured")
	}
}
