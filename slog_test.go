package slog

import (
	"fmt"
	"testing"
	"time"
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

	time.Sleep(100 * time.Millisecond)

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

func TestDebug(t *testing.T) {
	fmt.Println("======================DEBUG======================")
	testAll("DEBUG", "")
}

func TestInfo(t *testing.T) {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("======================INFO======================")
	testAll("INFO", "")
}

func TestNotice(t *testing.T) {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("======================NOTICE======================")
	testAll("NOTICE", "")
}

func TestWarn(t *testing.T) {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("======================WARN======================")
	testAll("WARN", "")
}

func TestError(t *testing.T) {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("======================ERROR======================")
	testAll("ERROR", "")
}

func TestFatal(t *testing.T) {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("======================FATAL======================")
	testAll("FATAL", "")
}

func TestGetLogPath(t *testing.T) {
	fullPath := getLogPath()
	expectFullPath := "testing"

	if fullPath != expectFullPath {
		t.Errorf("log full path is: %s, expect: %s", fullPath, expectFullPath)
	}
}

func TestLevel(t *testing.T) {
	t.Log("test log level")
	logger := GetLogger()

	if lv, _ := parseLevel(LvNameDebug); lv != Debug {
		t.Errorf("parse log level name error, level: %s", LvNameDebug)
	}
	if lv, _ := parseLevel(LvNameInfo); lv != Info {
		t.Errorf("parse log level name error, level: %s", LvNameInfo)
	}
	if lv, _ := parseLevel(LvNameNotice); lv != Notice {
		t.Errorf("parse log level name error, level: %s", LvNameNotice)
	}
	if lv, _ := parseLevel(LvNameWarn); lv != Warn {
		t.Errorf("parse log level name error, level: %s", LvNameWarn)
	}
	if lv, _ := parseLevel(LvNameError); lv != Error {
		t.Errorf("parse log level name error, level: %s", LvNameError)
	}
	if lv, _ := parseLevel(LvNameFatal); lv != Fatal {
		t.Errorf("parse log level name error, level: %s", LvNameFatal)
	}
	for _, level := range []string{LvNameDebug, LvNameInfo, LvNameNotice, LvNameWarn, LvNameError, LvNameFatal} {
		logger.SetLevel(level)
		if logger.GetLevel() != level {
			t.Errorf("test log level error: %s", level)
		}
	}
}

func TestMakeAbbrPath(t *testing.T) {
	abbrPath := makeAbbrPath("github.com/kuun/src/test")
	if abbrPath != "g/k/s/test" {
		t.Errorf("make abbravitated path error, abbr path: %s", abbrPath)
	}
}

func TestClose(t *testing.T) {
	Close()
}
