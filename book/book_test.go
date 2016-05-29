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
	book, err := Open(os.Getenv("TEST_BOOK"))
	if err != nil {
		t.Fail()
	}
	/*
		for k, v := range book.Positions {
			fmt.Println(k, v)
			//break
		}
	*/
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

func TestOpenSaveOpen(t *testing.T) {
	if os.Getenv("TEST_BOOK") == "" {
		t.SkipNow()
	}
	opened, err := Open(os.Getenv("TEST_BOOK"))
	if err != nil {
		t.Error(err)
	}
	dest := "/tmp/savedbook.bin"
	if err := opened.Save(dest); err != nil {
		t.Error(err)
	}
	saved, err := Open(dest)
	if err != nil {
		t.Error(err)
	}
	if len(opened.Positions) != len(saved.Positions) {
		t.Error(len(opened.Positions), "!=", len(saved.Positions))
	}
	for k, v := range opened.Positions {
		if v2, ok := saved.Positions[k]; !ok {
			t.Error("key missing")
		} else {
			if fmt.Sprint(v) != fmt.Sprint(v2) {
				t.Error(fmt.Sprint(v), "!=", fmt.Sprint(v2))
			}
		}
	}
}
