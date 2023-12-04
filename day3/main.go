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

func checkAdjacent(grid [][]rune, x, y int) bool {
	directions := []struct {
		dx, dy int
	}{
		{-1, 0}, {1, 0},
		{0, -1}, {0, 1},
		{-1, -1}, {-1, 1},
		{1, -1}, {1, 1},
	}

	for _, d := range directions {
		newX, newY := x+d.dx, y+d.dy
		if newX >= 0 && newY >= 0 && newX < len(grid) && newY < len(grid[0]) {
			if isSymbol(grid[newX][newY]) {
				return true
			}
		}
	}
	return false
}

// Extracts the full number around position (x, y)

func extractFullNumber(grid [][]rune, x, y int) (string, int) {
	start := y
	for start > 0 && isDigit(grid[x][start-1]) {
		start--
	}

	end := y
	for end < len(grid[x])-1 && isDigit(grid[x][end+1]) {
		end++
	}

	return string(grid[x][start : end+1]), start
}

func sumAdjacentNumbers(input []string) int {
	var grid [][]rune
	for _, line := range input {
		grid = append(grid, []rune(line))
	}

	processed := make(map[string]bool)
	var sum int

	for x, row := range grid {
		for y, ch := range row {
			if isDigit(ch) && checkAdjacent(grid, x, y) {
				fullNumber, startPos := extractFullNumber(grid, x, y)
				numberKey := fmt.Sprintf("%s-%d-%d", fullNumber, x, startPos)

				if !processed[numberKey] {
					n, err := strconv.Atoi(fullNumber)
					if err != nil {
						fmt.Printf("Error converting string to int: %v\n", err)
						continue
					}
					fmt.Printf("Found number: %d at position [%d,%d]\n", n, x, startPos)
					sum += n

					// Mark this full number with its starting position as processed
					processed[numberKey] = true
				}
			}
		}
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
