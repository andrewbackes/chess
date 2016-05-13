package game

import (
	"strings"
	"testing"
)

func TestLoadRootPos(t *testing.T) {
	root := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	_, err := FromFEN(root)
	if err != nil {
		t.Fail()
	}
}

// integration test
func TestFENEncodeDecode(t *testing.T) {
	fens := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"4k3/8/8/8/8/8/4P3/4K3 w - - 5 39",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2",
	}
	for _, fen := range fens {
		t.Log("In:  ", fen)
		g, err := FromFEN(fen)
		out := g.FEN()
		t.Log("Out: ", out)
		if err != nil || fen != out {
			t.Error("Do not match.")
		}
	}
}

func TestFENenPassant(t *testing.T) {
	fen := "rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2"
	g, _ := FromFEN(fen)
	if *g.history.enPassant != C6 {
		t.Fail()
	}
}

func TestFENCastlingRights(t *testing.T) {
	fen := "rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2"
	g, _ := FromFEN(fen)
	if !g.history.castlingRights[White][shortSide] || !g.history.castlingRights[Black][shortSide] ||
		!g.history.castlingRights[White][longSide] || !g.history.castlingRights[Black][longSide] {
		t.Fail()
	}
}

func TestFENWhitesMove(t *testing.T) {
	g := NewGame()
	fen := g.FEN()
	player := strings.Split(fen, " ")[1]
	if player != "w" {
		t.Fail()
	}
}

func TestFENBlacksMove(t *testing.T) {
	g := NewGame()
	g.MakeMove(Move("e2e4"))
	fen := g.FEN()
	player := strings.Split(fen, " ")[1]
	if player != "b" {
		t.Fail()
	}
}
