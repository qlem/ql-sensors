package main

import (
	"container/list"
)

/*
SENSOR TYPE VALUE
Press <q> to quit
*/

const (
	SensorLabelElem1 = "SENSOR "
	SensorLabelElem2 = "TYPE "
	SensorLabelElem3 = "VALUE"
	QuitLabel        = "Press <q> to quit"
)

func getQuitLabel(cols []int) string {
	length := cols[0] + cols[1] + cols[2]
	label := make([]byte, length)
	for i := 0; i < cols[0]+cols[1]+cols[2]-len(QuitLabel); i++ {
		label[i] = ' '
	}
	for i, c := range QuitLabel {
		label[i+length-len(QuitLabel)] = byte(c)
	}
	return string(label)
}

func getLabel(cols []int) string {
	length := cols[0] + cols[1] + cols[2]
	label := make([]byte, length)
	for i, c := range SensorLabelElem1 {
		label[i] = byte(c)
	}
	for i := len(SensorLabelElem1); i < cols[0]; i++ {
		label[i] = ' '
	}
	for i, c := range SensorLabelElem2 {
		label[i+cols[0]] = byte(c)
	}
	for i := cols[0] + len(SensorLabelElem2); i < length-len(SensorLabelElem3); i++ {
		label[i] = ' '
	}
	for i, c := range SensorLabelElem3 {
		label[i+length-len(SensorLabelElem3)] = byte(c)
	}
	return string(label)
}

func printValues(term *Terminal, sensors *list.List, cols []int, y *int) {
	for e := sensors.Front(); e != nil; e = e.Next() {
		sensor := e.Value.(*Sensor)
		setBold(term)
		printWindow(term, *y, 0, sensor.name)
		unsetBold(term)
		for t := sensor.inputs.Front(); t != nil; t = t.Next() {
			input := t.Value.(*Input)
			if sensor.inputs.Len() > 1 {
				*y++
				printWindow(term, *y, 2, input.label)
			}
			printWindow(term, *y, cols[0], input.type_)
			offset := cols[0] + cols[1] + cols[2] - len(input.value)
			printWindow(term, *y, offset, input.value)
		}
		*y++
	}
}

func computeColumnsWidth(sensors *list.List, cols []int) {
	cols[0] = len(SensorLabelElem1)
	cols[1] = len(SensorLabelElem2)
	cols[2] = len(SensorLabelElem3)
	for e := sensors.Front(); e != nil; e = e.Next() {
		sensor := e.Value.(*Sensor)
		if len(sensor.name)+1 > cols[0] {
			cols[0] = len(sensor.name) + 1
		}
		for t := sensor.inputs.Front(); t != nil; t = t.Next() {
			input := t.Value.(*Input)
			if len(input.label)+3 > cols[0] {
				cols[0] = len(input.label) + 3
			}
			if len(input.type_)+1 > cols[1] {
				cols[1] = len(input.type_) + 1
			}
			if len(input.value) > cols[2] {
				cols[2] = len(input.value)
			}
		}
	}
}

func print_(term *Terminal, sensors *list.List) {
	cols := make([]int, 3)
	computeColumnsWidth(sensors, cols)
	setReverse(term)
	printWindow(term, 0, 0, getLabel(cols))
	unsetReverse(term)
	y := 1
	printValues(term, sensors, cols, &y)
	setReverse(term)
	printWindow(term, y, 0, getQuitLabel(cols))
	unsetReverse(term)
}
