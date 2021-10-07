package main

import (
	"context"
	"math/rand"
	"testing"
	"time"
)

func TestSomething(t *testing.T) {
	testCases := []struct {
		name       string
		state      GameState
		validMoves []Move
	}{
		{
			name:       "two foods, two snakes, different distances",
			validMoves: []Move{Left, Rignt},
			state: GameState{
				Board: Board{
					Width: 5, Height: 5,
					Snakes: []Battlesnake{
						{Body: []Coord{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}}},
						{Body: []Coord{{X: 4, Y: 2}, {X: 4, Y: 1}, {X: 4, Y: 0}}},
					},
					Food: []Coord{{0, 0}, {2, 0}, {4, 4}},
				},
				You: Battlesnake{
					Body: []Coord{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}},
					Head: Coord{X: 1, Y: 0},
				},
			},
		},
		{
			name:       "two foods, one snake, same distances",
			validMoves: []Move{Left, Rignt, Down},
			state: GameState{
				Board: Board{
					Width: 5, Height: 5,
					Snakes: []Battlesnake{
						{Body: []Coord{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 1, Y: 3}}},
					},
					Food: []Coord{{0, 0}, {2, 0}},
				},
				You: Battlesnake{
					Body: []Coord{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 1, Y: 3}},
					Head: Coord{X: 1, Y: 1},
				},
			},
		},
		{
			name:       "3 foods, one snake, same distances, one blocked by body",
			validMoves: []Move{Left, Rignt, Down},
			state: GameState{
				Board: Board{
					Width: 5, Height: 5,
					Snakes: []Battlesnake{
						{Body: []Coord{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 1, Y: 3}}},
					},
					Food: []Coord{{0, 0}, {0, 2}, {2, 0}},
				},
				You: Battlesnake{
					Body: []Coord{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 1, Y: 3}},
					Head: Coord{X: 1, Y: 1},
				},
			},
		},
		{
			name:       "3 foods, one snake, 1 food in next move",
			validMoves: []Move{Down},
			state: GameState{
				Board: Board{
					Width: 5, Height: 5,
					Snakes: []Battlesnake{
						{Body: []Coord{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 1, Y: 3}}},
					},
					Food: []Coord{{0, 0}, {1, 0}, {0, 2}},
				},
				You: Battlesnake{
					Body: []Coord{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 1, Y: 3}},
					Head: Coord{X: 1, Y: 1},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Run every test with a different seed
			seed := time.Now().UnixNano()
			rand.Seed(seed)

			move := hungry(context.Background(), tc.state)
			valid := false
			for _, expected := range tc.validMoves {
				if move == expected.String() {
					valid = true
					break
				}
			}

			if !valid {
				t.Errorf("got invalid move '%s', seed: %d", move, seed)
			}
		})
	}
}
