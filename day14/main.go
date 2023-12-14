package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	emptySpace  = "."
	roundedRock = "O"
	cubeRock    = "#"
)

func ReadInputGrid(path string) [][]string {
	fmt.Println("Reading from input file...", path)

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var grid [][]string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		grid = append(grid, strings.Split(sc.Text(), ""))
	}

	return grid
}

type Coordinate struct {
	x, y int
}

func moveRocksNorth(grid [][]string) [][]string {
	moved := true
	for moved {
		moved = false
		for x := 0; x < len(grid[0]); x++ {
			for y := 0; y < len(grid); y++ {
				if grid[y][x] == roundedRock {
					highestPosition := findHighestPosition(grid, x, y)
					if highestPosition != y {
						grid[highestPosition][x] = roundedRock
						grid[y][x] = emptySpace
						moved = true
					}
				}
			}
		}
	}
	return grid
}

func findHighestPosition(grid [][]string, x, y int) int {
	for newY := y - 1; newY >= 0; newY-- {
		if grid[newY][x] != emptySpace {
			return newY + 1
		}
	}
	return 0
}

func CalculateLoad(grid [][]string) int {
	load := 0
	for y, row := range grid {
		for _, cell := range row {
			if cell == roundedRock {
				load += len(grid) - y
			}
		}
	}
	return load
}
func PrintGrid(grid [][]string) {
	for _, row := range grid {
		fmt.Println(strings.Join(row, ""))
	}
}

func main() {
	inputGrid := ReadInputGrid("input.txt")
	fmt.Println("inputGrid:")
	PrintGrid(inputGrid)
	fmt.Println("")
	fmt.Println("tiltedGrid:")
	tiltedGrid := moveRocksNorth(inputGrid)
	PrintGrid(tiltedGrid)
	load := CalculateLoad(tiltedGrid)
	fmt.Println("Total load on the north support beams:", load)
}
