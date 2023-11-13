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

func Display(vm *core.Core) string {
	var sb strings.Builder
	sb.WriteString(uiS0)
	stack := vm.GetStackArray()
	registerMap := vm.GetRegisterMap()
	for i := 4; i >= 0; i-- {
		regKey := core.RegisterKeys[i]
		registerStr := fmt.Sprintf("%s: %04d", regKey, registerMap[regKey])
		stackValue := ""
		if i < len(stack) {
			stackValue = stack[i].GetString()
		}
		stackStr := fmt.Sprintf("%s %1d: %-*s", stackAliases[i], i, 58, stackValue)
		msgStr := ""
		switch i {
		case 0:
			msgStr = fmt.Sprintf("%-5s %s", "Last:", vm.LastInput)
			break
		case 1:
			msgStr = fmt.Sprintf("%-5s %s", "Err:", vm.Message)
			break
		case 2:
			msgStr = fmt.Sprintf("%-5s %s", "Mode:", vm.GetMode())
			break
		}
		sb.WriteString(fmt.Sprintf(" ║ %-*s ║ %-*.40s ║ %-*.30s ║\n", regWidth, registerStr, stackWidth, stackStr, msgWidth, msgStr))
	}
	sb.WriteString(uiS3)
	end := int(math.Min(float64(len(vm.Console)), float64(scrHeight)))
	for i := 0; i < end; i++ {
		sb.WriteString(fmt.Sprintf(" ║%-*.92s║\n", scrWidth, vm.Console[i]))
	}
	for i := scrHeight - end; i > 0; i-- {
		sb.WriteString(fmt.Sprintf(" ║%-*s║\n", scrWidth, ""))

	}
	sb.WriteString(uiS4)
	sb.WriteString(uiIn)
	sb.WriteString(uiS5)
	sb.WriteString(fmt.Sprintf("  \033[2A \x1b[31m28z\033[0m [%s]> ", vm.Prompt))
	vm.Prompt = ""
	return sb.String()
}
