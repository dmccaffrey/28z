package main

import (
	"fmt"
	"strings"
	"os"
	"errors"
	"math"
)

type UnaryFunc func(x StackData) (StackData, error)
type BinaryFunc func(x StackData, y StackData) (StackData, error)

var uFuncs = map[string]UnaryFunc {
	"@sin": func (x StackData) (StackData, error) {
		x.flt = math.Sin(x.flt)
		return x, nil
	},
	"@cos": func (x StackData) (StackData, error) {
		x.flt = math.Cos(x.flt)
		return x, nil
	},
	"@tan": func (x StackData) (StackData, error) {
		x.flt = math.Tan(x.flt)
		return x, nil
	},
	"@log": func (x StackData) (StackData, error) {
		x.flt = math.Log10(x.flt)
		return x, nil
	},
	"@ln": func (x StackData) (StackData, error) {
		x.flt = math.Log(x.flt)
		return x, nil
	},
	"@logb": func (x StackData) (StackData, error) {
		x.flt = math.Logb(x.flt)
		return x, nil
	},
}

var bFuncs = map[string]BinaryFunc {
}

type EnvState struct {
	err string
	overflow bool
	stack []StackData
	regs []StackData
	console string
	writer TermWriter
}

func (s EnvState) Display(instruction string) {
	content := fmt.Sprintf("Current Instruction=%s\n\n", instruction)
	content += fmt.Sprintf("%s\t\t%*s\n", "Registers", 21, "Stack")
	content += fmt.Sprintf("%s\t\t%s\n", strings.Repeat("-", 28), strings.Repeat("-", 28))
	end := MaxStackLen-1
	for i := end; i >= 0; i-- {
		stackEntry := DefaultStackData()
		if (i < len(s.stack)) {
			stackIndex := len(s.stack) - i - 1
			stackEntry = s.stack[stackIndex]
		}
		content += fmt.Sprintf("R%s: %*s (%c)\t\t%d: %*s (%c)\n", string(i + 65), 20, s.regs[i].ToString(), s.regs[i].dataType, i, 20, stackEntry.ToString(), stackEntry.dataType)
	}
	if len(s.err) != 0 {
		content += fmt.Sprintf("\n<Err: %s>", s.err)
	}
	content += fmt.Sprintf("\n| %s", s.console)
	content += "\n\n: "
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
			lines := strings.Split(x.str, "|")
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
		s.applyUnaryFunc(func(x StackData) (StackData, error) {
			s.console += x.ToString() + " "
			return x, nil
		})
	case 'l':
		for ; s.regs[1].flt>0; s.regs[1].flt -= 1.0 {
			lines := strings.Split(s.regs[0].str, "|")
			for i := range(lines) {
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
	state := EnvState{ "", false, []StackData{}, make([]StackData, MaxStackLen), "", writer}
	for i := range(state.regs) {
		state.regs[i] = DefaultStackData()
	}
	return state
}