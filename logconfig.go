package core

import (
	"errors"
	"fmt"
)

type LogConfig struct {
	LogPath   string `yaml:"logPath"`
	LogFile   string `yaml:"logFile"`
	LogLevel  int    `yaml:"logLevel"`
	ConsLevel int    `yaml:"consLevel"`
	MaxFiles  int    `yaml:"maxFiles"` // limit the number of log files under `logPath`
	MaxSize   int64  `yaml:"maxSize"`  // limit size of a log file (KByte)
}

func (cfg *LogConfig) PrintData() {
	fmt.Println("LogPath  ", cfg.LogPath)
	fmt.Println("LogFile  ", cfg.LogFile)
	fmt.Println("LogLevel ", GetLogLevelText(cfg.LogLevel))
	fmt.Println("ConsLevel", GetLogLevelText(cfg.ConsLevel))
	fmt.Println("MaxFiles ", cfg.MaxFiles)
	fmt.Println("MaxSize  ", cfg.MaxSize)
}
func (cfg *LogConfig) String() string {
	str := fmt.Sprintf("Logging config: "+
		"LogPath = %s, LogFile = %s, LogLevel = %s, ConsLevel = %s, MaxFiles = %d, MaxSize = %d.",
		cfg.LogPath, cfg.LogFile, GetLogLevelText(cfg.LogLevel), GetLogLevelText(cfg.ConsLevel), cfg.MaxFiles, cfg.MaxSize)
	return str
}

func GetDefaultConfig(name string) *LogConfig {
	cfg := LogConfig{
		LogPath:   "",
		LogFile:   name,
		LogLevel:  LogLevelInfo,
		ConsLevel: LogLevelError,
		MaxFiles:  8,
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
	if cfg.MaxSize < 0 || cfg.MaxSize >= 128*1024 {
		cfg.MaxSize = 1024
	}
	return err
}
