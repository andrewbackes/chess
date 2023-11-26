// Package position is for working with chess positions. It holds the state of
// a chess game at a particular move.
package position

import (
	"strings"
	"time"

	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/board"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
)

// Position represents the state of a game during a player's turn.
type Position struct {
	// bitBoard has one bitBoard per player per color.
	// TODO(andrewbackes): When bitboard was changed to a map vs [2][6]uint64 array,
	// the time to run unit tests increased from around 1 minute to around 7 minutes.
	// It should probably be changed back. The change was made because piece.Type None
	// was moved to the from the the None, Pawn, ..., King iota for the const piece.Type.
	// When Pawn was first it was very easy to use arrays, but with None first, it messed
	// everything up. To fix it, just make a bitBoards struct that is backed by the array
	// with the appropriate getters and setters. Also remember, if that doesn't fix it,
	// then it has something to do with move.Move being changed from a string to a
	// struct with many more members. Check the v1.0 tag for the well performing unit
	// tests.
	bitBoard       map[piece.Color]map[piece.Type]uint64
	MoveNumber     int    `json:"moveNumber" bson:"moveNumber"`
	FiftyMoveCount uint64 `json:"fiftyMoveCount,omitempty" bson:"fiftyMoveCount,omitempty"`
	// ThreeFoldCount keeps track of how many times a certain position has been seen in the game so far.
	ThreeFoldCount map[Hash]int                        `json:"threeFoldCount,omitempty" bson:"threeFoldCount,omitempty"`
	EnPassant      square.Square                       `json:"enPassant,omitempty" bson:"enPassant,omitempty"`
	CastlingRights map[piece.Color]map[board.Side]bool `json:"castlingRights" bson:"castlingRights"`
	ActiveColor    piece.Color                         `json:"activeColor" bson:"activeColor"`
	// MovesLeft in the time control.
	MovesLeft map[piece.Color]int           `json:"movesLeft" bson:"movesLeft"`
	Clocks    map[piece.Color]time.Duration `json:"clock" bson:"clock"`
	LastMove  move.Move                     `json:"lastMove"`
}

func (p *Position) MailBox() string {
	return BitBoards(p.bitBoard).MailBox()
}

// NewCastlingRights returns castling rights set to their default
// settings.
func NewCastlingRights() map[piece.Color]map[board.Side]bool {
	return map[piece.Color]map[board.Side]bool{
		piece.White: {board.ShortSide: true, board.LongSide: true},
		piece.Black: {board.ShortSide: true, board.LongSide: true},
	}
}

// New returns a game board in the opening position. If you want
// a blank board, use Clear().
func New() *Position {
	p := &Position{
		MoveNumber:     1,
		bitBoard:       newBitboards(),
		ActiveColor:    piece.White,
		EnPassant:      square.NoSquare,
		CastlingRights: NewCastlingRights(),
		FiftyMoveCount: 0,
		ThreeFoldCount: make(map[Hash]int),
		MovesLeft:      make(map[piece.Color]int),
		Clocks:         make(map[piece.Color]time.Duration),
		LastMove:       move.Null,
	}
	p.Reset()
	return p
}

// Copy makes an exact copy of the position.
func Copy(p *Position) *Position {
	n := &Position{
		bitBoard:       newBitboards(),
		MoveNumber:     p.MoveNumber,
		ActiveColor:    p.ActiveColor,
		EnPassant:      p.EnPassant,
		FiftyMoveCount: p.FiftyMoveCount,
		CastlingRights: make(map[piece.Color]map[board.Side]bool),
		ThreeFoldCount: make(map[Hash]int),
		MovesLeft:      make(map[piece.Color]int),
		Clocks:         make(map[piece.Color]time.Duration),
		LastMove:       p.LastMove,
	}
	for k, v := range p.ThreeFoldCount {
		n.ThreeFoldCount[k] = v
	}
	for _, color := range piece.Colors {
		for k, v := range p.bitBoard[color] {
			n.bitBoard[color][k] = v
		}
		n.CastlingRights[color] = make(map[board.Side]bool)
		for _, side := range board.Sides {
			n.CastlingRights[color][side] = p.CastlingRights[color][side]
		}
		n.MovesLeft[color] = p.MovesLeft[color]
		n.Clocks[color] = p.Clocks[color]
	}
	return n
}

