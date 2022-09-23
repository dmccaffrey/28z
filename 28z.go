package main

import (
	"fmt"
	"strings"
	"io"
	"os"
	"unsafe"
	"syscall"
	"bufio"
)

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
	//out, size := AcquireTty()
	return TermWriter{0, io.Writer(os.Stdout), windowSize{0,0}, true}
}

func DebugTermWriter() TermWriter {
	return TermWriter{0, io.Writer(os.Stdout), windowSize{0,0}, false}
}

func (w *TermWriter) Publish(content string) {
	if w.interactive {
		var clear = fmt.Sprintf("%c[%dA%c[2K", 27, 1, 27)
		_, err := fmt.Fprint(w.output, strings.Repeat(clear, w.lastLineCount))
		if err != nil {
			fmt.Printf("Failed to reset output: err=%s", err.Error())
		}
	}

	// +1 for the newlilne caused by input
	w.lastLineCount = strings.Count(content, "\n") + 1
	fmt.Fprintf(w.output, "%s", content)
}

func AcquireTty() (io.Writer, windowSize) {
	out, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Failed to acquire TTY")
		return io.Writer(os.Stdout), windowSize{0,0}
	}
	writer := io.Writer(out)
	var winsize windowSize
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL,
		out.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&winsize)))
	fmt.Fprint(writer, "\033[H\033[2J")
	return writer, winsize
}

func main() {
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
