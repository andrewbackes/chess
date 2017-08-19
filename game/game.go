// Package game handles chess games. You can create timed and untimed games.
package game

import (
	"fmt"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position"
	"github.com/andrewbackes/chess/position/board"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
	"time"
)

// Game represents a chess game.
type Game struct {
	Tags      map[string]string
	control   map[piece.Color]TimeControl
	Positions []*position.Position
}

// New returns a fresh game with all of the pieces in the
// opening position.
func New() *Game {
	return &Game{
		control:   nil,
		Tags:      make(map[string]string),
		Positions: []*position.Position{position.New()},
	}
}

// NewTimedGame does the same thing as NewGame() but sets the
// time control to what is specified.
func NewTimedGame(control map[piece.Color]TimeControl) *Game {
	g := New()
	g.control = control
	g.Position().Clocks[piece.White] = control[piece.White].Time
	g.Position().MovesLeft[piece.White] = control[piece.White].Moves
	g.Position().Clocks[piece.Black] = control[piece.Black].Time
	g.Position().MovesLeft[piece.Black] = control[piece.Black].Moves
	return g
}

// Position returns the current position of the game.
func (G *Game) Position() *position.Position {
	return G.Positions[len(G.Positions)-1]
}

// ActiveColor returns the color of the player whos turn it is.
func (G *Game) ActiveColor() piece.Color {
	return G.Position().ActiveColor
}

// QuickMove makes the specified move without checking the legality
// of the move or the status of the game post move.
func (G *Game) QuickMove(m move.Move) {
	from, to, movingPiece, capturedPiece := G.decompose(m)
	G.makeMove(m, from, to, movingPiece, capturedPiece)
}

// MakeMove makes the specified move on the game position. Game state information
// such as the en passant square, castling rights, 50 move rule count are also adjusted.
// The game status after the given move is made is returned.
func (G *Game) MakeMove(m move.Move) (GameStatus, error) {
	from, to, movingPiece, capturedPiece := G.decompose(m)
	if G.illegalMove(movingPiece, m) {
		return G.illegalMoveStatus(), fmt.Errorf("%s illegal move %s", G.Position().ActiveColor, m)
	}
	G.makeMove(m, from, to, movingPiece, capturedPiece)
	if G.Position().Clocks[movingPiece.Color] < 0 {
		return map[piece.Color]GameStatus{piece.White: WhiteTimedOut, piece.Black: BlackTimedOut}[movingPiece.Color], nil
	}
	G.Position().Clocks[movingPiece.Color] += G.control[movingPiece.Color].Increment
	if G.Position().MovesLeft[movingPiece.Color] <= 0 && G.control[movingPiece.Color].Repeating {
		G.Position().MovesLeft[movingPiece.Color] = G.control[movingPiece.Color].Moves
	}
	return G.Status(), nil
}

func (G *Game) decompose(m move.Move) (from, to square.Square, movingPiece, capturedPiece piece.Piece) {
	from, to = m.From(), m.To()
	movingPiece = G.Position().OnSquare(from)
	capturedPiece = G.Position().OnSquare(to)
	return
}

func (G *Game) makeMove(m move.Move, from, to square.Square, movingPiece, capturedPiece piece.Piece) {
	newPos := G.Position().MakeMove(m)
	G.Positions = append(G.Positions, newPos)
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
	if G.Position().FiftyMoveCount >= 100 { // we keep track of it in half moves, start at 0
		return FiftyMoveRule
	}
	if G.Position().InsufficientMaterial() {
		return InsufficientMaterial
	}
	return InProgress
}

func (G *Game) illegalMove(p piece.Piece, m move.Move) bool {
	if p.Color == piece.Neither || p.Type == piece.None {
		return true
	}
	n := move.Move{Source: m.Source, Destination: m.Destination, Promote: m.Promote}
	if _, exists := G.LegalMoves()[n]; !exists {
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
	hash := G.Position().Polyglot()
	if G.Position().ThreeFoldCount[hash] >= 3 {
		return true
	}
	return false
}

// Check returns whether or not the specified color is in check.
func (G *Game) Check(color piece.Color) bool {
	return G.Position().Check(color)
}

// EnPassant returns a pointer to a square or nil if there is not
// en passant square.
func (G *Game) EnPassant() square.Square {
	return G.Position().EnPassant
}

// LegalMoves returns only the legal moves that can be made.
func (G *Game) LegalMoves() map[move.Move]struct{} {
	return G.Position().LegalMoves()
}

// String puts the Board into a pretty print-able format.
func (G Game) String() string {
	castles := [][]string{{"K", "Q"}, {"k", "q"}}
	rights := ""
	for c := piece.White; c <= piece.Black; c++ {
		for s := board.ShortSide; s <= board.LongSide; s++ {
			if G.Position().CastlingRights[c][s] {
				rights += castles[c][s]
			}
		}
	}
	enpass := "None"
	if G.Position().EnPassant <= square.LastSquare {
		enpass = fmt.Sprint(G.Position().EnPassant)
	}
	str := "   +---+---+---+---+---+---+---+---+\n"
	for i := 63; i >= 0; i-- {
		p := G.Position().OnSquare(square.Square(i))
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
				str += fmt.Sprint("   50 Move Rule:    ", G.Position().FiftyMoveCount)
			case 1:
				str += fmt.Sprint("   White's Clock:   ", G.Position().Clocks[piece.White], " (", G.Position().MovesLeft[piece.White], " moves)")
			case 0:
				str += fmt.Sprint("   Black's Clock:   ", G.Position().Clocks[piece.Black], " (", G.Position().MovesLeft[piece.Black], " moves)")
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
	return G.Position().Clocks[player]
}

// MovesLeft returns the number of moves left until time control.
func (G *Game) MovesLeft(player piece.Color) int {
	return G.Position().MovesLeft[player]
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
