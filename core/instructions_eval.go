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
	return _eval(x.GetSequence(), core)
}

func _eval(sequence []CoreValue, core *Core) InstructionResult {
	end := len(sequence) - 1

	Logger.Printf("Evaluating sequence: value=%s\n", sequence)
	for i := end; i >= 0; i-- {
		val := sequence[i]
		switch val.GetType() {
		case InstructionType:
			Logger.Printf("Evaluating instruction: value=%s\n", val.GetString())
			if !val.(InstructionValue).CheckArgs(core) {
				return InstructionResult{true, "Too few arguments to instruction"}
			}
			result := val.(InstructionValue).Eval(core)
			if result != successResult {
				return result
			}
			break

		case ReferenceType:
			Logger.Printf("Dereferencing value: value=%s\n", val)
			_eval(val.(ReferenceValue).Dereference(core).GetSequence(), core)
			break

		default:
			Logger.Printf("Pushing value: value=%s\n", val)
			core.Push(val)
			break
		}
	}
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

func stream(core *Core) InstructionResult {
	x := consumeOne(core)
	core.NewStack()
	for i := 0; i < 2760; i++ {
		core.Push(FloatValue{value: float64(core.Ram[i])})
		core.Push(x)
		eval(core)
		core.Ram[i] = byte(consumeOne(core).GetFloat())
	}
	core.DropStack()
	return successResult
}

func reduce(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	core.NewStack()
	values := y.GetSequence()
	offset := 1
	lastResult := values[0]
	for ; offset < len(values); offset++ {
		core.Push(lastResult)
		core.Push(values[offset])
		core.Push(x)
		eval(core)
		lastResult = *core.Pop()
	}
	core.DropStack()
	core.Push(lastResult)
	return successResult
}
