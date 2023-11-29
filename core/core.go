package core

const (
	Running ExecutionMode = 1
	Storing               = 2
	Halted                = 3
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
	ExecutionMode    int8
	RegisterFunction func(*Core) int
	Core             struct {
		stackStack         Stack[Stack[CoreValue]]
		prevStack          *Stack[CoreValue]
		Ram                []byte
		Regs               Registers
		Error              CoreValue
		interactiveHandler InteractiveHandler
	}
	InteractiveHandler interface {
		GetInput(*Core) (bool, string)
		Display(*Core)
		Output(string)
		Prompt(*Core, string) (bool, string)
		Clear()
	}
	Registers struct {
		State       StateRegister
		Mode        ExecutionMode
		LoopCounter int16
	}
	StateRegister struct {
		ResultFlag bool
		BreakFlag  bool
		PromptFlag bool
	}
)

func NewCore() Core {
	core := Core{}
	core.stackStack.Push(NewCoreValueStack())
	core.Regs.Mode = Running
	core.Ram = make([]byte, 8192)
	return core
}

func (c *Core) currentStack() *Stack[CoreValue] {
	return c.stackStack.Peek()
}

func (c *Core) NewStack() {
	c.prevStack = c.stackStack.Peek()
	c.stackStack.Push(NewCoreValueStack())
}

func (c *Core) DropStack() {
	c.stackStack.Pop()
	if c.stackStack.top == nil {
		c.NewStack()
	}

	if c.stackStack.top.prev != nil {
		c.prevStack = &c.stackStack.top.prev.value

	} else {
		c.prevStack = nil
	}
}

func (c *Core) ProcessRaw(input string) {
	Logger.Printf("Processing raw input: input=%s\n", input)
	c.unsetError()
	if input == "" {
		Logger.Printf("Error: Empty input provided")
		return
	}
	if input == "run" || input == ">" {
		c.Regs.Mode = Running
	}

	runReference := false
	if input[0] == '%' {
		runReference = true
	}

	value := RawToInstruction(input)
	if value.GetType() == InstructionType {
		c.ProcessInstruction(value.(InstructionValue))
		return
	}

	value = RawToImmediateCoreValue(input)
	if value.GetType() == DefaultType {
		Logger.Printf("Error: Not a valid input: input=%s\n", input)
		c.setError("Not a valid input")
		return
	}

	if c.Regs.Mode == Storing {
		Logger.Printf("Storing value: input=%s, value=%s\n", input, value)
		c.Push(value)
		return
	}

	if value.GetType() == ReferenceType {
		value = value.(ReferenceValue).Dereference(c)
	}

	if runReference {
		result := _eval(value.GetSequence(), c)
		if result.error {
			Logger.Printf("Error: Failed to eval sequence: err=%s, seq=%s", result.message, value)
			return
		}
		return
	}

	c.Push(value)
}

func (c *Core) ProcessInstruction(instruction InstructionValue) {
	Logger.Printf("Processing instruction: value=%s\n", instruction.value.description)
	impl := instruction.value
	if !impl.IsValid() {
		Logger.Printf("Error: Instruction is not valid: value=%s", instruction.value.description)
		c.setError("Not a valid instruction")
		return
	}

	if c.Regs.Mode == Storing {
		c.Push(instruction)
		Logger.Printf("Pushed instruction: value=%s", instruction.value.description)
		return
	}

	if impl.argCount > c.currentStack().Len() {
		Logger.Printf("Error: Too few arguments for instruction: value=%s, have=%d, required=%d",
			instruction.value.description, c.currentStack().Len(), instruction.value.argCount)
		c.setError("Too few arguments")
		return
	}

	Logger.Printf("Executing instruction: value=%s", instruction.value.description)
	result := impl.impl(c)
	if result.error {
		Logger.Printf("Error: Error evaluating instruction: value=%s, err=%s", instruction.value.description, result.message)
		c.setError(result.message)
	}
}

func (c *Core) setError(error string) {
	c.Error = StringValue{value: error}
	c.Regs.Mode = Halted
}

func (c *Core) unsetError() {
	c.Error = DefaultValue{}
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
	c.DropStack()
}

func (c *Core) GetStackArray() []CoreValue {
	return c.currentStack().ToArray()
}

func (c *Core) Store(key CoreValue, value CoreValue) {
	if key.GetType() == FloatType {
		index := key.GetInt()
		if index < 0 || index >= len(c.Ram) {
			Logger.Printf("Error: store out of range for RAM: index=%d", index)
			c.setError("RAM offset too large")
			return
		}
		c.Ram[index] = byte(value.GetInt())
		return
	}
	if key.GetType() == StringType {
		Variables[key.GetString()] = value
		return
	}
	Logger.Printf("Invalid key type: %s", key)
	c.setError("Invalid key type")
}

func (c *Core) ShouldBreak() bool {
	shouldBreak := c.Regs.State.BreakFlag
	c.Regs.State.BreakFlag = false
	return shouldBreak
}

func (c *Core) StackCount() int {
	return c.currentStack().length
}

func (c *Core) StackDepth() int {
	return c.stackStack.length
}

func (c *Core) SetInteractiveHandler(handler InteractiveHandler) {
	c.interactiveHandler = handler
}

func (c *Core) Mainloop() {
	c.interactiveHandler.Display(c)
	run := true
	for run {
		shouldContinue, input := c.interactiveHandler.GetInput(c)
		run = shouldContinue
		if run {
			c.ProcessRaw(input)
			c.interactiveHandler.Display(c)
		}
	}
}
