package core

func ceval(core *Core) InstructionResult {
	if core.Regs.State.ResultFlag {
		eval(core)
	} else {
		_ = consumeOne(core)
	}
	return successResult
}

func ceval2(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	if core.Regs.State.ResultFlag {
		_eval(y.GetSequence(), core)
	} else {
		_eval(x.GetSequence(), core)
	}
	return successResult
}

func eval(core *Core) InstructionResult {
	x := consumeOne(core)
	return _eval(x.GetSequence(), core)
}

func _eval(sequence []CoreValue, core *Core) InstructionResult {
	end := len(sequence) - 1

	Logger.Printf("Evaluating sequence: len=%d, value=%s\n", len(sequence), sequence)
	for i := end; i >= 0; i-- {
		val := sequence[i]
		switch val.GetType() {
		case InstructionType:
			Logger.Printf("[%d] Evaluating instruction: value=%s\n", i, val.GetString())
			if !val.(InstructionValue).CheckArgs(core) {
				return InstructionResult{true, "Too few arguments to instruction"}
			}
			result := val.(InstructionValue).Eval(core)
			if result != successResult {
				return result
			}
			break

		case ReferenceType:
			Logger.Printf("[%d] Dereferencing value: value=%s\n", i, val)
			core.Push(val.(ReferenceValue).Dereference(core))
			break

		default:
			Logger.Printf("[%d] Pushing value: value=%s\n", i, val)
			core.Push(val)
			break
		}
		if core.ShouldBreak() {
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
		result := core.Pop()
		if result != nil {
			results[i] = *result
		} else {
			results[i] = DefaultValue{}
			Logger.Printf("Error: no result during apply")
		}
	}
	core.DropStack()
	core.Push(SequenceValue{value: results})
	return successResult
}

func each(core *Core) InstructionResult {
	x, y := consumeTwo(core)
	core.NewStack()
	for i, value := range y.GetSequence() {
		core.Push(FloatValue{value: float64(i)})
		core.Push(value)
		_eval(x.GetSequence(), core)
		if core.ShouldBreak() {
			break
		}
	}
	core.DropStack()
	return successResult
}

func stream(core *Core) InstructionResult {
	x := consumeOne(core)
	core.NewStack()
	for i := 0; i < 2760; i++ {
		core.Push(FloatValue{value: float64(core.Ram[i])})
		core.Push(x)
		eval(core)
		if core.ShouldBreak() {
			break
		}
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
