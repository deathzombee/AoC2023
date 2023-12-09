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
	seedRanges := readFile("input.txt")
	individualSeeds := readFile2("input.txt")

	// Process individual seeds
	minLocationIndividual := processSeeds(individualSeeds)
	fmt.Println("Minimum location number (individual seeds):", minLocationIndividual)

	// Process seed ranges
	minLocationRange := processSeedRanges(seedRanges)
	fmt.Println("Minimum location number (seed ranges):", minLocationRange)
}

func parseSeeds(line string) []int {
	var seeds []int
	for _, s := range strings.Split(strings.TrimPrefix(line, "seeds: "), " ") {
		seed, err := strconv.Atoi(s)
		if err != nil {
			continue // Skip invalid numbers
		}
		seeds = append(seeds, seed)
	}
	return seeds
}

func parseMapping(line string, currentMappings *[]NumberRangeMapping) {
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

func getMappingSlice(line string) *[]NumberRangeMapping {
	switch line {
	case "seed-to-soil map:":
		return &seedToSoilMappings
	case "soil-to-fertilizer map:":
		return &soilToFertilizerMappings
	case "fertilizer-to-water map:":
		return &fertilizerToWaterMappings
	case "water-to-light map:":
		return &waterToLightMappings
	case "light-to-temperature map:":
		return &lightToTemperatureMappings
	case "temperature-to-humidity map:":
		return &temperatureToHumidityMappings
	case "humidity-to-location map:":
		return &humidityToLocationMappings
	}
	return nil
}

func readFile(filename string) []NumberRangeMapping {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var seedRanges []NumberRangeMapping
	var currentMappings *[]NumberRangeMapping

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasSuffix(line, "map:") {
			currentMappings = getMappingSlice(line)
			*currentMappings = []NumberRangeMapping{}
			continue
		}

		if currentMappings == nil {
			// Assuming logic for parsing seed ranges is different
			// Specific to your seed range parsing needs
			seedRanges = append(seedRanges, parseSeedRanges(line)...)
		} else {
			parseMapping(line, currentMappings)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return seedRanges
}

func parseSeedRanges(line string) []NumberRangeMapping {
	var seedRanges []NumberRangeMapping
	parts := strings.Split(strings.TrimPrefix(line, "seeds: "), " ")

	for i := 0; i < len(parts); i += 2 {
		srcStart, err := strconv.Atoi(parts[i])
		if err != nil {
			continue
		}
		length, err := strconv.Atoi(parts[i+1])
		if err != nil {
			continue
		}
		seedRanges = append(seedRanges, NumberRangeMapping{
			srcStart: srcStart,
			length:   length,
		})
	}
	return seedRanges
}

func readFile2(filename string) []int {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var seeds []int
	var currentMappings *[]NumberRangeMapping

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasSuffix(line, "map:") {
			currentMappings = getMappingSlice(line)
			*currentMappings = []NumberRangeMapping{}
			continue
		}

		if currentMappings == nil {
			seeds = append(seeds, parseSeeds(line)...)
		} else {
			parseMapping(line, currentMappings)
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

func processSeedRanges(ranges []NumberRangeMapping) int {
	minLocation := int(^uint(0) >> 1) // Max int value

	for _, rangeItem := range ranges {
		for i := 0; i < rangeItem.length; i++ {
			seed := rangeItem.srcStart + i
			location := findLocationForSeed(seed)
			if location < minLocation {
				minLocation = location
			}
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
