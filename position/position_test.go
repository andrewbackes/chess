package position

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/board"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
)

func TestMailBox(t *testing.T) {
	testCases := []struct {
		name     string
		position testPosition
		want     string
	}{
		{"Empty", testPosition{}, "                                                                "},
		{"Initial", nil, "RNBKQBNRPPPPPPPP                                pppppppprnbkqbnr"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.position.Position()
			mbx := p.MailBox()
			if mbx != tc.want {
				t.Logf("Position:\n%v", p)
				t.Errorf("*Position.MailBox() =\n\t\"%s\"\n, want\n\t\"%s\"", mbx, tc.want)
			}
		})
	}
}

func piecesOnSquare(b *Position, s square.Square) int {
	count := 0
	for c := piece.White; c <= piece.Black; c++ {
		for p := piece.Pawn; p <= piece.King; p++ {
			if (b.bitBoard[c][p] & (1 << s)) != 0 {
				count++
			}
		}
	}
	return count
}

func changedBitBoards(before, after *Position) map[piece.Piece]struct{} {
	changed := make(map[piece.Piece]struct{})

	for c := range before.bitBoard {
		for p := range before.bitBoard[piece.Color(c)] {
			if before.bitBoard[piece.Color(c)][p] != after.bitBoard[piece.Color(c)][p] {
				changed[piece.New(piece.Color(c), piece.Type(p))] = struct{}{}
			}
		}
	}
	return changed
}

func TestMovePawn(t *testing.T) {
	beforeMove := New()
	afterMove := beforeMove.MakeMove(move.Parse("e2e4"))
	changed := changedBitBoards(beforeMove, afterMove)
	t.Log("Changed: ", changed)
	if _, c := changed[piece.New(piece.White, piece.Pawn)]; !c || len(changed) != 1 {
		t.Fail()
	}
}

func TestMoveKnight(t *testing.T) {
	beforeMove := New()
	afterMove := beforeMove.MakeMove(move.Parse("b1c3"))
	changed := changedBitBoards(beforeMove, afterMove)
	t.Log("Changed: ", changed)
	if _, c := changed[piece.New(piece.White, piece.Knight)]; !c || len(changed) != 1 {
		t.Fail()
	}
}

// Stores information about various position changes.
type positionChanges struct {
	squares   squareChanges
	castling  castlingChanges
	enPassant *square.Square
}

// Stores information about changes on a chess board for testing purposes. It considers, that here can be multiple piece types on one square.
// Values are:
// 	false - Piece was removed from square.
// 	true  - Piece was added to square.
type squareChanges map[square.Square]pieceChanges
type pieceChanges map[piece.Piece]bool

// Stores information about castling rights changes.
// Values are:
// 	false - Castling right for color and side was changed from true to false
// 	true  - Castling right for color and side was changed from false to true (should never happen).
type castlingChanges map[piece.Color]map[board.Side]bool

func squarePtr(sq square.Square) *square.Square {
	return &sq
}

// Returns board changes between two positions.
func changedSquares(before, after *Position) squareChanges {
	changed := make(squareChanges)
	for sq := square.Square(0); sq <= square.LastSquare; sq += 1 {
		for _, col := range piece.Colors {
			for tpe := piece.Pawn; tpe <= piece.King; tpe += 1 {
				sqb := uint64(1 << sq)
				bsqpc := (before.bitBoard[col][tpe] & sqb) > 0
				asqpc := (after.bitBoard[col][tpe] & sqb) > 0
				if bsqpc != asqpc {
					if _, ok := changed[sq]; !ok {
						changed[sq] = pieceChanges{}
					}
					changed[sq][piece.New(col, tpe)] = asqpc && !bsqpc
				}
			}
		}
	}

	return changed
}

// Returns castling changes between two positions.
func changedCastlingRights(before, after *Position) castlingChanges {
	changed := make(castlingChanges)
	for _, col := range piece.Colors {
		for _, side := range board.Sides {
			if before.CastlingRights[col][side] != after.CastlingRights[col][side] {
				if _, ok := changed[col]; !ok {
					changed[col] = map[board.Side]bool{}
				}
				changed[col][side] = after.CastlingRights[col][side]
			}
		}
	}

	if len(changed) == 0 {
		return nil
	}
	return changed
}

func TestPositionMakeMove(t *testing.T) {
	testCases := []struct {
		name        string
		p           testPosition
		pc          positionChanger
		m           move.Move
		wantChanged positionChanges
	}{
		// Move.
		{"InitialBoard-ActiveWhite-WhitePawn-e2e3",
			InitialTestPosition, nil,
			move.Move{square.E2, square.E3, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.E2: pieceChanges{piece.New(piece.White, piece.Pawn): false},
					square.E3: pieceChanges{piece.New(piece.White, piece.Pawn): true},
				},
			},
		},
		{"InitialBoard-ActiveBlack-BlackPawn-e7e6",
			InitialTestPosition, active(piece.Black),
			move.Move{square.E7, square.E6, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.E7: pieceChanges{piece.New(piece.Black, piece.Pawn): false},
					square.E6: pieceChanges{piece.New(piece.Black, piece.Pawn): true},
				},
			},
		},
		{"InitialBoard-ActiveWhite-WhitePawn-e2e3-Time5",
			InitialTestPosition, nil,
			move.Move{square.E2, square.E3, piece.None, 5},
			positionChanges{
				squares: squareChanges{
					square.E2: pieceChanges{piece.New(piece.White, piece.Pawn): false},
					square.E3: pieceChanges{piece.New(piece.White, piece.Pawn): true},
				},
			},
		},
		{"InitialBoard-ActiveBlack-BlackPawn-e7e6-Time5",
			InitialTestPosition, active(piece.Black),
			move.Move{square.E7, square.E6, piece.None, 5},
			positionChanges{
				squares: squareChanges{
					square.E7: pieceChanges{piece.New(piece.Black, piece.Pawn): false},
					square.E6: pieceChanges{piece.New(piece.Black, piece.Pawn): true},
				},
			},
		},
		{"InitialBoard-ActiveWhite-WhitePawn-e2e4-StoreEnPassant",
			InitialTestPosition, nil,
			move.Move{square.E2, square.E4, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.E2: pieceChanges{piece.New(piece.White, piece.Pawn): false},
					square.E4: pieceChanges{piece.New(piece.White, piece.Pawn): true},
				},
				enPassant: squarePtr(square.E3),
			},
		},
		{"InitialBoard-ActiveBlack-BlackPawn-e7e5-StoreEnPassant",
			InitialTestPosition, active(piece.Black),
			move.Move{square.E7, square.E5, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.E7: pieceChanges{piece.New(piece.Black, piece.Pawn): false},
					square.E5: pieceChanges{piece.New(piece.Black, piece.Pawn): true},
				},
				enPassant: squarePtr(square.E6),
			},
		},
		{"InitialBoard-ActiveWhite-WhiteKnight-b1c3",
			InitialTestPosition, nil,
			move.Move{square.B1, square.C3, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.B1: pieceChanges{piece.New(piece.White, piece.Knight): false},
					square.C3: pieceChanges{piece.New(piece.White, piece.Knight): true},
				},
			},
		},
		{"InitialBoard-ActiveBlack-BlackKnight-b8c6",
			InitialTestPosition, active(piece.Black),
			move.Move{square.B8, square.C6, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.B8: pieceChanges{piece.New(piece.Black, piece.Knight): false},
					square.C6: pieceChanges{piece.New(piece.Black, piece.Knight): true},
				},
			},
		},
		{"TwoPawnsAtTwoOpositeKings-ActiveWhite-WhiteKing-e1f2",
			TwoPawnsAtTwoOpositeKingsTestPosition, nil,
			move.Move{square.E1, square.F2, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.E1: pieceChanges{piece.New(piece.White, piece.King): false},
					square.F2: pieceChanges{piece.New(piece.White, piece.King): true},
				},
				castling: castlingChanges{piece.White: map[board.Side]bool{board.ShortSide: false, board.LongSide: false}},
			},
		},
		{"TwoPawnsAtTwoOpositeKings-ActiveBlack-BlackKing-e8f7",
			TwoPawnsAtTwoOpositeKingsTestPosition, active(piece.Black),
			move.Move{square.E8, square.F7, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.E8: pieceChanges{piece.New(piece.Black, piece.King): false},
					square.F7: pieceChanges{piece.New(piece.Black, piece.King): true},
				},
				castling: castlingChanges{piece.Black: map[board.Side]bool{board.ShortSide: false, board.LongSide: false}},
			},
		},
		// Capture.
		{"TwoPawnsAtTwoOpositeKings-ActiveWhite-WhiteKing-e1e2-CaptureBlackPawn",
			TwoPawnsAtTwoOpositeKingsTestPosition, nil,
			move.Move{square.E1, square.E2, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.E1: pieceChanges{piece.New(piece.White, piece.King): false},
					square.E2: pieceChanges{
						piece.New(piece.Black, piece.Pawn): false,
						piece.New(piece.White, piece.King): true,
					},
				},
				castling: castlingChanges{piece.White: map[board.Side]bool{board.ShortSide: false, board.LongSide: false}},
			},
		},
		{"TwoPawnsAtTwoOpositeKings-ActiveBlack-BlackKing-e8e7-CaptureWhitePawn",
			TwoPawnsAtTwoOpositeKingsTestPosition, active(piece.Black),
			move.Move{square.E8, square.E7, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.E8: pieceChanges{piece.New(piece.Black, piece.King): false},
					square.E7: pieceChanges{
						piece.New(piece.White, piece.Pawn): false,
						piece.New(piece.Black, piece.King): true,
					},
				},
				castling: castlingChanges{piece.Black: map[board.Side]bool{board.ShortSide: false, board.LongSide: false}},
			},
		},
		// Castle.
		{"TwoKingsFourRooks-ActiveWhite-WhiteKing-e1g1-ShortCastle",
			TwoKingsFourRooksTestPosition, nil,
			move.Move{square.E1, square.G1, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.E1: pieceChanges{piece.New(piece.White, piece.King): false},
					square.G1: pieceChanges{piece.New(piece.White, piece.King): true},
					square.H1: pieceChanges{piece.New(piece.White, piece.Rook): false},
					square.F1: pieceChanges{piece.New(piece.White, piece.Rook): true},
				},
				castling: castlingChanges{piece.White: map[board.Side]bool{board.ShortSide: false, board.LongSide: false}},
			},
		},
		{"TwoKingsFourRooks-ActiveBlack-BlackKing-e8g8-ShortCastle",
			TwoKingsFourRooksTestPosition, active(piece.Black),
			move.Move{square.E8, square.G8, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.E8: pieceChanges{piece.New(piece.Black, piece.King): false},
					square.G8: pieceChanges{piece.New(piece.Black, piece.King): true},
					square.H8: pieceChanges{piece.New(piece.Black, piece.Rook): false},
					square.F8: pieceChanges{piece.New(piece.Black, piece.Rook): true},
				},
				castling: castlingChanges{piece.Black: map[board.Side]bool{board.ShortSide: false, board.LongSide: false}},
			},
		},
		{"TwoKingsFourRooks-ActiveWhite-WhiteKing-e1c1-LongCastle",
			TwoKingsFourRooksTestPosition, nil,
			move.Move{square.E1, square.C1, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.E1: pieceChanges{piece.New(piece.White, piece.King): false},
					square.C1: pieceChanges{piece.New(piece.White, piece.King): true},
					square.A1: pieceChanges{piece.New(piece.White, piece.Rook): false},
					square.D1: pieceChanges{piece.New(piece.White, piece.Rook): true},
				},
				castling: castlingChanges{piece.White: map[board.Side]bool{board.ShortSide: false, board.LongSide: false}},
			},
		},
		{"TwoKingsFourRooks-ActiveBlack-BlackKing-e8c8-LongCastle",
			TwoKingsFourRooksTestPosition, active(piece.Black),
			move.Move{square.E8, square.C8, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.E8: pieceChanges{piece.New(piece.Black, piece.King): false},
					square.C8: pieceChanges{piece.New(piece.Black, piece.King): true},
					square.A8: pieceChanges{piece.New(piece.Black, piece.Rook): false},
					square.D8: pieceChanges{piece.New(piece.Black, piece.Rook): true},
				},
				castling: castlingChanges{piece.Black: map[board.Side]bool{board.ShortSide: false, board.LongSide: false}},
			},
		},
		// Castling rights.
		{"TwoKingsFourRooks-ActiveWhite-WhiteRook-a1a2",
			TwoKingsFourRooksTestPosition, nil,
			move.Move{square.A1, square.A2, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.A1: pieceChanges{piece.New(piece.White, piece.Rook): false},
					square.A2: pieceChanges{piece.New(piece.White, piece.Rook): true},
				},
				castling: castlingChanges{piece.White: map[board.Side]bool{board.LongSide: false}},
			},
		},
		{"TwoKingsFourRooks-ActiveBlack-BlackRook-a8a7",
			TwoKingsFourRooksTestPosition, active(piece.Black),
			move.Move{square.A8, square.A7, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.A8: pieceChanges{piece.New(piece.Black, piece.Rook): false},
					square.A7: pieceChanges{piece.New(piece.Black, piece.Rook): true},
				},
				castling: castlingChanges{piece.Black: map[board.Side]bool{board.LongSide: false}},
			},
		},
		{"TwoKingsFourRooks-ActiveWhite-WhiteRook-h1h3",
			TwoKingsFourRooksTestPosition, nil,
			move.Move{square.H1, square.H3, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.H1: pieceChanges{piece.New(piece.White, piece.Rook): false},
					square.H3: pieceChanges{piece.New(piece.White, piece.Rook): true},
				},
				castling: castlingChanges{piece.White: map[board.Side]bool{board.ShortSide: false}},
			},
		},
		{"TwoKingsFourRooks-ActiveBlack-BlackRook-h8h6",
			TwoKingsFourRooksTestPosition, active(piece.Black),
			move.Move{square.H8, square.H6, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.H8: pieceChanges{piece.New(piece.Black, piece.Rook): false},
					square.H6: pieceChanges{piece.New(piece.Black, piece.Rook): true},
				},
				castling: castlingChanges{piece.Black: map[board.Side]bool{board.ShortSide: false}},
			},
		},
		// En-passant.
		{"EnPassantCapture-ActiveWhite-EnPassantOnB6-WhitePawn-b5c6-CaptureBlackPawnEnPassant",
			EnPassantCaptureTestPosition, enPassant(square.C6),
			move.Move{square.B5, square.C6, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.B5: pieceChanges{piece.New(piece.White, piece.Pawn): false},
					square.C6: pieceChanges{piece.New(piece.White, piece.Pawn): true},
					square.C5: pieceChanges{piece.New(piece.Black, piece.Pawn): false},
				},
			},
		},
		{"EnPassantCapture-ActiveBlack-EnPassantOnB6-BlackPawn-g4f3-CaptureWhitePawnEnPassant",
			EnPassantCaptureTestPosition, multi(active(piece.Black), enPassant(square.F3)),
			move.Move{square.G4, square.F3, piece.None, 0},
			positionChanges{
				squares: squareChanges{
					square.G4: pieceChanges{piece.New(piece.Black, piece.Pawn): false},
					square.F3: pieceChanges{piece.New(piece.Black, piece.Pawn): true},
					square.F4: pieceChanges{piece.New(piece.White, piece.Pawn): false},
				},
			},
		},
		// Promotion.
		{"Promotion-ActiveWhite-WhitePawn-b7b8-PromoteToQueen",
			PromotionTestPosition, nil,
			move.Move{square.B7, square.B8, piece.Queen, 0},
			positionChanges{
				squares: squareChanges{
					square.B7: pieceChanges{piece.New(piece.White, piece.Pawn): false},
					square.B8: pieceChanges{piece.New(piece.White, piece.Queen): true},
				},
			},
		},
		{"Promotion-ActiveBlack-BlackPawn-b2b1-PromoteToQueen",
			PromotionTestPosition, active(piece.Black),
			move.Move{square.B2, square.B1, piece.Queen, 0},
			positionChanges{
				squares: squareChanges{
					square.B2: pieceChanges{piece.New(piece.Black, piece.Pawn): false},
					square.B1: pieceChanges{piece.New(piece.Black, piece.Queen): true},
				},
			},
		},
		{"Promotion-ActiveWhite-WhitePawn-b7a8-CaptureAndPromoteToRook",
			PromotionTestPosition, nil,
			move.Move{square.B7, square.A8, piece.Rook, 0},
			positionChanges{
				squares: squareChanges{
					square.B7: pieceChanges{piece.New(piece.White, piece.Pawn): false},
					square.A8: pieceChanges{
						piece.New(piece.White, piece.Rook): true,
						piece.New(piece.Black, piece.Rook): false,
					},
				},
				castling: castlingChanges{piece.Black: map[board.Side]bool{board.LongSide: false}},
			},
		},
		{"Promotion-ActiveBlack-BlackPawn-b2a1-CaptureAndPromoteToRook",
			PromotionTestPosition, active(piece.Black),
			move.Move{square.B2, square.A1, piece.Rook, 0},
			positionChanges{
				squares: squareChanges{
					square.B2: pieceChanges{piece.New(piece.Black, piece.Pawn): false},
					square.A1: pieceChanges{
						piece.New(piece.Black, piece.Rook): true,
						piece.New(piece.White, piece.Rook): false,
					},
				},
				castling: castlingChanges{piece.White: map[board.Side]bool{board.LongSide: false}},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Prepare position.
			beforeMove, err := testCasePosition(tc.p, tc.pc)
			if err != nil {
				t.Fatal(err)
			}

			// Catch panics.
			defer func() {
				if err := recover(); err != nil {
					t.Logf("Position:\n%v", beforeMove)
					t.Errorf("*Position.MakeMove(%v) should not panic, but panicked with: %v", tc.m, err)
				}
			}()

			{ // Check if test move is a legal move.
				legalMoves := beforeMove.LegalMoves()
				moveNoDuration := tc.m
				moveNoDuration.Duration = time.Duration(0)
				if _, ok := legalMoves[moveNoDuration]; !ok {
					t.Logf("Position:\n%v", beforeMove)
					t.Fatalf("Move %v is not in Position.LegalMoves() result: %v", moveNoDuration, legalMoves)
				}
			}

			// Make move.
			afterMove := beforeMove.MakeMove(tc.m)

			{ // Check bitBoards changes.
				changed := changedSquares(beforeMove, afterMove)
				if !reflect.DeepEqual(changed, tc.wantChanged.squares) {
					t.Logf("Position:\n%v", beforeMove)
					t.Errorf("After *Position.MakeMove(%v), board changes are %v, want %v", tc.m, changed, tc.wantChanged.squares)
				}
			}

			{ // Check castling rights changes, if any.
				changed := changedCastlingRights(beforeMove, afterMove)
				if !reflect.DeepEqual(changed, tc.wantChanged.castling) {
					t.Logf("Position:\n%v", beforeMove)
					t.Errorf("After *Position.MakeMove(%v), castling changes are %v, want %v", tc.m, changed, tc.wantChanged.castling)
				}
			}

			{ // Check EnPassant after move.
				want := square.NoSquare
				if tc.wantChanged.enPassant != nil {
					want = *tc.wantChanged.enPassant
				}
				if afterMove.EnPassant != want {
					t.Logf("Position:\n%v", beforeMove)
					t.Errorf("After *Position.MakeMove(%v), EnPassant is %v, want %v", tc.m, afterMove.EnPassant, want)
				}
			}

			// Check ActiveColor after move.
			if afterMove.ActiveColor != (beforeMove.ActiveColor+1)%2 {
				t.Logf("Position:\n%v", beforeMove)
				t.Fatalf("After *Position.MakeMove(%v), board ActiveColor is %v, want %v", tc.m, afterMove.ActiveColor, (beforeMove.ActiveColor+1)%2)
			}

			// Do additional checks only if ActiveColor before move is a valid color.
			if beforeMove.ActiveColor == piece.White || beforeMove.ActiveColor == piece.Black {
				// Check Clocks for both colors after move.
				wantClock := beforeMove.Clocks[beforeMove.ActiveColor] - tc.m.Duration
				if afterMove.Clocks[beforeMove.ActiveColor] != wantClock {
					t.Logf("Position:\n%v", beforeMove)
					t.Errorf("After *Position.MakeMove(%v), Clock[%v] is %v, want %v", tc.m, beforeMove.ActiveColor, afterMove.Clocks[beforeMove.ActiveColor], wantClock)
				}
				if afterMove.Clocks[afterMove.ActiveColor] != beforeMove.Clocks[afterMove.ActiveColor] {
					t.Logf("Position:\n%v", beforeMove)
					t.Errorf("After *Position.MakeMove(%v), Clock[%v] is %v, want %v", tc.m, afterMove.ActiveColor, afterMove.Clocks[afterMove.ActiveColor], beforeMove.Clocks[afterMove.ActiveColor])
				}

				// Check MovesLeft for both colors after move.
				wantML := beforeMove.MovesLeft[beforeMove.ActiveColor] - 1
				if afterMove.MovesLeft[beforeMove.ActiveColor] != wantML {
					t.Logf("Position:\n%v", beforeMove)
					t.Errorf("After *Position.MakeMove(%v), MovesLeft[%v] is %v, want %v", tc.m, beforeMove.ActiveColor, afterMove.MovesLeft[beforeMove.ActiveColor], wantML)
				}
				if afterMove.MovesLeft[afterMove.ActiveColor] != beforeMove.MovesLeft[afterMove.ActiveColor] {
					t.Logf("Position:\n%v", beforeMove)
					t.Errorf("After *Position.MakeMove(%v), MovesLeft[%v] is %v, want %v", tc.m, afterMove.ActiveColor, afterMove.MovesLeft[afterMove.ActiveColor], beforeMove.MovesLeft[afterMove.ActiveColor])
				}

				// Check MoveNumber after move.
				if beforeMove.ActiveColor == piece.Black {
					wantMN := beforeMove.MoveNumber + 1
					if afterMove.MoveNumber != wantMN {
						t.Logf("Position:\n%v", beforeMove)
						t.Errorf("After *Position.MakeMove(%v), MoveNumber is %v, want %v", tc.m, afterMove.MoveNumber, wantMN)
					}
				} else if afterMove.MoveNumber != beforeMove.MoveNumber {
					t.Logf("Position:\n%v", beforeMove)
					t.Errorf("After *Position.MakeMove(%v), MoveNumber is %v, want %v", tc.m, afterMove.MoveNumber, beforeMove.MoveNumber)
				}
			}

			// Check LastMove after move.
			if afterMove.LastMove != tc.m {
				t.Logf("Position:\n%v", beforeMove)
				t.Errorf("After *Position.MakeMove(%v), LastMove is %v, want %v", tc.m, afterMove.LastMove, tc.m)
			}

			// Check ThreeFoldCount update after move.
			hash := afterMove.Polyglot()
			beforeHC := 0
			if hc, exists := beforeMove.ThreeFoldCount[hash]; exists {
				beforeHC = hc
			}
			afterHC := afterMove.ThreeFoldCount[hash]
			if afterHC != beforeHC+1 {
				t.Logf("Position:\n%v", beforeMove)
				t.Errorf("After *Position.MakeMove(%v), ThreeFoldCount[hash] is %v, want %v", tc.m, afterHC, beforeHC+1)
			}
		})
	}
}

