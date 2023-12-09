package main

import (
	"AoC2023/utils"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type card struct {
	cardNum          int
	value            int
	winningNumbers   []int
	scratchedNumbers []int
}

func parseCard(line string) *card {
	// split the line into a slice of strings by whitespace
	// remove the : from the card number
	// convert the card number to an int
	cNum, _ := strconv.Atoi(strings.Fields(line)[1][:len(strings.Fields(line)[1])-1])
	cWinNums := parseNumbers(strings.Fields(line)[2:12])
	cScratchedNums := parseNumbers(strings.Fields(line)[13:])
	c := card{cardNum: cNum}
	c.winningNumbers = cWinNums
	c.scratchedNumbers = cScratchedNums
	return &c
}

// take in a slice whitespace separated numbers and return a slice of ints
func parseNumbers(numbers []string) []int {
	var nums []int
	for _, n := range numbers {
		num, _ := strconv.Atoi(n)
		nums = append(nums, num)
	}
	//sort.Ints(nums)
	slices.Sort(nums)
	return nums
}

// take in a card and print the card number, winning numbers, and scratched numbers
func printCard(c *card) {
	fmt.Printf("cardNum: %d\n", c.cardNum)
	fmt.Printf("winningNumbers: %d\n", c.winningNumbers)
	fmt.Printf("scratchedNumbers: %d\n", c.scratchedNumbers)
}

// count the number of scratched numbers on a card that are also on the winning numbers list
func countWinningNumbers(c *card) int {

	var count int
	for _, n := range c.scratchedNumbers {
		if slices.Contains(c.winningNumbers, n) {
			count++
		}
	}
	return count

}
func main() {
	var cards []*card
	lines, err := utils.ReadLines("input.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	for _, line := range lines {
		card := parseCard(line)
		// case if we have arguments
		if len(os.Args) > 1 && os.Args[1] == "print" {
			printCard(card)
		}

		//card value is the number of winning numbers on the card - 1 as an exponent of 2
		if countWinningNumbers(card) == 0 {
			card.value = 0
		} else {
			card.value = 1 << (countWinningNumbers(card) - 1)
		}
		cards = append(cards, card)

	}
	// add up the values of all the cards
	var total int
	for _, c := range cards {
		total += c.value
	}
	fmt.Println("total:", total)
}
