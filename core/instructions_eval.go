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
		core.EvalSequence(y.GetSequence())
	} else {
		core.EvalSequence(x.GetSequence())
	}
	return successResult
}

func eval(core *Core) InstructionResult {
	x := consumeOne(core)
	core.EvalSequence(x.GetSequence())
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
		core.EvalSequence(x.GetSequence())
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