func TestPositionPut(t *testing.T) {
	testCases := []struct {
		name        string
		sq          square.Square
		pc          piece.Piece
		wantChanged squareChanges
	}{
		// Pawns.
		{"White-Pawn-A3", square.A3, piece.New(piece.White, piece.Pawn),
			squareChanges{square.A3: pieceChanges{piece.New(piece.White, piece.Pawn): true}},
		},
		{"Black-Pawn-A6", square.A6, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.A6: pieceChanges{piece.New(piece.Black, piece.Pawn): true}},
		},
		{"White-Pawn-A2", square.A2, piece.New(piece.White, piece.Pawn),
			squareChanges{},
		},
		{"Black-Pawn-A7", square.A7, piece.New(piece.Black, piece.Pawn),
			squareChanges{},
		},
		{"White-Pawn-A1", square.A1, piece.New(piece.White, piece.Pawn),
			squareChanges{square.A1: pieceChanges{
				piece.New(piece.White, piece.Rook): false,
				piece.New(piece.White, piece.Pawn): true,
			}},
		},
		{"Black-Pawn-A8", square.A8, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.A8: pieceChanges{
				piece.New(piece.Black, piece.Rook): false,
				piece.New(piece.Black, piece.Pawn): true,
			}},
		},
		{"White-Pawn-B1", square.B1, piece.New(piece.White, piece.Pawn),
			squareChanges{square.B1: pieceChanges{
				piece.New(piece.White, piece.Knight): false,
				piece.New(piece.White, piece.Pawn):   true,
			}},
		},
		{"Black-Pawn-B8", square.B8, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.B8: pieceChanges{
				piece.New(piece.Black, piece.Knight): false,
				piece.New(piece.Black, piece.Pawn):   true,
			}},
		},
		{"White-Pawn-C1", square.C1, piece.New(piece.White, piece.Pawn),
			squareChanges{square.C1: pieceChanges{
				piece.New(piece.White, piece.Bishop): false,
				piece.New(piece.White, piece.Pawn):   true,
			}},
		},
		{"Black-Pawn-C8", square.C8, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.C8: pieceChanges{
				piece.New(piece.Black, piece.Bishop): false,
				piece.New(piece.Black, piece.Pawn):   true,
			}},
		},
		{"White-Pawn-D1", square.D1, piece.New(piece.White, piece.Pawn),
			squareChanges{square.D1: pieceChanges{
				piece.New(piece.White, piece.Queen): false,
				piece.New(piece.White, piece.Pawn):  true,
			}},
		},
		{"Black-Pawn-D8", square.D8, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.D8: pieceChanges{
				piece.New(piece.Black, piece.Queen): false,
				piece.New(piece.Black, piece.Pawn):  true,
			}},
		},
		{"White-Pawn-E1", square.E1, piece.New(piece.White, piece.Pawn),
			squareChanges{square.E1: pieceChanges{
				piece.New(piece.White, piece.King): false,
				piece.New(piece.White, piece.Pawn): true,
			}},
		},
		{"Black-Pawn-E8", square.E8, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.E8: pieceChanges{
				piece.New(piece.Black, piece.King): false,
				piece.New(piece.Black, piece.Pawn): true,
			}},
		},
		{"White-Pawn-H6", square.H6, piece.New(piece.White, piece.Pawn),
			squareChanges{square.H6: pieceChanges{piece.New(piece.White, piece.Pawn): true}},
		},
		{"Black-Pawn-H3", square.H3, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.H3: pieceChanges{piece.New(piece.Black, piece.Pawn): true}},
		},
		{"White-Pawn-H7", square.H7, piece.New(piece.White, piece.Pawn),
			squareChanges{square.H7: pieceChanges{
				piece.New(piece.Black, piece.Pawn): false,
				piece.New(piece.White, piece.Pawn): true,
			}},
		},
		{"Black-Pawn-H2", square.H2, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.H2: pieceChanges{
				piece.New(piece.White, piece.Pawn): false,
				piece.New(piece.Black, piece.Pawn): true,
			}},
		},
		{"White-Pawn-H8", square.H8, piece.New(piece.White, piece.Pawn),
			squareChanges{square.H8: pieceChanges{
				piece.New(piece.Black, piece.Rook): false,
				piece.New(piece.White, piece.Pawn): true,
			}},
		},
		{"Black-Pawn-H1", square.H1, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.H1: pieceChanges{
				piece.New(piece.White, piece.Rook): false,
				piece.New(piece.Black, piece.Pawn): true,
			}},
		},
		{"White-Pawn-G8", square.G8, piece.New(piece.White, piece.Pawn),
			squareChanges{square.G8: pieceChanges{
				piece.New(piece.Black, piece.Knight): false,
				piece.New(piece.White, piece.Pawn):   true,
			}},
		},
		{"Black-Pawn-G1", square.G1, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.G1: pieceChanges{
				piece.New(piece.White, piece.Knight): false,
				piece.New(piece.Black, piece.Pawn):   true,
			}},
		},
		{"White-Pawn-F8", square.F8, piece.New(piece.White, piece.Pawn),
			squareChanges{square.F8: pieceChanges{
				piece.New(piece.Black, piece.Bishop): false,
				piece.New(piece.White, piece.Pawn):   true,
			}},
		},
		{"Black-Pawn-F1", square.F1, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.F1: pieceChanges{
				piece.New(piece.White, piece.Bishop): false,
				piece.New(piece.Black, piece.Pawn):   true,
			}},
		},
		{"White-Pawn-E8", square.E8, piece.New(piece.White, piece.Pawn),
			squareChanges{square.E8: pieceChanges{
				piece.New(piece.Black, piece.King): false,
				piece.New(piece.White, piece.Pawn): true,
			}},
		},
		{"Black-Pawn-E1", square.E1, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.E1: pieceChanges{
				piece.New(piece.White, piece.King): false,
				piece.New(piece.Black, piece.Pawn): true,
			}},
		},
		{"White-Pawn-D8", square.D8, piece.New(piece.White, piece.Pawn),
			squareChanges{square.D8: pieceChanges{
				piece.New(piece.Black, piece.Queen): false,
				piece.New(piece.White, piece.Pawn):  true,
			}},
		},
		{"Black-Pawn-D1", square.D1, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.D1: pieceChanges{
				piece.New(piece.White, piece.Queen): false,
				piece.New(piece.Black, piece.Pawn):  true,
			}},
		},

		// Rooks.
		{"White-Rook-A3", square.A3, piece.New(piece.White, piece.Rook),
			squareChanges{square.A3: pieceChanges{piece.New(piece.White, piece.Rook): true}},
		},
		{"Black-Rook-A6", square.A6, piece.New(piece.Black, piece.Rook),
			squareChanges{square.A6: pieceChanges{piece.New(piece.Black, piece.Rook): true}},
		},
		{"White-Rook-A2", square.A2, piece.New(piece.White, piece.Rook),
			squareChanges{square.A2: pieceChanges{
				piece.New(piece.White, piece.Pawn): false,
				piece.New(piece.White, piece.Rook): true,
			}},
		},
		{"Black-Rook-A7", square.A7, piece.New(piece.Black, piece.Rook),
			squareChanges{square.A7: pieceChanges{
				piece.New(piece.Black, piece.Pawn): false,
				piece.New(piece.Black, piece.Rook): true,
			}},
		},
		{"White-Rook-A1", square.A1, piece.New(piece.White, piece.Rook),
			squareChanges{},
		},
		{"Black-Rook-A8", square.A8, piece.New(piece.Black, piece.Rook),
			squareChanges{},
		},
		{"White-Rook-B1", square.B1, piece.New(piece.White, piece.Rook),
			squareChanges{square.B1: pieceChanges{
				piece.New(piece.White, piece.Knight): false,
				piece.New(piece.White, piece.Rook):   true,
			}},
		},
		{"Black-Rook-B8", square.B8, piece.New(piece.Black, piece.Rook),
			squareChanges{square.B8: pieceChanges{
				piece.New(piece.Black, piece.Knight): false,
				piece.New(piece.Black, piece.Rook):   true,
			}},
		},
		{"White-Rook-C1", square.C1, piece.New(piece.White, piece.Rook),
			squareChanges{square.C1: pieceChanges{
				piece.New(piece.White, piece.Bishop): false,
				piece.New(piece.White, piece.Rook):   true,
			}},
		},
		{"Black-Rook-C8", square.C8, piece.New(piece.Black, piece.Rook),
			squareChanges{square.C8: pieceChanges{
				piece.New(piece.Black, piece.Bishop): false,
				piece.New(piece.Black, piece.Rook):   true,
			}},
		},
		{"White-Rook-D1", square.D1, piece.New(piece.White, piece.Rook),
			squareChanges{square.D1: pieceChanges{
				piece.New(piece.White, piece.Queen): false,
				piece.New(piece.White, piece.Rook):  true,
			}},
		},
		{"Black-Rook-D8", square.D8, piece.New(piece.Black, piece.Rook),
			squareChanges{square.D8: pieceChanges{
				piece.New(piece.Black, piece.Queen): false,
				piece.New(piece.Black, piece.Rook):  true,
			}},
		},
		{"White-Rook-E1", square.E1, piece.New(piece.White, piece.Rook),
			squareChanges{square.E1: pieceChanges{
				piece.New(piece.White, piece.King): false,
				piece.New(piece.White, piece.Rook): true,
			}},
		},
		{"Black-Rook-E8", square.E8, piece.New(piece.Black, piece.Rook),
			squareChanges{square.E8: pieceChanges{
				piece.New(piece.Black, piece.King): false,
				piece.New(piece.Black, piece.Rook): true,
			}},
		},
		{"White-Rook-H6", square.H6, piece.New(piece.White, piece.Rook),
			squareChanges{square.H6: pieceChanges{piece.New(piece.White, piece.Rook): true}},
		},
		{"Black-Rook-H3", square.H3, piece.New(piece.Black, piece.Rook),
			squareChanges{square.H3: pieceChanges{piece.New(piece.Black, piece.Rook): true}},
		},
		{"White-Rook-H7", square.H7, piece.New(piece.White, piece.Rook),
			squareChanges{square.H7: pieceChanges{
				piece.New(piece.Black, piece.Pawn): false,
				piece.New(piece.White, piece.Rook): true,
			}},
		},
		{"Black-Rook-H2", square.H2, piece.New(piece.Black, piece.Rook),
			squareChanges{square.H2: pieceChanges{
				piece.New(piece.White, piece.Pawn): false,
				piece.New(piece.Black, piece.Rook): true,
			}},
		},
		{"White-Rook-H8", square.H8, piece.New(piece.White, piece.Rook),
			squareChanges{square.H8: pieceChanges{
				piece.New(piece.Black, piece.Rook): false,
				piece.New(piece.White, piece.Rook): true,
			}},
		},
		{"Black-Rook-H1", square.H1, piece.New(piece.Black, piece.Rook),
			squareChanges{square.H1: pieceChanges{
				piece.New(piece.White, piece.Rook): false,
				piece.New(piece.Black, piece.Rook): true,
			}},
		},
		{"White-Rook-G8", square.G8, piece.New(piece.White, piece.Rook),
			squareChanges{square.G8: pieceChanges{
				piece.New(piece.Black, piece.Knight): false,
				piece.New(piece.White, piece.Rook):   true,
			}},
		},
		{"Black-Rook-G1", square.G1, piece.New(piece.Black, piece.Rook),
			squareChanges{square.G1: pieceChanges{
				piece.New(piece.White, piece.Knight): false,
				piece.New(piece.Black, piece.Rook):   true,
			}},
		},
		{"White-Rook-F8", square.F8, piece.New(piece.White, piece.Rook),
			squareChanges{square.F8: pieceChanges{
				piece.New(piece.Black, piece.Bishop): false,
				piece.New(piece.White, piece.Rook):   true,
			}},
		},
		{"Black-Rook-F1", square.F1, piece.New(piece.Black, piece.Rook),
			squareChanges{square.F1: pieceChanges{
				piece.New(piece.White, piece.Bishop): false,
				piece.New(piece.Black, piece.Rook):   true,
			}},
		},
		{"White-Rook-E8", square.E8, piece.New(piece.White, piece.Rook),
			squareChanges{square.E8: pieceChanges{
				piece.New(piece.Black, piece.King): false,
				piece.New(piece.White, piece.Rook): true,
			}},
		},
		{"Black-Rook-E1", square.E1, piece.New(piece.Black, piece.Rook),
			squareChanges{square.E1: pieceChanges{
				piece.New(piece.White, piece.King): false,
				piece.New(piece.Black, piece.Rook): true,
			}},
		},
		{"White-Rook-D8", square.D8, piece.New(piece.White, piece.Rook),
			squareChanges{square.D8: pieceChanges{
				piece.New(piece.Black, piece.Queen): false,
				piece.New(piece.White, piece.Rook):  true,
			}},
		},
		{"Black-Rook-D1", square.D1, piece.New(piece.Black, piece.Rook),
			squareChanges{square.D1: pieceChanges{
				piece.New(piece.White, piece.Queen): false,
				piece.New(piece.Black, piece.Rook):  true,
			}},
		},

		// Knights.
		{"White-Knight-B3", square.B3, piece.New(piece.White, piece.Knight),
			squareChanges{square.B3: pieceChanges{piece.New(piece.White, piece.Knight): true}},
		},
		{"Black-Knight-B6", square.B6, piece.New(piece.Black, piece.Knight),
			squareChanges{square.B6: pieceChanges{piece.New(piece.Black, piece.Knight): true}},
		},
		{"White-Knight-B2", square.B2, piece.New(piece.White, piece.Knight),
			squareChanges{square.B2: pieceChanges{
				piece.New(piece.White, piece.Pawn):   false,
				piece.New(piece.White, piece.Knight): true,
			}},
		},
		{"Black-Knight-B7", square.B7, piece.New(piece.Black, piece.Knight),
			squareChanges{square.B7: pieceChanges{
				piece.New(piece.Black, piece.Pawn):   false,
				piece.New(piece.Black, piece.Knight): true,
			}},
		},
		{"White-Knight-A1", square.A1, piece.New(piece.White, piece.Knight),
			squareChanges{square.A1: pieceChanges{
				piece.New(piece.White, piece.Rook):   false,
				piece.New(piece.White, piece.Knight): true,
			}},
		},
		{"Black-Knight-A8", square.A8, piece.New(piece.Black, piece.Knight),
			squareChanges{square.A8: pieceChanges{
				piece.New(piece.Black, piece.Rook):   false,
				piece.New(piece.Black, piece.Knight): true,
			}},
		},
		{"White-Knight-B1", square.B1, piece.New(piece.White, piece.Knight),
			squareChanges{},
		},
		{"Black-Knight-B8", square.B8, piece.New(piece.Black, piece.Knight),
			squareChanges{},
		},
		{"White-Knight-C1", square.C1, piece.New(piece.White, piece.Knight),
			squareChanges{square.C1: pieceChanges{
				piece.New(piece.White, piece.Bishop): false,
				piece.New(piece.White, piece.Knight): true,
			}},
		},
		{"Black-Knight-C8", square.C8, piece.New(piece.Black, piece.Knight),
			squareChanges{square.C8: pieceChanges{
				piece.New(piece.Black, piece.Bishop): false,
				piece.New(piece.Black, piece.Knight): true,
			}},
		},
		{"White-Knight-D1", square.D1, piece.New(piece.White, piece.Knight),
			squareChanges{square.D1: pieceChanges{
				piece.New(piece.White, piece.Queen):  false,
				piece.New(piece.White, piece.Knight): true,
			}},
		},
		{"Black-Knight-D8", square.D8, piece.New(piece.Black, piece.Knight),
			squareChanges{square.D8: pieceChanges{
				piece.New(piece.Black, piece.Queen):  false,
				piece.New(piece.Black, piece.Knight): true,
			}},
		},
		{"White-Knight-E1", square.E1, piece.New(piece.White, piece.Knight),
			squareChanges{square.E1: pieceChanges{
				piece.New(piece.White, piece.King):   false,
				piece.New(piece.White, piece.Knight): true,
			}},
		},
		{"Black-Knight-E8", square.E8, piece.New(piece.Black, piece.Knight),
			squareChanges{square.E8: pieceChanges{
				piece.New(piece.Black, piece.King):   false,
				piece.New(piece.Black, piece.Knight): true,
			}},
		},
		{"White-Knight-G6", square.G6, piece.New(piece.White, piece.Knight),
			squareChanges{square.G6: pieceChanges{piece.New(piece.White, piece.Knight): true}},
		},
		{"Black-Knight-G3", square.G3, piece.New(piece.Black, piece.Knight),
			squareChanges{square.G3: pieceChanges{piece.New(piece.Black, piece.Knight): true}},
		},
		{"White-Knight-G7", square.G7, piece.New(piece.White, piece.Knight),
			squareChanges{square.G7: pieceChanges{
				piece.New(piece.Black, piece.Pawn):   false,
				piece.New(piece.White, piece.Knight): true,
			}},
		},
		{"Black-Knight-G2", square.G2, piece.New(piece.Black, piece.Knight),
			squareChanges{square.G2: pieceChanges{
				piece.New(piece.White, piece.Pawn):   false,
				piece.New(piece.Black, piece.Knight): true,
			}},
		},
		{"White-Knight-H8", square.H8, piece.New(piece.White, piece.Knight),
			squareChanges{square.H8: pieceChanges{
				piece.New(piece.Black, piece.Rook):   false,
				piece.New(piece.White, piece.Knight): true,
			}},
		},
		{"Black-Knight-H1", square.H1, piece.New(piece.Black, piece.Knight),
			squareChanges{square.H1: pieceChanges{
				piece.New(piece.White, piece.Rook):   false,
				piece.New(piece.Black, piece.Knight): true,
			}},
		},
		{"White-Knight-G8", square.G8, piece.New(piece.White, piece.Knight),
			squareChanges{square.G8: pieceChanges{
				piece.New(piece.Black, piece.Knight): false,
				piece.New(piece.White, piece.Knight): true,
			}},
		},
		{"Black-Knight-G1", square.G1, piece.New(piece.Black, piece.Knight),
			squareChanges{square.G1: pieceChanges{
				piece.New(piece.White, piece.Knight): false,
				piece.New(piece.Black, piece.Knight): true,
			}},
		},
		{"White-Knight-F8", square.F8, piece.New(piece.White, piece.Knight),
			squareChanges{square.F8: pieceChanges{
				piece.New(piece.Black, piece.Bishop): false,
				piece.New(piece.White, piece.Knight): true,
			}},
		},
		{"Black-Knight-F1", square.F1, piece.New(piece.Black, piece.Knight),
			squareChanges{square.F1: pieceChanges{
				piece.New(piece.White, piece.Bishop): false,
				piece.New(piece.Black, piece.Knight): true,
			}},
		},
		{"White-Knight-E8", square.E8, piece.New(piece.White, piece.Knight),
			squareChanges{square.E8: pieceChanges{
				piece.New(piece.Black, piece.King):   false,
				piece.New(piece.White, piece.Knight): true,
			}},
		},
		{"Black-Knight-E1", square.E1, piece.New(piece.Black, piece.Knight),
			squareChanges{square.E1: pieceChanges{
				piece.New(piece.White, piece.King):   false,
				piece.New(piece.Black, piece.Knight): true,
			}},
		},
		{"White-Knight-D8", square.D8, piece.New(piece.White, piece.Knight),
			squareChanges{square.D8: pieceChanges{
				piece.New(piece.Black, piece.Queen):  false,
				piece.New(piece.White, piece.Knight): true,
			}},
		},
		{"Black-Knight-D1", square.D1, piece.New(piece.Black, piece.Knight),
			squareChanges{square.D1: pieceChanges{
				piece.New(piece.White, piece.Queen):  false,
				piece.New(piece.Black, piece.Knight): true,
			}},
		},

		// Bishops.
		{"White-Bishop-C3", square.C3, piece.New(piece.White, piece.Bishop),
			squareChanges{square.C3: pieceChanges{piece.New(piece.White, piece.Bishop): true}},
		},
		{"Black-Bishop-C6", square.C6, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.C6: pieceChanges{piece.New(piece.Black, piece.Bishop): true}},
		},
		{"White-Bishop-C2", square.C2, piece.New(piece.White, piece.Bishop),
			squareChanges{square.C2: pieceChanges{
				piece.New(piece.White, piece.Pawn):   false,
				piece.New(piece.White, piece.Bishop): true,
			}},
		},
		{"Black-Bishop-C7", square.C7, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.C7: pieceChanges{
				piece.New(piece.Black, piece.Pawn):   false,
				piece.New(piece.Black, piece.Bishop): true,
			}},
		},
		{"White-Bishop-A1", square.A1, piece.New(piece.White, piece.Bishop),
			squareChanges{square.A1: pieceChanges{
				piece.New(piece.White, piece.Rook):   false,
				piece.New(piece.White, piece.Bishop): true,
			}},
		},
		{"Black-Bishop-A8", square.A8, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.A8: pieceChanges{
				piece.New(piece.Black, piece.Rook):   false,
				piece.New(piece.Black, piece.Bishop): true,
			}},
		},
		{"White-Bishop-B1", square.B1, piece.New(piece.White, piece.Bishop),
			squareChanges{square.B1: pieceChanges{
				piece.New(piece.White, piece.Knight): false,
				piece.New(piece.White, piece.Bishop): true,
			}},
		},
		{"Black-Bishop-B8", square.B8, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.B8: pieceChanges{
				piece.New(piece.Black, piece.Knight): false,
				piece.New(piece.Black, piece.Bishop): true,
			}},
		},
		{"White-Bishop-C1", square.C1, piece.New(piece.White, piece.Bishop),
			squareChanges{},
		},
		{"Black-Bishop-C8", square.C8, piece.New(piece.Black, piece.Bishop),
			squareChanges{},
		},
		{"White-Bishop-D1", square.D1, piece.New(piece.White, piece.Bishop),
			squareChanges{square.D1: pieceChanges{
				piece.New(piece.White, piece.Queen):  false,
				piece.New(piece.White, piece.Bishop): true,
			}},
		},
		{"Black-Bishop-D8", square.D8, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.D8: pieceChanges{
				piece.New(piece.Black, piece.Queen):  false,
				piece.New(piece.Black, piece.Bishop): true,
			}},
		},
		{"White-Bishop-E1", square.E1, piece.New(piece.White, piece.Bishop),
			squareChanges{square.E1: pieceChanges{
				piece.New(piece.White, piece.King):   false,
				piece.New(piece.White, piece.Bishop): true,
			}},
		},
		{"Black-Bishop-E8", square.E8, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.E8: pieceChanges{
				piece.New(piece.Black, piece.King):   false,
				piece.New(piece.Black, piece.Bishop): true,
			}},
		},
		{"White-Bishop-F6", square.F6, piece.New(piece.White, piece.Bishop),
			squareChanges{square.F6: pieceChanges{piece.New(piece.White, piece.Bishop): true}},
		},
		{"Black-Bishop-F3", square.F3, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.F3: pieceChanges{piece.New(piece.Black, piece.Bishop): true}},
		},
		{"White-Bishop-F7", square.F7, piece.New(piece.White, piece.Bishop),
			squareChanges{square.F7: pieceChanges{
				piece.New(piece.Black, piece.Pawn):   false,
				piece.New(piece.White, piece.Bishop): true,
			}},
		},
		{"Black-Bishop-F2", square.F2, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.F2: pieceChanges{
				piece.New(piece.White, piece.Pawn):   false,
				piece.New(piece.Black, piece.Bishop): true,
			}},
		},
		{"White-Bishop-H8", square.H8, piece.New(piece.White, piece.Bishop),
			squareChanges{square.H8: pieceChanges{
				piece.New(piece.Black, piece.Rook):   false,
				piece.New(piece.White, piece.Bishop): true,
			}},
		},
		{"Black-Bishop-H1", square.H1, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.H1: pieceChanges{
				piece.New(piece.White, piece.Rook):   false,
				piece.New(piece.Black, piece.Bishop): true,
			}},
		},
		{"White-Bishop-G8", square.G8, piece.New(piece.White, piece.Bishop),
			squareChanges{square.G8: pieceChanges{
				piece.New(piece.Black, piece.Knight): false,
				piece.New(piece.White, piece.Bishop): true,
			}},
		},
		{"Black-Bishop-G1", square.G1, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.G1: pieceChanges{
				piece.New(piece.White, piece.Knight): false,
				piece.New(piece.Black, piece.Bishop): true,
			}},
		},
		{"White-Bishop-F8", square.F8, piece.New(piece.White, piece.Bishop),
			squareChanges{square.F8: pieceChanges{
				piece.New(piece.Black, piece.Bishop): false,
				piece.New(piece.White, piece.Bishop): true,
			}},
		},
		{"Black-Bishop-F1", square.F1, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.F1: pieceChanges{
				piece.New(piece.White, piece.Bishop): false,
				piece.New(piece.Black, piece.Bishop): true,
			}},
		},
		{"White-Bishop-E8", square.E8, piece.New(piece.White, piece.Bishop),
			squareChanges{square.E8: pieceChanges{
				piece.New(piece.Black, piece.King):   false,
				piece.New(piece.White, piece.Bishop): true,
			}},
		},
		{"Black-Bishop-E1", square.E1, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.E1: pieceChanges{
				piece.New(piece.White, piece.King):   false,
				piece.New(piece.Black, piece.Bishop): true,
			}},
		},
		{"White-Bishop-D8", square.D8, piece.New(piece.White, piece.Bishop),
			squareChanges{square.D8: pieceChanges{
				piece.New(piece.Black, piece.Queen):  false,
				piece.New(piece.White, piece.Bishop): true,
			}},
		},
		{"Black-Bishop-D1", square.D1, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.D1: pieceChanges{
				piece.New(piece.White, piece.Queen):  false,
				piece.New(piece.Black, piece.Bishop): true,
			}},
		},

		// Queens.
		{"White-Queen-D3", square.D3, piece.New(piece.White, piece.Queen),
			squareChanges{square.D3: pieceChanges{piece.New(piece.White, piece.Queen): true}},
		},
		{"Black-Queen-D6", square.D6, piece.New(piece.Black, piece.Queen),
			squareChanges{square.D6: pieceChanges{piece.New(piece.Black, piece.Queen): true}},
		},
		{"White-Queen-D2", square.D2, piece.New(piece.White, piece.Queen),
			squareChanges{square.D2: pieceChanges{
				piece.New(piece.White, piece.Pawn):  false,
				piece.New(piece.White, piece.Queen): true,
			}},
		},
		{"Black-Queen-D7", square.D7, piece.New(piece.Black, piece.Queen),
			squareChanges{square.D7: pieceChanges{
				piece.New(piece.Black, piece.Pawn):  false,
				piece.New(piece.Black, piece.Queen): true,
			}},
		},
		{"White-Queen-A1", square.A1, piece.New(piece.White, piece.Queen),
			squareChanges{square.A1: pieceChanges{
				piece.New(piece.White, piece.Rook):  false,
				piece.New(piece.White, piece.Queen): true,
			}},
		},
		{"Black-Queen-A8", square.A8, piece.New(piece.Black, piece.Queen),
			squareChanges{square.A8: pieceChanges{
				piece.New(piece.Black, piece.Rook):  false,
				piece.New(piece.Black, piece.Queen): true,
			}},
		},
		{"White-Queen-B1", square.B1, piece.New(piece.White, piece.Queen),
			squareChanges{square.B1: pieceChanges{
				piece.New(piece.White, piece.Knight): false,
				piece.New(piece.White, piece.Queen):  true,
			}},
		},
		{"Black-Queen-B8", square.B8, piece.New(piece.Black, piece.Queen),
			squareChanges{square.B8: pieceChanges{
				piece.New(piece.Black, piece.Knight): false,
				piece.New(piece.Black, piece.Queen):  true,
			}},
		},
		{"White-Queen-C1", square.C1, piece.New(piece.White, piece.Queen),
			squareChanges{square.C1: pieceChanges{
				piece.New(piece.White, piece.Bishop): false,
				piece.New(piece.White, piece.Queen):  true,
			}},
		},
		{"Black-Queen-C8", square.C8, piece.New(piece.Black, piece.Queen),
			squareChanges{square.C8: pieceChanges{
				piece.New(piece.Black, piece.Bishop): false,
				piece.New(piece.Black, piece.Queen):  true,
			}},
		},
		{"White-Queen-D1", square.D1, piece.New(piece.White, piece.Queen),
			squareChanges{},
		},
		{"Black-Queen-D8", square.D8, piece.New(piece.Black, piece.Queen),
			squareChanges{},
		},
		{"White-Queen-E1", square.E1, piece.New(piece.White, piece.Queen),
			squareChanges{square.E1: pieceChanges{
				piece.New(piece.White, piece.King):  false,
				piece.New(piece.White, piece.Queen): true,
			}},
		},
		{"Black-Queen-E8", square.E8, piece.New(piece.Black, piece.Queen),
			squareChanges{square.E8: pieceChanges{
				piece.New(piece.Black, piece.King):  false,
				piece.New(piece.Black, piece.Queen): true,
			}},
		},
		{"White-Queen-D6", square.D6, piece.New(piece.White, piece.Queen),
			squareChanges{square.D6: pieceChanges{piece.New(piece.White, piece.Queen): true}},
		},
		{"Black-Queen-D3", square.D3, piece.New(piece.Black, piece.Queen),
			squareChanges{square.D3: pieceChanges{piece.New(piece.Black, piece.Queen): true}},
		},
		{"White-Queen-D7", square.D7, piece.New(piece.White, piece.Queen),
			squareChanges{square.D7: pieceChanges{
				piece.New(piece.Black, piece.Pawn):  false,
				piece.New(piece.White, piece.Queen): true,
			}},
		},
		{"Black-Queen-D2", square.D2, piece.New(piece.Black, piece.Queen),
			squareChanges{square.D2: pieceChanges{
				piece.New(piece.White, piece.Pawn):  false,
				piece.New(piece.Black, piece.Queen): true,
			}},
		},
		{"White-Queen-H8", square.H8, piece.New(piece.White, piece.Queen),
			squareChanges{square.H8: pieceChanges{
				piece.New(piece.Black, piece.Rook):  false,
				piece.New(piece.White, piece.Queen): true,
			}},
		},
		{"Black-Queen-H1", square.H1, piece.New(piece.Black, piece.Queen),
			squareChanges{square.H1: pieceChanges{
				piece.New(piece.White, piece.Rook):  false,
				piece.New(piece.Black, piece.Queen): true,
			}},
		},
		{"White-Queen-G8", square.G8, piece.New(piece.White, piece.Queen),
			squareChanges{square.G8: pieceChanges{
				piece.New(piece.Black, piece.Knight): false,
				piece.New(piece.White, piece.Queen):  true,
			}},
		},
		{"Black-Queen-G1", square.G1, piece.New(piece.Black, piece.Queen),
			squareChanges{square.G1: pieceChanges{
				piece.New(piece.White, piece.Knight): false,
				piece.New(piece.Black, piece.Queen):  true,
			}},
		},
		{"White-Queen-F8", square.F8, piece.New(piece.White, piece.Queen),
			squareChanges{square.F8: pieceChanges{
				piece.New(piece.Black, piece.Bishop): false,
				piece.New(piece.White, piece.Queen):  true,
			}},
		},
		{"Black-Queen-F1", square.F1, piece.New(piece.Black, piece.Queen),
			squareChanges{square.F1: pieceChanges{
				piece.New(piece.White, piece.Bishop): false,
				piece.New(piece.Black, piece.Queen):  true,
			}},
		},
		{"White-Queen-E8", square.E8, piece.New(piece.White, piece.Queen),
			squareChanges{square.E8: pieceChanges{
				piece.New(piece.Black, piece.King):  false,
				piece.New(piece.White, piece.Queen): true,
			}},
		},
		{"Black-Queen-E1", square.E1, piece.New(piece.Black, piece.Queen),
			squareChanges{square.E1: pieceChanges{
				piece.New(piece.White, piece.King):  false,
				piece.New(piece.Black, piece.Queen): true,
			}},
		},
		{"White-Queen-D8", square.D8, piece.New(piece.White, piece.Queen),
			squareChanges{square.D8: pieceChanges{
				piece.New(piece.Black, piece.Queen): false,
				piece.New(piece.White, piece.Queen): true,
			}},
		},
		{"Black-Queen-D1", square.D1, piece.New(piece.Black, piece.Queen),
			squareChanges{square.D1: pieceChanges{
				piece.New(piece.White, piece.Queen): false,
				piece.New(piece.Black, piece.Queen): true,
			}},
		},

		// Kings.
		{"White-King-E3", square.E3, piece.New(piece.White, piece.King),
			squareChanges{square.E3: pieceChanges{piece.New(piece.White, piece.King): true}},
		},
		{"Black-King-E6", square.E6, piece.New(piece.Black, piece.King),
			squareChanges{square.E6: pieceChanges{piece.New(piece.Black, piece.King): true}},
		},
		{"White-King-E2", square.E2, piece.New(piece.White, piece.King),
			squareChanges{square.E2: pieceChanges{
				piece.New(piece.White, piece.Pawn): false,
				piece.New(piece.White, piece.King): true,
			}},
		},
		{"Black-King-E7", square.E7, piece.New(piece.Black, piece.King),
			squareChanges{square.E7: pieceChanges{
				piece.New(piece.Black, piece.Pawn): false,
				piece.New(piece.Black, piece.King): true,
			}},
		},
		{"White-King-A1", square.A1, piece.New(piece.White, piece.King),
			squareChanges{square.A1: pieceChanges{
				piece.New(piece.White, piece.Rook): false,
				piece.New(piece.White, piece.King): true,
			}},
		},
		{"Black-King-A8", square.A8, piece.New(piece.Black, piece.King),
			squareChanges{square.A8: pieceChanges{
				piece.New(piece.Black, piece.Rook): false,
				piece.New(piece.Black, piece.King): true,
			}},
		},
		{"White-King-B1", square.B1, piece.New(piece.White, piece.King),
			squareChanges{square.B1: pieceChanges{
				piece.New(piece.White, piece.Knight): false,
				piece.New(piece.White, piece.King):   true,
			}},
		},
		{"Black-King-B8", square.B8, piece.New(piece.Black, piece.King),
			squareChanges{square.B8: pieceChanges{
				piece.New(piece.Black, piece.Knight): false,
				piece.New(piece.Black, piece.King):   true,
			}},
		},
		{"White-King-C1", square.C1, piece.New(piece.White, piece.King),
			squareChanges{square.C1: pieceChanges{
				piece.New(piece.White, piece.Bishop): false,
				piece.New(piece.White, piece.King):   true,
			}},
		},
		{"Black-King-C8", square.C8, piece.New(piece.Black, piece.King),
			squareChanges{square.C8: pieceChanges{
				piece.New(piece.Black, piece.Bishop): false,
				piece.New(piece.Black, piece.King):   true,
			}},
		},
		{"White-King-D1", square.D1, piece.New(piece.White, piece.King),
			squareChanges{square.D1: pieceChanges{
				piece.New(piece.White, piece.Queen): false,
				piece.New(piece.White, piece.King):  true,
			}},
		},
		{"Black-King-D8", square.D8, piece.New(piece.Black, piece.King),
			squareChanges{square.D8: pieceChanges{
				piece.New(piece.Black, piece.Queen): false,
				piece.New(piece.Black, piece.King):  true,
			}},
		},
		{"White-King-E1", square.E1, piece.New(piece.White, piece.King),
			squareChanges{},
		},
		{"Black-King-E8", square.E8, piece.New(piece.Black, piece.King),
			squareChanges{},
		},
		{"White-King-E6", square.E6, piece.New(piece.White, piece.King),
			squareChanges{square.E6: pieceChanges{piece.New(piece.White, piece.King): true}},
		},
		{"Black-King-E3", square.E3, piece.New(piece.Black, piece.King),
			squareChanges{square.E3: pieceChanges{piece.New(piece.Black, piece.King): true}},
		},
		{"White-King-E7", square.E7, piece.New(piece.White, piece.King),
			squareChanges{square.E7: pieceChanges{
				piece.New(piece.Black, piece.Pawn): false,
				piece.New(piece.White, piece.King): true,
			}},
		},
		{"Black-King-E2", square.E2, piece.New(piece.Black, piece.King),
			squareChanges{square.E2: pieceChanges{
				piece.New(piece.White, piece.Pawn): false,
				piece.New(piece.Black, piece.King): true,
			}},
		},
		{"White-King-H8", square.H8, piece.New(piece.White, piece.King),
			squareChanges{square.H8: pieceChanges{
				piece.New(piece.Black, piece.Rook): false,
				piece.New(piece.White, piece.King): true,
			}},
		},
		{"Black-King-H1", square.H1, piece.New(piece.Black, piece.King),
			squareChanges{square.H1: pieceChanges{
				piece.New(piece.White, piece.Rook): false,
				piece.New(piece.Black, piece.King): true,
			}},
		},
		{"White-King-G8", square.G8, piece.New(piece.White, piece.King),
			squareChanges{square.G8: pieceChanges{
				piece.New(piece.Black, piece.Knight): false,
				piece.New(piece.White, piece.King):   true,
			}},
		},
		{"Black-King-G1", square.G1, piece.New(piece.Black, piece.King),
			squareChanges{square.G1: pieceChanges{
				piece.New(piece.White, piece.Knight): false,
				piece.New(piece.Black, piece.King):   true,
			}},
		},
		{"White-King-F8", square.F8, piece.New(piece.White, piece.King),
			squareChanges{square.F8: pieceChanges{
				piece.New(piece.Black, piece.Bishop): false,
				piece.New(piece.White, piece.King):   true,
			}},
		},
		{"Black-King-F1", square.F1, piece.New(piece.Black, piece.King),
			squareChanges{square.F1: pieceChanges{
				piece.New(piece.White, piece.Bishop): false,
				piece.New(piece.Black, piece.King):   true,
			}},
		},
		{"White-King-E8", square.E8, piece.New(piece.White, piece.King),
			squareChanges{square.E8: pieceChanges{
				piece.New(piece.Black, piece.King): false,
				piece.New(piece.White, piece.King): true,
			}},
		},
		{"Black-King-E1", square.E1, piece.New(piece.Black, piece.King),
			squareChanges{square.E1: pieceChanges{
				piece.New(piece.White, piece.King): false,
				piece.New(piece.Black, piece.King): true,
			}},
		},
		{"White-King-D8", square.D8, piece.New(piece.White, piece.King),
			squareChanges{square.D8: pieceChanges{
				piece.New(piece.Black, piece.Queen): false,
				piece.New(piece.White, piece.King):  true,
			}},
		},
		{"Black-King-D1", square.D1, piece.New(piece.Black, piece.King),
			squareChanges{square.D1: pieceChanges{
				piece.New(piece.White, piece.Queen): false,
				piece.New(piece.Black, piece.King):  true,
			}},
		},

		// Non-standard squares and pieces.
		{"White-Pawn-NoSquare", square.NoSquare, piece.New(piece.White, piece.Pawn),
			squareChanges{},
		},
		{"Black-Pawn-NoSquare", square.NoSquare, piece.New(piece.Black, piece.Pawn),
			squareChanges{},
		},
		{"White-Pawn-Square(100)", square.Square(100), piece.New(piece.White, piece.Pawn),
			squareChanges{},
		},
		{"Black-Pawn-Square(100)", square.Square(100), piece.New(piece.Black, piece.Pawn),
			squareChanges{},
		},
		{"White-None-A1", square.A1, piece.New(piece.White, piece.None),
			squareChanges{square.A1: pieceChanges{piece.New(piece.White, piece.Rook): false}},
		},
		{"Black-None-A8", square.A8, piece.New(piece.Black, piece.None),
			squareChanges{square.A8: pieceChanges{piece.New(piece.Black, piece.Rook): false}},
		},
		{"White-Type(10)-A1", square.A1, piece.New(piece.White, piece.Type(10)),
			squareChanges{square.A1: pieceChanges{piece.New(piece.White, piece.Rook): false}},
		},
		{"Black-Type(10)-A8", square.A8, piece.New(piece.Black, piece.Type(10)),
			squareChanges{square.A8: pieceChanges{piece.New(piece.Black, piece.Rook): false}},
		},
		{"NoColor-Pawn-A1", square.A1, piece.New(piece.NoColor, piece.Pawn),
			squareChanges{square.A1: pieceChanges{piece.New(piece.White, piece.Rook): false}},
		},
		{"NoColor-Pawn-A8", square.A8, piece.New(piece.NoColor, piece.Pawn),
			squareChanges{square.A8: pieceChanges{piece.New(piece.Black, piece.Rook): false}},
		},
		{"Color(5)-Pawn-A1", square.A1, piece.New(piece.Color(5), piece.Pawn),
			squareChanges{square.A1: pieceChanges{piece.New(piece.White, piece.Rook): false}},
		},
		{"Color(5)-Pawn-A8", square.A8, piece.New(piece.Color(5), piece.Pawn),
			squareChanges{square.A8: pieceChanges{piece.New(piece.Black, piece.Rook): false}},
		},
		{"NoColor-None-A1", square.A1, piece.New(piece.NoColor, piece.None),
			squareChanges{square.A1: pieceChanges{piece.New(piece.White, piece.Rook): false}},
		},
		{"NoColor-None-A8", square.A8, piece.New(piece.NoColor, piece.None),
			squareChanges{square.A8: pieceChanges{piece.New(piece.Black, piece.Rook): false}},
		},
		{"NoColor-Type(10)-A1", square.A1, piece.New(piece.NoColor, piece.Type(10)),
			squareChanges{square.A1: pieceChanges{piece.New(piece.White, piece.Rook): false}},
		},
		{"NoColor-Type(10)-A8", square.A8, piece.New(piece.NoColor, piece.Type(10)),
			squareChanges{square.A8: pieceChanges{piece.New(piece.Black, piece.Rook): false}},
		},
		{"Color(5)-None-A1", square.A1, piece.New(piece.Color(5), piece.None),
			squareChanges{square.A1: pieceChanges{piece.New(piece.White, piece.Rook): false}},
		},
		{"Color(5)-None-A8", square.A8, piece.New(piece.Color(5), piece.None),
			squareChanges{square.A8: pieceChanges{piece.New(piece.Black, piece.Rook): false}},
		},
		{"Color(5)-Type(10)-A1", square.A1, piece.New(piece.Color(5), piece.Type(10)),
			squareChanges{square.A1: pieceChanges{piece.New(piece.White, piece.Rook): false}},
		},
		{"Color(5)-Type(10)-A8", square.A8, piece.New(piece.Color(5), piece.Type(10)),
			squareChanges{square.A8: pieceChanges{piece.New(piece.Black, piece.Rook): false}},
		},
		{"NoColor-None-NoSquare", square.NoSquare, piece.New(piece.NoColor, piece.None),
			squareChanges{},
		},
		{"NoColor-None-NoSquare", square.NoSquare, piece.New(piece.NoColor, piece.None),
			squareChanges{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := New()
			defer func() {
				if err := recover(); err != nil {
					t.Errorf("For initial board, *Position.Put(%v, %v) should not panic, but panicked with: %v", tc.pc, tc.sq, err)
				}
			}()
			p.Put(tc.pc, tc.sq)
			ch := changedSquares(New(), p)
			if !reflect.DeepEqual(ch, tc.wantChanged) {
				t.Errorf("For initial board and *Position.Put(%v, %v), board changes are %v, want %v", tc.pc, tc.sq, ch, tc.wantChanged)
			}
		})
	}
}

