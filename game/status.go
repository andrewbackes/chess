package game

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

func (g GameStatus) String() string {
	return map[GameStatus]string{
		InProgress:           "In progress",
		BlackCheckmated:      "White checkmated Black",
		WhiteCheckmated:      "Black checkmated White",
		BlackTimedOut:        "Black ran out of time",
		WhiteTimedOut:        "White ran out of time",
		BlackResigned:        "Black resigned",
		WhiteResigned:        "White resigned",
		BlackIllegalMove:     "Black made an illegal move",
		WhiteIllegalMove:     "White made an illegal move",
		Threefold:            "Draw by threefold repetition",
		FiftyMoveRule:        "Draw by Fifty move rule",
		Stalemate:            "Draw by stalemate",
		InsufficientMaterial: "Draw by insufficient material",
		WhiteWon:             "White won",
		BlackWon:             "Black won",
		Draw:                 "Draw",
	}[g]
}

const (
	WhiteWon GameStatus = (BlackCheckmated | BlackTimedOut | BlackResigned | BlackIllegalMove)
	BlackWon GameStatus = (WhiteCheckmated | WhiteTimedOut | WhiteResigned | WhiteIllegalMove)
	Draw     GameStatus = (Threefold | FiftyMoveRule | Stalemate | InsufficientMaterial)
)
