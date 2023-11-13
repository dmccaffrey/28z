package core

import (
	"fmt"
	"strconv"
	"strings"
)

type CoreValueType int

const (
	FloatType       CoreValueType = 1
	StringType                    = 2
	SequenceType                  = 3
	InstructionType               = 4
	ReferenceType                 = 5
	DefaultType                   = 6
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
		DefaultValue
		value float64
	}
	StringValue struct {
		DefaultValue
		value string
	}
	SequenceValue struct {
		DefaultValue
		value []CoreValue
	}
	InstructionValue struct {
		DefaultValue
		value Instruction
	}
	ReferenceValue struct {
		DefaultValue
		value string
	}
)

// Default
func (d DefaultValue) GetFloat() float64 {
	return 0
}

func (d DefaultValue) GetString() string {
	return "nil"
}

func (d DefaultValue) GetInt() int {
	return 0
}

func (d DefaultValue) GetType() CoreValueType {
	return DefaultType
}

func (d DefaultValue) GetSequence() []CoreValue {
	return []CoreValue{}
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

func (f FloatValue) GetSequence() []CoreValue {
	return []CoreValue{f}
}

// String
func (s StringValue) GetString() string {
	return s.value
}

func (s StringValue) GetType() CoreValueType {
	return StringType
}

func (s StringValue) GetSequence() []CoreValue {
	return []CoreValue{s}
}

// Sequence
func (s SequenceValue) GetString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%d]:", len(s.value)))
	for _, v := range s.value {
		if v.GetType() == InstructionType {
			sb.WriteString("%")

		} else if v.GetType() == SequenceType {
			sb.WriteString(fmt.Sprintf("[%d]", len(v.GetSequence())))

		} else {
			sb.WriteString(v.GetString())
		}
		sb.WriteByte(',')
		if sb.Len() > 40 {
			sb.WriteString("...")
			break
		}
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

func (s InstructionValue) GetSequence() []CoreValue {
	return []CoreValue{s}
}

func (s InstructionValue) Eval(core *Core) InstructionResult {
	return s.value.impl(core)
}

// Reference
func (r ReferenceValue) GetType() CoreValueType {
	return ReferenceType
}

func (r ReferenceValue) GetSequence() []CoreValue {
	return []CoreValue{r}
}

func (r ReferenceValue) Dereference(core *Core) CoreValue {
	reg, ok := core.GetRegisterMap()[r.value]
	if ok {
		return FloatValue{value: float64(reg)}
	}
	variable, ok := Variables[r.value]
	if ok {
		return variable
	}
	variable, ok = Programs[r.value]
	if ok {
		return variable
	}
	return DefaultValue{}
}

func (r ReferenceValue) GetString() string {
	return "REF/" + r.value
}
