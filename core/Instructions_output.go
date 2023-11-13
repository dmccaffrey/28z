package core

import (
	"math"
	"slices"
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

	results := make([]float64, 92)
	step := (end.GetFloat() - start.GetFloat()) / 92
	for col := 0; col < 92; col++ {
		x := float64(col) * step
		core.Push(FloatValue{value: float64(x)})
		_eval(f.GetSequence(), core)
		result := consumeOne(core)
		results[col] = result.GetFloat()
	}
	min := slices.Min(results)
	max := slices.Max(results)
	if max == 0 {
		max = 1
	}
	if min < 0 && max < 0 {
		min = math.Abs(min)
		max = math.Abs(max)
	}
	for col := 0; col < 92; col++ {
		result := results[col]
		row := scale(result, min, max, 31)
		core.Ram[xyToOffset(col, row)] = 19

		row = scale(0, min, max, 31)
		core.Ram[xyToOffset(col, row)] = 8
	}
	core.Push(FloatValue{value: float64(min)})
	core.Push(FloatValue{value: float64(max)})
	render(core)
	return successResult
}

func scale(value float64, min float64, max float64, bound int) int {
	var scaled float64
	if min >= 0 {
		scaled = (value - min) / max
	} else {
		scaled = (value + math.Abs(min)) / (math.Abs(min) + max)
	}
	return bound - int(math.Round(scaled*float64(bound)))
}

func xyToOffset(x int, y int) int {
	return y*92 + x
}
