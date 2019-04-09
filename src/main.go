package main

import (
	"bufio"
	"container/list"
	"os"
	"strconv"
)

func getOffset(sensors *list.List) int {
	size := 0
	for e := sensors.Front(); e != nil; e = e.Next() {
		sensor := e.Value.(*sensor)
		if len(sensor.name) > size {
			size = len(sensor.name)
		}
		for t := sensor.temps.Front(); t != nil; t = t.Next() {
			temp := t.Value.(*temp)
			if len(temp.label)+2 > size {
				size = len(temp.label) + 2
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
		for t := sensor.temps.Front(); t != nil; t = t.Next() {
			temp := t.Value.(*temp)
			if temp.label != "" {
				*y++
				printWindow(term, *y, 2, temp.label)
			}
			printWindow(term, *y, offset, strconv.FormatFloat(temp.value, 'f', 1, 32))
		}
		*y += 2
	}
}

func loop(term *terminal, sensors *list.List) {
	reader := bufio.NewReader(os.Stdin)
	run := true
	for run {
		buffer := make([]byte, 10)
		if _, err := reader.Read(buffer); err != nil {
		}
		if buffer[0] == 'q' {
			run = false
		}

		y := 0
		windowClear(term)
		refreshSensorList(sensors)
		printSensorValues(term, sensors, &y)
		printWindow(term, y, 0, "Press <q> for quit")
		windowRefresh(term)
	}
}

func main() {
	term := new(terminal)
	sensors := list.New()

	saveInitialConfig(term)
	initTerminal()
	initWindow(term, 20, 30)
	nonCanonicalMode()
	loop(term, sensors)
	deleteWindow(term)
	endWindow()
	resetTerminal(term)
}
