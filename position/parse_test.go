package position

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
)

func TestParseMove(t *testing.T) {
	for _, group := range testParseMoveGroups {
		t.Run(group.Name, func(t *testing.T) {
			for _, tc := range group.TestCases {
				t.Run(tc.Name, func(t *testing.T) {
					// Get position.
					p, err := testCasePosition(group.Position, tc.positionChangerFunc)
					if err != nil {
						t.Fatalf("Position preparation error: %s", err)
					}

					// Call ParseMove function with test case input on Position.
					m, err := p.ParseMove(tc.Move)

					gotMoveString := ""
					if m != move.Null {
						gotMoveString = m.String()
					}

					// Compare results with expected values.
					if (tc.WantError == nil && gotMoveString != tc.Want) || !sameError(err, tc.WantError) {
						t.Errorf("*Position.ParseMove(%s):\n\tgot \n(%s, %#v),\n\twant \n(%s, %#v)\n ", tc.Move, m, err, tc.Want, tc.WantError)
						//t.Logf("got move: %s, %#v, move.Null: == %v, deepequal %v, %s %#v", m.String(), m, m == move.Null, reflect.DeepEqual(m, move.Null), move.Null.String(), move.Null)
					}
				})
			}
		})
	}
}

func sameError(errA, errB error) bool {
	if errA != nil && errB != nil {
		return errA.Error() == errB.Error()
	} else if errA == nil && errB == nil {
		return true
	}
	return false
}

func testCasePosition(p testPosition, pc positionChanger) (*Position, error) {
	res := p.Position()

	// Apply additional changes to Position.
	if pc != nil {
		if changedTestPosition, err := pc(*res); err != nil {
			return nil, err
		} else {
			res = &changedTestPosition
		}
	}
	return res, nil
}

type testParseMoveTestCases struct {
	Name      string
	Position  testPosition
	TestCases []testParseMoveTestCase
}

type testPosition map[square.Square]piece.Piece

// Returns a new `Position` structure filled with pieces on squares as defined in `p`.
func (p testPosition) Position() *Position {
	// Add pieces to new Position structure.
	res := New()
	res.Clear()
	for sq, pc := range p {
		res.QuickPut(pc, sq)
	}
	return res
}

type testParseMoveTestCase struct {
	Name                string
	positionChangerFunc positionChanger
	Move                string
	Want                string
	WantError           error
}

type positionChanger func(Position) (Position, error)

