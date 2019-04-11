package main

import (
	"container/list"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Input struct {
	number int
	type_  string
	label  string
	value  string
}

type Sensor struct {
	name   string
	inputs *list.List
}

func toInt(str string) int {
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
		sensor := e.Value.(*Sensor)
		if sensor.name == sensorName {
			return true
		}
	}
	return false
}

func sensorContainsInput(sensors *list.List, sensorName string, inputNumber int) bool {
	for e := sensors.Front(); e != nil; e = e.Next() {
		sensor := e.Value.(*Sensor)
		if sensor.name == sensorName {
			for t := sensor.inputs.Front(); t != nil; t = t.Next() {
				input := t.Value.(*Input)
				if input.number == inputNumber {
					return true
				}
			}
		}
	}
	return false
}

func getSensorFromList(sensors *list.List, sensorName string) *Sensor {
	for e := sensors.Front(); e != nil; e = e.Next() {
		sensor := e.Value.(*Sensor)
		if sensor.name == sensorName {
			return sensor
		}
	}
	return nil
}

func getInputFromSensor(sensors *list.List, sensorName string, inputNumber int) *Input {
	for e := sensors.Front(); e != nil; e = e.Next() {
		sensor := e.Value.(*Sensor)
		if sensor.name == sensorName {
			for t := sensor.inputs.Front(); t != nil; t = t.Next() {
				input := t.Value.(*Input)
				if input.number == inputNumber {
					return input
				}
			}
		}
	}
	return nil
}

func addSensorToList(sensors *list.List, sensorName string) {
	sensor := new(Sensor)
	sensor.name = sensorName
	sensor.inputs = list.New()
	sensors.PushBack(sensor)
}

func addInputToSensor(sensors *list.List, sensorName string, inputNumber int, inputType string) {
	input := new(Input)
	input.number = inputNumber
	input.type_ = inputType
	sensor := getSensorFromList(sensors, sensorName)
	sensor.inputs.PushBack(input)
}

func setValue(input *Input, rawValue string) {
	switch input.type_ {
	case "temp":
		dec := toInt(rawValue)
		float := float64(dec) / 1000
		str := strconv.FormatFloat(float, 'f', 1, 32)
		input.value = str
		break
	default:
		input.value = rawValue
	}
}

func refreshSensorValues(files []os.FileInfo, path string, sensors *list.List) {

	validFile := regexp.MustCompile(`^([a-z]+)([0-9]+)_(input|label)$`)

	name := getContentFile(path + "/name")

	if !containsSensor(sensors, name) {
		addSensorToList(sensors, name)
	}

	for _, file := range files {

		if validFile.MatchString(file.Name()) {
			fileName := path + "/" + file.Name()
			content := getContentFile(fileName)
			subs := validFile.FindSubmatch([]byte(file.Name()))
			number := toInt(string(subs[2]))
			type_ := string(subs[1])
			fileType := string(subs[3])

			if !sensorContainsInput(sensors, name, number) {
				addInputToSensor(sensors, name, number, type_)
			}

			input := getInputFromSensor(sensors, name, number)
			switch fileType {
			case "input":
				setValue(input, content)
				break
			case "label":
				input.label = content
				break
			}
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
