package main

import (
	"AoC2023/utils"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type Check struct {
	cond  uint8
	v     uint8
	value int
	dest  string
}

type Rule struct {
	name string
	chk  []Check
}

type Part map[uint8]int

func (p Part) String() string {
	return fmt.Sprintf("x=%d, m=%d, a=%d, s=%d", p['x'], p['m'], p['a'], p['s'])
}

func (r Rule) String() string {
	return fmt.Sprintf("%s: %v", r.name, r.chk)
}

func (c Check) String() string {
	if c.cond == 'D' {
		return fmt.Sprintf("else: %s", c.dest)
	}
	return fmt.Sprintf("%c%c%d => %s", c.v, c.cond, c.value, c.dest)
}

func (r Rule) apply(prt Part) string {
	for _, ch := range r.chk {
		switch ch.cond {
		case '<':
			if prt[ch.v] < ch.value {
				return ch.dest
			}
		case '>':
			if prt[ch.v] > ch.value {
				return ch.dest
			}
		case 'D':
			return ch.dest
		}
	}
	panic("no check found")
}

func parseRule(line string) Rule {
	// Regex to match the rule format
	ruleRegex := regexp.MustCompile(`(\w+)\{((?:[xmas][<=>]\d+:\w+,?\s?)+)(\w+)}`)
	match := ruleRegex.FindStringSubmatch(line)
	if match == nil {
		panic("invalid rule format")
	}

	name := match[1]
	checkString := match[2]
	defaultWorkflow := match[3]

	cRegex := regexp.MustCompile(`([xmas])([<=>])(\d+):(\w+)`)
	cMatches := cRegex.FindAllStringSubmatch(checkString, -1)

	var checks []Check
	for _, prt := range cMatches {
		subject := prt[1][0]
		cond := prt[2][0]
		value, _ := strconv.Atoi(prt[3])
		dest := prt[4]
		checks = append(checks, Check{
			cond:  cond,
			v:     subject,
			value: value,
			dest:  dest,
		})
	}
	// Add the default workflow as the last check
	checks = append(checks, Check{
		cond: 'D', // Using 'D' to represent the default case
		dest: defaultWorkflow,
	})

	return Rule{name: name, chk: checks}
}

func parseParts(line string) Part {
	partRegex := regexp.MustCompile(`\{([xmas]=\d+,?)+}`)
	match := partRegex.FindStringSubmatch(line)
	if match == nil {
		panic("invalid part format")
	}

	keyValueRegex := regexp.MustCompile(`([xmas])=(\d+)`)
	keyValueMatches := keyValueRegex.FindAllStringSubmatch(match[0], -1)

	p := make(Part)
	for _, kv := range keyValueMatches {
		key := kv[1][0]
		value, _ := strconv.Atoi(kv[2])
		p[key] = value
	}
	return p
}

func run(rules map[string]Rule, parts Part) int {
	var pc = "in"
	for {
		rule := rules[pc]
		label := rule.apply(parts)
		if label == "R" {
			return 0
		}
		if label == "A" {
			return parts['x'] + parts['m'] + parts['a'] + parts['s']
		}
		pc = label
	}
}

func Part1(input []string) int {
	var workflowInput []string
	var partInput []string
	partSection := false

	for _, line := range input {
		if line == "" {
			partSection = true
			continue
		}

		if partSection {
			partInput = append(partInput, line)
		} else {
			workflowInput = append(workflowInput, line)
		}
	}

	var rules = make(map[string]Rule)
	for _, line := range workflowInput {
		rule := parseRule(line)
		rules[rule.name] = rule
	}

	var parts []Part
	for _, line := range partInput {
		p := parseParts(line)
		parts = append(parts, p)
	}

	var res int
	for _, part := range parts {
		res += run(rules, part)
	}
	return res
}

func Part2(input []string) int {
	// Implement Part 2 logic here
	return 0
}

func main() {
	start := time.Now()

	input, _ := utils.ReadLines("input.txt")
	fmt.Println("part1: ", Part1(input))
	fmt.Println(time.Since(start))

	start = time.Now()
	fmt.Println("part2: ", Part2(input))
	fmt.Println(time.Since(start))
}
