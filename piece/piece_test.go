package piece

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPiecePrint(t *testing.T) {
	p := New(White, Pawn)
	if fmt.Sprint(p) != "P" {
		t.Fail()
	}
}

func TestColorUnmarshalJson(t *testing.T) {
	blob := `["White"]`
	var c []Color
	json.Unmarshal([]byte(blob), &c)
	assert.Equal(t, []Color{White}, c)
}

func TestColorMarshalJson(t *testing.T) {
	j := []Color{White}
	result, err := json.Marshal(j)
	assert.Equal(t, nil, err)
	assert.Equal(t, `["White"]`, string(result))
}
