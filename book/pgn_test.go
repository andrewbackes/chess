package book

import (
	//"fmt"
	"fmt"
	"github.com/andrewbackes/chess"
	"strings"
	"testing"
)

func TestFromPGN(t *testing.T) {
	input := `[Event "one"]
[Round "1"]
[Result "1-0"]

1. e2e4 e7e5 2. d1h5 e8e7 3. h5e5 1-0

[Event "two"]
[Round "2"]
[Result "1/2-1/2"]

1. e2e4 e7e5 2. d1h5 e8e7 1/2-1/2
`
	pgns, _ := chess.ReadPGN(strings.NewReader(input))
	book, _ := FromPGN(pgns, 4)
	fmt.Println(book)
}
