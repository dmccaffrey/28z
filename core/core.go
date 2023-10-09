package core

import (
	"strconv"
	"strings"
)

const (
	Running ExecutionMode = 1000
	Storing               = 2000
	Halted                = 3000
)

const (
	Reg_LoopC string = "LOOPC"
	Reg_Flags        = "FLAGS"
	Reg_State        = "STATE"
	Reg_Depth        = "DEPTH"
	Reg_Count        = "COUNT"
)

var boolToI = map[bool]int{false: 0, true: 1}
var RegisterKeys = []string{Reg_Count, Reg_Depth, Reg_State, Reg_Flags, Reg_LoopC}

type (
	ExecutionMode    int
	RegisterFunction func(*Core) int
	Core             struct {
		VarMap      map[string]CoreValue
		stackStack  Stack[Stack[CoreValue]]
		Message     string
		Mode        ExecutionMode
		Console     []string
		LastInput   string
		Ram         []byte
		loopCounter int
		resultFlag  bool
	}
)

func NewCore() Core {
	core := Core{}
	core.VarMap = map[string]CoreValue{}
	core.stackStack.Push(NewCoreValueStack())
	core.Mode = Running
	core.Console = make([]string, 0)
	core.Ram = make([]byte, 8192)
	return core
}

func (c *Core) currentStack() *Stack[CoreValue] {
	return c.stackStack.Peek()
}

func (c *Core) NewStack() {
	c.stackStack.Push(NewCoreValueStack())
}

func (c *Core) DropStack() {
	c.stackStack.Pop()
}

func (c *Core) ProcessRaw(input string) {
	c.LastInput = input
	if input == "run" || input == ">" || input == "}" || input == ")" {
		c.Mode = Running
	}

	if len(input) > 1 {
		if input[0] == '\'' {
			input = strings.TrimPrefix(input, "'")
			c.Push(StringValue{value: input})
			return
		}
		if input[0] == '$' {
			input = strings.TrimPrefix(input, "$")
			reg, ok := c.GetRegisterMap()[input]
			if ok {
				c.Push(FloatValue{value: float64(reg)})
				return
			}
			variable, ok := c.VarMap[input]
			if ok {
				c.Push(variable)
				return
			}
		}
	}
	result, err := strconv.ParseFloat(input, 64)
	if err == nil {
		c.Push(FloatValue{value: result})
		return
	}
	c.ProcessInstruction(input)
}

func (c *Core) ProcessInstruction(instruction string) {
	impl := instructionMap[instruction]
	if !impl.IsValid() {
		c.setError("Not a valid instruciton")
		return
	}

	if c.Mode == Storing {
		c.Push(InstructionValue{value: impl})
		return
	}

	if impl.argCount > c.currentStack().Len() {
		c.setError("Too few arguments")
		return
	}

	result := impl.impl(c)
	if result.error {
		c.setError(result.message)
	}
}

func (c *Core) ProcessCoreValue(value CoreValue) {
	if value.GetType() == StringType {
		impl := instructionMap[value.GetString()]
		if impl.IsValid() {
			c.ProcessInstruction(value.GetString())
		}

	} else {
		c.Push(value)
	}
}

func (c *Core) setError(error string) {
	c.Message = error
	c.Mode = Halted
}

func (c *Core) GetStackString() string {
	result := "Stack: "
	stack := c.currentStack().ToArray()
	for i := 0; i < len(stack); i++ {
		if stack[i] == nil {
			continue
		}
		str := stack[i].GetString()
		result += str + ", "
	}

	return result
}

func (c *Core) Push(value CoreValue) {
	c.currentStack().Push(value)
}

func (c *Core) Pop() *CoreValue {
	value := c.currentStack().Pop()
	return value
}

func (c *Core) ClearStack() {
	for c.stackStack.length != 0 {
		c.DropStack()
	}
	c.NewStack()
}

func (c *Core) GetStackArray() []CoreValue {
	return c.currentStack().ToArray()
}

func (c *Core) GetRegisterMap() map[string]int {
	return map[string]int{
		Reg_LoopC: getLoopCounter(c),
		Reg_Flags: zeroFunc(c),
		Reg_State: coreState(c),
		Reg_Depth: StackDepth(c),
		Reg_Count: StackCount(c),
	}
}

func (c *Core) WriteLine(line string) {
	c.Console = append(c.Console, line)
}

func (c *Core) ClearConsole() {
	c.Console = make([]string, 0)
}

func zeroFunc(core *Core) int {
	return 0
}

func StackCount(core *Core) int {
	return core.currentStack().length
}

func StackDepth(core *Core) int {
	return core.stackStack.length
}

func coreState(core *Core) int {
	return int(core.Mode) + boolToI[core.resultFlag]
}

func getLoopCounter(core *Core) int {
	return core.loopCounter
}

func (c *Core) GetResultFlag() bool {
	return c.resultFlag
}

func (c *Core) SetResultFlag(result bool) {
	c.resultFlag = result
}