func TestPositionQuickPut(t *testing.T) {
	testCases := []struct {
		name        string
		sq          square.Square
		pc          piece.Piece
		wantChanged squareChanges
	}{
		// Pawns.
		{"White-Pawn-A3", square.A3, piece.New(piece.White, piece.Pawn),
			squareChanges{square.A3: pieceChanges{piece.New(piece.White, piece.Pawn): true}},
		},
		{"Black-Pawn-A6", square.A6, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.A6: pieceChanges{piece.New(piece.Black, piece.Pawn): true}},
		},
		{"White-Pawn-A2", square.A2, piece.New(piece.White, piece.Pawn),
			squareChanges{},
		},
		{"Black-Pawn-A7", square.A7, piece.New(piece.Black, piece.Pawn),
			squareChanges{},
		},
		{"White-Pawn-A1", square.A1, piece.New(piece.White, piece.Pawn),
			squareChanges{square.A1: pieceChanges{piece.New(piece.White, piece.Pawn): true}},
		},
		{"Black-Pawn-A8", square.A8, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.A8: pieceChanges{piece.New(piece.Black, piece.Pawn): true}},
		},
		{"White-Pawn-B1", square.B1, piece.New(piece.White, piece.Pawn),
			squareChanges{square.B1: pieceChanges{piece.New(piece.White, piece.Pawn): true}},
		},
		{"Black-Pawn-B8", square.B8, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.B8: pieceChanges{piece.New(piece.Black, piece.Pawn): true}},
		},
		{"White-Pawn-C1", square.C1, piece.New(piece.White, piece.Pawn),
			squareChanges{square.C1: pieceChanges{piece.New(piece.White, piece.Pawn): true}},
		},
		{"Black-Pawn-C8", square.C8, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.C8: pieceChanges{piece.New(piece.Black, piece.Pawn): true}},
		},
		{"White-Pawn-D1", square.D1, piece.New(piece.White, piece.Pawn),
			squareChanges{square.D1: pieceChanges{piece.New(piece.White, piece.Pawn): true}},
		},
		{"Black-Pawn-D8", square.D8, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.D8: pieceChanges{piece.New(piece.Black, piece.Pawn): true}},
		},
		{"White-Pawn-E1", square.E1, piece.New(piece.White, piece.Pawn),
			squareChanges{square.E1: pieceChanges{piece.New(piece.White, piece.Pawn): true}},
		},
		{"Black-Pawn-E8", square.E8, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.E8: pieceChanges{piece.New(piece.Black, piece.Pawn): true}},
		},
		{"White-Pawn-H6", square.H6, piece.New(piece.White, piece.Pawn),
			squareChanges{square.H6: pieceChanges{piece.New(piece.White, piece.Pawn): true}},
		},
		{"Black-Pawn-H3", square.H3, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.H3: pieceChanges{piece.New(piece.Black, piece.Pawn): true}},
		},
		{"White-Pawn-H7", square.H7, piece.New(piece.White, piece.Pawn),
			squareChanges{square.H7: pieceChanges{piece.New(piece.White, piece.Pawn): true}},
		},
		{"Black-Pawn-H2", square.H2, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.H2: pieceChanges{piece.New(piece.Black, piece.Pawn): true}},
		},
		{"White-Pawn-H8", square.H8, piece.New(piece.White, piece.Pawn),
			squareChanges{square.H8: pieceChanges{piece.New(piece.White, piece.Pawn): true}},
		},
		{"Black-Pawn-H1", square.H1, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.H1: pieceChanges{piece.New(piece.Black, piece.Pawn): true}},
		},
		{"White-Pawn-G8", square.G8, piece.New(piece.White, piece.Pawn),
			squareChanges{square.G8: pieceChanges{piece.New(piece.White, piece.Pawn): true}},
		},
		{"Black-Pawn-G1", square.G1, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.G1: pieceChanges{piece.New(piece.Black, piece.Pawn): true}},
		},
		{"White-Pawn-F8", square.F8, piece.New(piece.White, piece.Pawn),
			squareChanges{square.F8: pieceChanges{piece.New(piece.White, piece.Pawn): true}},
		},
		{"Black-Pawn-F1", square.F1, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.F1: pieceChanges{piece.New(piece.Black, piece.Pawn): true}},
		},
		{"White-Pawn-E8", square.E8, piece.New(piece.White, piece.Pawn),
			squareChanges{square.E8: pieceChanges{piece.New(piece.White, piece.Pawn): true}},
		},
		{"Black-Pawn-E1", square.E1, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.E1: pieceChanges{piece.New(piece.Black, piece.Pawn): true}},
		},
		{"White-Pawn-D8", square.D8, piece.New(piece.White, piece.Pawn),
			squareChanges{square.D8: pieceChanges{piece.New(piece.White, piece.Pawn): true}},
		},
		{"Black-Pawn-D1", square.D1, piece.New(piece.Black, piece.Pawn),
			squareChanges{square.D1: pieceChanges{piece.New(piece.Black, piece.Pawn): true}},
		},

		// Rooks.
		{"White-Rook-A3", square.A3, piece.New(piece.White, piece.Rook),
			squareChanges{square.A3: pieceChanges{piece.New(piece.White, piece.Rook): true}},
		},
		{"Black-Rook-A6", square.A6, piece.New(piece.Black, piece.Rook),
			squareChanges{square.A6: pieceChanges{piece.New(piece.Black, piece.Rook): true}},
		},
		{"White-Rook-A2", square.A2, piece.New(piece.White, piece.Rook),
			squareChanges{square.A2: pieceChanges{piece.New(piece.White, piece.Rook): true}},
		},
		{"Black-Rook-A7", square.A7, piece.New(piece.Black, piece.Rook),
			squareChanges{square.A7: pieceChanges{piece.New(piece.Black, piece.Rook): true}},
		},
		{"White-Rook-A1", square.A1, piece.New(piece.White, piece.Rook),
			squareChanges{},
		},
		{"Black-Rook-A8", square.A8, piece.New(piece.Black, piece.Rook),
			squareChanges{},
		},
		{"White-Rook-B1", square.B1, piece.New(piece.White, piece.Rook),
			squareChanges{square.B1: pieceChanges{piece.New(piece.White, piece.Rook): true}},
		},
		{"Black-Rook-B8", square.B8, piece.New(piece.Black, piece.Rook),
			squareChanges{square.B8: pieceChanges{piece.New(piece.Black, piece.Rook): true}},
		},
		{"White-Rook-C1", square.C1, piece.New(piece.White, piece.Rook),
			squareChanges{square.C1: pieceChanges{piece.New(piece.White, piece.Rook): true}},
		},
		{"Black-Rook-C8", square.C8, piece.New(piece.Black, piece.Rook),
			squareChanges{square.C8: pieceChanges{piece.New(piece.Black, piece.Rook): true}},
		},
		{"White-Rook-D1", square.D1, piece.New(piece.White, piece.Rook),
			squareChanges{square.D1: pieceChanges{piece.New(piece.White, piece.Rook): true}},
		},
		{"Black-Rook-D8", square.D8, piece.New(piece.Black, piece.Rook),
			squareChanges{square.D8: pieceChanges{piece.New(piece.Black, piece.Rook): true}},
		},
		{"White-Rook-E1", square.E1, piece.New(piece.White, piece.Rook),
			squareChanges{square.E1: pieceChanges{piece.New(piece.White, piece.Rook): true}},
		},
		{"Black-Rook-E8", square.E8, piece.New(piece.Black, piece.Rook),
			squareChanges{square.E8: pieceChanges{piece.New(piece.Black, piece.Rook): true}},
		},
		{"White-Rook-H6", square.H6, piece.New(piece.White, piece.Rook),
			squareChanges{square.H6: pieceChanges{piece.New(piece.White, piece.Rook): true}},
		},
		{"Black-Rook-H3", square.H3, piece.New(piece.Black, piece.Rook),
			squareChanges{square.H3: pieceChanges{piece.New(piece.Black, piece.Rook): true}},
		},
		{"White-Rook-H7", square.H7, piece.New(piece.White, piece.Rook),
			squareChanges{square.H7: pieceChanges{piece.New(piece.White, piece.Rook): true}},
		},
		{"Black-Rook-H2", square.H2, piece.New(piece.Black, piece.Rook),
			squareChanges{square.H2: pieceChanges{piece.New(piece.Black, piece.Rook): true}},
		},
		{"White-Rook-H8", square.H8, piece.New(piece.White, piece.Rook),
			squareChanges{square.H8: pieceChanges{piece.New(piece.White, piece.Rook): true}},
		},
		{"Black-Rook-H1", square.H1, piece.New(piece.Black, piece.Rook),
			squareChanges{square.H1: pieceChanges{piece.New(piece.Black, piece.Rook): true}},
		},
		{"White-Rook-G8", square.G8, piece.New(piece.White, piece.Rook),
			squareChanges{square.G8: pieceChanges{piece.New(piece.White, piece.Rook): true}},
		},
		{"Black-Rook-G1", square.G1, piece.New(piece.Black, piece.Rook),
			squareChanges{square.G1: pieceChanges{piece.New(piece.Black, piece.Rook): true}},
		},
		{"White-Rook-F8", square.F8, piece.New(piece.White, piece.Rook),
			squareChanges{square.F8: pieceChanges{piece.New(piece.White, piece.Rook): true}},
		},
		{"Black-Rook-F1", square.F1, piece.New(piece.Black, piece.Rook),
			squareChanges{square.F1: pieceChanges{piece.New(piece.Black, piece.Rook): true}},
		},
		{"White-Rook-E8", square.E8, piece.New(piece.White, piece.Rook),
			squareChanges{square.E8: pieceChanges{piece.New(piece.White, piece.Rook): true}},
		},
		{"Black-Rook-E1", square.E1, piece.New(piece.Black, piece.Rook),
			squareChanges{square.E1: pieceChanges{piece.New(piece.Black, piece.Rook): true}},
		},
		{"White-Rook-D8", square.D8, piece.New(piece.White, piece.Rook),
			squareChanges{square.D8: pieceChanges{piece.New(piece.White, piece.Rook): true}},
		},
		{"Black-Rook-D1", square.D1, piece.New(piece.Black, piece.Rook),
			squareChanges{square.D1: pieceChanges{piece.New(piece.Black, piece.Rook): true}},
		},

		// Knights.
		{"White-Knight-B3", square.B3, piece.New(piece.White, piece.Knight),
			squareChanges{square.B3: pieceChanges{piece.New(piece.White, piece.Knight): true}},
		},
		{"Black-Knight-B6", square.B6, piece.New(piece.Black, piece.Knight),
			squareChanges{square.B6: pieceChanges{piece.New(piece.Black, piece.Knight): true}},
		},
		{"White-Knight-B2", square.B2, piece.New(piece.White, piece.Knight),
			squareChanges{square.B2: pieceChanges{piece.New(piece.White, piece.Knight): true}},
		},
		{"Black-Knight-B7", square.B7, piece.New(piece.Black, piece.Knight),
			squareChanges{square.B7: pieceChanges{piece.New(piece.Black, piece.Knight): true}},
		},
		{"White-Knight-A1", square.A1, piece.New(piece.White, piece.Knight),
			squareChanges{square.A1: pieceChanges{piece.New(piece.White, piece.Knight): true}},
		},
		{"Black-Knight-A8", square.A8, piece.New(piece.Black, piece.Knight),
			squareChanges{square.A8: pieceChanges{piece.New(piece.Black, piece.Knight): true}},
		},
		{"White-Knight-B1", square.B1, piece.New(piece.White, piece.Knight),
			squareChanges{},
		},
		{"Black-Knight-B8", square.B8, piece.New(piece.Black, piece.Knight),
			squareChanges{},
		},
		{"White-Knight-C1", square.C1, piece.New(piece.White, piece.Knight),
			squareChanges{square.C1: pieceChanges{piece.New(piece.White, piece.Knight): true}},
		},
		{"Black-Knight-C8", square.C8, piece.New(piece.Black, piece.Knight),
			squareChanges{square.C8: pieceChanges{piece.New(piece.Black, piece.Knight): true}},
		},
		{"White-Knight-D1", square.D1, piece.New(piece.White, piece.Knight),
			squareChanges{square.D1: pieceChanges{piece.New(piece.White, piece.Knight): true}},
		},
		{"Black-Knight-D8", square.D8, piece.New(piece.Black, piece.Knight),
			squareChanges{square.D8: pieceChanges{piece.New(piece.Black, piece.Knight): true}},
		},
		{"White-Knight-E1", square.E1, piece.New(piece.White, piece.Knight),
			squareChanges{square.E1: pieceChanges{piece.New(piece.White, piece.Knight): true}},
		},
		{"Black-Knight-E8", square.E8, piece.New(piece.Black, piece.Knight),
			squareChanges{square.E8: pieceChanges{piece.New(piece.Black, piece.Knight): true}},
		},
		{"White-Knight-G6", square.G6, piece.New(piece.White, piece.Knight),
			squareChanges{square.G6: pieceChanges{piece.New(piece.White, piece.Knight): true}},
		},
		{"Black-Knight-G3", square.G3, piece.New(piece.Black, piece.Knight),
			squareChanges{square.G3: pieceChanges{piece.New(piece.Black, piece.Knight): true}},
		},
		{"White-Knight-G7", square.G7, piece.New(piece.White, piece.Knight),
			squareChanges{square.G7: pieceChanges{piece.New(piece.White, piece.Knight): true}},
		},
		{"Black-Knight-G2", square.G2, piece.New(piece.Black, piece.Knight),
			squareChanges{square.G2: pieceChanges{piece.New(piece.Black, piece.Knight): true}},
		},
		{"White-Knight-H8", square.H8, piece.New(piece.White, piece.Knight),
			squareChanges{square.H8: pieceChanges{piece.New(piece.White, piece.Knight): true}},
		},
		{"Black-Knight-H1", square.H1, piece.New(piece.Black, piece.Knight),
			squareChanges{square.H1: pieceChanges{piece.New(piece.Black, piece.Knight): true}},
		},
		{"White-Knight-G8", square.G8, piece.New(piece.White, piece.Knight),
			squareChanges{square.G8: pieceChanges{piece.New(piece.White, piece.Knight): true}},
		},
		{"Black-Knight-G1", square.G1, piece.New(piece.Black, piece.Knight),
			squareChanges{square.G1: pieceChanges{piece.New(piece.Black, piece.Knight): true}},
		},
		{"White-Knight-F8", square.F8, piece.New(piece.White, piece.Knight),
			squareChanges{square.F8: pieceChanges{piece.New(piece.White, piece.Knight): true}},
		},
		{"Black-Knight-F1", square.F1, piece.New(piece.Black, piece.Knight),
			squareChanges{square.F1: pieceChanges{piece.New(piece.Black, piece.Knight): true}},
		},
		{"White-Knight-E8", square.E8, piece.New(piece.White, piece.Knight),
			squareChanges{square.E8: pieceChanges{piece.New(piece.White, piece.Knight): true}},
		},
		{"Black-Knight-E1", square.E1, piece.New(piece.Black, piece.Knight),
			squareChanges{square.E1: pieceChanges{piece.New(piece.Black, piece.Knight): true}},
		},
		{"White-Knight-D8", square.D8, piece.New(piece.White, piece.Knight),
			squareChanges{square.D8: pieceChanges{piece.New(piece.White, piece.Knight): true}},
		},
		{"Black-Knight-D1", square.D1, piece.New(piece.Black, piece.Knight),
			squareChanges{square.D1: pieceChanges{piece.New(piece.Black, piece.Knight): true}},
		},

		// Bishops.
		{"White-Bishop-C3", square.C3, piece.New(piece.White, piece.Bishop),
			squareChanges{square.C3: pieceChanges{piece.New(piece.White, piece.Bishop): true}},
		},
		{"Black-Bishop-C6", square.C6, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.C6: pieceChanges{piece.New(piece.Black, piece.Bishop): true}},
		},
		{"White-Bishop-C2", square.C2, piece.New(piece.White, piece.Bishop),
			squareChanges{square.C2: pieceChanges{piece.New(piece.White, piece.Bishop): true}},
		},
		{"Black-Bishop-C7", square.C7, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.C7: pieceChanges{piece.New(piece.Black, piece.Bishop): true}},
		},
		{"White-Bishop-A1", square.A1, piece.New(piece.White, piece.Bishop),
			squareChanges{square.A1: pieceChanges{piece.New(piece.White, piece.Bishop): true}},
		},
		{"Black-Bishop-A8", square.A8, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.A8: pieceChanges{piece.New(piece.Black, piece.Bishop): true}},
		},
		{"White-Bishop-B1", square.B1, piece.New(piece.White, piece.Bishop),
			squareChanges{square.B1: pieceChanges{piece.New(piece.White, piece.Bishop): true}},
		},
		{"Black-Bishop-B8", square.B8, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.B8: pieceChanges{piece.New(piece.Black, piece.Bishop): true}},
		},
		{"White-Bishop-C1", square.C1, piece.New(piece.White, piece.Bishop),
			squareChanges{},
		},
		{"Black-Bishop-C8", square.C8, piece.New(piece.Black, piece.Bishop),
			squareChanges{},
		},
		{"White-Bishop-D1", square.D1, piece.New(piece.White, piece.Bishop),
			squareChanges{square.D1: pieceChanges{piece.New(piece.White, piece.Bishop): true}},
		},
		{"Black-Bishop-D8", square.D8, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.D8: pieceChanges{piece.New(piece.Black, piece.Bishop): true}},
		},
		{"White-Bishop-E1", square.E1, piece.New(piece.White, piece.Bishop),
			squareChanges{square.E1: pieceChanges{piece.New(piece.White, piece.Bishop): true}},
		},
		{"Black-Bishop-E8", square.E8, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.E8: pieceChanges{piece.New(piece.Black, piece.Bishop): true}},
		},
		{"White-Bishop-F6", square.F6, piece.New(piece.White, piece.Bishop),
			squareChanges{square.F6: pieceChanges{piece.New(piece.White, piece.Bishop): true}},
		},
		{"Black-Bishop-F3", square.F3, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.F3: pieceChanges{piece.New(piece.Black, piece.Bishop): true}},
		},
		{"White-Bishop-F7", square.F7, piece.New(piece.White, piece.Bishop),
			squareChanges{square.F7: pieceChanges{piece.New(piece.White, piece.Bishop): true}},
		},
		{"Black-Bishop-F2", square.F2, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.F2: pieceChanges{piece.New(piece.Black, piece.Bishop): true}},
		},
		{"White-Bishop-H8", square.H8, piece.New(piece.White, piece.Bishop),
			squareChanges{square.H8: pieceChanges{piece.New(piece.White, piece.Bishop): true}},
		},
		{"Black-Bishop-H1", square.H1, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.H1: pieceChanges{piece.New(piece.Black, piece.Bishop): true}},
		},
		{"White-Bishop-G8", square.G8, piece.New(piece.White, piece.Bishop),
			squareChanges{square.G8: pieceChanges{piece.New(piece.White, piece.Bishop): true}},
		},
		{"Black-Bishop-G1", square.G1, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.G1: pieceChanges{piece.New(piece.Black, piece.Bishop): true}},
		},
		{"White-Bishop-F8", square.F8, piece.New(piece.White, piece.Bishop),
			squareChanges{square.F8: pieceChanges{piece.New(piece.White, piece.Bishop): true}},
		},
		{"Black-Bishop-F1", square.F1, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.F1: pieceChanges{piece.New(piece.Black, piece.Bishop): true}},
		},
		{"White-Bishop-E8", square.E8, piece.New(piece.White, piece.Bishop),
			squareChanges{square.E8: pieceChanges{piece.New(piece.White, piece.Bishop): true}},
		},
		{"Black-Bishop-E1", square.E1, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.E1: pieceChanges{piece.New(piece.Black, piece.Bishop): true}},
		},
		{"White-Bishop-D8", square.D8, piece.New(piece.White, piece.Bishop),
			squareChanges{square.D8: pieceChanges{piece.New(piece.White, piece.Bishop): true}},
		},
		{"Black-Bishop-D1", square.D1, piece.New(piece.Black, piece.Bishop),
			squareChanges{square.D1: pieceChanges{piece.New(piece.Black, piece.Bishop): true}},
		},

		// Queens.
		{"White-Queen-D3", square.D3, piece.New(piece.White, piece.Queen),
			squareChanges{square.D3: pieceChanges{piece.New(piece.White, piece.Queen): true}},
		},
		{"Black-Queen-D6", square.D6, piece.New(piece.Black, piece.Queen),
			squareChanges{square.D6: pieceChanges{piece.New(piece.Black, piece.Queen): true}},
		},
		{"White-Queen-D2", square.D2, piece.New(piece.White, piece.Queen),
			squareChanges{square.D2: pieceChanges{piece.New(piece.White, piece.Queen): true}},
		},
		{"Black-Queen-D7", square.D7, piece.New(piece.Black, piece.Queen),
			squareChanges{square.D7: pieceChanges{piece.New(piece.Black, piece.Queen): true}},
		},
		{"White-Queen-A1", square.A1, piece.New(piece.White, piece.Queen),
			squareChanges{square.A1: pieceChanges{piece.New(piece.White, piece.Queen): true}},
		},
		{"Black-Queen-A8", square.A8, piece.New(piece.Black, piece.Queen),
			squareChanges{square.A8: pieceChanges{piece.New(piece.Black, piece.Queen): true}},
		},
		{"White-Queen-B1", square.B1, piece.New(piece.White, piece.Queen),
			squareChanges{square.B1: pieceChanges{piece.New(piece.White, piece.Queen): true}},
		},
		{"Black-Queen-B8", square.B8, piece.New(piece.Black, piece.Queen),
			squareChanges{square.B8: pieceChanges{piece.New(piece.Black, piece.Queen): true}},
		},
		{"White-Queen-C1", square.C1, piece.New(piece.White, piece.Queen),
			squareChanges{square.C1: pieceChanges{piece.New(piece.White, piece.Queen): true}},
		},
		{"Black-Queen-C8", square.C8, piece.New(piece.Black, piece.Queen),
			squareChanges{square.C8: pieceChanges{piece.New(piece.Black, piece.Queen): true}},
		},
		{"White-Queen-D1", square.D1, piece.New(piece.White, piece.Queen),
			squareChanges{},
		},
		{"Black-Queen-D8", square.D8, piece.New(piece.Black, piece.Queen),
			squareChanges{},
		},
		{"White-Queen-E1", square.E1, piece.New(piece.White, piece.Queen),
			squareChanges{square.E1: pieceChanges{piece.New(piece.White, piece.Queen): true}},
		},
		{"Black-Queen-E8", square.E8, piece.New(piece.Black, piece.Queen),
			squareChanges{square.E8: pieceChanges{piece.New(piece.Black, piece.Queen): true}},
		},
		{"White-Queen-D6", square.D6, piece.New(piece.White, piece.Queen),
			squareChanges{square.D6: pieceChanges{piece.New(piece.White, piece.Queen): true}},
		},
		{"Black-Queen-D3", square.D3, piece.New(piece.Black, piece.Queen),
			squareChanges{square.D3: pieceChanges{piece.New(piece.Black, piece.Queen): true}},
		},
		{"White-Queen-D7", square.D7, piece.New(piece.White, piece.Queen),
			squareChanges{square.D7: pieceChanges{piece.New(piece.White, piece.Queen): true}},
		},
		{"Black-Queen-D2", square.D2, piece.New(piece.Black, piece.Queen),
			squareChanges{square.D2: pieceChanges{piece.New(piece.Black, piece.Queen): true}},
		},
		{"White-Queen-H8", square.H8, piece.New(piece.White, piece.Queen),
			squareChanges{square.H8: pieceChanges{piece.New(piece.White, piece.Queen): true}},
		},
		{"Black-Queen-H1", square.H1, piece.New(piece.Black, piece.Queen),
			squareChanges{square.H1: pieceChanges{piece.New(piece.Black, piece.Queen): true}},
		},
		{"White-Queen-G8", square.G8, piece.New(piece.White, piece.Queen),
			squareChanges{square.G8: pieceChanges{piece.New(piece.White, piece.Queen): true}},
		},
		{"Black-Queen-G1", square.G1, piece.New(piece.Black, piece.Queen),
			squareChanges{square.G1: pieceChanges{piece.New(piece.Black, piece.Queen): true}},
		},
		{"White-Queen-F8", square.F8, piece.New(piece.White, piece.Queen),
			squareChanges{square.F8: pieceChanges{piece.New(piece.White, piece.Queen): true}},
		},
		{"Black-Queen-F1", square.F1, piece.New(piece.Black, piece.Queen),
			squareChanges{square.F1: pieceChanges{piece.New(piece.Black, piece.Queen): true}},
		},
		{"White-Queen-E8", square.E8, piece.New(piece.White, piece.Queen),
			squareChanges{square.E8: pieceChanges{piece.New(piece.White, piece.Queen): true}},
		},
		{"Black-Queen-E1", square.E1, piece.New(piece.Black, piece.Queen),
			squareChanges{square.E1: pieceChanges{piece.New(piece.Black, piece.Queen): true}},
		},
		{"White-Queen-D8", square.D8, piece.New(piece.White, piece.Queen),
			squareChanges{square.D8: pieceChanges{piece.New(piece.White, piece.Queen): true}},
		},
		{"Black-Queen-D1", square.D1, piece.New(piece.Black, piece.Queen),
			squareChanges{square.D1: pieceChanges{piece.New(piece.Black, piece.Queen): true}},
		},

		// Kings.
		{"White-King-E3", square.E3, piece.New(piece.White, piece.King),
			squareChanges{square.E3: pieceChanges{piece.New(piece.White, piece.King): true}},
		},
		{"Black-King-E6", square.E6, piece.New(piece.Black, piece.King),
			squareChanges{square.E6: pieceChanges{piece.New(piece.Black, piece.King): true}},
		},
		{"White-King-E2", square.E2, piece.New(piece.White, piece.King),
			squareChanges{square.E2: pieceChanges{piece.New(piece.White, piece.King): true}},
		},
		{"Black-King-E7", square.E7, piece.New(piece.Black, piece.King),
			squareChanges{square.E7: pieceChanges{piece.New(piece.Black, piece.King): true}},
		},
		{"White-King-A1", square.A1, piece.New(piece.White, piece.King),
			squareChanges{square.A1: pieceChanges{piece.New(piece.White, piece.King): true}},
		},
		{"Black-King-A8", square.A8, piece.New(piece.Black, piece.King),
			squareChanges{square.A8: pieceChanges{piece.New(piece.Black, piece.King): true}},
		},
		{"White-King-B1", square.B1, piece.New(piece.White, piece.King),
			squareChanges{square.B1: pieceChanges{piece.New(piece.White, piece.King): true}},
		},
		{"Black-King-B8", square.B8, piece.New(piece.Black, piece.King),
			squareChanges{square.B8: pieceChanges{piece.New(piece.Black, piece.King): true}},
		},
		{"White-King-C1", square.C1, piece.New(piece.White, piece.King),
			squareChanges{square.C1: pieceChanges{piece.New(piece.White, piece.King): true}},
		},
		{"Black-King-C8", square.C8, piece.New(piece.Black, piece.King),
			squareChanges{square.C8: pieceChanges{piece.New(piece.Black, piece.King): true}},
		},
		{"White-King-D1", square.D1, piece.New(piece.White, piece.King),
			squareChanges{square.D1: pieceChanges{piece.New(piece.White, piece.King): true}},
		},
		{"Black-King-D8", square.D8, piece.New(piece.Black, piece.King),
			squareChanges{square.D8: pieceChanges{piece.New(piece.Black, piece.King): true}},
		},
		{"White-King-E1", square.E1, piece.New(piece.White, piece.King),
			squareChanges{},
		},
		{"Black-King-E8", square.E8, piece.New(piece.Black, piece.King),
			squareChanges{},
		},
		{"White-King-E6", square.E6, piece.New(piece.White, piece.King),
			squareChanges{square.E6: pieceChanges{piece.New(piece.White, piece.King): true}},
		},
		{"Black-King-E3", square.E3, piece.New(piece.Black, piece.King),
			squareChanges{square.E3: pieceChanges{piece.New(piece.Black, piece.King): true}},
		},
		{"White-King-E7", square.E7, piece.New(piece.White, piece.King),
			squareChanges{square.E7: pieceChanges{piece.New(piece.White, piece.King): true}},
		},
		{"Black-King-E2", square.E2, piece.New(piece.Black, piece.King),
			squareChanges{square.E2: pieceChanges{piece.New(piece.Black, piece.King): true}},
		},
		{"White-King-H8", square.H8, piece.New(piece.White, piece.King),
			squareChanges{square.H8: pieceChanges{piece.New(piece.White, piece.King): true}},
		},
		{"Black-King-H1", square.H1, piece.New(piece.Black, piece.King),
			squareChanges{square.H1: pieceChanges{piece.New(piece.Black, piece.King): true}},
		},
		{"White-King-G8", square.G8, piece.New(piece.White, piece.King),
			squareChanges{square.G8: pieceChanges{piece.New(piece.White, piece.King): true}},
		},
		{"Black-King-G1", square.G1, piece.New(piece.Black, piece.King),
			squareChanges{square.G1: pieceChanges{piece.New(piece.Black, piece.King): true}},
		},
		{"White-King-F8", square.F8, piece.New(piece.White, piece.King),
			squareChanges{square.F8: pieceChanges{piece.New(piece.White, piece.King): true}},
		},
		{"Black-King-F1", square.F1, piece.New(piece.Black, piece.King),
			squareChanges{square.F1: pieceChanges{piece.New(piece.Black, piece.King): true}},
		},
		{"White-King-E8", square.E8, piece.New(piece.White, piece.King),
			squareChanges{square.E8: pieceChanges{piece.New(piece.White, piece.King): true}},
		},
		{"Black-King-E1", square.E1, piece.New(piece.Black, piece.King),
			squareChanges{square.E1: pieceChanges{piece.New(piece.Black, piece.King): true}},
		},
		{"White-King-D8", square.D8, piece.New(piece.White, piece.King),
			squareChanges{square.D8: pieceChanges{piece.New(piece.White, piece.King): true}},
		},
		{"Black-King-D1", square.D1, piece.New(piece.Black, piece.King),
			squareChanges{square.D1: pieceChanges{piece.New(piece.Black, piece.King): true}},
		},

		// Non-standard squares and pieces.
		{"White-Pawn-NoSquare", square.NoSquare, piece.New(piece.White, piece.Pawn),
			squareChanges{},
		},
		{"Black-Pawn-NoSquare", square.NoSquare, piece.New(piece.Black, piece.Pawn),
			squareChanges{},
		},
		{"White-Pawn-Square(100)", square.Square(100), piece.New(piece.White, piece.Pawn),
			squareChanges{},
		},
		{"Black-Pawn-Square(100)", square.Square(100), piece.New(piece.Black, piece.Pawn),
			squareChanges{},
		},
		{"White-None-A1", square.A1, piece.New(piece.White, piece.None),
			squareChanges{},
		},
		{"Black-None-A8", square.A8, piece.New(piece.Black, piece.None),
			squareChanges{},
		},
		{"White-Type(10)-A1", square.A1, piece.New(piece.White, piece.Type(10)),
			squareChanges{},
		},
		{"Black-Type(10)-A8", square.A8, piece.New(piece.Black, piece.Type(10)),
			squareChanges{},
		},
		{"NoColor-Pawn-A1", square.A1, piece.New(piece.NoColor, piece.Pawn),
			squareChanges{},
		},
		{"NoColor-Pawn-A8", square.A8, piece.New(piece.NoColor, piece.Pawn),
			squareChanges{},
		},
		{"Color(5)-Pawn-A1", square.A1, piece.New(piece.Color(5), piece.Pawn),
			squareChanges{},
		},
		{"Color(5)-Pawn-A8", square.A8, piece.New(piece.Color(5), piece.Pawn),
			squareChanges{},
		},
		{"NoColor-None-A1", square.A1, piece.New(piece.NoColor, piece.None),
			squareChanges{},
		},
		{"NoColor-None-A8", square.A8, piece.New(piece.NoColor, piece.None),
			squareChanges{},
		},
		{"NoColor-Type(10)-A1", square.A1, piece.New(piece.NoColor, piece.Type(10)),
			squareChanges{},
		},
		{"NoColor-Type(10)-A8", square.A8, piece.New(piece.NoColor, piece.Type(10)),
			squareChanges{},
		},
		{"Color(5)-None-A1", square.A1, piece.New(piece.Color(5), piece.None),
			squareChanges{},
		},
		{"Color(5)-None-A8", square.A8, piece.New(piece.Color(5), piece.None),
			squareChanges{},
		},
		{"Color(5)-Type(10)-A1", square.A1, piece.New(piece.Color(5), piece.Type(10)),
			squareChanges{},
		},
		{"Color(5)-Type(10)-A8", square.A8, piece.New(piece.Color(5), piece.Type(10)),
			squareChanges{},
		},
		{"NoColor-None-NoSquare", square.NoSquare, piece.New(piece.NoColor, piece.None),
			squareChanges{},
		},
		{"NoColor-None-Square(100)", square.Square(100), piece.New(piece.NoColor, piece.None),
			squareChanges{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := New()
			defer func() {
				if err := recover(); err != nil {
					t.Errorf("For initial board, *Position.QuickPut(%v, %v) should not panic, but panicked with: %v", tc.pc, tc.sq, err)
				}
			}()
			p.QuickPut(tc.pc, tc.sq)
			ch := changedSquares(New(), p)
			if !reflect.DeepEqual(ch, tc.wantChanged) {
				t.Errorf("For initial board and *Position.QuickPut(%v, %v), board changes are %v, want %v", tc.pc, tc.sq, ch, tc.wantChanged)
			}
		})
	}
}

