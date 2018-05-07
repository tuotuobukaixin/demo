package models

import (
	"math"
)

type SysConfig struct {
	EventSerialNumber int64 `json:"eventSerialNumber"`
	AlarmSerialNumber int64 `json:"alarmSerialNumber"`
}

func LoadSystemIDFromPath(systemId int64) *SysConfig {
	sysConfig := SysConfig{}
	sysConfig.EventSerialNumber = systemId * int64(math.Pow(10, 13))
	sysConfig.AlarmSerialNumber = systemId * int64(math.Pow(10, 13))

	return &sysConfig
}
