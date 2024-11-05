package main

import (
	"AoC2023/utils"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func parseLineToFloat64s(line string) ([]float64, error) {
	fields := strings.Fields(line)[1:]
	var floats []float64
	for _, field := range fields {
		value, err := strconv.Atoi(field)
		if err != nil {
			return nil, fmt.Errorf("error parsing float: %w", err)
		}
		floats = append(floats, float64(value))
	}
	return floats, nil
}
func concatenateLineElements(line string) (float64, error) {
	elements := strings.Fields(line)[1:]
	concatenated := strings.Join(elements, "")

	concatenatedInt, err := strconv.Atoi(concatenated)
	if err != nil {
		return 0, fmt.Errorf("error concatenating elements: %w", err)
	}

	return float64(concatenatedInt), nil
}

func part1(lines []string) (int, error) {
	raceTimes, err := parseLineToFloat64s(lines[0])
	if err != nil {
		return 0, err
	}

	raceDistances, err := parseLineToFloat64s(lines[1])
	if err != nil {
		return 0, err
	}

	var holdTimeProduct = 1
	for i := 0; i < len(raceTimes); i++ {
		holdTimeProduct *= holdTimeRange(raceTimes[i], raceDistances[i])
	}
	return holdTimeProduct, nil
}
func part2(lines []string) (int, error) {
	raceTime, err := concatenateLineElements(lines[0])
	if err != nil {
		return 0, err
	}
	raceDistance, err := concatenateLineElements(lines[1])
	if err != nil {
		return 0, err
	}

	return holdTimeRange(raceTime, raceDistance), nil
}

// use the quadratic formula to find the range of valid hold times
func holdTimeRange(totalTime float64, distance float64) int {
	var a float64 = 1
	var b = -totalTime
	var c = distance
	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return 0
	}
	root1 := (-b + math.Sqrt(discriminant)) / (2 * a)
	root2 := (-b - math.Sqrt(discriminant)) / (2 * a)
	// the hold times are the roots of the quadratic equation
	minHoldTime := math.Ceil(min(root1, root2))
	maxHoldTime := math.Floor(max(root1, root2))
	// verify that the hold times are valid (i.e. they don't exceed the total time)
	minHoldTime = max(minHoldTime, 0)
	maxHoldTime = min(maxHoldTime, totalTime)
	return int(maxHoldTime - minHoldTime + 1)

}

func main() {
	filename := "input.txt"
	lines, err := utils.ReadLines(filename)
	if err != nil {
		fmt.Printf("Error reading the file: %v\n", err)
		return
	}

	p1, err := part1(lines)
	if err != nil {
		fmt.Printf("Error parsing input: %v\n", err)
		return
	}

	p2, err := part2(lines)
	if err != nil {
		fmt.Printf("Error parsing input: %v\n", err)
		return
	}
	fmt.Printf("Part 1: %v\n", p1)
	fmt.Printf("Part 2: %v\n", p2)
}
