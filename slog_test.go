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

func testAll(out io.Writer, level Level, prefix string, flag int) {
	logger := New(out, level, prefix, flag)

	testComputer := computer{"mycomputer", 4}

	format := "my computer is %s, core num is %v\n"

	logger.Debug("my computer is ", testComputer.name, ", core num is ", testComputer.core)
	logger.Debugf(format, testComputer.name, testComputer.core)
	logger.Debugln("my computer is", testComputer.name, ", core num is", testComputer.core)

	logger.Info("my computer is ", testComputer.name, ", core num is ", testComputer.core)
	logger.Infof(format, testComputer.name, testComputer.core)
	logger.Infoln("my computer is", testComputer.name, ", core num is", testComputer.core)

	logger.Notice("my computer is ", testComputer.name, ", core num is ", testComputer.core)
	logger.Noticef(format, testComputer.name, testComputer.core)
	logger.Noticeln("my computer is", testComputer.name, ", core num is", testComputer.core)

	logger.Warn("my computer is ", testComputer.name, ", core num is ", testComputer.core)
	logger.Warnf(format, testComputer.name, testComputer.core)
	logger.Warnln("my computer is", testComputer.name, ", core num is", testComputer.core)

	logger.Error("my computer is ", testComputer.name, ", core num is ", testComputer.core)
	logger.Errorf(format, testComputer.name, testComputer.core)
	logger.Errorln("my computer is", testComputer.name, ", core num is", testComputer.core)

	/*logger.Fatal("my computer is ", testComputer.name, ", core num is ", testComputer.core)*/
	/*logger.Fatalf(format, testComputer.name, testComputer.core)*/
	/*logger.Fatalln("my computer is", testComputer.name, ", core num is", testComputer.core)*/

	/*logger.Panic("my computer is ", testComputer.name, ", core num is ", testComputer.core)*/
	/*logger.Panicf(format, testComputer.name, testComputer.core)*/
	/*logger.Panicln("my computer is", testComputer.name, ", core num is", testComputer.core)*/
}

func TestOut(t *testing.T) {
	// test write to file
	file, _ := os.OpenFile("/tmp/test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer file.Close()
	/*file, _:= os.Open("/tmp/test.log")*/
	testAll(file, DEBUG, "", Ldate|Ltime|Lmicroseconds|Lshortfile)
}

func TestDebug(t *testing.T) {
	fmt.Println("========================================DEBUG========================================")
	testAll(os.Stdout, DEBUG, "", Ldate|Ltime|Lmicroseconds|Lshortfile)
}

func TestInfo(t *testing.T) {
	fmt.Println("========================================INFO========================================")
	testAll(os.Stdout, INFO, "", Ldate|Ltime|Lmicroseconds|Lshortfile)
}

func TestNotice(t *testing.T) {
	fmt.Println("========================================NOTICE========================================")
	testAll(os.Stdout, NOTICE, "", Ldate|Ltime|Lmicroseconds|Lshortfile)
}

func TestWarn(t *testing.T) {
	fmt.Println("========================================WARN========================================")
	testAll(os.Stdout, WARN, "", Ldate|Ltime|Lmicroseconds|Lshortfile)
}

func TestError(t *testing.T) {
	fmt.Println("========================================ERROR========================================")
	testAll(os.Stdout, ERROR, "", Ldate|Ltime|Lmicroseconds|Lshortfile)
}

func TestFatal(t *testing.T) {
	fmt.Println("========================================FATAL========================================")
	testAll(os.Stdout, FATAL, "", Ldate|Ltime|Lmicroseconds|Lshortfile)
}

func TestDefaultLog(t *testing.T) {
	fmt.Println("========================================Default DEBUG========================================")

	levels := []Level{DEBUG, INFO, NOTICE, WARN, ERROR, FATAL}
	format := "my computer is %s, core num is %v\n"
	for _, level := range levels {
		fmt.Printf("Log Level: %s\n", levelName[level])
		testComputer := computer{"mycomputer", 4}
		SetLevel(level)

		Debug("my computer is ", testComputer.name, ", core num is ", testComputer.core)
		Debugf(format, testComputer.name, testComputer.core)
		Debugln("my computer is", testComputer.name, ", core num is", testComputer.core)

		Info("my computer is ", testComputer.name, ", core num is ", testComputer.core)
		Infof(format, testComputer.name, testComputer.core)
		Infoln("my computer is", testComputer.name, ", core num is", testComputer.core)

		Notice("my computer is ", testComputer.name, ", core num is ", testComputer.core)
		Noticef(format, testComputer.name, testComputer.core)
		Noticeln("my computer is", testComputer.name, ", core num is", testComputer.core)

		Warn("my computer is ", testComputer.name, ", core num is ", testComputer.core)
		Warnf(format, testComputer.name, testComputer.core)
		Warnln("my computer is", testComputer.name, ", core num is", testComputer.core)

		Error("my computer is ", testComputer.name, ", core num is ", testComputer.core)
		Errorf(format, testComputer.name, testComputer.core)
		Errorln("my computer is", testComputer.name, ", core num is", testComputer.core)
	}
}
