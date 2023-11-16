package core

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var RawData map[string][]byte = make(map[string][]byte)
var Programs map[string]SequenceValue = make(map[string]SequenceValue)
var Symbols []rune = []rune{}

func LoadRom() error {
	err := filepath.Walk("rom/", loadFile)
	loadSymbols()
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
		n := 0
		for _, val := range data {
			if val != 10 {
				data[n] = val - 32
				n++
			}
		}
		RawData[name] = []byte(data[:n])
	}

	return nil
}

func loadSymbols() {
	raw, err := os.ReadFile("rom/symbols.set")
	if err != nil {
		log.Fatal("Could not read symbols")
	}
	Symbols = bytes.Runes(raw)
	fmt.Printf("Symbols=%s", string(Symbols[:]))
}

func convertToSequence(offset int, inputs []string) (int, SequenceValue) {
	values := []CoreValue{}
	for ; offset < len(inputs); offset++ {
		input := strings.TrimLeft(inputs[offset], " \t")
		if input == "" {
			continue
		}
		if input[0] == '#' {
			continue
		}
		if input == "<" {
			newOffset, value := convertToSequence(offset+1, inputs)
			offset = newOffset
			values = append([]CoreValue{value}, values...)

		} else if input == ">" {
			return offset, SequenceValue{value: values}

		} else {
			value, err := RawToCoreValue(input, nil)
			if err != nil {
				Logger.Printf("Error: sequence contains invalid value: offset=%d, inputs=%s...", offset, inputs[0])
			}
			values = append([]CoreValue{value}, values...)
		}
	}
	return offset, SequenceValue{value: values}
}
