package logger

import (
	"time"
)

type Logger struct {
	outputs []Output
}

func NewLogger(outputs ...Output) *Logger {
	return &Logger{
		outputs: outputs,
	}
}

func (l *Logger) Log(level Level, message string) {
	entry := LogEntry{
		Time:    time.Now(),
		Level:   level,
		Message: message,
	}

	for _, output := range l.outputs {
		_ = output.Write(entry)
	}
}

func (l *Logger) Debug(message string) {
	l.Log(Debug, message)
}

func (l *Logger) Info(message string) {
	l.Log(Info, message)
}

func (l *Logger) Warn(message string) {
	l.Log(Warn, message)
}

func (l *Logger) Error(message string) {
	l.Log(Error, message)
}
