package main

import (
	"strconv"
	"testing"
)

func TestFindNearstFood(t *testing.T) {
	testCases := []struct {
		head      Coord
		foods     []Coord
		hazards   []Coord
		wantCoord Coord
		wantSafe  bool
	}{
		{
			head:      Coord{0, 0},
			foods:     []Coord{{1, 1}, {2, 2}},
			wantCoord: Coord{1, 1},
			wantSafe:  true,
		},
		{ // Not a real game state, but still a valid test
			head:      Coord{2, 2},
			foods:     []Coord{{1, 1}, {2, 2}},
			wantCoord: Coord{2, 2},
			wantSafe:  true,
		},
		{
			head:      Coord{0, 0},
			foods:     []Coord{{1, 1}, {2, 2}},
			hazards:   []Coord{{1, 1}, {2, 2}},
			wantCoord: Coord{0, 0},
			wantSafe:  false,
		},
		{
			head:      Coord{0, 0},
			foods:     []Coord{{1, 1}, {10, 2}},
			hazards:   []Coord{{1, 1}, {2, 2}},
			wantCoord: Coord{10, 2},
			wantSafe:  true,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, safe := findNearstFood(tc.head, tc.foods, tc.hazards)
			if safe != tc.wantSafe {
				t.Errorf("expecting safe %t, got %t", tc.wantSafe, safe)
			}

			if got != tc.wantCoord {
				t.Errorf("expecting: %#v, got: %#v", tc.wantCoord, got)
			}
		})
	}
}
