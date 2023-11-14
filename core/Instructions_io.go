package core

import (
	"math"
	"slices"
	"strings"
	"time"
)

func print(core *Core) InstructionResult {
	x := consumeOne(core)
	if x.GetType() == ReferenceType {
		x = x.(ReferenceValue).Dereference(core)
	}
	core.WriteLine(x.GetString())
	core.interactiveHandler.Output(core)
	return successResult
}

func clearBuffer(core *Core) InstructionResult {
	core.ClearConsole()
	return successResult
}

func render(core *Core) InstructionResult {
	core.ClearConsole()
	for r := 0; r < 30; r++ {
		var sb strings.Builder
		for c := 0; c < 92; c++ {
			value := int(core.Ram[92*r+c]) % (len(Symbols))
			sb.WriteRune(Symbols[value])
		}
		core.WriteLine(sb.String())
	}
	return successResult
}

func show(core *Core) InstructionResult {
	render(core)
	core.interactiveHandler.Output(core)

	return successResult
}

func sleep(core *Core) InstructionResult {
	x := consumeOne(core)
	time.Sleep(time.Duration(x.GetFloat()) * time.Millisecond)
	return successResult
}

func graph(core *Core) InstructionResult {
	f := consumeOne(core)
	end, start := consumeTwo(core)

	y := -1
	step := (end.GetFloat() - start.GetFloat()) / 92
	if start.GetFloat() < 0 && end.GetFloat() >= 0 {
		span := end.GetFloat() + math.Abs(start.GetFloat())
		step = span / 92
		y = int(math.Round(math.Abs(start.GetFloat()) / step))

	} else if start.GetFloat() < 0 && end.GetFloat() < 0 {
		span := math.Abs(start.GetFloat()) + math.Abs(end.GetFloat())
		step = span / 92
	}

	if y != -1.0 {
		for row := 0; row < 32; row++ {
			if core.Ram[xyToOffset(y, row)] == 0 {
				core.Ram[xyToOffset(y, row)] = 2
			}
		}
	}

	results := make([]float64, 92)

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
		offset := xyToOffset(col, row)
		if core.Ram[offset] == 0 {
			core.Ram[offset] = 8
		}
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

func prompt(core *Core) InstructionResult {
	x := consumeOne(core)
	core.Prompt = x.GetString()
	core.interactiveHandler.Output(core)
	valid, input := core.interactiveHandler.Input(core)
	if valid {
		value := RawToImmediateValue(input, core)
		core.Push(value)

	} else {
		core.Push(DefaultValue{})
	}
	return successResult
}