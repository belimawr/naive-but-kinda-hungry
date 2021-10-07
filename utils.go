package main

// adjacentPoints returns the adjacent points to any point
// the returned slice follows the order:
// up, down, left, right
func adjacentPoints(head Coord) []Coord {
	moves := []Coord{}

	moves = append(moves, Coord{X: head.X, Y: head.Y + 1}) // up
	moves = append(moves, Coord{X: head.X, Y: head.Y - 1}) // down
	moves = append(moves, Coord{X: head.X - 1, Y: head.Y}) // left
	moves = append(moves, Coord{X: head.X + 1, Y: head.Y}) // right

	return moves
}
