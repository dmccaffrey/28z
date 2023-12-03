package ui

import (
	"bufio"
	"dmccaffrey/28z/core"
	"fmt"
	"io"
	"os"
	"strings"
)

type Interactive28z struct {
	reader    *bufio.Reader
	writer    io.Writer
	lastInput string
	prompt    string
	message   string
	console   []string
}

func NewInteractive28z() Interactive28z {
	z := Interactive28z{}
	z.reader = bufio.NewReader(os.Stdin)
	z.writer = io.Writer(os.Stdout)
	z.console = make([]string, 0)
	z.lastInput = ""
	z.message = ""
	z.prompt = ""
	return z
}

func (z *Interactive28z) GetInput(vm *core.Core) (bool, string) {
	input, err := z.reader.ReadString('\n')
	if err != nil {
		core.Logger.Printf("Error: Failed to read input\n")
		z.message = "Failed to read input"
		return true, ""
	}
	input = strings.TrimSuffix(input, "\n")
	if input == "exit" {
		return false, ""
	}
	z.lastInput = input
	return true, input
}

func (z *Interactive28z) Display(vm *core.Core) {
	fmt.Print("\033[H\033[2J")
	z.writer.Write([]byte(z.DisplayDebugUi(vm)))
	z.prompt = ""
}

func (z *Interactive28z) Output(line string) {
	z.console = append(z.console, line)
}

func (z *Interactive28z) Prompt(vm *core.Core, prompt string) {
	z.prompt = prompt
}

func (z *Interactive28z) Clear() {
	z.console = make([]string, 0)
}
