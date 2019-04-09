package main

import (
	"container/list"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
)

type fan struct {
	number int
	label  string
	value  int
}

type temp struct {
	number int
	label  string
	value  float64
}

type sensor struct {
	name  string
	temps *list.List
}

func getNumberTemp(str string) int {
	number, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	return number
}

func getContentFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(content[:len(content)-1])
}

func containsSensor(sensors *list.List, sensorName string) bool {
	for e := sensors.Front(); e != nil; e = e.Next() {
		sensor := e.Value.(*sensor)
		if sensor.name == sensorName {
			return true
		}
	}
	return false
}

func sensorContainsTemp(sensors *list.List, sensorName string, tempNumber int) bool {
	for e := sensors.Front(); e != nil; e = e.Next() {
		sensor := e.Value.(*sensor)
		if sensor.name == sensorName {
			for t := sensor.temps.Front(); t != nil; t = t.Next() {
				temp := t.Value.(*temp)
				if temp.number == tempNumber {
					return true
				}
			}
		}
	}
	return false
}

func getSensorFromList(sensors *list.List, sensorName string) *sensor {
	for e := sensors.Front(); e != nil; e = e.Next() {
		sensor := e.Value.(*sensor)
		if sensor.name == sensorName {
			return sensor
		}
	}
	return nil
}

func getTempFromSensor(sensors *list.List, sensorName string, tempNumber int) *temp {
	for e := sensors.Front(); e != nil; e = e.Next() {
		sensor := e.Value.(*sensor)
		if sensor.name == sensorName {
			for t := sensor.temps.Front(); t != nil; t = t.Next() {
				temp := t.Value.(*temp)
				if temp.number == tempNumber {
					return temp
				}
			}
		}
	}
	return nil
}

func addSensorToList(sensors *list.List, sensorName string) {
	sensor := new(sensor)
	sensor.name = sensorName
	sensor.temps = list.New()
	sensors.PushBack(sensor)
}

func addTempToSensor(sensors *list.List, sensorName string, tempNumber int) {
	temp := new(temp)
	temp.number = tempNumber
	sensor := getSensorFromList(sensors, sensorName)
	sensor.temps.PushBack(temp)
}

func addFanToSensor(sensors *list.List, sensorName string, fanNumber int) {
	fan := new(fan)
	fan.number = fanNumber
	sensor := getSensorFromList(sensors, sensorName)
	sensor.temps.PushBack(fan)
}

func refreshSensorValues(files []os.FileInfo, path string, sensors *list.List) {

	validTempInputFile := regexp.MustCompile(`^temp[0-9]_input$`)
	validTempLabelFile := regexp.MustCompile(`^temp[0-9]_label$`)
	tempInputNumber := regexp.MustCompile(`^temp(\d?)_input$`)
	tempLabelNumber := regexp.MustCompile(`^temp(\d?)_label$`)

	/* validFanInputFile := regexp.MustCompile(`^fan[0-9]_input$`)
	validFanLabelFile := regexp.MustCompile(`^fan[0-9]_label$`)
	fanInputNumber := regexp.MustCompile(`^fan(\d?)_input$`)
	fanLabelNumber := regexp.MustCompile(`^fan(\d?)_label$`) */

	name := getContentFile(path + "/name")

	if !containsSensor(sensors, name) {
		addSensorToList(sensors, name)
	}

	for _, file := range files {

		if validTempInputFile.MatchString(file.Name()) {
			tempFilePath := path + "/" + file.Name()
			content := getContentFile(tempFilePath)
			number := getNumberTemp(string(tempInputNumber.FindSubmatch([]byte(file.Name()))[1]))

			if !sensorContainsTemp(sensors, name, number) {
				addTempToSensor(sensors, name, number)
			}

			temp := getTempFromSensor(sensors, name, number)
			value, err := strconv.ParseFloat(content, 64)
			if err != nil {
				log.Fatal(err)
			}
			temp.value = value / 1000
		}

		if validTempLabelFile.MatchString(file.Name()) {
			labelFilePath := path + "/" + file.Name()
			label := getContentFile(labelFilePath)
			number := getNumberTemp(string(tempLabelNumber.FindSubmatch([]byte(file.Name()))[1]))

			if !sensorContainsTemp(sensors, name, number) {
				addTempToSensor(sensors, name, number)
			}
			temp := getTempFromSensor(sensors, name, number)
			temp.label = label
		}
	}
}

func refreshSensorList(sensors *list.List) {

	links, err := ioutil.ReadDir("/sys/class/hwmon")
	if err != nil {
		log.Fatal(err)
	}

	for _, sensor := range links {
		path := "/sys/class/hwmon/" + sensor.Name()
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}
		refreshSensorValues(files, path, sensors)
	}
}
