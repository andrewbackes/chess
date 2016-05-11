package game

import (
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
	for i, san := range moves {
		move, err := g.ParseMove(san)
		if err != nil {
			t.Error(err)
		}
		s := g.MakeMove(move)
		if (s != InProgress && i+1 < len(moves)) || (i+1 >= len(moves) && s != Threefold) {
			for _, fen := range g.history.fen {
				t.Log(fen)
			}
			t.Error("half-move", i, "(", san, ") ended with status", s)
		}
	}
}

func TestStalemate(t *testing.T) {

}

func TestTimedOut(t *testing.T) {

}

func TestInsufMaterial(t *testing.T) {

}

func TestFiftyMoveRule(t *testing.T) {

}
