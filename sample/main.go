package main

import (
	"fmt"
	"github.com/iftsoft/core/log"
	"time"
)

func main() {
	fmt.Println("-------BEGIN------------")

	logCfg := log.LogConfig{
		LogPath:   "logs",
		LogFile:   "sample",
		LogLevel:  log.LogLevelTrace,
		ConsLevel: log.LogLevelError,
		MaxFiles:  4,
		DelFiles:  1,
		MaxSize:   1024,
	}
	log.StartFileLogger(&logCfg)
	out := log.GetLogAgent(log.LogLevelTrace, "APP")
	out.Info("Start application")

	out.Info("Stop application")
	time.Sleep(time.Second)
	log.StopFileLogger()
	fmt.Println("-------END------------")
}
