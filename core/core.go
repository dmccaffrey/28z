package core

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

var ModeNameMap = map[ExecutionMode]string{
	Running: "Running",
	Storing: "Storing",
	Halted:  "Halted",
}
var boolToI = map[bool]int{false: 0, true: 1}
var RegisterKeys = []string{Reg_Count, Reg_Depth, Reg_State, Reg_Flags, Reg_LoopC}

type (
	ExecutionMode    int
	RegisterFunction func(*Core) int
	Core             struct {
		stackStack         Stack[Stack[CoreValue]]
		prevStack          *Stack[CoreValue]
		Message            string
		Mode               ExecutionMode
		Console            []string
		LastInput          string
		Ram                []byte
		loopCounter        int
		resultFlag         bool
		breakFlag          bool
		Prompt             string
		interactiveHandler *InteractiveHandler
	}
	InteractiveHandler struct {
		Input  func(*Core) (bool, string)
		Output func(*Core)
	}
)

func NewCore() Core {
	core := Core{}
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
	c.LastInput = input
	if input == "" {
		Logger.Printf("Error: Empty input provided")
		return
	}
	if input == "run" || input == ">" {
		c.Mode = Running
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

	if c.Mode == Storing {
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

	if c.Mode == Storing {
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
	c.Message = error
	c.Mode = Halted
}

func (c *Core) unsetError() {
	c.Message = "OK"
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

func (c *Core) GetMode() string {
	return ModeNameMap[c.Mode]
}

func (c *Core) Store(key CoreValue, value CoreValue) {
	if key.GetType() == FloatType {
		if key.GetInt() >= len(c.Ram) {
			c.setError("RAM offset too large")
			return
		}
		c.Ram[key.GetInt()] = byte(value.GetInt())
		return
	}
	if key.GetType() == StringType {
		Variables[key.GetString()] = value
		return
	}
	c.setError("Invalid key type")
}

func (c *Core) ShouldBreak() bool {
	shouldBreak := c.breakFlag
	c.breakFlag = false
	return shouldBreak
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

func (c *Core) SetInteractiveHandler(handler *InteractiveHandler) {
	c.interactiveHandler = handler
}

func (c *Core) Mainloop() {
	c.interactiveHandler.Output(c)
	run := true
	for run {
		shouldContinue, input := c.interactiveHandler.Input(c)
		run = shouldContinue
		if run {
			c.ProcessRaw(input)
			c.interactiveHandler.Output(c)
		}
	}
}
