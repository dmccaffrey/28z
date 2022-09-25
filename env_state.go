package main

import (
	"fmt"
	"strings"
	"os"
	"errors"
	"math"
	"time"
)

type UnaryFunc func(x StackData) (StackData, error)
type BinaryFunc func(x StackData, y StackData) (StackData, error)

type EnvState struct {
	err string
	overflow bool
	stack []StackData
	regs [4]StackData
	console string
	graph [graphW][graphH]bool
	writer TermWriter
}

func (s EnvState) Display(instruction string) {
	content := fmt.Sprintf("\n  ╓%s╖\n", strings.Repeat("─", 92))
	content += fmt.Sprintf("  ║ \x1b[31m28z\033[0m ┇ Current Instruction = %-*s ║\n", 62, instruction)
	content += fmt.Sprintf("  ╟%s╥%s╢\n", strings.Repeat("─", 46), strings.Repeat("─", 45))
	//content += fmt.Sprintf("  ╟%s╢\n", strings.Repeat("─", 94))
	content += fmt.Sprintf("  ║ %-*s║ %-*s║\n", 45, "Registers", 44, "Stack")
	content += fmt.Sprintf("  ╟%s╫%s╢\n", strings.Repeat("┄", 46), strings.Repeat("┄", 45))
	end := MaxStackLen-1
	for i := end; i >= 0; i-- {
		stackEntry := DefaultStackData()
		if (i < len(s.stack)) {
			stackIndex := len(s.stack) - i - 1
			stackEntry = s.stack[stackIndex]
		}
		registerStr := fmt.Sprintf("R%s: (%c) %-*s", string(i + 65), s.regs[i].dataType, 20, s.regs[i].ToString())
		stackStr := fmt.Sprintf("%d:", i)
		if stackEntry.dataType != Nil {
			stackStr = fmt.Sprintf("%d: (%c) %-*s", i, stackEntry.dataType, 35, stackEntry.ToString())
		}
		content += fmt.Sprintf("  ║ %-*s║ %-*s║\n", 45, registerStr, 44, stackStr)
	}
	content += fmt.Sprintf("  ╟%s╨%s╢\n", strings.Repeat("─", 46), strings.Repeat("─", 45))
	if s.err != "" {
		content += fmt.Sprintf("  ║ 🯀 %-*s║\n", 89, s.err)

	} else {
		content += fmt.Sprintf("  ║ 🮱 %-*s║\n", 89, "OK")
	}
	content += fmt.Sprintf("  ╟%s╢\n", strings.Repeat("─", 92))
	lines := strings.Split(s.console, "\n")
	for _,v := range lines {
		content += fmt.Sprintf("  ║%-*s║\n", 92, v)
	}
	content += fmt.Sprintf("  ╟%s╜\n", strings.Repeat("─", 92))
	content += "  ║\n  ╙─🮥 <instruction> 🮴: "
	s.writer.Publish(content)
}

