package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type NumberRangeMapping struct {
	srcStart  int
	destStart int
	length    int
}

var seedToSoilMappings []NumberRangeMapping
var soilToFertilizerMappings []NumberRangeMapping
var fertilizerToWaterMappings []NumberRangeMapping
var waterToLightMappings []NumberRangeMapping
var lightToTemperatureMappings []NumberRangeMapping
var temperatureToHumidityMappings []NumberRangeMapping
var humidityToLocationMappings []NumberRangeMapping

func main() {
	seeds := readFile("input.txt")

	minLocation := processSeeds(seeds)

	fmt.Println("Minimum location number:", minLocation)
}

func readFile(filename string) []int {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)

	var seeds []int
	var currentMappings *[]NumberRangeMapping

	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line is a header for a new mapping
		if strings.HasSuffix(line, "map:") {
			switch line {
			case "seed-to-soil map:":
				currentMappings = &seedToSoilMappings
			case "soil-to-fertilizer map:":
				currentMappings = &soilToFertilizerMappings
			case "fertilizer-to-water map:":
				currentMappings = &fertilizerToWaterMappings
			case "water-to-light map:":
				currentMappings = &waterToLightMappings
			case "light-to-temperature map:":
				currentMappings = &lightToTemperatureMappings
			case "temperature-to-humidity map:":
				currentMappings = &temperatureToHumidityMappings
			case "humidity-to-location map:":
				currentMappings = &humidityToLocationMappings
			}
			*currentMappings = []NumberRangeMapping{} // Initialize the current mapping slice
			continue
		}

		// Process the line based on its type
		if currentMappings == nil { // Seeds line
			for _, s := range strings.Split(strings.TrimPrefix(line, "seeds: "), " ") {
				seed, err := strconv.Atoi(s)
				if err != nil {
					continue // Skip invalid numbers
				}
				seeds = append(seeds, seed)
			}
		} else { // Mapping line
			parts := strings.Split(line, " ")
			if len(parts) == 3 {
				srcStart, _ := strconv.Atoi(parts[1])
				destStart, _ := strconv.Atoi(parts[0])
				length, _ := strconv.Atoi(parts[2])
				*currentMappings = append(*currentMappings, NumberRangeMapping{
					srcStart:  srcStart,
					destStart: destStart,
					length:    length,
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return seeds
}

func processSeeds(seeds []int) int {
	minLocation := int(^uint(0) >> 1) // Max int value

	for _, seed := range seeds {
		location := findLocationForSeed(seed)
		if location < minLocation {
			minLocation = location
		}
	}

	return minLocation
}

func findLocationForSeed(seed int) int {
	soil := getMappedValue(seed, seedToSoilMappings)
	fertilizer := getMappedValue(soil, soilToFertilizerMappings)
	water := getMappedValue(fertilizer, fertilizerToWaterMappings)
	light := getMappedValue(water, waterToLightMappings)
	temperature := getMappedValue(light, lightToTemperatureMappings)
	humidity := getMappedValue(temperature, temperatureToHumidityMappings)
	location := getMappedValue(humidity, humidityToLocationMappings)

	return location
}

func getMappedValue(num int, mappings []NumberRangeMapping) int {
	for _, mapping := range mappings {
		if num >= mapping.srcStart && num < mapping.srcStart+mapping.length {
			return mapping.destStart + (num - mapping.srcStart)
		}
	}
	return num // If not mapped, it maps to itself
}
