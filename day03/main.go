package main

import (
	"AoC2023/utils"
	"fmt"
	"strconv"
	"unicode"
)

// isDigit checks if a rune is a digit.
func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

// isSymbol checks if a rune is a symbol (not a digit and not '.').
func isSymbol(ch rune) bool {
	return !unicode.IsDigit(ch) && ch != '.'
}

// extractFullNumber extracts the full number around a given position in the grid.
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

// checkAdjacent checks if a digit has a symbol in adjacent cells.
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

// sumAdjacentNumbers calculates the sum of adjacent numbers in the grid.
func sumAdjacentNumbers(grid [][]rune) int {
	processed := make(map[string]struct{})
	var sum int

	for x, row := range grid {
		for y, ch := range row {
			if isDigit(ch) && checkAdjacent(grid, x, y) {
				fullNumber, startPos := extractFullNumber(grid, x, y)
				numberKey := fmt.Sprintf("%s-%d-%d", fullNumber, x, startPos)

				if _, exists := processed[numberKey]; !exists {
					n, err := strconv.Atoi(fullNumber)
					if err != nil {
						fmt.Printf("Error converting string to int: %v\n", err)
						continue
					}
					sum += n
					processed[numberKey] = struct{}{}
				}
			}
		}
	}

	return sum
}

// extractNumbersAndPositions extracts numbers and their positions from the grid.
func extractNumbersAndPositions(grid [][]rune) map[string][][2]int {
	numPositions := make(map[string][][2]int)
	for x, row := range grid {
		for y, ch := range row {
			if isDigit(ch) {
				numStr, _ := extractFullNumber(grid, x, y)
				numPositions[numStr] = append(numPositions[numStr], [2]int{x, y})
			}
		}
	}
	return numPositions
}

// findStarPositions finds positions of '*' in the grid.
func findStarPositions(grid [][]rune) [][2]int {
	var starPositions [][2]int
	for x, row := range grid {
		for y, ch := range row {
			if ch == '*' {
				starPositions = append(starPositions, [2]int{x, y})
			}
		}
	}
	return starPositions
}

// calculateGearRatios finds adjacent numbers to '*' symbols and calculates their product.
func calculateGearRatios(starPositions [][2]int, numPositions map[string][][2]int) int {
	var gearRatioSum int
	for _, pos := range starPositions {
		x, y := pos[0], pos[1]
		var adjacentNums []int
		for numStr, positions := range numPositions {
			for _, p := range positions {
				if abs(p[0]-x) <= 1 && abs(p[1]-y) <= 1 {
					num, err := strconv.Atoi(numStr)
					if err != nil {
						fmt.Printf("Error converting string to int: %v\n", err)
						continue
					}
					adjacentNums = append(adjacentNums, num)
					break
				}
			}
		}
		if len(adjacentNums) == 2 {
			gearRatioSum += adjacentNums[0] * adjacentNums[1]
		}
	}
	return gearRatioSum
}

// abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// readAndPrepareGrid reads the file and prepares the grid of runes.
func readAndPrepareGrid(filename string) ([][]rune, error) {
	lines, err := utils.ReadLines(filename)
	if err != nil {
		return nil, err
	}

	var grid [][]rune
	for _, line := range lines {
		grid = append(grid, []rune(line))
	}
	return grid, nil
}

func main() {
	filename := "input.txt"

	grid, err := readAndPrepareGrid(filename)
	if err != nil {
		fmt.Printf("Error reading the file: %v\n", err)
		return
	}

	numPositions := extractNumbersAndPositions(grid)
	starPositions := findStarPositions(grid)
	gearRatioSum := calculateGearRatios(starPositions, numPositions)
	sum := sumAdjacentNumbers(grid)
	fmt.Printf("Sum of all adjacent numbers: %d\n", sum)
	fmt.Printf("Sum of all gear ratios: %d\n", gearRatioSum)
}
