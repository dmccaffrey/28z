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
	core      *core.Core
	reader    *bufio.Reader
	writer    io.Writer
	lastInput string
	prompt    string
	message   string
	console   []string
}

func NewInteractive28z(vm *core.Core) *Interactive28z {
	z := Interactive28z{}
	z.core = vm
	z.reader = bufio.NewReader(os.Stdin)
	z.writer = io.Writer(os.Stdout)
	z.console = make([]string, 0, 30)
	z.lastInput = ""
	z.message = ""
	z.prompt = ""
	return &z
}

func (z *Interactive28z) Display() {
	fmt.Print("\033[H\033[2J")
	z.writer.Write(z.GenerateDebugUi())
}

func (z *Interactive28z) Output(line string) {
	z.console = append(z.console, line)
}

func (z *Interactive28z) Prompt(prompt string) {
	z.prompt = prompt
}

func (z *Interactive28z) Clear() {
	z.console = z.console[:0]
}

func (z *Interactive28z) Run() {
	go z.HandleEvents()
	z.PollInput()
}

func (z *Interactive28z) HandleEvents() {
	z.Display()
	for {
		select {
		case message := <-z.core.Control:
			switch message.Command {
			case core.Stop:
				core.Logger.Println("Stopping interactive UI event handler")
				return
			case core.Clear:
				z.Clear()
			case core.Prompt:
				z.Prompt(message.Arg)
			case core.Output:
				z.Output(message.Arg)
			case core.StateUpdated:
				z.Display()
			}
		}
	}
}

func (z *Interactive28z) PollInput() {
	for {
		input, err := z.reader.ReadString('\n')
		if err != nil {
			core.Logger.Printf("Failed to read input: err=%s\n", err.Error())
			continue
		}
		input = strings.TrimSuffix(input, "\n")

		if input == "exit" {
			z.core.Control <- core.CommandMessage{Command: core.Stop, Arg: ""}
			return
		}
		z.lastInput = input
		z.core.Input <- input
	}
}