func TestPutOnOccSquare(t *testing.T) {
	b := New()
	b.Clear()
	b.QuickPut(piece.New(piece.White, piece.Pawn), square.E2)
	b.Put(piece.New(piece.Black, piece.Queen), square.E2)
	if b.bitBoard[piece.White][piece.Pawn] != 0 {
		t.Fail()
	}
}

type testFindGroup struct {
	Name      string
	Position  testPosition
	TestCases []testFindTestCase
}

type testFindTestCase struct {
	Name  string
	Piece piece.Piece
	Want  map[square.Square]struct{}
}

func TestFind(t *testing.T) {
	testFindGroups := []testFindGroup{
		{
			"ZeroOrOneResult",
			testPosition{
				square.E1: piece.New(piece.White, piece.King),
				square.E2: piece.New(piece.White, piece.Pawn),
				square.E8: piece.New(piece.Black, piece.King),
				square.D7: piece.New(piece.Black, piece.Pawn),
			},
			[]testFindTestCase{
				{"WhiteQueen", piece.New(piece.White, piece.Queen), map[square.Square]struct{}{}},
				{"BlackKnight", piece.New(piece.Black, piece.Knight), map[square.Square]struct{}{}},
				{"WhiteBishop", piece.New(piece.White, piece.Bishop), map[square.Square]struct{}{}},
				{"BlackRook", piece.New(piece.Black, piece.Rook), map[square.Square]struct{}{}},
				{"WhiteKnight", piece.New(piece.White, piece.Knight), map[square.Square]struct{}{}},
				{"WhiteKing", piece.New(piece.White, piece.King), map[square.Square]struct{}{square.E1: struct{}{}}},
				{"WhitePawn", piece.New(piece.White, piece.Pawn), map[square.Square]struct{}{square.E2: struct{}{}}},
			},
		},
		{
			"TwoOrFourResults",
			testPosition{
				square.A2: piece.New(piece.White, piece.Pawn),
				square.B2: piece.New(piece.White, piece.Pawn),
				square.C2: piece.New(piece.White, piece.Pawn),
				square.D2: piece.New(piece.White, piece.Pawn),
				square.G7: piece.New(piece.Black, piece.Pawn),
				square.H7: piece.New(piece.Black, piece.Pawn),
			},
			[]testFindTestCase{
				{"WhitePawns", piece.New(piece.White, piece.Pawn), map[square.Square]struct{}{square.D2: struct{}{}, square.C2: struct{}{}, square.B2: struct{}{}, square.A2: struct{}{}}},
				{"BlackPawns", piece.New(piece.Black, piece.Pawn), map[square.Square]struct{}{square.H7: struct{}{}, square.G7: struct{}{}}},
			},
		},
		{
			"EightOrOneResults",
			testPosition(nil), // Initial chess board.
			[]testFindTestCase{
				{"WhitePawns", piece.New(piece.White, piece.Pawn), map[square.Square]struct{}{square.H2: struct{}{}, square.G2: struct{}{}, square.F2: struct{}{}, square.E2: struct{}{}, square.D2: struct{}{}, square.C2: struct{}{}, square.B2: struct{}{}, square.A2: struct{}{}}},
				{"BlackPawns", piece.New(piece.Black, piece.Pawn), map[square.Square]struct{}{square.H7: struct{}{}, square.G7: struct{}{}, square.F7: struct{}{}, square.E7: struct{}{}, square.D7: struct{}{}, square.C7: struct{}{}, square.B7: struct{}{}, square.A7: struct{}{}}},
				{"WhiteKing", piece.New(piece.White, piece.King), map[square.Square]struct{}{square.E1: struct{}{}}},
			},
		},
		{
			"FindNonStandardPieces",
			testPosition(nil), // Initial chess board.
			[]testFindTestCase{
				{"WhiteTypeNone", piece.New(piece.White, piece.None), map[square.Square]struct{}{}},
				{"WhiteType(piece.King+1)", piece.New(piece.White, piece.Type(piece.King+1)), map[square.Square]struct{}{}},
				{"BlackTypeNone", piece.New(piece.Black, piece.None), map[square.Square]struct{}{}},
				{"BlackType(piece.King+1)", piece.New(piece.Black, piece.Type(piece.King+1)), map[square.Square]struct{}{}},
				{"NoColorTypeNone", piece.New(piece.NoColor, piece.None), map[square.Square]struct{}{}},
				{"NoColorTypeKing", piece.New(piece.NoColor, piece.King), map[square.Square]struct{}{}},
				{"NoColorType(piece.King+1)", piece.New(piece.NoColor, piece.Type(piece.King+1)), map[square.Square]struct{}{}},
				{"Color(5)TypeNone", piece.New(piece.Color(5), piece.None), map[square.Square]struct{}{}},
				{"Color(5)TypeQueen", piece.New(piece.Color(5), piece.Queen), map[square.Square]struct{}{}},
				{"Color(5)Type(piece.King+1)", piece.New(piece.Color(5), piece.Type(piece.King+1)), map[square.Square]struct{}{}},
			},
		},
	}
	for _, group := range testFindGroups {
		t.Run(group.Name, func(t *testing.T) {
			// Get position.
			p, err := testCasePosition(group.Position, nil)
			if err != nil {
				t.Fatalf("Position preparation error: %s", err)
			}

			for _, tc := range group.TestCases {
				t.Run(tc.Name, func(t *testing.T) {
					// Call Find function with test case input on Position.
					s := p.Find(tc.Piece)

					// Compare results with expected values.
					if !reflect.DeepEqual(s, tc.Want) {
						t.Errorf("*Position.Find(%s): got \n%#v,\n\twant \n%#v\n", tc.Piece, s, tc.Want)
					}
				})
			}
		})
	}
}

