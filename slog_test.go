package slog

import (
	"fmt"
	"os"
	"testing"
)

type computer struct {
	name string
	core int
}

func testAll(level string, prefix string) {
	logger := GetLogger()
	logger.SetLevel(level)
	testComputer := computer{"mycomputer", 4}
	//fmt.Printf("logger: %#v\n", logger)

	format := "my computer is %s, core num is %v"

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
}

func TestOut(t *testing.T) {
	// test write to file
	file, _ := os.OpenFile("/tmp/test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer file.Close()
	/*file, _:= os.Open("/tmp/test.log")*/
	testAll("DEBUG", "")
}

func TestDebug(t *testing.T) {
	fmt.Println("========================================DEBUG========================================")
	testAll("DEBUG", "")
}

func TestInfo(t *testing.T) {
	fmt.Println("========================================INFO========================================")
	testAll("INFO", "")
}

func TestNotice(t *testing.T) {
	fmt.Println("========================================NOTICE========================================")
	testAll("NOTICE", "")
}

func TestWarn(t *testing.T) {
	fmt.Println("========================================WARN========================================")
	testAll("WARN", "")
}

func TestError(t *testing.T) {
	fmt.Println("========================================ERROR========================================")
	testAll("ERROR", "")
}

func TestFatal(t *testing.T) {
	fmt.Println("========================================FATAL========================================")
	testAll("FATAL", "")
}

//func TestGetLogPath(t *testing.T) {
//	fullPath, shortPath := getLogPath()
//
//	expectFullPath := "github.com/kuun/slog"
//	expectShortPath := "g/k/slog"
//	if fullPath !=  expectFullPath{
//		t.Errorf("log full path is: %s, expect: %s", fullPath, expectFullPath)
//	}
//	if shortPath != expectShortPath {
//		t.Errorf("log short path is: %s, expect: %s", shortPath, expectShortPath)
//	}
//}

func TestLevel(t *testing.T) {
	t.Log("test log level")
	logger := GetLogger()

	if parseLevel(lvNameDebug) != lvDebug {
		t.Errorf("parse log level name error, level: %s", lvNameDebug)
	}
	if parseLevel(lvNameInfo) != lvInfo {
		t.Errorf("parse log level name error, level: %s", lvNameInfo)
	}
	if parseLevel(lvNameNotice) != lvNotice {
		t.Errorf("parse log level name error, level: %s", lvNameNotice)
	}
	if parseLevel(lvNameWarn) != lvWarn {
		t.Errorf("parse log level name error, level: %s", lvNameWarn)
	}
	if parseLevel(lvNameError) != lvError {
		t.Errorf("parse log level name error, level: %s", lvNameError)
	}
	if parseLevel(lvNameFatal) != lvFatal {
		t.Errorf("parse log level name error, level: %s", lvNameFatal)
	}
	for _, level := range []string{lvNameDebug, lvNameInfo, lvNameNotice, lvNameWarn, lvNameError, lvNameFatal} {
		logger.SetLevel(level)
		if logger.GetLevel() != level {
			t.Errorf("test log level error: %s", level)
		}
	}
	logger.SetLevel("DEBUG")
}
