package main

import (
	"fmt"

	"paas_lager"
	"paas_lager/lager"
)

func main() {
	paas_lager.Init(paas_lager.Config{
		LoggerLevel:   "DEBUG",
		LoggerFile:    "",
		EnableRsyslog: false,
		LogFormatText: false,
	})

	logger := paas_lager.NewLogger("example")

	logger.Infof("Hi %s, system is starting up ...", "paas-bot")

	logger.Debug("check-info", lager.Data{
		"info": "something",
	})

	err := fmt.Errorf("Oops, error occurred!")
	logger.Warn("failed-to-do-somthing", err, lager.Data{
		"info": "something",
	})

	err = fmt.Errorf("This is an error")
	logger.Error("failed-to-do-somthing", err)

	logger.Info("shutting-down")
}
