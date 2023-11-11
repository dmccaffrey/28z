package core

import "math"

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

func mmap(core *Core) InstructionResult {
	x := consumeOne(core)
	bytes := RawData[x.GetString()]
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
	for k := range Programs {
		core.WriteLine("Program: " + k)
	}
	for k := range RawData {
		core.WriteLine("Data: " + k)
	}
	return successResult
}
