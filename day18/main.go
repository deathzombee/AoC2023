package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type Point struct {
	X, Y int
}

type Interval struct {
	Start, End Point
}

func processInstructions(instructions string) ([]Interval, []Point, int, int, int, int) {
	var allIntervals []Interval
	x, y := 0, 0 // Starting point
	points := []Point{{0, 0}}
	current := Point{0, 0}
	minX, minY, maxX, maxY := 0, 0, 0, 0

	for _, instruction := range strings.Split(instructions, ", ") {
		var dir string
		var steps int
		fmt.Sscanf(instruction, "%s %d", &dir, &steps)

		switch dir {
		case "R":
			newInterval := Interval{Start: Point{x, y}, End: Point{x + steps, y}}
			allIntervals = append(allIntervals, newInterval)
			x += steps
			current = Point{x, y}
			if x > maxX {
				maxX = x
			}
		case "L":
			newInterval := Interval{Start: Point{x, y}, End: Point{x - steps, y}}
			allIntervals = append(allIntervals, newInterval)
			x -= steps
			current = Point{x, y}
			if x < minX {
				minX = x
			}
		case "U":
			ynew := y - steps
			newInterval := Interval{Start: Point{x, y}, End: Point{x, ynew}}
			allIntervals = append(allIntervals, newInterval)
			y = ynew
			current = Point{x, y}
			if ynew < minY {
				minY = ynew
			}
		case "D":
			ynew := y + steps
			newInterval := Interval{Start: Point{x, y}, End: Point{x, ynew}}
			allIntervals = append(allIntervals, newInterval)
			y = ynew
			current = Point{x, y}
			if ynew > maxY {
				maxY = ynew
			}
		}
		points = append(points, current)
	}

	return allIntervals, points, minX, minY, maxX, maxY
}

func calculatePathPerimeter(points []Point) int {
	perimeter := 0
	for i := 0; i < len(points)-1; i++ {
		perimeter += distance(points[i], points[i+1])
	}
	return perimeter
}

func distance(a, b Point) int {
	return int(math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y)))
}
func calculateArea(intervals []Interval) float64 {
	var vertices []Point
	for _, interval := range intervals {
		if len(vertices) == 0 || vertices[len(vertices)-1] != interval.Start {
			vertices = append(vertices, interval.Start)
		}
		vertices = append(vertices, interval.End)
	}

	n := len(vertices)
	if n < 3 { // A polygon must have at least 3 vertices
		return 0
	}

	var area float64
	for i := 0; i < n-1; i++ {
		area += (float64(vertices[i].X) * float64(vertices[i+1].Y)) - (float64(vertices[i+1].X) * float64(vertices[i].Y))
	}
	// Closing the polygon
	area += (float64(vertices[n-1].X) * float64(vertices[0].Y)) - (float64(vertices[0].X) * float64(vertices[n-1].Y))

	return math.Abs(area) / 2
}
func main() {
	// Example instructions
	instructions := "R 6, D 5, L 2, D 2, R 2, D 2, L 5, U 2, L 1, U 2, R 2, U 3, L 2, U 2"
	t := time.Now()
	allIntervals, points, _, _, _, _ := processInstructions(instructions)
	per := calculatePathPerimeter(points)
	fmt.Println("Run time:", time.Since(t))
	fmt.Printf("Perimeter of the path is: %d\n", per)
	t2 := time.Now()
	area := calculateArea(allIntervals)
	total := per/2 + int(area) + 1
	fmt.Println("Run time:", time.Since(t2))
	fmt.Printf("Total is: %d\n", total)
	inside := total - per
	fmt.Printf("Filled is: %d\n", inside)
}
