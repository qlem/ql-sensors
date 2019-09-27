package main

import (
	"container/list"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
)

const SensorDir = "/sys/class/hwmon"

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

func setValue(input *Input, rawValue string) {
	if input.type_ == "temp" && rawValue != "N/A" {
		dec := toInt(rawValue)
		float := float64(dec) / 1000
		str := strconv.FormatFloat(float, 'f', 0, 32)
		input.value = str
	} else {
		input.value = rawValue
	}
}

func getContentFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		// Handle the read error from a corrupted input file
		return "N/A"
	}
	return string(content[:len(content)-1])
}

func addInput(sensor *Sensor, inputNumber int, inputType string) *Input {
	input := new(Input)
	input.number = inputNumber
	input.type_ = inputType
	sensor.inputs.PushBack(input)
	return input
}

func getInput(sensor *Sensor, inputNumber int, inputType string) *Input {
	for t := sensor.inputs.Front(); t != nil; t = t.Next() {
		input := t.Value.(*Input)
		if input.number == inputNumber && input.type_ == inputType {
			return input
		}
	}
	return addInput(sensor, inputNumber, inputType)
}

func addSensor(sensors *list.List, sensorName string) *Sensor {
	sensor := new(Sensor)
	sensor.name = sensorName
	sensor.inputs = list.New()
	sensors.PushBack(sensor)
	return sensor
}

func getSensor(sensors *list.List, sensorName string) *Sensor {
	for e := sensors.Front(); e != nil; e = e.Next() {
		sensor := e.Value.(*Sensor)
		if sensor.name == sensorName {
			return sensor
		}
	}
	return addSensor(sensors, sensorName)
}

func refreshSensorValues(files []os.FileInfo, path string, sensors *list.List) {

	validFile := regexp.MustCompile(`^([a-z]+)([0-9]+)_(input|label)$`)
	name := getContentFile(path + "/name")
	sensor := getSensor(sensors, name)

	for _, file := range files {

		if validFile.MatchString(file.Name()) {
			fileName := path + "/" + file.Name()
			content := getContentFile(fileName)
			subs := validFile.FindSubmatch([]byte(file.Name()))
			number := toInt(string(subs[2]))
			inputType := string(subs[1])
			fileType := string(subs[3])

			input := getInput(sensor, number, inputType)
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

	links, err := ioutil.ReadDir(SensorDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, sensor := range links {
		path := SensorDir + "/" + sensor.Name()
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}
		refreshSensorValues(files, path, sensors)
	}
}
