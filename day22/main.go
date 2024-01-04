package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Brick struct {
	x, y, z    int
	x2, y2, z2 int
}

func (b *Brick) String() string {
	return fmt.Sprintf("%d,%d,%d~%d,%d,%d", b.x, b.y, b.z, b.x2, b.y2, b.z2)
}

func (b *Brick) collide(b2 *Brick) bool {
	switch {
	case b.z > b2.z2 || b.z2 < b2.z:
		// No overlap in Z-axis
		return false
	case b.y > b2.y2 || b.y2 < b2.y:
		// No overlap in Y-axis
		return false
	case b.x > b2.x2 || b.x2 < b2.x:
		// No overlap in X-axis
		return false
	default:
		// Overlap in all axes
		return true
	}
}
func makeBricks(input []string) ([]Brick, int, bool) {
	var bricks []Brick
	for _, row := range input {
		items := strings.FieldsFunc(row, func(r rune) bool {
			return r == ',' || r == '~'
		})

		// Convert strings to integers
		nums := make([]int, 6)
		var err error
		for i, item := range items {
			nums[i], err = strconv.Atoi(item)
			if err != nil {
				fmt.Printf("Error converting string to int: %v\n", err)
				return nil, -1, true
			}
		}

		brick := Brick{
			x:  nums[0],
			y:  nums[1],
			z:  nums[2],
			x2: nums[3],
			y2: nums[4],
			z2: nums[5],
		}

		bricks = append(bricks, brick)
	}
	slices.SortFunc(bricks, func(a, b Brick) int {
		return min(a.z, a.z2) - min(b.z, b.z2)
	})
	stack(bricks)

	return bricks, 0, false
}

func stack(bricks []Brick) int {
	m := make(map[int]int)
	for i := 0; i < len(bricks); i++ {
		for bricks[i].z > 1 {
			brick := bricks[i]
			brick.z--
			brick.z2--
			collide := false
			for j := i - 1; j > -1; j-- {
				if bricks[j].collide(&brick) {
					collide = true
					break
				}
			}
			if collide {
				break
			}
			m[i]++
			bricks[i] = brick
		}
	}
	// return the delta in the z direction
	return len(m)
}

func part1(input []string) int {
	b, i, done := makeBricks(input)
	if done {
		return i
	}
	// see if removable
	re := 0
	for i := 0; i < len(b); i++ {
		c := make([]Brick, len(b))
		copy(c, b)
		c = append(c[:i], c[i+1:]...)
		delta := stack(c)
		if delta == 0 {
			re++
		}
	}
	return re
}

func part2(input []string) int {
	b, i, done := makeBricks(input)
	if done {
		return i
	}
	sum := 0
	// see if removable
	for i := 0; i < len(b); i++ {
		c := make([]Brick, len(b))
		copy(c, b)
		c = append(c[:i], c[i+1:]...)
		delta := stack(c)
		if delta > 0 {
			sum += delta
		}
	}
	return sum
}
func main() {
	in, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	input := strings.Split(string(in), "\n")
	part1 := part1(input)
	fmt.Println("The answer to part1:", part1)
	part2 := part2(input)
	fmt.Println("The answer to part2:", part2)
}
