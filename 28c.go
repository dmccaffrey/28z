package main

import (
	"fmt"
	"strings"
	"io"
	"os"
	"strconv"
	"errors"
	"unsafe"
	"syscall"
	"bufio"
)

const(
	Flt byte = 'f'
	Hex      = 'h'
	Oct      = 'o'
	Bin      = 'b'
	Str      = 's'
	Nil      = '0'
)

const(
	MaxStackLen = 4
)

var registerMap = map[string]int {
	"A": 0,
	"RA": 0,
	"B": 1,
	"RB": 1,
	"C": 2,
	"RC": 2,
	"D": 3,
	"RD": 3,
}

type UnaryFunc func(x StackData) (StackData, error)
type BinaryFunc func(x StackData, y StackData) (StackData, error)

type StackData struct {
	dataType byte
	str string
	num int64
	flt float64
}

func DefaultStackData() StackData {
	return StackData{Nil, "", 0, 0}
}

func (d *StackData) Parse(s string) error {
	switch s[0] {
	case '"':
		s = strings.TrimSuffix(s, "\"")
		s = strings.TrimPrefix(s, "\"")
		d.dataType = Str
		d.str = s
	case 'x':
		d.dataType = Hex
		s = s[1:]
		res, err := strconv.ParseInt(s, 16, 64)
		if err == nil {
			d.num = res
		}
		return err
	case 'o':
		d.dataType = Oct
		s = s[1:]
		res, err := strconv.ParseInt(s, 8, 64)
		if err == nil {
			d.num = res
		}
		return err
	default:
		d.dataType = Flt
		res, err := strconv.ParseFloat(s, 64)
		if err == nil {
			d.flt = res
		}
		return err
	}
	return nil
}

func (d StackData) Plus(input StackData) (StackData, error) {
	result := d
	if d.dataType != input.dataType {
		return result, errors.New("Operand data types must match")
	}
	switch d.dataType {
	case Flt:
		result.flt = d.flt + input.flt
	case Hex:
		fallthrough
	case Oct:
		fallthrough
	case Bin:
		result.num = d.num + input.num
	default:
		result.str = d.str + input.str
	}
	return result, nil
}

func (d StackData) Mult(input StackData) (StackData, error) {
	if d.dataType == Str {
		return StackData{}, errors.New("Multiplication not defined for strings")
	}
	result := d
	if d.dataType != input.dataType {
		return result, errors.New("Operand data types must match")
	}
	switch d.dataType {
	case Flt:
		result.flt = d.flt * input.flt
	case Hex:
		fallthrough
	case Oct:
		fallthrough
	case Bin:
		result.num = d.num * input.num
	}
	return result, nil
}

func (d StackData) Div(input StackData) (StackData, error) {
	if d.dataType == Str {
		return StackData{}, errors.New("Division not defined for strings")
	}
	result := d
	if d.dataType != input.dataType {
		return result, errors.New("Operand data types must match")
	}
	switch d.dataType {
	case Flt:
		if input.flt == 0 {
			return input, errors.New("Division by zero is not defined")
		}
		result.flt = d.flt / input.flt
	case Hex:
		fallthrough
	case Oct:
		fallthrough
	case Bin:
		result.num = d.num / input.num
	}
	return result, nil
}

func (d StackData) Minus(input StackData) (StackData, error) {
	if d.dataType == Str {
		return StackData{}, errors.New("Subtraction not defined for strings")
	}
	return d.Plus(input.ChS())
}

func (d StackData) ChS() StackData {
	d.num = - d.num
	d.flt = - d.flt
	return d
}

func (d StackData) ToString() string {
	switch d.dataType {
	case Str:
		return fmt.Sprintf("%s", d.str)
	case Flt:
		return fmt.Sprintf("%.8E", d.flt)
	case Hex:
		return fmt.Sprintf("%014x", d.num)
	case Oct:
		return fmt.Sprintf("%014o", d.num)
	}
	return ""
}


