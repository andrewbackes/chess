package book

import (
	"fmt"
	"os"
	"testing"
)

func TestOpenBook(t *testing.T) {
	if os.Getenv("TEST_BOOK") == "" {
		t.SkipNow()
	}
	book, err := FromBin(os.Getenv("TEST_BOOK"))
	if err != nil {
		t.Fail()
	}
	for k, v := range book.Positions {
		fmt.Println(k, v)
		break
	}
	fmt.Println(len(book.Positions))
}

func TestBitsMasks(t *testing.T) {
	move := uint16(65535)
	for i := uint(0); i < 5; i++ {
		g := bits(move, i)
		if g != 7 {
			t.Error(g)
		}
	}
}
