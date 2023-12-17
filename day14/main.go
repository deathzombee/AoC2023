package main

import (
	"AoC2023/utils"
	"errors"
	"fmt"
	"time"
)

type Direction func([][]int) [][]int

func readGrid(s []string) (grid [][]int, err error) {
	grid = make([][]int, len(s))
	for y, line := range s {
		grid[y] = make([]int, len(line))
		for x := 0; x < len(line); x++ {
			switch line[x] {
			case 'O':
				grid[y][x] = 2
			case '#':
				grid[y][x] = 1
			case '.':
				grid[y][x] = 0
			default:
				return grid, errors.New("invalid input")
			}
		}
	}
	return grid, err
}
func north(grid [][]int) [][]int {
	for x := range grid[0] {
		stopPoint := 0
		for y, value := range grid {
			if value[x] == 1 {
				stopPoint = y + 1
			} else if value[x] == 2 {
				if y != stopPoint {
					grid[stopPoint][x], grid[y][x] = grid[y][x], 0
				}
				stopPoint++
			}
		}
	}
	return grid
}
func south(grid [][]int) [][]int {
	for x := range grid[0] {
		stopPoint := len(grid) - 1
		for y := len(grid) - 1; y >= 0; y-- {
			if grid[y][x] == 1 {
				stopPoint = y - 1
			} else if grid[y][x] == 2 {
				if y != stopPoint {
					grid[stopPoint][x], grid[y][x] = grid[y][x], 0
				}
				stopPoint--
			}
		}
	}
	return grid
}
func west(grid [][]int) [][]int {
	for y := range grid {
		stopPoint := 0
		for x, value := range grid[y] {
			if value == 1 {
				stopPoint = x + 1
			} else if value == 2 {
				if x != stopPoint {
					grid[y][stopPoint], grid[y][x] = grid[y][x], 0
				}
				stopPoint++
			}
		}
	}
	return grid
}
func east(grid [][]int) [][]int {
	for y := range grid {
		stopPoint := len(grid[0]) - 1
		for x := len(grid[0]) - 1; x >= 0; x-- {
			if grid[y][x] == 1 {
				stopPoint = x - 1
			} else if grid[y][x] == 2 {
				if x != stopPoint {
					grid[y][stopPoint], grid[y][x] = grid[y][x], 0
				}
				stopPoint--
			}
		}
	}
	return grid
}

func sequence(values []int) (frequency, offset int, ok bool) {

	if len(values) < 10 {
		return 0, 0, false
	}

	for offset = 0; offset < len(values); offset++ {
		for frequency = 1; offset+2*frequency < len(values); frequency++ {
			isRepeating := true
			for i := 0; i < frequency; i++ {
				if offset+2*frequency+i >= len(values) ||
					values[offset+i] != values[offset+frequency+i] ||
					values[offset+i] != values[offset+2*frequency+i] {
					isRepeating = false
					break
				}
			}
			if isRepeating {
				return frequency, offset, true
			}
		}
	}

	return 0, 0, false
}

func loadCalc(s []string, transform Direction) (int, error) {
	var result int
	if grid, err := readGrid(s); err != nil {
		return 0, err
	} else {
		grid = transform(grid)
		var numRocks int
		for y := 0; y < len(grid); y++ {
			numRocks = 0
			for x := 0; x < len(grid[y]); x++ {
				if grid[y][x] == 2 {
					numRocks++
				}
			}
			result += numRocks * (len(grid) - y)
		}
	}
	return result, nil
}

func cycledLoad(s []string, numCycles int) (int, error) {
	var result int

	if grid, err := readGrid(s); err != nil {
		return 0, err
	} else {
		var results []int

		for {
			grid = north(grid)
			grid = west(grid)
			grid = south(grid)
			grid = east(grid)

			numRocks := 0
			result = 0

			for y := 0; y < len(grid); y++ {
				numRocks = 0
				for x := 0; x < len(grid[y]); x++ {
					if grid[y][x] > 1 {
						numRocks++
					}
				}
				result += numRocks * (len(grid) - y)
			}
			results = append(results, result)
			if frequency, offset, ok := sequence(results); ok {
				mod := (numCycles - offset) % frequency
				result = results[offset+mod-1]
				break
			}
		}
	}
	return result, nil
}
func main() {

	output, _ := utils.ReadLines("input.txt")
	time1 := time.Now()
	part1, _ := loadCalc(output, north)
	fmt.Println("p1 time:", time.Since(time1))
	fmt.Println("part 1:", part1)
	time2 := time.Now()
	part2, _ := cycledLoad(output, 1000000000)
	fmt.Println("p2 time:", time.Since(time2))
	fmt.Println("part 2:", part2)

}