// Equals compares two positions for the same active player, position of the
// pieces, castling rights and en passant. It does not compare how much time
// is left on the clock or how many moves are left in the time control.
func (p *Position) Equals(q *Position) bool {
	if p.ActiveColor != q.ActiveColor {
		return false
	}
	if p.EnPassant != q.EnPassant {
		return false
	}
	for _, color := range piece.Colors {
		for k := range p.bitBoard[color] {
			if p.bitBoard[color][k] != q.bitBoard[color][k] {
				return false
			}
		}
		for _, side := range board.Sides {
			if p.CastlingRights[color][side] != q.CastlingRights[color][side] {
				return false
			}
		}
	}
	return true
}

func (p *Position) GetCastlingRights() map[piece.Color]map[board.Side]bool {
	return p.CastlingRights
}

func (p *Position) GetActiveColor() piece.Color {
	return p.ActiveColor
}

func (p *Position) GetEnPassant() square.Square {
	return p.EnPassant
}

func (p *Position) GetFiftyMoveCount() uint64 {
	return p.FiftyMoveCount
}

func (p *Position) GetMoveNumber() int {
	return p.MoveNumber
}

/*
// String puts the Board into a pretty print-able format.
func (p Position) String() (str string) {
	str += "+---+---+---+---+---+---+---+---+\n"
	for i := 1; i <= 64; i++ {
		sq := square.Square(64 - i)
		str += "|"
		noPiece := true
		for c := range p.bitBoard {
			for j := range p.bitBoard[c] {
				if ((1 << sq) & p.bitBoard[c][j]) != 0 {
					str += fmt.Sprint(" ", p.OnSquare(sq), " ")
					noPiece = false
				}
			}
		}
		if noPiece {
			str += "   "
		}
		if sq%8 == 0 {
			str += "|\n"
			str += "+---+---+---+---+---+---+---+---+"
			if sq < square.LastSquare {
				str += "\n"
			}
		}
	}
	return
}
*/
// Clear empties the Board.
func (p *Position) Clear() {
	p.bitBoard = newBitboards()
}

// Reset puts the pieces in the new game position.
func (p *Position) Reset() {
	// puts the pieces in their starting/newgame positions
	for color := piece.Color(0); color < 2; color = color + 1 {
		//Pawns first:
		p.bitBoard[color][piece.Pawn] = 255 << (8 + (color * 8 * 5))
		//Then the rest of the pieces:
		p.bitBoard[color][piece.Knight] = (1 << (square.B1 + square.Square(color*8*7))) ^ (1 << (square.G1 + square.Square(color*8*7)))
		p.bitBoard[color][piece.Bishop] = (1 << (square.C1 + square.Square(color*8*7))) ^ (1 << (square.F1 + square.Square(color*8*7)))
		p.bitBoard[color][piece.Rook] = (1 << (square.A1 + square.Square(color*8*7))) ^ (1 << (square.H1 + square.Square(color*8*7)))
		p.bitBoard[color][piece.Queen] = (1 << (square.D1 + square.Square(color*8*7)))
		p.bitBoard[color][piece.King] = (1 << (square.E1 + square.Square(color*8*7)))
	}
}

// OnSquare returns the piece that is on the specified square.
func (p *Position) OnSquare(s square.Square) piece.Piece {
	return BitBoards(p.bitBoard).OnSquare(s)
}

// Occupied returns a bitBoard with all of the specified colors pieces.
func (p *Position) occupied(c piece.Color) uint64 {
	var mask uint64
	for pc := piece.Pawn; pc <= piece.King; pc++ {
		if c == piece.BothColors {
			mask |= p.bitBoard[piece.White][pc] | p.bitBoard[piece.Black][pc]
		} else {
			mask |= p.bitBoard[c][pc]
		}
	}
	return mask
}

