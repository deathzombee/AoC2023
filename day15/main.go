package main

import (
	_ "embed"
	"fmt"
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
		//fmt.Printf("%s: %d\n", v, hashValue)
		added += hashValue
	}
	return added
}

func main() {
	intructions := input
	t1 := time.Now()
	p1 := hashList(intructions)
	fmt.Println("p1 time:", time.Since(t1))
	fmt.Println("p1:", p1)
}
