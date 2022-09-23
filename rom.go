package main

import (
	"math"
	"path/filepath"
	"os"
	"strings"
	"fmt"
)

const(
	Flt byte = 'f'
	Hex      = 'h'
	Oct      = 'o'
	Bin      = 'b'
	Str      = 's'
	Nil      = '0'
)

const(
	MaxStackLen = 4
)

var registerMap = map[string]int {
	"A": 0,
		"RA": 0,
		"B": 1,
		"RB": 1,
		"C": 2,
		"RC": 2,
		"D": 3,
		"RD": 3,
	}

var constsMap = map[string]float64 {
	"$pi": math.Pi,
	"$tau": math.Pi * 2,
	"$e": math.E,
	"$phi": math.Phi,
	//"$maxf": math.MaxFloat64, ????
	"$maxf": math.MaxInt64,
}

var progsMap = map[string]string {}


func loadRom() {
	err := filepath.Walk("rom/", loadFile)
	fmt.Printf("\nprogsMap=%s", progsMap)
	if err != nil {
		fmt.Printf("\nFailed to read program: err=%s", err.Error())
	}
}

func loadFile(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	fmt.Printf("\nLoading: path=%s", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	name := strings.Replace(path, "rom/", "", -1)
	name = strings.Replace(name, ".28", "", -1)
	name = strings.ToUpper(name)
	prog := strings.Replace(string(data), "\n", "|", -1)
	fmt.Printf("Loaded: name=%s, prog=%s", name, prog)
	progsMap[name] = prog
	return nil
}
