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
	var sequence []CoreValue
	if x.GetType() != SequenceType {
		sequence = x.GetSequence()
	} else {
		sequence = x.(SequenceValue).value
	}
	for core.loopCounter != 0 {
		_eval(sequence, core)
		decrement(core)
		if core.ShouldBreak() {
			break
		}
	}
	return successResult
}