func (p *Position) decompose(m move.Move) (from, to square.Square, movingPiece, capturedPiece piece.Piece) {
	return m.From(), m.To(), p.OnSquare(m.From()), p.OnSquare(m.To())
}

// MakeMove attempts to make the given move no matter legality or validity.
// It does not change game state such as en passant or castling rights.
// What ever move you specify will attempt to be made. If it is illegal
// or invalid you will get undetermined behavior.
func (p *Position) MakeMove(m move.Move) *Position {
	q := Copy(p)
	from, to, movingPiece, capturedPiece := q.decompose(m)
	q.adjustMoveCounter(movingPiece, capturedPiece)
	q.adjustCastlingRights(movingPiece, from, to)
	q.adjustEnPassant(movingPiece, from, to)
	q.adjustBoard(m, from, to, movingPiece, capturedPiece)
	q.Clocks[q.ActiveColor] -= m.Duration
	q.MovesLeft[q.ActiveColor]--
	q.ActiveColor = (q.ActiveColor + 1) % 2
	if q.ActiveColor == piece.White {
		q.MoveNumber++
	}
	q.LastMove = m
	q.adjustThreeFoldCounter()
	return q
}

func (p *Position) adjustThreeFoldCounter() {
	hash := p.Polyglot()
	if p.FiftyMoveCount == 0 {
		p.ThreeFoldCount = make(map[Hash]int)
	}
	if count, exists := p.ThreeFoldCount[hash]; exists {
		count++
		p.ThreeFoldCount[hash] = count
	} else {
		p.ThreeFoldCount[hash] = 1
	}
}

func (p *Position) adjustMoveCounter(movingPiece, capturedPiece piece.Piece) {
	if capturedPiece.Type != piece.None || movingPiece.Type == piece.Pawn {
		p.FiftyMoveCount = 0
	} else {
		p.FiftyMoveCount++
	}
}

func (p *Position) adjustEnPassant(movingPiece piece.Piece, from, to square.Square) {
	if movingPiece.Type == piece.Pawn {
		p.EnPassant = square.NoSquare
		if int(from)-int(to) == 16 || int(from)-int(to) == -16 {
			s := square.Square(int(from) + []int{8, -8}[movingPiece.Color])
			p.EnPassant = s
		}
	} else {
		p.EnPassant = square.NoSquare
	}
}

func (p *Position) adjustCastlingRights(movingPiece piece.Piece, from, to square.Square) {
	for side := board.ShortSide; side <= board.LongSide; side++ {
		if movingPiece.Type == piece.King || //King moves
			(movingPiece.Type == piece.Rook &&
				from == [2][2]square.Square{{square.H1, square.A1}, {square.H8, square.A8}}[movingPiece.Color][side]) {
			p.CastlingRights[movingPiece.Color][side] = false
		}
		if to == [2][2]square.Square{{square.H8, square.A8}, {square.H1, square.A1}}[movingPiece.Color][side] {
			p.CastlingRights[[]piece.Color{piece.Black, piece.White}[movingPiece.Color]][side] = false
		}
	}
}

