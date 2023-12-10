package main

import (
	"AoC2023/utils"
	"fmt"
)

func main() {
	content, _ := utils.ReadLines("input.txt") // Adjust the path to the file

	lines := content
	grid := make([][]rune, len(lines))
	for i, line := range lines {
		grid[i] = []rune(line)
	}

	pipeTypes := map[rune][]string{
		'|': {"n", "s"},
		'-': {"w", "e"},
		'L': {"n", "e"},
		'J': {"n", "w"},
		'7': {"s", "w"},
		'F': {"s", "e"},
		'S': {"n", "s", "w", "e"},
	}

	directions := map[string][3]interface{}{
		"n": {-1, 0, "s"},
		"s": {1, 0, "n"},
		"w": {0, -1, "e"},
		"e": {0, 1, "w"},
	}

	var start [2]int
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == 'S' {
				start = [2]int{i, j}
				break
			}
		}
	}

	encounteredPlaces := make(map[[2]int]int)
	searchQueue := [][2]interface{}{{start, 0}}

	for len(searchQueue) > 0 {
		item := searchQueue[0]
		searchQueue = searchQueue[1:]

		current := item[0].([2]int)
		distance := item[1].(int)

		if _, found := encounteredPlaces[current]; found {
			continue
		}

		encounteredPlaces[current] = distance
		i, j := current[0], current[1]
		availableDirections := pipeTypes[grid[i][j]]

		for _, dir := range availableDirections {
			dirData := directions[dir]
			di, dj := dirData[0].(int), dirData[1].(int)
			opposite := dirData[2].(string)

			newI, newJ := i+di, j+dj
			if newI < 0 || newI >= len(grid) || newJ < 0 || newJ >= len(grid[newI]) {
				continue
			}

			target := grid[newI][newJ]
			if _, exists := pipeTypes[target]; !exists {
				continue
			}

			targetDirections := pipeTypes[target]
			for _, tDir := range targetDirections {
				if tDir == opposite {
					searchQueue = append(searchQueue, [2]interface{}{[2]int{newI, newJ}, distance + 1})
					break
				}
			}
		}
	}

	maxDistance := 0
	for _, dist := range encounteredPlaces {
		if dist > maxDistance {
			maxDistance = dist
		}
	}

	fmt.Println("Maximum distance from start:", maxDistance)

	// Determine the type of pipe at the start position
	getPipeType := func(i, j int) rune {
		var reachableDirs []string
		for dir, data := range directions {
			di, dj := data[0].(int), data[1].(int)
			if i+di < 0 || i+di >= len(grid) || j+dj < 0 || j+dj >= len(grid[i+di]) {
				continue
			}
			if _, ok := encounteredPlaces[[2]int{i + di, j + dj}]; !ok {
				continue
			}
			reachableDirs = append(reachableDirs, dir)
		}
		for pType, dirs := range pipeTypes {
			if len(reachableDirs) == len(dirs) {
				matches := true
				for _, rDir := range reachableDirs {
					found := false
					for _, d := range dirs {
						if rDir == d {
							found = true
							break
						}
					}
					if !found {
						matches = false
						break
					}
				}
				if matches {
					return pType
				}
			}
		}
		return ' '
	}

	// Replace the start position with the determined pipe type
	grid[start[0]][start[1]] = getPipeType(start[0], start[1])

	// Count the inside of the loop
	insideCount := 0
	for i := range grid {
		norths := 0
		for j := range grid[i] {
			if _, ok := encounteredPlaces[[2]int{i, j}]; ok {
				if containsDirection(pipeTypes[grid[i][j]], 'n') {
					norths++
				}
				continue
			}
			if norths%2 != 0 {
				grid[i][j] = 'I'
				insideCount++
			} else {
				grid[i][j] = 'O'
			}
		}
	}

	fmt.Println("Inside count:", insideCount)
}

// containsDirection checks if a slice of strings contains a specific direction.
func containsDirection(slice []string, dir rune) bool {
	for _, s := range slice {
		if rune(s[0]) == dir {
			return true
		}
	}
	return false
}
