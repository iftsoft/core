package main

import (
	"fmt"
	"github.com/iftsoft/core"
)

type AppConfig struct {
	Logging core.LogConfig `yaml:"logging"`
}

func main() {
	fmt.Println("-------BEGIN------------")

	config := &AppConfig{}
	err := core.ReadYamlFile("config.yml", config)
	if err != nil {
		fmt.Println(err)
	} else {
		core.StartFileLogger(&config.Logging)
	}
	log := core.GetLogAgent(core.LogLevelTrace, "APP")
	log.Info("Start application")
	log.Info("Config %+v", config)

	log.Info("Stop application")
	core.StopFileLogger()
	fmt.Println("-------END------------")
}
