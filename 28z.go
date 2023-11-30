package main

import (
	"dmccaffrey/28z/core"
	"dmccaffrey/28z/ui"
	"flag"
	"fmt"
)

func main() {
	//enableDebug := flag.Bool("debug", false, "Enable debug output")
	help := flag.Bool("help", false, "Output help documentation")
	eval := flag.String("eval", "", "Specify a reference to evaluate on start")
	flag.Parse()
	if *help {
		OutputHelpDocumentation()
		return
	}

	core.LogToFile()

	core.Logger.Printf("Initializing ROM\n")
	err := core.LoadRom()
	if err != nil {
		fmt.Printf("Failed to load ROM: %s\n", err.Error())
		return
	}

	core.Logger.Printf("Initializing instruction map\n")
	core.InitializeInstructionMap()

	z := ui.NewInteractive28z()
	core.Logger.Printf("Initializing core\n")
	vm := core.NewCore(&z)
	if *eval != "" {
		core.Logger.Printf("Evaluating initial input: input=%s\n", *eval)
		vm.ProcessRaw(*eval)
	}
	core.Logger.Printf("Starting main loop\n")
	vm.Mainloop()
}

func OutputHelpDocumentation() {
	core.OutputInstructionHelpDoc()
}
