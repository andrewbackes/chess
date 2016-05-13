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
	"fmt"
	"strings"
	"time"
)

type Game struct {
	tags    map[string]string
	control [2]TimeControl
	board   Board
	history GameHistory
	status  GameStatus
}

type GameHistory struct {
	fen            []string
	move           []Move
	fiftyMoveCount uint64
	enPassant      *Square
	castlingRights [2][2]bool
}

// NewGame returns a fresh game with all of the pieces in the
// opening position.
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

// NewTimedGame does the same thing as NewGame() but sets the
// time control to what is specified.
func NewTimedGame(control [2]TimeControl) *Game {
	g := NewGame()
	g.control = control
	g.control[White].Reset()
	g.control[Black].Reset()
	return g
}

// PlayerToMove returns the color of the player whos turn it is.
func (G *Game) PlayerToMove() Color {
	return Color(len(G.history.move) % 2)
}

// QuickMove makes the specified move without checking the legality
// of the move or the status of the game post move.
func (G *Game) QuickMove(m Move) {
	from, to := getSquares(m)
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
func (G *Game) MakeTimedMove(m Move, timeTaken time.Duration) GameStatus {
	color := G.PlayerToMove()
	G.control[color].clock -= timeTaken
	if G.control[color].clock <= 0 {
		return map[Color]GameStatus{White: WhiteTimedOut, Black: BlackTimedOut}[color]
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
func (G *Game) MakeMove(m Move) GameStatus {
	from, to := getSquares(m)
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
	activeColor := G.PlayerToMove()
	check, stale := G.isInCheck(activeColor), len(G.LegalMoves()) == 0
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
	if G.insufficientMaterial() {
		return InsufficientMaterial
	}
	return InProgress
}

func (G *Game) illegalMove(p Piece, m Move) bool {
	if p.Color == Neither || p.Type == None {
		return true
	}
	if _, legal := G.LegalMoves()[m]; !legal {
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

func (G *Game) adjustMoveCounter(movingPiece, capturedPiece Piece) {
	if capturedPiece.Type != None || movingPiece.Type == Pawn {
		G.history.fiftyMoveCount = 0
	} else {
		G.history.fiftyMoveCount++
	}
}

func (G *Game) adjustEnPassant(movingPiece Piece, from, to Square) {
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

func (G *Game) adjustCastlingRights(movingPiece Piece, from, to Square) {
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

func (G *Game) insufficientMaterial() bool {
	/*
		BUG!
		TODO:
		  	-(Any number of additional bishops of either color on the same color of square due to underpromotion do not affect the situation.)
	*/

	loneKing := []bool{
		G.board.occupied(White)&G.board.bitBoard[White][King] == G.board.occupied(White),
		G.board.occupied(Black)&G.board.bitBoard[Black][King] == G.board.occupied(Black)}

	if !loneKing[White] && !loneKing[Black] {
		return false
	}

	for color := White; color <= Black; color++ {
		otherColor := []Color{Black, White}[color]
		if loneKing[color] {
			// King vs King:
			if loneKing[otherColor] {
				return true
			}
			// King vs King & Knight
			if popcount(G.board.bitBoard[otherColor][Knight]) == 1 {
				mask := G.board.bitBoard[otherColor][King] | G.board.bitBoard[otherColor][Knight]
				occuppied := G.board.occupied(otherColor)
				if occuppied&mask == occuppied {
					return true
				}
			}
			// King vs King & Bishop
			if popcount(G.board.bitBoard[otherColor][Bishop]) == 1 {
				mask := G.board.bitBoard[otherColor][King] | G.board.bitBoard[otherColor][Bishop]
				occuppied := G.board.occupied(otherColor)
				if occuppied&mask == occuppied {
					return true
				}
			}
		}
		// King vs King & oppoSite bishop
		kingBishopMask := G.board.bitBoard[color][King] | G.board.bitBoard[color][Bishop]
		if (G.board.occupied(color)&kingBishopMask == G.board.occupied(color)) && (popcount(G.board.bitBoard[color][Bishop]) == 1) {
			mask := G.board.bitBoard[otherColor][King] | G.board.bitBoard[otherColor][Bishop]
			occuppied := G.board.occupied(otherColor)
			if (occuppied&mask == occuppied) && (popcount(G.board.bitBoard[otherColor][Bishop]) == 1) {
				color1 := bitscan(G.board.bitBoard[color][Bishop]) % 2
				color2 := bitscan(G.board.bitBoard[otherColor][Bishop]) % 2
				if color1 == color2 {
					return true
				}
			}
		}

	}
	return false
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
