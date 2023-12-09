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
	hits             int
	winningNumbers   []int
	scratchedNumbers []int
}

func parseCard(line string) *card {
	// split the line into a slice of strings by whitespace
	// remove the : from the card number
	// convert the card number to an int
	// case if string.field(line)[12] == "|"
	// use this to split the slice into winning numbers and scratched numbers
	cNum, _ := strconv.Atoi(strings.Fields(line)[1][:len(strings.Fields(line)[1])-1])
	c := card{cardNum: cNum}
	if strings.Fields(line)[12] == "|" {
		cWinNums := parseNumbers(strings.Fields(line)[2:12])
		cScratchedNums := parseNumbers(strings.Fields(line)[13:])
		c.winningNumbers = cWinNums
		c.scratchedNumbers = cScratchedNums
	} else {
		cWinNums := parseNumbers(strings.Fields(line)[2:7])
		cScratchedNums := parseNumbers(strings.Fields(line)[8:])
		c.winningNumbers = cWinNums
		c.scratchedNumbers = cScratchedNums
	}
	return &c
}

// take in a slice whitespace separated numbers and return a slice of ints
func parseNumbers(numbers []string) []int {
	var nums []int
	for _, n := range numbers {
		num, _ := strconv.Atoi(n)
		nums = append(nums, num)
	}
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
	c.hits = count
	return count

}

// Recursive function to process a card and its copies
func processCard(c *card, remainingCards []*card) int {
	if c.cardNum > len(remainingCards) || c.hits == 0 {
		return 0
	}

	totalExtra := 0
	for i := 0; i < c.hits; i++ {
		nextCardIndex := c.cardNum - 1 + i + 1
		if nextCardIndex < len(remainingCards) {
			totalExtra += 1 + processCard(remainingCards[nextCardIndex], remainingCards)
		}
	}

	return totalExtra
}

// Function to calculate the total number of cards including extras
func totalCardsIncludingExtras(originalCards []*card) int {
	total := len(originalCards) // Start with the count of original cards

	for _, c := range originalCards {
		total += processCard(c, originalCards)
	}

	return total
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
		if len(os.Args) > 1 && os.Args[1] == "print" {
			printCard(card)
		}
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
	fmt.Println("Total value:", total)

	// Calculate total number of cards including extras
	totalIncludingExtras := totalCardsIncludingExtras(cards)
	fmt.Println("Total scratchcards including extras:", totalIncludingExtras)
}