type EnvState struct {
	err string
	overflow bool
	stack []StackData
	regs []int64
	console string
}

func (s EnvState) Display() string {
	content := ""
	content += fmt.Sprintf("%s %*s\n", "Registers", 17, "Stack")
	content += fmt.Sprintf("%s %*s\n", strings.Repeat("-", 12), 29, strings.Repeat("-", 20))
	end := MaxStackLen-1
	for i := end; i >= 0; i-- {
		suffix := byte('?')
		stackString := ""
		if (i < len(s.stack)) {
			stackIndex := len(s.stack) - i - 1
			stackString = s.stack[stackIndex].ToString()
			suffix = s.stack[stackIndex].dataType
			if len(stackString) > 14 {
				stackString = stackString[:11] + "..."
			}
		}
		content += fmt.Sprintf("R%s: %08x %*s %d: %*s (%c)\n", string(i + 65), s.regs[i], 8, "", i, 14, stackString, suffix)
	}
	if len(s.err) != 0 {
		content += fmt.Sprintf("\n<Err: %s>", s.err)
	}
	content += fmt.Sprintf("\n| %s", s.console)
	content += "\n\n: "
	return content
}

func (s *EnvState) Parse(input string) {
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
	case 'd':
		_ = s.Pop()
	case 's':
		s.applyBinaryFunc(func(x StackData, y StackData) (StackData, error) {
			var reg int
			if (x.dataType == Str) {
				reg = registerMap[x.str]
			} else {
				reg = int(x.num)
			}
			if (reg < len(s.regs)-1) {
				if y.dataType == Flt {
					y.num = int64(y.flt)
				}
				s.regs[reg] = y.num
			} else {
				return DefaultStackData(), errors.New(fmt.Sprintf("Invalid register: reg=%d", reg))
			}
			return DefaultStackData(), nil
		})
	case 'r':
		s.applyUnaryFunc(func(x StackData) (StackData, error) {
			var reg int
			if (x.dataType == Str) {
				reg = registerMap[x.str]
			} else {
				reg = int(x.num)
			}
			result := DefaultStackData()
			if (reg < len(s.regs)-1) {
				result.dataType = Hex
				result.num = s.regs[reg]
				return result, nil
			}
			return result, errors.New(fmt.Sprintf("Invalid register: reg=%d", reg))
		})
	case 'p':
		s.applyUnaryFunc(func(x StackData) (StackData, error) {
			s.console += x.ToString() + " "
			result := DefaultStackData()
			result.dataType = Nil
			return result, nil
		})
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
			return
		}
		s.Push(data)
	}
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

func NewEnvState() EnvState {
	return EnvState{ "", false, []StackData{}, make([]int64, MaxStackLen), ""}
}

type windowSize struct {
	rows    uint16
	cols    uint16
}

type TermWriter struct {
	lastLineCount int
	output io.Writer
	winSize windowSize
}

func (w *TermWriter) Publish(content string) {
	var clear = fmt.Sprintf("%c[%dA%c[2K", 27, 1, 27)
	_, _ = fmt.Fprint(w.output, strings.Repeat(clear, w.lastLineCount))

	// +1 for the newlilne caused by input
	w.lastLineCount = strings.Count(content, "\n") + 1
	bytes := []byte(content)
	w.output.Write(bytes)
}

func (w *TermWriter) AcquireTty() {
	out, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if err != nil {
		fmt.Fprintf(w.output, "Failed to acquire TTY")
		return
	}
	w.output = io.Writer(out)
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL,
		out.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&w.winSize)))
	fmt.Fprint(w.output, "\033[H\033[2J")
}

func main() {
	writer := TermWriter{0, io.Writer(os.Stdout), windowSize{0, 0}}
	writer.AcquireTty()
	state := NewEnvState()
	in := bufio.NewReader(os.Stdin)
	for {
		writer.Publish(state.Display())
		//fmt.Scanf("%s", &input)
		input, _ := in.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")
		if input == "exit" {
			break
		}
		state.Parse(input)
	}
}
