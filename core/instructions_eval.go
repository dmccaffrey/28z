package core

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
			if val.GetType() == ReferenceType {
				val = val.(ReferenceValue).Dereference(core)
			}
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

func apply(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	results := make([]CoreValue, len(y.GetSequence()))
	core.NewStack()
	for i, value := range y.GetSequence() {
		core.Push(value)
		core.Push(x)
		eval(core)
		results[i] = *core.Pop()
	}
	core.DropStack()
	core.Push(SequenceValue{value: results})
	return successResult
}
