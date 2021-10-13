package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func JSONPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func gameMapToString(m BoardMap) string {
	b := &strings.Builder{}
	size := len(m)

	for y := size - 1; y >= 0; y-- {
		for x := 0; x < size; x++ {
			var token string
			switch val := m[x][y]; val {
			case SnakeHead:
				token = "H"
			case SnakeBody:
				token = "*"
			case SnakeTail:
				token = "T"
			case Food:
				token = "F"
			case LookAheadHead:
				token = "+"
			default:
				token = fmt.Sprintf("%5d", val)
			}
			fmt.Fprintf(b, "%5s", token)
		}
		fmt.Fprint(b, "\n")
	}

	return b.String()
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
