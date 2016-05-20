package chess

import (
	"testing"
)

func testPGNs() []*PGN {
	pgns := []*PGN{
		&PGN{Tags: map[string]string{
			"Event": "keep", "WhiteElo": "2200", "BlackElo": "2700", "Round": "1"}},
		&PGN{Tags: map[string]string{
			"Event": "remove", "WhiteElo": "2400", "BlackElo": "2500", "Round": "2"}},
		&PGN{Tags: map[string]string{
			"Event": "keep", "WhiteElo": "2600", "BlackElo": "2300", "Round": "3"}},
	}
	return pgns
}

func TestFilterEquals(t *testing.T) {
	pgns := testPGNs()
	expected := pgns[1]
	filtered := FilterPGNs(pgns, NewTagFilter("Event==remove"))
	if len(filtered) != 1 || filtered[0] != expected {
		t.Log(filtered)
		t.Fail()
	}
}

func TestFilterNotEquals(t *testing.T) {
	pgns := testPGNs()
	filtered := FilterPGNs(pgns, NewTagFilter("Event!=remove"))
	if len(filtered) != 2 || filtered[0].Tags["Round"] != "1" || filtered[1].Tags["Round"] != "3" {
		t.Log(filtered)
		t.Fail()
	}
}

func TestFiltGT(t *testing.T) {
	pgns := testPGNs()
	f := NewTagFilter("BlackElo>2500")
	filtered := FilterPGNs(pgns, f)
	if len(filtered) != 1 || filtered[0].Tags["Round"] != "1" {
		t.Fail()
	}
}

func TestFiltGTEQ(t *testing.T) {
	pgns := testPGNs()
	f := NewTagFilter("BlackElo>=2500")
	filtered := FilterPGNs(pgns, f)
	if len(filtered) != 2 || filtered[0].Tags["Round"] != "1" || filtered[1].Tags["Round"] != "2" {
		t.Fail()
	}
}
