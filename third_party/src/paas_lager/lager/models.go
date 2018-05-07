package lager

import "encoding/json"

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

func FormatLogLevel(x LogLevel) string {
	var level string
	switch x {
	case DEBUG:
		level = "debug"
	case INFO:
		level = "info"
	case WARN:
		level = "warn"
	case ERROR:
		level = "error"
	case FATAL:
		level = "fatal"
	}
	return level
}

func (x LogLevel) MarshalJSON() ([]byte, error) {
	// var level string
	var level = FormatLogLevel(x)
	return json.Marshal(level)
}

/*
func (x LogLevel) MarshalJSON() ([]byte, error) {
	var level string
	switch x {
	case DEBUG:
		level = "debug"
	case INFO:
		level = "info"
	case WARN:
		level = "warn"
	case ERROR:
		level = "error"
	case FATAL:
		level = "fatal"
	}
	return json.Marshal(level)
}
*/

type Data map[string]interface{}

type LogFormat struct {
	Timestamp string   `json:"timestamp"`
	Source    string   `json:"source"`
	Message   string   `json:"message"`
	LogLevel  LogLevel `json:"log_level"`
	Data      Data     `json:"data"`
	ProcessID int      `json:"process_id"`
	File      string   `json:"file"`
	LineNo    int      `json:"lineno"`
	Method    string   `json:"method"`
}

func (log LogFormat) ToJSON() ([]byte, error) {
	return json.Marshal(log)
}
