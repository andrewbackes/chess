package piece

import "testing"

func TestColorString(t *testing.T) {
	testCases := []struct {
		name string
		c    Color
		want string
	}{
		{"White", White, "White"},
		{"Black", Black, "Black"},
		{"NoColor", NoColor, ""},
		{"Color(4)", Color(4), ""},
		//{"Color(-1)", Color(-1), ""}, // Build error: constant -1 overflows Color.
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if s := tc.c.String(); s != tc.want {
				t.Errorf("%#v.String() = %s, want %s", tc.c, s, tc.want)
			}
		})
	}
}

const benchValidColorCount = 2 // White, Black.

func BenchmarkColorString(b *testing.B) {
	b.Run("Valid-Colors", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Is Color(0), Color(1), Color(0), Color(1), ...
			s := Color(i % benchValidColorCount).String()
			if s == "" {
				b.Fatal("Valid color should not have a empty string output")
			}
		}
	})
	b.Run("Invalid-Colors", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Is Color(2), Color(3), Color(2), Color(3), ...
			s := Color(benchValidColorCount + i%benchValidColorCount).String()
			if s != "" {
				b.Fatal("Invalid color should have a empty string output")
			}
		}
	})
}
