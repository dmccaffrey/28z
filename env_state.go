package main

import (
	"strings"
	"os"
	"time"
)

const(
	ramSize = 4096
)

type Ram [ramSize]byte
type Registers [4]StackData

type UnaryFunc func(x StackData) (StackData, error)
type BinaryFunc func(x StackData, y StackData) (StackData, error)

type EnvState struct {
	err string
	overflow bool
	stack []StackData
	regs Registers
	console string
	ram Ram
	prompt string
	writer TermWriter
}

func (s *EnvState) Parse(input string, userInput bool) bool {
	if len(input) == 0 {
		return true
	}
	s.err = ""
	s.prompt = ""
	defer Display(*s, input, userInput)
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
	case '|':
		s.applyBinaryFunc(func(x StackData, y StackData) (StackData, error) {
			return y.Or(x)
		})
	case '&':
		s.applyBinaryFunc(func(x StackData, y StackData) (StackData, error) {
			return y.And(x)
		})
	case '`':
		if s.regs[registerMap["RC"]].dataType == Flt && s.regs[registerMap["RC"]].flt == 1 {
			return s.Parse(input[1:], userInput)
		}
	case '@':
		ufun, ok := uFuncs[input]
		if ok {
			s.applyUnaryFunc(ufun)
			break
		}
		bfun, ok := bFuncs[input]
		if ok {
			s.applyBinaryFunc(bfun)
			break
		}
		s.err = "Unknown function"
	case '$':
		if len(input) < 2 {
			s.err = "Invalid constant"
			return true
		}
		s.Push(StackData{Flt, "", constsMap[input]})
	case 'c':
		if strings.Contains(input, "m") {
			s.ram = Ram{}
		}
		if strings.Contains(input, "c") {
			s.console = ""
		}
		if strings.Contains(input, "r") {
			s.regs[RA] = DefaultStackData()
			s.regs[RB] = DefaultStackData()
			s.regs[RC] = DefaultStackData()
			s.regs[RD] = DefaultStackData()
		}
		if strings.Contains(input, "s") {
			s.stack = []StackData{}
		}
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
			x.Eval(s.Parse)
			return DefaultStackData(), nil
		})
	case 's':
		s.applyBinaryFunc(func(x StackData, y StackData) (StackData, error) {
			return Store(x, y, &s.regs, &s.ram)
		})
	case 'r':
		s.applyUnaryFunc(func(x StackData) (StackData, error) {
			return Recall(x, &s.regs, &s.ram)
		})
	case 'p':
		s.applyUnaryFunc(func(x StackData) (StackData, error) {
			output := strings.Replace(x.ToString(false), `\n`, "\n", -1)
			s.console += output + " "
			return x, nil
		})
	case 'g':
		RenderGraph(&s.console, s.ram)
	case 'l':
		s.regs[1].Loop(s.regs[0], s.Parse)
	case 'q':
		s.applyUnaryFunc(func(x StackData) (StackData, error) {
			s.prompt = x.str
			Display(*s, input, true)
			response, err := getInput()
			for err != nil {
				s.err = err.Error()
				Display(*s, response, true)
				response, err = getInput()
			}
			if response == "debug" {
				debugLoop(s)
			}
			s.Parse(response, true)
			return DefaultStackData(), nil
		})
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
			if x.dataType == Str {
				if x.str == y.str {
					res = 1.0
				}
			} else if x.dataType != Nil {
				if x.flt == y.flt {
					res = 1.0
				}
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
