package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const DefaultStepCount = 64

type position struct {
	x, y int
}

type vertex struct {
	pos   position
	plot  bool // true for garden plot, false for rock
	steps int
}

type graph struct {
	vertices      [][]vertex
	width, height int
}

func parseInput(filePath string) (*graph, position, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, position{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	grid := make([][]vertex, 0)
	startPos := position{}
	y := 0

	for scanner.Scan() {
		line := scanner.Text()
		row := make([]vertex, len(line))
		for x, c := range line {
			plot := c == '.' || c == 'S'
			if c == 'S' {
				startPos = position{x, y}
			}
			row[x] = vertex{position{x, y}, plot, -1}
		}
		grid = append(grid, row)
		y++
	}

	if err := scanner.Err(); err != nil {
		return nil, position{}, err
	}

	return &graph{grid, len(grid[0]), len(grid)}, startPos, nil
}

func (g *graph) getNeighbors(pos position) []position {
	directions := []position{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	neighbors := make([]position, 0, 4)

	for _, dir := range directions {
		newPos := position{pos.x + dir.x, pos.y + dir.y}
		if newPos.x >= 0 && newPos.x < g.width && newPos.y >= 0 && newPos.y < g.height && g.vertices[newPos.y][newPos.x].plot {
			neighbors = append(neighbors, newPos)
		}
	}
	return neighbors
}
func bfs(g *graph, start position, steps int) int {
	current := make(map[position]struct{})
	next := make(map[position]struct{})

	// Add the start position to the current set
	current[start] = struct{}{}

	for i := 0; i < steps; i++ {
		// Clear the next set for the next iteration
		for k := range next {
			delete(next, k)
		}

		for pos := range current {
			for _, neighbor := range g.getNeighbors(pos) {
				// Check if the neighbor has already been scheduled for the next step
				if _, found := next[neighbor]; !found {
					next[neighbor] = struct{}{}
				}
			}
		}

		// Swap current and next for the next iteration
		current, next = next, current
	}

	// 'current' now holds the unique positions reached exactly at 'steps' steps
	return len(current)
}

func main() {
	inputFile := "input.txt"
	g, startPos, err := parseInput(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	t := time.Now()
	stepCount := DefaultStepCount // Change this variable to set a custom step count
	result := bfs(g, startPos, stepCount)
	fmt.Println("Number of garden plots reachable in", stepCount, "steps:", result)
	fmt.Println("Time taken:", time.Since(t))
}
