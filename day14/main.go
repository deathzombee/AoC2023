package main

import (
	"AoC2023/utils"
	"errors"
	"fmt"
	"time"
)

var getSeqCount int

type Direction func([][]int) [][]int

func readGrid(s []string) (table [][]int, err error) {
	table = make([][]int, len(s))
	for y, line := range s {
		table[y] = make([]int, len(line))
		for x := 0; x < len(line); x++ {
			switch line[x] {
			case 'O':
				table[y][x] = 2
			case '#':
				table[y][x] = 1
			case '.':
				table[y][x] = 0
			default:
				return table, errors.New("invalid input")
			}
		}
	}
	return table, err
}
func north(table [][]int) [][]int {
	for x := range table[0] {
		stopPoint := 0
		for y, value := range table {
			if value[x] == 1 {
				stopPoint = y + 1
			} else if value[x] == 2 {
				if y != stopPoint {
					table[stopPoint][x], table[y][x] = table[y][x], 0
				}
				stopPoint++
			}
		}
	}
	return table
}
func south(table [][]int) [][]int {
	for x := range table[0] {
		stopPoint := len(table) - 1
		for y := len(table) - 1; y >= 0; y-- {
			if table[y][x] == 1 {
				stopPoint = y - 1
			} else if table[y][x] == 2 {
				if y != stopPoint {
					table[stopPoint][x], table[y][x] = table[y][x], 0
				}
				stopPoint--
			}
		}
	}
	return table
}
func west(table [][]int) [][]int {
	for y := range table {
		stopPoint := 0
		for x, value := range table[y] {
			if value == 1 {
				stopPoint = x + 1
			} else if value == 2 {
				if x != stopPoint {
					table[y][stopPoint], table[y][x] = table[y][x], 0
				}
				stopPoint++
			}
		}
	}
	return table
}
func east(table [][]int) [][]int {
	for y := range table {
		stopPoint := len(table[0]) - 1
		for x := len(table[0]) - 1; x >= 0; x-- {
			if table[y][x] == 1 {
				stopPoint = x - 1
			} else if table[y][x] == 2 {
				if x != stopPoint {
					table[y][stopPoint], table[y][x] = table[y][x], 0
				}
				stopPoint--
			}
		}
	}
	return table
}

func sequence(values []int) (frequency, offset int, ok bool) {
	getSeqCount++

	if len(values) < 10 {
		return 0, 0, false
	}

	for offset = 0; offset < len(values); offset++ {
		for frequency = 1; offset+2*frequency < len(values); frequency++ {
			isRepeating := true
			for i := 0; i < frequency; i++ {
				if offset+2*frequency+i >= len(values) ||
					values[offset+i] != values[offset+frequency+i] ||
					values[offset+i] != values[offset+2*frequency+i] {
					isRepeating = false
					break
				}
			}
			if isRepeating {
				return frequency, offset, true
			}
		}
	}

	return 0, 0, false
}

func loadCalc(s []string, transform Direction) (int, error) {
	var result int
	if table, err := readGrid(s); err != nil {
		return 0, err
	} else {
		table = transform(table)
		var numRocks int
		for y := 0; y < len(table); y++ {
			numRocks = 0
			for x := 0; x < len(table[y]); x++ {
				if table[y][x] == 2 {
					numRocks++
				}
			}
			result += numRocks * (len(table) - y)
		}
	}
	return result, nil
}

func cycledLoad(s []string, numCycles int) (int, error) {
	var result int

	if table, err := readGrid(s); err != nil {
		return 0, err
	} else {
		var results []int

		for {
			table = north(table)
			table = west(table)
			table = south(table)
			table = east(table)

			numRocks := 0
			result = 0

			for y := 0; y < len(table); y++ {
				numRocks = 0
				for x := 0; x < len(table[y]); x++ {
					if table[y][x] > 1 {
						numRocks++
					}
				}
				result += numRocks * (len(table) - y)
			}
			results = append(results, result)
			if frequency, offset, ok := sequence(results); ok {
				mod := (numCycles - offset) % frequency
				fmt.Println("frequency:", frequency, "offset:", offset, "ok:", ok, "mod:", mod, "numCycles:", numCycles)
				result = results[offset+mod-1]
				fmt.Println("index:", offset+mod-1)
				break
			}
		}
	}
	return result, nil
}
func main() {

	output, _ := utils.ReadLines("input.txt")
	time1 := time.Now()
	part1, _ := loadCalc(output, north)
	fmt.Println("p1 time:", time.Since(time1))
	fmt.Println("part 1:", part1)
	time2 := time.Now()
	part2, _ := cycledLoad(output, 1000000000)
	fmt.Println("p2 time:", time.Since(time2))
	fmt.Println("part 2:", part2)
	fmt.Println("getSeqCount:", getSeqCount)

}
