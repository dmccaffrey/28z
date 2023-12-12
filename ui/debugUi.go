package ui

import (
	"bytes"
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

var uiS0 = fmt.Sprintf(" ╓%s╥%s╥%s╖\n", strings.Repeat("─", regWidth+2), strings.Repeat("─", stackWidth+2), strings.Repeat("─", msgWidth+2))
var uiS1 = fmt.Sprintf(" ╟%s╥%s╥%s╢\n", strings.Repeat("─", regWidth+2), strings.Repeat("─", stackWidth+2), strings.Repeat("─", msgWidth+2))
var uiS3 = fmt.Sprintf(" ╟%s╨%s╨%s╢\n", strings.Repeat("─", regWidth+2), strings.Repeat("─", stackWidth+2), strings.Repeat("─", msgWidth+2))
var uiS4 = fmt.Sprintf(" ╟%s╢\n", strings.Repeat("─", 92))
var uiIn = fmt.Sprintf(" ║  %*s ║\n", 89, "")
var uiS5 = fmt.Sprintf(" ╙─%s╜\n", strings.Repeat("─", 91))
var startDim = "\033[2m"
var endDim = "\033[0m"

var stackAliases = []string{"(x)", "(y)", "(z)", "   ", "   "}

var lastUiUpdate = time.Now().Local()

func (z *Interactive28z) GenerateDebugUi() []byte {
	var bb bytes.Buffer

	bb.WriteString("\033[H\033[2J")

	if z.prompt != "" {
		bb.WriteString(startDim)
	}

	bb.WriteString(uiS0)
	stack := z.core.GetStackArray()
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
			msgStr = fmt.Sprintf("%-5s %s", "LAST:", z.lastInput)
			regStr = fmt.Sprintf("%-6s %s", "STATE:", StateToString(z.core.Regs.State))
			break
		case 1:
			msgStr = fmt.Sprintf("%-5s %s", "ERR:", z.message)
			regStr = fmt.Sprintf("%-6s %03d", "LOOPC:", z.core.Regs.Mode)
			break
		case 2:
			msgStr = fmt.Sprintf("%-5s %08d", "TICK:", z.core.Ticks)
			regStr = fmt.Sprintf("%-6s %03d", "MODE:", z.core.Regs.Mode)
			break
		case 3:
			regStr = fmt.Sprintf("%-6s %03d", "DEPTH:", z.core.StackDepth())
			break
		case 4:
			regStr = fmt.Sprintf("%-6s %03d", "COUNT:", z.core.StackCount())
			break
		}
		bb.WriteString(fmt.Sprintf(" ║ %-*s ║ %-*.40s ║ %-*.30s ║\n", regWidth, regStr, stackWidth, stackStr, msgWidth, msgStr))
	}
	bb.WriteString(uiS3)
	end := int(math.Min(float64(len(z.console)), float64(scrHeight)))
	for i := 0; i < end; i++ {
		bb.WriteString(fmt.Sprintf(" ║%-*.92s║\n", scrWidth, z.console[i]))
	}
	for i := scrHeight - end; i > 0; i-- {
		bb.WriteString(fmt.Sprintf(" ║%-*s║\n", scrWidth, ""))

	}
	bb.WriteString(uiS4)
	bb.WriteString(uiIn)
	bb.WriteString(uiS5)

	promptLine := " > "
	if z.prompt != "" {
		bb.WriteString(endDim)
		promptLine = fmt.Sprintf("\0331 | Requested input: %s > \0330", z.prompt)
	}

	bb.WriteString(fmt.Sprintf("  \033[2A \x1b[31m28z\033[0m %s %s", promptLine, string(z.runes)))

	return bb.Bytes()
}

var b2i = map[bool]int8{false: 0, true: 1}

func StateToString(r core.StateRegister) string {
	return fmt.Sprintf("%d%d%d", b2i[r.ResultFlag], b2i[r.BreakFlag], b2i[r.PromptFlag])
}
