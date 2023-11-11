package core

import (
	"fmt"
)

var successResult = InstructionResult{false, ""}

type InstructionImpl func(*Core) InstructionResult

type Instruction struct {
	description string
	argCount    int
	resultCount int
	impl        InstructionImpl
	usage       string
}

type InstructionResult struct {
	error   bool
	message string
}

func (i *Instruction) IsValid() bool {
	return i.description != ""
}

var instructionMap = map[string]Instruction{
	"+":        {"Add x and y", 2, 1, add, "6 ⤶ 2 ⤶ + ⤶ ⤒8"},
	"-":        {"Subtract x from y", 2, 1, subtract, "6 ⤶ 2 ⤶ - ⤶ ⤒4"},
	"*":        {"Multiply y by x", 2, 1, multiply, "6 ⤶ 2 ⤶ * ⤶ ⤒12"},
	"/":        {"Divide y by x", 2, 1, divide, "6 ⤶ 2 ⤶ / ⤶ ⤒3"},
	"mod":      {"y modulus by x", 2, 1, modulus, "6 ⤶ 2 ⤶ / ⤶ ⤒0"},
	"inverse":  {"Inverts x", 1, 1, inverse, ""},
	"<":        {"Define sequence", 0, 0, defineSequence, ""},
	">":        {"Define sequence", 0, 0, reduceSequence, ""},
	"this":     {"Refer to the current sequence", 0, 1, this, ""},
	"eval":     {"Evaluate x", 1, 0, nil, ""},
	"consume":  {"Pop from previous stack and push to current", 0, 1, consume, ""},
	"produce":  {"Pop from this stack and push to previous", 1, 0, produce, ""},
	"apply":    {"Evalue x against all entries in y", 2, 1, apply, ""},
	"reduce":   {"Use x to reduce y to a single value", 2, 1, reduce, ""},
	"enter":    {"Enter function", 0, 0, enter, ""},
	"end":      {"Return from function", 0, 0, end, ""},
	"store":    {"Store y into x", 2, 1, store, "2 ⤶ 'a ⤶ store ⤶ ⤒2; y⥗a"},
	"put":      {"Put y into x", 2, 0, put, "2 ⤶ 'a ⤶ store ⤶ y⥗a"},
	"exchange": {"Exchange y and the value in var x", 2, 1, exchange, "3 ⤶ 'a ⤶ exchange ⤶ ⤒a 3⥗a"},
	"recall":   {"Recall x", 1, 1, recall, "'a ⤶ recall ⤶ ⤒a"},
	"purge":    {"Purge x", 1, 0, purge, "'a ⤶ purge ⤶ undefined⥗a"},
	"drop":     {"Drop x", 1, 0, drop, "drop ⤶"},
	"swap":     {"Swap x and y", 2, 2, swap, "swap ⤶ ⤒x,y"},
	"clear":    {"Clear stack", 0, 0, clear, "clear ⤶"},
	"collect":  {"Collect stack into x", 1, 1, collect, ""},
	"expand":   {"Expand x into the stack", 1, -1, expand, ""},
	"print":    {"Print x", 1, 0, print, "'Hello world ⤶ print ⤶ Hello world⥱Console"},
	"render":   {"Render RAM as buffer", 0, 0, render, ""},
	"graph":    {"Graph a sequence", 3, 0, graph, ""},
	"status":   {"Display status", 0, 0, nil, ""},
	"files":    {"List availabel files in ROM", 0, 0, files, "files ⤶ [files]⥱Console"},
	"mmap":     {"Map a file to RAM", 1, 0, mmap, "'rom/file.raw ⤶ mmap ⤶ file.byes⥱RAM"},
	"repeat":   {"Execute x repeatedly", 4, 0, repeat, "0 ⤶ < ⤶'f ⤶ repeat ⤶"},
	"<=":       {"Set the result flag to 1 if y <= x", 2, 0, lessThan, ""},
	">=":       {"Set the result flag to 1 if y >= x", 2, 0, greaterThan, ""},
	"==":       {"Set the result flag to 1 if x = y", 2, 0, equals, ""},
	"!=":       {"Set the result flag to 1 if x != y", 2, 0, notEquals, ""},
	"ceval":    {"Conditionally evaluate x if result flag is 1", 1, 0, ceval, ""},
	"setloop":  {"Set loop counter to x", 1, 0, setLoop, ""},
	"dec":      {"Decrement the loop register", 0, 0, decrement, "dec"},
	"loop":     {"Execute x if the loop counter is not zero", 0, 0, loopNotZero, ""},
	"halt":     {"Halt execution", 0, 0, halt, "halt ⤶"},
}

var currentSequence = []CoreValue{}

func InitializeInstructionMap() {
	value := instructionMap["eval"]
	value.impl = eval
	instructionMap["eval"] = value
}

func OutputInstructionHelpDoc() {
	fmt.Printf("## Supported instructions\n\n")
	for k, v := range instructionMap {
		fmt.Printf("### %s\n- Description: %s\n- Arg count: %d\n- Result count: %d\n- Usage: %s\n\n", k, v.description, v.argCount, v.resultCount, v.usage)
	}
}

func halt(core *Core) InstructionResult {
	core.Mode = Halted
	return successResult
}

func repeat(core *Core) InstructionResult {
	return successResult
}

func consumeOne(core *Core) CoreValue {
	return *core.currentStack().Pop()
}

func consumeTwo(core *Core) (CoreValue, CoreValue) {
	return *core.currentStack().Pop(), *core.currentStack().Pop()
}
