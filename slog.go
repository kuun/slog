// A simple log wrapper for std log, support log level.
// usage reforence golang std log

package slog

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/kuun/slog/writer"
)

type Empty struct {
}

// Level is log level
type Level int

// log level definition
const (
	Debug = iota
	Info
	Notice
	Warn
	Error
	Fatal
)

// log level name
const (
	LvNameDebug  = "DEBUG"
	LvNameInfo   = "INFO"
	LvNameNotice = "NOTICE"
	LvNameWarn   = "WARN"
	LvNameError  = "ERROR"
	LvNameFatal  = "FATAL"
)

func (lv Level) String() string {
	switch lv {
	case Debug:
		return LvNameDebug
	case Info:
		return LvNameInfo
	case Notice:
		return LvNameNotice
	case Warn:
		return LvNameWarn
	case Error:
		return LvNameError
	case Fatal:
		return LvNameFatal
	default:
		return "UNKOWN"
	}
}

func parseLevel(strLevel string) (Level, bool) {
	switch strLevel {
	case LvNameDebug:
		return Debug, true
	case LvNameInfo:
		return Info, true
	case LvNameNotice:
		return Notice, true
	case LvNameWarn:
		return Warn, true
	case LvNameError:
		return Error, true
	case LvNameFatal:
		return Fatal, true
	default:
		return -1, false
	}
}

// Logger is the interface of slog, the interface likes the std log package
type Logger interface {
	// GetLevel gets log level of the logger.
	GetLevel() string
	// SetLevel sets log level of the logger.
	SetLevel(lv string) error

	// Is*Enabled tests the current log level is enabled or not
	IsDebugEnabled() bool
	IsInfoEnabled() bool
	IsNoticeEnabled() bool
	IsWarnEnabled() bool
	IsErrorEnabled() bool
	IsFatalEnabled() bool

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
var loggers = make(map[string]Logger)

// all log writers, indexed by writer name
var writers = make(map[string]writer.LogWriter)

var conf configration

func init() {
	var data []byte
	var err error

	confFile := os.Getenv("SLOG_CONF_FILE")
	if confFile != "" {
		if data, err = ioutil.ReadFile(confFile); err != nil {
			fmt.Printf("slog read file '%s' error: %s", confFile, err)
			goto FAIL
		}
		if err = json.Unmarshal(data, &conf); err != nil {
			fmt.Printf("slog parse config file '%s' error: %s", confFile, err)
			goto FAIL
		}
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
		conf.Loggers = append(conf.Loggers, logConf{
			Pattern: "*",
			Level:   LvNameDebug,
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
	case LvNameDebug, LvNameInfo, LvNameNotice, LvNameWarn, LvNameError, LvNameFatal:
		return nil
	default:
		return errors.New("unkown log level name: " + level)
	}
}

// GetLogger returns a logger assosiates to caller's package,
// the logger path is the caller's package path
func GetLogger(pkgInfo interface{}) Logger {
	fullPath := reflect.TypeOf(pkgInfo).PkgPath()
	abbrPath := makeAbbrPath(fullPath)
	if logger, ok := loggers[fullPath]; ok {
		return logger
	}
	return doGetLogger(fullPath, abbrPath)
}

// GetLoggerWithPath returns a logger with a specific path, recommands to use
// GetLogger instead of GetLoggerWithPath, unless the package path includes
// string 'src'.
func GetLoggerWithPath(path string) Logger {
	return doGetLogger(path, makeAbbrPath(path))
}

// Close flushes all data to files closes, this func should be called
// before application exits.
func Close() {
	for _, wr := range writers {
		wr.Close()
	}
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

func getLogPath() (fullPath string) {
	_, file, _, _ := runtime.Caller(2)
	lastSlashPos := strings.LastIndexByte(file, '/')
	srcPos := strings.Index(file, "/src/")
	return file[srcPos+5 : lastSlashPos]
}

func makeAbbrPath(fullPath string) (abbrPath string) {
	dirs := strings.Split(fullPath, "/")
	count := len(dirs)
	for i, dir := range dirs {
		if i != count-1 {
			abbrPath += dir[0:1]
			abbrPath += "/"
		} else {
			abbrPath += dir
		}
	}
	return abbrPath
}

func isWildMatch(pattern, str string) bool {
	if pattern == "*" {
		return true
	}
	patternLen := len(pattern)
	if pattern[patternLen-1] == '*' {
		return strings.HasPrefix(str, pattern[:patternLen-1])
	}
	return pattern == str
}
