package main

import (
	"strings"
	"time"
	"strconv"
	"errors"
	"fmt"
	"math"
)

type StackData struct {
	dataType byte
	str string
	flt float64
}

func DefaultStackData() StackData {
	return StackData{Nil, "", 0}
}

func (d *StackData) Parse(s string) error {
	switch s[0] {
	case 'q':
		time.Sleep(3 * time.Second)
		return nil
	case '"':
		s = strings.TrimSuffix(s, "\"")
		s = strings.TrimPrefix(s, "\"")
		d.dataType = Str
		d.str = s
	case 'x':
		d.dataType = Hex
		s = s[1:]
		res, err := strconv.ParseInt(s, 16, 64)
		if err == nil {
			d.flt = float64(res)
		}
		return err
	case 'o':
		d.dataType = Oct
		s = s[1:]
		res, err := strconv.ParseInt(s, 8, 64)
		if err == nil {
			d.flt = float64(res)
		}
		return err
	default:
		d.dataType = Flt
		res, err := strconv.ParseFloat(s, 64)
		if err == nil {
			d.flt = res
		}
		return err
	}
	return nil
}

func (d StackData) Plus(input StackData) (StackData, error) {
	result := d
	if d.dataType != input.dataType {
		return result, errors.New("Operand data types must match")
	}
	switch d.dataType {
	case Str:
		result.str = d.str + input.str
	default:
		result.flt = d.flt + input.flt
	}
	return result, nil
}

func (d StackData) Mult(input StackData) (StackData, error) {
	if d.dataType == Str || input.dataType == Str {
		return DefaultStackData(), errors.New("Multiplication not defined for strings")
	}
	d.flt *= input.flt
	return d, nil
}

func (d StackData) Div(input StackData) (StackData, error) {
	if d.dataType == Str || input.dataType == Str {
		return DefaultStackData(), errors.New("Division not defined for strings")
	}
	if input.flt == 0 {
		return DefaultStackData(), errors.New("Division by zero is not defined")
	}
	d.flt /= input.flt
	return d, nil
}

func (d StackData) Minus(input StackData) (StackData, error) {
	if d.dataType == Str || input.dataType == Str {
		return StackData{}, errors.New("Subtraction not defined for strings")
	}
	return d.Plus(input.ChS())
}

func (d StackData) Pow(input StackData) (StackData, error) {
	if d.dataType == Str || input.dataType == Str {
		return DefaultStackData(), errors.New("Power not defined for strings")
	}
	d.flt = math.Pow(d.flt, input.flt)
	return d, nil
}

func (d StackData) ChS() StackData {
	d.flt = - d.flt
	return d
}

func (d StackData) ToString() string {
	stackStr := "?"
	switch d.dataType {
	case Str:
		stackStr = fmt.Sprintf("%s", d.str)
	case Flt:
		stackStr = fmt.Sprintf("%.13E", d.flt)
	case Hex:
		stackStr = fmt.Sprintf("%019x", d.flt)
	case Oct:
		stackStr = fmt.Sprintf("%019o", d.flt)
	case Nil:
		stackStr = ""
	}
	if len(stackStr) > 35 {
		stackStr = stackStr[:32] + "..."
	}
	return stackStr
}
