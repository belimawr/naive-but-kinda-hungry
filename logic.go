package main

import (
	"context"
	"math/rand"
	"sort"

	"github.com/rs/zerolog"
)

// info returns the Battlesnake appearance and author permissions.
// For customization options, see https://docs.battlesnake.com/references/personalization
func info(ctx context.Context) BattlesnakeInfoResponse {
	zerolog.Ctx(ctx).Info().Msg("Sending snake info")
	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "Tiago Queiroz",
		Color:      "#b30000",
		Head:       "smart-caterpillar",
		Tail:       "round-bum",
	}
}

// start is called everytime the Battlesnake is entered into a game.
// The provided GameState contains information about the game that's about to be played.
func start(ctx context.Context, state GameState) {
	zerolog.Ctx(ctx).Info().Msg("Starting game!")
}

// end is called when a a game has ended.
func end(ctx context.Context, state GameState) {
	zerolog.Ctx(ctx).Info().Msg("The game has ended")
}

// mode is called on every turn and returns the Battlesnake's next move
func move(ctx context.Context, state GameState) BattlesnakeMoveResponse {
	logger := zerolog.Ctx(ctx)

	possibleMoves := safeMoves(state, true, true, true)
	safeMovesCount := 0
	for _, isSafe := range possibleMoves {
		if isSafe {
			safeMovesCount++
		}
	}

	if safeMovesCount == 0 {
		logger.Debug().Msg("no safe moves, trying a less restrictive approach") //
		nextMove := randomEmptySquare(ctx, state)
		logger.Info().Msgf("MOVE: %s", nextMove)

		return BattlesnakeMoveResponse{
			Move: nextMove,
		}
	}

	nextMove := findNextMove(ctx, state, possibleMoves)
	logger.Info().Msgf("MOVE: %s", nextMove)

	return BattlesnakeMoveResponse{
		Move: nextMove,
	}
}

func randomEmptySquare(ctx context.Context, state GameState) string {
	zerolog.Ctx(ctx).Debug().Msg("Random empty square")

	possibleMoves := safeMoves(state, true, false, true)

	return randomMove(ctx, state, possibleMoves)
}

func findNextMove(
	ctx context.Context,
	state GameState,
	possibleMoves map[string]bool,
) string {

	//	me := state.You

	// Find food
	// if me.Health < 75 {
	//	zerolog.Ctx(ctx).Debug().Msg("find food!")
	return gotoNearstFood(ctx, state, possibleMoves)
	// }

	// loop up and down
	// if possibleMoves["up"] {
	// 	return "up"
	// }

	// if possibleMoves["left"] {
	// 	return "left"
	// }

	// if possibleMoves["down"] {
	// 	return "down"
	// }

	// if possibleMoves["right"] {
	// 	return "right"
	// }

	//	return "down"
}

// onHarzard returns true if any part of a snake is into
// a hazard square
func onHarzard(me Battlesnake, hazard []Coord) bool {
	for _, sauce := range hazard {
		if me.Head == sauce {
			return true
		}
	}

	return false
}

func gotoNearstFood(ctx context.Context, state GameState, possibleMoves map[string]bool) string {
	logger := zerolog.Ctx(ctx)

	myHead := state.You.Head
	// If there is no food, move randomly
	if len(state.Board.Food) == 0 {
		return randomMove(ctx, state, possibleMoves)
	}

	food, safe := findNearstFood(state.You.Head, state.Board.Food, state.Board.Hazards)
	if !safe {
		logger.Info().Msg("NO SAFE FOOD")
		return randomMove(ctx, state, possibleMoves)
	}

	switch {
	case myHead.X > food.X:
		if possibleMoves["left"] {
			return "left"
		}
		fallthrough

	case myHead.X < food.X:
		if possibleMoves["right"] {
			return "right"
		}
		fallthrough

	case myHead.Y > food.Y:
		if possibleMoves["down"] {
			return "down"
		}
		fallthrough

	case myHead.Y < food.Y:
		if possibleMoves["up"] {
			return "up"
		}
		fallthrough

	default:
		return randomMove(ctx, state, possibleMoves)
	}
}

func movesToSlice(possibleMoves map[string]bool) []string {
	safeMovesSlice := []string{}
	for move, isSafe := range possibleMoves {
		if isSafe {
			safeMovesSlice = append(safeMovesSlice, move)
		}
	}

	return safeMovesSlice
}

