// Package book is for working with polyglot opening books.
package book

import (
	"encoding/binary"
	"github.com/andrewbackes/chess/position"
	"os"
	"sort"
)

// Book is a polyglot opening book loaded into memory.
type Book struct {
	Positions map[uint64][]Entry
	seed      uint64
}

// New makes a new empty opening book
func New() *Book {
	return &Book{
		Positions: make(map[uint64][]Entry),
	}
}

// Entry is a weighted move in an internally loaded opening book.
type Entry struct {
	Move   string
	Weight uint16
	Learn  uint32
}

// PolyglotEntry is a line in a polyglot opening book.
type PolyglotEntry struct {
	Key   uint64
	Move  uint16
	Score uint16
	Learn uint32
}

type byWeight []Entry

func (m byWeight) Len() int      { return len(m) }
func (m byWeight) Swap(i, j int) { t := m[i]; m[i] = m[j]; m[j] = t }
func (m byWeight) Less(i, j int) bool {
	if m[i].Weight > m[j].Weight {
		return true
	} else if m[i].Weight == m[j].Weight {
		if m[i].Move > m[j].Move {
			return true
		}
	}
	return false
}

// Save saves the opening book in binary format.
func (b *Book) Save(filename string) error {
	var file *os.File
	var err error
	if _, er := os.Stat(filename); os.IsNotExist(er) {
		// file doesnt exist
	} else if er == nil {
		// file does exist
		os.Remove(filename)
	}
	file, err = os.Create(filename)
	defer file.Close()
	if err != nil {
		return err
	}
	for key, moves := range b.Positions {
		for _, entry := range moves {
			e := PolyglotEntry{
				Key:   key,
				Move:  encodedMove(entry.Move),
				Score: entry.Weight,
				Learn: entry.Learn,
			}
			err = binary.Write(file, binary.BigEndian, &e)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Open loads a polyglot opening book into memory.
// filename is the full path to a .bin opening book.
func Open(filename string) (*Book, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	book := &Book{
		Positions: make(map[uint64][]Entry),
	}
	entry := PolyglotEntry{}
	key := uint64(0)
	for {
		err := binary.Read(file, binary.BigEndian, &entry)
		if err != nil {
			// EOF
			break
		}
		if entry.Move != 0 {
			mv := Entry{
				Move:   decodeMove(entry.Move),
				Weight: entry.Score,
				Learn:  entry.Learn,
			}
			book.Positions[entry.Key] = append(book.Positions[entry.Key], mv)
		}
		if entry.Key != key {
			sort.Sort(byWeight(book.Positions[key]))
			key = entry.Key
		}
	}
	return book, nil
}

func encodedMove(move string) uint16 {
	mv := move
	switch mv {
	case "e1g1":
		mv = "e1h1"
	case "e1c1":
		mv = "e1a1"
	case "e8g8":
		mv = "e8h8"
	case "e8c8":
		mv = "e8a8"
	}
	from, to := position.Split(position.Move(mv))
	fromFile, fromRank := indexToFR(int(from))
	toFile, toRank := indexToFR(int(to))
	var promo uint16
	if len(mv) > 4 {
		promo = map[string]uint16{"": 0, "k": 1, "b": 2, "r": 3, "q": 4}[string(mv[4])]
	}
	return (promo << 12) + (uint16(fromRank) << 9) + (uint16(fromFile) << 6) + (uint16(toRank) << 3) + (uint16(toFile))
}

func indexToFR(index int) (file int, row int) {
	// 0  --> h1 --> 7,0
	// 7  --> a1 --> 0,0 (row,file)
	// 63 --> a8 --> 0,7
	// 56 --> h8 --> 7,7
	row = index / 8
	file = 7 - (index % 8)
	return
}

/*

bits                meaning
===================================
0,1,2               to file
3,4,5               to row
6,7,8               from file
9,10,11             from row
12,13,14            promotion piece

"promotion piece" is encoded as follows
none       0
knight     1
bishop     2
rook       3
queen      4

white short      e1h1
white long       e1a1
black short      e8h8
black long       e8a8

*/
func decodeMove(move uint16) string {
	fromRank := bits(move, 3)
	fromFile := bits(move, 2)
	from := position.NewSquare(uint(fromFile+1), uint(fromRank+1))

	toRank := bits(move, 1)
	toFile := bits(move, 0)
	to := position.NewSquare(uint(toFile+1), uint(toRank+1))

	promo := bits(move, 4)
	var promoStr string
	if promo < 4 {
		promoStr = []string{"", "k", "b", "r", "q"}[promo]
	}
	mv := from.String() + to.String() + promoStr
	switch mv {
	case "e1h1":
		return "e1g1"
	case "e1a1":
		return "e1c1"
	case "e8h8":
		return "e8g8"
	case "e8a8":
		return "e8c8"
	}
	return (mv)
}

func bits(move uint16, group uint) uint16 {
	mask := uint16(7)
	mask = mask << (3 * group)
	return (move & mask) >> (3 * group)
}
