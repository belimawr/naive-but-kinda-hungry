package main

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

// func shortestPath(src, dest Node, distances map[int]int) []Node {
// 	path := []Node{}

// 	current := src
// 	for distances[current.ID] != 0 {
// 		fmt.Printf("Distance: %s -> %d\n", mapNode2[current.ID], distances[current.ID])
// 		fmt.Println("Condition: ", distances[current.ID], " != 0")
// 		fmt.Println("Current: ", current)

// 		path = append(path, current)
// 		fmt.Println("Path: ", path)

// 		adjacent := current.Children[0]

// 		childIdx := 0
// 		for childIdx < len(current.Children) {
// 			fmt.Println("distances[adjacent.ID]", distances[adjacent.ID])
// 			if distances[adjacent.ID] == distances[current.ID]+1 {
// 				current = adjacent
// 				childIdx = 0
// 				break
// 			}
// 		}
// 	}

// 	return path
// }

var mapNameToID2 = map[string]int{
	"A": 1,
	"B": 2,
	"C": 3,
	"D": 4,
	"E": 5,
	"F": 6,
	"G": 7,
	"H": 8,
	"I": 9,
	"J": 10,
}

var mapNode2 = map[int]string{
	1:  "A",
	2:  "B",
	3:  "C",
	4:  "D",
	5:  "E",
	6:  "F",
	7:  "G",
	8:  "H",
	9:  "I",
	10: "J",
}

type GameMap [][]int

// const Unitialised = math.MaxInt
// const SnakeBody = Unitialised - 1
// const SnakeHead = SnakeBody - 1
// const Food = SnakeHead - 1

const Unitialised = -1
const SnakeBody = -2
const SnakeHead = -3
const Food = -4

// NewGameMap returns a initialised NewGameMap
// of the desired size
func NewGameMap(size int) GameMap {
	board := make(GameMap, size)
	for i := range board {
		board[i] = make([]int, size)
		for j := range board[i] {
			board[i][j] = Unitialised
		}
	}

	return board
}

// MarkSnakes marks all snakes into the board
func MarkSnakes(board GameMap, state GameState) {
	for _, snake := range state.Board.Snakes {
		for i, p := range snake.Body {
			if i == 0 { // TODO: Optmise it, if necessary
				board[p.X][p.Y] = SnakeHead
				continue
			}
			board[p.X][p.Y] = SnakeBody
		}
	}
}

func CalculateDistanceFromFood(board GameMap, state GameState) {
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
			// if adjacent has not benn processed or it's not a snake
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
