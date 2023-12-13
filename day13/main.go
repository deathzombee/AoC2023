package main

import (
	_ "embed"
	"fmt"
	"strings"
	"sync"
	"time"
)

//go:embed input.txt
var in string

type result struct {
	gridIndex int
	value     int
}

func rotateCCW(grid []string) []string {
	if len(grid) == 0 {
		return grid
	}

	rows := len(grid)
	cols := len(grid[0])
	var rotatedGrid []string

	for col := 0; col < cols; col++ {
		var newRow strings.Builder
		for row := 0; row < rows; row++ {
			newRow.WriteByte(grid[row][col])
		}
		rotatedGrid = append(rotatedGrid, newRow.String())
	}
	return rotatedGrid
}

func process(lines []string, old int) int {
	for i := 0; i < len(lines)-1; i++ {
		if i == old {
			continue
		}

		if lines[i] == lines[i+1] {
			for j := 0; i-j >= 0 && i+1+j < len(lines); j++ {
				if lines[i-j] != lines[i+1+j] {
					goto next
				}
			}
			return i + 1
		}
	next:
	}
	return 0
}

func main() {
	input := strings.Split(in, "\n")
	total := 0
	var wg sync.WaitGroup
	resultsChan := make(chan result) // Using an unbuffered channel

	t := time.Now()
	first := 0 // Starting index of a new batch

	// Function to process a batch of lines
	processBatch := func(start, end int) {
		grid := input[start:end]
		wg.Add(1)
		go func(grid []string, idx int) {
			defer wg.Done()
			v := process(grid, -1)
			h := process(rotateCCW(grid), -1)
			resultsChan <- result{gridIndex: idx, value: v*100 + h}
		}(grid, start)
	}

	for i, line := range input {
		if line == "" || i == len(input)-1 {
			var end int
			if i == len(input)-1 { // i.e. last line
				end = i + 1
			} else {
				end = i
			}

			if first < end {
				processBatch(first, end)
			}
			first = i + 1 // Update the start of the next batch
		}
	}
	// Closing the results channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collecting results from the goroutines
	for res := range resultsChan {
		total += res.value
	}

	elapsed := time.Since(t)
	fmt.Println(elapsed)
	fmt.Printf("Final Total: %d\n", total)
}
