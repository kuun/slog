package slog

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/kuun/slog/buffer"
	"github.com/kuun/slog/writer"
)

var levelChars = [6]byte{'D', 'I', 'N', 'W', 'E', 'F'}

type loggerImpl struct {
	level    Level
	writers  []writer.LogWriter
	fullPath string
	abbrPath string

	// log.Logger.Output callPath
	callPath int

	debugImpl  func(l *loggerImpl, level Level, v ...interface{})
	debugfImpl func(l *loggerImpl, level Level, format string, v ...interface{})

	infoImpl  func(l *loggerImpl, level Level, v ...interface{})
	infofImpl func(l *loggerImpl, level Level, format string, v ...interface{})

	noticeImpl  func(l *loggerImpl, level Level, v ...interface{})
	noticefImpl func(l *loggerImpl, level Level, format string, v ...interface{})

	warnImpl  func(l *loggerImpl, level Level, v ...interface{})
	warnfImpl func(l *loggerImpl, level Level, format string, v ...interface{})

	errorImpl  func(l *loggerImpl, level Level, v ...interface{})
	errorfImpl func(l *loggerImpl, level Level, format string, v ...interface{})

	fatalImpl  func(l *loggerImpl, level Level, v ...interface{})
	fatalfImpl func(l *loggerImpl, level Level, format string, v ...interface{})
}

// func suffix is "Y" is valid implements
// func suffix is "N" is empty implements

func printImplY(l *loggerImpl, level Level, v ...interface{}) {
	buff := l.header(level)
	buff.WriteString(fmt.Sprint(v...))
	buff.WriteByte('\n')
	for _, wr := range l.writers {
		wr.Write(buff)
	}
}

func printN(l *loggerImpl, level Level, v ...interface{}) {
}

func printfImplY(l *loggerImpl, level Level, format string, v ...interface{}) {
	buff := l.header(level)
	buff.WriteString(fmt.Sprintf(format, v...))
	buff.WriteByte('\n')
	for _, wr := range l.writers {
		wr.Write(buff)
	}
}

func printfImplN(l *loggerImpl, level Level, format string, v ...interface{}) {
}

func (l *loggerImpl) GetLevel() string {
	return l.level.String()
}

func (l *loggerImpl) SetLevel(level string) error {
	if lv, ok := parseLevel(level); ok {
		l.level = lv
	} else {
		return errors.New("unkown log level")
	}
	// init all print func
	l.debugImpl = printN
	l.debugfImpl = printfImplN

	l.infoImpl = printN
	l.infofImpl = printfImplN

	l.noticeImpl = printN
	l.noticefImpl = printfImplN

	l.warnImpl = printN
	l.warnfImpl = printfImplN

	l.errorImpl = printN
	l.errorfImpl = printfImplN

	// log that level is PANIC and FATAL must be output
	switch l.level {
	case Debug:
		l.debugImpl = printImplY
		l.debugfImpl = printfImplY
		fallthrough
	case Info:
		l.infoImpl = printImplY
		l.infofImpl = printfImplY
		fallthrough
	case Notice:
		l.noticeImpl = printImplY
		l.noticefImpl = printfImplY
		fallthrough
	case Warn:
		l.warnImpl = printImplY
		l.warnfImpl = printfImplY
		fallthrough
	case Error:
		l.errorImpl = printImplY
		l.errorfImpl = printfImplY
	case Fatal:
		l.fatalImpl = printImplY
		l.fatalfImpl = printfImplY
	}
	return nil
}

func (l *loggerImpl) Above(lv Level) bool {
	return lv >= l.level
}

// Debug
func (l *loggerImpl) Debug(v ...interface{}) {
	var level Level = Debug
	l.debugImpl(l, level, v...)
}

func (l *loggerImpl) Debugf(format string, v ...interface{}) {
	var level Level = Debug
	l.debugfImpl(l, level, format, v...)
}

// Info
func (l *loggerImpl) Info(v ...interface{}) {
	var level Level = Info
	l.infoImpl(l, level, v...)
}

func (l *loggerImpl) Infof(format string, v ...interface{}) {
	var level Level = Info
	l.infofImpl(l, level, format, v...)
}

// Notice
func (l *loggerImpl) Notice(v ...interface{}) {
	var level Level = Notice
	l.noticeImpl(l, level, v...)
}

func (l *loggerImpl) Noticef(format string, v ...interface{}) {
	var level Level = Notice
	l.noticefImpl(l, level, format, v...)
}

// Warn
func (l *loggerImpl) Warn(v ...interface{}) {
	var level Level = Warn
	l.warnImpl(l, level, v...)
}

func (l *loggerImpl) Warnf(format string, v ...interface{}) {
	var level Level = Warn
	l.warnfImpl(l, level, format, v...)
}

// Error
func (l *loggerImpl) Error(v ...interface{}) {
	var level Level = Error
	l.errorImpl(l, level, v...)
}

func (l *loggerImpl) Errorf(format string, v ...interface{}) {
	var level Level = Error
	l.errorfImpl(l, level, format, v...)
}

// Fatal
func (l *loggerImpl) Fatal(v ...interface{}) {
	var level Level = Fatal
	l.fatalImpl(l, level, v...)
	os.Exit(1)
}

func (l *loggerImpl) Fatalf(format string, v ...interface{}) {
	var level Level = Fatal
	l.fatalfImpl(l, level, format, v...)
	os.Exit(1)
}

//
// these codes are from github/golang/glog
//

// header formats a log header and returns a buffer containing the formatted
// header and the user's file and line number.
// Log lines have this form:
// 	L mmdd hh:mm:ss.uuuuuu threadid file:line] msg...
// where the fields are defined as follows:
// 	L                A single character, representing the log level (eg 'I' for INFO)
// 	mm               The month (zero padded; ie May is '05')
// 	dd               The day (zero padded)
// 	hh:mm:ss.uuuuuu  Time in hours, minutes and fractional seconds
// 	threadid         The space-padded thread ID as returned by GetTID()
// 	file             The file name
// 	line             The line number
// 	msg              The user-supplied message
func (l *loggerImpl) header(lv Level) *buffer.Buffer {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return l.formatHeader(lv, file, line)
}

// formatHeader formats a log header using the provided file name and line number.
func (l *loggerImpl) formatHeader(lv Level, file string, line int) *buffer.Buffer {
	now := time.Now()
	if line < 0 {
		line = 0 // not a real line number, but acceptable to someDigits
	}
	buf := buffer.GetBuffer()

	// Avoid Fprintf, for speed. The format is so simple that we can do it quickly by hand.
	// It's worth about 3X. Fprintf is hard.
	_, month, day := now.Date()
	hour, minute, second := now.Clock()
	// Lmmdd hh:mm:ss.uuuuuu threadid file:line]
	buf.Tmp[0] = levelChars[lv]
	buf.Tmp[1] = ' '
	buf.TwoDigits(2, int(month))
	buf.TwoDigits(4, day)
	buf.Tmp[6] = ' '
	buf.TwoDigits(7, hour)
	buf.Tmp[9] = ':'
	buf.TwoDigits(10, minute)
	buf.Tmp[12] = ':'
	buf.TwoDigits(13, second)
	buf.Tmp[15] = '.'
	buf.NDigits(6, 16, now.Nanosecond()/1000, '0')
	buf.Tmp[22] = ' '
	buf.Write(buf.Tmp[:23])
	buf.WriteString(file)
	buf.Tmp[0] = ':'
	n := buf.SomeDigits(1, line)
	buf.Tmp[n+1] = ']'
	buf.Tmp[n+2] = ' '
	buf.Write(buf.Tmp[:n+3])
	return buf
}
