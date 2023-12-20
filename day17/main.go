package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"time"
)

const (
	planeVertical = iota
	planeHorizontal
	planeUndecided // special plane for start position
	infinity       = 1 << 30
)

type graph struct {
	vertices      []vertex
	width, height int
}

type position struct {
	x, y int
}

type vertex struct {
	pos            position
	dir            int
	visited        bool
	entropy        int
	ajustedEntropy int
	totalEntropy   int
	heapIndex      int
}

type priorityQueue []*vertex

// priorityQueue implements heap.Interface and holds Vertices.
// The items are ordered by totalEntropy, the lower the better.

// Len returns the length of the priorityQueue
func (pq *priorityQueue) Len() int { return len(*pq) }

// Less returns true if the item at index i is less than the item at index j
func (pq *priorityQueue) Less(i, j int) bool {
	return (*pq)[i].totalEntropy < (*pq)[j].totalEntropy
}

// Swap swaps the items at index i and j
func (pq *priorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].heapIndex = i
	(*pq)[j].heapIndex = j
}

// Push pushes an item to the priorityQueue
func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*vertex)
	item.heapIndex = n
	*pq = append(*pq, item)
}

// Pop pops an item from the priorityQueue
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	//  avoid memory leak by clearing the reference
	old[n-1] = nil
	// the index should be decremented to avoid out of bounds
	item.heapIndex = -1
	*pq = old[0 : n-1]
	return item
}

// Update modifies the priority and value of an Item in the queue.
func (pq *priorityQueue) Update(item *vertex) {
	heap.Fix(pq, item.heapIndex)
}

func (g *graph) getEdges(u *vertex, minSteps int, maxSteps int) []*vertex {
	// there are at most 6 edges (2 for each direction) and they are the vertices that are at most maxSteps away
	edges := make([]*vertex, 0, 6)

	if u.dir == planeHorizontal || u.dir == planeUndecided {
		for entropy, dy := 0, 1; dy <= maxSteps; dy++ {
			v := g.getVertex(u.pos.x, u.pos.y+dy, planeVertical)
			if v != nil {
				entropy += v.entropy
				if dy >= minSteps {
					v.ajustedEntropy = entropy
					edges = append(edges, v)
				}
			}
		}
		for entropy, dy := 0, 1; dy <= maxSteps; dy++ {
			v := g.getVertex(u.pos.x, u.pos.y-dy, planeVertical)
			if v != nil {
				entropy += v.entropy
				if dy >= minSteps {
					v.ajustedEntropy = entropy
					edges = append(edges, v)
				}
			}
		}
	}

	if u.dir == planeVertical || u.dir == planeUndecided {
		for entropy, dx := 0, 1; dx <= maxSteps; dx++ {
			v := g.getVertex(u.pos.x+dx, u.pos.y, planeHorizontal)
			if v != nil {
				entropy += v.entropy
				if dx >= minSteps {
					v.ajustedEntropy = entropy
					edges = append(edges, v)
				}
			}
		}
		for entropy, dx := 0, 1; dx <= maxSteps; dx++ {
			v := g.getVertex(u.pos.x-dx, u.pos.y, planeHorizontal)
			if v != nil {
				entropy += v.entropy
				if dx >= minSteps {
					v.ajustedEntropy = entropy
					edges = append(edges, v)
				}
			}
		}
	}

	return edges
}

func (g *graph) getVertex(x int, y int, plane int) *vertex {
	if x < 0 || y < 0 || y >= g.height || x >= g.width {
		return nil
	}
	//for a normal graph, the index would be y*g.width+x
	// ours is double wide to accommodate both planes
	// if (1,1) in a 2x2 matrix would be index 3, in our case it's 6
	//if it's in the horizontal plane it's 7
	return &g.vertices[y*2*g.width+x*2+plane]
}

func parseInput(input []byte) [][]int {
	input = bytes.TrimSpace(input)
	lines := bytes.Split(input, []byte("\n"))
	grid := make([][]int, len(lines))
	for i, line := range lines {
		grid[i] = make([]int, len(line))
		for j, c := range line {
			grid[i][j] = int(c) - '0'
		}
	}
	return grid
}

func graphFromGrid(grid [][]int) graph {
	graph := graph{}
	vertices := make([]vertex, 0, len(grid)*len(grid[0])*2)
	graph.height = len(grid)
	for y := range grid {
		graph.width = len(grid[y])
		for x := range grid[y] {
			vertices = append(vertices, vertex{
				pos:          position{x: x, y: y},
				dir:          planeVertical,
				entropy:      grid[y][x],
				totalEntropy: infinity,
			})
			vertices = append(vertices, vertex{
				pos:          position{x: x, y: y},
				dir:          planeHorizontal,
				entropy:      grid[y][x],
				totalEntropy: infinity,
			})
		}
	}
	graph.vertices = vertices
	return graph
}

// this is a modified version of dijkstra's algorithm that finds the lowest entropy path
func dijkstra(grid [][]int, minSteps int, maxSteps int) int {
	graph := graphFromGrid(grid)
	vertices := graph.vertices
	vertices[0].totalEntropy = 0
	vertices[0].dir = planeUndecided
	pq := make(priorityQueue, len(vertices))
	for i := 0; i < len(vertices); i++ {
		vertices[i].heapIndex = i
		pq[i] = &vertices[i]
	}
	heap.Init(&pq)
	var u *vertex
	var edges = &vertices[len(vertices)-1]
	for {
		u = heap.Pop(&pq).(*vertex)
		if u.pos.x == edges.pos.x && u.pos.y == edges.pos.y {
			break
		}
		u.visited = true
		for _, edges := range graph.getEdges(u, minSteps, maxSteps) {
			if u.totalEntropy+edges.ajustedEntropy < edges.totalEntropy {
				edges.totalEntropy = u.totalEntropy + edges.ajustedEntropy
				pq.Update(edges)
			}
		}
	}
	return u.totalEntropy
}

func part1(grid [][]int) int {
	return dijkstra(grid, 1, 3)
}
func part2(grid [][]int) int {
	return dijkstra(grid, 4, 10)
}

func main() {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	t := time.Now()
	grid := parseInput(input)
	fmt.Println("Part 1:", part1(grid), "run time:", time.Since(t))
	t = time.Now()
	fmt.Println("Part 2:", part2(grid), "run time:", time.Since(t))
}
