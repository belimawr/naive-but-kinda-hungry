package main

type Move int

const (
	Up Move = iota
	Down
	Left
	Rignt
)

var moveNames = []string{"up", "down", "left", "right"}

// String returns the string representation of the Move
func (m Move) String() string {
	if Up <= m && m <= Rignt {
		return moveNames[m]
	}
	return "unknown"
}

type MoveCoord struct {
	Move
	Coord
}
