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
	file, _ := os.Open(os.Getenv("TEST_BOOK"))
	defer file.Close()
	book, err := Read(file)
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
	file, _ := os.Open(os.Getenv("TEST_BOOK"))
	defer file.Close()
	opened, err := Read(file)
	if err != nil {
		t.Error(err)
	}
	dest := "/tmp/savedbook.bin"
	if err := opened.Save(dest); err != nil {
		t.Error(err)
	}
	d, _ := os.Open(dest)
	defer file.Close()
	saved, err := Read(d)
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
