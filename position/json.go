package position

import (
	"encoding/json"
)

// MarshalJSON turns a position into a JSON byte slice.
func (p *Position) MarshalJSON() ([]byte, error) {
	type Alias Position
	return json.Marshal(&struct {
		Board string
		*Alias
	}{
		Board: BitBoards(p.bitBoard).MailBox(),
		Alias: (*Alias)(p),
	})
}
