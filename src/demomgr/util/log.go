package util

import (
	"paas_lager"
	"paas_lager/lager"
)

//log var
var LOGGER lager.Logger

func init() {

	paas_lager.Init(paas_lager.Config{
		LoggerLevel:   "DEBUG",
		LoggerFile:    "app.log",
		EnableRsyslog: false,
	})

	LOGGER = paas_lager.NewLogger("demomgr")
}
