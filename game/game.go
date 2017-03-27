// Package game handles chess games. You can create timed and untimed games.
package game

import (
	"fmt"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/polyglot"
	"github.com/andrewbackes/chess/position"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
	"time"
)

// Game represents a chess game.
type Game struct {
	Tags          map[string]string
	control       [2]TimeControl
	Position      *position.Position
	Moves         []move.Move
	positionCache map[uint64]int
}

// New returns a fresh game with all of the pieces in the
// opening position.
func New() *Game {
	return &Game{
		control:       [2]TimeControl{},
		Tags:          make(map[string]string),
		Position:      position.New(),
		positionCache: make(map[uint64]int),
	}
}

// NewTimedGame does the same thing as NewGame() but sets the
// time control to what is specified.
func NewTimedGame(control [2]TimeControl) *Game {
	g := New()
	g.control = control
	g.control[piece.White].Reset()
	g.control[piece.Black].Reset()
	return g
}

// ActiveColor returns the color of the player whos turn it is.
func (G *Game) ActiveColor() piece.Color {
	return G.Position.ActiveColor
}

// QuickMove makes the specified move without checking the legality
// of the move or the status of the game post move.
func (G *Game) QuickMove(m *move.Move) {
	from, to, movingPiece, capturedPiece := G.decompose(m)
	G.makeMove(m, from, to, movingPiece, capturedPiece)
}

// MakeTimedMove does the same thing as MakeMove but also adds the duration
// of the move to the player's clock. If the player goes over time, then
// the TimedOut game status is returned. In that case, the move is not
// added to the game history.
func (G *Game) MakeTimedMove(m *move.Move, timeTaken time.Duration) GameStatus {
	color := G.ActiveColor()
	G.control[color].clock -= timeTaken
	if G.control[color].clock <= 0 {
		return map[piece.Color]GameStatus{piece.White: WhiteTimedOut, piece.Black: BlackTimedOut}[color]
	}
	status := G.MakeMove(m)
	G.control[color].clock += G.control[color].Increment
	if G.control[color].movesLeft <= 0 && G.control[color].Repeating {
		G.control[color].Reset()
	}
	return status
}

// MakeMove makes the specified move on the game position. Game state information
// such as the en passant square, castling rights, 50 move rule count are also adjusted.
// The game status after the given move is made is returned.
func (G *Game) MakeMove(m *move.Move) GameStatus {
	from, to, movingPiece, capturedPiece := G.decompose(m)
	if G.illegalMove(movingPiece, m) {
		defer func() { G.Moves = append(G.Moves, *m) }()
		return G.illegalMoveStatus()
	}
	G.makeMove(m, from, to, movingPiece, capturedPiece)
	return G.Status()
}

func (G *Game) decompose(m *move.Move) (from, to square.Square, movingPiece, capturedPiece piece.Piece) {
	from, to = m.From(), m.To()
	movingPiece = G.Position.OnSquare(from)
	capturedPiece = G.Position.OnSquare(to)
	return
}

func (G *Game) makeMove(m *move.Move, from, to square.Square, movingPiece, capturedPiece piece.Piece) {
	G.Position.MakeMove(m)
	G.Moves = append(G.Moves, *m)
	G.cachePosition()
}

// Status returns the game's status.
func (G *Game) Status() GameStatus {
	activeColor := G.ActiveColor()
	check, stale := G.Check(activeColor), len(G.LegalMoves()) == 0
	if stale && check {
		return []GameStatus{WhiteCheckmated, BlackCheckmated}[activeColor]
	}
	if stale {
		return Stalemate
	}
	if G.threeFold() {
		return Threefold
	}
	if G.Position.FiftyMoveCount >= 100 { // we keep track of it in half moves, start at 0
		return FiftyMoveRule
	}
	if G.Position.InsufficientMaterial() {
		return InsufficientMaterial
	}
	return InProgress
}

func (G *Game) illegalMove(p piece.Piece, m *move.Move) bool {
	if p.Color == piece.Neither || p.Type == piece.None {
		return true
	}
	if _, legal := G.LegalMoves()[*m]; !legal {
		return true
	}
	return false
}

func (G *Game) illegalMoveStatus() GameStatus {
	if G.ActiveColor() == piece.White {
		return WhiteIllegalMove
	}
	return BlackIllegalMove
}

// TODO(andrewbackes): threeFold detection should not have to go through all of the move history.
// BUG(andrewbackes): starting FEN is not considered when calculating threefold.
func (G *Game) threeFold() bool {
	hash := polyglot.Encode(G.Position)
	if G.positionCache[hash] >= 3 {
		return true
	}
	return false
}

func (G *Game) cachePosition() {
	hash := polyglot.Encode(G.Position)
	c := G.positionCache[hash]
	c++
	G.positionCache[hash] = c
}

// Check returns whether or not the specified color is in check.
func (G *Game) Check(color piece.Color) bool {
	return G.Position.Check(color)
}

// EnPassant returns a pointer to a square or nil if there is not
// en passant square.
func (G *Game) EnPassant() square.Square {
	return G.Position.EnPassant
}

// LegalMoves returns only the legal moves that can be made.
func (G *Game) LegalMoves() map[move.Move]struct{} {
	return G.Position.LegalMoves()
}

// String puts the Board into a pretty print-able format.
func (G Game) String() string {
	castles := [][]string{{"K", "Q"}, {"k", "q"}}
	rights := ""
	for c := piece.White; c <= piece.Black; c++ {
		for s := position.ShortSide; s <= position.LongSide; s++ {
			if G.Position.CastlingRights[c][s] {
				rights += castles[c][s]
			}
		}
	}
	enpass := "None"
	if G.Position.EnPassant <= square.LastSquare {
		enpass = fmt.Sprint(G.Position.EnPassant)
	}
	str := "   +---+---+---+---+---+---+---+---+\n"
	for i := 63; i >= 0; i-- {
		p := G.Position.OnSquare(square.Square(i))
		if i%8 == 7 {
			str += fmt.Sprint(" ", (i/8)+1, " ")
		}
		str += "| " + fmt.Sprint(p) + " "
		if i%8 == 0 {
			str += "|"
			switch i / 8 {
			case 7:
				str += fmt.Sprint("   Active Color:    ", []string{"White", "Black"}[G.ActiveColor()])
			case 5:
				str += fmt.Sprint("   En Passant:      ", enpass)
			case 4:
				str += fmt.Sprint("   Castling Rights: ", rights)
			case 3:
				str += fmt.Sprint("   50 Move Rule:    ", G.Position.FiftyMoveCount)
			case 1:
				str += fmt.Sprint("   White's Clock:   ", G.control[0].clock, " (", G.control[0].movesLeft, " moves)")
			case 0:
				str += fmt.Sprint("   Black's Clock:   ", G.control[1].clock, " (", G.control[1].movesLeft, " moves)")
			}
			str += "\n   +---+---+---+---+---+---+---+---+\n"

		}
	}
	str += "     A   B   C   D   E   F   G   H\n"
	return str
}

// Clock returns the time left for the current player. It does not update
// until after a move is made.
func (G *Game) Clock(player piece.Color) time.Duration {
	return G.control[player].clock
}

// MovesLeft returns the number of moves left until time control.
func (G *Game) MovesLeft(player piece.Color) int64 {
	return G.control[player].movesLeft
}

// Result returns a human readable
func (G *Game) Result() string {
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
