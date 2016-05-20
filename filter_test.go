package chess

import (
	"testing"
)

func TestFilter(t *testing.T) {
	pgns := []*PGN{
		&PGN{Tags: map[string]string{"Event": "keep"}},
		&PGN{Tags: map[string]string{"Event": "remove"}},
		&PGN{Tags: map[string]string{"Event": "keep"}},
	}
	expected := pgns[1]
	filtered := FilterPGNs(pgns, NewTagFilter("Event==remove"))
	if len(filtered) != 1 || filtered[0] != expected {
		t.Log(filtered)
		t.Fail()
	}
}
