package logger

import "time"

type LogEntry struct {
	Time    time.Time
	Level   Level
	Message string
}
