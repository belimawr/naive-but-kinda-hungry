package main

type BoardMap [][]int

const (
	Unitialised   = -1
	SnakeBody     = -2
	SnakeHead     = -3
	SnakeTail     = -4
	LookAheadHead = -5
	Sauce         = -6
	Food          = 0
)

// NewBoardMap returns a initialised NewBoardMap
// of the desired size
func NewBoardMap(size int) BoardMap {
	board := make(BoardMap, size)
	for i := range board {
		board[i] = make([]int, size)
		for j := range board[i] {
			board[i][j] = Unitialised
		}
	}

	return board
}

// MarkObstacles marks everything we need to be awere of
// on the board:
// - Snakes: head, body and tails
// - Possible snakes next head position
// - Hazard sauce
func MarkObstacles(board BoardMap, state GameState) {
	for _, snake := range state.Board.Snakes {
		board[snake.Head.X][snake.Head.Y] = SnakeHead

		// Look ahead
		for _, p := range adjacentPoints(snake.Head) {
			if isOutOfBounds(p, state.Board.Height, state.Board.Width) {
				continue
			}

			if board[p.X][p.Y] == Unitialised {
				board[p.X][p.Y] = LookAheadHead
			}
		}

		for i := 1; i < len(snake.Body)-1; i++ {
			p := snake.Body[i]
			board[p.X][p.Y] = SnakeBody
		}

		tail := snake.Body[len(snake.Body)-1]
		board[tail.X][tail.Y] = SnakeTail
	}
}

// CalculateDistanceFromFood sets the distances from any empty
// cell in the board to the nearst food. Food cells are marked as 0.
// Snakes already marked on the board are left untouched.
func CalculateDistanceFromFood(board BoardMap, state GameState) {
	bfs := make(chan Coord, state.Board.Height*state.Board.Width)

	// Put all the food into bfs and mark their distance as 0
	for _, food := range state.Board.Food {
		bfs <- food
		board[food.X][food.Y] = 0
	}

	for len(bfs) > 0 {
		current := <-bfs

		for _, adjacent := range adjacentPoints(current) {
			if isOutOfBounds(adjacent, state.Board.Width, state.Board.Width) {
				continue
			}

			x, y := adjacent.X, adjacent.Y
			// skip points that have been processed or are snakes
			if board[x][y] != Unitialised {
				continue
			}

			// calculate it's distance
			board[x][y] = board[current.X][current.Y] + 1
			// add it to the bfs
			bfs <- adjacent
		}
	}
}

func isOutOfBounds(p Coord, h, w int) bool {
	x, y := p.X, p.Y
	return x < 0 || y < 0 || x >= w || y >= h
}

type Node struct {
	// We need something to differentiate the nodes.
	ID       int
	Children []Node
}

func NewNode(id int, children ...Node) Node {
	return Node{
		ID:       id,
		Children: children,
	}
}

// BFS implements Breadth-first search for Node
// Caveat: we use a buffered channel, so there is a small chance
// of deadlock!
func BFS(root Node, isGoal func(Node) bool) Node {
	queue := make(chan Node, 100) // TODO: not hardcode it
	visited := map[int]struct{}{} //int{}      // map[NodeID]distance

	// Add the root node and mark it as visited
	queue <- root
	visited[root.ID] = struct{}{}

	for len(queue) > 0 {
		current := <-queue

		if isGoal(current) {
			return current
		}

		for _, child := range current.Children {
			if _, isVisited := visited[child.ID]; !isVisited {
				visited[child.ID] = struct{}{} // visited[current.ID] + 1
				queue <- child
			}
		}
	}

	return Node{}
}

// TODO: Do I need to keep this?
func BFSDistances(root Node) map[int]int {
	queue := make(chan Node, 100) // TODO: not hardcode it
	distances := map[int]int{}    // map[NodeID]distance

	// Add the root node and mark it as visited
	queue <- root
	distances[root.ID] = 0

	for len(queue) > 0 {
		current := <-queue

		for _, child := range current.Children {
			if _, isVisited := distances[child.ID]; !isVisited {
				distances[child.ID] = distances[current.ID] + 1
				queue <- child
			}
		}
	}

	return distances
}
