package main

import (
	"AoC2023/utils"
	"fmt"
	"time"
)

const (
	UP = iota
	RIGHT
	DOWN
	LEFT
	UPDOWN
	LEFTRIGHT
)

// Pos struct
type Pos struct {
	X, Y int
}

// nextDir function optimized
func nextDir(dir int, c uint8) int {
	switch dir {
	case RIGHT:
		if c == '.' || c == '-' {
			return RIGHT
		}
		if c == '|' {
			return UPDOWN
		}
		if c == '/' {
			return UP
		}
		return DOWN
	case LEFT:
		if c == '.' || c == '-' {
			return LEFT
		}
		if c == '|' {
			return UPDOWN
		}
		if c == '/' {
			return DOWN
		}
		return UP
	case UP:
		if c == '.' || c == '|' {
			return UP
		}
		if c == '-' {
			return LEFTRIGHT
		}
		if c == '/' {
			return RIGHT
		}
		return LEFT
	case DOWN:
		if c == '.' || c == '|' {
			return DOWN
		}
		if c == '-' {
			return LEFTRIGHT
		}
		if c == '/' {
			return LEFT
		}
		return RIGHT
	}
	panic("invalid state")
}

type state struct {
	pos Pos
	dir int
}

// Combined solve and isValidPos functions
func solve(grid [][]uint8, current state) int {
	var todo = make([]state, 0, 100) // Pre-allocate based on expected size
	visited := make(map[state]bool)
	energized := make(map[Pos]bool)

	todo = append(todo, current)
	for len(todo) > 0 {
		s := todo[0]
		todo = todo[1:]

		if visited[s] || !isValid(s.pos, grid) {
			continue
		}
		visited[s] = true
		energized[s.pos] = true

		x, y := s.pos.X, s.pos.Y
		switch nextDir(s.dir, grid[y][x]) {
		case UP:
			todo = append(todo, state{pos: Pos{X: x, Y: y - 1}, dir: UP})
		case RIGHT:
			todo = append(todo, state{pos: Pos{X: x + 1, Y: y}, dir: RIGHT})
		case DOWN:
			todo = append(todo, state{pos: Pos{X: x, Y: y + 1}, dir: DOWN})
		case LEFT:
			todo = append(todo, state{pos: Pos{X: x - 1, Y: y}, dir: LEFT})
		case UPDOWN:
			todo = append(todo, state{pos: Pos{X: x, Y: y - 1}, dir: UP})
			todo = append(todo, state{pos: Pos{X: x, Y: y + 1}, dir: DOWN})
		case LEFTRIGHT:
			todo = append(todo, state{pos: Pos{X: x - 1, Y: y}, dir: LEFT})
			todo = append(todo, state{pos: Pos{X: x + 1, Y: y}, dir: RIGHT})
		}
	}
	return len(energized)
}

func isValid(pos Pos, grid [][]uint8) bool {
	return pos.Y >= 0 && pos.Y < len(grid) && pos.X >= 0 && pos.X < len(grid[pos.Y])
}

func Part1(grid [][]uint8) int {
	return solve(grid, state{pos: Pos{X: 0, Y: 0}, dir: RIGHT})
}

func Part2(grid [][]uint8) int {

	var res int
	maxX, maxY := len(grid[0])-1, len(grid)-1
	for x := 0; x <= maxX; x++ {
		res = max(res, solve(grid, state{pos: Pos{X: x, Y: 0}, dir: DOWN}))
		res = max(res, solve(grid, state{pos: Pos{X: x, Y: maxY}, dir: UP}))
	}
	for y := 0; y <= maxY; y++ {
		res = max(res, solve(grid, state{pos: Pos{X: 0, Y: y}, dir: RIGHT}))
		res = max(res, solve(grid, state{pos: Pos{X: maxX, Y: y}, dir: LEFT}))
	}

	return res
}

func buildMatrixCharFromString(input []string) [][]uint8 {
	grid := make([][]uint8, len(input))
	for i, line := range input {
		grid[i] = make([]uint8, len(line))
		for j, c := range line {
			grid[i][j] = uint8(c)
		}
	}
	return grid
}

func main() {
	input, _ := utils.ReadLines("input.txt")
	grid := buildMatrixCharFromString(input)
	fmt.Println("--2023 day 16 solution--")
	start := time.Now()
	fmt.Println("Part1: ", Part1(grid))
	fmt.Println("Elapsed: ", time.Since(start))

	start = time.Now()
	fmt.Println("Part2: ", Part2(grid))
	fmt.Println("Elapsed: ", time.Since(start))
}
