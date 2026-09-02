package logger

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileOutput struct {
	file *os.File
}

func NewFileOutput(path string) (*FileOutput, error) {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}

	return &FileOutput{
		file: file,
	}, nil
}

func (f *FileOutput) Write(entry LogEntry) error {
	_, err := fmt.Fprintf(
		f.file,
		"%s [%s] %s\n",
		entry.Time.Format("2006-01-02 15:04:05"),
		entry.Level,
		entry.Message,
	)

	return err
}

func (f *FileOutput) Close() error {
	return f.file.Close()
}
