package main

import (
	"fmt"
	"strings"
	"io"
	"os"
	"bufio"
)

var clear = fmt.Sprintf("%c[%dA%c[2K", 27, 1, 27)

type windowSize struct {
	rows    uint16
	cols    uint16
}

type TermWriter struct {
	lastLineCount int
	output io.Writer
	winSize windowSize
	interactive bool
}

func InteractiveTermWriter() TermWriter {
	return TermWriter{0, io.Writer(os.Stdout), windowSize{0,0}, true}
}

func DebugTermWriter() TermWriter {
	return TermWriter{0, io.Writer(os.Stdout), windowSize{0,0}, false}
}

func (w *TermWriter) Publish(content string) {
	if w.interactive {
		w.output.Write([]byte(strings.Repeat(clear, w.lastLineCount)))
		w.output.Write([]byte("\033[2J"))
	}

	// +1 for the newlilne caused by input
	w.lastLineCount = strings.Count(content, "\n") + 1
	w.output.Write([]byte(content))
}

func main() {
	loadRom()
	//os.Exit(0)
	//state := NewEnvState(DebugTermWriter())
	state := NewEnvState(InteractiveTermWriter())
	state.Display("")
	in := bufio.NewReader(os.Stdin)
	for {
		input, err := in.ReadString('\n')
		if err == io.EOF {
			return
		}
		input = strings.TrimSuffix(input, "\n")
		if input == "exit" {
			break
		}
		if !state.Parse(input) {
			break
		}
		state.Display(input)
	}
	state.Display("exit")
}
