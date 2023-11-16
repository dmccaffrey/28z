package core

import (
	"math"
)

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
	xVal := Variables[x.GetString()]
	if xVal == nil {
		return InstructionResult{true, "Variable not set"}
	}
	Variables[x.GetString()] = y
	core.Push(xVal)
	return successResult
}

func recall(core *Core) InstructionResult {
	x := consumeOne(core)
	if x.GetType() == FloatType && int(x.GetFloat()) < len(core.Ram) {
		core.Push(FloatValue{value: float64(core.Ram[int(x.GetFloat())])})
		return successResult
	}
	val := Variables[x.GetString()]
	if val == nil {
		return InstructionResult{true, "Variable not set"}
	}
	core.Push(val)
	return successResult
}

func purge(core *Core) InstructionResult {
	x := consumeOne(core)
	Variables[x.GetString()] = nil
	return successResult
}

func mmap(core *Core) InstructionResult {
	x := consumeOne(core)
	bytes := RawData[x.GetString()]
	if bytes == nil {
		return InstructionResult{true, "File not found"}
	}
	len := int(math.Min(float64(len(bytes)), float64(len(core.Ram))))
	for i := 0; i < len; i++ {
		value := bytes[i]
		core.Ram[i] = value
	}
	return successResult
}

func files(core *Core) InstructionResult {
	for k := range Programs {
		core.WriteLine("Program: " + k)
	}
	for k := range RawData {
		core.WriteLine("Data: " + k)
	}
	return successResult
}

func zero(core *Core) InstructionResult {
	for i := range core.Ram {
		core.Ram[i] = 0
	}
	return successResult
}
