package logger

import (
	"fmt"
)

type ConsoleOutput struct{}

func NewConsoleOutput() *ConsoleOutput {
	return &ConsoleOutput{}
}

func (c *ConsoleOutput) Write(entry LogEntry) error {
	fmt.Printf(
		"%s [%s] %s\n",
		"["+entry.Time.Format("2006-01-02 15:04:05")+"]",
		entry.Level.String(),
		entry.Message,
	)

	return nil
}
