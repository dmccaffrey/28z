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
	"v+":       {"Addition for all permuations of x and y", 2, 1, vplus, ""},
	"-":        {"Subtract x from y", 2, 1, subtract, "6 ⤶ 2 ⤶ - ⤶ ⤒4"},
	"*":        {"Multiply y by x", 2, 1, multiply, "6 ⤶ 2 ⤶ * ⤶ ⤒12"},
	"/":        {"Divide y by x", 2, 1, divide, "6 ⤶ 2 ⤶ / ⤶ ⤒3"},
	"mod":      {"y modulus by x", 2, 1, modulus, "6 ⤶ 2 ⤶ / ⤶ ⤒0"},
	"inverse":  {"Inverts x", 1, 1, inverse, ""},
	"sin":      {"sin of x", 1, 1, sin, ""},
	"cos":      {"cos of x", 1, 1, cos, ""},
	"rand":     {"Generate random between 0 and 1", 0, 1, random, ""},
	"<":        {"Define sequence", 0, 0, defineSequence, "< ⤶"},
	">":        {"Define sequence", 0, 0, reduceSequence, "> ⤶"},
	"this":     {"Refer to the current sequence", 0, 1, this, "this ⤶"},
	"eval":     {"Evaluate x", 1, 0, nil, ""},
	"consume":  {"Pop from previous stack and push to current", 0, 1, consume, "consume ⤶"},
	"produce":  {"Pop from this stack and push to previous", 1, 0, produce, "produce ⤶"},
	"apply":    {"Evalue x against all entries in y to modify y", 2, 1, apply, "apply ⤶"},
	"each":     {"Evaluate x against all entries in y", 2, 0, each, ""},
	"reduce":   {"Use x to reduce y to a single value", 2, 1, reduce, "reduce ⤶"},
	"enter":    {"Enter function, creating a new stack", 0, 0, enter, "enter ⤶"},
	"end":      {"Return from function, dropping the stack", 0, 0, end, "end ⤶"},
	"store":    {"Store y into x, preserving y", 2, 1, store, "2 ⤶ 'a ⤶ put ⤶ ⤒2; y⥗a"},
	"move":     {"Store y into x", 2, 0, move, "2 ⤶ 'a ⤶ asref ⤶ y⥗a"},
	"exchange": {"Exchange y and the value in var x", 2, 1, exchange, "3 ⤶ 'a ⤶ exchange ⤶ ⤒a 3⥗a"},
	"get":      {"Dereference x, preserving x", 1, 1, get, "'a ⤶ get ⤶ ⤒a"},
	"deref":    {"Derefernce x, replacing x", 1, 1, deref, "'a ⤶ deref ⤶ ⤒a"},
	"purge":    {"Deallocate that reference x", 1, 0, purge, "'a ⤶ purge ⤶ undefined⥗a"},
	"drop":     {"Drop x", 1, 0, drop, "drop ⤶"},
	"swap":     {"Swap x and y", 2, 2, swap, "swap ⤶ ⤒x,y"},
	"clear":    {"Clear stack", 0, 0, clear, "clear ⤶"},
	"collect":  {"Collect stack into x", 1, 1, collect, "1 ⤶ 2 ⤶ collect ⤶ ⤒[2]:1,2"},
	"pair":     {"Collect x and y into x", 2, 1, pair, "1 ⤶ 2 ⤶ collect ⤶ ⤒[2]:1,2"},
	"expand":   {"Expand x into the stack", 1, -1, expand, "⤒[2]:1,2 | expand ⤶ ⤒1 ⤒2"},
	"dup":      {"Duplicates x on the stack", 1, 2, duplicate, ""},
	"print":    {"Print x", 1, 0, print, "'Hello world ⤶ print ⤶ Hello world⥱Console"},
	"clearbuf": {"Clear the output buffer", 0, 0, clearBuffer, ""},
	"render":   {"Render RAM as buffer", 0, 0, render, "render ⤶"},
	"show":     {"Render and pause", 0, 0, show, ""},
	"graph":    {"Graph a sequence", 3, 0, graph, "graph ⤶"},
	"prompt":   {"Prompt the user for a value", 1, 1, prompt, "'Enter x ⤶ prompt ⤶"},
	"status":   {"Display status", 0, 0, nil, ""},
	"files":    {"List availabel files in ROM", 0, 0, files, "files ⤶ [files]⥱Console"},
	"mmap":     {"Map a file to RAM", 1, 0, mmap, "'rom/file.raw ⤶ mmap ⤶ file.byes⥱RAM"},
	"stream":   {"Apply x to renderable RAM", 1, 0, stream, ""},
	"zero":     {"Zero RAM", 0, 0, zero, ""},
	"repeat":   {"Execute x repeatedly", 1, 0, repeat, "0 ⤶ < ⤶'f ⤶ repeat ⤶"},
	"<=":       {"Set the result flag to 1 if y <= x", 2, 0, lessThan, ""},
	">=":       {"Set the result flag to 1 if y >= x", 2, 0, greaterThan, ""},
	"==":       {"Set the result flag to 1 if x = y", 2, 0, equals, ""},
	"!=":       {"Set the result flag to 1 if x != y", 2, 0, notEquals, ""},
	"unset":    {"Sets the result flat to 0", 0, 0, unset, ""},
	"ceval":    {"Conditionally evaluate x if result flag is 1", 1, 0, ceval, "⤒<sequence> | ceval ⤶"},
	"ceval2":   {"Conditionally evaluate y if result flag is 1, otherwise evaluate x", 2, 0, ceval2, "⤒<sequence>, ⤒<sequence> | ceval2 ⤶"},
	"generate": {"Evaluate a pair where y is the last input and x is the generator", 1, 1, generate, "⤒<pair> ⤶ generate ⤶ ⤒<pair>, ⤒<result>"},
	"setloop":  {"Set loop counter to x", 1, 0, setLoop, "5 ⤶ setloop ⤶"},
	"dec":      {"Decrement the loop register", 0, 0, decrement, "dec"},
	"loop":     {"Execute x if the loop counter is not zero", 0, 0, loopNotZero, "5 ⤶ setloop ⤶ ⤒<sequence> | loop ⤶"},
	"halt":     {"Halt execution", 0, 0, halt, "halt ⤶"},
	"sleep":    {"Sleep for x ms", 1, 0, sleep, ""},
	"inspect":  {"Write a raw object to file", 1, 0, inspect, ""},
	"stop":     {"Stop the current loop", 0, 0, stop, ""},
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
	core.Regs.Mode = Halted
	return successResult
}

func stop(core *Core) InstructionResult {
	core.Regs.State.BreakFlag = true
	return successResult
}

func consumeOne(core *Core) CoreValue {
	if *&core.currentStack().length == 0 {
		return DefaultValue{}
	}
	return *core.currentStack().Pop()
}

func consumeTwo(core *Core) (CoreValue, CoreValue) {
	return *core.currentStack().Pop(), *core.currentStack().Pop()
}
