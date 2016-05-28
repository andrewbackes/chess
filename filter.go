package chess

import (
	"errors"
	"strconv"
	"strings"
)

// TagFilter is used to filter something like a PGN or Game based on
// some constraints on an assiciated tag.
//
// Possible Operators:
//      ">=", "<=", "!=", "==", "=", ">", "<"
type TagFilter struct {
	Tag      string
	Operator string
	Operand  string
}

// NewTagFilter makes a new TagFilter from a string.
func NewTagFilter(filter string) (TagFilter, error) {
	operators := []string{">=", "<=", "!=", "==", "=", ">", "<"}
	for _, op := range operators {
		if strings.Contains(filter, op) {
			split := strings.Split(filter, op)
			if len(split) >= 2 {
				return TagFilter{
					Tag:      split[0],
					Operator: op,
					Operand:  split[1],
				}, nil
			}
		}
	}
	return TagFilter{}, errors.New("could not parse filter")
}

// FilterPGNs filters a slice of PGNs based on a slice of filters.
func FilterPGNs(pgns []*PGN, filters ...TagFilter) []*PGN {
	var filtered []*PGN
	filterOut := func(pgn *PGN) bool {
		for _, f := range filters {
			if !satisfiesFilter(pgn, f) {
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

func satisfiesFilter(pgn *PGN, filter TagFilter) bool {
	if filter.Operator == "==" || filter.Operator == "=" {
		return pgn.Tags[filter.Tag] == filter.Operand
	}
	if filter.Operator == "!=" {
		return pgn.Tags[filter.Tag] != filter.Operand
	}
	pgnVal, err := strconv.Atoi(pgn.Tags[filter.Tag])
	if err != nil {
		return false
	}
	filterVal, err := strconv.Atoi(filter.Operand)
	if err != nil {
		return false
	}
	switch filter.Operator {
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
