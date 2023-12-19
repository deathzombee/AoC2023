package main

import (
	_ "embed"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// sick we can embed files didn't know that
//
//go:embed input.txt
var inputDay string

// Interval struct to hold the interval
type Interval struct{ Min, Max int }

var EmptyInterval = Interval{0, -1}

// transformer takes in a seed interval and a rule interval and returns the transformed interval and the remaining intervals
func transformer(seedInterval, ruleInterval Interval, destStart int) (transformed Interval, remaining []Interval) {
	// Define the start and end points of the seed interval and rule interval
	seedStart, seedEnd := seedInterval.Min, seedInterval.Max
	ruleStart, ruleEnd := ruleInterval.Min, ruleInterval.Max

	// Calculate the offset for transformation
	offset := destStart - ruleStart
	//print the seedstart seedend rulestart ruleend and offset
	//fmt.Printf("seedStart: %d, seedEnd: %d, ruleStart: %d, ruleEnd: %d, deststart %d, offset %d\n", seedStart, seedEnd, ruleStart, ruleEnd, destStart, offset)

	// Check if there is no overlap between the seed interval and rule interval
	if seedEnd < ruleStart || seedStart > ruleEnd {
		// No overlap, so return the seed interval as is
		return EmptyInterval, []Interval{seedInterval}
	}

	// Handle the case where the seed interval is completely within the rule interval
	if seedStart >= ruleStart && seedEnd <= ruleEnd {
		return Interval{Min: seedStart + offset, Max: seedEnd + offset}, remaining
	}

	// Handle the case where the seed interval partially overlaps the rule interval
	if seedEnd >= ruleStart && seedEnd <= ruleEnd {
		// Add non-overlapping part to remaining intervals if it exists
		if seedStart < ruleStart {
			remaining = append(remaining, Interval{Min: seedStart, Max: ruleStart - 1})
		}
		return Interval{Min: ruleStart + offset, Max: seedEnd + offset}, remaining
	}

	if seedStart >= ruleStart && seedStart <= ruleEnd {
		// Add non-overlapping part to remaining intervals if it exists
		if seedEnd > ruleEnd {
			remaining = append(remaining, Interval{Min: ruleEnd + 1, Max: seedEnd})
		}
		return Interval{Min: seedStart + offset, Max: ruleEnd + offset}, remaining
	}

	// Handle the case where the rule interval is completely within the seed interval
	if seedStart <= ruleStart && ruleEnd <= seedEnd {
		// Add non-overlapping parts to remaining intervals if they exist
		if seedStart < ruleStart {
			remaining = append(remaining, Interval{Min: seedStart, Max: ruleStart - 1})
		}
		if ruleEnd < seedEnd {
			remaining = append(remaining, Interval{Min: ruleEnd + 1, Max: seedEnd})
		}
		return Interval{Min: ruleStart + offset, Max: ruleEnd + offset}, remaining
	}

	// If none of the above conditions are met, this is an unexpected case
	panic("transformer: unexpected interval configuration")
}

type Rule struct {
	i    Interval
	dest int
}

func applyOneRule(rule Rule, seeds []Interval) (transformed []Interval, other []Interval) {
	//fmt.Printf("Applying Rule: %+v to Seeds: %+v\n", rule, seeds)
	for _, seed := range seeds {
		if seed.Max < rule.i.Min || seed.Min > rule.i.Max {
			other = append(other, seed)
		} else {
			t, o := transformer(seed, rule.i, rule.dest)
			transformed = append(transformed, t)
			other = append(other, o...)
			//fmt.Printf("Rule Applied Successfully: %+v, Transformed: %+v, Remaining: %+v\n", rule, t, o)
		}
	}
	return transformed, other
}

func applyTransformationRules(rules []Rule, interval Interval) []Interval {
	var resultingIntervals []Interval
	intervalsToProcess := []Interval{interval}

	for _, rule := range rules {
		var transformedIntervals []Interval
		var nextIntervalsToProcess []Interval

		for _, currentInterval := range intervalsToProcess {
			transformed, remaining := applyOneRule(rule, []Interval{currentInterval})
			transformedIntervals = append(transformedIntervals, transformed...)
			nextIntervalsToProcess = append(nextIntervalsToProcess, remaining...)
		}

		intervalsToProcess = nextIntervalsToProcess
		if len(transformedIntervals) > 0 {
			resultingIntervals = append(resultingIntervals, transformedIntervals...)
		}
	}

	return append(resultingIntervals, intervalsToProcess...)
}

func processSeedIntervals(inputParts []string, seedIntervals []Interval) int {
	// Initialize the result to the maximum possible integer value
	minSeedLocation := math.MaxInt64

	for _, part := range inputParts {
		// Split the part into lines and parse each line as a rule
		lines := strings.Split(part, "\n")[1:] // Skip the first line as it's a header
		var rules []Rule

		for _, line := range lines {
			var destStart, srcStart, length int
			_, err := fmt.Sscanf(line, "%d %d %d", &destStart, &srcStart, &length)
			if err != nil {
				return 0
			}
			rules = append(rules, Rule{Interval{srcStart, srcStart + length - 1}, destStart})
		}

		// Apply all rules to each seed interval
		var transformedSeeds []Interval
		for _, seedInterval := range seedIntervals {
			transformedSeeds = append(transformedSeeds, applyTransformationRules(rules, seedInterval)...)
		}

		seedIntervals = transformedSeeds
	}

	// Find the minimum value across all the final seed intervals
	//fmt.Printf("Final Seed Intervals: ")
	for _, interval := range seedIntervals {
		//fmt.Printf("%+v ", interval)
		if interval.Min < minSeedLocation {
			minSeedLocation = interval.Min
		}
	}
	//fmt.Printf("\n")
	return minSeedLocation
}

func Part1(input string) int {
	input = strings.TrimSuffix(input, "\n")
	parts := strings.Split(input, "\n\n")
	numbers, _ := strings.CutPrefix(parts[0], "seeds: ")

	var seeds []Interval
	for _, v := range strings.Fields(numbers) {
		n, _ := strconv.Atoi(v)
		seeds = append(seeds, Interval{n, n})
	}
	return processSeedIntervals(parts[1:], seeds)
}

func Part2(input string) int {
	input = strings.TrimSuffix(input, "\n")
	parts := strings.Split(input, "\n\n")
	numbers, _ := strings.CutPrefix(parts[0], "seeds: ")

	var values []int
	for _, v := range strings.Fields(numbers) {
		if n, err := strconv.Atoi(v); err == nil {
			values = append(values, n)
		}
	}

	var seeds []Interval
	for i := 0; i < len(values); i += 2 {
		seeds = append(seeds, Interval{values[i], values[i] + values[i+1] - 1})
	}
	return processSeedIntervals(parts[1:], seeds)
}

func main() {
	start := time.Now()
	fmt.Println("part1: ", Part1(inputDay))
	fmt.Println(time.Since(start))
	fmt.Printf("\n")
	start = time.Now()
	fmt.Println("part2: ", Part2(inputDay))
	fmt.Println(time.Since(start))
}
