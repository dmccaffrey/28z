package core

import (
	"strings"
)

func print(core *Core) InstructionResult {
	x := consumeOne(core)
	if x.GetType() == ReferenceType {
		x = x.(ReferenceValue).Dereference(core)
	}
	core.WriteLine(x.GetString())
	return successResult
}

func render(core *Core) InstructionResult {
	core.ClearConsole()
	for r := 0; r < 30; r++ {
		var sb strings.Builder
		for c := 0; c < 92; c++ {
			value := int(core.Ram[92*r+c]) % (len(Symbols) - 1)
			sb.WriteRune(Symbols[value])
		}
		core.WriteLine(sb.String())
	}
	return successResult
}

func graph(core *Core) InstructionResult {
	f := consumeOne(core)
	end, start := consumeTwo(core)

	step := 92 / (end.GetInt() - start.GetInt())
	for r := 0; r < 92; r++ {
		core.Push(FloatValue{value: float64(r * step)})
		core.Push(f)
		eval(core)
		core.Ram[r+92*(consumeOne(core).GetInt()/32)] = 'x'

	}
	render(core)
	return successResult
}
