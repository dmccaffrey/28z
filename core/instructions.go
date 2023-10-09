package core

import (
	"fmt"
	"math"
	"strings"
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
	"+":   {"Add x and y", 2, 1, add, "6 ⤶ 2 ⤶ + ⤶ ⤒8"},
	"-":   {"Subtract x from y", 2, 1, subtract, "6 ⤶ 2 ⤶ - ⤶ ⤒4"},
	"*":   {"Multiply y by x", 2, 1, multiply, "6 ⤶ 2 ⤶ * ⤶ ⤒12"},
	"/":   {"Divide y by x", 2, 1, divide, "6 ⤶ 2 ⤶ / ⤶ ⤒3"},
	"mod": {"y modulus by x", 2, 1, modulus, "6 ⤶ 2 ⤶ / ⤶ ⤒0"},

	"{":    {"Define function", 1, 0, define, ""},
	"}":    {"Reduce function", 0, 0, reduce, ""},
	"<":    {"Define sequence", 0, 0, defineSequence, ""},
	">":    {"Define sequence", 0, 0, reduceSequence, ""},
	"this": {"Refer to the current sequence", 0, 1, this, ""},
	"eval": {"Evaluate x", 1, 0, nil, ""},

	"enter": {"Enter function", 1, 0, enter, ""},
	"end":   {"Return from function", 1, 0, end, ""},

	"store":    {"Store y into x", 2, 1, store, "2 ⤶ 'a ⤶ store ⤶ ⤒2; y⥗a"},
	"put":      {"Put y into x", 2, 0, put, "2 ⤶ 'a ⤶ store ⤶ y⥗a"},
	"exchange": {"Exchange y and the value in var x", 2, 1, exchange, "3 ⤶ 'a ⤶ exchange ⤶ ⤒a 3⥗a"},
	"recall":   {"Recall x", 1, 1, recall, "'a ⤶ recall ⤶ ⤒a"},
	"purge":    {"Purge x", 1, 0, purge, "'a ⤶ purge ⤶ undefined⥗a"},

	"drop":    {"Drop x", 1, 0, drop, "drop ⤶"},
	"swap":    {"Swap x and y", 2, 2, swap, "swap ⤶ ⤒x,y"},
	"clear":   {"Clear stack", 0, 0, clear, "clear ⤶"},
	"collect": {"Collect stack into x", 1, 1, nil, ""},

	"print":  {"Print x", 1, 0, print, "'Hello world ⤶ print ⤶ Hello world⥱Console"},
	"render": {"Render RAM as buffer", 0, 0, render, ""},
	"graph":  {"Graph a sequence", 1, 0, nil, ""},
	"status": {"Display status", 0, 0, nil, ""},

	"files": {"List availabel files in ROM", 0, 0, files, "files ⤶ [files]⥱Console"},
	"mmap":  {"Map a file to RAM", 1, 0, mmap, "'rom/file.raw ⤶ mmap ⤶ file.byes⥱RAM"},

	"repeat": {"Execute x repeatedly", 4, 0, repeat, "0 ⤶ < ⤶'f ⤶ repeat ⤶"},
	"<=":     {"Set the result flag to 1 if y <= x", 2, 0, lessThan, ""},
	">=":     {"Set the result flag to 1 if y >= x", 2, 0, greaterThan, ""},
	"==":     {"Set the result flag to 1 if x = y", 2, 0, equals, ""},
	"!=":     {"Set the result flag to 1 if x != y", 2, 0, notEquals, ""},
	"ceval":  {"Conditionally evaluate x if result flag is 1", 1, 0, ceval, ""},

	"setloop": {"Set loop counter to x", 1, 0, setLoop, ""},
	"dec":     {"Decrement the loop register", 0, 0, decrement, "dec"},
	"loop":    {"Execute x if the loop counter is not zero", 0, 0, loopNotZero, ""},

	"halt": {"Halt execution", 0, 0, halt, "halt ⤶"},
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

func add(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	if x.GetType() == FloatType {
		core.Push(FloatValue{value: y.GetFloat() + x.GetFloat()})
		return successResult

	}
	if x.GetType() == StringType {
		core.Push(StringValue{value: y.GetString() + x.GetString()})
		return successResult
	}
	return InstructionResult{true, "Unexpected operands"}
}

func subtract(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	if x.GetType() == FloatType {
		core.Push(FloatValue{value: y.GetFloat() - x.GetFloat()})
		return successResult
	}
	if x.GetType() == StringType {
		core.Push(y)
		return successResult
	}
	return InstructionResult{true, "Unexpected operands"}
}

func multiply(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	if x.GetType() == FloatType {
		core.Push(FloatValue{value: y.GetFloat() * x.GetFloat()})
		return successResult
	}
	if x.GetType() == StringType {
		core.Push(y)
		return successResult
	}
	return InstructionResult{true, "Unexpected operands"}
}

func divide(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	if x.GetType() == FloatType {
		core.Push(FloatValue{value: y.GetFloat() / x.GetFloat()})
		return successResult
	}
	if x.GetType() == StringType {
		core.Push(y)
		return successResult
	}
	return InstructionResult{true, "Unexpected operands"}
}

func modulus(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	if x.GetType() == FloatType {
		core.Push(FloatValue{value: float64(y.GetInt() % x.GetInt())})
		return successResult
	}
	if x.GetType() == StringType {
		core.Push(y)
		return successResult
	}
	return InstructionResult{true, "Unexpected operands"}
}

func define(core *Core) InstructionResult {
	argCount := consumeOne(core)
	core.NewStack()
	core.Push(argCount)
	core.Push(StringValue{value: "enter"})
	core.Mode = Storing
	return successResult
}

func reduce(core *Core) InstructionResult {
	if core.currentStack().length >= 2 && core.stackStack.length > 1 {
		steps := core.currentStack().ToArray()
		value := SequenceValue{value: steps}
		core.DropStack()
		core.Push(value)
	}
	return successResult
}

func enter(core *Core) InstructionResult {
	argCount := consumeOne(core).GetInt()
	if argCount > StackCount(core) {
		return InstructionResult{true, "Not enough aguments"}
	}
	args := make([]CoreValue, argCount)
	for i := argCount - 1; i >= 0; i-- {
		args[i] = consumeOne(core)
	}
	core.NewStack()
	for _, v := range args {
		core.Push(v)
	}
	return successResult
}

func end(core *Core) InstructionResult {
	resultCount := consumeOne(core).GetInt()
	args := make([]CoreValue, resultCount)
	for i := resultCount - 1; i >= 0; i-- {
		args[i] = consumeOne(core)
	}
	core.DropStack()
	for _, v := range args {
		core.Push(v)
	}
	return successResult
}

func defineSequence(core *Core) InstructionResult {
	core.NewStack()
	core.Mode = Storing
	return successResult
}

func reduceSequence(core *Core) InstructionResult {
	steps := core.currentStack().ToArray()
	value := SequenceValue{value: steps}
	core.DropStack()
	core.Mode = Running
	core.Push(value)
	return successResult
}

func this(core *Core) InstructionResult {
	core.Push(SequenceValue{value: currentSequence})
	return successResult
}

func ceval(core *Core) InstructionResult {
	if core.GetResultFlag() {
		eval(core)
	} else {
		_ = consumeOne(core)
	}
	return successResult
}

func eval(core *Core) InstructionResult {
	x := consumeOne(core)
	lastSequence := currentSequence
	currentSequence = x.GetSequence()
	end := len(x.GetSequence()) - 1
	for i := range x.GetSequence() {
		val := x.GetSequence()[end-i]
		if val.GetType() != InstructionType {
			core.Push(val)

		} else {
			if !val.(InstructionValue).CheckArgs(core) {
				return InstructionResult{true, "Too few arguments to instruction"}
			}
			result := val.(InstructionValue).Eval(core)
			if result != successResult {
				return result
			}
		}
	}
	currentSequence = lastSequence
	return successResult
}

func store(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	core.Store(x, y)
	core.Push(y)
	return successResult
}

func put(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	core.Store(x, y)
	return successResult
}

func exchange(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	xVal := core.VarMap[x.GetString()]
	if xVal == nil {
		return InstructionResult{true, "Variable not set"}
	}
	core.VarMap[x.GetString()] = y
	core.Push(xVal)
	return successResult
}

func recall(core *Core) InstructionResult {
	x := consumeOne(core)
	val := core.VarMap[x.GetString()]
	if val == nil {
		return InstructionResult{true, "Variable not set"}
	}
	core.Push(val)
	return successResult
}

func purge(core *Core) InstructionResult {
	x := consumeOne(core)
	core.VarMap[x.GetString()] = nil
	return successResult
}

func drop(core *Core) InstructionResult {
	consumeOne(core)
	return successResult
}

func swap(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	core.Push(x)
	core.Push(y)
	return successResult
}

func clear(core *Core) InstructionResult {
	core.ClearStack()
	return successResult
}

func print(core *Core) InstructionResult {
	x := consumeOne(core)
	if x.GetType() == ReferenceType {
		x = x.(ReferenceValue).Dereference(core)
	}
	core.WriteLine(x.GetString())
	return successResult
}

func render(core *Core) InstructionResult {
	core.ClearConsole()
	for r := 0; r < 30; r++ {
		var sb strings.Builder
		for c := 0; c < 92; c++ {
			value := core.Ram[92*r+c]
			if value > 31 && value < 127 {
				sb.WriteByte(value)

			} else {
				sb.WriteByte('-')
			}
		}
		core.WriteLine(sb.String())
	}
	return successResult
}

func halt(core *Core) InstructionResult {
	core.Mode = Halted
	return successResult
}

func mmap(core *Core) InstructionResult {
	x := consumeOne(core)
	bytes := Rom[x.GetString()]
	if bytes == nil {
		return InstructionResult{true, "File not found"}
	}
	len := int(math.Min(float64(len(bytes)), float64(len(core.Ram))))
	for i := 0; i < len; i++ {
		core.Ram[i] = bytes[i]
	}
	return successResult
}

func files(core *Core) InstructionResult {
	for k := range Rom {
		core.WriteLine(k)
	}
	return successResult
}

func repeat(core *Core) InstructionResult {
	return successResult
}

func lessThan(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	if x.GetFloat() <= y.GetFloat() {
		core.SetResultFlag(true)
	} else {
		core.SetResultFlag(false)
	}
	return successResult
}

func greaterThan(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	if x.GetFloat() >= y.GetFloat() {
		core.SetResultFlag(true)
	} else {
		core.SetResultFlag(false)
	}
	return successResult
}

func equals(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	if x.GetFloat() == y.GetFloat() {
		core.SetResultFlag(true)
	} else {
		core.SetResultFlag(false)
	}
	return successResult
}

func notEquals(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	if x.GetFloat() != y.GetFloat() {
		core.SetResultFlag(true)
	} else {
		core.SetResultFlag(false)
	}
	return successResult
}

func setLoop(core *Core) InstructionResult {
	x := consumeOne(core)
	core.loopCounter = x.GetInt()
	return successResult
}

func decrement(core *Core) InstructionResult {
	core.loopCounter -= 1
	return successResult
}

func loopNotZero(core *Core) InstructionResult {
	x := consumeOne(core)
	for core.loopCounter != 0 {
		core.Push(x)
		eval(core)
	}
	return successResult
}

func consumeOne(core *Core) CoreValue {
	return *core.currentStack().Pop()
}

func consumeTwo(core *Core) (CoreValue, CoreValue) {
	return *core.currentStack().Pop(), *core.currentStack().Pop()
}
