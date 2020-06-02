package read_file

import (
	//		"fmt"
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ParseCPU(dataFile string, components string) {
	var buffer bytes.Buffer

	lines := strings.Split(dataFile, "\n")
	lines = lines[2 : len(lines)-1]

	for index, line := range lines {
		regular, errRegular := regexp.Compile(`\s+`)
		check(errRegular)
		line = regular.ReplaceAllString(line, " ")
		words := strings.Split(line, " ")

		for index, word := range words {
			if index != len(words)-1 {
				buffer.WriteString(word)
				buffer.WriteString(";")
			} else {
				buffer.WriteString(word)
			}
		}

		if index != len(lines)-1 {
			buffer.WriteString("\n")
		}
	}

	pathFile := "stats/statsCPU-" + components + ".csv"
	errWrite := ioutil.WriteFile(pathFile, buffer.Bytes(), 0644)
	check(errWrite)

	buffer.Reset()

}

func ParseMemory(dataFile string, components string) {
	var buffer bytes.Buffer

	lines := strings.Split(dataFile, "\n")
	lines = lines[1 : len(lines)-1]

	for index, line := range lines {
		regular, errRegular := regexp.Compile(`\s+`)
		check(errRegular)
		line = regular.ReplaceAllString(line, " ")
		line = strings.TrimSpace(line)
		words := strings.Split(line, " ")

		for index, word := range words {
			if index != len(words)-1 {
				buffer.WriteString(word)
				buffer.WriteString(";")
			} else {
				buffer.WriteString(word)
			}
		}

		if index != len(lines)-1 {
			buffer.WriteString("\n")
		}
	}

	pathFile := "stats/statsMem-" + components + ".csv"
	errWrite := ioutil.WriteFile(pathFile, buffer.Bytes(), 0644)
	check(errWrite)

	buffer.Reset()

}
