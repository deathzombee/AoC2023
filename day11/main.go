package main

import (
	"AoC2023/utils"
	"fmt"
	"time"
)

type Pos struct {
	X, Y int
}

func (p Pos) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func parseInput(input []string) ([][]rune, []Pos) {
	lines := make([][]rune, len(input))
	var galaxies []Pos
	for i, line := range input {
		lines[i] = []rune(line)
		for j, char := range line {
			if char == '#' {
				galaxies = append(galaxies, Pos{X: j, Y: i})
			}
		}
	}
	return lines, galaxies
}

func emtlc(lines [][]rune) ([]int, []int) {
	var emptyLines, emptyColumns []int
	for i, line := range lines {
		if bemtc(line) {
			emptyLines = append(emptyLines, i)
		}
	}
	for i := 0; i < len(lines[0]); i++ {
		if bemtl(lines, i) {
			emptyColumns = append(emptyColumns, i)
		}
	}
	return emptyLines, emptyColumns
}

func bemtc(line []rune) bool {
	for _, char := range line {
		if char == '#' {
			return false
		}
	}
	return true
}

func bemtl(lines [][]rune, columnIndex int) bool {
	for _, line := range lines {
		if line[columnIndex] == '#' {
			return false
		}
	}
	return true
}

func expgal(galaxies []Pos, emptyLines []int, emptyColumns []int, factor int) []Pos {
	var expanded []Pos
	for _, galaxy := range galaxies {
		expanded = append(expanded, scale(galaxy, emptyLines, emptyColumns, factor))
	}
	return expanded
}

func scale(p Pos, emptyLines []int, emptyColumns []int, factor int) Pos {
	var addX, addY int
	for _, l := range emptyLines {
		if l < p.Y {
			addY += factor - 1
		}
	}
	for _, c := range emptyColumns {
		if c < p.X {
			addX += factor - 1
		}
	}
	return Pos{X: p.X + addX, Y: p.Y + addY}
}

// ManhattanDistance calculates the Manhattan distance between two positions.
func ManhattanDistance(p1, p2 Pos) int {
	return abs(p1.X-p2.X) + abs(p1.Y-p2.Y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solve(input []string, factor int) int {
	lines, galaxies := parseInput(input)

	emptyLines, emptyColumns := emtlc(lines)

	expandedGalaxies := expgal(galaxies, emptyLines, emptyColumns, factor)

	var totalDistance int
	for i := 0; i < len(expandedGalaxies); i++ {
		for j := i + 1; j < len(expandedGalaxies); j++ {
			totalDistance += ManhattanDistance(expandedGalaxies[i], expandedGalaxies[j])
		}
	}

	return totalDistance
}
func part1(input []string) int {

	return solve(input, 2)
}
func part2(input []string) int {
	return solve(input, 1000000)
}

func main() {
	input, _ := utils.ReadLines("input.txt")
	t := time.Now()
	fmt.Println("Part 1:", part1(input), "run time:", time.Since(t))
	t = time.Now()
	fmt.Println("Part 2:", part2(input), "run time:", time.Since(t))

}
