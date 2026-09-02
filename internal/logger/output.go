package logger

type Output interface {
	Write(entry LogEntry) error
}