func randomMove(
	ctx context.Context,
	state GameState,
	possibleMoves map[string]bool,
) string {

	safeMovesSlice := movesToSlice(possibleMoves)

	if len(safeMovesSlice) == 0 {
		zerolog.Ctx(ctx).Info().Msg("NO SAFE MOVES! Trying to go into sauce")

		desperateMovesSlice := movesToSlice(safeMoves(state, true, false, false))
		if len(desperateMovesSlice) != 0 {
			return desperateMovesSlice[rand.Intn(len(desperateMovesSlice))]
		}

		zerolog.Ctx(ctx).Info().Msg("NO SAFE MOVES! Going down")
		return "down"
	}

	zerolog.Ctx(ctx).Info().Msg("Random move!")
	return safeMovesSlice[rand.Intn(len(safeMovesSlice))]
}

// adjacentPoints returns the adjacent points to any point
func adjacentPoints(head Coord) []Coord {
	moves := []Coord{}
	moves = append(moves, Coord{X: head.X + 1, Y: head.Y})
	moves = append(moves, Coord{X: head.X - 1, Y: head.Y})
	moves = append(moves, Coord{X: head.X, Y: head.Y + 1})
	moves = append(moves, Coord{X: head.X, Y: head.Y - 1})

	return moves
}

// safeMoves returns all safe moves.
// - allowHeadColision: Allows for head collisions when we're bigger
// - predictMoves: makes unsae any move that could hit any
//   snake's head on the next move
func safeMoves(
	state GameState,
	allowHeadColision, predictMoves bool,
	avoidHazard bool,
) map[string]bool {

	possibleMoves := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	me := state.You
	myHead := me.Head
	boardWidth := state.Board.Width
	boardHeight := state.Board.Height

	// Avoid Walls
	if myHead.X == boardWidth-1 {
		possibleMoves["right"] = false
	}

	if myHead.X == 0 {
		possibleMoves["left"] = false
	}

	if myHead.Y == boardHeight-1 {
		possibleMoves["up"] = false
	}

	if myHead.Y == 0 {
		possibleMoves["down"] = false
	}

	// Avoid hitting other noGoCoords
	noGoCoords := map[Coord]struct{}{}
	// Add ourselves
	for _, c := range me.Body {
		noGoCoords[c] = struct{}{}
	}

	// Add all other snakes and allow for head-to-head
	// if we can win
	for _, s := range state.Board.Snakes {
		if me.ID == s.ID {
			continue
		}

		// Add the whole snake's body
		for _, p := range s.Body {
			noGoCoords[p] = struct{}{}
		}

		if allowHeadColision {
			// If we win on head to head collisions
			// remove the other snake's head from noGoCoords
			if me.Health > s.Health {
				delete(noGoCoords, s.Head)
			}
		}
	}

	if predictMoves {
		// Add possible next moves fom all snakes
		// but ourselves
		for _, s := range state.Board.Snakes {
			if me.ID == s.ID {
				continue
			}

			for _, p := range adjacentPoints(s.Head) {
				noGoCoords[p] = struct{}{}
			}
		}
	}

	if avoidHazard {
		for _, sauce := range state.Board.Hazards {
			noGoCoords[sauce] = struct{}{}
		}
	}

	// For each possible move, verify which ones are safe
	for p := range noGoCoords {
		for _, m := range []string{"up", "down", "left", "right"} {
			switch m {
			case "up":
				nextHead := myHead
				nextHead.Y++
				if p == nextHead {
					possibleMoves["up"] = false
				}
				break

			case "down":
				nextHead := myHead
				nextHead.Y--
				if p == nextHead {
					possibleMoves["down"] = false
				}
				break

			case "left":
				nextHead := myHead
				nextHead.X--
				if p == nextHead {
					possibleMoves["left"] = false
				}
				break

			case "right":
				nextHead := myHead
				nextHead.X++
				if p == nextHead {
					possibleMoves["right"] = false
				}
				break
			}
		}
	}

	return possibleMoves
}

// findNearstFood returns the nearst food that is not in
// a hazard sauce
func findNearstFood(head Coord, foods []Coord, hazard []Coord) (Coord, bool) {
	hazardSet := map[Coord]struct{}{}
	for _, h := range hazard {
		hazardSet[h] = struct{}{}
	}

	foodSet := map[Coord]struct{}{}
	for _, h := range foods {
		foodSet[h] = struct{}{}
	}

	safeFood := []struct {
		p    Coord
		dist float64
	}{}

	for food := range foodSet {
		if _, notSafeFood := hazardSet[food]; notSafeFood {
			continue
		}

		p := struct {
			p    Coord
			dist float64
		}{
			p:    food,
			dist: dist(head, food),
		}

		safeFood = append(safeFood, p)
	}

	if len(safeFood) == 0 {
		return Coord{}, false
	}

	sort.Slice(safeFood, func(i, j int) bool {
		return safeFood[i].dist < safeFood[j].dist
	})

	return safeFood[0].p, true
}
