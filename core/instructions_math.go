package core

import "math"

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

func inverse(core *Core) InstructionResult {
	x := consumeOne(core)
	val := x.GetFloat()
	if val == 0 {
		return InstructionResult{true, "Divide by zero"}
	}
	core.Push(FloatValue{value: 1 / val})
	return successResult
}

func sin(core *Core) InstructionResult {
	x := consumeOne(core)
	result := math.Sin(x.GetFloat())
	core.Push(FloatValue{value: result})
	return successResult
}

func cos(core *Core) InstructionResult {
	x := consumeOne(core)
	result := math.Cos(x.GetFloat())
	core.Push(FloatValue{value: result})
	return successResult
}
