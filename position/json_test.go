package position

import (
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		position testPosition
		want     string
	}{
		{"Empty", testPosition{}, `{"Board":"                                                                ","moveNumber":1,"enPassant":64,"castlingRights":{"0":{"0":true,"1":true},"1":{"0":true,"1":true}},"activeColor":"White","movesLeft":{},"clock":{},"lastMove":{"source":64,"destination":64}}`},
		{"Initial", nil, `{"Board":"RNBKQBNRPPPPPPPP                                pppppppprnbkqbnr","moveNumber":1,"enPassant":64,"castlingRights":{"0":{"0":true,"1":true},"1":{"0":true,"1":true}},"activeColor":"White","movesLeft":{},"clock":{},"lastMove":{"source":64,"destination":64}}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.position.Position()
			jsonBytes, err := p.MarshalJSON()
			json := string(jsonBytes)
			if json != tc.want || err != nil {
				t.Logf("Position:\n%v", p)
				t.Errorf("*Position.MarshalJSON() =\n\t([]byte(\"%s\"), %v)\n, want\n\t([]byte(\"%s\"), %v)", json, err, tc.want, error(nil))
			}
		})
	}
}
