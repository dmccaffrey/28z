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

var reader *bufio.Reader = bufio.NewReader(os.Stdin)

func main() {
	//enableDebug := flag.Bool("debug", false, "Enable debug output")
	help := flag.Bool("help", false, "Output help documentation")
	flag.Parse()
	if *help {
		OutputHelpDocumentation()
		return
	}
	vm := core.NewCore()
	writer := io.Writer(os.Stdout)

	display(&vm, writer)
	input, _ := getConsoleInput()
	for input != "exit" {
		if input == "run" || input == ">>" {
			vm.Mode = core.Running
		}
		vm.ProcessRaw(input)
		fmt.Printf("msg=%s", vm.Message)
		//stack := vm.GetStackString()
		//fmt.Println(stack)
		display(&vm, writer)
		input, _ = getConsoleInput()
	}
}

func getConsoleInput() (string, error) {
	input, err := reader.ReadString('\n')
	input = strings.TrimSuffix(input, "\n")
	return input, err
}

func clearConsole() {
	fmt.Print("\033[H\033[2J")
}

func display(vm *core.Core, writer io.Writer) {
	clearConsole()
	writer.Write([]byte(ui.Display(vm)))
}

func OutputHelpDocumentation() {
	core.OutputInstructionHelpDoc()
}
