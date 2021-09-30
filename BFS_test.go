package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestBFSTree(t *testing.T) {
	n10 := NewNode(10)
	n9 := NewNode(9)
	n8 := NewNode(8)
	n7 := NewNode(7)
	n6 := NewNode(6, n9, n10)
	n5 := NewNode(5)
	n4 := NewNode(4, n7, n8)
	n3 := NewNode(3, n5, n6)
	n2 := NewNode(2)
	root := NewNode(1, n2, n3, n4)

	iteration := 1
	fn := func(n Node) bool {
		if n.ID != iteration {
			t.Fatalf("did not expect node %d on iteration %d", n.ID, iteration)
		}

		iteration++
		return false
	}

	BFS(root, fn)
}

func TestBFSGraph(t *testing.T) {
	j := NewNode(10)
	i := NewNode(9)
	h := NewNode(8)
	g := NewNode(7)
	f := NewNode(6, i, j)
	e := NewNode(5)
	d := NewNode(4, g, h, j)
	c := NewNode(3, e, f)
	b := NewNode(2)
	a := NewNode(1, b, c, g, d)

	iteration := 0
	expectedOrder := []Node{a, b, c, g, d, e, f, h, j, i}
	fn := func(n Node) bool {
		if n.ID != expectedOrder[iteration].ID {
			t.Fatalf("did not expect node %q on iteration %d", mapNode[n.ID], iteration+1)
		}

		iteration++
		return false
	}

	BFS(a, fn)
}

func TestBFSDistances(t *testing.T) {
	n10 := NewNode(10)
	n9 := NewNode(9)
	n8 := NewNode(8)
	n7 := NewNode(7)
	n6 := NewNode(6, n9, n10)
	n5 := NewNode(5)
	n4 := NewNode(4, n7, n8)
	n3 := NewNode(3, n5, n6)
	n2 := NewNode(2)
	root := NewNode(1, n2, n3, n4)

	want := map[int]int{
		mapNameToID["A"]: 0,
		mapNameToID["C"]: 1,
		mapNameToID["H"]: 2,
		mapNameToID["F"]: 2,
		mapNameToID["G"]: 2,
		mapNameToID["I"]: 3,
		mapNameToID["J"]: 3,
		mapNameToID["B"]: 1,
		mapNameToID["D"]: 1,
		mapNameToID["E"]: 2,
	}
	got := BFSDistances(root)

	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
		for k, v := range got {
			fmt.Println(mapNode[k], v)
		}
	}
}

// func TestShortestPath(t *testing.T) {
// 	j := NewNode(10)
// 	i := NewNode(9)
// 	h := NewNode(8)
// 	g := NewNode(7)
// 	f := NewNode(6, i, j)
// 	e := NewNode(5)
// 	d := NewNode(4, g, h)
// 	c := NewNode(3, e, f)
// 	b := NewNode(2)
// 	a := NewNode(1, b, c, d)

// 	distances := BFSDistances(a)
// 	got := shortestPath(a, h, distances)

// 	fmt.Print(got)
// }

var mapNameToID = map[string]int{
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

var mapNode = map[int]string{
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

func TestMakrSnakes(t *testing.T) {
	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{
				{
					Body: []Coord{{1, 0}, {1, 1}, {1, 2}},
				},
			},
		},
	}
	expected := GameMap{
		[]int{Unitialised, Unitialised, Unitialised}, // x = 0
		[]int{SnakeHead, SnakeBody, SnakeBody},       // x = 1
		[]int{Unitialised, Unitialised, Unitialised}, // x = 2
	}

	board := NewGameMap(3)

	MarkSnakes(board, state)

	if !reflect.DeepEqual(board, expected) {
		t.Errorf("expecting:\n%#v\ngot:\n%#v", expected, board)
	}
}

func TestCalculateDistanceFromFood(t *testing.T) {
	state := GameState{
		Board: Board{
			Width:  5,
			Height: 5,
			Snakes: []Battlesnake{
				{
					Body: []Coord{{1, 0}, {1, 1}, {1, 2}},
				},
			},
			Food: []Coord{{0, 0}, {4, 4}},
		},
	}

	// Visual representation is diferent from the data representation
	expected := GameMap{
		[]int{0, 1, 2, 3, 4},
		[]int{-3, -2, -2, 4, 3},
		[]int{6, 5, 4, 3, 2},
		[]int{5, 4, 3, 2, 1},
		[]int{4, 3, 2, 1, 0},
	}

	board := NewGameMap(5)
	oldBoard := printGameMap(board) // For debugging

	MarkSnakes(board, state)
	CalculateDistanceFromFood(board, state)

	if !reflect.DeepEqual(expected, board) {
		t.Errorf("original board:\n%s\ndid not expect:\n%s\n%s\n",
			oldBoard,
			printGameMap(board),
			printGameMap(expected))
	}
}

func JSONPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func printGameMap(m GameMap) string {
	b := &strings.Builder{}
	for x := len(m) - 1; x >= 0; x-- {
		for y := range m[x] {
			fmt.Fprintf(b, "%5d", m[x][y])
		}
		fmt.Fprintln(b, "")
	}

	return b.String()
}
