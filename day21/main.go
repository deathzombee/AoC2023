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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

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

// copy the graph i times in the x and y direction
func expandGraph(g *graph, factor int) *graph {
	newWidth := g.width * factor
	newHeight := g.height * factor
	expandedVertices := make([][]vertex, newHeight)

	for y := 0; y < newHeight; y++ {
		expandedVertices[y] = make([]vertex, newWidth)
		for x := 0; x < newWidth; x++ {
			// Copy the vertex from the original graph, adjusting for the expansion factor
			originalVertex := g.vertices[y%g.height][x%g.width]
			expandedVertices[y][x] = vertex{position{x, y}, originalVertex.plot, -1}
		}
	}

	return &graph{expandedVertices, newWidth, newHeight}
}

func calculateCoefficients(points []int) (int, int, int) {
	a := (points[2] + points[0] - 2*points[1]) / 2
	b := points[1] - points[0] - a
	c := points[0]
	fmt.Println("p[0]", points[0], "p[1]", points[1], "p[2]", points[2])
	return a, b, c
}

func predictValue(a, b, c, n int) int {
	return a*n*n + b*n + c
}
func part2(g *graph, startPos position) int {
	factor := 5 // Expansion factor
	expandedGraph := expandGraph(g, factor)

	// Determine step values based on the size of the expanded graph
	size := len(g.vertices)
	fmt.Println("size", size)
	half := size / 2
	fmt.Println("half", half)
	steps := []int{half, half + size, half + 2*size}
	fmt.Println("startPos", startPos)
	startPos = position{(startPos.x * factor) + 2, (startPos.y * factor) + 2}
	fmt.Println("startPos", startPos)
	fmt.Println("len of expandedGraph", len(expandedGraph.vertices))

	// Gather data points

	dataPoints := make([]int, len(steps))
	for i, step := range steps {
		dataPoints[i] = bfs(expandedGraph, startPos, step)
		fmt.Printf("Steps: %d, Reachable Plots: %d\n", step, dataPoints[i])
	}

	// Calculate coefficients for polynomial regression
	a, b, c := calculateCoefficients(dataPoints)
	fmt.Println("Coefficients: ", a, b, c)

	// Predict the value for a large number of steps
	largeStepCount := (26501365 - startPos.x) / size // Example large step count that equals a grid of size 131 with starting position 65,65
	predictedValue := predictValue(a, b, c, largeStepCount)
	printGraph(expandedGraph)
	return predictedValue
}

// print expanded graph
func printGraph(g *graph) {
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			if g.vertices[y][x].plot {
				fmt.Print(".")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
}
func main() {
	inputFile := "input.txt"
	g, startPos, err := parseInput(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	t := time.Now()
	stepCount := DefaultStepCount
	result := bfs(g, startPos, stepCount)
	fmt.Println("Number of garden plots reachable in", stepCount, "steps:", result)
	fmt.Println("Time taken:", time.Since(t))
	t2 := time.Now()
	result2 := part2(g, startPos)
	fmt.Println("Number of garden plots reachable in 26501365 steps:", result2)
	fmt.Println("Time taken:", time.Since(t2))
}
