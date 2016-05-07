package writer

import "github.com/kuun/slog/buffer"

// Type is log writer type
type Type string

const (
	// FILE represents file log writer
	FILE = "FILE"
)

// LogWriter is the interface of log writer, used to write log to somewhere
type LogWriter interface {
	SetName(name string)
	GetName() string
	GetType() Type
	Write(buff *buffer.Buffer)
	Run()
	Close()
}
