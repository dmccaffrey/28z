package core

import (
	"os"
	"path/filepath"
	"strings"
)

var RawData map[string][]byte = make(map[string][]byte)
var Programs map[string]SequenceValue = make(map[string]SequenceValue)

func LoadRom() error {
	err := filepath.Walk("rom/", loadFile)
	return err
}

func loadFile(path string, info os.FileInfo, err error) error {
	if info.IsDir() || err != nil {
		return nil
	}
	fileName := filepath.Base(path)
	if strings.HasPrefix(fileName, ".") {
		return nil
	}
	name := strings.Replace(path, "rom/", "", -1)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	if strings.HasSuffix(fileName, ".28") {
		name = strings.Replace(name, ".28", "", -1)
		inputs := strings.Split(string(data[:]), "\n")
		_, result := convertToSequence(0, inputs)
		Programs[name] = result
		return nil
	}

	if strings.HasSuffix(fileName, ".raw") {
		RawData[name] = data
	}

	return nil
}

func convertToSequence(offset int, inputs []string) (int, SequenceValue) {
	values := []CoreValue{}
	for ; offset < len(inputs); offset++ {
		input := strings.TrimLeft(inputs[offset], " \t")
		if input == "<" {
			newOffset, value := convertToSequence(offset+1, inputs)
			offset = newOffset
			values = append([]CoreValue{value}, values...)

		} else if input == ">" {
			return offset + 1, SequenceValue{value: values}

		} else {
			value := RawToCoreValue(input)
			values = append([]CoreValue{value}, values...)
		}
	}
	return offset, SequenceValue{value: values}
}
