package main

import (
	"AoC2023/utils"
	"fmt"
	"time"
)

// returns the accumulation of columns left of the mirror plus rows above the mirror times 100
func acc(grid []string, swap bool) int {
	vert := false
	// try horizontal match first
	mid := midPoint(intVals(grid, false), swap) // Horizontal
	// if no horizontal match, try vertical
	if mid == -1 {
		vert = true
		mid = midPoint(intVals(grid, true), swap) // Vertical
	}

	if vert {
		// columns left of the mirror
		return mid + 1
	}
	// rows above the mirror for horizontal match * 100
	return 100 * (mid + 1)
}

// midPoint returns the index of the middle number in the list of numbers
func midPoint(numbers []int, swap bool) int {
	var mid, start int
	for i := 0; i < len(numbers)-1; i++ {
		start = i
		var swapFound bool
		for end := i + 1; start >= 0 && end < len(numbers); end++ {
			xor := numbers[start] ^ numbers[end]
			valid := bW(xor, swap)
			swapFound = swapFound || (swap && valid && xor > 0)

			if valid {
				mid = i
			} else {
				mid = -1
				break
			}
			start--
		}
		if mid != -1 && (!swap || swapFound) {
			return mid
		}
	}
	return -1
}

// bW returns true if xor is 0 and swap is false or if xor is a power of 2 and swap is true
// otherwise returns false
func bW(xor int, swap bool) bool {
	return (xor == 0 && !swap) || (swap && (xor&(xor-1)) == 0)
}

// intVals returns a list of ints representing the binary values of the grids
// big endian
func intVals(g []string, vert bool) []int {
	var ints []int
	if !vert {
		for _, grid := range g {
			num := hNum(grid)
			ints = append(ints, num)
		}
	} else {
		for i := 0; i < len(g[0]); i++ {
			num := vNum(g, i)
			ints = append(ints, num)
		}
	}
	return ints
}

func hNum(g string) int {
	var num int
	for i, char := range g {
		if char == '#' {
			num += 1 << i
		}
	}
	return num
}

func vNum(g []string, index int) int {
	var num int
	for i, grid := range g {
		if grid[index] == '#' {
			num += 1 << i
		}
	}
	return num
}
func main() {
	lines, _ := utils.ReadLines("input.txt")

	var sum, sum2 int
	var g []string
	var t, t2 time.Duration
	for i, line := range lines {
		if len(line) != 0 {
			g = append(g, line)
		}
		if len(line) == 0 || i == len(lines)-1 {
			start := time.Now()
			sum += acc(g, false)
			t += time.Since(start)

			start = time.Now()
			sum2 += acc(g, true)
			t2 += time.Since(start)
			g = []string{}
		}
	}

	fmt.Println("Part 1: ", sum, " time:", t)
	fmt.Println("Part 2: ", sum2, " time:", t2)
}
