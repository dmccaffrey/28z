package core

func enter(core *Core) InstructionResult {
	core.NewStack()
	return successResult
}

func end(core *Core) InstructionResult {
	core.DropStack()
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

func consume(core *Core) InstructionResult {
	if core.prevStack != nil {
		x := core.prevStack.Pop()
		if x != nil {
			core.Push(*x)
			return successResult
		}
		return InstructionResult{true, "Previous stack is empty"}
	}
	return InstructionResult{true, "No previous stack"}
}

func produce(core *Core) InstructionResult {
	x := consumeOne(core)
	if core.prevStack != nil {
		core.prevStack.Push(x)
	}
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

func collect(core *Core) InstructionResult {
	value := SequenceValue{value: core.GetStackArray()}
	core.ClearStack()
	core.Push(value)
	return successResult
}

func expand(core *Core) InstructionResult {
	x := consumeOne(core)
	values := x.GetSequence()
	for i := range values {
		core.Push(values[len(values)-i-1])
	}
	return successResult
}
