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
		values := make([]CoreValue, len(inputs))
		for i, input := range inputs {
			values[len(inputs)-i-1] = RawToCoreValue(input)
		}
		Programs[name] = SequenceValue{value: values}
		return nil
	}

	if strings.HasSuffix(fileName, ".raw") {
		RawData[name] = data
	}

	return nil
}
