package polyglot

import (
	"github.com/andrewbackes/chess/fen"
	"github.com/andrewbackes/chess/piece"
	"testing"
)

func TestIndexToFR(t *testing.T) {
	tests := [][]int{
		[]int{0, 7, 0},
		[]int{7, 0, 0},
		[]int{63, 0, 7},
		[]int{56, 7, 7},
	}
	for _, test := range tests {
		f, r := indexToFR(test[0])
		if f != test[1] || r != test[2] {
			t.Error(test, f, r)
		}
	}
}

func TestPieceToPG(t *testing.T) {
	p := piece.New(piece.White, piece.Knight)
	if pieceToPG(p) != 3 {
		t.Fail()
	}
	p = piece.New(piece.Black, piece.Queen)
	if pieceToPG(p) != 8 {
		t.Fail()
	}
	p = piece.New(piece.Black, piece.King)
	if pieceToPG(p) != 10 {
		t.Fail()
	}
	p = piece.New(piece.White, piece.King)
	if pieceToPG(p) != 11 {
		t.Fail()
	}
}

func TestPolyglotHash(t *testing.T) {
	type test struct {
		FEN string
		key uint64
	}

	tests := []test{
		test{
			FEN: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			key: 0x463b96181691fc9c,
		},

		test{
			FEN: "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			key: 0x823c9b50fd114196,
		},

		test{
			FEN: "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2",
			key: 0x0756b94461c50fb0,
		},

		test{
			FEN: "rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 2",
			key: 0x662fafb965db29d4,
		},

		test{
			FEN: "rnbqkbnr/ppp1p1pp/8/3pPp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 3",
			key: 0x22a48b5a8e47ff78,
		},

		test{
			FEN: "rnbqkbnr/ppp1p1pp/8/3pPp2/8/8/PPPPKPPP/RNBQ1BNR b kq - 0 3",
			key: 0x652a607ca3f242c1,
		},

		test{
			FEN: "rnbq1bnr/ppp1pkpp/8/3pPp2/8/8/PPPPKPPP/RNBQ1BNR w - - 0 4",
			key: 0x00fdd303c946bdd9,
		},

		test{
			FEN: "rnbqkbnr/p1pppppp/8/8/PpP4P/8/1P1PPPP1/RNBQKBNR b KQkq c3 0 3",
			key: 0x3c8123ea7b067637,
		},

		test{
			FEN: "rnbqkbnr/p1pppppp/8/8/P6P/R1p5/1P1PPPP1/1NBQKBNR b Kkq - 0 4",
			key: 0x5c3f9b829b279560,
		},
	}

	for i, p := range tests {
		pos, err := fen.Decode(p.FEN)
		if err != nil {
			t.Error(err)
		}
		//got := fmt.Sprintf("%s", g.Polyglot())
		got := Encode(pos)
		if got != Hash(p.key) {
			t.Log("test index:", i)
			t.Log("got", got, "wanted", p.key)
			t.Fail()
		}
	}

}
