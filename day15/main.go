package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//go:embed input.txt
var input string

// Holiday ASCII String Helper Alogrithm
func hash(s string) int32 {
	var current int32 = 0
	for _, c := range s {
		//fmt.Println(c)
		current += c
		current = current * 17
		current = current % 256
	}
	return current
}
func hashList(s string) int32 {
	var added int32 = 0
	list := strings.Split(s, ",")
	for _, v := range list {
		hashValue := hash(v)
		added += hashValue
	}
	return added
}
func hashmap(s string) map[int32][]string {
	m := make(map[int32][]string)
	list := strings.Split(s, ",")

	for _, v := range list {
		label := v
		value := ""

		if idx := strings.Index(v, "="); idx != -1 {
			label = v[:idx]
			value = v[idx+1:]
			hashValue := hash(label)
			found := false
			for i, val := range m[hashValue] {
				if strings.HasPrefix(val, label+" ") {
					// Replace existing lens
					m[hashValue][i] = label + " " + value
					found = true
					break
				}
			}
			if !found {
				// Add new lens
				m[hashValue] = append(m[hashValue], label+" "+value)
			}
		} else if idx := strings.Index(v, "-"); idx != -1 {
			label = v[:idx]
			hashValue := hash(label)
			// Remove the lens and shift others forward
			var newBox []string
			for _, val := range m[hashValue] {
				if !strings.HasPrefix(val, label+" ") {
					newBox = append(newBox, val)
				}
			}
			m[hashValue] = newBox
		}
	}
	return m
}
func focusingPower(m map[int32][]string) int {
	totalPower := 0

	for box, lenses := range m {
		for i, lens := range lenses {
			parts := strings.Fields(lens)
			if len(parts) < 2 {
				continue // Skip if the lens does not have a focal length
			}
			focalLength, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Error converting focal length:", parts[1])
				continue
			}
			power := (int(box) + 1) * (i + 1) * focalLength
			totalPower += power
		}
	}
	return totalPower
}
func main() {
	intructions := input
	t1 := time.Now()
	p1 := hashList(intructions)
	fmt.Println("p1 time:", time.Since(t1))
	fmt.Println("p1:", p1)
	t2 := time.Now()
	p2 := hashmap(intructions)
	totalPower := focusingPower(p2)
	fmt.Println("p2 time:", time.Since(t2))
	fmt.Println("Total Focusing Power:", totalPower)
}