// printBitBoards is a helper for diagnosing issues.
func (b *Position) printBitBoards() {
	for c := range b.bitBoard {
		for j := range b.bitBoard[c] {
			fmt.Println(piece.New(piece.Color(c), piece.Type(j)))
			fmt.Println(BitBoard(b.bitBoard[c][j]))
		}
	}
}

/*
// TODO(andrewbackes): add more advanced insufficient material checks.
func TestInsufMaterial(t *testing.T) {
	fens := []string{
		"8/8/4kb2/8/8/3K4/8/8 w - - 0 1",
		"8/8/4k3/8/6N1/3K4/8/8 w - - 0 1",
	}
	for _, fen := range fens {
		b, _ := FromFEN(fen)
		if b.InsufficientMaterial() != true {
			t.Fail()
		}
	}
}
*/

func TestKandBvKandOpB(t *testing.T) {
	b := New()
	b.Clear()
	b.QuickPut(piece.New(piece.White, piece.Bishop), square.A1)
	b.QuickPut(piece.New(piece.Black, piece.Bishop), square.B1)
	b.QuickPut(piece.New(piece.White, piece.King), square.A8)
	b.QuickPut(piece.New(piece.Black, piece.King), square.H1)
	if b.InsufficientMaterial() != false {
		t.Fail()
	}
}

