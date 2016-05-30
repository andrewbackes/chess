package chess

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
