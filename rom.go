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

var uFuncs = map[string]UnaryFunc {
	"@sin": func (x StackData) (StackData, error) {
		x.flt = math.Sin(x.flt)
		return x, nil
	},
		"@cos": func (x StackData) (StackData, error) {
			x.flt = math.Cos(x.flt)
			return x, nil
		},
		"@tan": func (x StackData) (StackData, error) {
			x.flt = math.Tan(x.flt)
			return x, nil
		},
		"@log": func (x StackData) (StackData, error) {
			x.flt = math.Log10(x.flt)
			return x, nil
		},
		"@ln": func (x StackData) (StackData, error) {
			x.flt = math.Log(x.flt)
			return x, nil
		},
		"@logb": func (x StackData) (StackData, error) {
			x.flt = math.Logb(x.flt)
			return x, nil
		},
}

var bFuncs = map[string]BinaryFunc {
}

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
