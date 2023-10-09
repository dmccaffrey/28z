package core

import "fmt"

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
	"+": {"Add x and y", 2, 1, add, "6 ⤶ 2 ⤶ + ⤶ ⤒8"},
	"-": {"Subtract x from y", 2, 1, subtract, "6 ⤶ 2 ⤶ - ⤶ ⤒4"},
	"*": {"Multiply y by x", 2, 1, multiply, "6 ⤶ 2 ⤶ * ⤶ ⤒12"},
	"/": {"Divide y by x", 2, 1, divide, "6 ⤶ 2 ⤶ / ⤶ ⤒3"},

	"<<": {"Define function", 1, 0, define, ""},
	">>": {"Reduce function", 0, 0, reduce, ""},

	"enter": {"Enter function", 1, 0, enter, ""},
	"end":   {"Return from function", 1, 0, end, ""},

	"store":    {"Store y into x", 2, 1, store, "2 ⤶ 'a ⤶ store ⤶ ⤒2; y⥗a"},
	"put":      {"Put y into x", 2, 0, put, "2 ⤶ 'a ⤶ store ⤶ y⥗a"},
	"exchange": {"Exchange y and the value in var x", 2, 1, exchange, "3 ⤶ 'a ⤶ exchange ⤶ ⤒a 3⥗a"},
	"recall":   {"Recall x", 1, 1, recall, "'a ⤶ recall ⤶ ⤒a"},
	"purge":    {"Purge x", 1, 0, purge, "'a ⤶ purge ⤶ undefined⥗a"},
	"eval":     {"Evaluate x", 1, 0, nil, ""},

	"drop":    {"Drop x", 1, 0, drop, "drop ⤶"},
	"swap":    {"Swap x and y", 2, 2, swap, "swap ⤶ ⤒x,y"},
	"clear":   {"Clear stack", 0, 0, clear, "clear ⤶"},
	"collect": {"Collect stack into x", 1, 1, nil, ""},

	"print":  {"Print x", 1, 0, print, "'Hello world ⤶ print ⤶ Hello world⥱Console"},
	"graph":  {"Render graph", 0, 0, nil, ""},
	"status": {"Display status", 0, 0, nil, ""},

	"halt": {"Halt execution", 0, 0, halt, "halt ⤶"},
}

func OutputInstructionHelpDoc() {
	fmt.Printf("## Supported instructions\n\n")
	for k, v := range instructionMap {
		fmt.Printf("### %s\nDescription: %s\nArg count: %d\nReult count: %d\nUsage: %s\n\n", k, v.description, v.argCount, v.resultCount, v.usage)
	}
}

func add(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	if x.GetType() == FloatType {
		core.Push(FloatValue{y.GetFloat() + x.GetFloat()})
		return successResult

	}
	if x.GetType() == StringType {
		core.Push(StringValue{y.GetString() + x.GetString()})
		return successResult
	}
	return InstructionResult{true, "Unexpected operands"}
}

func subtract(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	if x.GetType() == FloatType {
		core.Push(FloatValue{y.GetFloat() - x.GetFloat()})
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
		core.Push(FloatValue{y.GetFloat() * x.GetFloat()})
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
		core.Push(FloatValue{y.GetFloat() / x.GetFloat()})
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
	core.Push(StringValue{"enter"})
	core.Mode = Storing
	return successResult
}

func reduce(core *Core) InstructionResult {
	result := ""
	if core.currentStack().length >= 2 && core.stackStack.length > 1 {
		steps := core.currentStack().ToArray()
		for _, v := range steps {
			result += v.GetString() + ";"
		}
		core.DropStack()
		core.Push(StringValue{result})
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

func store(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	core.VarMap[x.GetString()] = y
	core.Push(y)
	return successResult
}

func put(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	core.VarMap[x.GetString()] = y
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
	core.WriteLine(x.GetString())
	return successResult
}

func halt(core *Core) InstructionResult {
	core.Mode = Halted
	return successResult
}

func consumeOne(core *Core) CoreValue {
	return *core.currentStack().Pop()
}

func consumeTwo(core *Core) (CoreValue, CoreValue) {
	return *core.currentStack().Pop(), *core.currentStack().Pop()
}
