package main

import (
	"container/list"
	"fmt"
	"log"
	"strings"
)

/*
CPUs USAGE              TOTAL [100]
[100] [100] [100] [100] [100] [100]
[100] [100]
MEM      USED     TOTAL           %
 ram       8.00    16.00G     [ 50]
 swap      0.00     8.00G     [  0]
SENSOR            TYPE        VALUE
                  Press <q> to quit
*/

const (
	TotalWidth          = 35
	SensorOffset1       = 2
	SensorOffset2       = 19
	CoreLength          = 5
	CpuTotalValueOffset = 31
	CpuLabel            = "CPUs USAGE              TOTAL [   ]"
	MemLabel            = "MEM      USED     TOTAL           %"
	SensorLabel         = "SENSOR            TYPE        VALUE"
	QuitLabel           = "                  Press <q> to quit"
)

func printSensorsValues(term *Terminal, sensors *list.List, y *int) {
	setReverse(term)
	printWindow(term, *y, 0, SensorLabel)
	unsetReverse(term)
	*y++
	// TODO truncate if sensor name / type / value > column
	for e := sensors.Front(); e != nil; e = e.Next() {
		sensor := e.Value.(*Sensor)
		setBold(term)
		printWindow(term, *y, 0, sensor.name)
		unsetBold(term)
		for t := sensor.inputs.Front(); t != nil; t = t.Next() {
			input := t.Value.(*Input)
			if sensor.inputs.Len() > 1 {
				*y++
				printWindow(term, *y, SensorOffset1, input.label)
			}
			printWindow(term, *y, SensorOffset2, input.type_)
			offset := TotalWidth - len(input.value)
			printWindow(term, *y, offset, input.value)
		}
		*y++
	}
}

func printMemInfo(term *Terminal, y *int) {
	x := 0
	setReverse(term)
	printWindow(term, *y, x, MemLabel)
	unsetReverse(term)

	*y++
	printWindow(term, *y, x, " ram       8.00    16.00G     [ 50]")
	*y++
	printWindow(term, *y, x, " swap      0.00     8.00G     [  0]")
	*y++
}

func printCpuUsage(term *Terminal, cpus []Cpu, y *int) {
	setReverse(term)
	printWindow(term, 0, 0, CpuLabel)
	var usage strings.Builder
	if _, err := fmt.Fprintf(&usage, "%3d", cpus[0].usage); err != nil {
		log.Fatal(err)
	}
	printWindow(term, 0, CpuTotalValueOffset, usage.String())
	unsetReverse(term)
	x := 0
	*y = 1
	for i := 1; i < len(cpus); i++ {
		printWindow(term, *y, x, "[   ]")
		x++
		usage.Reset()
		if _, err := fmt.Fprintf(&usage, "%3d", cpus[i].usage); err != nil {
			log.Fatal(err)
		}
		printWindow(term, *y, x, usage.String())
		if x+CoreLength > TotalWidth {
			*y++
			x = 0
		} else {
			x += CoreLength
		}
	}
	*y++
}

func printValues(term *Terminal, sensors *list.List, cpus []Cpu) {
	y := 0
	printCpuUsage(term, cpus, &y)
	printMemInfo(term, &y)
	printSensorsValues(term, sensors, &y)
	setReverse(term)
	printWindow(term, y, 0, QuitLabel)
	unsetReverse(term)
}