var testParseMoveGroups []testParseMoveTestCases = []testParseMoveTestCases{
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
		[]testParseMoveTestCase{
			{"InvalidMove-NullMove-White", active(piece.White), "0000", "", nil},
			{"InvalidMove-NullMove-Black", active(piece.Black), "0000", "", nil},
			{"InvalidMove-EmptyMove-White", active(piece.White), "", "", errors.New("could not parse ''")},
			{"InvalidMove-EmptyMove-Black", active(piece.White), "", "", errors.New("could not parse ''")},
			{"PCN-InvalidMove-FromOutOfBoundsToOutOfBoundsMove-White", active(piece.White), "j9i0", "", errors.New("could not parse 'j9i0'")},
			{"PCN-InvalidMove-FromOutOfBoundsToOutOfBoundsMove-Black", active(piece.Black), "j9i0", "", errors.New("could not parse 'j9i0'")},
			{"PCN-InvalidMove-FromOutOfBounds-White", active(piece.White), "j9e3", "", errors.New("could not parse 'j9e3'")},
			{"PCN-InvalidMove-FromOutOfBounds-Black", active(piece.Black), "j9e3", "", errors.New("could not parse 'j9e3'")},
			{"PCN-InvalidMove-ToOutOfBoundsMove-White", active(piece.White), "e2i0", "", errors.New("could not parse 'e2i0'")},
			{"PCN-InvalidMove-ToOutOfBoundsMove-Black", active(piece.Black), "e2i0", "", errors.New("could not parse 'e2i0'")},
			{"PCN-InvalidMove-White", active(piece.White), "d1c4", "d1c4", nil},
			{"PCN-InvalidMove-Black", active(piece.Black), "d1c4", "d1c4", nil},
			{"PCN-ValidMove-White", active(piece.White), "b1c3", "b1c3", nil},
			{"PCN-ValidMove-Black", active(piece.Black), "b1c3", "b1c3", nil},
			{"SAN-Move-White", active(piece.White), "Kd2", "e1d2", nil},
			{"SAN-Move-Black", active(piece.Black), "Kd7", "e8d7", nil},
			{"SAN-Move-InvalidCheck-White", active(piece.White), "Kd2+", "e1d2", nil},
			{"SAN-Move-InvalidCheck-Black", active(piece.Black), "Kd7+", "e8d7", nil},
			{"SAN-Move-InvalidMate-White", active(piece.White), "Kd2#", "e1d2", nil},
			{"SAN-Move-InvalidMate-Black", active(piece.Black), "Kd7#", "e8d7", nil},
			{"SAN-Move-Check-White", active(piece.White), "Ra8+", "a1a8", nil},
			{"SAN-Move-Check-Black", active(piece.Black), "Rh1+", "h8h1", nil},
			{"SAN-Move-InvalidMissingCheck-White", active(piece.White), "Ra8", "a1a8", nil},
			{"SAN-Move-InvalidMissingCheck-Black", active(piece.Black), "Rh1", "h8h1", nil},
			{"SAN-Move-Mate-White", active(piece.White), "Bc6#", "d5c6", nil},
			{"SAN-Move-Mate-Black", active(piece.Black), "Bc3#", "d4c3", nil},
			{"SAN-Move-InvalidMissingMate-White", active(piece.White), "Bc6", "d5c6", nil},
			{"SAN-Move-InvalidMissingMate-Black", active(piece.Black), "Bc3", "d4c3", nil},
			{"SAN-Move-InvalidSuffixCharacter-White", active(piece.White), "Kd2!", "", errors.New("could not parse 'Kd2!'")},
			{"SAN-Move-InvalidSuffixCharacter-Black", active(piece.Black), "Kd7@", "", errors.New("could not parse 'Kd7@'")},
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
		[]testParseMoveTestCase{
			{"PCN-PawnMove-NoPromo-White", active(piece.White), "d2d3", "d2d3", nil},
			{"PCN-PawnMove-NoPromo-Black", active(piece.Black), "d7d6", "d7d6", nil},
			{"PCN-PawnMove-InvalidPromo-White", active(piece.White), "d2d3Q", "d2d3q", nil},
			{"PCN-PawnMove-InvalidPromo-Black", active(piece.Black), "d7d6q", "d7d6q", nil},
			{"PCN-RookMove-NoPromo-White", active(piece.White), "g2g3", "g2g3", nil},
			{"PCN-RookMove-NoPromo-Black", active(piece.Black), "g7g6", "g7g6", nil},
			{"PCN-RookMove-InvalidPromo-White", active(piece.White), "g2g3Q", "g2g3q", nil},
			{"PCN-RookMove-InvalidPromo-Black", active(piece.Black), "g7g6q", "g7g6q", nil},
			{"PCN-InvalidMove-NoPromo-White", active(piece.White), "e3e2", "e3e2", nil},
			{"PCN-InvalidMove-NoPromo-Black", active(piece.Black), "e6e7", "e6e7", nil},
			{"PCN-InvalidMove-Promo-White", active(piece.White), "e3e2q", "e3e2q", nil},
			{"PCN-InvalidMove-Promo-Black", active(piece.Black), "e6e7Q", "e6e7q", nil},
			{"PCN-OpponentHomeRankMove-PawnMove-PromoQueen-White", active(piece.White), "b7b8q", "b7b8q", nil},
			{"PCN-OpponentHomeRankMove-PawnMove-PromoQueen-Black", active(piece.Black), "b2b1Q", "b2b1q", nil},
			{"PCN-OpponentHomeRankMove-PawnMove-PromoRook-White", active(piece.White), "b7b8r", "b7b8r", nil},
			{"PCN-OpponentHomeRankMove-PawnMove-PromoRook-Black", active(piece.Black), "b2b1R", "b2b1r", nil},
			{"PCN-OpponentHomeRankMove-PawnMove-PromoBishop-White", active(piece.White), "b7b8B", "b7b8b", nil},
			{"PCN-OpponentHomeRankMove-PawnMove-PromoBishop-Black", active(piece.Black), "b2b1b", "b2b1b", nil},
			{"PCN-OpponentHomeRankMove-PawnMove-PromoKnight-White", active(piece.White), "b7b8N", "b7b8n", nil},
			{"PCN-OpponentHomeRankMove-PawnMove-PromoKnight-Black", active(piece.Black), "b2b1n", "b2b1n", nil},
			{"PCN-OpponentHomeRankMove-PawnMove-InvalidNoPromo-White", active(piece.White), "b7b8", "b7b8q", nil},
			{"PCN-OpponentHomeRankMove-PawnMove-InvalidNoPromo-Black", active(piece.Black), "b2b1", "b2b1q", nil},
			{"PCN-OpponentHomeRankMove-PawnMove-InvalidDestinationOccupied-PromoQueen-White", active(piece.White), "h7h8q", "h7h8q", nil},
			{"PCN-OpponentHomeRankMove-PawnMove-InvalidDestinationOccupied-PromoQueen-Black", active(piece.Black), "h2h1Q", "h2h1q", nil},
			{"PCN-OpponentHomeRankMove-PawnMove-InvalidDestinationOccupied-NoPromo-White", active(piece.White), "h7h8", "h7h8q", nil},
			{"PCN-OpponentHomeRankMove-PawnMove-InvalidDestinationOccupied-NoPromo-Black", active(piece.Black), "h2h1", "h2h1q", nil},
			{"PCN-OpponentHomeRankMove-RookMove-InvalidPromoQueen-White", active(piece.White), "c7c8Q", "c7c8q", nil},
			{"PCN-OpponentHomeRankMove-RookMove-InvalidPromoQueen-Black", active(piece.Black), "c2c1q", "c2c1q", nil},
			{"PCN-OpponentHomeRankMove-RookMove-NoPromo-White", active(piece.White), "c7c8", "c7c8", nil},
			{"PCN-OpponentHomeRankMove-RookMove-NoPromo-Black", active(piece.Black), "c2c1", "c2c1", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-PromoQueen-White", active(piece.White), "b7c8q", "b7c8q", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-PromoQueen-Black", active(piece.Black), "b2c1Q", "b2c1q", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-PromoRook-White", active(piece.White), "b7c8r", "b7c8r", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-PromoRook-Black", active(piece.Black), "b2c1R", "b2c1r", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-PromoBishop-White", active(piece.White), "b7c8B", "b7c8b", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-PromoBishop-Black", active(piece.Black), "b2c1b", "b2c1b", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-PromoKnight-White", active(piece.White), "b7c8N", "b7c8n", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-PromoKnight-Black", active(piece.Black), "b2c1n", "b2c1n", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-InvalidNoPromo-White", active(piece.White), "b7c8", "b7c8q", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-InvalidNoPromo-Black", active(piece.Black), "b2c1", "b2c1q", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-InvalidDestinationNotOccupied-PromoQueen-White", active(piece.White), "b7a8q", "b7a8q", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-InvalidDestinationNotOccupied-PromoQueen-Black", active(piece.Black), "b2a1Q", "b2a1q", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-InvalidDestinationNotOccupied-PromoRook-White", active(piece.White), "b7a8r", "b7a8r", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-InvalidDestinationNotOccupied-PromoRook-Black", active(piece.Black), "b2a1R", "b2a1r", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-InvalidDestinationNotOccupied-PromoBishop-White", active(piece.White), "b7a8B", "b7a8b", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-InvalidDestinationNotOccupied-PromoBishop-Black", active(piece.Black), "b2a1b", "b2a1b", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-InvalidDestinationNotOccupied-PromoKnight-White", active(piece.White), "b7a8N", "b7a8n", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-InvalidDestinationNotOccupied-PromoKnight-Black", active(piece.Black), "b2a1n", "b2a1n", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-InvalidDestinationNotOccupied-InvalidNoPromo-White", active(piece.White), "b7a8", "b7a8q", nil},
			{"PCN-OpponentHomeRankMove-PawnCaptureMove-InvalidDestinationNotOccupied-InvalidNoPromo-Black", active(piece.Black), "b2a1", "b2a1q", nil},
			{"PCN-OpponentHomeRankMove-InvalidSourceNoPiece-PromoQueen-White", active(piece.White), "a7a8Q", "a7a8q", nil},
			{"PCN-OpponentHomeRankMove-InvalidSourceNoPiece-PromoQueen-Black", active(piece.Black), "a2a1q", "a2a1q", nil},
			{"PCN-OpponentHomeRankMove-InvalidSourceNoPiece-NoPromo-White", active(piece.White), "a7a8", "a7a8", nil},
			{"PCN-OpponentHomeRankMove-InvalidSourceNoPiece-NoPromo-Black", active(piece.Black), "a2a1", "a2a1", nil},
			{"PCN-OpponentHomeRankMove-InvalidSourceOpponentPawn-PromoQueen-White", active(piece.White), "d7d8q", "d7d8q", nil},
			{"PCN-OpponentHomeRankMove-InvalidSourceOpponentPawn-PromoQueen-Black", active(piece.Black), "d2d1Q", "d2d1q", nil},
			{"PCN-OpponentHomeRankMove-InvalidSourceOpponentPawn-NoPromo-White", active(piece.White), "d7d8", "d7d8", nil},
			{"PCN-OpponentHomeRankMove-InvalidSourceOpponentPawn-NoPromo-Black", active(piece.Black), "d2d1", "d2d1", nil},
			{"PCN-OpponentHomeRankMove-InvalidSourceOpponentRook-PromoQueen-White", active(piece.White), "g7g8q", "g7g8q", nil},
			{"PCN-OpponentHomeRankMove-InvalidSourceOpponentRook-PromoQueen-Black", active(piece.Black), "g2g1Q", "g2g1q", nil},
			{"PCN-OpponentHomeRankMove-InvalidSourceOpponentRook-NoPromo-White", active(piece.White), "g7g8", "g7g8", nil},
			{"PCN-OpponentHomeRankMove-InvalidSourceOpponentRook-NoPromo-Black", active(piece.Black), "g2g1", "g2g1", nil},
			{"SAN-PawnMove-NoPromo-White", active(piece.White), "d3", "d2d3", nil},
			{"SAN-PawnMove-NoPromo-Black", active(piece.Black), "d6", "d7d6", nil},
			{"SAN-PawnMove-InvalidPromo-White", active(piece.White), "d3=Q", "d2d3q", nil},
			{"SAN-PawnMove-InvalidPromo-Black", active(piece.Black), "d6=q", "d7d6q", nil},
			{"SAN-RookMove-NoPromo-White", active(piece.White), "Rg3", "g2g3", nil},
			{"SAN-RookMove-NoPromo-Black", active(piece.Black), "Rg6", "g7g6", nil},
			{"SAN-RookMove-InvalidPromo-White", active(piece.White), "Rg3=Q", "g2g3q", nil},
			{"SAN-RookMove-InvalidPromo-Black", active(piece.Black), "Rg6=q", "g7g6q", nil},
			{"SAN-OpponentHomeRankMove-PawnMove-PromoQueen-White", active(piece.White), "b8=q", "b7b8q", nil},
			{"SAN-OpponentHomeRankMove-PawnMove-PromoQueen-Black", active(piece.Black), "b1=Q", "b2b1q", nil},
			{"SAN-OpponentHomeRankMove-PawnMove-PromoRook-White", active(piece.White), "b8=r", "b7b8r", nil},
			{"SAN-OpponentHomeRankMove-PawnMove-PromoRook-Black", active(piece.Black), "b1=R", "b2b1r", nil},
			{"SAN-OpponentHomeRankMove-PawnMove-PromoBishop-White", active(piece.White), "b8=B", "b7b8b", nil},
			{"SAN-OpponentHomeRankMove-PawnMove-PromoBishop-Black", active(piece.Black), "b1=b", "b2b1b", nil},
			{"SAN-OpponentHomeRankMove-PawnMove-PromoKnight-White", active(piece.White), "b8=N", "b7b8n", nil},
			{"SAN-OpponentHomeRankMove-PawnMove-PromoKnight-Black", active(piece.Black), "b1=n", "b2b1n", nil},
			{"SAN-OpponentHomeRankMove-PawnMove-InvalidNoPromo-White", active(piece.White), "b8", "b7b8q", nil},
			{"SAN-OpponentHomeRankMove-PawnMove-InvalidNoPromo-Black", active(piece.Black), "b1", "b2b1q", nil},
			{"SAN-OpponentHomeRankMove-PawnMove-InvalidDestinationOccupied-PromoQueen-White", active(piece.White), "h8=q", "", errors.New("could not find source square of 'h8=q'")},
			{"SAN-OpponentHomeRankMove-PawnMove-InvalidDestinationOccupied-PromoQueen-Black", active(piece.Black), "h1=Q", "", errors.New("could not find source square of 'h1=Q'")},
			{"SAN-OpponentHomeRankMove-PawnMove-InvalidDestinationOccupied-NoPromo-White", active(piece.White), "h8", "", errors.New("could not find source square of 'h8'")},
			{"SAN-OpponentHomeRankMove-PawnMove-InvalidDestinationOccupied-NoPromo-Black", active(piece.Black), "h1", "", errors.New("could not find source square of 'h1'")},
			{"SAN-OpponentHomeRankMove-RookMove-InvalidPromoQueen-InvalidNoCapture-White", active(piece.White), "Rc8=Q", "c7c8q", nil},
			{"SAN-OpponentHomeRankMove-RookMove-InvalidPromoQueen-InvalidNoCapture-Black", active(piece.Black), "Rc1=q", "c2c1q", nil},
			{"SAN-OpponentHomeRankMove-RookCaptureMove-InvalidPromoQueen-White", active(piece.White), "Rxc8=Q", "c7c8q", nil},
			{"SAN-OpponentHomeRankMove-RookCaptureMove-InvalidPromoQueen-Black", active(piece.Black), "Rxc1=q", "c2c1q", nil},
			{"SAN-OpponentHomeRankMove-RookMove-NoPromo-InvalidNoCapture-White", active(piece.White), "Rc8", "c7c8", nil},
			{"SAN-OpponentHomeRankMove-RookMove-NoPromo-InvalidNoCapture-Black", active(piece.Black), "Rc1", "c2c1", nil},
			{"SAN-OpponentHomeRankMove-RookCaptureMove-NoPromo-Check-White", active(piece.White), "Rxc8+", "c7c8", nil},
			{"SAN-OpponentHomeRankMove-RookCaptureMove-NoPromo-Check-Black", active(piece.Black), "Rxc1+", "c2c1", nil},
			{"SAN-OpponentHomeRankMove-RookCaptureMove-NoPromo-InvalidMissingCheck-White", active(piece.White), "Rxc8", "c7c8", nil},
			{"SAN-OpponentHomeRankMove-RookCaptureMove-NoPromo-InvalidMissingCheck-Black", active(piece.Black), "Rxc1", "c2c1", nil},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-PromoQueen-Check-White", active(piece.White), "bxc8=q+", "b7c8q", nil},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-PromoQueen-Check-Black", active(piece.Black), "bxc1=Q+", "b2c1q", nil},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-PromoQueen-InvalidMissingCheck-White", active(piece.White), "bxc8=q", "b7c8q", nil},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-PromoQueen-InvalidMissingCheck-Black", active(piece.Black), "bxc1=Q", "b2c1q", nil},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-PromoQueen-InvalidMate-White", active(piece.White), "bxc8=q#", "b7c8q", nil},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-PromoQueen-InvalidMate-Black", active(piece.Black), "bxc1=Q#", "b2c1q", nil},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-PromoRook-Check-White", active(piece.White), "bxc8=r+", "b7c8r", nil},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-PromoRook-Check-Black", active(piece.Black), "bxc1=R+", "b2c1r", nil},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-PromoRook-InvalidMissingCheck-White", active(piece.White), "bxc8=r", "b7c8r", nil},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-PromoRook-InvalidMissingCheck-Black", active(piece.Black), "bxc1=R", "b2c1r", nil},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-PromoBishop-White", active(piece.White), "bxc8=B", "b7c8b", nil},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-PromoBishop-Black", active(piece.Black), "bxc1=b", "b2c1b", nil},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-PromoKnight-White", active(piece.White), "bxc8=N", "b7c8n", nil},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-PromoKnight-Black", active(piece.Black), "bxc1=n", "b2c1n", nil},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-InvalidNoPromo-White", active(piece.White), "bxc8", "b7c8q", nil},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-InvalidNoPromo-Black", active(piece.Black), "bxc1", "b2c1q", nil},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-InvalidDestinationNotOccupied-PromoQueen-White", active(piece.White), "bxa8=q", "", errors.New("could not find source square of 'bxa8=q'")},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-InvalidDestinationNotOccupied-PromoQueen-Black", active(piece.Black), "bxa1=Q", "", errors.New("could not find source square of 'bxa1=Q'")},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-InvalidDestinationNotOccupied-InvalidNoPromo-White", active(piece.White), "bxa8", "", errors.New("could not find source square of 'bxa8'")},
			{"SAN-OpponentHomeRankMove-PawnCaptureMove-InvalidDestinationNotOccupied-InvalidNoPromo-Black", active(piece.Black), "bxa1", "", errors.New("could not find source square of 'bxa1'")},
			{"SAN-OpponentHomeRankMove-InvalidSourceNoPiece-PromoQueen-White", active(piece.White), "a8=Q", "", errors.New("could not find source square of 'a8=Q'")},
			{"SAN-OpponentHomeRankMove-InvalidSourceNoPiece-PromoQueen-Black", active(piece.Black), "a1=q", "", errors.New("could not find source square of 'a1=q'")},
			{"SAN-OpponentHomeRankMove-InvalidSourceNoPiece-NoPromo-White", active(piece.White), "a8", "", errors.New("could not find source square of 'a8'")},
			{"SAN-OpponentHomeRankMove-InvalidSourceNoPiece-NoPromo-Black", active(piece.Black), "a1", "", errors.New("could not find source square of 'a1'")},
			{"SAN-OpponentHomeRankMove-InvalidSourceOpponentPawn-PromoQueen-White", active(piece.White), "d8=q", "", errors.New("could not find source square of 'd8=q'")},
			{"SAN-OpponentHomeRankMove-InvalidSourceOpponentPawn-PromoQueen-Black", active(piece.Black), "d1=Q", "", errors.New("could not find source square of 'd1=Q'")},
			{"SAN-OpponentHomeRankMove-InvalidSourceOpponentPawn-NoPromo-White", active(piece.White), "d8", "", errors.New("could not find source square of 'd8'")},
			{"SAN-OpponentHomeRankMove-InvalidSourceOpponentPawn-NoPromo-Black", active(piece.Black), "d1", "", errors.New("could not find source square of 'd1'")},
			{"SAN-OpponentHomeRankMove-InvalidSourceOpponentRook-PromoQueen-White", active(piece.White), "g8=q", "", errors.New("could not find source square of 'g8=q'")},
			{"SAN-OpponentHomeRankMove-InvalidSourceOpponentRook-PromoQueen-Black", active(piece.Black), "g1=Q", "", errors.New("could not find source square of 'g1=Q'")},
			{"SAN-OpponentHomeRankMove-InvalidSourceOpponentRook-NoPromo-White", active(piece.White), "g8", "", errors.New("could not find source square of 'g8'")},
			{"SAN-OpponentHomeRankMove-InvalidSourceOpponentRook-NoPromo-Black", active(piece.Black), "g1", "", errors.New("could not find source square of 'g1'")},
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
		[]testParseMoveTestCase{
			{"SAN-Move-White", active(piece.White), "e3", "e2e3", nil},
			{"SAN-Move-Black", active(piece.Black), "e6", "e7e6", nil},
			{"SAN-Move-InvalidRedundantType-White", active(piece.White), "Pe3", "e2e3", nil},
			{"SAN-Move-InvalidRedundantType-Black", active(piece.Black), "Pe6", "e7e6", nil},
			{"SAN-Move-InvalidRedundantFile-White", active(piece.White), "ee3", "e2e3", nil},
			{"SAN-Move-InvalidRedundantFile-Black", active(piece.Black), "ee6", "e7e6", nil},
			{"SAN-Move-InvalidRedundantRank-White", active(piece.White), "2e3", "e2e3", nil},
			{"SAN-Move-InvalidRedundantRank-Black", active(piece.Black), "7e6", "e7e6", nil},
			{"SAN-Move-InvalidRedundantTypeFileRank-White", active(piece.White), "Pe2e3", "e2e3", nil},
			{"SAN-Move-InvalidRedundantTypeFileRank-Black", active(piece.Black), "Pe7e6", "e7e6", nil},
			{"SAN-Move-InvalidWrongRedundantType-White", active(piece.White), "Ke3", "", errors.New("could not find source square of 'Ke3'")},
			{"SAN-Move-InvalidWrongRedundantType-Black", active(piece.Black), "Ke6", "", errors.New("could not find source square of 'Ke6'")},
			{"SAN-Move-InvalidWrongRedundantFile-White", active(piece.White), "ae3", "", errors.New("could not find source square of 'ae3'")},
			{"SAN-Move-InvalidWrongRedundantFile-Black", active(piece.Black), "ae6", "", errors.New("could not find source square of 'ae6'")},
			{"SAN-Move-InvalidWrongRedundantRank-White", active(piece.White), "1e3", "", errors.New("could not find source square of '1e3'")},
			{"SAN-Move-InvalidWrongRedundantRank-Black", active(piece.Black), "8e6", "", errors.New("could not find source square of '8e6'")},
			{"SAN-Move-InvalidWrongRedundantTypeFileRank-White", active(piece.White), "Ka1e3", "a1e3", nil},
			{"SAN-Move-InvalidWrongRedundantTypeFileRank-Black", active(piece.Black), "Ka8e6", "a8e6", nil},
			{"SAN-Move-InvalidNoPawnOrigin-White", active(piece.White), "a6", "", errors.New("could not find source square of 'a6'")},
			{"SAN-Move-InvalidNoPawnOrigin-Black", active(piece.Black), "a3", "", errors.New("could not find source square of 'a3'")},
			{"SAN-Move-InvalidCapture-White", active(piece.White), "xe3", "e2e3", nil},
			{"SAN-Move-InvalidCapture-Black", active(piece.Black), "xe6", "e7e6", nil},
			{"SAN-Move-InvalidFile-White", active(piece.White), "de3", "", errors.New("could not find source square of 'de3'")},
			{"SAN-Move-InvalidFile-Black", active(piece.Black), "de6", "", errors.New("could not find source square of 'de6'")},
			{"SAN-Move-InvalidRank-White", active(piece.White), "1e3", "", errors.New("could not find source square of '1e3'")},
			{"SAN-Move-InvalidRank-Black", active(piece.Black), "8e6", "", errors.New("could not find source square of '8e6'")},
			{"SAN-Move-InvalidFileRank-White", active(piece.White), "Pe3e3", "e3e3", nil},
			{"SAN-Move-InvalidFileRank-Black", active(piece.Black), "Pe6e6", "e6e6", nil},
			{"SAN-Move-InvalidPromo-White", active(piece.White), "e3=Q", "e2e3q", nil},
			{"SAN-Move-InvalidPromo-Black", active(piece.Black), "e6=q", "e7e6q", nil},
			{"SAN-Move-InvalidCheck-White", active(piece.White), "e3+", "e2e3", nil},
			{"SAN-Move-InvalidMate-Black", active(piece.Black), "e6#", "e7e6", nil},
			{"SAN-Move-InvalidPromoCheck-White", active(piece.White), "e3=Q+", "e2e3q", nil},
			{"SAN-Move-InvalidPromoMate-Black", active(piece.Black), "e6=q#", "e7e6q", nil},
			{"SAN-LongMove-White", active(piece.White), "e4", "e2e4", nil},
			{"SAN-LongMove-Black", active(piece.Black), "e5", "e7e5", nil},
			{"SAN-LongMove-InvalidRedundantType-White", active(piece.White), "Pe4", "e2e4", nil},
			{"SAN-LongMove-InvalidRedundantType-Black", active(piece.Black), "Pe5", "e7e5", nil},
			{"SAN-LongMove-InvalidRedundantFile-White", active(piece.White), "ee4", "e2e4", nil},
			{"SAN-LongMove-InvalidRedundantFile-Black", active(piece.Black), "ee5", "e7e5", nil},
			{"SAN-LongMove-InvalidRedundantRank-White", active(piece.White), "2e4", "e2e4", nil},
			{"SAN-LongMove-InvalidRedundantRank-Black", active(piece.Black), "7e5", "e7e5", nil},
			{"SAN-LongMove-InvalidRedundantTypeFileRank-White", active(piece.White), "Pe2e4", "e2e4", nil},
			{"SAN-LongMove-InvalidRedundantTypeFileRank-Black", active(piece.Black), "Pe7e5", "e7e5", nil},
			{"SAN-LongMove-InvalidWrongRedundantType-White", active(piece.White), "Ke4", "", errors.New("could not find source square of 'Ke4'")},
			{"SAN-LongMove-InvalidWrongRedundantType-Black", active(piece.Black), "Ke5", "", errors.New("could not find source square of 'Ke5'")},
			{"SAN-LongMove-InvalidWrongRedundantFile-White", active(piece.White), "ae4", "", errors.New("could not find source square of 'ae4'")},
			{"SAN-LongMove-InvalidWrongRedundantFile-Black", active(piece.Black), "ae5", "", errors.New("could not find source square of 'ae5'")},
			{"SAN-LongMove-InvalidWrongRedundantRank-White", active(piece.White), "1e4", "", errors.New("could not find source square of '1e4'")},
			{"SAN-LongMove-InvalidWrongRedundantRank-Black", active(piece.Black), "8e5", "", errors.New("could not find source square of '8e5'")},
			{"SAN-LongMove-InvalidWrongRedundantTypeFileRank-White", active(piece.White), "Ka1e4", "a1e4", nil},
			{"SAN-LongMove-InvalidWrongRedundantTypeFileRank-Black", active(piece.Black), "Ka8e5", "a8e5", nil},
			{"SAN-LongMove-InvalidNoPawn-White", active(piece.White), "d4", "", errors.New("could not find source square of 'd4'")},
			{"SAN-LongMove-InvalidNoPawn-Black", active(piece.Black), "d5", "", errors.New("could not find source square of 'd5'")},
			{"SAN-LongMove-InvalidPawnBlocked-White", active(piece.White), "f4", "", errors.New("could not find source square of 'f4'")},
			{"SAN-LongMove-InvalidPawnBlocked-Black", active(piece.Black), "f5", "", errors.New("could not find source square of 'f5'")},
			{"SAN-LongMove-InvalidPawnNotHome-White", active(piece.White), "h5", "", errors.New("could not find source square of 'h5'")},
			{"SAN-LongMove-InvalidPawnNotHome-Black", active(piece.Black), "h4", "", errors.New("could not find source square of 'h4'")},
			{"SAN-LongMove-InvalidPromo-White", active(piece.White), "e4=R", "e2e4r", nil},
			{"SAN-LongMove-InvalidPromo-Black", active(piece.Black), "e5=r", "e7e5r", nil},
			{"SAN-LongMove-InvalidMate-White", active(piece.White), "e4#", "e2e4", nil},
			{"SAN-LongMove-InvalidCheck-Black", active(piece.Black), "e5+", "e7e5", nil},
			{"SAN-LongMove-InvalidPromoMate-White", active(piece.White), "e4=R#", "e2e4r", nil},
			{"SAN-LongMove-InvalidPromoCheck-Black", active(piece.Black), "e5=r+", "e7e5r", nil},
			{"SAN-LongMove-InvalidCapture-White", active(piece.White), "xe4", "e2e4", nil},
			{"SAN-LongMove-InvalidCapture-Black", active(piece.Black), "xe5", "e7e5", nil},
			{"SAN-CaptureMove-White", active(piece.White), "exd3", "e2d3", nil},
			{"SAN-CaptureMove-Black", active(piece.Black), "exd6", "e7d6", nil},
			{"SAN-CaptureMove-InvalidMissingOrigin-White", active(piece.White), "xd3", "e2d3", nil},
			{"SAN-CaptureMove-InvalidMissingOrigin-Black", active(piece.Black), "xd6", "e7d6", nil},
			{"SAN-CaptureMove-InvalidMissingCapture-White", active(piece.White), "ed3", "e2d3", nil},
			{"SAN-CaptureMove-InvalidMissingCapture-Black", active(piece.Black), "ed6", "e7d6", nil},
			{"SAN-CaptureMove-InvalidMissingOriginCapture-White", active(piece.White), "d3", "e2d3", nil},
			{"SAN-CaptureMove-InvalidMissingOriginCapture-Black", active(piece.Black), "d6", "e7d6", nil},
			{"SAN-CaptureMove-InvalidRedundantType-White", active(piece.White), "Pexd3", "e2d3", nil},
			{"SAN-CaptureMove-InvalidRedundantType-Black", active(piece.Black), "Pexd6", "e7d6", nil},
			{"SAN-CaptureMove-InvalidRedundantRank-White", active(piece.White), "e2xd3", "e2d3", nil},
			{"SAN-CaptureMove-InvalidRedundantRank-Black", active(piece.Black), "e7xd6", "e7d6", nil},
			{"SAN-CaptureMove-InvalidRedundantTypeRank-White", active(piece.White), "Pe2xd3", "e2d3", nil},
			{"SAN-CaptureMove-InvalidRedundantTypeRank-Black", active(piece.Black), "Pe7xd6", "e7d6", nil},
			{"SAN-CaptureMove-InvalidWrongRedundantType-White", active(piece.White), "Kexd3", "", errors.New("could not find source square of 'Kexd3'")},
			{"SAN-CaptureMove-InvalidWrongRedundantType-Black", active(piece.Black), "Kexd6", "", errors.New("could not find source square of 'Kexd6'")},
			{"SAN-CaptureMove-InvalidWrongRedundantRank-White", active(piece.White), "a2xd3", "a2d3", nil},
			{"SAN-CaptureMove-InvalidWrongRedundantRank-Black", active(piece.Black), "a7xd6", "a7d6", nil},
			{"SAN-CaptureMove-InvalidWrongRedundantTypeRank-White", active(piece.White), "Ka2xd3", "a2d3", nil},
			{"SAN-CaptureMove-InvalidWrongRedundantTypeRank-Black", active(piece.Black), "Ka7xd6", "a7d6", nil},
			{"SAN-CaptureMove-Ambiguous-White", active(piece.White), "gxf3", "g2f3", nil},
			{"SAN-CaptureMove-Ambiguous-Black", active(piece.Black), "exf6", "e7f6", nil},
			{"SAN-CaptureMove-Ambiguous-InvalidMissingOrigin-White", active(piece.White), "xf3", "", errors.New("could not find source square of 'xf3'")},
			{"SAN-CaptureMove-Ambiguous-InvalidMissingOrigin-Black", active(piece.Black), "xf6", "", errors.New("could not find source square of 'xf6'")},
			{"SAN-CaptureMove-Ambiguous-InvalidMissingCapture-White", active(piece.White), "gf3", "g2f3", nil},
			{"SAN-CaptureMove-Ambiguous-InvalidMissingCapture-Black", active(piece.Black), "ef6", "e7f6", nil},
			{"SAN-CaptureMove-Ambiguous-InvalidMissingOriginAndCapture-White", active(piece.White), "f3", "", errors.New("could not find source square of 'f3'")},
			{"SAN-CaptureMove-Ambiguous-InvalidMissingOriginAndCapture-Black", active(piece.Black), "f6", "", errors.New("could not find source square of 'f6'")},
			{"SAN-CaptureMove-InvalidNoEnemyPieceToCapture-White", active(piece.White), "gxh3", "", errors.New("could not find source square of 'gxh3'")},
			{"SAN-CaptureMove-InvalidNoEnemyPieceToCapture-Black", active(piece.Black), "gxh6", "", errors.New("could not find source square of 'gxh6'")},
			{"SAN-CaptureMove-InvalidNoOriginPawn-White", active(piece.White), "cxd3", "", errors.New("could not find source square of 'cxd3'")},
			{"SAN-CaptureMove-InvalidNoOriginPawn-Black", active(piece.Black), "cxd6", "", errors.New("could not find source square of 'cxd6'")},
			{"SAN-CaptureMove-InvalidPromo-White", active(piece.White), "exd3=B", "e2d3b", nil},
			{"SAN-CaptureMove-InvalidPromo-Black", active(piece.Black), "exd6=b", "e7d6b", nil},
			{"SAN-CaptureMove-InvalidCheck-White", active(piece.White), "exd3+", "e2d3", nil},
			{"SAN-CaptureMove-InvalidMate-Black", active(piece.Black), "exd6#", "e7d6", nil},
			{"SAN-CaptureMove-InvalidPromoCheck-White", active(piece.White), "exd3=B+", "e2d3b", nil},
			{"SAN-CaptureMove-InvalidPromoMate-Black", active(piece.Black), "exd6=b#", "e7d6b", nil},
			{"SAN-PromoMove-PromoQueen-White", active(piece.White), "h8=Q", "h7h8q", nil},
			{"SAN-PromoMove-PromoQueen-Black", active(piece.Black), "a1=q", "a2a1q", nil},
			{"SAN-PromoMove-PromoRook-White", active(piece.White), "h8=R", "h7h8r", nil},
			{"SAN-PromoMove-PromoRook-Black", active(piece.Black), "a1=r", "a2a1r", nil},
			{"SAN-PromoMove-PromoBishop-White", active(piece.White), "h8=B", "h7h8b", nil},
			{"SAN-PromoMove-PromoBishop-Black", active(piece.Black), "a1=b", "a2a1b", nil},
			{"SAN-PromoMove-PromoKnight-White", active(piece.White), "h8=n", "h7h8n", nil},
			{"SAN-PromoMove-PromoKnight-Black", active(piece.Black), "a1=N", "a2a1n", nil},
			{"SAN-PromoMove-InvalidCapture-White", active(piece.White), "xh8=q", "h7h8q", nil},
			{"SAN-PromoMove-InvalidCapture-Black", active(piece.Black), "xa1=Q", "a2a1q", nil},
			{"SAN-PromoMove-InvalidCaptureCheck-White", active(piece.White), "xh8=q+", "h7h8q", nil},
			{"SAN-PromoMove-InvalidCaptureMate-Black", active(piece.Black), "xa1=Q#", "a2a1q", nil},
			{"SAN-CapturePromoMove-PromoBishop-White", active(piece.White), "hxg8=B", "h7g8b", nil},
			{"SAN-CapturePromoMove-PromoKnight-Black", active(piece.Black), "axb1=N", "a2b1n", nil},
			{"SAN-CapturePromoMove-PromoQueen-Mate-White", active(piece.White), "hxg8=Q#", "h7g8q", nil},
			{"SAN-CapturePromoMove-PromoQueen-Check-Black", active(piece.Black), "axb1=q+", "a2b1q", nil},
			{"SAN-CapturePromoMove-PromoBishop-InvalidCheck-White", active(piece.White), "hxg8=B+", "h7g8b", nil},
			{"SAN-CapturePromoMove-PromoKnight-InvalidMate-Black", active(piece.Black), "axb1=n#", "a2b1n", nil},
			{"SAN-CapturePromoMove-PromoQueen-InvalidMissingMate-White", active(piece.White), "hxg8=Q", "h7g8q", nil},
			{"SAN-CapturePromoMove-PromoQueen-InvalidMissingCheck-Black", active(piece.Black), "axb1=q", "a2b1q", nil},
			{"SAN-CaptureMove-EnPassantOnC6-White", multi(active(piece.White), enPassant(square.C6)), "bxc6", "b5c6", nil},
			{"SAN-CaptureMove-EnPassantOnC3-Black", multi(active(piece.Black), enPassant(square.C3)), "bxc3", "b4c3", nil},
			{"SAN-CaptureMove-InvalidEnPassantNotOnC6-White", active(piece.White), "bxc6", "", errors.New("could not find source square of 'bxc6'")},
			{"SAN-CaptureMove-InvalidEnPassantNotOnC3-Black", active(piece.Black), "bxc3", "", errors.New("could not find source square of 'bxc3'")},
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
		[]testParseMoveTestCase{
			{"SAN-Move-White", active(piece.White), "Ng5", "h3g5", nil},
			{"SAN-Move-Black", active(piece.Black), "Ng4", "h6g4", nil},
			{"SAN-Move-InvalidMissingType-White", active(piece.White), "g5", "", errors.New("could not find source square of 'g5'")},
			{"SAN-Move-InvalidMissingType-Black", active(piece.Black), "g4", "", errors.New("could not find source square of 'g4'")},
			{"SAN-Move-InvalidRedundantFile-White", active(piece.White), "Nhg5", "h3g5", nil},
			{"SAN-Move-InvalidRedundantFile-Black", active(piece.Black), "Nhg4", "h6g4", nil},
			{"SAN-Move-InvalidRedundantRank-White", active(piece.White), "N3g5", "h3g5", nil},
			{"SAN-Move-InvalidRedundantRank-Black", active(piece.Black), "N6g4", "h6g4", nil},
			{"SAN-Move-InvalidRedundantFileRank-White", active(piece.White), "Nh3g5", "h3g5", nil},
			{"SAN-Move-InvalidRedundantFileRank-Black", active(piece.Black), "Nh6g4", "h6g4", nil},
			{"SAN-Move-InvalidWrongType-White", active(piece.White), "Kg5", "", errors.New("could not find source square of 'Kg5'")},
			{"SAN-Move-InvalidWrongType-Black", active(piece.Black), "Kg4", "", errors.New("could not find source square of 'Kg4'")},
			{"SAN-Move-InvalidWrongRedundantFile-White", active(piece.White), "Nag5", "", errors.New("could not find source square of 'Nag5'")},
			{"SAN-Move-InvalidWrongRedundantFile-Black", active(piece.Black), "Nag4", "", errors.New("could not find source square of 'Nag4'")},
			{"SAN-Move-InvalidWrongRedundantRank-White", active(piece.White), "N1g5", "", errors.New("could not find source square of 'N1g5'")},
			{"SAN-Move-InvalidWrongRedundantRank-Black", active(piece.Black), "N8g4", "", errors.New("could not find source square of 'N8g4'")},
			{"SAN-Move-InvalidWrongRedundantFileRank-White", active(piece.White), "Na1g5", "a1g5", nil},
			{"SAN-Move-InvalidWrongRedundantFileRank-Black", active(piece.Black), "Na8g4", "a8g4", nil},
			{"SAN-Move-InvalidNoKnightOrigin-White", active(piece.White), "Ng2", "", errors.New("could not find source square of 'Ng2'")},
			{"SAN-Move-InvalidNoKnightOrigin-Black", active(piece.Black), "Ng7", "", errors.New("could not find source square of 'Ng7'")},
			{"SAN-Move-InvalidCapture-White", active(piece.White), "Nxg5", "h3g5", nil},
			{"SAN-Move-InvalidCapture-Black", active(piece.Black), "Nxg4", "h6g4", nil},
			{"SAN-Move-InvalidFile-White", active(piece.White), "Ngg5", "", errors.New("could not find source square of 'Ngg5'")},
			{"SAN-Move-InvalidFile-Black", active(piece.Black), "Ngg4", "", errors.New("could not find source square of 'Ngg4'")},
			{"SAN-Move-InvalidRank-White", active(piece.White), "N5g5", "", errors.New("could not find source square of 'N5g5'")},
			{"SAN-Move-InvalidRank-Black", active(piece.Black), "N4g4", "", errors.New("could not find source square of 'N4g4'")},
			{"SAN-Move-InvalidFileRank-White", active(piece.White), "Ng5g5", "g5g5", nil},
			{"SAN-Move-InvalidFileRank-Black", active(piece.Black), "Ng4g4", "g4g4", nil},
			{"SAN-Move-InvalidPromo-White", active(piece.White), "Ng5=Q", "h3g5q", nil},
			{"SAN-Move-InvalidPromo-Black", active(piece.Black), "Ng4=q", "h6g4q", nil},
			{"SAN-Move-InvalidMate-White", active(piece.White), "Ng5#", "h3g5", nil},
			{"SAN-Move-InvalidCheck-Black", active(piece.Black), "Ng4+", "h6g4", nil},
			{"SAN-Move-InvalidPromoMate-White", active(piece.White), "Ng5=Q#", "h3g5q", nil},
			{"SAN-Move-InvalidPromoCheck-Black", active(piece.Black), "Ng4=q+", "h6g4q", nil},
			{"SAN-AmbiguousMove-SpecifyFile-White", active(piece.White), "Nbc7", "b5c7", nil},
			{"SAN-AmbiguousMove-SpecifyFile-Black", active(piece.Black), "Nbc2", "b4c2", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMissingType-White", active(piece.White), "bc7", "", errors.New("could not find source square of 'bc7'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMissingType-Black", active(piece.Black), "bc2", "", errors.New("could not find source square of 'bc2'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMissingFile-White", active(piece.White), "Nc7", "", errors.New("could not find source square of 'Nc7'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMissingFile-Black", active(piece.Black), "Nc2", "", errors.New("could not find source square of 'Nc2'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidRedundantRank-White", active(piece.White), "Nb5c7", "b5c7", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidRedundantRank-Black", active(piece.Black), "Nb4c2", "b4c2", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongType-White", active(piece.White), "Kbc7", "", errors.New("could not find source square of 'Kbc7'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongType-Black", active(piece.Black), "Kbc2", "", errors.New("could not find source square of 'Kbc2'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongFile-White", active(piece.White), "Nac7", "", errors.New("could not find source square of 'Nac7'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongFile-Black", active(piece.Black), "Nac2", "", errors.New("could not find source square of 'Nac2'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongRedundantRank-White", active(piece.White), "Na5c7", "a5c7", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongRedundantRank-Black", active(piece.Black), "Na4c2", "a4c2", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidCapture-White", active(piece.White), "Nbxc7", "b5c7", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidCapture-Black", active(piece.Black), "Nbxc2", "b4c2", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidPromo-White", active(piece.White), "Nbc7=r", "b5c7r", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidPromo-Black", active(piece.Black), "Nbc2=R", "b4c2r", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidCheck-White", active(piece.White), "Nbc7+", "b5c7", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMate-Black", active(piece.Black), "Nbc2#", "b4c2", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidPromoCheck-White", active(piece.White), "Nbc7=r+", "b5c7r", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidPromoMate-Black", active(piece.Black), "Nbc2=R#", "b4c2r", nil},
			{"SAN-AmbiguousMove-SpecifyRank-White", active(piece.White), "N1a3", "b1a3", nil},
			{"SAN-AmbiguousMove-SpecifyRank-Black", active(piece.Black), "N8a6", "b8a6", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMissingType-White", active(piece.White), "1a3", "", errors.New("could not find source square of '1a3'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMissingType-Black", active(piece.Black), "8a6", "", errors.New("could not find source square of '8a6'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMissingRank-White", active(piece.White), "Na3", "", errors.New("could not find source square of 'Na3'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMissingRank-Black", active(piece.Black), "Na6", "", errors.New("could not find source square of 'Na6'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidRedundantFile-White", active(piece.White), "Nb1a3", "b1a3", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidRedundantFile-Black", active(piece.Black), "Nb8a6", "b8a6", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongType-White", active(piece.White), "K1a3", "", errors.New("could not find source square of 'K1a3'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongType-Black", active(piece.Black), "K8a6", "", errors.New("could not find source square of 'K8a6'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongRank-White", active(piece.White), "N2a3", "", errors.New("could not find source square of 'N2a3'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongRank-Black", active(piece.Black), "N7a6", "", errors.New("could not find source square of 'N7a6'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongRedundantFile-White", active(piece.White), "Na1a3", "a1a3", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongRedundantFile-Black", active(piece.Black), "Na8a6", "a8a6", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidCapture-White", active(piece.White), "N1xa3", "b1a3", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidCapture-Black", active(piece.Black), "N8xa6", "b8a6", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidPromo-White", active(piece.White), "N1a3=B", "b1a3b", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidPromo-Black", active(piece.Black), "N8a6=b", "b8a6b", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMate-White", active(piece.White), "N1a3#", "b1a3", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidCheck-Black", active(piece.Black), "N8a6+", "b8a6", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidPromoMate-White", active(piece.White), "N1a3=B#", "b1a3b", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidPromoCheck-Black", active(piece.Black), "N8a6=b+", "b8a6b", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-White", active(piece.White), "Nd1c3", "d1c3", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-Black", active(piece.Black), "Nd4c6", "d4c6", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingType-White", active(piece.White), "d1c3", "d1c3", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingType-Black", active(piece.Black), "d4c6", "d4c6", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingFile-White", active(piece.White), "N1c3", "", errors.New("could not find source square of 'N1c3'")},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingFile-Black", active(piece.Black), "N4c6", "", errors.New("could not find source square of 'N4c6'")},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingRank-White", active(piece.White), "Ndc3", "", errors.New("could not find source square of 'Ndc3'")},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingRank-Black", active(piece.Black), "Ndc6", "", errors.New("could not find source square of 'Ndc6'")},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingFileRank-White", active(piece.White), "Nc3", "", errors.New("could not find source square of 'Nc3'")},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingFileRank-Black", active(piece.Black), "Nc6", "", errors.New("could not find source square of 'Nc6'")},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongType-White", active(piece.White), "Kd1c3", "d1c3", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongType-Black", active(piece.Black), "Kd4c6", "d4c6", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongFile-White", active(piece.White), "Nc1c3", "c1c3", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongFile-Black", active(piece.Black), "Nc4c6", "c4c6", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongRank-White", active(piece.White), "Nd3c3", "d3c3", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongRank-Black", active(piece.Black), "Nd6c6", "d6c6", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongFileRank-White", active(piece.White), "Nc3c3", "c3c3", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongFileRank-Black", active(piece.Black), "Nc6c6", "c6c6", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidCapture-White", active(piece.White), "Nd1xc3", "d1c3", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidCapture-Black", active(piece.Black), "Nd4xc6", "d4c6", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidPromo-White", active(piece.White), "Nd1c3=n", "d1c3n", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidPromo-Black", active(piece.Black), "Nd4c6=N", "d4c6n", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidCheck-White", active(piece.White), "Nd1c3+", "d1c3", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMate-Black", active(piece.Black), "Nd4c6#", "d4c6", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidPromoCheck-White", active(piece.White), "Nd1c3=n+", "d1c3n", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidPromoMate-Black", active(piece.Black), "Nd4c6=N#", "d4c6n", nil},
			{"SAN-CaptureMove-White", active(piece.White), "Nxg1", "h3g1", nil},
			{"SAN-CaptureMove-Black", active(piece.Black), "Nxg8", "h6g8", nil},
			{"SAN-CaptureMove-InvalidMissingType-White", active(piece.White), "xg1", "", errors.New("could not find source square of 'xg1'")},
			{"SAN-CaptureMove-InvalidMissingType-Black", active(piece.Black), "xg8", "", errors.New("could not find source square of 'xg8'")},
			{"SAN-CaptureMove-InvalidMissingCapture-White", active(piece.White), "Ng1", "h3g1", nil},
			{"SAN-CaptureMove-InvalidMissingCapture-Black", active(piece.Black), "Ng8", "h6g8", nil},
			{"SAN-CaptureMove-InvalidRedundantFile-White", active(piece.White), "Nhxg1", "h3g1", nil},
			{"SAN-CaptureMove-InvalidRedundantFile-Black", active(piece.Black), "Nhxg8", "h6g8", nil},
			{"SAN-CaptureMove-InvalidRedundantRank-White", active(piece.White), "N3xg1", "h3g1", nil},
			{"SAN-CaptureMove-InvalidRedundantRank-Black", active(piece.Black), "N6xg8", "h6g8", nil},
			{"SAN-CaptureMove-InvalidRedundantFileRank-White", active(piece.White), "Nh3xg1", "h3g1", nil},
			{"SAN-CaptureMove-InvalidRedundantFileRank-Black", active(piece.Black), "Nh6xg8", "h6g8", nil},
			{"SAN-CaptureMove-InvalidWrongType-White", active(piece.White), "Qxg1", "", errors.New("could not find source square of 'Qxg1'")},
			{"SAN-CaptureMove-InvalidWrongType-Black", active(piece.Black), "Qxg8", "", errors.New("could not find source square of 'Qxg8'")},
			{"SAN-CaptureMove-InvalidWrongRedundantFile-White", active(piece.White), "Ngxg1", "", errors.New("could not find source square of 'Ngxg1'")},
			{"SAN-CaptureMove-InvalidWrongRedundantFile-Black", active(piece.Black), "Ngxg8", "", errors.New("could not find source square of 'Ngxg8'")},
			{"SAN-CaptureMove-InvalidWrongRedundantRank-White", active(piece.White), "N2xg1", "", errors.New("could not find source square of 'N2xg1'")},
			{"SAN-CaptureMove-InvalidWrongRedundantRank-Black", active(piece.Black), "N7xg8", "", errors.New("could not find source square of 'N7xg8'")},
			{"SAN-CaptureMove-InvalidWrongRedundantFileRank-White", active(piece.White), "Ng2xg1", "g2g1", nil},
			{"SAN-CaptureMove-InvalidWrongRedundantFileRank-Black", active(piece.Black), "Ng7xg8", "g7g8", nil},
			{"SAN-CaptureMove-InvalidWrongPromo-White", active(piece.White), "Nxg1=q", "h3g1q", nil},
			{"SAN-CaptureMove-InvalidWrongPromo-Black", active(piece.Black), "Nxg8=Q", "h6g8q", nil},
			{"SAN-CaptureMove-InvalidWrongMate-White", active(piece.White), "Nxg1#", "h3g1", nil},
			{"SAN-CaptureMove-InvalidWrongCheck-Black", active(piece.Black), "Nxg8+", "h6g8", nil},
			{"SAN-CaptureMove-InvalidWrongPromoMate-White", active(piece.White), "Nxg1=q#", "h3g1q", nil},
			{"SAN-CaptureMove-InvalidWrongPromoCheck-Black", active(piece.Black), "Nxg8=Q+", "h6g8q", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-White", active(piece.White), "Ncxa5", "c4a5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-Black", active(piece.Black), "Ncxa4", "c5a4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingType-White", active(piece.White), "cxa5", "", errors.New("could not find source square of 'cxa5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingType-Black", active(piece.Black), "cxa4", "", errors.New("could not find source square of 'cxa4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingFile-White", active(piece.White), "Nxa5", "", errors.New("could not find source square of 'Nxa5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingFile-Black", active(piece.Black), "Nxa4", "", errors.New("could not find source square of 'Nxa4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingCapture-White", active(piece.White), "Nca5", "c4a5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingCapture-Black", active(piece.Black), "Nca4", "c5a4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidRedundantRank-White", active(piece.White), "Nc4xa5", "c4a5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidRedundantRank-Black", active(piece.Black), "Nc5xa4", "c5a4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongType-White", active(piece.White), "Qcxa5", "", errors.New("could not find source square of 'Qcxa5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongType-Black", active(piece.Black), "Qcxa4", "", errors.New("could not find source square of 'Qcxa4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongFile-White", active(piece.White), "Naxa5", "", errors.New("could not find source square of 'Naxa5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongFile-Black", active(piece.Black), "Naxa4", "", errors.New("could not find source square of 'Naxa4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongRedundantRank-White", active(piece.White), "Nc5xa5", "c5a5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongRedundantRank-Black", active(piece.Black), "Nc4xa4", "c4a4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidPromo-White", active(piece.White), "Ncxa5=R", "c4a5r", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidPromo-Black", active(piece.Black), "Ncxa4=r", "c5a4r", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidCheck-White", active(piece.White), "Ncxa5+", "c4a5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMate-Black", active(piece.Black), "Ncxa4#", "c5a4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidPromoCheck-White", active(piece.White), "Ncxa5=R+", "c4a5r", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidPromoMate-Black", active(piece.Black), "Ncxa4=r#", "c5a4r", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-White", active(piece.White), "N3xa5", "b3a5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-Black", active(piece.Black), "N2xa4", "b2a4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingType-White", active(piece.White), "3xa5", "", errors.New("could not find source square of '3xa5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingType-Black", active(piece.Black), "2xa4", "", errors.New("could not find source square of '2xa4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingRank-White", active(piece.White), "Nxa5", "", errors.New("could not find source square of 'Nxa5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingRank-Black", active(piece.Black), "Nxa4", "", errors.New("could not find source square of 'Nxa4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingCapture-White", active(piece.White), "N3a5", "b3a5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingCapture-Black", active(piece.Black), "N2a4", "b2a4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidRedundantFile-White", active(piece.White), "Nb3xa5", "b3a5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidRedundantFile-Black", active(piece.Black), "Nb2xa4", "b2a4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongType-White", active(piece.White), "Q3xa5", "", errors.New("could not find source square of 'Q3xa5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongType-Black", active(piece.Black), "Q2xa4", "", errors.New("could not find source square of 'Q2xa4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongRank-White", active(piece.White), "N5xa5", "", errors.New("could not find source square of 'N5xa5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongRank-Black", active(piece.Black), "N4xa4", "", errors.New("could not find source square of 'N4xa4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongRedundantFile-White", active(piece.White), "Na3xa5", "a3a5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongRedundantFile-Black", active(piece.Black), "Na2xa4", "a2a4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidPromo-White", active(piece.White), "N3xa5=b", "b3a5b", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidPromo-Black", active(piece.Black), "N2xa4=B", "b2a4b", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMate-White", active(piece.White), "N3xa5#", "b3a5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidCheck-Black", active(piece.Black), "N2xa4+", "b2a4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidPromoMate-White", active(piece.White), "N3xa5=b#", "b3a5b", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidPromoCheck-Black", active(piece.Black), "N2xa4=B+", "b2a4b", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-White", active(piece.White), "Nb3xc5", "b3c5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-Black", active(piece.Black), "Nb2xc4", "b2c4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingType-White", active(piece.White), "b3xc5", "b3c5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingType-Black", active(piece.Black), "b2xc4", "b2c4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingFile-White", active(piece.White), "N3xc5", "", errors.New("could not find source square of 'N3xc5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingFile-Black", active(piece.Black), "N2xc4", "", errors.New("could not find source square of 'N2xc4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingRank-White", active(piece.White), "Nbxc5", "", errors.New("could not find source square of 'Nbxc5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingRank-Black", active(piece.Black), "Nbxc4", "", errors.New("could not find source square of 'Nbxc4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingFileRank-White", active(piece.White), "Nxc5", "", errors.New("could not find source square of 'Nxc5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingFileRank-Black", active(piece.Black), "Nxc4", "", errors.New("could not find source square of 'Nxc4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingCapture-White", active(piece.White), "Nb3c5", "b3c5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingCapture-Black", active(piece.Black), "Nb2c4", "b2c4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongType-White", active(piece.White), "Qb3xc5", "b3c5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongType-Black", active(piece.Black), "Qb2xc4", "b2c4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongFile-White", active(piece.White), "Nc3xc5", "c3c5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongFile-Black", active(piece.Black), "Nc2xc4", "c2c4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongRank-White", active(piece.White), "Nb5xc5", "b5c5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongRank-Black", active(piece.Black), "Nb4xc4", "b4c4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongFileRank-White", active(piece.White), "Nc5xc5", "c5c5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongFileRank-Black", active(piece.Black), "Nc4xc4", "c4c4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidPromo-White", active(piece.White), "Nb3xc5=N", "b3c5n", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidPromo-Black", active(piece.Black), "Nb2xc4=n", "b2c4n", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidCheck-White", active(piece.White), "Nb3xc5+", "b3c5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMate-Black", active(piece.Black), "Nb2xc4#", "b2c4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidPromoCheck-White", active(piece.White), "Nb3xc5=N+", "b3c5n", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidPromoMate-Black", active(piece.Black), "Nb2xc4=n#", "b2c4n", nil},
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
		[]testParseMoveTestCase{
			{"SAN-Move-White", active(piece.White), "Be3", "c1e3", nil},
			{"SAN-Move-Black", active(piece.Black), "Be6", "c8e6", nil},
			{"SAN-Move-InvalidMissingType-White", active(piece.White), "e3", "", errors.New("could not find source square of 'e3'")},
			{"SAN-Move-InvalidMissingType-Black", active(piece.Black), "e6", "", errors.New("could not find source square of 'e6'")},
			{"SAN-Move-InvalidRedundantFile-White", active(piece.White), "Bce3", "c1e3", nil},
			{"SAN-Move-InvalidRedundantFile-Black", active(piece.Black), "Bce6", "c8e6", nil},
			{"SAN-Move-InvalidRedundantRank-White", active(piece.White), "B1e3", "c1e3", nil},
			{"SAN-Move-InvalidRedundantRank-Black", active(piece.Black), "B8e6", "c8e6", nil},
			{"SAN-Move-InvalidRedundantFileRank-White", active(piece.White), "Bc1e3", "c1e3", nil},
			{"SAN-Move-InvalidRedundantFileRank-Black", active(piece.Black), "Bc8e6", "c8e6", nil},
			{"SAN-Move-InvalidWrongType-White", active(piece.White), "Ke3", "", errors.New("could not find source square of 'Ke3'")},
			{"SAN-Move-InvalidWrongType-Black", active(piece.Black), "Ke6", "", errors.New("could not find source square of 'Ke6'")},
			{"SAN-Move-InvalidWrongRedundantFile-White", active(piece.White), "Bee3", "", errors.New("could not find source square of 'Bee3'")},
			{"SAN-Move-InvalidWrongRedundantFile-Black", active(piece.Black), "Bee6", "", errors.New("could not find source square of 'Bee6'")},
			{"SAN-Move-InvalidWrongRedundantRank-White", active(piece.White), "B3e3", "", errors.New("could not find source square of 'B3e3'")},
			{"SAN-Move-InvalidWrongRedundantRank-Black", active(piece.Black), "B6e6", "", errors.New("could not find source square of 'B6e6'")},
			{"SAN-Move-InvalidWrongRedundantFileRank-White", active(piece.White), "Be3e3", "e3e3", nil},
			{"SAN-Move-InvalidWrongRedundantFileRank-Black", active(piece.Black), "Be6e6", "e6e6", nil},
			{"SAN-Move-InvalidNoBishopOrigin-White", active(piece.White), "Bf2", "", errors.New("could not find source square of 'Bf2'")},
			{"SAN-Move-InvalidNoBishopOrigin-Black", active(piece.Black), "Bf7", "", errors.New("could not find source square of 'Bf7'")},
			{"SAN-Move-InvalidCapture-White", active(piece.White), "Bxe3", "c1e3", nil},
			{"SAN-Move-InvalidCapture-Black", active(piece.Black), "Bxe6", "c8e6", nil},
			{"SAN-Move-InvalidFile-White", active(piece.White), "Bee3", "", errors.New("could not find source square of 'Bee3'")},
			{"SAN-Move-InvalidFile-Black", active(piece.Black), "Bee6", "", errors.New("could not find source square of 'Bee6'")},
			{"SAN-Move-InvalidRank-White", active(piece.White), "B3e3", "", errors.New("could not find source square of 'B3e3'")},
			{"SAN-Move-InvalidRank-Black", active(piece.Black), "B6e6", "", errors.New("could not find source square of 'B6e6'")},
			{"SAN-Move-InvalidFileRank-White", active(piece.White), "Be3e3", "e3e3", nil},
			{"SAN-Move-InvalidFileRank-Black", active(piece.Black), "Be6e6", "e6e6", nil},
			{"SAN-Move-InvalidPromo-White", active(piece.White), "Be3=Q", "c1e3q", nil},
			{"SAN-Move-InvalidPromo-Black", active(piece.Black), "Be6=q", "c8e6q", nil},
			{"SAN-Move-InvalidCheck-White", active(piece.White), "Be3+", "c1e3", nil},
			{"SAN-Move-InvalidMate-Black", active(piece.Black), "Be6#", "c8e6", nil},
			{"SAN-Move-InvalidPromoCheck-White", active(piece.White), "Be3=Q+", "c1e3q", nil},
			{"SAN-Move-InvalidPromoMate-Black", active(piece.Black), "Be6=q#", "c8e6q", nil},
			{"SAN-AmbiguousMove-SpecifyFile-White", active(piece.White), "Bfg4", "f3g4", nil},
			{"SAN-AmbiguousMove-SpecifyFile-Black", active(piece.Black), "Bfg5", "f6g5", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMissingType-White", active(piece.White), "fg4", "", errors.New("could not find source square of 'fg4'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMissingType-Black", active(piece.Black), "fg5", "", errors.New("could not find source square of 'fg5'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMissingFile-White", active(piece.White), "Bg4", "", errors.New("could not find source square of 'Bg4'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMissingFile-Black", active(piece.Black), "Bg5", "", errors.New("could not find source square of 'Bg5'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidRedundantRank-White", active(piece.White), "Bf3g4", "f3g4", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidRedundantRank-Black", active(piece.Black), "Bf6g5", "f6g5", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongType-White", active(piece.White), "Kfg4", "", errors.New("could not find source square of 'Kfg4'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongType-Black", active(piece.Black), "Kfg5", "", errors.New("could not find source square of 'Kfg5'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongFile-White", active(piece.White), "Bgg4", "", errors.New("could not find source square of 'Bgg4'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongFile-Black", active(piece.Black), "Bgg5", "", errors.New("could not find source square of 'Bgg5'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongRedundantRank-White", active(piece.White), "Bf4g4", "f4g4", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongRedundantRank-Black", active(piece.Black), "Bf5g5", "f5g5", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidCapture-White", active(piece.White), "Bfxg4", "f3g4", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidCapture-Black", active(piece.Black), "Bfxg5", "f6g5", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidPromo-White", active(piece.White), "Bfg4=r", "f3g4r", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidPromo-Black", active(piece.Black), "Bfg5=R", "f6g5r", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMate-White", active(piece.White), "Bfg4#", "f3g4", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidCheck-Black", active(piece.Black), "Bfg5+", "f6g5", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidPromoMate-White", active(piece.White), "Bfg4=r#", "f3g4r", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidPromoCheck-Black", active(piece.Black), "Bfg5=R+", "f6g5r", nil},
			{"SAN-AmbiguousMove-SpecifyRank-White", active(piece.White), "B1d2", "c1d2", nil},
			{"SAN-AmbiguousMove-SpecifyRank-Black", active(piece.Black), "B8d7", "c8d7", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMissingType-White", active(piece.White), "1d2", "", errors.New("could not find source square of '1d2'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMissingType-Black", active(piece.Black), "8d7", "", errors.New("could not find source square of '8d7'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMissingRank-White", active(piece.White), "Bd2", "", errors.New("could not find source square of 'Bd2'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMissingRank-Black", active(piece.Black), "Bd7", "", errors.New("could not find source square of 'Bd7'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidRedundantFile-White", active(piece.White), "Bc1d2", "c1d2", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidRedundantFile-Black", active(piece.Black), "Bc8d7", "c8d7", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongType-White", active(piece.White), "K1d2", "", errors.New("could not find source square of 'K1d2'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongType-Black", active(piece.Black), "K8d7", "", errors.New("could not find source square of 'K8d7'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongRank-White", active(piece.White), "B2d2", "", errors.New("could not find source square of 'B2d2'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongRank-Black", active(piece.Black), "B7d7", "", errors.New("could not find source square of 'B7d7'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongRedundantFile-White", active(piece.White), "Bd1d2", "d1d2", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongRedundantFile-Black", active(piece.Black), "Bd8d7", "d8d7", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidCapture-White", active(piece.White), "B1xd2", "c1d2", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidCapture-Black", active(piece.Black), "B8xd7", "c8d7", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidPromo-White", active(piece.White), "B1d2=B", "c1d2b", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidPromo-Black", active(piece.Black), "B8d7=b", "c8d7b", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidCheck-White", active(piece.White), "B1d2+", "c1d2", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMate-Black", active(piece.Black), "B8d7#", "c8d7", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidPromoCheck-White", active(piece.White), "B1d2=B+", "c1d2b", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidPromoMate-Black", active(piece.Black), "B8d7=b#", "c8d7b", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-White", active(piece.White), "Bh1g2", "h1g2", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-Black", active(piece.Black), "Bh8g7", "h8g7", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingType-White", active(piece.White), "h1g2", "h1g2", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingType-Black", active(piece.Black), "h8g7", "h8g7", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingFile-White", active(piece.White), "B1g2", "", errors.New("could not find source square of 'B1g2'")},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingFile-Black", active(piece.Black), "B8g7", "", errors.New("could not find source square of 'B8g7'")},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingRank-White", active(piece.White), "Bhg2", "", errors.New("could not find source square of 'Bhg2'")},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingRank-Black", active(piece.Black), "Bhg7", "", errors.New("could not find source square of 'Bhg7'")},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingFileRank-White", active(piece.White), "Bg2", "", errors.New("could not find source square of 'Bg2'")},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingFileRank-Black", active(piece.Black), "Bg7", "", errors.New("could not find source square of 'Bg7'")},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongType-White", active(piece.White), "Kh1g2", "h1g2", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongType-Black", active(piece.Black), "Kh8g7", "h8g7", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongFile-White", active(piece.White), "Bg1g2", "g1g2", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongFile-Black", active(piece.Black), "Bg8g7", "g8g7", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongRank-White", active(piece.White), "Bh2g2", "h2g2", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongRank-Black", active(piece.Black), "Bh7g7", "h7g7", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongFileRank-White", active(piece.White), "Bg2g2", "g2g2", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongFileRank-Black", active(piece.Black), "Bg7g7", "g7g7", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidCapture-White", active(piece.White), "Bh1xg2", "h1g2", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidCapture-Black", active(piece.Black), "Bh8xg7", "h8g7", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidPromo-White", active(piece.White), "Bh1g2=n", "h1g2n", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidPromo-Black", active(piece.Black), "Bh8g7=N", "h8g7n", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMate-White", active(piece.White), "Bh1g2#", "h1g2", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidCheck-Black", active(piece.Black), "Bh8g7+", "h8g7", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidPromoMate-White", active(piece.White), "Bh1g2=n#", "h1g2n", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidPromoCheck-Black", active(piece.Black), "Bh8g7=N+", "h8g7n", nil},
			{"SAN-CaptureMove-White", active(piece.White), "Bxf6", "e7f6", nil},
			{"SAN-CaptureMove-Black", active(piece.Black), "Bxf3", "e2f3", nil},
			{"SAN-CaptureMove-InvalidMissingType-White", active(piece.White), "xf6", "", errors.New("could not find source square of 'xf6'")},
			{"SAN-CaptureMove-InvalidMissingType-Black", active(piece.Black), "xf3", "", errors.New("could not find source square of 'xf3'")},
			{"SAN-CaptureMove-InvalidMissingCapture-White", active(piece.White), "Bf6", "e7f6", nil},
			{"SAN-CaptureMove-InvalidMissingCapture-Black", active(piece.Black), "Bf3", "e2f3", nil},
			{"SAN-CaptureMove-InvalidRedundantFile-White", active(piece.White), "Bexf6", "e7f6", nil},
			{"SAN-CaptureMove-InvalidRedundantFile-Black", active(piece.Black), "Bexf3", "e2f3", nil},
			{"SAN-CaptureMove-InvalidRedundantRank-White", active(piece.White), "B7xf6", "e7f6", nil},
			{"SAN-CaptureMove-InvalidRedundantRank-Black", active(piece.Black), "B2xf3", "e2f3", nil},
			{"SAN-CaptureMove-InvalidRedundantFileRank-White", active(piece.White), "Be7xf6", "e7f6", nil},
			{"SAN-CaptureMove-InvalidRedundantFileRank-Black", active(piece.Black), "Be2xf3", "e2f3", nil},
			{"SAN-CaptureMove-InvalidWrongType-White", active(piece.White), "Kxf6", "", errors.New("could not find source square of 'Kxf6'")},
			{"SAN-CaptureMove-InvalidWrongType-Black", active(piece.Black), "Kxf3", "", errors.New("could not find source square of 'Kxf3'")},
			{"SAN-CaptureMove-InvalidWrongRedundantFile-White", active(piece.White), "Bfxf6", "", errors.New("could not find source square of 'Bfxf6'")},
			{"SAN-CaptureMove-InvalidWrongRedundantFile-Black", active(piece.Black), "Bfxf3", "", errors.New("could not find source square of 'Bfxf3'")},
			{"SAN-CaptureMove-InvalidWrongRedundantRank-White", active(piece.White), "B6xf6", "", errors.New("could not find source square of 'B6xf6'")},
			{"SAN-CaptureMove-InvalidWrongRedundantRank-Black", active(piece.Black), "B3xf3", "", errors.New("could not find source square of 'B3xf3'")},
			{"SAN-CaptureMove-InvalidWrongRedundantFileRank-White", active(piece.White), "Bf6xf6", "f6f6", nil},
			{"SAN-CaptureMove-InvalidWrongRedundantFileRank-Black", active(piece.Black), "Bf3xf3", "f3f3", nil},
			{"SAN-CaptureMove-InvalidWrongPromo-White", active(piece.White), "Bxf6=q", "e7f6q", nil},
			{"SAN-CaptureMove-InvalidWrongPromo-Black", active(piece.Black), "Bxf3=Q", "e2f3q", nil},
			{"SAN-CaptureMove-InvalidWrongCheck-White", active(piece.White), "Bxf6+", "e7f6", nil},
			{"SAN-CaptureMove-InvalidWrongMate-Black", active(piece.Black), "Bxf3#", "e2f3", nil},
			{"SAN-CaptureMove-InvalidWrongPromoCheck-White", active(piece.White), "Bxf6=q+", "e7f6q", nil},
			{"SAN-CaptureMove-InvalidWrongPromoMate-Black", active(piece.Black), "Bxf3=Q#", "e2f3q", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-White", active(piece.White), "Baxb4", "a3b4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-Black", active(piece.Black), "Baxb5", "a6b5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingType-White", active(piece.White), "axb4", "", errors.New("could not find source square of 'axb4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingType-Black", active(piece.Black), "axb5", "", errors.New("could not find source square of 'axb5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingFile-White", active(piece.White), "Bxb4", "", errors.New("could not find source square of 'Bxb4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingFile-Black", active(piece.Black), "Bxb5", "", errors.New("could not find source square of 'Bxb5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingCapture-White", active(piece.White), "Bab4", "a3b4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingCapture-Black", active(piece.Black), "Bab5", "a6b5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidRedundantRank-White", active(piece.White), "Ba3xb4", "a3b4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidRedundantRank-Black", active(piece.Black), "Ba6xb5", "a6b5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongType-White", active(piece.White), "Kaxb4", "", errors.New("could not find source square of 'Kaxb4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongType-Black", active(piece.Black), "Kaxb5", "", errors.New("could not find source square of 'Kaxb5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongFile-White", active(piece.White), "Bbxb4", "", errors.New("could not find source square of 'Bbxb4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongFile-Black", active(piece.Black), "Bbxb5", "", errors.New("could not find source square of 'Bbxb5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongRedundantRank-White", active(piece.White), "Ba4xb4", "a4b4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongRedundantRank-Black", active(piece.Black), "Ba5xb5", "a5b5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidPromo-White", active(piece.White), "Baxb4=R", "a3b4r", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidPromo-Black", active(piece.Black), "Baxb5=r", "a6b5r", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMate-White", active(piece.White), "Baxb4#", "a3b4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidCheck-Black", active(piece.Black), "Baxb5+", "a6b5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidPromoMate-White", active(piece.White), "Baxb4=R#", "a3b4r", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidPromoCheck-Black", active(piece.Black), "Baxb5=r+", "a6b5r", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-White", active(piece.White), "B1xe2", "f1e2", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-Black", active(piece.Black), "B8xe7", "f8e7", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingType-White", active(piece.White), "1xe2", "", errors.New("could not find source square of '1xe2'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingType-Black", active(piece.Black), "8xe7", "", errors.New("could not find source square of '8xe7'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingRank-White", active(piece.White), "Bxe2", "", errors.New("could not find source square of 'Bxe2'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingRank-Black", active(piece.Black), "Bxe7", "", errors.New("could not find source square of 'Bxe7'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingCapture-White", active(piece.White), "B1e2", "f1e2", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingCapture-Black", active(piece.Black), "B8e7", "f8e7", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidRedundantFile-White", active(piece.White), "Bf1xe2", "f1e2", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidRedundantFile-Black", active(piece.Black), "Bf8xe7", "f8e7", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongType-White", active(piece.White), "Q1xe2", "", errors.New("could not find source square of 'Q1xe2'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongType-Black", active(piece.Black), "Q8xe7", "", errors.New("could not find source square of 'Q8xe7'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongRank-White", active(piece.White), "B2xe2", "", errors.New("could not find source square of 'B2xe2'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongRank-Black", active(piece.Black), "B7xe7", "", errors.New("could not find source square of 'B7xe7'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongRedundantFile-White", active(piece.White), "Be1xe2", "e1e2", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongRedundantFile-Black", active(piece.Black), "Be8xe7", "e8e7", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidPromo-White", active(piece.White), "B1xe2=b", "f1e2b", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidPromo-Black", active(piece.Black), "B8xe7=B", "f8e7b", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidCheck-White", active(piece.White), "B1xe2+", "f1e2", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMate-Black", active(piece.Black), "B8xe7#", "f8e7", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidPromoCheck-White", active(piece.White), "B1xe2=b+", "f1e2b", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidPromoMate-Black", active(piece.Black), "B8xe7=B#", "f8e7b", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-White", active(piece.White), "Ba1xb2", "a1b2", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-Black", active(piece.Black), "Ba8xb7", "a8b7", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingType-White", active(piece.White), "a1xb2", "a1b2", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingType-Black", active(piece.Black), "a8xb7", "a8b7", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingFile-White", active(piece.White), "B1xb2", "", errors.New("could not find source square of 'B1xb2'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingFile-Black", active(piece.Black), "B8xb7", "", errors.New("could not find source square of 'B8xb7'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingRank-White", active(piece.White), "Baxb2", "", errors.New("could not find source square of 'Baxb2'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingRank-Black", active(piece.Black), "Baxb7", "", errors.New("could not find source square of 'Baxb7'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingFileRank-White", active(piece.White), "Bxb2", "", errors.New("could not find source square of 'Bxb2'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingFileRank-Black", active(piece.Black), "Bxb7", "", errors.New("could not find source square of 'Bxb7'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingCapture-White", active(piece.White), "Ba1b2", "a1b2", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingCapture-Black", active(piece.Black), "Ba8b7", "a8b7", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongType-White", active(piece.White), "Ka1xb2", "a1b2", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongType-Black", active(piece.Black), "Ka8xb7", "a8b7", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongFile-White", active(piece.White), "Bb1xb2", "b1b2", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongFile-Black", active(piece.Black), "Bb8xb7", "b8b7", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongRank-White", active(piece.White), "Ba2xb2", "a2b2", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongRank-Black", active(piece.Black), "Ba7xb7", "a7b7", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongFileRank-White", active(piece.White), "Bb2xb2", "b2b2", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongFileRank-Black", active(piece.Black), "Bb7xb7", "b7b7", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidPromo-White", active(piece.White), "Ba1xb2=N", "a1b2n", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidPromo-Black", active(piece.Black), "Ba8xb7=n", "a8b7n", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMate-White", active(piece.White), "Ba1xb2#", "a1b2", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidCheck-Black", active(piece.Black), "Ba8xb7+", "a8b7", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidPromoMate-White", active(piece.White), "Ba1xb2=N#", "a1b2n", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidPromoCheck-Black", active(piece.Black), "Ba8xb7=n+", "a8b7n", nil},
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
		[]testParseMoveTestCase{
			{"SAN-Move-White", active(piece.White), "Rg1", "g2g1", nil},
			{"SAN-Move-Black", active(piece.Black), "Rg8", "g7g8", nil},
			{"SAN-Move-InvalidMissingType-White", active(piece.White), "g1", "", errors.New("could not find source square of 'g1'")},
			{"SAN-Move-InvalidMissingType-Black", active(piece.Black), "g8", "", errors.New("could not find source square of 'g8'")},
			{"SAN-Move-InvalidRedundantFile-White", active(piece.White), "Rgg1", "g2g1", nil},
			{"SAN-Move-InvalidRedundantFile-Black", active(piece.Black), "Rgg8", "g7g8", nil},
			{"SAN-Move-InvalidRedundantRank-White", active(piece.White), "R2g1", "g2g1", nil},
			{"SAN-Move-InvalidRedundantRank-Black", active(piece.Black), "R7g8", "g7g8", nil},
			{"SAN-Move-InvalidRedundantFileRank-White", active(piece.White), "Rg2g1", "g2g1", nil},
			{"SAN-Move-InvalidRedundantFileRank-Black", active(piece.Black), "Rg7g8", "g7g8", nil},
			{"SAN-Move-InvalidWrongType-White", active(piece.White), "Kg1", "", errors.New("could not find source square of 'Kg1'")},
			{"SAN-Move-InvalidWrongType-Black", active(piece.Black), "Kg8", "", errors.New("could not find source square of 'Kg8'")},
			{"SAN-Move-InvalidWrongRedundantFile-White", active(piece.White), "Rhg1", "", errors.New("could not find source square of 'Rhg1'")},
			{"SAN-Move-InvalidWrongRedundantFile-Black", active(piece.Black), "Rhg8", "", errors.New("could not find source square of 'Rhg8'")},
			{"SAN-Move-InvalidWrongRedundantRank-White", active(piece.White), "R8g1", "", errors.New("could not find source square of 'R8g1'")},
			{"SAN-Move-InvalidWrongRedundantRank-Black", active(piece.Black), "R1g8", "", errors.New("could not find source square of 'R1g8'")},
			{"SAN-Move-InvalidWrongRedundantFileRank-White", active(piece.White), "Rh8g1", "h8g1", nil},
			{"SAN-Move-InvalidWrongRedundantFileRank-Black", active(piece.Black), "Rh1g8", "h1g8", nil},
			{"SAN-Move-InvalidNoRookOrigin-White", active(piece.White), "Re1", "", errors.New("could not find source square of 'Re1'")},
			{"SAN-Move-InvalidNoRookOrigin-Black", active(piece.Black), "Re8", "", errors.New("could not find source square of 'Re8'")},
			{"SAN-Move-InvalidCapture-White", active(piece.White), "Rxg1", "g2g1", nil},
			{"SAN-Move-InvalidCapture-Black", active(piece.Black), "Rxg8", "g7g8", nil},
			{"SAN-Move-InvalidFile-White", active(piece.White), "Rhg1", "", errors.New("could not find source square of 'Rhg1'")},
			{"SAN-Move-InvalidFile-Black", active(piece.Black), "Rhg8", "", errors.New("could not find source square of 'Rhg8'")},
			{"SAN-Move-InvalidRank-White", active(piece.White), "R1g1", "", errors.New("could not find source square of 'R1g1'")},
			{"SAN-Move-InvalidRank-Black", active(piece.Black), "R8g8", "", errors.New("could not find source square of 'R8g8'")},
			{"SAN-Move-InvalidFileRank-White", active(piece.White), "Rh1g1", "h1g1", nil},
			{"SAN-Move-InvalidFileRank-Black", active(piece.Black), "Rh8g8", "h8g8", nil},
			{"SAN-Move-InvalidPromo-White", active(piece.White), "Rg1=Q", "g2g1q", nil},
			{"SAN-Move-InvalidPromo-Black", active(piece.Black), "Rg8=q", "g7g8q", nil},
			{"SAN-Move-InvalidMate-White", active(piece.White), "Rg1#", "g2g1", nil},
			{"SAN-Move-InvalidCheck-Black", active(piece.Black), "Rg8+", "g7g8", nil},
			{"SAN-Move-InvalidPromoMate-White", active(piece.White), "Rg1=Q#", "g2g1q", nil},
			{"SAN-Move-InvalidPromoCheck-Black", active(piece.Black), "Rg8=q+", "g7g8q", nil},
			{"SAN-AmbiguousMove-SpecifyFile-White", active(piece.White), "Rce2", "c2e2", nil},
			{"SAN-AmbiguousMove-SpecifyFile-Black", active(piece.Black), "Rce7", "c7e7", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMissingType-White", active(piece.White), "ce2", "", errors.New("could not find source square of 'ce2'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMissingType-Black", active(piece.Black), "ce7", "", errors.New("could not find source square of 'ce7'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMissingFile-White", active(piece.White), "Re2", "", errors.New("could not find source square of 'Re2'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMissingFile-Black", active(piece.Black), "Re7", "", errors.New("could not find source square of 'Re7'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidRedundantRank-White", active(piece.White), "Rc2e2", "c2e2", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidRedundantRank-Black", active(piece.Black), "Rc7e7", "c7e7", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongType-White", active(piece.White), "Qce2", "", errors.New("could not find source square of 'Qce2'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongType-Black", active(piece.Black), "Qce7", "", errors.New("could not find source square of 'Qce7'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongFile-White", active(piece.White), "Rae2", "", errors.New("could not find source square of 'Rae2'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongFile-Black", active(piece.Black), "Rae7", "", errors.New("could not find source square of 'Rae7'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongRedundantRank-White", active(piece.White), "Rc8e2", "c8e2", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongRedundantRank-Black", active(piece.Black), "Rc1e7", "c1e7", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidCapture-White", active(piece.White), "Rcxe2", "c2e2", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidCapture-Black", active(piece.Black), "Rcxe7", "c7e7", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidPromo-White", active(piece.White), "Rce2=r", "c2e2r", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidPromo-Black", active(piece.Black), "Rce7=R", "c7e7r", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidCheck-White", active(piece.White), "Rce2+", "c2e2", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMate-Black", active(piece.Black), "Rce7#", "c7e7", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidPromoCheck-White", active(piece.White), "Rce2=r+", "c2e2r", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidPromoMate-Black", active(piece.Black), "Rce7=R#", "c7e7r", nil},
			{"SAN-AmbiguousMove-SpecifyRank-White", active(piece.White), "R2g3", "g2g3", nil},
			{"SAN-AmbiguousMove-SpecifyRank-Black", active(piece.Black), "R7g6", "g7g6", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMissingType-White", active(piece.White), "2g3", "", errors.New("could not find source square of '2g3'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMissingType-Black", active(piece.Black), "7g6", "", errors.New("could not find source square of '7g6'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMissingRank-White", active(piece.White), "Rg3", "", errors.New("could not find source square of 'Rg3'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMissingRank-Black", active(piece.Black), "Rg6", "", errors.New("could not find source square of 'Rg6'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidRedundantFile-White", active(piece.White), "Rg2g3", "g2g3", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidRedundantFile-Black", active(piece.Black), "Rg7g6", "g7g6", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongType-White", active(piece.White), "K2g3", "", errors.New("could not find source square of 'K2g3'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongType-Black", active(piece.Black), "K7g6", "", errors.New("could not find source square of 'K7g6'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongRank-White", active(piece.White), "R1g3", "", errors.New("could not find source square of 'R1g3'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongRank-Black", active(piece.Black), "R8g6", "", errors.New("could not find source square of 'R8g6'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongRedundantFile-White", active(piece.White), "Rh2g3", "h2g3", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongRedundantFile-Black", active(piece.Black), "Rh7g6", "h7g6", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidCapture-White", active(piece.White), "R2xg3", "g2g3", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidCapture-Black", active(piece.Black), "R7xg6", "g7g6", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidPromo-White", active(piece.White), "R2g3=B", "g2g3b", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidPromo-Black", active(piece.Black), "R7g6=b", "g7g6b", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMate-White", active(piece.White), "R2g3#", "g2g3", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidCheck-Black", active(piece.Black), "R7g6+", "g7g6", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidPromoMate-White", active(piece.White), "R2g3=B#", "g2g3b", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidPromoCheck-Black", active(piece.Black), "R7g6=b+", "g7g6b", nil},
			{"SAN-CaptureMove-White", active(piece.White), "Rxh6", "h3h6", nil},
			{"SAN-CaptureMove-Black", active(piece.Black), "Rxh3", "h6h3", nil},
			{"SAN-CaptureMove-InvalidMissingType-White", active(piece.White), "xh6", "", errors.New("could not find source square of 'xh6'")},
			{"SAN-CaptureMove-InvalidMissingType-Black", active(piece.Black), "xh3", "", errors.New("could not find source square of 'xh3'")},
			{"SAN-CaptureMove-InvalidMissingCapture-White", active(piece.White), "Rh6", "h3h6", nil},
			{"SAN-CaptureMove-InvalidMissingCapture-Black", active(piece.Black), "Rh3", "h6h3", nil},
			{"SAN-CaptureMove-InvalidRedundantFile-White", active(piece.White), "Rhxh6", "h3h6", nil},
			{"SAN-CaptureMove-InvalidRedundantFile-Black", active(piece.Black), "Rhxh3", "h6h3", nil},
			{"SAN-CaptureMove-InvalidRedundantRank-White", active(piece.White), "R3xh6", "h3h6", nil},
			{"SAN-CaptureMove-InvalidRedundantRank-Black", active(piece.Black), "R6xh3", "h6h3", nil},
			{"SAN-CaptureMove-InvalidRedundantFileRank-White", active(piece.White), "Rh3xh6", "h3h6", nil},
			{"SAN-CaptureMove-InvalidRedundantFileRank-Black", active(piece.Black), "Rh6xh3", "h6h3", nil},
			{"SAN-CaptureMove-InvalidWrongType-White", active(piece.White), "Kxh6", "", errors.New("could not find source square of 'Kxh6'")},
			{"SAN-CaptureMove-InvalidWrongType-Black", active(piece.Black), "Kxh3", "", errors.New("could not find source square of 'Kxh3'")},
			{"SAN-CaptureMove-InvalidWrongRedundantFile-White", active(piece.White), "Rgxh6", "", errors.New("could not find source square of 'Rgxh6'")},
			{"SAN-CaptureMove-InvalidWrongRedundantFile-Black", active(piece.Black), "Rgxh3", "", errors.New("could not find source square of 'Rgxh3'")},
			{"SAN-CaptureMove-InvalidWrongRedundantRank-White", active(piece.White), "R1xh6", "", errors.New("could not find source square of 'R1xh6'")},
			{"SAN-CaptureMove-InvalidWrongRedundantRank-Black", active(piece.Black), "R8xh3", "", errors.New("could not find source square of 'R8xh3'")},
			{"SAN-CaptureMove-InvalidWrongRedundantFileRank-White", active(piece.White), "Rg1xh6", "g1h6", nil},
			{"SAN-CaptureMove-InvalidWrongRedundantFileRank-Black", active(piece.Black), "Rg8xh3", "g8h3", nil},
			{"SAN-CaptureMove-InvalidPromo-White", active(piece.White), "Rxh6=q", "h3h6q", nil},
			{"SAN-CaptureMove-InvalidPromo-Black", active(piece.Black), "Rxh3=Q", "h6h3q", nil},
			{"SAN-CaptureMove-InvalidMate-White", active(piece.White), "Rxh6#", "h3h6", nil},
			{"SAN-CaptureMove-InvalidCheck-Black", active(piece.Black), "Rxh3+", "h6h3", nil},
			{"SAN-CaptureMove-InvalidPromoMate-White", active(piece.White), "Rxh6=q#", "h3h6q", nil},
			{"SAN-CaptureMove-InvalidPromoCheck-Black", active(piece.Black), "Rxh3=Q+", "h6h3q", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-White", active(piece.White), "Rbxa3", "b3a3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-Black", active(piece.Black), "Rbxa6", "b6a6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingType-White", active(piece.White), "bxa3", "", errors.New("could not find source square of 'bxa3'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingType-Black", active(piece.Black), "bxa6", "", errors.New("could not find source square of 'bxa6'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingFile-White", active(piece.White), "Rxa3", "", errors.New("could not find source square of 'Rxa3'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingFile-Black", active(piece.Black), "Rxa6", "", errors.New("could not find source square of 'Rxa6'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingCapture-White", active(piece.White), "Rba3", "b3a3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingCapture-Black", active(piece.Black), "Rba6", "b6a6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidRedundantRank-White", active(piece.White), "Rb3xa3", "b3a3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidRedundantRank-Black", active(piece.Black), "Rb6xa6", "b6a6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongType-White", active(piece.White), "Kbxa3", "", errors.New("could not find source square of 'Kbxa3'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongType-Black", active(piece.Black), "Kbxa6", "", errors.New("could not find source square of 'Kbxa6'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongFile-White", active(piece.White), "Raxa3", "", errors.New("could not find source square of 'Raxa3'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongFile-Black", active(piece.Black), "Rcxa6", "", errors.New("could not find source square of 'Rcxa6'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongRedundantRank-White", active(piece.White), "Rb2xa3", "b2a3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongRedundantRank-Black", active(piece.Black), "Rb7xa6", "b7a6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidPromo-White", active(piece.White), "Rbxa3=R", "b3a3r", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidPromo-Black", active(piece.Black), "Rbxa6=r", "b6a6r", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidCheck-White", active(piece.White), "Rbxa3+", "b3a3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMate-Black", active(piece.Black), "Rbxa6#", "b6a6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidPromoCheck-White", active(piece.White), "Rbxa3=R+", "b3a3r", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidPromoMate-Black", active(piece.Black), "Rbxa6=r#", "b6a6r", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-White", active(piece.White), "R3xb4", "b3b4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-Black", active(piece.Black), "R4xb5", "b4b5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingType-White", active(piece.White), "3xb4", "", errors.New("could not find source square of '3xb4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingType-Black", active(piece.Black), "4xb5", "", errors.New("could not find source square of '4xb5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingRank-White", active(piece.White), "Rxb4", "", errors.New("could not find source square of 'Rxb4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingRank-Black", active(piece.Black), "Rxb5", "", errors.New("could not find source square of 'Rxb5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingCapture-White", active(piece.White), "R3b4", "b3b4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingCapture-Black", active(piece.Black), "R4b5", "b4b5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidRedundantFile-White", active(piece.White), "Rb3xb4", "b3b4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidRedundantFile-Black", active(piece.Black), "Rb4xb5", "b4b5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongType-White", active(piece.White), "K3xb4", "", errors.New("could not find source square of 'K3xb4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongType-Black", active(piece.Black), "K4xb5", "", errors.New("could not find source square of 'K4xb5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongRank-White", active(piece.White), "R1xb4", "", errors.New("could not find source square of 'R1xb4'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongRank-Black", active(piece.Black), "R8xb5", "", errors.New("could not find source square of 'R8xb5'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongRedundantFile-White", active(piece.White), "Ra3xb4", "a3b4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongRedundantFile-Black", active(piece.Black), "Ra4xb5", "a4b5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidPromo-White", active(piece.White), "R3xb4=b", "b3b4b", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidPromo-Black", active(piece.Black), "R4xb5=B", "b4b5b", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMate-White", active(piece.White), "R3xb4#", "b3b4", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidCheck-Black", active(piece.Black), "R4xb5+", "b4b5", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidPromoMate-White", active(piece.White), "R3xb4=b#", "b3b4b", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidPromoCheck-Black", active(piece.Black), "R4xb5=B+", "b4b5b", nil},
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
		[]testParseMoveTestCase{
			{"SAN-Move-White", active(piece.White), "Qc1", "b1c1", nil},
			{"SAN-Move-Black", active(piece.Black), "Qc8", "b8c8", nil},
			{"SAN-Move-InvalidMissingType-White", active(piece.White), "c1", "", errors.New("could not find source square of 'c1'")},
			{"SAN-Move-InvalidMissingType-Black", active(piece.Black), "c8", "", errors.New("could not find source square of 'c8'")},
			{"SAN-Move-InvalidRedundantFile-White", active(piece.White), "Qbc1", "b1c1", nil},
			{"SAN-Move-InvalidRedundantFile-Black", active(piece.Black), "Qbc8", "b8c8", nil},
			{"SAN-Move-InvalidRedundantRank-White", active(piece.White), "Q1c1", "b1c1", nil},
			{"SAN-Move-InvalidRedundantRank-Black", active(piece.Black), "Q8c8", "b8c8", nil},
			{"SAN-Move-InvalidRedundantFileRank-White", active(piece.White), "Qb1c1", "b1c1", nil},
			{"SAN-Move-InvalidRedundantFileRank-Black", active(piece.Black), "Qb8c8", "b8c8", nil},
			{"SAN-Move-InvalidWrongType-White", active(piece.White), "Kc1", "", errors.New("could not find source square of 'Kc1'")},
			{"SAN-Move-InvalidWrongType-Black", active(piece.Black), "Kc8", "", errors.New("could not find source square of 'Kc8'")},
			{"SAN-Move-InvalidWrongRedundantFile-White", active(piece.White), "Qac1", "", errors.New("could not find source square of 'Qac1'")},
			{"SAN-Move-InvalidWrongRedundantFile-Black", active(piece.Black), "Qac8", "", errors.New("could not find source square of 'Qac8'")},
			{"SAN-Move-InvalidWrongRedundantRank-White", active(piece.White), "Q8c1", "", errors.New("could not find source square of 'Q8c1'")},
			{"SAN-Move-InvalidWrongRedundantRank-Black", active(piece.Black), "Q1c8", "", errors.New("could not find source square of 'Q1c8'")},
			{"SAN-Move-InvalidWrongRedundantFileRank-White", active(piece.White), "Qa8c1", "a8c1", nil},
			{"SAN-Move-InvalidWrongRedundantFileRank-Black", active(piece.Black), "Qa1c8", "a1c8", nil},
			{"SAN-Move-InvalidNoQueenOrigin-White", active(piece.White), "Qf1", "", errors.New("could not find source square of 'Qf1'")},
			{"SAN-Move-InvalidNoQueenOrigin-Black", active(piece.Black), "Qf8", "", errors.New("could not find source square of 'Qf8'")},
			{"SAN-Move-InvalidCapture-White", active(piece.White), "Qxc1", "b1c1", nil},
			{"SAN-Move-InvalidCapture-Black", active(piece.Black), "Qxc8", "b8c8", nil},
			{"SAN-Move-InvalidFile-White", active(piece.White), "Qcc1", "", errors.New("could not find source square of 'Qcc1'")},
			{"SAN-Move-InvalidFile-Black", active(piece.Black), "Qcc8", "", errors.New("could not find source square of 'Qcc8'")},
			{"SAN-Move-InvalidRank-White", active(piece.White), "Q2c1", "", errors.New("could not find source square of 'Q2c1'")},
			{"SAN-Move-InvalidRank-Black", active(piece.Black), "Q7c8", "", errors.New("could not find source square of 'Q7c8'")},
			{"SAN-Move-InvalidFileRank-White", active(piece.White), "Qc1c1", "c1c1", nil},
			{"SAN-Move-InvalidFileRank-Black", active(piece.Black), "Qc8c8", "c8c8", nil},
			{"SAN-Move-InvalidPromo-White", active(piece.White), "Qc1=Q", "b1c1q", nil},
			{"SAN-Move-InvalidPromo-Black", active(piece.Black), "Qc8=q", "b8c8q", nil},
			{"SAN-Move-InvalidCheck-White", active(piece.White), "Qc1+", "b1c1", nil},
			{"SAN-Move-InvalidMate-Black", active(piece.Black), "Qc8#", "b8c8", nil},
			{"SAN-Move-InvalidPromoCheck-White", active(piece.White), "Qc1=Q+", "b1c1q", nil},
			{"SAN-Move-InvalidPromoMate-Black", active(piece.Black), "Qc8=q#", "b8c8q", nil},
			{"SAN-AmbiguousMove-SpecifyFile-White", active(piece.White), "Qbb2", "b1b2", nil},
			{"SAN-AmbiguousMove-SpecifyFile-Black", active(piece.Black), "Qbb7", "b8b7", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMissingType-White", active(piece.White), "bb2", "", errors.New("could not find source square of 'bb2'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMissingType-Black", active(piece.Black), "bb7", "", errors.New("could not find source square of 'bb7'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMissingFile-White", active(piece.White), "Qb2", "", errors.New("could not find source square of 'Qb2'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMissingFile-Black", active(piece.Black), "Qb7", "", errors.New("could not find source square of 'Qb7'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidRedundantRank-White", active(piece.White), "Qb1b2", "b1b2", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidRedundantRank-Black", active(piece.Black), "Qb8b7", "b8b7", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongType-White", active(piece.White), "Kbb2", "", errors.New("could not find source square of 'Kbb2'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongType-Black", active(piece.Black), "Kbb7", "", errors.New("could not find source square of 'Kbb7'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongFile-White", active(piece.White), "Qhb2", "", errors.New("could not find source square of 'Qhb2'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongFile-Black", active(piece.Black), "Qhb7", "", errors.New("could not find source square of 'Qhb7'")},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongRedundantRank-White", active(piece.White), "Qb8b2", "b8b2", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidWrongRedundantRank-Black", active(piece.Black), "Qb1b7", "b1b7", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidCapture-White", active(piece.White), "Qbxb2", "b1b2", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidCapture-Black", active(piece.Black), "Qbxb7", "b8b7", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidPromo-White", active(piece.White), "Qbb2=r", "b1b2r", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidPromo-Black", active(piece.Black), "Qbb7=R", "b8b7r", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidMate-White", active(piece.White), "Qbb2#", "b1b2", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidCheck-Black", active(piece.Black), "Qbb7+", "b8b7", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidPromoMate-White", active(piece.White), "Qbb2=r#", "b1b2r", nil},
			{"SAN-AmbiguousMove-SpecifyFile-InvalidPromoCheck-Black", active(piece.Black), "Qbb7=R+", "b8b7r", nil},
			{"SAN-AmbiguousMove-SpecifyRank-White", active(piece.White), "Q2b2", "a2b2", nil},
			{"SAN-AmbiguousMove-SpecifyRank-Black", active(piece.Black), "Q7b7", "a7b7", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMissingType-White", active(piece.White), "2b2", "", errors.New("could not find source square of '2b2'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMissingType-Black", active(piece.Black), "7b7", "", errors.New("could not find source square of '7b7'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMissingRank-White", active(piece.White), "Qb2", "", errors.New("could not find source square of 'Qb2'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMissingRank-Black", active(piece.Black), "Qb7", "", errors.New("could not find source square of 'Qb7'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidRedundantFile-White", active(piece.White), "Qa2b2", "a2b2", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidRedundantFile-Black", active(piece.Black), "Qa7b7", "a7b7", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongType-White", active(piece.White), "K2b2", "", errors.New("could not find source square of 'K2b2'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongType-Black", active(piece.Black), "K7b7", "", errors.New("could not find source square of 'K7b7'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongRank-White", active(piece.White), "Q8b2", "", errors.New("could not find source square of 'Q8b2'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongRank-Black", active(piece.Black), "Q1b7", "", errors.New("could not find source square of 'Q1b7'")},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongRedundantFile-White", active(piece.White), "Qh2b2", "h2b2", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidWrongRedundantFile-Black", active(piece.Black), "Qh7b7", "h7b7", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidCapture-White", active(piece.White), "Q2xb2", "a2b2", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidCapture-Black", active(piece.Black), "Q7xb7", "a7b7", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidPromo-White", active(piece.White), "Q2b2=B", "a2b2b", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidPromo-Black", active(piece.Black), "Q7b7=b", "a7b7b", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidCheck-White", active(piece.White), "Q2b2+", "a2b2", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidMate-Black", active(piece.Black), "Q7b7#", "a7b7", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidPromoCheck-White", active(piece.White), "Q2b2=B+", "a2b2b", nil},
			{"SAN-AmbiguousMove-SpecifyRank-InvalidPromoMate-Black", active(piece.Black), "Q7b7=b#", "a7b7b", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-White", active(piece.White), "Qa1b2", "a1b2", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-Black", active(piece.Black), "Qa8b7", "a8b7", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingType-White", active(piece.White), "a1b2", "a1b2", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingType-Black", active(piece.Black), "a8b7", "a8b7", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingFile-White", active(piece.White), "Q1b2", "", errors.New("could not find source square of 'Q1b2'")},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingFile-Black", active(piece.Black), "Q8b7", "", errors.New("could not find source square of 'Q8b7'")},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingRank-White", active(piece.White), "Qab2", "", errors.New("could not find source square of 'Qab2'")},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingRank-Black", active(piece.Black), "Qab7", "", errors.New("could not find source square of 'Qab7'")},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingFileRank-White", active(piece.White), "Qb2", "", errors.New("could not find source square of 'Qb2'")},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMissingFileRank-Black", active(piece.Black), "Qb7", "", errors.New("could not find source square of 'Qb7'")},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongType-White", active(piece.White), "Ka1b2", "a1b2", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongType-Black", active(piece.Black), "Ka8b7", "a8b7", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongFile-White", active(piece.White), "Qh1b2", "h1b2", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongFile-Black", active(piece.Black), "Qh8b7", "h8b7", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongRank-White", active(piece.White), "Qa8b2", "a8b2", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongRank-Black", active(piece.Black), "Qa1b7", "a1b7", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongFileRank-White", active(piece.White), "Qh8b2", "h8b2", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidWrongFileRank-Black", active(piece.Black), "Qh1b7", "h1b7", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidCapture-White", active(piece.White), "Qa1xb2", "a1b2", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidCapture-Black", active(piece.Black), "Qa8xb7", "a8b7", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidPromo-White", active(piece.White), "Qa1b2=n", "a1b2n", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidPromo-Black", active(piece.Black), "Qa8b7=N", "a8b7n", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidMate-White", active(piece.White), "Qa1b2#", "a1b2", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidCheck-Black", active(piece.Black), "Qa8b7+", "a8b7", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidPromoMate-White", active(piece.White), "Qa1b2=n#", "a1b2n", nil},
			{"SAN-AmbiguousMove-SpecifyFileRank-InvalidPromoCheck-Black", active(piece.Black), "Qa8b7=N+", "a8b7n", nil},
			{"SAN-CaptureMove-White", active(piece.White), "Qxa3", "a2a3", nil},
			{"SAN-CaptureMove-Black", active(piece.Black), "Qxa6", "a7a6", nil},
			{"SAN-CaptureMove-InvalidMissingType-White", active(piece.White), "xa3", "", errors.New("could not find source square of 'xa3'")},
			{"SAN-CaptureMove-InvalidMissingType-Black", active(piece.Black), "xa6", "", errors.New("could not find source square of 'xa6'")},
			{"SAN-CaptureMove-InvalidMissingCapture-White", active(piece.White), "Qa3", "a2a3", nil},
			{"SAN-CaptureMove-InvalidMissingCapture-Black", active(piece.Black), "Qa6", "a7a6", nil},
			{"SAN-CaptureMove-InvalidRedundantFile-White", active(piece.White), "Qaxa3", "a2a3", nil},
			{"SAN-CaptureMove-InvalidRedundantFile-Black", active(piece.Black), "Qaxa6", "a7a6", nil},
			{"SAN-CaptureMove-InvalidRedundantRank-White", active(piece.White), "Q2xa3", "a2a3", nil},
			{"SAN-CaptureMove-InvalidRedundantRank-Black", active(piece.Black), "Q7xa6", "a7a6", nil},
			{"SAN-CaptureMove-InvalidRedundantFileRank-White", active(piece.White), "Qa2xa3", "a2a3", nil},
			{"SAN-CaptureMove-InvalidRedundantFileRank-Black", active(piece.Black), "Qa7xa6", "a7a6", nil},
			{"SAN-CaptureMove-InvalidWrongType-White", active(piece.White), "Kxa3", "", errors.New("could not find source square of 'Kxa3'")},
			{"SAN-CaptureMove-InvalidWrongType-Black", active(piece.Black), "Kxa6", "", errors.New("could not find source square of 'Kxa6'")},
			{"SAN-CaptureMove-InvalidWrongRedundantFile-White", active(piece.White), "Qhxa3", "", errors.New("could not find source square of 'Qhxa3'")},
			{"SAN-CaptureMove-InvalidWrongRedundantFile-Black", active(piece.Black), "Qhxa6", "", errors.New("could not find source square of 'Qhxa6'")},
			{"SAN-CaptureMove-InvalidWrongRedundantRank-White", active(piece.White), "Q8xa3", "", errors.New("could not find source square of 'Q8xa3'")},
			{"SAN-CaptureMove-InvalidWrongRedundantRank-Black", active(piece.Black), "Q1xa6", "", errors.New("could not find source square of 'Q1xa6'")},
			{"SAN-CaptureMove-InvalidWrongRedundantFileRank-White", active(piece.White), "Qh8xa3", "h8a3", nil},
			{"SAN-CaptureMove-InvalidWrongRedundantFileRank-Black", active(piece.Black), "Qh1xa6", "h1a6", nil},
			{"SAN-CaptureMove-InvalidPromo-White", active(piece.White), "Qxa3=q", "a2a3q", nil},
			{"SAN-CaptureMove-InvalidPromo-Black", active(piece.Black), "Qxa6=Q", "a7a6q", nil},
			{"SAN-CaptureMove-InvalidMate-White", active(piece.White), "Qxa3#", "a2a3", nil},
			{"SAN-CaptureMove-InvalidCheck-Black", active(piece.Black), "Qxa6+", "a7a6", nil},
			{"SAN-CaptureMove-InvalidPromoMate-White", active(piece.White), "Qxa3=q#", "a2a3q", nil},
			{"SAN-CaptureMove-InvalidPromoCheck-Black", active(piece.Black), "Qxa6=Q+", "a7a6q", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-White", active(piece.White), "Qhxh3", "h4h3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-Black", active(piece.Black), "Qhxh6", "h5h6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingType-White", active(piece.White), "hxh3", "", errors.New("could not find source square of 'hxh3'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingType-Black", active(piece.Black), "hxh6", "", errors.New("could not find source square of 'hxh6'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingFile-White", active(piece.White), "Qxh3", "", errors.New("could not find source square of 'Qxh3'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingFile-Black", active(piece.Black), "Qxh6", "", errors.New("could not find source square of 'Qxh6'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingCapture-White", active(piece.White), "Qhh3", "h4h3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMissingCapture-Black", active(piece.Black), "Qhh6", "h5h6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidRedundantRank-White", active(piece.White), "Qh4xh3", "h4h3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidRedundantRank-Black", active(piece.Black), "Qh5xh6", "h5h6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongType-White", active(piece.White), "Khxh3", "", errors.New("could not find source square of 'Khxh3'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongType-Black", active(piece.Black), "Khxh6", "", errors.New("could not find source square of 'Khxh6'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongFile-White", active(piece.White), "Qaxh3", "", errors.New("could not find source square of 'Qaxh3'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongFile-Black", active(piece.Black), "Qaxh6", "", errors.New("could not find source square of 'Qaxh6'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongRedundantRank-White", active(piece.White), "Qh8xh3", "h8h3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidWrongRedundantRank-Black", active(piece.Black), "Qh1xh6", "h1h6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidPromo-White", active(piece.White), "Qhxh3=R", "h4h3r", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidPromo-Black", active(piece.Black), "Qhxh6=r", "h5h6r", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidCheck-White", active(piece.White), "Qhxh3+", "h4h3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidMate-Black", active(piece.Black), "Qhxh6#", "h5h6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidPromoCheck-White", active(piece.White), "Qhxh3=R+", "h4h3r", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFile-InvalidPromoMate-Black", active(piece.Black), "Qhxh6=r#", "h5h6r", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-White", active(piece.White), "Q3xh3", "g3h3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-Black", active(piece.Black), "Q6xh6", "g6h6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingType-White", active(piece.White), "3xh3", "", errors.New("could not find source square of '3xh3'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingType-Black", active(piece.Black), "6xh6", "", errors.New("could not find source square of '6xh6'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingRank-White", active(piece.White), "Qxh3", "", errors.New("could not find source square of 'Qxh3'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingRank-Black", active(piece.Black), "Qxh6", "", errors.New("could not find source square of 'Qxh6'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingCapture-White", active(piece.White), "Q3h3+", "g3h3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMissingCapture-Black", active(piece.Black), "Q6h6+", "g6h6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidRedundantFile-White", active(piece.White), "Qg3xh3", "g3h3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidRedundantFile-Black", active(piece.Black), "Qg6xh6", "g6h6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongType-White", active(piece.White), "K3xh3", "", errors.New("could not find source square of 'K3xh3'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongType-Black", active(piece.Black), "K6xh6", "", errors.New("could not find source square of 'K6xh6'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongRank-White", active(piece.White), "Q8xh3", "", errors.New("could not find source square of 'Q8xh3'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongRank-Black", active(piece.Black), "Q1xh6", "", errors.New("could not find source square of 'Q1xh6'")},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongRedundantFile-White", active(piece.White), "Qa3xh3", "a3h3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidWrongRedundantFile-Black", active(piece.Black), "Qa6xh6", "a6h6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidPromo-White", active(piece.White), "Q3xh3=b", "g3h3b", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidPromo-Black", active(piece.Black), "Q6xh6=B", "g6h6b", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidMate-White", active(piece.White), "Q3xh3#", "g3h3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidCheck-Black", active(piece.Black), "Q6xh6+", "g6h6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidPromoMate-White", active(piece.White), "Q3xh3=b#", "g3h3b", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyRank-InvalidPromoCheck-Black", active(piece.Black), "Q6xh6=B+", "g6h6b", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-White", active(piece.White), "Qg4xh3", "g4h3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-Black", active(piece.Black), "Qg5xh6", "g5h6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingType-White", active(piece.White), "g4xh3", "g4h3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingType-Black", active(piece.Black), "g5xh6", "g5h6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingFile-White", active(piece.White), "Q4xh3", "", errors.New("could not find source square of 'Q4xh3'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingFile-Black", active(piece.Black), "Q5xh6", "", errors.New("could not find source square of 'Q5xh6'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingRank-White", active(piece.White), "Qgxh3", "", errors.New("could not find source square of 'Qgxh3'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingRank-Black", active(piece.Black), "Qgxh6", "", errors.New("could not find source square of 'Qgxh6'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingFileRank-White", active(piece.White), "Qxh3", "", errors.New("could not find source square of 'Qxh3'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingFileRank-Black", active(piece.Black), "Qxh6", "", errors.New("could not find source square of 'Qxh6'")},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingCapture-White", active(piece.White), "Qg4h3", "g4h3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMissingCapture-Black", active(piece.Black), "Qg5h6", "g5h6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongType-White", active(piece.White), "Kg4xh3", "g4h3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongType-Black", active(piece.Black), "Kg5xh6", "g5h6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongFile-White", active(piece.White), "Qa4xh3", "a4h3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongFile-Black", active(piece.Black), "Qa5xh6", "a5h6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongRank-White", active(piece.White), "Qg8xh3", "g8h3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongRank-Black", active(piece.Black), "Qg1xh6", "g1h6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongFileRank-White", active(piece.White), "Qa8xh3", "a8h3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidWrongFileRank-Black", active(piece.Black), "Qa1xh6", "a1h6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidPromo-White", active(piece.White), "Qg4xh3=N", "g4h3n", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidPromo-Black", active(piece.Black), "Qg5xh6=n", "g5h6n", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidCheck-White", active(piece.White), "Qg4xh3+", "g4h3", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidMate-Black", active(piece.Black), "Qg5xh6#", "g5h6", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidPromoCheck-White", active(piece.White), "Qg4xh3=N+", "g4h3n", nil},
			{"SAN-AmbiguousCaptureMove-SpecifyFileRank-InvalidPromoMate-Black", active(piece.Black), "Qg5xh6=n#", "g5h6n", nil},
		},
	},
	{"King",
		map[square.Square]piece.Piece{
			square.E1: piece.New(piece.White, piece.King),
			square.H1: piece.New(piece.White, piece.Rook),
			square.F7: piece.New(piece.White, piece.Rook),

			square.E8: piece.New(piece.Black, piece.King),
			square.A8: piece.New(piece.Black, piece.Rook),
			square.D2: piece.New(piece.Black, piece.Rook),

			// r . . . k . . . 8
			// . . . . . R . . 7
			// . . . . . . . . 6
			// . . . . . . . . 5
			// . . . . . . . . 4
			// . . . . . . . . 3
			// . . . r . . . . 2
			// . . . . K . . R 1
			// a b c d e f g h
		},
		[]testParseMoveTestCase{
			{"SAN-Castling-KingSide-White", active(piece.White), "O-O", "e1g1", nil},
			{"SAN-InvalidCastling-KingSide-Black", active(piece.Black), "O-O", "e8g8", nil},
			{"SAN-Castling-KingSide-KingMove2Places-White", active(piece.White), "Kg1", "e1g1", nil},
			{"SAN-InvalidCastling-KingSide-KingMove2Places-Black", active(piece.Black), "Kg8", "", errors.New("could not find source square of 'Kg8'")},
			{"SAN-InvalidCastling-KingSide-KingMoveOnRookPlace-White", active(piece.White), "Kh1", "", errors.New("could not find source square of 'Kh1'")},
			{"SAN-InvalidCastling-KingSide-KingMoveOnRookPlace-Black", active(piece.Black), "Kh8", "", errors.New("could not find source square of 'Kh8'")},
			{"SAN-InvalidCastling-QueenSide-White", active(piece.White), "O-O-O", "e1c1", nil},
			{"SAN-Castling-QueenSide-Black", active(piece.Black), "O-O-O", "e8c8", nil},
			{"SAN-InvalidCastling-QueenSide-KingMove2Places-White", active(piece.White), "Kc1", "", errors.New("could not find source square of 'Kc1'")},
			{"SAN-Castling-QueenSide-KingMove2Places-Black", active(piece.Black), "Kc8", "e8c8", nil},
			{"SAN-InvalidCastling-QueenSide-KingMoveOnRookPlace-White", active(piece.White), "Ka1", "", errors.New("could not find source square of 'Ka1'")},
			{"SAN-InvalidCastling-QueenSide-KingMoveOnRookPlace-Black", active(piece.Black), "Ka8", "", errors.New("could not find source square of 'Ka8'")},
			{"SAN-Move-White", active(piece.White), "Kf1", "e1f1", nil},
			{"SAN-Move-Black", active(piece.Black), "Kd8", "e8d8", nil},
			{"SAN-Move-InvalidMissingType-White", active(piece.White), "f1", "", errors.New("could not find source square of 'f1'")},
			{"SAN-Move-InvalidMissingType-Black", active(piece.Black), "d8", "", errors.New("could not find source square of 'd8'")},
			{"SAN-Move-InvalidRedundantFile-White", active(piece.White), "Kef1", "e1f1", nil},
			{"SAN-Move-InvalidRedundantFile-Black", active(piece.Black), "Ked8", "e8d8", nil},
			{"SAN-Move-InvalidRedundantRank-White", active(piece.White), "K1f1", "e1f1", nil},
			{"SAN-Move-InvalidRedundantRank-Black", active(piece.Black), "K8d8", "e8d8", nil},
			{"SAN-Move-InvalidRedundantFileRank-White", active(piece.White), "Ke1f1", "e1f1", nil},
			{"SAN-Move-InvalidRedundantFileRank-Black", active(piece.Black), "Ke8d8", "e8d8", nil},
			{"SAN-Move-InvalidWrongType-White", active(piece.White), "Qf1", "", errors.New("could not find source square of 'Qf1'")},
			{"SAN-Move-InvalidWrongType-Black", active(piece.Black), "Qd8", "", errors.New("could not find source square of 'Qd8'")},
			{"SAN-Move-InvalidWrongRedundantFile-White", active(piece.White), "Kaf1", "", errors.New("could not find source square of 'Kaf1'")},
			{"SAN-Move-InvalidWrongRedundantFile-Black", active(piece.Black), "Kad8", "", errors.New("could not find source square of 'Kad8'")},
			{"SAN-Move-InvalidWrongRedundantRank-White", active(piece.White), "K8f1", "", errors.New("could not find source square of 'K8f1'")},
			{"SAN-Move-InvalidWrongRedundantRank-Black", active(piece.Black), "K1d8", "", errors.New("could not find source square of 'K1d8'")},
			{"SAN-Move-InvalidWrongRedundantFileRank-White", active(piece.White), "Ka8f1", "a8f1", nil},
			{"SAN-Move-InvalidWrongRedundantFileRank-Black", active(piece.Black), "Ka1d8", "a1d8", nil},
			{"SAN-Move-InvalidNoKingOrigin-White", active(piece.White), "Kh2", "", errors.New("could not find source square of 'Kh2'")},
			{"SAN-Move-InvalidNoKingOrigin-Black", active(piece.Black), "Ka7", "", errors.New("could not find source square of 'Ka7'")},
			{"SAN-Move-InvalidCapture-White", active(piece.White), "Kxf1", "e1f1", nil},
			{"SAN-Move-InvalidCapture-Black", active(piece.Black), "Kxd8", "e8d8", nil},
			{"SAN-Move-InvalidFile-White", active(piece.White), "Kff1", "", errors.New("could not find source square of 'Kff1'")},
			{"SAN-Move-InvalidFile-Black", active(piece.Black), "Kdd8", "", errors.New("could not find source square of 'Kdd8'")},
			{"SAN-Move-InvalidRank-White", active(piece.White), "K2f1", "", errors.New("could not find source square of 'K2f1'")},
			{"SAN-Move-InvalidRank-Black", active(piece.Black), "K7d8", "", errors.New("could not find source square of 'K7d8'")},
			{"SAN-Move-InvalidFileRank-White", active(piece.White), "Kf2f1", "f2f1", nil},
			{"SAN-Move-InvalidFileRank-Black", active(piece.Black), "Kd7d8", "d7d8", nil},
			{"SAN-Move-InvalidPromo-White", active(piece.White), "Kf1=Q", "e1f1q", nil},
			{"SAN-Move-InvalidPromo-Black", active(piece.Black), "Kd8=q", "e8d8q", nil},
			{"SAN-Move-InvalidCheck-White", active(piece.White), "Kf1+", "e1f1", nil},
			{"SAN-Move-InvalidMate-Black", active(piece.Black), "Kd8#", "e8d8", nil},
			{"SAN-Move-InvalidPromoCheck-White", active(piece.White), "Kf1=Q+", "e1f1q", nil},
			{"SAN-Move-InvalidPromoMate-Black", active(piece.Black), "Kd8=q#", "e8d8q", nil},
			{"SAN-CaptureMove-White", active(piece.White), "Kxd2", "e1d2", nil},
			{"SAN-CaptureMove-Black", active(piece.Black), "Kxf7", "e8f7", nil},
			{"SAN-CaptureMove-InvalidMissingType-White", active(piece.White), "xd2", "", errors.New("could not find source square of 'xd2'")},
			{"SAN-CaptureMove-InvalidMissingType-Black", active(piece.Black), "xf7", "", errors.New("could not find source square of 'xf7'")},
			{"SAN-CaptureMove-InvalidRedundantFile-White", active(piece.White), "Kexd2", "e1d2", nil},
			{"SAN-CaptureMove-InvalidRedundantFile-Black", active(piece.Black), "Kexf7", "e8f7", nil},
			{"SAN-CaptureMove-InvalidRedundantRank-White", active(piece.White), "K1xd2", "e1d2", nil},
			{"SAN-CaptureMove-InvalidRedundantRank-Black", active(piece.Black), "K8xf7", "e8f7", nil},
			{"SAN-CaptureMove-InvalidRedundantFileRank-White", active(piece.White), "Ke1xd2", "e1d2", nil},
			{"SAN-CaptureMove-InvalidRedundantFileRank-Black", active(piece.Black), "Ke8xf7", "e8f7", nil},
			{"SAN-CaptureMove-InvalidWrongType-White", active(piece.White), "Qxd2", "", errors.New("could not find source square of 'Qxd2'")},
			{"SAN-CaptureMove-InvalidWrongType-Black", active(piece.Black), "Qxf7", "", errors.New("could not find source square of 'Qxf7'")},
			{"SAN-CaptureMove-InvalidWrongRedundantFile-White", active(piece.White), "Kaxd2", "", errors.New("could not find source square of 'Kaxd2'")},
			{"SAN-CaptureMove-InvalidWrongRedundantFile-Black", active(piece.Black), "Kaxf7", "", errors.New("could not find source square of 'Kaxf7'")},
			{"SAN-CaptureMove-InvalidWrongRedundantRank-White", active(piece.White), "K8xd2", "", errors.New("could not find source square of 'K8xd2'")},
			{"SAN-CaptureMove-InvalidWrongRedundantRank-Black", active(piece.Black), "K1xf7", "", errors.New("could not find source square of 'K1xf7'")},
			{"SAN-CaptureMove-InvalidWrongRedundantFileRank-White", active(piece.White), "Ka8xd2", "a8d2", nil},
			{"SAN-CaptureMove-InvalidWrongRedundantFileRank-Black", active(piece.Black), "Ka1xf7", "a1f7", nil},
			{"SAN-CaptureMove-InvalidPromo-White", active(piece.White), "Kxd2=q", "e1d2q", nil},
			{"SAN-CaptureMove-InvalidPromo-Black", active(piece.Black), "Kxf7=Q", "e8f7q", nil},
			{"SAN-CaptureMove-InvalidCheck-White", active(piece.White), "Kxd2+", "e1d2", nil},
			{"SAN-CaptureMove-InvalidMate-Black", active(piece.Black), "Kxf7#", "e8f7", nil},
			{"SAN-CaptureMove-InvalidPromoCheck-White", active(piece.White), "Kxd2=q+", "e1d2q", nil},
			{"SAN-CaptureMove-InvalidPromoMate-Black", active(piece.Black), "Kxf7=Q#", "e8f7q", nil},
		},
	},
}

// Returns a positionChanger function, which sets the active color.
func active(color piece.Color) positionChanger {
	return func(inPos Position) (outPos Position, outErr error) {
		outPos = *Copy(&inPos)
		outPos.ActiveColor = color
		return outPos, nil
	}
}

// Returns a positionChanger function, which sets the en passant property to a given square.
func enPassant(sq square.Square) positionChanger {
	return func(inPos Position) (outPos Position, outErr error) {
		outPos = *Copy(&inPos)
		outPos.EnPassant = sq
		return outPos, nil
	}
}

// Returns a positionChanger function, which changes a piece at position.
func pos(sq square.Square, pc piece.Piece) positionChanger {
	return func(inPos Position) (outPos Position, outErr error) {
		outPos = *Copy(&inPos)
		outPos.Put(pc, sq)
		return outPos, nil
	}
}

// Returns a positionChanger function, which applies all the given positionChangers.
func multi(positionChangers ...positionChanger) positionChanger {
	return func(inPos Position) (outPos Position, outErr error) {
		outPos = *Copy(&inPos)
		for _, pc := range positionChangers {
			if pc != nil {
				if changedPosition, err := pc(outPos); err != nil {
					return outPos, fmt.Errorf("positionChanger error: %v", err)
				} else {
					outPos = changedPosition
				}
			}
		}
		return outPos, nil
	}
}

func BenchmarkParseMove(b *testing.B) {
	benchGroups := []struct {
		Name       string
		filterFunc func(tc testParseMoveTestCase) bool
	}{
		{"all-testcases", func(tc testParseMoveTestCase) bool { return true }},
		{"all-non-error-testcases", func(tc testParseMoveTestCase) bool {
			if tc.WantError != nil {
				return false
			}
			return true
		}},
		{"all-valid-non-error-testcases", func(tc testParseMoveTestCase) bool {
			if strings.Contains(tc.Name, "Invalid") || tc.WantError != nil {
				return false
			}
			return true
		}},
		{"PCN-testcases", func(tc testParseMoveTestCase) bool {
			if !strings.Contains(tc.Name, "PCN") {
				return false
			}
			return true
		}},
		{"PCN-non-error-testcases", func(tc testParseMoveTestCase) bool {
			if !strings.Contains(tc.Name, "PCN") || tc.WantError != nil {
				return false
			}
			return true
		}},
		{"PCN-valid-non-error-testcases", func(tc testParseMoveTestCase) bool {
			if !strings.Contains(tc.Name, "PCN") || strings.Contains(tc.Name, "Invalid") || tc.WantError != nil {
				return false
			}
			return true
		}},
		{"SAN-testcases", func(tc testParseMoveTestCase) bool {
			if !strings.Contains(tc.Name, "SAN") {
				return false
			}
			return true
		}},
		{"SAN-non-error-testcases", func(tc testParseMoveTestCase) bool {
			if !strings.Contains(tc.Name, "SAN") || tc.WantError != nil {
				return false
			}
			return true
		}},
		{"SAN-valid-non-error-testcases", func(tc testParseMoveTestCase) bool {
			if !strings.Contains(tc.Name, "SAN") || strings.Contains(tc.Name, "Invalid") || tc.WantError != nil {
				return false
			}
			return true
		}},
	}

	// Loop through benchmark groups.
	for _, bg := range benchGroups {
		// Get benchmark cases from test cases dataset using filter function.
		benches, err := benchmarkParseMoveBenches(bg.filterFunc)
		if err != nil {
			b.Fatal("Error getting benchmark cases:", err)
		}

		//if bcm, err := json.MarshalIndent(benches[:2], "", "  "); err != nil {
		//	b.Log("error:", err)
		//} else {
		//	b.Log("size:", len(benches))
		//	b.Log(string(bcm))
		//}
		b.ResetTimer()
		// Run test group benchmarks.
		b.Run(bg.Name, func(b *testing.B) {
			benchmarkParseMove(b, benches)
		})
	}
}

type benchParseMove struct {
	Name     string
	position Position
	Move     string
}

// Returns a `[]benchParseMove` filled with data from `testParseMoveGroups` filtered by the function provided in input.
// The input filter function gets a `testParseMoveTestCase` structure as input and if the function returns true, the test case's piece positions with move string is added to returned result.
func benchmarkParseMoveBenches(filterFunc func(tc testParseMoveTestCase) bool) ([]benchParseMove, error) {
	res := make([]benchParseMove, 0)
	for _, group := range testParseMoveGroups {
		for _, tc := range group.TestCases {
			if filterFunc(tc) {
				p, err := testCasePosition(group.Position, tc.positionChangerFunc)
				if err != nil {
					return nil, err
				}
				res = append(res, benchParseMove{tc.Name, *p, tc.Move})
			}
		}
	}
	return res, nil
}

// Benchmarks Position.ParseMove function on array of benchmarks.
func benchmarkParseMove(b *testing.B, benches []benchParseMove) {
	for i := 0; i < b.N; i++ {
		// Get benchmark.
		bpm := benches[i%len(benches)]
		// Call ParseMove function with test case input on Position.
		m, err := bpm.position.ParseMove(bpm.Move)
		// Use result to be sure not skip the ParseMove function due to optimizations.
		_, _ = m, err
	}
}