func (p *Position) adjustBoard(m move.Move, from, to square.Square, movingPiece, capturedPiece piece.Piece) {
	// Remove captured piece:
	if capturedPiece.Type != piece.None {
		p.bitBoard[capturedPiece.Color][capturedPiece.Type] ^= (1 << to)
	}

	// Move piece:
	p.bitBoard[movingPiece.Color][movingPiece.Type] ^= ((1 << from) | (1 << to))

	// Castle:
	if movingPiece.Type == piece.King {
		if from == (square.E1+square.Square(56*uint8(movingPiece.Color))) && (to == square.G1+square.Square(56*uint8(movingPiece.Color))) {
			p.bitBoard[movingPiece.Color][piece.Rook] ^= (1 << (square.H1 + square.Square(56*movingPiece.Color))) | (1 << (square.F1 + square.Square(56*movingPiece.Color)))
		} else if from == square.E1+square.Square(56*uint8(movingPiece.Color)) && (to == square.C1+square.Square(56*uint8(movingPiece.Color))) {
			p.bitBoard[movingPiece.Color][piece.Rook] ^= (1 << (square.A1 + square.Square(56*(movingPiece.Color)))) | (1 << (square.D1 + square.Square(56*(movingPiece.Color))))
		}
	}

	if movingPiece.Type == piece.Pawn {
		// Handle en Passant capture:
		// capturedPiece just means the piece on the destination square
		if (int(to)-int(from))%8 != 0 && capturedPiece.Type == piece.None {
			if movingPiece.Color == piece.White {
				p.bitBoard[piece.Black][piece.Pawn] ^= (1 << (to - 8))
			} else if movingPiece.Color == piece.Black {
				p.bitBoard[piece.White][piece.Pawn] ^= (1 << (to + 8))
			}
		}
		// Handle Promotions:
		if m.Promote != piece.None {
			p.bitBoard[movingPiece.Color][movingPiece.Type] ^= (1 << to) // remove piece.Pawn
			p.bitBoard[movingPiece.Color][m.Promote] ^= (1 << to)        // add promoted piece
		}
	}
}

// Put places a piece on the square and removes any other piece
// that may be on that square.
func (p *Position) Put(pp piece.Piece, s square.Square) {
	pc := p.OnSquare(s)
	if pc.Type != piece.None {
		p.bitBoard[pc.Color][pc.Type] ^= (1 << s)
	}
	p.bitBoard[pp.Color][pp.Type] |= (1 << s)
}

// QuickPut places a piece on the square without removing
// any piece that may already be on that square.
func (p *Position) QuickPut(pc piece.Piece, s square.Square) {
	p.bitBoard[pc.Color][pc.Type] |= (1 << s)
}

// Find returns the squares that hold the specified piece.
func (p *Position) Find(pc piece.Piece) map[square.Square]struct{} {
	s := make(map[square.Square]struct{})
	bits := p.bitBoard[pc.Color][pc.Type]
	for bits != 0 {
		sq := bitscan(bits)
		s[square.Square(sq)] = struct{}{}
		bits ^= (1 << sq)
	}
	return s
}

func (p *Position) InsufficientMaterial() bool {
	/*
		BUG!
		TODO:
		  	-(Any number of additional bishops of either color on the same color of square due to underpromotion do not affect the situation.)
	*/
	loneKing := []bool{
		p.occupied(piece.White)&p.bitBoard[piece.White][piece.King] == p.occupied(piece.White),
		p.occupied(piece.Black)&p.bitBoard[piece.Black][piece.King] == p.occupied(piece.Black)}

	if !loneKing[piece.White] && !loneKing[piece.Black] {
		return false
	}

	for color := piece.White; color <= piece.Black; color++ {
		otherColor := []piece.Color{piece.Black, piece.White}[color]
		if loneKing[color] {
			// King vs King:
			if loneKing[otherColor] {
				return true
			}
			// King vs King & Knight
			if popcount(p.bitBoard[otherColor][piece.Knight]) == 1 {
				mask := p.bitBoard[otherColor][piece.King] | p.bitBoard[otherColor][piece.Knight]
				occuppied := p.occupied(otherColor)
				if occuppied&mask == occuppied {
					return true
				}
			}
			// King vs King & Bishop
			if popcount(p.bitBoard[otherColor][piece.Bishop]) == 1 {
				mask := p.bitBoard[otherColor][piece.King] | p.bitBoard[otherColor][piece.Bishop]
				occuppied := p.occupied(otherColor)
				if occuppied&mask == occuppied {
					return true
				}
			}
		}
		// King vs King & oppoSite bishop
		kingBishopMask := p.bitBoard[color][piece.King] | p.bitBoard[color][piece.Bishop]
		if (p.occupied(color)&kingBishopMask == p.occupied(color)) && (popcount(p.bitBoard[color][piece.Bishop]) == 1) {
			mask := p.bitBoard[otherColor][piece.King] | p.bitBoard[otherColor][piece.Bishop]
			occuppied := p.occupied(otherColor)
			if (occuppied&mask == occuppied) && (popcount(p.bitBoard[otherColor][piece.Bishop]) == 1) {
				color1 := bitscan(p.bitBoard[color][piece.Bishop]) % 2
				color2 := bitscan(p.bitBoard[otherColor][piece.Bishop]) % 2
				if color1 == color2 {
					return true
				}
			}
		}

	}
	return false
}

