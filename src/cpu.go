package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

const CpuFile = "/proc/stat"

type Cpu struct {
	prevIdle  int
	prevTotal int
	usage     int
}

func computeUsage(tokens [][]byte, cpus []Cpu, index int) {
	idle := toInt(string(tokens[3])) + toInt(string(tokens[4]))
	total := 0
	// TODO remove guest and guest_nice
	for _, stat := range tokens {
		total += toInt(string(stat))
	}
	diffIdle := idle - (cpus)[index].prevIdle
	diffTotal := total - (cpus)[index].prevTotal
	cpus[index].usage = (1000 * (diffTotal - diffIdle) / diffTotal) / 10
	cpus[index].prevIdle = idle
	cpus[index].prevTotal = total
}

func countDigit(line string, i *int) int {
	count := 0
	for line[*i] >= '0' && line[*i] <= '9' {
		count++
		*i++
	}
	return count
}

func parseLine(line string) [][]byte {
	j := 0
	tmp := 0
	tokens := make([][]byte, 10)
	for i := 0; i < len(line); i++ {
		tmp = i
		if line[i] >= '0' && line[i] <= '9' && (i > 0 && line[i-1] == ' ') {
			size := countDigit(line, &i)
			tokens[j] = make([]byte, size)
			for k := 0; k < size; k++ {
				tokens[j][k] = line[tmp]
				tmp++
			}
			j++
		}
	}
	return tokens
}

func refreshCpus(cpus []Cpu) []Cpu {
	file, err := os.Open(CpuFile)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)
	i := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if !strings.HasPrefix(line, "cpu") {
			break
		}
		tokens := parseLine(line)
		if i >= len(cpus) {
			cpus = append(cpus, Cpu{
				prevIdle:  0,
				prevTotal: 0,
				usage:     0,
			})
		}
		computeUsage(tokens, cpus, i)
		i++
	}
	return cpus
}
