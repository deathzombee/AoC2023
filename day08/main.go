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

func ghostNav(m map[string][]string, dir []int) ([]int, error) {
	forwardMap := make(map[string]string)
	zEarly := make(map[string]bool)

	// Building the fast forward map
	for s := range m {
		w := s
		for i, t := range dir {
			if i > 0 && strings.HasSuffix(w, "Z") {
				zEarly[s] = true
				break
			}
			w = m[w][t]
		}
		forwardMap[s] = w
	}

	// Start network: nodes ending with 'A'
	network := make(map[string]bool)
	for node := range m {
		if strings.HasSuffix(node, "A") {
			network[node] = true
		}
	}

	// Calculating factors
	var factors []int
	for start := range network {
		w := start
		j := 0
		t := map[string]int{}

		for {
			if _, f := t[w]; f {
				break // Exit if we've been here before
			}
			t[w] = j

			if strings.HasSuffix(w, "Z") && j > 0 {
				factors = append(factors, j)
				break
			}

			w = forwardMap[w]
			j++
		}
	}

	return factors, nil
}
func part1(m map[string][]string, dir []int) int {
	node, directionCount := navigate(m, dir)
	fmt.Println(node, directionCount)
	return directionCount
}
func part2(m map[string][]string, dir []int) int {
	directionCount2, _ := ghostNav(m, dir)
	return utils.LCM(directionCount2...) * len(dir)
}
func main() {
	input, _ := utils.ReadLines("input.txt")
	t := time.Now()
	m, directions := parseInput(input)
	p1 := part1(m, directions)
	fmt.Printf("Time p1: %v\n", time.Since(t))
	fmt.Printf("Part 1: %v\n", p1)
	t = time.Now()
	p2 := part2(m, directions)
	fmt.Printf("Time p2: %v\n", time.Since(t))
	fmt.Printf("Part 2: %v\n", p2)
}
