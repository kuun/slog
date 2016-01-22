// A simple log wrapper for std log, support log level.
// usage reforence golang std log

package slog

import (
	"fmt"
	"io"
	stdLog "log"
	"os"
)

// Level is log level
type Level int

// log level definition
const (
	DEBUG = iota
	INFO
	NOTICE
	WARN
	ERROR
	FATAL
	_PANIC // need not access it
)

var levelName = []string{"[DEBUG] ", "[INFO] ", "[NOTICE] ", "[WARN] ", "[ERROR] ", "[FATAL] ", "[PANIC] "}

// From std log.
// These flags define which text to prefix to each log entry generated by the Logger.
const (
	// Bits or'ed together to control what's printed. There is no control over the
	// order they appear (the order listed here) or the format they present (as
	// described in the comments).  A colon appears after these items:
	//	2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
	Ldate         = 1 << iota     // the date: 2009/01/23
	Ltime                         // the time: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

type Slogger struct {
	logger *stdLog.Logger
	level  Level
	// log.Logger.Output callPath
	callPath int

	debugImpl   func(l *Slogger, level Level, v ...interface{})
	debugfImpl  func(l *Slogger, level Level, format string, v ...interface{})
	debuglnImpl func(l *Slogger, level Level, v ...interface{})

	infoImpl   func(l *Slogger, level Level, v ...interface{})
	infofImpl  func(l *Slogger, level Level, format string, v ...interface{})
	infolnImpl func(l *Slogger, level Level, v ...interface{})

	noticeImpl   func(l *Slogger, level Level, v ...interface{})
	noticefImpl  func(l *Slogger, level Level, format string, v ...interface{})
	noticelnImpl func(l *Slogger, level Level, v ...interface{})

	warnImpl   func(l *Slogger, level Level, v ...interface{})
	warnfImpl  func(l *Slogger, level Level, format string, v ...interface{})
	warnlnImpl func(l *Slogger, level Level, v ...interface{})

	errorImpl   func(l *Slogger, level Level, v ...interface{})
	errorfImpl  func(l *Slogger, level Level, format string, v ...interface{})
	errorlnImpl func(l *Slogger, level Level, v ...interface{})

	fatalImpl   func(l *Slogger, level Level, v ...interface{})
	fatalfImpl  func(l *Slogger, level Level, format string, v ...interface{})
	fatallnImpl func(l *Slogger, level Level, v ...interface{})

	panicImpl   func(l *Slogger, level Level, v ...interface{})
	panicfImpl  func(l *Slogger, level Level, format string, v ...interface{})
	paniclnImpl func(l *Slogger, level Level, v ...interface{})
}

// New will create a new Slogger object
// level's value: "DEBUG", "INFO", "NOTICE", "WARN", "ERROR", "FATAL"
func New(out io.Writer, level string, prefix string, flag int) *Slogger {
	l := &Slogger{
		logger:   stdLog.New(out, prefix, flag),
		callPath: 3,
	}

	lv := parseLevel(level)
	l.SetLevel(lv)
	return l
}

func parseLevel(strLevel string) Level {
	switch strLevel {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "NOTICE":
		return NOTICE
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		panic("unkown log level: " + strLevel)
	}
}

// func suffix is "Y" is valid implements
// func suffix is "N" is empty implements

func printImplY(l *Slogger, level Level, v ...interface{}) {
	l.logger.Output(l.callPath, levelName[level]+fmt.Sprint(v...))
}

func printImplN(l *Slogger, level Level, v ...interface{}) {
}

func printfImplY(l *Slogger, level Level, format string, v ...interface{}) {
	l.logger.Output(l.callPath, fmt.Sprintf(levelName[level]+format, v...))
}

func printfImplN(l *Slogger, level Level, format string, v ...interface{}) {
}

func printlnImplY(l *Slogger, level Level, v ...interface{}) {
	l.logger.Output(l.callPath, levelName[level]+fmt.Sprintln(v...))
}

func printlnImplN(l *Slogger, level Level, v ...interface{}) {
}

func (l *Slogger) SetLevel(level Level) {
	l.level = level
	// init all print func
	l.debugImpl = printImplN
	l.debugfImpl = printfImplN
	l.debuglnImpl = printlnImplN

	l.infoImpl = printImplN
	l.infofImpl = printfImplN
	l.infolnImpl = printlnImplN

	l.noticeImpl = printImplN
	l.noticefImpl = printfImplN
	l.noticelnImpl = printlnImplN

	l.warnImpl = printImplN
	l.warnfImpl = printfImplN
	l.warnlnImpl = printlnImplN

	l.errorImpl = printImplN
	l.errorfImpl = printfImplN
	l.errorlnImpl = printlnImplN

	// log that level is PANIC and FATAL must be output
	switch level {
	case DEBUG:
		l.debugImpl = printImplY
		l.debugfImpl = printfImplY
		l.debuglnImpl = printlnImplY
		fallthrough
	case INFO:
		l.infoImpl = printImplY
		l.infofImpl = printfImplY
		l.infolnImpl = printlnImplY
		fallthrough
	case NOTICE:
		l.noticeImpl = printImplY
		l.noticefImpl = printfImplY
		l.noticelnImpl = printlnImplY
		fallthrough
	case WARN:
		l.warnImpl = printImplY
		l.warnfImpl = printfImplY
		l.warnlnImpl = printlnImplY
		fallthrough
	case ERROR:
		l.errorImpl = printImplY
		l.errorfImpl = printfImplY
		l.errorlnImpl = printlnImplY
	case _PANIC, FATAL:

	}
}

// Debug
func (l *Slogger) Debug(v ...interface{}) {
	var level Level = DEBUG
	l.debugImpl(l, level, v...)
}

func (l *Slogger) Debugf(format string, v ...interface{}) {
	var level Level = DEBUG
	l.debugfImpl(l, level, format, v...)
}

func (l *Slogger) Debugln(v ...interface{}) {
	var level Level = DEBUG
	l.debuglnImpl(l, level, v...)
}

// Info
func (l *Slogger) Info(v ...interface{}) {
	var level Level = INFO
	l.infoImpl(l, level, v...)
}

func (l *Slogger) Infof(format string, v ...interface{}) {
	var level Level = INFO
	l.infofImpl(l, level, format, v...)
}

func (l *Slogger) Infoln(v ...interface{}) {
	var level Level = INFO
	l.infolnImpl(l, level, v...)
}

// Notice
func (l *Slogger) Notice(v ...interface{}) {
	var level Level = NOTICE
	l.noticeImpl(l, level, v...)
}

func (l *Slogger) Noticef(format string, v ...interface{}) {
	var level Level = NOTICE
	l.noticefImpl(l, level, format, v...)
}

func (l *Slogger) Noticeln(v ...interface{}) {
	var level Level = NOTICE
	l.noticelnImpl(l, level, v...)
}

// Warn
func (l *Slogger) Warn(v ...interface{}) {
	var level Level = WARN
	l.warnImpl(l, level, v...)
}

func (l *Slogger) Warnf(format string, v ...interface{}) {
	var level Level = WARN
	l.warnfImpl(l, level, format, v...)
}

func (l *Slogger) Warnln(v ...interface{}) {
	var level Level = WARN
	l.warnlnImpl(l, level, v...)
}

// Error
func (l *Slogger) Error(v ...interface{}) {
	var level Level = ERROR
	l.errorImpl(l, level, v...)
}

func (l *Slogger) Errorf(format string, v ...interface{}) {
	var level Level = ERROR
	l.errorfImpl(l, level, format, v...)
}

func (l *Slogger) Errorln(v ...interface{}) {
	var level Level = ERROR
	l.errorlnImpl(l, level, v...)
}

// Fatal
func (l *Slogger) Fatal(v ...interface{}) {
	l.logger.Output(l.callPath, levelName[FATAL]+fmt.Sprint(v...))
	os.Exit(1)
}

func (l *Slogger) Fatalf(format string, v ...interface{}) {
	l.logger.Output(l.callPath, fmt.Sprintf(levelName[FATAL]+format, v...))
	os.Exit(1)
}

func (l *Slogger) Fatalln(v ...interface{}) {
	l.logger.Output(l.callPath, levelName[FATAL]+fmt.Sprintln(v...))
	os.Exit(1)
}

// Panic
func (l *Slogger) Panic(v ...interface{}) {
	s := levelName[_PANIC] + fmt.Sprintln(v...)
	l.logger.Output(l.callPath, s)
	panic(s)
}

func (l *Slogger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(levelName[_PANIC]+format, v...)
	l.logger.Output(l.callPath, s)
	panic(s)
}

func (l *Slogger) Panicln(v ...interface{}) {
	s := levelName[_PANIC] + fmt.Sprintln(v...)
	l.logger.Output(l.callPath, s)
	panic(s)
}

// default logger
var dlog *Slogger

func init() {
	dlog = New(os.Stderr, "DEBUG", "", LstdFlags|Lshortfile)
	dlog.callPath = 4
}

func SetLevel(level Level) {
	dlog.SetLevel(level)
}

func Debug(v ...interface{}) {
	dlog.Debug(v...)
}

func Debugf(fmt string, v ...interface{}) {
	dlog.Debugf(fmt, v...)
}

func Debugln(v ...interface{}) {
	dlog.Debugln(v...)
}

func Info(v ...interface{}) {
	dlog.Info(v...)
}

func Infof(fmt string, v ...interface{}) {
	dlog.Infof(fmt, v...)
}

func Infoln(v ...interface{}) {
	dlog.Infoln(v...)
}
func Notice(v ...interface{}) {
	dlog.Notice(v...)
}

func Noticef(fmt string, v ...interface{}) {
	dlog.Noticef(fmt, v...)
}

func Noticeln(v ...interface{}) {
	dlog.Noticeln(v...)
}

func Warn(v ...interface{}) {
	dlog.Warn(v...)
}

func Warnf(fmt string, v ...interface{}) {
	dlog.Warnf(fmt, v...)
}

func Warnln(v ...interface{}) {
	dlog.Warnln(v...)
}

func Error(v ...interface{}) {
	dlog.Error(v...)
}

func Errorf(fmt string, v ...interface{}) {
	dlog.Errorf(fmt, v...)
}

func Errorln(v ...interface{}) {
	dlog.Errorln(v...)
}

func Fatal(v ...interface{}) {
	dlog.Fatal(v...)
}

func Fatalf(fmt string, v ...interface{}) {
	dlog.Fatalf(fmt, v...)
}

func Fatalln(v ...interface{}) {
	dlog.Fatalln(v...)
}

func Panic(v ...interface{}) {
	dlog.Panic(v...)
}

func Panicf(fmt string, v ...interface{}) {
	dlog.Panicf(fmt, v...)
}

func Panicln(v ...interface{}) {
	dlog.Panicln(v...)
}
