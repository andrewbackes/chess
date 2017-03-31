package board

type Side uint

const (
	ShortSide, kingSide = Side(0), Side(0)
	LongSide, queenSide = Side(1), Side(1)
)

// Sides is used to range through the sides of the board.
var Sides = [2]Side{ShortSide, LongSide}
