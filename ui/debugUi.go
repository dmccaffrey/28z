package ui

import (
	"dmccaffrey/28z/core"
	"fmt"
	"math"
	"strings"
	"time"
)

var regWidth = 13
var stackWidth = 41
var msgWidth = 30
var scrWidth = 92
var scrHeight = 30

var uiS0 = fmt.Sprintf(" ╓%s╥%s╥%s╖\n", strings.Repeat("─", regWidth+2), strings.Repeat("─", stackWidth+2), strings.Repeat("┄", msgWidth+2))
var uiS1 = fmt.Sprintf(" ╟%s╥%s╥%s╢\n", strings.Repeat("─", regWidth+2), strings.Repeat("─", stackWidth+2), strings.Repeat("┄", msgWidth+2))
var uiS3 = fmt.Sprintf(" ╟%s╨%s╨%s╢\n", strings.Repeat("─", regWidth+2), strings.Repeat("─", stackWidth+2), strings.Repeat("┄", msgWidth+2))
var uiS4 = fmt.Sprintf(" ╟%s╢\n", strings.Repeat("─", 92))
var uiIn = fmt.Sprintf(" ║  %*s ║\n", 89, "")
var uiS5 = fmt.Sprintf(" ╙─%s╜\n", strings.Repeat("─", 91))

var stackAliases = []string{"(x)", "(y)", "(z)", "   ", "   "}

var lastUiUpdate = time.Now().Local()

func (z Interactive28z) DisplayDebugUi(vm *core.Core) string {
	var sb strings.Builder

	if z.prompt != "" {
		sb.WriteString("\033[2m")
	}

	sb.WriteString(uiS0)
	stack := vm.GetStackArray()
	for i := 4; i >= 0; i-- {
		stackValue := ""
		if i < len(stack) {
			stackValue = stack[i].GetString()
		}
		stackStr := fmt.Sprintf("%s %1d: %-*s", stackAliases[i], i, 58, stackValue)
		msgStr := ""
		regStr := ""
		switch i {
		case 0:
			msgStr = fmt.Sprintf("%-5s %s", "Last:", z.lastInput)
			regStr = fmt.Sprintf("%-6s %s", "STATE:", StateToString(vm.Regs.State))
			break
		case 1:
			msgStr = fmt.Sprintf("%-5s %s", "Err:", z.message)
			regStr = fmt.Sprintf("%-6s %d", "LOOPC:", vm.Regs.Mode)
			break
		case 2:
			regStr = fmt.Sprintf("%-6s %d", "MODE:", vm.Regs.Mode)
			break
		case 3:
			regStr = fmt.Sprintf("%-6s %d", "DEPTH:", vm.StackDepth())
			break
		case 4:
			regStr = fmt.Sprintf("%-6s %d", "COUNT:", vm.StackCount())
			break
		}
		sb.WriteString(fmt.Sprintf(" ║ %-*s ║ %-*.40s ║ %-*.30s ║\n", regWidth, regStr, stackWidth, stackStr, msgWidth, msgStr))
	}
	sb.WriteString(uiS3)
	end := int(math.Min(float64(len(z.console)), float64(scrHeight)))
	for i := 0; i < end; i++ {
		sb.WriteString(fmt.Sprintf(" ║%-*.92s║\n", scrWidth, z.console[i]))
	}
	for i := scrHeight - end; i > 0; i-- {
		sb.WriteString(fmt.Sprintf(" ║%-*s║\n", scrWidth, ""))

	}
	sb.WriteString(uiS4)
	sb.WriteString(uiIn)
	sb.WriteString(uiS5)

	promptLine := " > "
	if z.prompt != "" {
		sb.WriteString("\033[0m")
		promptLine = fmt.Sprintf("\0331 | Requested input: %s > \0330", z.prompt)
	}

	sb.WriteString(fmt.Sprintf("  \033[2A \x1b[31m28z\033[0m %s", promptLine))
	return sb.String()
}

var b2i = map[bool]int8{false: 0, true: 1}

func StateToString(r core.StateRegister) string {
	return fmt.Sprintf("%d%d%d", b2i[r.ResultFlag], b2i[r.BreakFlag], b2i[r.PromptFlag])
}
