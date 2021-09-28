package main

import "math"

func dist(a, b Coord) float64 {
	xDelta := float64(a.X - b.X)
	yDelta := float64(a.Y - b.Y)
	dist := math.Sqrt(math.Pow(xDelta, 2) + math.Pow(yDelta, 2))
	return dist
}
