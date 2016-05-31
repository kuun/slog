package writer

import (
	"fmt"
	"os"

	"github.com/kuun/slog/buffer"
)

type fileWriter struct {
	wType     Type                // writer type
	name      string              // writer name
	file      *os.File            // log file
	cacheChn  chan *buffer.Buffer // cache buffers that will be writing.
	flushDone chan bool           // all cached buffers are flush to file
	isRunning bool                // if writer's writing gorotine is running
}

const fileWriterCache = 50

// NewFileWriter creates a new file log writer
func NewFileWriter(name, fileName string) (wr LogWriter, err error) {
	var file *os.File
	switch name {
	case "STDOUT":
		file = os.Stdout
	case "STDERR":
		file = os.Stderr
	default:
		file, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666|os.ModeAppend)
		if err != nil {
			return nil, err
		}
	}
	return &fileWriter{
		wType:     FILE,
		name:      name,
		file:      file,
		cacheChn:  make(chan *buffer.Buffer, fileWriterCache),
		flushDone: make(chan bool),
	}, nil
}

func (writer *fileWriter) SetName(name string) {
	writer.name = name
}

func (writer *fileWriter) GetName() string {
	return writer.name
}

func (writer *fileWriter) GetType() Type {
	return writer.wType
}

func (writer *fileWriter) Write(buff *buffer.Buffer) {
	writer.cacheChn <- buff
}

// Run starts a go routine to write buffers to file
func (writer *fileWriter) Run() {
	if writer.isRunning {
		return
	}
	go func() {
		for {
			buff := <-writer.cacheChn
			if buff == nil {
				writer.flushDone <- true
				continue
			}
			if _, err := writer.file.Write(buff.Bytes()); err != nil {
				fmt.Fprintf(os.Stderr, "write log error: %s\n", err)
			}
			buffer.PutBuffer(buff)
		}
	}()
	writer.isRunning = true
}

func (writer *fileWriter) Close() {
	writer.Write(nil)
	_ = <-writer.flushDone
	if writer.name == "STDOUT" || writer.name == "STDERR" {
		return
	}
	writer.file.Close()
}
