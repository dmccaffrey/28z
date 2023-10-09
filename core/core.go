package core

import (
	"strconv"
	"strings"
)

const (
	Running ExecutionMode = 0
	Storing               = 1
	Halted                = 2
)

const (
	Reg_LoopC string = "LOOPC"
	Reg_Flags        = "FLAGS"
	Reg_State        = "STATE"
	Reg_Depth        = "DEPTH"
	Reg_Count        = "COUNT"
)

var RegisterKeys = []string{Reg_Count, Reg_Depth, Reg_State, Reg_Flags, Reg_LoopC}

type (
	ExecutionMode    int
	RegisterFunction func(*Core) int
	Core             struct {
		VarMap     map[string]CoreValue
		stackStack Stack[Stack[CoreValue]]
		Message    string
		Mode       ExecutionMode
		Console    []string
		LastInput  string
	}
)

func NewCore() Core {
	core := Core{}
	core.VarMap = map[string]CoreValue{}
	core.stackStack.Push(NewCoreValueStack())
	core.Mode = Running
	core.Console = make([]string, 0)
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
	if c.Mode == Storing {
		c.Push(StringValue{input})
		return
	}
	if len(input) > 1 && input[0] == '\'' {
		input = strings.TrimSuffix(input, "'")
		input = strings.TrimPrefix(input, "'")
		if len(input) != 0 {
			c.Push(StringValue{input})
			return
		}
	}
	result, err := strconv.ParseFloat(input, 64)
	if err == nil {
		c.Push(FloatValue{result})
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
	if impl.argCount > c.currentStack().Len() {
		c.setError("Too few arguments")
		return
	}
	result := impl.impl(c)
	if result.error {
		c.setError(result.message)
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
		Reg_LoopC: zeroFunc(c),
		Reg_Flags: zeroFunc(c),
		Reg_State: coreState(c),
		Reg_Depth: StackDepth(c),
		Reg_Count: StackCount(c),
	}
}

func (c *Core) WriteLine(line string) {
	c.Console = append(c.Console, line)
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
	return int(core.Mode)
}
