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

	logger.Debug("#Debug# ", "my computer is ", testComputer.name, ", core num is ", testComputer.core, "\n")
	logger.Debugf("#Debugf# "+format, testComputer.name, testComputer.core)
	logger.Debugln("#Debugln# ", "my computer is ", testComputer.name, ", core num is ", testComputer.core)

	logger.Info("#Info# ", "my computer is ", testComputer.name, ", core num is ", testComputer.core, "\n")
	logger.Infof("#Infof# "+format, testComputer.name, testComputer.core)
	logger.Infoln("#Infoln# ", "my computer is ", testComputer.name, ", core num is ", testComputer.core)

	logger.Notice("#Notice# ", "my computer is ", testComputer.name, ", core num is ", testComputer.core, "\n")
	logger.Noticef("#Noticef# "+format, testComputer.name, testComputer.core)
	logger.Noticeln("#Noticeln# ", "my computer is ", testComputer.name, ", core num is ", testComputer.core)

	logger.Warn("#Warn# ", "my computer is ", testComputer.name, ", core num is ", testComputer.core, "\n")
	logger.Warnf("#Warnf# "+format, testComputer.name, testComputer.core)
	logger.Warnln("#Warnln# ", "my computer is ", testComputer.name, ", core num is ", testComputer.core)

	logger.Error("#Error# ", "my computer is ", testComputer.name, ", core num is ", testComputer.core, "\n")
	logger.Errorf("#Errorf# "+format, testComputer.name, testComputer.core)
	logger.Errorln("#Errorln# ", "my computer is ", testComputer.name, ", core num is ", testComputer.core)

	/*logger.Fatal("#Fatal# ", "my computer is ", testComputer.name, ", core num is ", testComputer.core, "\n")*/
	/*logger.Fatalf("#Fatalf# " + format, testComputer.name, testComputer.core)*/
	/*logger.Fatalln("#Fatalln# ", "my computer is ", testComputer.name, ", core num is ", testComputer.core)*/

	/*logger.Panic("#Panic# ", "my computer is ", testComputer.name, ", core num is ", testComputer.core, "\n")*/
	/*logger.Panicf("#Panicf# " + format, testComputer.name, testComputer.core)*/
	/*logger.Panicln("#Panicln# ", "my computer is ", testComputer.name, ", core num is ", testComputer.core)*/
}

func TestOut(t *testing.T) {
	// test write to file
	file, _ := os.OpenFile("/tmp/test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
