package core

import (
	"strconv"
	"strings"
)

func RawToCoreValue(input string) CoreValue {
	value := RawToImmediateCoreValue(input)
	if value.GetType() != DefaultType {
		return value
	}
	return RawToInstruction(input)
}

func RawToImmediateCoreValue(input string) CoreValue {
	if input == "" {
		return DefaultValue{}
	}

	if len(input) > 1 {
		switch input[0] {
		case '\'':
			input = strings.TrimPrefix(input, "'")
			return StringValue{value: input}
		case '$':
			input = strings.TrimPrefix(input, "$")
			return ReferenceValue{value: input}
		case '%':
			input = strings.TrimPrefix(input, "%")
			ref := ReferenceValue{value: input}
			return ref
		}
	}

	result, err := strconv.ParseFloat(input, 64)
	if err == nil {
		return FloatValue{value: result}
	}

	return DefaultValue{}
}

func RawToInstruction(input string) CoreValue {
	Logger.Printf("Parsing raw to core: input=%s\n", input)
	if input == "" {
		return DefaultValue{}
	}

	instruction, ok := instructionMap[input]
	if ok {
		return InstructionValue{value: instruction}
	}

	return DefaultValue{}
}
