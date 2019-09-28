package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

const (
	MemFile            = "/proc/meminfo"
	MemTotalPrefix     = "MemTotal"
	MemFreePrefix      = "MemFree"
	BuffersPrefix      = "Buffers"
	CachedPrefix       = "Cached"
	SReclaimablePrefix = "SReclaimable"
	ShmemPrefix        = "Shmem"
	SwapTotalPrefix    = "SwapTotal"
	SwapFreePrefix     = "SwapFree"
	KbConversion       = 1048576.0
)

type Mem struct {
	memTotal     float64
	memUsed      float64
	memFree      float64
	buffers      float64
	cached       float64
	sReclaimable float64
	shMem        float64
	memAmount    float64
	swapTotal    float64
	swapUsed     float64
	swapFree     float64
	swapAmount   float64
}

type FooFunc func(mem *Mem, line string)

type Func_ struct {
	func_  FooFunc
	prefix string
}

func getMemValue(line string) float64 {
	value := 0.0
	for i, c := range line {
		if c >= '0' && c <= '9' {
			length := countDigit(line, i)
			raw := make([]byte, length)
			for j := 0; j < len(raw); j++ {
				raw[j] = line[i]
				i++
			}
			float := float64(toInt(string(raw)))
			return float / KbConversion
		}
	}
	return value
}

func memTotal(mem *Mem, line string) {
	mem.memTotal = getMemValue(line)
}

func memFree(mem *Mem, line string) {
	mem.memFree = getMemValue(line)
}

func memBuffered(mem *Mem, line string) {
	mem.buffers = getMemValue(line)
}

func memCached(mem *Mem, line string) {
	mem.cached = getMemValue(line)
}

func memSReclaimable(mem *Mem, line string) {
	mem.sReclaimable = getMemValue(line)
}

func memShmem(mem *Mem, line string) {
	mem.shMem = getMemValue(line)
}

func swapTotal(mem *Mem, line string) {
	mem.swapTotal = getMemValue(line)
}

func swapFree(mem *Mem, line string) {
	mem.swapFree = getMemValue(line)
}

func computeMemValues(mem *Mem) {
	mem.cached = mem.cached + mem.sReclaimable - mem.shMem
	mem.memUsed = mem.memTotal - mem.memFree - mem.buffers - mem.cached
	mem.memAmount = (mem.memUsed * 100) / mem.memTotal
	mem.swapUsed = mem.swapTotal - mem.swapFree
	mem.swapAmount = (mem.swapUsed * 100) / mem.swapTotal
}

func refreshMem(mem *Mem) {

	array := make([]Func_, 8)
	array[0].func_ = memTotal
	array[0].prefix = MemTotalPrefix
	array[1].func_ = memFree
	array[1].prefix = MemFreePrefix
	array[2].func_ = swapTotal
	array[2].prefix = SwapTotalPrefix
	array[3].func_ = swapFree
	array[3].prefix = SwapFreePrefix
	array[4].func_ = memBuffered
	array[4].prefix = BuffersPrefix
	array[5].func_ = memCached
	array[5].prefix = CachedPrefix
	array[6].func_ = memSReclaimable
	array[6].prefix = SReclaimablePrefix
	array[7].func_ = memShmem
	array[7].prefix = ShmemPrefix

	file, err := os.Open(MemFile)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}
		for _, elem := range array {
			if strings.HasPrefix(line, elem.prefix) {
				elem.func_(mem, line)
			}
		}
	}
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
	computeMemValues(mem)
}
