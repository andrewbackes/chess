// Package book is for working with polyglot opening books.
package book

import (
	"encoding/binary"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/chess/position/square"
	"io"
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
	Move   move.Move
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
		if m[i].Move.String() > m[j].Move.String() {
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
		sort.Sort(byWeight(moves))
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

// Read loads a polyglot opening book into memory.
// binFile is the already opened polyglot bin file.
func Read(binFile io.Reader) (*Book, error) {
	book := &Book{
		Positions: make(map[uint64][]Entry),
	}
	entry := PolyglotEntry{}
	key := uint64(0)
	for {
		err := binary.Read(binFile, binary.BigEndian, &entry)
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

func encodedMove(mv move.Move) uint16 {
	switch mv.String() {
	case "e1g1":
		mv = move.Parse("e1h1")
	case "e1c1":
		mv = move.Parse("e1a1")
	case "e8g8":
		mv = move.Parse("e8h8")
	case "e8c8":
		mv = move.Parse("e8a8")
	}
	from, to := mv.From(), mv.To()
	fromFile, fromRank := indexToFR(int(from))
	toFile, toRank := indexToFR(int(to))
	return (uint16(mv.Promote) << 12) + (uint16(fromRank) << 9) + (uint16(fromFile) << 6) + (uint16(toRank) << 3) + (uint16(toFile))
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
func decodeMove(m uint16) move.Move {
	fromRank := bits(m, 3)
	fromFile := bits(m, 2)
	from := square.New(uint(fromFile+1), uint(fromRank+1))

	toRank := bits(m, 1)
	toFile := bits(m, 0)
	to := square.New(uint(toFile+1), uint(toRank+1))

	promo := bits(m, 4)
	var promoStr string
	if promo < 4 {
		promoStr = []string{"", "k", "b", "r", "q"}[promo]
	}
	mv := from.String() + to.String() + promoStr
	switch mv {
	case "e1h1":
		return move.Parse("e1g1")
	case "e1a1":
		return move.Parse("e1c1")
	case "e8h8":
		return move.Parse("e8g8")
	case "e8a8":
		return move.Parse("e8c8")
	}
	return move.Parse(mv)
}

func bits(move uint16, group uint) uint16 {
	mask := uint16(7)
	mask = mask << (3 * group)
	return (move & mask) >> (3 * group)
}
