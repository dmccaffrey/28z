package core

import (
	"strconv"
	"strings"
)

type CoreValueType int

const (
	FloatType       CoreValueType = 1
	StringType                    = 2
	SequenceType                  = 3
	InstructionType               = 4
	DefaultType                   = 5
)

type (
	CoreValue interface {
		GetType() CoreValueType
		GetFloat() float64
		GetString() string
		GetInt() int
		GetSequence() []CoreValue
	}
	DefaultValue struct{}
	FloatValue   struct {
		*DefaultValue
		value float64
	}
	StringValue struct {
		*DefaultValue
		value string
	}
	SequenceValue struct {
		*DefaultValue
		value []CoreValue
	}
	InstructionValue struct {
		*DefaultValue
		value Instruction
	}
)

// Default
func (d DefaultValue) GetFloat() float64 {
	return 0
}

func (d DefaultValue) GetString() string {
	return ""
}

func (d DefaultValue) GetInt() int {
	return 0
}

func (d DefaultValue) GetType() CoreValueType {
	return DefaultType
}

func (d DefaultValue) GetSequence() []CoreValue {
	return []CoreValue{d}
}

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
func (s StringValue) GetString() string {
	return s.value
}

func (s StringValue) GetType() CoreValueType {
	return StringType
}

// Sequence
func (s SequenceValue) GetString() string {
	var sb strings.Builder
	for _, v := range s.value {
		sb.WriteString(v.GetString())
		sb.WriteByte('|')
	}
	return sb.String()
}

func (s SequenceValue) GetType() CoreValueType {
	return SequenceType
}

func (s SequenceValue) GetSequence() []CoreValue {
	return s.value
}

// Instruction
func (s InstructionValue) GetString() string {
	return s.value.description
}

func (s InstructionValue) GetType() CoreValueType {
	return InstructionType
}

func (s InstructionValue) CheckArgs(core *Core) bool {
	if core.currentStack().length >= s.value.argCount {
		return true
	}
	return false
}

func (s InstructionValue) Eval(core *Core) InstructionResult {
	return s.value.impl(core)
}
