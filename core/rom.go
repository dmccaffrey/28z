package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var Rom map[string][]byte = make(map[string][]byte)

func LoadRom() error {
	err := filepath.Walk("rom/", loadFile)
	return err
}

func loadFile(path string, info os.FileInfo, err error) error {
	if info.IsDir() || err != nil {
		return nil
	}
	fileName := filepath.Base(path)
	if strings.HasPrefix(fileName, ".") || (!strings.HasSuffix(fileName, ".28") && !strings.HasSuffix(fileName, ".raw")) {
		return nil
	}
	fmt.Printf("\nLoading: path=%s", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	name := strings.Replace(path, "rom/", "", -1)
	name = strings.Replace(name, ".28", "", -1)
	Rom[name] = data
	return nil
}
