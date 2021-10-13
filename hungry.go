package main

import (
	"context"
	"fmt"
	"math/rand"
	"sort"

	"github.com/rs/zerolog"
)

// hungry uses the following strategy:
// 1. Map the distance to all food
// 2. Go to the nearst food via the shortest path
//   2.1 If there is more then 1 nearst food, pick one randomly
//   2.2 If there is more then 1 shortest path, pick one randomly
func hungry(ctx context.Context, state GameState) string {
	board := NewBoardMap(state.Board.Height) // We assume it's a square

	MarkObstacles(board, state)
	CalculateDistanceFromFood(board, state)

	head := state.You.Head
	// Possible moves are all the adjacent cells to the head
	possibleMoves := possibleMoves(board, head, state.Board.Height)

	// Put them all into a slice containing:
	// - the position
	// - the movement (up, down, left, right)
	// - the distance to the food
	dists := []struct {
		MoveCoord
		dist int
	}{}
	for _, move := range possibleMoves {
		if board[move.X][move.Y] < 0 {
			continue
		}

		pd := struct {
			MoveCoord
			dist int
		}{move, board[move.X][move.Y]}
		dists = append(dists, pd)
	}

	if len(dists) == 0 {
		// If cannot go the the nearst food, random move to a safe place
		zerolog.Ctx(ctx).Info().Msg("cannot go to food, random safe move")
		safe := []MoveCoord{}
		for _, move := range possibleMoves {
			if board[move.X][move.Y] < 0 {
				safe = append(safe, move)
			}
		}

		if len(safe) != 0 {
			return safe[rand.Intn(len(safe))].String()
		}

		zerolog.Ctx(ctx).Info().Msg("NO SAFE MOVES, going down")
		fmt.Println(gameMapToString(board))
		return "down"
	}

	// sort by distance to food
	sort.Slice(dists, func(i, j int) bool { return dists[i].dist < dists[j].dist })

	fmt.Println("Moves sorted by distance: ", dists)
	// put all moves that are equally further from the a food
	// into a slice
	closest := dists[0].dist
	moves := []string{}
	for _, move := range dists {
		if move.dist == closest {
			moves = append(moves, move.String())
			continue
		}
		break
	}

	return moves[rand.Intn(len(moves))]
}

func possibleMoves(board BoardMap, head Coord, size int) []MoveCoord {
	possibleMoves := []MoveCoord{}

	// for moves in: up, down, left, right
	for i, p := range adjacentPoints(head) {
		if isOutOfBounds(p, size, size) {
			continue
		}

		if isSnake(board[p.X][p.Y]) {
			continue
		}

		possibleMoves = append(possibleMoves, MoveCoord{Move(i), p})
	}

	return possibleMoves
}

func isSnake(v int) bool {
	if v == SnakeBody || v == SnakeHead {
		return true
	}
	return false
}
