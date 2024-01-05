package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
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
	start         position
	end           position
}

type intersectionGraph map[position][][2]interface{}

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

var (
	directions = []position{{0, -1}, {0, 1}, {1, 0}, {-1, 0}}
)

func parseInput(filePath string) (*graph, position, position, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, position{}, position{}, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	scanner := bufio.NewScanner(file)
	grid := make([][]vertex, 0)
	var startPos, endPos position
	startPosFound := false
	y := 0

	for scanner.Scan() {
		line := scanner.Text()
		row := make([]vertex, len(line))
		for x, c := range line {
			valid := c == '.' || c == '^' || c == '>' || c == 'v' || c == '<'
			if c == '.' {
				endPos = position{x, y} // Update end position every time a '.' is found
				if !startPosFound {
					startPos = position{x, y}
					startPosFound = true
				}
			}
			row[x] = vertex{position{x, y}, valid, c, -1, nil}
		}
		grid = append(grid, row)
		y++
	}

	if err := scanner.Err(); err != nil {
		return nil, position{}, position{}, err
	}

	if !startPosFound {
		return nil, position{}, position{}, fmt.Errorf("no starting position found")
	}
	return &graph{vertices: grid, width: len(grid[0]), height: len(grid), start: startPos, end: endPos}, startPos, endPos, nil
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

func gDist(cur position, dist int, seen map[position]bool, intersections []position, neighbors map[position][]position) (position, int) {
	for _, p := range intersections {
		if p == cur {
			return cur, dist
		}
	}

	for _, nb := range neighbors[cur] {
		if !seen[nb] {
			seen[cur] = true
			return gDist(nb, dist+1, seen, intersections, neighbors)
		}
	}

	return position{}, 0
}

func bfs(start, end position, score int, seen map[position]bool, gr intersectionGraph) []int {
	if start == end {
		return []int{score}
	}

	var scores []int

	for _, pair := range gr[start] {

		current := pair[0].(position)
		dist := pair[1].(int)
		if !seen[current] {
			seen[current] = true
			scores = append(scores, bfs(current, end, score+dist, seen, gr)...)
			delete(seen, current)
		}
	}

	return scores
}

func Part2(input *graph) int {
	intersections, neighbors := neighborsIntersections(input)
	gr := gGraph(intersections, neighbors)
	vals := bfs(input.start, input.end, 0, map[position]bool{input.start: true}, gr)

	return gMax(vals)
}

func gGraph(intersections []position, neighbors map[position][]position) intersectionGraph {
	gr := make(intersectionGraph)

	for _, i := range intersections {
		for _, n := range neighbors[i] {
			t, d := gDist(n, 1, map[position]bool{i: true}, intersections, neighbors)
			gr[i] = append(gr[i], [2]interface{}{t, d})
		}
	}
	return gr
}

func neighborsIntersections(g *graph) ([]position, map[position][]position) {
	vtxs := g.vertices
	intersections := []position{g.start, g.end}

	neighbors := make(map[position][]position)

	// Corrected loops to iterate over 2D slice
	for y, row := range vtxs {
		for x, vtx := range row {
			if vtx.valid {
				currentPos := position{x, y}
				neighbors[currentPos] = []position{}
				for _, dir := range directions {
					nb := position{x + dir.x, y + dir.y}
					if nb.x >= 0 && nb.x < g.width && nb.y >= 0 && nb.y < g.height && vtxs[nb.y][nb.x].valid {
						neighbors[currentPos] = append(neighbors[currentPos], nb)
					}
				}
				if len(neighbors[currentPos]) > 2 {
					intersections = append(intersections, currentPos)
				}
			}
		}
	}
	return intersections, neighbors
}

func gMax(result []int) int {
	s := result[0]
	for _, score := range result {
		if score > s {
			s = score
		}
	}
	return s
}

func printGraph(g *graph, path []position) {
	m := make(map[position]bool)
	for _, pos := range path {
		m[pos] = true
	}

	for _, row := range g.vertices {
		for _, v := range row {
			if v.isSlope() {
				fmt.Printf("%c ", v.char)
			} else if m[v.pos] {
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
	g, startPos, _, err := parseInput("input.txt")
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	t1 := time.Now()
	p1, _ := part1(g, startPos)
	fmt.Println("Longest hike with slopes considered:", p1, time.Since(t1))
	//printGraph(g, path)
	t2 := time.Now()
	p2 := Part2(g)
	fmt.Println("Longest hike with slopes not considered:", p2, time.Since(t2))
}