func TestString(t *testing.T) {

	testCases := []struct {
		Name     string
		Position testPosition
		Changes  positionChanger
		Want     string
	}{
		{"InitialPosition", nil, nil, `   a b c d e f g h
 ┌─────────────────┐
8│ r n b q k b n r │8
7│ p p p p p p p p │7
6│ . . . . . . . . │6
5│ . . . . . . . . │5
4│ . . . . . . . . │4
3│ . . . . . . . . │3
2│ P P P P P P P P │2
1│ R N B Q K B N R │1
 └─────────────────┘
   a b c d e f g h

MoveNumber: 1
ActiveColor: White
CastlingRights:
  White: O-O-O O-O
  Black: O-O-O O-O
EnPassant:
LastMove:
FiftyMoveCount: 0
ThreeFoldCount: 0
MovesLeft:
  White: 0
  Black: 0
Clocks:
  White: 0s
  Black: 0s`,
		},
		{"InitialPosition-ShortPawnMove-e2e3", nil, makeMove("e2e3"), `   a b c d e f g h
 ┌─────────────────┐
8│ r n b q k b n r │8
7│ p p p p p p p p │7
6│ . . . . . . . . │6
5│ . . . . . . . . │5
4│ . . . . . . . . │4
3│ . . . . P . . . │3
2│ P P P P . P P P │2
1│ R N B Q K B N R │1
 └─────────────────┘
   a b c d e f g h

MoveNumber: 1
ActiveColor: Black
CastlingRights:
  White: O-O-O O-O
  Black: O-O-O O-O
EnPassant:
LastMove: e2e3
FiftyMoveCount: 0
ThreeFoldCount: 1
MovesLeft:
  White: -1
  Black: 0
Clocks:
  White: 0s
  Black: 0s`,
		},
		{"InitialPosition-LongPawnMove-e2e4", nil, makeMove("e2e4"), `   a b c d e f g h
 ┌─────────────────┐
8│ r n b q k b n r │8
7│ p p p p p p p p │7
6│ . . . . . . . . │6
5│ . . . . . . . . │5
4│ . . . . P . . . │4
3│ . . . . . . . . │3
2│ P P P P . P P P │2
1│ R N B Q K B N R │1
 └─────────────────┘
   a b c d e f g h

MoveNumber: 1
ActiveColor: Black
CastlingRights:
  White: O-O-O O-O
  Black: O-O-O O-O
EnPassant: e3
LastMove: e2e4
FiftyMoveCount: 0
ThreeFoldCount: 1
MovesLeft:
  White: -1
  Black: 0
Clocks:
  White: 0s
  Black: 0s`,
		},
		{"InitialPosition-2xLongPawnMove-e2e4-e7e5", nil,
			multi(
				makeMove("e2e4"),
				makeMove("e7e5"),
			),
			`   a b c d e f g h
 ┌─────────────────┐
8│ r n b q k b n r │8
7│ p p p p . p p p │7
6│ . . . . . . . . │6
5│ . . . . p . . . │5
4│ . . . . P . . . │4
3│ . . . . . . . . │3
2│ P P P P . P P P │2
1│ R N B Q K B N R │1
 └─────────────────┘
   a b c d e f g h

MoveNumber: 2
ActiveColor: White
CastlingRights:
  White: O-O-O O-O
  Black: O-O-O O-O
EnPassant: e6
LastMove: e7e5
FiftyMoveCount: 0
ThreeFoldCount: 1
MovesLeft:
  White: -1
  Black: -1
Clocks:
  White: 0s
  Black: 0s`,
		},
		{"InitialPosition-CastlingRights", nil,
			multi(
				castling(piece.White, board.ShortSide, false),
				castling(piece.White, board.LongSide, false),
				castling(piece.Black, board.LongSide, false),
			),
			`   a b c d e f g h
 ┌─────────────────┐
8│ r n b q k b n r │8
7│ p p p p p p p p │7
6│ . . . . . . . . │6
5│ . . . . . . . . │5
4│ . . . . . . . . │4
3│ . . . . . . . . │3
2│ P P P P P P P P │2
1│ R N B Q K B N R │1
 └─────────────────┘
   a b c d e f g h

MoveNumber: 1
ActiveColor: White
CastlingRights:
  White:
  Black: O-O
EnPassant:
LastMove:
FiftyMoveCount: 0
ThreeFoldCount: 0
MovesLeft:
  White: 0
  Black: 0
Clocks:
  White: 0s
  Black: 0s`,
		},
		{"InitialPosition-Clocks", nil,
			multi(
				clock(piece.White, 10*time.Minute),
				clock(piece.Black, 5*time.Second),
			),
			`   a b c d e f g h
 ┌─────────────────┐
8│ r n b q k b n r │8
7│ p p p p p p p p │7
6│ . . . . . . . . │6
5│ . . . . . . . . │5
4│ . . . . . . . . │4
3│ . . . . . . . . │3
2│ P P P P P P P P │2
1│ R N B Q K B N R │1
 └─────────────────┘
   a b c d e f g h

MoveNumber: 1
ActiveColor: White
CastlingRights:
  White: O-O-O O-O
  Black: O-O-O O-O
EnPassant:
LastMove:
FiftyMoveCount: 0
ThreeFoldCount: 0
MovesLeft:
  White: 0
  Black: 0
Clocks:
  White: 10m0s
  Black: 5s`,
		},
		{"InitialPosition-FiftyMoveCount-10", nil, func(p Position) (Position, error) {
			out := *Copy(&p)
			out.FiftyMoveCount = 10
			return out, nil
		},
			`   a b c d e f g h
 ┌─────────────────┐
8│ r n b q k b n r │8
7│ p p p p p p p p │7
6│ . . . . . . . . │6
5│ . . . . . . . . │5
4│ . . . . . . . . │4
3│ . . . . . . . . │3
2│ P P P P P P P P │2
1│ R N B Q K B N R │1
 └─────────────────┘
   a b c d e f g h

MoveNumber: 1
ActiveColor: White
CastlingRights:
  White: O-O-O O-O
  Black: O-O-O O-O
EnPassant:
LastMove:
FiftyMoveCount: 10
ThreeFoldCount: 0
MovesLeft:
  White: 0
  Black: 0
Clocks:
  White: 0s
  Black: 0s`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Get position.
			p, err := testCasePosition(tc.Position, tc.Changes)
			if err != nil {
				t.Fatalf("Position preparation error: %s", err)
			}

			if s := p.String(); s != tc.Want {
				t.Errorf("Position.String() =\n%s\n.\nWant\n%s\n.", s, tc.Want)
			}
		})
	}
}

