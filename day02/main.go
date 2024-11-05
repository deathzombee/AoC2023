package main

import (
	"AoC2023/utils"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	lines, err := utils.ReadLines("input.txt")
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	redLimit, greenLimit, blueLimit := 12, 13, 14
	var sumOfGameIDs int
	var sumOfGamePowers int

	for _, line := range lines {
		possible, _ := isGamePossible(line, redLimit, greenLimit, blueLimit)
		if possible {
			gameID := getGameID(line)
			sumOfGameIDs += gameID
			//fmt.Printf("Game %d is possible with %d subsets\n", gameID, subsetCount)
		}
		power := gamePower(line)
		sumOfGamePowers += power
		//fmt.Printf("Game %d has power %d\n", getGameID(line), power)
	}

	fmt.Printf("Sum of game IDs: %d\n", sumOfGameIDs)
	fmt.Printf("Sum of game powers: %d\n", sumOfGamePowers)
}

func getGameID(line string) int {
	parts := strings.Split(line, ":")
	gameID, err := strconv.Atoi(strings.TrimSpace(strings.Split(parts[0], " ")[1]))
	if err != nil {
		log.Fatalf("getGameID: %v", err)
	}
	return gameID
}

func isGamePossible(line string, redLimit, greenLimit, blueLimit int) (bool, int) {
	parts := strings.Split(line, ": ")
	subsets := strings.Split(parts[1], ";")
	for _, subset := range subsets {
		cubes := strings.Split(strings.TrimSpace(subset), ",")
		for _, cube := range cubes {
			colorCount := strings.Split(strings.TrimSpace(cube), " ")
			count, err := strconv.Atoi(colorCount[0])
			if err != nil {
				log.Fatalf("isGamePossible: %v", err)
			}

			switch colorCount[1] {
			case "red":
				if count > redLimit {
					return false, len(subsets)
				}
			case "green":
				if count > greenLimit {
					return false, len(subsets)
				}
			case "blue":
				if count > blueLimit {
					return false, len(subsets)
				}
			}
		}
	}
	return true, len(subsets)
}
func gamePower(line string) int {
	var redCount, greenCount, blueCount int
	parts := strings.Split(line, ": ")
	subsets := strings.Split(parts[1], ";")
	for _, subset := range subsets {
		cubes := strings.Split(strings.TrimSpace(subset), ",")
		for _, cube := range cubes {
			colorCount := strings.Split(strings.TrimSpace(cube), " ")
			count, err := strconv.Atoi(colorCount[0])

			if err != nil {
				log.Fatalf("gamePower: %v", err)
			}

			//keep track of max count of each color in each game
			switch colorCount[1] {
			case "red":
				if count > redCount {
					redCount = count
				}
			case "green":
				if count > greenCount {
					greenCount = count
				}
			case "blue":
				if count > blueCount {
					blueCount = count
				}
			}

			//fmt.Printf("colorCount: %v\n", colorCount)

		}

	}

	return redCount * greenCount * blueCount
}
