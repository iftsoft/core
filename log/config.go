package log

import (
	"bytes"
	"errors"
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

func GetDefaultConfig(name string) *LogConfig {
	cfg := LogConfig{
		LogPath:   "",
		LogFile:   name,
		LogLevel:  LogLevelInfo,
		ConsLevel: LogLevelError,
		MaxFiles:  8,
		DelFiles:  1,
		MaxSize:   1024,
	}
	return &cfg
}

func checkLogConfig(cfg *LogConfig) (err error) {
	if cfg == nil {
		return errors.New("Logging: config is not set")
	}
	if cfg.LogFile == "" {
		return errors.New("Logging: file name is not set")
	}
	if cfg.LogPath == "" {
		cfg.LogPath = "."
	}
	if cfg.LogLevel < LogLevelEmpty || cfg.LogLevel >= LogLevelMax {
		cfg.LogLevel = LogLevelInfo
	}
	if cfg.ConsLevel < LogLevelEmpty || cfg.ConsLevel >= LogLevelMax {
		cfg.ConsLevel = LogLevelError
	}
	if cfg.MaxFiles < 0 || cfg.MaxFiles >= 1024 {
		cfg.MaxFiles = 8
	}
	if cfg.DelFiles < 0 || cfg.DelFiles >= cfg.MaxFiles {
		cfg.DelFiles = 1
	}
	if cfg.MaxSize < 0 || cfg.MaxSize >= 128*1024 {
		cfg.MaxSize = 1024
	}
	return err
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
