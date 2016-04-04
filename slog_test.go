package slog

import (
	"fmt"
	"io"
	"os"
	"testing"
)

type computer struct {
	name string
	core int
}

func testAll(out io.Writer, level string, prefix string) {
	logger := GetLogger()

	testComputer := computer{"mycomputer", 4}

	format := "my computer is %s, core num is %v\n"

	logger.Debug("my computer is ", testComputer.name, ", core num is ", testComputer.core)
	logger.Debugf(format, testComputer.name, testComputer.core)

	logger.Info("my computer is ", testComputer.name, ", core num is ", testComputer.core)
	logger.Infof(format, testComputer.name, testComputer.core)

	logger.Notice("my computer is ", testComputer.name, ", core num is ", testComputer.core)
	logger.Noticef(format, testComputer.name, testComputer.core)

	logger.Warn("my computer is ", testComputer.name, ", core num is ", testComputer.core)
	logger.Warnf(format, testComputer.name, testComputer.core)

	logger.Error("my computer is ", testComputer.name, ", core num is ", testComputer.core)
	logger.Errorf(format, testComputer.name, testComputer.core)

	/*logger.Fatal("my computer is ", testComputer.name, ", core num is ", testComputer.core)*/
	/*logger.Fatalf(format, testComputer.name, testComputer.core)*/

	/*logger.Panic("my computer is ", testComputer.name, ", core num is ", testComputer.core)*/
	/*logger.Panicf(format, testComputer.name, testComputer.core)*/
}

func TestOut(t *testing.T) {
	// test write to file
	file, _ := os.OpenFile("/tmp/test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer file.Close()
	/*file, _:= os.Open("/tmp/test.log")*/
	testAll(file, "DEBUG", "")
}

func TestDebug(t *testing.T) {
	fmt.Println("========================================DEBUG========================================")
	testAll(os.Stdout, "DEBUG", "")
}

func TestInfo(t *testing.T) {
	fmt.Println("========================================INFO========================================")
	testAll(os.Stdout, "INFO", "")
}

func TestNotice(t *testing.T) {
	fmt.Println("========================================NOTICE========================================")
	testAll(os.Stdout, "NOTICE", "")
}

func TestWarn(t *testing.T) {
	fmt.Println("========================================WARN========================================")
	testAll(os.Stdout, "WARN", "")
}

func TestError(t *testing.T) {
	fmt.Println("========================================ERROR========================================")
	testAll(os.Stdout, "ERROR", "")
}

func TestFatal(t *testing.T) {
	fmt.Println("========================================FATAL========================================")
	testAll(os.Stdout, "FATAL", "")
}

