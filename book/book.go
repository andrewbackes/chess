// Package book is for working with polyglot opening books.
package book

import (
	"encoding/binary"
	"os"
)

// Book is a polyglot opening book loaded into memory.
type Book struct {
	Positions map[uint64][]Move
}

// Move is a weighted move in an internally loaded opening book.
type Move struct {
	Move   string
	Weight uint16
}

// PolyglotEntry is a line in a polyglot opening book.
type PolyglotEntry struct {
	Key   uint64
	Move  uint16
	Score uint16
	Learn uint32
}

// FromBin loads a polyglot opening book into memory.
// filename is the full path to a .bin opening book.
func FromBin(filename string) (*Book, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	book := &Book{
		Positions: make(map[uint64][]Move),
	}
	entry := PolyglotEntry{}
	for {
		err := binary.Read(file, binary.BigEndian, &entry)
		if err != nil {
			// EOF
			break
		}
		if entry.Move != 0 {
			mv := Move{
				Move:   decodeMove(entry.Move),
				Weight: entry.Score,
			}
			book.Positions[entry.Key] = append(book.Positions[entry.Key], mv)
		}
	}
	return book, nil
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
	/*
		toFile := bits(move, 0)
		toRank := bits(move, 1)
		fromFile := bits(move, 2)
		fromRank := bits(move, 3)
		promo := bits(move, 4)
	*/
	return "0000"
}

func bits(move uint16, group uint) uint16 {
	mask := uint16(7)
	mask = mask << (3 * group)
	return (move & mask) >> (3 * group)
}
