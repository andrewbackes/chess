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

// Colors can be used to loop through the colors via range.
var Colors = [2]Color{White, Black}

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
