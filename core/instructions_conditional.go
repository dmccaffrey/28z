package core

func lessThan(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	if x.GetFloat() <= y.GetFloat() {
		core.Regs.State.ResultFlag = true
	} else {
		core.Regs.State.ResultFlag = false
	}
	return successResult
}

func greaterThan(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	if x.GetFloat() >= y.GetFloat() {
		core.Regs.State.ResultFlag = true
	} else {
		core.Regs.State.ResultFlag = false
	}
	return successResult
}

func equals(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	if x.GetType() == StringType {
		if x.GetString() == y.GetString() {
			core.Regs.State.ResultFlag = true
			return successResult
		}
		core.Regs.State.ResultFlag = false
		return successResult
	}
	if x.GetFloat() == y.GetFloat() {
		core.Regs.State.ResultFlag = true
		return successResult
	}
	core.Regs.State.ResultFlag = false
	return successResult
}

func notEquals(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	if x.GetFloat() != y.GetFloat() {
		core.Regs.State.ResultFlag = true
	} else {
		core.Regs.State.ResultFlag = false
	}
	return successResult
}

func unset(core *Core) InstructionResult {
	core.Regs.State.ResultFlag = false
	return successResult
}
