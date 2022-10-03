package main

import(
	"time"
	"strings"
	"fmt"
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

func Display(s EnvState, instruction string, alwaysUpdate bool) {
	if !alwaysUpdate && time.Since(lastUiUpdate).Milliseconds() < 150 {
		return
	}
	lastUiUpdate = time.Now().Local()

	var sb strings.Builder
	sb.WriteString(uiS0)
	sb.WriteString(fmt.Sprintf("  ║ \x1b[31m28z\033[0m ┇ Current Instruction = %-*s ║\n", 62, instruction))
	sb.WriteString(uiS1)
	sb.WriteString(headers)
	sb.WriteString(uiS2)
	end := MaxStackLen-1
	for i := end; i >= 0; i-- {
		stackEntry := DefaultStackData()
		if (i < len(s.stack)) {
			stackIndex := len(s.stack) - i - 1
			stackEntry = s.stack[stackIndex]
		}
		registerStr := fmt.Sprintf("R%s: (%c) %-*s", string(i + 65), s.regs[i].dataType, 20, s.regs[i].ToString())
		stackStr := fmt.Sprintf("%d:", i)
		if stackEntry.dataType != Nil {
			stackStr = fmt.Sprintf("%d: (%c) %-*s", i, stackEntry.dataType, 35, stackEntry.ToString())
		}
		sb.WriteString(fmt.Sprintf("  ║ %-*s║ %-*s║\n", 45, registerStr, 44, stackStr))
	}
	sb.WriteString(uiS3)
	if s.err != "" {
		sb.WriteString(fmt.Sprintf("  ║ Status: 🯀 %-*s║\n", 81, s.err))

	} else {
		sb.WriteString(fmt.Sprintf("  ║ Status: 🮱 %-*s║\n", 81, "OK"))
	}
	sb.WriteString(uiS4)
	lines := strings.Split(s.console, "\n")
	for _,v := range lines {
		//max := math.Min(float64(len(v)), 92)
		//sb.WriteString(fmt.Sprintf("  ║%-*s║\n", 92, v[:int(max)]))
		sb.WriteString(fmt.Sprintf("  ║%-*s║\n", 92, v))
	}
	for i := 36-len(lines); i>0; i-- {
		sb.WriteString(fmt.Sprintf("  ║%-*s║\n", 92, ""))

	}
	sb.WriteString(uiS4)
	sb.WriteString(uiIn)
	sb.WriteString(uiS5)
	sb.WriteString("    \033[2A> ")
	s.writer.Publish(sb.String())
}
