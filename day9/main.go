package main

import (
	"AoC2023/utils"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func parseLineToInts(line string) ([]int, error) {
	// split the line into a slice of strings by whitespace
	strValues := strings.Fields(line)
	// convert the slice of strings to a slice of ints
	intValues := make([]int, len(strValues))
	for i, str := range strValues {
		var err error
		intValues[i], err = strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
	}
	return intValues, nil
}

// make a slice of ints that is the difference between each element in the slice
func makeDifferences(values []int) []int {
	differences := make([]int, len(values)-1)
	for i := 0; i < len(values)-1; i++ {
		differences[i] = values[i+1] - values[i]
	}
	return differences
}

func allZeros(values []int) bool {
	for _, v := range values {
		if v != 0 {
			return false
		}
	}
	return true
}

// Parameters:
//
//	values - The initial sequence of integers.
//
// Returns:
//
//	The next extrapolated value in the sequence.
func extrapolateSequence(values []int) int {
	differenceSequences := [][]int{values}

	for {
		lastSequence := differenceSequences[len(differenceSequences)-1]
		differences := makeDifferences(lastSequence)
		differenceSequences = append(differenceSequences, differences)
		if allZeros(differences) {
			break
		}
	}
	return calculateNextValue(differenceSequences)
}

// using the difference sequences, calculate the next value in the sequence
func calculateNextValue(differenceSequences [][]int) int {
	for i := len(differenceSequences) - 2; i >= 0; i-- {
		lastValue := differenceSequences[i][len(differenceSequences[i])-1]
		diffValue := differenceSequences[i+1][len(differenceSequences[i+1])-1]
		differenceSequences[i] = append(differenceSequences[i], lastValue+diffValue)
	}
	return differenceSequences[0][len(differenceSequences[0])-1]
}

func main() {
	lines, err := utils.ReadLines("input.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	totalNext := 0
	totalPrevious := 0
	for _, line := range lines {
		values, err := parseLineToInts(line)
		if err != nil {
			fmt.Println("Error parsing line:", err)
			continue
		}
		nextValue := extrapolateSequence(values)
		// reverse the slice and extrapolate again
		slices.Reverse(values)
		previousValue := extrapolateSequence(values)
		totalNext += nextValue
		totalPrevious += previousValue
	}

	fmt.Println("Total sum of extrapolated next values:", totalNext)
	fmt.Println("Total sum of extrapolated previous values:", totalPrevious)
}
