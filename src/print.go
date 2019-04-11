package main

import (
	"container/list"
)

func computeColumnWidths(sensors *list.List) (width1 int, width2 int, width3 int) {
	width1 = len("SENSOR ")
	width2 = len("TYPE ")
	width3 = len("VALUE")
	for e := sensors.Front(); e != nil; e = e.Next() {
		sensor := e.Value.(*Sensor)
		if len(sensor.name)+1 > width1 {
			width1 = len(sensor.name) + 1
		}
		for t := sensor.inputs.Front(); t != nil; t = t.Next() {
			input := t.Value.(*Input)
			if len(input.label)+3 > width1 {
				width1 = len(input.label) + 3
			}
			if len(input.type_)+1 > width2 {
				width2 = len(input.type_) + 1
			}
			if len(input.value) > width3 {
				width3 = len(input.value)
			}
		}
	}
	return width1, width2, width3
}

func computeFirstRow(width1 int, width2 int, width3 int) string {
	row := make([]byte, width1+width2+width3)
	for i := range row {
		row[i] = ' '
	}
	for i, c := range "SENSOR" {
		row[i] = byte(c)
	}
	i := width1
	for _, c := range "TYPE" {
		row[i] = byte(c)
		i++
	}
	offset := width3 - len("VALUE")
	i = width1 + width2 + offset
	for _, c := range "VALUE" {
		row[i] = byte(c)
		i++
	}
	return string(row)
}

func computeLastRow(width1 int, width2 int, width3 int) string {
	str := "press <q> to exit"
	firstRowLen := width1 + width2 + width3
	offset := firstRowLen - len(str)
	row := make([]byte, width1+width2+width3)
	for i := range row {
		row[i] = ' '
	}
	i := offset
	for _, c := range str {
		row[i] = byte(c)
		i++
	}
	return string(row)
}

func printValues(term *Terminal, sensors *list.List) {
	width1, width2, width3 := computeColumnWidths(sensors)
	setReverse(term)
	printWindow(term, 0, 0, computeFirstRow(width1, width2, width3))
	unsetReverse(term)
	y := 1
	for e := sensors.Front(); e != nil; e = e.Next() {
		sensor := e.Value.(*Sensor)
		setBold(term)
		printWindow(term, y, 0, sensor.name)
		unsetBold(term)
		for t := sensor.inputs.Front(); t != nil; t = t.Next() {
			input := t.Value.(*Input)
			if input.label != "" {
				y++
				printWindow(term, y, 2, input.label)
			}
			printWindow(term, y, width1, input.type_)
			offset := width3 - len(input.value)
			printWindow(term, y, width1+width2+offset, input.value)
		}
		y++
	}
	setReverse(term)
	printWindow(term, y, 0, computeLastRow(width1, width2, width3))
	unsetReverse(term)
}
