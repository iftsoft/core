package main

import (
	"fmt"
	"github.com/iftsoft/core"
	"time"
)

func main() {
	fmt.Println("-------BEGIN------------")

	logCfg := core.LogConfig{
		LogPath:   "logs",
		LogFile:   "sample",
		LogLevel:  core.LogLevelTrace,
		ConsLevel: core.LogLevelError,
		MaxFiles:  4,
		DelFiles:  1,
		MaxSize:   1024,
	}
	core.StartFileLogger(&logCfg)
	log := core.GetLogAgent(core.LogLevelTrace, "APP")
	log.Info("Start application")

	log.Info("Stop application")
	time.Sleep(time.Second)
	core.StopFileLogger()
	fmt.Println("-------END------------")
}
