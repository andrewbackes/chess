/*******************************************************************************

 TODO:
 	-Test the fen save/load by loading a file full of FEN and make sure the
 	same fen comes out after loading.
 	-check state before each modifying function and functions that print to
 	the screen
 	-double check 3 fold and 50 move rules.
 	-optimize 3 fold.
 	-some of the Game data members can be omitted from json.

*******************************************************************************/

package game

import (
	"time"
)

type TimeControl struct {
	Time      time.Duration
	BonusTime time.Duration
	Moves     int64
	Repeating bool
	clock     time.Duration
	movesLeft int64
}

func (t *TimeControl) Reset() {
	t.clock = t.Time
	t.movesLeft = t.Moves
}

func NewTimedGame(control [2]TimeControl) *Game {
	g := NewGame()
	g.control = control
	return g
}

func NewGame() *Game {
	return &Game{
		control: [2]TimeControl{},
		tags:    make(map[string]string),
		board:   NewBoard(),
		history: GameHistory{
			fen:            make([]string, 0),
			move:           make([]Move, 0),
			castlingRights: [2][2]bool{{true, true}, {true, true}},
		},
	}
}

type GameHistory struct {
	fen            []string
	move           []Move
	fiftyMoveCount uint64
	enPassant      *Square
	castlingRights [2][2]bool
}

type Game struct {
	tags    map[string]string
	control [2]TimeControl
	board   Board
	history GameHistory
	status  GameStatus
}

// PlayerToMove returns the color of the player whos turn it is.
func (G *Game) PlayerToMove() Color {
	return Color(len(G.history.move) % 2)
}

func (G *Game) MakeTimedMove(m Move, timeTaken time.Duration) GameStatus {
	color := G.PlayerToMove()
	G.control[color].clock += timeTaken
	status := G.MakeMove(m)
	if G.control[color].clock > G.control[color].Time {
		status = map[Color]GameStatus{White: WhiteTimedOut, Black: BlackTimedOut}[color]
	}
	return status
}

// MakeMove applies the specified move to the game board and adjusts any
// reprecusions. Possible reprecusions include en passant and castling.
// The status of the game after this move is made is returned.
func (G *Game) MakeMove(m Move) GameStatus {
	G.history.fen = append(G.history.fen, G.FEN())
	G.history.move = append(G.history.move, m)
	G.history.fiftyMoveCount++
	from, to := getSquares(m)
	movingPiece := G.board.OnSquare(from)
	if G.illegalMove(movingPiece, from, to) {
		return G.illegalMoveStatus()
	}

	capturedPiece := G.board.OnSquare(to)
	if capturedPiece.Type != None || movingPiece.Type == Pawn {
		G.history.fiftyMoveCount = 0
	}
	G.handleCastlingRights(movingPiece, from, to)
	G.handleEnPassant(movingPiece, from, to)
	G.board.MakeMove(m)
	return InProgress
}

func (G *Game) illegalMove(movingPiece Piece, from, to Square) bool {
	if movingPiece.Color == Neither || movingPiece.Type == None {
		return true
	}
	return false
}

func (G *Game) illegalMoveStatus() GameStatus {
	if G.PlayerToMove() == White {
		return WhiteIllegalMove
	}
	return BlackIllegalMove
}

func (G *Game) handleEnPassant(movingPiece Piece, from, to Square) {
	if movingPiece.Type == Pawn {
		G.history.enPassant = nil
		if int(from)-int(to) == 16 || int(from)-int(to) == -16 {
			s := Square(int(from) + []int{8, -8}[movingPiece.Color])
			G.history.enPassant = &s
		}
	} else {
		G.history.enPassant = nil
	}
}

func (G *Game) handleCastlingRights(movingPiece Piece, from, to Square) {
	for side := ShortSide; side <= LongSide; side++ {
		if movingPiece.Type == King || //King moves
			(movingPiece.Type == Rook && from == Square([2][2]uint8{{H1, A1}, {H8, A8}}[movingPiece.Color][side])) {
			G.history.castlingRights[movingPiece.Color][side] = false
		}
		if to == Square([2][2]uint8{{H8, A8}, {H1, A1}}[movingPiece.Color][side]) {
			G.history.castlingRights[[]Color{Black, White}[movingPiece.Color]][side] = false
		}
	}
}
