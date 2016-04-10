package writer

import (
	"fmt"
	"github.com/kuun/slog/buffer"
	"os"
)

type fileWriter struct {
	wType     WriterType          // writer type
	name      string              // writer name
	file      *os.File            // log file
	cacheChn  chan *buffer.Buffer // cache buffers that will be writing.
	isRunning bool                // if writer's writing gorotine is running
}

const fileWriterCache = 50

func NewFileWriter(name, fileName string) (wr *fileWriter, err error) {
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
		wType:    FILE,
		name:     name,
		file:     file,
		cacheChn: make(chan *buffer.Buffer, fileWriterCache),
	}, nil
}

func (writer *fileWriter) SetName(name string) {
	writer.name = name
}

func (writer *fileWriter) GetName() string {
	return writer.name
}

func (writer *fileWriter) GetType() WriterType {
	return writer.wType
}

func (writer *fileWriter) Write(buff *buffer.Buffer) {
	writer.cacheChn <- buff
}

func (writer *fileWriter) Run() {
	if writer.isRunning {
		return
	}
	go func() {
		for {
			buff := <-writer.cacheChn
			if _, err := writer.file.Write(buff.Bytes()); err != nil {
				fmt.Fprintf(os.Stderr, "write log error: %s\n", err)
			}
			buffer.PutBuffer(buff)
		}
	}()
	writer.isRunning = true
}

func (writer *fileWriter) Close() {
	if writer.name == "STDOUT" || writer.name == "STDERR" {
		return
	}
	writer.file.Close()
}
