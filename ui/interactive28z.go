package ui

import (
	"dmccaffrey/28z/core"
	"strings"
	"time"

	"github.com/mattn/go-tty"
)

type Interactive28z struct {
	core         *core.Core
	tty          *tty.TTY
	lastInput    string
	prompt       string
	message      string
	console      []string
	runes        []rune
	ticker       time.Ticker
	input        chan rune
	lastUiUpdate time.Time
	run          bool
}

func NewInteractive28z(vm *core.Core) *Interactive28z {
	z := Interactive28z{}
	z.core = vm
	var err error
	z.tty, err = tty.Open()
	if err != nil {
		panic(err)
	}
	z.console = make([]string, 0, 30)
	z.lastInput = ""
	z.message = ""
	z.prompt = ""
	z.runes = make([]rune, 0, 128)
	z.ticker = *time.NewTicker(1 * time.Second)
	z.input = make(chan rune)
	z.run = true
	return &z
}

func (z *Interactive28z) Display() {
	z.tty.Output().Write(z.GenerateDebugUi())
	z.lastUiUpdate = time.Now()
}

func (z *Interactive28z) Output(line string) {
	z.console = append(z.console, line)
}

func (z *Interactive28z) Prompt(prompt string) {
	z.prompt = prompt
	z.Display()
}

func (z *Interactive28z) Clear() {
	z.console = z.console[:0]
}

func (z *Interactive28z) Run() {
	defer z.tty.Close()
	z.Display()
	go z.PollRune()
	z.HandleEvents()
}

func (z *Interactive28z) HandleEvents() {
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
		case r := <-z.input:
			if !z.HandleRune(r) {
				z.exit()
				return
			}
			z.Display()
		case <-z.ticker.C:
			z.Display()
		}
	}
}

func (z *Interactive28z) exit() {
	z.run = false
	z.core.Halt()
	z.ticker.Stop()
	z.tty.Output().WriteString("\n\n  Bailing out, you are on your own. Good Luck.\n")
}

func (z *Interactive28z) ReadLine() string {
	input, err := z.tty.ReadString()
	if err != nil {
		core.Logger.Printf("Failed to read input: err=%s\n", err.Error())
		return ""
	}
	return strings.TrimSuffix(input, "\n")
}

func (z *Interactive28z) PollRune() {
	for z.run {
		r, err := z.tty.ReadRune()
		if err != nil {
			core.Logger.Printf("Failed to read input: err=%s\n", err.Error())
			return
		}
		core.Logger.Printf("Read input: %s %d\n", string(r), r)
		z.input <- r
	}
}

func (z *Interactive28z) HandleRune(r rune) bool {
	switch r {
	case 127:
		if len(z.runes) <= 0 {
			return true
		}
		z.runes = z.runes[:len(z.runes)-1]
	case 13:
		input := string(z.runes)
		z.runes = z.runes[:0]
		if input == "exit" {

			return false
		}
		z.lastInput = input
		z.core.Input <- input
		z.prompt = ""
	default:
		z.runes = append(z.runes, r)
	}
	return true
}
