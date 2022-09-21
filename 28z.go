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
	"math"
	"time"
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

var constsMap = map[string]float64 {
	"$pi": math.Pi,
	"$tau": math.Pi * 2,
	"$e": math.E,
	"$phi": math.Phi,
	"$maxf": math.MaxFloat64,
}


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

type StackData struct {
	dataType byte
	str string
	flt float64
}

func DefaultStackData() StackData {
	return StackData{Nil, "", 0}
}

func (d *StackData) Parse(s string) error {
	switch s[0] {
	case 'q':
		time.Sleep(3 * time.Second)
		return nil
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
			d.flt = float64(res)
		}
		return err
	case 'o':
		d.dataType = Oct
		s = s[1:]
		res, err := strconv.ParseInt(s, 8, 64)
		if err == nil {
			d.flt = float64(res)
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
	case Str:
		result.str = d.str + input.str
	default:
		result.flt = d.flt + input.flt
	}
	return result, nil
}

func (d StackData) Mult(input StackData) (StackData, error) {
	if d.dataType == Str || input.dataType == Str {
		return DefaultStackData(), errors.New("Multiplication not defined for strings")
	}
	d.flt *= input.flt
	return d, nil
}

func (d StackData) Div(input StackData) (StackData, error) {
	if d.dataType == Str || input.dataType == Str {
		return DefaultStackData(), errors.New("Division not defined for strings")
	}
	if input.flt == 0 {
		return DefaultStackData(), errors.New("Division by zero is not defined")
	}
	d.flt /= input.flt
	return d, nil
}

func (d StackData) Minus(input StackData) (StackData, error) {
	if d.dataType == Str || input.dataType == Str {
		return StackData{}, errors.New("Subtraction not defined for strings")
	}
	return d.Plus(input.ChS())
}

func (d StackData) Pow(input StackData) (StackData, error) {
	if d.dataType == Str || input.dataType == Str {
		return DefaultStackData(), errors.New("Power not defined for strings")
	}
	d.flt = math.Pow(d.flt, input.flt)
	return d, nil
}


func (d StackData) ChS() StackData {
	d.flt = - d.flt
	return d
}

func (d StackData) ToString() string {
	stackStr := "?"
	switch d.dataType {
	case Str:
		stackStr = fmt.Sprintf("%s", d.str)
	case Flt:
		stackStr = fmt.Sprintf("%.8E", d.flt)
	case Hex:
		stackStr = fmt.Sprintf("%014x", d.flt)
	case Oct:
		stackStr = fmt.Sprintf("%014o", d.flt)
	case Nil:
		stackStr = ""
	}
	if len(stackStr) > 14 {
		stackStr = stackStr[:11] + "..."
	}
	return stackStr
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
	content += fmt.Sprintf("%s\t\t%*s\n", "Registers", 17, "Stack")
	content += fmt.Sprintf("%s\t\t%s\n", strings.Repeat("-", 22), strings.Repeat("-", 22))
	end := MaxStackLen-1
	for i := end; i >= 0; i-- {
		stackEntry := DefaultStackData()
		if (i < len(s.stack)) {
			stackIndex := len(s.stack) - i - 1
			stackEntry = s.stack[stackIndex]
		}
		content += fmt.Sprintf("R%s: %*s (%c)\t\t%d: %*s (%c)\n", string(i + 65), 14, s.regs[i].ToString(), s.regs[i].dataType, i, 14, stackEntry.ToString(), stackEntry.dataType)
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
			s.Parse(input[1:])
		}
	case '@':
		if len(input) < 3 {
			s.err = "Invalid function"
			return true
		}
		if input[1] != 2 {
			s.applyUnaryFunc(uFuncs[input])
		} else if input[1] == '2' {
			s.applyBinaryFunc(bFuncs[input])
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
			reg := len(s.regs)
			if (x.dataType == Str) {
				reg = registerMap[strings.ToUpper(x.str)]
			}
			if (reg < len(s.regs)) {
				s.regs[reg] = y
			} else {
				return y, errors.New(fmt.Sprintf("Invalid register: reg=%d", reg))
			}
			return DefaultStackData(), nil
		})
	case 'r':
		s.applyUnaryFunc(func(x StackData) (StackData, error) {
			var reg int
			if (x.dataType == Str) {
				reg = registerMap[strings.ToUpper(x.str)]
			} else {
				reg = int(x.flt)
			}
			if (reg < len(s.regs)) {
				return s.regs[reg], nil
			}
			return DefaultStackData(), errors.New(fmt.Sprintf("Invalid register: reg=%d", reg))
		})
	case 'p':
		s.applyUnaryFunc(func(x StackData) (StackData, error) {
			s.console += x.ToString() + " "
			return x, nil
		})
	case 'l':
		for ; s.regs[1].flt>0; s.regs[1].flt-- {
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
				s.console += ">"
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

type windowSize struct {
	rows    uint16
	cols    uint16
}

type TermWriter struct {
	lastLineCount int
	output io.Writer
	winSize windowSize
	interactive bool
}

func InteractiveTermWriter() TermWriter {
	//out, size := AcquireTty()
	return TermWriter{0, io.Writer(os.Stdout), windowSize{0,0}, true}
}

func DebugTermWriter() TermWriter {
	return TermWriter{0, io.Writer(os.Stdout), windowSize{0,0}, false}
}

func (w *TermWriter) Publish(content string) {
	if w.interactive {
		var clear = fmt.Sprintf("%c[%dA%c[2K", 27, 1, 27)
		_, err := fmt.Fprint(w.output, strings.Repeat(clear, w.lastLineCount))
		if err != nil {
			fmt.Printf("Failed to reset output: err=%s", err.Error())
		}
	}

	// +1 for the newlilne caused by input
	w.lastLineCount = strings.Count(content, "\n") + 1
	fmt.Fprintf(w.output, "%s", content)
}

func AcquireTty() (io.Writer, windowSize) {
	out, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Failed to acquire TTY")
		return io.Writer(os.Stdout), windowSize{0,0}
	}
	writer := io.Writer(out)
	var winsize windowSize
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL,
		out.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&winsize)))
	fmt.Fprint(writer, "\033[H\033[2J")
	return writer, winsize
}

func main() {
	state := NewEnvState(DebugTermWriter())
	//state := NewEnvState(InteractiveTermWriter())
	state.Display("")
	in := bufio.NewReader(os.Stdin)
	for {
		input, err := in.ReadString('\n')
		if err == io.EOF {
			return
		}
		input = strings.TrimSuffix(input, "\n")
		if input == "exit" {
			break
		}
		if !state.Parse(input) {
			break
		}
		state.Display(input)
	}
	state.Display("exit")
}
