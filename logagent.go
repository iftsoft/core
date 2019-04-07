package core

import (
	"fmt"
	"runtime"
	"time"
)

// log level
const (
	LogLevelEmpty int = iota
	LogLevelPanic
	LogLevelError
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
	LogLevelDump
	LogLevelTrace
	LogLevelMax
)

var logLevelNames = [LogLevelMax]string{
	"EMPTY", "PANIC", "ERROR", "WARN ", "INFO ", "DEBUG", "DUMP ", "TRACE",
}

func GetLogLevelText(level int) string {
	if level >= 0 && level < LogLevelMax {
		return logLevelNames[level]
	}
	return "UNDEF"
}

func GetLogAgent(level int, title string) *LogAgent {
	if level >= LogLevelMax {
		level = LogLevelEmpty
	}
	return &LogAgent{level, title}
}

type LogAgent struct {
	logLevel int
	modTitle string
}

func (log *LogAgent) Init(level int, title string) {
	log.logLevel = level
	log.modTitle = title
}

func (log *LogAgent) SetLevel(level int) {
	log.logLevel = level
}

func (log LogAgent) IsTrace() bool {
	return log.logLevel >= LogLevelTrace
}
func (log LogAgent) IsDump() bool {
	return log.logLevel >= LogLevelDump
}
func (log LogAgent) IsDebug() bool {
	return log.logLevel >= LogLevelDebug
}
func (log LogAgent) IsInfo() bool {
	return log.logLevel >= LogLevelInfo
}
func (log LogAgent) IsWarn() bool {
	return log.logLevel >= LogLevelWarn
}
func (log LogAgent) IsError() bool {
	return log.logLevel >= LogLevelError
}
func (log LogAgent) IsPanic() bool {
	return log.logLevel >= LogLevelPanic
}
func (log LogAgent) IsEmpty() bool {
	return log.logLevel == LogLevelEmpty
}

func (log *LogAgent) Trace(format string, args ...interface{}) {
	if log != nil && log.IsTrace() {
		log.formatLine(LogLevelTrace, fmt.Sprintf(format, args...))
	}
}
func (log *LogAgent) Dump(format string, args ...interface{}) {
	if log != nil && log.IsDump() {
		log.formatLine(LogLevelDump, fmt.Sprintf(format, args...))
	}
}
func (log *LogAgent) Debug(format string, args ...interface{}) {
	if log != nil && log.IsDebug() {
		log.formatLine(LogLevelDebug, fmt.Sprintf(format, args...))
	}
}
func (log *LogAgent) Info(format string, args ...interface{}) {
	if log != nil && log.IsInfo() {
		log.formatLine(LogLevelInfo, fmt.Sprintf(format, args...))
	}
}
func (log *LogAgent) Warn(format string, args ...interface{}) {
	if log != nil && log.IsWarn() {
		log.formatLine(LogLevelWarn, fmt.Sprintf(format, args...))
	}
}
func (log *LogAgent) Error(format string, args ...interface{}) {
	if log != nil && log.IsError() {
		log.formatLine(LogLevelError, fmt.Sprintf(format, args...))
	}
}
func (log *LogAgent) Panic(format string, args ...interface{}) {
	if log != nil && log.IsPanic() {
		text := fmt.Sprintf(format, args...)
		log.formatLine(LogLevelPanic, TraceCallStack(text, 2))
	}
}

func (log *LogAgent) formatLine(level int, text string) {
	t := time.Now()
	moment := t.Format("2006-01-02 15:04:05.999999")
	for size := len(moment); size < 26; size++ {
		moment += "0"
	}
	gid := GetGID()
	var mesg string
	switch level {
	case LogLevelDebug, LogLevelError, LogLevelPanic:
		pc, file, line, ok := runtime.Caller(2)
		if ok {
			mesg = fmt.Sprintf("%s [%s %s %s] %s:%d %s() %s\n",
				moment, logLevelNames[level], gid, log.modTitle, file, line, runtime.FuncForPC(pc).Name(), text)
		} else {
			mesg = fmt.Sprintf("%s [%s %s %s] %s\n", moment, logLevelNames[level], gid, log.modTitle, text)
		}
	case LogLevelTrace, LogLevelDump, LogLevelInfo, LogLevelWarn:
		mesg = fmt.Sprintf("%s [%s %s %s] %s\n", moment, logLevelNames[level], gid, log.modTitle, text)
	}
	LogToFile(level, mesg)
}

func (log *LogAgent) PanicRecover() {
	if r := recover(); r != nil {
		if log != nil {
			log.Warn("Panic Recovered: %+v", r)
		}
	}
}
