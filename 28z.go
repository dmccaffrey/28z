package main

import (
	"bufio"
	"dmccaffrey/28z/core"
	"dmccaffrey/28z/ui"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Interactive28z struct {
	reader *bufio.Reader
	writer io.Writer
}

func NewInteractive28z() Interactive28z {
	z := Interactive28z{}
	z.reader = bufio.NewReader(os.Stdin)
	z.writer = io.Writer(os.Stdout)
	return z
}

func main() {
	//enableDebug := flag.Bool("debug", false, "Enable debug output")
	help := flag.Bool("help", false, "Output help documentation")
	flag.Parse()
	if *help {
		OutputHelpDocumentation()
		return
	}

	err := core.LoadRom()
	if err != nil {
		fmt.Printf("Failed to load ROM: %s\n", err.Error())
		return
	}

	core.InitializeInstructionMap()

	z := NewInteractive28z()
	vm := core.NewCore()
	vm.Mainloop(&core.InteractiveHandler{Input: z.input, Output: z.output})
}

func (z *Interactive28z) input(vm *core.Core) (bool, string) {
	input, err := z.reader.ReadString('\n')
	if err != nil {
		vm.Message = "Failed to read input"
		return true, ""
	}
	input = strings.TrimSuffix(input, "\n")
	if input == "exit" {
		return false, ""
	}
	return true, input
}

func (z *Interactive28z) output(vm *core.Core) {
	fmt.Print("\033[H\033[2J")
	z.writer.Write([]byte(ui.Display(vm)))
}

func OutputHelpDocumentation() {
	core.OutputInstructionHelpDoc()
}