func TestCapture(t *testing.T) {
	b := New()
	b.Clear()
	b.QuickPut(piece.New(piece.White, piece.King), square.E4)
	b.QuickPut(piece.New(piece.Black, piece.King), square.A8)
	b.QuickPut(piece.New(piece.Black, piece.Pawn), square.E5)
	result := b.MakeMove(move.Parse("e4e5"))
	for c := piece.White; c <= piece.Black; c++ {
		for p := piece.Pawn; p < piece.King; p++ {
			if result.bitBoard[c][p] != 0 {
				t.Log(BitBoard(result.bitBoard[c][p]))
				t.Fail()
			}
		}
	}
	if result.bitBoard[piece.White][piece.King] == 0 || result.bitBoard[piece.Black][piece.King] == 0 {
		t.Log("==0")
		t.Fail()
	}
}

func TestShortCastle(t *testing.T) {
	b := New()
	b.Clear()
	b.QuickPut(piece.New(piece.Black, piece.King), square.A8)
	b.QuickPut(piece.New(piece.White, piece.King), square.E1)
	b.QuickPut(piece.New(piece.White, piece.Rook), square.H1)
	result := b.MakeMove(move.Parse("e1g1"))
	if result.bitBoard[piece.White][piece.King] != (1<<square.G1) ||
		result.bitBoard[piece.White][piece.Rook] != (1<<square.F1) {
		t.Fail()
	}
}

func TestLongCastle(t *testing.T) {
	b := New()
	b.Clear()
	b.QuickPut(piece.New(piece.Black, piece.King), square.H8)
	b.QuickPut(piece.New(piece.White, piece.King), square.E1)
	b.QuickPut(piece.New(piece.White, piece.Rook), square.A1)
	result := b.MakeMove(move.Parse("e1c1"))
	if result.bitBoard[piece.White][piece.King] != (1<<square.C1) ||
		result.bitBoard[piece.White][piece.Rook] != (1<<square.D1) {
		fmt.Println(BitBoard(result.bitBoard[piece.White][piece.King]))
		fmt.Println("--Rook:--")
		fmt.Println(BitBoard(result.bitBoard[piece.White][piece.Rook]))
		t.Fail()
	}
}

func TestBitBoardPrint(t *testing.T) {
	b := New()
	expected := `00000000
00000000
00000000
00000000
00000000
00000000
11111111
00000000
`
	got := BitBoard(b.bitBoard[piece.White][piece.Pawn]).String()
	if got != expected {
		t.Error(got)
	}

}

func TestCheck(t *testing.T) {
	onlyWhiteCheckTestPosition := testPosition{
		square.E1: piece.New(piece.White, piece.King),
		square.E8: piece.New(piece.Black, piece.King),
		square.E7: piece.New(piece.Black, piece.Rook),
	}
	onlyBlackCheckTestPosition := testPosition{
		square.E1: piece.New(piece.White, piece.King),
		square.E2: piece.New(piece.White, piece.Rook),
		square.E8: piece.New(piece.Black, piece.King),
	}
	bothChecksTestPosition := testPosition{
		square.E1: piece.New(piece.White, piece.King),
		square.G6: piece.New(piece.White, piece.Bishop),
		square.E8: piece.New(piece.Black, piece.King),
		square.C3: piece.New(piece.Black, piece.Bishop),
	}

	testCases := []struct {
		Name     string
		Position testPosition
		col      piece.Color
		want     bool
	}{
		{"InitialTestPosition-White", InitialTestPosition, piece.White, false},
		{"InitialTestPosition-Black", InitialTestPosition, piece.Black, false},
		{"InitialTestPosition-NoColor", InitialTestPosition, piece.NoColor, false},
		{"InitialTestPosition-Color(5)", InitialTestPosition, piece.Color(5), false},
		{"onlyWhiteCheckTestPosition-White", onlyWhiteCheckTestPosition, piece.White, true},
		{"onlyWhiteCheckTestPosition-Black", onlyWhiteCheckTestPosition, piece.Black, false},
		{"onlyBlackCheckTestPosition-White", onlyBlackCheckTestPosition, piece.White, false},
		{"onlyBlackCheckTestPosition-Black", onlyBlackCheckTestPosition, piece.Black, true},
		{"bothChecksTestPosition-White", bothChecksTestPosition, piece.White, true},
		{"bothChecksTestPosition-Black", bothChecksTestPosition, piece.Black, true},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			p := tc.Position.Position()
			defer func() {
				if err := recover(); err != nil {
					t.Logf("Position:\n%v", p)
					t.Errorf("Position.Check(%#v) should not panic, but panicked with: %v", tc.col, err)
				}
			}()
			if ch := p.Check(tc.col); ch != tc.want {
				t.Logf("Position:\n%v", p)
				t.Errorf("Position.Check(%v) = %v, want %v", tc.col, ch, tc.want)
			}
		})
	}
}

func TestSAN(t *testing.T) {
	for _, group := range testSANGroups {
		t.Run(group.Name, func(t *testing.T) {
			for _, tc := range group.TestCases {
				t.Run(tc.Name, func(t *testing.T) {
					// Get position.
					p, err := testCasePosition(group.Position, tc.positionChangerFunc)
					if err != nil {
						t.Fatalf("Position preparation error: %s", err)
					}

					// Call ParseMove function with test case input on Position.
					s := p.SAN(move.Parse(tc.Move))

					// Compare results with expected values.
					if tc.Want != s {
						t.Errorf("*Position.SAN(%s): got \n%s,\n\twant \n%s\n", tc.Move, s, tc.Want)
					}
				})
			}
		})
	}
}

type testSANGroup struct {
	Name      string
	Position  testPosition
	TestCases []testSANTestCase
}

type testSANTestCase struct {
	Name                string
	positionChangerFunc positionChanger
	Move                string
	Want                string
}

