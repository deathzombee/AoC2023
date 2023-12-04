package main

import (
	"AoC2023/utils"
	"fmt"
	"strconv"
	"unicode"
)

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func isSymbol(ch rune) bool {
	return !unicode.IsDigit(ch) && ch != '.'
}

func checkAdjacent(grid [][]rune, x, y int) (bool, rune) {
	rows := len(grid)
	if rows == 0 {
		fmt.Println("Grid is empty")
		return false, ' ' // or handle the error as appropriate
	}

	for _, row := range grid {
		if len(row) != rows {
			fmt.Printf("Inconsistent row length found: %d expected, %d found\n", rows, len(row))
			return false, ' ' // or handle the error as appropriate
		}
	}

	cols := len(grid[0])

	directions := []struct{ dx, dy int }{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for _, dir := range directions {
		newX, newY := x+dir.dx, y+dir.dy
		if newX >= 0 && newY >= 0 && newX < rows && newY < cols {
			if isSymbol(grid[newX][newY]) {
				return true, grid[newX][newY]
			}
		}
	}
	return false, ' ' // No adjacent symbol found
}

func extractFullNumber(grid [][]rune, x, y, dx, dy int) string {
	extract := func(dx, dy int) string {
		ix, iy := x, y
		number := ""
		for {
			ix += dx
			iy += dy
			if ix < 0 || iy < 0 || ix >= len(grid) || iy >= len(grid[0]) || !isDigit(grid[ix][iy]) {
				break
			}
			number += string(grid[ix][iy])
		}
		return number
	}

	// Reverse part1 and concatenate with current digit and part2
	part1 := extract(-dx, -dy) // Extract backward
	part2 := extract(dx, dy)   // Extract forward
	reversedPart1 := reverseString(part1)
	fullNumber := reversedPart1 + string(grid[x][y]) + part2
	return fullNumber
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func findAdjacentNumbers(input []string) []string {
	var grid [][]rune
	for _, line := range input {
		grid = append(grid, []rune(line))
	}

	processed := make(map[string]struct{})
	var fullNumbers []string

	for i, row := range grid {
		for j, ch := range row {
			if isDigit(ch) {
				adjacent, symbol := checkAdjacent(grid, i, j)
				if adjacent {
					// Extract full number horizontally and vertically
					fullNumberH := extractFullNumber(grid, i, j, 0, 1) // Horizontal
					fullNumberV := extractFullNumber(grid, i, j, 1, 0) // Vertical

					if _, ok := processed[fullNumberH]; !ok && fullNumberH != string(ch) {
						fmt.Printf("Adjacent number (Horizontal) at [%d,%d] to symbol '%c': %s\n", i, j, symbol, fullNumberH)
						fullNumbers = append(fullNumbers, fullNumberH)
						processed[fullNumberH] = struct{}{}
					}
					if _, ok := processed[fullNumberV]; !ok && fullNumberV != string(ch) {
						fmt.Printf("Adjacent number (Vertical) at [%d,%d] to symbol '%c': %s\n", i, j, symbol, fullNumberV)
						fullNumbers = append(fullNumbers, fullNumberV)
						processed[fullNumberV] = struct{}{}
					}
				}
			}
		}
	}
	return fullNumbers
}

func sumAdjacentNumbers(input []string) int {
	fullNumbers := findAdjacentNumbers(input)
	var sum int
	for _, num := range fullNumbers {
		n, err := strconv.Atoi(num)
		if err != nil {
			fmt.Printf("Error converting string to int: %v\n", err)
			return 0
		}
		sum += n
	}
	return sum
}
func main() {
	filename := "input.txt"

	lines, err := utils.ReadLines(filename)
	if err != nil {
		fmt.Printf("Error reading the file: %v\n", err)
		return
	}
	sum := sumAdjacentNumbers(lines)
	fmt.Printf("Sum of all adjacent numbers: %d\n", sum)
}
