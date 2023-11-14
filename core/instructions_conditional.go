package core

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

func unset(core *Core) InstructionResult {
	core.SetResultFlag(false)
	return successResult
}