var testSANGroups = []testSANGroup{
	{"Misc",
		map[square.Square]piece.Piece{
			square.E1: piece.New(piece.White, piece.King),
			square.F6: piece.New(piece.White, piece.Queen),
			square.A1: piece.New(piece.White, piece.Rook),
			square.D5: piece.New(piece.White, piece.Bishop),
			square.E4: piece.New(piece.White, piece.Pawn),

			square.E8: piece.New(piece.Black, piece.King),
			square.F3: piece.New(piece.Black, piece.Queen),
			square.H8: piece.New(piece.Black, piece.Rook),
			square.D4: piece.New(piece.Black, piece.Bishop),
			square.E5: piece.New(piece.White, piece.Pawn),

			// . . . . k . . r 8
			// . . . . . . . . 7
			// . . . . . Q . . 6
			// . . . B p . . . 5
			// . . . b P . . . 4
			// . . . . . q . . 3
			// . . . . . . . . 2
			// R . . . K . . . 1
			// a b c d e f g h
		},
		[]testSANTestCase{
			{"IllegalMove-White", active(piece.White), "e1e3", ""},
			{"IllegalMove-Black", active(piece.Black), "e8e6", ""},
			{"Move-White", active(piece.White), "e1d2", "Kd2"},
			{"Move-Black", active(piece.Black), "e8d7", "Kd7"},
			{"Move-Check-White", active(piece.White), "a1a8", "Ra8+"},
			{"Move-Check-Black", active(piece.Black), "h8h1", "Rh1+"},
			{"Move-Mate-White", active(piece.White), "d5c6", "Bc6#"},
			{"Move-Mate-Black", active(piece.Black), "d4c3", "Bc3#"},
		},
	},
	{"Promo",
		map[square.Square]piece.Piece{
			square.E1: piece.New(piece.White, piece.King),
			square.C1: piece.New(piece.White, piece.Rook),
			square.C7: piece.New(piece.White, piece.Rook),
			square.G2: piece.New(piece.White, piece.Rook),
			square.H1: piece.New(piece.White, piece.Rook),
			square.B7: piece.New(piece.White, piece.Pawn),
			square.D2: piece.New(piece.White, piece.Pawn),
			square.H7: piece.New(piece.White, piece.Pawn),

			square.E8: piece.New(piece.Black, piece.King),
			square.C2: piece.New(piece.Black, piece.Rook),
			square.C8: piece.New(piece.Black, piece.Rook),
			square.G7: piece.New(piece.Black, piece.Rook),
			square.H8: piece.New(piece.Black, piece.Rook),
			square.B2: piece.New(piece.Black, piece.Pawn),
			square.D7: piece.New(piece.Black, piece.Pawn),
			square.H2: piece.New(piece.Black, piece.Pawn),

			// . . r . k . . r 8
			// . P R p . . r P 7
			// . . . . . . . . 6
			// . . . . . . . . 5
			// . . . . . . . . 4
			// . . . . . . . . 3
			// . p r P . . R p 2
			// . . R . K . . R 1
			// a b c d e f g h
		},
		[]testSANTestCase{
			{"PawnMove-NoPromo-White", active(piece.White), "d2d3", "d3"},
			{"PawnMove-NoPromo-Black", active(piece.Black), "d7d6", "d6"},
			{"RookMove-NoPromo-White", active(piece.White), "g2g3", "Rg3"},
			{"RookMove-NoPromo-Black", active(piece.Black), "g7g6", "Rg6"},
			{"OpponentHomeRankMove-PawnMove-PromoQueen-White", active(piece.White), "b7b8q", "b8=Q"},
			{"OpponentHomeRankMove-PawnMove-PromoQueen-Black", active(piece.Black), "b2b1q", "b1=Q"},
			{"OpponentHomeRankMove-PawnMove-PromoRook-White", active(piece.White), "b7b8r", "b8=R"},
			{"OpponentHomeRankMove-PawnMove-PromoRook-Black", active(piece.Black), "b2b1r", "b1=R"},
			{"OpponentHomeRankMove-PawnMove-PromoBishop-White", active(piece.White), "b7b8b", "b8=B"},
			{"OpponentHomeRankMove-PawnMove-PromoBishop-Black", active(piece.Black), "b2b1b", "b1=B"},
			{"OpponentHomeRankMove-PawnMove-PromoKnight-White", active(piece.White), "b7b8n", "b8=N"},
			{"OpponentHomeRankMove-PawnMove-PromoKnight-Black", active(piece.Black), "b2b1n", "b1=N"},
			{"OpponentHomeRankMove-RookCaptureMove-NoPromo-Check-White", active(piece.White), "c7c8", "Rxc8+"},
			{"OpponentHomeRankMove-RookCaptureMove-NoPromo-Check-Black", active(piece.Black), "c2c1", "Rxc1+"},
			{"OpponentHomeRankMove-PawnCaptureMove-PromoQueen-Check-White", active(piece.White), "b7c8q", "bxc8=Q+"},
			{"OpponentHomeRankMove-PawnCaptureMove-PromoQueen-Check-Black", active(piece.Black), "b2c1q", "bxc1=Q+"},
			{"OpponentHomeRankMove-PawnCaptureMove-PromoRook-Check-White", active(piece.White), "b7c8r", "bxc8=R+"},
			{"OpponentHomeRankMove-PawnCaptureMove-PromoRook-Check-Black", active(piece.Black), "b2c1r", "bxc1=R+"},
			{"OpponentHomeRankMove-PawnCaptureMove-PromoBishop-White", active(piece.White), "b7c8b", "bxc8=B"},
			{"OpponentHomeRankMove-PawnCaptureMove-PromoBishop-Black", active(piece.Black), "b2c1b", "bxc1=B"},
			{"OpponentHomeRankMove-PawnCaptureMove-PromoKnight-White", active(piece.White), "b7c8n", "bxc8=N"},
			{"OpponentHomeRankMove-PawnCaptureMove-PromoKnight-Black", active(piece.Black), "b2c1n", "bxc1=N"},
		},
	},
	{"Pawn",
		map[square.Square]piece.Piece{
			square.G1: piece.New(piece.White, piece.King),
			square.F6: piece.New(piece.White, piece.Queen),
			square.D6: piece.New(piece.White, piece.Rook),
			square.B1: piece.New(piece.White, piece.Knight),
			square.B5: piece.New(piece.White, piece.Pawn),
			square.C4: piece.New(piece.White, piece.Pawn),
			square.E2: piece.New(piece.White, piece.Pawn),
			square.F2: piece.New(piece.White, piece.Pawn),
			square.G2: piece.New(piece.White, piece.Pawn),
			square.H3: piece.New(piece.White, piece.Pawn),
			square.H7: piece.New(piece.White, piece.Pawn),

			square.E8: piece.New(piece.Black, piece.King),
			square.F3: piece.New(piece.Black, piece.Queen),
			square.D3: piece.New(piece.Black, piece.Rook),
			square.G8: piece.New(piece.Black, piece.Knight),
			square.A2: piece.New(piece.Black, piece.Pawn),
			square.B4: piece.New(piece.Black, piece.Pawn),
			square.C5: piece.New(piece.Black, piece.Pawn),
			square.F7: piece.New(piece.Black, piece.Pawn),
			square.E7: piece.New(piece.Black, piece.Pawn),
			square.G7: piece.New(piece.Black, piece.Pawn),
			square.H6: piece.New(piece.Black, piece.Pawn),

			// . . . . k . n . 8
			// . . . . p p p P 7
			// . . . R . Q . p 6
			// . P p . . . . . 5
			// . p P . . . . . 4
			// . . . r . q . P 3
			// p . . . P P P . 2
			// . N . . . . K . 1
			// a b c d e f g h
		},
		[]testSANTestCase{
			{"Move-White", active(piece.White), "e2e3", "e3"},
			{"Move-Black", active(piece.Black), "e7e6", "e6"},
			{"LongMove-White", active(piece.White), "e2e4", "e4"},
			{"LongMove-Black", active(piece.Black), "e7e5", "e5"},
			{"CaptureMove-White", active(piece.White), "e2d3", "exd3"},
			{"CaptureMove-Black", active(piece.Black), "e7d6", "exd6"},
			{"CaptureMove-Ambiguous-White", active(piece.White), "g2f3", "gxf3"},
			{"CaptureMove-Ambiguous-Black", active(piece.Black), "e7f6", "exf6"},
			{"PromoMove-PromoQueen-White", active(piece.White), "h7h8q", "h8=Q"},
			{"PromoMove-PromoQueen-Black", active(piece.Black), "a2a1q", "a1=Q"},
			{"PromoMove-PromoRook-White", active(piece.White), "h7h8r", "h8=R"},
			{"PromoMove-PromoRook-Black", active(piece.Black), "a2a1r", "a1=R"},
			{"PromoMove-PromoBishop-White", active(piece.White), "h7h8b", "h8=B"},
			{"PromoMove-PromoBishop-Black", active(piece.Black), "a2a1b", "a1=B"},
			{"PromoMove-PromoKnight-White", active(piece.White), "h7h8n", "h8=N"},
			{"PromoMove-PromoKnight-Black", active(piece.Black), "a2a1n", "a1=N"},
			{"CapturePromoMove-PromoBishop-White", active(piece.White), "h7g8b", "hxg8=B"},
			{"CapturePromoMove-PromoKnight-Black", active(piece.Black), "a2b1n", "axb1=N"},
			{"CapturePromoMove-PromoQueen-Mate-White", active(piece.White), "h7g8q", "hxg8=Q#"},
			{"CapturePromoMove-PromoQueen-Check-Black", active(piece.Black), "a2b1q", "axb1=Q+"},
			{"CaptureMove-EnPassantOnC6-White", multi(active(piece.White), enPassant(square.C6)), "b5c6", "bxc6"},
			{"CaptureMove-EnPassantOnC3-Black", multi(active(piece.Black), enPassant(square.C3)), "b4c3", "bxc3"},
		},
	},
	{"Knight",
		map[square.Square]piece.Piece{
			square.H1: piece.New(piece.White, piece.King),
			square.A4: piece.New(piece.White, piece.Knight),
			square.B1: piece.New(piece.White, piece.Knight),
			square.B3: piece.New(piece.White, piece.Knight),
			square.B5: piece.New(piece.White, piece.Knight),
			square.B7: piece.New(piece.White, piece.Knight),
			square.C4: piece.New(piece.White, piece.Knight),
			square.D1: piece.New(piece.White, piece.Knight),
			square.D3: piece.New(piece.White, piece.Knight),
			square.D5: piece.New(piece.White, piece.Knight),
			square.G8: piece.New(piece.White, piece.Knight),
			square.H3: piece.New(piece.White, piece.Knight),

			square.H8: piece.New(piece.Black, piece.King),
			square.A5: piece.New(piece.Black, piece.Knight),
			square.B2: piece.New(piece.Black, piece.Knight),
			square.B4: piece.New(piece.Black, piece.Knight),
			square.B6: piece.New(piece.Black, piece.Knight),
			square.B8: piece.New(piece.Black, piece.Knight),
			square.C5: piece.New(piece.Black, piece.Knight),
			square.D2: piece.New(piece.Black, piece.Knight),
			square.D4: piece.New(piece.Black, piece.Knight),
			square.D8: piece.New(piece.Black, piece.Knight),
			square.G1: piece.New(piece.Black, piece.Knight),
			square.H6: piece.New(piece.Black, piece.Knight),

			// . n . n . . N k 8
			// . N . . . . . . 7
			// . n . . . . . n 6
			// n N n N . . . . 5
			// N n N n . . . . 4
			// . N . N . . . N 3
			// . n . n . . . . 2
			// . N . N . . n K 1
			// a b c d e f g h
		},
		[]testSANTestCase{
			{"Move-White", active(piece.White), "h3g5", "Ng5"},
			{"Move-Black", active(piece.Black), "h6g4", "Ng4"},
			{"AmbiguousMove-SpecifyFile-White", active(piece.White), "b5c7", "Nbc7"},
			{"AmbiguousMove-SpecifyFile-Black", active(piece.Black), "b4c2", "Nbc2"},
			{"AmbiguousMove-SpecifyRank-White", active(piece.White), "b1a3", "N1a3"},
			{"AmbiguousMove-SpecifyRank-Black", active(piece.Black), "b8a6", "N8a6"},
			{"AmbiguousMove-SpecifyFileRank-White", active(piece.White), "d1c3", "Nd1c3"},
			{"AmbiguousMove-SpecifyFileRank-Black", active(piece.Black), "d4c6", "Nd4c6"},
			{"CaptureMove-White", active(piece.White), "h3g1", "Nxg1"},
			{"CaptureMove-Black", active(piece.Black), "h6g8", "Nxg8"},
			{"AmbiguousCaptureMove-SpecifyFile-White", active(piece.White), "c4a5", "Ncxa5"},
			{"AmbiguousCaptureMove-SpecifyFile-Black", active(piece.Black), "c5a4", "Ncxa4"},
			{"AmbiguousCaptureMove-SpecifyRank-White", active(piece.White), "b3a5", "N3xa5"},
			{"AmbiguousCaptureMove-SpecifyRank-Black", active(piece.Black), "b2a4", "N2xa4"},
			{"AmbiguousCaptureMove-SpecifyFileRank-White", active(piece.White), "b3c5", "Nb3xc5"},
			{"AmbiguousCaptureMove-SpecifyFileRank-Black", active(piece.Black), "b2c4", "Nb2xc4"},
		},
	},
	{"Bishop",
		map[square.Square]piece.Piece{
			square.E1: piece.New(piece.White, piece.King),
			square.A1: piece.New(piece.White, piece.Bishop),
			square.A3: piece.New(piece.White, piece.Bishop),
			square.B5: piece.New(piece.White, piece.Bishop),
			square.B7: piece.New(piece.White, piece.Bishop),
			square.C1: piece.New(piece.White, piece.Bishop),
			square.C3: piece.New(piece.White, piece.Bishop),
			square.E7: piece.New(piece.White, piece.Bishop),
			square.F1: piece.New(piece.White, piece.Bishop),
			square.F3: piece.New(piece.White, piece.Bishop),
			square.H1: piece.New(piece.White, piece.Bishop),
			square.H3: piece.New(piece.White, piece.Bishop),

			square.E8: piece.New(piece.Black, piece.King),
			square.A6: piece.New(piece.Black, piece.Bishop),
			square.A8: piece.New(piece.Black, piece.Bishop),
			square.B2: piece.New(piece.Black, piece.Bishop),
			square.B4: piece.New(piece.Black, piece.Bishop),
			square.C6: piece.New(piece.Black, piece.Bishop),
			square.C8: piece.New(piece.Black, piece.Bishop),
			square.E2: piece.New(piece.Black, piece.Bishop),
			square.F6: piece.New(piece.Black, piece.Bishop),
			square.F8: piece.New(piece.Black, piece.Bishop),
			square.H6: piece.New(piece.Black, piece.Bishop),
			square.H8: piece.New(piece.Black, piece.Bishop),

			// b . b . k b . b 8
			// . B . . B . . . 7
			// b . b . . b . b 6
			// . B . . . . . . 5
			// . b . . . . . . 4
			// B . B . . B . B 3
			// . b . . b . . . 2
			// B . B . K B . B 1
			// a b c d e f g h
		},
		[]testSANTestCase{
			{"Move-White", active(piece.White), "c1e3", "Be3"},
			{"Move-Black", active(piece.Black), "c8e6", "Be6"},
			{"AmbiguousMove-SpecifyFile-White", active(piece.White), "f3g4", "Bfg4"},
			{"AmbiguousMove-SpecifyFile-Black", active(piece.Black), "f6g5", "Bfg5"},
			{"AmbiguousMove-SpecifyRank-White", active(piece.White), "c1d2", "B1d2"},
			{"AmbiguousMove-SpecifyRank-Black", active(piece.Black), "c8d7", "B8d7"},
			{"AmbiguousMove-SpecifyFileRank-White", active(piece.White), "h1g2", "Bh1g2"},
			{"AmbiguousMove-SpecifyFileRank-Black", active(piece.Black), "h8g7", "Bh8g7"},
			{"CaptureMove-White", active(piece.White), "e7f6", "Bxf6"},
			{"CaptureMove-Black", active(piece.Black), "e2f3", "Bxf3"},
			{"AmbiguousCaptureMove-SpecifyFile-White", active(piece.White), "a3b4", "Baxb4"},
			{"AmbiguousCaptureMove-SpecifyFile-Black", active(piece.Black), "a6b5", "Baxb5"},
			{"AmbiguousCaptureMove-SpecifyRank-White", active(piece.White), "f1e2", "B1xe2"},
			{"AmbiguousCaptureMove-SpecifyRank-Black", active(piece.Black), "f8e7", "B8xe7"},
			{"AmbiguousCaptureMove-SpecifyFileRank-White", active(piece.White), "a1b2", "Ba1xb2"},
			{"AmbiguousCaptureMove-SpecifyFileRank-Black", active(piece.Black), "a8b7", "Ba8xb7"},
		},
	},
	{"Rook",
		map[square.Square]piece.Piece{
			square.D1: piece.New(piece.White, piece.King),
			square.A2: piece.New(piece.White, piece.Rook),
			square.A6: piece.New(piece.White, piece.Rook),
			square.B1: piece.New(piece.White, piece.Rook),
			square.B3: piece.New(piece.White, piece.Rook),
			square.B5: piece.New(piece.White, piece.Rook),
			square.B7: piece.New(piece.White, piece.Rook),
			square.C2: piece.New(piece.White, piece.Rook),
			square.G2: piece.New(piece.White, piece.Rook),
			square.G4: piece.New(piece.White, piece.Rook),
			square.H3: piece.New(piece.White, piece.Rook),

			square.D8: piece.New(piece.Black, piece.King),
			square.A3: piece.New(piece.Black, piece.Rook),
			square.A7: piece.New(piece.Black, piece.Rook),
			square.B2: piece.New(piece.Black, piece.Rook),
			square.B4: piece.New(piece.Black, piece.Rook),
			square.B6: piece.New(piece.Black, piece.Rook),
			square.B8: piece.New(piece.Black, piece.Rook),
			square.C7: piece.New(piece.Black, piece.Rook),
			square.G5: piece.New(piece.Black, piece.Rook),
			square.G7: piece.New(piece.Black, piece.Rook),
			square.H6: piece.New(piece.Black, piece.Rook),

			// . r . k . . . . 8
			// r R r . . . r . 7
			// R r . . . . . r 6
			// . R . . . . r . 5
			// . r . . . . R . 4
			// r R . . . . . R 3
			// R r R . . . R . 2
			// . R . K . . . . 1
			// a b c d e f g h
		},
		[]testSANTestCase{
			{"Move-White", active(piece.White), "g2g1", "Rg1"},
			{"Move-Black", active(piece.Black), "g7g8", "Rg8"},
			{"AmbiguousMove-SpecifyFile-White", active(piece.White), "c2e2", "Rce2"},
			{"AmbiguousMove-SpecifyFile-Black", active(piece.Black), "c7e7", "Rce7"},
			{"AmbiguousMove-SpecifyRank-White", active(piece.White), "g2g3", "R2g3"},
			{"AmbiguousMove-SpecifyRank-Black", active(piece.Black), "g7g6", "R7g6"},
			{"CaptureMove-White", active(piece.White), "h3h6", "Rxh6"},
			{"CaptureMove-Black", active(piece.Black), "h6h3", "Rxh3"},
			{"AmbiguousCaptureMove-SpecifyFile-White", active(piece.White), "b3a3", "Rbxa3"},
			{"AmbiguousCaptureMove-SpecifyFile-Black", active(piece.Black), "b6a6", "Rbxa6"},
			{"AmbiguousCaptureMove-SpecifyRank-White", active(piece.White), "b3b4", "R3xb4"},
			{"AmbiguousCaptureMove-SpecifyRank-Black", active(piece.Black), "b4b5", "R4xb5"},
		},
	},
	{"Queen",
		map[square.Square]piece.Piece{
			square.E1: piece.New(piece.White, piece.King),
			square.A1: piece.New(piece.White, piece.Queen),
			square.A2: piece.New(piece.White, piece.Queen),
			square.B1: piece.New(piece.White, piece.Queen),
			square.G3: piece.New(piece.White, piece.Queen),
			square.G4: piece.New(piece.White, piece.Queen),
			square.H4: piece.New(piece.White, piece.Queen),
			square.A6: piece.New(piece.White, piece.Pawn),
			square.F6: piece.New(piece.White, piece.Pawn),
			square.H6: piece.New(piece.White, piece.Pawn),

			square.E8: piece.New(piece.Black, piece.King),
			square.A7: piece.New(piece.Black, piece.Queen),
			square.A8: piece.New(piece.Black, piece.Queen),
			square.B8: piece.New(piece.Black, piece.Queen),
			square.G5: piece.New(piece.Black, piece.Queen),
			square.G6: piece.New(piece.Black, piece.Queen),
			square.H5: piece.New(piece.Black, piece.Queen),
			square.A3: piece.New(piece.Black, piece.Pawn),
			square.F3: piece.New(piece.Black, piece.Pawn),
			square.H3: piece.New(piece.Black, piece.Pawn),

			// q q . . k . . . 8
			// q . . . . . . . 7
			// P . . . . P q P 6
			// . . . . . . q q 5
			// . . . . . . Q Q 4
			// p . . . . p Q p 3
			// Q . . . . . . . 2
			// Q Q . . K . . . 1
			// a b c d e f g h
		},
		[]testSANTestCase{
			{"Move-White", active(piece.White), "b1c1", "Qc1"},
			{"Move-Black", active(piece.Black), "b8c8", "Qc8"},
			{"AmbiguousMove-SpecifyFile-White", active(piece.White), "b1b2", "Qbb2"},
			{"AmbiguousMove-SpecifyFile-Black", active(piece.Black), "b8b7", "Qbb7"},
			{"AmbiguousMove-SpecifyRank-White", active(piece.White), "a2b2", "Q2b2"},
			{"AmbiguousMove-SpecifyRank-Black", active(piece.Black), "a7b7", "Q7b7"},
			{"AmbiguousMove-SpecifyFileRank-White", active(piece.White), "a1b2", "Qa1b2"},
			{"AmbiguousMove-SpecifyFileRank-Black", active(piece.Black), "a8b7", "Qa8b7"},
			{"CaptureMove-White", active(piece.White), "a2a3", "Qxa3"},
			{"CaptureMove-Black", active(piece.Black), "a7a6", "Qxa6"},
			{"AmbiguousCaptureMove-SpecifyFile-White", active(piece.White), "h4h3", "Qhxh3"},
			{"AmbiguousCaptureMove-SpecifyFile-Black", active(piece.Black), "h5h6", "Qhxh6"},
			{"AmbiguousCaptureMove-SpecifyRank-White", active(piece.White), "g3h3", "Q3xh3"},
			{"AmbiguousCaptureMove-SpecifyRank-Black", active(piece.Black), "g6h6", "Q6xh6"},
			{"AmbiguousCaptureMove-SpecifyFileRank-White", active(piece.White), "g4h3", "Qg4xh3"},
			{"AmbiguousCaptureMove-SpecifyFileRank-Black", active(piece.Black), "g5h6", "Qg5xh6"},
		},
	},
	{"King",
		map[square.Square]piece.Piece{
			square.E1: piece.New(piece.White, piece.King),
			square.A1: piece.New(piece.White, piece.Rook),
			square.H1: piece.New(piece.White, piece.Rook),

			square.E8: piece.New(piece.Black, piece.King),
			square.A8: piece.New(piece.Black, piece.Rook),
			square.H8: piece.New(piece.Black, piece.Rook),

			// r . . . k . . r 8
			// . . . . . . . . 7
			// . . . . . . . . 6
			// . . . . . . . . 5
			// . . . . . . . . 4
			// . . . . . . . . 3
			// . . . . . . . . 2
			// R . . . K . . R 1
			// a b c d e f g h
		},
		[]testSANTestCase{
			{"Castling-KingSide-White", active(piece.White), "e1g1", "O-O"},
			{"Castling-KingSide-Black", active(piece.Black), "e8g8", "O-O"},
			{"Castling-QueenSide-White", active(piece.White), "e1c1", "O-O-O"},
			{"Castling-QueenSide-Black", active(piece.Black), "e8c8", "O-O-O"},
			{"Move-White", active(piece.White), "e1d2", "Kd2"},
			{"Move-Black", active(piece.Black), "e8d7", "Kd7"},
			// Note: To capture, a new piece is added to the board, using the pos positionChangerFunc.
			{"CaptureMove-White", multi(pos(square.D2, piece.New(piece.Black, piece.Rook)), active(piece.White)), "e1d2", "Kxd2"},
			{"CaptureMove-Black", multi(pos(square.D7, piece.New(piece.White, piece.Rook)), active(piece.Black)), "e8d7", "Kxd7"},
		},
	},
}

// Various boards for testing purposes.
var (
	InitialTestPosition                   = testPosition(nil)
	TwoPawnsAtTwoOpositeKingsTestPosition = testPosition{
		square.E1: piece.New(piece.White, piece.King),
		square.E7: piece.New(piece.White, piece.Pawn),
		square.E8: piece.New(piece.Black, piece.King),
		square.E2: piece.New(piece.Black, piece.Pawn),
	}
	TwoKingsFourRooksTestPosition = testPosition{
		square.E1: piece.New(piece.White, piece.King),
		square.A1: piece.New(piece.White, piece.Rook),
		square.H1: piece.New(piece.White, piece.Rook),
		square.E8: piece.New(piece.Black, piece.King),
		square.A8: piece.New(piece.Black, piece.Rook),
		square.H8: piece.New(piece.Black, piece.Rook),
	}
	EnPassantCaptureTestPosition = testPosition{
		square.E1: piece.New(piece.White, piece.King),
		square.B5: piece.New(piece.White, piece.Pawn),
		square.F4: piece.New(piece.White, piece.Pawn),
		square.E8: piece.New(piece.Black, piece.King),
		square.C5: piece.New(piece.Black, piece.Pawn),
		square.G4: piece.New(piece.Black, piece.Pawn),
	}
	PromotionTestPosition = testPosition{
		square.E1: piece.New(piece.White, piece.King),
		square.A1: piece.New(piece.White, piece.Rook),
		square.B7: piece.New(piece.White, piece.Pawn),
		square.C6: piece.New(piece.White, piece.Pawn),
		square.G2: piece.New(piece.White, piece.Pawn),
		square.E8: piece.New(piece.Black, piece.King),
		square.A8: piece.New(piece.Black, piece.Rook),
		square.B2: piece.New(piece.Black, piece.Pawn),
		square.C3: piece.New(piece.Black, piece.Pawn),
		square.G7: piece.New(piece.Black, piece.Pawn),
	}
)

func BenchmarkPositionPut(b *testing.B) {
	newSquares := [square.LastSquare + 1]piece.Piece{}
	tpe, col := piece.Type(0), piece.Color(0)
	for sq := square.Square(0); sq <= square.LastSquare; sq += 1 {
		newSquares[sq] = piece.New(col, tpe)
		tpe, col = (tpe+1)%(piece.King+1), (col+1)%(piece.Black+1)
	}

	pos := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pos = New()
		for sq := square.Square(0); sq <= square.LastSquare; sq += 1 {
			pos.Put(newSquares[sq], sq)
		}
	}
}

func BenchmarkPositionQuickPut(b *testing.B) {
	newSquares := [square.LastSquare + 1]piece.Piece{}
	tpe, col := piece.Type(0), piece.Color(0)
	for sq := square.Square(0); sq <= square.LastSquare; sq += 1 {
		newSquares[sq] = piece.New(col, tpe)
		tpe, col = (tpe+1)%(piece.King+1), (col+1)%(piece.Black+1)
	}

	pos := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pos = New()
		for sq := square.Square(0); sq <= square.LastSquare; sq += 1 {
			pos.QuickPut(newSquares[sq], sq)
		}
	}
}
