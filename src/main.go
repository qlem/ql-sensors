package main

import (
	"bufio"
	"container/list"
	"os"
	"strconv"
	"time"
)

func getOffset(sensors *list.List) int {
	size := 0
	for e := sensors.Front(); e != nil; e = e.Next() {
		sensor := e.Value.(*sensor)
		if len(sensor.name) > size {
			size = len(sensor.name)
		}
		for t := sensor.inputs.Front(); t != nil; t = t.Next() {
			input := t.Value.(*input)
			if len(input.label)+2 > size {
				size = len(input.label) + 2
			}
		}
	}
	return size
}

func printSensorValues(term *terminal, sensors *list.List, y *int) {
	offset := getOffset(sensors) + 2
	for e := sensors.Front(); e != nil; e = e.Next() {
		sensor := e.Value.(*sensor)
		setBold(term)
		printWindow(term, *y, 0, sensor.name)
		unsetBold(term)
		for t := sensor.inputs.Front(); t != nil; t = t.Next() {
			input := t.Value.(*input)
			if input.label != "" {
				*y++
				printWindow(term, *y, 2, input.label)
			}
			switch input.type_ {
			case TEMP:
				value := float64(input.value) / 1000
				str := strconv.FormatFloat(value, 'f', 1, 32)
				printWindow(term, *y, offset, str)
				break
			default:
				printWindow(term, *y, offset, strconv.Itoa(input.value))
			}
		}
		*y++
	}
}

func loop(term *terminal, sensors *list.List) {
	reader := bufio.NewReader(os.Stdin)
	run := true
	ref := time.Unix(0, 0)
	for run {
		buffer := make([]byte, 10)
		if _, err := reader.Read(buffer); err != nil {
		}
		if buffer[0] == 'q' {
			run = false
		}

		now := time.Now()
		delta := now.Sub(ref)
		if delta.Seconds() >= 1 {
			y := 0
			windowClear(term)
			refreshSensorList(sensors)
			printSensorValues(term, sensors, &y)
			y++
			printWindow(term, y, 0, "Press <q> for quit")
			windowRefresh(term)
			ref = time.Now()
		}
	}
}

func main() {
	term := new(terminal)
	sensors := list.New()

	saveInitialConfig(term)
	initTerminal()
	initWindow(term, 0, 0)
	nonCanonicalMode()
	loop(term, sensors)
	deleteWindow(term)
	endWindow()
	resetTerminal(term)
}
