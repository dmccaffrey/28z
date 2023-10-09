package main

import (
	"math"
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

var progsMap = map[string]string {
	"MOD": "|?>|`;|;|\"RD|s|\"RD|r|;|\"RC|s|\"RC|r|\"RD|r|\"RC|r|/|i|*|-|",
	"PYTHAG": "|2|^|;|2|^|+|1|2|/|^|",
	"IN2MM": "|25.4|*",
	"SEQR": "|1|;|/|;|1|;|/|+|1|;|/|",
}
