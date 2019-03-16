package log

import (
	"bytes"
	"fmt"
	"runtime"
)

type LogConfig struct {
	LogPath   string
	LogFile   string
	LogLevel  int
	ConsLevel int
	MaxFiles  int   // limit the number of log files under `logPath`
	DelFiles  int   // number of files deleted when reaching the limit of the number of log files
	MaxSize   int64 // limit size of a log file (KByte)
}

func (cfg *LogConfig) PrintData() {
	fmt.Println("LogPath  ", cfg.LogPath)
	fmt.Println("LogFile  ", cfg.LogFile)
	fmt.Println("LogLevel ", GetLogLevelText(cfg.LogLevel))
	fmt.Println("ConsLevel", GetLogLevelText(cfg.ConsLevel))
	fmt.Println("MaxFiles ", cfg.MaxFiles)
	fmt.Println("DelFiles ", cfg.DelFiles)
	fmt.Println("MaxSize  ", cfg.MaxSize)
}
func (cfg *LogConfig) String() string {
	str := fmt.Sprintf("Logging config: "+
		"LogPath = %s, LogFile = %s, LogLevel = %s, ConsLevel = %s, MaxFiles = %d, DelFiles = %d, MaxSize = %d.",
		cfg.LogPath, cfg.LogFile, GetLogLevelText(cfg.LogLevel), GetLogLevelText(cfg.ConsLevel), cfg.MaxFiles, cfg.DelFiles, cfg.MaxSize)
	return str
}

// Get go routine ID as a string
func GetGID() string {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	return string(b)
}

// Extract error text from Error object
func GetErrorText(err error) string {
	if err != nil {
		return err.Error()
	}
	return "Success"
}

// Get Call Stack Trace as a string
func TraceCallStack(text string, i int) string {
	//	i := 2
	for {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		text += fmt.Sprintf("\n%s:%d %s", file, line, runtime.FuncForPC(pc).Name())
		i++
	}
	return text
}