func (s *EnvState) Parse(input string) bool {
	if len(input) == 0 {
		return true
	}
	s.err = ""
	if input == "exit" {
		os.Exit(0)
	}
	if input == "break" {
		return false
	}
	switch input[0] {
	case '+':
		s.applyBinaryFunc(func(x StackData, y StackData) (StackData, error) {
			return y.Plus(x)
		})
	case '-':
		s.applyBinaryFunc(func(x StackData, y StackData) (StackData, error) {
			return y.Minus(x)
		})
	case '/':
		s.applyBinaryFunc(func(x StackData, y StackData) (StackData, error) {
			return y.Div(x)
		})
	case '*':
		s.applyBinaryFunc(func(x StackData, y StackData) (StackData, error) {
			return y.Mult(x)
		})
	case '^':
		s.applyBinaryFunc(func(x StackData, y StackData) (StackData, error) {
			return y.Pow(x)
		})
	case '`':
		if s.regs[registerMap["RC"]].dataType == Flt && s.regs[registerMap["RC"]].flt == 1 {
			return s.Parse(input[1:])
		}
	case '@':
		if len(input) < 3 {
			s.err = "Invalid function"
			return true
		}
		if input[1] != 2 {
			fun, ok := uFuncs[input]
			if ok {
				s.applyUnaryFunc(fun)
			} else {
				s.err = "Unknown function"
			}
		} else if input[1] == '2' {
			fun, ok := bFuncs[input]
			if ok {
				s.applyBinaryFunc(fun)
			} else {
				s.err = "Unknown function"
			}
		}
	case '$':
		if len(input) < 2 {
			s.err = "Invalid constant"
			return true
		}
		s.Push(StackData{Flt, "", constsMap[input]})
	case 'd':
		s.applyUnaryFunc(func(x StackData) (StackData, error) {
			return DefaultStackData(), nil
		})
	case 'i':
		s.applyUnaryFunc(func(x StackData) (StackData, error){
			x.flt = float64(int(x.flt))
			return x, nil
		})
	case 'e':
		s.applyUnaryFunc(func(x StackData) (StackData, error) {
			var lines []string
			if strings.Contains(x.str, "_") {
				lines = strings.Split(x.str, "_")
			} else {
				lines = strings.Split(x.str, "|")
			}
			for i := range(lines) {
				s.Parse(lines[i])
			}
			return DefaultStackData(), nil
		})
	case 's':
		s.applyBinaryFunc(func(x StackData, y StackData) (StackData, error) {
			if (x.dataType == Str) {
				reg, ok := registerMap[strings.ToUpper(x.str)]
				if ok {
					s.regs[reg] = y

				} else {
					return y, errors.New(fmt.Sprintf("Invalid register: reg=%d", reg))
				}
			}
			return DefaultStackData(), nil
		})
	case 'r':
		s.applyUnaryFunc(func(x StackData) (StackData, error) {
			if (x.dataType == Str) {
				reg, ok := registerMap[strings.ToUpper(x.str)]
				if ok {
					return s.regs[reg], nil
				}
				prog, ok := progsMap[strings.ToUpper(x.str)]
				if ok {
					return StackData{Str, prog, 0.0}, nil
				}

			} else {
				reg := int(x.flt)
				if reg < len(s.regs) {
					return s.regs[reg], nil
				}
			}
			return DefaultStackData(), errors.New(fmt.Sprintf("Invalid register or program: input=%d", x.str))
		})
	case 'p':
		if strings.Contains(input, "c") {
			s.console = ""
		}
		s.applyUnaryFunc(func(x StackData) (StackData, error) {
			output := strings.Replace(x.ToString(), `\n`, "\n", -1)
			s.console += output + " "
			return x, nil
		})
	case 'g':
		if strings.Contains(input, "p") {
			s.console = ""
			for r:=0; r<graphH; r++ {
				for c:=0; c<graphW; c++ {
					if s.graph[c][r] {
						s.console += "█"
					} else {
						s.console += "░"
					}
				}
				s.console += "\n"
			}
		}
		if strings.Contains(input, "c") {
			s.graph = [graphW][graphH]bool{}
		}
		s.applyUnaryFunc(func(x StackData) (StackData, error){
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
			xPt := int(s.regs[registerMap["RB"]].flt)
			if xPt < graphW {
				s.graph[xPt][yPt] = true
			}
			return x, nil
		})
	case 'l':
		for ; s.regs[1].flt>0; s.regs[1].flt -= 1.0 {
			lines := strings.Split(s.regs[0].str, "|")
			for i := range(lines) {
				if s.err != "" {
					return true
				}
				if len(lines[i]) == 0 {
					continue
				}
				res := s.Parse(lines[i])
				s.Display(lines[i])
				if !res {
					return true
				}
			}
		}
	case '?':
		if len(input) < 2 {
			s.err = "Invalid comparison"
			return true
		}
		if len(s.stack) < 2 {
			s.err = "Comparison requires two parameters"
			return true
		}
		x := s.stack[len(s.stack)-1]
		y := s.stack[len(s.stack)-2]
		res := 0.0
		switch(input[1]) {
		case '=':
			if x.flt == y.flt {
				res = 1.0
			}
		case '<':
			if x.flt < y.flt {
				res = 1.0
			}
		case '>':
			if x.flt > y.flt {
				res = 1.0
			}
		}
		s.regs[registerMap["RC"]].dataType = Flt
		s.regs[registerMap["RC"]].flt = res
	case ';':
		if len(s.stack) >= 2 {
			x := s.Pop()
			y := s.Pop()
			s.Push(x)
			s.Push(y)
		}
	case ':':
		if len(s.stack) == 4 {
			x := s.Pop()
			y := s.Pop()
			z := s.Pop()
			t := s.Pop()
			s.Push(y)
			s.Push(x)
			s.Push(t)
			s.Push(z)
		}
	case '~':
		s.applyUnaryFunc(func(x StackData) (StackData, error){
			time.Sleep(time.Duration(x.flt) * time.Millisecond)
			return DefaultStackData(), nil
		})
	default:
		data := StackData{}
		err := data.Parse(input)
		if err != nil {
			s.err = err.Error()
			return true
		}
		s.Push(data)
	}
	return true
}

func (s *EnvState) applyBinaryFunc(f BinaryFunc) {
	if len(s.stack) < 2 {
		s.err = "Operation requires two operands"
		return
	}
	x := s.Pop()
	y := s.Pop()
	r, err := f(x, y)
	if err == nil {
		if r.dataType != Nil {
			s.Push(r)
		}
		return
	}
	s.Push(y)
	s.Push(x)
	s.err = err.Error()
}

func (s *EnvState) applyUnaryFunc(f UnaryFunc) {
	if len(s.stack) < 1 {
		s.err = "Operation requires one operand"
		return
	}
	x := s.Pop()
	r, err := f(x)
	if err == nil {
		if r.dataType != Nil {
			s.Push(r)
		}
		return
	}
	s.Push(x)
	s.err = err.Error()
}

func (s *EnvState) PushRaw(input string) {
	input = strings.Replace(input, ",", "", -1)
	data := StackData{}
	err := data.Parse(input)
	if err != nil {
		s.err = err.Error()
		return
	}
	s.Push(data)
}

func (s *EnvState) Push(input StackData) {
	s.stack = append(s.stack, input)
	stackLen := len(s.stack)
	if stackLen > MaxStackLen {
		s.stack = s.stack[1 :]
	}
}

func (s *EnvState) Pop() StackData {
	val := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return val
}

func NewEnvState(writer TermWriter) EnvState {
	state := EnvState{ writer: writer}
	for i := range(state.regs) {
		state.regs[i] = DefaultStackData()
	}
	return state
}
