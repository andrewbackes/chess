package piece

import (
	"encoding/json"
	"strings"
)

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

const COLOR_COUNT = 2

// Colors can be used to loop through the colors via range.
var Colors = [COLOR_COUNT]Color{White, Black}

// OtherColor can be used to get opponent's color. E.g. `oponentColor := piece.OtherColor[position.ActiveColor]`.
var OtherColor = [COLOR_COUNT]Color{Black, White}

func (c Color) String() string {
	return map[Color]string{
		White: "White",
		Black: "Black",
	}[c]
}

func (c Color) MarshalJSON() ([]byte, error) {
	return json.Marshal((c).String())
}

func (c *Color) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	case "white":
		*c = White
	case "black":
		*c = Black
	}
	return nil
}
