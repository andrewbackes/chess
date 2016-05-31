package pgn

import (
	"strconv"
	"strings"
)

// Filterer is an interface used by the Filter function to decide if a PGN
// should be kept or not.
type Filterer interface {
	Include(*PGN) bool
}

// TagFilter is used to filter a PGN based on some constraints on an assiciated tag.
//
// Possible Operators:
//      ">=", "<=", "!=", "==", "=", ">", "<"
type TagFilter struct {
	Tag      string
	Operator string
	Operand  string
}

// NewTagFilter makes a new Filter from a string.
func NewTagFilter(filter string) TagFilter {
	operators := []string{">=", "<=", "!=", "==", "=", ">", "<"}
	for _, op := range operators {
		if strings.Contains(filter, op) {
			split := strings.Split(filter, op)
			if len(split) >= 2 {
				return TagFilter{
					Tag:      split[0],
					Operator: op,
					Operand:  split[1],
				}
			}
		}
	}
	return TagFilter{}
}

// Filter filters a slice of PGNs based on a slice of filters.
func Filter(pgns []*PGN, filters ...Filterer) []*PGN {
	var filtered []*PGN
	filterOut := func(pgn *PGN) bool {
		for _, f := range filters {
			if !f.Include(pgn) {
				return true
			}
		}
		return false
	}
	for _, pgn := range pgns {
		if !filterOut(pgn) {
			filtered = append(filtered, pgn)
		}
	}
	return filtered
}

// Include makes TagFilter impelement the Filterer interface. It decides if a
// PGN meets the requirements of a filter.
func (t TagFilter) Include(pgn *PGN) bool {
	if t.Operator == "==" || t.Operator == "=" {
		return pgn.Tags[t.Tag] == t.Operand
	}
	if t.Operator == "!=" {
		return pgn.Tags[t.Tag] != t.Operand
	}
	pgnVal, err := strconv.Atoi(pgn.Tags[t.Tag])
	if err != nil {
		return false
	}
	filterVal, err := strconv.Atoi(t.Operand)
	if err != nil {
		return false
	}
	switch t.Operator {
	case ">=":
		return pgnVal >= filterVal
	case "<=":
		return pgnVal <= filterVal
	case ">":
		return pgnVal > filterVal
	case "<":
		return pgnVal < filterVal
	}
	return false
}
