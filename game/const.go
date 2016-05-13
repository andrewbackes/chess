package game

// Color is the color of a piece or square.
type Color uint8

// Possible colors of pieces.
const (
	White      Color = 0
	Black      Color = 1
	Both       Color = 2
	BothColors Color = 2
	Neither    Color = 2
	NoColor    Color = 2
)

// PieceType is a player's piece.
type PieceType uint8

// Possible pieces.
const (
	Pawn PieceType = iota
	Knight
	Bishop
	Rook
	Queen
	King
	None
)
const NoPiece PieceType = None

type Move string

type GameStatus uint16

const (
	InProgress           GameStatus = 1 << iota
	BlackCheckmated                 //2
	WhiteCheckmated                 //4
	BlackTimedOut                   //8
	WhiteTimedOut                   //16
	BlackResigned                   //32
	WhiteResigned                   //64
	BlackIllegalMove                //128
	WhiteIllegalMove                //256
	Threefold                       //512
	FiftyMoveRule                   //1024
	Stalemate                       //2048
	InsufficientMaterial            //4096
)

const (
	WhiteWon GameStatus = (BlackCheckmated | BlackTimedOut | BlackResigned | BlackIllegalMove)
	BlackWon GameStatus = (WhiteCheckmated | WhiteTimedOut | WhiteResigned | WhiteIllegalMove)
	Draw     GameStatus = (Threefold | FiftyMoveRule | Stalemate | InsufficientMaterial)
)

const (
	ShortSide, KingSide uint = 0, 0
	LongSide, QueenSide uint = 1, 1
)

// Square on the board.
type Square uint8

// Squares on the board.
const (
	H1 Square = iota
	G1
	F1
	E1
	D1
	C1
	B1
	A1
	H2
	G2
	F2
	E2
	D2
	C2
	B2
	A2
	H3
	G3
	F3
	E3
	D3
	C3
	B3
	A3
	H4
	G4
	F4
	E4
	D4
	C4
	B4
	A4
	H5
	G5
	F5
	E5
	D5
	C5
	B5
	A5
	H6
	G6
	F6
	E6
	D6
	C6
	B6
	A6
	H7
	G7
	F7
	E7
	D7
	C7
	B7
	A7
	H8
	G8
	F8
	E8
	D8
	C8
	B8
	A8
)

const LastSquare Square = A8
