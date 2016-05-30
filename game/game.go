package chess

import (
	"fmt"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position"
	"strings"
	"time"
)

// Game represents a chess game.
type Game struct {
	Tags          map[string]string
	control       [2]TimeControl
	Position      position.Position
	Moves         []position.Move
	positionCache map[position.Simple]int
}

// NewGame returns a fresh game with all of the pieces in the
// opening position.
func NewGame() *Game {
	return &Game{
		control:       [2]TimeControl{},
		Tags:          make(map[string]string),
		Position:      position.New(),
		positionCache: make(map[position.Simple]int),
	}
}

// NewTimedGame does the same thing as NewGame() but sets the
// time control to what is specified.
func NewTimedGame(control [2]TimeControl) *Game {
	g := NewGame()
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
func (G *Game) QuickMove(m board.Move) {
	from, to, movingPiece, capturedPiece = G.decompose(m)
	G.makeMove(m, from, to, movingPiece, capturedPiece)
}

// MakeTimedMove does the same thing as MakeMove but also adds the duration
// of the move to the player's clock. If the player goes over time, then
// the TimedOut game status is returned. In that case, the move is not
// added to the game history.
func (G *Game) MakeTimedMove(m board.Move, timeTaken time.Duration) GameStatus {
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

// MakeMove makes the specified move on the game board. Game state information
// such as the en passant square, castling rights, 50 move rule count are also adjusted.
// The game status after the given move is made is returned.
func (G *Game) MakeMove(m board.Move) GameStatus {
	from, to, movingPiece, capturedPiece = G.decompose(m)
	if G.illegalMove(movingPiece, m) {
		defer func() { G.Moves = append(G.Moves, m) }()
		return G.illegalMoveStatus()
	}
	G.makeMove(m, from, to, movingPiece, capturedPiece)
	return G.Status()
}

func (G *Game) decompose(m board.Move) (from, to position.Move, movingPiece, capturedPiece position.Move) {
	from, to = board.Split(m)
	movingPiece = G.board.OnSquare(from)
	capturedPiece = G.board.OnSquare(to)
	return
}

func (G *Game) makeMove(m position.Move, from, to position.Move, movingPiece, capturedPiece position.Move) {
	G.adjustMoveCounter(movingPiece, capturedPiece)
	G.adjustCastlingRights(movingPiece, from, to)
	G.adjustEnPassant(movingPiece, from, to)
	G.board.MakeMove(m)
	G.Moves = append(G.Moves, m)
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
	if G.Position.FfiftyMoveCount >= 100 { // we keep track of it in half moves, start at 0
		return FiftyMoveRule
	}
	if G.board.InsufficientMaterial() {
		return InsufficientMaterial
	}
	return InProgress
}

func (G *Game) illegalMove(p piece.Piece, m board.Move) bool {
	if p.Color == piece.Neither || p.Type == piece.None {
		return true
	}
	if _, legal := G.LegalMoves()[m]; !legal {
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

func (G *Game) adjustMoveCounter(movingPiece, capturedPiece piece.Piece) {
	if capturedPiece.Type != piece.None || movingPiece.Type == piece.Pawn {
		G.Position.FfiftyMoveCount = 0
	} else {
		G.Position.FfiftyMoveCount++
	}
}

func (G *Game) adjustEnPassant(movingPiece piece.Piece, from, to board.Square) {
	if movingPiece.Type == piece.Pawn {
		G.Position.EnPassant = nil
		if int(from)-int(to) == 16 || int(from)-int(to) == -16 {
			s := board.Square(int(from) + []int{8, -8}[movingPiece.Color])
			G.Position.EnPassant = &s
		}
	} else {
		G.Position.EnPassant = nil
	}
}

func (G *Game) adjustCastlingRights(movingPiece piece.Piece, from, to board.Square) {
	for side := board.ShortSide; side <= board.LongSide; side++ {
		if movingPiece.Type == piece.King || //King moves
			(movingPiece.Type == piece.Rook &&
				from == [2][2]board.Square{{board.H1, board.A1}, {board.H8, board.A8}}[movingPiece.Color][side]) {
			G.Position.CastlingRights[movingPiece.Color][side] = false
		}
		if to == [2][2]board.Square{{board.H8, board.A8}, {board.H1, board.A1}}[movingPiece.Color][side] {
			G.Position.CastlingRights[[]piece.Color{piece.Black, piece.White}[movingPiece.Color]][side] = false
		}
	}
}

// TODO(andrewbackes): threeFold detection should not have to go through all of the move history.
// BUG(andrewbackes): starting FEN is not considered when calculating threefold.
func (G *Game) threeFold() bool {
	if G.positionCache[position.Simplify(G.Position)] >= 3 {
		return true
	}
	return false
}

func (G *Game) cachePosition() {
	c := G.positionCache[position.Simplify(G.Position)]
	c++
	G.positionCache[position.Simplify(G.Position)] = c
}

// Check returns whether or not the specified color is in check.
func (G *Game) Check(color piece.Color) bool {
	return G.Postition.Check(color)
}

// EnPassant returns a pointer to a square or nil if there is not
// en passant square.
func (G *Game) EnPassant() board.Square {
	return G.Position.EnPassant
}

// LegalMoves returns only the legal moves that can be made.
func (G *Game) LegalMoves() map[board.Move]struct{} {
	toMove := G.ActiveColor()
	return G.board.LegalMoves(toMove, G.Position.EnPassant, G.Position.CastlingRights)
}

// String puts the Board into a pretty print-able format.
func (G Game) String() string {
	castles := [][]string{{"K", "Q"}, {"k", "q"}}
	rights := ""
	for c := piece.White; c <= piece.Black; c++ {
		for s := board.ShortSide; s <= board.LongSide; s++ {
			if G.Position.CastlingRights[c][s] {
				rights += castles[c][s]
			}
		}
	}
	str := "   +---+---+---+---+---+---+---+---+\n"
	for i := 63; i >= 0; i-- {
		p := G.board.OnSquare(board.Square(i))
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
				str += fmt.Sprint("   En Passant:      ", G.Position.EnPassant)
			case 4:
				str += fmt.Sprint("   Castling Rights: ", rights)
			case 3:
				str += fmt.Sprint("   50 Move Rule:    ", G.Position.FfiftyMoveCount)
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
