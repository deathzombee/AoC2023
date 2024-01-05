package main

import (
	"bufio"
	"fmt"
	"os"
)

type position struct {
	x, y int
}

type vertex struct {
	pos    position
	valid  bool
	char   rune
	steps  int
	parent *position
}

type graph struct {
	vertices      [][]vertex
	width, height int
}

func (v *vertex) isSlope() bool {

	return v.char == '^' || v.char == '>' || v.char == 'v' || v.char == '<'

}

func (v *vertex) slopeDirection() position {
	switch v.char {
	case '^':
		return position{0, -1}
	case '>':
		return position{1, 0}
	case 'v':
		return position{0, 1}
	case '<':
		return position{-1, 0}
	}
	return position{0, 0}
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
	var startPos position
	startPosFound := false
	y := 0

	for scanner.Scan() {
		line := scanner.Text()
		row := make([]vertex, len(line))
		for x, c := range line {
			valid := c == '.' || c == '^' || c == '>' || c == 'v' || c == '<'
			if c == '.' && !startPosFound {
				startPos = position{x, y}
				startPosFound = true
			}
			row[x] = vertex{position{x, y}, valid, c, -1, nil}
		}
		grid = append(grid, row)
		y++
	}

	if err := scanner.Err(); err != nil {
		return nil, position{}, err
	}

	if !startPosFound {
		return nil, position{}, fmt.Errorf("no starting position found")
	}

	return &graph{vertices: grid, width: len(grid[0]), height: len(grid)}, startPos, nil
}

func (g *graph) getNeighbors(pos position) []position {
	var directions []position

	currentTile := g.vertices[pos.y][pos.x]
	if currentTile.isSlope() {
		directions = []position{currentTile.slopeDirection()}
	} else {
		directions = []position{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	}

	neighbors := make([]position, 0, len(directions))
	for _, dir := range directions {
		newPos := position{pos.x + dir.x, pos.y + dir.y}
		if newPos.x >= 0 && newPos.x < g.width && newPos.y >= 0 && newPos.y < g.height && g.vertices[newPos.y][newPos.x].valid {
			neighbors = append(neighbors, newPos)
		}
	}
	return neighbors
}
func (g *graph) getNeighbors2(pos position) []position {
	var directions []position

	directions = []position{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	neighbors := make([]position, 0, len(directions))
	for _, dir := range directions {
		newPos := position{pos.x + dir.x, pos.y + dir.y}
		if newPos.x >= 0 && newPos.x < g.width && newPos.y >= 0 && newPos.y < g.height && g.vertices[newPos.y][newPos.x].valid {
			neighbors = append(neighbors, newPos)
		}
	}
	return neighbors
}
func dfs(g *graph, pos position, visited map[position]bool, steps int, maxSteps *int, furthestPos *position) {
	if visited[pos] {
		return
	}

	visited[pos] = true
	g.vertices[pos.y][pos.x].steps = steps

	if steps > *maxSteps {
		*maxSteps = steps
		*furthestPos = pos
	}

	for _, neighbor := range g.getNeighbors(pos) {
		if !visited[neighbor] {
			g.vertices[neighbor.y][neighbor.x].parent = &pos
			dfs(g, neighbor, visited, steps+1, maxSteps, furthestPos)
		}
	}

	visited[pos] = false
}

func part1(g *graph, startPos position) (int, []position) {
	maxSteps := 0
	visited := make(map[position]bool)
	furthestPos := startPos
	dfs(g, startPos, visited, 0, &maxSteps, &furthestPos)

	var path []position
	for pos := furthestPos; pos != startPos; pos = *g.vertices[pos.y][pos.x].parent {
		path = append([]position{pos}, path...)
	}
	path = append([]position{startPos}, path...)

	return maxSteps, path
}

func printGraph(g *graph, path []position) {
	pathMap := make(map[position]bool)
	for _, pos := range path {
		pathMap[pos] = true
	}

	for _, row := range g.vertices {
		for _, v := range row {
			if v.isSlope() {
				fmt.Printf("%c ", v.char)
			} else if pathMap[v.pos] {
				fmt.Print("O ") // Print the slope character
			} else if v.valid {
				fmt.Print(". ")
			} else {
				fmt.Print("# ")
			}
		}
		fmt.Println()
	}
}

func main() {
	g, startPos, err := parseInput("input.txt")
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	p1, _ := part1(g, startPos)
	fmt.Println("Longest hike with slopes considered:", p1)
	//printGraph(g, path)
}
