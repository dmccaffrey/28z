package core

import (
	"strconv"
)

type CoreValueType int

const (
	FloatType  CoreValueType = 1
	StringType               = 2
)

type (
	CoreValue interface {
		GetType() CoreValueType
		GetFloat() float64
		GetString() string
		GetInt() int
	}
	FloatValue struct {
		value float64
	}
	StringValue struct {
		value string
	}
)

// Float
func (f FloatValue) GetFloat() float64 {
	return f.value
}

func (f FloatValue) GetString() string {
	return strconv.FormatFloat(f.value, 'E', -1, 64)
}

func (f FloatValue) GetInt() int {
	return int(f.value)
}

func (f FloatValue) GetType() CoreValueType {
	return FloatType
}

// String
func (s StringValue) GetFloat() float64 {
	return 0
}

func (s StringValue) GetString() string {
	return s.value
}

func (s StringValue) GetInt() int {
	return 0
}

func (s StringValue) GetType() CoreValueType {
	return StringType
}
