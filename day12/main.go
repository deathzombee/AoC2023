package main

import (
	"AoC2023/utils"
	"fmt"
	"strings"
	"time"
)

type key struct {
	length, within, remaining int
}

// recursive function to find the number of ways to arrange the vents
func possibleWays(cache map[key]uint64, s []byte, within, remaining []int) uint64 {
	if len(s) == 0 {
		switch {
		case within == nil && len(remaining) == 0:
			return 1
		case within != nil && len(remaining) == 1 && within[0] == remaining[0]:
			return 1
		default:
			return 0
		}
	}

	if within != nil && len(remaining) == 0 {
		return 0
	}

	cacheKey := key{len(s), 0, len(remaining)}
	if within != nil {
		cacheKey.within = within[0]
	}
	if val, found := cache[cacheKey]; found {
		return val
	}

	var ways uint64
	switch s[0] {
	case '.':
		if within != nil && within[0] != remaining[0] {
			ways = 0
		} else if within != nil {
			ways = possibleWays(cache, s[1:], nil, remaining[1:])
		} else {
			ways = possibleWays(cache, s[1:], nil, remaining)
		}
	case '#':
		if within != nil {
			newWithin := within[0] + 1
			ways = possibleWays(cache, s[1:], []int{newWithin}, remaining)
		} else {
			ways = possibleWays(cache, s[1:], []int{1}, remaining)
		}
	case '?':
		if within != nil {
			ways = possibleWays(cache, s[1:], []int{within[0] + 1}, remaining)
			if within[0] == remaining[0] {
				ways += possibleWays(cache, s[1:], nil, remaining[1:])
			}
		} else {
			ways = possibleWays(cache, s[1:], []int{1}, remaining) +
				possibleWays(cache, s[1:], nil, remaining)
		}
	}

	cache[cacheKey] = ways
	return ways
}

func main() {
	input, _ := utils.ReadLines("input.txt")

	var totalP1, totalP2 uint64
	for _, line := range input {
		parts := strings.Split(line, " ")
		numsStr := strings.Split(parts[1], ",")
		nums := make([]int, len(numsStr))
		for i, n := range numsStr {
			_, err := fmt.Sscanf(n, "%d", &nums[i])
			if err != nil {
				return
			}
		}

		newVents := strings.Repeat(parts[0]+"?", 4) + parts[0]
		newNums := make([]int, 0, len(nums)*5)
		for i := 0; i < 5; i++ {
			newNums = append(newNums, nums...)
		}

		cache := make(map[key]uint64)
		p1 := possibleWays(cache, []byte(parts[0]), nil, nums)
		cache = make(map[key]uint64)
		p2 := possibleWays(cache, []byte(newVents), nil, newNums)
		totalP1 += p1
		totalP2 += p2
	}
	start := time.Now()

	fmt.Printf("Total arrangements (Part 1): %d\n", totalP1)
	fmt.Println("Time taken p1: ", time.Since(start))
	start = time.Now()
	fmt.Printf("Total arrangements (Part 2): %d\n", totalP2)
	fmt.Printf("Time taken p2: %v\n", time.Since(start))
}
