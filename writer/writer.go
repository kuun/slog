package writer

import "github.com/kuun/slog/buffer"

type WriterType string

const (
	FILE   = "FILE"
)

type LogWriter interface {
	SetName(name string)
	GetName() string
	GetType() WriterType
	Write(buff *buffer.Buffer)
	Run()
	Close()
}



