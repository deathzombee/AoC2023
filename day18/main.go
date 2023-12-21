package main

import (
	"AoC2023/utils"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"
)

type point struct {
	X, Y int
}

type interval struct {
	Start, End point
}

func processInput(input []string) (string, string) {
	var directions []string
	var colors []string

	for _, line := range input {
		//parts := strings.Split(line, " ")
		ps := strings.Fields(line)
		prt1 := ps[0] + " " + ps[1]
		prt2 := ps[2]
		parts := []string{prt1, prt2}
		var directionPart string
		var colorPart string

		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part[0] == 'R' || part[0] == 'L' || part[0] == 'U' || part[0] == 'D' {
				directionPart += part + " "
				// add the integer part to the directions before the parenthesis
			} else if strings.HasPrefix(part, "(#") && strings.HasSuffix(part, ")") {
				colorPart = part
			}
		}
		if directionPart != "" {
			directions = append(directions, directionPart)
		}

		if colorPart != "" {
			// Extract the hexadecimal part from the colorPart
			colorHex := strings.TrimPrefix(strings.TrimSuffix(colorPart, ")"), "(#")

			// Extract the first 5 characters of the hexadecimal as a number
			colorDigits := colorHex[:5]
			colorNumber, err := strconv.ParseInt(colorDigits, 16, 64)
			if err != nil {
				fmt.Printf("Error parsing colorDigits: %v\n", err)
				continue
			}

			// Extract the last digit of the hexadecimal
			directionDigit := colorHex[5]

			// Map the directionDigit to a direction
			var mappedDirection string
			switch directionDigit {
			case '0':
				mappedDirection = "R"
			case '1':
				mappedDirection = "D"
			case '2':
				mappedDirection = "L"
			case '3':
				mappedDirection = "U"
			default:
				fmt.Printf("Invalid direction digit: %c\n", directionDigit)
				continue
			}

			// Add the mapped direction and number to the colors
			colors = append(colors, mappedDirection+" "+strconv.FormatInt(colorNumber, 10))
		}
	}

	// Join the colors slice into a single string
	colorsStr := strings.Join(colors, ", ")
	directionsStr := strings.Join(directions, ", ")

	return directionsStr, colorsStr
}
func processInstructions(instructions string) ([]interval, []point) {
	var allIntervals []interval
	x, y := 0, 0 // Starting point
	points := []point{{0, 0}}
	current := point{0, 0}

	for _, instruction := range strings.Split(instructions, ", ") {
		var dir string
		var steps int

		n, err := fmt.Sscanf(instruction, "%s %d", &dir, &steps)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error parsing instruction: %v\n", err)
			continue // Skip this instruction and continue with the next one
		}
		if n != 2 {
			fmt.Printf("Invalid instruction format: %s\n", instruction)
			continue
		}

		switch dir {
		case "R":
			newInterval := interval{Start: point{x, y}, End: point{x + steps, y}}
			allIntervals = append(allIntervals, newInterval)
			x += steps

			current = point{x, y}
		case "L":
			newInterval := interval{Start: point{x, y}, End: point{x - steps, y}}
			allIntervals = append(allIntervals, newInterval)
			x -= steps
			current = point{x, y}
		case "U":
			ynew := y - steps
			newInterval := interval{Start: point{x, y}, End: point{x, ynew}}
			allIntervals = append(allIntervals, newInterval)
			y = ynew
			current = point{x, y}
		case "D":
			ynew := y + steps
			newInterval := interval{Start: point{x, y}, End: point{x, ynew}}
			allIntervals = append(allIntervals, newInterval)
			y = ynew
			current = point{x, y}
		}

		points = append(points, current)
	}

	return allIntervals, points
}

func calculatePathPerimeter(points []point) int {
	perimeter := 0
	for i := 0; i < len(points)-1; i++ {
		perimeter += distance(points[i], points[i+1])
	}
	return perimeter
}

func calculateArea(intervals []interval) float64 {
	var vertices []point
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

func distance(a, b point) int {
	return int(math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y)))
}

func part1(input string) int {
	allIntervals, points := processInstructions(input)
	per := calculatePathPerimeter(points)
	area := calculateArea(allIntervals)
	picksArea := per/2 + int(area) + 1
	return picksArea
}
func part2(input string) int {
	allIntervals, points := processInstructions(input)
	per := calculatePathPerimeter(points)
	area := calculateArea(allIntervals)
	picksArea := per/2 + int(area) + 1
	return picksArea
}
func main() {
	input, _ := utils.ReadLines("input.txt")
	instr, cc := processInput(input)
	t := time.Now()
	p1 := part1(instr)
	fmt.Println("Part 1:", time.Since(t))
	fmt.Println("Area:", p1)
	t2 := time.Now()
	p2 := part2(cc)
	fmt.Println("Part 2:", time.Since(t2))
	fmt.Println("Area:", p2)
}
