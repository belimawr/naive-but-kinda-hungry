package main

import (
	"fmt"
	"reflect"
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
func TestMakrSnakes(t *testing.T) {
	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{
				{
					Body: []Coord{{1, 0}, {1, 1}, {1, 2}},
					Head: Coord{1, 0},
				},
			},
		},
	}
	expected := BoardMap{
		[]int{Unitialised, Unitialised, Unitialised}, // x = 0
		[]int{SnakeHead, SnakeBody, SnakeTail},       // x = 1
		[]int{Unitialised, Unitialised, Unitialised}, // x = 2
	}

	board := NewBoardMap(3)

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
					Body: []Coord{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}},
					Head: Coord{1, 0},
				},

				{
					Body: []Coord{{X: 4, Y: 2}, {X: 4, Y: 1}, {X: 4, Y: 0}},
					Head: Coord{4, 2},
				},
			},
			Food: []Coord{{0, 0}, {4, 4}},
		},
	}

	// Visual representation is diferent from the data representation
	expected := BoardMap{
		[]int{0, 1, 2, 3, 4},    // X = 0
		[]int{-3, -2, -4, 4, 3}, // X = 1
		[]int{6, 5, 4, 3, 2},    // X = 2
		[]int{5, 4, 3, 2, 1},    // X = 3
		[]int{-4, -2, -3, 1, 0}, // X = 4
	}

	board := NewBoardMap(5)

	MarkSnakes(board, state)

	CalculateDistanceFromFood(board, state)

	if !reflect.DeepEqual(expected, board) {
		t.Errorf("expecting:\n%s\ngot:\n%s\n",
			gameMapToString(expected),
			gameMapToString(board))
	}
	fmt.Println(gameMapToString(board))
}
