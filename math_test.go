package main

import (
	"fmt"
	"math"
	"testing"
)

func TestDist(t *testing.T) {
	// The differene between two the expected result
	// and the current result must be smaller than this delta
	delta := 0.0001
	testCases := []struct {
		a    Coord
		b    Coord
		want float64
	}{
		{
			a:    Coord{3, 2},
			b:    Coord{9, 7},
			want: 7.8102,
		},
		{
			a:    Coord{4, 3},
			b:    Coord{9, 4},
			want: 5.099,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.want), func(t *testing.T) {
			got, want := dist(tc.a, tc.b), tc.want
			if math.Abs(got-want) > delta {
				t.Errorf("expecting %f, got: %f", want, got)
			}
		})
	}
}
