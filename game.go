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

package chess

import (
	"fmt"
	"github.com/andrewbackes/chess/board"
	"github.com/andrewbackes/chess/piece"
	"strings"
	"time"
)

// Game represents a chess game.
type Game struct {
	tags    map[string]string
	control [2]TimeControl
	board   board.Board
	history gameHistory
	status  GameStatus
}

type gameHistory struct {
	fen            []string
	move           []board.Move
	fiftyMoveCount uint64
	enPassant      *board.Square
	castlingRights [2][2]bool
}

// NewGame returns a fresh game with all of the pieces in the
// opening position.
func NewGame() *Game {
	return &Game{
		control: [2]TimeControl{},
		tags:    make(map[string]string),
		board:   board.New(),
		history: gameHistory{
			fen:            make([]string, 0),
			move:           make([]board.Move, 0),
			castlingRights: [2][2]bool{{true, true}, {true, true}},
		},
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

// Board returns the Board object used for the game.
func (G *Game) Board() *board.Board {
	return &G.board
}

// MoveHistory returns a slice of every move made so far in the game.
func (G *Game) MoveHistory() []board.Move {
	return G.history.move
}

// ActiveColor returns the color of the player whos turn it is.
func (G *Game) ActiveColor() piece.Color {
	return piece.Color(len(G.history.move) % 2)
}

// QuickMove makes the specified move without checking the legality
// of the move or the status of the game post move.
func (G *Game) QuickMove(m board.Move) {
	from, to := board.Split(m)
	movingPiece := G.board.OnSquare(from)
	capturedPiece := G.board.OnSquare(to)
	G.adjustMoveCounter(movingPiece, capturedPiece)
	G.adjustCastlingRights(movingPiece, from, to)
	G.adjustEnPassant(movingPiece, from, to)
	G.board.MakeMove(m)
	G.history.move = append(G.history.move, m)
	G.history.fen = append(G.history.fen, G.FEN())
}

// MakeTimedMove does the same thing as MakeMove but also adds the duration
// of the move to the player's clock. If the player goes over time, then
// the TimedOut game status is returned. In that case, the move is not
// added to the game history.
//
// TODO(andrewbackes): add bonus time
// TODO(andrewbackes): reset clock if move limit reached
func (G *Game) MakeTimedMove(m board.Move, timeTaken time.Duration) GameStatus {
	color := G.ActiveColor()
	G.control[color].clock -= timeTaken
	if G.control[color].clock <= 0 {
		return map[piece.Color]GameStatus{piece.White: WhiteTimedOut, piece.Black: BlackTimedOut}[color]
	}
	status := G.MakeMove(m)
	G.control[color].clock += G.control[color].Increment
	if G.control[color].movesLeft <= 0 && G.control[color].Repeating {
		fmt.Println(G.control[color].movesLeft, G.control[color].Repeating)
		G.control[color].Reset()
	}
	return status
}

// MakeMove makes the specified move on the game board. Game state information
// such as the en passant square, castling rights, 50 move rule count are also adjusted.
// The game status after the given move is made is returned.
func (G *Game) MakeMove(m board.Move) GameStatus {
	from, to := board.Split(m)
	movingPiece := G.board.OnSquare(from)
	capturedPiece := G.board.OnSquare(to)
	if G.illegalMove(movingPiece, m) {
		defer func() { G.history.move = append(G.history.move, m) }()
		return G.illegalMoveStatus()
	}
	G.adjustMoveCounter(movingPiece, capturedPiece)
	G.adjustCastlingRights(movingPiece, from, to)
	G.adjustEnPassant(movingPiece, from, to)
	G.board.MakeMove(m)
	G.history.move = append(G.history.move, m)
	G.history.fen = append(G.history.fen, G.FEN())
	return G.gameStatus()
}

func (G *Game) gameStatus() GameStatus {
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
	if G.history.fiftyMoveCount >= 100 { // we keep track of it in half moves, start at 0
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
		G.history.fiftyMoveCount = 0
	} else {
		G.history.fiftyMoveCount++
	}
}

func (G *Game) adjustEnPassant(movingPiece piece.Piece, from, to board.Square) {
	if movingPiece.Type == piece.Pawn {
		G.history.enPassant = nil
		if int(from)-int(to) == 16 || int(from)-int(to) == -16 {
			s := board.Square(int(from) + []int{8, -8}[movingPiece.Color])
			G.history.enPassant = &s
		}
	} else {
		G.history.enPassant = nil
	}
}

func (G *Game) adjustCastlingRights(movingPiece piece.Piece, from, to board.Square) {
	for side := board.ShortSide; side <= board.LongSide; side++ {
		if movingPiece.Type == piece.King || //King moves
			(movingPiece.Type == piece.Rook &&
				from == [2][2]board.Square{{board.H1, board.A1}, {board.H8, board.A8}}[movingPiece.Color][side]) {
			G.history.castlingRights[movingPiece.Color][side] = false
		}
		if to == [2][2]board.Square{{board.H8, board.A8}, {board.H1, board.A1}}[movingPiece.Color][side] {
			G.history.castlingRights[[]piece.Color{piece.Black, piece.White}[movingPiece.Color]][side] = false
		}
	}
}

// TODO(andrewbackes): threeFold detection should not have to go through all of the move history.
func (G *Game) threeFold() bool {
	for i := 0; i < len(G.history.fen); i++ {
		fen := G.history.fen[i]
		fenSplit := strings.Split(fen, " ")
		fenPrefix := fenSplit[0] + " " + fenSplit[1] + " " + fenSplit[2] + " " + fenSplit[3]
		for j := i + 1; j < len(G.history.fen); j++ {
			if strings.HasPrefix(G.history.fen[j], fenPrefix) {
				for k := j + 1; k < len(G.history.fen); k++ {
					if strings.HasPrefix(G.history.fen[k], fenPrefix) {
						return true
					}
				}
			}
		}
	}
	return false
}

// Check returns whether or not the specified color is in check.
func (G *Game) Check(color piece.Color) bool {
	return G.board.Check(color)
	/*
		opponent := []piece.Color{piece.Black, piece.White}[color]
		s := G.board.Find(piece.New(color, piece.King))
		for sq := range s {
			return G.board.Threatened(sq, opponent)
		}
		return false
	*/
}

// EnPassant returns a pointer to a square or nil if there is not
// en passant square.
func (G *Game) EnPassant() *board.Square {
	return G.history.enPassant
}

// LegalMoves returns only the legal moves that can be made.
func (G *Game) LegalMoves() map[board.Move]struct{} {
	toMove := G.ActiveColor()
	return G.board.LegalMoves(toMove, G.history.enPassant, G.history.castlingRights)
}
