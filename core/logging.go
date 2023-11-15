package core

import (
	"io"
	"log"
	"os"
)

var Logger = log.New(io.Discard, "", 0)

func LogToFile() {
	logfile, err := os.Create("28z.log")
	if err != nil {
		panic("Could not open log file")
	}
	Logger.SetOutput(logfile)
}
