// A simple log wrapper for std log, support log level.
// usage reforence golang std log

package slog

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/kuun/slog/writer"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

// Level is log level
type level int

func (lv level) String() string {
	switch lv {
	case lvDebug:
		return "DEBUG"
	case lvInfo:
		return "INFO"
	case lvNotice:
		return "NOTICE"
	case lvWarn:
		return "WARN"
	case lvError:
		return "ERROR"
	case lvFatal:
		return "FATAL"
	default:
		return "UNKOWN"
	}
}

// log level definition
const (
	lvDebug = iota
	lvInfo
	lvNotice
	lvWarn
	lvError
	lvFatal
	lvPanic // need not access it
)

// log level name
const (
	lvNameDebug  = "DEBUG"
	lvNameInfo   = "INFO"
	lvNameNotice = "NOTICE"
	lvNameWarn   = "WARN"
	lvNameError  = "ERROR"
	lvNameFatal  = "FATAL"
)

type Logger interface {
	GetLevel() string
	SetLevel(lv string)

	Debug(v ...interface{})
	Debugf(fmt string, v ...interface{})

	Info(v ...interface{})
	Infof(fmt string, v ...interface{})

	Notice(v ...interface{})
	Noticef(fmt string, v ...interface{})

	Warn(v ...interface{})
	Warnf(fmt string, v ...interface{})

	Error(v ...interface{})
	Errorf(fmt string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(fmt string, v ...interface{})
}

type writerConf struct {
	// Name is log writer name
	Name string `json:"name"`
	// Type is log writer type, valid value: "FILE"
	// note: type "STD" is used only by slog, user can't use it
	Type string `json:"type"`
	// File is a log file, valid only when the Type is "FILE"
	File string `json:"file"`
}

type logConf struct {
	Pattern string   `json:"pattern"`
	Level   string   `json:"level"`
	Writers []string `json:"writers"`
}

type configration struct {
	Writers []writerConf `json:"writers"`
	Loggers []logConf    `json:"loggers"`
}

// all loggers, indexed by logger full path
var loggers map[string]Logger = make(map[string]Logger)

// all log writers, indexed by writer name
var writers map[string]writer.LogWriter = make(map[string]writer.LogWriter)

var confFile string

var conf configration

func init() {
	var data []byte
	var err error

	flag.StringVar(&confFile, "slog", "", "config file for slog")
	flag.Parse()
	if confFile == "" {
		flag.Usage()
		goto FAIL
	}
	if data, err = ioutil.ReadFile(confFile); err != nil {
		fmt.Printf("slog read file '%s' error: %s", confFile, err)
		goto FAIL
	}
	if err = json.Unmarshal(data, &conf); err != nil {
		fmt.Printf("slog parse config file '%s' error: %s", confFile, err)
		goto FAIL
	}
	if err = initWriters(); err != nil {
		fmt.Printf("slog init writers error: %s", err)
		goto FAIL
	}
	if err = verifyConf(); err != nil {
		fmt.Printf("config file is invalid: '%s'", err)
		goto FAIL
	}
	return
FAIL:
	os.Exit(1)
}

func initWriters() error {
	for _, wrConf := range conf.Writers {
		if wrConf.Type == writer.FILE {
			return errors.New("not valid writer type: " + wrConf.Type)
		}
		wr, err := writer.NewFileWriter(wrConf.Name, wrConf.File)
		if err != nil {
			return err
		}
		if writers[wrConf.Name] != nil {
			return errors.New("writer name is duplicated: " + wrConf.Name)
		}
		writers[wrConf.Name] = wr
	}
	return nil
}

func verifyConf() error {
	if len(conf.Loggers) == 0 {
		conf.Loggers = append(conf.Loggers,logConf{
			Pattern: "*",
			Level: lvNameDebug,
			Writers: []string{"STDOUT"},
		})
	}
	for _, logger := range conf.Loggers {
		if err := verifyLogLevel(logger.Level); err != nil {
			return err
		}
		for _, writerName := range logger.Writers {
			if wr := writers[writerName]; wr == nil {
				switch writerName {
				case "STDOUT", "STDERR":
					writers[writerName], _ = writer.NewFileWriter(writerName, "")
				default:
					return errors.New("can't find log writer: " + writerName)
				}
			}
		}
	}

	return nil
}

func verifyLogLevel(level string) error {
	switch level {
	case lvNameDebug, lvNameInfo, lvNameNotice, lvNameWarn, lvNameError, lvNameFatal:
		return nil
	default:
		return errors.New("unkown log level name: " + level)
	}
}

func GetLogger() Logger {
	fullPath, abbrPath := getLogPath()
	return doGetLogger(fullPath, abbrPath)
}

func GetLoggerWithPath(path string) Logger {
	return doGetLogger(path, path)
}

func doGetLogger(fullPath, abbrPath string) Logger {
	for _, logConf := range conf.Loggers {
		if isWildMatch(logConf.Pattern, fullPath) {
			var logger Logger
			logger = loggers[fullPath]
			if logger == nil {
				logger = &loggerImpl{
					fullPath: fullPath,
					abbrPath: abbrPath,
					writers:  getLogWriters(logConf.Writers),
				}
				logger.SetLevel(logConf.Level)
				loggers[fullPath] = logger
			}
			return logger
		}
	}
	panic("should not arrive here!")
	return nil
}

func getLogWriters(writerNames []string) []writer.LogWriter {
	size := len(writerNames)
	wrs := make([]writer.LogWriter, 0, size)
	for _, wrConf := range writerNames {
		wr := writers[wrConf]
		wr.Run()
		wrs = append(wrs, wr)
	}
	return wrs
}

func getLogPath() (fullPath, abbrPath string) {
	_, file, _, _ := runtime.Caller(2)
	lastSlashPos := strings.LastIndexByte(file, '/')
	srcPos := strings.Index(file, "/src/")
	fullPath = file[srcPos+5 : lastSlashPos]

	dirs := strings.Split(fullPath, "/")
	count := len(dirs)
	abbrPath = ""
	for i, dir := range dirs {
		if i != count-1 {
			abbrPath += dir[0:1]
			abbrPath += "/"
		} else {
			abbrPath += dir
		}
	}
	return fullPath, abbrPath
}

func isWildMatch(pattern, str string) bool {
	if pattern == "*" {
		return true
	}
	patternLen := len(pattern)
	if pattern[patternLen-1] == '*' {
		return strings.HasPrefix(str, pattern[:patternLen-1])
	} else {
		return pattern == str
	}
}

func parseLevel(strLevel string) level {
	switch strLevel {
	case "DEBUG":
		return lvDebug
	case "INFO":
		return lvInfo
	case "NOTICE":
		return lvNotice
	case "WARN":
		return lvWarn
	case "ERROR":
		return lvError
	case "FATAL":
		return lvFatal
	default:
		panic("unkown log level: " + strLevel)
	}
}
