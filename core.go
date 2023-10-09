package main

import (
	"errors"
	"strings"
	"math"
	"fmt"
)

type Instruction string

type Parser func (instruction string) (bool)

func (d StackData) Eval(p Parser) error {
	if d.dataType != Str {
		return errors.New("Invalid type for eval")
	}
	separator := " "
	if strings.Contains(d.str, "_") {
		separator = "_"
	}
	for _,v := range strings.Split(d.str, separator) {
		if !p(v) {
			errors.New("Evaluation stopped due to errors")
		}
	}
	return nil
}

func (d *StackData) Loop(f StackData, p Parser) error {
	for ; d.flt>0; d.flt -= 1.0 {
		err := f.Eval(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func GraphPoint(x StackData, rb StackData, graph *Graph) (StackData, error) {
	if x.flt > 1.0 || x.flt < -1.0 {
		return x, errors.New("Graph value must be between -1 and 1")
	}
	scaled := (graphH / 2) * x.flt
	scaled += graphH / 2
	yPt := int(math.Round(scaled))
	if yPt > graphH-1 {
		yPt = graphH-1
	} else if yPt < 0 {
		yPt = 0
	}
	xPt := int(rb.flt)
	if xPt < graphW {
		graph[xPt][yPt] = true
	}
	return x, nil
}

func RenderGraph(console *string, graph Graph) {
	*console = ""
	for r:=0; r<graphH; r++ {
		for c:=0; c<graphW; c++ {
			if graph[c][r] {
				*console += "█"
			} else {
				*console += "░"
			}
		}
		*console += "\n"
	}
}

func Store(x StackData, y StackData, regs *Registers) (StackData, error) {
	if (x.dataType == Str) {
		reg, ok := registerMap[strings.ToUpper(x.str)]
		if ok {
			regs[reg] = y

		} else {
			return y, errors.New(fmt.Sprintf("Invalid register: reg=%d", reg))
		}
	}
	return DefaultStackData(), nil
}

func Recall(x StackData, regs *Registers) (StackData, error) {
	if (x.dataType == Str) {
		reg, ok := registerMap[strings.ToUpper(x.str)]
		if ok {
			return regs[reg], nil
		}
		prog, ok := progsMap[strings.ToUpper(x.str)]
		if ok {
			return StackData{Str, prog, 0.0}, nil
		}

	} else {
		reg := int(x.flt)
		if reg < len(regs) {
			return regs[reg], nil
		}
	}
	return DefaultStackData(), errors.New(fmt.Sprintf("Invalid register or program: input=%d", x.str))
}
