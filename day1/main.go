package main

import (
	"AoC2023/utils"
	"fmt"
	"log"
	"strconv"
	"unicode"
)

var digitMap = map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

var validNumbersSet map[string]struct{}
var maxNumberWordLength int

func preprocessNumberWords() {
	validNumbersSet = make(map[string]struct{})
	maxNumberWordLength = 0
	for word := range digitMap {
		validNumbersSet[word] = struct{}{}
		if len(word) > maxNumberWordLength {
			maxNumberWordLength = len(word)
		}
	}
}

func startsWithNumber(str string) int {
	maxLength := min(len(str), maxNumberWordLength)
	for length := 1; length <= maxLength; length++ {
		if _, ok := validNumbersSet[str[:length]]; ok {
			return digitMap[str[:length]]
		}
	}
	return -1
}

func endsWithNumber(str string) int {
	maxLength := min(len(str), maxNumberWordLength)
	for length := 1; length <= maxLength; length++ {
		if _, ok := validNumbersSet[str[len(str)-length:]]; ok {
			return digitMap[str[len(str)-length:]]
		}
	}
	return -1
}

func getCalibrationValue1(s string) int {
	firstDigit, lastDigit := -1, -1
	for _, r := range s {
		if unicode.IsDigit(r) {
			if firstDigit == -1 {
				firstDigit = int(r - '0')
			}
			lastDigit = int(r - '0')
		}
	}
	if firstDigit == -1 || lastDigit == -1 {
		return 0 // No digits found in the line
	}
	value, _ := strconv.Atoi(fmt.Sprintf("%d%d", firstDigit, lastDigit))
	return value
}

func getCalibrationValue2(s string) int {
	first, last := -1, -1
	index := 0

	for index < len(s) {
		if first < 0 {
			if unicode.IsDigit(rune(s[index])) {
				first = int(s[index] - '0')
			} else if num := startsWithNumber(s[index:]); num != -1 {
				first = num
			}
		}

		if last < 0 {
			if unicode.IsDigit(rune(s[len(s)-1-index])) {
				last = int(s[len(s)-1-index] - '0')
			} else if num := endsWithNumber(s[:len(s)-index]); num != -1 {
				last = num
			}
		}

		if first >= 0 && last >= 0 {
			break
		}
		index++
	}

	if first == -1 || last == -1 {
		return 0
	}
	return first*10 + last
}

func part1() {
	lines, err := utils.ReadLines("input.txt")
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}
	sum := 0
	for _, line := range lines {
		calibrationValue := getCalibrationValue1(line)
		//fmt.Printf("Line: %s, Calibration Value: %d\n", line, calibrationValue)
		sum += calibrationValue
	}

	fmt.Println("Part 1 total sum of calibration values:", sum)
}

func part2() {
	lines, err := utils.ReadLines("input.txt")
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}
	sum := 0
	for _, line := range lines {
		calibrationValue := getCalibrationValue2(line)
		//fmt.Printf("Line: %s, Calibration Value: %d\n", line, calibrationValue)
		sum += calibrationValue
	}

	fmt.Println("Part 2 total sum of calibration values:", sum)
}

func main() {
	preprocessNumberWords() // Preprocess number words once before processing lines
	part1()
	part2()
}
