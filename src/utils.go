package main

import (
	"log"
	"strconv"
)

func countDigit(line string, i int) int {
	count := 0
	for i < len(line) && (line[i] >= '0' && line[i] <= '9') {
		count++
		i++
	}
	return count
}

func toInt(str string) int {
	number, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	return number
}
