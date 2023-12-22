package main

import (
	"AoC2023/utils"
	"container/list"
	"fmt"
	"strconv"
	"strings"
)

type Pulse struct {
	Sender string
	Node   string
	State  bool
}

func netsim(flipflops map[string]bool, conjunctions map[string]map[string]bool, network map[string][]string) (int, int) {
	lowPulses := 0
	highPulses := 0

	for i := 1; i <= 1000; i++ {
		q := list.New()
		q.PushBack(Pulse{Sender: "node", Node: "broadcaster", State: false})

		for q.Len() > 0 {
			element := q.Front()
			pulse := element.Value.(Pulse)
			q.Remove(element)

			if pulse.State {
				highPulses++
			} else {
				lowPulses++
			}

			_, isFlipFlop := flipflops[pulse.Node]
			_, isConjunction := conjunctions[pulse.Node]

			switch {
			case isFlipFlop:
				if !pulse.State {
					flipflops[pulse.Node] = !flipflops[pulse.Node]
					newState := flipflops[pulse.Node]
					for _, nv := range network[pulse.Node] {
						q.PushBack(Pulse{Sender: pulse.Node, Node: nv, State: newState})
					}
				}
			case isConjunction:
				conjunctions[pulse.Node][pulse.Sender] = pulse.State
				newState := !getAll(conjunctions[pulse.Node])
				for _, nv := range network[pulse.Node] {
					q.PushBack(Pulse{Sender: pulse.Node, Node: nv, State: newState})
				}
			case pulse.Node == "broadcaster":
				for _, nv := range network[pulse.Node] {
					q.PushBack(Pulse{Sender: pulse.Node, Node: nv, State: pulse.State})
				}
			}
		}
	}
	return lowPulses, highPulses
}

func getAll(m map[string]bool) bool {
	for _, v := range m {
		if !v {
			return false
		}
	}
	return true
}

func buildBinaryString(module string, graph map[string][]string) string {
	bin := ""
	for {
		if nextModules, ok := graph["%"+module]; ok {
			switch len(nextModules) {
			case 2:
				bin = "1" + bin
			case 1:
				// Add '1' if the next module is not a flip-flop, '0' otherwise.
				if _, exists := graph["%"+nextModules[0]]; !exists {
					bin = "1" + bin
				} else {
					bin = "0" + bin
				}
			default:
				bin = "0" + bin
			}

			// Find the next valid module to continue.
			var nextValidModule string
			for _, nextModule := range nextModules {
				if _, ok := graph["%"+nextModule]; ok {
					nextValidModule = nextModule
					break
				}
			}

			if nextValidModule == "" {
				break
			}
			module = nextValidModule
		} else {
			break
		}
	}
	return bin
}

func calculateLCM(numbers []int) int {
	lcm := numbers[0]
	for _, number := range numbers[1:] {
		lcm = lcmVal(lcm, number)
	}
	return lcm
}

func lcmVal(a, b int) int {
	return a / gcd(a, b) * b
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func Part1(data []string) int {
	lines := data
	network := make(map[string][]string)
	flipflops := make(map[string]bool)
	conjunctions := make(map[string]map[string]bool)

	for _, line := range lines {
		parts := strings.Split(line, " -> ")
		destinations := strings.Split(parts[1], ", ")

		if parts[0] == "broadcaster" {
			network[parts[0]] = destinations
		} else {
			network[parts[0][1:]] = destinations
		}

		switch parts[0][0] {
		case '%':
			flipflops[parts[0][1:]] = false
		case '&':
			conjunctions[parts[0][1:]] = make(map[string]bool)
		}
	}
	for server, clients := range network {
		for _, client := range clients {
			if _, ok := conjunctions[client]; ok {
				conjunctions[client][server] = false
			}
		}
	}

	lowPulses, highPulses := netsim(flipflops, conjunctions, network)
	return lowPulses * highPulses
}

// Part2 calculates the least common multiple of binary numbers derived from a graph.
func Part2(data []string) int {
	graph := make(map[string][]string)

	// Parse each line and populate the graph.
	for _, line := range data {
		parts := strings.Split(line, " -> ")
		graph[parts[0]] = strings.Split(parts[1], ", ")
	}

	var results []int
	for _, m := range graph["broadcaster"] {
		binaryString := buildBinaryString(m, graph)
		value, _ := strconv.ParseInt(binaryString, 2, 64)
		results = append(results, int(value))
	}

	return calculateLCM(results)
}

func main() {
	data, _ := utils.ReadLines("input.txt")

	fmt.Println("Part 1:", Part1(data))
	fmt.Println("Part 2:", Part2(data))
}
