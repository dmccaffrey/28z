package ui

import (
	"dmccaffrey/28z/core"
	"fmt"
	"math"
	"strings"
	"time"
)

var uiS0 = fmt.Sprintf("\n  ╓%s╖\n", strings.Repeat("─", 92))
var uiS1 = fmt.Sprintf("  ╟%s╥%s╢\n", strings.Repeat("─", 46), strings.Repeat("─", 45))
var headers = fmt.Sprintf("  ║ %-*s║ %-*s║\n", 45, "Registers", 44, "Stack")
var uiS2 = fmt.Sprintf("  ╟%s╫%s╢\n", strings.Repeat("┄", 46), strings.Repeat("┄", 45))
var uiS3 = fmt.Sprintf("  ╟%s╨%s╢\n", strings.Repeat("─", 46), strings.Repeat("─", 45))
var uiS4 = fmt.Sprintf("  ╟%s╢\n", strings.Repeat("─", 92))
var uiIn = fmt.Sprintf("  ║  %*s🮴 ║\n", 88, "")
var uiS5 = fmt.Sprintf("  ╙─%s╜\n", strings.Repeat("─", 91))

var lastUiUpdate = time.Now().Local()

func Display(vm *core.Core) string {
	var sb strings.Builder
	sb.WriteString(uiS0)
	sb.WriteString(fmt.Sprintf("  ║ \x1b[31m28z\033[0m ┇ Current Instruction = %-*s ║\n", 62, vm.LastInput))
	sb.WriteString(uiS1)
	sb.WriteString(headers)
	sb.WriteString(uiS2)
	stack := vm.GetStackArray()
	registerMap := vm.GetRegisterMap()
	for i := 4; i >= 0; i-- {
		regKey := core.RegisterKeys[i]
		registerStr := fmt.Sprintf("R%s: %023d", regKey, registerMap[regKey])
		stackValue := ""
		if i < len(stack) {
			stackValue = stack[i].GetString()
		}
		stackStr := fmt.Sprintf("%02d: %-*s", i, 38, stackValue)
		sb.WriteString(fmt.Sprintf("  ║ %-*s║ %-*s║\n", 45, registerStr, 44, stackStr))
	}
	sb.WriteString(uiS3)
	sb.WriteString(fmt.Sprintf("  ║ Status: 🯀 %-*s║\n", 81, vm.Message))

	sb.WriteString(uiS4)
	end := int(math.Min(float64(len(vm.Console)), 36))
	for i := 0; i < end; i++ {
		sb.WriteString(fmt.Sprintf("  ║%-*s║\n", 92, vm.Console[i]))
	}
	for i := 36 - end; i > 0; i-- {
		sb.WriteString(fmt.Sprintf("  ║%-*s║\n", 92, ""))

	}
	sb.WriteString(uiS4)
	sb.WriteString(uiIn)
	sb.WriteString(uiS5)
	sb.WriteString(fmt.Sprintf("    \033[2A %s > ", ""))
	return sb.String()
}
