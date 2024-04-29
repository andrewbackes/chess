package piece

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestPiecePrint(t *testing.T) {
	testCases := []struct {
		name string
		p    Piece
		want string
	}{
		{"White-Pawn", New(White, Pawn), "P"},
		{"Black-Pawn", New(Black, Pawn), "p"},
		{"NoColor-Pawn", New(NoColor, Pawn), "p"},
		{"Color(4)-Pawn", New(Color(4), Pawn), "p"},
		{"White-None", New(White, None), " "},
		{"Black-None", New(Black, None), " "},
		{"NoColor-None", New(NoColor, None), " "},
		{"Color(4)-None", New(Color(4), None), " "},
		{"White-Rook", New(White, Rook), "R"},
		{"Black-Queen", New(Black, Queen), "q"},
		{"NoColor-Knight", New(NoColor, Knight), "n"},
		{"Color(4)-King", New(Color(4), King), "k"},
		{"White-Type(10)", New(White, Type(10)), ""},
		{"Black-Type(10)", New(Black, Type(10)), ""},
		{"NoColor-Type(10)", New(NoColor, Type(10)), ""},
		{"Color(4)-Type(10)", New(Color(4), Type(10)), ""},
		//{"Color(-1)-Type(-1)", New(Color(-1), Type(-1)), ""}, // Build error: constant -1 overflows Color/Type.
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if s := fmt.Sprint(tc.p); s != tc.want {
				t.Errorf("%#v.String() = %#v, want %#v", tc.p, s, tc.want)
			}
		})
	}
}

func TestColorUnmarshalJson(t *testing.T) {
	blob := `["White"]`
	want := []Color{White}
	var c []Color
	err := json.Unmarshal([]byte(blob), &c)
	if err != nil || !reflect.DeepEqual(c, want) {
		t.Fatalf("json.Unmarshal(%#v, &c) = %#v, want %#v", blob, err, error(nil))
	} else if !reflect.DeepEqual(c, want) {
		t.Fatalf("json.Unmarshal(%#v, &c); c = %#v, want %#v", blob, c, want)
	}
}

func TestColorMarshalJson(t *testing.T) {
	j := []Color{White}
	want := `["White"]`
	result, err := json.Marshal(j)
	if err != nil || string(result) != want {
		t.Fatalf("json.Marshal(%#v) = (%s, %#v), want (%s, %#v)", j, result, err, want, error(nil))
	}
}

func TestTypeString(t *testing.T) {
	testCases := []struct {
		name string
		ty   Type
		want string
	}{
		{"Pawn", Pawn, "p"},
		{"Rook", Rook, "r"},
		{"King", King, "k"},
		{"None", None, " "},
		{"Type(10)", Type(10), ""},
		//{"Type(-1)", Type(-1), ""}, // Build error: constant -1 overflows Type.
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if s := tc.ty.String(); s != tc.want {
				t.Errorf("%#v.String() = %#v, want %#v", tc.ty, s, tc.want)
			}
		})
	}
}

func TestPieceFigurine(t *testing.T) {
	testCases := []struct {
		name string
		p    Piece
		want string
	}{
		{"White-Pawn", New(White, Pawn), "♙"},
		{"Black-Pawn", New(Black, Pawn), "♟"},
		{"NoColor-Pawn", New(NoColor, Pawn), " "},
		{"Black-None", New(Black, None), " "},
		{"NoColor-None", New(NoColor, None), " "},
		{"Color(4)-Pawn", New(Color(4), Pawn), "\x00"},
		{"Black-Type(10)", New(Black, Type(10)), "\x00"},
		{"Color(4)-Type(10)", New(Color(4), Type(10)), "\x00"},
		//{"Color(-1)-Type(-1)", New(Color(-1), Type(-1)), ""}, // Build error: constant -1 overflows Color/Type.
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if f := tc.p.Figurine(); f != tc.want {
				t.Errorf("%#v.Figurine() = %#v, want %#v", tc.p, f, tc.want)
			}
		})
	}
}

const benchValidTypesCount = 7 // None, Pawn, ..., King.

func BenchmarkTypeString(b *testing.B) {
	b.Run("Valid-Types", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Is Type(0), Type(1), ..., Type(6), Type(0), Type(1), ...
			s := Type(i % benchValidTypesCount).String()
			if s == "" {
				b.Fatal("Valid type should not have a empty string output")
			}
		}
	})
	b.Run("Invalid-Types", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Is Type(7), Type(8), ..., Type(13), Type(7), Type(8), ...
			s := Type(benchValidTypesCount + i%benchValidTypesCount).String()
			if s != "" {
				b.Fatal("Invalid type should have a empty string output")
			}
		}
	})
}

func BenchmarkPieceFigurine(b *testing.B) {
	b.Run("Valid-Colors_Valid-Types", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// In the piece following colors and types are present.
			// Color(0), Color(1), Color(2), Color(0), Color(1), Color(2), ...
			// Type(0), Type(1), ..., Type(6), Type(0), Type(1), ...
			p := Piece{Color(i % benchValidColorCount), Type(i % benchValidTypesCount)}
			s := p.Figurine()
			if s == "\x00" {
				b.Fatalf("Figurine of piece with valid color and type %#v should not have an empty string output", p)
			}
		}
	})
	b.Run("Invalid-Colors_Valid-Types", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// In the piece following colors and types are present.
			// Color(3), Color(4), Color(5), Color(3), Color(4), Color(5), ...
			// Type(0), Type(1), ..., Type(6), Type(0), Type(1), ...
			p := Piece{Color(benchValidColorCount + i%benchValidColorCount), Type(i % benchValidTypesCount)}
			s := p.Figurine()
			if s != "\x00" && s != " " {
				b.Fatalf("Figurine of piece with invalid color and valid type %#v, should have an empty string output, got: %#v", p, s)
			}
		}
	})
	b.Run("Valid-Colors_Invalid-Types", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// In the piece following colors and types are present.
			// Color(0), Color(1), Color(2), Color(0), Color(1), Color(2), ...
			// Type(7), Type(8), ..., Type(13), Type(7), Type(8), ...
			p := Piece{Color(i % benchValidColorCount), Type(benchValidTypesCount + i%benchValidTypesCount)}
			s := p.Figurine()
			if s != "\x00" && s != " " {
				b.Fatalf("Figurine of piece with valid color and invalid type %#v should have an empty string output, got: %#v", p, s)
			}
		}
	})
	b.Run("Invalid-Colors_Invalid-Types", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// In the piece following colors and types are present.
			// Color(3), Color(4), Color(5), Color(0), Color(1), Color(2), ...
			// Type(7), Type(8), ..., Type(13), Type(7), Type(8), ...
			p := Piece{Color(benchValidColorCount + i%benchValidColorCount), Type(benchValidTypesCount + i%benchValidTypesCount)}
			s := p.Figurine()
			if s != "\x00" && s != " " {
				b.Fatalf("Figurine of piece with invalid color and type %#v should have an empty string output, got: %#v", p, s)
			}
		}
	})
}