// Check returns whether or not the specified color is in check.
func (p *Position) Check(color piece.Color) bool {
	opponent := []piece.Color{piece.Black, piece.White}[color]
	kingsq := square.Square(bitscan(p.bitBoard[color][piece.King]))
	return p.Threatened(kingsq, opponent)
}

// Returns a string that is a SAN representation of the input move from the current position.
// If the input move is not a legal move, empty string is returned.
func (p Position) SAN(m move.Move) string {
	legalMoves := p.LegalMoves()
	if _, exists := legalMoves[m]; !exists {
		return ""
	}

	movingPiece := p.OnSquare(m.From())

	// Castling.
	if movingPiece.Type == piece.King {
		fromFile, toFile := int(m.From())%8, int(m.To())%8
		if fromFile == int(square.E1)%8 {
			if toFile == int(square.G1)%8 {
				return "O-O"
			}
			if toFile == int(square.C1)%8 {
				return "O-O-O"
			}
		}
	}

	mF, mR := m.Source.Algebraic()[0], m.Source.Algebraic()[1]
	otherMoves, sameRank, sameFile := false, false, false
	// Loop through all legal moves and check if there are some other moves with the same piece type to destination square from the same file or rank.
	for lm, _ := range legalMoves {
		// Check if the current legal move is from the same source square as input move.
		if lm.Source == m.Source {
			// It is the same source square, do not consider it as the other move we are searching for.
			continue
		}
		// Current legal move is from some other source square as the one in input.

		// Check if the move is with the same piece type.
		if p.OnSquare(lm.Source).Type != movingPiece.Type {
			// Piece type differs, this is not the move we are searching for.
			continue
		}
		// Current legal move is from other source square with the same piece type as the input move.

		// Check if the move is to the same destination as the input move.
		if lm.Destination != m.Destination {
			// Destination differs, do not continue.
			continue
		}
		// Current legal move is from other source square to the same destination, with the same piece type as the input move.
		otherMoves = true

		// Mark if the other move is from the same file or rank as the input move.
		lmF, lmR := lm.Source.Algebraic()[0], lm.Source.Algebraic()[1]
		if mR == lmR {
			sameRank = true
		}
		if mF == lmF {
			sameFile = true
		}

		// If there are other moves from the same file and rank, we don't need to look for another moves.
		if sameFile && sameRank {
			break
		}
	}

	// Parts of SAN for output.
	var pieceUpper, fromFile, fromRank, capture, promotion, check string

	if movingPiece.Type != piece.Pawn {
		pieceUpper = strings.ToUpper(movingPiece.Type.String())
	}

	if otherMoves {
		if !sameFile {
			fromFile = string(mF)
		} else if !sameRank {
			fromRank = string(mR)
		} else {
			fromFile = string(mF)
			fromRank = string(mR)
		}
	}

	capturedPiece := p.OnSquare(m.To())
	if p.EnPassant == m.To() {
		capturedPiece = piece.New((p.ActiveColor+1)%2, piece.Pawn)
	}
	if capturedPiece.Type != piece.None {
		capture = "x"
		if movingPiece.Type == piece.Pawn && fromFile == "" {
			fromFile = string(m.Source.Algebraic()[0])
		}
	}
	if m.Promote != piece.None {
		promotion = "=" + strings.ToUpper(m.Promote.String())
	}
	nextPosition := p.MakeMove(m)
	if nextPosition.Check(nextPosition.ActiveColor) {
		if len(nextPosition.LegalMoves()) == 0 {
			check = "#"
		} else {
			check = "+"
		}
	}

	return pieceUpper + fromFile + fromRank + capture + m.To().String() + promotion + check
}
