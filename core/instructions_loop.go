package core

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
