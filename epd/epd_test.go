package chess

import (
	"testing"
)

func TestParseEPD(t *testing.T) {
	test := "1k1r4/pp1b1R2/3q2pp/4p3/2B5/4Q3/PPP2B2/2K5 b - - bm Qd1+; id \"BK.01\";"
	epd, err := Decode(test)
	t.Log(epd)
	if err != nil {
		t.Fail()
	}
	if epd.Operations[0].Code != "bm" ||
		epd.Operations[0].Operand != "Qd1+" ||
		epd.Operations[1].Code != "id" ||
		epd.Operations[1].Operand != "\"BK.01\"" {
		t.Fail()
	}
}
