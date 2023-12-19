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
func generate(input []string) (map[string]Rule, []Part) {
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
	return rules, parts
}
func Part1(rules map[string]Rule, parts []Part) int {
	var res int
	for _, part := range parts {
		res += run(rules, part)
	}
	return res
}

func cloneMap(original map[rune][2]int) map[rune][2]int {
	c := make(map[rune][2]int, 4)
	for k, v := range original {
		c[k] = v
	}
	return c
}

func count(workflows map[string]Rule, workflow string, values map[rune][2]int) int {
	if workflow == "R" {
		return 0
	} else if workflow == "A" {
		product := 1
		for _, v := range values {
			product *= v[1] - v[0] + 1
		}
		return product
	}

	total := 0
	for _, r := range workflows[workflow].chk {
		// low high for this variable
		v := values[rune(r.v)]

		// initialize true and false ranges
		var tv [2]int
		var fv [2]int
		switch r.cond {
		case '<':
			// if the condition is <, then the true range is from the low to the value - 1
			// and the false range is from the value to the high
			tv = [2]int{v[0], r.value - 1}
			fv = [2]int{r.value, v[1]}
		case '>':
			// if the condition is >, then the true range is from the value + 1 to the high
			// and the false range is from the low to the value
			tv = [2]int{r.value + 1, v[1]}
			fv = [2]int{v[0], r.value}
		default:
			total += count(workflows, r.dest, values)
			continue
		}

		// branch for true into its own recursion with a copy of the values
		if tv[0] <= tv[1] {
			v2 := cloneMap(values)
			v2[rune(r.v)] = tv
			total += count(workflows, r.dest, v2)
		}

		// false value keeps going through this rules checks
		// unless its low is greater than its high
		if fv[0] > fv[1] {
			break
		}
		// update the values for the next check
		values[rune(r.v)] = fv
	}
	// total keeps track of the product between calls
	return total
}
func main() {
	vals := map[rune][2]int{
		'x': {1, 4000},
		'm': {1, 4000},
		'a': {1, 4000},
		's': {1, 4000}}
	start := time.Now()
	input, _ := utils.ReadLines("input.txt")
	ru, pa := generate(input)
	fmt.Println("input parsing:", time.Since(start))
	start = time.Now()
	p1 := Part1(ru, pa)
	fmt.Println("part1: ", p1, time.Since(start))
	start = time.Now()
	p2 := count(ru, "in", vals)
	fmt.Println("part2: ", p2, time.Since(start))
}
