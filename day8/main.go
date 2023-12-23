package main

import (
	"AoC2023/utils"
	"fmt"
	"strings"
	"time"
)

func parseInput(input []string) (map[string][]string, []int) {
	m := make(map[string][]string)
	var dir []int
	directions := false

	for _, line := range input {
		switch {
		case line == "":
			directions = true
		case directions:
			p := strings.Split(line, " = ")
			name := p[0]
			split := strings.Split(strings.Trim(p[1], "()"), ", ")

			for _, neighbor := range split {
				m[name] = append(m[name], neighbor)
			}

		default:
			for _, d := range line {
				switch d {
				case 'L':
					dir = append(dir, 0)
				case 'R':
					dir = append(dir, 1)
				}
			}

		}
	}

	return m, dir
}
func navigate(m map[string][]string, dir []int) (string, int) {
	node := "AAA"
	i := 0
	c := 0

	for node != "ZZZ" {
		// Determine the next node based on the current direction
		n := m[node][dir[i]]
		node = n
		c++

		// Move to next direction, loop back if at end of dir
		i = (i + 1) % len(dir)
	}

	return node, c
}

func main() {
	input, _ := utils.ReadLines("input.txt")
	t := time.Now()
	m, directions := parseInput(input)
	node, directionCount := navigate(m, directions)
	fmt.Println(node, directionCount)
	fmt.Println("part1 time:", time.Since(t))
}
